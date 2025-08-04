package fireflyMCP

import (
	"testing"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestMapRecurrenceRepetitionToRecurrenceRepetition(t *testing.T) {
	tests := []struct {
		name     string
		input    *client.RecurrenceRepetition
		expected RecurrenceRepetition
	}{
		{
			name: "Basic repetition",
			input: &client.RecurrenceRepetition{
				Id:          ptr("123"),
				Type:        client.Weekly,
				Moment:      "Monday",
				Skip:        ptr(int32(0)),
				Weekend:     ptr(int32(1)),
				Description: ptr("Weekly on Monday"),
			},
			expected: RecurrenceRepetition{
				Id:          "123",
				Type:        "weekly",
				Moment:      "Monday",
				Skip:        0,
				Weekend:     1,
				Description: ptr("Weekly on Monday"),
			},
		},
		{
			name: "Repetition with nil values",
			input: &client.RecurrenceRepetition{
				Id:          ptr("456"),
				Type:        client.Monthly,
				Moment:      "15",
				Skip:        nil,
				Weekend:     nil,
				Description: nil,
			},
			expected: RecurrenceRepetition{
				Id:          "456",
				Type:        "monthly",
				Moment:      "15",
				Skip:        0,
				Weekend:     0,
				Description: nil,
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: RecurrenceRepetition{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := mapRecurrenceRepetitionToRecurrenceRepetition(tt.input)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}

func TestMapRecurrenceTransactionToRecurrenceTransaction(t *testing.T) {
	tests := []struct {
		name     string
		input    *client.RecurrenceTransaction
		expected RecurrenceTransaction
	}{
		{
			name: "Complete transaction",
			input: &client.RecurrenceTransaction{
				Id:              ptr("789"),
				Description:     ptr("Monthly rent"),
				Amount:          ptr("1250.00"),
				CurrencyCode:    ptr("USD"),
				CurrencySymbol:  ptr("$"),
				CategoryId:      ptr("10"),
				CategoryName:    ptr("Housing"),
				BudgetId:        ptr("20"),
				BudgetName:      ptr("Monthly Budget"),
				SourceId:        ptr("1"),
				SourceName:      ptr("Checking Account"),
				DestinationId:   ptr("2"),
				DestinationName: ptr("Landlord"),
			},
			expected: RecurrenceTransaction{
				Id:              "789",
				Description:     "Monthly rent",
				Amount:          "1250.00",
				CurrencyCode:    "USD",
				CurrencySymbol:  "$",
				CategoryId:      ptr("10"),
				CategoryName:    ptr("Housing"),
				BudgetId:        ptr("20"),
				BudgetName:      ptr("Monthly Budget"),
				SourceId:        "1",
				SourceName:      "Checking Account",
				DestinationId:   "2",
				DestinationName: "Landlord",
			},
		},
		{
			name: "Transaction without category and budget",
			input: &client.RecurrenceTransaction{
				Id:              ptr("999"),
				Description:     ptr("Utility payment"),
				Amount:          ptr("75.50"),
				CurrencyCode:    ptr("EUR"),
				CurrencySymbol:  ptr("€"),
				CategoryId:      nil,
				CategoryName:    nil,
				BudgetId:        nil,
				BudgetName:      nil,
				SourceId:        ptr("3"),
				SourceName:      ptr("Main Account"),
				DestinationId:   ptr("4"),
				DestinationName: ptr("Utility Company"),
			},
			expected: RecurrenceTransaction{
				Id:              "999",
				Description:     "Utility payment",
				Amount:          "75.50",
				CurrencyCode:    "EUR",
				CurrencySymbol:  "€",
				CategoryId:      nil,
				CategoryName:    nil,
				BudgetId:        nil,
				BudgetName:      nil,
				SourceId:        "3",
				SourceName:      "Main Account",
				DestinationId:   "4",
				DestinationName: "Utility Company",
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: RecurrenceTransaction{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := mapRecurrenceTransactionToRecurrenceTransaction(tt.input)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}

func TestMapRecurrenceToRecurrence(t *testing.T) {
	now := time.Now()
	laterDate := now.Add(365 * 24 * time.Hour)
	repeatUntil := now.Add(730 * 24 * time.Hour)

	// Create openapi Date types
	nowDate := openapi_types.Date{Time: now}
	laterDateDate := openapi_types.Date{Time: laterDate}
	repeatUntilDate := openapi_types.Date{Time: repeatUntil}

	tests := []struct {
		name     string
		input    *client.RecurrenceSingle
		expected Recurrence
	}{
		{
			name: "Complete recurrence",
			input: &client.RecurrenceSingle{
				Data: client.RecurrenceRead{
					Id: "100",
					Attributes: client.Recurrence{
						Active:          ptr(true),
						ApplyRules:      ptr(false),
						Description:     ptr("Monthly rent payment"),
						FirstDate:       &nowDate,
						LatestDate:      &laterDateDate,
						Notes:           ptr("Apartment rent"),
						NrOfRepetitions: ptr(int32(12)),
						RepeatUntil:     &repeatUntilDate,
						Title:           ptr("Monthly Rent"),
						Type:            ptr(client.RecurrenceTransactionType("withdrawal")),
						Repetitions: &[]client.RecurrenceRepetition{
							{
								Id:          ptr("101"),
								Type:        client.Monthly,
								Moment:      "1",
								Skip:        ptr(int32(0)),
								Weekend:     ptr(int32(1)),
								Description: ptr("First of month"),
							},
						},
						Transactions: &[]client.RecurrenceTransaction{
							{
								Id:              ptr("102"),
								Description:     ptr("Rent payment"),
								Amount:          ptr("1250.00"),
								CurrencyCode:    ptr("USD"),
								CurrencySymbol:  ptr("$"),
								CategoryId:      ptr("10"),
								CategoryName:    ptr("Housing"),
								BudgetId:        ptr("20"),
								BudgetName:      ptr("Monthly Budget"),
								SourceId:        ptr("1"),
								SourceName:      ptr("Checking Account"),
								DestinationId:   ptr("2"),
								DestinationName: ptr("Landlord"),
							},
						},
					},
				},
			},
			expected: Recurrence{
				Id:              "100",
				Type:            "withdrawal",
				Title:           "Monthly Rent",
				Description:     "Monthly rent payment",
				FirstDate:       now,
				LatestDate:      &laterDate,
				RepeatUntil:     &repeatUntil,
				NrOfRepetitions: ptr(12),
				ApplyRules:      false,
				Active:          true,
				Notes:           ptr("Apartment rent"),
				Repetitions: []RecurrenceRepetition{
					{
						Id:          "101",
						Type:        "monthly",
						Moment:      "1",
						Skip:        0,
						Weekend:     1,
						Description: ptr("First of month"),
					},
				},
				Transactions: []RecurrenceTransaction{
					{
						Id:              "102",
						Description:     "Rent payment",
						Amount:          "1250.00",
						CurrencyCode:    "USD",
						CurrencySymbol:  "$",
						CategoryId:      ptr("10"),
						CategoryName:    ptr("Housing"),
						BudgetId:        ptr("20"),
						BudgetName:      ptr("Monthly Budget"),
						SourceId:        "1",
						SourceName:      "Checking Account",
						DestinationId:   "2",
						DestinationName: "Landlord",
					},
				},
			},
		},
		{
			name: "Minimal recurrence",
			input: &client.RecurrenceSingle{
				Data: client.RecurrenceRead{
					Id: "200",
					Attributes: client.Recurrence{
						Active:          ptr(false),
						ApplyRules:      ptr(true),
						Description:     ptr("Simple recurrence"),
						FirstDate:       &nowDate,
						LatestDate:      nil,
						Notes:           nil,
						NrOfRepetitions: nil,
						RepeatUntil:     nil,
						Title:           ptr("Simple"),
						Type:            ptr(client.RecurrenceTransactionType("deposit")),
						Repetitions:     nil,
						Transactions:    nil,
					},
				},
			},
			expected: Recurrence{
				Id:              "200",
				Type:            "deposit",
				Title:           "Simple",
				Description:     "Simple recurrence",
				FirstDate:       now,
				LatestDate:      nil,
				RepeatUntil:     nil,
				NrOfRepetitions: nil,
				ApplyRules:      true,
				Active:          false,
				Notes:           nil,
				Repetitions:     []RecurrenceRepetition{},
				Transactions:    []RecurrenceTransaction{},
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: Recurrence{},
		},
		{
			name: "Nil data attributes",
			input: &client.RecurrenceSingle{
				Data: client.RecurrenceRead{
					Id:         "300",
					Attributes: client.Recurrence{},
				},
			},
			expected: Recurrence{
				Id:           "300",
				FirstDate:    time.Time{},
				Repetitions:  []RecurrenceRepetition{},
				Transactions: []RecurrenceTransaction{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := mapRecurrenceToRecurrence(tt.input)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}

func TestMapRecurrenceArrayToRecurrenceList(t *testing.T) {
	now := time.Now()
	nowDate := openapi_types.Date{Time: now}

	tests := []struct {
		name     string
		input    *client.RecurrenceArray
		expected RecurrenceList
	}{
		{
			name: "Array with multiple recurrences",
			input: &client.RecurrenceArray{
				Data: []client.RecurrenceRead{
					{
						Id: "1",
						Attributes: client.Recurrence{
							Title:       ptr("Rent"),
							Type:        ptr(client.RecurrenceTransactionType("withdrawal")),
							Description: ptr("Monthly rent"),
							FirstDate:   &nowDate,
							Active:      ptr(true),
							ApplyRules:  ptr(false),
						},
					},
					{
						Id: "2",
						Attributes: client.Recurrence{
							Title:       ptr("Salary"),
							Type:        ptr(client.RecurrenceTransactionType("deposit")),
							Description: ptr("Monthly salary"),
							FirstDate:   &nowDate,
							Active:      ptr(true),
							ApplyRules:  ptr(true),
						},
					},
				},
				Meta: &client.Meta{
					Pagination: &client.MetaPagination{
						Count:       ptr(2),
						CurrentPage: ptr(1),
						PerPage:     ptr(10),
						Total:       ptr(2),
						TotalPages:  ptr(1),
					},
				},
			},
			expected: RecurrenceList{
				Data: []Recurrence{
					{
						Id:           "1",
						Title:        "Rent",
						Type:         "withdrawal",
						Description:  "Monthly rent",
						FirstDate:    now,
						Active:       true,
						ApplyRules:   false,
						Repetitions:  []RecurrenceRepetition{},
						Transactions: []RecurrenceTransaction{},
					},
					{
						Id:           "2",
						Title:        "Salary",
						Type:         "deposit",
						Description:  "Monthly salary",
						FirstDate:    now,
						Active:       true,
						ApplyRules:   true,
						Repetitions:  []RecurrenceRepetition{},
						Transactions: []RecurrenceTransaction{},
					},
				},
				Pagination: Pagination{
					Count:       2,
					CurrentPage: 1,
					PerPage:     10,
					Total:       2,
					TotalPages:  1,
				},
			},
		},
		{
			name: "Empty array",
			input: &client.RecurrenceArray{
				Data: []client.RecurrenceRead{},
				Meta: &client.Meta{
					Pagination: &client.MetaPagination{
						Count:       ptr(0),
						CurrentPage: ptr(1),
						PerPage:     ptr(10),
						Total:       ptr(0),
						TotalPages:  ptr(0),
					},
				},
			},
			expected: RecurrenceList{
				Data: []Recurrence{},
				Pagination: Pagination{
					Count:       0,
					CurrentPage: 1,
					PerPage:     10,
					Total:       0,
					TotalPages:  0,
				},
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: RecurrenceList{Data: []Recurrence{}},
		},
		{
			name: "Nil meta",
			input: &client.RecurrenceArray{
				Data: []client.RecurrenceRead{
					{
						Id: "3",
						Attributes: client.Recurrence{
							Title:      ptr("Test"),
							Type:       ptr(client.RecurrenceTransactionType("transfer")),
							FirstDate:  &nowDate,
							Active:     ptr(true),
							ApplyRules: ptr(false),
						},
					},
				},
				Meta: nil,
			},
			expected: RecurrenceList{
				Data: []Recurrence{
					{
						Id:           "3",
						Title:        "Test",
						Type:         "transfer",
						FirstDate:    now,
						Active:       true,
						ApplyRules:   false,
						Repetitions:  []RecurrenceRepetition{},
						Transactions: []RecurrenceTransaction{},
					},
				},
				Pagination: Pagination{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := mapRecurrenceArrayToRecurrenceList(tt.input)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}

// Helper function
func ptr[T any](v T) *T {
	return &v
}
