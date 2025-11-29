package fireflyMCP

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestConfig holds test configuration
type TestConfig struct {
	ServerURL string
	APIToken  string
	Timeout   time.Duration
}

// loadTestConfig loads test configuration from environment or config file
func loadTestConfig(t *testing.T) *TestConfig {
	// Try to load from environment first
	// Support both new (FIREFLY_MCP_*) and legacy (FIREFLY_TEST_*) env vars
	serverURL := os.Getenv("FIREFLY_MCP_SERVER_URL")
	if serverURL == "" {
		serverURL = os.Getenv("FIREFLY_TEST_URL") // Legacy support
	}

	apiToken := os.Getenv("FIREFLY_MCP_API_TOKEN")
	if apiToken == "" {
		apiToken = os.Getenv("FIREFLY_TEST_TOKEN") // Legacy support
	}

	if serverURL == "" || apiToken == "" {
		// Fallback to config file
		config, err := LoadConfig("../../config.yaml")
		if err != nil {
			t.Skipf("Skipping integration tests: no test config available (%v)", err)
		}
		serverURL = config.Server.URL
		apiToken = config.API.Token
	}

	if serverURL == "" || apiToken == "" {
		t.Skip("Skipping integration tests: Set FIREFLY_MCP_SERVER_URL and FIREFLY_MCP_API_TOKEN environment variables (or legacy FIREFLY_TEST_URL and FIREFLY_TEST_TOKEN)")
	}

	return &TestConfig{
		ServerURL: serverURL,
		APIToken:  apiToken,
		Timeout:   30 * time.Second,
	}
}

// createTestServer creates a test MCP server instance
func createTestServer(t *testing.T, testConfig *TestConfig) *FireflyMCPServer {
	config := &Config{
		Server: struct {
			URL string `yaml:"url" mapstructure:"url"`
		}{URL: testConfig.ServerURL},
		API: struct {
			Token string `yaml:"token" mapstructure:"token"`
		}{Token: testConfig.APIToken},
		Client: struct {
			Timeout int `yaml:"timeout" mapstructure:"timeout"`
		}{Timeout: int(testConfig.Timeout.Seconds())},
		Limits: struct {
			Accounts     int `yaml:"accounts" mapstructure:"accounts"`
			Transactions int `yaml:"transactions" mapstructure:"transactions"`
			Categories   int `yaml:"categories" mapstructure:"categories"`
			Budgets      int `yaml:"budgets" mapstructure:"budgets"`
		}{
			Accounts:     10,
			Transactions: 5,
			Categories:   10,
			Budgets:      10,
		},
		MCP: struct {
			Name         string `yaml:"name" mapstructure:"name"`
			Version      string `yaml:"version" mapstructure:"version"`
			Instructions string `yaml:"instructions" mapstructure:"instructions"`
		}{
			Name:         "firefly-iii-mcp-test",
			Version:      "1.0.0-test",
			Instructions: "Test MCP server for Firefly III",
		},
	}

	server, err := NewFireflyMCPServer(config)
	require.NoError(t, err, "Failed to create test server")
	return server
}

// mockTransport implements mcp.Transport for testing
type mockTransport struct {
	requests  []interface{}
	responses []interface{}
}

func (m *mockTransport) Start(ctx context.Context) error {
	return nil
}

func (m *mockTransport) Close() error {
	return nil
}

func (m *mockTransport) Send(ctx context.Context, message interface{}) error {
	m.requests = append(m.requests, message)
	return nil
}

func (m *mockTransport) Receive(ctx context.Context) (interface{}, error) {
	if len(m.responses) == 0 {
		return nil, fmt.Errorf("no more responses")
	}
	response := m.responses[0]
	m.responses = m.responses[1:]
	return response, nil
}
