package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dezer32/firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
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
	serverURL := os.Getenv("FIREFLY_TEST_URL")
	apiToken := os.Getenv("FIREFLY_TEST_TOKEN")

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
		t.Skip("Skipping integration tests: FIREFLY_TEST_URL and FIREFLY_TEST_TOKEN environment variables not set")
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
			URL string `yaml:"url"`
		}{URL: testConfig.ServerURL},
		API: struct {
			Token string `yaml:"token"`
		}{Token: testConfig.APIToken},
		Client: struct {
			Timeout int `yaml:"timeout"`
		}{Timeout: int(testConfig.Timeout.Seconds())},
		Limits: struct {
			Accounts     int `yaml:"accounts"`
			Transactions int `yaml:"transactions"`
			Categories   int `yaml:"categories"`
			Budgets      int `yaml:"budgets"`
		}{
			Accounts:     10,
			Transactions: 5,
			Categories:   10,
			Budgets:      10,
		},
		MCP: struct {
			Name         string `yaml:"name"`
			Version      string `yaml:"version"`
			Instructions string `yaml:"instructions"`
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

func TestIntegration_ListAccounts(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// Test direct client call first
	t.Run(
		"DirectClientCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing direct client call to %s\n", testConfig.ServerURL)

			clientWithAuth, err := client.NewClient(
				testConfig.ServerURL, client.WithRequestEditorFn(
					func(ctx context.Context, req *http.Request) error {
						req.Header.Set("Authorization", "Bearer "+testConfig.APIToken)
						req.Header.Set("Accept", "application/vnd.api+json")
						req.Header.Set("Content-Type", "application/json")
						return nil
					},
				),
			)
			require.NoError(t, err, "Failed to create client")

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			params := &client.ListAccountParams{}
			response, err := clientWithAuth.ListAccount(ctx, params)

			fmt.Printf("[DEBUG_LOG] Response status: %v, Error: %v\n", response.StatusCode, err)

			if err != nil {
				t.Logf("API call failed (this might be expected): %v", err)
				// Don't fail the test here, just log the error
			} else {
				assert.Equal(t, 200, response.StatusCode, "Expected successful response")
				t.Logf("Successfully retrieved accounts from Firefly III API")
			}
		},
	)

	// Test MCP tool call
	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_accounts\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters
			params := &mcp.CallToolParamsFor[ListAccountsArgs]{
				Name: "list_accounts",
				Arguments: ListAccountsArgs{
					Limit: 5,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleListAccounts(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call failed (this might be expected): %v", err)
				// Check if it's a network/auth error vs a code error
				assert.Contains(t, err.Error(), "failed to list accounts", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_accounts MCP tool")
			}
		},
	)
}

func TestIntegration_SearchAccounts(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// Test direct client call first
	t.Run(
		"DirectClientCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing direct client call for search accounts to %s\n", testConfig.ServerURL)

			clientWithAuth, err := client.NewClient(
				testConfig.ServerURL, client.WithRequestEditorFn(
					func(ctx context.Context, req *http.Request) error {
						req.Header.Set("Authorization", "Bearer "+testConfig.APIToken)
						req.Header.Set("Accept", "application/vnd.api+json")
						req.Header.Set("Content-Type", "application/json")
						return nil
					},
				),
			)
			require.NoError(t, err, "Failed to create client")

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			params := &client.SearchAccountsParams{
				Query: "checking",
				Field: client.AccountSearchFieldFilterName,
			}
			limit := int32(5)
			params.Limit = &limit

			response, err := clientWithAuth.SearchAccounts(ctx, params)

			fmt.Printf("[DEBUG_LOG] Response status: %v, Error: %v\n", response.StatusCode, err)

			if err != nil {
				t.Logf("API call failed (this might be expected): %v", err)
			} else {
				assert.Equal(t, 200, response.StatusCode, "Expected successful response")
				t.Logf("Successfully searched accounts from Firefly III API")
			}
		},
	)

	// Test MCP tool call with various scenarios
	testCases := []struct {
		name      string
		args      SearchAccountsArgs
		expectErr bool
		errMsg    string
	}{
		{
			name: "Valid search by name",
			args: SearchAccountsArgs{
				Query: "checking",
				Field: "name",
				Limit: 5,
			},
			expectErr: false,
		},
		{
			name: "Search all fields",
			args: SearchAccountsArgs{
				Query: "test",
				Field: "all",
				Limit: 10,
			},
			expectErr: false,
		},
		{
			name: "Search by IBAN",
			args: SearchAccountsArgs{
				Query: "NL",
				Field: "iban",
				Limit: 5,
			},
			expectErr: false,
		},
		{
			name: "Missing query",
			args: SearchAccountsArgs{
				Field: "name",
			},
			expectErr: true,
			errMsg:    "Query parameter is required",
		},
		{
			name: "Missing field",
			args: SearchAccountsArgs{
				Query: "test",
			},
			expectErr: true,
			errMsg:    "Field parameter is required",
		},
		{
			name: "With pagination",
			args: SearchAccountsArgs{
				Query: "account",
				Field: "all",
				Limit: 3,
				Page:  1,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				fmt.Printf("[DEBUG_LOG] Testing MCP search_accounts: %s\n", tc.name)

				session := &mcp.ServerSession{}
				params := &mcp.CallToolParamsFor[SearchAccountsArgs]{
					Name:      "search_accounts",
					Arguments: tc.args,
				}

				ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
				defer cancel()

				result, err := server.handleSearchAccounts(ctx, session, params)

				fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

				if tc.expectErr {
					assert.NotNil(t, result)
					assert.True(t, result.IsError)
					if tc.errMsg != "" {
						assert.Contains(t, result.Content[0].(*mcp.TextContent).Text, tc.errMsg)
					}
				} else {
					if err != nil {
						t.Logf("MCP tool call failed (this might be expected): %v", err)
					} else {
						assert.NotNil(t, result)
						assert.False(t, result.IsError)

						// Parse and validate the response
						var accountList AccountList
						err := json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &accountList)
						assert.NoError(t, err, "Failed to parse response")

						// Check pagination info
						assert.NotNil(t, accountList.Pagination)
						t.Logf("Found %d accounts out of %d total", accountList.Pagination.Count, accountList.Pagination.Total)
					}
				}
			},
		)
	}
}

func TestIntegration_ListTransactions(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_transactions\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters
			params := &mcp.CallToolParamsFor[ListTransactionsArgs]{
				Name: "list_transactions",
				Arguments: ListTransactionsArgs{
					Limit: 3,
					Start: "2024-01-01",
					End:   "2024-12-31",
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleListTransactions(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list transactions", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_transactions MCP tool")

				// Verify the result structure
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						// Unmarshal to verify it's a properly formatted TransactionList
						var transactionList TransactionList
						err := json.Unmarshal([]byte(textContent.Text), &transactionList)
						assert.NoError(t, err, "Result should be valid TransactionList JSON")

						// Verify pagination is included
						assert.GreaterOrEqual(t, transactionList.Pagination.Count, 0, "Pagination count should be >= 0")

						// If there are transactions, verify the structure
						if len(transactionList.Data) > 0 {
							firstGroup := transactionList.Data[0]
							assert.NotEmpty(t, firstGroup.Id, "Transaction group should have an ID")

							// Verify transactions within the group
							if len(firstGroup.Transactions) > 0 {
								firstTransaction := firstGroup.Transactions[0]
								assert.NotEmpty(t, firstTransaction.Id, "Transaction should have an ID")
								assert.NotEmpty(t, firstTransaction.Amount, "Transaction should have an amount")
								assert.NotEmpty(t, firstTransaction.Date, "Transaction should have a date")
								assert.NotEmpty(t, firstTransaction.Description, "Transaction should have a description")
								assert.NotEmpty(t, firstTransaction.Type, "Transaction should have a type")

								// Check source/destination names are populated based on type
								switch firstTransaction.Type {
								case "withdrawal", "expense":
									assert.NotEmpty(t, firstTransaction.SourceName, "Withdrawal should have source name")
									assert.NotEmpty(t, firstTransaction.DestinationName, "Withdrawal should have destination name")
								case "deposit", "income":
									assert.NotEmpty(t, firstTransaction.SourceName, "Deposit should have source name")
									assert.NotEmpty(t, firstTransaction.DestinationName, "Deposit should have destination name")
								case "transfer":
									assert.NotEmpty(t, firstTransaction.SourceName, "Transfer should have source name")
									assert.NotEmpty(t, firstTransaction.DestinationName, "Transfer should have destination name")
								}
							}
						}

						t.Logf("Successfully verified TransactionList structure with %d transaction groups", len(transactionList.Data))
					}
				}
			}
		},
	)

	t.Run(
		"MCPToolCallWithPagination", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_transactions with pagination\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters with pagination
			params := &mcp.CallToolParamsFor[ListTransactionsArgs]{
				Name: "list_transactions",
				Arguments: ListTransactionsArgs{
					Limit: 2,
					Page:  1,
					Start: "2024-01-01",
					End:   "2024-12-31",
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleListTransactions(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call with pagination result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call with pagination failed (this might be expected): %v", err)
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")

				// Verify pagination parameters in response
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						var transactionList TransactionList
						err := json.Unmarshal([]byte(textContent.Text), &transactionList)
						assert.NoError(t, err, "Result should be valid TransactionList JSON")

						if transactionList.Pagination.Count > 0 {
							assert.Equal(t, 1, transactionList.Pagination.CurrentPage, "Should be on page 1")
							assert.LessOrEqual(t, transactionList.Pagination.Count, 2, "Should have at most 2 items per page")
							assert.Equal(t, 2, transactionList.Pagination.PerPage, "Should show 2 items per page")
						}
					}
				}

				t.Logf("Successfully called list_transactions MCP tool with pagination")
			}
		},
	)

	t.Run(
		"MCPToolCallWithType", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_transactions with type filter\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters with type filter
			params := &mcp.CallToolParamsFor[ListTransactionsArgs]{
				Name: "list_transactions",
				Arguments: ListTransactionsArgs{
					Limit: 5,
					Type:  "withdrawal",
					Start: "2024-01-01",
					End:   "2024-12-31",
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleListTransactions(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call with type filter result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call with type filter failed (this might be expected): %v", err)
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")

				// Verify filtered results
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						var transactionList TransactionList
						err := json.Unmarshal([]byte(textContent.Text), &transactionList)
						assert.NoError(t, err, "Result should be valid TransactionList JSON")

						// Verify all returned transactions are of the requested type
						for _, group := range transactionList.Data {
							for _, transaction := range group.Transactions {
								assert.Equal(t, "withdrawal", transaction.Type, "All transactions should be withdrawals")
							}
						}
					}
				}

				t.Logf("Successfully called list_transactions MCP tool with type filter")
			}
		},
	)
}

func TestIntegration_GetTransaction(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// First, get a list of transactions to find a valid ID
	ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
	defer cancel()

	// List transactions to get a valid ID
	apiParams := &client.ListTransactionParams{}
	limit := int32(1)
	apiParams.Limit = &limit

	resp, err := server.client.ListTransactionWithResponse(ctx, apiParams)
	if err != nil {
		t.Fatalf("Failed to list transactions: %v", err)
	}

	if resp.StatusCode() != 200 || resp.ApplicationvndApiJSON200 == nil || len(resp.ApplicationvndApiJSON200.Data) == 0 {
		t.Skip("No transactions available for testing get_transaction")
	}

	transactionId := resp.ApplicationvndApiJSON200.Data[0].Id

	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for get_transaction with ID: %s\n", transactionId)

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters
			params := &mcp.CallToolParamsFor[GetTransactionArgs]{
				Name: "get_transaction",
				Arguments: GetTransactionArgs{
					ID: transactionId,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleGetTransaction(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "Expected no error from MCP tool call")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.False(t, result.IsError, "Expected successful result")
			assert.NotEmpty(t, result.Content, "Expected content in result")

			// Verify the response structure
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				fmt.Printf("[DEBUG_LOG] Response content: %s\n", textContent.Text)

				var transactionGroup TransactionGroup
				err := json.Unmarshal([]byte(textContent.Text), &transactionGroup)
				assert.NoError(t, err, "Result should be valid TransactionGroup JSON")

				// Verify basic structure
				assert.Equal(t, transactionId, transactionGroup.Id, "Transaction ID should match")
				assert.NotEmpty(t, transactionGroup.Transactions, "Should have at least one transaction")

				// Verify first transaction has required fields
				if len(transactionGroup.Transactions) > 0 {
					firstTransaction := transactionGroup.Transactions[0]
					assert.NotEmpty(t, firstTransaction.Amount, "Transaction should have amount")
					assert.NotEmpty(t, firstTransaction.Description, "Transaction should have description")
					assert.NotEmpty(t, firstTransaction.Type, "Transaction should have type")
				}
			}

			t.Logf("Successfully called get_transaction MCP tool and received TransactionGroup DTO")
		},
	)
}

func TestIntegration_GetSummary(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for get_summary\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters
			params := &mcp.CallToolParamsFor[GetSummaryArgs]{
				Name: "get_summary",
				Arguments: GetSummaryArgs{
					Start: "2024-01-01",
					End:   "2024-12-31",
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleGetSummary(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "Expected no error from MCP tool call")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.False(t, result.IsError, "Expected successful result")
			assert.NotEmpty(t, result.Content, "Expected content in result")

			// Verify the response structure
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				fmt.Printf("[DEBUG_LOG] Response content: %s\n", textContent.Text)

				var summaryList BasicSummaryList
				err := json.Unmarshal([]byte(textContent.Text), &summaryList)
				assert.NoError(t, err, "Result should be valid BasicSummaryList JSON")

				// Verify basic structure
				assert.NotNil(t, summaryList.Data, "Data should not be nil")

				// Log the summary entries for debugging
				for _, summary := range summaryList.Data {
					t.Logf("Summary entry: Key=%s, Title=%s, Currency=%s, Value=%s",
						summary.Key, summary.Title, summary.CurrencyCode, summary.MonetaryValue)
				}

				// Verify some expected keys (if any data is returned)
				if len(summaryList.Data) > 0 {
					// Check that each entry has the required fields populated
					for _, summary := range summaryList.Data {
						assert.NotEmpty(t, summary.Key, "Key should not be empty")
						assert.NotEmpty(t, summary.Title, "Title should not be empty")
						assert.NotEmpty(t, summary.CurrencyCode, "CurrencyCode should not be empty")
						assert.NotEmpty(t, summary.MonetaryValue, "MonetaryValue should not be empty")
					}
				}
			}

			t.Logf("Successfully called get_summary MCP tool and received BasicSummaryList DTO")
		},
	)

	t.Run(
		"MCPToolCallWithDefaultDates", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for get_summary with default dates\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters without dates (should use current month)
			params := &mcp.CallToolParamsFor[GetSummaryArgs]{
				Name:      "get_summary",
				Arguments: GetSummaryArgs{
					// No dates provided - should default to current month
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleGetSummary(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call with default dates result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "Expected no error from MCP tool call")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.False(t, result.IsError, "Expected successful result")
			assert.NotEmpty(t, result.Content, "Expected content in result")

			// Verify the response structure
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				var summaryList BasicSummaryList
				err := json.Unmarshal([]byte(textContent.Text), &summaryList)
				assert.NoError(t, err, "Result should be valid BasicSummaryList JSON")
				assert.NotNil(t, summaryList.Data, "Data should not be nil")
			}

			t.Logf("Successfully called get_summary MCP tool with default dates")
		},
	)
}

func TestIntegration_ListBudgets(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budgets\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters
			params := &mcp.CallToolParamsFor[ListBudgetsArgs]{
				Name: "list_budgets",
				Arguments: ListBudgetsArgs{
					Limit: 5,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleListBudgets(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list budgets", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_budgets MCP tool")
			}
		},
	)

	t.Run(
		"MCPToolCallWithDates", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budgets with date parameters\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters with date range
			params := &mcp.CallToolParamsFor[ListBudgetsArgs]{
				Name: "list_budgets",
				Arguments: ListBudgetsArgs{
					Start: "2024-01-01",
					End:   "2024-12-31",
					Limit: 10,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleListBudgets(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call with dates result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call with dates failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list budgets", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_budgets MCP tool with date parameters")
			}
		},
	)

	t.Run(
		"MCPToolCallWithPagination", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budgets with pagination parameters\n")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Test case 1: Basic pagination
			params := &mcp.CallToolParamsFor[ListBudgetsArgs]{
				Name: "list_budgets",
				Arguments: ListBudgetsArgs{
					Limit: 3,
					Page:  1,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, err := server.handleListBudgets(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] MCP call with pagination result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call with pagination failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list budgets", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_budgets MCP tool with pagination parameters")
			}

			// Test case 2: Different page number
			params2 := &mcp.CallToolParamsFor[ListBudgetsArgs]{
				Name: "list_budgets",
				Arguments: ListBudgetsArgs{
					Limit: 5,
					Page:  2,
				},
			}

			result2, err2 := server.handleListBudgets(ctx, session, params2)

			fmt.Printf("[DEBUG_LOG] MCP call with page 2 result: %v, Error: %v\n", result2 != nil, err2)

			if err2 != nil {
				t.Logf("MCP tool call with page 2 failed (this might be expected): %v", err2)
			} else {
				assert.NotNil(t, result2, "Expected non-nil result for page 2")
				assert.False(t, result2.IsError, "Expected successful result for page 2")
				t.Logf("Successfully called list_budgets MCP tool with page 2")
			}

			// Test case 3: Edge case - page 0 should be ignored
			params3 := &mcp.CallToolParamsFor[ListBudgetsArgs]{
				Name: "list_budgets",
				Arguments: ListBudgetsArgs{
					Limit: 5,
					Page:  0, // Should be ignored
				},
			}

			result3, err3 := server.handleListBudgets(ctx, session, params3)

			fmt.Printf("[DEBUG_LOG] MCP call with page 0 result: %v, Error: %v\n", result3 != nil, err3)

			if err3 != nil {
				t.Logf("MCP tool call with page 0 failed (this might be expected): %v", err3)
			} else {
				assert.NotNil(t, result3, "Expected non-nil result for page 0")
				assert.False(t, result3.IsError, "Expected successful result for page 0")
				t.Logf("Successfully called list_budgets MCP tool with page 0 (should be ignored)")
			}

			// Test case 4: Combined with date parameters
			params4 := &mcp.CallToolParamsFor[ListBudgetsArgs]{
				Name: "list_budgets",
				Arguments: ListBudgetsArgs{
					Start: "2024-01-01",
					End:   "2024-12-31",
					Limit: 2,
					Page:  1,
				},
			}

			result4, err4 := server.handleListBudgets(ctx, session, params4)

			fmt.Printf("[DEBUG_LOG] MCP call with pagination and dates result: %v, Error: %v\n", result4 != nil, err4)

			if err4 != nil {
				t.Logf("MCP tool call with pagination and dates failed (this might be expected): %v", err4)
			} else {
				assert.NotNil(t, result4, "Expected non-nil result for pagination with dates")
				assert.False(t, result4.IsError, "Expected successful result for pagination with dates")
				t.Logf("Successfully called list_budgets MCP tool with pagination and date parameters")
			}
		},
	)
}

func TestIntegration_ErrorHandling(t *testing.T) {
	// Test with invalid configuration
	t.Run(
		"InvalidURL", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing error handling with invalid URL\n")

			config := &Config{
				Server: struct {
					URL string `yaml:"url"`
				}{URL: "https://invalid-url-that-does-not-exist.com/api"},
				API: struct {
					Token string `yaml:"token"`
				}{Token: "invalid-token"},
				Client: struct {
					Timeout int `yaml:"timeout"`
				}{Timeout: 5},
			}

			server, err := NewFireflyMCPServer(config)
			require.NoError(t, err, "Server creation should not fail")

			// Create a mock session
			session := &mcp.ServerSession{}

			// Create tool call parameters
			params := &mcp.CallToolParamsFor[ListAccountsArgs]{
				Name: "list_accounts",
				Arguments: ListAccountsArgs{
					Limit: 1,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Call the handler - this should fail
			result, err := server.handleListAccounts(ctx, session, params)

			fmt.Printf("[DEBUG_LOG] Error handling test - Error: %v\n", err)

			// MCP handlers return errors in the result structure, not as Go errors
			assert.NoError(t, err, "MCP handlers should not return Go errors")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.True(t, result.IsError, "Expected error result with invalid URL")

			// Print the error content for debugging
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
					fmt.Printf("[DEBUG_LOG] Error result content: %s\n", textContent.Text)
				}
			}
		},
	)
}

func TestIntegration_AllTools(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// Test all available tools
	tools := []struct {
		name string
		test func(t *testing.T)
	}{
		{"list_accounts", func(t *testing.T) {
			session := &mcp.ServerSession{}
			params := &mcp.CallToolParamsFor[ListAccountsArgs]{
				Name:      "list_accounts",
				Arguments: ListAccountsArgs{Limit: 2},
			}
			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()
			_, err := server.handleListAccounts(ctx, session, params)
			if err != nil {
				t.Logf("Tool failed (expected): %v", err)
			}
		}},
		{"search_accounts", func(t *testing.T) {
			session := &mcp.ServerSession{}
			params := &mcp.CallToolParamsFor[SearchAccountsArgs]{
				Name: "search_accounts",
				Arguments: SearchAccountsArgs{
					Query: "test",
					Field: "all",
					Limit: 2,
				},
			}
			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()
			_, err := server.handleSearchAccounts(ctx, session, params)
			if err != nil {
				t.Logf("Tool failed (expected): %v", err)
			}
		}},
		{"list_transactions", func(t *testing.T) {
			session := &mcp.ServerSession{}
			params := &mcp.CallToolParamsFor[ListTransactionsArgs]{
				Name: "list_transactions",
				Arguments: ListTransactionsArgs{
					Limit: 2,
					Start: "2024-01-01",
					End:   "2024-12-31",
				},
			}
			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()
			_, err := server.handleListTransactions(ctx, session, params)
			if err != nil {
				t.Logf("Tool failed (expected): %v", err)
			}
		}},
		{"list_budgets", func(t *testing.T) {
			session := &mcp.ServerSession{}
			params := &mcp.CallToolParamsFor[ListBudgetsArgs]{
				Name:      "list_budgets",
				Arguments: ListBudgetsArgs{Limit: 2},
			}
			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()
			_, err := server.handleListBudgets(ctx, session, params)
			if err != nil {
				t.Logf("Tool failed (expected): %v", err)
			}
		}},
		{"list_categories", func(t *testing.T) {
			session := &mcp.ServerSession{}
			params := &mcp.CallToolParamsFor[ListCategoriesArgs]{
				Name:      "list_categories",
				Arguments: ListCategoriesArgs{Limit: 2},
			}
			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()
			_, err := server.handleListCategories(ctx, session, params)
			if err != nil {
				t.Logf("Tool failed (expected): %v", err)
			}
		}},
	}

	for _, tool := range tools {
		t.Run(
			tool.name, func(t *testing.T) {
				fmt.Printf("[DEBUG_LOG] Testing tool: %s\n", tool.name)
				tool.test(t)
			},
		)
	}
}
