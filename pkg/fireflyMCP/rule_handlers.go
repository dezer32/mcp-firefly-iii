package fireflyMCP

import (
	"context"
	"fmt"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Rule Group argument types

type ListRuleGroupsArgs struct {
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of rule groups to return"`
	Page  int `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type GetRuleGroupArgs struct {
	ID string `json:"id" jsonschema:"Rule group ID (required)"`
}

type CreateRuleGroupArgs struct {
	RuleGroupStoreRequest
}

type UpdateRuleGroupArgs struct {
	ID string `json:"id" jsonschema:"Rule group ID (required)"`
	RuleGroupUpdateRequest
}

type DeleteRuleGroupArgs struct {
	ID string `json:"id" jsonschema:"Rule group ID (required)"`
}

type ListRulesByGroupArgs struct {
	ID    string `json:"id" jsonschema:"Rule group ID (required)"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of rules to return"`
	Page  int    `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type TestRuleGroupArgs struct {
	ID       string  `json:"id" jsonschema:"Rule group ID to test (required)"`
	Start    string  `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End      string  `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Accounts []int64 `json:"accounts,omitempty" jsonschema:"Limit to these account IDs"`
}

type TriggerRuleGroupArgs struct {
	ID       string  `json:"id" jsonschema:"Rule group ID to trigger (required)"`
	Start    string  `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End      string  `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Accounts []int64 `json:"accounts,omitempty" jsonschema:"Limit to these account IDs"`
}

// Rule argument types

type ListRulesArgs struct {
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of rules to return"`
	Page  int `json:"page,omitempty" jsonschema:"Page number for pagination (default: 1)"`
}

type GetRuleArgs struct {
	ID string `json:"id" jsonschema:"Rule ID (required)"`
}

type CreateRuleArgs struct {
	RuleStoreRequest
}

type UpdateRuleArgs struct {
	ID string `json:"id" jsonschema:"Rule ID (required)"`
	RuleUpdateRequest
}

type DeleteRuleArgs struct {
	ID string `json:"id" jsonschema:"Rule ID (required)"`
}

type TestRuleArgs struct {
	ID       string  `json:"id" jsonschema:"Rule ID to test (required)"`
	Start    string  `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End      string  `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Accounts []int64 `json:"accounts,omitempty" jsonschema:"Limit to these account IDs"`
}

type TriggerRuleArgs struct {
	ID       string  `json:"id" jsonschema:"Rule ID to trigger (required)"`
	Start    string  `json:"start,omitempty" jsonschema:"Start date (YYYY-MM-DD)"`
	End      string  `json:"end,omitempty" jsonschema:"End date (YYYY-MM-DD)"`
	Accounts []int64 `json:"accounts,omitempty" jsonschema:"Limit to these account IDs"`
}

// Rule Group handlers

func (s *FireflyMCPServer) handleListRuleGroups(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListRuleGroupsArgs,
) (*mcp.CallToolResult, any, error) {
	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.ListRuleGroupParams{}

	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	page := int32(args.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	resp, err := apiClient.ListRuleGroupWithResponse(ctx, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error listing rule groups: %v", err))
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	ruleGroupList := mapRuleGroupArrayToRuleGroupList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(ruleGroupList)
}

func (s *FireflyMCPServer) handleGetRuleGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args GetRuleGroupArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule group ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.GetRuleGroupParams{}
	resp, err := apiClient.GetRuleGroupWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error getting rule group: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule group not found")
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	var ruleGroup *RuleGroup
	if resp.ApplicationvndApiJSON200 != nil {
		ruleGroup = mapRuleGroupReadToRuleGroup(&resp.ApplicationvndApiJSON200.Data)
	}
	return newSuccessResult(ruleGroup)
}

func (s *FireflyMCPServer) handleCreateRuleGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args CreateRuleGroupArgs,
) (*mcp.CallToolResult, any, error) {
	if args.Title == "" {
		return newErrorResult("Title is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.StoreRuleGroupParams{}
	body := mapRuleGroupStoreRequestToAPI(&args.RuleGroupStoreRequest)

	resp, err := apiClient.StoreRuleGroupWithResponse(ctx, apiParams, body)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error creating rule group: %v", err))
	}

	if resp.StatusCode() == 422 {
		return newErrorResult(fmt.Sprintf("Validation error: %s", string(resp.Body)))
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	var ruleGroup *RuleGroup
	if resp.ApplicationvndApiJSON200 != nil {
		ruleGroup = mapRuleGroupReadToRuleGroup(&resp.ApplicationvndApiJSON200.Data)
	}
	return newSuccessResult(ruleGroup)
}

func (s *FireflyMCPServer) handleUpdateRuleGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args UpdateRuleGroupArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule group ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.UpdateRuleGroupParams{}
	body := mapRuleGroupUpdateRequestToAPI(&args.RuleGroupUpdateRequest)

	resp, err := apiClient.UpdateRuleGroupWithResponse(ctx, args.ID, apiParams, body)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error updating rule group: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule group not found")
	}

	if resp.StatusCode() == 422 {
		return newErrorResult(fmt.Sprintf("Validation error: %s", string(resp.Body)))
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	var ruleGroup *RuleGroup
	if resp.ApplicationvndApiJSON200 != nil {
		ruleGroup = mapRuleGroupReadToRuleGroup(&resp.ApplicationvndApiJSON200.Data)
	}
	return newSuccessResult(ruleGroup)
}

func (s *FireflyMCPServer) handleDeleteRuleGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args DeleteRuleGroupArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule group ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.DeleteRuleGroupParams{}
	resp, err := apiClient.DeleteRuleGroupWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error deleting rule group: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule group not found")
	}

	if resp.StatusCode() != 204 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	return newSuccessResult(map[string]string{"status": "deleted", "id": args.ID})
}

func (s *FireflyMCPServer) handleListRulesByGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListRulesByGroupArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule group ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.ListRuleByGroupParams{}

	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	page := int32(args.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	resp, err := apiClient.ListRuleByGroupWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error listing rules by group: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule group not found")
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	ruleList := mapRuleArrayToRuleList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(ruleList)
}

func (s *FireflyMCPServer) handleTestRuleGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args TestRuleGroupArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule group ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.TestRuleGroupParams{}

	// Parse date filters
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Invalid start date format: %v", err))
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Invalid end date format: %v", err))
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	if len(args.Accounts) > 0 {
		apiParams.Accounts = &args.Accounts
	}

	resp, err := apiClient.TestRuleGroupWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error testing rule group: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule group not found")
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	// Return matched transactions
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(transactionList)
}

func (s *FireflyMCPServer) handleTriggerRuleGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args TriggerRuleGroupArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule group ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.FireRuleGroupParams{}

	// Parse date filters
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Invalid start date format: %v", err))
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Invalid end date format: %v", err))
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	if len(args.Accounts) > 0 {
		apiParams.Accounts = &args.Accounts
	}

	resp, err := apiClient.FireRuleGroupWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error triggering rule group: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule group not found")
	}

	if resp.StatusCode() != 204 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	return newSuccessResult(map[string]string{"status": "triggered", "id": args.ID, "message": "Rule group execution started asynchronously"})
}

// Rule handlers

func (s *FireflyMCPServer) handleListRules(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ListRulesArgs,
) (*mcp.CallToolResult, any, error) {
	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.ListRuleParams{}

	if args.Limit > 0 {
		limit := int32(args.Limit)
		apiParams.Limit = &limit
	}

	page := int32(args.Page)
	if page == 0 {
		page = 1
	}
	apiParams.Page = &page

	resp, err := apiClient.ListRuleWithResponse(ctx, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error listing rules: %v", err))
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	ruleList := mapRuleArrayToRuleList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(ruleList)
}

func (s *FireflyMCPServer) handleGetRule(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args GetRuleArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.GetRuleParams{}
	resp, err := apiClient.GetRuleWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error getting rule: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule not found")
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	var rule *Rule
	if resp.ApplicationvndApiJSON200 != nil {
		rule = mapRuleReadToRule(&resp.ApplicationvndApiJSON200.Data)
	}
	return newSuccessResult(rule)
}

func (s *FireflyMCPServer) handleCreateRule(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args CreateRuleArgs,
) (*mcp.CallToolResult, any, error) {
	// Validate required fields
	if args.Title == "" {
		return newErrorResult("Title is required")
	}
	if args.RuleGroupId == "" && (args.RuleGroupTitle == nil || *args.RuleGroupTitle == "") {
		return newErrorResult("Rule group ID or title is required")
	}
	if args.Trigger == "" {
		return newErrorResult("Trigger type is required (store-journal or update-journal)")
	}
	if len(args.Triggers) == 0 {
		return newErrorResult("At least one trigger condition is required")
	}
	if len(args.Actions) == 0 {
		return newErrorResult("At least one action is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.StoreRuleParams{}
	body := mapRuleStoreRequestToAPI(&args.RuleStoreRequest)

	resp, err := apiClient.StoreRuleWithResponse(ctx, apiParams, body)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error creating rule: %v", err))
	}

	if resp.StatusCode() == 422 {
		return newErrorResult(fmt.Sprintf("Validation error: %s", string(resp.Body)))
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	var rule *Rule
	if resp.ApplicationvndApiJSON200 != nil {
		rule = mapRuleReadToRule(&resp.ApplicationvndApiJSON200.Data)
	}
	return newSuccessResult(rule)
}

func (s *FireflyMCPServer) handleUpdateRule(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args UpdateRuleArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.UpdateRuleParams{}
	body := mapRuleUpdateRequestToAPI(&args.RuleUpdateRequest)

	resp, err := apiClient.UpdateRuleWithResponse(ctx, args.ID, apiParams, body)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error updating rule: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule not found")
	}

	if resp.StatusCode() == 422 {
		return newErrorResult(fmt.Sprintf("Validation error: %s", string(resp.Body)))
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	var rule *Rule
	if resp.ApplicationvndApiJSON200 != nil {
		rule = mapRuleReadToRule(&resp.ApplicationvndApiJSON200.Data)
	}
	return newSuccessResult(rule)
}

func (s *FireflyMCPServer) handleDeleteRule(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args DeleteRuleArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.DeleteRuleParams{}
	resp, err := apiClient.DeleteRuleWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error deleting rule: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule not found")
	}

	if resp.StatusCode() != 204 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	return newSuccessResult(map[string]string{"status": "deleted", "id": args.ID})
}

func (s *FireflyMCPServer) handleTestRule(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args TestRuleArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.TestRuleParams{}

	// Parse date filters
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Invalid start date format: %v", err))
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Invalid end date format: %v", err))
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	if len(args.Accounts) > 0 {
		apiParams.Accounts = &args.Accounts
	}

	resp, err := apiClient.TestRuleWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error testing rule: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule not found")
	}

	if resp.StatusCode() != 200 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	// Return matched transactions
	transactionList := mapTransactionArrayToTransactionList(resp.ApplicationvndApiJSON200)
	return newSuccessResult(transactionList)
}

func (s *FireflyMCPServer) handleTriggerRule(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args TriggerRuleArgs,
) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return newErrorResult("Rule ID is required")
	}

	apiClient, err := s.getClient(ctx, req)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to get API client: %v", err))
	}

	apiParams := &client.FireRuleParams{}

	// Parse date filters
	if args.Start != "" {
		startDate, err := time.Parse("2006-01-02", args.Start)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Invalid start date format: %v", err))
		}
		date := openapi_types.Date{Time: startDate}
		apiParams.Start = &date
	}

	if args.End != "" {
		endDate, err := time.Parse("2006-01-02", args.End)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Invalid end date format: %v", err))
		}
		date := openapi_types.Date{Time: endDate}
		apiParams.End = &date
	}

	if len(args.Accounts) > 0 {
		apiParams.Accounts = &args.Accounts
	}

	resp, err := apiClient.FireRuleWithResponse(ctx, args.ID, apiParams)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Error triggering rule: %v", err))
	}

	if resp.StatusCode() == 404 {
		return newErrorResult("Rule not found")
	}

	if resp.StatusCode() != 204 {
		return newErrorResult(fmt.Sprintf("API error: %d - %s", resp.StatusCode(), string(resp.Body)))
	}

	return newSuccessResult(map[string]string{"status": "triggered", "id": args.ID, "message": "Rule execution started asynchronously"})
}
