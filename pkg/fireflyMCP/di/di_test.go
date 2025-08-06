package di

import (
	"testing"

	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP"
)

func TestInitializeServer(t *testing.T) {
	// Create test configuration
	config := &fireflyMCP.Config{}
	config.Server.URL = "https://test.firefly.example.com"
	config.API.Token = "test-token"
	config.Client.Timeout = 30
	config.MCP.Name = "test-server"
	config.MCP.Version = "1.0.0"
	config.MCP.Instructions = "Test instructions"
	config.Limits.Accounts = 100
	config.Limits.Transactions = 100
	config.Limits.Categories = 100
	config.Limits.Budgets = 100

	// Initialize server using DI
	server, err := InitializeServer(config)
	if err != nil {
		t.Fatalf("Failed to initialize server: %v", err)
	}

	// Verify server is properly initialized
	if server == nil {
		t.Fatal("Server should not be nil")
	}

	if server.Server == nil {
		t.Error("MCP Server should not be nil")
	}

	if server.Client == nil {
		t.Error("Firefly client should not be nil")
	}

	if server.Config == nil {
		t.Error("Config should not be nil")
	}

	if server.Handlers == nil {
		t.Error("Handlers should not be nil")
	}

	if server.Chain == nil {
		t.Error("Middleware chain should not be nil")
	}

	// Verify configuration is properly set
	if server.Config.Server.URL != config.Server.URL {
		t.Errorf("Server URL mismatch: got %s, want %s", 
			server.Config.Server.URL, config.Server.URL)
	}

	if server.Config.API.Token != config.API.Token {
		t.Errorf("API token mismatch: got %s, want %s",
			server.Config.API.Token, config.API.Token)
	}

	// Verify handlers are properly initialized
	if server.Handlers.GetAccountHandlers() == nil {
		t.Error("Account handlers should not be nil")
	}

	if server.Handlers.GetTransactionHandlers() == nil {
		t.Error("Transaction handlers should not be nil")
	}

	if server.Handlers.GetBudgetHandlers() == nil {
		t.Error("Budget handlers should not be nil")
	}

	if server.Handlers.GetCategoryHandlers() == nil {
		t.Error("Category handlers should not be nil")
	}

	if server.Handlers.GetTagHandlers() == nil {
		t.Error("Tag handlers should not be nil")
	}

	if server.Handlers.GetInsightHandlers() == nil {
		t.Error("Insight handlers should not be nil")
	}

	if server.Handlers.GetBillHandlers() == nil {
		t.Error("Bill handlers should not be nil")
	}

	if server.Handlers.GetRecurrenceHandlers() == nil {
		t.Error("Recurrence handlers should not be nil")
	}
}

func TestProviders(t *testing.T) {
	config := &fireflyMCP.Config{}
	config.Server.URL = "https://test.firefly.example.com"
	config.API.Token = "test-token"
	config.Client.Timeout = 30
	config.MCP.Name = "test-server"
	config.MCP.Version = "1.0.0"

	t.Run("ProvideHTTPClient", func(t *testing.T) {
		client := ProvideHTTPClient(config)
		if client == nil {
			t.Fatal("HTTP client should not be nil")
		}
		if client.Timeout != 30000000000 { // 30 seconds in nanoseconds
			t.Errorf("Unexpected timeout: %v", client.Timeout)
		}
	})

	t.Run("ProvideMCPServer", func(t *testing.T) {
		server := ProvideMCPServer(config)
		if server == nil {
			t.Fatal("MCP server should not be nil")
		}
	})

	t.Run("ProvideMiddlewareChain", func(t *testing.T) {
		chain := ProvideMiddlewareChain(config)
		if chain == nil {
			t.Fatal("Middleware chain should not be nil")
		}
	})

	t.Run("ProvideFireflyClient", func(t *testing.T) {
		httpClient := ProvideHTTPClient(config)
		client, err := ProvideFireflyClient(config, httpClient)
		if err != nil {
			t.Fatalf("Failed to create Firefly client: %v", err)
		}
		if client == nil {
			t.Fatal("Firefly client should not be nil")
		}
	})
}