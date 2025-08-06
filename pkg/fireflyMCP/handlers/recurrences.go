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

// RecurrenceHandlers handles all recurrence-related MCP tools
type RecurrenceHandlers struct {
	*BaseHandler
	client *client.ClientWithResponses
}

// NewRecurrenceHandlers creates a new RecurrenceHandlers instance
func NewRecurrenceHandlers(ctx HandlerContext) *RecurrenceHandlers {
	return &RecurrenceHandlers{
		BaseHandler: NewBaseHandler(ctx),
		client:      ctx.GetClient().(*client.ClientWithResponses),
	}
}

// ListRecurrencesArgs defines arguments for listing recurrences
type ListRecurrencesArgs struct {
	Limit int `json:"limit,omitempty" mcp:"Maximum number of recurrences to return"`
	Page  int `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// GetRecurrenceArgs defines arguments for getting a specific recurrence
type GetRecurrenceArgs struct {
	ID string `json:"id" mcp:"Recurrence ID"`
}

// ListRecurrenceTransactionsArgs defines arguments for listing recurrence transactions
type ListRecurrenceTransactionsArgs struct {
	ID    string `json:"id" mcp:"Recurrence ID"`
	Type  string `json:"type,omitempty" mcp:"Filter by transaction type"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// HandleListRecurrences handles the list_recurrences tool
func (h *RecurrenceHandlers) HandleListRecurrences(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListRecurrencesArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListRecurrenceParams{}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	page := int32(params.Arguments.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	resp, err := h.client.ListRecurrenceWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error listing recurrences"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	recurrenceList := mappers.MapRecurrenceArrayToRecurrenceList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(recurrenceList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleGetRecurrence handles the get_recurrence tool
func (h *RecurrenceHandlers) HandleGetRecurrence(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GetRecurrenceArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.GetRecurrenceParams{}

	resp, err := h.client.GetRecurrenceWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return h.HandleError(err, "Error getting recurrence"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	var recurrence interface{}
	if resp.ApplicationvndApiJSON200 != nil {
		recurrence = mappers.MapRecurrenceToRecurrence(&resp.ApplicationvndApiJSON200.Data)
	}

	result, _ := json.MarshalIndent(recurrence, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleListRecurrenceTransactions handles the list_recurrence_transactions tool
func (h *RecurrenceHandlers) HandleListRecurrenceTransactions(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListRecurrenceTransactionsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListTransactionByRecurrenceParams{}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	page := int32(params.Arguments.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	if params.Arguments.Type != "" {
		filter := client.TransactionTypeFilter(params.Arguments.Type)
		apiParams.Type = &filter
	}

	// Set date filters if provided
	if params.Arguments.Start != "" {
		startDate, err := time.Parse("2006-01-02", params.Arguments.Start)
		if err != nil {
			return h.HandleError(err, "Invalid start date format"), nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if params.Arguments.End != "" {
		endDate, err := time.Parse("2006-01-02", params.Arguments.End)
		if err != nil {
			return h.HandleError(err, "Invalid end date format"), nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	resp, err := h.client.ListTransactionByRecurrenceWithResponse(ctx, params.Arguments.ID, apiParams)
	if err != nil {
		return h.HandleError(err, fmt.Sprintf("Error listing recurrence transactions for %s", params.Arguments.ID)), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	transactionList := mappers.MapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(transactionList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// RegisterTools registers all recurrence-related tools
func (h *RecurrenceHandlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_recurrences",
			Description: "List all recurrences in Firefly III",
		}, h.HandleListRecurrences,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "get_recurrence",
			Description: "Get details of a specific recurrence",
		}, h.HandleGetRecurrence,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_recurrence_transactions",
			Description: "List transactions created by a specific recurrence",
		}, h.HandleListRecurrenceTransactions,
	)
}