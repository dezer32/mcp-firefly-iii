package builders_test

import (
	"fmt"
	"log"

	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/builders"
)

func ExampleListAccountParamsBuilder() {
	// Build parameters for listing asset accounts with pagination
	builder := builders.NewListAccountParamsBuilder()
	builder.WithAssetAccounts()
	builder.WithLimit(25)
	builder.WithPage(1)
	params, err := builder.Build()
	
	if err != nil {
		log.Fatal(err)
	}
	
	// Use params with the API client
	fmt.Printf("Type: %v, Limit: %v, Page: %v\n", 
		params.Type != nil, 
		params.Limit != nil, 
		params.Page != nil)
	// Output: Type: true, Limit: true, Page: true
}

func ExampleSearchAccountParamsBuilder() {
	// Search for accounts with "savings" in the name
	builder := builders.NewSearchAccountParamsBuilder("savings")
	builder.WithField("name")
	builder.WithLimit(10)
	params, err := builder.Build()
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Query: %s, Field: %s\n", params.Query, params.Field)
	// Output: Query: savings, Field: name
}

func ExampleListTransactionParamsBuilder() {
	// List withdrawals for the current month
	builder := builders.NewListTransactionParamsBuilder()
	builder.WithWithdrawals()
	builder.WithCurrentMonth()
	builder.WithLimit(50)
	params, err := builder.Build()
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Has Type: %v, Has Date Range: %v\n",
		params.Type != nil,
		params.Start != nil && params.End != nil)
	// Output: Has Type: true, Has Date Range: true
}

func ExampleDateRangeBuilder_WithLastNDays() {
	// Create a date range for the last 7 days
	builder := builders.NewDateRangeBuilder().
		WithLastNDays(7)
	
	if err := builder.Validate(); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Date range created for last 7 days")
	// Output: Date range created for last 7 days
}

func ExamplePaginationBuilder() {
	// Create pagination parameters
	builder := builders.NewPaginationBuilder().
		WithLimit(100).
		WithPage(3)
	
	if err := builder.Validate(); err != nil {
		log.Fatal(err)
	}
	
	limit := builder.GetLimit()
	page := builder.GetPage()
	
	fmt.Printf("Limit: %d, Page: %d\n", *limit, *page)
	// Output: Limit: 100, Page: 3
}

func ExampleCompositeBuilder() {
	// Combine multiple builders for validation
	composite := builders.NewCompositeBuilder().
		Add(builders.NewPaginationBuilder().WithLimit(25)).
		Add(builders.NewDateRangeBuilder().WithCurrentMonth())
	
	if err := composite.Validate(); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("All builders validated successfully")
	// Output: All builders validated successfully
}

func ExampleGetSummaryParamsBuilder() {
	// Get summary for a specific date range
	builder := builders.NewGetSummaryParamsBuilder()
	builder.WithDateRange("2024-01-01", "2024-12-31")
	params, err := builder.Build()
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Start: %s, End: %s\n",
		params.Start.Time.Format("2006-01-02"),
		params.End.Time.Format("2006-01-02"))
	// Output: Start: 2024-01-01, End: 2024-12-31
}

func ExampleExpenseCategoryInsightsParamsBuilder() {
	// Get expense insights for specific accounts
	builder := builders.NewExpenseCategoryInsightsParamsBuilder()
	builder.WithCurrentMonth()
	builder.WithAccounts(1, 2, 3)
	params, err := builder.Build()
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Has date range: %v, Account count: %d\n",
		params.Start.Time.IsZero() == false,
		len(*params.Accounts))
	// Output: Has date range: true, Account count: 3
}