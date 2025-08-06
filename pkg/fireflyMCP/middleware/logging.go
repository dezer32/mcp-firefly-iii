package middleware

import (
	"fmt"
	"log"
	"time"
)

// LoggingMiddleware provides request/response logging
type LoggingMiddleware struct {
	logger *log.Logger
	level  LogLevel
}

// LogLevel defines the logging level
type LogLevel int

const (
	// LogLevelDebug logs everything
	LogLevelDebug LogLevel = iota
	// LogLevelInfo logs info and above
	LogLevelInfo
	// LogLevelWarn logs warnings and errors
	LogLevelWarn
	// LogLevelError logs only errors
	LogLevelError
)

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *log.Logger, level LogLevel) *LoggingMiddleware {
	if logger == nil {
		logger = log.Default()
	}
	return &LoggingMiddleware{
		logger: logger,
		level:  level,
	}
}

// Process implements the Middleware interface
func (l *LoggingMiddleware) Process(next Handler) Handler {
	return func(req *ToolRequest) (*ToolResponse, error) {
		// Log request
		if l.level <= LogLevelInfo {
			requestID, _ := GetRequestID(req.Context)
			l.logger.Printf("[INFO] Request: tool=%s, request_id=%s, time=%s",
				req.ToolName, requestID, req.StartTime.Format(time.RFC3339))
		}

		if l.level <= LogLevelDebug {
			l.logger.Printf("[DEBUG] Request arguments: %+v", req.Arguments)
			l.logger.Printf("[DEBUG] Request metadata: %+v", req.Metadata)
		}

		// Call next handler
		resp, err := next(req)

		// Calculate duration
		duration := time.Since(req.StartTime)

		// Log response
		if err != nil {
			if l.level <= LogLevelError {
				requestID, _ := GetRequestID(req.Context)
				l.logger.Printf("[ERROR] Response: tool=%s, request_id=%s, duration=%v, error=%v",
					req.ToolName, requestID, duration, err)
			}
		} else if resp.IsError {
			if l.level <= LogLevelWarn {
				requestID, _ := GetRequestID(req.Context)
				l.logger.Printf("[WARN] Response: tool=%s, request_id=%s, duration=%v, is_error=true",
					req.ToolName, requestID, duration)
			}
		} else {
			if l.level <= LogLevelInfo {
				requestID, _ := GetRequestID(req.Context)
				l.logger.Printf("[INFO] Response: tool=%s, request_id=%s, duration=%v, success=true",
					req.ToolName, requestID, duration)
			}
		}

		if l.level <= LogLevelDebug && resp != nil {
			l.logger.Printf("[DEBUG] Response metadata: %+v", resp.Metadata)
		}

		return resp, err
	}
}

// Name returns the middleware name
func (l *LoggingMiddleware) Name() string {
	return "logging"
}

// RequestLoggingMiddleware logs detailed request information
type RequestLoggingMiddleware struct {
	logger    *log.Logger
	logBody   bool
	maxBodySize int
}

// NewRequestLoggingMiddleware creates a new request logging middleware
func NewRequestLoggingMiddleware(logger *log.Logger, logBody bool) *RequestLoggingMiddleware {
	if logger == nil {
		logger = log.Default()
	}
	return &RequestLoggingMiddleware{
		logger:      logger,
		logBody:     logBody,
		maxBodySize: 1024, // Limit body logging to 1KB
	}
}

// Process implements the Middleware interface
func (r *RequestLoggingMiddleware) Process(next Handler) Handler {
	return func(req *ToolRequest) (*ToolResponse, error) {
		// Generate request ID if not present
		requestID, ok := GetRequestID(req.Context)
		if !ok {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
			req.Context = WithRequestID(req.Context, requestID)
		}

		// Log request details
		r.logger.Printf("[REQUEST] ID=%s Tool=%s Time=%s",
			requestID, req.ToolName, req.StartTime.Format(time.RFC3339))

		if r.logBody && req.Arguments != nil {
			argStr := fmt.Sprintf("%+v", req.Arguments)
			if len(argStr) > r.maxBodySize {
				argStr = argStr[:r.maxBodySize] + "..."
			}
			r.logger.Printf("[REQUEST_BODY] ID=%s Args=%s", requestID, argStr)
		}

		// Call next handler
		resp, err := next(req)

		// Log response
		duration := time.Since(req.StartTime)
		if err != nil {
			r.logger.Printf("[RESPONSE] ID=%s Duration=%v Error=%v",
				requestID, duration, err)
		} else if resp != nil {
			r.logger.Printf("[RESPONSE] ID=%s Duration=%v IsError=%v",
				requestID, duration, resp.IsError)
		}

		return resp, err
	}
}

// Name returns the middleware name
func (r *RequestLoggingMiddleware) Name() string {
	return "request_logging"
}