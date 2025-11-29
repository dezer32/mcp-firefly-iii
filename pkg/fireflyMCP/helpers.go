package fireflyMCP

import (
	"encoding/json"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// newErrorResult creates a standardized MCP error response.
func newErrorResult(msg string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
		IsError: true,
	}, nil, nil
}

// newSuccessResult creates a standardized MCP success response with JSON-formatted data.
// Returns an error result if JSON marshaling fails.
func newSuccessResult(data interface{}) (*mcp.CallToolResult, any, error) {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return newErrorResult("Failed to marshal response: " + err.Error())
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil, nil
}

// parseOptionalDate parses a date string in YYYY-MM-DD format.
// Returns nil if the input string is empty.
// Returns an error if the date format is invalid.
func parseOptionalDate(dateStr string) (*openapi_types.Date, error) {
	if dateStr == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &openapi_types.Date{Time: parsed}, nil
}
