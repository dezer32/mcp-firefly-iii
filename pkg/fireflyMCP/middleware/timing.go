package middleware

import (
	"log"
	"time"
)

// TimingMiddleware measures request processing time
type TimingMiddleware struct {
	logger    *log.Logger
	threshold time.Duration // Log slow requests above this threshold
}

// NewTimingMiddleware creates a new timing middleware
func NewTimingMiddleware(logger *log.Logger, threshold time.Duration) *TimingMiddleware {
	if logger == nil {
		logger = log.Default()
	}
	if threshold == 0 {
		threshold = 1 * time.Second // Default to 1 second
	}
	return &TimingMiddleware{
		logger:    logger,
		threshold: threshold,
	}
}

// Process implements the Middleware interface
func (t *TimingMiddleware) Process(next Handler) Handler {
	return func(req *ToolRequest) (*ToolResponse, error) {
		start := time.Now()

		// Call next handler
		resp, err := next(req)

		// Calculate duration
		duration := time.Since(start)

		// Add timing to response metadata
		if resp == nil {
			resp = &ToolResponse{
				Metadata: make(map[string]interface{}),
			}
		}
		if resp.Metadata == nil {
			resp.Metadata = make(map[string]interface{})
		}
		resp.Duration = duration
		resp.Metadata["duration_ms"] = duration.Milliseconds()
		resp.Metadata["duration_ns"] = duration.Nanoseconds()

		// Log slow requests
		if duration > t.threshold {
			requestID, _ := GetRequestID(req.Context)
			t.logger.Printf("[SLOW_REQUEST] tool=%s, request_id=%s, duration=%v (threshold=%v)",
				req.ToolName, requestID, duration, t.threshold)
		}

		return resp, err
	}
}

// Name returns the middleware name
func (t *TimingMiddleware) Name() string {
	return "timing"
}

// MetricsMiddleware collects metrics about request processing
type MetricsMiddleware struct {
	metrics *RequestMetrics
}

// RequestMetrics holds request processing metrics
type RequestMetrics struct {
	TotalRequests   int64
	SuccessRequests int64
	ErrorRequests   int64
	TotalDuration   time.Duration
	MinDuration     time.Duration
	MaxDuration     time.Duration
	ToolMetrics     map[string]*ToolMetrics
}

// ToolMetrics holds metrics for a specific tool
type ToolMetrics struct {
	TotalCalls    int64
	SuccessCalls  int64
	ErrorCalls    int64
	TotalDuration time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration
	LastCalled    time.Time
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: &RequestMetrics{
			ToolMetrics: make(map[string]*ToolMetrics),
			MinDuration: time.Duration(1<<63 - 1), // Max duration as initial min
		},
	}
}

// Process implements the Middleware interface
func (m *MetricsMiddleware) Process(next Handler) Handler {
	return func(req *ToolRequest) (*ToolResponse, error) {
		start := time.Now()

		// Call next handler
		resp, err := next(req)

		// Calculate duration
		duration := time.Since(start)

		// Update global metrics
		m.updateGlobalMetrics(duration, err, resp)

		// Update tool-specific metrics
		m.updateToolMetrics(req.ToolName, duration, err, resp)

		return resp, err
	}
}

// updateGlobalMetrics updates global request metrics
func (m *MetricsMiddleware) updateGlobalMetrics(duration time.Duration, err error, resp *ToolResponse) {
	m.metrics.TotalRequests++
	m.metrics.TotalDuration += duration

	if err == nil && (resp == nil || !resp.IsError) {
		m.metrics.SuccessRequests++
	} else {
		m.metrics.ErrorRequests++
	}

	if duration < m.metrics.MinDuration {
		m.metrics.MinDuration = duration
	}
	if duration > m.metrics.MaxDuration {
		m.metrics.MaxDuration = duration
	}
}

// updateToolMetrics updates tool-specific metrics
func (m *MetricsMiddleware) updateToolMetrics(toolName string, duration time.Duration, err error, resp *ToolResponse) {
	toolMetrics, exists := m.metrics.ToolMetrics[toolName]
	if !exists {
		toolMetrics = &ToolMetrics{
			MinDuration: time.Duration(1<<63 - 1), // Max duration as initial min
		}
		m.metrics.ToolMetrics[toolName] = toolMetrics
	}

	toolMetrics.TotalCalls++
	toolMetrics.TotalDuration += duration
	toolMetrics.LastCalled = time.Now()

	if err == nil && (resp == nil || !resp.IsError) {
		toolMetrics.SuccessCalls++
	} else {
		toolMetrics.ErrorCalls++
	}

	if duration < toolMetrics.MinDuration {
		toolMetrics.MinDuration = duration
	}
	if duration > toolMetrics.MaxDuration {
		toolMetrics.MaxDuration = duration
	}
}

// GetMetrics returns the current metrics
func (m *MetricsMiddleware) GetMetrics() *RequestMetrics {
	return m.metrics
}

// GetToolMetrics returns metrics for a specific tool
func (m *MetricsMiddleware) GetToolMetrics(toolName string) *ToolMetrics {
	return m.metrics.ToolMetrics[toolName]
}

// Name returns the middleware name
func (m *MetricsMiddleware) Name() string {
	return "metrics"
}

// ResetMetrics resets all metrics
func (m *MetricsMiddleware) ResetMetrics() {
	m.metrics = &RequestMetrics{
		ToolMetrics: make(map[string]*ToolMetrics),
		MinDuration: time.Duration(1<<63 - 1),
	}
}