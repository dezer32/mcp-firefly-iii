package mappers

import (
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"
)

// Old implementation for comparison
func mapAccountArrayToAccountListOld(accountArray *client.AccountArray) *dto.AccountList {
	if accountArray == nil {
		return nil
	}

	accountList := &dto.AccountList{
		Data: make([]dto.Account, len(accountArray.Data)),
	}

	// Map account data
	for i, accountRead := range accountArray.Data {
		account := dto.Account{
			Id:     accountRead.Id,
			Active: accountRead.Attributes.Active != nil && *accountRead.Attributes.Active,
			Name:   accountRead.Attributes.Name,
			Notes:  accountRead.Attributes.Notes,
			Type:   string(accountRead.Attributes.Type),
		}

		accountList.Data[i] = account
	}

	// Map pagination
	if accountArray.Meta.Pagination != nil {
		pagination := accountArray.Meta.Pagination
		accountList.Pagination = dto.Pagination{
			Count:       GetIntValue(pagination.Count),
			Total:       GetIntValue(pagination.Total),
			CurrentPage: GetIntValue(pagination.CurrentPage),
			PerPage:     GetIntValue(pagination.PerPage),
			TotalPages:  GetIntValue(pagination.TotalPages),
		}
	}

	return accountList
}

// createTestAccountArray creates test data for benchmarking
func createTestAccountArray(size int) *client.AccountArray {
	data := make([]client.AccountRead, size)
	active := true
	accountType := client.ShortAccountTypeProperty("asset")
	
	for i := 0; i < size; i++ {
		data[i] = client.AccountRead{
			Id:   "acc-" + string(rune(i)),
			Type: "accounts",
			Attributes: client.Account{
				Active:    &active,
				Name:      "Test Account " + string(rune(i)),
				Notes:     strPtr("Test notes"),
				Type:      accountType,
				CreatedAt: timePtr(time.Now()),
				UpdatedAt: timePtr(time.Now()),
			},
		}
	}

	count := size
	total := size
	currentPage := 1
	perPage := size
	totalPages := 1

	return &client.AccountArray{
		Data: data,
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
}

// Helper functions for creating test data
func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// BenchmarkMapAccountArrayOld benchmarks the old implementation
func BenchmarkMapAccountArrayOld(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run("Size"+string(rune(size)), func(b *testing.B) {
			testData := createTestAccountArray(size)
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				_ = mapAccountArrayToAccountListOld(testData)
			}
		})
	}
}

// BenchmarkMapAccountArrayGeneric benchmarks the new generic implementation
func BenchmarkMapAccountArrayGeneric(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run("Size"+string(rune(size)), func(b *testing.B) {
			testData := createTestAccountArray(size)
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				_ = MapAccountArrayToAccountList(testData)
			}
		})
	}
}

// BenchmarkMapCategoryArrayOld benchmarks old category mapping
func BenchmarkMapCategoryArrayOld(b *testing.B) {
	// Create test category array
	categoryArray := &client.CategoryArray{
		Data: make([]client.CategoryRead, 100),
	}
	
	for i := 0; i < 100; i++ {
		categoryArray.Data[i] = client.CategoryRead{
			Id:   "cat-" + string(rune(i)),
			Type: "categories",
			Attributes: client.Category{
				Name:  "Category " + string(rune(i)),
				Notes: strPtr("Notes"),
			},
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		categoryList := &dto.CategoryList{
			Data: make([]dto.Category, 0),
		}

		// Old implementation inline
		for _, categoryRead := range categoryArray.Data {
			category := dto.Category{
				Id:    categoryRead.Id,
				Name:  categoryRead.Attributes.Name,
				Notes: categoryRead.Attributes.Notes,
			}
			categoryList.Data = append(categoryList.Data, category)
		}
	}
}

// BenchmarkMapCategoryArrayGeneric benchmarks new generic category mapping
func BenchmarkMapCategoryArrayGeneric(b *testing.B) {
	// Create test category array
	categoryArray := &client.CategoryArray{
		Data: make([]client.CategoryRead, 100),
	}
	
	for i := 0; i < 100; i++ {
		categoryArray.Data[i] = client.CategoryRead{
			Id:   "cat-" + string(rune(i)),
			Type: "categories",
			Attributes: client.Category{
				Name:  "Category " + string(rune(i)),
				Notes: strPtr("Notes"),
			},
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapCategoryArrayToCategoryList(categoryArray)
	}
}

// BenchmarkMemoryAllocation compares memory allocation between implementations
func BenchmarkMemoryAllocation(b *testing.B) {
	testData := createTestAccountArray(100)
	
	b.Run("OldImplementation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = mapAccountArrayToAccountListOld(testData)
		}
	})
	
	b.Run("GenericImplementation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = MapAccountArrayToAccountList(testData)
		}
	})
}