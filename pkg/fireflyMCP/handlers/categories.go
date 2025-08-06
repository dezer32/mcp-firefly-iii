package handlers

import (
	"context"
	"encoding/json"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/mappers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CategoryHandlers handles all category-related MCP tools
type CategoryHandlers struct {
	*BaseHandler
	client *client.ClientWithResponses
}

// NewCategoryHandlers creates a new CategoryHandlers instance
func NewCategoryHandlers(ctx HandlerContext) *CategoryHandlers {
	return &CategoryHandlers{
		BaseHandler: NewBaseHandler(ctx),
		client:      ctx.GetClient().(*client.ClientWithResponses),
	}
}

// ListCategoriesArgs defines arguments for listing categories
type ListCategoriesArgs struct {
	Limit int `json:"limit,omitempty" mcp:"Maximum number of categories to return"`
	Page  int `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// HandleListCategories handles the list_categories tool
func (h *CategoryHandlers) HandleListCategories(
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

	resp, err := h.client.ListCategoryWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error listing categories"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	categoryList := mappers.MapCategoryArrayToCategoryList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(categoryList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// RegisterTools registers all category-related tools
func (h *CategoryHandlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_categories",
			Description: "List all categories in Firefly III",
		}, h.HandleListCategories,
	)
}