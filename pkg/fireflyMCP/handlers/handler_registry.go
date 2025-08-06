package handlers

import (
	"github.com/dezer32/mcp-firefly-iii/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HandlerRegistryImpl implements the Registry interface
type HandlerRegistryImpl struct {
	AccountHandlers     *AccountHandlers
	TransactionHandlers *TransactionHandlers
	BudgetHandlers      *BudgetHandlers
	CategoryHandlers    *CategoryHandlers
	TagHandlers         *TagHandlers
	InsightHandlers     *InsightHandlers
	BillHandlers        *BillHandlers
	RecurrenceHandlers  *RecurrenceHandlers
}

// ServerContext adapts FireflyMCPServer to HandlerContext interface
type ServerContext struct {
	Client *client.ClientWithResponses
	Config interface{}
}

// GetClient returns the Firefly III API client
func (s *ServerContext) GetClient() interface{} {
	return s.Client
}

// GetConfig returns the application configuration
func (s *ServerContext) GetConfig() interface{} {
	return s.Config
}

// NewHandlerRegistry creates a new handler registry with all handlers
func NewHandlerRegistry(client *client.ClientWithResponses, config interface{}) *HandlerRegistryImpl {
	ctx := &ServerContext{
		Client: client,
		Config: config,
	}

	return &HandlerRegistryImpl{
		AccountHandlers:     NewAccountHandlers(ctx),
		TransactionHandlers: NewTransactionHandlers(ctx),
		BudgetHandlers:      NewBudgetHandlers(ctx),
		CategoryHandlers:    NewCategoryHandlers(ctx),
		TagHandlers:         NewTagHandlers(ctx),
		InsightHandlers:     NewInsightHandlers(ctx),
		BillHandlers:        NewBillHandlers(ctx),
		RecurrenceHandlers:  NewRecurrenceHandlers(ctx),
	}
}

// RegisterAll registers all available MCP tools with the server
func (r *HandlerRegistryImpl) RegisterAll(server *mcp.Server) {
	// Register account tools
	r.AccountHandlers.RegisterTools(server)
	
	// Register transaction tools
	r.TransactionHandlers.RegisterTools(server)
	
	// Register budget tools
	r.BudgetHandlers.RegisterTools(server)
	
	// Register category tools
	r.CategoryHandlers.RegisterTools(server)
	
	// Register tag tools
	r.TagHandlers.RegisterTools(server)
	
	// Register insight tools
	r.InsightHandlers.RegisterTools(server)
	
	// Register bill tools
	r.BillHandlers.RegisterTools(server)
	
	// Register recurrence tools
	r.RecurrenceHandlers.RegisterTools(server)
}

// GetAccountHandlers returns the account handlers instance
func (r *HandlerRegistryImpl) GetAccountHandlers() *AccountHandlers {
	return r.AccountHandlers
}

// GetTransactionHandlers returns the transaction handlers instance
func (r *HandlerRegistryImpl) GetTransactionHandlers() *TransactionHandlers {
	return r.TransactionHandlers
}

// GetBudgetHandlers returns the budget handlers instance
func (r *HandlerRegistryImpl) GetBudgetHandlers() *BudgetHandlers {
	return r.BudgetHandlers
}

// GetCategoryHandlers returns the category handlers instance
func (r *HandlerRegistryImpl) GetCategoryHandlers() *CategoryHandlers {
	return r.CategoryHandlers
}

// GetTagHandlers returns the tag handlers instance
func (r *HandlerRegistryImpl) GetTagHandlers() *TagHandlers {
	return r.TagHandlers
}

// GetInsightHandlers returns the insight handlers instance
func (r *HandlerRegistryImpl) GetInsightHandlers() *InsightHandlers {
	return r.InsightHandlers
}

// GetBillHandlers returns the bill handlers instance
func (r *HandlerRegistryImpl) GetBillHandlers() *BillHandlers {
	return r.BillHandlers
}

// GetRecurrenceHandlers returns the recurrence handlers instance
func (r *HandlerRegistryImpl) GetRecurrenceHandlers() *RecurrenceHandlers {
	return r.RecurrenceHandlers
}