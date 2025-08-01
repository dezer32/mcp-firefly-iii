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
}

type GetAccountArgs struct {
	ID string `json:"id" mcp:"Account ID"`
}

type ListTransactionsArgs struct {
	Type  string `json:"type,omitempty" mcp:"Filter by transaction type"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
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

	// Budget tools
	mcp.AddTool(
		s.server, &mcp.Tool{
			Name:        "list_budgets",
			Description: "List all budgets in Firefly III",
		}, s.handleListBudgets,
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
}

// Tool handlers

func (s *FireflyMCPServer) handleListAccounts(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListAccountsArgs]) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListAccountParams{}

	if params.Arguments.Type != "" {
		filter := client.AccountTypeFilter(params.Arguments.Type)
		apiParams.Type = &filter
	}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
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

	// Format response
	result, _ := json.MarshalIndent(resp.ApplicationvndApiJSON200, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleGetAccount(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetAccountArgs]) (*mcp.CallToolResultFor[struct{}], error) {
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

	result, _ := json.MarshalIndent(resp.ApplicationvndApiJSON200, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleListTransactions(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListTransactionsArgs]) (*mcp.CallToolResultFor[struct{}], error) {
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

	result, _ := json.MarshalIndent(resp.ApplicationvndApiJSON200, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleGetTransaction(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetTransactionArgs]) (*mcp.CallToolResultFor[struct{}], error) {
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

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %d", resp.StatusCode())},
			},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(resp.ApplicationvndApiJSON200, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func (s *FireflyMCPServer) handleListBudgets(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListBudgetsArgs]) (*mcp.CallToolResultFor[struct{}], error) {
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

func (s *FireflyMCPServer) handleListCategories(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ListCategoriesArgs]) (*mcp.CallToolResultFor[struct{}], error) {
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

func (s *FireflyMCPServer) handleGetSummary(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GetSummaryArgs]) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.GetBasicSummaryParams{}

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

	result, _ := json.MarshalIndent(resp.ApplicationvndApiJSON200, "", "  ")
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
