package fireflyMCP

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapTransactionUpdateRequestToAPI_EmptyRequest(t *testing.T) {
	req := &TransactionUpdateRequest{}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result)
	assert.Nil(t, result.ApplyRules)
	assert.Nil(t, result.FireWebhooks)
	assert.Nil(t, result.GroupTitle)
	assert.Nil(t, result.Transactions)
}

func TestMapTransactionUpdateRequestToAPI_WithBooleanFields(t *testing.T) {
	req := &TransactionUpdateRequest{
		ApplyRules:   true,
		FireWebhooks: true,
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.ApplyRules)
	assert.True(t, *result.ApplyRules)
	assert.NotNil(t, result.FireWebhooks)
	assert.True(t, *result.FireWebhooks)
}

func TestMapTransactionUpdateRequestToAPI_WithGroupTitle(t *testing.T) {
	req := &TransactionUpdateRequest{
		GroupTitle: "Test Group",
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.GroupTitle)
	assert.Equal(t, "Test Group", *result.GroupTitle)
}

func TestMapTransactionUpdateRequestToAPI_WithTransaction(t *testing.T) {
	amount := "100.00"
	sourceId := "1"
	destId := "2"
	categoryName := "Food"
	notes := "Test notes"

	req := &TransactionUpdateRequest{
		Transactions: []TransactionSplitRequest{
			{
				Type:            "withdrawal",
				Date:            "2024-01-15",
				Amount:          amount,
				Description:     "Test transaction",
				SourceId:        &sourceId,
				DestinationId:   &destId,
				CategoryName:    &categoryName,
				Notes:           &notes,
				Tags:            []string{"tag1", "tag2"},
			},
		},
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.Transactions)
	assert.Len(t, *result.Transactions, 1)

	txn := (*result.Transactions)[0]
	assert.NotNil(t, txn.Type)
	assert.Equal(t, "withdrawal", string(*txn.Type))
	assert.NotNil(t, txn.Date)
	assert.Equal(t, 2024, txn.Date.Year())
	assert.Equal(t, 1, int(txn.Date.Month()))
	assert.Equal(t, 15, txn.Date.Day())
	assert.NotNil(t, txn.Amount)
	assert.Equal(t, "100.00", *txn.Amount)
	assert.NotNil(t, txn.Description)
	assert.Equal(t, "Test transaction", *txn.Description)
	assert.Equal(t, &sourceId, txn.SourceId)
	assert.Equal(t, &destId, txn.DestinationId)
	assert.Equal(t, &categoryName, txn.CategoryName)
	assert.Equal(t, &notes, txn.Notes)
	assert.NotNil(t, txn.Tags)
	assert.Equal(t, []string{"tag1", "tag2"}, *txn.Tags)
}

func TestMapTransactionUpdateRequestToAPI_WithRFC3339Date(t *testing.T) {
	req := &TransactionUpdateRequest{
		Transactions: []TransactionSplitRequest{
			{
				Date: "2024-01-15T10:30:00Z",
			},
		},
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.Transactions)
	assert.Len(t, *result.Transactions, 1)

	txn := (*result.Transactions)[0]
	assert.NotNil(t, txn.Date)
	assert.Equal(t, 2024, txn.Date.Year())
	assert.Equal(t, 1, int(txn.Date.Month()))
	assert.Equal(t, 15, txn.Date.Day())
}

func TestMapTransactionUpdateRequestToAPI_WithOrder(t *testing.T) {
	order := 5
	req := &TransactionUpdateRequest{
		Transactions: []TransactionSplitRequest{
			{
				Order: &order,
			},
		},
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.Transactions)
	txn := (*result.Transactions)[0]
	assert.NotNil(t, txn.Order)
	assert.Equal(t, int32(5), *txn.Order)
}

func TestMapTransactionUpdateRequestToAPI_WithCurrencyFields(t *testing.T) {
	currencyId := "1"
	currencyCode := "USD"
	foreignAmount := "50.00"
	foreignCurrencyId := "2"
	foreignCurrencyCode := "EUR"

	req := &TransactionUpdateRequest{
		Transactions: []TransactionSplitRequest{
			{
				CurrencyId:          &currencyId,
				CurrencyCode:        &currencyCode,
				ForeignAmount:       &foreignAmount,
				ForeignCurrencyId:   &foreignCurrencyId,
				ForeignCurrencyCode: &foreignCurrencyCode,
			},
		},
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.Transactions)
	txn := (*result.Transactions)[0]
	assert.Equal(t, &currencyId, txn.CurrencyId)
	assert.Equal(t, &currencyCode, txn.CurrencyCode)
	assert.Equal(t, &foreignAmount, txn.ForeignAmount)
	assert.Equal(t, &foreignCurrencyId, txn.ForeignCurrencyId)
	assert.Equal(t, &foreignCurrencyCode, txn.ForeignCurrencyCode)
}

func TestMapTransactionUpdateRequestToAPI_WithBillFields(t *testing.T) {
	billId := "10"
	billName := "Monthly Rent"

	req := &TransactionUpdateRequest{
		Transactions: []TransactionSplitRequest{
			{
				BillId:   &billId,
				BillName: &billName,
			},
		},
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.Transactions)
	txn := (*result.Transactions)[0]
	assert.Equal(t, &billId, txn.BillId)
	assert.Equal(t, &billName, txn.BillName)
}

func TestMapTransactionUpdateRequestToAPI_WithReconciled(t *testing.T) {
	reconciled := true

	req := &TransactionUpdateRequest{
		Transactions: []TransactionSplitRequest{
			{
				Reconciled: &reconciled,
			},
		},
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.Transactions)
	txn := (*result.Transactions)[0]
	assert.NotNil(t, txn.Reconciled)
	assert.True(t, *txn.Reconciled)
}

func TestMapTransactionUpdateRequestToAPI_MultipleSplits(t *testing.T) {
	amount1 := "50.00"
	amount2 := "30.00"

	req := &TransactionUpdateRequest{
		GroupTitle: "Split Transaction",
		Transactions: []TransactionSplitRequest{
			{
				Amount:      amount1,
				Description: "First split",
			},
			{
				Amount:      amount2,
				Description: "Second split",
			},
		},
	}

	result := mapTransactionUpdateRequestToAPI(req)

	assert.NotNil(t, result.GroupTitle)
	assert.Equal(t, "Split Transaction", *result.GroupTitle)
	assert.NotNil(t, result.Transactions)
	assert.Len(t, *result.Transactions, 2)

	txn1 := (*result.Transactions)[0]
	assert.Equal(t, "50.00", *txn1.Amount)
	assert.Equal(t, "First split", *txn1.Description)

	txn2 := (*result.Transactions)[1]
	assert.Equal(t, "30.00", *txn2.Amount)
	assert.Equal(t, "Second split", *txn2.Description)
}
