# Firefly III Go Client

This project provides a Go client library for the [Firefly III](https://www.firefly-iii.org/) personal finance manager API, generated using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

## Features

- **Complete API Coverage**: Generated from the official Firefly III OpenAPI specification (v6.2.21)
- **Type Safety**: Fully typed Go structs for all API models and responses
- **HTTP Client Integration**: Built on standard Go HTTP client with customizable transport
- **Authentication Support**: Bearer token authentication built-in
- **Response Handling**: Both raw HTTP responses and parsed response objects
- **Comprehensive**: Supports all Firefly III API endpoints including:
  - Accounts management
  - Transactions and transfers
  - Budgets and categories
  - Bills and recurring transactions
  - Rules and rule groups
  - Reports and insights
  - User and system management

## Installation

```bash
go get github.com/dezer32/firefly-iii
```

## Quick Start

### Basic Setup

```go
package main

import (
    "context"
    "log"
    "net/http"
    
    "github.com/dezer32/firefly-iii/pkg/client"
)

// AuthTransport adds Bearer token authentication
type AuthTransport struct {
    Token string
    Base  http.RoundTripper
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    req.Header.Set("Authorization", "Bearer "+t.Token)
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    
    base := t.Base
    if base == nil {
        base = http.DefaultTransport
    }
    return base.RoundTrip(req)
}

func main() {
    // Create authenticated HTTP client
    httpClient := &http.Client{
        Transport: &AuthTransport{
            Token: "your-api-token-here",
        },
    }

    // Create Firefly III client
    fireflyClient, err := client.NewClientWithResponses(
        "https://your-firefly-iii-instance.com",
        client.WithHTTPClient(httpClient),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    ctx := context.Background()

    // Get system information
    aboutResp, err := fireflyClient.GetAboutWithResponse(ctx, &client.GetAboutParams{})
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    if aboutResp.StatusCode() == 200 {
        log.Println("Successfully connected to Firefly III!")
    }
}
```

### API Token Setup

1. Log into your Firefly III instance
2. Go to Options → Profile → OAuth
3. Create a new Personal Access Token
4. Copy the token and use it in your Go application

## Usage Examples

### List Accounts

```go
accountsResp, err := fireflyClient.ListAccountWithResponse(ctx, &client.ListAccountParams{
    Limit: &[]int32{10}[0], // Get first 10 accounts
    Type:  &[]client.AccountTypeFilter{client.AccountTypeFilter("asset")}[0],
})

if err != nil {
    log.Printf("Error: %v", err)
    return
}

if accountsResp.StatusCode() == 200 && accountsResp.JSON200 != nil {
    for _, account := range accountsResp.JSON200.Data {
        fmt.Printf("Account: %s (ID: %s)\n", 
            *account.Attributes.Name, 
            account.Id)
    }
}
```

### Create a Transaction

```go
transactionData := client.TransactionStore{
    Transactions: []client.TransactionSplitStore{
        {
            Type:               client.TransactionTypeProperty("withdrawal"),
            Description:        "Coffee purchase",
            Amount:             "4.50",
            SourceId:           &sourceAccountId,
            DestinationName:    &[]string{"Coffee Shop"}[0],
            CategoryName:       &[]string{"Food & Drinks"}[0],
        },
    },
}

transactionResp, err := fireflyClient.StoreTransactionWithResponse(
    ctx, 
    &client.StoreTransactionParams{}, 
    transactionData,
)

if err != nil {
    log.Printf("Error creating transaction: %v", err)
    return
}

if transactionResp.StatusCode() == 200 {
    fmt.Println("Transaction created successfully!")
}
```

### Get Budget Information

```go
budgetsResp, err := fireflyClient.ListBudgetWithResponse(ctx, &client.ListBudgetParams{
    Limit: &[]int32{20}[0],
})

if err != nil {
    log.Printf("Error: %v", err)
    return
}

if budgetsResp.StatusCode() == 200 && budgetsResp.JSON200 != nil {
    for _, budget := range budgetsResp.JSON200.Data {
        fmt.Printf("Budget: %s\n", *budget.Attributes.Name)
        if budget.Attributes.Spent != nil {
            fmt.Printf("  Spent: %s\n", *budget.Attributes.Spent)
        }
    }
}
```

## Client Types

The library provides two client types:

### Basic Client
```go
client, err := client.NewClient("https://your-instance.com", options...)
```
Returns raw `*http.Response` objects that you need to parse manually.

### Client with Responses (Recommended)
```go
client, err := client.NewClientWithResponses("https://your-instance.com", options...)
```
Returns typed response objects with parsed JSON data accessible via `.JSON200`, `.JSON400`, etc.

## Error Handling

The client provides detailed error information:

```go
resp, err := fireflyClient.ListAccountWithResponse(ctx, params)
if err != nil {
    log.Printf("Request failed: %v", err)
    return
}

switch resp.StatusCode() {
case 200:
    // Success - use resp.JSON200
    accounts := resp.JSON200.Data
case 401:
    // Unauthorized - check your API token
    log.Println("Authentication failed")
case 422:
    // Validation error - check resp.JSON422 for details
    if resp.JSON422 != nil {
        log.Printf("Validation errors: %+v", resp.JSON422.Errors)
    }
default:
    log.Printf("Unexpected status: %d", resp.StatusCode())
}
```

## Configuration Options

### Custom HTTP Client
```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &AuthTransport{Token: "your-token"},
}

client, err := client.NewClientWithResponses(
    serverURL,
    client.WithHTTPClient(httpClient),
)
```

### Request Editors
```go
// Add custom headers to all requests
requestEditor := func(ctx context.Context, req *http.Request) error {
    req.Header.Set("X-Custom-Header", "value")
    return nil
}

resp, err := fireflyClient.ListAccountWithResponse(
    ctx, 
    params, 
    requestEditor,
)
```

## Development

### Regenerating the Client

The project includes a `//go:generate` directive for easy client regeneration. You can regenerate the client in two ways:

#### Using go generate (Recommended)

```bash
# Regenerate the client using the built-in generate directive
go generate ./pkg/client
```

This uses the `//go:generate` directive in `pkg/client/generate.go` to automatically regenerate the client code.

#### Manual Generation

If you need to regenerate the client manually or with custom parameters:

1. Update the OpenAPI spec file in `resources/`
2. Run the generation command:

```bash
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
    -generate client,models,embedded-spec \
    -package client \
    -o pkg/client/client.go \
    resources/firefly-iii-6.2.21-v1.yaml
```

### Running the Example

```bash
# Edit example/main.go with your Firefly III URL and API token
go run example/main.go
```

## API Documentation

For detailed API documentation, refer to:
- [Firefly III API Documentation](https://api-docs.firefly-iii.org/)
- [Generated client code](pkg/client/client.go) - contains all available methods and types

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the same terms as Firefly III.

## Support

- [Firefly III Documentation](https://docs.firefly-iii.org/)
- [Firefly III GitHub](https://github.com/firefly-iii/firefly-iii)
- [oapi-codegen Documentation](https://github.com/oapi-codegen/oapi-codegen)