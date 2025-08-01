package fireflyMCP

import (
	"testing"
	"time"

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

func TestMapTransactionArrayToTransactionList(t *testing.T) {
	// Test with nil input
	result := mapTransactionArrayToTransactionList(nil)
	assert.Nil(t, result)

	// Test with empty transaction array
	emptyArray := &client.TransactionArray{
		Data: []client.TransactionRead{},
		Meta: client.Meta{},
	}
	result = mapTransactionArrayToTransactionList(emptyArray)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data)

	// Test with sample data
	reconciled := true
	groupTitle := "Group Transaction"
	journalId := "123"
	currencyCode := "USD"
	destinationId := "456"
	sourceId := "789"
	tags := []string{"tag1", "tag2"}
	notes := "Test notes"
	count := 1
	total := 1
	currentPage := 1
	perPage := 10
	totalPages := 1

	transactionArray := &client.TransactionArray{
		Data: []client.TransactionRead{
			{
				Id: "1",
				Attributes: client.Transaction{
					GroupTitle: &groupTitle,
					Transactions: []client.TransactionSplit{
						{
							TransactionJournalId: &journalId,
							Amount:               "100.50",
							BillId:               nil,
							BillName:             nil,
							BudgetId:             nil,
							BudgetName:           nil,
							CategoryId:           nil,
							CategoryName:         nil,
							CurrencyCode:         &currencyCode,
							Date:                 time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							Description:          "Test transaction",
							DestinationId:        &destinationId,
							DestinationName:      &[]string{"Destination Account"}[0],
							DestinationType:      &[]client.AccountTypeProperty{client.AccountTypePropertyAssetAccount}[0],
							Notes:                &notes,
							Reconciled:           &reconciled,
							SourceId:             &sourceId,
							SourceName:           &[]string{"Source Account"}[0],
							Tags:                 &tags,
							Type:                 client.TransactionSplitType_Withdrawal,
						},
					},
				},
				Type: "transactions",
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

	result = mapTransactionArrayToTransactionList(transactionArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	group := result.Data[0]
	assert.Equal(t, "1", group.Id)
	assert.Equal(t, groupTitle, group.GroupTitle)
	assert.Len(t, group.Transactions, 1)

	transaction := group.Transactions[0]
	assert.Equal(t, journalId, transaction.Id)
	assert.Equal(t, "100.50", transaction.Amount)
	assert.Nil(t, transaction.BillId)
	assert.Nil(t, transaction.BillName)
	assert.Nil(t, transaction.BudgetId)
	assert.Nil(t, transaction.BudgetName)
	assert.Nil(t, transaction.CategoryId)
	assert.Nil(t, transaction.CategoryName)
	assert.Equal(t, currencyCode, transaction.CurrencyCode)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), transaction.Date)
	assert.Equal(t, "Test transaction", transaction.Description)
	assert.Equal(t, destinationId, transaction.DestinationId)
	assert.Equal(t, "Destination Account", transaction.DestinationName)
	assert.Equal(t, "Asset", transaction.DestinationType)
	assert.Equal(t, &notes, transaction.Notes)
	assert.True(t, transaction.Reconciled)
	assert.Equal(t, sourceId, transaction.SourceId)
	assert.Equal(t, "Source Account", transaction.SourceName)
	assert.Len(t, transaction.Tags, 2)
	assert.Equal(t, []string{"tag1", "tag2"}, transaction.Tags)
	assert.Equal(t, "withdrawal", transaction.Type)

	// Verify pagination
	assert.Equal(t, count, result.Pagination.Count)
	assert.Equal(t, total, result.Pagination.Total)
	assert.Equal(t, currentPage, result.Pagination.CurrentPage)
	assert.Equal(t, perPage, result.Pagination.PerPage)
	assert.Equal(t, totalPages, result.Pagination.TotalPages)
}

func TestMapTransactionArrayToTransactionList_MultipleTransactions(t *testing.T) {
	// Test with multiple transactions in a group
	reconciled := false
	groupTitle := "Split Transaction"
	journalId1 := "123"
	journalId2 := "124"
	currencyCode := "EUR"
	destinationId := "456"
	sourceId := "789"

	transactionArray := &client.TransactionArray{
		Data: []client.TransactionRead{
			{
				Id: "1",
				Attributes: client.Transaction{
					GroupTitle: &groupTitle,
					Transactions: []client.TransactionSplit{
						{
							TransactionJournalId: &journalId1,
							Amount:               "50.00",
							CurrencyCode:         &currencyCode,
							Date:                 time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							Description:          "Split 1",
							DestinationId:        &destinationId,
							DestinationName:      &[]string{"Expense Account 1"}[0],
							DestinationType:      &[]client.AccountTypeProperty{client.AccountTypePropertyExpenseAccount}[0],
							Reconciled:           &reconciled,
							SourceId:             &sourceId,
							SourceName:           &[]string{"Source Account"}[0],
							Type:                 client.TransactionSplitType_Withdrawal,
						},
						{
							TransactionJournalId: &journalId2,
							Amount:               "75.00",
							CurrencyCode:         &currencyCode,
							Date:                 time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							Description:          "Split 2",
							DestinationId:        &destinationId,
							DestinationName:      &[]string{"Expense Account 2"}[0],
							DestinationType:      &[]client.AccountTypeProperty{client.AccountTypePropertyExpenseAccount}[0],
							Reconciled:           &reconciled,
							SourceId:             &sourceId,
							SourceName:           &[]string{"Source Account"}[0],
							Type:                 client.TransactionSplitType_Withdrawal,
						},
					},
				},
				Type: "transactions",
			},
		},
		Meta: client.Meta{},
	}

	result := mapTransactionArrayToTransactionList(transactionArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	group := result.Data[0]
	assert.Equal(t, "1", group.Id)
	assert.Equal(t, groupTitle, group.GroupTitle)
	assert.Len(t, group.Transactions, 2)

	// Verify first transaction
	transaction1 := group.Transactions[0]
	assert.Equal(t, journalId1, transaction1.Id)
	assert.Equal(t, "50.00", transaction1.Amount)
	assert.Equal(t, "Split 1", transaction1.Description)
	assert.Equal(t, "Expense Account 1", transaction1.DestinationName)
	assert.False(t, transaction1.Reconciled)

	// Verify second transaction
	transaction2 := group.Transactions[1]
	assert.Equal(t, journalId2, transaction2.Id)
	assert.Equal(t, "75.00", transaction2.Amount)
	assert.Equal(t, "Split 2", transaction2.Description)
	assert.Equal(t, "Expense Account 2", transaction2.DestinationName)
	assert.False(t, transaction2.Reconciled)
}

func TestMapTransactionArrayToTransactionList_NilFields(t *testing.T) {
	// Test with nil optional fields
	transactionArray := &client.TransactionArray{
		Data: []client.TransactionRead{
			{
				Id: "1",
				Attributes: client.Transaction{
					GroupTitle: nil, // nil group title
					Transactions: []client.TransactionSplit{
						{
							TransactionJournalId: nil, // nil journal ID
							Amount:               "100.00",
							CurrencyCode:         nil, // nil currency code
							Date:                 time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							Description:          "Test transaction",
							DestinationId:        nil,
							DestinationName:      nil,
							DestinationType:      nil,
							Notes:                nil, // nil notes
							Reconciled:           nil, // nil reconciled
							SourceId:             nil,
							SourceName:           nil,
							Tags:                 nil, // nil tags
							Type:                 client.TransactionSplitType_Deposit,
						},
					},
				},
				Type: "transactions",
			},
		},
		Meta: client.Meta{},
	}

	result := mapTransactionArrayToTransactionList(transactionArray)

	// Verify the mapping with nil fields
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	group := result.Data[0]
	assert.Equal(t, "1", group.Id)
	assert.Equal(t, "", group.GroupTitle) // Should be empty string when nil
	assert.Len(t, group.Transactions, 1)

	transaction := group.Transactions[0]
	assert.Equal(t, "", transaction.Id) // Should be empty string when nil
	assert.Equal(t, "100.00", transaction.Amount)
	assert.Equal(t, "", transaction.CurrencyCode)     // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationId)    // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationName)  // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationType)  // Should be empty string when nil
	assert.Nil(t, transaction.Notes)                  // Should remain nil
	assert.False(t, transaction.Reconciled)           // Should be false when nil
	assert.Equal(t, "", transaction.SourceId)         // Should be empty string when nil
	assert.Equal(t, "", transaction.SourceName)       // Should be empty string when nil
	assert.Empty(t, transaction.Tags)                 // Should be empty slice when nil
	assert.Equal(t, "deposit", transaction.Type)
}

func TestGetAccountTypeValue(t *testing.T) {
	// Test with nil pointer
	assert.Equal(t, client.AccountTypeProperty(""), getAccountTypeValue(nil))

	// Test with valid pointer
	value := client.AccountTypeProperty_Asset
	assert.Equal(t, client.AccountTypeProperty_Asset, getAccountTypeValue(&value))
}
