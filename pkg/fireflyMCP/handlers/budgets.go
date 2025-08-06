package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/mappers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// BudgetHandlers handles all budget-related MCP tools
type BudgetHandlers struct {
	*BaseHandler
	client *client.ClientWithResponses
}

// NewBudgetHandlers creates a new BudgetHandlers instance
func NewBudgetHandlers(ctx HandlerContext) *BudgetHandlers {
	return &BudgetHandlers{
		BaseHandler: NewBaseHandler(ctx),
		client:      ctx.GetClient().(*client.ClientWithResponses),
	}
}

// ListBudgetsArgs defines arguments for listing budgets
type ListBudgetsArgs struct {
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of budgets to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// ListBudgetLimitsArgs defines arguments for listing budget limits
type ListBudgetLimitsArgs struct {
	ID    string `json:"id" mcp:"Budget ID"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
}

// ListBudgetTransactionsArgs defines arguments for listing budget transactions
type ListBudgetTransactionsArgs struct {
	ID    string `json:"id" mcp:"Budget ID"`
	Type  string `json:"type,omitempty" mcp:"Filter by transaction type"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// HandleListBudgets handles the list_budgets tool
func (h *BudgetHandlers) HandleListBudgets(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListBudgetsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListBudgetParams{}

	// Set date range
	var startDate, endDate openapi_types.Date
	if params.Arguments.Start != "" {
		parsedStart, err := time.Parse("2006-01-02", params.Arguments.Start)
		if err != nil {
			return h.HandleError(err, "Invalid start date format"), nil
		}
		startDate = openapi_types.Date{Time: parsedStart}
		apiParams.Start = &startDate
	}

	if params.Arguments.End != "" {
		parsedEnd, err := time.Parse("2006-01-02", params.Arguments.End)
		if err != nil {
			return h.HandleError(err, "Invalid end date format"), nil
		}
		endDate = openapi_types.Date{Time: parsedEnd}
		apiParams.End = &endDate
	}

	// Use defaults if not provided
	if apiParams.Start == nil {
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = openapi_types.Date{Time: startOfMonth}
		apiParams.Start = &startDate
	}

	if apiParams.End == nil {
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		endOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 23, 59, 59, 999999999, now.Location())
		endDate = openapi_types.Date{Time: endOfMonth}
		apiParams.End = &endDate
	}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	resp, err := h.client.ListBudgetWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error listing budgets"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	budgetList := mappers.MapBudgetArrayToBudgetList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(budgetList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleListBudgetLimits handles the list_budget_limits tool
func (h *BudgetHandlers) HandleListBudgetLimits(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListBudgetLimitsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListBudgetLimitByBudgetParams{}

	// Set date range
	var startDate, endDate openapi_types.Date
	if params.Arguments.Start != "" {
		parsedStart, err := time.Parse("2006-01-02", params.Arguments.Start)
		if err != nil {
			return h.HandleError(err, "Invalid start date format"), nil
		}
		startDate = openapi_types.Date{Time: parsedStart}
		apiParams.Start = &startDate
	}

	if params.Arguments.End != "" {
		parsedEnd, err := time.Parse("2006-01-02", params.Arguments.End)
		if err != nil {
			return h.HandleError(err, "Invalid end date format"), nil
		}
		endDate = openapi_types.Date{Time: parsedEnd}
		apiParams.End = &endDate
	}

	// Use defaults if not provided
	if apiParams.Start == nil {
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = openapi_types.Date{Time: startOfMonth}
		apiParams.Start = &startDate
	}

	if apiParams.End == nil {
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		endOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 23, 59, 59, 999999999, now.Location())
		endDate = openapi_types.Date{Time: endOfMonth}
		apiParams.End = &endDate
	}

	resp, err := h.client.ListBudgetLimitByBudgetWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return h.HandleError(err, "Error listing budget limits"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	budgetLimitList := mappers.MapBudgetLimitArrayToBudgetLimitList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(budgetLimitList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleListBudgetTransactions handles the list_budget_transactions tool
func (h *BudgetHandlers) HandleListBudgetTransactions(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListBudgetTransactionsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListTransactionByBudgetParams{}

	// Set date range
	var startDate, endDate openapi_types.Date
	if params.Arguments.Start != "" {
		parsedStart, err := time.Parse("2006-01-02", params.Arguments.Start)
		if err != nil {
			return h.HandleError(err, "Invalid start date format"), nil
		}
		startDate = openapi_types.Date{Time: parsedStart}
		apiParams.Start = &startDate
	}

	if params.Arguments.End != "" {
		parsedEnd, err := time.Parse("2006-01-02", params.Arguments.End)
		if err != nil {
			return h.HandleError(err, "Invalid end date format"), nil
		}
		endDate = openapi_types.Date{Time: parsedEnd}
		apiParams.End = &endDate
	}

	// Use defaults if not provided
	if apiParams.Start == nil {
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = openapi_types.Date{Time: startOfMonth}
		apiParams.Start = &startDate
	}

	if apiParams.End == nil {
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		endOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 23, 59, 59, 999999999, now.Location())
		endDate = openapi_types.Date{Time: endOfMonth}
		apiParams.End = &endDate
	}

	if params.Arguments.Type != "" {
		filter := client.TransactionTypeFilter(params.Arguments.Type)
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

	resp, err := h.client.ListTransactionByBudgetWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return h.HandleError(err, fmt.Sprintf("Error listing transactions for budget %s", params.Arguments.ID)), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	transactionList := mappers.MapTransactionArrayToTransactionGroupList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(transactionList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// RegisterTools registers all budget-related tools
func (h *BudgetHandlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_budgets",
			Description: "List all budgets in Firefly III",
		}, h.HandleListBudgets,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_budget_limits",
			Description: "List limits for a specific budget",
		}, h.HandleListBudgetLimits,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_budget_transactions",
			Description: "List transactions for a specific budget",
		}, h.HandleListBudgetTransactions,
	)
}