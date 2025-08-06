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
	GetCount() int
}

// Validatable is an interface for entities that can validate themselves
type Validatable interface {
	Validate() error
}

// Ensure Account implements MCPEntity
var _ MCPEntity = Account{}

// GetID returns the account ID (implements MCPEntity)
func (a Account) GetID() string {
	return a.GetId()
}

// GetName is already implemented in dto.go

// Ensure Budget implements MCPEntity
var _ MCPEntity = Budget{}

// GetID returns the budget ID (implements MCPEntity)
func (b Budget) GetID() string {
	return b.GetId()
}

// GetName is already implemented in dto.go

// Ensure Category implements MCPEntity
var _ MCPEntity = Category{}

// GetID returns the category ID (implements MCPEntity)
func (c Category) GetID() string {
	return c.GetId()
}

// GetName is already implemented in dto.go

// Ensure Tag implements MCPEntity
var _ MCPEntity = Tag{}

// GetID returns the tag ID (implements MCPEntity)
func (t Tag) GetID() string {
	return t.GetId()
}

// GetName returns the tag name (implements MCPEntity)
func (t Tag) GetName() string {
	return t.GetTag()
}

// Ensure Bill implements MCPEntity
var _ MCPEntity = Bill{}

// GetID returns the bill ID (implements MCPEntity)
func (b Bill) GetID() string {
	return b.GetId()
}

// GetName is already implemented in dto.go

// Ensure Transaction implements MCPEntity
var _ MCPEntity = Transaction{}

// GetID returns the transaction ID (implements MCPEntity)
func (t Transaction) GetID() string {
	return t.GetId()
}

// GetName returns the transaction description (implements MCPEntity)
func (t Transaction) GetName() string {
	return t.GetDescription()
}

// Ensure TransactionGroup implements MCPEntity
var _ MCPEntity = TransactionGroup{}

// GetID returns the transaction group ID (implements MCPEntity)
func (g TransactionGroup) GetID() string {
	return g.GetId()
}

// GetName returns the transaction group title (implements MCPEntity)
func (g TransactionGroup) GetName() string {
	return g.GetGroupTitle()
}

// Ensure Recurrence implements MCPEntity
var _ MCPEntity = Recurrence{}

// GetID returns the recurrence ID (implements MCPEntity)
func (r Recurrence) GetID() string {
	return r.GetId()
}

// GetName returns the recurrence title (implements MCPEntity)
func (r Recurrence) GetName() string {
	return r.GetTitle()
}

// Ensure BasicSummary implements MCPEntity
var _ MCPEntity = BasicSummary{}

// GetID returns the basic summary key as ID (implements MCPEntity)
func (s BasicSummary) GetID() string {
	return s.GetKey()
}

// GetName returns the basic summary title (implements MCPEntity)
func (s BasicSummary) GetName() string {
	return s.GetTitle()
}

// Ensure InsightCategoryEntry implements MCPEntity
var _ MCPEntity = InsightCategoryEntry{}

// GetID returns the insight category entry ID (implements MCPEntity)
func (e InsightCategoryEntry) GetID() string {
	return e.GetId()
}

// GetName is already implemented in dto.go

// Ensure AccountList implements Pageable
var _ Pageable = AccountList{}

// GetCount returns the number of items in the list (implements Pageable)
func (l AccountList) GetCount() int {
	return l.GetPagination().GetCount()
}

// Ensure BudgetList implements Pageable
var _ Pageable = BudgetList{}

// GetCount returns the number of items in the list (implements Pageable)
func (l BudgetList) GetCount() int {
	return l.GetPagination().GetCount()
}

// Ensure CategoryList implements Pageable
var _ Pageable = CategoryList{}

// GetCount returns the number of items in the list (implements Pageable)
func (l CategoryList) GetCount() int {
	return l.GetPagination().GetCount()
}

// Ensure TransactionList implements Pageable
var _ Pageable = TransactionList{}

// GetCount returns the number of items in the list (implements Pageable)
func (l TransactionList) GetCount() int {
	return l.GetPagination().GetCount()
}

// Validation implementations

// Validate validates an Account
func (a Account) Validate() error {
	if a.GetId() == "" {
		return fmt.Errorf("account ID is required")
	}
	if a.GetName() == "" {
		return fmt.Errorf("account name is required")
	}
	if a.GetType() == "" {
		return fmt.Errorf("account type is required")
	}
	return nil
}

// Validate validates a Budget
func (b Budget) Validate() error {
	if b.GetId() == "" {
		return fmt.Errorf("budget ID is required")
	}
	if b.GetName() == "" {
		return fmt.Errorf("budget name is required")
	}
	return nil
}

// Validate validates a Category
func (c Category) Validate() error {
	if c.GetId() == "" {
		return fmt.Errorf("category ID is required")
	}
	if c.GetName() == "" {
		return fmt.Errorf("category name is required")
	}
	return nil
}

// Validate validates a Tag
func (t Tag) Validate() error {
	if t.GetId() == "" {
		return fmt.Errorf("tag ID is required")
	}
	if t.GetTag() == "" {
		return fmt.Errorf("tag name is required")
	}
	return nil
}

// Validate validates a Bill
func (b Bill) Validate() error {
	if b.GetId() == "" {
		return fmt.Errorf("bill ID is required")
	}
	if b.GetName() == "" {
		return fmt.Errorf("bill name is required")
	}
	if b.GetAmountMin() == "" {
		return fmt.Errorf("bill minimum amount is required")
	}
	if b.GetAmountMax() == "" {
		return fmt.Errorf("bill maximum amount is required")
	}
	if b.GetCurrencyCode() == "" {
		return fmt.Errorf("bill currency code is required")
	}
	return nil
}

// Validate validates a Transaction
func (t Transaction) Validate() error {
	if t.GetId() == "" {
		return fmt.Errorf("transaction ID is required")
	}
	if t.GetAmount() == "" {
		return fmt.Errorf("transaction amount is required")
	}
	if t.GetDescription() == "" {
		return fmt.Errorf("transaction description is required")
	}
	if t.GetType() == "" {
		return fmt.Errorf("transaction type is required")
	}
	if t.GetSourceId() == "" {
		return fmt.Errorf("transaction source ID is required")
	}
	if t.GetDestinationId() == "" {
		return fmt.Errorf("transaction destination ID is required")
	}
	return nil
}

// Validate validates a TransactionGroup
func (g TransactionGroup) Validate() error {
	if g.GetId() == "" {
		return fmt.Errorf("transaction group ID is required")
	}
	if len(g.GetTransactions()) == 0 {
		return fmt.Errorf("transaction group must have at least one transaction")
	}
	for i, t := range g.GetTransactions() {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("transaction %d: %w", i, err)
		}
	}
	return nil
}

// Validate validates a Recurrence
func (r Recurrence) Validate() error {
	if r.GetId() == "" {
		return fmt.Errorf("recurrence ID is required")
	}
	if r.GetTitle() == "" {
		return fmt.Errorf("recurrence title is required")
	}
	if r.GetType() == "" {
		return fmt.Errorf("recurrence type is required")
	}
	if len(r.GetRepetitions()) == 0 {
		return fmt.Errorf("recurrence must have at least one repetition")
	}
	if len(r.GetTransactions()) == 0 {
		return fmt.Errorf("recurrence must have at least one transaction")
	}
	return nil
}

// Validate validates a BasicSummary
func (s BasicSummary) Validate() error {
	if s.GetKey() == "" {
		return fmt.Errorf("summary key is required")
	}
	if s.GetTitle() == "" {
		return fmt.Errorf("summary title is required")
	}
	return nil
}

// Validate validates an InsightCategoryEntry
func (e InsightCategoryEntry) Validate() error {
	if e.GetId() == "" {
		return fmt.Errorf("insight category entry ID is required")
	}
	if e.GetName() == "" {
		return fmt.Errorf("insight category entry name is required")
	}
	return nil
}

// Validate validates a Pagination
func (p Pagination) Validate() error {
	if p.GetPerPage() <= 0 {
		return fmt.Errorf("per page must be positive")
	}
	if p.GetCurrentPage() <= 0 {
		return fmt.Errorf("current page must be positive")
	}
	if p.GetTotalPages() < 0 {
		return fmt.Errorf("total pages cannot be negative")
	}
	if p.GetTotal() < 0 {
		return fmt.Errorf("total cannot be negative")
	}
	if p.GetCount() < 0 {
		return fmt.Errorf("count cannot be negative")
	}
	return nil
}

// Validate validates an AccountList
func (l AccountList) Validate() error {
	if err := l.GetPagination().Validate(); err != nil {
		return fmt.Errorf("pagination: %w", err)
	}
	for i, a := range l.GetData() {
		if err := a.Validate(); err != nil {
			return fmt.Errorf("account %d: %w", i, err)
		}
	}
	return nil
}

// Validate validates a BudgetList
func (l BudgetList) Validate() error {
	if err := l.GetPagination().Validate(); err != nil {
		return fmt.Errorf("pagination: %w", err)
	}
	for i, b := range l.GetData() {
		if err := b.Validate(); err != nil {
			return fmt.Errorf("budget %d: %w", i, err)
		}
	}
	return nil
}

// Validate validates a CategoryList
func (l CategoryList) Validate() error {
	if err := l.GetPagination().Validate(); err != nil {
		return fmt.Errorf("pagination: %w", err)
	}
	for i, c := range l.GetData() {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("category %d: %w", i, err)
		}
	}
	return nil
}

// Validate validates a TransactionList
func (l TransactionList) Validate() error {
	if err := l.GetPagination().Validate(); err != nil {
		return fmt.Errorf("pagination: %w", err)
	}
	for i, g := range l.GetData() {
		if err := g.Validate(); err != nil {
			return fmt.Errorf("transaction group %d: %w", i, err)
		}
	}
	return nil
}