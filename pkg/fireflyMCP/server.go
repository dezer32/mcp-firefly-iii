package fireflyMCP

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
func NewServer(config *Config) (*FireflyMCPServer, error) {
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: time.Duration(config.Client.Timeout) * time.Second,
	}

	// Create Firefly III client with request editor for authentication
	fireflyClient, err := client.NewClientWithResponses(
		config.Server.URL,
		client.WithHTTPClient(httpClient),
		client.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				req.Header.Set("Authorization", "Bearer "+config.API.Token)
				return nil
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firefly III client: %w", err)
	}

	// Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    config.MCP.Name,
			Version: config.MCP.Version,
		}, nil,
	)

	// Create handler registry
	handlerRegistry := handlers.NewHandlerRegistry(fireflyClient, config)

	server := &FireflyMCPServer{
		Server:   mcpServer,
		Client:   fireflyClient,
		Config:   config,
		Handlers: handlerRegistry,
		Chain:    nil, // No middleware chain in legacy constructor
	}

	// Register all tools through the handler registry
	handlerRegistry.RegisterAll(mcpServer)

	return server, nil
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