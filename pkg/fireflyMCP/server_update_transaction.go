package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// handleUpdateTransaction updates an existing transaction in Firefly III
func (s *FireflyMCPServer) handleUpdateTransaction(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args UpdateTransactionArgs,
) (*mcp.CallToolResult, any, error) {
	// Validate required ID
	if args.ID == "" {
		return newErrorResult("Error: transaction ID is required")
	}

	// Validate each transaction if provided
	for i, txn := range args.Transactions {
		// Validate transaction type if provided
		if txn.Type != "" {
			validTypes := map[string]bool{
				"withdrawal": true,
				"deposit":    true,
				"transfer":   true,
			}
			if !validTypes[txn.Type] {
				return newErrorResult(fmt.Sprintf(
					"Error: transaction[%d].type must be one of: withdrawal, deposit, transfer", i))
			}
		}

		// Validate date format if provided
		if txn.Date != "" {
			if _, err := time.Parse("2006-01-02", txn.Date); err != nil {
				if _, err := time.Parse(time.RFC3339, txn.Date); err != nil {
					return newErrorResult(fmt.Sprintf(
						"Error: transaction[%d].date must be in format YYYY-MM-DD or RFC3339", i))
				}
			}
		}
	}

	// Get API client
	apiClient, err := s.getClient(ctx)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	// Convert DTO to API model
	apiRequest := mapTransactionUpdateRequestToAPI(&args.TransactionUpdateRequest)

	// Call the API
	resp, err := apiClient.UpdateTransactionWithResponse(ctx, args.ID, &client.UpdateTransactionParams{}, *apiRequest)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error updating transaction: %v", err))
	}

	// Check content type
	contentType := resp.HTTPResponse.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "application/vnd.api+json") {
		bodyPreview := string(resp.Body)
		if len(bodyPreview) > 500 {
			bodyPreview = bodyPreview[:500] + "..."
		}
		return newErrorResult(fmt.Sprintf("Error: Expected JSON response but got %s (status: %d)",
			contentType, resp.StatusCode()))
	}

	// Handle response codes
	switch resp.StatusCode() {
	case 200:
		var transactionSingle client.TransactionSingle

		if resp.ApplicationvndApiJSON200 != nil {
			transactionSingle = *resp.ApplicationvndApiJSON200
		} else {
			if len(resp.Body) == 0 {
				return newErrorResult(fmt.Sprintf("Error: empty response body for status %d", resp.StatusCode()))
			}

			if err := json.Unmarshal(resp.Body, &transactionSingle); err != nil {
				bodyPreview := string(resp.Body)
				if len(bodyPreview) > 200 {
					bodyPreview = bodyPreview[:200] + "..."
				}
				return newErrorResult(fmt.Sprintf("Error parsing response: %v (body: %s)", err, bodyPreview))
			}
		}

		transactionGroup := mapTransactionReadToTransactionGroup(&transactionSingle.Data)
		return newSuccessResult(transactionGroup)

	case 404:
		return newErrorResult("Error: Transaction not found")

	case 422:
		errorMsg := "Validation error"
		if resp.JSON422 != nil && resp.JSON422.Message != nil {
			errorMsg = *resp.JSON422.Message
		}
		return newErrorResult(fmt.Sprintf("Validation error: %s", errorMsg))

	case 400:
		return newErrorResult("Bad request: invalid data provided")

	default:
		return newErrorResult(fmt.Sprintf("Error: unexpected status code %d", resp.StatusCode()))
	}
}

// mapTransactionUpdateRequestToAPI converts DTO to API model for update
func mapTransactionUpdateRequestToAPI(req *TransactionUpdateRequest) *client.UpdateTransactionJSONRequestBody {
	apiReq := &client.UpdateTransactionJSONRequestBody{}

	// Set boolean fields only if true
	if req.ApplyRules {
		apiReq.ApplyRules = &req.ApplyRules
	}
	if req.FireWebhooks {
		apiReq.FireWebhooks = &req.FireWebhooks
	}
	if req.GroupTitle != "" {
		apiReq.GroupTitle = &req.GroupTitle
	}

	// Map transactions if provided
	if len(req.Transactions) > 0 {
		apiTransactions := make([]client.TransactionSplitUpdate, len(req.Transactions))

		for i, txn := range req.Transactions {
			apiTxn := client.TransactionSplitUpdate{}

			// Map type if provided
			if txn.Type != "" {
				txnType := client.TransactionTypeProperty(txn.Type)
				apiTxn.Type = &txnType
			}

			// Map date if provided
			if txn.Date != "" {
				parsedDate, err := time.Parse("2006-01-02", txn.Date)
				if err != nil {
					parsedDate, _ = time.Parse(time.RFC3339, txn.Date)
				}
				apiTxn.Date = &parsedDate
			}

			// Map amount if provided
			if txn.Amount != "" {
				apiTxn.Amount = &txn.Amount
			}

			// Map description if provided
			if txn.Description != "" {
				apiTxn.Description = &txn.Description
			}

			// Map account fields
			apiTxn.SourceId = txn.SourceId
			apiTxn.SourceName = txn.SourceName
			apiTxn.DestinationId = txn.DestinationId
			apiTxn.DestinationName = txn.DestinationName

			// Map categorization fields
			apiTxn.CategoryId = txn.CategoryId
			apiTxn.CategoryName = txn.CategoryName
			apiTxn.BudgetId = txn.BudgetId
			apiTxn.BudgetName = txn.BudgetName

			// Map tags
			if len(txn.Tags) > 0 {
				tags := make([]string, len(txn.Tags))
				copy(tags, txn.Tags)
				apiTxn.Tags = &tags
			}

			// Map currency fields
			apiTxn.CurrencyId = txn.CurrencyId
			apiTxn.CurrencyCode = txn.CurrencyCode
			apiTxn.ForeignAmount = txn.ForeignAmount
			apiTxn.ForeignCurrencyId = txn.ForeignCurrencyId
			apiTxn.ForeignCurrencyCode = txn.ForeignCurrencyCode

			// Map bill fields
			apiTxn.BillId = txn.BillId
			apiTxn.BillName = txn.BillName

			// Map notes and reconciled
			apiTxn.Notes = txn.Notes
			apiTxn.Reconciled = txn.Reconciled

			// Map order
			if txn.Order != nil {
				order := int32(*txn.Order)
				apiTxn.Order = &order
			}

			apiTransactions[i] = apiTxn
		}
		apiReq.Transactions = &apiTransactions
	}

	return apiReq
}
