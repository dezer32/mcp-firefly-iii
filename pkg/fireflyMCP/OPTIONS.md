# Functional Options Pattern Documentation

## Overview

The Firefly III MCP server now supports the functional options pattern for flexible configuration. This pattern provides a clean and extensible way to configure both the server and client instances.

## Benefits

- **Backward Compatibility**: The old `NewServer(config)` constructor still works
- **Flexible Configuration**: Mix and match configuration options as needed
- **Type Safety**: Compile-time validation of options
- **Extensibility**: Easy to add new options without breaking existing code
- **Default Values**: Sensible defaults for all optional settings

## Server Configuration

### Basic Usage

```go
server, err := fireflyMCP.NewServerWithOptions(
    fireflyMCP.WithBaseURL("https://firefly.example.com"),
    fireflyMCP.WithAPIToken("your-api-token"),
)
```

### Available Server Options

#### API Configuration
- `WithBaseURL(url string)` - Sets the Firefly III base URL
- `WithAPIToken(token string)` - Sets the API authentication token
- `WithRequestEditor(editor)` - Adds custom request editor functions

#### HTTP Client Configuration
- `WithHTTPClient(client *http.Client)` - Uses a custom HTTP client
- `WithTimeout(timeout time.Duration)` - Sets request timeout
- `WithConnectionPool(maxIdle, maxPerHost int)` - Configures connection pooling

#### MCP Configuration
- `WithMCPInfo(name, version string)` - Sets MCP server name and version

#### Middleware Configuration
- `WithLogging(enabled bool, level LogLevel)` - Configures logging
- `WithMetrics(enabled bool)` - Enables metrics collection
- `WithRecovery(enabled bool)` - Enables panic recovery
- `WithTracing(enabled bool)` - Enables distributed tracing
- `WithMiddleware(mw Middleware)` - Adds custom middleware

#### Performance Configuration
- `WithRateLimit(limit, burst int)` - Sets rate limiting
- `WithCache(enabled bool, ttl time.Duration)` - Configures caching

#### Configuration Object
- `WithConfig(config *Config)` - Applies existing configuration

### Advanced Example

```go
server, err := fireflyMCP.NewServerWithOptions(
    // API configuration
    fireflyMCP.WithBaseURL("https://firefly.example.com"),
    fireflyMCP.WithAPIToken("your-api-token"),
    
    // HTTP client configuration
    fireflyMCP.WithTimeout(45*time.Second),
    fireflyMCP.WithConnectionPool(200, 20),
    
    // Middleware configuration
    fireflyMCP.WithLogging(true, middleware.LogLevelDebug),
    fireflyMCP.WithMetrics(true),
    fireflyMCP.WithRecovery(true),
    
    // Performance configuration
    fireflyMCP.WithRateLimit(150, 15),
    fireflyMCP.WithCache(true, 10*time.Minute),
)
```

## Client Configuration

### Basic Usage

```go
client, err := fireflyMCP.NewFireflyClient(
    fireflyMCP.WithClientBaseURL("https://firefly.example.com"),
    fireflyMCP.WithClientAPIToken("your-api-token"),
)
```

### Available Client Options

- `WithClientHTTPClient(client *http.Client)` - Custom HTTP client
- `WithClientTimeout(timeout time.Duration)` - Request timeout
- `WithClientBaseURL(url string)` - Firefly III base URL
- `WithClientAPIToken(token string)` - API token
- `WithClientRequestEditor(editor)` - Request editor function
- `WithClientRetry(count int, wait time.Duration)` - Retry configuration
- `WithClientUserAgent(userAgent string)` - Custom User-Agent header

### Client Example

```go
client, err := fireflyMCP.NewFireflyClient(
    fireflyMCP.WithClientBaseURL("https://firefly.example.com"),
    fireflyMCP.WithClientAPIToken("your-api-token"),
    fireflyMCP.WithClientTimeout(30*time.Second),
    fireflyMCP.WithClientRetry(3, 2*time.Second),
    fireflyMCP.WithClientUserAgent("my-app/1.0"),
)
```

## Backward Compatibility

The old constructor still works and internally uses the functional options:

```go
config := &fireflyMCP.Config{}
config.Server.URL = "https://firefly.example.com"
config.API.Token = "your-api-token"
config.Client.Timeout = 30

// This still works (deprecated but supported)
server, err := fireflyMCP.NewServer(config)
```

## Error Handling

All options return errors if invalid values are provided:

```go
server, err := fireflyMCP.NewServerWithOptions(
    fireflyMCP.WithTimeout(-1), // Error: invalid timeout
)
if err != nil {
    // Handle error: ErrInvalidTimeout
}
```

## Common Errors

- `ErrNilHTTPClient` - HTTP client cannot be nil
- `ErrInvalidTimeout` - Timeout must be greater than 0
- `ErrEmptyAPIToken` - API token cannot be empty
- `ErrEmptyBaseURL` - Base URL cannot be empty
- `ErrNilRequestEditor` - Request editor cannot be nil
- `ErrNilMiddleware` - Middleware cannot be nil
- `ErrInvalidRateLimit` - Rate limit must be positive
- `ErrNilConfig` - Config cannot be nil
- `ErrInvalidRetryCount` - Retry count cannot be negative

## Custom Middleware

You can add custom middleware using the functional options:

```go
customMiddleware := middleware.MiddlewareFunc(func(next middleware.Handler) middleware.Handler {
    return middleware.HandlerFunc(func(ctx context.Context, method string, params interface{}) (interface{}, error) {
        // Custom logic before
        log.Printf("Calling method: %s", method)
        
        // Call next handler
        result, err := next.Handle(ctx, method, params)
        
        // Custom logic after
        log.Printf("Method completed: %s", method)
        
        return result, err
    })
})

server, err := fireflyMCP.NewServerWithOptions(
    fireflyMCP.WithBaseURL("https://firefly.example.com"),
    fireflyMCP.WithAPIToken("your-api-token"),
    fireflyMCP.WithMiddleware(customMiddleware),
)
```

## Migration Guide

### From Old Constructor

Before:
```go
config := &Config{...}
server, err := NewServer(config)
```

After:
```go
server, err := NewServerWithOptions(
    WithConfig(config), // Use existing config
    // Add additional options as needed
    WithLogging(true),
    WithMetrics(true),
)
```

### From Manual Configuration

Before:
```go
httpClient := &http.Client{Timeout: 30 * time.Second}
client, err := client.NewClientWithResponses(url, 
    client.WithHTTPClient(httpClient),
    client.WithRequestEditorFn(authEditor),
)
```

After:
```go
client, err := NewFireflyClient(
    WithClientBaseURL(url),
    WithClientAPIToken(token),
    WithClientTimeout(30*time.Second),
)
```

## Best Practices

1. **Start Simple**: Begin with minimal required options
2. **Add As Needed**: Add additional options based on requirements
3. **Use Defaults**: Rely on sensible defaults when possible
4. **Error Handling**: Always check for errors from option functions
5. **Custom Middleware**: Add middleware in the correct order (recovery first, then logging, etc.)

## Testing

The functional options pattern makes testing easier:

```go
func TestMyFunction(t *testing.T) {
    // Create test server with specific configuration
    server, err := NewServerWithOptions(
        WithBaseURL("http://test-server"),
        WithAPIToken("test-token"),
        WithLogging(false), // Disable logging in tests
        WithTimeout(1*time.Second), // Short timeout for tests
    )
    
    // Test your functionality
}
```

## Future Extensions

The functional options pattern makes it easy to add new features:

- Prometheus metrics exporter
- OpenTelemetry tracing
- Circuit breaker pattern
- Advanced retry strategies
- Request/response transformers
- Authentication providers
- Multi-tenant support

## Summary

The functional options pattern provides a powerful and flexible way to configure the Firefly III MCP server while maintaining backward compatibility. It allows for clean, readable configuration code that can grow with your needs.