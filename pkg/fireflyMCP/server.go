package fireflyMCP

import (
	"context"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/handlers"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/middleware"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FireflyMCPServer represents the MCP server for Firefly III
type FireflyMCPServer struct {
	Server   *mcp.Server
	Client   *client.ClientWithResponses
	Config   *Config
	Handlers *handlers.HandlerRegistryImpl
	Chain    *middleware.Chain
}

// NewServer creates a new FireflyMCPServer instance
// Deprecated: Use NewServerWithOptions for more flexible configuration
func NewServer(config *Config) (*FireflyMCPServer, error) {
	// Use the new functional options pattern internally
	return NewServerWithOptions(WithConfig(config))
}

// Run starts the MCP server with the given transport
func (s *FireflyMCPServer) Run(ctx context.Context, transport mcp.Transport) error {
	return s.Server.Run(ctx, transport)
}

// GetHandlers returns the handler registry for testing purposes
func (s *FireflyMCPServer) GetHandlers() *handlers.HandlerRegistryImpl {
	return s.Handlers
}

// GetClient returns the Firefly III client for testing purposes
func (s *FireflyMCPServer) GetClient() *client.ClientWithResponses {
	return s.Client
}

// GetConfig returns the configuration for testing purposes
func (s *FireflyMCPServer) GetConfig() *Config {
	return s.Config
}