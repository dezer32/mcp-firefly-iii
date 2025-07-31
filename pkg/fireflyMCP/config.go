package fireflyMCP

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the MCP server configuration
type Config struct {
	Server struct {
		URL string `yaml:"url"`
	} `yaml:"server"`
	API struct {
		Token string `yaml:"token"`
	} `yaml:"api"`
	Client struct {
		Timeout int `yaml:"timeout"`
	} `yaml:"client"`
	Limits struct {
		Accounts     int `yaml:"accounts"`
		Transactions int `yaml:"transactions"`
		Categories   int `yaml:"categories"`
		Budgets      int `yaml:"budgets"`
	} `yaml:"limits"`
	MCP struct {
		Name         string `yaml:"name"`
		Version      string `yaml:"version"`
		Instructions string `yaml:"instructions"`
	} `yaml:"mcp"`
}

// LoadConfig loads configuration from YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if config.Client.Timeout == 0 {
		config.Client.Timeout = 30
	}
	if config.Limits.Accounts == 0 {
		config.Limits.Accounts = 100
	}
	if config.Limits.Transactions == 0 {
		config.Limits.Transactions = 100
	}
	if config.Limits.Categories == 0 {
		config.Limits.Categories = 100
	}
	if config.Limits.Budgets == 0 {
		config.Limits.Budgets = 100
	}
	if config.MCP.Name == "" {
		config.MCP.Name = "firefly-iii-mcp"
	}
	if config.MCP.Version == "" {
		config.MCP.Version = "1.0.0"
	}
	if config.MCP.Instructions == "" {
		config.MCP.Instructions = "MCP server for Firefly III personal finance management"
	}

	return &config, nil
}

// GetTimeout returns the client timeout as a duration
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.Client.Timeout) * time.Second
}
