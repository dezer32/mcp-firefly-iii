package mappers

import (
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
	"github.com/stretchr/testify/assert"
)

func TestMapArrayToList_Accounts(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		result := MapAccountArrayToAccountList(nil)
		assert.Nil(t, result)
	})

	t.Run("empty array returns empty list", func(t *testing.T) {
		input := &client.AccountArray{
			Data: []client.AccountRead{},
		}
		
		result := MapAccountArrayToAccountList(input)
		assert.NotNil(t, result)
		assert.Empty(t, result.Data)
	})

	t.Run("maps data correctly", func(t *testing.T) {
		active := true
		accountType := client.ShortAccountTypeProperty("asset")
		notes := "Test notes"
		
		input := &client.AccountArray{
			Data: []client.AccountRead{
				{
					Id:   "acc-1",
					Type: "accounts",
					Attributes: client.Account{
						Active:    &active,
						Name:      "Test Account 1",
						Notes:     &notes,
						Type:      accountType,
						CreatedAt: timePtr(time.Now()),
						UpdatedAt: timePtr(time.Now()),
					},
				},
				{
					Id:   "acc-2",
					Type: "accounts",
					Attributes: client.Account{
						Active:    nil, // Test nil handling
						Name:      "Test Account 2",
						Notes:     nil,
						Type:      accountType,
						CreatedAt: timePtr(time.Now()),
						UpdatedAt: timePtr(time.Now()),
					},
				},
			},
		}
		
		result := MapAccountArrayToAccountList(input)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		
		// Check first account
		assert.Equal(t, "acc-1", result.Data[0].Id)
		assert.True(t, result.Data[0].Active)
		assert.Equal(t, "Test Account 1", result.Data[0].Name)
		assert.Equal(t, &notes, result.Data[0].Notes)
		assert.Equal(t, "asset", result.Data[0].Type)
		
		// Check second account with nil values
		assert.Equal(t, "acc-2", result.Data[1].Id)
		assert.False(t, result.Data[1].Active) // nil should become false
		assert.Equal(t, "Test Account 2", result.Data[1].Name)
		assert.Nil(t, result.Data[1].Notes)
	})

	t.Run("maps pagination correctly", func(t *testing.T) {
		count := 10
		total := 50
		currentPage := 2
		perPage := 10
		totalPages := 5
		
		input := &client.AccountArray{
			Data: []client.AccountRead{},
			Meta: client.Meta{
				Pagination: &struct {
					Count       *int    `json:"count,omitempty"`
					CurrentPage *int    `json:"current_page,omitempty"`
					PerPage     *int    `json:"per_page,omitempty"`
					Total       *int    `json:"total,omitempty"`
					TotalPages  *int    `json:"total_pages,omitempty"`
				}{
					Count:       &count,
					Total:       &total,
					CurrentPage: &currentPage,
					PerPage:     &perPage,
					TotalPages:  &totalPages,
				},
			},
		}
		
		result := MapAccountArrayToAccountList(input)
		assert.NotNil(t, result)
		assert.Equal(t, 10, result.Pagination.Count)
		assert.Equal(t, 50, result.Pagination.Total)
		assert.Equal(t, 2, result.Pagination.CurrentPage)
		assert.Equal(t, 10, result.Pagination.PerPage)
		assert.Equal(t, 5, result.Pagination.TotalPages)
	})

	t.Run("handles nil pagination", func(t *testing.T) {
		input := &client.AccountArray{
			Data: []client.AccountRead{},
			Meta: client.Meta{
				Pagination: nil,
			},
		}
		
		result := MapAccountArrayToAccountList(input)
		assert.NotNil(t, result)
		assert.Equal(t, dto.Pagination{}, result.Pagination)
	})
}

func TestMapArrayToList_Categories(t *testing.T) {
	t.Run("maps categories correctly", func(t *testing.T) {
		notes := "Category notes"
		
		input := &client.CategoryArray{
			Data: []client.CategoryRead{
				{
					Id:   "cat-1",
					Type: "categories",
					Attributes: client.Category{
						Name:  "Category 1",
						Notes: &notes,
					},
				},
				{
					Id:   "cat-2",
					Type: "categories",
					Attributes: client.Category{
						Name:  "Category 2",
						Notes: nil,
					},
				},
			},
		}
		
		result := MapCategoryArrayToCategoryList(input)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		
		assert.Equal(t, "cat-1", result.Data[0].Id)
		assert.Equal(t, "Category 1", result.Data[0].Name)
		assert.Equal(t, &notes, result.Data[0].Notes)
		
		assert.Equal(t, "cat-2", result.Data[1].Id)
		assert.Equal(t, "Category 2", result.Data[1].Name)
		assert.Nil(t, result.Data[1].Notes)
	})
}

func TestMapArrayToList_Budgets(t *testing.T) {
	t.Run("maps budgets with spent data", func(t *testing.T) {
		active := true
		sum := "100.50"
		currencyCode := "USD"
		spent := []client.BudgetSpent{
			{
				Sum:          &sum,
				CurrencyCode: &currencyCode,
			},
		}
		
		input := &client.BudgetArray{
			Data: []client.BudgetRead{
				{
					Id:   "budget-1",
					Type: "budgets",
					Attributes: client.Budget{
						Active: &active,
						Name:   "Test Budget",
						Notes:  strPtr("Budget notes"),
						Spent:  &spent,
					},
				},
			},
		}
		
		result := MapBudgetArrayToBudgetList(input)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 1)
		
		assert.Equal(t, "budget-1", result.Data[0].Id)
		assert.True(t, result.Data[0].Active)
		assert.Equal(t, "Test Budget", result.Data[0].Name)
		assert.Equal(t, "100.50", result.Data[0].Spent.Sum)
		assert.Equal(t, "USD", result.Data[0].Spent.CurrencyCode)
	})

	t.Run("handles nil spent data", func(t *testing.T) {
		active := false
		
		input := &client.BudgetArray{
			Data: []client.BudgetRead{
				{
					Id:   "budget-2",
					Type: "budgets",
					Attributes: client.Budget{
						Active: &active,
						Name:   "Test Budget 2",
						Notes:  nil,
						Spent:  nil,
					},
				},
			},
		}
		
		result := MapBudgetArrayToBudgetList(input)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, 1)
		
		assert.Equal(t, "budget-2", result.Data[0].Id)
		assert.False(t, result.Data[0].Active)
		assert.Equal(t, "0", result.Data[0].Spent.Sum)
		assert.Equal(t, "", result.Data[0].Spent.CurrencyCode)
	})
}

func TestMapPaginationToDTO(t *testing.T) {
	t.Run("maps pagination correctly", func(t *testing.T) {
		count := 25
		total := 100
		currentPage := 3
		perPage := 25
		totalPages := 4
		
		input := struct {
			Count       *int
			Total       *int
			CurrentPage *int
			PerPage     *int
			TotalPages  *int
		}{
			Count:       &count,
			Total:       &total,
			CurrentPage: &currentPage,
			PerPage:     &perPage,
			TotalPages:  &totalPages,
		}
		
		result := MapPaginationToDTO(input)
		assert.Equal(t, 25, result.Count)
		assert.Equal(t, 100, result.Total)
		assert.Equal(t, 3, result.CurrentPage)
		assert.Equal(t, 25, result.PerPage)
		assert.Equal(t, 4, result.TotalPages)
	})

	t.Run("handles nil pagination", func(t *testing.T) {
		result := MapPaginationToDTO(nil)
		assert.Equal(t, dto.Pagination{}, result)
	})

	t.Run("handles nil fields", func(t *testing.T) {
		input := struct {
			Count       *int
			Total       *int
			CurrentPage *int
			PerPage     *int
			TotalPages  *int
		}{
			Count:       nil,
			Total:       nil,
			CurrentPage: nil,
			PerPage:     nil,
			TotalPages:  nil,
		}
		
		result := MapPaginationToDTO(input)
		assert.Equal(t, 0, result.Count)
		assert.Equal(t, 0, result.Total)
		assert.Equal(t, 0, result.CurrentPage)
		assert.Equal(t, 0, result.PerPage)
		assert.Equal(t, 0, result.TotalPages)
	})
}