package fireflyMCP

import (
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	openapi_types "github.com/oapi-codegen/runtime/types"
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
					Type:   client.ShortAccountTypePropertyAsset,
				},
				Type: "accounts",
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
					Type:   client.ShortAccountTypePropertyLiability,
				},
				Type: "accounts",
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
					Type:   client.ShortAccountTypePropertyExpense,
				},
				Type: "accounts",
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

func TestMapAccountSingleToAccount(t *testing.T) {
	// Test with nil input
	result := mapAccountSingleToAccount(nil)
	assert.Nil(t, result)

	// Test with sample data
	active := true
	notes := "Test account notes"

	accountSingle := &client.AccountSingle{
		Data: client.AccountRead{
			Id: "1",
			Attributes: client.Account{
				Active: &active,
				Name:   "Test Account",
				Notes:  &notes,
				Type:   client.ShortAccountTypePropertyAsset,
			},
			Type: "accounts",
		},
	}

	result = mapAccountSingleToAccount(accountSingle)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Id)
	assert.True(t, result.Active)
	assert.Equal(t, "Test Account", result.Name)
	assert.Equal(t, &notes, result.Notes)
	assert.Equal(t, "asset", result.Type)
}

func TestMapAccountSingleToAccount_InactiveAccount(t *testing.T) {
	// Test with inactive account
	active := false

	accountSingle := &client.AccountSingle{
		Data: client.AccountRead{
			Id: "2",
			Attributes: client.Account{
				Active: &active,
				Name:   "Inactive Account",
				Type:   client.ShortAccountTypePropertyLiability,
			},
			Type: "accounts",
		},
	}

	result := mapAccountSingleToAccount(accountSingle)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Equal(t, "2", result.Id)
	assert.False(t, result.Active)
	assert.Equal(t, "Inactive Account", result.Name)
	assert.Nil(t, result.Notes)
	assert.Equal(t, "liability", result.Type)
}

func TestMapAccountSingleToAccount_NilActiveField(t *testing.T) {
	// Test with nil Active field (should default to false)
	accountSingle := &client.AccountSingle{
		Data: client.AccountRead{
			Id: "3",
			Attributes: client.Account{
				Active: nil, // nil pointer
				Name:   "Account with nil active",
				Type:   client.ShortAccountTypePropertyExpense,
			},
			Type: "accounts",
		},
	}

	result := mapAccountSingleToAccount(accountSingle)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Equal(t, "3", result.Id)
	assert.False(t, result.Active) // Should be false when Active is nil
	assert.Equal(t, "Account with nil active", result.Name)
	assert.Equal(t, "expense", result.Type)
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
							Type:                 client.Withdrawal,
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
	assert.Equal(t, "Asset account", transaction.DestinationType)
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
							Type:                 client.Withdrawal,
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
							Type:                 client.Withdrawal,
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
							Type:                 client.Deposit,
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
	assert.Equal(t, "", transaction.CurrencyCode)    // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationId)   // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationName) // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationType) // Should be empty string when nil
	assert.Nil(t, transaction.Notes)                 // Should remain nil
	assert.False(t, transaction.Reconciled)          // Should be false when nil
	assert.Equal(t, "", transaction.SourceId)        // Should be empty string when nil
	assert.Equal(t, "", transaction.SourceName)      // Should be empty string when nil
	assert.Empty(t, transaction.Tags)                // Should be empty slice when nil
	assert.Equal(t, "deposit", transaction.Type)
}

func TestGetAccountTypeValue(t *testing.T) {
	// Test with nil pointer
	assert.Equal(t, client.AccountTypeProperty(""), getAccountTypeValue(nil))

	// Test with valid pointer
	value := client.AccountTypePropertyAssetAccount
	assert.Equal(t, client.AccountTypePropertyAssetAccount, getAccountTypeValue(&value))
}

func TestMapTransactionReadToTransactionGroup(t *testing.T) {
	// Test with nil input
	result := mapTransactionReadToTransactionGroup(nil)
	assert.Nil(t, result)

	// Test with sample data
	reconciled := true
	groupTitle := "Group Transaction"
	journalId := "123"
	currencyCode := "USD"
	destinationId := "456"
	sourceId := "789"
	tags := []string{"tag1", "tag2"}
	notes := "Test notes"

	transactionRead := &client.TransactionRead{
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
					Type:                 client.Withdrawal,
				},
			},
		},
		Type: "transactions",
	}

	result = mapTransactionReadToTransactionGroup(transactionRead)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Id)
	assert.Equal(t, groupTitle, result.GroupTitle)
	assert.Len(t, result.Transactions, 1)

	transaction := result.Transactions[0]
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
	assert.Equal(t, "Asset account", transaction.DestinationType)
	assert.Equal(t, &notes, transaction.Notes)
	assert.True(t, transaction.Reconciled)
	assert.Equal(t, sourceId, transaction.SourceId)
	assert.Equal(t, "Source Account", transaction.SourceName)
	assert.Len(t, transaction.Tags, 2)
	assert.Equal(t, []string{"tag1", "tag2"}, transaction.Tags)
	assert.Equal(t, "withdrawal", transaction.Type)
}

func TestMapTransactionReadToTransactionGroup_MultipleTransactions(t *testing.T) {
	// Test with multiple transactions in a group
	reconciled := false
	groupTitle := "Split Transaction"
	journalId1 := "123"
	journalId2 := "124"
	currencyCode := "EUR"
	destinationId := "456"
	sourceId := "789"

	transactionRead := &client.TransactionRead{
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
					Type:                 client.Withdrawal,
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
					Type:                 client.Withdrawal,
				},
			},
		},
		Type: "transactions",
	}

	result := mapTransactionReadToTransactionGroup(transactionRead)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Id)
	assert.Equal(t, groupTitle, result.GroupTitle)
	assert.Len(t, result.Transactions, 2)

	// Verify first transaction
	transaction1 := result.Transactions[0]
	assert.Equal(t, journalId1, transaction1.Id)
	assert.Equal(t, "50.00", transaction1.Amount)
	assert.Equal(t, "Split 1", transaction1.Description)
	assert.Equal(t, "Expense Account 1", transaction1.DestinationName)
	assert.False(t, transaction1.Reconciled)

	// Verify second transaction
	transaction2 := result.Transactions[1]
	assert.Equal(t, journalId2, transaction2.Id)
	assert.Equal(t, "75.00", transaction2.Amount)
	assert.Equal(t, "Split 2", transaction2.Description)
	assert.Equal(t, "Expense Account 2", transaction2.DestinationName)
	assert.False(t, transaction2.Reconciled)
}

func TestMapTransactionReadToTransactionGroup_NilFields(t *testing.T) {
	// Test with nil optional fields
	transactionRead := &client.TransactionRead{
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
					Type:                 client.Deposit,
				},
			},
		},
		Type: "transactions",
	}

	result := mapTransactionReadToTransactionGroup(transactionRead)

	// Verify the mapping with nil fields
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Id)
	assert.Equal(t, "", result.GroupTitle) // Should be empty string when nil
	assert.Len(t, result.Transactions, 1)

	transaction := result.Transactions[0]
	assert.Equal(t, "", transaction.Id) // Should be empty string when nil
	assert.Equal(t, "100.00", transaction.Amount)
	assert.Equal(t, "", transaction.CurrencyCode)    // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationId)   // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationName) // Should be empty string when nil
	assert.Equal(t, "", transaction.DestinationType) // Should be empty string when nil
	assert.Nil(t, transaction.Notes)                 // Should remain nil
	assert.False(t, transaction.Reconciled)          // Should be false when nil
	assert.Equal(t, "", transaction.SourceId)        // Should be empty string when nil
	assert.Equal(t, "", transaction.SourceName)      // Should be empty string when nil
	assert.Empty(t, transaction.Tags)                // Should be empty slice when nil
	assert.Equal(t, "deposit", transaction.Type)
}

func TestMapTransactionReadToTransactionGroup_EmptyTransactions(t *testing.T) {
	// Test with empty transactions slice
	groupTitle := "Empty Group"

	transactionRead := &client.TransactionRead{
		Id: "1",
		Attributes: client.Transaction{
			GroupTitle:   &groupTitle,
			Transactions: []client.TransactionSplit{}, // empty slice
		},
		Type: "transactions",
	}

	result := mapTransactionReadToTransactionGroup(transactionRead)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Id)
	assert.Equal(t, groupTitle, result.GroupTitle)
	assert.Empty(t, result.Transactions)
}

func TestMapBasicSummaryToBasicSummaryList_Success(t *testing.T) {
	// Test with normal data
	key1 := "balance-in-EUR"
	title1 := "Balance (EUR)"
	currencyCode1 := "EUR"
	monetaryValue1 := "1234.56"

	key2 := "spent-in-USD"
	title2 := "Spent (USD)"
	currencyCode2 := "USD"
	monetaryValue2 := "-500.00"

	basicSummary := &client.BasicSummary{
		"balance": client.BasicSummaryEntry{
			Key:           &key1,
			Title:         &title1,
			CurrencyCode:  &currencyCode1,
			MonetaryValue: &monetaryValue1,
		},
		"spent": client.BasicSummaryEntry{
			Key:           &key2,
			Title:         &title2,
			CurrencyCode:  &currencyCode2,
			MonetaryValue: &monetaryValue2,
		},
	}

	result := mapBasicSummaryToBasicSummaryList(basicSummary)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)

	// Check that both entries are mapped correctly
	foundBalance := false
	foundSpent := false

	for _, summary := range result.Data {
		if summary.Key == key1 {
			foundBalance = true
			assert.Equal(t, title1, summary.Title)
			assert.Equal(t, currencyCode1, summary.CurrencyCode)
			assert.Equal(t, monetaryValue1, summary.MonetaryValue)
		}
		if summary.Key == key2 {
			foundSpent = true
			assert.Equal(t, title2, summary.Title)
			assert.Equal(t, currencyCode2, summary.CurrencyCode)
			assert.Equal(t, monetaryValue2, summary.MonetaryValue)
		}
	}

	assert.True(t, foundBalance, "Balance entry should be found")
	assert.True(t, foundSpent, "Spent entry should be found")
}

func TestMapBasicSummaryToBasicSummaryList_EmptyMap(t *testing.T) {
	// Test with empty map
	emptySummary := &client.BasicSummary{}
	result := mapBasicSummaryToBasicSummaryList(emptySummary)

	assert.NotNil(t, result)
	assert.Empty(t, result.Data)
}

func TestMapBasicSummaryToBasicSummaryList_NilInput(t *testing.T) {
	// Test with nil input
	result := mapBasicSummaryToBasicSummaryList(nil)

	assert.NotNil(t, result)
	assert.Empty(t, result.Data)
}

func TestMapBasicSummaryToBasicSummaryList_NilValues(t *testing.T) {
	// Test with nil values in entries
	basicSummary := &client.BasicSummary{
		"entry1": client.BasicSummaryEntry{
			Key:           nil,
			Title:         nil,
			CurrencyCode:  nil,
			MonetaryValue: nil,
		},
		"entry2": client.BasicSummaryEntry{
			Key:           nil,
			Title:         nil,
			CurrencyCode:  nil,
			MonetaryValue: nil,
		},
	}

	result := mapBasicSummaryToBasicSummaryList(basicSummary)

	// Verify the mapping handles nil values gracefully
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)

	// All values should be empty strings
	for _, summary := range result.Data {
		assert.Equal(t, "", summary.Key)
		assert.Equal(t, "", summary.Title)
		assert.Equal(t, "", summary.CurrencyCode)
		assert.Equal(t, "", summary.MonetaryValue)
	}
}

func TestMapBasicSummaryToBasicSummaryList_InvalidData(t *testing.T) {
	// Test with various edge cases
	key := "net-worth-in-EUR"
	title := "Net Worth (EUR)"
	currencyCode := "EUR"
	monetaryValue := "0.00"

	basicSummary := &client.BasicSummary{
		"net_worth": client.BasicSummaryEntry{
			Key:           &key,
			Title:         &title,
			CurrencyCode:  &currencyCode,
			MonetaryValue: &monetaryValue,
		},
	}

	result := mapBasicSummaryToBasicSummaryList(basicSummary)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	summary := result.Data[0]
	assert.Equal(t, key, summary.Key)
	assert.Equal(t, title, summary.Title)
	assert.Equal(t, currencyCode, summary.CurrencyCode)
	assert.Equal(t, monetaryValue, summary.MonetaryValue)
}

func TestMapInsightGroupToDTO(t *testing.T) {
	// Test with nil input
	result := mapInsightGroupToDTO(nil)
	assert.NotNil(t, result)
	assert.Empty(t, result.Entries)

	// Test with empty insight group
	emptyGroup := &client.InsightGroup{}
	result = mapInsightGroupToDTO(emptyGroup)
	assert.NotNil(t, result)
	assert.Empty(t, result.Entries)

	// Test with sample data
	id1 := "1"
	name1 := "Groceries"
	difference1 := "-150.50"
	currencyCode1 := "EUR"

	id2 := "2"
	name2 := "Transport"
	difference2 := "-75.00"
	currencyCode2 := "EUR"

	insightGroup := &client.InsightGroup{
		{
			Id:           &id1,
			Name:         &name1,
			Difference:   &difference1,
			CurrencyCode: &currencyCode1,
		},
		{
			Id:           &id2,
			Name:         &name2,
			Difference:   &difference2,
			CurrencyCode: &currencyCode2,
		},
	}

	result = mapInsightGroupToDTO(insightGroup)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Entries, 2)

	// Check first entry
	entry1 := result.Entries[0]
	assert.Equal(t, id1, entry1.Id)
	assert.Equal(t, name1, entry1.Name)
	assert.Equal(t, difference1, entry1.Amount)
	assert.Equal(t, currencyCode1, entry1.CurrencyCode)

	// Check second entry
	entry2 := result.Entries[1]
	assert.Equal(t, id2, entry2.Id)
	assert.Equal(t, name2, entry2.Name)
	assert.Equal(t, difference2, entry2.Amount)
	assert.Equal(t, currencyCode2, entry2.CurrencyCode)
}

func TestMapInsightGroupToDTO_NilValues(t *testing.T) {
	// Test with nil values in entries
	insightGroup := &client.InsightGroup{
		{
			Id:           nil,
			Name:         nil,
			Difference:   nil,
			CurrencyCode: nil,
		},
	}

	result := mapInsightGroupToDTO(insightGroup)

	// Verify the mapping handles nil values gracefully
	assert.NotNil(t, result)
	assert.Len(t, result.Entries, 1)

	entry := result.Entries[0]
	assert.Equal(t, "", entry.Id)
	assert.Equal(t, "", entry.Name)
	assert.Equal(t, "", entry.Amount)
	assert.Equal(t, "", entry.CurrencyCode)
}

func TestMapInsightTotalToDTO(t *testing.T) {
	// Test with nil input
	result := mapInsightTotalToDTO(nil)
	assert.NotNil(t, result)
	assert.Empty(t, result.Entries)

	// Test with empty insight total
	emptyTotal := &client.InsightTotal{}
	result = mapInsightTotalToDTO(emptyTotal)
	assert.NotNil(t, result)
	assert.Empty(t, result.Entries)

	// Test with sample data
	difference1 := "-500.00"
	currencyCode1 := "EUR"

	difference2 := "-250.50"
	currencyCode2 := "USD"

	insightTotal := &client.InsightTotal{
		{
			Difference:   &difference1,
			CurrencyCode: &currencyCode1,
		},
		{
			Difference:   &difference2,
			CurrencyCode: &currencyCode2,
		},
	}

	result = mapInsightTotalToDTO(insightTotal)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Entries, 2)

	// Check first entry
	entry1 := result.Entries[0]
	assert.Equal(t, difference1, entry1.Amount)
	assert.Equal(t, currencyCode1, entry1.CurrencyCode)

	// Check second entry
	entry2 := result.Entries[1]
	assert.Equal(t, difference2, entry2.Amount)
	assert.Equal(t, currencyCode2, entry2.CurrencyCode)
}

func TestMapInsightTotalToDTO_NilValues(t *testing.T) {
	// Test with nil values in entries
	insightTotal := &client.InsightTotal{
		{
			Difference:   nil,
			CurrencyCode: nil,
		},
	}

	result := mapInsightTotalToDTO(insightTotal)

	// Verify the mapping handles nil values gracefully
	assert.NotNil(t, result)
	assert.Len(t, result.Entries, 1)

	entry := result.Entries[0]
	assert.Equal(t, "", entry.Amount)
	assert.Equal(t, "", entry.CurrencyCode)
}

func TestMapInsightTotalToDTO_MultipleCurrencies(t *testing.T) {
	// Test with multiple currencies
	differenceEUR := "-1000.00"
	currencyCodeEUR := "EUR"

	differenceUSD := "-750.00"
	currencyCodeUSD := "USD"

	differenceGBP := "-500.00"
	currencyCodeGBP := "GBP"

	insightTotal := &client.InsightTotal{
		{
			Difference:   &differenceEUR,
			CurrencyCode: &currencyCodeEUR,
		},
		{
			Difference:   &differenceUSD,
			CurrencyCode: &currencyCodeUSD,
		},
		{
			Difference:   &differenceGBP,
			CurrencyCode: &currencyCodeGBP,
		},
	}

	result := mapInsightTotalToDTO(insightTotal)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Entries, 3)

	// Check all entries
	assert.Equal(t, differenceEUR, result.Entries[0].Amount)
	assert.Equal(t, currencyCodeEUR, result.Entries[0].CurrencyCode)

	assert.Equal(t, differenceUSD, result.Entries[1].Amount)
	assert.Equal(t, currencyCodeUSD, result.Entries[1].CurrencyCode)

	assert.Equal(t, differenceGBP, result.Entries[2].Amount)
	assert.Equal(t, currencyCodeGBP, result.Entries[2].CurrencyCode)
}

func TestMapBudgetLimitArrayToBudgetLimitList(t *testing.T) {
	// Test with nil input
	result := mapBudgetLimitArrayToBudgetLimitList(nil)
	assert.Nil(t, result)

	// Test with empty budget limit array
	emptyArray := &client.BudgetLimitArray{
		Data: []client.BudgetLimitRead{},
		Meta: client.Meta{},
	}
	result = mapBudgetLimitArrayToBudgetLimitList(emptyArray)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data)

	// Test with sample data
	amount := "500.00"
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
	currencyCode := "EUR"
	currencySymbol := "€"
	spent := "150.50"
	budgetId := "5"
	count := 1
	total := 1
	currentPage := 1
	perPage := 10
	totalPages := 1

	budgetLimitArray := &client.BudgetLimitArray{
		Data: []client.BudgetLimitRead{
			{
				Id: "10",
				Attributes: client.BudgetLimit{
					Amount:         amount,
					BudgetId:       &budgetId,
					CurrencyCode:   &currencyCode,
					CurrencySymbol: &currencySymbol,
					End:            end,
					Start:          start,
					Spent:          &spent,
				},
				Type: "budget_limits",
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

	result = mapBudgetLimitArrayToBudgetLimitList(budgetLimitArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	budgetLimit := result.Data[0]
	assert.Equal(t, "10", budgetLimit.Id)
	assert.Equal(t, amount, budgetLimit.Amount)
	assert.Equal(t, start, budgetLimit.Start)
	assert.Equal(t, end, budgetLimit.End)
	assert.Equal(t, "5", budgetLimit.BudgetId)
	assert.Equal(t, currencyCode, budgetLimit.CurrencyCode)
	assert.Equal(t, currencySymbol, budgetLimit.CurrencySymbol)

	// Verify spent data
	assert.Len(t, budgetLimit.Spent, 1)
	assert.Equal(t, spent, budgetLimit.Spent[0].Sum)
	assert.Equal(t, currencyCode, budgetLimit.Spent[0].CurrencyCode)
	assert.Equal(t, currencySymbol, budgetLimit.Spent[0].CurrencySymbol)

	// Verify pagination
	assert.Equal(t, count, result.Pagination.Count)
	assert.Equal(t, total, result.Pagination.Total)
	assert.Equal(t, currentPage, result.Pagination.CurrentPage)
	assert.Equal(t, perPage, result.Pagination.PerPage)
	assert.Equal(t, totalPages, result.Pagination.TotalPages)
}

func TestMapBudgetLimitArrayToBudgetLimitList_NilFields(t *testing.T) {
	// Test with nil optional fields
	amount := "1000.00"
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	budgetId := "10"

	budgetLimitArray := &client.BudgetLimitArray{
		Data: []client.BudgetLimitRead{
			{
				Id: "20",
				Attributes: client.BudgetLimit{
					Amount:         amount,
					BudgetId:       &budgetId,
					CurrencyCode:   nil,
					CurrencySymbol: nil,
					End:            end,
					Start:          start,
					Spent:          nil,
				},
				Type: "budget_limits",
			},
		},
		Meta: client.Meta{
			Pagination: nil,
		},
	}

	result := mapBudgetLimitArrayToBudgetLimitList(budgetLimitArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	budgetLimit := result.Data[0]
	assert.Equal(t, "20", budgetLimit.Id)
	assert.Equal(t, amount, budgetLimit.Amount)
	assert.Equal(t, "10", budgetLimit.BudgetId)
	assert.Equal(t, "", budgetLimit.CurrencyCode)
	assert.Equal(t, "", budgetLimit.CurrencySymbol)

	// Verify spent data is empty when nil
	assert.Empty(t, budgetLimit.Spent)

	// Verify pagination has zero values when nil
	assert.Equal(t, 0, result.Pagination.Count)
	assert.Equal(t, 0, result.Pagination.Total)
	assert.Equal(t, 0, result.Pagination.CurrentPage)
	assert.Equal(t, 0, result.Pagination.PerPage)
	assert.Equal(t, 0, result.Pagination.TotalPages)
}

func TestMapBudgetLimitArrayToBudgetLimitList_MultipleLimits(t *testing.T) {
	// Test with multiple budget limits
	amount1 := "300.00"
	amount2 := "600.00"
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
	currencyCode := "USD"
	currencySymbol := "$"
	spent1 := "100.00"
	spent2 := "450.75"
	budgetId1 := "15"
	budgetId2 := "16"
	count := 2
	total := 2
	currentPage := 1
	perPage := 10
	totalPages := 1

	budgetLimitArray := &client.BudgetLimitArray{
		Data: []client.BudgetLimitRead{
			{
				Id: "30",
				Attributes: client.BudgetLimit{
					Amount:         amount1,
					BudgetId:       &budgetId1,
					CurrencyCode:   &currencyCode,
					CurrencySymbol: &currencySymbol,
					End:            end,
					Start:          start,
					Spent:          &spent1,
				},
				Type: "budget_limits",
			},
			{
				Id: "31",
				Attributes: client.BudgetLimit{
					Amount:         amount2,
					BudgetId:       &budgetId2,
					CurrencyCode:   &currencyCode,
					CurrencySymbol: &currencySymbol,
					End:            end,
					Start:          start,
					Spent:          &spent2,
				},
				Type: "budget_limits",
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

	result := mapBudgetLimitArrayToBudgetLimitList(budgetLimitArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)

	// Verify first budget limit
	budgetLimit1 := result.Data[0]
	assert.Equal(t, "30", budgetLimit1.Id)
	assert.Equal(t, amount1, budgetLimit1.Amount)
	assert.Equal(t, "15", budgetLimit1.BudgetId)
	assert.Len(t, budgetLimit1.Spent, 1)
	assert.Equal(t, spent1, budgetLimit1.Spent[0].Sum)

	// Verify second budget limit
	budgetLimit2 := result.Data[1]
	assert.Equal(t, "31", budgetLimit2.Id)
	assert.Equal(t, amount2, budgetLimit2.Amount)
	assert.Equal(t, "16", budgetLimit2.BudgetId)
	assert.Len(t, budgetLimit2.Spent, 1)
	assert.Equal(t, spent2, budgetLimit2.Spent[0].Sum)

	// Verify pagination
	assert.Equal(t, count, result.Pagination.Count)
	assert.Equal(t, total, result.Pagination.Total)
}

func TestMapTagArrayToTagList(t *testing.T) {
	// Test with nil input
	result := mapTagArrayToTagList(nil)
	assert.Nil(t, result)

	// Test with empty tag array
	emptyArray := &client.TagArray{
		Data: []client.TagRead{},
		Meta: client.Meta{},
	}
	result = mapTagArrayToTagList(emptyArray)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data)

	// Test with sample data
	description := "Test tag description"
	tagDate := openapi_types.Date{Time: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)}
	createdAt := time.Date(2023, 11, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 11, 15, 15, 30, 0, 0, time.UTC)
	count := 2
	total := 2
	currentPage := 1
	perPage := 10
	totalPages := 1

	tagArray := &client.TagArray{
		Data: []client.TagRead{
			{
				Id: "1",
				Attributes: client.TagModel{
					Tag:         "groceries",
					Description: &description,
					Date:        &tagDate,
					CreatedAt:   &createdAt,
					UpdatedAt:   &updatedAt,
				},
			},
			{
				Id: "2",
				Attributes: client.TagModel{
					Tag: "vacation",
					// No optional fields for this tag
				},
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

	result = mapTagArrayToTagList(tagArray)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)

	// Verify first tag with all fields
	tag1 := result.Data[0]
	assert.Equal(t, "1", tag1.Id)
	assert.Equal(t, "groceries", tag1.Tag)
	assert.NotNil(t, tag1.Description)
	assert.Equal(t, description, *tag1.Description)
	assert.NotNil(t, tag1.Date)
	assert.Equal(t, "2023-12-01", *tag1.Date)
	assert.Equal(t, createdAt, tag1.CreatedAt)
	assert.Equal(t, updatedAt, tag1.UpdatedAt)

	// Verify second tag with minimal fields
	tag2 := result.Data[1]
	assert.Equal(t, "2", tag2.Id)
	assert.Equal(t, "vacation", tag2.Tag)
	assert.Nil(t, tag2.Description)
	assert.Nil(t, tag2.Date)

	// Verify pagination
	assert.Equal(t, count, result.Pagination.Count)
	assert.Equal(t, total, result.Pagination.Total)
	assert.Equal(t, currentPage, result.Pagination.CurrentPage)
	assert.Equal(t, perPage, result.Pagination.PerPage)
	assert.Equal(t, totalPages, result.Pagination.TotalPages)
}

func TestMapTagArrayToTagList_NilPagination(t *testing.T) {
	// Test with nil pagination
	tagArray := &client.TagArray{
		Data: []client.TagRead{
			{
				Id: "1",
				Attributes: client.TagModel{
					Tag: "test",
				},
			},
		},
		Meta: client.Meta{
			Pagination: nil,
		},
	}

	result := mapTagArrayToTagList(tagArray)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "test", result.Data[0].Tag)

	// Verify pagination has default zero values
	assert.Equal(t, 0, result.Pagination.Count)
	assert.Equal(t, 0, result.Pagination.Total)
}
