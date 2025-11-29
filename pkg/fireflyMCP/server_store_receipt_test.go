package fireflyMCP

import (
	"context"
	"testing"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

// Ensure mcp is used (for TextContent type assertion)
var _ = mcp.TextContent{}

func TestValidateReceiptRequest(t *testing.T) {
	tests := []struct {
		name          string
		req           *ReceiptStoreRequest
		expectedError string
	}{
		{
			name: "Valid request",
			req: &ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "2024-11-29",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectedError: "",
		},
		{
			name: "Empty store_name",
			req: &ReceiptStoreRequest{
				StoreName:   "",
				ReceiptDate: "2024-11-29",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectedError: "store_name is required and must not be empty",
		},
		{
			name: "Whitespace store_name",
			req: &ReceiptStoreRequest{
				StoreName:   "   ",
				ReceiptDate: "2024-11-29",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectedError: "store_name is required and must not be empty",
		},
		{
			name: "Empty receipt_date",
			req: &ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectedError: "receipt_date is required",
		},
		{
			name: "Invalid receipt_date format",
			req: &ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "29-11-2024",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectedError: "receipt_date must be in YYYY-MM-DD format",
		},
		{
			name: "Empty items array",
			req: &ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "2024-11-29",
				Items:       []ReceiptItemRequest{},
			},
			expectedError: "items array is required and must contain at least one item",
		},
		{
			name: "Item missing amount",
			req: &ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "2024-11-29",
				Items: []ReceiptItemRequest{
					{Amount: "", Description: "Item 1"},
				},
			},
			expectedError: "items[0].amount is required",
		},
		{
			name: "Item missing description",
			req: &ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "2024-11-29",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: ""},
				},
			},
			expectedError: "items[0].description is required",
		},
		{
			name: "Second item invalid",
			req: &ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "2024-11-29",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
					{Amount: "20.00", Description: ""},
				},
			},
			expectedError: "items[1].description is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReceiptRequest(tt.req)
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestValidateTotalAmount(t *testing.T) {
	tests := []struct {
		name          string
		items         []ReceiptItemRequest
		totalAmount   string
		expectedError string
	}{
		{
			name: "Valid total",
			items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
				{Amount: "20.00", Description: "Item 2"},
			},
			totalAmount:   "30.00",
			expectedError: "",
		},
		{
			name: "Valid total with decimals",
			items: []ReceiptItemRequest{
				{Amount: "10.50", Description: "Item 1"},
				{Amount: "20.25", Description: "Item 2"},
			},
			totalAmount:   "30.75",
			expectedError: "",
		},
		{
			name: "Total mismatch",
			items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
				{Amount: "20.00", Description: "Item 2"},
			},
			totalAmount:   "50.00",
			expectedError: "total amount validation failed",
		},
		{
			name: "Invalid total format",
			items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
			},
			totalAmount:   "invalid",
			expectedError: "invalid total_amount format",
		},
		{
			name: "Invalid item amount format",
			items: []ReceiptItemRequest{
				{Amount: "abc", Description: "Item 1"},
			},
			totalAmount:   "10.00",
			expectedError: "invalid amount format for items[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTotalAmount(tt.items, tt.totalAmount)
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestMapReceiptToTransactionStoreRequest(t *testing.T) {
	sourceId := "1"
	sourceName := "Test Account"
	categoryId := "10"
	categoryName := "Groceries"
	defaultCategoryName := "Default Category"
	currencyCode := "USD"
	notes := "Receipt notes"
	itemNotes := "Item specific notes"
	budgetId := "5"

	t.Run("Basic mapping", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:   "Test Store",
			ReceiptDate: "2024-11-29",
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
				{Amount: "20.00", Description: "Item 2"},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, "Receipt: Test Store - 2024-11-29", result.GroupTitle)
		assert.Len(t, result.Transactions, 2)
		assert.Equal(t, "withdrawal", result.Transactions[0].Type)
		assert.Equal(t, "withdrawal", result.Transactions[1].Type)
		assert.Equal(t, "2024-11-29", result.Transactions[0].Date)
		assert.Equal(t, "Test Store", *result.Transactions[0].DestinationName)
	})

	t.Run("Source account mapping", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:         "Test Store",
			ReceiptDate:       "2024-11-29",
			SourceAccountId:   &sourceId,
			SourceAccountName: &sourceName,
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, &sourceId, result.Transactions[0].SourceId)
		assert.Equal(t, &sourceName, result.Transactions[0].SourceName)
	})

	t.Run("Category fallback to default", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:           "Test Store",
			ReceiptDate:         "2024-11-29",
			DefaultCategoryName: &defaultCategoryName,
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, &defaultCategoryName, result.Transactions[0].CategoryName)
	})

	t.Run("Item category overrides default", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:           "Test Store",
			ReceiptDate:         "2024-11-29",
			DefaultCategoryName: &defaultCategoryName,
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1", CategoryName: &categoryName},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, &categoryName, result.Transactions[0].CategoryName)
	})

	t.Run("Item category ID takes precedence over name", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:   "Test Store",
			ReceiptDate: "2024-11-29",
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1", CategoryId: &categoryId, CategoryName: &categoryName},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, &categoryId, result.Transactions[0].CategoryId)
		assert.Nil(t, result.Transactions[0].CategoryName)
	})

	t.Run("Currency code propagation", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:    "Test Store",
			ReceiptDate:  "2024-11-29",
			CurrencyCode: &currencyCode,
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
				{Amount: "20.00", Description: "Item 2"},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, &currencyCode, result.Transactions[0].CurrencyCode)
		assert.Equal(t, &currencyCode, result.Transactions[1].CurrencyCode)
	})

	t.Run("Receipt notes on first item only", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:   "Test Store",
			ReceiptDate: "2024-11-29",
			Notes:       &notes,
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
				{Amount: "20.00", Description: "Item 2"},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, &notes, result.Transactions[0].Notes)
		assert.Nil(t, result.Transactions[1].Notes)
	})

	t.Run("Item notes override receipt notes", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:   "Test Store",
			ReceiptDate: "2024-11-29",
			Notes:       &notes,
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1", Notes: &itemNotes},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, &itemNotes, result.Transactions[0].Notes)
	})

	t.Run("Budget mapping", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:   "Test Store",
			ReceiptDate: "2024-11-29",
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1", BudgetId: &budgetId},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, &budgetId, result.Transactions[0].BudgetId)
	})

	t.Run("Order assignment", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:   "Test Store",
			ReceiptDate: "2024-11-29",
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
				{Amount: "20.00", Description: "Item 2"},
				{Amount: "30.00", Description: "Item 3"},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.Equal(t, 0, *result.Transactions[0].Order)
		assert.Equal(t, 1, *result.Transactions[1].Order)
		assert.Equal(t, 2, *result.Transactions[2].Order)
	})

	t.Run("ApplyRules and FireWebhooks flags", func(t *testing.T) {
		req := &ReceiptStoreRequest{
			StoreName:    "Test Store",
			ReceiptDate:  "2024-11-29",
			ApplyRules:   true,
			FireWebhooks: true,
			Items: []ReceiptItemRequest{
				{Amount: "10.00", Description: "Item 1"},
			},
		}

		result := mapReceiptToTransactionStoreRequest(req)

		assert.True(t, result.ApplyRules)
		assert.True(t, result.FireWebhooks)
	})
}

func TestMergeTags(t *testing.T) {
	tests := []struct {
		name        string
		receiptTags []string
		itemTags    []string
		expected    []string
	}{
		{
			name:        "Both empty",
			receiptTags: nil,
			itemTags:    nil,
			expected:    nil,
		},
		{
			name:        "Only receipt tags",
			receiptTags: []string{"tag1", "tag2"},
			itemTags:    nil,
			expected:    []string{"tag1", "tag2"},
		},
		{
			name:        "Only item tags",
			receiptTags: nil,
			itemTags:    []string{"tag1", "tag2"},
			expected:    []string{"tag1", "tag2"},
		},
		{
			name:        "Merge without duplicates",
			receiptTags: []string{"tag1", "tag2"},
			itemTags:    []string{"tag3", "tag4"},
			expected:    []string{"tag1", "tag2", "tag3", "tag4"},
		},
		{
			name:        "Merge with duplicates",
			receiptTags: []string{"tag1", "tag2"},
			itemTags:    []string{"tag2", "tag3"},
			expected:    []string{"tag1", "tag2", "tag3"},
		},
		{
			name:        "All duplicates",
			receiptTags: []string{"tag1", "tag2"},
			itemTags:    []string{"tag1", "tag2"},
			expected:    []string{"tag1", "tag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeTags(tt.receiptTags, tt.itemTags)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleStoreReceipt_Validation(t *testing.T) {
	tests := []struct {
		name          string
		args          ReceiptStoreRequest
		expectedError string
	}{
		{
			name: "Empty store_name",
			args: ReceiptStoreRequest{
				StoreName:   "",
				ReceiptDate: "2024-11-29",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectedError: "store_name is required",
		},
		{
			name: "Invalid date format",
			args: ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "29-11-2024",
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
				},
			},
			expectedError: "receipt_date must be in YYYY-MM-DD format",
		},
		{
			name: "Empty items",
			args: ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "2024-11-29",
				Items:       []ReceiptItemRequest{},
			},
			expectedError: "items array is required",
		},
		{
			name: "Total amount mismatch",
			args: ReceiptStoreRequest{
				StoreName:   "Test Store",
				ReceiptDate: "2024-11-29",
				TotalAmount: stringPtr("100.00"),
				Items: []ReceiptItemRequest{
					{Amount: "10.00", Description: "Item 1"},
					{Amount: "20.00", Description: "Item 2"},
				},
			},
			expectedError: "total amount validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &FireflyMCPServer{
				client: &client.ClientWithResponses{},
			}

			result, _, err := server.handleStoreReceipt(context.Background(), nil, tt.args)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, result.IsError)
			assert.Len(t, result.Content, 1)

			textContent, ok := result.Content[0].(*mcp.TextContent)
			assert.True(t, ok)
			assert.Contains(t, textContent.Text, tt.expectedError)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
