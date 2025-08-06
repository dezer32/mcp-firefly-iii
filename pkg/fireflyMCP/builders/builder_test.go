package builders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaginationBuilder(t *testing.T) {
	tests := []struct {
		name      string
		build     func() *PaginationBuilder
		wantLimit *int32
		wantPage  *int32
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid pagination",
			build: func() *PaginationBuilder {
				return NewPaginationBuilder().
					WithLimit(10).
					WithPage(2)
			},
			wantLimit: int32Ptr(10),
			wantPage:  int32Ptr(2),
			wantError: false,
		},
		{
			name: "invalid limit (zero)",
			build: func() *PaginationBuilder {
				return NewPaginationBuilder().
					WithLimit(0)
			},
			wantLimit: nil,
			wantPage:  nil,
			wantError: true,
			errorMsg:  "limit must be greater than 0",
		},
		{
			name: "invalid limit (too large)",
			build: func() *PaginationBuilder {
				return NewPaginationBuilder().
					WithLimit(1001)
			},
			wantLimit: nil,
			wantPage:  nil,
			wantError: true,
			errorMsg:  "limit cannot exceed 1000",
		},
		{
			name: "invalid page (zero)",
			build: func() *PaginationBuilder {
				return NewPaginationBuilder().
					WithPage(0)
			},
			wantLimit: nil,
			wantPage:  nil,
			wantError: true,
			errorMsg:  "page must be greater than 0",
		},
		{
			name: "only limit set",
			build: func() *PaginationBuilder {
				return NewPaginationBuilder().
					WithLimit(25)
			},
			wantLimit: int32Ptr(25),
			wantPage:  nil,
			wantError: false,
		},
		{
			name: "only page set",
			build: func() *PaginationBuilder {
				return NewPaginationBuilder().
					WithPage(3)
			},
			wantLimit: nil,
			wantPage:  int32Ptr(3),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.build()
			err := builder.Validate()

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLimit, builder.GetLimit())
				assert.Equal(t, tt.wantPage, builder.GetPage())
			}
		})
	}
}

func TestDateRangeBuilder(t *testing.T) {
	tests := []struct {
		name      string
		build     func() *DateRangeBuilder
		wantError bool
		errorMsg  string
		validate  func(t *testing.T, b *DateRangeBuilder)
	}{
		{
			name: "valid date range",
			build: func() *DateRangeBuilder {
				return NewDateRangeBuilder().
					WithStartDate("2024-01-01").
					WithEndDate("2024-01-31")
			},
			wantError: false,
			validate: func(t *testing.T, b *DateRangeBuilder) {
				require.NotNil(t, b.GetStart())
				require.NotNil(t, b.GetEnd())
				assert.Equal(t, "2024-01-01", b.GetStart().Time.Format("2006-01-02"))
				assert.Equal(t, "2024-01-31", b.GetEnd().Time.Format("2006-01-02"))
			},
		},
		{
			name: "invalid start date format",
			build: func() *DateRangeBuilder {
				return NewDateRangeBuilder().
					WithStartDate("2024/01/01")
			},
			wantError: true,
			errorMsg:  "invalid start date format",
		},
		{
			name: "invalid end date format",
			build: func() *DateRangeBuilder {
				return NewDateRangeBuilder().
					WithEndDate("31-01-2024")
			},
			wantError: true,
			errorMsg:  "invalid end date format",
		},
		{
			name: "start date after end date",
			build: func() *DateRangeBuilder {
				return NewDateRangeBuilder().
					WithStartDate("2024-02-01").
					WithEndDate("2024-01-01")
			},
			wantError: true,
			errorMsg:  "start date must be before end date",
		},
		{
			name: "with date range helper",
			build: func() *DateRangeBuilder {
				return NewDateRangeBuilder().
					WithDateRange("2024-03-01", "2024-03-31")
			},
			wantError: false,
			validate: func(t *testing.T, b *DateRangeBuilder) {
				require.NotNil(t, b.GetStart())
				require.NotNil(t, b.GetEnd())
				assert.Equal(t, "2024-03-01", b.GetStart().Time.Format("2006-01-02"))
				assert.Equal(t, "2024-03-31", b.GetEnd().Time.Format("2006-01-02"))
			},
		},
		{
			name: "with current month",
			build: func() *DateRangeBuilder {
				return NewDateRangeBuilder().
					WithCurrentMonth()
			},
			wantError: false,
			validate: func(t *testing.T, b *DateRangeBuilder) {
				require.NotNil(t, b.GetStart())
				require.NotNil(t, b.GetEnd())
				
				now := time.Now()
				startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				assert.Equal(t, startOfMonth.Format("2006-01-02"), b.GetStart().Time.Format("2006-01-02"))
			},
		},
		{
			name: "with last N days",
			build: func() *DateRangeBuilder {
				return NewDateRangeBuilder().
					WithLastNDays(7)
			},
			wantError: false,
			validate: func(t *testing.T, b *DateRangeBuilder) {
				require.NotNil(t, b.GetStart())
				require.NotNil(t, b.GetEnd())
				
				duration := b.GetEnd().Time.Sub(b.GetStart().Time)
				assert.InDelta(t, 7*24*time.Hour, duration, float64(24*time.Hour))
			},
		},
		{
			name: "invalid last N days (zero)",
			build: func() *DateRangeBuilder {
				return NewDateRangeBuilder().
					WithLastNDays(0)
			},
			wantError: true,
			errorMsg:  "days must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.build()
			err := builder.Validate()

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, builder)
				}
			}
		})
	}
}

func TestSearchBuilder(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		field     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid search with query only",
			query:     "test account",
			field:     "",
			wantError: false,
		},
		{
			name:      "valid search with field",
			query:     "test",
			field:     "name",
			wantError: false,
		},
		{
			name:      "empty query",
			query:     "",
			field:     "name",
			wantError: true,
			errorMsg:  "search query cannot be empty",
		},
		{
			name:      "invalid field",
			query:     "test",
			field:     "invalid_field",
			wantError: true,
			errorMsg:  "invalid search field",
		},
		{
			name:      "all valid fields",
			query:     "test",
			field:     "all",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchBuilder(tt.query)
			if tt.field != "" {
				builder.WithField(tt.field)
			}
			
			err := builder.Validate()
			
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.query, builder.GetQuery())
				if tt.field == "" {
					assert.Equal(t, "all", string(builder.GetField()))
				} else {
					assert.Equal(t, tt.field, string(builder.GetField()))
				}
			}
		})
	}
}

func TestCompositeBuilder(t *testing.T) {
	t.Run("valid composite", func(t *testing.T) {
		composite := NewCompositeBuilder().
			Add(NewPaginationBuilder().WithLimit(10)).
			Add(NewDateRangeBuilder().WithCurrentMonth())
		
		err := composite.Validate()
		assert.NoError(t, err)
	})
	
	t.Run("invalid builder in composite", func(t *testing.T) {
		composite := NewCompositeBuilder().
			Add(NewPaginationBuilder().WithLimit(0)). // Invalid limit
			Add(NewDateRangeBuilder().WithCurrentMonth())
		
		err := composite.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "builder 0 validation failed")
	})
	
	t.Run("nil builder", func(t *testing.T) {
		composite := NewCompositeBuilder().
			Add(nil)
		
		err := composite.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot add nil builder")
	})
	
	t.Run("build method", func(t *testing.T) {
		composite := NewCompositeBuilder().
			Add(NewPaginationBuilder().WithLimit(25))
		
		err := composite.Build()
		assert.NoError(t, err)
	})
}

// Helper function to create int32 pointer
func int32Ptr(i int32) *int32 {
	return &i
}