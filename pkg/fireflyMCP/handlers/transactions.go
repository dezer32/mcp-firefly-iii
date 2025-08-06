package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// TransactionHandlers handles all transaction-related MCP tools
type TransactionHandlers struct {
	*BaseHandler
	client *client.ClientWithResponses
}

// NewTransactionHandlers creates a new TransactionHandlers instance
func NewTransactionHandlers(ctx HandlerContext) *TransactionHandlers {
	return &TransactionHandlers{
		BaseHandler: NewBaseHandler(ctx),
		client:      ctx.GetClient().(*client.ClientWithResponses),
	}
}

// ListTransactionsArgs defines arguments for listing transactions
type ListTransactionsArgs struct {
	Type  string `json:"type,omitempty" mcp:"Filter by transaction type"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// GetTransactionArgs defines arguments for getting a specific transaction
type GetTransactionArgs struct {
	ID string `json:"id" mcp:"Transaction ID"`
}

// SearchTransactionsArgs defines arguments for searching transactions
type SearchTransactionsArgs struct {
	Query string `json:"query" mcp:"The search query"`
	Limit int32  `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int32  `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
}

// RegisterTools registers all transaction-related tools
func (h *TransactionHandlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_transactions",
			Description: "List transactions in Firefly III",
		}, h.HandleListTransactions,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "get_transaction",
			Description: "Get details of a specific transaction",
		}, h.HandleGetTransaction,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "search_transactions",
			Description: "Search for transactions by keyword",
		}, h.HandleSearchTransactions,
	)
}

// HandleListTransactions handles the list_transactions tool
func (h *TransactionHandlers) HandleListTransactions(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListTransactionsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListTransactionParams{}

	if params.Arguments.Type != "" {
		transType := client.TransactionTypeFilter(params.Arguments.Type)
		apiParams.Type = &transType
	}

	if params.Arguments.Start != "" {
		start, err := time.Parse("2006-01-02", params.Arguments.Start)
		if err == nil {
			startDate := openapi_types.Date{Time: start}
			apiParams.Start = &startDate
		}
	}

	if params.Arguments.End != "" {
		end, err := time.Parse("2006-01-02", params.Arguments.End)
		if err == nil {
			endDate := openapi_types.Date{Time: end}
			apiParams.End = &endDate
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

	resp, err := h.client.ListTransactionWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error listing transactions"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// TODO: Use mapper from mappers package when available
	result, _ := json.Marshal(resp.ApplicationvndApiJSON200)
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleGetTransaction handles the get_transaction tool
func (h *TransactionHandlers) HandleGetTransaction(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GetTransactionArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	resp, err := h.client.GetTransactionWithResponse(ctx, params.Arguments.ID, &client.GetTransactionParams{})
	if err != nil {
		return h.HandleError(err, "Error getting transaction"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// TODO: Use mapper from mappers package when available
	result, _ := json.Marshal(resp.ApplicationvndApiJSON200)
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleSearchTransactions handles the search_transactions tool
func (h *TransactionHandlers) HandleSearchTransactions(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[SearchTransactionsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.SearchTransactionsParams{
		Query: params.Arguments.Query,
	}

	if params.Arguments.Limit > 0 {
		apiParams.Limit = &params.Arguments.Limit
	}

	if params.Arguments.Page > 0 {
		apiParams.Page = &params.Arguments.Page
	}

	resp, err := h.client.SearchTransactionsWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error searching transactions"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// TODO: Use mapper from mappers package when available
	result, _ := json.Marshal(resp.ApplicationvndApiJSON200)
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}