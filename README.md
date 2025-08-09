# Firefly III MCP Server

This is a Model Context Protocol (MCP) server implementation for Firefly III personal finance management system, built using the [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk).

## Features

The MCP server provides the following tools for interacting with Firefly III:

### Account Management
- `list_accounts` - List all accounts with optional filtering by type and limit
- `get_account` - Get detailed information about a specific account
- `search_accounts` - Search for accounts by name, IBAN, or other fields

### Transaction Management  
- `list_transactions` - List transactions with optional filtering by type, date range, and limit
- `get_transaction` - Get detailed information about a specific transaction
- `search_transactions` - Search for transactions by keyword
- `store_transaction` - Create a new transaction with support for splits, categorization, and rules
- `store_transactions_bulk` - Create multiple transaction groups in a single operation (up to 100 at once)

### Budget Management
- `list_budgets` - List all budgets with optional limit
- `list_budget_limits` - List budget limits for a specific budget with optional date range
- `list_budget_transactions` - List transactions for a specific budget with optional filters

### Category Management
- `list_categories` - List all categories with optional limit

### Tag Management
- `list_tags` - List all tags with optional pagination

### Financial Summary
- `get_summary` - Get basic financial summary with optional date range

### Expense Insights
- `expense_category_insights` - Get expense insights grouped by category for a date range
- `expense_total_insights` - Get total expense trends for a date range

## Configuration

The server uses a YAML configuration file (`config.yaml`) with the following structure:

```yaml
# Firefly III API Configuration
server:
  url: "https://your-firefly-instance.com/api"

api:
  token: "your-api-token-here"

client:
  timeout: 30 # timeout in seconds

# API call limits
limits:
  accounts: 10
  transactions: 5
  categories: 10
  budgets: 10

# MCP server configuration
mcp:
  name: "firefly-iii-mcp"
  version: "1.0.0"
  instructions: "MCP server for Firefly III personal finance management"
```

## Setup

1. **Configure Firefly III API Access**
   - Obtain an API token from your Firefly III instance
   - Update the `server.url` and `api.token` in `config.yaml`

2. **Build the Server**
   ```bash
   go build ./cmd/mcp-server
   ```

3. **Run the Server**
   ```bash
   ./mcp-server [config-file]
   ```
   
   If no config file is specified, it defaults to `config.yaml`.

## Usage

The server communicates over stdin/stdout using the MCP protocol. It can be integrated with MCP-compatible clients.

### Store Transaction Parameters

The `store_transaction` tool creates new transactions in Firefly III. It accepts the following parameters:

#### Request Structure
- `error_if_duplicate_hash` (boolean, optional) - Break if transaction already exists based on hash
- `apply_rules` (boolean, optional) - Whether to apply Firefly III rules when submitting
- `fire_webhooks` (boolean, optional) - Whether to fire webhooks (default: true)
- `group_title` (string, optional) - Title for split transactions
- `transactions` (array, required) - Array of transaction splits

#### Transaction Split Parameters

Each transaction in the `transactions` array requires:

**Required Fields:**
- `type` (string) - Transaction type: `withdrawal`, `deposit`, or `transfer`
- `date` (string) - Transaction date in `YYYY-MM-DD` format or ISO 8601 datetime
- `amount` (string) - Transaction amount as a positive decimal string
- `description` (string) - Transaction description

**Account Fields (at least one source/destination required):**
- `source_id` (string, optional) - Source account ID
- `source_name` (string, optional) - Source account name (creates new if doesn't exist)
- `destination_id` (string, optional) - Destination account ID  
- `destination_name` (string, optional) - Destination account name (creates new if doesn't exist)

**Categorization Fields (optional):**
- `category_id` (string) - Category ID
- `category_name` (string) - Category name (creates new if doesn't exist)
- `budget_id` (string) - Budget ID
- `budget_name` (string) - Budget name
- `tags` (array of strings) - Transaction tags

**Currency Fields (optional):**
- `currency_id` (string) - Currency ID
- `currency_code` (string) - Currency code (e.g., "USD", "EUR")
- `foreign_amount` (string) - Amount in foreign currency
- `foreign_currency_id` (string) - Foreign currency ID
- `foreign_currency_code` (string) - Foreign currency code

**Other Fields (optional):**
- `bill_id` (string) - Bill ID
- `bill_name` (string) - Bill name
- `piggy_bank_id` (string) - Piggy bank ID
- `piggy_bank_name` (string) - Piggy bank name
- `notes` (string) - Transaction notes
- `reconciled` (boolean) - Whether transaction is reconciled
- `order` (integer) - Order in the transaction split list

### Store Transactions Bulk Parameters

The `store_transactions_bulk` tool creates multiple transaction groups in Firefly III in a single operation. It's useful for batch importing transactions or creating multiple related transactions at once.

#### Request Structure
- `transaction_groups` (array, required) - Array of transaction groups to create (max 100)
- `delay_ms` (integer, optional) - Delay in milliseconds between API calls to avoid rate limiting (default: 100)

Each item in `transaction_groups` is a complete `store_transaction` request with the same parameters as described above.

#### Response Structure
The tool returns a detailed response showing the result of each transaction group creation:

```json
{
  "results": [
    {
      "index": 0,
      "success": true,
      "transaction_group": { /* created transaction details */ }
    },
    {
      "index": 1,
      "success": false,
      "error": "Validation error: Invalid category"
    }
  ],
  "summary": {
    "total": 2,
    "successful": 1,
    "failed": 1
  }
}
```

### Tool Examples

#### List Accounts
```json
{
  "name": "list_accounts",
  "arguments": {
    "type": "asset",
    "limit": 5
  }
}
```

#### Get Account Details
```json
{
  "name": "get_account", 
  "arguments": {
    "id": "123"
  }
}
```

#### Get Expense Category Insights
```json
{
  "name": "expense_category_insights",
  "arguments": {
    "start": "2024-01-01",
    "end": "2024-12-31",
    "accounts": ["1", "2"] // optional
  }
}
```

#### Get Expense Total Insights
```json
{
  "name": "expense_total_insights",
  "arguments": {
    "start": "2024-01-01",
    "end": "2024-12-31",
    "accounts": ["1", "2"] // optional
  }
}
```

#### Search Accounts
```json
{
  "name": "search_accounts",
  "arguments": {
    "query": "checking",
    "field": "name",
    "limit": 5
  }
}
```

Field options: `all`, `iban`, `name`, `number`, `id`

#### List Transactions
```json
{
  "name": "list_transactions",
  "arguments": {
    "type": "withdrawal",
    "start": "2024-01-01",
    "end": "2024-01-31",
    "limit": 10
  }
}
```

#### Search Transactions
```json
{
  "name": "search_transactions",
  "arguments": {
    "query": "groceries",
    "limit": 10,
    "page": 1,
    "start": "2024-01-01",
    "end": "2024-12-31"
  }
}
```

#### Store Transaction
```json
{
  "name": "store_transaction",
  "arguments": {
    "apply_rules": true,
    "fire_webhooks": true,
    "transactions": [
      {
        "type": "withdrawal",
        "date": "2024-01-15",
        "amount": "45.99",
        "description": "Grocery shopping",
        "source_id": "1",
        "destination_name": "Local Supermarket",
        "category_name": "Groceries",
        "tags": ["food", "weekly-shopping"]
      }
    ]
  }
}
```

#### Store Transactions Bulk
```json
{
  "name": "store_transactions_bulk",
  "arguments": {
    "transaction_groups": [
      {
        "group_title": "Groceries Week 1",
        "apply_rules": true,
        "transactions": [
          {
            "type": "withdrawal",
            "date": "2024-01-15",
            "amount": "45.99",
            "description": "Grocery shopping",
            "source_id": "1",
            "destination_name": "Local Supermarket",
            "category_name": "Groceries"
          }
        ]
      },
      {
        "group_title": "Utilities January",
        "transactions": [
          {
            "type": "withdrawal",
            "date": "2024-01-20",
            "amount": "120.00",
            "description": "Electric bill",
            "source_id": "1",
            "destination_name": "Power Company",
            "category_name": "Utilities"
          }
        ]
      }
    ],
    "delay_ms": 100
  }
}
```

#### List Tags
```json
{
  "name": "list_tags",
  "arguments": {
    "limit": 10,
    "page": 1
  }
}
```

#### Get Financial Summary
```json
{
  "name": "get_summary",
  "arguments": {
    "start": "2024-01-01",
    "end": "2024-01-31"
  }
}
```

#### List Budget Limits
```json
{
  "name": "list_budget_limits",
  "arguments": {
    "id": "1",
    "start": "2024-01-01",
    "end": "2024-12-31"
  }
}
```

#### List Budget Transactions
```json
{
  "name": "list_budget_transactions",
  "arguments": {
    "id": "1",
    "type": "withdrawal",
    "start": "2024-01-01",
    "end": "2024-12-31",
    "limit": 10,
    "page": 1
  }
}
```

## Architecture

The implementation consists of:

- **`cmd/mcp-server/main.go`** - Server entry point
- **`pkg/fireflyMCP/config.go`** - Configuration management
- **`pkg/fireflyMCP/server.go`** - MCP server implementation with tool handlers
- **`pkg/client/`** - Auto-generated Firefly III API client

## Authentication

The server uses Bearer token authentication with the Firefly III API. The token is automatically added to all API requests via a request editor function.

## Error Handling

All tools include proper error handling for:
- API connection errors
- Authentication failures
- Invalid parameters
- HTTP error responses

Errors are returned as MCP tool results with appropriate error messages.

## Development

To extend the server with additional tools:

1. Define argument types in `server.go`
2. Register the tool in `registerTools()`
3. Implement the handler function following the existing patterns
4. Update this documentation

## Dependencies

- [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk) - MCP protocol implementation
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - OpenAPI client generation
- Firefly III API client (auto-generated)