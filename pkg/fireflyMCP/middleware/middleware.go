// Package middleware provides a chain of responsibility pattern implementation
// for cross-cutting concerns in the Firefly III MCP server.
//
// The middleware package enables request/response interception, context enrichment,
// and composable processing pipelines for MCP tool calls.
package middleware

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolRequest represents an incoming MCP tool request
type ToolRequest struct {
	// ToolName is the name of the tool being called
	ToolName string
	// Arguments contains the tool-specific arguments
	Arguments interface{}
	// Context carries request-scoped values
	Context context.Context
	// Metadata contains additional request information
	Metadata map[string]interface{}
	// StartTime records when the request was received
	StartTime time.Time
}

// ToolResponse represents an MCP tool response
type ToolResponse struct {
	// Result is the actual tool response
	Result interface{}
	// IsError indicates if the response is an error
	IsError bool
	// Duration is the time taken to process the request
	Duration time.Duration
	// Metadata contains additional response information
	Metadata map[string]interface{}
}

// Handler defines the signature for tool handlers
type Handler func(req *ToolRequest) (*ToolResponse, error)

// Middleware defines the middleware interface
// Middleware functions wrap handlers to provide additional functionality
type Middleware interface {
	// Process wraps a handler with middleware logic
	Process(next Handler) Handler
	// Name returns the middleware name for logging/debugging
	Name() string
}

// Chain manages a sequence of middleware
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new middleware chain
func NewChain(middlewares ...Middleware) *Chain {
	return &Chain{
		middlewares: middlewares,
	}
}

// Then applies the middleware chain to a handler
func (c *Chain) Then(handler Handler) Handler {
	// Apply middleware in reverse order so they execute in the order added
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i].Process(handler)
	}
	return handler
}

// Append adds middleware to the end of the chain
func (c *Chain) Append(middlewares ...Middleware) *Chain {
	newMiddlewares := make([]Middleware, len(c.middlewares)+len(middlewares))
	copy(newMiddlewares, c.middlewares)
	copy(newMiddlewares[len(c.middlewares):], middlewares)
	return &Chain{middlewares: newMiddlewares}
}

// Prepend adds middleware to the beginning of the chain
func (c *Chain) Prepend(middlewares ...Middleware) *Chain {
	newMiddlewares := make([]Middleware, len(middlewares)+len(c.middlewares))
	copy(newMiddlewares, middlewares)
	copy(newMiddlewares[len(middlewares):], c.middlewares)
	return &Chain{middlewares: newMiddlewares}
}

// ContextKey is a type for context keys
type ContextKey string

const (
	// ContextKeyRequestID is the context key for request ID
	ContextKeyRequestID ContextKey = "request_id"
	// ContextKeyUserID is the context key for user ID
	ContextKeyUserID ContextKey = "user_id"
	// ContextKeyTraceID is the context key for trace ID
	ContextKeyTraceID ContextKey = "trace_id"
	// ContextKeySpanID is the context key for span ID
	ContextKeySpanID ContextKey = "span_id"
)

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ContextKeyRequestID).(string)
	return id, ok
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(string)
	return id, ok
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// GetTraceID retrieves the trace ID from context
func GetTraceID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ContextKeyTraceID).(string)
	return id, ok
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ContextKeyTraceID, traceID)
}

// GetSpanID retrieves the span ID from context
func GetSpanID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ContextKeySpanID).(string)
	return id, ok
}

// WithSpanID adds a span ID to the context
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, ContextKeySpanID, spanID)
}

// HandlerAdapter adapts MCP handlers to the middleware chain
type HandlerAdapter struct {
	chain   *Chain
	handler Handler
}

// NewHandlerAdapter creates a new handler adapter with middleware chain
func NewHandlerAdapter(chain *Chain, handler Handler) *HandlerAdapter {
	return &HandlerAdapter{
		chain:   chain,
		handler: chain.Then(handler),
	}
}

// Handle processes a request through the middleware chain
func (a *HandlerAdapter) Handle(toolName string, args interface{}) (*mcp.CallToolResultFor[interface{}], error) {
	req := &ToolRequest{
		ToolName:  toolName,
		Arguments: args,
		Context:   context.Background(),
		Metadata:  make(map[string]interface{}),
		StartTime: time.Now(),
	}

	resp, err := a.handler(req)
	if err != nil {
		return &mcp.CallToolResultFor[interface{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil
	}

	if resp.IsError {
		return &mcp.CallToolResultFor[interface{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error processing request"},
			},
			IsError: true,
		}, nil
	}

	// Convert response to MCP format
	// This is a simplified conversion - actual implementation would depend on response type
	return &mcp.CallToolResultFor[interface{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Success"},
		},
		IsError: false,
	}, nil
}