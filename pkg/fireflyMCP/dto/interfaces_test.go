package dto

import (
	"testing"
	"time"
)

func TestMCPEntityImplementations(t *testing.T) {
	tests := []struct {
		name       string
		entity     MCPEntity
		expectedID string
		expectedName string
	}{
		{
			name: "Account implements MCPEntity",
			entity: &Account{
				Id:   "acc-123",
				Name: "Test Account",
			},
			expectedID:   "acc-123",
			expectedName: "Test Account",
		},
		{
			name: "Budget implements MCPEntity",
			entity: &Budget{
				Id:   "bud-456",
				Name: "Monthly Budget",
			},
			expectedID:   "bud-456",
			expectedName: "Monthly Budget",
		},
		{
			name: "Category implements MCPEntity",
			entity: &Category{
				Id:   "cat-789",
				Name: "Groceries",
			},
			expectedID:   "cat-789",
			expectedName: "Groceries",
		},
		{
			name: "Tag implements MCPEntity",
			entity: &Tag{
				Id:  "tag-101",
				Tag: "Important",
			},
			expectedID:   "tag-101",
			expectedName: "Important",
		},
		{
			name: "Bill implements MCPEntity",
			entity: &Bill{
				Id:   "bill-202",
				Name: "Internet Bill",
			},
			expectedID:   "bill-202",
			expectedName: "Internet Bill",
		},
		{
			name: "Transaction implements MCPEntity",
			entity: &Transaction{
				Id:          "trans-303",
				Description: "Coffee purchase",
			},
			expectedID:   "trans-303",
			expectedName: "Coffee purchase",
		},
		{
			name: "TransactionGroup implements MCPEntity",
			entity: &TransactionGroup{
				Id:         "tg-404",
				GroupTitle: "Weekly expenses",
			},
			expectedID:   "tg-404",
			expectedName: "Weekly expenses",
		},
		{
			name: "Recurrence implements MCPEntity",
			entity: &Recurrence{
				Id:    "rec-505",
				Title: "Monthly rent",
			},
			expectedID:   "rec-505",
			expectedName: "Monthly rent",
		},
		{
			name: "BasicSummary implements MCPEntity",
			entity: &BasicSummary{
				Key:   "total_spent",
				Title: "Total Spent",
			},
			expectedID:   "total_spent",
			expectedName: "Total Spent",
		},
		{
			name: "InsightCategoryEntry implements MCPEntity",
			entity: &InsightCategoryEntry{
				Id:   "ins-606",
				Name: "Food expenses",
			},
			expectedID:   "ins-606",
			expectedName: "Food expenses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entity.GetID(); got != tt.expectedID {
				t.Errorf("GetID() = %v, want %v", got, tt.expectedID)
			}
			if got := tt.entity.GetName(); got != tt.expectedName {
				t.Errorf("GetName() = %v, want %v", got, tt.expectedName)
			}
		})
	}
}

func TestPageableImplementations(t *testing.T) {
	tests := []struct {
		name         string
		pageable     Pageable
		expectedCount int
		expectedPage  int
		expectedTotal int
	}{
		{
			name: "AccountList implements Pageable",
			pageable: &AccountList{
				Data: []Account{
					{Id: "1", Name: "Account 1"},
					{Id: "2", Name: "Account 2"},
				},
				Pagination: Pagination{
					CurrentPage: 1,
					TotalPages:  5,
					Total:       10,
					PerPage:     2,
				},
			},
			expectedCount: 2,
			expectedPage:  1,
			expectedTotal: 10,
		},
		{
			name: "BudgetList implements Pageable",
			pageable: &BudgetList{
				Data: []Budget{
					{Id: "1", Name: "Budget 1"},
				},
				Pagination: Pagination{
					CurrentPage: 2,
					TotalPages:  3,
					Total:       5,
					PerPage:     2,
				},
			},
			expectedCount: 1,
			expectedPage:  2,
			expectedTotal: 5,
		},
		{
			name: "CategoryList implements Pageable",
			pageable: &CategoryList{
				Data: []Category{
					{Id: "1", Name: "Category 1"},
					{Id: "2", Name: "Category 2"},
					{Id: "3", Name: "Category 3"},
				},
				Pagination: Pagination{
					CurrentPage: 1,
					TotalPages:  1,
					Total:       3,
					PerPage:     10,
				},
			},
			expectedCount: 3,
			expectedPage:  1,
			expectedTotal: 3,
		},
		{
			name: "TagList implements Pageable",
			pageable: &TagList{
				Data: []Tag{
					{Id: "1", Tag: "Tag 1"},
				},
				Pagination: Pagination{
					CurrentPage: 1,
					TotalPages:  2,
					Total:       3,
					PerPage:     2,
				},
			},
			expectedCount: 1,
			expectedPage:  1,
			expectedTotal: 3,
		},
		{
			name: "BillList implements Pageable",
			pageable: &BillList{
				Data: []Bill{
					{Id: "1", Name: "Bill 1"},
					{Id: "2", Name: "Bill 2"},
				},
				Pagination: Pagination{
					CurrentPage: 1,
					TotalPages:  1,
					Total:       2,
					PerPage:     10,
				},
			},
			expectedCount: 2,
			expectedPage:  1,
			expectedTotal: 2,
		},
		{
			name: "RecurrenceList implements Pageable",
			pageable: &RecurrenceList{
				Data: []Recurrence{
					{Id: "1", Title: "Recurrence 1"},
				},
				Pagination: Pagination{
					CurrentPage: 3,
					TotalPages:  5,
					Total:       25,
					PerPage:     5,
				},
			},
			expectedCount: 1,
			expectedPage:  3,
			expectedTotal: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pageable.GetCount(); got != tt.expectedCount {
				t.Errorf("GetCount() = %v, want %v", got, tt.expectedCount)
			}
			
			pagination := tt.pageable.GetPagination()
			if pagination.CurrentPage != tt.expectedPage {
				t.Errorf("GetPagination().CurrentPage = %v, want %v", pagination.CurrentPage, tt.expectedPage)
			}
			if pagination.Total != tt.expectedTotal {
				t.Errorf("GetPagination().Total = %v, want %v", pagination.Total, tt.expectedTotal)
			}
			
			// Test that GetData returns correct number of items
			data := tt.pageable.GetData()
			if len(data) != tt.expectedCount {
				t.Errorf("GetData() returned %v items, want %v", len(data), tt.expectedCount)
			}
			
			// Verify each item implements MCPEntity
			for i, item := range data {
				if item == nil {
					t.Errorf("GetData()[%d] is nil", i)
					continue
				}
				// Test interface methods work
				_ = item.GetID()
				_ = item.GetName()
			}
		})
	}
}

func TestValidatableImplementations(t *testing.T) {
	tests := []struct {
		name      string
		validatable Validatable
		wantErr   bool
		errContains string
	}{
		// Valid cases
		{
			name: "Valid Account",
			validatable: &Account{
				Id:   "acc-1",
				Name: "Test Account",
				Type: "asset",
			},
			wantErr: false,
		},
		{
			name: "Valid Budget",
			validatable: &Budget{
				Id:   "bud-1",
				Name: "Monthly Budget",
			},
			wantErr: false,
		},
		{
			name: "Valid Category",
			validatable: &Category{
				Id:   "cat-1",
				Name: "Food",
			},
			wantErr: false,
		},
		{
			name: "Valid Tag",
			validatable: &Tag{
				Id:  "tag-1",
				Tag: "Important",
			},
			wantErr: false,
		},
		{
			name: "Valid Bill",
			validatable: &Bill{
				Id:         "bill-1",
				Name:       "Internet",
				AmountMin:  "50.00",
				AmountMax:  "60.00",
				RepeatFreq: "monthly",
			},
			wantErr: false,
		},
		{
			name: "Valid Transaction",
			validatable: &Transaction{
				Id:            "trans-1",
				Amount:        "100.00",
				Description:   "Purchase",
				Type:          "withdrawal",
				SourceId:      "acc-1",
				DestinationId: "acc-2",
			},
			wantErr: false,
		},
		{
			name: "Valid TransactionGroup",
			validatable: &TransactionGroup{
				Id: "tg-1",
				Transactions: []Transaction{
					{
						Id:            "trans-1",
						Amount:        "100.00",
						Description:   "Purchase",
						Type:          "withdrawal",
						SourceId:      "acc-1",
						DestinationId: "acc-2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid Recurrence",
			validatable: &Recurrence{
				Id:    "rec-1",
				Title: "Monthly Rent",
				Type:  "withdrawal",
				Repetitions: []RecurrenceRepetition{
					{Id: "rep-1", Type: "monthly"},
				},
				Transactions: []RecurrenceTransaction{
					{Id: "trans-1", Description: "Rent"},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid TransactionStoreRequest",
			validatable: &TransactionStoreRequest{
				Transactions: []TransactionSplitRequest{
					{
						Type:            "withdrawal",
						Date:            "2024-01-01",
						Amount:          "100.00",
						Description:     "Test",
						SourceId:        ptrString("acc-1"),
						DestinationName: ptrString("Store"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid Pagination",
			validatable: &Pagination{
				CurrentPage: 1,
				TotalPages:  5,
				PerPage:     10,
			},
			wantErr: false,
		},
		
		// Invalid cases
		{
			name: "Invalid Account - missing ID",
			validatable: &Account{
				Name: "Test Account",
				Type: "asset",
			},
			wantErr:     true,
			errContains: "ID cannot be empty",
		},
		{
			name: "Invalid Account - missing Name",
			validatable: &Account{
				Id:   "acc-1",
				Type: "asset",
			},
			wantErr:     true,
			errContains: "name cannot be empty",
		},
		{
			name: "Invalid Account - missing Type",
			validatable: &Account{
				Id:   "acc-1",
				Name: "Test Account",
			},
			wantErr:     true,
			errContains: "type cannot be empty",
		},
		{
			name: "Invalid Budget - missing ID",
			validatable: &Budget{
				Name: "Monthly Budget",
			},
			wantErr:     true,
			errContains: "ID cannot be empty",
		},
		{
			name: "Invalid Category - missing Name",
			validatable: &Category{
				Id: "cat-1",
			},
			wantErr:     true,
			errContains: "name cannot be empty",
		},
		{
			name: "Invalid Tag - missing Tag name",
			validatable: &Tag{
				Id: "tag-1",
			},
			wantErr:     true,
			errContains: "name cannot be empty",
		},
		{
			name: "Invalid Bill - missing RepeatFreq",
			validatable: &Bill{
				Id:        "bill-1",
				Name:      "Internet",
				AmountMin: "50.00",
				AmountMax: "60.00",
			},
			wantErr:     true,
			errContains: "repeat frequency cannot be empty",
		},
		{
			name: "Invalid Transaction - missing Amount",
			validatable: &Transaction{
				Id:            "trans-1",
				Description:   "Purchase",
				Type:          "withdrawal",
				SourceId:      "acc-1",
				DestinationId: "acc-2",
			},
			wantErr:     true,
			errContains: "amount cannot be empty",
		},
		{
			name: "Invalid TransactionGroup - no transactions",
			validatable: &TransactionGroup{
				Id: "tg-1",
			},
			wantErr:     true,
			errContains: "must contain at least one transaction",
		},
		{
			name: "Invalid TransactionGroup - invalid transaction",
			validatable: &TransactionGroup{
				Id: "tg-1",
				Transactions: []Transaction{
					{
						Id: "trans-1",
						// Missing required fields
					},
				},
			},
			wantErr:     true,
			errContains: "transaction 0 validation failed",
		},
		{
			name: "Invalid Recurrence - no repetitions",
			validatable: &Recurrence{
				Id:    "rec-1",
				Title: "Monthly Rent",
				Type:  "withdrawal",
				Transactions: []RecurrenceTransaction{
					{Id: "trans-1"},
				},
			},
			wantErr:     true,
			errContains: "must have at least one repetition",
		},
		{
			name: "Invalid TransactionSplitRequest - invalid type",
			validatable: &TransactionSplitRequest{
				Type:            "invalid",
				Date:            "2024-01-01",
				Amount:          "100.00",
				Description:     "Test",
				SourceId:        ptrString("acc-1"),
				DestinationName: ptrString("Store"),
			},
			wantErr:     true,
			errContains: "must be one of",
		},
		{
			name: "Invalid TransactionSplitRequest - missing source",
			validatable: &TransactionSplitRequest{
				Type:          "withdrawal",
				Date:          "2024-01-01",
				Amount:        "100.00",
				Description:   "Test",
				DestinationId: ptrString("acc-2"),
			},
			wantErr:     true,
			errContains: "source_id or source_name must be provided",
		},
		{
			name: "Invalid Pagination - invalid PerPage",
			validatable: &Pagination{
				CurrentPage: 1,
				TotalPages:  5,
				PerPage:     0,
			},
			wantErr:     true,
			errContains: "per_page must be greater than 0",
		},
		{
			name: "Invalid Pagination - current page > total pages",
			validatable: &Pagination{
				CurrentPage: 10,
				TotalPages:  5,
				PerPage:     10,
			},
			wantErr:     true,
			errContains: "current_page cannot be greater than total_pages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validatable.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestGetDataReturnsCorrectTypes(t *testing.T) {
	// Test that GetData() returns the correct entity types
	t.Run("AccountList returns Account entities", func(t *testing.T) {
		list := &AccountList{
			Data: []Account{
				{Id: "1", Name: "Account 1"},
				{Id: "2", Name: "Account 2"},
			},
		}
		
		data := list.GetData()
		for i, entity := range data {
			// Type assertion to verify it's an Account
			acc, ok := entity.(*Account)
			if !ok {
				t.Errorf("GetData()[%d] is not *Account", i)
			}
			if acc.Id != list.Data[i].Id {
				t.Errorf("GetData()[%d].Id = %v, want %v", i, acc.Id, list.Data[i].Id)
			}
		}
	})
}

// Helper functions
func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrInt(i int) *int {
	return &i
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || len(substr) < len(s) && containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 1; i < len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}