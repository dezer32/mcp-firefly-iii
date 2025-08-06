package dto

import (
	"strings"
	"testing"
	"time"
)

func TestNewDTOFactory(t *testing.T) {
	tests := []struct {
		name            string
		options         []FactoryOption
		wantValidation  bool
		wantNilHandling bool
	}{
		{
			name:            "Default factory",
			options:         nil,
			wantValidation:  true,
			wantNilHandling: true,
		},
		{
			name:            "Factory with validation disabled",
			options:         []FactoryOption{WithValidation(false)},
			wantValidation:  false,
			wantNilHandling: true,
		},
		{
			name:            "Factory with nil handling disabled",
			options:         []FactoryOption{WithNilHandling(false)},
			wantValidation:  true,
			wantNilHandling: false,
		},
		{
			name:            "Factory with both options",
			options:         []FactoryOption{WithValidation(false), WithNilHandling(false)},
			wantValidation:  false,
			wantNilHandling: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewDTOFactory(tt.options...)
			if factory == nil {
				t.Fatal("NewDTOFactory returned nil")
			}
			
			// Test internal state by trying to create invalid DTOs
			// If validation is disabled, invalid DTOs should be created
			_, err := factory.CreateAccount("", "", "", false, nil)
			if tt.wantValidation && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.wantValidation && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name          string
		factory       DTOFactory
		id            string
		accountName   string
		accountType   string
		active        bool
		notes         *string
		wantErr       bool
		errContains   string
	}{
		{
			name:        "Valid account",
			factory:     NewDTOFactory(WithValidation(true)),
			id:          "acc-1",
			accountName: "Test Account",
			accountType: "asset",
			active:      true,
			notes:       ptrString("Some notes"),
			wantErr:     false,
		},
		{
			name:        "Valid account with nil notes",
			factory:     NewDTOFactory(WithValidation(true)),
			id:          "acc-2",
			accountName: "Test Account 2",
			accountType: "expense",
			active:      false,
			notes:       nil,
			wantErr:     false,
		},
		{
			name:        "Invalid account - missing ID",
			factory:     NewDTOFactory(WithValidation(true)),
			id:          "",
			accountName: "Test Account",
			accountType: "asset",
			active:      true,
			notes:       nil,
			wantErr:     true,
			errContains: "ID cannot be empty",
		},
		{
			name:        "Invalid account - missing name",
			factory:     NewDTOFactory(WithValidation(true)),
			id:          "acc-3",
			accountName: "",
			accountType: "asset",
			active:      true,
			notes:       nil,
			wantErr:     true,
			errContains: "name cannot be empty",
		},
		{
			name:        "Invalid account bypassed with validation disabled",
			factory:     NewDTOFactory(WithValidation(false)),
			id:          "",
			accountName: "",
			accountType: "",
			active:      false,
			notes:       nil,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := tt.factory.CreateAccount(tt.id, tt.accountName, tt.accountType, tt.active, tt.notes)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CreateAccount() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			
			if !tt.wantErr && account != nil {
				if account.Id != tt.id {
					t.Errorf("Account.Id = %v, want %v", account.Id, tt.id)
				}
				if account.Name != tt.accountName {
					t.Errorf("Account.Name = %v, want %v", account.Name, tt.accountName)
				}
				if account.Type != tt.accountType {
					t.Errorf("Account.Type = %v, want %v", account.Type, tt.accountType)
				}
				if account.Active != tt.active {
					t.Errorf("Account.Active = %v, want %v", account.Active, tt.active)
				}
			}
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	factory := NewDTOFactory(WithValidation(true))
	now := time.Now()
	
	tests := []struct {
		name         string
		id           string
		amount       string
		description  string
		sourceId     string
		sourceName   string
		destId       string
		destName     string
		txType       string
		currencyCode string
		destType     string
		date         time.Time
		wantErr      bool
		errContains  string
	}{
		{
			name:         "Valid transaction",
			id:           "trans-1",
			amount:       "100.00",
			description:  "Test transaction",
			sourceId:     "acc-1",
			sourceName:   "Source Account",
			destId:       "acc-2",
			destName:     "Dest Account",
			txType:       "withdrawal",
			currencyCode: "USD",
			destType:     "expense",
			date:         now,
			wantErr:      false,
		},
		{
			name:         "Invalid transaction - missing amount",
			id:           "trans-2",
			amount:       "",
			description:  "Test transaction",
			sourceId:     "acc-1",
			sourceName:   "Source Account",
			destId:       "acc-2",
			destName:     "Dest Account",
			txType:       "withdrawal",
			currencyCode: "USD",
			destType:     "expense",
			date:         now,
			wantErr:      true,
			errContains:  "amount cannot be empty",
		},
		{
			name:         "Invalid transaction - missing description",
			id:           "trans-3",
			amount:       "100.00",
			description:  "",
			sourceId:     "acc-1",
			sourceName:   "Source Account",
			destId:       "acc-2",
			destName:     "Dest Account",
			txType:       "withdrawal",
			currencyCode: "USD",
			destType:     "expense",
			date:         now,
			wantErr:      true,
			errContains:  "description cannot be empty",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := factory.CreateTransaction(
				tt.id, tt.amount, tt.description,
				tt.sourceId, tt.sourceName,
				tt.destId, tt.destName,
				tt.txType, tt.currencyCode, tt.destType,
				tt.date,
			)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CreateTransaction() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			
			if !tt.wantErr && transaction != nil {
				if transaction.Id != tt.id {
					t.Errorf("Transaction.Id = %v, want %v", transaction.Id, tt.id)
				}
				if transaction.Amount != tt.amount {
					t.Errorf("Transaction.Amount = %v, want %v", transaction.Amount, tt.amount)
				}
				if transaction.Description != tt.description {
					t.Errorf("Transaction.Description = %v, want %v", transaction.Description, tt.description)
				}
				if len(transaction.Tags) != 0 {
					t.Errorf("Transaction.Tags should be empty array, got %v", transaction.Tags)
				}
			}
		})
	}
}

func TestCreateTransactionGroup(t *testing.T) {
	factory := NewDTOFactory(WithValidation(true))
	
	validTransaction := Transaction{
		Id:            "trans-1",
		Amount:        "100.00",
		Description:   "Test",
		Type:          "withdrawal",
		SourceId:      "acc-1",
		DestinationId: "acc-2",
	}
	
	tests := []struct {
		name         string
		id           string
		groupTitle   string
		transactions []Transaction
		wantErr      bool
		errContains  string
	}{
		{
			name:         "Valid transaction group",
			id:           "tg-1",
			groupTitle:   "Test Group",
			transactions: []Transaction{validTransaction},
			wantErr:      false,
		},
		{
			name:         "Invalid transaction group - no transactions",
			id:           "tg-2",
			groupTitle:   "Test Group",
			transactions: []Transaction{},
			wantErr:      true,
			errContains:  "must contain at least one transaction",
		},
		{
			name:         "Invalid transaction group - nil transactions",
			id:           "tg-3",
			groupTitle:   "Test Group",
			transactions: nil,
			wantErr:      true,
			errContains:  "must contain at least one transaction",
		},
		{
			name:       "Invalid transaction group - invalid transaction",
			id:         "tg-4",
			groupTitle: "Test Group",
			transactions: []Transaction{
				{Id: "", Amount: "", Description: ""},
			},
			wantErr:     true,
			errContains: "transaction 0 validation failed",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group, err := factory.CreateTransactionGroup(tt.id, tt.groupTitle, tt.transactions)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransactionGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CreateTransactionGroup() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			
			if !tt.wantErr && group != nil {
				if group.Id != tt.id {
					t.Errorf("TransactionGroup.Id = %v, want %v", group.Id, tt.id)
				}
				if group.GroupTitle != tt.groupTitle {
					t.Errorf("TransactionGroup.GroupTitle = %v, want %v", group.GroupTitle, tt.groupTitle)
				}
			}
		})
	}
}

func TestCreateLists(t *testing.T) {
	factory := NewDTOFactory(WithValidation(true))
	
	validPagination := Pagination{
		CurrentPage: 1,
		TotalPages:  1,
		PerPage:     10,
		Total:       2,
		Count:       2,
	}
	
	invalidPagination := Pagination{
		CurrentPage: 0, // Invalid
		TotalPages:  1,
		PerPage:     10,
	}
	
	t.Run("CreateAccountList", func(t *testing.T) {
		accounts := []Account{
			{Id: "1", Name: "Account 1", Type: "asset", Active: true},
			{Id: "2", Name: "Account 2", Type: "expense", Active: false},
		}
		
		list, err := factory.CreateAccountList(accounts, validPagination)
		if err != nil {
			t.Errorf("CreateAccountList() unexpected error: %v", err)
		}
		if list == nil {
			t.Fatal("CreateAccountList() returned nil")
		}
		if len(list.Data) != 2 {
			t.Errorf("CreateAccountList() Data length = %v, want 2", len(list.Data))
		}
		
		// Test with invalid pagination
		_, err = factory.CreateAccountList(accounts, invalidPagination)
		if err == nil {
			t.Error("CreateAccountList() expected error for invalid pagination")
		}
		
		// Test with invalid account
		invalidAccounts := []Account{
			{Id: "", Name: "Invalid", Type: "asset"},
		}
		_, err = factory.CreateAccountList(invalidAccounts, validPagination)
		if err == nil {
			t.Error("CreateAccountList() expected error for invalid account")
		}
		
		// Test with nil accounts
		list, err = factory.CreateAccountList(nil, validPagination)
		if err != nil {
			t.Errorf("CreateAccountList() unexpected error with nil accounts: %v", err)
		}
		if len(list.Data) != 0 {
			t.Errorf("CreateAccountList() should initialize empty array for nil accounts")
		}
	})
	
	t.Run("CreateBudgetList", func(t *testing.T) {
		budgets := []Budget{
			{Id: "1", Name: "Budget 1", Active: true},
			{Id: "2", Name: "Budget 2", Active: false},
		}
		
		list, err := factory.CreateBudgetList(budgets, validPagination)
		if err != nil {
			t.Errorf("CreateBudgetList() unexpected error: %v", err)
		}
		if list == nil {
			t.Fatal("CreateBudgetList() returned nil")
		}
		if len(list.Data) != 2 {
			t.Errorf("CreateBudgetList() Data length = %v, want 2", len(list.Data))
		}
		
		// Test with nil budgets
		list, err = factory.CreateBudgetList(nil, validPagination)
		if err != nil {
			t.Errorf("CreateBudgetList() unexpected error with nil budgets: %v", err)
		}
		if len(list.Data) != 0 {
			t.Errorf("CreateBudgetList() should initialize empty array for nil budgets")
		}
	})
	
	t.Run("CreateCategoryList", func(t *testing.T) {
		categories := []Category{
			{Id: "1", Name: "Category 1"},
			{Id: "2", Name: "Category 2"},
		}
		
		list, err := factory.CreateCategoryList(categories, validPagination)
		if err != nil {
			t.Errorf("CreateCategoryList() unexpected error: %v", err)
		}
		if list == nil {
			t.Fatal("CreateCategoryList() returned nil")
		}
		if len(list.Data) != 2 {
			t.Errorf("CreateCategoryList() Data length = %v, want 2", len(list.Data))
		}
	})
	
	t.Run("CreateTagList", func(t *testing.T) {
		tags := []Tag{
			{Id: "1", Tag: "Tag 1"},
			{Id: "2", Tag: "Tag 2"},
		}
		
		list, err := factory.CreateTagList(tags, validPagination)
		if err != nil {
			t.Errorf("CreateTagList() unexpected error: %v", err)
		}
		if list == nil {
			t.Fatal("CreateTagList() returned nil")
		}
		if len(list.Data) != 2 {
			t.Errorf("CreateTagList() Data length = %v, want 2", len(list.Data))
		}
	})
}

func TestCreateBill(t *testing.T) {
	factory := NewDTOFactory(WithValidation(true))
	now := time.Now()
	
	tests := []struct {
		name        string
		id          string
		billName    string
		amountMin   string
		amountMax   string
		repeatFreq  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "Valid bill",
			id:         "bill-1",
			billName:   "Internet",
			amountMin:  "50.00",
			amountMax:  "60.00",
			repeatFreq: "monthly",
			wantErr:    false,
		},
		{
			name:        "Invalid bill - missing name",
			id:          "bill-2",
			billName:    "",
			amountMin:   "50.00",
			amountMax:   "60.00",
			repeatFreq:  "monthly",
			wantErr:     true,
			errContains: "name cannot be empty",
		},
		{
			name:        "Invalid bill - missing repeat frequency",
			id:          "bill-3",
			billName:    "Internet",
			amountMin:   "50.00",
			amountMax:   "60.00",
			repeatFreq:  "",
			wantErr:     true,
			errContains: "repeat frequency cannot be empty",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bill, err := factory.CreateBill(
				tt.id, tt.billName, tt.amountMin, tt.amountMax,
				tt.repeatFreq, "USD", now, 0, true, nil, nil, nil,
			)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CreateBill() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			
			if !tt.wantErr && bill != nil {
				if bill.Id != tt.id {
					t.Errorf("Bill.Id = %v, want %v", bill.Id, tt.id)
				}
				if bill.Name != tt.billName {
					t.Errorf("Bill.Name = %v, want %v", bill.Name, tt.billName)
				}
			}
		})
	}
}

func TestCreateInsightEntries(t *testing.T) {
	factory := NewDTOFactory(WithValidation(false)) // No validation implemented for insights yet
	
	t.Run("CreateInsightCategoryEntry", func(t *testing.T) {
		entry, err := factory.CreateInsightCategoryEntry("ins-1", "Food", "150.00", "USD")
		if err != nil {
			t.Errorf("CreateInsightCategoryEntry() unexpected error: %v", err)
		}
		if entry == nil {
			t.Fatal("CreateInsightCategoryEntry() returned nil")
		}
		if entry.Id != "ins-1" {
			t.Errorf("InsightCategoryEntry.Id = %v, want ins-1", entry.Id)
		}
		if entry.Name != "Food" {
			t.Errorf("InsightCategoryEntry.Name = %v, want Food", entry.Name)
		}
		if entry.Amount != "150.00" {
			t.Errorf("InsightCategoryEntry.Amount = %v, want 150.00", entry.Amount)
		}
		if entry.CurrencyCode != "USD" {
			t.Errorf("InsightCategoryEntry.CurrencyCode = %v, want USD", entry.CurrencyCode)
		}
	})
	
	t.Run("CreateInsightTotalEntry", func(t *testing.T) {
		entry, err := factory.CreateInsightTotalEntry("500.00", "EUR")
		if err != nil {
			t.Errorf("CreateInsightTotalEntry() unexpected error: %v", err)
		}
		if entry == nil {
			t.Fatal("CreateInsightTotalEntry() returned nil")
		}
		if entry.Amount != "500.00" {
			t.Errorf("InsightTotalEntry.Amount = %v, want 500.00", entry.Amount)
		}
		if entry.CurrencyCode != "EUR" {
			t.Errorf("InsightTotalEntry.CurrencyCode = %v, want EUR", entry.CurrencyCode)
		}
	})
}

func TestCreateBasicSummary(t *testing.T) {
	factory := NewDTOFactory(WithValidation(false)) // No validation implemented for summaries yet
	
	summary, err := factory.CreateBasicSummary("total_spent", "Total Spent", "USD", "1250.50")
	if err != nil {
		t.Errorf("CreateBasicSummary() unexpected error: %v", err)
	}
	if summary == nil {
		t.Fatal("CreateBasicSummary() returned nil")
	}
	if summary.Key != "total_spent" {
		t.Errorf("BasicSummary.Key = %v, want total_spent", summary.Key)
	}
	if summary.Title != "Total Spent" {
		t.Errorf("BasicSummary.Title = %v, want Total Spent", summary.Title)
	}
	if summary.CurrencyCode != "USD" {
		t.Errorf("BasicSummary.CurrencyCode = %v, want USD", summary.CurrencyCode)
	}
	if summary.MonetaryValue != "1250.50" {
		t.Errorf("BasicSummary.MonetaryValue = %v, want 1250.50", summary.MonetaryValue)
	}
}

func TestCreateRecurrence(t *testing.T) {
	factory := NewDTOFactory(WithValidation(false)) // Validation skipped for recurrences
	now := time.Now()
	
	recurrence, err := factory.CreateRecurrence(
		"rec-1", "withdrawal", "Monthly Rent", "Apartment rent",
		now, true, false,
	)
	
	if err != nil {
		t.Errorf("CreateRecurrence() unexpected error: %v", err)
	}
	if recurrence == nil {
		t.Fatal("CreateRecurrence() returned nil")
	}
	if recurrence.Id != "rec-1" {
		t.Errorf("Recurrence.Id = %v, want rec-1", recurrence.Id)
	}
	if recurrence.Type != "withdrawal" {
		t.Errorf("Recurrence.Type = %v, want withdrawal", recurrence.Type)
	}
	if recurrence.Title != "Monthly Rent" {
		t.Errorf("Recurrence.Title = %v, want Monthly Rent", recurrence.Title)
	}
	if recurrence.Description != "Apartment rent" {
		t.Errorf("Recurrence.Description = %v, want Apartment rent", recurrence.Description)
	}
	if !recurrence.Active {
		t.Error("Recurrence.Active should be true")
	}
	if recurrence.ApplyRules {
		t.Error("Recurrence.ApplyRules should be false")
	}
	if recurrence.Repetitions == nil {
		t.Error("Recurrence.Repetitions should be initialized as empty array")
	}
	if recurrence.Transactions == nil {
		t.Error("Recurrence.Transactions should be initialized as empty array")
	}
}

func TestDefaultFactory(t *testing.T) {
	// Test that DefaultFactory is properly initialized
	if DefaultFactory == nil {
		t.Fatal("DefaultFactory is nil")
	}
	
	// Test that it has validation enabled by default
	_, err := DefaultFactory.CreateAccount("", "", "", false, nil)
	if err == nil {
		t.Error("DefaultFactory should have validation enabled by default")
	}
}

func TestFactoryWithNilHandling(t *testing.T) {
	factoryWithNil := NewDTOFactory(WithNilHandling(true))
	factoryWithoutNil := NewDTOFactory(WithNilHandling(false))
	
	notes := ptrString("test notes")
	
	// Both should handle non-nil values the same way
	acc1, _ := factoryWithNil.CreateAccount("1", "Test", "asset", true, notes)
	acc2, _ := factoryWithoutNil.CreateAccount("1", "Test", "asset", true, notes)
	
	if acc1.Notes != acc2.Notes {
		t.Error("Both factories should handle non-nil values the same way")
	}
	
	// Test with nil values
	acc3, _ := factoryWithNil.CreateAccount("2", "Test2", "asset", true, nil)
	acc4, _ := factoryWithoutNil.CreateAccount("2", "Test2", "asset", true, nil)
	
	// Both should preserve nil
	if acc3.Notes != nil {
		t.Error("Factory with nil handling should preserve nil")
	}
	if acc4.Notes != nil {
		t.Error("Factory without nil handling should preserve nil")
	}
}