// Package di provides dependency injection setup using Wire
package di

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/handlers"
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/middleware"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ProvideHTTPClient creates an HTTP client with configured timeout
func ProvideHTTPClient(config *fireflyMCP.Config) *http.Client {
	return &http.Client{
		Timeout: time.Duration(config.Client.Timeout) * time.Second,
	}
}

// ProvideFireflyClient creates a Firefly III API client
func ProvideFireflyClient(config *fireflyMCP.Config, httpClient *http.Client) (*client.ClientWithResponses, error) {
	fireflyClient, err := client.NewClientWithResponses(
		config.Server.URL,
		client.WithHTTPClient(httpClient),
		client.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				req.Header.Set("Authorization", "Bearer "+config.API.Token)
				return nil
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firefly III client: %w", err)
	}
	return fireflyClient, nil
}

// ProvideMCPServer creates an MCP server instance
func ProvideMCPServer(config *fireflyMCP.Config) *mcp.Server {
	return mcp.NewServer(
		&mcp.Implementation{
			Name:    config.MCP.Name,
			Version: config.MCP.Version,
		}, nil,
	)
}

// ProvideHandlerContext creates a handler context
func ProvideHandlerContext(client *client.ClientWithResponses, config *fireflyMCP.Config) handlers.HandlerContext {
	return &handlers.ServerContext{
		Client: client,
		Config: config,
	}
}

// ProvideAccountHandlers creates account handlers
func ProvideAccountHandlers(ctx handlers.HandlerContext) *handlers.AccountHandlers {
	return handlers.NewAccountHandlers(ctx)
}

// ProvideTransactionHandlers creates transaction handlers
func ProvideTransactionHandlers(ctx handlers.HandlerContext) *handlers.TransactionHandlers {
	return handlers.NewTransactionHandlers(ctx)
}

// ProvideBudgetHandlers creates budget handlers
func ProvideBudgetHandlers(ctx handlers.HandlerContext) *handlers.BudgetHandlers {
	return handlers.NewBudgetHandlers(ctx)
}

// ProvideCategoryHandlers creates category handlers
func ProvideCategoryHandlers(ctx handlers.HandlerContext) *handlers.CategoryHandlers {
	return handlers.NewCategoryHandlers(ctx)
}

// ProvideTagHandlers creates tag handlers
func ProvideTagHandlers(ctx handlers.HandlerContext) *handlers.TagHandlers {
	return handlers.NewTagHandlers(ctx)
}

// ProvideInsightHandlers creates insight handlers
func ProvideInsightHandlers(ctx handlers.HandlerContext) *handlers.InsightHandlers {
	return handlers.NewInsightHandlers(ctx)
}

// ProvideBillHandlers creates bill handlers
func ProvideBillHandlers(ctx handlers.HandlerContext) *handlers.BillHandlers {
	return handlers.NewBillHandlers(ctx)
}

// ProvideRecurrenceHandlers creates recurrence handlers
func ProvideRecurrenceHandlers(ctx handlers.HandlerContext) *handlers.RecurrenceHandlers {
	return handlers.NewRecurrenceHandlers(ctx)
}

// ProvideHandlerRegistry creates a handler registry with all handlers
func ProvideHandlerRegistry(
	accountHandlers *handlers.AccountHandlers,
	transactionHandlers *handlers.TransactionHandlers,
	budgetHandlers *handlers.BudgetHandlers,
	categoryHandlers *handlers.CategoryHandlers,
	tagHandlers *handlers.TagHandlers,
	insightHandlers *handlers.InsightHandlers,
	billHandlers *handlers.BillHandlers,
	recurrenceHandlers *handlers.RecurrenceHandlers,
) *handlers.HandlerRegistryImpl {
	return &handlers.HandlerRegistryImpl{
		AccountHandlers:     accountHandlers,
		TransactionHandlers: transactionHandlers,
		BudgetHandlers:      budgetHandlers,
		CategoryHandlers:    categoryHandlers,
		TagHandlers:         tagHandlers,
		InsightHandlers:     insightHandlers,
		BillHandlers:        billHandlers,
		RecurrenceHandlers:  recurrenceHandlers,
	}
}

// ProvideMiddlewareChain creates the middleware chain
func ProvideMiddlewareChain(config *fireflyMCP.Config) *middleware.Chain {
	// Create middleware instances based on configuration
	middlewares := []middleware.Middleware{}

	// Always add recovery middleware first
	middlewares = append(middlewares, middleware.NewRecoveryMiddleware(nil, true))

	// Add logging middleware with INFO level by default
	middlewares = append(middlewares, middleware.NewLoggingMiddleware(nil, middleware.LogLevelInfo))

	// Add timing middleware
	middlewares = append(middlewares, middleware.NewTimingMiddleware(nil, 1*time.Second))

	// Add metrics middleware
	middlewares = append(middlewares, middleware.NewMetricsMiddleware())

	return middleware.NewChain(middlewares...)
}

// ProvideFireflyMCPServer creates the main server instance
func ProvideFireflyMCPServer(
	mcpServer *mcp.Server,
	client *client.ClientWithResponses,
	config *fireflyMCP.Config,
	registry *handlers.HandlerRegistryImpl,
	chain *middleware.Chain,
) *fireflyMCP.FireflyMCPServer {
	// Register all tools through the handler registry
	registry.RegisterAll(mcpServer)

	return &fireflyMCP.FireflyMCPServer{
		Server:   mcpServer,
		Client:   client,
		Config:   config,
		Handlers: registry,
		Chain:    chain,
	}
}