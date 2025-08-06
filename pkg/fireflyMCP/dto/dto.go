package dto

import (
	"encoding/json"
	"time"
)

// Pagination represents paginated response metadata with immutable fields
type Pagination struct {
	count       int
	total       int
	currentPage int
	perPage     int
	totalPages  int
}

// GetCount returns the count
func (p Pagination) GetCount() int { return p.count }

// GetTotal returns the total
func (p Pagination) GetTotal() int { return p.total }

// GetCurrentPage returns the current page
func (p Pagination) GetCurrentPage() int { return p.currentPage }

// GetPerPage returns the per page count
func (p Pagination) GetPerPage() int { return p.perPage }

// GetTotalPages returns the total pages
func (p Pagination) GetTotalPages() int { return p.totalPages }

// MarshalJSON implements json.Marshaler
func (p Pagination) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Count       int `json:"count"`
		Total       int `json:"total"`
		CurrentPage int `json:"current_page"`
		PerPage     int `json:"per_page"`
		TotalPages  int `json:"total_pages"`
	}
	return json.Marshal(&Alias{
		Count:       p.count,
		Total:       p.total,
		CurrentPage: p.currentPage,
		PerPage:     p.perPage,
		TotalPages:  p.totalPages,
	})
}

// Spent represents spent amount with immutable fields
type Spent struct {
	sum          string
	currencyCode string
}

// GetSum returns the sum
func (s Spent) GetSum() string { return s.sum }

// GetCurrencyCode returns the currency code
func (s Spent) GetCurrencyCode() string { return s.currencyCode }

// MarshalJSON implements json.Marshaler
func (s Spent) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Sum          string `json:"sum"`
		CurrencyCode string `json:"currency_code"`
	}
	return json.Marshal(&Alias{
		Sum:          s.sum,
		CurrencyCode: s.currencyCode,
	})
}

// Budget represents a budget with immutable fields
type Budget struct {
	id     string
	active bool
	name   string
	notes  interface{}
	spent  Spent
}

// GetId returns the budget ID
func (b Budget) GetId() string { return b.id }

// GetActive returns the active status
func (b Budget) GetActive() bool { return b.active }

// GetName returns the budget name
func (b Budget) GetName() string { return b.name }

// GetNotes returns the budget notes
func (b Budget) GetNotes() interface{} { return b.notes }

// GetSpent returns the spent amount
func (b Budget) GetSpent() Spent { return b.spent }

// MarshalJSON implements json.Marshaler
func (b Budget) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id     string      `json:"id"`
		Active bool        `json:"active"`
		Name   string      `json:"name"`
		Notes  interface{} `json:"notes"`
		Spent  Spent       `json:"spent"`
	}
	return json.Marshal(&Alias{
		Id:     b.id,
		Active: b.active,
		Name:   b.name,
		Notes:  b.notes,
		Spent:  b.spent,
	})
}

// BudgetList represents a list of budgets with immutable fields
type BudgetList struct {
	data       []Budget
	pagination Pagination
}

// GetData returns the budget data
func (l BudgetList) GetData() []Budget {
	// Return a copy to maintain immutability
	result := make([]Budget, len(l.data))
	copy(result, l.data)
	return result
}

// GetPagination returns the pagination info
func (l BudgetList) GetPagination() Pagination { return l.pagination }

// MarshalJSON implements json.Marshaler
func (l BudgetList) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Data       []Budget   `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
	return json.Marshal(&Alias{
		Data:       l.data,
		Pagination: l.pagination,
	})
}

// Category represents a category with immutable fields
type Category struct {
	id    string
	name  string
	notes interface{}
}

// GetId returns the category ID
func (c Category) GetId() string { return c.id }

// GetName returns the category name
func (c Category) GetName() string { return c.name }

// GetNotes returns the category notes
func (c Category) GetNotes() interface{} { return c.notes }

// MarshalJSON implements json.Marshaler
func (c Category) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id    string      `json:"id"`
		Name  string      `json:"name"`
		Notes interface{} `json:"notes"`
	}
	return json.Marshal(&Alias{
		Id:    c.id,
		Name:  c.name,
		Notes: c.notes,
	})
}

// CategoryList represents a list of categories with immutable fields
type CategoryList struct {
	data       []Category
	pagination Pagination
}

// GetData returns the category data
func (l CategoryList) GetData() []Category {
	// Return a copy to maintain immutability
	result := make([]Category, len(l.data))
	copy(result, l.data)
	return result
}

// GetPagination returns the pagination info
func (l CategoryList) GetPagination() Pagination { return l.pagination }

// MarshalJSON implements json.Marshaler
func (l CategoryList) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Data       []Category `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
	return json.Marshal(&Alias{
		Data:       l.data,
		Pagination: l.pagination,
	})
}

// Account represents an account with immutable fields
type Account struct {
	id     string
	active bool
	name   string
	notes  *string
	typ    string
}

// GetId returns the account ID
func (a Account) GetId() string { return a.id }

// GetActive returns the active status
func (a Account) GetActive() bool { return a.active }

// GetName returns the account name
func (a Account) GetName() string { return a.name }

// GetNotes returns the account notes
func (a Account) GetNotes() *string { return a.notes }

// GetType returns the account type
func (a Account) GetType() string { return a.typ }

// MarshalJSON implements json.Marshaler
func (a Account) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id     string  `json:"id"`
		Active bool    `json:"active"`
		Name   string  `json:"name"`
		Notes  *string `json:"notes"`
		Type   string  `json:"type"`
	}
	return json.Marshal(&Alias{
		Id:     a.id,
		Active: a.active,
		Name:   a.name,
		Notes:  a.notes,
		Type:   a.typ,
	})
}

// AccountList represents a list of accounts with immutable fields
type AccountList struct {
	data       []Account
	pagination Pagination
}

// GetData returns the account data
func (l AccountList) GetData() []Account {
	// Return a copy to maintain immutability
	result := make([]Account, len(l.data))
	copy(result, l.data)
	return result
}

// GetPagination returns the pagination info
func (l AccountList) GetPagination() Pagination { return l.pagination }

// MarshalJSON implements json.Marshaler
func (l AccountList) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Data       []Account  `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
	return json.Marshal(&Alias{
		Data:       l.data,
		Pagination: l.pagination,
	})
}

// Transaction represents a transaction with immutable fields
type Transaction struct {
	id              string
	amount          string
	billId          interface{}
	billName        interface{}
	budgetId        *string
	budgetName      *string
	categoryId      *string
	categoryName    *string
	currencyCode    string
	date            time.Time
	description     string
	destinationId   string
	destinationName string
	destinationType string
	notes           *string
	reconciled      bool
	sourceId        string
	sourceName      string
	tags            []string
	typ             string
}

// GetId returns the transaction ID
func (t Transaction) GetId() string { return t.id }

// GetAmount returns the transaction amount
func (t Transaction) GetAmount() string { return t.amount }

// GetBillId returns the bill ID
func (t Transaction) GetBillId() interface{} { return t.billId }

// GetBillName returns the bill name
func (t Transaction) GetBillName() interface{} { return t.billName }

// GetBudgetId returns the budget ID
func (t Transaction) GetBudgetId() *string { return t.budgetId }

// GetBudgetName returns the budget name
func (t Transaction) GetBudgetName() *string { return t.budgetName }

// GetCategoryId returns the category ID
func (t Transaction) GetCategoryId() *string { return t.categoryId }

// GetCategoryName returns the category name
func (t Transaction) GetCategoryName() *string { return t.categoryName }

// GetCurrencyCode returns the currency code
func (t Transaction) GetCurrencyCode() string { return t.currencyCode }

// GetDate returns the transaction date
func (t Transaction) GetDate() time.Time { return t.date }

// GetDescription returns the transaction description
func (t Transaction) GetDescription() string { return t.description }

// GetDestinationId returns the destination ID
func (t Transaction) GetDestinationId() string { return t.destinationId }

// GetDestinationName returns the destination name
func (t Transaction) GetDestinationName() string { return t.destinationName }

// GetDestinationType returns the destination type
func (t Transaction) GetDestinationType() string { return t.destinationType }

// GetNotes returns the transaction notes
func (t Transaction) GetNotes() *string { return t.notes }

// GetReconciled returns the reconciled status
func (t Transaction) GetReconciled() bool { return t.reconciled }

// GetSourceId returns the source ID
func (t Transaction) GetSourceId() string { return t.sourceId }

// GetSourceName returns the source name
func (t Transaction) GetSourceName() string { return t.sourceName }

// GetTags returns the transaction tags
func (t Transaction) GetTags() []string {
	// Return a copy to maintain immutability
	result := make([]string, len(t.tags))
	copy(result, t.tags)
	return result
}

// GetType returns the transaction type
func (t Transaction) GetType() string { return t.typ }

// MarshalJSON implements json.Marshaler
func (t Transaction) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id              string      `json:"id"`
		Amount          string      `json:"amount"`
		BillId          interface{} `json:"bill_id"`
		BillName        interface{} `json:"bill_name"`
		BudgetId        *string     `json:"budget_id"`
		BudgetName      *string     `json:"budget_name"`
		CategoryId      *string     `json:"category_id"`
		CategoryName    *string     `json:"category_name"`
		CurrencyCode    string      `json:"currency_code"`
		Date            time.Time   `json:"date"`
		Description     string      `json:"description"`
		DestinationId   string      `json:"destination_id"`
		DestinationName string      `json:"destination_name"`
		DestinationType string      `json:"destination_type"`
		Notes           *string     `json:"notes"`
		Reconciled      bool        `json:"reconciled"`
		SourceId        string      `json:"source_id"`
		SourceName      string      `json:"source_name"`
		Tags            []string    `json:"tags"`
		Type            string      `json:"type"`
	}
	return json.Marshal(&Alias{
		Id:              t.id,
		Amount:          t.amount,
		BillId:          t.billId,
		BillName:        t.billName,
		BudgetId:        t.budgetId,
		BudgetName:      t.budgetName,
		CategoryId:      t.categoryId,
		CategoryName:    t.categoryName,
		CurrencyCode:    t.currencyCode,
		Date:            t.date,
		Description:     t.description,
		DestinationId:   t.destinationId,
		DestinationName: t.destinationName,
		DestinationType: t.destinationType,
		Notes:           t.notes,
		Reconciled:      t.reconciled,
		SourceId:        t.sourceId,
		SourceName:      t.sourceName,
		Tags:            t.tags,
		Type:            t.typ,
	})
}

// TransactionGroup represents a transaction group with immutable fields
type TransactionGroup struct {
	id           string
	groupTitle   string
	transactions []Transaction
}

// GetId returns the transaction group ID
func (g TransactionGroup) GetId() string { return g.id }

// GetGroupTitle returns the group title
func (g TransactionGroup) GetGroupTitle() string { return g.groupTitle }

// GetTransactions returns the transactions
func (g TransactionGroup) GetTransactions() []Transaction {
	// Return a copy to maintain immutability
	result := make([]Transaction, len(g.transactions))
	copy(result, g.transactions)
	return result
}

// MarshalJSON implements json.Marshaler
func (g TransactionGroup) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id           string        `json:"id"`
		GroupTitle   string        `json:"group_title"`
		Transactions []Transaction `json:"transactions"`
	}
	return json.Marshal(&Alias{
		Id:           g.id,
		GroupTitle:   g.groupTitle,
		Transactions: g.transactions,
	})
}

// TransactionList represents a list of transaction groups with immutable fields
type TransactionList struct {
	data       []TransactionGroup
	pagination Pagination
}

// GetData returns the transaction group data
func (l TransactionList) GetData() []TransactionGroup {
	// Return a copy to maintain immutability
	result := make([]TransactionGroup, len(l.data))
	copy(result, l.data)
	return result
}

// GetPagination returns the pagination info
func (l TransactionList) GetPagination() Pagination { return l.pagination }

// MarshalJSON implements json.Marshaler
func (l TransactionList) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Data       []TransactionGroup `json:"data"`
		Pagination Pagination         `json:"pagination"`
	}
	return json.Marshal(&Alias{
		Data:       l.data,
		Pagination: l.pagination,
	})
}

// Additional helper functions for backward compatibility

// NewPaginationFromValues creates a Pagination from individual values
func NewPaginationFromValues(count, total, currentPage, perPage, totalPages int) Pagination {
	return NewPaginationBuilder().
		WithCount(count).
		WithTotal(total).
		WithCurrentPage(currentPage).
		WithPerPage(perPage).
		WithTotalPages(totalPages).
		Build()
}

// NewSpentFromValues creates a Spent from individual values
func NewSpentFromValues(sum, currencyCode string) Spent {
	return Spent{
		sum:          sum,
		currencyCode: currencyCode,
	}
}

// BasicSummary represents a basic summary with immutable fields
type BasicSummary struct {
	key           string
	title         string
	currencyCode  string
	monetaryValue string
}

// GetKey returns the key
func (s BasicSummary) GetKey() string { return s.key }

// GetTitle returns the title
func (s BasicSummary) GetTitle() string { return s.title }

// GetCurrencyCode returns the currency code
func (s BasicSummary) GetCurrencyCode() string { return s.currencyCode }

// GetMonetaryValue returns the monetary value
func (s BasicSummary) GetMonetaryValue() string { return s.monetaryValue }

// MarshalJSON implements json.Marshaler
func (s BasicSummary) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Key           string `json:"key"`
		Title         string `json:"title"`
		CurrencyCode  string `json:"currency_code"`
		MonetaryValue string `json:"monetary_value"`
	}
	return json.Marshal(&Alias{
		Key:           s.key,
		Title:         s.title,
		CurrencyCode:  s.currencyCode,
		MonetaryValue: s.monetaryValue,
	})
}

// BasicSummaryList represents a list of basic summaries
type BasicSummaryList struct {
	data []BasicSummary
}

// GetData returns the summary data
func (l BasicSummaryList) GetData() []BasicSummary {
	// Return a copy to maintain immutability
	result := make([]BasicSummary, len(l.data))
	copy(result, l.data)
	return result
}

// MarshalJSON implements json.Marshaler
func (l BasicSummaryList) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Data []BasicSummary `json:"data"`
	}
	return json.Marshal(&Alias{
		Data: l.data,
	})
}

// InsightCategoryEntry represents an insight category entry with immutable fields
type InsightCategoryEntry struct {
	id           string
	name         string
	amount       string
	currencyCode string
}

// GetId returns the ID
func (e InsightCategoryEntry) GetId() string { return e.id }

// GetName returns the name
func (e InsightCategoryEntry) GetName() string { return e.name }

// GetAmount returns the amount
func (e InsightCategoryEntry) GetAmount() string { return e.amount }

// GetCurrencyCode returns the currency code
func (e InsightCategoryEntry) GetCurrencyCode() string { return e.currencyCode }

// MarshalJSON implements json.Marshaler
func (e InsightCategoryEntry) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	}
	return json.Marshal(&Alias{
		Id:           e.id,
		Name:         e.name,
		Amount:       e.amount,
		CurrencyCode: e.currencyCode,
	})
}

// InsightTotalEntry represents an insight total entry with immutable fields
type InsightTotalEntry struct {
	amount       string
	currencyCode string
}

// GetAmount returns the amount
func (e InsightTotalEntry) GetAmount() string { return e.amount }

// GetCurrencyCode returns the currency code
func (e InsightTotalEntry) GetCurrencyCode() string { return e.currencyCode }

// MarshalJSON implements json.Marshaler
func (e InsightTotalEntry) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	}
	return json.Marshal(&Alias{
		Amount:       e.amount,
		CurrencyCode: e.currencyCode,
	})
}

// InsightCategoryResponse represents an insight category response
type InsightCategoryResponse struct {
	entries []InsightCategoryEntry
}

// GetEntries returns the entries
func (r InsightCategoryResponse) GetEntries() []InsightCategoryEntry {
	// Return a copy to maintain immutability
	result := make([]InsightCategoryEntry, len(r.entries))
	copy(result, r.entries)
	return result
}

// MarshalJSON implements json.Marshaler
func (r InsightCategoryResponse) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Entries []InsightCategoryEntry `json:"entries"`
	}
	return json.Marshal(&Alias{
		Entries: r.entries,
	})
}

// InsightTotalResponse represents an insight total response
type InsightTotalResponse struct {
	entries []InsightTotalEntry
}

// GetEntries returns the entries
func (r InsightTotalResponse) GetEntries() []InsightTotalEntry {
	// Return a copy to maintain immutability
	result := make([]InsightTotalEntry, len(r.entries))
	copy(result, r.entries)
	return result
}

// MarshalJSON implements json.Marshaler
func (r InsightTotalResponse) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Entries []InsightTotalEntry `json:"entries"`
	}
	return json.Marshal(&Alias{
		Entries: r.entries,
	})
}

// BudgetSpent represents budget spent amount
type BudgetSpent struct {
	Sum            string `json:"sum"`
	CurrencyCode   string `json:"currency_code"`
	CurrencySymbol string `json:"currency_symbol"`
}

// BudgetLimit represents a budget limit
type BudgetLimit struct {
	Id             string        `json:"id"`
	Amount         string        `json:"amount"`
	Start          time.Time     `json:"start"`
	End            time.Time     `json:"end"`
	BudgetId       string        `json:"budget_id"`
	CurrencyCode   string        `json:"currency_code"`
	CurrencySymbol string        `json:"currency_symbol"`
	Spent          []BudgetSpent `json:"spent"`
}

// BudgetLimitList represents a list of budget limits
type BudgetLimitList struct {
	Data       []BudgetLimit `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

// Tag represents a tag with immutable fields
type Tag struct {
	id          string
	tag         string
	description *string
}

// GetId returns the tag ID
func (t Tag) GetId() string { return t.id }

// GetTag returns the tag name
func (t Tag) GetTag() string { return t.tag }

// GetDescription returns the tag description
func (t Tag) GetDescription() *string { return t.description }

// MarshalJSON implements json.Marshaler
func (t Tag) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id          string  `json:"id"`
		Tag         string  `json:"tag"`
		Description *string `json:"description"`
	}
	return json.Marshal(&Alias{
		Id:          t.id,
		Tag:         t.tag,
		Description: t.description,
	})
}

// TagList represents a list of tags
type TagList struct {
	Data       []Tag      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// PaidDate represents a paid date
type PaidDate struct {
	date                 *time.Time
	transactionGroupId   *string
	transactionJournalId *string
}

// GetDate returns the date
func (p PaidDate) GetDate() *time.Time {
	return p.date
}

// GetTransactionGroupId returns the transaction group ID
func (p PaidDate) GetTransactionGroupId() *string {
	return p.transactionGroupId
}

// GetTransactionJournalId returns the transaction journal ID
func (p PaidDate) GetTransactionJournalId() *string {
	return p.transactionJournalId
}

// MarshalJSON implements json.Marshaler
func (p PaidDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Date                 *time.Time `json:"date"`
		TransactionGroupId   *string    `json:"transaction_group_id"`
		TransactionJournalId *string    `json:"transaction_journal_id"`
	}{
		Date:                 p.date,
		TransactionGroupId:   p.transactionGroupId,
		TransactionJournalId: p.transactionJournalId,
	})
}

// Bill represents a bill with immutable fields
type Bill struct {
	id                string
	active            bool
	name              string
	amountMin         string
	amountMax         string
	date              time.Time
	repeatFreq        string
	skip              int
	currencyCode      string
	notes             *string
	nextExpectedMatch *time.Time
	paidDates         []PaidDate
}

// GetId returns the bill ID
func (b Bill) GetId() string { return b.id }

// GetActive returns the active status
func (b Bill) GetActive() bool { return b.active }

// GetName returns the bill name
func (b Bill) GetName() string { return b.name }

// GetAmountMin returns the minimum amount
func (b Bill) GetAmountMin() string { return b.amountMin }

// GetAmountMax returns the maximum amount
func (b Bill) GetAmountMax() string { return b.amountMax }

// GetDate returns the bill date
func (b Bill) GetDate() time.Time { return b.date }

// GetRepeatFreq returns the repeat frequency
func (b Bill) GetRepeatFreq() string { return b.repeatFreq }

// GetSkip returns the skip value
func (b Bill) GetSkip() int { return b.skip }

// GetCurrencyCode returns the currency code
func (b Bill) GetCurrencyCode() string { return b.currencyCode }

// GetNotes returns the bill notes
func (b Bill) GetNotes() *string { return b.notes }

// GetNextExpectedMatch returns the next expected match date
func (b Bill) GetNextExpectedMatch() *time.Time { return b.nextExpectedMatch }

// GetPaidDates returns the paid dates
func (b Bill) GetPaidDates() []PaidDate {
	// Return a copy to maintain immutability
	result := make([]PaidDate, len(b.paidDates))
	copy(result, b.paidDates)
	return result
}

// MarshalJSON implements json.Marshaler
func (b Bill) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id                string     `json:"id"`
		Active            bool       `json:"active"`
		Name              string     `json:"name"`
		AmountMin         string     `json:"amount_min"`
		AmountMax         string     `json:"amount_max"`
		Date              time.Time  `json:"date"`
		RepeatFreq        string     `json:"repeat_freq"`
		Skip              int        `json:"skip"`
		CurrencyCode      string     `json:"currency_code"`
		Notes             *string    `json:"notes"`
		NextExpectedMatch *time.Time `json:"next_expected_match"`
		PaidDates         []PaidDate `json:"paid_dates"`
	}
	return json.Marshal(&Alias{
		Id:                b.id,
		Active:            b.active,
		Name:              b.name,
		AmountMin:         b.amountMin,
		AmountMax:         b.amountMax,
		Date:              b.date,
		RepeatFreq:        b.repeatFreq,
		Skip:              b.skip,
		CurrencyCode:      b.currencyCode,
		Notes:             b.notes,
		NextExpectedMatch: b.nextExpectedMatch,
		PaidDates:         b.paidDates,
	})
}

// BillList represents a list of bills
type BillList struct {
	Data       []Bill     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// RecurrenceRepetition represents a recurrence repetition
type RecurrenceRepetition struct {
	Id          string  `json:"id"`
	Type        string  `json:"type"`
	Moment      string  `json:"moment"`
	Skip        int     `json:"skip"`
	Weekend     int     `json:"weekend"`
	Description *string `json:"description"`
}

// RecurrenceTransaction represents a recurrence transaction
type RecurrenceTransaction struct {
	Id              string  `json:"id"`
	Description     string  `json:"description"`
	Amount          string  `json:"amount"`
	CurrencyCode    string  `json:"currency_code"`
	CategoryId      *string `json:"category_id"`
	CategoryName    *string `json:"category_name"`
	BudgetId        *string `json:"budget_id"`
	BudgetName      *string `json:"budget_name"`
	SourceId        string  `json:"source_id"`
	SourceName      string  `json:"source_name"`
	DestinationId   string  `json:"destination_id"`
	DestinationName string  `json:"destination_name"`
}

// Recurrence represents a recurrence with immutable fields
type Recurrence struct {
	id              string
	typ             string
	title           string
	description     string
	firstDate       time.Time
	latestDate      *time.Time
	repeatUntil     *time.Time
	nrOfRepetitions *int
	applyRules      bool
	active          bool
	notes           *string
	repetitions     []RecurrenceRepetition
	transactions    []RecurrenceTransaction
}

// GetId returns the recurrence ID
func (r Recurrence) GetId() string { return r.id }

// GetType returns the recurrence type
func (r Recurrence) GetType() string { return r.typ }

// GetTitle returns the recurrence title
func (r Recurrence) GetTitle() string { return r.title }

// GetDescription returns the recurrence description
func (r Recurrence) GetDescription() string { return r.description }

// GetFirstDate returns the first date
func (r Recurrence) GetFirstDate() time.Time { return r.firstDate }

// GetLatestDate returns the latest date
func (r Recurrence) GetLatestDate() *time.Time { return r.latestDate }

// GetRepeatUntil returns the repeat until date
func (r Recurrence) GetRepeatUntil() *time.Time { return r.repeatUntil }

// GetNrOfRepetitions returns the number of repetitions
func (r Recurrence) GetNrOfRepetitions() *int { return r.nrOfRepetitions }

// GetApplyRules returns the apply rules flag
func (r Recurrence) GetApplyRules() bool { return r.applyRules }

// GetActive returns the active status
func (r Recurrence) GetActive() bool { return r.active }

// GetNotes returns the recurrence notes
func (r Recurrence) GetNotes() *string { return r.notes }

// GetRepetitions returns the repetitions
func (r Recurrence) GetRepetitions() []RecurrenceRepetition {
	// Return a copy to maintain immutability
	result := make([]RecurrenceRepetition, len(r.repetitions))
	copy(result, r.repetitions)
	return result
}

// GetTransactions returns the transactions
func (r Recurrence) GetTransactions() []RecurrenceTransaction {
	// Return a copy to maintain immutability
	result := make([]RecurrenceTransaction, len(r.transactions))
	copy(result, r.transactions)
	return result
}

// MarshalJSON implements json.Marshaler
func (r Recurrence) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Id              string                  `json:"id"`
		Type            string                  `json:"type"`
		Title           string                  `json:"title"`
		Description     string                  `json:"description"`
		FirstDate       time.Time               `json:"first_date"`
		LatestDate      *time.Time              `json:"latest_date"`
		RepeatUntil     *time.Time              `json:"repeat_until"`
		NrOfRepetitions *int                    `json:"nr_of_repetitions"`
		ApplyRules      bool                    `json:"apply_rules"`
		Active          bool                    `json:"active"`
		Notes           *string                 `json:"notes"`
		Repetitions     []RecurrenceRepetition  `json:"repetitions"`
		Transactions    []RecurrenceTransaction `json:"transactions"`
	}
	return json.Marshal(&Alias{
		Id:              r.id,
		Type:            r.typ,
		Title:           r.title,
		Description:     r.description,
		FirstDate:       r.firstDate,
		LatestDate:      r.latestDate,
		RepeatUntil:     r.repeatUntil,
		NrOfRepetitions: r.nrOfRepetitions,
		ApplyRules:      r.applyRules,
		Active:          r.active,
		Notes:           r.notes,
		Repetitions:     r.repetitions,
		Transactions:    r.transactions,
	})
}

// RecurrenceList represents a list of recurrences
type RecurrenceList struct {
	Data       []Recurrence `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

// TransactionStoreRequest represents the request body for creating a new transaction
// Note: This is a request DTO and doesn't need to be immutable
type TransactionStoreRequest struct {
	ErrorIfDuplicateHash bool                      `json:"error_if_duplicate_hash" mcp:"Break if transaction with same hash already exists (default: false)"`
	ApplyRules           bool                      `json:"apply_rules" mcp:"Whether to apply processing rules when creating transaction (default: false)"`
	FireWebhooks         bool                      `json:"fire_webhooks" mcp:"Whether to fire webhooks for this transaction (default: true)"`
	GroupTitle           string                    `json:"group_title" mcp:"Title for the transaction group (for split transactions)"`
	Transactions         []TransactionSplitRequest `json:"transactions" mcp:"Array of transactions to create (required, at least one)"`
}

// TransactionSplitRequest represents a single transaction in a transaction group
// Note: This is a request DTO and doesn't need to be immutable
type TransactionSplitRequest struct {
	Type                string   `json:"type" mcp:"Transaction type: withdrawal, deposit, transfer (required)"`
	Date                string   `json:"date" mcp:"Transaction date (YYYY-MM-DD or YYYY-MM-DD HH:MM:SS) (required)"`
	Amount              string   `json:"amount" mcp:"Transaction amount as string (e.g. '100.00') (required)"`
	Description         string   `json:"description" mcp:"Transaction description (required)"`
	SourceId            *string  `json:"source_id,omitempty" mcp:"Source account ID (use either source_id or source_name)"`
	SourceName          *string  `json:"source_name,omitempty" mcp:"Source account name (use either source_id or source_name)"`
	DestinationId       *string  `json:"destination_id,omitempty" mcp:"Destination account ID (use either destination_id or destination_name)"`
	DestinationName     *string  `json:"destination_name,omitempty" mcp:"Destination account name (use either destination_id or destination_name)"`
	CategoryId          *string  `json:"category_id,omitempty" mcp:"Category ID (use either category_id or category_name)"`
	CategoryName        *string  `json:"category_name,omitempty" mcp:"Category name (use either category_id or category_name)"`
	BudgetId            *string  `json:"budget_id,omitempty" mcp:"Budget ID (use either budget_id or budget_name)"`
	BudgetName          *string  `json:"budget_name,omitempty" mcp:"Budget name (use either budget_id or budget_name)"`
	Tags                []string `json:"tags,omitempty" mcp:"Array of tag names to attach to transaction"`
	CurrencyId          *string  `json:"currency_id,omitempty" mcp:"Currency ID for the transaction"`
	CurrencyCode        *string  `json:"currency_code,omitempty" mcp:"Currency code (e.g. 'USD', 'EUR')"`
	ForeignAmount       *string  `json:"foreign_amount,omitempty" mcp:"Amount in foreign currency as string"`
	ForeignCurrencyId   *string  `json:"foreign_currency_id,omitempty" mcp:"Foreign currency ID"`
	ForeignCurrencyCode *string  `json:"foreign_currency_code,omitempty" mcp:"Foreign currency code (e.g. 'USD', 'EUR')"`
	BillId              *string  `json:"bill_id,omitempty" mcp:"Bill ID to link this transaction to"`
	BillName            *string  `json:"bill_name,omitempty" mcp:"Bill name to link this transaction to"`
	PiggyBankId         *string  `json:"piggy_bank_id,omitempty" mcp:"Piggy bank ID for savings transfers"`
	PiggyBankName       *string  `json:"piggy_bank_name,omitempty" mcp:"Piggy bank name for savings transfers"`
	Notes               *string  `json:"notes,omitempty" mcp:"Additional notes or comments for the transaction"`
	Reconciled          *bool    `json:"reconciled,omitempty" mcp:"Whether the transaction has been reconciled (default: false)"`
	Order               *int     `json:"order,omitempty" mcp:"Order of this split in the transaction group"`
}