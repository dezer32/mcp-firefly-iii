package fireflyMCP

import (
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestMapBillToBill(t *testing.T) {
	// Test with nil input
	result := mapBillToBill(nil)
	assert.Nil(t, result)

	// Test with sample data
	active := true
	amountMin := "50.00"
	amountMax := "100.00"
	billDate := time.Now()
	repeatFreq := "monthly"
	skip := int32(0)
	currencyCode := "USD"
	currencySymbol := "$"
	notes := "Test bill notes"
	nextExpectedMatch := time.Now().AddDate(0, 1, 0)

	billRead := &client.BillRead{
		Id: "1",
		Attributes: client.Bill{
			Active:            &active,
			Name:              "Test Bill",
			AmountMin:         amountMin,
			AmountMax:         amountMax,
			Date:              billDate,
			RepeatFreq:        client.BillRepeatFrequency(repeatFreq),
			Skip:              &skip,
			CurrencyCode:      &currencyCode,
			CurrencySymbol:    &currencySymbol,
			Notes:             &notes,
			NextExpectedMatch: &nextExpectedMatch,
			PaidDates: &[]struct {
				Date                 *time.Time `json:"date,omitempty"`
				TransactionGroupId   *string    `json:"transaction_group_id,omitempty"`
				TransactionJournalId *string    `json:"transaction_journal_id,omitempty"`
			}{
				{
					Date:                 &billDate,
					TransactionGroupId:   strPtr("123"),
					TransactionJournalId: strPtr("456"),
				},
			},
		},
		Type: "bills",
	}

	result = mapBillToBill(billRead)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Id)
	assert.True(t, result.Active)
	assert.Equal(t, "Test Bill", result.Name)
	assert.Equal(t, amountMin, result.AmountMin)
	assert.Equal(t, amountMax, result.AmountMax)
	assert.Equal(t, billDate, result.Date)
	assert.Equal(t, repeatFreq, result.RepeatFreq)
	assert.Equal(t, int(skip), result.Skip)
	assert.Equal(t, currencyCode, result.CurrencyCode)
	// CurrencySymbol is not part of Bill DTO
	assert.Equal(t, &notes, result.Notes)
	assert.Equal(t, &nextExpectedMatch, result.NextExpectedMatch)

	// Verify paid dates
	assert.Len(t, result.PaidDates, 1)
	assert.Equal(t, &billDate, result.PaidDates[0].Date)
	assert.Equal(t, strPtr("123"), result.PaidDates[0].TransactionGroupId)
	assert.Equal(t, strPtr("456"), result.PaidDates[0].TransactionJournalId)
}

func TestMapBillToBill_MinimalData(t *testing.T) {
	// Test with minimal required data
	billRead := &client.BillRead{
		Id: "2",
		Attributes: client.Bill{
			Name:       "Minimal Bill",
			AmountMin:  "10.00",
			AmountMax:  "20.00",
			Date:       time.Now(),
			RepeatFreq: "weekly",
		},
		Type: "bills",
	}

	result2 := mapBillToBill(billRead)

	// Verify the mapping
	assert.NotNil(t, result2)
	assert.Equal(t, "2", result2.Id)
	assert.False(t, result2.Active) // default false
	assert.Equal(t, "Minimal Bill", result2.Name)
	assert.Equal(t, "10.00", result2.AmountMin)
	assert.Equal(t, "20.00", result2.AmountMax)
	assert.Equal(t, "weekly", result2.RepeatFreq)
	assert.Equal(t, 0, result2.Skip) // default 0
	assert.Empty(t, result2.CurrencyCode)
	// CurrencySymbol is not part of Bill DTO
	assert.Nil(t, result2.Notes)
	assert.Nil(t, result2.NextExpectedMatch)
	assert.Empty(t, result2.PaidDates)
}

func TestMapBillArrayToBillList(t *testing.T) {
	// Test with nil input
	result := mapBillArrayToBillList(nil)
	assert.Nil(t, result)

	// Test with empty bill array
	emptyArray := &client.BillArray{
		Data: []client.BillRead{},
		Meta: client.Meta{},
	}
	result = mapBillArrayToBillList(emptyArray)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data)

	// Test with sample data
	active := true
	count := 1
	total := 1
	currentPage := 1
	perPage := 10
	totalPages := 1

	billArray := &client.BillArray{
		Data: []client.BillRead{
			{
				Id: "1",
				Attributes: client.Bill{
					Active:     &active,
					Name:       "Test Bill",
					AmountMin:  "50.00",
					AmountMax:  "100.00",
					Date:       time.Now(),
					RepeatFreq: "monthly",
				},
				Type: "bills",
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

	result = mapBillArrayToBillList(billArray)

	// Verify the mapping
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	bill := result.Data[0]
	assert.Equal(t, "1", bill.Id)
	assert.True(t, bill.Active)
	assert.Equal(t, "Test Bill", bill.Name)

	// Verify pagination
	assert.Equal(t, count, result.Pagination.Count)
	assert.Equal(t, total, result.Pagination.Total)
	assert.Equal(t, currentPage, result.Pagination.CurrentPage)
	assert.Equal(t, perPage, result.Pagination.PerPage)
	assert.Equal(t, totalPages, result.Pagination.TotalPages)
}

func TestMapBillArrayToBillList_NilPagination(t *testing.T) {
	// Test with nil pagination
	billArray := &client.BillArray{
		Data: []client.BillRead{
			{
				Id: "1",
				Attributes: client.Bill{
					Name:       "Test Bill",
					AmountMin:  "50.00",
					AmountMax:  "100.00",
					Date:       time.Now(),
					RepeatFreq: "monthly",
				},
				Type: "bills",
			},
		},
		Meta: client.Meta{
			Pagination: nil,
		},
	}

	result := mapBillArrayToBillList(billArray)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)

	// Verify pagination has default zero values
	assert.Equal(t, 0, result.Pagination.Count)
	assert.Equal(t, 0, result.Pagination.Total)
}

// Helper function for creating string pointers
func strPtr(s string) *string {
	return &s
}
