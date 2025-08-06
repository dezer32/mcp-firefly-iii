package fireflyMCP_test

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/middleware"
)

func ExampleNewServerWithOptions_basic() {
	// Create a server with basic configuration
	server, err := fireflyMCP.NewServerWithOptions(
		fireflyMCP.WithBaseURL("https://firefly.example.com"),
		fireflyMCP.WithAPIToken("your-api-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	// Use the server
	_ = server
}

func ExampleNewServerWithOptions_withCustomHTTPClient() {
	// Create a custom HTTP client
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 20,
		},
	}
	
	// Create server with custom HTTP client
	server, err := fireflyMCP.NewServerWithOptions(
		fireflyMCP.WithBaseURL("https://firefly.example.com"),
		fireflyMCP.WithAPIToken("your-api-token"),
		fireflyMCP.WithHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = server
}

func ExampleNewServerWithOptions_withMiddleware() {
	// Create server with custom middleware configuration
	server, err := fireflyMCP.NewServerWithOptions(
		fireflyMCP.WithBaseURL("https://firefly.example.com"),
		fireflyMCP.WithAPIToken("your-api-token"),
		fireflyMCP.WithLogging(true, middleware.LogLevelDebug),
		fireflyMCP.WithMetrics(true),
		fireflyMCP.WithRecovery(true),
		fireflyMCP.WithTracing(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = server
}

func ExampleNewServerWithOptions_withRateLimiting() {
	// Create server with rate limiting
	server, err := fireflyMCP.NewServerWithOptions(
		fireflyMCP.WithBaseURL("https://firefly.example.com"),
		fireflyMCP.WithAPIToken("your-api-token"),
		fireflyMCP.WithRateLimit(100, 10), // 100 requests per second with burst of 10
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = server
}

func ExampleNewServerWithOptions_withCaching() {
	// Create server with caching enabled
	server, err := fireflyMCP.NewServerWithOptions(
		fireflyMCP.WithBaseURL("https://firefly.example.com"),
		fireflyMCP.WithAPIToken("your-api-token"),
		fireflyMCP.WithCache(true, 10*time.Minute), // Cache for 10 minutes
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = server
}

func ExampleNewServerWithOptions_withCustomMiddleware() {
	// Create a custom middleware
	customMiddleware := middleware.MiddlewareFunc(func(next middleware.Handler) middleware.Handler {
		return middleware.HandlerFunc(func(ctx context.Context, method string, params interface{}) (interface{}, error) {
			// Custom logic before handling
			log.Printf("Handling method: %s", method)
			
			// Call the next handler
			result, err := next.Handle(ctx, method, params)
			
			// Custom logic after handling
			log.Printf("Method %s completed", method)
			
			return result, err
		})
	})
	
	// Create server with custom middleware
	server, err := fireflyMCP.NewServerWithOptions(
		fireflyMCP.WithBaseURL("https://firefly.example.com"),
		fireflyMCP.WithAPIToken("your-api-token"),
		fireflyMCP.WithMiddleware(customMiddleware),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = server
}

func ExampleNewServerWithOptions_withExistingConfig() {
	// Load existing configuration
	config := &fireflyMCP.Config{}
	config.Server.URL = "https://firefly.example.com"
	config.API.Token = "your-api-token"
	config.Client.Timeout = 30
	config.MCP.Name = "firefly-mcp"
	config.MCP.Version = "1.0.0"
	config.Limits.Accounts = 100
	config.Limits.Transactions = 100
	
	// Create server with existing config and additional options
	server, err := fireflyMCP.NewServerWithOptions(
		fireflyMCP.WithConfig(config),
		fireflyMCP.WithLogging(true, middleware.LogLevelInfo),
		fireflyMCP.WithMetrics(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = server
}

func ExampleNewServerWithOptions_complete() {
	// Create server with all options configured
	server, err := fireflyMCP.NewServerWithOptions(
		// API configuration
		fireflyMCP.WithBaseURL("https://firefly.example.com"),
		fireflyMCP.WithAPIToken("your-api-token"),
		
		// HTTP client configuration
		fireflyMCP.WithTimeout(45*time.Second),
		fireflyMCP.WithConnectionPool(200, 20),
		
		// MCP configuration
		fireflyMCP.WithMCPInfo("custom-firefly-mcp", "2.0.0"),
		
		// Middleware configuration
		fireflyMCP.WithLogging(true, middleware.LogLevelInfo),
		fireflyMCP.WithMetrics(true),
		fireflyMCP.WithRecovery(true),
		fireflyMCP.WithTracing(false),
		
		// Performance configuration
		fireflyMCP.WithRateLimit(150, 15),
		fireflyMCP.WithCache(true, 10*time.Minute),
		
		// Custom request editor
		fireflyMCP.WithRequestEditor(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-Custom-Header", "custom-value")
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = server
}

func ExampleNewFireflyClient_basic() {
	// Create a Firefly III client with basic configuration
	client, err := fireflyMCP.NewFireflyClient(
		fireflyMCP.WithClientBaseURL("https://firefly.example.com"),
		fireflyMCP.WithClientAPIToken("your-api-token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	// Use the client
	_ = client
}

func ExampleNewFireflyClient_withRetry() {
	// Create client with retry configuration
	client, err := fireflyMCP.NewFireflyClient(
		fireflyMCP.WithClientBaseURL("https://firefly.example.com"),
		fireflyMCP.WithClientAPIToken("your-api-token"),
		fireflyMCP.WithClientRetry(5, 2*time.Second), // Retry 5 times with 2 second wait
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = client
}

func ExampleNewFireflyClient_withCustomHeaders() {
	// Create client with custom headers
	client, err := fireflyMCP.NewFireflyClient(
		fireflyMCP.WithClientBaseURL("https://firefly.example.com"),
		fireflyMCP.WithClientAPIToken("your-api-token"),
		fireflyMCP.WithClientUserAgent("my-app/1.0"),
		fireflyMCP.WithClientRequestEditor(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-Request-ID", "unique-id-123")
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	
	_ = client
}

// Example of backward compatibility with the old constructor
func ExampleNewServer_backwardCompatibility() {
	// Create configuration
	config := &fireflyMCP.Config{}
	config.Server.URL = "https://firefly.example.com"
	config.API.Token = "your-api-token"
	config.Client.Timeout = 30
	config.MCP.Name = "firefly-mcp"
	config.MCP.Version = "1.0.0"
	config.Limits.Accounts = 100
	config.Limits.Transactions = 100
	
	// Use the old constructor (still works, but deprecated)
	server, err := fireflyMCP.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	
	// The old constructor now internally uses the new functional options
	_ = server
}