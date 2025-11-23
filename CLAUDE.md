# Firefly III MCP Server - Claude Development Guide

**Note**: This project uses [bd (beads)](https://github.com/steveyegge/beads) for issue tracking. Use `bd` commands instead of markdown TODOs. See AGENTS.md for workflow details.

## Project Overview

This is a Model Context Protocol (MCP) server implementation for Firefly III personal finance management system. The server provides a bridge between AI assistants (like Claude) and Firefly III instances, enabling programmatic access to financial data through a standardized protocol.

## Architecture

### Core Components

#### API Layer (`pkg/client/`)
- Auto-generated OpenAPI client for Firefly III REST API
- Generated from `resources/firefly-iii-6.2.21-v1.yaml` OpenAPI specification
- Provides type-safe Go bindings for all Firefly III endpoints
- Uses `oapi-codegen` for code generation

#### Data Transfer Objects (`pkg/fireflyMCP/dto.go`)
- Simplified DTOs for cleaner API responses
- Focuses on essential data fields
- Includes pagination support
- Types: `Budget`, `Category`, `Account`, `Transaction`, `TransactionGroup`
- **IMPORTANT**: When returning data (Account, TransactionGroup, etc.), always use the existing DTOs and map values from the API responses to these DTOs. Never return raw API responses directly.

#### Service Layer (`pkg/fireflyMCP/server.go`)
- MCP server implementation using `modelcontextprotocol/go-sdk`
- Tool handlers for each MCP function
- Data mapping between Firefly III API and simplified DTOs
- Error handling and response formatting

#### Configuration (`pkg/fireflyMCP/config.go`)
- YAML-based configuration management
- Environment variable support for sensitive data
- Default values for optional settings
- Timeout and rate limit configuration

### Data Model

The server exposes simplified data models:

```go
// Core entities
Account {
  Id, Active, Name, Notes, Type
}

Transaction {
  Id, Amount, Date, Description, 
  SourceName, DestinationName, Type,
  CategoryName, BudgetName, Tags
}

Budget {
  Id, Active, Name, Notes, Spent
}

Category {
  Id, Name, Notes
}

// Insight entities
InsightCategoryEntry {
  Id, Name, Amount, CurrencyCode
}

InsightTotalEntry {
  Amount, CurrencyCode
}

// All responses include pagination
Pagination {
  Count, Total, CurrentPage, PerPage, TotalPages
}
```

### API Endpoints

MCP tools map to Firefly III endpoints:

- `list_accounts` → GET /api/v1/accounts
- `get_account` → GET /api/v1/accounts/{id}
- `search_accounts` → GET /api/v1/search/accounts
- `list_transactions` → GET /api/v1/transactions
- `get_transaction` → GET /api/v1/transactions/{id}
- `search_transactions` → GET /api/v1/search/transactions
- `list_budgets` → GET /api/v1/budgets
- `list_categories` → GET /api/v1/categories
- `get_summary` → GET /api/v1/summary/basic
- `expense_category_insights` → GET /api/v1/insight/expense/category
- `expense_total_insights` → GET /api/v1/insight/expense/total

### Technology Stack

- **Language**: Go 1.24.5
- **MCP SDK**: `modelcontextprotocol/go-sdk` v0.2.0
- **HTTP Client**: Generated via `oapi-codegen`
- **Testing**: `stretchr/testify`
- **Configuration**: `gopkg.in/yaml.v3`

### Project Structure

```
firefly-iii/
├── cmd/mcp-server/         # Main entry point
│   └── main.go
├── pkg/
│   ├── client/            # Auto-generated API client
│   │   ├── client.go
│   │   └── generate.go
│   └── fireflyMCP/        # MCP server implementation
│       ├── config.go      # Configuration management
│       ├── dto.go         # Data transfer objects
│       ├── server.go      # MCP server & handlers
│       ├── integration_test.go  # Integration tests
│       └── mapper_test.go       # Unit tests for mappers
├── resources/             # OpenAPI specifications
├── example/              # Example usage
├── config.yaml          # Configuration file
├── README.md           # User documentation
├── TESTING.md          # Testing documentation
└── CLAUDE.md          # This file
```

## Testing Strategy

### Unit Tests (`mapper_test.go`)
- Tests for data mapping functions
- Validates DTO transformations
- Edge case handling (nil values, empty data)
- No external dependencies

### Integration Tests (`integration_test.go`)
- Real API calls to Firefly III instances
- End-to-end MCP tool testing
- Error handling verification
- Configuration via environment variables or config file

### Test Execution
```bash
# Unit tests only
go test ./pkg/fireflyMCP -run "^Test[^Integration]"

# Integration tests (requires live Firefly III)
go test ./pkg/fireflyMCP -run TestIntegration

# All tests
go test ./pkg/fireflyMCP -v
```

## Development Commands

### Building
```bash
# Build the MCP server
go build -o mcp-server ./cmd/mcp-server

# Build with version info
go build -ldflags "-X main.Version=1.0.0" ./cmd/mcp-server
```

### Running
```bash
# Run with default config.yaml
./mcp-server

# Run with custom config
./mcp-server /path/to/custom-config.yaml

# Run directly without building
go run ./cmd/mcp-server
```

### Code Generation
```bash
# Regenerate API client from OpenAPI spec
cd pkg/client
go generate
```

### Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./pkg/fireflyMCP

# Run specific test
go test -v -run TestMapBudgetArrayToBudgetList ./pkg/fireflyMCP
```

### Development Tools
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Update dependencies
go mod tidy

# Vendor dependencies
go mod vendor
```

## Environment Setup

### Prerequisites
1. Go 1.24.5 or later
2. Access to a Firefly III instance
3. Valid API token from Firefly III

### Configuration Setup
1. Copy `config.yaml.example` to `config.yaml`
2. Update server URL and API token
3. Adjust limits and timeouts as needed

### Environment Variables
For sensitive data, use environment variables:
```bash
export FIREFLY_TEST_URL="https://your-instance.com/api"
export FIREFLY_TEST_TOKEN="your-api-token"
```

## Development Guidelines

### Code Style
- Follow standard Go conventions
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and testable

### Error Handling
- Return MCP error responses, not Go errors
- Include helpful error messages
- Log errors for debugging
- Handle nil pointers gracefully

### Testing Requirements
- Write unit tests for new functions
- Add integration tests for new MCP tools
- Maintain >80% code coverage
- Test error conditions

### API Client Updates
When Firefly III API changes:
1. Update OpenAPI spec in `resources/`
2. Regenerate client: `cd pkg/client && go generate`
3. Update DTOs if needed
4. Update mappers and tests
5. Test against live instance

### Adding New MCP Tools
1. Define argument struct in `server.go`
2. Create handler function
3. Register tool in `registerTools()`
4. Add mapper if needed
5. Write unit and integration tests
6. Update README documentation

## Security Considerations

### API Token Management
- Never commit tokens to version control
- Use environment variables in production
- Rotate tokens regularly
- Use read-only tokens when possible

### Input Validation
- Validate all user inputs
- Sanitize date formats
- Check numeric ranges
- Handle missing required fields

### Network Security
- Use HTTPS for all API calls
- Verify SSL certificates
- Handle timeouts appropriately
- Implement retry logic with backoff

## Future Considerations

### Potential Enhancements
1. **Write Operations**: Add support for creating/updating transactions
2. **Bulk Operations**: Support batch processing for better performance
3. **Webhooks**: Implement real-time data synchronization
4. **Caching**: Add intelligent caching for frequently accessed data
5. **Advanced Filtering**: Support complex query parameters
6. **Multi-Currency**: Enhanced multi-currency support
7. **Attachments**: Support for transaction attachments
8. **Rules Engine**: Integration with Firefly III rules

### Performance Optimizations
1. Connection pooling for HTTP client
2. Concurrent API calls where appropriate
3. Response compression
4. Pagination optimization
5. Query result caching

### Monitoring & Observability
1. Prometheus metrics integration
2. OpenTelemetry tracing
3. Structured logging
4. Health check endpoints
5. Performance benchmarks

### MCP Protocol Extensions
1. Streaming responses for large datasets
2. Progress indicators for long operations
3. Cancellation support
4. Resource management
5. Event notifications

## Known Issues & Limitations

### Current Limitations
1. Read-only operations (no write support)
2. Limited to 9 core MCP tools
3. No support for attachments
4. Basic pagination only
5. No real-time updates

### API Compatibility
- Tested with Firefly III v6.2.21
- May have issues with older versions
- JSON unmarshaling errors with some field types
- Currency ID field type inconsistencies

### Performance Considerations
- Each MCP call makes synchronous API request
- No connection pooling
- Limited by Firefly III API rate limits
- No caching mechanism

## Debugging Tips

### Common Issues
1. **Authentication Failures**: Check API token validity
2. **Network Errors**: Verify Firefly III URL accessibility
3. **JSON Errors**: Check for API response format changes
4. **Timeout Issues**: Increase timeout in config
5. **Empty Responses**: Verify data exists in Firefly III

### Debug Logging
Add debug statements:
```go
fmt.Printf("[DEBUG] API URL: %s\n", config.Server.URL)
fmt.Printf("[DEBUG] Response: %+v\n", resp)
```

### Testing Individual Tools
```bash
# Test via stdin (manual)
echo '{"method":"tools/call","params":{"name":"list_accounts","arguments":{"limit":5}}}' | ./mcp-server

# Use MCP client for testing
mcp-client call firefly-iii-mcp list_accounts '{"limit": 5}'
```

## Contributing Guidelines

### Code Contributions
1. Fork the repository
2. Create feature branch
3. Write tests first (TDD)
4. Implement feature
5. Run all tests
6. Submit pull request

### Documentation Updates
1. Update relevant .md files
2. Add code comments
3. Include examples
4. Update CHANGELOG

### Review Checklist
- [ ] Tests pass
- [ ] Code formatted
- [ ] Documentation updated
- [ ] No security issues
- [ ] Backwards compatible