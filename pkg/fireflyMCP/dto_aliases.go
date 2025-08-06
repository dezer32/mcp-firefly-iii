package fireflyMCP

import "github.com/dezer32/mcp-firefly-iii/pkg/fireflyMCP/dto"

// Type aliases for backward compatibility
type (
	Pagination               = dto.Pagination
	Spent                    = dto.Spent
	Budget                   = dto.Budget
	BudgetList               = dto.BudgetList
	Category                 = dto.Category
	CategoryList             = dto.CategoryList
	Account                  = dto.Account
	AccountList              = dto.AccountList
	Transaction              = dto.Transaction
	TransactionGroup         = dto.TransactionGroup
	TransactionList          = dto.TransactionList
	BasicSummary             = dto.BasicSummary
	BasicSummaryList         = dto.BasicSummaryList
	InsightCategoryEntry     = dto.InsightCategoryEntry
	InsightTotalEntry        = dto.InsightTotalEntry
	InsightCategoryResponse  = dto.InsightCategoryResponse
	InsightTotalResponse     = dto.InsightTotalResponse
	BudgetSpent              = dto.BudgetSpent
	BudgetLimit              = dto.BudgetLimit
	BudgetLimitList          = dto.BudgetLimitList
	Tag                      = dto.Tag
	TagList                  = dto.TagList
	PaidDate                 = dto.PaidDate
	Bill                     = dto.Bill
	BillList                 = dto.BillList
	RecurrenceRepetition     = dto.RecurrenceRepetition
	RecurrenceTransaction    = dto.RecurrenceTransaction
	Recurrence               = dto.Recurrence
	RecurrenceList           = dto.RecurrenceList
	TransactionStoreRequest  = dto.TransactionStoreRequest
	TransactionSplitRequest  = dto.TransactionSplitRequest
)