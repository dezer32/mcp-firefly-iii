package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HTTPServer wraps the MCP server with HTTP transport.
type HTTPServer struct {
	mcpServer  *FireflyMCPServer
	httpServer *http.Server
	config     *Config
	logger     *slog.Logger
}

// NewHTTPServer creates a new HTTP server for the MCP server.
func NewHTTPServer(mcpServer *FireflyMCPServer, config *Config, logger *slog.Logger) *HTTPServer {
	if logger == nil {
		logger = slog.Default()
	}

	return &HTTPServer{
		mcpServer: mcpServer,
		config:    config,
		logger:    logger,
	}
}

// Start starts the HTTP server and blocks until the context is cancelled.
func (s *HTTPServer) Start(ctx context.Context) error {
	// Create Streamable HTTP handler from MCP SDK
	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return s.mcpServer.MCPServer()
		},
		&mcp.StreamableHTTPOptions{
			SessionTimeout: time.Duration(s.config.HTTP.SessionTimeout) * time.Second,
			Logger:         s.logger,
		},
	)

	// Build middleware chain (order matters: outer -> inner)
	// Request flow: logging -> rate limit -> CORS -> auth -> handler
	var h http.Handler = handler

	// Auth middleware (innermost - runs last before handler)
	h = BearerAuthMiddleware(s.config.HTTP.AuthToken, s.logger)(h)

	// CORS middleware
	h = CORSMiddleware(s.config.HTTP.AllowedOrigins, s.logger)(h)

	// Rate limiting middleware
	h = RateLimitMiddleware(s.config.HTTP.RateLimit, s.config.HTTP.RateBurst, s.logger)(h)

	// Logging middleware (outermost - runs first)
	h = RequestLoggingMiddleware(s.logger)(h)

	// Create router with health endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)
	mux.Handle("/", h)

	addr := fmt.Sprintf("%s:%d", s.config.HTTP.Host, s.config.HTTP.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(s.config.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.HTTP.IdleTimeout) * time.Second,
	}

	s.logger.Info("starting HTTP server",
		"addr", addr,
		"auth_enabled", s.config.HTTP.AuthToken != "",
		"rate_limit", s.config.HTTP.RateLimit,
		"rate_burst", s.config.HTTP.RateBurst)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	// Wait for context cancellation or error
	select {
	case <-ctx.Done():
		s.logger.Info("shutting down HTTP server...")
		return s.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// Shutdown gracefully shuts down the HTTP server.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	// Create shutdown context with 30 second timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("shutdown error", "error", err)
		return err
	}

	s.logger.Info("HTTP server stopped")
	return nil
}

// handleHealth returns the health status of the server.
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// handleReady returns whether the server is ready to accept requests.
func (s *HTTPServer) handleReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
