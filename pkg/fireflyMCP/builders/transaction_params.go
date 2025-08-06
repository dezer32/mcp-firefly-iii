package builders

import (
	"fmt"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
)

// ListTransactionParamsBuilder provides fluent interface for building ListTransactionParams
type ListTransactionParamsBuilder struct {
	params *client.ListTransactionParams
	*PaginationBuilder
	*DateRangeBuilder
	errors []error
}

// NewListTransactionParamsBuilder creates a new builder for ListTransactionParams
func NewListTransactionParamsBuilder() *ListTransactionParamsBuilder {
	return &ListTransactionParamsBuilder{
		params:            &client.ListTransactionParams{},
		PaginationBuilder: NewPaginationBuilder(),
		DateRangeBuilder:  NewDateRangeBuilder(),
		errors:            []error{},
	}
}

// WithType sets the transaction type filter
func (b *ListTransactionParamsBuilder) WithType(transactionType string) *ListTransactionParamsBuilder {
	validTypes := map[string]bool{
		"all":         true,
		"withdrawal":  true,
		"withdrawals": true,
		"expense":     true,
		"deposit":     true,
		"deposits":    true,
		"income":      true,
		"transfer":    true,
		"transfers":   true,
		"opening_balance": true,
		"reconciliation": true,
	}
	
	if !validTypes[transactionType] {
		b.errors = append(b.errors, fmt.Errorf("invalid transaction type: %s", transactionType))
		return b
	}
	
	filter := client.TransactionTypeFilter(transactionType)
	b.params.Type = &filter
	return b
}

// WithWithdrawals sets filter to withdrawals only
func (b *ListTransactionParamsBuilder) WithWithdrawals() *ListTransactionParamsBuilder {
	return b.WithType("withdrawal")
}

// WithDeposits sets filter to deposits only
func (b *ListTransactionParamsBuilder) WithDeposits() *ListTransactionParamsBuilder {
	return b.WithType("deposit")
}

// WithTransfers sets filter to transfers only
func (b *ListTransactionParamsBuilder) WithTransfers() *ListTransactionParamsBuilder {
	return b.WithType("transfer")
}

// Build constructs the final ListTransactionParams
func (b *ListTransactionParamsBuilder) Build() (*client.ListTransactionParams, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	
	// Apply date range
	b.params.Start = b.DateRangeBuilder.GetStart()
	b.params.End = b.DateRangeBuilder.GetEnd()
	
	// Apply pagination
	b.params.Limit = b.PaginationBuilder.GetLimit()
	b.params.Page = b.PaginationBuilder.GetPage()
	
	return b.params, nil
}

// Validate checks if the builder state is valid
func (b *ListTransactionParamsBuilder) Validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("validation errors: %v", b.errors)
	}
	
	// Validate date range
	if err := b.DateRangeBuilder.Validate(); err != nil {
		return err
	}
	
	// Validate pagination
	if err := b.PaginationBuilder.Validate(); err != nil {
		return err
	}
	
	return nil
}

// SearchTransactionParamsBuilder provides fluent interface for building SearchTransactionsParams
type SearchTransactionParamsBuilder struct {
	params *client.SearchTransactionsParams
	*SearchBuilder
	*PaginationBuilder
	*DateRangeBuilder
	errors []error
}

// NewSearchTransactionParamsBuilder creates a new builder for SearchTransactionsParams
func NewSearchTransactionParamsBuilder(query string) *SearchTransactionParamsBuilder {
	return &SearchTransactionParamsBuilder{
		params:            &client.SearchTransactionsParams{Query: query},
		SearchBuilder:     NewSearchBuilder(query),
		PaginationBuilder: NewPaginationBuilder(),
		DateRangeBuilder:  NewDateRangeBuilder(),
		errors:            []error{},
	}
}

// Build constructs the final SearchTransactionsParams
func (b *SearchTransactionParamsBuilder) Build() (*client.SearchTransactionsParams, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	
	// Apply pagination (note: SearchTransactionsParams uses different types)
	if limit := b.PaginationBuilder.GetLimit(); limit != nil {
		b.params.Limit = limit
	}
	if page := b.PaginationBuilder.GetPage(); page != nil {
		b.params.Page = page
	}
	
	// Date range is not available in SearchTransactionsParams based on the handler code
	// but we keep the builder for consistency and future extensibility
	
	return b.params, nil
}

// Validate checks if the builder state is valid
func (b *SearchTransactionParamsBuilder) Validate() error {
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
	
	// Validate date range if dates are set
	if err := b.DateRangeBuilder.Validate(); err != nil {
		return err
	}
	
	return nil
}