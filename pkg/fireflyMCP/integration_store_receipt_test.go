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
var _ = mcp.TextContent{}

func TestIntegrationStoreReceipt(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// First, get an asset account to use as source
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
	uniqueStoreName := fmt.Sprintf("Test Store %d", time.Now().Unix())

	testCases := []struct {
		name           string
		args           ReceiptStoreRequest
		expectErr      bool
		validateResult func(*testing.T, *TransactionGroup)
	}{
		{
			name: "BasicReceipt",
			args: ReceiptStoreRequest{
				StoreName:       uniqueStoreName + " Basic",
				ReceiptDate:     currentDate,
				SourceAccountId: &assetAccountId,
				Items: []ReceiptItemRequest{
					{Amount: "10.50", Description: "Item 1"},
					{Amount: "20.25", Description: "Item 2"},
					{Amount: "5.75", Description: "Item 3"},
				},
			},
			expectErr: false,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 3)
				assert.Contains(t, result.GroupTitle, "Receipt:")
				assert.Contains(t, result.GroupTitle, uniqueStoreName+" Basic")
				// All transactions should be withdrawals
				for _, txn := range result.Transactions {
					assert.Equal(t, "withdrawal", txn.Type)
				}
			},
		},
		{
			name: "ReceiptWithTotalValidation",
			args: ReceiptStoreRequest{
				StoreName:       uniqueStoreName + " With Total",
				ReceiptDate:     currentDate,
				SourceAccountId: &assetAccountId,
				TotalAmount:     getStringPtr("36.50"),
				Items: []ReceiptItemRequest{
					{Amount: "10.50", Description: "Item 1"},
					{Amount: "20.25", Description: "Item 2"},
					{Amount: "5.75", Description: "Item 3"},
				},
			},
			expectErr: false,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 3)
			},
		},
		{
			name: "ReceiptWithCategories",
			args: ReceiptStoreRequest{
				StoreName:           uniqueStoreName + " With Categories",
				ReceiptDate:         currentDate,
				SourceAccountId:     &assetAccountId,
				DefaultCategoryName: getStringPtr("Groceries"),
				Items: []ReceiptItemRequest{
					{Amount: "15.00", Description: "Food item"},
					{Amount: "25.00", Description: "Non-food item", CategoryName: getStringPtr("Household")},
				},
			},
			expectErr: false,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 2)
				// First item should have default category
				if result.Transactions[0].CategoryName != nil {
					assert.Equal(t, "Groceries", *result.Transactions[0].CategoryName)
				}
				// Second item should have overridden category
				if result.Transactions[1].CategoryName != nil {
					assert.Equal(t, "Household", *result.Transactions[1].CategoryName)
				}
			},
		},
		{
			name: "ReceiptWithTags",
			args: ReceiptStoreRequest{
				StoreName:       uniqueStoreName + " With Tags",
				ReceiptDate:     currentDate,
				SourceAccountId: &assetAccountId,
				Tags:            []string{"weekly-shopping", "groceries"},
				Items: []ReceiptItemRequest{
					{Amount: "30.00", Description: "Tagged item", Tags: []string{"special"}},
				},
			},
			expectErr: false,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 1)
				// Should have merged tags
				txn := result.Transactions[0]
				assert.Contains(t, txn.Tags, "weekly-shopping")
				assert.Contains(t, txn.Tags, "groceries")
				assert.Contains(t, txn.Tags, "special")
			},
		},
		{
			name: "ReceiptWithNotes",
			args: ReceiptStoreRequest{
				StoreName:       uniqueStoreName + " With Notes",
				ReceiptDate:     currentDate,
				SourceAccountId: &assetAccountId,
				Notes:           getStringPtr("Receipt level notes"),
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "First item"},
					{Amount: "20.00", Description: "Second item", Notes: getStringPtr("Item specific notes")},
				},
			},
			expectErr: false,
			validateResult: func(t *testing.T, result *TransactionGroup) {
				assert.NotEmpty(t, result.Id)
				assert.Len(t, result.Transactions, 2)
				// First item should have receipt notes
				if result.Transactions[0].Notes != nil {
					assert.Equal(t, "Receipt level notes", *result.Transactions[0].Notes)
				}
				// Second item should have item-specific notes
				if result.Transactions[1].Notes != nil {
					assert.Equal(t, "Item specific notes", *result.Transactions[1].Notes)
				}
			},
		},
		{
			name: "ValidationError_MismatchedTotal",
			args: ReceiptStoreRequest{
				StoreName:       uniqueStoreName + " Mismatched",
				ReceiptDate:     currentDate,
				SourceAccountId: &assetAccountId,
				TotalAmount:     getStringPtr("100.00"), // Wrong total
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
					{Amount: "20.00", Description: "Item 2"},
				},
			},
			expectErr: true,
		},
		{
			name: "ValidationError_EmptyStoreName",
			args: ReceiptStoreRequest{
				StoreName:       "",
				ReceiptDate:     currentDate,
				SourceAccountId: &assetAccountId,
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectErr: true,
		},
		{
			name: "ValidationError_InvalidDate",
			args: ReceiptStoreRequest{
				StoreName:       uniqueStoreName + " Invalid Date",
				ReceiptDate:     "29-11-2024", // Wrong format
				SourceAccountId: &assetAccountId,
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectErr: true,
		},
		{
			name: "ValidationError_EmptyItems",
			args: ReceiptStoreRequest{
				StoreName:       uniqueStoreName + " Empty Items",
				ReceiptDate:     currentDate,
				SourceAccountId: &assetAccountId,
				Items:           []ReceiptItemRequest{},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
			defer cancel()

			result, _, err := server.handleStoreReceipt(ctx, nil, tc.args)
			require.NoError(t, err, "Handler should not return error")

			if tc.expectErr {
				assert.True(t, result.IsError, "Expected error result")
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
						assert.Contains(t, textContent.Text, "Error:", "Expected error message")
					}
				}
			} else {
				if result.IsError {
					if len(result.Content) > 0 {
						if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
							t.Logf("Error response: %s", textContent.Text)
						}
					}
				}

				assert.False(t, result.IsError, "Expected success result")
				require.Len(t, result.Content, 1, "Expected one content item")

				textContent, ok := result.Content[0].(*mcp.TextContent)
				require.True(t, ok, "Expected text content")

				t.Logf("Response text: %s", textContent.Text)

				var transactionGroup TransactionGroup
				err = json.Unmarshal([]byte(textContent.Text), &transactionGroup)
				require.NoError(t, err, "Failed to parse response JSON")

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

func TestIntegrationStoreReceipt_EndToEnd(t *testing.T) {
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

	// Create a receipt
	currentDate := time.Now().Format("2006-01-02")
	uniqueStoreName := fmt.Sprintf("E2E Test Store %d", time.Now().Unix())

	storeArgs := ReceiptStoreRequest{
		StoreName:           uniqueStoreName,
		ReceiptDate:         currentDate,
		SourceAccountId:     &assetAccountId,
		TotalAmount:         getStringPtr("87.50"),
		DefaultCategoryName: getStringPtr("Shopping"),
		Tags:                []string{"e2e-test", "receipt"},
		Notes:               getStringPtr("End-to-end receipt test"),
		Items: []ReceiptItemRequest{
			{Amount: "45.00", Description: "Main purchase"},
			{Amount: "25.50", Description: "Secondary item", CategoryName: getStringPtr("Household")},
			{Amount: "17.00", Description: "Additional item", Tags: []string{"discount"}},
		},
	}

	// Create receipt
	createResult, _, err := server.handleStoreReceipt(ctx, nil, storeArgs)
	require.NoError(t, err)
	require.False(t, createResult.IsError)

	// Parse created transaction group
	textContent, ok := createResult.Content[0].(*mcp.TextContent)
	require.True(t, ok)

	var createdGroup TransactionGroup
	err = json.Unmarshal([]byte(textContent.Text), &createdGroup)
	require.NoError(t, err)
	require.NotEmpty(t, createdGroup.Id)

	// Verify group structure
	assert.Contains(t, createdGroup.GroupTitle, "Receipt:")
	assert.Contains(t, createdGroup.GroupTitle, uniqueStoreName)
	assert.Len(t, createdGroup.Transactions, 3)

	// Verify all transactions are withdrawals to the same store
	for _, txn := range createdGroup.Transactions {
		assert.Equal(t, "withdrawal", txn.Type)
		assert.Equal(t, uniqueStoreName, txn.DestinationName)
	}

	// Verify transaction through GET
	getArgs := GetTransactionArgs{
		ID: createdGroup.Id,
	}

	getResult, _, err := server.handleGetTransaction(ctx, nil, getArgs)
	require.NoError(t, err)
	require.False(t, getResult.IsError)

	// Parse fetched transaction
	textContent, ok = getResult.Content[0].(*mcp.TextContent)
	require.True(t, ok)

	var fetchedGroup TransactionGroup
	err = json.Unmarshal([]byte(textContent.Text), &fetchedGroup)
	require.NoError(t, err)

	// Verify the fetched group matches what we created
	assert.Equal(t, createdGroup.Id, fetchedGroup.Id)
	assert.Len(t, fetchedGroup.Transactions, 3)
	assert.Equal(t, createdGroup.GroupTitle, fetchedGroup.GroupTitle)
}

func getStringPtr(s string) *string {
	return &s
}
