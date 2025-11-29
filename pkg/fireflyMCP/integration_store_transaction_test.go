package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ensure mcp is used (for TextContent type assertion)
var _ mcp.Content = &mcp.TextContent{}

func TestIntegrationStoreTransaction(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// First, get an asset account to use in transactions
	ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
	defer cancel()

	// List asset accounts to get valid source account
	apiParams := &client.ListAccountParams{}
	accountType := client.AccountTypeFilter("asset")
	apiParams.Type = &accountType
	limit := int32(1)
	apiParams.Limit = &limit

	resp, err := server.client.ListAccountWithResponse(ctx, apiParams)
	require.NoError(t, err, "Failed to list accounts")
	require.Equal(t, 200, resp.StatusCode(), "Failed to get accounts")
	require.NotNil(t, resp.ApplicationvndApiJSON200)
	require.NotEmpty(t, resp.ApplicationvndApiJSON200.Data, "No asset accounts available for testing")

	assetAccountId := resp.ApplicationvndApiJSON200.Data[0].Id

	// Test data preparation
	currentDate := time.Now().Format("2006-01-02")
	uniqueDescription := fmt.Sprintf("Integration test transaction %d", time.Now().Unix())

	testCases := []struct {
		name           string
		args           TransactionStoreRequest
		expectErr      bool
		expectedStatus int
		validateResult func(*testing.T, *TransactionGroup)
	}{
		{
			name: "CreateWithdrawalTransaction",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:            "withdrawal",
						Date:            currentDate,
						Amount:          "25.50",
						Description:     uniqueDescription + " - withdrawal",
						SourceId:        &assetAccountId,
						DestinationName: strPtr("Test Grocery Store"),
					},
				},
			},
			expectErr:      false,
			expectedStatus: 201,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 1)
				assert.Equal(t, "withdrawal", result.Transactions[0].Type)
				// Check amount starts with expected value (API may return more decimal places)
				assert.Contains(t, result.Transactions[0].Amount, "25.5")
				assert.Contains(t, result.Transactions[0].Description, uniqueDescription)
			},
		},
		{
			name: "CreateDepositTransaction",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:          "deposit",
						Date:          currentDate,
						Amount:        "1000.00",
						Description:   uniqueDescription + " - deposit",
						SourceName:    strPtr("Test Employer"),
						DestinationId: &assetAccountId,
					},
				},
			},
			expectErr:      false,
			expectedStatus: 201,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 1)
				assert.Equal(t, "deposit", result.Transactions[0].Type)
				// Check amount starts with expected value (API may return more decimal places)
				assert.Contains(t, result.Transactions[0].Amount, "1000")
			},
		},
		{
			name: "CreateSplitTransaction",
			args: TransactionStoreRequest{
				GroupTitle: "Split transaction test",
				Transactions: []TransactionSplitRequest{
					{
						Type:            "withdrawal",
						Date:            currentDate,
						Amount:          "50.00",
						Description:     uniqueDescription + " - split part 1",
						SourceId:        &assetAccountId,
						DestinationName: strPtr("Test Store 1"),
					},
					{
						Type:            "withdrawal",
						Date:            currentDate,
						Amount:          "30.00",
						Description:     uniqueDescription + " - split part 2",
						SourceId:        &assetAccountId,
						DestinationName: strPtr("Test Store 2"),
					},
				},
			},
			expectErr:      false,
			expectedStatus: 201,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 2)
				assert.Equal(t, "Split transaction test", result.GroupTitle)
				// Check amounts start with expected values (API may return more decimal places)
				assert.Contains(t, result.Transactions[0].Amount, "50")
				assert.Contains(t, result.Transactions[1].Amount, "30")
			},
		},
		{
			name: "CreateWithOptionalFields",
			args: TransactionStoreRequest{
				ApplyRules:   false,
				FireWebhooks: true,
				Transactions: []TransactionSplitRequest{
					{
						Type:            "withdrawal",
						Date:            currentDate,
						Amount:          "75.00",
						Description:     uniqueDescription + " - with optional fields",
						SourceId:        &assetAccountId,
						DestinationName: strPtr("Test Store with Tags"),
						CategoryName:    strPtr("Groceries"),
						Tags:            []string{"integration-test", "automated"},
						Notes:           strPtr("This is a test transaction with optional fields"),
					},
				},
			},
			expectErr:      false,
			expectedStatus: 201,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 1)
				txn := result.Transactions[0]
				// Check amount starts with expected value (API may return more decimal places)
				assert.Contains(t, txn.Amount, "75")
				if txn.CategoryName != nil {
					assert.Equal(t, "Groceries", *txn.CategoryName)
				}
				assert.Contains(t, txn.Tags, "integration-test")
				assert.Contains(t, txn.Tags, "automated")
				if txn.Notes != nil {
					assert.Equal(t, "This is a test transaction with optional fields", *txn.Notes)
				}
			},
		},
		{
			name: "ValidationError_MissingRequiredFields",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type: "withdrawal",
						// Missing required fields: date, amount, description
					},
				},
			},
			expectErr:      true,
			expectedStatus: 0, // Won't reach API call
		},
		{
			name: "ValidationError_InvalidTransactionType",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:        "invalid_type",
						Date:        currentDate,
						Amount:      "50.00",
						Description: "Invalid transaction type test",
					},
				},
			},
			expectErr:      true,
			expectedStatus: 0, // Won't reach API call
		},
		{
			name: "ValidationError_EmptyTransactions",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{},
			},
			expectErr:      true,
			expectedStatus: 0, // Won't reach API call
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			result, _, err := server.handleStoreTransaction(ctx, nil, tc.args)
			require.NoError(t, err, "Handler should not return error")

			if tc.expectErr {
				assert.True(t, result.IsError, "Expected error result")
				// For validation errors, check error message
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						assert.Contains(t, textContent.Text, "Error:", "Expected error message")
					}
				}
			} else {
				// Debug: if marked as error, print the error content
				if result.IsError {
					if len(result.Content) > 0 {
						if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
							t.Logf("Error response: %s", textContent.Text)
						}
					}
				}

				assert.False(t, result.IsError, "Expected success result")
				require.Len(t, result.Content, 1, "Expected one content item")

				// Parse the JSON response
				textContent, ok := result.Content[0].(*mcp.TextContent)
				require.True(t, ok, "Expected text content")

				// Debug: print the actual response text
				t.Logf("Response text: %s", textContent.Text)

				var transactionGroup TransactionGroup
				err = json.Unmarshal([]byte(textContent.Text), &transactionGroup)
				require.NoError(t, err, "Failed to parse response JSON")

				// Validate the result
				if tc.validateResult != nil {
					tc.validateResult(t, &transactionGroup)
				}

				// Verify transaction was created by fetching it
				if transactionGroup.Id != "" {
					getArgs := GetTransactionArgs{
						ID: transactionGroup.Id,
					}

					getResult, _, err := server.handleGetTransaction(ctx, nil, getArgs)
					require.NoError(t, err, "Failed to get created transaction")
					assert.False(t, getResult.IsError, "Error getting created transaction")
				}
			}
		})
	}
}

// Test creating a transaction and then verifying it through GET
func TestIntegrationStoreTransaction_EndToEnd(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// Get an asset account
	ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
	defer cancel()

	apiParams := &client.ListAccountParams{}
	accountType := client.AccountTypeFilter("asset")
	apiParams.Type = &accountType
	limit := int32(1)
	apiParams.Limit = &limit

	resp, err := server.client.ListAccountWithResponse(ctx, apiParams)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode())
	require.NotEmpty(t, resp.ApplicationvndApiJSON200.Data)

	assetAccountId := resp.ApplicationvndApiJSON200.Data[0].Id

	// Create a transaction
	currentDate := time.Now().Format("2006-01-02")
	uniqueDescription := fmt.Sprintf("End-to-end test %d", time.Now().Unix())

	storeArgs := TransactionStoreRequest{
		Transactions: []TransactionSplitRequest{
			{
				Type:            "withdrawal",
				Date:            currentDate,
				Amount:          "123.45",
				Description:     uniqueDescription,
				SourceId:        &assetAccountId,
				DestinationName: strPtr("End-to-end Test Store"),
				CategoryName:    strPtr("Testing"),
				Tags:            []string{"e2e-test"},
			},
		},
	}

	// Create transaction
	createResult, _, err := server.handleStoreTransaction(ctx, nil, storeArgs)
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	// Parse created transaction
	textContent, ok := createResult.Content[0].(*mcp.TextContent)
	require.True(t, ok)

	var createdTransaction TransactionGroup
	err = json.Unmarshal([]byte(textContent.Text), &createdTransaction)
	require.NoError(t, err)
	require.NotEmpty(t, createdTransaction.Id)

	// Verify transaction through GET
	getArgs := GetTransactionArgs{
		ID: createdTransaction.Id,
	}

	getResult, _, err := server.handleGetTransaction(ctx, nil, getArgs)
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	// Parse fetched transaction
	textContent, ok = getResult.Content[0].(*mcp.TextContent)
	require.True(t, ok)

	var fetchedTransaction TransactionGroup
	err = json.Unmarshal([]byte(textContent.Text), &fetchedTransaction)
	require.NoError(t, err)

	// Verify the fetched transaction matches what we created
	assert.Equal(t, createdTransaction.Id, fetchedTransaction.Id)
	assert.Len(t, fetchedTransaction.Transactions, 1)
	// Check amount starts with expected value (API may return more decimal places)
	assert.Contains(t, fetchedTransaction.Transactions[0].Amount, "123.45")
	assert.Equal(t, uniqueDescription, fetchedTransaction.Transactions[0].Description)
	assert.Contains(t, fetchedTransaction.Transactions[0].Tags, "e2e-test")
}

// Helper functions are defined in bill_mapper_test.go
