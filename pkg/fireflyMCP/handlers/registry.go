package handlers

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Registry defines the interface for managing MCP tool handlers
type Registry interface {
	// RegisterAll registers all available MCP tools with the server
	RegisterAll(server *mcp.Server)
}

// HandlerContext provides dependencies needed by handlers
type HandlerContext interface {
	// GetClient returns the Firefly III API client
	GetClient() interface{}
	// GetConfig returns the application configuration
	GetConfig() interface{}
}

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	Context HandlerContext
}

// NewBaseHandler creates a new BaseHandler instance
func NewBaseHandler(ctx HandlerContext) *BaseHandler {
	return &BaseHandler{
		Context: ctx,
	}
}

// HandleError creates a standardized error response
func (h *BaseHandler) HandleError(err error, message string) *mcp.CallToolResultFor[struct{}] {
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: message + ": " + err.Error()},
		},
		IsError: true,
	}
}

// HandleAPIError creates a standardized API error response
func (h *BaseHandler) HandleAPIError(statusCode int) *mcp.CallToolResultFor[struct{}] {
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "API error: " + string(rune(statusCode))},
		},
		IsError: true,
	}
}