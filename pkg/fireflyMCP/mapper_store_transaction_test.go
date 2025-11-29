package fireflyMCP

import (
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestMapTransactionStoreRequestToAPI(t *testing.T) {
	t.Run("minimal required fields", func(t *testing.T) {
		// Test with minimum required fields
		req := &TransactionStoreRequest{
			Transactions: []TransactionSplitRequest{
				{
					Type:        "withdrawal",
					Date:        "2024-01-15",
					Amount:      "50.00",
					Description: "Test transaction",
				},
			},
		}

		result := mapTransactionStoreRequestToAPI(req)
		assert.NotNil(t, result)
		assert.Len(t, result.Transactions, 1)
		assert.Equal(t, client.TransactionTypeProperty("withdrawal"), result.Transactions[0].Type)
		expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
		assert.Equal(t, expectedDate, result.Transactions[0].Date)
		assert.Equal(t, "50.00", result.Transactions[0].Amount)
		assert.Equal(t, "Test transaction", result.Transactions[0].Description)

		// Verify optional fields are nil
		assert.Nil(t, result.ErrorIfDuplicateHash)
		assert.Nil(t, result.ApplyRules)
		assert.Nil(t, result.FireWebhooks)
		assert.Nil(t, result.GroupTitle)
	})

	t.Run("all optional fields", func(t *testing.T) {
		// Test with all optional fields
		errorIfDup := true
		applyRules := false
		fireWebhooks := true
		groupTitle := "Split transaction"
		sourceId := "1"
		sourceName := "Checking Account"
		destId := "2"
		destName := "Groceries"
		categoryId := "3"
		categoryName := "Food"
		budgetId := "4"
		budgetName := "Monthly Budget"
		currencyId := "5"
		currencyCode := "USD"
		foreignAmount := "45.00"
		foreignCurrencyId := "6"
		foreignCurrencyCode := "EUR"
		billId := "7"
		billName := "Monthly Rent"
		piggyBankId := "8"
		piggyBankName := "Vacation Fund"
		notes := "Test notes"
		reconciled := true
		order := 1

		req := &TransactionStoreRequest{
			ErrorIfDuplicateHash: errorIfDup,
			ApplyRules:           applyRules,
			FireWebhooks:         fireWebhooks,
			GroupTitle:           groupTitle,
			Transactions: []TransactionSplitRequest{
				{
					Type:                "withdrawal",
					Date:                "2024-01-15T10:30:00Z",
					Amount:              "50.00",
					Description:         "Test transaction with all fields",
					SourceId:            &sourceId,
					SourceName:          &sourceName,
					DestinationId:       &destId,
					DestinationName:     &destName,
					CategoryId:          &categoryId,
					CategoryName:        &categoryName,
					BudgetId:            &budgetId,
					BudgetName:          &budgetName,
					Tags:                []string{"tag1", "tag2"},
					CurrencyId:          &currencyId,
					CurrencyCode:        &currencyCode,
					ForeignAmount:       &foreignAmount,
					ForeignCurrencyId:   &foreignCurrencyId,
					ForeignCurrencyCode: &foreignCurrencyCode,
					BillId:              &billId,
					BillName:            &billName,
					PiggyBankId:         &piggyBankId,
					PiggyBankName:       &piggyBankName,
					Notes:               &notes,
					Reconciled:          &reconciled,
					Order:               &order,
				},
			},
		}

		result := mapTransactionStoreRequestToAPI(req)
		assert.NotNil(t, result)

		// Verify top-level optional fields
		// Note: These are only set if they are true in the mapper
		assert.NotNil(t, result.ErrorIfDuplicateHash)
		assert.Equal(t, errorIfDup, *result.ErrorIfDuplicateHash)
		// applyRules is false so should be nil
		assert.Nil(t, result.ApplyRules)
		assert.NotNil(t, result.FireWebhooks)
		assert.Equal(t, fireWebhooks, *result.FireWebhooks)
		assert.NotNil(t, result.GroupTitle)
		assert.Equal(t, groupTitle, *result.GroupTitle)

		// Verify transaction fields
		assert.Len(t, result.Transactions, 1)
		txn := result.Transactions[0]

		assert.Equal(t, client.TransactionTypeProperty("withdrawal"), txn.Type)
		expectedDate, _ := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")
		assert.Equal(t, expectedDate, txn.Date)
		assert.Equal(t, "50.00", txn.Amount)
		assert.Equal(t, "Test transaction with all fields", txn.Description)

		// Verify all optional fields
		assert.NotNil(t, txn.SourceId)
		assert.Equal(t, sourceId, *txn.SourceId)
		assert.NotNil(t, txn.SourceName)
		assert.Equal(t, sourceName, *txn.SourceName)
		assert.NotNil(t, txn.DestinationId)
		assert.Equal(t, destId, *txn.DestinationId)
		assert.NotNil(t, txn.DestinationName)
		assert.Equal(t, destName, *txn.DestinationName)
		assert.NotNil(t, txn.CategoryId)
		assert.Equal(t, categoryId, *txn.CategoryId)
		assert.NotNil(t, txn.CategoryName)
		assert.Equal(t, categoryName, *txn.CategoryName)
		assert.NotNil(t, txn.BudgetId)
		assert.Equal(t, budgetId, *txn.BudgetId)
		assert.NotNil(t, txn.BudgetName)
		assert.Equal(t, budgetName, *txn.BudgetName)
		assert.NotNil(t, txn.Tags)
		assert.Len(t, *txn.Tags, 2)
		assert.Equal(t, "tag1", (*txn.Tags)[0])
		assert.Equal(t, "tag2", (*txn.Tags)[1])
		assert.NotNil(t, txn.CurrencyId)
		assert.Equal(t, currencyId, *txn.CurrencyId)
		assert.NotNil(t, txn.CurrencyCode)
		assert.Equal(t, currencyCode, *txn.CurrencyCode)
		assert.NotNil(t, txn.ForeignAmount)
		assert.Equal(t, foreignAmount, *txn.ForeignAmount)
		assert.NotNil(t, txn.ForeignCurrencyId)
		assert.Equal(t, foreignCurrencyId, *txn.ForeignCurrencyId)
		assert.NotNil(t, txn.ForeignCurrencyCode)
		assert.Equal(t, foreignCurrencyCode, *txn.ForeignCurrencyCode)
		assert.NotNil(t, txn.BillId)
		assert.Equal(t, billId, *txn.BillId)
		assert.NotNil(t, txn.BillName)
		assert.Equal(t, billName, *txn.BillName)
		assert.NotNil(t, txn.PiggyBankId)
		assert.Equal(t, int32(8), *txn.PiggyBankId)
		assert.NotNil(t, txn.PiggyBankName)
		assert.Equal(t, piggyBankName, *txn.PiggyBankName)
		assert.NotNil(t, txn.Notes)
		assert.Equal(t, notes, *txn.Notes)
		assert.NotNil(t, txn.Reconciled)
		assert.Equal(t, reconciled, *txn.Reconciled)
		assert.NotNil(t, txn.Order)
		assert.Equal(t, int32(1), *txn.Order)
	})

	t.Run("multiple transactions (split)", func(t *testing.T) {
		// Test with multiple transactions in a split
		groupTitle := "Split transaction"
		req := &TransactionStoreRequest{
			GroupTitle: groupTitle,
			Transactions: []TransactionSplitRequest{
				{
					Type:        "withdrawal",
					Date:        "2024-01-15",
					Amount:      "30.00",
					Description: "Part 1",
				},
				{
					Type:        "withdrawal",
					Date:        "2024-01-15",
					Amount:      "20.00",
					Description: "Part 2",
				},
				{
					Type:        "withdrawal",
					Date:        "2024-01-15",
					Amount:      "10.00",
					Description: "Part 3",
				},
			},
		}

		result := mapTransactionStoreRequestToAPI(req)
		assert.NotNil(t, result)
		assert.NotNil(t, result.GroupTitle)
		assert.Equal(t, groupTitle, *result.GroupTitle)
		assert.Len(t, result.Transactions, 3)

		// Verify each transaction
		expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
		for i, expectedAmount := range []string{"30.00", "20.00", "10.00"} {
			txn := result.Transactions[i]
			assert.Equal(t, client.TransactionTypeProperty("withdrawal"), txn.Type)
			assert.Equal(t, expectedDate, txn.Date)
			assert.Equal(t, expectedAmount, txn.Amount)
			assert.Equal(t, "Part "+string(rune('1'+i)), txn.Description)
		}
	})

	t.Run("nil values handling", func(t *testing.T) {
		// Test that nil values are handled correctly
		req := &TransactionStoreRequest{
			Transactions: []TransactionSplitRequest{
				{
					Type:        "deposit",
					Date:        "2024-01-15",
					Amount:      "100.00",
					Description: "Salary",
					// All optional fields are nil
					SourceId:            nil,
					SourceName:          nil,
					DestinationId:       nil,
					DestinationName:     nil,
					CategoryId:          nil,
					CategoryName:        nil,
					BudgetId:            nil,
					BudgetName:          nil,
					Tags:                nil,
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
				},
			},
		}

		result := mapTransactionStoreRequestToAPI(req)
		assert.NotNil(t, result)
		assert.Len(t, result.Transactions, 1)

		txn := result.Transactions[0]
		// Verify all optional fields are nil
		assert.Nil(t, txn.SourceId)
		assert.Nil(t, txn.SourceName)
		assert.Nil(t, txn.DestinationId)
		assert.Nil(t, txn.DestinationName)
		assert.Nil(t, txn.CategoryId)
		assert.Nil(t, txn.CategoryName)
		assert.Nil(t, txn.BudgetId)
		assert.Nil(t, txn.BudgetName)
		assert.Nil(t, txn.Tags)
		assert.Nil(t, txn.CurrencyId)
		assert.Nil(t, txn.CurrencyCode)
		assert.Nil(t, txn.ForeignAmount)
		assert.Nil(t, txn.ForeignCurrencyId)
		assert.Nil(t, txn.ForeignCurrencyCode)
		assert.Nil(t, txn.BillId)
		assert.Nil(t, txn.BillName)
		assert.Nil(t, txn.PiggyBankId)
		assert.Nil(t, txn.PiggyBankName)
		assert.Nil(t, txn.Notes)
		assert.Nil(t, txn.Reconciled)
		// Order is always set by the mapper (defaults to index when nil in request)
		assert.NotNil(t, txn.Order)
		assert.Equal(t, int32(0), *txn.Order)
	})

	t.Run("empty tags array", func(t *testing.T) {
		// Test that empty tags array is handled correctly
		req := &TransactionStoreRequest{
			Transactions: []TransactionSplitRequest{
				{
					Type:        "transfer",
					Date:        "2024-01-15",
					Amount:      "500.00",
					Description: "Transfer between accounts",
					Tags:        []string{}, // Empty tags array
				},
			},
		}

		result := mapTransactionStoreRequestToAPI(req)
		assert.NotNil(t, result)
		assert.Len(t, result.Transactions, 1)

		txn := result.Transactions[0]
		// Empty tags array should result in nil Tags field
		assert.Nil(t, txn.Tags)
	})

	t.Run("piggy bank ID conversion", func(t *testing.T) {
		// Test piggy bank ID string to int32 conversion
		testCases := []struct {
			input    string
			expected int32
		}{
			{"123", 123},
			{"0", 0},
			{"invalid", 0}, // Invalid strings should convert to 0
			{"", 0},        // Empty string should convert to 0
		}

		for _, tc := range testCases {
			piggyBankId := tc.input
			req := &TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:        "withdrawal",
						Date:        "2024-01-15",
						Amount:      "25.00",
						Description: "Test piggy bank ID conversion",
						PiggyBankId: &piggyBankId,
					},
				},
			}

			result := mapTransactionStoreRequestToAPI(req)
			assert.NotNil(t, result)
			assert.Len(t, result.Transactions, 1)

			txn := result.Transactions[0]
			assert.NotNil(t, txn.PiggyBankId)
			assert.Equal(t, tc.expected, *txn.PiggyBankId)
		}
	})
}
