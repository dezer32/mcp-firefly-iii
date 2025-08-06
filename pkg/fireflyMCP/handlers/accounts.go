package handlers

import (
	"context"
	"encoding/json"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/mappers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AccountHandlers handles all account-related MCP tools
type AccountHandlers struct {
	*BaseHandler
	client *client.ClientWithResponses
}

// NewAccountHandlers creates a new AccountHandlers instance
func NewAccountHandlers(ctx HandlerContext) *AccountHandlers {
	return &AccountHandlers{
		BaseHandler: NewBaseHandler(ctx),
		client:      ctx.GetClient().(*client.ClientWithResponses),
	}
}

// ListAccountsArgs defines arguments for listing accounts
type ListAccountsArgs struct {
	Type  string `json:"type,omitempty" mcp:"Filter by account type (asset, expense, revenue, etc.)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of accounts to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// GetAccountArgs defines arguments for getting a specific account
type GetAccountArgs struct {
	ID string `json:"id" mcp:"Account ID"`
}

// SearchAccountsArgs defines arguments for searching accounts
type SearchAccountsArgs struct {
	Query string `json:"query" mcp:"The search query"`
	Field string `json:"field" mcp:"The account field(s) to search in (all, iban, name, number, id)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of accounts to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// HandleListAccounts handles the list_accounts tool
func (h *AccountHandlers) HandleListAccounts(
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

	resp, err := h.client.ListAccountWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error listing accounts"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	accountList := mappers.MapAccountArrayToAccountList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(accountList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleGetAccount handles the get_account tool
func (h *AccountHandlers) HandleGetAccount(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GetAccountArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	resp, err := h.client.GetAccountWithResponse(ctx, params.Arguments.ID, &client.GetAccountParams{})
	if err != nil {
		return h.HandleError(err, "Error getting account"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	if resp.ApplicationvndApiJSON200 != nil {
		account := mappers.MapAccountReadToAccount(&resp.ApplicationvndApiJSON200.Data)
		result, _ := json.MarshalIndent(account, "", "  ")
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(result)},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Account not found"},
		},
		IsError: true,
	}, nil
}

// HandleSearchAccounts handles the search_accounts tool
func (h *AccountHandlers) HandleSearchAccounts(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[SearchAccountsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.SearchAccountsParams{
		Query: params.Arguments.Query,
	}
	
	if params.Arguments.Field != "" {
		field := client.AccountSearchFieldFilter(params.Arguments.Field)
		apiParams.Field = field
	}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	resp, err := h.client.SearchAccountsWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error searching accounts"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	accountList := mappers.MapAccountArrayToAccountList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(accountList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// RegisterTools registers all account-related tools
func (h *AccountHandlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_accounts",
			Description: "List all accounts in Firefly III",
		}, h.HandleListAccounts,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "get_account",
			Description: "Get details of a specific account",
		}, h.HandleGetAccount,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "search_accounts",
			Description: "Search for accounts by name, IBAN, or other fields",
		}, h.HandleSearchAccounts,
	)
}