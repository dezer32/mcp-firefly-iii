package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/dezer32/firefly-iii/pkg/client"
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
	Type  string `json:"type,omitempty" mcp:"Filter by account type (asset, expense, revenue, etc.)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of accounts to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

type GetAccountArgs struct {
	ID string `json:"id" mcp:"Account ID"`
}

type ListTransactionsArgs struct {
	Type  string `json:"type,omitempty" mcp:"Filter by transaction type"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

type GetTransactionArgs struct {
	ID string `json:"id" mcp:"Transaction ID"`
}

type ListBudgetsArgs struct {
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of budgets to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

type ListCategoriesArgs struct {
	Limit int `json:"limit,omitempty" mcp:"Maximum number of categories to return"`
	Page  int `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

type GetSummaryArgs struct {
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
}

type SearchAccountsArgs struct {
	Query string `json:"query" mcp:"The search query"`
	Field string `json:"field" mcp:"The account field(s) to search in (all, iban, name, number, id)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of accounts to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

type SearchTransactionsArgs struct {
	Query string `json:"query" mcp:"The search query"`
	Limit int32  `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int32  `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
}

type ExpenseCategoryInsightsArgs struct {
	Start    string   `json:"start" mcp:"Start date (YYYY-MM-DD)"`
	End      string   `json:"end" mcp:"End date (YYYY-MM-DD)"`
	Accounts []string `json:"accounts,omitempty" mcp:"Account IDs to include in results"`
}

type ExpenseTotalInsightsArgs struct {
	Start    string   `json:"start" mcp:"Start date (YYYY-MM-DD)"`
	End      string   `json:"end" mcp:"End date (YYYY-MM-DD)"`
	Accounts []string `json:"accounts,omitempty" mcp:"Account IDs to include in results"`
}

type ListBudgetLimitsArgs struct {
	ID    string `json:"id" mcp:"Budget ID"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
}

type ListBudgetTransactionsArgs struct {
	ID    string `json:"id" mcp:"Budget ID"`
	Type  string `json:"type,omitempty" mcp:"Filter by transaction type"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// NewFireflyMCPServer creates a new Firefly III MCP server
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
}

// Tool handlers

func (s *FireflyMCPServer) handleListAccounts(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListAccountsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListAccountParams{}

	if params.Arguments.Type != "" {
		filter := client.AccountTypeFilter(params.Arguments.Type)
		apiParams.Type = &filter
	}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	resp, err := s.client.ListAccountWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing accounts: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO
	accountList := mapAccountArrayToAccountList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(accountList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleGetAccount(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GetAccountArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	if params.Arguments.ID == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Account ID is required"},
			},
			IsError: true,
		}, nil
	}

	apiParams := &client.GetAccountParams{}
	resp, err := s.client.GetAccountWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting account: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO
	account := mapAccountSingleToAccount(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(account, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleSearchAccounts(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[SearchAccountsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Validate required arguments
	if params.Arguments.Query == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Query parameter is required"},
			},
			IsError: true,
		}, nil
	}

	if params.Arguments.Field == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Field parameter is required"},
			},
			IsError: true,
		}, nil
	}

	// Build API parameters
	apiParams := &client.SearchAccountsParams{
		Query: params.Arguments.Query,
		Field: client.AccountSearchFieldFilter(params.Arguments.Field),
	}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	// Call the API
	resp, err := s.client.SearchAccountsWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error searching accounts: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO - reuse existing mapper since response type is AccountArray
	accountList := mapAccountArrayToAccountList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(accountList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleListTransactions(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListTransactionsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListTransactionParams{}

	if params.Arguments.Type != "" {
		filter := client.TransactionTypeFilter(params.Arguments.Type)
		apiParams.Type = &filter
	}

	if params.Arguments.Start != "" {
		if startDate, err := time.Parse("2006-01-02", params.Arguments.Start); err == nil {
			date := openapi_types.Date{Time: startDate}
			apiParams.Start = &date
		}
	}

	if params.Arguments.End != "" {
		if endDate, err := time.Parse("2006-01-02", params.Arguments.End); err == nil {
			date := openapi_types.Date{Time: endDate}
			apiParams.End = &date
		}
	}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	resp, err := s.client.ListTransactionWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing transactions: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(transactionList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleGetTransaction(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GetTransactionArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	if params.Arguments.ID == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Transaction ID is required"},
			},
			IsError: true,
		}, nil
	}

	apiParams := &client.GetTransactionParams{}
	resp, err := s.client.GetTransactionWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting transaction: %v", err)},
			},
			IsError: true,
		}, nil
	}

	// Map the response to TransactionGroup DTO
	if resp.StatusCode() == 404 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Transaction not found"},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	transactionGroup := mapTransactionReadToTransactionGroup(&resp.ApplicationvndApiJSON200.Data)

	result, _ := json.MarshalIndent(transactionGroup, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleSearchTransactions(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[SearchTransactionsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Validate required arguments
	if params.Arguments.Query == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Query parameter is required"},
			},
			IsError: true,
		}, nil
	}

	// Build API parameters
	apiParams := &client.SearchTransactionsParams{
		Query: params.Arguments.Query,
		Limit: &params.Arguments.Limit,
		Page:  &params.Arguments.Page,
	}

	// Call the API
	resp, err := s.client.SearchTransactionsWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error searching transactions: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO - reuse existing mapper since response type is TransactionArray
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(transactionList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleListBudgets(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListBudgetsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
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
	if params.Arguments.Start != "" {
		if customStartDate, err := time.Parse("2006-01-02", params.Arguments.Start); err == nil {
			date := openapi_types.Date{Time: customStartDate}
			apiParams.Start = &date
		}
	}

	// Handle end date if provided
	if params.Arguments.End != "" {
		if endDate, err := time.Parse("2006-01-02", params.Arguments.End); err == nil {
			date := openapi_types.Date{Time: endDate}
			apiParams.End = &date
		}
	}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	resp, err := s.client.ListBudgetWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing budgets: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map the response to BudgetList DTO
	budgetList := mapBudgetArrayToBudgetList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(budgetList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleListCategories(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListCategoriesArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListCategoryParams{}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	resp, err := s.client.ListCategoryWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing categories: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map the response to CategoryList DTO
	categoryList := mapCategoryArrayToCategoryList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(categoryList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleGetSummary(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GetSummaryArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
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

	if params.Arguments.Start != "" {
		if startDate, err := time.Parse("2006-01-02", params.Arguments.Start); err == nil {
			date := openapi_types.Date{Time: startDate}
			apiParams.Start = date
		}
	}

	if params.Arguments.End != "" {
		if endDate, err := time.Parse("2006-01-02", params.Arguments.End); err == nil {
			date := openapi_types.Date{Time: endDate}
			apiParams.End = date
		}
	}

	resp, err := s.client.GetBasicSummaryWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting summary: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO
	summaryList := mapBasicSummaryToBasicSummaryList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(summaryList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
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
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ExpenseCategoryInsightsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Validate required dates
	if params.Arguments.Start == "" || params.Arguments.End == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Start and End dates are required"},
			},
			IsError: true,
		}, nil
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", params.Arguments.Start)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
			},
			IsError: true,
		}, nil
	}

	endDate, err := time.Parse("2006-01-02", params.Arguments.End)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
			},
			IsError: true,
		}, nil
	}

	// Build API parameters
	apiParams := &client.InsightExpenseCategoryParams{
		Start: openapi_types.Date{Time: startDate},
		End:   openapi_types.Date{Time: endDate},
	}

	// Convert account IDs from strings to int64
	if len(params.Arguments.Accounts) > 0 {
		accounts := make([]int64, len(params.Arguments.Accounts))
		for i, accStr := range params.Arguments.Accounts {
			var accID int64
			if _, err := fmt.Sscanf(accStr, "%d", &accID); err != nil {
				return &mcp.CallToolResultFor[struct{}]{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Invalid account ID: %s", accStr)},
					},
					IsError: true,
				}, nil
			}
			accounts[i] = accID
		}
		apiParams.Accounts = &accounts
	}

	// Call the API
	resp, err := s.client.InsightExpenseCategoryWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting expense category insights: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO
	insightResponse := mapInsightGroupToDTO(resp.JSON200)
	result, _ := json.MarshalIndent(insightResponse, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// handleExpenseTotalInsights returns total expense insights
func (s *FireflyMCPServer) handleExpenseTotalInsights(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ExpenseTotalInsightsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Validate required dates
	if params.Arguments.Start == "" || params.Arguments.End == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Start and End dates are required"},
			},
			IsError: true,
		}, nil
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", params.Arguments.Start)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
			},
			IsError: true,
		}, nil
	}

	endDate, err := time.Parse("2006-01-02", params.Arguments.End)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
			},
			IsError: true,
		}, nil
	}

	// Build API parameters
	apiParams := &client.InsightExpenseTotalParams{
		Start: openapi_types.Date{Time: startDate},
		End:   openapi_types.Date{Time: endDate},
	}

	// Convert account IDs from strings to int64
	if len(params.Arguments.Accounts) > 0 {
		accounts := make([]int64, len(params.Arguments.Accounts))
		for i, accStr := range params.Arguments.Accounts {
			var accID int64
			if _, err := fmt.Sscanf(accStr, "%d", &accID); err != nil {
				return &mcp.CallToolResultFor[struct{}]{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Invalid account ID: %s", accStr)},
					},
					IsError: true,
				}, nil
			}
			accounts[i] = accID
		}
		apiParams.Accounts = &accounts
	}

	// Call the API
	resp, err := s.client.InsightExpenseTotalWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting expense total insights: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO
	insightResponse := mapInsightTotalToDTO(resp.JSON200)
	result, _ := json.MarshalIndent(insightResponse, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// handleListBudgetLimits returns budget limits for a specific budget
func (s *FireflyMCPServer) handleListBudgetLimits(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListBudgetLimitsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Validate required budget ID
	if params.Arguments.ID == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Budget ID is required"},
			},
			IsError: true,
		}, nil
	}

	// Build API parameters
	apiParams := &client.ListBudgetLimitByBudgetParams{}

	// Parse optional start date
	if params.Arguments.Start != "" {
		startDate, err := time.Parse("2006-01-02", params.Arguments.Start)
		if err != nil {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
				},
				IsError: true,
			}, nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	// Parse optional end date
	if params.Arguments.End != "" {
		endDate, err := time.Parse("2006-01-02", params.Arguments.End)
		if err != nil {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
				},
				IsError: true,
			}, nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	// Call the API
	resp, err := s.client.ListBudgetLimitByBudgetWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing budget limits: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO
	budgetLimitList := mapBudgetLimitArrayToBudgetLimitList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(budgetLimitList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleListBudgetTransactions(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListBudgetTransactionsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Validate required budget ID
	if params.Arguments.ID == "" {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Budget ID is required"},
			},
			IsError: true,
		}, nil
	}

	// Build API parameters
	apiParams := &client.ListTransactionByBudgetParams{}

	// Set pagination parameters
	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	// Parse optional start date
	if params.Arguments.Start != "" {
		startDate, err := time.Parse("2006-01-02", params.Arguments.Start)
		if err != nil {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
				},
				IsError: true,
			}, nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	// Parse optional end date
	if params.Arguments.End != "" {
		endDate, err := time.Parse("2006-01-02", params.Arguments.End)
		if err != nil {
			return &mcp.CallToolResultFor[struct{}]{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
				},
				IsError: true,
			}, nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	// Set transaction type filter if provided
	if params.Arguments.Type != "" {
		typeFilter := client.TransactionTypeFilter(params.Arguments.Type)
		apiParams.Type = &typeFilter
	}

	// Call the API
	resp, err := s.client.ListTransactionByBudgetWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing budget transactions: %v", err)},
			},
			IsError: true,
		}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	// Map response to DTO
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(transactionList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
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

// fixCurrencyIdFields converts numeric currency_id values to strings in JSON response
// This fixes the JSON unmarshaling error where API returns numbers but structs expect strings
func fixCurrencyIdFields(jsonStr string) string {
	// Pattern to match "currency_id": <number> and convert to "currency_id": "<number>"
	re := regexp.MustCompile(`"currency_id":\s*(\d+)`)
	return re.ReplaceAllString(jsonStr, `"currency_id": "$1"`)
}
