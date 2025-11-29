package fireflyMCP

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// handleStoreReceipt creates a group of transactions representing a shopping receipt
func (s *FireflyMCPServer) handleStoreReceipt(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ReceiptStoreRequest,
) (*mcp.CallToolResult, any, error) {
	// Step 1: Validate required fields
	if err := validateReceiptRequest(&args); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	// Step 2: Validate total amount if provided
	if args.TotalAmount != nil && *args.TotalAmount != "" {
		if err := validateTotalAmount(args.Items, *args.TotalAmount); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
	}

	// Step 3: Map receipt to TransactionStoreRequest
	txnRequest := mapReceiptToTransactionStoreRequest(&args)

	// Step 4: Delegate to existing store_transaction handler
	return s.handleStoreTransaction(ctx, req, *txnRequest)
}

// validateReceiptRequest validates the receipt request fields
func validateReceiptRequest(req *ReceiptStoreRequest) error {
	// Validate store_name
	if strings.TrimSpace(req.StoreName) == "" {
		return fmt.Errorf("store_name is required and must not be empty")
	}

	// Validate receipt_date format
	if req.ReceiptDate == "" {
		return fmt.Errorf("receipt_date is required")
	}
	if _, err := time.Parse("2006-01-02", req.ReceiptDate); err != nil {
		return fmt.Errorf("receipt_date must be in YYYY-MM-DD format")
	}

	// Validate items array
	if len(req.Items) == 0 {
		return fmt.Errorf("items array is required and must contain at least one item")
	}

	// Validate each item
	for i, item := range req.Items {
		if strings.TrimSpace(item.Amount) == "" {
			return fmt.Errorf("items[%d].amount is required", i)
		}
		if strings.TrimSpace(item.Description) == "" {
			return fmt.Errorf("items[%d].description is required", i)
		}
	}

	return nil
}

// validateTotalAmount checks if item amounts sum to the expected total
func validateTotalAmount(items []ReceiptItemRequest, expectedTotal string) error {
	// Parse expected total
	expected := new(big.Rat)
	if _, ok := expected.SetString(expectedTotal); !ok {
		return fmt.Errorf("invalid total_amount format: %s", expectedTotal)
	}

	// Sum all item amounts
	calculated := new(big.Rat)
	for i, item := range items {
		amount := new(big.Rat)
		if _, ok := amount.SetString(item.Amount); !ok {
			return fmt.Errorf("invalid amount format for items[%d]: %s", i, item.Amount)
		}
		calculated.Add(calculated, amount)
	}

	// Compare
	if expected.Cmp(calculated) != 0 {
		expectedFloat, _ := expected.Float64()
		calculatedFloat, _ := calculated.Float64()
		return fmt.Errorf(
			"total amount validation failed: expected %.2f, calculated %.2f",
			expectedFloat, calculatedFloat,
		)
	}

	return nil
}

// mapReceiptToTransactionStoreRequest converts a receipt to a transaction store request
func mapReceiptToTransactionStoreRequest(req *ReceiptStoreRequest) *TransactionStoreRequest {
	// Generate group title
	groupTitle := fmt.Sprintf("Receipt: %s - %s", req.StoreName, req.ReceiptDate)

	// Build transactions array
	transactions := make([]TransactionSplitRequest, len(req.Items))

	for i, item := range req.Items {
		txn := TransactionSplitRequest{
			Type:            "withdrawal",
			Date:            req.ReceiptDate,
			Amount:          item.Amount,
			Description:     item.Description,
			DestinationName: &req.StoreName,
		}

		// Source account (receipt-level)
		if req.SourceAccountId != nil {
			txn.SourceId = req.SourceAccountId
		}
		if req.SourceAccountName != nil {
			txn.SourceName = req.SourceAccountName
		}

		// Category: item-level overrides receipt default
		if item.CategoryId != nil {
			txn.CategoryId = item.CategoryId
		} else if item.CategoryName != nil {
			txn.CategoryName = item.CategoryName
		} else if req.DefaultCategoryId != nil {
			txn.CategoryId = req.DefaultCategoryId
		} else if req.DefaultCategoryName != nil {
			txn.CategoryName = req.DefaultCategoryName
		}

		// Budget (item-level only)
		if item.BudgetId != nil {
			txn.BudgetId = item.BudgetId
		}
		if item.BudgetName != nil {
			txn.BudgetName = item.BudgetName
		}

		// Currency (receipt-level)
		if req.CurrencyCode != nil {
			txn.CurrencyCode = req.CurrencyCode
		}

		// Tags: merge receipt-level and item-level, deduplicate
		txn.Tags = mergeTags(req.Tags, item.Tags)

		// Notes: item-level, or receipt notes on first item only
		if item.Notes != nil {
			txn.Notes = item.Notes
		} else if i == 0 && req.Notes != nil {
			txn.Notes = req.Notes
		}

		// Order
		order := i
		txn.Order = &order

		transactions[i] = txn
	}

	return &TransactionStoreRequest{
		GroupTitle:   groupTitle,
		ApplyRules:   req.ApplyRules,
		FireWebhooks: req.FireWebhooks,
		Transactions: transactions,
	}
}

// mergeTags merges two tag slices and removes duplicates
func mergeTags(receiptTags, itemTags []string) []string {
	if len(receiptTags) == 0 && len(itemTags) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var result []string

	for _, tag := range receiptTags {
		if !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}
	for _, tag := range itemTags {
		if !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}

	return result
}
