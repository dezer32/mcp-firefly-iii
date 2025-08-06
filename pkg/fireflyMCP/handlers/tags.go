package handlers

import (
	"context"
	"encoding/json"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/mappers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TagHandlers handles all tag-related MCP tools
type TagHandlers struct {
	*BaseHandler
	client *client.ClientWithResponses
}

// NewTagHandlers creates a new TagHandlers instance
func NewTagHandlers(ctx HandlerContext) *TagHandlers {
	return &TagHandlers{
		BaseHandler: NewBaseHandler(ctx),
		client:      ctx.GetClient().(*client.ClientWithResponses),
	}
}

// ListTagsArgs defines arguments for listing tags
type ListTagsArgs struct {
	Limit int `json:"limit,omitempty" mcp:"Maximum number of tags to return"`
	Page  int `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// HandleListTags handles the list_tags tool
func (h *TagHandlers) HandleListTags(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListTagsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	apiParams := &client.ListTagParams{}

	if params.Arguments.Limit > 0 {
		limit := int32(params.Arguments.Limit)
		apiParams.Limit = &limit
	}

	if params.Arguments.Page > 0 {
		page := int32(params.Arguments.Page)
		apiParams.Page = &page
	}

	resp, err := h.client.ListTagWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error listing tags"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	tagList := mappers.MapTagArrayToTagList(resp.ApplicationvndApiJSON200)
	result, _ := json.MarshalIndent(tagList, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// RegisterTools registers all tag-related tools
func (h *TagHandlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "list_tags",
			Description: "List all tags in Firefly III",
		}, h.HandleListTags,
	)
}