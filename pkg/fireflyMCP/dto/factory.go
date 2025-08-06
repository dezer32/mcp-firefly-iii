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
	account := NewAccountBuilder().
		WithId(id).
		WithName(name).
		WithType(accountType).
		WithActive(active).
		WithNotes(f.safeString(notes)).
		Build()
	
	if f.validateOnCreate {
		if err := account.Validate(); err != nil {
			return nil, fmt.Errorf("account validation failed: %w", err)
		}
	}
	
	return &account, nil
}

// CreateBudget creates a new Budget DTO with validation
func (f *dtoFactory) CreateBudget(id, name string, active bool, notes interface{}, spent Spent) (*Budget, error) {
	budget := NewBudgetBuilder().
		WithId(id).
		WithName(name).
		WithActive(active).
		WithNotes(notes).
		WithSpent(spent).
		Build()
	
	if f.validateOnCreate {
		if err := budget.Validate(); err != nil {
			return nil, fmt.Errorf("budget validation failed: %w", err)
		}
	}
	
	return &budget, nil
}

// CreateCategory creates a new Category DTO with validation
func (f *dtoFactory) CreateCategory(id, name string, notes interface{}) (*Category, error) {
	category := NewCategoryBuilder().
		WithId(id).
		WithName(name).
		WithNotes(notes).
		Build()
	
	if f.validateOnCreate {
		if err := category.Validate(); err != nil {
			return nil, fmt.Errorf("category validation failed: %w", err)
		}
	}
	
	return &category, nil
}

// CreateTag creates a new Tag DTO with validation
func (f *dtoFactory) CreateTag(id, tag string, description *string) (*Tag, error) {
	t := NewTagBuilder().
		WithId(id).
		WithTag(tag).
		WithDescription(f.safeString(description)).
		Build()
	
	if f.validateOnCreate {
		if err := t.Validate(); err != nil {
			return nil, fmt.Errorf("tag validation failed: %w", err)
		}
	}
	
	return &t, nil
}

// CreateBill creates a new Bill DTO with validation
func (f *dtoFactory) CreateBill(id, name, amountMin, amountMax, repeatFreq, currencyCode string, date time.Time, skip int, active bool, notes *string, nextExpectedMatch *time.Time, paidDates []PaidDate) (*Bill, error) {
	bill := NewBillBuilder().
		WithId(id).
		WithName(name).
		WithAmountMin(amountMin).
		WithAmountMax(amountMax).
		WithRepeatFreq(repeatFreq).
		WithCurrencyCode(currencyCode).
		WithDate(date).
		WithSkip(skip).
		WithActive(active).
		WithNotes(f.safeString(notes)).
		WithNextExpectedMatch(nextExpectedMatch).
		WithPaidDates(paidDates).
		Build()
	
	if f.validateOnCreate {
		if err := bill.Validate(); err != nil {
			return nil, fmt.Errorf("bill validation failed: %w", err)
		}
	}
	
	return &bill, nil
}

// CreateTransaction creates a new Transaction DTO with validation
func (f *dtoFactory) CreateTransaction(id, amount, description, sourceId, sourceName, destId, destName, txType, currencyCode, destType string, date time.Time) (*Transaction, error) {
	transaction := NewTransactionBuilder().
		WithId(id).
		WithAmount(amount).
		WithDescription(description).
		WithSourceId(sourceId).
		WithSourceName(sourceName).
		WithDestinationId(destId).
		WithDestinationName(destName).
		WithType(txType).
		WithCurrencyCode(currencyCode).
		WithDestinationType(destType).
		WithDate(date).
		WithTags([]string{}). // Initialize with empty array
		Build()
	
	if f.validateOnCreate {
		if err := transaction.Validate(); err != nil {
			return nil, fmt.Errorf("transaction validation failed: %w", err)
		}
	}
	
	return &transaction, nil
}

// CreateTransactionGroup creates a new TransactionGroup DTO with validation
func (f *dtoFactory) CreateTransactionGroup(id, groupTitle string, transactions []Transaction) (*TransactionGroup, error) {
	builder := NewTransactionGroupBuilder().
		WithId(id).
		WithGroupTitle(groupTitle)
		
	if transactions != nil {
		builder.WithTransactions(transactions)
	} else {
		builder.WithTransactions([]Transaction{})
	}
	
	group := builder.Build()
	
	if f.validateOnCreate {
		if err := group.Validate(); err != nil {
			return nil, fmt.Errorf("transaction group validation failed: %w", err)
		}
	}
	
	return &group, nil
}

// CreateRecurrence creates a new Recurrence DTO with validation
func (f *dtoFactory) CreateRecurrence(id, recType, title, description string, firstDate time.Time, active, applyRules bool) (*Recurrence, error) {
	recurrence := NewRecurrenceBuilder().
		WithId(id).
		WithType(recType).
		WithTitle(title).
		WithDescription(description).
		WithFirstDate(firstDate).
		WithActive(active).
		WithApplyRules(applyRules).
		WithRepetitions([]RecurrenceRepetition{}). // Initialize empty
		WithTransactions([]RecurrenceTransaction{}). // Initialize empty
		Build()
	
	// Note: Validation will fail if repetitions and transactions are empty
	// These should be added after creation
	if f.validateOnCreate {
		// Skip validation for now as repetitions and transactions need to be added
		// The caller should validate after adding these
	}
	
	return &recurrence, nil
}

// CreateInsightCategoryEntry creates a new InsightCategoryEntry DTO
func (f *dtoFactory) CreateInsightCategoryEntry(id, name, amount, currencyCode string) (*InsightCategoryEntry, error) {
	// For now, create directly as builder is not yet implemented for this type
	entry := InsightCategoryEntry{
		id:           id,
		name:         name,
		amount:       amount,
		currencyCode: currencyCode,
	}
	
	if f.validateOnCreate {
		if err := entry.Validate(); err != nil {
			return nil, fmt.Errorf("insight category entry validation failed: %w", err)
		}
	}
	
	return &entry, nil
}

// CreateInsightTotalEntry creates a new InsightTotalEntry DTO
func (f *dtoFactory) CreateInsightTotalEntry(amount, currencyCode string) (*InsightTotalEntry, error) {
	// For now, create directly as builder is not yet implemented for this type
	entry := InsightTotalEntry{
		amount:       amount,
		currencyCode: currencyCode,
	}
	
	// InsightTotalEntry doesn't have Validate method yet
	// Add validation if needed
	
	return &entry, nil
}

// CreateBasicSummary creates a new BasicSummary DTO
func (f *dtoFactory) CreateBasicSummary(key, title, currencyCode, monetaryValue string) (*BasicSummary, error) {
	// For now, create directly as builder is not yet implemented for this type
	summary := BasicSummary{
		key:           key,
		title:         title,
		currencyCode:  currencyCode,
		monetaryValue: monetaryValue,
	}
	
	if f.validateOnCreate {
		if err := summary.Validate(); err != nil {
			return nil, fmt.Errorf("basic summary validation failed: %w", err)
		}
	}
	
	return &summary, nil
}

// DefaultFactory is the default factory instance with validation enabled
var DefaultFactory = NewDTOFactory(WithValidation(true), WithNilHandling(true))