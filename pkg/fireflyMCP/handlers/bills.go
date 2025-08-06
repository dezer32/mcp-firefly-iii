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

// BillHandlers handles all bill-related MCP tools
type BillHandlers struct {
	*BaseHandler
	client *client.ClientWithResponses
}

// NewBillHandlers creates a new BillHandlers instance
func NewBillHandlers(ctx HandlerContext) *BillHandlers {
	return &BillHandlers{
		BaseHandler: NewBaseHandler(ctx),
		client:      ctx.GetClient().(*client.ClientWithResponses),
	}
}

// ListBillsArgs defines arguments for listing bills
type ListBillsArgs struct {
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of bills to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// GetBillArgs defines arguments for getting a specific bill
type GetBillArgs struct {
	ID    string `json:"id" mcp:"Bill ID"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
}

// ListBillTransactionsArgs defines arguments for listing bill transactions
type ListBillTransactionsArgs struct {
	ID    string `json:"id" mcp:"Bill ID"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// HandleListBills handles the list_bills tool
func (h *BillHandlers) HandleListBills(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListBillsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListBillParams{}

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

	resp, err := h.client.ListBillWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error listing bills"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	billList := mappers.MapBillArrayToBillList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(billList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleGetBill handles the get_bill tool
func (h *BillHandlers) HandleGetBill(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GetBillArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.GetBillParams{}

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

	resp, err := h.client.GetBillWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return h.HandleError(err, fmt.Sprintf("Error getting bill %s", params.Arguments.ID)), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	if resp.ApplicationvndApiJSON200 != nil {
		bill := mappers.MapBillReadToBill(&resp.ApplicationvndApiJSON200.Data)
		result, _ := json.MarshalIndent(bill, "", "  ")
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(result)},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Bill not found"},
		},
		IsError: true,
	}, nil
}

// HandleListBillTransactions handles the list_bill_transactions tool
func (h *BillHandlers) HandleListBillTransactions(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListBillTransactionsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListTransactionByBillParams{}

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

	resp, err := h.client.ListTransactionByBillWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return h.HandleError(err, fmt.Sprintf("Error listing transactions for bill %s", params.Arguments.ID)), nil
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

// RegisterTools registers all bill-related tools
func (h *BillHandlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_bills",
			Description: "List all bills in Firefly III",
		}, h.HandleListBills,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "get_bill",
			Description: "Get details of a specific bill",
		}, h.HandleGetBill,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_bill_transactions",
			Description: "List transactions for a specific bill",
		}, h.HandleListBillTransactions,
	)
}