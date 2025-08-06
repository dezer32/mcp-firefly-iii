// Package builders provides fluent interfaces for constructing API parameters
package builders

import (
	"errors"
	"fmt"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Builder is the base interface for all parameter builders
type Builder interface {
	// Build constructs and validates the final parameters
	Build() error
	// Validate performs validation on the current state
	Validate() error
}

// PaginationBuilder provides fluent interface for pagination parameters
type PaginationBuilder struct {
	limit *int32
	page  *int32
	errors []error
}

// NewPaginationBuilder creates a new PaginationBuilder
func NewPaginationBuilder() *PaginationBuilder {
	return &PaginationBuilder{
		errors: []error{},
	}
}

// WithLimit sets the limit for pagination
func (b *PaginationBuilder) WithLimit(limit int) *PaginationBuilder {
	if limit <= 0 {
		b.errors = append(b.errors, errors.New("limit must be greater than 0"))
		return b
	}
	if limit > 1000 {
		b.errors = append(b.errors, errors.New("limit cannot exceed 1000"))
		return b
	}
	l := int32(limit)
	b.limit = &l
	return b
}

// WithPage sets the page number for pagination
func (b *PaginationBuilder) WithPage(page int) *PaginationBuilder {
	if page <= 0 {
		b.errors = append(b.errors, errors.New("page must be greater than 0"))
		return b
	}
	p := int32(page)
	b.page = &p
	return b
}

// GetLimit returns the limit value
func (b *PaginationBuilder) GetLimit() *int32 {
	return b.limit
}

// GetPage returns the page value
func (b *PaginationBuilder) GetPage() *int32 {
	return b.page
}

// Build constructs and validates the pagination parameters
func (b *PaginationBuilder) Build() error {
	return b.Validate()
}

// Validate checks if the builder state is valid
func (b *PaginationBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	return nil
}

// DateRangeBuilder provides fluent interface for date range parameters
type DateRangeBuilder struct {
	start  *openapi_types.Date
	end    *openapi_types.Date
	errors []error
}

// NewDateRangeBuilder creates a new DateRangeBuilder
func NewDateRangeBuilder() *DateRangeBuilder {
	return &DateRangeBuilder{
		errors: []error{},
	}
}

// WithStartDate sets the start date
func (b *DateRangeBuilder) WithStartDate(date string) *DateRangeBuilder {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		b.errors = append(b.errors, fmt.Errorf("invalid start date format: %w", err))
		return b
	}
	d := openapi_types.Date{Time: parsedDate}
	b.start = &d
	return b
}

// WithEndDate sets the end date
func (b *DateRangeBuilder) WithEndDate(date string) *DateRangeBuilder {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		b.errors = append(b.errors, fmt.Errorf("invalid end date format: %w", err))
		return b
	}
	d := openapi_types.Date{Time: parsedDate}
	b.end = &d
	return b
}

// WithDateRange sets both start and end dates
func (b *DateRangeBuilder) WithDateRange(start, end string) *DateRangeBuilder {
	return b.WithStartDate(start).WithEndDate(end)
}

// WithCurrentMonth sets the date range to the current month
func (b *DateRangeBuilder) WithCurrentMonth() *DateRangeBuilder {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := now.AddDate(0, 1, 0)
	endOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 23, 59, 59, 999999999, now.Location())
	
	b.start = &openapi_types.Date{Time: startOfMonth}
	b.end = &openapi_types.Date{Time: endOfMonth}
	return b
}

// WithLastNDays sets the date range to the last N days
func (b *DateRangeBuilder) WithLastNDays(days int) *DateRangeBuilder {
	if days <= 0 {
		b.errors = append(b.errors, errors.New("days must be greater than 0"))
		return b
	}
	now := time.Now()
	start := now.AddDate(0, 0, -days)
	
	b.start = &openapi_types.Date{Time: start}
	b.end = &openapi_types.Date{Time: now}
	return b
}

// GetStart returns the start date
func (b *DateRangeBuilder) GetStart() *openapi_types.Date {
	return b.start
}

// GetEnd returns the end date
func (b *DateRangeBuilder) GetEnd() *openapi_types.Date {
	return b.end
}

// Build constructs and validates the date range parameters
func (b *DateRangeBuilder) Build() error {
	return b.Validate()
}

// Validate checks if the builder state is valid
func (b *DateRangeBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	
	// Validate date range logic
	if b.start != nil && b.end != nil {
		if b.start.Time.After(b.end.Time) {
			return errors.New("start date must be before end date")
		}
	}
	
	return nil
}

// SearchBuilder provides fluent interface for search parameters
type SearchBuilder struct {
	query  string
	field  string
	errors []error
}

// NewSearchBuilder creates a new SearchBuilder
func NewSearchBuilder(query string) *SearchBuilder {
	sb := &SearchBuilder{
		query:  query,
		errors: []error{},
	}
	
	if query == "" {
		sb.errors = append(sb.errors, errors.New("search query cannot be empty"))
	}
	
	return sb
}

// WithField sets the field to search in
func (b *SearchBuilder) WithField(field string) *SearchBuilder {
	validFields := map[string]bool{
		"all":    true,
		"iban":   true,
		"name":   true,
		"number": true,
		"id":     true,
	}
	
	if !validFields[field] {
		b.errors = append(b.errors, fmt.Errorf("invalid search field: %s", field))
		return b
	}
	
	b.field = field
	return b
}

// GetQuery returns the search query
func (b *SearchBuilder) GetQuery() string {
	return b.query
}

// GetField returns the search field
func (b *SearchBuilder) GetField() client.AccountSearchFieldFilter {
	if b.field == "" {
		return "all"
	}
	return client.AccountSearchFieldFilter(b.field)
}

// Build constructs and validates the search parameters
func (b *SearchBuilder) Build() error {
	return b.Validate()
}

// Validate checks if the builder state is valid
func (b *SearchBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	return nil
}

// CompositeBuilder combines multiple builders
type CompositeBuilder struct {
	builders []Builder
	errors   []error
}

// NewCompositeBuilder creates a new CompositeBuilder
func NewCompositeBuilder() *CompositeBuilder {
	return &CompositeBuilder{
		builders: []Builder{},
		errors:   []error{},
	}
}

// Add adds a builder to the composite
func (b *CompositeBuilder) Add(builder Builder) *CompositeBuilder {
	if builder == nil {
		b.errors = append(b.errors, errors.New("cannot add nil builder"))
		return b
	}
	b.builders = append(b.builders, builder)
	return b
}

// Build constructs and validates all builders
func (b *CompositeBuilder) Build() error {
	return b.Validate()
}

// Validate checks if all builders are valid
func (b *CompositeBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("composite builder errors: %v", b.errors)
	}
	
	for i, builder := range b.builders {
		if err := builder.Validate(); err != nil {
			return fmt.Errorf("builder %d validation failed: %w", i, err)
		}
	}
	
	return nil
}