package fireflyMCP

import (
	"testing"

	"github.com/dezer32/firefly-iii/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestMapBudgetArrayToBudgetList(t *testing.T) {
	// Test with nil input
	result := mapBudgetArrayToBudgetList(nil)
	assert.Nil(t, result)

	// Test with empty budget array
	emptyArray := &client.BudgetArray{
		Data: []client.BudgetRead{},
		Meta: client.Meta{},
	}
	result = mapBudgetArrayToBudgetList(emptyArray)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data)

	// Test with sample data
	active := true
	sum := "100.50"
	currencyCode := "USD"
	notes := "Test budget notes"
	count := 1
	total := 1
	currentPage := 1
	perPage := 10
	totalPages := 1

	budgetArray := &client.BudgetArray{
		Data: []client.BudgetRead{
			{
				Id: "1",
				Attributes: client.Budget{
					Active: &active,
					Name:   "Test Budget",
					Notes:  &notes,
					Spent: &[]client.BudgetSpent{
						{
							Sum:          &sum,
							CurrencyCode: &currencyCode,
						},
					},
				},
				Type: "budgets",
			},
		},
		Meta: client.Meta{
			Pagination: &struct {
				Count       *int `json:"count,omitempty"`
				CurrentPage *int `json:"current_page,omitempty"`
				PerPage     *int `json:"per_page,omitempty"`
				Total       *int `json:"total,omitempty"`
				TotalPages  *int `json:"total_pages,omitempty"`
			}{
				Count:       &count,
				Total:       &total,
				CurrentPage: &currentPage,
				PerPage:     &perPage,
				TotalPages:  &totalPages,
			},
		},
	}

	result = mapBudgetArrayToBudgetList(budgetArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	budget := result.Data[0]
	assert.Equal(t, "1", budget.Id)
	assert.True(t, budget.Active)
	assert.Equal(t, "Test Budget", budget.Name)
	assert.Equal(t, &notes, budget.Notes)

	// Verify spent data
	assert.Equal(t, sum, budget.Spent.Sum)
	assert.Equal(t, currencyCode, budget.Spent.CurrencyCode)

	// Verify pagination
	assert.Equal(t, count, result.Pagination.Count)
	assert.Equal(t, total, result.Pagination.Total)
	assert.Equal(t, currentPage, result.Pagination.CurrentPage)
	assert.Equal(t, perPage, result.Pagination.PerPage)
	assert.Equal(t, totalPages, result.Pagination.TotalPages)
}

func TestMapBudgetArrayToBudgetList_MultipleSpentItems(t *testing.T) {
	// Test that only the first spent item is used when multiple exist
	active := true
	firstSum := "100.50"
	secondSum := "200.75"
	firstCurrency := "USD"
	secondCurrency := "EUR"

	budgetArray := &client.BudgetArray{
		Data: []client.BudgetRead{
			{
				Id: "1",
				Attributes: client.Budget{
					Active: &active,
					Name:   "Test Budget",
					Spent: &[]client.BudgetSpent{
						{
							Sum:          &firstSum,
							CurrencyCode: &firstCurrency,
						},
						{
							Sum:          &secondSum,
							CurrencyCode: &secondCurrency,
						},
					},
				},
				Type: "budgets",
			},
		},
		Meta: client.Meta{},
	}

	result := mapBudgetArrayToBudgetList(budgetArray)

	// Verify only the first spent item is used
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	budget := result.Data[0]
	assert.Equal(t, firstSum, budget.Spent.Sum)
	assert.Equal(t, firstCurrency, budget.Spent.CurrencyCode)
	// Verify second item is NOT used
	assert.NotEqual(t, secondSum, budget.Spent.Sum)
	assert.NotEqual(t, secondCurrency, budget.Spent.CurrencyCode)
}

func TestGetStringValue(t *testing.T) {
	// Test with nil pointer
	assert.Equal(t, "", getStringValue(nil))

	// Test with valid pointer
	value := "test"
	assert.Equal(t, "test", getStringValue(&value))
}

func TestGetIntValue(t *testing.T) {
	// Test with nil pointer
	assert.Equal(t, 0, getIntValue(nil))

	// Test with valid pointer
	value := 42
	assert.Equal(t, 42, getIntValue(&value))
}

func TestMapAccountArrayToAccountList(t *testing.T) {
	// Test with nil input
	result := mapAccountArrayToAccountList(nil)
	assert.Nil(t, result)

	// Test with empty account array
	emptyArray := &client.AccountArray{
		Data: []client.AccountRead{},
		Meta: client.Meta{},
	}
	result = mapAccountArrayToAccountList(emptyArray)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data)

	// Test with sample data
	active := true
	notes := "Test account notes"
	count := 1
	total := 1
	currentPage := 1
	perPage := 10
	totalPages := 1

	accountArray := &client.AccountArray{
		Data: []client.AccountRead{
			{
				Id: "1",
				Attributes: client.Account{
					Active: &active,
					Name:   "Test Account",
					Notes:  &notes,
				},
				Type: "asset",
			},
		},
		Meta: client.Meta{
			Pagination: &struct {
				Count       *int `json:"count,omitempty"`
				CurrentPage *int `json:"current_page,omitempty"`
				PerPage     *int `json:"per_page,omitempty"`
				Total       *int `json:"total,omitempty"`
				TotalPages  *int `json:"total_pages,omitempty"`
			}{
				Count:       &count,
				Total:       &total,
				CurrentPage: &currentPage,
				PerPage:     &perPage,
				TotalPages:  &totalPages,
			},
		},
	}

	result = mapAccountArrayToAccountList(accountArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	account := result.Data[0]
	assert.Equal(t, "1", account.Id)
	assert.True(t, account.Active)
	assert.Equal(t, "Test Account", account.Name)
	assert.Equal(t, &notes, account.Notes)
	assert.Equal(t, "asset", account.Type)

	// Verify pagination
	assert.Equal(t, count, result.Pagination.Count)
	assert.Equal(t, total, result.Pagination.Total)
	assert.Equal(t, currentPage, result.Pagination.CurrentPage)
	assert.Equal(t, perPage, result.Pagination.PerPage)
	assert.Equal(t, totalPages, result.Pagination.TotalPages)
}

func TestMapAccountArrayToAccountList_InactiveAccount(t *testing.T) {
	// Test with inactive account
	active := false

	accountArray := &client.AccountArray{
		Data: []client.AccountRead{
			{
				Id: "2",
				Attributes: client.Account{
					Active: &active,
					Name:   "Inactive Account",
				},
				Type: "liability",
			},
		},
		Meta: client.Meta{},
	}

	result := mapAccountArrayToAccountList(accountArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	account := result.Data[0]
	assert.Equal(t, "2", account.Id)
	assert.False(t, account.Active)
	assert.Equal(t, "Inactive Account", account.Name)
	assert.Nil(t, account.Notes)
	assert.Equal(t, "liability", account.Type)
}

func TestMapAccountArrayToAccountList_NilActiveField(t *testing.T) {
	// Test with nil Active field (should default to false)
	accountArray := &client.AccountArray{
		Data: []client.AccountRead{
			{
				Id: "3",
				Attributes: client.Account{
					Active: nil, // nil pointer
					Name:   "Account with nil active",
				},
				Type: "expense",
			},
		},
		Meta: client.Meta{},
	}

	result := mapAccountArrayToAccountList(accountArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	account := result.Data[0]
	assert.Equal(t, "3", account.Id)
	assert.False(t, account.Active) // Should be false when Active is nil
	assert.Equal(t, "Account with nil active", account.Name)
	assert.Equal(t, "expense", account.Type)
}
