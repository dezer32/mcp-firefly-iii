package fireflyMCP

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BulkTransactionStoreRequest represents the request for creating multiple transaction groups
type BulkTransactionStoreRequest struct {
	TransactionGroups []TransactionStoreRequest `json:"transaction_groups" mcp:"Array of transaction groups to create (required, at least one)"`
	DelayMs           int                       `json:"delay_ms,omitempty" mcp:"Delay in milliseconds between API calls to avoid rate limiting (default: 100)"`
}

// BulkTransactionStoreResponse represents the response for bulk transaction creation
type BulkTransactionStoreResponse struct {
	Results []TransactionGroupResult `json:"results"`
	Summary BulkSummary              `json:"summary"`
}

// TransactionGroupResult represents the result of creating a single transaction group
type TransactionGroupResult struct {
	Index            int               `json:"index"`
	Success          bool              `json:"success"`
	TransactionGroup *TransactionGroup `json:"transaction_group,omitempty"`
	Error            string            `json:"error,omitempty"`
}

// BulkSummary provides a summary of the bulk operation results
type BulkSummary struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

// handleStoreTransactionsBulk creates multiple transaction groups in Firefly III
func (s *FireflyMCPServer) handleStoreTransactionsBulk(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[BulkTransactionStoreRequest],
) (*mcp.CallToolResultFor[struct{}], error) {
	// Validate input
	if len(params.Arguments.TransactionGroups) == 0 {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: transaction_groups array is required and must not be empty"},
			},
			IsError: true,
		}, nil
	}

	// Limit batch size to prevent excessive API calls
	const maxBatchSize = 100
	if len(params.Arguments.TransactionGroups) > maxBatchSize {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: batch size exceeds maximum of %d transaction groups", maxBatchSize),
				},
			},
			IsError: true,
		}, nil
	}

	// Set default delay if not specified
	delayMs := params.Arguments.DelayMs
	if delayMs <= 0 {
		delayMs = 100 // Default 100ms delay between API calls
	}

	// Initialize response
	response := BulkTransactionStoreResponse{
		Results: make([]TransactionGroupResult, 0, len(params.Arguments.TransactionGroups)),
		Summary: BulkSummary{
			Total:      len(params.Arguments.TransactionGroups),
			Successful: 0,
			Failed:     0,
		},
	}

	// Process each transaction group sequentially
	for i, group := range params.Arguments.TransactionGroups {
		// Add delay between API calls (except for the first one)
		if i > 0 && delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}

		// Create individual result
		result := TransactionGroupResult{
			Index: i,
		}

		// Call the existing single transaction handler
		singleParams := &mcp.CallToolParamsFor[TransactionStoreRequest]{
			Arguments: group,
		}

		// Execute the single transaction creation
		singleResult, err := s.handleStoreTransaction(ctx, ss, singleParams)

		// Process the result
		if err != nil {
			// Actual error in handler (rare, as handler returns errors in result)
			result.Success = false
			result.Error = fmt.Sprintf("Internal error: %v", err)
			response.Summary.Failed++
		} else if singleResult.IsError {
			// Transaction creation failed
			result.Success = false
			// Extract error message from content
			if len(singleResult.Content) > 0 {
				if textContent, ok := singleResult.Content[0].(*mcp.TextContent); ok {
					result.Error = textContent.Text
				} else {
					result.Error = "Unknown error occurred"
				}
			}
			response.Summary.Failed++
		} else {
			// Transaction creation succeeded
			result.Success = true
			response.Summary.Successful++

			// Parse the transaction group from the response
			if len(singleResult.Content) > 0 {
				if textContent, ok := singleResult.Content[0].(*mcp.TextContent); ok {
					var txGroup TransactionGroup
					if err := json.Unmarshal([]byte(textContent.Text), &txGroup); err == nil {
						result.TransactionGroup = &txGroup
					}
				}
			}
		}

		response.Results = append(response.Results, result)

		// Check for context cancellation
		select {
		case <-ctx.Done():
			// Context cancelled, stop processing and return partial results
			response.Summary.Failed += (response.Summary.Total - len(response.Results))
			break
		default:
			// Continue processing
		}
	}

	// Marshal response to JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return &mcp.CallToolResultFor[struct{}]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error marshaling response: %v", err)},
			},
			IsError: true,
		}, nil
	}

	// Determine if overall operation should be marked as error
	// Only mark as error if ALL transactions failed
	isError := response.Summary.Successful == 0 && response.Summary.Failed > 0

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
		IsError: isError,
	}, nil
}
