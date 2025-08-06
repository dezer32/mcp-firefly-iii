package middleware

import (
	"fmt"
	"log"
	"runtime/debug"
)

// RecoveryMiddleware recovers from panics and converts them to errors
type RecoveryMiddleware struct {
	logger         *log.Logger
	printStackTrace bool
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware(logger *log.Logger, printStackTrace bool) *RecoveryMiddleware {
	if logger == nil {
		logger = log.Default()
	}
	return &RecoveryMiddleware{
		logger:         logger,
		printStackTrace: printStackTrace,
	}
}

// Process implements the Middleware interface
func (r *RecoveryMiddleware) Process(next Handler) Handler {
	return func(req *ToolRequest) (resp *ToolResponse, err error) {
		// Defer panic recovery
		defer func() {
			if rec := recover(); rec != nil {
				// Log the panic
				requestID, _ := GetRequestID(req.Context)
				r.logger.Printf("[PANIC] Recovered from panic in tool=%s, request_id=%s: %v",
					req.ToolName, requestID, rec)

				if r.printStackTrace {
					r.logger.Printf("[PANIC] Stack trace:\n%s", debug.Stack())
				}

				// Convert panic to error response
				err = fmt.Errorf("internal error: %v", rec)
				resp = &ToolResponse{
					IsError:  true,
					Metadata: map[string]interface{}{
						"panic":      true,
						"panic_value": fmt.Sprintf("%v", rec),
					},
				}
			}
		}()

		// Call next handler
		return next(req)
	}
}

// Name returns the middleware name
func (r *RecoveryMiddleware) Name() string {
	return "recovery"
}

// SafeExecutionMiddleware ensures safe execution with timeout and panic recovery
type SafeExecutionMiddleware struct {
	logger        *log.Logger
	panicHandler  func(interface{})
}

// NewSafeExecutionMiddleware creates a new safe execution middleware
func NewSafeExecutionMiddleware(logger *log.Logger) *SafeExecutionMiddleware {
	if logger == nil {
		logger = log.Default()
	}
	return &SafeExecutionMiddleware{
		logger: logger,
	}
}

// SetPanicHandler sets a custom panic handler
func (s *SafeExecutionMiddleware) SetPanicHandler(handler func(interface{})) {
	s.panicHandler = handler
}

// Process implements the Middleware interface
func (s *SafeExecutionMiddleware) Process(next Handler) Handler {
	return func(req *ToolRequest) (resp *ToolResponse, err error) {
		// Channel to receive result
		done := make(chan struct{})

		// Execute in goroutine for safety
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					requestID, _ := GetRequestID(req.Context)
					s.logger.Printf("[SAFE_EXEC] Panic recovered for tool=%s, request_id=%s: %v",
						req.ToolName, requestID, rec)

					if s.panicHandler != nil {
						s.panicHandler(rec)
					}

					err = fmt.Errorf("execution failed: %v", rec)
					resp = &ToolResponse{
						IsError: true,
						Metadata: map[string]interface{}{
							"error_type": "panic",
							"error":      fmt.Sprintf("%v", rec),
						},
					}
				}
				close(done)
			}()

			// Execute the handler
			resp, err = next(req)
		}()

		// Wait for completion or context cancellation
		select {
		case <-done:
			return resp, err
		case <-req.Context.Done():
			return &ToolResponse{
				IsError: true,
				Metadata: map[string]interface{}{
					"error_type": "timeout",
					"error":      "request context cancelled",
				},
			}, req.Context.Err()
		}
	}
}

// Name returns the middleware name
func (s *SafeExecutionMiddleware) Name() string {
	return "safe_execution"
}