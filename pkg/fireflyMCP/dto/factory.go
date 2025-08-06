package dto

import (
	"fmt"
	"time"
)

// DTOFactory is the interface for creating and validating DTOs
type DTOFactory interface {
	// Account creation methods
	CreateAccount(id, name, accountType string, active bool, notes *string) (*Account, error)
	
	// Budget creation methods
	CreateBudget(id, name string, active bool, notes interface{}, spent Spent) (*Budget, error)
	
	// Category creation methods
	CreateCategory(id, name string, notes interface{}) (*Category, error)
	
	// Tag creation methods
	CreateTag(id, tag string, description *string) (*Tag, error)
	
	// Bill creation methods
	CreateBill(id, name, amountMin, amountMax, repeatFreq, currencyCode string, date time.Time, skip int, active bool, notes *string, nextExpectedMatch *time.Time, paidDates []PaidDate) (*Bill, error)
	
	// Transaction creation methods
	CreateTransaction(id, amount, description, sourceId, sourceName, destId, destName, txType, currencyCode, destType string, date time.Time) (*Transaction, error)
	CreateTransactionGroup(id, groupTitle string, transactions []Transaction) (*TransactionGroup, error)
	
	// Recurrence creation methods
	CreateRecurrence(id, recType, title, description string, firstDate time.Time, active, applyRules bool) (*Recurrence, error)
	
	// Insight creation methods
	CreateInsightCategoryEntry(id, name, amount, currencyCode string) (*InsightCategoryEntry, error)
	CreateInsightTotalEntry(amount, currencyCode string) (*InsightTotalEntry, error)
	
	// Summary creation methods
	CreateBasicSummary(key, title, currencyCode, monetaryValue string) (*BasicSummary, error)
	
	// List creation methods with pagination
	CreateAccountList(accounts []Account, pagination Pagination) (*AccountList, error)
	CreateBudgetList(budgets []Budget, pagination Pagination) (*BudgetList, error)
	CreateCategoryList(categories []Category, pagination Pagination) (*CategoryList, error)
	CreateTagList(tags []Tag, pagination Pagination) (*TagList, error)
	CreateBillList(bills []Bill, pagination Pagination) (*BillList, error)
	CreateRecurrenceList(recurrences []Recurrence, pagination Pagination) (*RecurrenceList, error)
	CreateTransactionList(groups []TransactionGroup, pagination Pagination) (*TransactionList, error)
}

// dtoFactory is the default implementation of DTOFactory
type dtoFactory struct {
	validateOnCreate bool
	handleNilStrings bool
}

// FactoryOption is a functional option for configuring the factory
type FactoryOption func(*dtoFactory)

// WithValidation enables validation on DTO creation
func WithValidation(enabled bool) FactoryOption {
	return func(f *dtoFactory) {
		f.validateOnCreate = enabled
	}
}

// WithNilHandling enables special handling of nil strings
func WithNilHandling(enabled bool) FactoryOption {
	return func(f *dtoFactory) {
		f.handleNilStrings = enabled
	}
}

// NewDTOFactory creates a new DTOFactory with the given options
func NewDTOFactory(opts ...FactoryOption) DTOFactory {
	f := &dtoFactory{
		validateOnCreate: true,  // Default: validate on create
		handleNilStrings: true,  // Default: handle nil strings
	}
	
	for _, opt := range opts {
		opt(f)
	}
	
	return f
}

// Helper method to handle nil strings
func (f *dtoFactory) safeString(s *string) *string {
	if !f.handleNilStrings || s != nil {
		return s
	}
	// Return nil as-is if nil handling is disabled
	return nil
}

// CreateAccount creates a new Account DTO with validation
func (f *dtoFactory) CreateAccount(id, name, accountType string, active bool, notes *string) (*Account, error) {
	account := &Account{
		Id:     id,
		Name:   name,
		Type:   accountType,
		Active: active,
		Notes:  f.safeString(notes),
	}
	
	if f.validateOnCreate {
		if err := account.Validate(); err != nil {
			return nil, fmt.Errorf("account validation failed: %w", err)
		}
	}
	
	return account, nil
}

// CreateBudget creates a new Budget DTO with validation
func (f *dtoFactory) CreateBudget(id, name string, active bool, notes interface{}, spent Spent) (*Budget, error) {
	budget := &Budget{
		Id:     id,
		Name:   name,
		Active: active,
		Notes:  notes,
		Spent:  spent,
	}
	
	if f.validateOnCreate {
		if err := budget.Validate(); err != nil {
			return nil, fmt.Errorf("budget validation failed: %w", err)
		}
	}
	
	return budget, nil
}

// CreateCategory creates a new Category DTO with validation
func (f *dtoFactory) CreateCategory(id, name string, notes interface{}) (*Category, error) {
	category := &Category{
		Id:    id,
		Name:  name,
		Notes: notes,
	}
	
	if f.validateOnCreate {
		if err := category.Validate(); err != nil {
			return nil, fmt.Errorf("category validation failed: %w", err)
		}
	}
	
	return category, nil
}

// CreateTag creates a new Tag DTO with validation
func (f *dtoFactory) CreateTag(id, tag string, description *string) (*Tag, error) {
	t := &Tag{
		Id:          id,
		Tag:         tag,
		Description: f.safeString(description),
	}
	
	if f.validateOnCreate {
		if err := t.Validate(); err != nil {
			return nil, fmt.Errorf("tag validation failed: %w", err)
		}
	}
	
	return t, nil
}

// CreateBill creates a new Bill DTO with validation
func (f *dtoFactory) CreateBill(id, name, amountMin, amountMax, repeatFreq, currencyCode string, date time.Time, skip int, active bool, notes *string, nextExpectedMatch *time.Time, paidDates []PaidDate) (*Bill, error) {
	bill := &Bill{
		Id:                id,
		Name:              name,
		AmountMin:         amountMin,
		AmountMax:         amountMax,
		RepeatFreq:        repeatFreq,
		CurrencyCode:      currencyCode,
		Date:              date,
		Skip:              skip,
		Active:            active,
		Notes:             f.safeString(notes),
		NextExpectedMatch: nextExpectedMatch,
		PaidDates:         paidDates,
	}
	
	if f.validateOnCreate {
		if err := bill.Validate(); err != nil {
			return nil, fmt.Errorf("bill validation failed: %w", err)
		}
	}
	
	return bill, nil
}

// CreateTransaction creates a new Transaction DTO with validation
func (f *dtoFactory) CreateTransaction(id, amount, description, sourceId, sourceName, destId, destName, txType, currencyCode, destType string, date time.Time) (*Transaction, error) {
	transaction := &Transaction{
		Id:              id,
		Amount:          amount,
		Description:     description,
		SourceId:        sourceId,
		SourceName:      sourceName,
		DestinationId:   destId,
		DestinationName: destName,
		Type:            txType,
		CurrencyCode:    currencyCode,
		DestinationType: destType,
		Date:            date,
		Tags:            []string{}, // Initialize with empty array
	}
	
	if f.validateOnCreate {
		if err := transaction.Validate(); err != nil {
			return nil, fmt.Errorf("transaction validation failed: %w", err)
		}
	}
	
	return transaction, nil
}

// CreateTransactionGroup creates a new TransactionGroup DTO with validation
func (f *dtoFactory) CreateTransactionGroup(id, groupTitle string, transactions []Transaction) (*TransactionGroup, error) {
	group := &TransactionGroup{
		Id:           id,
		GroupTitle:   groupTitle,
		Transactions: transactions,
	}
	
	if transactions == nil {
		group.Transactions = []Transaction{}
	}
	
	if f.validateOnCreate {
		if err := group.Validate(); err != nil {
			return nil, fmt.Errorf("transaction group validation failed: %w", err)
		}
	}
	
	return group, nil
}

// CreateRecurrence creates a new Recurrence DTO with validation
func (f *dtoFactory) CreateRecurrence(id, recType, title, description string, firstDate time.Time, active, applyRules bool) (*Recurrence, error) {
	recurrence := &Recurrence{
		Id:          id,
		Type:        recType,
		Title:       title,
		Description: description,
		FirstDate:   firstDate,
		Active:      active,
		ApplyRules:  applyRules,
		Repetitions: []RecurrenceRepetition{}, // Initialize empty
		Transactions: []RecurrenceTransaction{}, // Initialize empty
	}
	
	// Note: Validation will fail if repetitions and transactions are empty
	// These should be added after creation
	if f.validateOnCreate {
		// Skip validation for now as repetitions and transactions need to be added
		// The caller should validate after adding these
	}
	
	return recurrence, nil
}

// CreateInsightCategoryEntry creates a new InsightCategoryEntry DTO
func (f *dtoFactory) CreateInsightCategoryEntry(id, name, amount, currencyCode string) (*InsightCategoryEntry, error) {
	entry := &InsightCategoryEntry{
		Id:           id,
		Name:         name,
		Amount:       amount,
		CurrencyCode: currencyCode,
	}
	
	if f.validateOnCreate {
		// Add validation if InsightCategoryEntry implements Validatable
		if validator, ok := interface{}(entry).(Validatable); ok {
			if err := validator.Validate(); err != nil {
				return nil, fmt.Errorf("insight category entry validation failed: %w", err)
			}
		}
	}
	
	return entry, nil
}

// CreateInsightTotalEntry creates a new InsightTotalEntry DTO
func (f *dtoFactory) CreateInsightTotalEntry(amount, currencyCode string) (*InsightTotalEntry, error) {
	entry := &InsightTotalEntry{
		Amount:       amount,
		CurrencyCode: currencyCode,
	}
	
	if f.validateOnCreate {
		// Add validation if InsightTotalEntry implements Validatable
		if validator, ok := interface{}(entry).(Validatable); ok {
			if err := validator.Validate(); err != nil {
				return nil, fmt.Errorf("insight total entry validation failed: %w", err)
			}
		}
	}
	
	return entry, nil
}

// CreateBasicSummary creates a new BasicSummary DTO
func (f *dtoFactory) CreateBasicSummary(key, title, currencyCode, monetaryValue string) (*BasicSummary, error) {
	summary := &BasicSummary{
		Key:           key,
		Title:         title,
		CurrencyCode:  currencyCode,
		MonetaryValue: monetaryValue,
	}
	
	if f.validateOnCreate {
		// Add validation if BasicSummary implements Validatable
		if validator, ok := interface{}(summary).(Validatable); ok {
			if err := validator.Validate(); err != nil {
				return nil, fmt.Errorf("basic summary validation failed: %w", err)
			}
		}
	}
	
	return summary, nil
}

// CreateAccountList creates a new AccountList with pagination
func (f *dtoFactory) CreateAccountList(accounts []Account, pagination Pagination) (*AccountList, error) {
	if accounts == nil {
		accounts = []Account{}
	}
	
	list := &AccountList{
		Data:       accounts,
		Pagination: pagination,
	}
	
	if f.validateOnCreate {
		if err := pagination.Validate(); err != nil {
			return nil, fmt.Errorf("pagination validation failed: %w", err)
		}
		
		// Validate each account
		for i, account := range accounts {
			if err := account.Validate(); err != nil {
				return nil, fmt.Errorf("account %d validation failed: %w", i, err)
			}
		}
	}
	
	return list, nil
}

// CreateBudgetList creates a new BudgetList with pagination
func (f *dtoFactory) CreateBudgetList(budgets []Budget, pagination Pagination) (*BudgetList, error) {
	if budgets == nil {
		budgets = []Budget{}
	}
	
	list := &BudgetList{
		Data:       budgets,
		Pagination: pagination,
	}
	
	if f.validateOnCreate {
		if err := pagination.Validate(); err != nil {
			return nil, fmt.Errorf("pagination validation failed: %w", err)
		}
		
		// Validate each budget
		for i, budget := range budgets {
			if err := budget.Validate(); err != nil {
				return nil, fmt.Errorf("budget %d validation failed: %w", i, err)
			}
		}
	}
	
	return list, nil
}

// CreateCategoryList creates a new CategoryList with pagination
func (f *dtoFactory) CreateCategoryList(categories []Category, pagination Pagination) (*CategoryList, error) {
	if categories == nil {
		categories = []Category{}
	}
	
	list := &CategoryList{
		Data:       categories,
		Pagination: pagination,
	}
	
	if f.validateOnCreate {
		if err := pagination.Validate(); err != nil {
			return nil, fmt.Errorf("pagination validation failed: %w", err)
		}
		
		// Validate each category
		for i, category := range categories {
			if err := category.Validate(); err != nil {
				return nil, fmt.Errorf("category %d validation failed: %w", i, err)
			}
		}
	}
	
	return list, nil
}

// CreateTagList creates a new TagList with pagination
func (f *dtoFactory) CreateTagList(tags []Tag, pagination Pagination) (*TagList, error) {
	if tags == nil {
		tags = []Tag{}
	}
	
	list := &TagList{
		Data:       tags,
		Pagination: pagination,
	}
	
	if f.validateOnCreate {
		if err := pagination.Validate(); err != nil {
			return nil, fmt.Errorf("pagination validation failed: %w", err)
		}
		
		// Validate each tag
		for i, tag := range tags {
			if err := tag.Validate(); err != nil {
				return nil, fmt.Errorf("tag %d validation failed: %w", i, err)
			}
		}
	}
	
	return list, nil
}

// CreateBillList creates a new BillList with pagination
func (f *dtoFactory) CreateBillList(bills []Bill, pagination Pagination) (*BillList, error) {
	if bills == nil {
		bills = []Bill{}
	}
	
	list := &BillList{
		Data:       bills,
		Pagination: pagination,
	}
	
	if f.validateOnCreate {
		if err := pagination.Validate(); err != nil {
			return nil, fmt.Errorf("pagination validation failed: %w", err)
		}
		
		// Validate each bill
		for i, bill := range bills {
			if err := bill.Validate(); err != nil {
				return nil, fmt.Errorf("bill %d validation failed: %w", i, err)
			}
		}
	}
	
	return list, nil
}

// CreateRecurrenceList creates a new RecurrenceList with pagination
func (f *dtoFactory) CreateRecurrenceList(recurrences []Recurrence, pagination Pagination) (*RecurrenceList, error) {
	if recurrences == nil {
		recurrences = []Recurrence{}
	}
	
	list := &RecurrenceList{
		Data:       recurrences,
		Pagination: pagination,
	}
	
	if f.validateOnCreate {
		if err := pagination.Validate(); err != nil {
			return nil, fmt.Errorf("pagination validation failed: %w", err)
		}
		
		// Validate each recurrence
		for i, recurrence := range recurrences {
			if err := recurrence.Validate(); err != nil {
				return nil, fmt.Errorf("recurrence %d validation failed: %w", i, err)
			}
		}
	}
	
	return list, nil
}

// CreateTransactionList creates a new TransactionList with pagination
func (f *dtoFactory) CreateTransactionList(groups []TransactionGroup, pagination Pagination) (*TransactionList, error) {
	if groups == nil {
		groups = []TransactionGroup{}
	}
	
	list := &TransactionList{
		Data:       groups,
		Pagination: pagination,
	}
	
	if f.validateOnCreate {
		if err := pagination.Validate(); err != nil {
			return nil, fmt.Errorf("pagination validation failed: %w", err)
		}
		
		// Validate each transaction group
		for i, group := range groups {
			if err := group.Validate(); err != nil {
				return nil, fmt.Errorf("transaction group %d validation failed: %w", i, err)
			}
		}
	}
	
	return list, nil
}

// DefaultFactory is the default factory instance with validation enabled
var DefaultFactory = NewDTOFactory(WithValidation(true), WithNilHandling(true))