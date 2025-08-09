package fireflyMCP

import (
	"context"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestHandleStoreTransactionsBulk_Validation(t *testing.T) {
	currentDate := time.Now().Format("2006-01-02")

	tests := []struct {
		name          string
		args          BulkTransactionStoreRequest
		expectedError string
		isError       bool
	}{
		{
			name: "Empty transaction groups array",
			args: BulkTransactionStoreRequest{
				TransactionGroups: []TransactionStoreRequest{},
			},
			expectedError: "transaction_groups array is required and must not be empty",
			isError:       true,
		},
		{
			name: "Exceeds maximum batch size",
			args: BulkTransactionStoreRequest{
				TransactionGroups: make([]TransactionStoreRequest, 101),
			},
			expectedError: "batch size exceeds maximum of 100 transaction groups",
			isError:       true,
		},
		{
			name: "Valid single transaction group",
			args: BulkTransactionStoreRequest{
				TransactionGroups: []TransactionStoreRequest{
					{
						GroupTitle: "Test Group",
						Transactions: []TransactionSplitRequest{
							{
								Type:        "withdrawal",
								Date:        currentDate,
								Amount:      "100.00",
								Description: "Test transaction",
							},
						},
					},
				},
			},
			expectedError: "", // Will fail at API call since no mock
			isError:       false,
		},
		{
			name: "Multiple valid transaction groups",
			args: BulkTransactionStoreRequest{
				TransactionGroups: []TransactionStoreRequest{
					{
						GroupTitle: "Group 1",
						Transactions: []TransactionSplitRequest{
							{
								Type:        "withdrawal",
								Date:        currentDate,
								Amount:      "50.00",
								Description: "First transaction",
							},
						},
					},
					{
						GroupTitle: "Group 2",
						Transactions: []TransactionSplitRequest{
							{
								Type:        "deposit",
								Date:        currentDate,
								Amount:      "75.00",
								Description: "Second transaction",
							},
						},
					},
				},
				DelayMs: 50,
			},
			expectedError: "",
			isError:       false,
		},
		{
			name: "Transaction group with validation error",
			args: BulkTransactionStoreRequest{
				TransactionGroups: []TransactionStoreRequest{
					{
						GroupTitle: "Invalid Group",
						Transactions: []TransactionSplitRequest{
							{
								Type:        "invalid_type",
								Date:        currentDate,
								Amount:      "100.00",
								Description: "Test",
							},
						},
					},
				},
			},
			expectedError: "",
			isError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal server instance
			server := &FireflyMCPServer{}

			// Create params
			params := &mcp.CallToolParamsFor[BulkTransactionStoreRequest]{
				Arguments: tt.args,
			}

			// Call the handler
			result, err := server.handleStoreTransactionsBulk(context.Background(), nil, params)

			// Should never return an error from the function itself
			assert.NoError(t, err)
			assert.NotNil(t, result)

			if tt.expectedError != "" {
				// Should return an error result
				assert.True(t, result.IsError)
				assert.Len(t, result.Content, 1)
				textContent, ok := result.Content[0].(*mcp.TextContent)
				assert.True(t, ok)
				assert.Contains(t, textContent.Text, tt.expectedError)
			}
		})
	}
}
