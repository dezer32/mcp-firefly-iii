# Integration Testing for Firefly III MCP Server

This document describes the integration testing approach for the Firefly III MCP server that makes real API calls to external systems.

## Overview

The integration tests verify that the MCP server can successfully communicate with a real Firefly III instance over HTTP. These tests make actual API calls to external systems, ensuring the complete integration works end-to-end.

## Test Structure

### Test Files
- `pkg/fireflyMCP/integration_test.go` - Main integration test suite

### Test Categories

1. **Direct Client Tests** - Test the underlying HTTP client directly
2. **MCP Tool Tests** - Test MCP tool handlers that wrap the client calls
3. **Error Handling Tests** - Test behavior with invalid configurations
4. **Comprehensive Tool Tests** - Test all available MCP tools

## Configuration

### Environment Variables (Recommended)
```bash
export FIREFLY_TEST_URL="https://your-firefly-instance.com/api"
export FIREFLY_TEST_TOKEN="your-api-token"
```

### Configuration File Fallback
If environment variables are not set, tests will attempt to load configuration from `config.yaml`:
```yaml
server:
  url: "https://your-firefly-instance.com/api"
api:
  token: "your-api-token"
```

## Running Tests

### Prerequisites
1. Ensure you have a working Firefly III instance
2. Generate an API token in Firefly III
3. Set environment variables or update config.yaml

### Commands

```bash
# Run all integration tests
go test -v ./pkg/fireflyMCP -run TestIntegration

# Run specific test
go test -v ./pkg/fireflyMCP -run TestIntegration_ListAccounts

# Run with timeout
go test -v -timeout 60s ./pkg/fireflyMCP -run TestIntegration
```

## Test Results Analysis

### Successful Test Output
```
[DEBUG_LOG] Testing direct client call to https://funds.appservice.space/api
[DEBUG_LOG] Response status: 200, Error: <nil>
Successfully retrieved accounts from Firefly III API
```

### What the Tests Verify

1. **Real API Communication**: Tests make actual HTTP requests to Firefly III
2. **Authentication**: Verifies API token authentication works
3. **MCP Protocol Integration**: Tests that MCP tools correctly call external APIs
4. **Error Handling**: Tests behavior with invalid URLs and tokens
5. **Multiple Endpoints**: Tests various Firefly III API endpoints

## Available Test Functions

### TestIntegration_ListAccounts
- Tests both direct client calls and MCP tool calls
- Verifies account listing functionality
- Checks HTTP response codes and error handling

### TestIntegration_ListTransactions
- Tests transaction listing with date ranges
- Verifies MCP tool parameter handling
- Tests external API integration for transactions

### TestIntegration_GetSummary
- Tests financial summary retrieval
- Verifies date range parameter handling
- Tests summary data processing

### TestIntegration_ListBudgets
- Tests budget listing functionality
- Verifies MCP tool parameter handling for budgets
- Tests external API integration for budget data

### TestIntegration_ErrorHandling
- Tests behavior with invalid URLs
- Verifies proper error propagation
- Tests timeout handling

### TestIntegration_AllTools
- Comprehensive test of all MCP tools
- Verifies each tool can make external calls
- Tests tool registration and parameter handling

## Debugging

### Debug Logging
Tests include debug logging with `[DEBUG_LOG]` prefix:
```go
fmt.Printf("[DEBUG_LOG] Testing direct client call to %s\n", testConfig.ServerURL)
```

### Common Issues

1. **Authentication Errors**: Check API token validity
2. **Network Timeouts**: Increase timeout in test configuration
3. **Invalid URLs**: Verify Firefly III instance is accessible
4. **Missing Dependencies**: Run `go mod tidy`

## Test Configuration Details

### TestConfig Structure
```go
type TestConfig struct {
    ServerURL string        // Firefly III API base URL
    APIToken  string        // Authentication token
    Timeout   time.Duration // Request timeout
}
```

### Default Limits
- Accounts: 10
- Transactions: 5
- Categories: 10
- Budgets: 10

## Integration with CI/CD

### Environment Setup
```bash
# In CI/CD pipeline
export FIREFLY_TEST_URL="https://test-firefly.example.com/api"
export FIREFLY_TEST_TOKEN="${FIREFLY_API_TOKEN}"
```

### Test Execution
```bash
# Skip integration tests if no config available
go test -v ./pkg/fireflyMCP -run TestIntegration || echo "Integration tests skipped"
```

## Security Considerations

1. **API Tokens**: Never commit real API tokens to version control
2. **Test Data**: Use test instances, not production data
3. **Rate Limiting**: Be mindful of API rate limits during testing
4. **Network Security**: Ensure test environments are properly secured

## Expected Behavior

### Successful Integration
- Tests make real HTTP requests to Firefly III
- API responses are properly parsed and handled
- MCP tools return structured data
- Error conditions are properly handled

### Test Outcomes
- **PASS**: External API calls succeed and return expected data
- **SKIP**: No configuration available (expected in some environments)
- **FAIL**: Network issues, authentication problems, or API changes

## Maintenance

### Updating Tests
1. Add new test functions for new MCP tools
2. Update API endpoint tests when Firefly III API changes
3. Adjust timeouts based on network conditions
4. Update authentication methods as needed

### Monitoring
- Monitor test execution times for performance regression
- Track API response formats for compatibility
- Watch for authentication token expiration
- Monitor external service availability

This integration testing approach ensures that the MCP server reliably communicates with external Firefly III instances in real-world scenarios.