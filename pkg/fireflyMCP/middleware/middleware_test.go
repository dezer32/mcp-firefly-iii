package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

// TestMiddleware is a test middleware that records calls
type TestMiddleware struct {
	Name_     string
	Called    bool
	CallOrder int
	Counter   *int
}

func (t *TestMiddleware) Process(next Handler) Handler {
	return func(req *ToolRequest) (*ToolResponse, error) {
		t.Called = true
		if t.Counter != nil {
			*t.Counter++
			t.CallOrder = *t.Counter
		}
		
		// Add marker to metadata
		if req.Metadata == nil {
			req.Metadata = make(map[string]interface{})
		}
		req.Metadata[t.Name_] = "processed"
		
		// Call next
		resp, err := next(req)
		
		// Add marker to response
		if resp != nil {
			if resp.Metadata == nil {
				resp.Metadata = make(map[string]interface{})
			}
			resp.Metadata[t.Name_] = "processed"
		}
		
		return resp, err
	}
}

func (t *TestMiddleware) Name() string {
	return t.Name_
}

func TestChainExecution(t *testing.T) {
	counter := 0
	m1 := &TestMiddleware{Name_: "m1", Counter: &counter}
	m2 := &TestMiddleware{Name_: "m2", Counter: &counter}
	m3 := &TestMiddleware{Name_: "m3", Counter: &counter}
	
	chain := NewChain(m1, m2, m3)
	
	handlerCalled := false
	handler := func(req *ToolRequest) (*ToolResponse, error) {
		handlerCalled = true
		
		// Verify all middleware processed the request
		if req.Metadata["m1"] != "processed" {
			t.Error("m1 did not process request")
		}
		if req.Metadata["m2"] != "processed" {
			t.Error("m2 did not process request")
		}
		if req.Metadata["m3"] != "processed" {
			t.Error("m3 did not process request")
		}
		
		return &ToolResponse{
			Result: "success",
		}, nil
	}
	
	wrapped := chain.Then(handler)
	
	req := &ToolRequest{
		ToolName: "test",
		Context:  context.Background(),
		Metadata: make(map[string]interface{}),
	}
	
	resp, err := wrapped(req)
	
	// Verify execution
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("Handler was not called")
	}
	if !m1.Called || !m2.Called || !m3.Called {
		t.Error("Not all middleware were called")
	}
	
	// Verify order
	if m1.CallOrder != 1 || m2.CallOrder != 2 || m3.CallOrder != 3 {
		t.Errorf("Incorrect call order: m1=%d, m2=%d, m3=%d", 
			m1.CallOrder, m2.CallOrder, m3.CallOrder)
	}
	
	// Verify response
	if resp.Metadata["m1"] != "processed" || 
	   resp.Metadata["m2"] != "processed" || 
	   resp.Metadata["m3"] != "processed" {
		t.Error("Not all middleware processed the response")
	}
}

func TestChainAppend(t *testing.T) {
	m1 := &TestMiddleware{Name_: "m1"}
	m2 := &TestMiddleware{Name_: "m2"}
	m3 := &TestMiddleware{Name_: "m3"}
	
	chain := NewChain(m1)
	chain = chain.Append(m2, m3)
	
	if len(chain.middlewares) != 3 {
		t.Errorf("Expected 3 middleware, got %d", len(chain.middlewares))
	}
	
	// Verify order
	if chain.middlewares[0].Name() != "m1" ||
	   chain.middlewares[1].Name() != "m2" ||
	   chain.middlewares[2].Name() != "m3" {
		t.Error("Incorrect middleware order after append")
	}
}

func TestChainPrepend(t *testing.T) {
	m1 := &TestMiddleware{Name_: "m1"}
	m2 := &TestMiddleware{Name_: "m2"}
	m3 := &TestMiddleware{Name_: "m3"}
	
	chain := NewChain(m3)
	chain = chain.Prepend(m1, m2)
	
	if len(chain.middlewares) != 3 {
		t.Errorf("Expected 3 middleware, got %d", len(chain.middlewares))
	}
	
	// Verify order
	if chain.middlewares[0].Name() != "m1" ||
	   chain.middlewares[1].Name() != "m2" ||
	   chain.middlewares[2].Name() != "m3" {
		t.Error("Incorrect middleware order after prepend")
	}
}

func TestContextFunctions(t *testing.T) {
	ctx := context.Background()
	
	// Test RequestID
	ctx = WithRequestID(ctx, "req-123")
	id, ok := GetRequestID(ctx)
	if !ok || id != "req-123" {
		t.Errorf("RequestID not set correctly: got %s, ok=%v", id, ok)
	}
	
	// Test UserID
	ctx = WithUserID(ctx, "user-456")
	id, ok = GetUserID(ctx)
	if !ok || id != "user-456" {
		t.Errorf("UserID not set correctly: got %s, ok=%v", id, ok)
	}
	
	// Test TraceID
	ctx = WithTraceID(ctx, "trace-789")
	id, ok = GetTraceID(ctx)
	if !ok || id != "trace-789" {
		t.Errorf("TraceID not set correctly: got %s, ok=%v", id, ok)
	}
	
	// Test SpanID
	ctx = WithSpanID(ctx, "span-012")
	id, ok = GetSpanID(ctx)
	if !ok || id != "span-012" {
		t.Errorf("SpanID not set correctly: got %s, ok=%v", id, ok)
	}
	
	// Test missing values
	emptyCtx := context.Background()
	_, ok = GetRequestID(emptyCtx)
	if ok {
		t.Error("GetRequestID should return false for missing value")
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	var logOutput strings.Builder
	logger := log.New(&logOutput, "", 0)
	
	recovery := NewRecoveryMiddleware(logger, false)
	
	// Handler that panics
	panicHandler := func(req *ToolRequest) (*ToolResponse, error) {
		panic("test panic")
	}
	
	wrapped := recovery.Process(panicHandler)
	
	req := &ToolRequest{
		ToolName: "test",
		Context:  context.Background(),
	}
	
	resp, err := wrapped(req)
	
	// Verify panic was recovered
	if err == nil {
		t.Error("Expected error from panic recovery")
	}
	if !strings.Contains(err.Error(), "test panic") {
		t.Errorf("Error should contain panic message: %v", err)
	}
	if resp == nil || !resp.IsError {
		t.Error("Response should indicate error")
	}
	if resp.Metadata["panic"] != true {
		t.Error("Response metadata should indicate panic")
	}
	
	// Verify logging
	if !strings.Contains(logOutput.String(), "[PANIC]") {
		t.Error("Panic should be logged")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	var logOutput strings.Builder
	logger := log.New(&logOutput, "", 0)
	
	logging := NewLoggingMiddleware(logger, LogLevelDebug)
	
	handler := func(req *ToolRequest) (*ToolResponse, error) {
		return &ToolResponse{
			Result: "success",
		}, nil
	}
	
	wrapped := logging.Process(handler)
	
	ctx := WithRequestID(context.Background(), "test-123")
	req := &ToolRequest{
		ToolName:  "test_tool",
		Context:   ctx,
		Arguments: map[string]interface{}{"arg": "value"},
		Metadata:  map[string]interface{}{"meta": "data"},
		StartTime: time.Now(),
	}
	
	_, err := wrapped(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	output := logOutput.String()
	
	// Verify logging output
	if !strings.Contains(output, "[INFO] Request") {
		t.Error("Request should be logged at INFO level")
	}
	if !strings.Contains(output, "test_tool") {
		t.Error("Tool name should be logged")
	}
	if !strings.Contains(output, "test-123") {
		t.Error("Request ID should be logged")
	}
	if !strings.Contains(output, "[DEBUG]") {
		t.Error("Debug information should be logged")
	}
	if !strings.Contains(output, "[INFO] Response") {
		t.Error("Response should be logged")
	}
}

func TestTimingMiddleware(t *testing.T) {
	var logOutput strings.Builder
	logger := log.New(&logOutput, "", 0)
	
	// Set low threshold to ensure slow request logging
	timing := NewTimingMiddleware(logger, 1*time.Nanosecond)
	
	handler := func(req *ToolRequest) (*ToolResponse, error) {
		time.Sleep(10 * time.Millisecond)
		return &ToolResponse{}, nil
	}
	
	wrapped := timing.Process(handler)
	
	req := &ToolRequest{
		ToolName: "slow_tool",
		Context:  context.Background(),
	}
	
	resp, err := wrapped(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// Verify timing was added to response
	if resp.Duration == 0 {
		t.Error("Duration should be set in response")
	}
	if resp.Metadata["duration_ms"] == nil {
		t.Error("Duration in milliseconds should be in metadata")
	}
	if resp.Metadata["duration_ns"] == nil {
		t.Error("Duration in nanoseconds should be in metadata")
	}
	
	// Verify slow request was logged
	if !strings.Contains(logOutput.String(), "[SLOW_REQUEST]") {
		t.Error("Slow request should be logged")
	}
}

func TestMetricsMiddleware(t *testing.T) {
	metrics := NewMetricsMiddleware()
	
	successHandler := func(req *ToolRequest) (*ToolResponse, error) {
		return &ToolResponse{Result: "success"}, nil
	}
	
	errorHandler := func(req *ToolRequest) (*ToolResponse, error) {
		return nil, errors.New("test error")
	}
	
	// Process successful request
	wrappedSuccess := metrics.Process(successHandler)
	req := &ToolRequest{ToolName: "test_tool", Context: context.Background()}
	
	_, _ = wrappedSuccess(req)
	
	// Process error request
	wrappedError := metrics.Process(errorHandler)
	_, _ = wrappedError(req)
	
	// Process another successful request for different tool
	req2 := &ToolRequest{ToolName: "other_tool", Context: context.Background()}
	_, _ = wrappedSuccess(req2)
	
	// Verify global metrics
	globalMetrics := metrics.GetMetrics()
	if globalMetrics.TotalRequests != 3 {
		t.Errorf("Expected 3 total requests, got %d", globalMetrics.TotalRequests)
	}
	if globalMetrics.SuccessRequests != 2 {
		t.Errorf("Expected 2 success requests, got %d", globalMetrics.SuccessRequests)
	}
	if globalMetrics.ErrorRequests != 1 {
		t.Errorf("Expected 1 error request, got %d", globalMetrics.ErrorRequests)
	}
	
	// Verify tool-specific metrics
	toolMetrics := metrics.GetToolMetrics("test_tool")
	if toolMetrics == nil {
		t.Fatal("Tool metrics should exist for test_tool")
	}
	if toolMetrics.TotalCalls != 2 {
		t.Errorf("Expected 2 calls for test_tool, got %d", toolMetrics.TotalCalls)
	}
	if toolMetrics.SuccessCalls != 1 {
		t.Errorf("Expected 1 success call for test_tool, got %d", toolMetrics.SuccessCalls)
	}
	if toolMetrics.ErrorCalls != 1 {
		t.Errorf("Expected 1 error call for test_tool, got %d", toolMetrics.ErrorCalls)
	}
}

func TestHandlerAdapter(t *testing.T) {
	// Create a simple middleware that adds metadata
	testMiddleware := &TestMiddleware{Name_: "adapter_test"}
	chain := NewChain(testMiddleware)
	
	// Create handler
	handler := func(req *ToolRequest) (*ToolResponse, error) {
		if req.ToolName != "test_tool" {
			t.Errorf("Expected tool name 'test_tool', got %s", req.ToolName)
		}
		return &ToolResponse{
			Result:  "success",
			IsError: false,
		}, nil
	}
	
	// Create adapter
	adapter := NewHandlerAdapter(chain, handler)
	
	// Call through adapter
	result, err := adapter.Handle("test_tool", map[string]interface{}{"arg": "value"})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.IsError {
		t.Error("Result should not be an error")
	}
	
	// Verify middleware was called
	if !testMiddleware.Called {
		t.Error("Middleware should have been called")
	}
}

func TestErrorPropagation(t *testing.T) {
	// Middleware that checks for errors
	errorCheckMiddleware := &TestMiddleware{Name_: "error_check"}
	
	chain := NewChain(errorCheckMiddleware)
	
	// Handler that returns an error
	errorHandler := func(req *ToolRequest) (*ToolResponse, error) {
		return nil, fmt.Errorf("test error")
	}
	
	wrapped := chain.Then(errorHandler)
	
	req := &ToolRequest{
		ToolName: "error_test",
		Context:  context.Background(),
		Metadata: make(map[string]interface{}),
	}
	
	_, err := wrapped(req)
	
	// Verify error propagation
	if err == nil {
		t.Error("Error should be propagated")
	}
	if !strings.Contains(err.Error(), "test error") {
		t.Errorf("Error message should be preserved: %v", err)
	}
	
	// Verify middleware was still called
	if !errorCheckMiddleware.Called {
		t.Error("Middleware should be called even with errors")
	}
}

func TestEmptyChain(t *testing.T) {
	// Create empty chain
	chain := NewChain()
	
	handlerCalled := false
	handler := func(req *ToolRequest) (*ToolResponse, error) {
		handlerCalled = true
		return &ToolResponse{Result: "success"}, nil
	}
	
	wrapped := chain.Then(handler)
	
	req := &ToolRequest{
		ToolName: "test",
		Context:  context.Background(),
	}
	
	resp, err := wrapped(req)
	
	// Verify handler is called directly with empty chain
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("Handler should be called with empty chain")
	}
	if resp.Result != "success" {
		t.Error("Response should be from handler")
	}
}