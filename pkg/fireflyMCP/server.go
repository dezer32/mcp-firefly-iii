package fireflyMCP

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// FireflyMCPServer represents the MCP server for Firefly III
type FireflyMCPServer struct {
	server *mcp.Server
	client *client.ClientWithResponses
	config *Config
}

// Tool argument types
type ListAccountsArgs struct {
	Type  string `json:"type,omitempty" jsonschema:"Filter by account type (asset, expense, revenue, etc.)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of accounts to return"`
	Page  int    `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type GetAccountArgs struct {
	ID string `json:"id" jsonschema:"Account ID"`
}

type ListTransactionsArgs struct {
	Type  string `json:"type,omitempty" jsonschema:"Filter by transaction type"`
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type GetTransactionArgs struct {
	ID string `json:"id" jsonschema:"Transaction ID"`
}

type ListBudgetsArgs struct {
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of budgets to return"`
	Page  int    `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type ListCategoriesArgs struct {
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of categories to return"`
	Page  int `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type GetSummaryArgs struct {
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
}

type SearchAccountsArgs struct {
	Query string `json:"query" jsonschema:"The search query"`
	Field string `json:"field" jsonschema:"The account field(s) to search in (all, iban, name, number, id)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of accounts to return"`
	Page  int    `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type SearchTransactionsArgs struct {
	Query string `json:"query" jsonschema:"The search query"`
	Limit int32  `json:"limit,omitempty" jsonschema:"Maximum number of transactions to return"`
	Page  int32  `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
}

type ExpenseCategoryInsightsArgs struct {
	Start    string   `json:"start" jsonschema:"Start date (YYYY-MM-DD)"`
	End      string   `json:"end" jsonschema:"End date (YYYY-MM-DD)"`
	Accounts []string `json:"accounts,omitempty" jsonschema:"Account IDs to include in results"`
}

type ExpenseTotalInsightsArgs struct {
	Start    string   `json:"start" jsonschema:"Start date (YYYY-MM-DD)"`
	End      string   `json:"end" jsonschema:"End date (YYYY-MM-DD)"`
	Accounts []string `json:"accounts,omitempty" jsonschema:"Account IDs to include in results"`
}

type ListBudgetLimitsArgs struct {
	ID    string `json:"id" jsonschema:"Budget ID"`
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
}

type ListBudgetTransactionsArgs struct {
	ID    string `json:"id" jsonschema:"Budget ID"`
	Type  string `json:"type,omitempty" jsonschema:"Filter by transaction type"`
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type ListTagsArgs struct {
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of tags to return"`
	Page  int `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type ListBillsArgs struct {
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of bills to return"`
	Page  int    `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type GetBillArgs struct {
	ID    string `json:"id" jsonschema:"Bill ID"`
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD) for payment info"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD) for payment info"`
}

type ListBillTransactionsArgs struct {
	ID    string `json:"id" jsonschema:"Bill ID"`
	Type  string `json:"type,omitempty" jsonschema:"Filter by transaction type"`
	Start string `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type StoreTransactionArgs struct {
	TransactionStoreRequest
}

type UpdateTransactionArgs struct {
	ID string `json:"id" jsonschema:"Transaction group ID (required)"`
	TransactionUpdateRequest
}

func NewFireflyMCPServer(config *Config) (*FireflyMCPServer, error) {
	// Create HTTP client with authentication
	httpClient := &http.Client{
		Timeout: config.GetTimeout(),
	}

	// Create Firefly III client with request editor for authentication
	fireflyClient, err := client.NewClientWithResponses(
		config.Server.URL,
		client.WithHTTPClient(httpClient),
		client.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				req.Header.Set("Authorization", "Bearer "+config.API.Token)
				req.Header.Set("Accept", "application/vnd.api+json")
				req.Header.Set("Content-Type", "application/json")
				return nil
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firefly III client: %w", err)
	}

	// Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    config.MCP.Name,
			Version: config.MCP.Version,
		}, nil,
	)

	server := &FireflyMCPServer{
		server: mcpServer,
		client: fireflyClient,
		config: config,
	}

	// Register tools
	server.registerTools()

	return server, nil
}

// Run starts the MCP server with the given transport
func (s *FireflyMCPServer) Run(ctx context.Context, transport mcp.Transport) error {
	return s.server.Run(ctx, transport)
}

// registerTools registers all available MCP tools
func (s *FireflyMCPServer) registerTools() {
	// Account tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_accounts",
			Description: "List all accounts in Firefly III",
		}, s.handleListAccounts,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "get_account",
			Description: "Get details of a specific account",
		}, s.handleGetAccount,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "search_accounts",
			Description: "Search for accounts by name, IBAN, or other fields",
		}, s.handleSearchAccounts,
	)

	// Transaction tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_transactions",
			Description: "List transactions in Firefly III",
		}, s.handleListTransactions,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "get_transaction",
			Description: "Get details of a specific transaction",
		}, s.handleGetTransaction,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "search_transactions",
			Description: "Search for transactions by keyword",
		}, s.handleSearchTransactions,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "store_transaction",
			Description: "Create a new transaction in Firefly III",
		}, s.handleStoreTransaction,
	)
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "store_transactions_bulk",
			Description: "Create multiple transaction groups in Firefly III (up to 100 at once)",
		}, s.handleStoreTransactionsBulk,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name: "store_receipt",
			Description: "Create a group of transactions from a single shopping receipt. " +
				"All items are recorded as withdrawals from the source account to the store (expense account). " +
				"Use this to register multiple purchases from one store receipt at once. " +
				"Optionally validates that item amounts sum to the expected total.",
		}, s.handleStoreReceipt,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "update_transaction",
			Description: "Update an existing transaction in Firefly III",
		}, s.handleUpdateTransaction,
	)

	// Budget tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_budgets",
			Description: "List all budgets in Firefly III",
		}, s.handleListBudgets,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_budget_limits",
			Description: "List budget limits for a specific budget with optional date range",
		}, s.handleListBudgetLimits,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_budget_transactions",
			Description: "List transactions for a specific budget with optional filters",
		}, s.handleListBudgetTransactions,
	)

	// Category tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_categories",
			Description: "List all categories in Firefly III",
		}, s.handleListCategories,
	)

	// Tag tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_tags",
			Description: "List all tags in Firefly III",
		}, s.handleListTags,
	)

	// Summary tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "get_summary",
			Description: "Get basic financial summary from Firefly III",
		}, s.handleGetSummary,
	)

	// Insights tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "expense_category_insights",
			Description: "Get expense insights grouped by category for a date range",
		}, s.handleExpenseCategoryInsights,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "expense_total_insights",
			Description: "Get total expense insights for a date range",
		}, s.handleExpenseTotalInsights,
	)

	// Bill tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_bills",
			Description: "List all bills in Firefly III",
		}, s.handleListBills,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "get_bill",
			Description: "Get details of a specific bill",
		}, s.handleGetBill,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_bill_transactions",
			Description: "List transactions associated with a specific bill",
		}, s.handleListBillTransactions,
	)

	// Recurrence tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_recurrences",
			Description: "List all recurrences in Firefly III",
		}, s.handleListRecurrences,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "get_recurrence",
			Description: "Get details of a specific recurrence",
		}, s.handleGetRecurrence,
	)

	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_recurrence_transactions",
			Description: "List transactions created by a specific recurrence",
		}, s.handleListRecurrenceTransactions,
	)
}

// Tool handlers

func (s *FireflyMCPServer) handleListAccounts(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListAccountsArgs,
) (*mcp.CallToolResult, any, error) {
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

	resp, err := s.client.ListAccountWithResponse(ctx, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error listing accounts: %v", err))
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d", resp.StatusCode()))
	}

	// Map response to DTO
	accountList := mapAccountArrayToAccountList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(accountList)
}

func (s *FireflyMCPServer) handleGetAccount(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args GetAccountArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Account ID is required")
	}

	apiParams := &client.GetAccountParams{}
	resp, err := s.client.GetAccountWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error getting account: %v", err))
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d", resp.StatusCode()))
	}

	// Map response to DTO
	account := mapAccountSingleToAccount(resp.ApplicationvndApiJSON200)
	return newSuccessResult(account)
}

func (s *FireflyMCPServer) handleSearchAccounts(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args SearchAccountsArgs,
) (*mcp.CallToolResult, any, error) {
	// Validate required arguments
	if args.Query == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Query parameter is required"},
			},
			IsError: true,
		}, nil, nil
	}

	if args.Field == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Field parameter is required"},
			},
			IsError: true,
		}, nil, nil
	}

	// Build API parameters
	apiParams := &client.SearchAccountsParams{
		Query: args.Query,
		Field: client.AccountSearchFieldFilter(args.Field),
	}

	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	if args.Page > 0 {
		page := int32(args.Page)
		apiParams.Page = &page
	}

	// Call the API
	resp, err := s.client.SearchAccountsWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error searching accounts: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map response to DTO - reuse existing mapper since response type is AccountArray
	accountList := mapAccountArrayToAccountList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(accountList)
}

func (s *FireflyMCPServer) handleListTransactions(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListTransactionsArgs,
) (*mcp.CallToolResult, any, error) {
	apiParams := &client.ListTransactionParams{}

	if args.Type != "" {
		filter := client.TransactionTypeFilter(args.Type)
		apiParams.Type = &filter
	}

	if args.Start != "" {
		if startDate, err := time.Parse("2006-01-02", args.Start); err == nil {
			date := openapi_types.Date{Time: startDate}
			apiParams.Start = &date
		}
	}

	if args.End != "" {
		if endDate, err := time.Parse("2006-01-02", args.End); err == nil {
			date := openapi_types.Date{Time: endDate}
			apiParams.End = &date
		}
	}

	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	if args.Page > 0 {
		page := int32(args.Page)
		apiParams.Page = &page
	}

	resp, err := s.client.ListTransactionWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing transactions: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map response to DTO
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(transactionList)
}

func (s *FireflyMCPServer) handleGetTransaction(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args GetTransactionArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Transaction ID is required"},
			},
			IsError: true,
		}, nil, nil
	}

	apiParams := &client.GetTransactionParams{}
	resp, err := s.client.GetTransactionWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting transaction: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	// Map the response to TransactionGroup DTO
	if resp.StatusCode() == 404 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Transaction not found"},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	transactionGroup := mapTransactionReadToTransactionGroup(&resp.ApplicationvndApiJSON200.Data)
	return newSuccessResult(transactionGroup)
}

func (s *FireflyMCPServer) handleSearchTransactions(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args SearchTransactionsArgs,
) (*mcp.CallToolResult, any, error) {
	// Validate required arguments
	if args.Query == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Query parameter is required"},
			},
			IsError: true,
		}, nil, nil
	}

	// Build API parameters
	apiParams := &client.SearchTransactionsParams{
		Query: args.Query,
		Limit: &args.Limit,
		Page:  &args.Page,
	}

	// Call the API
	resp, err := s.client.SearchTransactionsWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error searching transactions: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map response to DTO - reuse existing mapper since response type is TransactionArray
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(transactionList)
}

func (s *FireflyMCPServer) handleListBudgets(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListBudgetsArgs,
) (*mcp.CallToolResult, any, error) {
	apiParams := &client.ListBudgetParams{}

	// Set default start date to first day of current month
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startDate := openapi_types.Date{Time: firstDayOfMonth}
	apiParams.Start = &startDate

	// Set default end date to last day of current month
	lastDayOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
	endDate := openapi_types.Date{Time: lastDayOfMonth}
	apiParams.End = &endDate

	// Handle custom start date if provided
	if args.Start != "" {
		if customStartDate, err := time.Parse("2006-01-02", args.Start); err == nil {
			date := openapi_types.Date{Time: customStartDate}
			apiParams.Start = &date
		}
	}

	// Handle end date if provided
	if args.End != "" {
		if endDate, err := time.Parse("2006-01-02", args.End); err == nil {
			date := openapi_types.Date{Time: endDate}
			apiParams.End = &date
		}
	}

	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	if args.Page > 0 {
		page := int32(args.Page)
		apiParams.Page = &page
	}

	resp, err := s.client.ListBudgetWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing budgets: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map the response to BudgetList DTO
	budgetList := mapBudgetArrayToBudgetList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(budgetList)
}

func (s *FireflyMCPServer) handleListCategories(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListCategoriesArgs,
) (*mcp.CallToolResult, any, error) {
	apiParams := &client.ListCategoryParams{}

	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	if args.Page > 0 {
		page := int32(args.Page)
		apiParams.Page = &page
	}

	resp, err := s.client.ListCategoryWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing categories: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map the response to CategoryList DTO
	categoryList := mapCategoryArrayToCategoryList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(categoryList)
}

func (s *FireflyMCPServer) handleListTags(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListTagsArgs,
) (*mcp.CallToolResult, any, error) {
	apiParams := &client.ListTagParams{}

	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	if args.Page > 0 {
		page := int32(args.Page)
		apiParams.Page = &page
	}

	resp, err := s.client.ListTagWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing tags: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map the response to TagList DTO
	tagList := mapTagArrayToTagList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(tagList)
}

func (s *FireflyMCPServer) handleGetSummary(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args GetSummaryArgs,
) (*mcp.CallToolResult, any, error) {
	apiParams := &client.GetBasicSummaryParams{}

	// Set default start date to first day of current month
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startDate := openapi_types.Date{Time: firstDayOfMonth}
	apiParams.Start = startDate

	// Set default end date to last day of current month
	lastDayOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
	endDate := openapi_types.Date{Time: lastDayOfMonth}
	apiParams.End = endDate

	if args.Start != "" {
		if startDate, err := time.Parse("2006-01-02", args.Start); err == nil {
			date := openapi_types.Date{Time: startDate}
			apiParams.Start = date
		}
	}

	if args.End != "" {
		if endDate, err := time.Parse("2006-01-02", args.End); err == nil {
			date := openapi_types.Date{Time: endDate}
			apiParams.End = date
		}
	}

	resp, err := s.client.GetBasicSummaryWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting summary: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map response to DTO
	summaryList := mapBasicSummaryToBasicSummaryList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(summaryList)
}

// mapBudgetArrayToBudgetList converts client.BudgetArray to BudgetList DTO
func mapBudgetArrayToBudgetList(budgetArray *client.BudgetArray) *BudgetList {
	if budgetArray == nil {
		return nil
	}

	budgetList := &BudgetList{
		Data: make([]Budget, len(budgetArray.Data)),
	}

	// Map budget data
	for i, budgetRead := range budgetArray.Data {
		budget := Budget{
			Id:     budgetRead.Id,
			Active: budgetRead.Attributes.Active != nil && *budgetRead.Attributes.Active,
			Name:   budgetRead.Attributes.Name,
			Notes:  budgetRead.Attributes.Notes,
		}

		// Map spent information - take only the first value
		if budgetRead.Attributes.Spent != nil && len(*budgetRead.Attributes.Spent) > 0 {
			firstSpent := (*budgetRead.Attributes.Spent)[0]
			budget.Spent = Spent{
				Sum:          getStringValue(firstSpent.Sum),
				CurrencyCode: getStringValue(firstSpent.CurrencyCode),
			}
		}

		budgetList.Data[i] = budget
	}

	// Map pagination
	if budgetArray.Meta.Pagination != nil {
		pagination := budgetArray.Meta.Pagination
		budgetList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return budgetList
}

// mapCategoryArrayToCategoryList converts client.CategoryArray to CategoryList DTO
func mapCategoryArrayToCategoryList(categoryArray *client.CategoryArray) *CategoryList {
	if categoryArray == nil {
		return nil
	}

	categoryList := &CategoryList{
		Data: make([]Category, len(categoryArray.Data)),
	}

	// Map category data
	for i, categoryRead := range categoryArray.Data {
		category := Category{
			Id:    categoryRead.Id,
			Name:  categoryRead.Attributes.Name,
			Notes: categoryRead.Attributes.Notes,
		}

		categoryList.Data[i] = category
	}

	// Map pagination
	if categoryArray.Meta.Pagination != nil {
		pagination := categoryArray.Meta.Pagination
		categoryList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return categoryList
}

// mapTagArrayToTagList converts client.TagArray to TagList DTO
func mapTagArrayToTagList(tagArray *client.TagArray) *TagList {
	if tagArray == nil {
		return nil
	}

	tagList := &TagList{
		Data: make([]Tag, len(tagArray.Data)),
	}

	// Map tag data
	for i, tagRead := range tagArray.Data {
		tagList.Data[i] = Tag{
			Id:          tagRead.Id,
			Tag:         tagRead.Attributes.Tag,
			Description: tagRead.Attributes.Description,
		}
	}

	// Map pagination
	if tagArray.Meta.Pagination != nil {
		pagination := tagArray.Meta.Pagination
		tagList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return tagList
}

// mapAccountArrayToAccountList converts client.AccountArray to AccountList DTO
func mapAccountArrayToAccountList(accountArray *client.AccountArray) *AccountList {
	if accountArray == nil {
		return nil
	}

	accountList := &AccountList{
		Data: make([]Account, len(accountArray.Data)),
	}

	// Map account data
	for i, accountRead := range accountArray.Data {
		account := Account{
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
		accountList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return accountList
}

// mapAccountSingleToAccount converts client.AccountSingle to Account DTO
func mapAccountSingleToAccount(accountSingle *client.AccountSingle) *Account {
	if accountSingle == nil {
		return nil
	}

	return &Account{
		Id:     accountSingle.Data.Id,
		Active: accountSingle.Data.Attributes.Active != nil && *accountSingle.Data.Attributes.Active,
		Name:   accountSingle.Data.Attributes.Name,
		Notes:  accountSingle.Data.Attributes.Notes,
		Type:   string(accountSingle.Data.Attributes.Type),
	}
}

// mapTransactionArrayToTransactionList converts client.TransactionArray to TransactionList DTO
func mapTransactionArrayToTransactionList(transactionArray *client.TransactionArray) *TransactionList {
	if transactionArray == nil {
		return nil
	}

	transactionList := &TransactionList{
		Data: make([]TransactionGroup, len(transactionArray.Data)),
	}

	// Map transaction data
	for i, transactionRead := range transactionArray.Data {
		transactionList.Data[i] = *mapTransactionReadToTransactionGroup(&transactionRead)
	}

	// Map pagination
	if transactionArray.Meta.Pagination != nil {
		pagination := transactionArray.Meta.Pagination
		transactionList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return transactionList
}

// mapTransactionReadToTransactionGroup converts client.TransactionRead to TransactionGroup DTO
func mapTransactionReadToTransactionGroup(transactionRead *client.TransactionRead) *TransactionGroup {
	if transactionRead == nil {
		return nil
	}

	group := &TransactionGroup{
		Id:           transactionRead.Id,
		GroupTitle:   getStringValue(transactionRead.Attributes.GroupTitle),
		Transactions: make([]Transaction, len(transactionRead.Attributes.Transactions)),
	}

	// Map individual transactions within the group
	for i, split := range transactionRead.Attributes.Transactions {
		transaction := Transaction{
			Id:              getStringValue(split.TransactionJournalId),
			Amount:          split.Amount,
			BillId:          split.BillId,
			BillName:        split.BillName,
			BudgetId:        split.BudgetId,
			BudgetName:      split.BudgetName,
			CategoryId:      split.CategoryId,
			CategoryName:    split.CategoryName,
			CurrencyCode:    getStringValue(split.CurrencyCode),
			Date:            split.Date,
			Description:     split.Description,
			DestinationId:   getStringValue(split.DestinationId),
			DestinationName: getStringValue(split.DestinationName),
			DestinationType: string(getAccountTypeValue(split.DestinationType)),
			Notes:           split.Notes,
			Reconciled:      split.Reconciled != nil && *split.Reconciled,
			SourceId:        getStringValue(split.SourceId),
			SourceName:      getStringValue(split.SourceName),
			Type:            string(split.Type),
		}

		// Handle tags
		if split.Tags != nil && len(*split.Tags) > 0 {
			transaction.Tags = *split.Tags
		} else {
			transaction.Tags = []string{}
		}

		group.Transactions[i] = transaction
	}

	return group
}

// mapBasicSummaryToBasicSummaryList converts client.BasicSummary to BasicSummaryList DTO
func mapBasicSummaryToBasicSummaryList(basicSummary *client.BasicSummary) *BasicSummaryList {
	if basicSummary == nil {
		return &BasicSummaryList{
			Data: []BasicSummary{},
		}
	}

	summaryList := &BasicSummaryList{
		Data: make([]BasicSummary, 0, len(*basicSummary)),
	}

	// Convert the map to a slice of BasicSummary DTOs
	for _, entry := range *basicSummary {
		summary := BasicSummary{
			Key:           getStringValue(entry.Key),
			Title:         getStringValue(entry.Title),
			CurrencyCode:  getStringValue(entry.CurrencyCode),
			MonetaryValue: getStringValue(entry.MonetaryValue),
		}
		summaryList.Data = append(summaryList.Data, summary)
	}

	return summaryList
}

// handleExpenseCategoryInsights returns expense insights grouped by category
func (s *FireflyMCPServer) handleExpenseCategoryInsights(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ExpenseCategoryInsightsArgs,
) (*mcp.CallToolResult, any, error) {
	// Validate required dates
	if args.Start == "" || args.End == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Start and End dates are required"},
			},
			IsError: true,
		}, nil, nil
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", args.Start)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	endDate, err := time.Parse("2006-01-02", args.End)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	// Build API parameters
	apiParams := &client.InsightExpenseCategoryParams{
		Start: openapi_types.Date{Time: startDate},
		End:   openapi_types.Date{Time: endDate},
	}

	// Convert account IDs from strings to int64
	if len(args.Accounts) > 0 {
		accounts := make([]int64, len(args.Accounts))
		for i, accStr := range args.Accounts {
			var accID int64
			if _, err := fmt.Sscanf(accStr, "%d", &accID); err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Invalid account ID: %s", accStr)},
					},
					IsError: true,
				}, nil, nil
			}
			accounts[i] = accID
		}
		apiParams.Accounts = &accounts
	}

	// Call the API
	resp, err := s.client.InsightExpenseCategoryWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting expense category insights: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map response to DTO
	insightResponse := mapInsightGroupToDTO(resp.JSON200)
	return newSuccessResult(insightResponse)
}

// handleExpenseTotalInsights returns total expense insights
func (s *FireflyMCPServer) handleExpenseTotalInsights(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ExpenseTotalInsightsArgs,
) (*mcp.CallToolResult, any, error) {
	// Validate required dates
	if args.Start == "" || args.End == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Start and End dates are required"},
			},
			IsError: true,
		}, nil, nil
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", args.Start)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	endDate, err := time.Parse("2006-01-02", args.End)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	// Build API parameters
	apiParams := &client.InsightExpenseTotalParams{
		Start: openapi_types.Date{Time: startDate},
		End:   openapi_types.Date{Time: endDate},
	}

	// Convert account IDs from strings to int64
	if len(args.Accounts) > 0 {
		accounts := make([]int64, len(args.Accounts))
		for i, accStr := range args.Accounts {
			var accID int64
			if _, err := fmt.Sscanf(accStr, "%d", &accID); err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Invalid account ID: %s", accStr)},
					},
					IsError: true,
				}, nil, nil
			}
			accounts[i] = accID
		}
		apiParams.Accounts = &accounts
	}

	// Call the API
	resp, err := s.client.InsightExpenseTotalWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting expense total insights: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map response to DTO
	insightResponse := mapInsightTotalToDTO(resp.JSON200)
	return newSuccessResult(insightResponse)
}

// handleListBudgetLimits returns budget limits for a specific budget
func (s *FireflyMCPServer) handleListBudgetLimits(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListBudgetLimitsArgs,
) (*mcp.CallToolResult, any, error) {
	// Validate required budget ID
	if args.ID == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Budget ID is required"},
			},
			IsError: true,
		}, nil, nil
	}

	// Build API parameters
	apiParams := &client.ListBudgetLimitByBudgetParams{}

	// Parse optional start date
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	// Parse optional end date
	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	// Call the API
	resp, err := s.client.ListBudgetLimitByBudgetWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing budget limits: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map response to DTO
	budgetLimitList := mapBudgetLimitArrayToBudgetLimitList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(budgetLimitList)
}

func (s *FireflyMCPServer) handleListBudgetTransactions(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListBudgetTransactionsArgs,
) (*mcp.CallToolResult, any, error) {
	// Validate required budget ID
	if args.ID == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Budget ID is required"},
			},
			IsError: true,
		}, nil, nil
	}

	// Build API parameters
	apiParams := &client.ListTransactionByBudgetParams{}

	// Set pagination parameters
	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	if args.Page > 0 {
		page := int32(args.Page)
		apiParams.Page = &page
	}

	// Parse optional start date
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	// Parse optional end date
	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	// Set transaction type filter if provided
	if args.Type != "" {
		typeFilter := client.TransactionTypeFilter(args.Type)
		apiParams.Type = &typeFilter
	}

	// Call the API
	resp, err := s.client.ListTransactionByBudgetWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing budget transactions: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil, nil
	}

	// Map response to DTO
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(transactionList)
}

// mapInsightGroupToDTO converts client.InsightGroup to InsightCategoryResponse DTO
func mapInsightGroupToDTO(group *client.InsightGroup) *InsightCategoryResponse {
	if group == nil {
		return &InsightCategoryResponse{
			Entries: []InsightCategoryEntry{},
		}
	}

	response := &InsightCategoryResponse{
		Entries: make([]InsightCategoryEntry, 0, len(*group)),
	}

	for _, entry := range *group {
		// Use difference as the amount since it represents the expense
		categoryEntry := InsightCategoryEntry{
			Id:           getStringValue(entry.Id),
			Name:         getStringValue(entry.Name),
			Amount:       getStringValue(entry.Difference),
			CurrencyCode: getStringValue(entry.CurrencyCode),
		}
		response.Entries = append(response.Entries, categoryEntry)
	}

	return response
}

// mapInsightTotalToDTO converts client.InsightTotal to InsightTotalResponse DTO
func mapInsightTotalToDTO(total *client.InsightTotal) *InsightTotalResponse {
	if total == nil {
		return &InsightTotalResponse{
			Entries: []InsightTotalEntry{},
		}
	}

	response := &InsightTotalResponse{
		Entries: make([]InsightTotalEntry, 0, len(*total)),
	}

	for _, entry := range *total {
		totalEntry := InsightTotalEntry{
			Amount:       getStringValue(entry.Difference),
			CurrencyCode: getStringValue(entry.CurrencyCode),
		}
		response.Entries = append(response.Entries, totalEntry)
	}

	return response
}

// mapBudgetLimitArrayToBudgetLimitList converts client.BudgetLimitArray to BudgetLimitList DTO
func mapBudgetLimitArrayToBudgetLimitList(budgetLimitArray *client.BudgetLimitArray) *BudgetLimitList {
	if budgetLimitArray == nil {
		return nil
	}

	budgetLimitList := &BudgetLimitList{
		Data: make([]BudgetLimit, len(budgetLimitArray.Data)),
	}

	// Map budget limit data
	for i, budgetLimitRead := range budgetLimitArray.Data {
		budgetLimit := BudgetLimit{
			Id:             budgetLimitRead.Id,
			Amount:         budgetLimitRead.Attributes.Amount,
			Start:          budgetLimitRead.Attributes.Start,
			End:            budgetLimitRead.Attributes.End,
			BudgetId:       getStringValue(budgetLimitRead.Attributes.BudgetId),
			CurrencyCode:   getStringValue(budgetLimitRead.Attributes.CurrencyCode),
			CurrencySymbol: getStringValue(budgetLimitRead.Attributes.CurrencySymbol),
			Spent:          make([]BudgetSpent, 0),
		}

		// Map spent data if available
		if budgetLimitRead.Attributes.Spent != nil {
			budgetSpent := BudgetSpent{
				Sum:            getStringValue(budgetLimitRead.Attributes.Spent),
				CurrencyCode:   getStringValue(budgetLimitRead.Attributes.CurrencyCode),
				CurrencySymbol: getStringValue(budgetLimitRead.Attributes.CurrencySymbol),
			}
			budgetLimit.Spent = append(budgetLimit.Spent, budgetSpent)
		}

		budgetLimitList.Data[i] = budgetLimit
	}

	// Map pagination
	if budgetLimitArray.Meta.Pagination != nil {
		pagination := budgetLimitArray.Meta.Pagination
		budgetLimitList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return budgetLimitList
}

// mapBillToBill converts a single Bill from API to DTO
func mapBillToBill(billRead *client.BillRead) *Bill {
	if billRead == nil {
		return nil
	}

	bill := &Bill{
		Id:                billRead.Id,
		Active:            billRead.Attributes.Active != nil && *billRead.Attributes.Active,
		Name:              billRead.Attributes.Name,
		AmountMin:         billRead.Attributes.AmountMin,
		AmountMax:         billRead.Attributes.AmountMax,
		Date:              billRead.Attributes.Date,
		RepeatFreq:        string(billRead.Attributes.RepeatFreq),
		Skip:              0,
		CurrencyCode:      getStringValue(billRead.Attributes.CurrencyCode),
		Notes:             billRead.Attributes.Notes,
		NextExpectedMatch: billRead.Attributes.NextExpectedMatch,
		PaidDates:         []PaidDate{},
	}

	// Handle skip field
	if billRead.Attributes.Skip != nil {
		bill.Skip = int(*billRead.Attributes.Skip)
	}

	// Map paid dates if available
	if billRead.Attributes.PaidDates != nil {
		for _, pd := range *billRead.Attributes.PaidDates {
			paidDate := PaidDate{
				Date:                 pd.Date,
				TransactionGroupId:   pd.TransactionGroupId,
				TransactionJournalId: pd.TransactionJournalId,
			}
			bill.PaidDates = append(bill.PaidDates, paidDate)
		}
	}

	return bill
}

// mapBillArrayToBillList converts client.BillArray to BillList DTO
func mapBillArrayToBillList(billArray *client.BillArray) *BillList {
	if billArray == nil {
		return nil
	}

	billList := &BillList{
		Data: make([]Bill, 0, len(billArray.Data)),
	}

	// Map bill data
	for _, billRead := range billArray.Data {
		if mappedBill := mapBillToBill(&billRead); mappedBill != nil {
			billList.Data = append(billList.Data, *mappedBill)
		}
	}

	// Map pagination
	if billArray.Meta.Pagination != nil {
		pagination := billArray.Meta.Pagination
		billList.Pagination = Pagination{
			Count:       getIntValue(pagination.Count),
			Total:       getIntValue(pagination.Total),
			CurrentPage: getIntValue(pagination.CurrentPage),
			PerPage:     getIntValue(pagination.PerPage),
			TotalPages:  getIntValue(pagination.TotalPages),
		}
	}

	return billList
}

// handleListBills lists all bills in Firefly III
func (s *FireflyMCPServer) handleListBills(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListBillsArgs,
) (*mcp.CallToolResult, any, error) {
	// Prepare API parameters
	apiParams := &client.ListBillParams{}

	// Set pagination
	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}
	page := int32(args.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	// Set date filters if provided
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	// Call the API
	resp, err := s.client.ListBillWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing bills: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %s", string(resp.Body))},
			},
			IsError: true,
		}, nil, nil
	}

	// Map the response
	billList := mapBillArrayToBillList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(billList)
}

// handleGetBill gets a specific bill by ID
func (s *FireflyMCPServer) handleGetBill(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args GetBillArgs,
) (*mcp.CallToolResult, any, error) {
	// Prepare API parameters
	apiParams := &client.GetBillParams{}

	// Set date filters if provided
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	// Call the API
	resp, err := s.client.GetBillWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting bill: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %s", string(resp.Body))},
			},
			IsError: true,
		}, nil, nil
	}

	// Map the response
	var bill *Bill
	if resp.ApplicationvndApiJSON200 != nil {
		bill = mapBillToBill(&resp.ApplicationvndApiJSON200.Data)
	}
	return newSuccessResult(bill)
}

// handleListBillTransactions lists transactions associated with a specific bill
func (s *FireflyMCPServer) handleListBillTransactions(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListBillTransactionsArgs,
) (*mcp.CallToolResult, any, error) {
	// Prepare API parameters
	apiParams := &client.ListTransactionByBillParams{}

	// Set pagination
	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}
	page := int32(args.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	// Set date filters if provided
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	// Set transaction type filter if provided
	if args.Type != "" {
		typeFilter := client.TransactionTypeFilter(args.Type)
		apiParams.Type = &typeFilter
	}

	// Call the API
	resp, err := s.client.ListTransactionByBillWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing bill transactions: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %s", string(resp.Body))},
			},
			IsError: true,
		}, nil, nil
	}

	// Map the response
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(transactionList)
}

// getAccountTypeValue safely extracts AccountTypeProperty value, returns empty string if nil
func getAccountTypeValue(ptr *client.AccountTypeProperty) client.AccountTypeProperty {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// getIntValue safely extracts int value from pointer, returns 0 if nil
func getIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// getStringValue safely extracts string value from pointer, returns empty string if nil
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// getRecurrenceType converts RecurrenceTransactionType to string
func getRecurrenceType(t *client.RecurrenceTransactionType) string {
	if t == nil {
		return ""
	}
	return string(*t)
}

// fixCurrencyIdFields converts numeric currency_id values to strings in JSON response
// This fixes the JSON unmarshaling error where API returns numbers but structs expect strings
func fixCurrencyIdFields(jsonStr string) string {
	// Pattern to match "currency_id": <number> and convert to "currency_id": "<number>"
	re := regexp.MustCompile(`"currency_id":\s*(\d+)`)
	return re.ReplaceAllString(jsonStr, `"currency_id": "$1"`)
}
