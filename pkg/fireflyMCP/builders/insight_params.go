package builders

import (
	"fmt"
	"strconv"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
)

// GetSummaryParamsBuilder provides fluent interface for building GetBasicSummaryParams
type GetSummaryParamsBuilder struct {
	params *client.GetBasicSummaryParams
	*DateRangeBuilder
	errors []error
}

// NewGetSummaryParamsBuilder creates a new builder for GetBasicSummaryParams
func NewGetSummaryParamsBuilder() *GetSummaryParamsBuilder {
	builder := &GetSummaryParamsBuilder{
		params:           &client.GetBasicSummaryParams{},
		DateRangeBuilder: NewDateRangeBuilder(),
		errors:           []error{},
	}
	
	// Set default to current month if no dates specified
	builder.WithCurrentMonth()
	
	return builder
}

// Build constructs the final GetBasicSummaryParams
func (b *GetSummaryParamsBuilder) Build() (*client.GetBasicSummaryParams, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	
	// Apply date range
	if start := b.DateRangeBuilder.GetStart(); start != nil {
		b.params.Start = *start
	}
	if end := b.DateRangeBuilder.GetEnd(); end != nil {
		b.params.End = *end
	}
	
	return b.params, nil
}

// Validate checks if the builder state is valid
func (b *GetSummaryParamsBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	
	// Validate date range
	if err := b.DateRangeBuilder.Validate(); err != nil {
		return err
	}
	
	// Summary requires both start and end dates
	if b.DateRangeBuilder.GetStart() == nil || b.DateRangeBuilder.GetEnd() == nil {
		return fmt.Errorf("both start and end dates are required for summary")
	}
	
	return nil
}

// ExpenseCategoryInsightsParamsBuilder provides fluent interface for building expense category insights params
type ExpenseCategoryInsightsParamsBuilder struct {
	params *client.InsightExpenseCategoryParams
	*DateRangeBuilder
	accounts []int64
	errors   []error
}

// NewExpenseCategoryInsightsParamsBuilder creates a new builder for expense category insights
func NewExpenseCategoryInsightsParamsBuilder() *ExpenseCategoryInsightsParamsBuilder {
	return &ExpenseCategoryInsightsParamsBuilder{
		params:           &client.InsightExpenseCategoryParams{},
		DateRangeBuilder: NewDateRangeBuilder(),
		accounts:         []int64{},
		errors:           []error{},
	}
}

// WithAccount adds an account ID to include in results
func (b *ExpenseCategoryInsightsParamsBuilder) WithAccount(accountID int64) *ExpenseCategoryInsightsParamsBuilder {
	if accountID <= 0 {
		b.errors = append(b.errors, fmt.Errorf("account ID must be positive"))
		return b
	}
	b.accounts = append(b.accounts, accountID)
	return b
}

// WithAccountID adds an account ID (string) to include in results
func (b *ExpenseCategoryInsightsParamsBuilder) WithAccountID(accountID string) *ExpenseCategoryInsightsParamsBuilder {
	if accountID == "" {
		b.errors = append(b.errors, fmt.Errorf("account ID cannot be empty"))
		return b
	}
	id, err := strconv.ParseInt(accountID, 10, 64)
	if err != nil {
		b.errors = append(b.errors, fmt.Errorf("invalid account ID: %w", err))
		return b
	}
	return b.WithAccount(id)
}

// WithAccounts adds multiple account IDs to include in results
func (b *ExpenseCategoryInsightsParamsBuilder) WithAccounts(accountIDs ...int64) *ExpenseCategoryInsightsParamsBuilder {
	for _, id := range accountIDs {
		b.WithAccount(id)
	}
	return b
}

// Build constructs the final InsightExpenseCategoryParams
func (b *ExpenseCategoryInsightsParamsBuilder) Build() (*client.InsightExpenseCategoryParams, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	
	// Apply date range
	if start := b.DateRangeBuilder.GetStart(); start != nil {
		b.params.Start = *start
	}
	if end := b.DateRangeBuilder.GetEnd(); end != nil {
		b.params.End = *end
	}
	
	// Apply accounts if specified
	if len(b.accounts) > 0 {
		accountsParam := b.accounts
		b.params.Accounts = &accountsParam
	}
	
	return b.params, nil
}

// Validate checks if the builder state is valid
func (b *ExpenseCategoryInsightsParamsBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	
	// Validate date range
	if err := b.DateRangeBuilder.Validate(); err != nil {
		return err
	}
	
	// Insights require both start and end dates
	if b.DateRangeBuilder.GetStart() == nil || b.DateRangeBuilder.GetEnd() == nil {
		return fmt.Errorf("both start and end dates are required for insights")
	}
	
	return nil
}

// ExpenseTotalInsightsParamsBuilder provides fluent interface for building expense total insights params
type ExpenseTotalInsightsParamsBuilder struct {
	params *client.InsightExpenseTotalParams
	*DateRangeBuilder
	accounts []int64
	errors   []error
}

// NewExpenseTotalInsightsParamsBuilder creates a new builder for expense total insights
func NewExpenseTotalInsightsParamsBuilder() *ExpenseTotalInsightsParamsBuilder {
	return &ExpenseTotalInsightsParamsBuilder{
		params:           &client.InsightExpenseTotalParams{},
		DateRangeBuilder: NewDateRangeBuilder(),
		accounts:         []int64{},
		errors:           []error{},
	}
}

// WithAccount adds an account ID to include in results
func (b *ExpenseTotalInsightsParamsBuilder) WithAccount(accountID int64) *ExpenseTotalInsightsParamsBuilder {
	if accountID <= 0 {
		b.errors = append(b.errors, fmt.Errorf("account ID must be positive"))
		return b
	}
	b.accounts = append(b.accounts, accountID)
	return b
}

// WithAccountID adds an account ID (string) to include in results
func (b *ExpenseTotalInsightsParamsBuilder) WithAccountID(accountID string) *ExpenseTotalInsightsParamsBuilder {
	if accountID == "" {
		b.errors = append(b.errors, fmt.Errorf("account ID cannot be empty"))
		return b
	}
	id, err := strconv.ParseInt(accountID, 10, 64)
	if err != nil {
		b.errors = append(b.errors, fmt.Errorf("invalid account ID: %w", err))
		return b
	}
	return b.WithAccount(id)
}

// WithAccounts adds multiple account IDs to include in results
func (b *ExpenseTotalInsightsParamsBuilder) WithAccounts(accountIDs ...int64) *ExpenseTotalInsightsParamsBuilder {
	for _, id := range accountIDs {
		b.WithAccount(id)
	}
	return b
}

// Build constructs the final InsightExpenseTotalParams
func (b *ExpenseTotalInsightsParamsBuilder) Build() (*client.InsightExpenseTotalParams, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	
	// Apply date range
	if start := b.DateRangeBuilder.GetStart(); start != nil {
		b.params.Start = *start
	}
	if end := b.DateRangeBuilder.GetEnd(); end != nil {
		b.params.End = *end
	}
	
	// Apply accounts if specified
	if len(b.accounts) > 0 {
		accountsParam := b.accounts
		b.params.Accounts = &accountsParam
	}
	
	return b.params, nil
}

// Validate checks if the builder state is valid
func (b *ExpenseTotalInsightsParamsBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	
	// Validate date range
	if err := b.DateRangeBuilder.Validate(); err != nil {
		return err
	}
	
	// Insights require both start and end dates
	if b.DateRangeBuilder.GetStart() == nil || b.DateRangeBuilder.GetEnd() == nil {
		return fmt.Errorf("both start and end dates are required for insights")
	}
	
	return nil
}