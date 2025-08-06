package builders

import (
	"fmt"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
)

// ListAccountParamsBuilder provides fluent interface for building ListAccountParams
type ListAccountParamsBuilder struct {
	params *client.ListAccountParams
	*PaginationBuilder
	errors []error
}

// NewListAccountParamsBuilder creates a new builder for ListAccountParams
func NewListAccountParamsBuilder() *ListAccountParamsBuilder {
	return &ListAccountParamsBuilder{
		params:            &client.ListAccountParams{},
		PaginationBuilder: NewPaginationBuilder(),
		errors:            []error{},
	}
}

// WithType sets the account type filter
func (b *ListAccountParamsBuilder) WithType(accountType string) *ListAccountParamsBuilder {
	validTypes := map[string]bool{
		"asset":           true,
		"expense":         true,
		"revenue":         true,
		"liability":       true,
		"liabilities":     true,
		"initial_balance": true,
		"cash":            true,
	}
	
	if !validTypes[accountType] {
		b.errors = append(b.errors, fmt.Errorf("invalid account type: %s", accountType))
		return b
	}
	
	filter := client.AccountTypeFilter(accountType)
	b.params.Type = &filter
	return b
}

// WithAssetAccounts sets filter to asset accounts only
func (b *ListAccountParamsBuilder) WithAssetAccounts() *ListAccountParamsBuilder {
	return b.WithType("asset")
}

// WithExpenseAccounts sets filter to expense accounts only
func (b *ListAccountParamsBuilder) WithExpenseAccounts() *ListAccountParamsBuilder {
	return b.WithType("expense")
}

// WithRevenueAccounts sets filter to revenue accounts only
func (b *ListAccountParamsBuilder) WithRevenueAccounts() *ListAccountParamsBuilder {
	return b.WithType("revenue")
}

// Build constructs the final ListAccountParams
func (b *ListAccountParamsBuilder) Build() (*client.ListAccountParams, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	
	// Apply pagination
	b.params.Limit = b.PaginationBuilder.GetLimit()
	b.params.Page = b.PaginationBuilder.GetPage()
	
	return b.params, nil
}

// Validate checks if the builder state is valid
func (b *ListAccountParamsBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	
	// Validate pagination
	if err := b.PaginationBuilder.Validate(); err != nil {
		return err
	}
	
	return nil
}

// SearchAccountParamsBuilder provides fluent interface for building SearchAccountsParams
type SearchAccountParamsBuilder struct {
	params *client.SearchAccountsParams
	*SearchBuilder
	*PaginationBuilder
	errors []error
}

// NewSearchAccountParamsBuilder creates a new builder for SearchAccountsParams
func NewSearchAccountParamsBuilder(query string) *SearchAccountParamsBuilder {
	return &SearchAccountParamsBuilder{
		params:            &client.SearchAccountsParams{Query: query},
		SearchBuilder:     NewSearchBuilder(query),
		PaginationBuilder: NewPaginationBuilder(),
		errors:            []error{},
	}
}

// Build constructs the final SearchAccountsParams
func (b *SearchAccountParamsBuilder) Build() (*client.SearchAccountsParams, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	
	// Apply search field
	b.params.Field = b.SearchBuilder.GetField()
	
	// Apply pagination
	b.params.Limit = b.PaginationBuilder.GetLimit()
	b.params.Page = b.PaginationBuilder.GetPage()
	
	return b.params, nil
}

// Validate checks if the builder state is valid
func (b *SearchAccountParamsBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	
	// Validate search
	if err := b.SearchBuilder.Validate(); err != nil {
		return err
	}
	
	// Validate pagination
	if err := b.PaginationBuilder.Validate(); err != nil {
		return err
	}
	
	return nil
}