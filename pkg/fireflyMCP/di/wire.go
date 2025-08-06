//go:build wireinject
// +build wireinject

// Package di provides dependency injection using Wire
package di

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP"
	"github.com/google/wire"
)

// InitializeServer creates a fully wired FireflyMCPServer
func InitializeServer(config *fireflyMCP.Config) (*fireflyMCP.FireflyMCPServer, error) {
	wire.Build(
		ProvideHTTPClient,
		ProvideFireflyClient,
		ProvideMCPServer,
		ProvideHandlerContext,
		ProvideAccountHandlers,
		ProvideTransactionHandlers,
		ProvideBudgetHandlers,
		ProvideCategoryHandlers,
		ProvideTagHandlers,
		ProvideInsightHandlers,
		ProvideBillHandlers,
		ProvideRecurrenceHandlers,
		ProvideHandlerRegistry,
		ProvideMiddlewareChain,
		ProvideFireflyMCPServer,
	)
	return nil, nil
}