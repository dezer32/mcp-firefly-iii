package main

import (
	"context"
	"log"
	"os"

	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Load configuration
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := fireflyMCP.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create MCP server
	server, err := fireflyMCP.NewFireflyMCPServer(config)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Run the server over stdin/stdout
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
