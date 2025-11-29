package fireflyMCP

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromYAML(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  url: https://test.firefly.com/api
api:
  token: test-token-123
client:
  timeout: 60
limits:
  accounts: 200
  transactions: 300
  categories: 150
  budgets: 50
mcp:
  name: test-mcp
  version: 2.0.0
  instructions: Test instructions
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configFile)
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "https://test.firefly.com/api", config.Server.URL)
	assert.Equal(t, "test-token-123", config.API.Token)
	assert.Equal(t, 60, config.Client.Timeout)
	assert.Equal(t, 200, config.Limits.Accounts)
	assert.Equal(t, 300, config.Limits.Transactions)
	assert.Equal(t, 150, config.Limits.Categories)
	assert.Equal(t, 50, config.Limits.Budgets)
	assert.Equal(t, "test-mcp", config.MCP.Name)
	assert.Equal(t, "2.0.0", config.MCP.Version)
	assert.Equal(t, "Test instructions", config.MCP.Instructions)
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// Create a minimal config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  url: https://test.firefly.com/api
api:
  token: test-token-123
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configFile)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Check defaults are applied
	assert.Equal(t, 30, config.Client.Timeout)
	assert.Equal(t, 100, config.Limits.Accounts)
	assert.Equal(t, 100, config.Limits.Transactions)
	assert.Equal(t, 100, config.Limits.Categories)
	assert.Equal(t, 100, config.Limits.Budgets)
	assert.Equal(t, "firefly-iii-mcp", config.MCP.Name)
	assert.Equal(t, "1.0.0", config.MCP.Version)
	assert.Equal(t, "MCP server for Firefly III personal finance management", config.MCP.Instructions)
}

func TestLoadConfigFromEnvVars(t *testing.T) {
	// Set environment variables
	envVars := map[string]string{
		"FIREFLY_MCP_SERVER_URL":          "https://env.firefly.com/api",
		"FIREFLY_MCP_API_TOKEN":           "env-token-456",
		"FIREFLY_MCP_CLIENT_TIMEOUT":      "90",
		"FIREFLY_MCP_LIMITS_ACCOUNTS":     "250",
		"FIREFLY_MCP_LIMITS_TRANSACTIONS": "350",
		"FIREFLY_MCP_LIMITS_CATEGORIES":   "175",
		"FIREFLY_MCP_LIMITS_BUDGETS":      "75",
		"FIREFLY_MCP_MCP_NAME":            "env-mcp",
		"FIREFLY_MCP_MCP_VERSION":         "3.0.0",
		"FIREFLY_MCP_MCP_INSTRUCTIONS":    "Env instructions",
	}

	// Set env vars and clean up after test
	for key, value := range envVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	// Load config without a file (using only env vars)
	config, err := LoadConfig("")
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "https://env.firefly.com/api", config.Server.URL)
	assert.Equal(t, "env-token-456", config.API.Token)
	assert.Equal(t, 90, config.Client.Timeout)
	assert.Equal(t, 250, config.Limits.Accounts)
	assert.Equal(t, 350, config.Limits.Transactions)
	assert.Equal(t, 175, config.Limits.Categories)
	assert.Equal(t, 75, config.Limits.Budgets)
	assert.Equal(t, "env-mcp", config.MCP.Name)
	assert.Equal(t, "3.0.0", config.MCP.Version)
	assert.Equal(t, "Env instructions", config.MCP.Instructions)
}

func TestLoadConfigEnvOverridesYAML(t *testing.T) {
	// Create a config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  url: https://yaml.firefly.com/api
api:
  token: yaml-token
client:
  timeout: 30
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set environment variables to override
	os.Setenv("FIREFLY_MCP_SERVER_URL", "https://override.firefly.com/api")
	os.Setenv("FIREFLY_MCP_API_TOKEN", "override-token")
	os.Setenv("FIREFLY_MCP_CLIENT_TIMEOUT", "120")
	defer func() {
		os.Unsetenv("FIREFLY_MCP_SERVER_URL")
		os.Unsetenv("FIREFLY_MCP_API_TOKEN")
		os.Unsetenv("FIREFLY_MCP_CLIENT_TIMEOUT")
	}()

	config, err := LoadConfig(configFile)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Env vars should override YAML values
	assert.Equal(t, "https://override.firefly.com/api", config.Server.URL)
	assert.Equal(t, "override-token", config.API.Token)
	assert.Equal(t, 120, config.Client.Timeout)
}

func TestLoadConfigMissingRequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		configYAML  string
		errorString string
	}{
		{
			name: "missing server URL",
			configYAML: `
api:
  token: test-token
`,
			errorString: "server.url is required",
		},
		{
			name: "missing API token",
			configYAML: `
server:
  url: https://test.firefly.com/api
`,
			errorString: "api.token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configFile := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(configFile, []byte(tt.configYAML), 0644)
			require.NoError(t, err)

			config, err := LoadConfig(configFile)
			assert.Error(t, err)
			assert.Nil(t, config)
			assert.Contains(t, err.Error(), tt.errorString)
		})
	}
}

func TestLoadConfigInvalidValues(t *testing.T) {
	tests := []struct {
		name        string
		configYAML  string
		errorString string
	}{
		{
			name: "negative timeout",
			configYAML: `
server:
  url: https://test.firefly.com/api
api:
  token: test-token
client:
  timeout: -10
`,
			errorString: "client.timeout must be positive",
		},
		{
			name: "negative accounts limit",
			configYAML: `
server:
  url: https://test.firefly.com/api
api:
  token: test-token
limits:
  accounts: -5
`,
			errorString: "limits.accounts must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configFile := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(configFile, []byte(tt.configYAML), 0644)
			require.NoError(t, err)

			config, err := LoadConfig(configFile)
			assert.Error(t, err)
			assert.Nil(t, config)
			assert.Contains(t, err.Error(), tt.errorString)
		})
	}
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	// Try to load non-existent file without env vars (should fail validation)
	config, err := LoadConfig("/non/existent/config.yaml")
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "server.url is required")
}

func TestLoadConfigWithEnvVarsOnly(t *testing.T) {
	// Set minimal required env vars
	os.Setenv("FIREFLY_MCP_SERVER_URL", "https://env-only.firefly.com/api")
	os.Setenv("FIREFLY_MCP_API_TOKEN", "env-only-token")
	defer func() {
		os.Unsetenv("FIREFLY_MCP_SERVER_URL")
		os.Unsetenv("FIREFLY_MCP_API_TOKEN")
	}()

	// Load without config file
	config, err := LoadConfig("")
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "https://env-only.firefly.com/api", config.Server.URL)
	assert.Equal(t, "env-only-token", config.API.Token)
	// Defaults should be applied
	assert.Equal(t, 30, config.Client.Timeout)
	assert.Equal(t, 100, config.Limits.Accounts)
}

func TestGetTimeout(t *testing.T) {
	config := &Config{}
	config.Client.Timeout = 45

	timeout := config.GetTimeout()
	assert.Equal(t, "45s", timeout.String())
}

func TestLoadConfigPartialEnvOverride(t *testing.T) {
	// Create a config file with some values
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  url: https://yaml.firefly.com/api
api:
  token: yaml-token
client:
  timeout: 30
limits:
  accounts: 50
  transactions: 75
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Override only some values with env vars
	os.Setenv("FIREFLY_MCP_API_TOKEN", "env-override-token")
	os.Setenv("FIREFLY_MCP_LIMITS_ACCOUNTS", "150")
	defer func() {
		os.Unsetenv("FIREFLY_MCP_API_TOKEN")
		os.Unsetenv("FIREFLY_MCP_LIMITS_ACCOUNTS")
	}()

	config, err := LoadConfig(configFile)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Verify mixed values
	assert.Equal(t, "https://yaml.firefly.com/api", config.Server.URL) // From YAML
	assert.Equal(t, "env-override-token", config.API.Token)            // From env
	assert.Equal(t, 30, config.Client.Timeout)                         // From YAML
	assert.Equal(t, 150, config.Limits.Accounts)                       // From env
	assert.Equal(t, 75, config.Limits.Transactions)                    // From YAML
	assert.Equal(t, 100, config.Limits.Categories)                     // Default
	assert.Equal(t, 100, config.Limits.Budgets)                        // Default
}
