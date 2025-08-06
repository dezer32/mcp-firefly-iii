package dto

import "fmt"

// MCPEntity is a common interface for all MCP entities that have an ID and Name
type MCPEntity interface {
	GetID() string
	GetName() string
}

// Pageable is an interface for entities that support pagination
type Pageable interface {
	GetPagination() Pagination
	GetData() []MCPEntity
	GetCount() int
}

// Validatable is an interface for entities that can validate themselves
type Validatable interface {
	Validate() error
}

// Ensure Account implements MCPEntity
var _ MCPEntity = (*Account)(nil)

// GetID returns the account ID
func (a *Account) GetID() string {
	return a.Id
}

// GetName returns the account name
func (a *Account) GetName() string {
	return a.Name
}

// Ensure Budget implements MCPEntity
var _ MCPEntity = (*Budget)(nil)

// GetID returns the budget ID
func (b *Budget) GetID() string {
	return b.Id
}

// GetName returns the budget name
func (b *Budget) GetName() string {
	return b.Name
}

// Ensure Category implements MCPEntity
var _ MCPEntity = (*Category)(nil)

// GetID returns the category ID
func (c *Category) GetID() string {
	return c.Id
}

// GetName returns the category name
func (c *Category) GetName() string {
	return c.Name
}

// Ensure Tag implements MCPEntity
var _ MCPEntity = (*Tag)(nil)

// GetID returns the tag ID
func (t *Tag) GetID() string {
	return t.Id
}

// GetName returns the tag name
func (t *Tag) GetName() string {
	return t.Tag
}

// Ensure Bill implements MCPEntity
var _ MCPEntity = (*Bill)(nil)

// GetID returns the bill ID
func (b *Bill) GetID() string {
	return b.Id
}

// GetName returns the bill name
func (b *Bill) GetName() string {
	return b.Name
}

// Ensure Transaction implements MCPEntity
var _ MCPEntity = (*Transaction)(nil)

// GetID returns the transaction ID
func (t *Transaction) GetID() string {
	return t.Id
}

// GetName returns the transaction description
func (t *Transaction) GetName() string {
	return t.Description
}

// Ensure TransactionGroup implements MCPEntity
var _ MCPEntity = (*TransactionGroup)(nil)

// GetID returns the transaction group ID
func (tg *TransactionGroup) GetID() string {
	return tg.Id
}

// GetName returns the transaction group title
func (tg *TransactionGroup) GetName() string {
	return tg.GroupTitle
}

// Ensure Recurrence implements MCPEntity
var _ MCPEntity = (*Recurrence)(nil)

// GetID returns the recurrence ID
func (r *Recurrence) GetID() string {
	return r.Id
}

// GetName returns the recurrence title
func (r *Recurrence) GetName() string {
	return r.Title
}

// Ensure BasicSummary implements MCPEntity
var _ MCPEntity = (*BasicSummary)(nil)

// GetID returns the summary key as ID
func (bs *BasicSummary) GetID() string {
	return bs.Key
}

// GetName returns the summary title
func (bs *BasicSummary) GetName() string {
	return bs.Title
}

// Ensure InsightCategoryEntry implements MCPEntity
var _ MCPEntity = (*InsightCategoryEntry)(nil)

// GetID returns the insight category ID
func (ice *InsightCategoryEntry) GetID() string {
	return ice.Id
}

// GetName returns the insight category name
func (ice *InsightCategoryEntry) GetName() string {
	return ice.Name
}

// Pageable implementations for list types

// Ensure AccountList implements Pageable
var _ Pageable = (*AccountList)(nil)

// GetPagination returns the pagination info
func (al *AccountList) GetPagination() Pagination {
	return al.Pagination
}

// GetData returns the accounts as MCPEntity slice
func (al *AccountList) GetData() []MCPEntity {
	result := make([]MCPEntity, len(al.Data))
	for i := range al.Data {
		result[i] = &al.Data[i]
	}
	return result
}

// GetCount returns the number of items in the list
func (al *AccountList) GetCount() int {
	return len(al.Data)
}

// Ensure BudgetList implements Pageable
var _ Pageable = (*BudgetList)(nil)

// GetPagination returns the pagination info
func (bl *BudgetList) GetPagination() Pagination {
	return bl.Pagination
}

// GetData returns the budgets as MCPEntity slice
func (bl *BudgetList) GetData() []MCPEntity {
	result := make([]MCPEntity, len(bl.Data))
	for i := range bl.Data {
		result[i] = &bl.Data[i]
	}
	return result
}

// GetCount returns the number of items in the list
func (bl *BudgetList) GetCount() int {
	return len(bl.Data)
}

// Ensure CategoryList implements Pageable
var _ Pageable = (*CategoryList)(nil)

// GetPagination returns the pagination info
func (cl *CategoryList) GetPagination() Pagination {
	return cl.Pagination
}

// GetData returns the categories as MCPEntity slice
func (cl *CategoryList) GetData() []MCPEntity {
	result := make([]MCPEntity, len(cl.Data))
	for i := range cl.Data {
		result[i] = &cl.Data[i]
	}
	return result
}

// GetCount returns the number of items in the list
func (cl *CategoryList) GetCount() int {
	return len(cl.Data)
}

// Ensure TagList implements Pageable
var _ Pageable = (*TagList)(nil)

// GetPagination returns the pagination info
func (tl *TagList) GetPagination() Pagination {
	return tl.Pagination
}

// GetData returns the tags as MCPEntity slice
func (tl *TagList) GetData() []MCPEntity {
	result := make([]MCPEntity, len(tl.Data))
	for i := range tl.Data {
		result[i] = &tl.Data[i]
	}
	return result
}

// GetCount returns the number of items in the list
func (tl *TagList) GetCount() int {
	return len(tl.Data)
}

// Ensure BillList implements Pageable
var _ Pageable = (*BillList)(nil)

// GetPagination returns the pagination info
func (bl *BillList) GetPagination() Pagination {
	return bl.Pagination
}

// GetData returns the bills as MCPEntity slice
func (bl *BillList) GetData() []MCPEntity {
	result := make([]MCPEntity, len(bl.Data))
	for i := range bl.Data {
		result[i] = &bl.Data[i]
	}
	return result
}

// GetCount returns the number of items in the list
func (bl *BillList) GetCount() int {
	return len(bl.Data)
}

// Ensure RecurrenceList implements Pageable
var _ Pageable = (*RecurrenceList)(nil)

// GetPagination returns the pagination info
func (rl *RecurrenceList) GetPagination() Pagination {
	return rl.Pagination
}

// GetData returns the recurrences as MCPEntity slice
func (rl *RecurrenceList) GetData() []MCPEntity {
	result := make([]MCPEntity, len(rl.Data))
	for i := range rl.Data {
		result[i] = &rl.Data[i]
	}
	return result
}

// GetCount returns the number of items in the list
func (rl *RecurrenceList) GetCount() int {
	return len(rl.Data)
}

// Validatable implementations

// Validate validates an Account
func (a *Account) Validate() error {
	if a.Id == "" {
		return fmt.Errorf("account ID cannot be empty")
	}
	if a.Name == "" {
		return fmt.Errorf("account name cannot be empty")
	}
	if a.Type == "" {
		return fmt.Errorf("account type cannot be empty")
	}
	return nil
}

// Validate validates a Budget
func (b *Budget) Validate() error {
	if b.Id == "" {
		return fmt.Errorf("budget ID cannot be empty")
	}
	if b.Name == "" {
		return fmt.Errorf("budget name cannot be empty")
	}
	return nil
}

// Validate validates a Category
func (c *Category) Validate() error {
	if c.Id == "" {
		return fmt.Errorf("category ID cannot be empty")
	}
	if c.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	return nil
}

// Validate validates a Tag
func (t *Tag) Validate() error {
	if t.Id == "" {
		return fmt.Errorf("tag ID cannot be empty")
	}
	if t.Tag == "" {
		return fmt.Errorf("tag name cannot be empty")
	}
	return nil
}

// Validate validates a Bill
func (b *Bill) Validate() error {
	if b.Id == "" {
		return fmt.Errorf("bill ID cannot be empty")
	}
	if b.Name == "" {
		return fmt.Errorf("bill name cannot be empty")
	}
	if b.AmountMin == "" {
		return fmt.Errorf("bill minimum amount cannot be empty")
	}
	if b.AmountMax == "" {
		return fmt.Errorf("bill maximum amount cannot be empty")
	}
	if b.RepeatFreq == "" {
		return fmt.Errorf("bill repeat frequency cannot be empty")
	}
	return nil
}

// Validate validates a Transaction
func (t *Transaction) Validate() error {
	if t.Id == "" {
		return fmt.Errorf("transaction ID cannot be empty")
	}
	if t.Amount == "" {
		return fmt.Errorf("transaction amount cannot be empty")
	}
	if t.Description == "" {
		return fmt.Errorf("transaction description cannot be empty")
	}
	if t.Type == "" {
		return fmt.Errorf("transaction type cannot be empty")
	}
	if t.SourceId == "" {
		return fmt.Errorf("transaction source ID cannot be empty")
	}
	if t.DestinationId == "" {
		return fmt.Errorf("transaction destination ID cannot be empty")
	}
	return nil
}

// Validate validates a TransactionGroup
func (tg *TransactionGroup) Validate() error {
	if tg.Id == "" {
		return fmt.Errorf("transaction group ID cannot be empty")
	}
	if len(tg.Transactions) == 0 {
		return fmt.Errorf("transaction group must contain at least one transaction")
	}
	for i, t := range tg.Transactions {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("transaction %d validation failed: %w", i, err)
		}
	}
	return nil
}

// Validate validates a Recurrence
func (r *Recurrence) Validate() error {
	if r.Id == "" {
		return fmt.Errorf("recurrence ID cannot be empty")
	}
	if r.Title == "" {
		return fmt.Errorf("recurrence title cannot be empty")
	}
	if r.Type == "" {
		return fmt.Errorf("recurrence type cannot be empty")
	}
	if len(r.Repetitions) == 0 {
		return fmt.Errorf("recurrence must have at least one repetition")
	}
	if len(r.Transactions) == 0 {
		return fmt.Errorf("recurrence must have at least one transaction")
	}
	return nil
}

// Validate validates a TransactionStoreRequest
func (tsr *TransactionStoreRequest) Validate() error {
	if len(tsr.Transactions) == 0 {
		return fmt.Errorf("transaction store request must contain at least one transaction")
	}
	for i, t := range tsr.Transactions {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("transaction %d validation failed: %w", i, err)
		}
	}
	return nil
}

// Validate validates a TransactionSplitRequest
func (tsr *TransactionSplitRequest) Validate() error {
	if tsr.Type == "" {
		return fmt.Errorf("transaction type is required")
	}
	if tsr.Type != "withdrawal" && tsr.Type != "deposit" && tsr.Type != "transfer" {
		return fmt.Errorf("transaction type must be one of: withdrawal, deposit, transfer")
	}
	if tsr.Date == "" {
		return fmt.Errorf("transaction date is required")
	}
	if tsr.Amount == "" {
		return fmt.Errorf("transaction amount is required")
	}
	if tsr.Description == "" {
		return fmt.Errorf("transaction description is required")
	}
	
	// At least one source identifier must be provided
	if tsr.SourceId == nil && tsr.SourceName == nil {
		return fmt.Errorf("either source_id or source_name must be provided")
	}
	
	// At least one destination identifier must be provided
	if tsr.DestinationId == nil && tsr.DestinationName == nil {
		return fmt.Errorf("either destination_id or destination_name must be provided")
	}
	
	return nil
}

// Validate validates a Pagination
func (p *Pagination) Validate() error {
	if p.PerPage <= 0 {
		return fmt.Errorf("per_page must be greater than 0")
	}
	if p.CurrentPage <= 0 {
		return fmt.Errorf("current_page must be greater than 0")
	}
	if p.CurrentPage > p.TotalPages && p.TotalPages > 0 {
		return fmt.Errorf("current_page cannot be greater than total_pages")
	}
	return nil
}