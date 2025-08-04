package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_ListBills(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	t.Run("MCPToolCall", func(t *testing.T) {
		fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_bills\n")

		// Create a mock session
		session := &mcp.ServerSession{}

		// Create tool call parameters
		params := &mcp.CallToolParamsFor[ListBillsArgs]{
			Name: "list_bills",
			Arguments: ListBillsArgs{
				Limit: 5,
				Page:  1,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
		defer cancel()

		// Call the handler directly
		result, err := server.handleListBills(ctx, session, params)

		fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

		assert.NoError(t, err, "MCP handlers should not return Go errors")
		assert.NotNil(t, result, "Expected non-nil result")
		assert.False(t, result.IsError, "Expected successful result")

		// Parse and verify the response
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				var billList BillList
				err := json.Unmarshal([]byte(textContent.Text), &billList)
				assert.NoError(t, err, "Failed to unmarshal response")
				
				// Verify response structure
				assert.NotNil(t, billList.Data, "Expected non-nil data array")
				assert.GreaterOrEqual(t, billList.Pagination.Total, 0, "Expected non-negative total")
				
				// If there are bills, verify their structure
				if len(billList.Data) > 0 {
					bill := billList.Data[0]
					assert.NotEmpty(t, bill.Id, "Expected non-empty bill ID")
					assert.NotEmpty(t, bill.Name, "Expected non-empty bill name")
					assert.NotEmpty(t, bill.AmountMin, "Expected non-empty amount min")
					assert.NotEmpty(t, bill.AmountMax, "Expected non-empty amount max")
					assert.NotEmpty(t, bill.RepeatFreq, "Expected non-empty repeat frequency")
					t.Logf("Found bill: ID=%s, Name=%s, AmountMin=%s, AmountMax=%s", 
						bill.Id, bill.Name, bill.AmountMin, bill.AmountMax)
				}
			}
		}
	})

	t.Run("MCPToolCallWithDateRange", func(t *testing.T) {
		fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_bills with date range\n")

		// Create a mock session
		session := &mcp.ServerSession{}

		// Create tool call parameters with date range
		params := &mcp.CallToolParamsFor[ListBillsArgs]{
			Name: "list_bills",
			Arguments: ListBillsArgs{
				Start: "2024-01-01",
				End:   "2024-12-31",
				Limit: 10,
				Page:  1,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
		defer cancel()

		// Call the handler directly
		result, err := server.handleListBills(ctx, session, params)

		fmt.Printf("[DEBUG_LOG] MCP call with date range result: %v, Error: %v\n", result != nil, err)

		assert.NoError(t, err, "MCP handlers should not return Go errors")
		assert.NotNil(t, result, "Expected non-nil result")
		assert.False(t, result.IsError, "Expected successful result")
		t.Logf("Successfully called list_bills MCP tool with date range")
	})
}

func TestIntegration_GetBill(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// First, get a list of bills to find a valid ID
	ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
	defer cancel()

	// Get list of bills
	listParams := &mcp.CallToolParamsFor[ListBillsArgs]{
		Name: "list_bills",
		Arguments: ListBillsArgs{
			Limit: 1,
		},
	}

	listResult, err := server.handleListBills(ctx, &mcp.ServerSession{}, listParams)
	require.NoError(t, err)
	require.NotNil(t, listResult)
	require.False(t, listResult.IsError)

	// Parse the response to get a bill ID
	var billList BillList
	if len(listResult.Content) > 0 {
		if textContent, ok := listResult.Content[0].(*mcp.TextContent); ok {
			err := json.Unmarshal([]byte(textContent.Text), &billList)
			require.NoError(t, err)
		}
	}

	if len(billList.Data) == 0 {
		t.Skip("No bills found in the system")
	}

	billID := billList.Data[0].Id

	t.Run("MCPToolCall", func(t *testing.T) {
		fmt.Printf("[DEBUG_LOG] Testing MCP tool call for get_bill with ID: %s\n", billID)

		// Create a mock session
		session := &mcp.ServerSession{}

		// Create tool call parameters
		params := &mcp.CallToolParamsFor[GetBillArgs]{
			Name: "get_bill",
			Arguments: GetBillArgs{
				ID: billID,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
		defer cancel()

		// Call the handler directly
		result, err := server.handleGetBill(ctx, session, params)

		fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

		assert.NoError(t, err, "MCP handlers should not return Go errors")
		assert.NotNil(t, result, "Expected non-nil result")
		assert.False(t, result.IsError, "Expected successful result")

		// Parse and verify the response
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				var bill Bill
				err := json.Unmarshal([]byte(textContent.Text), &bill)
				assert.NoError(t, err, "Failed to unmarshal response")
				
				// Verify bill structure
				assert.Equal(t, billID, bill.Id, "Expected matching bill ID")
				assert.NotEmpty(t, bill.Name, "Expected non-empty bill name")
				assert.NotEmpty(t, bill.AmountMin, "Expected non-empty amount min")
				assert.NotEmpty(t, bill.AmountMax, "Expected non-empty amount max")
				assert.NotEmpty(t, bill.RepeatFreq, "Expected non-empty repeat frequency")
				t.Logf("Retrieved bill: ID=%s, Name=%s", bill.Id, bill.Name)
			}
		}
	})

	t.Run("MCPToolCallInvalidID", func(t *testing.T) {
		fmt.Printf("[DEBUG_LOG] Testing MCP tool call for get_bill with invalid ID\n")

		// Create a mock session
		session := &mcp.ServerSession{}

		// Create tool call parameters with invalid ID
		params := &mcp.CallToolParamsFor[GetBillArgs]{
			Name: "get_bill",
			Arguments: GetBillArgs{
				ID: "99999", // Assuming this ID doesn't exist
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
		defer cancel()

		// Call the handler directly
		result, err := server.handleGetBill(ctx, session, params)

		fmt.Printf("[DEBUG_LOG] MCP call with invalid ID result: %v, Error: %v\n", result != nil, err)

		assert.NoError(t, err, "MCP handlers should not return Go errors")
		assert.NotNil(t, result, "Expected non-nil result")
		assert.True(t, result.IsError, "Expected error result for invalid ID")
	})
}

func TestIntegration_ListBillTransactions(t *testing.T) {
	testConfig := loadTestConfig(t)
	server := createTestServer(t, testConfig)

	// First, get a list of bills to find a valid ID
	ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
	defer cancel()

	// Get list of bills
	listParams := &mcp.CallToolParamsFor[ListBillsArgs]{
		Name: "list_bills",
		Arguments: ListBillsArgs{
			Limit: 1,
		},
	}

	listResult, err := server.handleListBills(ctx, &mcp.ServerSession{}, listParams)
	require.NoError(t, err)
	require.NotNil(t, listResult)
	require.False(t, listResult.IsError)

	// Parse the response to get a bill ID
	var billList BillList
	if len(listResult.Content) > 0 {
		if textContent, ok := listResult.Content[0].(*mcp.TextContent); ok {
			err := json.Unmarshal([]byte(textContent.Text), &billList)
			require.NoError(t, err)
		}
	}

	if len(billList.Data) == 0 {
		t.Skip("No bills found in the system")
	}

	billID := billList.Data[0].Id

	t.Run("MCPToolCall", func(t *testing.T) {
		fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_bill_transactions with ID: %s\n", billID)

		// Create a mock session
		session := &mcp.ServerSession{}

		// Create tool call parameters
		params := &mcp.CallToolParamsFor[ListBillTransactionsArgs]{
			Name: "list_bill_transactions",
			Arguments: ListBillTransactionsArgs{
				ID:    billID,
				Limit: 5,
				Page:  1,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
		defer cancel()

		// Call the handler directly
		result, err := server.handleListBillTransactions(ctx, session, params)

		fmt.Printf("[DEBUG_LOG] MCP call result: %v, Error: %v\n", result != nil, err)

		assert.NoError(t, err, "MCP handlers should not return Go errors")
		assert.NotNil(t, result, "Expected non-nil result")
		// Note: This might return an error if the bill has no transactions, which is OK
		t.Logf("Called list_bill_transactions MCP tool")

		// Parse and verify the response if successful
		if !result.IsError && len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				var transactionList TransactionList
				err := json.Unmarshal([]byte(textContent.Text), &transactionList)
				assert.NoError(t, err, "Failed to unmarshal response")
				
				// Verify response structure
				assert.NotNil(t, transactionList.Data, "Expected non-nil data array")
				t.Logf("Found %d transaction groups for bill %s", len(transactionList.Data), billID)
			}
		}
	})

	t.Run("MCPToolCallWithFilters", func(t *testing.T) {
		fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_bill_transactions with filters\n")

		// Create a mock session
		session := &mcp.ServerSession{}

		// Create tool call parameters with filters
		params := &mcp.CallToolParamsFor[ListBillTransactionsArgs]{
			Name: "list_bill_transactions",
			Arguments: ListBillTransactionsArgs{
				ID:    billID,
				Type:  "withdrawal",
				Start: "2024-01-01",
				End:   "2024-12-31",
				Limit: 10,
				Page:  1,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
		defer cancel()

		// Call the handler directly
		result, err := server.handleListBillTransactions(ctx, session, params)

		fmt.Printf("[DEBUG_LOG] MCP call with filters result: %v, Error: %v\n", result != nil, err)

		assert.NoError(t, err, "MCP handlers should not return Go errors")
		assert.NotNil(t, result, "Expected non-nil result")
		t.Logf("Called list_bill_transactions MCP tool with filters")
	})

	t.Run("MCPToolCallMissingID", func(t *testing.T) {
		fmt.Printf("[DEBUG_LOG] Testing MCP tool call for list_bill_transactions with missing ID\n")

		// Create a mock session
		session := &mcp.ServerSession{}

		// Create tool call parameters without ID
		params := &mcp.CallToolParamsFor[ListBillTransactionsArgs]{
			Name: "list_bill_transactions",
			Arguments: ListBillTransactionsArgs{
				// ID is missing
				Limit: 5,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), testConfig.Timeout)
		defer cancel()

		// Call the handler directly
		result, err := server.handleListBillTransactions(ctx, session, params)

		fmt.Printf("[DEBUG_LOG] MCP call with missing ID result: %v, Error: %v\n", result != nil, err)

		assert.NoError(t, err, "MCP handlers should not return Go errors")
		assert.NotNil(t, result, "Expected non-nil result")
		assert.True(t, result.IsError, "Expected error result")

		// Verify error message
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				assert.Contains(t, textContent.Text, "Bill ID is required", "Expected required parameter error message")
			}
		}
	})
}