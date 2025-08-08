package fireflyMCP

import (
	"context"
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestHandleStoreTransaction_Validation(t *testing.T) {
	currentDate := time.Now().Format("2006-01-02")
	currentDateRFC3339 := time.Now().Format(time.RFC3339)
	
	tests := []struct {
		name          string
		args          TransactionStoreRequest
		expectedError string
	}{
		{
			name: "Empty transactions array",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{},
			},
			expectedError: "transactions array is required and must not be empty",
		},
		{
			name: "Missing type field",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Date:        currentDate,
						Amount:      "100.50",
						Description: "Test",
					},
				},
			},
			expectedError: "transaction[0].type is required",
		},
		{
			name: "Invalid transaction type",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:        "invalid",
						Date:        currentDate,
						Amount:      "100.50",
						Description: "Test",
					},
				},
			},
			expectedError: "transaction[0].type must be one of: withdrawal, deposit, transfer",
		},
		{
			name: "Missing date field",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:        "withdrawal",
						Amount:      "100.50",
						Description: "Test",
					},
				},
			},
			expectedError: "transaction[0].date is required",
		},
		{
			name: "Invalid date format",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:        "withdrawal",
						Date:        "15-01-2024", // Wrong format
						Amount:      "100.50",
						Description: "Test",
					},
				},
			},
			expectedError: "transaction[0].date must be in format YYYY-MM-DD or RFC3339",
		},
		{
			name: "Missing amount field",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:        "withdrawal",
						Date:        currentDate,
						Description: "Test",
					},
				},
			},
			expectedError: "transaction[0].amount is required",
		},
		{
			name: "Missing description field",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:   "withdrawal",
						Date:   currentDate,
						Amount: "100.50",
					},
				},
			},
			expectedError: "transaction[0].description is required",
		},
		{
			name: "Valid date in RFC3339 format",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:        "withdrawal",
						Date:        currentDateRFC3339,
						Amount:      "100.50",
						Description: "Test",
					},
				},
			},
			expectedError: "", // This should pass validation but fail on API call (no mock)
		},
		{
			name: "Multiple transactions with second invalid",
			args: TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:        "withdrawal",
						Date:        currentDate,
						Amount:      "50.00",
						Description: "First",
					},
					{
						Type:        "deposit",
						Date:        currentDate,
						Amount:      "", // Missing amount
						Description: "Second",
					},
				},
			},
			expectedError: "transaction[1].amount is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal server instance with a mock client
			server := &FireflyMCPServer{
				client: &client.ClientWithResponses{},
			}

			// Create params
			params := &mcp.CallToolParamsFor[TransactionStoreRequest]{
				Arguments: tt.args,
			}

			// Call the handler
			result, err := server.handleStoreTransaction(context.Background(), nil, params)

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
			} else {
				// For the valid case, it would fail on API call since we don't have a mock
				// but we're only testing validation here
				if result.IsError {
					textContent, ok := result.Content[0].(*mcp.TextContent)
					assert.True(t, ok)
					// Should fail with network/API error, not validation
					assert.NotContains(t, textContent.Text, "transaction[")
				}
			}
		})
	}
}

func TestMapTransactionStoreRequestToAPI_Coverage(t *testing.T) {
	currentDate := time.Now().Format("2006-01-02")
	
	// Test with nil values to improve coverage
	req := &TransactionStoreRequest{
		Transactions: []TransactionSplitRequest{
			{
				Type:        "withdrawal",
				Date:        currentDate,
				Amount:      "100.00",
				Description: "Test",
			},
		},
	}

	result := mapTransactionStoreRequestToAPI(req)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Transactions)
	assert.Len(t, result.Transactions, 1)

	// Test with empty optional fields
	req2 := &TransactionStoreRequest{
		ErrorIfDuplicateHash: false,
		ApplyRules:           false,
		FireWebhooks:         false,
		GroupTitle:           "",
		Transactions: []TransactionSplitRequest{
			{
				Type:                "deposit",
				Date:                currentDate,
				Amount:              "50.00",
				Description:         "Deposit test",
				SourceId:            nil,
				SourceName:          nil,
				DestinationId:       nil,
				DestinationName:     nil,
				CategoryId:          nil,
				CategoryName:        nil,
				BudgetId:            nil,
				BudgetName:          nil,
				CurrencyId:          nil,
				CurrencyCode:        nil,
				ForeignAmount:       nil,
				ForeignCurrencyId:   nil,
				ForeignCurrencyCode: nil,
				BillId:              nil,
				BillName:            nil,
				PiggyBankId:         nil,
				PiggyBankName:       nil,
				Notes:               nil,
				Reconciled:          nil,
				Order:               nil,
				Tags:                nil,
			},
		},
	}

	result2 := mapTransactionStoreRequestToAPI(req2)
	assert.NotNil(t, result2)
	assert.Nil(t, result2.ErrorIfDuplicateHash)
	assert.Nil(t, result2.ApplyRules)
	assert.Nil(t, result2.FireWebhooks)
	assert.Nil(t, result2.GroupTitle)
}
