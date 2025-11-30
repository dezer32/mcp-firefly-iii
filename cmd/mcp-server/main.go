package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// CLI flags (override config file and env vars)
	transport := flag.String("transport", "", "Transport type: stdio (default) or http")
	port := flag.Int("port", 0, "HTTP port (overrides config)")
	authToken := flag.String("auth-token", "", "Bearer token for HTTP authentication (overrides config)")
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	logLevel := flag.String("log-level", "info", "Log level: debug, info, warn, error")
	flag.Parse()

	// Setup logger
	logger := setupLogger(*logLevel)

	// Check if config file exists
	configFileExists := false
	if *configPath != "" {
		if _, err := os.Stat(*configPath); err == nil {
			configFileExists = true
			log.Printf("Loading configuration from file: %s", *configPath)
		} else if !os.IsNotExist(err) {
			log.Fatalf("Error accessing config file: %v", err)
		}
	}

	if !configFileExists {
		log.Printf("Config file not found, using environment variables and defaults")
	}

	config, err := fireflyMCP.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// CLI flags override config (highest priority)
	if *transport != "" {
		config.HTTP.Enabled = (*transport == "http")
	}
	if *port != 0 {
		config.HTTP.Port = *port
	}
	if *authToken != "" {
		config.HTTP.AuthToken = *authToken
	}

	log.Printf("Configuration loaded successfully")

	// Create MCP server
	server, err := fireflyMCP.NewFireflyMCPServer(config)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Run based on transport type
	if config.HTTP.Enabled {
		runHTTPServer(server, config, logger)
	} else {
		runStdioServer(server)
	}
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: logLevel}
	handler := slog.NewJSONHandler(os.Stderr, opts)
	return slog.New(handler)
}

func runStdioServer(server *fireflyMCP.FireflyMCPServer) {
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func runHTTPServer(server *fireflyMCP.FireflyMCPServer, config *fireflyMCP.Config, logger *slog.Logger) {
	// Create HTTP server
	httpServer := fireflyMCP.NewHTTPServer(server, config, logger)

	// Setup context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Start HTTP server (blocks until context is cancelled)
	if err := httpServer.Start(ctx); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
