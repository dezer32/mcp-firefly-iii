package fireflyMCP

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the MCP server configuration
type Config struct {
	Server struct {
		URL string `yaml:"url" mapstructure:"url"`
	} `yaml:"server" mapstructure:"server"`
	API struct {
		Token string `yaml:"token" mapstructure:"token"`
	} `yaml:"api" mapstructure:"api"`
	Client struct {
		Timeout int `yaml:"timeout" mapstructure:"timeout"`
	} `yaml:"client" mapstructure:"client"`
	Limits struct {
		Accounts     int `yaml:"accounts" mapstructure:"accounts"`
		Transactions int `yaml:"transactions" mapstructure:"transactions"`
		Categories   int `yaml:"categories" mapstructure:"categories"`
		Budgets      int `yaml:"budgets" mapstructure:"budgets"`
	} `yaml:"limits" mapstructure:"limits"`
	MCP struct {
		Name         string `yaml:"name" mapstructure:"name"`
		Version      string `yaml:"version" mapstructure:"version"`
		Instructions string `yaml:"instructions" mapstructure:"instructions"`
	} `yaml:"mcp" mapstructure:"mcp"`
	HTTP struct {
		Enabled        bool     `yaml:"enabled" mapstructure:"enabled"`
		Port           int      `yaml:"port" mapstructure:"port"`
		Host           string   `yaml:"host" mapstructure:"host"`
		ReadTimeout    int      `yaml:"read_timeout" mapstructure:"read_timeout"`
		WriteTimeout   int      `yaml:"write_timeout" mapstructure:"write_timeout"`
		IdleTimeout    int      `yaml:"idle_timeout" mapstructure:"idle_timeout"`
		SessionTimeout int      `yaml:"session_timeout" mapstructure:"session_timeout"`
		AllowedOrigins []string `yaml:"allowed_origins" mapstructure:"allowed_origins"`
		RateLimit      float64  `yaml:"rate_limit" mapstructure:"rate_limit"`
		RateBurst      int      `yaml:"rate_burst" mapstructure:"rate_burst"`
	} `yaml:"http" mapstructure:"http"`
}

// LoadConfig loads configuration from YAML file and environment variables
// Environment variables take precedence over YAML configuration
// Environment variables use the prefix FIREFLY_MCP_ and follow the pattern:
//
//	FIREFLY_MCP_SERVER_URL, FIREFLY_MCP_API_TOKEN, etc.
func LoadConfig(filename string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Configure environment variable support
	v.SetEnvPrefix("FIREFLY_MCP")
	v.AutomaticEnv()
	// Replace dots with underscores for nested config (e.g., server.url -> SERVER_URL)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind environment variables for nested fields
	bindEnvVars(v)

	// Try to read config file if it exists
	if filename != "" {
		// Check if file exists
		if _, err := os.Stat(filename); err == nil {
			v.SetConfigFile(filename)
			if err := v.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		} else if !os.IsNotExist(err) {
			// File exists but can't stat it
			return nil, fmt.Errorf("failed to access config file: %w", err)
		}
		// If file doesn't exist, continue with defaults and env vars only
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// bindEnvVars explicitly binds environment variables to config keys
// This is needed because Viper's AutomaticEnv doesn't automatically bind nested struct fields
func bindEnvVars(v *viper.Viper) {
	// Server config
	v.BindEnv("server.url")

	// API config
	v.BindEnv("api.token")

	// Client config
	v.BindEnv("client.timeout")

	// Limits config
	v.BindEnv("limits.accounts")
	v.BindEnv("limits.transactions")
	v.BindEnv("limits.categories")
	v.BindEnv("limits.budgets")

	// MCP config
	v.BindEnv("mcp.name")
	v.BindEnv("mcp.version")
	v.BindEnv("mcp.instructions")

	// HTTP config
	v.BindEnv("http.enabled")
	v.BindEnv("http.port")
	v.BindEnv("http.host")
	v.BindEnv("http.read_timeout")
	v.BindEnv("http.write_timeout")
	v.BindEnv("http.idle_timeout")
	v.BindEnv("http.session_timeout")
	v.BindEnv("http.allowed_origins")
	v.BindEnv("http.rate_limit")
	v.BindEnv("http.rate_burst")
}

// setDefaults configures default values for all configuration options
func setDefaults(v *viper.Viper) {
	// Client defaults
	v.SetDefault("client.timeout", 30)

	// Limits defaults
	v.SetDefault("limits.accounts", 100)
	v.SetDefault("limits.transactions", 100)
	v.SetDefault("limits.categories", 100)
	v.SetDefault("limits.budgets", 100)

	// MCP defaults
	v.SetDefault("mcp.name", "firefly-iii-mcp")
	v.SetDefault("mcp.version", "1.0.0")
	v.SetDefault("mcp.instructions", "MCP server for Firefly III personal finance management")

	// HTTP defaults
	v.SetDefault("http.enabled", false)
	v.SetDefault("http.port", 8080)
	v.SetDefault("http.host", "0.0.0.0")
	v.SetDefault("http.read_timeout", 30)
	v.SetDefault("http.write_timeout", 30)
	v.SetDefault("http.idle_timeout", 120)
	v.SetDefault("http.session_timeout", 300)
	v.SetDefault("http.allowed_origins", []string{"*"})
	v.SetDefault("http.rate_limit", 10.0)
	v.SetDefault("http.rate_burst", 20)
}

// validateConfig validates that required configuration fields are set
func validateConfig(config *Config) error {
	if config.Server.URL == "" {
		return fmt.Errorf("server.url is required (set via config file or FIREFLY_MCP_SERVER_URL)")
	}
	// api.token is required only for stdio mode (HTTP mode gets token from request header)
	if !config.HTTP.Enabled && config.API.Token == "" {
		return fmt.Errorf("api.token is required for stdio mode (set via config file or FIREFLY_MCP_API_TOKEN)")
	}
	if config.Client.Timeout <= 0 {
		return fmt.Errorf("client.timeout must be positive")
	}
	if config.Limits.Accounts <= 0 {
		return fmt.Errorf("limits.accounts must be positive")
	}
	if config.Limits.Transactions <= 0 {
		return fmt.Errorf("limits.transactions must be positive")
	}
	if config.Limits.Categories <= 0 {
		return fmt.Errorf("limits.categories must be positive")
	}
	if config.Limits.Budgets <= 0 {
		return fmt.Errorf("limits.budgets must be positive")
	}
	return nil
}

// GetTimeout returns the client timeout as a duration
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.Client.Timeout) * time.Second
}
