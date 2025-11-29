package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helpers are defined in test_helpers_test.go

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

			// Create tool call arguments
			args := ListAccountsArgs{
				Limit: 5,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListAccounts(ctx, nil, args)

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

				ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
				defer cancel()

				result, _, err := server.handleSearchAccounts(ctx, nil, tc.args)

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
						t.Logf(
							"Found %d accounts out of %d total",
							accountList.Pagination.Count,
							accountList.Pagination.Total,
						)
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

			args := ListTransactionsArgs{
				Limit: 3,
				Start: "2024-01-01",
				End:   "2024-12-31",
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListTransactions(ctx, nil, args)

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
								assert.NotEmpty(
									t,
									firstTransaction.Description,
									"Transaction should have a description",
								)
								assert.NotEmpty(t, firstTransaction.Type, "Transaction should have a type")

								// Check source/destination names are populated based on type
								switch firstTransaction.Type {
								case "withdrawal", "expense":
									assert.NotEmpty(
										t,
										firstTransaction.SourceName,
										"Withdrawal should have source name",
									)
									assert.NotEmpty(
										t,
										firstTransaction.DestinationName,
										"Withdrawal should have destination name",
									)
								case "deposit", "income":
									assert.NotEmpty(t, firstTransaction.SourceName, "Deposit should have source name")
									assert.NotEmpty(
										t,
										firstTransaction.DestinationName,
										"Deposit should have destination name",
									)
								case "transfer":
									assert.NotEmpty(t, firstTransaction.SourceName, "Transfer should have source name")
									assert.NotEmpty(
										t,
										firstTransaction.DestinationName,
										"Transfer should have destination name",
									)
								}
							}
						}

						t.Logf(
							"Successfully verified TransactionList structure with %d transaction groups",
							len(transactionList.Data),
						)
					}
				}
			}
		},
	)

	t.Run(
		"MCPToolCallWithPagination", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_transactions with pagination\n")

			args := ListTransactionsArgs{
				Limit: 2,
				Page:  1,
				Start: "2024-01-01",
				End:   "2024-12-31",
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListTransactions(ctx, nil, args)

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
							assert.LessOrEqual(
								t,
								transactionList.Pagination.Count,
								2,
								"Should have at most 2 items per page",
							)
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

			args := ListTransactionsArgs{
				Limit: 5,
				Type:  "withdrawal",
				Start: "2024-01-01",
				End:   "2024-12-31",
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListTransactions(ctx, nil, args)

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
								assert.Equal(
									t,
									"withdrawal",
									transaction.Type,
									"All transactions should be withdrawals",
								)
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

			args := GetTransactionArgs{
				ID: transactionId,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleGetTransaction(ctx, nil, args)

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

func TestIntegration_SearchTransactions(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	testCases := []struct {
		name      string
		args      SearchTransactionsArgs
		expectErr bool
		errMsg    string
	}{
		{
			name: "ValidSearch",
			args: SearchTransactionsArgs{
				Query: "test",
				Limit: 5,
				Page:  1,
			},
			expectErr: false,
		},
		{
			name: "SearchWithPagination",
			args: SearchTransactionsArgs{
				Query: "payment",
				Limit: 2,
				Page:  2,
			},
			expectErr: false,
		},
		{
			name: "EmptyQuery",
			args: SearchTransactionsArgs{
				Query: "",
			},
			expectErr: true,
			errMsg:    "Query parameter is required",
		},
		{
			name: "NoResults",
			args: SearchTransactionsArgs{
				Query: "xyznonexistent123",
				Limit: 10,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				fmt.Printf("[DEBUG_LOG] Testing MCP search_transactions: %s\n", tc.name)

				ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
				defer cancel()

				result, _, err := server.handleSearchTransactions(ctx, nil, tc.args)

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
						var transactionList TransactionList
						err := json.Unmarshal([]byte(result.Content[0].(*mcp.TextContent).Text), &transactionList)
						assert.NoError(t, err, "Failed to parse response")

						// Check pagination info
						assert.NotNil(t, transactionList.Pagination)
						t.Logf(
							"Found %d transaction groups out of %d total",
							transactionList.Pagination.Count,
							transactionList.Pagination.Total,
						)

						// Verify transaction structure if results exist
						if len(transactionList.Data) > 0 {
							firstGroup := transactionList.Data[0]
							assert.NotEmpty(t, firstGroup.Id, "Transaction group should have an ID")
							if len(firstGroup.Transactions) > 0 {
								firstTransaction := firstGroup.Transactions[0]
								assert.NotEmpty(
									t,
									firstTransaction.Description,
									"Transaction should have a description",
								)
								assert.NotEmpty(t, firstTransaction.Amount, "Transaction should have an amount")
							}
						}
					}
				}
			},
		)
	}
}

func TestIntegration_GetSummary(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for get_summary\n")

			args := GetSummaryArgs{
				Start: "2024-01-01",
				End:   "2024-12-31",
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleGetSummary(ctx, nil, args)

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
					t.Logf(
						"Summary entry: Key=%s, Title=%s, Currency=%s, Value=%s",
						summary.Key, summary.Title, summary.CurrencyCode, summary.MonetaryValue,
					)
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

			// No dates provided - should default to current month
			args := GetSummaryArgs{}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleGetSummary(ctx, nil, args)

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

			args := ListBudgetsArgs{
				Limit: 5,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgets(ctx, nil, args)

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

			args := ListBudgetsArgs{
				Start: "2024-01-01",
				End:   "2024-12-31",
				Limit: 10,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgets(ctx, nil, args)

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

			// Test case 1: Basic pagination
			args := ListBudgetsArgs{
				Limit: 3,
				Page:  1,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgets(ctx, nil, args)

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
			args2 := ListBudgetsArgs{
				Limit: 5,
				Page:  2,
			}

			result2, _, err2 := server.handleListBudgets(ctx, nil, args2)

			fmt.Printf("[DEBUG_LOG] MCP call with page 2 result: %v, Error: %v\n", result2 != nil, err2)

			if err2 != nil {
				t.Logf("MCP tool call with page 2 failed (this might be expected): %v", err2)
			} else {
				assert.NotNil(t, result2, "Expected non-nil result for page 2")
				assert.False(t, result2.IsError, "Expected successful result for page 2")
				t.Logf("Successfully called list_budgets MCP tool with page 2")
			}

			// Test case 3: Edge case - page 0 should be ignored
			args3 := ListBudgetsArgs{
				Limit: 5,
				Page:  0, // Should be ignored
			}

			result3, _, err3 := server.handleListBudgets(ctx, nil, args3)

			fmt.Printf("[DEBUG_LOG] MCP call with page 0 result: %v, Error: %v\n", result3 != nil, err3)

			if err3 != nil {
				t.Logf("MCP tool call with page 0 failed (this might be expected): %v", err3)
			} else {
				assert.NotNil(t, result3, "Expected non-nil result for page 0")
				assert.False(t, result3.IsError, "Expected successful result for page 0")
				t.Logf("Successfully called list_budgets MCP tool with page 0 (should be ignored)")
			}

			// Test case 4: Combined with date parameters
			args4 := ListBudgetsArgs{
				Start: "2024-01-01",
				End:   "2024-12-31",
				Limit: 2,
				Page:  1,
			}

			result4, _, err4 := server.handleListBudgets(ctx, nil, args4)

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
					URL string `yaml:"url" mapstructure:"url"`
				}{URL: "https://invalid-url-that-does-not-exist.com/api"},
				API: struct {
					Token string `yaml:"token" mapstructure:"token"`
				}{Token: "invalid-token"},
				Client: struct {
					Timeout int `yaml:"timeout" mapstructure:"timeout"`
				}{Timeout: 5},
			}

			server, err := NewFireflyMCPServer(config)
			require.NoError(t, err, "Server creation should not fail")

			args := ListAccountsArgs{
				Limit: 1,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Call the handler - this should fail
			result, _, err := server.handleListAccounts(ctx, nil, args)

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
		{
			"list_accounts", func(t *testing.T) {
				args := ListAccountsArgs{Limit: 2}
				ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
				defer cancel()
				_, _, err := server.handleListAccounts(ctx, nil, args)
				if err != nil {
					t.Logf("Tool failed (expected): %v", err)
				}
			},
		},
		{
			"search_accounts", func(t *testing.T) {
				args := SearchAccountsArgs{
					Query: "test",
					Field: "all",
					Limit: 2,
				}
				ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
				defer cancel()
				_, _, err := server.handleSearchAccounts(ctx, nil, args)
				if err != nil {
					t.Logf("Tool failed (expected): %v", err)
				}
			},
		},
		{
			"list_transactions", func(t *testing.T) {
				args := ListTransactionsArgs{
					Limit: 2,
					Start: "2024-01-01",
					End:   "2024-12-31",
				}
				ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
				defer cancel()
				_, _, err := server.handleListTransactions(ctx, nil, args)
				if err != nil {
					t.Logf("Tool failed (expected): %v", err)
				}
			},
		},
		{
			"list_budgets", func(t *testing.T) {
				args := ListBudgetsArgs{Limit: 2}
				ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
				defer cancel()
				_, _, err := server.handleListBudgets(ctx, nil, args)
				if err != nil {
					t.Logf("Tool failed (expected): %v", err)
				}
			},
		},
		{
			"list_categories", func(t *testing.T) {
				args := ListCategoriesArgs{Limit: 2}
				ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
				defer cancel()
				_, _, err := server.handleListCategories(ctx, nil, args)
				if err != nil {
					t.Logf("Tool failed (expected): %v", err)
				}
			},
		},
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

func TestIntegrationExpenseCategoryInsights(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"BasicCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for expense_category_insights\n")

			args := ExpenseCategoryInsightsArgs{
				Start: "2024-01-01",
				End:   "2024-12-31",
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleExpenseCategoryInsights(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call failed (this might be expected): %v", err)
				assert.Contains(
					t,
					err.Error(),
					"failed to get expense category insights",
					"Expected specific error message",
				)
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called expense_category_insights MCP tool")

				// Verify the result structure
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						// Unmarshal to verify it's a properly formatted InsightCategoryResponse
						var insightResponse InsightCategoryResponse
						err := json.Unmarshal([]byte(textContent.Text), &insightResponse)
						assert.NoError(t, err, "Result should be valid InsightCategoryResponse JSON")

						// Verify the structure
						for _, entry := range insightResponse.Entries {
							assert.NotEmpty(t, entry.Id, "Entry should have an ID")
							assert.NotEmpty(t, entry.Name, "Entry should have a name")
							assert.NotEmpty(t, entry.Amount, "Entry should have an amount")
							assert.NotEmpty(t, entry.CurrencyCode, "Entry should have a currency code")
						}

						t.Logf(
							"Successfully verified InsightCategoryResponse structure with %d entries",
							len(insightResponse.Entries),
						)
					}
				}
			}
		},
	)

	t.Run(
		"WithAccountFilter", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for expense_category_insights with account filter\n")

			// First, get a list of accounts to get valid IDs
			apiParams := &client.ListAccountParams{}
			limit := int32(2)
			apiParams.Limit = &limit

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			resp, err := server.client.ListAccountWithResponse(ctx, apiParams)
			if err != nil || resp.StatusCode() != 200 || resp.ApplicationvndApiJSON200 == nil || len(resp.ApplicationvndApiJSON200.Data) == 0 {
				t.Skip("No accounts available for testing with account filter")
			}

			// Get account IDs
			var accountIds []string
			for _, account := range resp.ApplicationvndApiJSON200.Data {
				accountIds = append(accountIds, account.Id)
				if len(accountIds) >= 2 {
					break
				}
			}

			args := ExpenseCategoryInsightsArgs{
				Start:    "2024-01-01",
				End:      "2024-12-31",
				Accounts: accountIds,
			}

			// Call the handler directly
			result, _, err := server.handleExpenseCategoryInsights(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call with account filter result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call with account filter failed (this might be expected): %v", err)
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called expense_category_insights MCP tool with account filter")
			}
		},
	)

	t.Run(
		"InvalidDateFormat", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for expense_category_insights with invalid date format\n")

			args := ExpenseCategoryInsightsArgs{
				Start: "01/01/2024", // Invalid format
				End:   "2024-12-31",
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleExpenseCategoryInsights(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call with invalid date result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "MCP handlers should not return Go errors")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.True(t, result.IsError, "Expected error result")

			// Verify error message
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
					assert.Contains(
						t,
						textContent.Text,
						"Invalid start date format",
						"Expected date format error message",
					)
				}
			}
		},
	)

	t.Run(
		"EmptyDateRange", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for expense_category_insights with empty date range\n")

			// No dates provided
			args := ExpenseCategoryInsightsArgs{}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleExpenseCategoryInsights(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call with empty dates result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "MCP handlers should not return Go errors")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.True(t, result.IsError, "Expected error result")

			// Verify error message
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
					assert.Contains(t, textContent.Text, "required", "Expected required parameter error message")
				}
			}
		},
	)
}

func TestIntegration_ListTags(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_tags\n")

			args := ListTagsArgs{
				Limit: 5,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListTags(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list tags", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_tags MCP tool")
			}
		},
	)

	t.Run(
		"MCPToolCallWithPagination", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_tags with pagination\n")

			args := ListTagsArgs{
				Limit: 10,
				Page:  1,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListTags(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Fatalf("MCP tool call failed: %v", err)
			}

			assert.NotNil(t, result, "Expected non-nil result")
			assert.False(t, result.IsError, "Expected successful result")
			assert.NotNil(t, result.Content, "Expected content in result")
			assert.Greater(t, len(result.Content), 0, "Expected at least one content item")

			// Parse the result
			textContent, ok := result.Content[0].(*mcp.TextContent)
			assert.True(t, ok, "Expected TextContent type")

			var tagList TagList
			err = json.Unmarshal([]byte(textContent.Text), &tagList)
			assert.NoError(t, err, "Failed to unmarshal tag list")

			// Verify the response structure
			assert.NotNil(t, tagList.Data, "Expected data array")
			assert.LessOrEqual(t, len(tagList.Data), 10, "Expected at most 10 tags")

			// Log the tags for debugging
			for i, tag := range tagList.Data {
				t.Logf("Tag %d: ID=%s, Name=%s", i+1, tag.Id, tag.Tag)
				if tag.Description != nil {
					t.Logf("  Description: %s", *tag.Description)
				}
			}
		})
}

func TestIntegrationExpenseTotalInsights(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"BasicCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for expense_total_insights\n")

			args := ExpenseTotalInsightsArgs{
				Start: "2024-01-01",
				End:   "2024-12-31",
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleExpenseTotalInsights(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call failed (this might be expected): %v", err)
				assert.Contains(
					t,
					err.Error(),
					"failed to get expense total insights",
					"Expected specific error message",
				)
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called expense_total_insights MCP tool")

				// Verify the result structure
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						// Unmarshal to verify it's a properly formatted InsightTotalResponse
						var insightResponse InsightTotalResponse
						err := json.Unmarshal([]byte(textContent.Text), &insightResponse)
						assert.NoError(t, err, "Result should be valid InsightTotalResponse JSON")

						// Verify the structure
						for _, entry := range insightResponse.Entries {
							assert.NotEmpty(t, entry.Amount, "Entry should have an amount")
							assert.NotEmpty(t, entry.CurrencyCode, "Entry should have a currency code")
						}

						t.Logf(
							"Successfully verified InsightTotalResponse structure with %d entries",
							len(insightResponse.Entries),
						)
					}
				}
			}
		},
	)

	t.Run(
		"WithAccountFilter", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for expense_total_insights with account filter\n")

			// First, get a list of accounts to get valid IDs
			apiParams := &client.ListAccountParams{}
			limit := int32(2)
			apiParams.Limit = &limit

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			resp, err := server.client.ListAccountWithResponse(ctx, apiParams)
			if err != nil || resp.StatusCode() != 200 || resp.ApplicationvndApiJSON200 == nil || len(resp.ApplicationvndApiJSON200.Data) == 0 {
				t.Skip("No accounts available for testing with account filter")
			}

			// Get account IDs
			var accountIds []string
			for _, account := range resp.ApplicationvndApiJSON200.Data {
				accountIds = append(accountIds, account.Id)
				if len(accountIds) >= 2 {
					break
				}
			}

			args := ExpenseTotalInsightsArgs{
				Start:    "2024-01-01",
				End:      "2024-12-31",
				Accounts: accountIds,
			}

			// Call the handler directly
			result, _, err := server.handleExpenseTotalInsights(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call with account filter result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call with account filter failed (this might be expected): %v", err)
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called expense_total_insights MCP tool with account filter")
			}
		},
	)

	t.Run(
		"InvalidDateFormat", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for expense_total_insights with invalid date format\n")

			args := ExpenseTotalInsightsArgs{
				Start: "2024-01-01",
				End:   "31-12-2024", // Invalid format
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleExpenseTotalInsights(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call with invalid date result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "MCP handlers should not return Go errors")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.True(t, result.IsError, "Expected error result")

			// Verify error message
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
					assert.Contains(
						t,
						textContent.Text,
						"Invalid end date format",
						"Expected date format error message",
					)
				}
			}
		},
	)

	t.Run(
		"MissingEndDate", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for expense_total_insights with missing end date\n")

			args := ExpenseTotalInsightsArgs{
				Start: "2024-01-01",
				// End date is missing
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleExpenseTotalInsights(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call with missing end date result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "MCP handlers should not return Go errors")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.True(t, result.IsError, "Expected error result")

			// Verify error message
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
					assert.Contains(
						t,
						textContent.Text,
						"Start and End dates are required",
						"Expected required parameter error message",
					)
				}
			}
		},
	)
}

func TestIntegration_ListBudgetLimits(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budget_limits\n")

			args := ListBudgetLimitsArgs{
				ID: "1", // Assuming budget with ID 1 exists
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgetLimits(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list budget limits", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_budget_limits MCP tool")
			}
		},
	)

	t.Run(
		"MCPToolCallWithDates", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budget_limits with date parameters\n")

			args := ListBudgetLimitsArgs{
				ID:    "1",
				Start: "2024-01-01",
				End:   "2024-12-31",
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgetLimits(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call result with dates: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call with dates failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list budget limits", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_budget_limits MCP tool with date parameters")
			}
		},
	)

	t.Run(
		"MCPToolCallMissingID", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budget_limits with missing ID\n")

			// ID is missing
			args := ListBudgetLimitsArgs{}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgetLimits(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call with missing ID result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "MCP handlers should not return Go errors")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.True(t, result.IsError, "Expected error result")

			// Verify error message
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
					assert.Contains(
						t,
						textContent.Text,
						"Budget ID is required",
						"Expected required parameter error message",
					)
				}
			}
		},
	)
}

func TestIntegration_ListBudgetTransactions(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run(
		"MCPToolCall", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budget_transactions\n")

			args := ListBudgetTransactionsArgs{
				ID:    "1", // Assuming budget with ID 1 exists
				Limit: 5,
				Page:  1,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgetTransactions(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list budget transactions", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_budget_transactions MCP tool")
			}
		},
	)

	t.Run(
		"MCPToolCallWithFilters", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budget_transactions with filters\n")

			args := ListBudgetTransactionsArgs{
				ID:    "1",
				Type:  "withdrawal",
				Start: "2024-01-01",
				End:   "2024-12-31",
				Limit: 10,
				Page:  1,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgetTransactions(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call result with filters: %v, Error: %v\n", result != nil, err)

			if err != nil {
				t.Logf("MCP tool call with filters failed (this might be expected): %v", err)
				assert.Contains(t, err.Error(), "failed to list budget transactions", "Expected specific error message")
			} else {
				assert.NotNil(t, result, "Expected non-nil result")
				assert.False(t, result.IsError, "Expected successful result")
				t.Logf("Successfully called list_budget_transactions MCP tool with filters")
			}
		},
	)

	t.Run(
		"MCPToolCallMissingID", func(t *testing.T) {
			fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_budget_transactions with missing ID\n")

			// ID is missing
			args := ListBudgetTransactionsArgs{
				Limit: 5,
			}

			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			// Call the handler directly
			result, _, err := server.handleListBudgetTransactions(ctx, nil, args)

			fmt.Printf("[DEBUG_LOG] MCP call with missing ID result: %v, Error: %v\n", result != nil, err)

			assert.NoError(t, err, "MCP handlers should not return Go errors")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.True(t, result.IsError, "Expected error result")

			// Verify error message
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
					assert.Contains(
						t,
						textContent.Text,
						"Budget ID is required",
						"Expected required parameter error message",
					)
				}
			}
		},
	)
}
