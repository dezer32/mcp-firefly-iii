package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// handleStoreTransaction creates a new transaction in Firefly III
func (s *FireflyMCPServer) handleStoreTransaction(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args TransactionStoreRequest,
) (*mcp.CallToolResult, any, error) {
	// Validate required fields
	if len(args.Transactions) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: transactions array is required and must not be empty"},
			},
			IsError: true,
		}, nil, nil
	}

	// Validate each transaction
	for i, txn := range args.Transactions {
		// Validate required fields
		if txn.Type == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: transaction[%d].type is required", i)},
				},
				IsError: true,
			}, nil, nil
		}
		if txn.Date == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: transaction[%d].date is required", i)},
				},
				IsError: true,
			}, nil, nil
		}
		if txn.Amount == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: transaction[%d].amount is required", i)},
				},
				IsError: true,
			}, nil, nil
		}
		if txn.Description == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: transaction[%d].description is required", i)},
				},
				IsError: true,
			}, nil, nil
		}

		// Validate transaction type
		validTypes := map[string]bool{
			"withdrawal": true,
			"deposit":    true,
			"transfer":   true,
		}
		if !validTypes[txn.Type] {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf(
							"Error: transaction[%d].type must be one of: withdrawal, deposit, transfer",
							i,
						),
					},
				},
				IsError: true,
			}, nil, nil
		}

		// Validate date format (basic check)
		if _, err := time.Parse("2006-01-02", txn.Date); err != nil {
			// Try parsing as datetime
			if _, err := time.Parse(time.RFC3339, txn.Date); err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf(
								"Error: transaction[%d].date must be in format YYYY-MM-DD or RFC3339",
								i,
							),
						},
					},
					IsError: true,
				}, nil, nil
			}
		}
	}

	// Convert DTO to API model
	apiRequest := mapTransactionStoreRequestToAPI(&args)

	// Get API client
	apiClient, err := s.getClient(ctx)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	// Call the API
	resp, err := apiClient.StoreTransactionWithResponse(ctx, &client.StoreTransactionParams{}, *apiRequest)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error creating transaction: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	// Debug: Check content type
	contentType := resp.HTTPResponse.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "application/vnd.api+json") {
		bodyPreview := string(resp.Body)
		if len(bodyPreview) > 500 {
			bodyPreview = bodyPreview[:500] + "..."
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: Expected JSON response but got %s (status: %d, body preview: %s)",
					contentType, resp.StatusCode(), bodyPreview)},
			},
			IsError: true,
		}, nil, nil
	}

	// Handle different response codes
	switch resp.StatusCode() {
	case 200, 201:
		// Success - convert response to DTO
		// The generated client only populates ApplicationvndApiJSON200 for status 200 when it can parse it
		// Sometimes the parsing fails, so we need to handle the raw body
		var transactionSingle client.TransactionSingle

		if resp.ApplicationvndApiJSON200 != nil {
			// Already parsed successfully
			transactionSingle = *resp.ApplicationvndApiJSON200
		} else {
			// Parse the body manually for both 200 and 201
			if len(resp.Body) == 0 {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Error: empty response body for status %d", resp.StatusCode())},
					},
					IsError: true,
				}, nil, nil
			}

			// Try to parse the response body
			if err := json.Unmarshal(resp.Body, &transactionSingle); err != nil {
				// Debug: Log the body content for debugging
				bodyPreview := string(resp.Body)
				if len(bodyPreview) > 200 {
					bodyPreview = bodyPreview[:200] + "..."
				}
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Error parsing response (status %d): %v (body: %s)",
							resp.StatusCode(), err, bodyPreview)},
					},
					IsError: true,
				}, nil, nil
			}
		}

		transactionGroup := mapTransactionReadToTransactionGroup(&transactionSingle.Data)

		// Marshal to JSON
		jsonData, err := json.MarshalIndent(transactionGroup, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error marshaling response: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(jsonData)},
			},
		}, nil, nil

	case 422:
		// Validation error
		errorMsg := "Validation error"
		if resp.JSON422 != nil && resp.JSON422.Message != nil {
			errorMsg = *resp.JSON422.Message
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Validation error: %s", errorMsg)},
			},
			IsError: true,
		}, nil, nil

	case 400:
		// Bad request
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Bad request: invalid data provided"},
			},
			IsError: true,
		}, nil, nil

	default:
		// Other errors
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: unexpected status code %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}
}

// mapTransactionStoreRequestToAPI converts DTO to API model
func mapTransactionStoreRequestToAPI(req *TransactionStoreRequest) *client.StoreTransactionJSONRequestBody {
	apiReq := &client.StoreTransactionJSONRequestBody{
		Transactions: make([]client.TransactionSplitStore, len(req.Transactions)),
	}

	// Only set boolean fields if they are true (to match API expectations)
	if req.ErrorIfDuplicateHash {
		apiReq.ErrorIfDuplicateHash = &req.ErrorIfDuplicateHash
	}
	if req.ApplyRules {
		apiReq.ApplyRules = &req.ApplyRules
	}
	if req.FireWebhooks {
		apiReq.FireWebhooks = &req.FireWebhooks
	}
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

		// Order is required by Firefly III API, default to index if not provided
		if txn.Order != nil {
			order := int32(*txn.Order)
			apiTxn.Order = &order
		} else {
			// Use the transaction index as the default order
			order := int32(i)
			apiTxn.Order = &order
		}

		apiReq.Transactions[i] = apiTxn
	}

	return apiReq
}
