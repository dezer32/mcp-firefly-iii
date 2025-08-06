# Middleware Package

The middleware package provides a chain of responsibility pattern implementation for cross-cutting concerns in the Firefly III MCP server.

## Overview

This package enables:
- Request/response interception
- Context enrichment and propagation
- Composable processing pipelines
- Centralized error handling
- Performance monitoring
- Request logging and tracing

## Architecture

### Core Components

1. **Middleware Interface**: Defines the contract for all middleware
2. **Chain**: Manages sequential execution of middleware
3. **Handler**: Core function signature for processing requests
4. **Context Management**: Utilities for context propagation

### Chain Pattern

The middleware chain follows the chain of responsibility pattern:

```
Request → Middleware1 → Middleware2 → Middleware3 → Handler
           ↓                ↓              ↓           ↓
        Response ← Response ← Response ← Response ←────┘
```

Each middleware can:
- Process the request before passing it to the next handler
- Process the response after the handler completes
- Short-circuit the chain by returning early
- Modify request/response data
- Add metadata and context

## Usage Examples

### Basic Setup

```go
package main

import (
    "log"
    "time"
    
    "github.com/firefly-iii/pkg/fireflyMCP/middleware"
)

func main() {
    // Create middleware instances
    logger := log.Default()
    
    recovery := middleware.NewRecoveryMiddleware(logger, true)
    logging := middleware.NewLoggingMiddleware(logger, middleware.LogLevelInfo)
    timing := middleware.NewTimingMiddleware(logger, 500*time.Millisecond)
    metrics := middleware.NewMetricsMiddleware()
    
    // Create middleware chain
    chain := middleware.NewChain(
        recovery,  // First: catch panics
        logging,   // Second: log requests
        timing,    // Third: measure timing
        metrics,   // Fourth: collect metrics
    )
    
    // Create your handler
    handler := func(req *middleware.ToolRequest) (*middleware.ToolResponse, error) {
        // Your business logic here
        return &middleware.ToolResponse{
            Result:  "Success",
            IsError: false,
        }, nil
    }
    
    // Apply middleware chain
    wrappedHandler := chain.Then(handler)
    
    // Use the wrapped handler
    req := &middleware.ToolRequest{
        ToolName:  "list_accounts",
        Arguments: map[string]interface{}{"limit": 10},
        Context:   context.Background(),
        Metadata:  make(map[string]interface{}),
        StartTime: time.Now(),
    }
    
    resp, err := wrappedHandler(req)
}
```

### Context Propagation

```go
// Add request ID to context
ctx := middleware.WithRequestID(context.Background(), "req-123")
ctx = middleware.WithUserID(ctx, "user-456")
ctx = middleware.WithTraceID(ctx, "trace-789")

req := &middleware.ToolRequest{
    ToolName:  "get_transaction",
    Context:   ctx,
    // ... other fields
}

// Retrieve from context in handler
handler := func(req *middleware.ToolRequest) (*middleware.ToolResponse, error) {
    requestID, _ := middleware.GetRequestID(req.Context)
    userID, _ := middleware.GetUserID(req.Context)
    traceID, _ := middleware.GetTraceID(req.Context)
    
    log.Printf("Processing request: %s for user: %s (trace: %s)", 
        requestID, userID, traceID)
    
    // Process request...
}
```

### Custom Middleware

```go
// Create custom middleware for authentication
type AuthMiddleware struct {
    apiKey string
}

func (a *AuthMiddleware) Process(next middleware.Handler) middleware.Handler {
    return func(req *middleware.ToolRequest) (*middleware.ToolResponse, error) {
        // Check authentication
        providedKey, ok := req.Metadata["api_key"].(string)
        if !ok || providedKey != a.apiKey {
            return &middleware.ToolResponse{
                IsError: true,
                Metadata: map[string]interface{}{
                    "error": "unauthorized",
                },
            }, fmt.Errorf("invalid API key")
        }
        
        // Add auth info to context
        req.Context = context.WithValue(req.Context, "authenticated", true)
        
        // Continue to next handler
        return next(req)
    }
}

func (a *AuthMiddleware) Name() string {
    return "auth"
}
```

### Dynamic Chain Modification

```go
// Start with basic chain
chain := middleware.NewChain(recovery, logging)

// Add middleware conditionally
if config.EnableMetrics {
    chain = chain.Append(metrics)
}

if config.EnableRateLimit {
    chain = chain.Append(rateLimit)
}

// Prepend critical middleware
chain = chain.Prepend(authentication)
```

### Integration with MCP Handlers

```go
// Adapt MCP handler to use middleware
adapter := middleware.NewHandlerAdapter(chain, mcpHandler)

// Register with MCP server
server.RegisterTool("list_accounts", func(args interface{}) (*mcp.CallToolResultFor[interface{}], error) {
    return adapter.Handle("list_accounts", args)
})
```

### Metrics Collection

```go
// Create metrics middleware
metrics := middleware.NewMetricsMiddleware()

// Use in chain
chain := middleware.NewChain(metrics)

// ... process requests ...

// Get metrics
globalMetrics := metrics.GetMetrics()
fmt.Printf("Total requests: %d\n", globalMetrics.TotalRequests)
fmt.Printf("Success rate: %.2f%%\n", 
    float64(globalMetrics.SuccessRequests)/float64(globalMetrics.TotalRequests)*100)
fmt.Printf("Average duration: %v\n", 
    globalMetrics.TotalDuration/time.Duration(globalMetrics.TotalRequests))

// Get tool-specific metrics
toolMetrics := metrics.GetToolMetrics("list_accounts")
if toolMetrics != nil {
    fmt.Printf("list_accounts calls: %d\n", toolMetrics.TotalCalls)
    fmt.Printf("list_accounts avg duration: %v\n", 
        toolMetrics.TotalDuration/time.Duration(toolMetrics.TotalCalls))
}
```

## Available Middleware

### LoggingMiddleware
Logs requests and responses with configurable levels:
- `LogLevelDebug`: Logs everything including request/response bodies
- `LogLevelInfo`: Logs basic request/response information
- `LogLevelWarn`: Logs warnings and errors
- `LogLevelError`: Logs only errors

### RecoveryMiddleware
Recovers from panics and converts them to error responses:
- Catches panics in handlers
- Logs stack traces (optional)
- Returns safe error responses

### TimingMiddleware
Measures request processing time:
- Tracks duration for each request
- Logs slow requests above threshold
- Adds timing metadata to responses

### MetricsMiddleware
Collects detailed metrics:
- Total request counts
- Success/error rates
- Min/max/average durations
- Per-tool statistics

### RequestLoggingMiddleware
Detailed request logging:
- Generates request IDs
- Logs request bodies (optional)
- Structured logging format

### SafeExecutionMiddleware
Ensures safe handler execution:
- Panic recovery
- Context cancellation support
- Custom panic handlers

## Best Practices

1. **Order Matters**: Place recovery middleware first to catch all panics
2. **Context Usage**: Use context for request-scoped values, not for configuration
3. **Metadata**: Use metadata for debugging information that shouldn't affect logic
4. **Error Handling**: Return errors from handlers, don't panic
5. **Performance**: Keep middleware lightweight, avoid blocking operations
6. **Testing**: Test middleware in isolation and as part of chains

## Testing

```go
func TestMiddlewareChain(t *testing.T) {
    // Create test middleware that modifies metadata
    testMiddleware := &TestMiddleware{}
    
    chain := middleware.NewChain(testMiddleware)
    
    handler := func(req *middleware.ToolRequest) (*middleware.ToolResponse, error) {
        // Verify middleware was applied
        if req.Metadata["test"] != "applied" {
            t.Error("Middleware not applied")
        }
        return &middleware.ToolResponse{}, nil
    }
    
    wrapped := chain.Then(handler)
    
    req := &middleware.ToolRequest{
        Metadata: make(map[string]interface{}),
    }
    
    _, err := wrapped(req)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
}
```

## Performance Considerations

- Middleware adds minimal overhead (typically <1μs per middleware)
- Use goroutine-safe implementations for concurrent requests
- Consider caching in middleware for expensive operations
- Monitor middleware performance using the MetricsMiddleware

## Future Enhancements

- [ ] Async middleware support
- [ ] Conditional middleware execution
- [ ] Middleware priorities
- [ ] Built-in rate limiting middleware
- [ ] Circuit breaker middleware
- [ ] Retry middleware with backoff
- [ ] Compression middleware
- [ ] Caching middleware