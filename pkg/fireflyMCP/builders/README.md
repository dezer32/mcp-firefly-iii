# Firefly III MCP Parameter Builders

The builders package provides a fluent interface for constructing API parameters for the Firefly III MCP server. This pattern improves code readability, ensures proper validation, and reduces errors when building complex parameter sets.

## Features

- **Fluent Interface**: Chain method calls for intuitive parameter construction
- **Built-in Validation**: Validates parameters at build time to catch errors early
- **Type Safety**: Leverages Go's type system to prevent invalid parameter combinations
- **Composability**: Combine multiple builders for complex parameter sets
- **Backward Compatible**: Can be used alongside existing direct parameter construction

## Core Builders

### PaginationBuilder

Handles pagination parameters for list operations:

```go
builder := NewPaginationBuilder().
    WithLimit(25).    // Set items per page
    WithPage(2)       // Set page number

// Validation happens automatically
if err := builder.Validate(); err != nil {
    // Handle validation error
}
```

### DateRangeBuilder

Manages date range parameters with convenient helpers:

```go
// Using specific dates
builder := NewDateRangeBuilder().
    WithStartDate("2024-01-01").
    WithEndDate("2024-01-31")

// Using convenience methods
builder := NewDateRangeBuilder().
    WithCurrentMonth()    // Sets to current month

builder := NewDateRangeBuilder().
    WithLastNDays(30)     // Last 30 days

// Using date range helper
builder := NewDateRangeBuilder().
    WithDateRange("2024-01-01", "2024-12-31")
```

### SearchBuilder

Constructs search parameters with field filtering:

```go
builder := NewSearchBuilder("search term").
    WithField("name")    // Search in specific field

// Available fields: all, iban, name, number, id
```

## Parameter Builders

### Account Parameters

#### ListAccountParamsBuilder

```go
// List all asset accounts with pagination
params, err := NewListAccountParamsBuilder().
    WithAssetAccounts().
    WithLimit(50).
    WithPage(1).
    Build()

// List expense accounts
params, err := NewListAccountParamsBuilder().
    WithExpenseAccounts().
    Build()

// List revenue accounts
params, err := NewListAccountParamsBuilder().
    WithRevenueAccounts().
    Build()

// Custom account type
params, err := NewListAccountParamsBuilder().
    WithType("liability").
    Build()
```

#### SearchAccountParamsBuilder

```go
// Search accounts by name
params, err := NewSearchAccountParamsBuilder("savings").
    WithField("name").
    WithLimit(10).
    Build()

// Search by IBAN
params, err := NewSearchAccountParamsBuilder("DE89").
    WithField("iban").
    Build()

// Search all fields
params, err := NewSearchAccountParamsBuilder("account").
    WithField("all").
    WithPage(2).
    Build()
```

### Transaction Parameters

#### ListTransactionParamsBuilder

```go
// List withdrawals for current month
params, err := NewListTransactionParamsBuilder().
    WithWithdrawals().
    WithCurrentMonth().
    WithLimit(100).
    Build()

// List deposits for specific date range
params, err := NewListTransactionParamsBuilder().
    WithDeposits().
    WithDateRange("2024-01-01", "2024-03-31").
    Build()

// List all transactions for last 7 days
params, err := NewListTransactionParamsBuilder().
    WithLastNDays(7).
    WithPage(1).
    Build()
```

#### SearchTransactionParamsBuilder

```go
// Search transactions
params, err := NewSearchTransactionParamsBuilder("grocery").
    WithLimit(20).
    WithPage(1).
    Build()
```

### Insight Parameters

#### GetSummaryParamsBuilder

```go
// Get summary for current month
params, err := NewGetSummaryParamsBuilder().
    Build()  // Defaults to current month

// Get summary for specific period
params, err := NewGetSummaryParamsBuilder().
    WithDateRange("2024-01-01", "2024-12-31").
    Build()

// Get summary for last 30 days
params, err := NewGetSummaryParamsBuilder().
    WithLastNDays(30).
    Build()
```

#### ExpenseCategoryInsightsParamsBuilder

```go
// Get expense insights by category
params, err := NewExpenseCategoryInsightsParamsBuilder().
    WithDateRange("2024-01-01", "2024-01-31").
    WithAccount("123").
    WithAccount("456").
    Build()

// Using multiple accounts at once
params, err := NewExpenseCategoryInsightsParamsBuilder().
    WithCurrentMonth().
    WithAccounts("123", "456", "789").
    Build()
```

#### ExpenseTotalInsightsParamsBuilder

```go
// Get total expense insights
params, err := NewExpenseTotalInsightsParamsBuilder().
    WithLastNDays(90).
    WithAccounts("123", "456").
    Build()
```

## Integration with Handlers

The builders can be used in handlers for cleaner parameter construction:

```go
// Before (direct parameter construction)
func HandleListAccounts(args ListAccountsArgs) {
    apiParams := &client.ListAccountParams{}
    
    if args.Type != "" {
        filter := client.AccountTypeFilter(args.Type)
        apiParams.Type = &filter
    }
    
    if args.Limit > 0 {
        limit := int32(args.Limit)
        apiParams.Limit = &limit
    }
    
    if args.Page > 0 {
        page := int32(args.Page)
        apiParams.Page = &page
    }
    
    // Use apiParams...
}

// After (using builders)
func HandleListAccounts(args ListAccountsArgs) {
    builder := NewListAccountParamsBuilder()
    
    if args.Type != "" {
        builder.WithType(args.Type)
    }
    
    if args.Limit > 0 {
        builder.WithLimit(args.Limit)
    }
    
    if args.Page > 0 {
        builder.WithPage(args.Page)
    }
    
    apiParams, err := builder.Build()
    if err != nil {
        // Handle validation error
        return
    }
    
    // Use apiParams...
}
```

## Composite Builder

For complex scenarios requiring multiple builders:

```go
composite := NewCompositeBuilder().
    Add(NewPaginationBuilder().WithLimit(25)).
    Add(NewDateRangeBuilder().WithCurrentMonth())

if err := composite.Validate(); err != nil {
    // Handle validation errors from any builder
}
```

## Error Handling

Builders accumulate errors during construction and validate them at build time:

```go
builder := NewListAccountParamsBuilder().
    WithType("invalid_type").    // Invalid type
    WithLimit(0).                 // Invalid limit
    WithPage(-1)                  // Invalid page

params, err := builder.Build()
if err != nil {
    // err will contain all validation errors
    fmt.Println(err)
    // Output: validation errors: [invalid account type: invalid_type, 
    //          limit must be greater than 0, page must be greater than 0]
}
```

## Validation Rules

### Pagination
- Limit must be > 0 and ≤ 1000
- Page must be > 0

### Date Ranges
- Dates must be in YYYY-MM-DD format
- Start date must be before end date
- Both dates required for insights and summaries

### Search
- Query cannot be empty
- Field must be one of: all, iban, name, number, id

### Account Types
- Must be one of: asset, expense, revenue, liability, liabilities, initial_balance, cash

### Transaction Types
- Must be one of: all, withdrawal, withdrawals, expense, deposit, deposits, income, transfer, transfers, opening_balance, reconciliation

## Best Practices

1. **Always check Build() errors**: The Build() method returns an error if validation fails
2. **Use convenience methods**: Prefer `WithAssetAccounts()` over `WithType("asset")`
3. **Chain for readability**: Take advantage of the fluent interface
4. **Validate early**: Call Build() as soon as you've set all parameters
5. **Reuse builders**: Builders can be reused after Build() is called

## Testing

The builders package includes comprehensive tests:

```bash
# Run all builder tests
go test ./pkg/fireflyMCP/builders -v

# Run with coverage
go test ./pkg/fireflyMCP/builders -cover

# Run specific test
go test ./pkg/fireflyMCP/builders -run TestPaginationBuilder
```

## Future Enhancements

- Support for more parameter types as the API grows
- Custom validation rules per deployment
- Builder presets for common scenarios
- Integration with middleware for automatic parameter building