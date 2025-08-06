package fireflyMCP

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerOptions(t *testing.T) {
	t.Run("WithHTTPClient", func(t *testing.T) {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		opt := WithHTTPClient(httpClient)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, httpClient, opts.httpClient)
	})
	
	t.Run("WithHTTPClient_Nil", func(t *testing.T) {
		opt := WithHTTPClient(nil)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrNilHTTPClient, err)
	})
	
	t.Run("WithTimeout", func(t *testing.T) {
		timeout := 45 * time.Second
		opt := WithTimeout(timeout)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, timeout, opts.timeout)
	})
	
	t.Run("WithTimeout_Invalid", func(t *testing.T) {
		opt := WithTimeout(0)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrInvalidTimeout, err)
	})
	
	t.Run("WithAPIToken", func(t *testing.T) {
		token := "test-token-123"
		opt := WithAPIToken(token)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, token, opts.apiToken)
	})
	
	t.Run("WithAPIToken_Empty", func(t *testing.T) {
		opt := WithAPIToken("")
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrEmptyAPIToken, err)
	})
	
	t.Run("WithBaseURL", func(t *testing.T) {
		url := "https://firefly.example.com"
		opt := WithBaseURL(url)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, url, opts.baseURL)
	})
	
	t.Run("WithBaseURL_Empty", func(t *testing.T) {
		opt := WithBaseURL("")
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrEmptyBaseURL, err)
	})
	
	t.Run("WithMCPInfo", func(t *testing.T) {
		name := "test-mcp"
		version := "2.0.0"
		opt := WithMCPInfo(name, version)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, name, opts.mcpName)
		assert.Equal(t, version, opts.mcpVersion)
	})
	
	t.Run("WithLogging", func(t *testing.T) {
		opt := WithLogging(true, middleware.LogLevelDebug)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.True(t, opts.enableLogging)
		assert.Equal(t, middleware.LogLevelDebug, opts.logLevel)
	})
	
	t.Run("WithMetrics", func(t *testing.T) {
		opt := WithMetrics(false)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.False(t, opts.enableMetrics)
	})
	
	t.Run("WithRecovery", func(t *testing.T) {
		opt := WithRecovery(false)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.False(t, opts.enableRecovery)
	})
	
	t.Run("WithTracing", func(t *testing.T) {
		opt := WithTracing(true)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.True(t, opts.enableTracing)
	})
	
	t.Run("WithRateLimit", func(t *testing.T) {
		opt := WithRateLimit(200, 20)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, 200, opts.rateLimit)
		assert.Equal(t, 20, opts.rateLimitBurst)
	})
	
	t.Run("WithRateLimit_Invalid", func(t *testing.T) {
		opt := WithRateLimit(0, 10)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrInvalidRateLimit, err)
	})
	
	t.Run("WithCache", func(t *testing.T) {
		ttl := 10 * time.Minute
		opt := WithCache(true, ttl)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.True(t, opts.cacheEnabled)
		assert.Equal(t, ttl, opts.cacheTTL)
	})
	
	t.Run("WithConnectionPool", func(t *testing.T) {
		opt := WithConnectionPool(200, 20)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, 200, opts.maxIdleConns)
		assert.Equal(t, 20, opts.maxConnsPerHost)
	})
	
	t.Run("WithConfig", func(t *testing.T) {
		config := &Config{}
		config.Server.URL = "https://test.com"
		config.API.Token = "test-token"
		config.Client.Timeout = 60
		config.MCP.Name = "test"
		config.MCP.Version = "1.0"
		opt := WithConfig(config)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, config, opts.configOverride)
		assert.Equal(t, "test-token", opts.apiToken)
		assert.Equal(t, "https://test.com", opts.baseURL)
		assert.Equal(t, 60*time.Second, opts.timeout)
	})
	
	t.Run("WithConfig_Nil", func(t *testing.T) {
		opt := WithConfig(nil)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrNilConfig, err)
	})
	
	t.Run("WithMiddleware", func(t *testing.T) {
		mw := middleware.NewLoggingMiddleware(nil, middleware.LogLevelInfo)
		opt := WithMiddleware(mw)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Len(t, opts.middlewares, 1)
	})
	
	t.Run("WithMiddleware_Nil", func(t *testing.T) {
		opt := WithMiddleware(nil)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrNilMiddleware, err)
	})
	
	t.Run("WithRequestEditor", func(t *testing.T) {
		editor := func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-Custom", "value")
			return nil
		}
		opt := WithRequestEditor(editor)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Len(t, opts.requestEditors, 1)
	})
	
	t.Run("WithRequestEditor_Nil", func(t *testing.T) {
		opt := WithRequestEditor(nil)
		
		opts := defaultServerOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrNilRequestEditor, err)
	})
}

func TestClientOptions(t *testing.T) {
	t.Run("WithClientHTTPClient", func(t *testing.T) {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		opt := WithClientHTTPClient(httpClient)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, httpClient, opts.httpClient)
	})
	
	t.Run("WithClientHTTPClient_Nil", func(t *testing.T) {
		opt := WithClientHTTPClient(nil)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrNilHTTPClient, err)
	})
	
	t.Run("WithClientTimeout", func(t *testing.T) {
		timeout := 45 * time.Second
		opt := WithClientTimeout(timeout)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, timeout, opts.timeout)
	})
	
	t.Run("WithClientTimeout_Invalid", func(t *testing.T) {
		opt := WithClientTimeout(0)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrInvalidTimeout, err)
	})
	
	t.Run("WithClientBaseURL", func(t *testing.T) {
		url := "https://firefly.example.com"
		opt := WithClientBaseURL(url)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, url, opts.baseURL)
	})
	
	t.Run("WithClientBaseURL_Empty", func(t *testing.T) {
		opt := WithClientBaseURL("")
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrEmptyBaseURL, err)
	})
	
	t.Run("WithClientAPIToken", func(t *testing.T) {
		token := "test-token-123"
		opt := WithClientAPIToken(token)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, token, opts.apiToken)
	})
	
	t.Run("WithClientAPIToken_Empty", func(t *testing.T) {
		opt := WithClientAPIToken("")
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrEmptyAPIToken, err)
	})
	
	t.Run("WithClientRetry", func(t *testing.T) {
		opt := WithClientRetry(5, 2*time.Second)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, 5, opts.retryCount)
		assert.Equal(t, 2*time.Second, opts.retryWait)
	})
	
	t.Run("WithClientRetry_Invalid", func(t *testing.T) {
		opt := WithClientRetry(-1, time.Second)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrInvalidRetryCount, err)
	})
	
	t.Run("WithClientUserAgent", func(t *testing.T) {
		userAgent := "custom-agent/1.0"
		opt := WithClientUserAgent(userAgent)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Equal(t, userAgent, opts.userAgent)
	})
	
	t.Run("WithClientRequestEditor", func(t *testing.T) {
		editor := func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-Custom", "value")
			return nil
		}
		opt := WithClientRequestEditor(editor)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.NoError(t, err)
		assert.Len(t, opts.requestEditors, 1)
	})
	
	t.Run("WithClientRequestEditor_Nil", func(t *testing.T) {
		opt := WithClientRequestEditor(nil)
		
		opts := defaultClientOptions()
		err := opt(opts)
		
		assert.Equal(t, ErrNilRequestEditor, err)
	})
}

func TestNewFireflyClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		client, err := NewFireflyClient(
			WithClientBaseURL("https://test.com"),
			WithClientAPIToken("test-token"),
		)
		
		require.NoError(t, err)
		require.NotNil(t, client)
	})
	
	t.Run("MissingBaseURL", func(t *testing.T) {
		client, err := NewFireflyClient(
			WithClientAPIToken("test-token"),
		)
		
		assert.Equal(t, ErrEmptyBaseURL, err)
		assert.Nil(t, client)
	})
	
	t.Run("MissingAPIToken", func(t *testing.T) {
		client, err := NewFireflyClient(
			WithClientBaseURL("https://test.com"),
		)
		
		assert.Equal(t, ErrEmptyAPIToken, err)
		assert.Nil(t, client)
	})
	
	t.Run("WithCustomHTTPClient", func(t *testing.T) {
		httpClient := &http.Client{
			Timeout: 60 * time.Second,
		}
		
		client, err := NewFireflyClient(
			WithClientBaseURL("https://test.com"),
			WithClientAPIToken("test-token"),
			WithClientHTTPClient(httpClient),
		)
		
		require.NoError(t, err)
		require.NotNil(t, client)
	})
	
	t.Run("WithAllOptions", func(t *testing.T) {
		client, err := NewFireflyClient(
			WithClientBaseURL("https://test.com"),
			WithClientAPIToken("test-token"),
			WithClientTimeout(45*time.Second),
			WithClientRetry(3, 2*time.Second),
			WithClientUserAgent("test-agent/1.0"),
			WithClientRequestEditor(func(ctx context.Context, req *http.Request) error {
				req.Header.Set("X-Custom", "value")
				return nil
			}),
		)
		
		require.NoError(t, err)
		require.NotNil(t, client)
	})
}

func TestBuildConfig(t *testing.T) {
	t.Run("WithConfigOverride", func(t *testing.T) {
		config := &Config{}
		config.Server.URL = "https://test.com"
		config.API.Token = "test-token"
		config.Client.Timeout = 60
		config.MCP.Name = "test"
		config.MCP.Version = "1.0"
		
		opts := &serverOptions{
			configOverride: config,
		}
		
		result := opts.buildConfig()
		assert.Equal(t, config, result)
	})
	
	t.Run("WithoutConfigOverride", func(t *testing.T) {
		opts := &serverOptions{
			baseURL:    "https://test.com",
			apiToken:   "test-token",
			timeout:    45 * time.Second,
			mcpName:    "test-mcp",
			mcpVersion: "2.0.0",
		}
		
		result := opts.buildConfig()
		assert.Equal(t, "https://test.com", result.Server.URL)
		assert.Equal(t, "test-token", result.API.Token)
		assert.Equal(t, 45, result.Client.Timeout)
		assert.Equal(t, "test-mcp", result.MCP.Name)
		assert.Equal(t, "2.0.0", result.MCP.Version)
	})
}

func TestBuildMiddlewareChain(t *testing.T) {
	t.Run("DefaultMiddlewares", func(t *testing.T) {
		opts := &serverOptions{
			enableRecovery: true,
			enableLogging:  true,
			enableMetrics:  true,
			logLevel:       middleware.LogLevelInfo,
		}
		
		chain := opts.buildMiddlewareChain()
		assert.NotNil(t, chain)
	})
	
	t.Run("DisabledMiddlewares", func(t *testing.T) {
		opts := &serverOptions{
			enableRecovery: false,
			enableLogging:  false,
			enableMetrics:  false,
		}
		
		chain := opts.buildMiddlewareChain()
		assert.NotNil(t, chain)
	})
	
	t.Run("WithCustomMiddlewares", func(t *testing.T) {
		customMw := middleware.NewTimingMiddleware(nil, 2*time.Second)
		
		opts := &serverOptions{
			enableRecovery: true,
			enableLogging:  true,
			enableMetrics:  false,
			logLevel:       middleware.LogLevelDebug,
			middlewares:    []middleware.Middleware{customMw},
		}
		
		chain := opts.buildMiddlewareChain()
		assert.NotNil(t, chain)
	})
}

func TestBackwardCompatibility(t *testing.T) {
	t.Run("NewServer_WithValidConfig", func(t *testing.T) {
		config := &Config{}
		config.Server.URL = "https://test.com"
		config.API.Token = "test-token"
		config.Client.Timeout = 30
		config.MCP.Name = "test"
		config.MCP.Version = "1.0"
		config.Limits.Accounts = 100
		config.Limits.Transactions = 100
		config.Limits.Categories = 100
		config.Limits.Budgets = 100
		
		server, err := NewServer(config)
		
		require.NoError(t, err)
		require.NotNil(t, server)
		assert.Equal(t, config, server.Config)
		assert.NotNil(t, server.Client)
		assert.NotNil(t, server.Server)
		assert.NotNil(t, server.Handlers)
	})
}

// Integration test example
func TestNewServerWithOptions_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	t.Run("CompleteServerSetup", func(t *testing.T) {
		server, err := NewServerWithOptions(
			WithBaseURL("https://firefly.example.com"),
			WithAPIToken("test-token"),
			WithTimeout(45*time.Second),
			WithMCPInfo("custom-mcp", "2.0.0"),
			WithLogging(true, middleware.LogLevelDebug),
			WithMetrics(true),
			WithRecovery(true),
			WithRateLimit(150, 15),
			WithCache(true, 10*time.Minute),
			WithConnectionPool(200, 20),
		)
		
		require.NoError(t, err)
		require.NotNil(t, server)
		assert.NotNil(t, server.Chain)
		assert.Equal(t, "custom-mcp", server.Config.MCP.Name)
		assert.Equal(t, "2.0.0", server.Config.MCP.Version)
	})
}