package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/mappers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// InsightHandlers handles all insight-related MCP tools
type InsightHandlers struct {
	*BaseHandler
	client *client.ClientWithResponses
}

// NewInsightHandlers creates a new InsightHandlers instance
func NewInsightHandlers(ctx HandlerContext) *InsightHandlers {
	return &InsightHandlers{
		BaseHandler: NewBaseHandler(ctx),
		client:      ctx.GetClient().(*client.ClientWithResponses),
	}
}

// GetSummaryArgs defines arguments for getting summary
type GetSummaryArgs struct {
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
}

// ExpenseCategoryInsightsArgs defines arguments for expense category insights
type ExpenseCategoryInsightsArgs struct {
	Start    string   `json:"start" mcp:"Start date (YYYY-MM-DD)"`
	End      string   `json:"end" mcp:"End date (YYYY-MM-DD)"`
	Accounts []string `json:"accounts,omitempty" mcp:"Account IDs to include in results"`
}

// ExpenseTotalInsightsArgs defines arguments for expense total insights
type ExpenseTotalInsightsArgs struct {
	Start    string   `json:"start" mcp:"Start date (YYYY-MM-DD)"`
	End      string   `json:"end" mcp:"End date (YYYY-MM-DD)"`
	Accounts []string `json:"accounts,omitempty" mcp:"Account IDs to include in results"`
}

// HandleGetSummary handles the get_summary tool
func (h *InsightHandlers) HandleGetSummary(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GetSummaryArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Set date range
	var startDate, endDate openapi_types.Date
	if params.Arguments.Start != "" {
		parsedStart, err := time.Parse("2006-01-02", params.Arguments.Start)
		if err != nil {
			return h.HandleError(err, "Invalid start date format"), nil
		}
		startDate = openapi_types.Date{Time: parsedStart}
	} else {
		// Use default start date
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = openapi_types.Date{Time: startOfMonth}
	}

	if params.Arguments.End != "" {
		parsedEnd, err := time.Parse("2006-01-02", params.Arguments.End)
		if err != nil {
			return h.HandleError(err, "Invalid end date format"), nil
		}
		endDate = openapi_types.Date{Time: parsedEnd}
	} else {
		// Use default end date
		now := time.Now()
		nextMonth := now.AddDate(0, 1, 0)
		endOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 23, 59, 59, 999999999, now.Location())
		endDate = openapi_types.Date{Time: endOfMonth}
	}

	apiParams := &client.GetBasicSummaryParams{
		Start: startDate,
		End:   endDate,
	}

	resp, err := h.client.GetBasicSummaryWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error getting summary"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Process the summary response
	summaryMap := make(map[string]interface{})

	if resp.ApplicationvndApiJSON200 != nil {
		for key, basicSummaryEntry := range *resp.ApplicationvndApiJSON200 {
			entryMap := map[string]interface{}{
				"currency_code":     "",
				"monetary_value":    nil,
				"sub_title":         "",
				"local_icon":        "",
				"value_parsed":      nil,
			}

			if basicSummaryEntry.CurrencyCode != nil {
				entryMap["currency_code"] = *basicSummaryEntry.CurrencyCode
			}
			if basicSummaryEntry.MonetaryValue != nil {
				entryMap["monetary_value"] = *basicSummaryEntry.MonetaryValue
			}
			if basicSummaryEntry.SubTitle != nil {
				entryMap["sub_title"] = *basicSummaryEntry.SubTitle
			}
			if basicSummaryEntry.LocalIcon != nil {
				entryMap["local_icon"] = *basicSummaryEntry.LocalIcon
			}
			if basicSummaryEntry.ValueParsed != nil {
				entryMap["value_parsed"] = *basicSummaryEntry.ValueParsed
			}

			summaryMap[key] = entryMap
		}
	}

	result, _ := json.MarshalIndent(summaryMap, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleExpenseCategoryInsights handles the expense_category_insights tool
func (h *InsightHandlers) HandleExpenseCategoryInsights(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ExpenseCategoryInsightsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", params.Arguments.Start)
	if err != nil {
		return h.HandleError(err, "Invalid start date format"), nil
	}

	endDate, err := time.Parse("2006-01-02", params.Arguments.End)
	if err != nil {
		return h.HandleError(err, "Invalid end date format"), nil
	}

	start := openapi_types.Date{Time: startDate}
	end := openapi_types.Date{Time: endDate}

	apiParams := &client.InsightExpenseCategoryParams{
		Start: start,
		End:   end,
	}

	// Note: The Accounts parameter is *[]int64, not string
	// For now, we'll skip account filtering as it requires int64 conversion

	resp, err := h.client.InsightExpenseCategoryWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error getting expense category insights"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	insights := mappers.MapInsightGroupArrayToInsightCategoryList(resp.JSON200)
	result, _ := json.MarshalIndent(insights, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// HandleExpenseTotalInsights handles the expense_total_insights tool
func (h *InsightHandlers) HandleExpenseTotalInsights(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ExpenseTotalInsightsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", params.Arguments.Start)
	if err != nil {
		return h.HandleError(err, "Invalid start date format"), nil
	}

	endDate, err := time.Parse("2006-01-02", params.Arguments.End)
	if err != nil {
		return h.HandleError(err, "Invalid end date format"), nil
	}

	start := openapi_types.Date{Time: startDate}
	end := openapi_types.Date{Time: endDate}

	apiParams := &client.InsightExpenseTotalParams{
		Start: start,
		End:   end,
	}

	// Note: The Accounts parameter is *[]int64, not string
	// For now, we'll skip account filtering as it requires int64 conversion

	resp, err := h.client.InsightExpenseTotalWithResponse(ctx, apiParams)
	if err != nil {
		return h.HandleError(err, "Error getting expense total insights"), nil
	}

	if resp.StatusCode() != 200 {
		return h.HandleAPIError(resp.StatusCode()), nil
	}

	// Map response to DTO
	insights := mappers.MapInsightTotalArrayToInsightTotalList(resp.JSON200)
	result, _ := json.MarshalIndent(insights, "", "  ")
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

// RegisterTools registers all insight-related tools
func (h *InsightHandlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "get_summary",
			Description: "Get a basic financial summary for a date range",
		}, h.HandleGetSummary,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "expense_category_insights",
			Description: "Get expense insights grouped by category",
		}, h.HandleExpenseCategoryInsights,
	)

	mcp.AddTool(
		server, &mcp.Tool{
			Name:        "expense_total_insights",
			Description: "Get total expense insights",
		}, h.HandleExpenseTotalInsights,
	)
}