// Package fireflyMCP provides functional options for configuring the Firefly MCP server and client
package fireflyMCP

import (
	"context"
	"net/http"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/handlers"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/middleware"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ServerOption represents a functional option for configuring the FireflyMCPServer
type ServerOption func(*serverOptions) error

// serverOptions holds all configurable options for the server
type serverOptions struct {
	// HTTP Client options
	httpClient     *http.Client
	timeout        time.Duration
	maxIdleConns   int
	maxConnsPerHost int
	
	// API options
	apiToken       string
	baseURL        string
	requestEditors []client.RequestEditorFn
	
	// MCP options
	mcpName        string
	mcpVersion     string
	
	// Middleware options
	middlewares    []middleware.Middleware
	enableLogging  bool
	enableMetrics  bool
	enableRecovery bool
	enableTracing  bool
	logLevel       middleware.LogLevel
	
	// Rate limiting
	rateLimit      int
	rateLimitBurst int
	
	// Cache options
	cacheEnabled   bool
	cacheTTL       time.Duration
	
	// Custom configuration override
	configOverride *Config
}

// defaultServerOptions returns default server options
func defaultServerOptions() *serverOptions {
	return &serverOptions{
		timeout:         30 * time.Second,
		maxIdleConns:    100,
		maxConnsPerHost: 10,
		mcpName:         "firefly-iii-mcp",
		mcpVersion:      "1.0.0",
		enableLogging:   true,
		enableMetrics:   true,
		enableRecovery:  true,
		enableTracing:   false,
		logLevel:        middleware.LogLevelInfo,
		rateLimit:       100,
		rateLimitBurst:  10,
		cacheEnabled:    false,
		cacheTTL:        5 * time.Minute,
		middlewares:     []middleware.Middleware{},
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ServerOption {
	return func(o *serverOptions) error {
		if client == nil {
			return ErrNilHTTPClient
		}
		o.httpClient = client
		return nil
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *serverOptions) error {
		if timeout <= 0 {
			return ErrInvalidTimeout
		}
		o.timeout = timeout
		return nil
	}
}

// WithAPIToken sets the Firefly III API token
func WithAPIToken(token string) ServerOption {
	return func(o *serverOptions) error {
		if token == "" {
			return ErrEmptyAPIToken
		}
		o.apiToken = token
		return nil
	}
}

// WithBaseURL sets the Firefly III base URL
func WithBaseURL(url string) ServerOption {
	return func(o *serverOptions) error {
		if url == "" {
			return ErrEmptyBaseURL
		}
		o.baseURL = url
		return nil
	}
}

// WithRequestEditor adds a request editor function to the client
func WithRequestEditor(editor client.RequestEditorFn) ServerOption {
	return func(o *serverOptions) error {
		if editor == nil {
			return ErrNilRequestEditor
		}
		o.requestEditors = append(o.requestEditors, editor)
		return nil
	}
}

// WithMCPInfo sets the MCP server name and version
func WithMCPInfo(name, version string) ServerOption {
	return func(o *serverOptions) error {
		if name != "" {
			o.mcpName = name
		}
		if version != "" {
			o.mcpVersion = version
		}
		return nil
	}
}

// WithMiddleware adds custom middleware to the chain
func WithMiddleware(mw middleware.Middleware) ServerOption {
	return func(o *serverOptions) error {
		if mw == nil {
			return ErrNilMiddleware
		}
		o.middlewares = append(o.middlewares, mw)
		return nil
	}
}

// WithLogging enables or disables logging middleware
func WithLogging(enabled bool, level ...middleware.LogLevel) ServerOption {
	return func(o *serverOptions) error {
		o.enableLogging = enabled
		if len(level) > 0 {
			o.logLevel = level[0]
		}
		return nil
	}
}

// WithMetrics enables or disables metrics collection
func WithMetrics(enabled bool) ServerOption {
	return func(o *serverOptions) error {
		o.enableMetrics = enabled
		return nil
	}
}

// WithRecovery enables or disables panic recovery
func WithRecovery(enabled bool) ServerOption {
	return func(o *serverOptions) error {
		o.enableRecovery = enabled
		return nil
	}
}

// WithTracing enables or disables distributed tracing
func WithTracing(enabled bool) ServerOption {
	return func(o *serverOptions) error {
		o.enableTracing = enabled
		return nil
	}
}

// WithRateLimit sets the rate limit for API calls
func WithRateLimit(limit, burst int) ServerOption {
	return func(o *serverOptions) error {
		if limit <= 0 || burst <= 0 {
			return ErrInvalidRateLimit
		}
		o.rateLimit = limit
		o.rateLimitBurst = burst
		return nil
	}
}

// WithCache enables caching with the specified TTL
func WithCache(enabled bool, ttl time.Duration) ServerOption {
	return func(o *serverOptions) error {
		o.cacheEnabled = enabled
		if ttl > 0 {
			o.cacheTTL = ttl
		}
		return nil
	}
}

// WithConnectionPool configures the HTTP connection pool
func WithConnectionPool(maxIdle, maxPerHost int) ServerOption {
	return func(o *serverOptions) error {
		if maxIdle > 0 {
			o.maxIdleConns = maxIdle
		}
		if maxPerHost > 0 {
			o.maxConnsPerHost = maxPerHost
		}
		return nil
	}
}

// WithConfig applies a configuration object
func WithConfig(config *Config) ServerOption {
	return func(o *serverOptions) error {
		if config == nil {
			return ErrNilConfig
		}
		o.configOverride = config
		// Apply config values to options
		if config.API.Token != "" {
			o.apiToken = config.API.Token
		}
		if config.Server.URL != "" {
			o.baseURL = config.Server.URL
		}
		if config.Client.Timeout > 0 {
			o.timeout = time.Duration(config.Client.Timeout) * time.Second
		}
		if config.MCP.Name != "" {
			o.mcpName = config.MCP.Name
		}
		if config.MCP.Version != "" {
			o.mcpVersion = config.MCP.Version
		}
		return nil
	}
}

// NewServerWithOptions creates a new FireflyMCPServer with functional options
func NewServerWithOptions(opts ...ServerOption) (*FireflyMCPServer, error) {
	// Start with default options
	options := defaultServerOptions()
	
	// Apply all provided options
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}
	
	// Validate required fields
	if options.apiToken == "" && (options.configOverride == nil || options.configOverride.API.Token == "") {
		return nil, ErrEmptyAPIToken
	}
	if options.baseURL == "" && (options.configOverride == nil || options.configOverride.Server.URL == "") {
		return nil, ErrEmptyBaseURL
	}
	
	// Build configuration from options
	config := options.buildConfig()
	
	// Create HTTP client if not provided
	httpClient := options.httpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: options.timeout,
			Transport: &http.Transport{
				MaxIdleConns:        options.maxIdleConns,
				MaxIdleConnsPerHost: options.maxConnsPerHost,
			},
		}
	}
	
	// Build request editors
	requestEditors := options.requestEditors
	if options.apiToken != "" {
		requestEditors = append(requestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+options.apiToken)
			return nil
		})
	}
	
	// Create Firefly III client
	clientOpts := []client.ClientOption{
		client.WithHTTPClient(httpClient),
	}
	for _, editor := range requestEditors {
		clientOpts = append(clientOpts, client.WithRequestEditorFn(editor))
	}
	
	fireflyClient, err := client.NewClientWithResponses(options.baseURL, clientOpts...)
	if err != nil {
		return nil, err
	}
	
	// Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    options.mcpName,
			Version: options.mcpVersion,
		}, nil,
	)
	
	// Build middleware chain
	chain := options.buildMiddlewareChain()
	
	// Create handler registry using the new pattern
	handlerRegistry := handlers.NewHandlerRegistry(fireflyClient, config)
	
	server := &FireflyMCPServer{
		Server:   mcpServer,
		Client:   fireflyClient,
		Config:   config,
		Handlers: handlerRegistry,
		Chain:    chain,
	}
	
	// Register all tools
	handlerRegistry.RegisterAll(mcpServer)
	
	return server, nil
}

// buildConfig creates a Config from serverOptions
func (o *serverOptions) buildConfig() *Config {
	if o.configOverride != nil {
		return o.configOverride
	}
	
	config := &Config{}
	config.Server.URL = o.baseURL
	config.API.Token = o.apiToken
	config.Client.Timeout = int(o.timeout.Seconds())
	config.MCP.Name = o.mcpName
	config.MCP.Version = o.mcpVersion
	config.Limits.Accounts = 100
	config.Limits.Transactions = 100
	config.Limits.Categories = 100
	config.Limits.Budgets = 100
	
	return config
}

// buildMiddlewareChain creates the middleware chain from options
func (o *serverOptions) buildMiddlewareChain() *middleware.Chain {
	var middlewares []middleware.Middleware
	
	// Add recovery middleware first if enabled
	if o.enableRecovery {
		middlewares = append(middlewares, middleware.NewRecoveryMiddleware(nil, true))
	}
	
	// Add logging middleware if enabled
	if o.enableLogging {
		middlewares = append(middlewares, middleware.NewLoggingMiddleware(nil, o.logLevel))
	}
	
	// Add timing middleware
	middlewares = append(middlewares, middleware.NewTimingMiddleware(nil, 1*time.Second))
	
	// Add metrics middleware if enabled
	if o.enableMetrics {
		middlewares = append(middlewares, middleware.NewMetricsMiddleware())
	}
	
	// Add custom middlewares
	middlewares = append(middlewares, o.middlewares...)
	
	return middleware.NewChain(middlewares...)
}

// ClientOption represents a functional option for configuring the Firefly III client
type ClientOption func(*clientOptions) error

// clientOptions holds all configurable options for the client
type clientOptions struct {
	httpClient     *http.Client
	timeout        time.Duration
	baseURL        string
	apiToken       string
	requestEditors []client.RequestEditorFn
	retryCount     int
	retryWait      time.Duration
	userAgent      string
}

// defaultClientOptions returns default client options
func defaultClientOptions() *clientOptions {
	return &clientOptions{
		timeout:    30 * time.Second,
		retryCount: 3,
		retryWait:  1 * time.Second,
		userAgent:  "firefly-iii-mcp/1.0.0",
	}
}

// WithClientHTTPClient sets a custom HTTP client for the Firefly client
func WithClientHTTPClient(httpClient *http.Client) ClientOption {
	return func(o *clientOptions) error {
		if httpClient == nil {
			return ErrNilHTTPClient
		}
		o.httpClient = httpClient
		return nil
	}
}

// WithClientTimeout sets the timeout for the Firefly client
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) error {
		if timeout <= 0 {
			return ErrInvalidTimeout
		}
		o.timeout = timeout
		return nil
	}
}

// WithClientBaseURL sets the base URL for the Firefly client
func WithClientBaseURL(url string) ClientOption {
	return func(o *clientOptions) error {
		if url == "" {
			return ErrEmptyBaseURL
		}
		o.baseURL = url
		return nil
	}
}

// WithClientAPIToken sets the API token for the Firefly client
func WithClientAPIToken(token string) ClientOption {
	return func(o *clientOptions) error {
		if token == "" {
			return ErrEmptyAPIToken
		}
		o.apiToken = token
		return nil
	}
}

// WithClientRequestEditor adds a request editor to the Firefly client
func WithClientRequestEditor(editor client.RequestEditorFn) ClientOption {
	return func(o *clientOptions) error {
		if editor == nil {
			return ErrNilRequestEditor
		}
		o.requestEditors = append(o.requestEditors, editor)
		return nil
	}
}

// WithClientRetry configures retry behavior for the client
func WithClientRetry(count int, wait time.Duration) ClientOption {
	return func(o *clientOptions) error {
		if count < 0 {
			return ErrInvalidRetryCount
		}
		o.retryCount = count
		if wait > 0 {
			o.retryWait = wait
		}
		return nil
	}
}

// WithClientUserAgent sets a custom User-Agent header
func WithClientUserAgent(userAgent string) ClientOption {
	return func(o *clientOptions) error {
		if userAgent != "" {
			o.userAgent = userAgent
		}
		return nil
	}
}

// NewFireflyClient creates a new Firefly III client with functional options
func NewFireflyClient(opts ...ClientOption) (*client.ClientWithResponses, error) {
	// Start with default options
	options := defaultClientOptions()
	
	// Apply all provided options
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}
	
	// Validate required fields
	if options.baseURL == "" {
		return nil, ErrEmptyBaseURL
	}
	if options.apiToken == "" {
		return nil, ErrEmptyAPIToken
	}
	
	// Create HTTP client if not provided
	httpClient := options.httpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: options.timeout,
		}
	}
	
	// Build request editors
	requestEditors := options.requestEditors
	
	// Add authorization header
	requestEditors = append(requestEditors, func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+options.apiToken)
		return nil
	})
	
	// Add user agent
	requestEditors = append(requestEditors, func(ctx context.Context, req *http.Request) error {
		req.Header.Set("User-Agent", options.userAgent)
		return nil
	})
	
	// Create client options
	clientOpts := []client.ClientOption{
		client.WithHTTPClient(httpClient),
	}
	for _, editor := range requestEditors {
		clientOpts = append(clientOpts, client.WithRequestEditorFn(editor))
	}
	
	// Create and return the client
	return client.NewClientWithResponses(options.baseURL, clientOpts...)
}