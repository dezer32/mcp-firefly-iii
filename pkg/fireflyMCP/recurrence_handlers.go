package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Recurrence argument types
type ListRecurrencesArgs struct {
	Limit int `json:"limit,omitempty" mcp:"Maximum number of recurrences to return"`
	Page  int `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

type GetRecurrenceArgs struct {
	ID string `json:"id" mcp:"Recurrence ID"`
}

type ListRecurrenceTransactionsArgs struct {
	ID    string `json:"id" mcp:"Recurrence ID"`
	Type  string `json:"type,omitempty" mcp:"Filter by transaction type"`
	Start string `json:"start,omitempty" mcp:"Start date (YYYY-MM-DD)"`
	End   string `json:"end,omitempty" mcp:"End date (YYYY-MM-DD)"`
	Limit int    `json:"limit,omitempty" mcp:"Maximum number of transactions to return"`
	Page  int    `json:"page,omitempty" mcp:"Page number for pagination (default: 1)"`
}

// handleListRecurrences lists all recurrences in Firefly III
func (s *FireflyMCPServer) handleListRecurrences(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListRecurrencesArgs,
) (*mcp.CallToolResult, struct{}, error) {
	// Prepare API parameters
	apiParams := &client.ListRecurrenceParams{}

	// Set pagination
	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}
	page := int32(args.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	// Call the API
	resp, err := s.client.ListRecurrenceWithResponse(ctx, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing recurrences: %v", err)},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %s", string(resp.Body))},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	// Map the response
	recurrenceList := mapRecurrenceArrayToRecurrenceList(resp.ApplicationvndApiJSON200)

	// Convert to JSON for response
	jsonData, err := json.Marshal(recurrenceList)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error marshaling response: %v", err)},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, struct{}{}, nil
}

// handleGetRecurrence gets a specific recurrence by ID
func (s *FireflyMCPServer) handleGetRecurrence(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args GetRecurrenceArgs,
) (*mcp.CallToolResult, struct{}, error) {
	// Prepare API parameters
	apiParams := &client.GetRecurrenceParams{}

	// Call the API
	resp, err := s.client.GetRecurrenceWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error getting recurrence: %v", err)},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %s", string(resp.Body))},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	// Map the response
	var recurrence *Recurrence
	if resp.ApplicationvndApiJSON200 != nil {
		recurrence = mapRecurrenceToRecurrence(&resp.ApplicationvndApiJSON200.Data)
	}

	// Convert to JSON for response
	jsonData, err := json.Marshal(recurrence)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error marshaling response: %v", err)},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, struct{}{}, nil
}

// handleListRecurrenceTransactions lists transactions created by a specific recurrence
func (s *FireflyMCPServer) handleListRecurrenceTransactions(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListRecurrenceTransactionsArgs,
) (*mcp.CallToolResult, struct{}, error) {
	// Prepare API parameters
	apiParams := &client.ListTransactionByRecurrenceParams{}

	// Set pagination
	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}
	page := int32(args.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	// Set transaction type filter if provided
	if args.Type != "" {
		filter := client.TransactionTypeFilter(args.Type)
		apiParams.Type = &filter
	}

	// Set date filters if provided
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid start date format: %v", err)},
				},
				IsError: true,
			}, struct{}{}, nil
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Invalid end date format: %v", err)},
				},
				IsError: true,
			}, struct{}{}, nil
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	// Call the API
	resp, err := s.client.ListTransactionByRecurrenceWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing recurrence transactions: %v", err)},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API error: %s", string(resp.Body))},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	// Map the response
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)

	// Convert to JSON for response
	jsonData, err := json.Marshal(transactionList)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error marshaling response: %v", err)},
			},
			IsError: true,
		}, struct{}{}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, struct{}{}, nil
}
