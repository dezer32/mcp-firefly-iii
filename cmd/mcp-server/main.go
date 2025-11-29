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

	// Check if config file exists
	configFileExists := false
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			configFileExists = true
			log.Printf("Loading configuration from file: %s", configPath)
		} else if !os.IsNotExist(err) {
			log.Fatalf("Error accessing config file: %v", err)
		}
	}

	if !configFileExists {
		log.Printf("Config file not found, using environment variables and defaults")
	}

	config, err := fireflyMCP.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Configuration loaded successfully")

	// Create MCP server
	server, err := fireflyMCP.NewFireflyMCPServer(config)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Run the server over stdin/stdout
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
