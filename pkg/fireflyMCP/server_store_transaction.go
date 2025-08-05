package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// handleStoreTransaction creates a new transaction in Firefly III
func (s *FireflyMCPServer) handleStoreTransaction(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[TransactionStoreRequest],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Validate required fields
	if len(params.Arguments.Transactions) == 0 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: transactions array is required and must not be empty"},
			},
			IsError: true,
		}, nil
	}

	// Validate each transaction
	for i, txn := range params.Arguments.Transactions {
		// Validate required fields
		if txn.Type == "" {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: transaction[%d].type is required", i)},
				},
				IsError: true,
			}, nil
		}
		if txn.Date == "" {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: transaction[%d].date is required", i)},
				},
				IsError: true,
			}, nil
		}
		if txn.Amount == "" {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: transaction[%d].amount is required", i)},
				},
				IsError: true,
			}, nil
		}
		if txn.Description == "" {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: transaction[%d].description is required", i)},
				},
				IsError: true,
			}, nil
		}

		// Validate transaction type
		validTypes := map[string]bool{
			"withdrawal": true,
			"deposit":    true,
			"transfer":   true,
		}
		if !validTypes[txn.Type] {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf(
							"Error: transaction[%d].type must be one of: withdrawal, deposit, transfer",
							i,
						),
					},
				},
				IsError: true,
			}, nil
		}

		// Validate date format (basic check)
		if _, err := time.Parse("2006-01-02", txn.Date); err != nil {
			// Try parsing as datetime
			if _, err := time.Parse(time.RFC3339, txn.Date); err != nil {
				return &mcp.CallToolResultFor[struct{}]{
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf(
								"Error: transaction[%d].date must be in format YYYY-MM-DD or RFC3339",
								i,
							),
						},
					},
					IsError: true,
				}, nil
			}
		}
	}

	// Convert DTO to API model
	apiRequest := mapTransactionStoreRequestToAPI(&params.Arguments)

	// Call the API
	resp, err := s.client.StoreTransactionWithResponse(ctx, &client.StoreTransactionParams{}, *apiRequest)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error creating transaction: %v", err)},
			},
			IsError: true,
		}, nil
	}

	// Handle different response codes
	switch resp.StatusCode() {
	case 200:
		// Success - convert response to DTO
		if resp.ApplicationvndApiJSON200 == nil {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: received empty response from server"},
				},
				IsError: true,
			}, nil
		}

		transactionGroup := mapTransactionReadToTransactionGroup(&resp.ApplicationvndApiJSON200.Data)

		// Marshal to JSON
		jsonData, err := json.MarshalIndent(transactionGroup, "", "  ")
		if err != nil {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error marshaling response: %v", err)},
				},
				IsError: true,
			}, nil
		}

		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(jsonData)},
			},
		}, nil

	case 422:
		// Validation error
		errorMsg := "Validation error"
		if resp.JSON422 != nil && resp.JSON422.Message != nil {
			errorMsg = *resp.JSON422.Message
		}
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Validation error: %s", errorMsg)},
			},
			IsError: true,
		}, nil

	case 400:
		// Bad request
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Bad request: invalid data provided"},
			},
			IsError: true,
		}, nil

	default:
		// Other errors
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: unexpected status code %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}
}

// mapTransactionStoreRequestToAPI converts DTO to API model
func mapTransactionStoreRequestToAPI(req *TransactionStoreRequest) *client.StoreTransactionJSONRequestBody {
	apiReq := &client.StoreTransactionJSONRequestBody{
		Transactions: make([]client.TransactionSplitStore, len(req.Transactions)),
	}

	apiReq.ErrorIfDuplicateHash = &req.ErrorIfDuplicateHash
	apiReq.ApplyRules = &req.ApplyRules
	apiReq.FireWebhooks = &req.FireWebhooks
	if req.GroupTitle != "" {
		apiReq.GroupTitle = &req.GroupTitle
	}

	// Map transactions
	for i, txn := range req.Transactions {
		// Parse date string to time.Time
		parsedDate, err := time.Parse("2006-01-02", txn.Date)
		if err != nil {
			// Try parsing as datetime
			parsedDate, err = time.Parse(time.RFC3339, txn.Date)
			if err != nil {
				// Default to current date if parsing fails
				parsedDate = time.Now()
			}
		}

		apiTxn := client.TransactionSplitStore{
			Type:        client.TransactionTypeProperty(txn.Type),
			Date:        parsedDate,
			Amount:      txn.Amount,
			Description: txn.Description,
		}

		// Map optional fields
		if txn.SourceId != nil {
			apiTxn.SourceId = txn.SourceId
		}
		if txn.SourceName != nil {
			apiTxn.SourceName = txn.SourceName
		}
		if txn.DestinationId != nil {
			apiTxn.DestinationId = txn.DestinationId
		}
		if txn.DestinationName != nil {
			apiTxn.DestinationName = txn.DestinationName
		}
		if txn.CategoryId != nil {
			apiTxn.CategoryId = txn.CategoryId
		}
		if txn.CategoryName != nil {
			apiTxn.CategoryName = txn.CategoryName
		}
		if txn.BudgetId != nil {
			apiTxn.BudgetId = txn.BudgetId
		}
		if txn.BudgetName != nil {
			apiTxn.BudgetName = txn.BudgetName
		}
		if len(txn.Tags) > 0 {
			tags := make([]string, len(txn.Tags))
			copy(tags, txn.Tags)
			apiTxn.Tags = &tags
		}
		if txn.CurrencyId != nil {
			apiTxn.CurrencyId = txn.CurrencyId
		}
		if txn.CurrencyCode != nil {
			apiTxn.CurrencyCode = txn.CurrencyCode
		}
		if txn.ForeignAmount != nil {
			apiTxn.ForeignAmount = txn.ForeignAmount
		}
		if txn.ForeignCurrencyId != nil {
			apiTxn.ForeignCurrencyId = txn.ForeignCurrencyId
		}
		if txn.ForeignCurrencyCode != nil {
			apiTxn.ForeignCurrencyCode = txn.ForeignCurrencyCode
		}
		if txn.BillId != nil {
			apiTxn.BillId = txn.BillId
		}
		if txn.BillName != nil {
			apiTxn.BillName = txn.BillName
		}
		if txn.PiggyBankId != nil {
			id := int32(0)
			if val, err := strconv.Atoi(*txn.PiggyBankId); err == nil {
				id = int32(val)
			}
			apiTxn.PiggyBankId = &id
		}
		if txn.PiggyBankName != nil {
			apiTxn.PiggyBankName = txn.PiggyBankName
		}
		if txn.Notes != nil {
			apiTxn.Notes = txn.Notes
		}
		if txn.Reconciled != nil {
			apiTxn.Reconciled = txn.Reconciled
		}
		if txn.Order != nil {
			order := int32(*txn.Order)
			apiTxn.Order = &order
		}

		apiReq.Transactions[i] = apiTxn
	}

	return apiReq
}
