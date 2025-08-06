package dto

import (
	"time"
)

// AccountBuilder builds an immutable Account
type AccountBuilder struct {
	id     string
	active bool
	name   string
	notes  *string
	typ    string
}

// NewAccountBuilder creates a new AccountBuilder
func NewAccountBuilder() *AccountBuilder {
	return &AccountBuilder{}
}

// WithId sets the account ID
func (b *AccountBuilder) WithId(id string) *AccountBuilder {
	b.id = id
	return b
}

// WithActive sets the active status
func (b *AccountBuilder) WithActive(active bool) *AccountBuilder {
	b.active = active
	return b
}

// WithName sets the account name
func (b *AccountBuilder) WithName(name string) *AccountBuilder {
	b.name = name
	return b
}

// WithNotes sets the account notes
func (b *AccountBuilder) WithNotes(notes *string) *AccountBuilder {
	b.notes = notes
	return b
}

// WithType sets the account type
func (b *AccountBuilder) WithType(typ string) *AccountBuilder {
	b.typ = typ
	return b
}

// Build creates an immutable Account
func (b *AccountBuilder) Build() Account {
	return Account{
		id:     b.id,
		active: b.active,
		name:   b.name,
		notes:  b.notes,
		typ:    b.typ,
	}
}

// BudgetBuilder builds an immutable Budget
type BudgetBuilder struct {
	id     string
	active bool
	name   string
	notes  interface{}
	spent  Spent
}

// NewBudgetBuilder creates a new BudgetBuilder
func NewBudgetBuilder() *BudgetBuilder {
	return &BudgetBuilder{}
}

// WithId sets the budget ID
func (b *BudgetBuilder) WithId(id string) *BudgetBuilder {
	b.id = id
	return b
}

// WithActive sets the active status
func (b *BudgetBuilder) WithActive(active bool) *BudgetBuilder {
	b.active = active
	return b
}

// WithName sets the budget name
func (b *BudgetBuilder) WithName(name string) *BudgetBuilder {
	b.name = name
	return b
}

// WithNotes sets the budget notes
func (b *BudgetBuilder) WithNotes(notes interface{}) *BudgetBuilder {
	b.notes = notes
	return b
}

// WithSpent sets the budget spent amount
func (b *BudgetBuilder) WithSpent(spent Spent) *BudgetBuilder {
	b.spent = spent
	return b
}

// Build creates an immutable Budget
func (b *BudgetBuilder) Build() Budget {
	return Budget{
		id:     b.id,
		active: b.active,
		name:   b.name,
		notes:  b.notes,
		spent:  b.spent,
	}
}

// CategoryBuilder builds an immutable Category
type CategoryBuilder struct {
	id    string
	name  string
	notes interface{}
}

// NewCategoryBuilder creates a new CategoryBuilder
func NewCategoryBuilder() *CategoryBuilder {
	return &CategoryBuilder{}
}

// WithId sets the category ID
func (b *CategoryBuilder) WithId(id string) *CategoryBuilder {
	b.id = id
	return b
}

// WithName sets the category name
func (b *CategoryBuilder) WithName(name string) *CategoryBuilder {
	b.name = name
	return b
}

// WithNotes sets the category notes
func (b *CategoryBuilder) WithNotes(notes interface{}) *CategoryBuilder {
	b.notes = notes
	return b
}

// Build creates an immutable Category
func (b *CategoryBuilder) Build() Category {
	return Category{
		id:    b.id,
		name:  b.name,
		notes: b.notes,
	}
}

// TransactionBuilder builds an immutable Transaction
type TransactionBuilder struct {
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

// NewTransactionBuilder creates a new TransactionBuilder
func NewTransactionBuilder() *TransactionBuilder {
	return &TransactionBuilder{
		tags: []string{},
	}
}

// WithId sets the transaction ID
func (b *TransactionBuilder) WithId(id string) *TransactionBuilder {
	b.id = id
	return b
}

// WithAmount sets the transaction amount
func (b *TransactionBuilder) WithAmount(amount string) *TransactionBuilder {
	b.amount = amount
	return b
}

// WithBillId sets the bill ID
func (b *TransactionBuilder) WithBillId(billId interface{}) *TransactionBuilder {
	b.billId = billId
	return b
}

// WithBillName sets the bill name
func (b *TransactionBuilder) WithBillName(billName interface{}) *TransactionBuilder {
	b.billName = billName
	return b
}

// WithBudgetId sets the budget ID
func (b *TransactionBuilder) WithBudgetId(budgetId *string) *TransactionBuilder {
	b.budgetId = budgetId
	return b
}

// WithBudgetName sets the budget name
func (b *TransactionBuilder) WithBudgetName(budgetName *string) *TransactionBuilder {
	b.budgetName = budgetName
	return b
}

// WithCategoryId sets the category ID
func (b *TransactionBuilder) WithCategoryId(categoryId *string) *TransactionBuilder {
	b.categoryId = categoryId
	return b
}

// WithCategoryName sets the category name
func (b *TransactionBuilder) WithCategoryName(categoryName *string) *TransactionBuilder {
	b.categoryName = categoryName
	return b
}

// WithCurrencyCode sets the currency code
func (b *TransactionBuilder) WithCurrencyCode(currencyCode string) *TransactionBuilder {
	b.currencyCode = currencyCode
	return b
}

// WithDate sets the transaction date
func (b *TransactionBuilder) WithDate(date time.Time) *TransactionBuilder {
	b.date = date
	return b
}

// WithDescription sets the transaction description
func (b *TransactionBuilder) WithDescription(description string) *TransactionBuilder {
	b.description = description
	return b
}

// WithDestinationId sets the destination ID
func (b *TransactionBuilder) WithDestinationId(destinationId string) *TransactionBuilder {
	b.destinationId = destinationId
	return b
}

// WithDestinationName sets the destination name
func (b *TransactionBuilder) WithDestinationName(destinationName string) *TransactionBuilder {
	b.destinationName = destinationName
	return b
}

// WithDestinationType sets the destination type
func (b *TransactionBuilder) WithDestinationType(destinationType string) *TransactionBuilder {
	b.destinationType = destinationType
	return b
}

// WithNotes sets the transaction notes
func (b *TransactionBuilder) WithNotes(notes *string) *TransactionBuilder {
	b.notes = notes
	return b
}

// WithReconciled sets the reconciled status
func (b *TransactionBuilder) WithReconciled(reconciled bool) *TransactionBuilder {
	b.reconciled = reconciled
	return b
}

// WithSourceId sets the source ID
func (b *TransactionBuilder) WithSourceId(sourceId string) *TransactionBuilder {
	b.sourceId = sourceId
	return b
}

// WithSourceName sets the source name
func (b *TransactionBuilder) WithSourceName(sourceName string) *TransactionBuilder {
	b.sourceName = sourceName
	return b
}

// WithTags sets the transaction tags
func (b *TransactionBuilder) WithTags(tags []string) *TransactionBuilder {
	b.tags = tags
	return b
}

// WithType sets the transaction type
func (b *TransactionBuilder) WithType(typ string) *TransactionBuilder {
	b.typ = typ
	return b
}

// Build creates an immutable Transaction
func (b *TransactionBuilder) Build() Transaction {
	// Create a copy of tags to ensure immutability
	tagsCopy := make([]string, len(b.tags))
	copy(tagsCopy, b.tags)
	
	return Transaction{
		id:              b.id,
		amount:          b.amount,
		billId:          b.billId,
		billName:        b.billName,
		budgetId:        b.budgetId,
		budgetName:      b.budgetName,
		categoryId:      b.categoryId,
		categoryName:    b.categoryName,
		currencyCode:    b.currencyCode,
		date:            b.date,
		description:     b.description,
		destinationId:   b.destinationId,
		destinationName: b.destinationName,
		destinationType: b.destinationType,
		notes:           b.notes,
		reconciled:      b.reconciled,
		sourceId:        b.sourceId,
		sourceName:      b.sourceName,
		tags:            tagsCopy,
		typ:             b.typ,
	}
}

// TransactionGroupBuilder builds an immutable TransactionGroup
type TransactionGroupBuilder struct {
	id           string
	groupTitle   string
	transactions []Transaction
}

// NewTransactionGroupBuilder creates a new TransactionGroupBuilder
func NewTransactionGroupBuilder() *TransactionGroupBuilder {
	return &TransactionGroupBuilder{
		transactions: []Transaction{},
	}
}

// WithId sets the transaction group ID
func (b *TransactionGroupBuilder) WithId(id string) *TransactionGroupBuilder {
	b.id = id
	return b
}

// WithGroupTitle sets the group title
func (b *TransactionGroupBuilder) WithGroupTitle(groupTitle string) *TransactionGroupBuilder {
	b.groupTitle = groupTitle
	return b
}

// WithTransactions sets the transactions
func (b *TransactionGroupBuilder) WithTransactions(transactions []Transaction) *TransactionGroupBuilder {
	b.transactions = transactions
	return b
}

// AddTransaction adds a single transaction to the group
func (b *TransactionGroupBuilder) AddTransaction(transaction Transaction) *TransactionGroupBuilder {
	b.transactions = append(b.transactions, transaction)
	return b
}

// Build creates an immutable TransactionGroup
func (b *TransactionGroupBuilder) Build() TransactionGroup {
	// Create a copy of transactions to ensure immutability
	transactionsCopy := make([]Transaction, len(b.transactions))
	copy(transactionsCopy, b.transactions)
	
	return TransactionGroup{
		id:           b.id,
		groupTitle:   b.groupTitle,
		transactions: transactionsCopy,
	}
}

// TransactionListBuilder builds an immutable TransactionList
type TransactionListBuilder struct {
	data       []TransactionGroup
	pagination Pagination
}

// NewTransactionListBuilder creates a new TransactionListBuilder
func NewTransactionListBuilder() *TransactionListBuilder {
	return &TransactionListBuilder{
		data: []TransactionGroup{},
	}
}

// WithData sets the transaction groups
func (b *TransactionListBuilder) WithData(data []TransactionGroup) *TransactionListBuilder {
	b.data = data
	return b
}

// AddTransactionGroup adds a single transaction group
func (b *TransactionListBuilder) AddTransactionGroup(group TransactionGroup) *TransactionListBuilder {
	b.data = append(b.data, group)
	return b
}

// WithPagination sets the pagination
func (b *TransactionListBuilder) WithPagination(pagination Pagination) *TransactionListBuilder {
	b.pagination = pagination
	return b
}

// Build creates an immutable TransactionList
func (b *TransactionListBuilder) Build() TransactionList {
	// Create a copy of data to ensure immutability
	dataCopy := make([]TransactionGroup, len(b.data))
	copy(dataCopy, b.data)
	
	return TransactionList{
		data:       dataCopy,
		pagination: b.pagination,
	}
}

// PaginationBuilder builds an immutable Pagination
type PaginationBuilder struct {
	count       int
	total       int
	currentPage int
	perPage     int
	totalPages  int
}

// NewPaginationBuilder creates a new PaginationBuilder
func NewPaginationBuilder() *PaginationBuilder {
	return &PaginationBuilder{}
}

// WithCount sets the count
func (b *PaginationBuilder) WithCount(count int) *PaginationBuilder {
	b.count = count
	return b
}

// WithTotal sets the total
func (b *PaginationBuilder) WithTotal(total int) *PaginationBuilder {
	b.total = total
	return b
}

// WithCurrentPage sets the current page
func (b *PaginationBuilder) WithCurrentPage(currentPage int) *PaginationBuilder {
	b.currentPage = currentPage
	return b
}

// WithPerPage sets the per page count
func (b *PaginationBuilder) WithPerPage(perPage int) *PaginationBuilder {
	b.perPage = perPage
	return b
}

// WithTotalPages sets the total pages
func (b *PaginationBuilder) WithTotalPages(totalPages int) *PaginationBuilder {
	b.totalPages = totalPages
	return b
}

// Build creates an immutable Pagination
func (b *PaginationBuilder) Build() Pagination {
	return Pagination{
		count:       b.count,
		total:       b.total,
		currentPage: b.currentPage,
		perPage:     b.perPage,
		totalPages:  b.totalPages,
	}
}

// TagBuilder builds an immutable Tag
type TagBuilder struct {
	id          string
	tag         string
	description *string
}

// NewTagBuilder creates a new TagBuilder
func NewTagBuilder() *TagBuilder {
	return &TagBuilder{}
}

// WithId sets the tag ID
func (b *TagBuilder) WithId(id string) *TagBuilder {
	b.id = id
	return b
}

// WithTag sets the tag name
func (b *TagBuilder) WithTag(tag string) *TagBuilder {
	b.tag = tag
	return b
}

// WithDescription sets the tag description
func (b *TagBuilder) WithDescription(description *string) *TagBuilder {
	b.description = description
	return b
}

// Build creates an immutable Tag
func (b *TagBuilder) Build() Tag {
	return Tag{
		id:          b.id,
		tag:         b.tag,
		description: b.description,
	}
}

// BillBuilder builds an immutable Bill
type BillBuilder struct {
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

// NewBillBuilder creates a new BillBuilder
func NewBillBuilder() *BillBuilder {
	return &BillBuilder{
		paidDates: []PaidDate{},
	}
}

// WithId sets the bill ID
func (b *BillBuilder) WithId(id string) *BillBuilder {
	b.id = id
	return b
}

// WithActive sets the active status
func (b *BillBuilder) WithActive(active bool) *BillBuilder {
	b.active = active
	return b
}

// WithName sets the bill name
func (b *BillBuilder) WithName(name string) *BillBuilder {
	b.name = name
	return b
}

// WithAmountMin sets the minimum amount
func (b *BillBuilder) WithAmountMin(amountMin string) *BillBuilder {
	b.amountMin = amountMin
	return b
}

// WithAmountMax sets the maximum amount
func (b *BillBuilder) WithAmountMax(amountMax string) *BillBuilder {
	b.amountMax = amountMax
	return b
}

// WithDate sets the bill date
func (b *BillBuilder) WithDate(date time.Time) *BillBuilder {
	b.date = date
	return b
}

// WithRepeatFreq sets the repeat frequency
func (b *BillBuilder) WithRepeatFreq(repeatFreq string) *BillBuilder {
	b.repeatFreq = repeatFreq
	return b
}

// WithSkip sets the skip value
func (b *BillBuilder) WithSkip(skip int) *BillBuilder {
	b.skip = skip
	return b
}

// WithCurrencyCode sets the currency code
func (b *BillBuilder) WithCurrencyCode(currencyCode string) *BillBuilder {
	b.currencyCode = currencyCode
	return b
}

// WithNotes sets the bill notes
func (b *BillBuilder) WithNotes(notes *string) *BillBuilder {
	b.notes = notes
	return b
}

// WithNextExpectedMatch sets the next expected match date
func (b *BillBuilder) WithNextExpectedMatch(nextExpectedMatch *time.Time) *BillBuilder {
	b.nextExpectedMatch = nextExpectedMatch
	return b
}

// WithPaidDates sets the paid dates
func (b *BillBuilder) WithPaidDates(paidDates []PaidDate) *BillBuilder {
	b.paidDates = paidDates
	return b
}

// AddPaidDate adds a single paid date
func (b *BillBuilder) AddPaidDate(paidDate PaidDate) *BillBuilder {
	b.paidDates = append(b.paidDates, paidDate)
	return b
}

// Build creates an immutable Bill
func (b *BillBuilder) Build() Bill {
	// Create a copy of paid dates to ensure immutability
	paidDatesCopy := make([]PaidDate, len(b.paidDates))
	copy(paidDatesCopy, b.paidDates)
	
	return Bill{
		id:                b.id,
		active:            b.active,
		name:              b.name,
		amountMin:         b.amountMin,
		amountMax:         b.amountMax,
		date:              b.date,
		repeatFreq:        b.repeatFreq,
		skip:              b.skip,
		currencyCode:      b.currencyCode,
		notes:             b.notes,
		nextExpectedMatch: b.nextExpectedMatch,
		paidDates:         paidDatesCopy,
	}
}

// RecurrenceBuilder builds an immutable Recurrence
type RecurrenceBuilder struct {
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

// NewRecurrenceBuilder creates a new RecurrenceBuilder
func NewRecurrenceBuilder() *RecurrenceBuilder {
	return &RecurrenceBuilder{
		repetitions:  []RecurrenceRepetition{},
		transactions: []RecurrenceTransaction{},
	}
}

// WithId sets the recurrence ID
func (b *RecurrenceBuilder) WithId(id string) *RecurrenceBuilder {
	b.id = id
	return b
}

// WithType sets the recurrence type
func (b *RecurrenceBuilder) WithType(typ string) *RecurrenceBuilder {
	b.typ = typ
	return b
}

// WithTitle sets the recurrence title
func (b *RecurrenceBuilder) WithTitle(title string) *RecurrenceBuilder {
	b.title = title
	return b
}

// WithDescription sets the recurrence description
func (b *RecurrenceBuilder) WithDescription(description string) *RecurrenceBuilder {
	b.description = description
	return b
}

// WithFirstDate sets the first date
func (b *RecurrenceBuilder) WithFirstDate(firstDate time.Time) *RecurrenceBuilder {
	b.firstDate = firstDate
	return b
}

// WithLatestDate sets the latest date
func (b *RecurrenceBuilder) WithLatestDate(latestDate *time.Time) *RecurrenceBuilder {
	b.latestDate = latestDate
	return b
}

// WithRepeatUntil sets the repeat until date
func (b *RecurrenceBuilder) WithRepeatUntil(repeatUntil *time.Time) *RecurrenceBuilder {
	b.repeatUntil = repeatUntil
	return b
}

// WithNrOfRepetitions sets the number of repetitions
func (b *RecurrenceBuilder) WithNrOfRepetitions(nrOfRepetitions *int) *RecurrenceBuilder {
	b.nrOfRepetitions = nrOfRepetitions
	return b
}

// WithApplyRules sets the apply rules flag
func (b *RecurrenceBuilder) WithApplyRules(applyRules bool) *RecurrenceBuilder {
	b.applyRules = applyRules
	return b
}

// WithActive sets the active status
func (b *RecurrenceBuilder) WithActive(active bool) *RecurrenceBuilder {
	b.active = active
	return b
}

// WithNotes sets the recurrence notes
func (b *RecurrenceBuilder) WithNotes(notes *string) *RecurrenceBuilder {
	b.notes = notes
	return b
}

// WithRepetitions sets the repetitions
func (b *RecurrenceBuilder) WithRepetitions(repetitions []RecurrenceRepetition) *RecurrenceBuilder {
	b.repetitions = repetitions
	return b
}

// WithTransactions sets the transactions
func (b *RecurrenceBuilder) WithTransactions(transactions []RecurrenceTransaction) *RecurrenceBuilder {
	b.transactions = transactions
	return b
}

// SpentBuilder builds an immutable Spent
type SpentBuilder struct {
	sum          string
	currencyCode string
}

// NewSpentBuilder creates a new SpentBuilder
func NewSpentBuilder() *SpentBuilder {
	return &SpentBuilder{}
}

// WithSum sets the sum
func (b *SpentBuilder) WithSum(sum string) *SpentBuilder {
	b.sum = sum
	return b
}

// WithCurrencyCode sets the currency code
func (b *SpentBuilder) WithCurrencyCode(currencyCode string) *SpentBuilder {
	b.currencyCode = currencyCode
	return b
}

// Build creates an immutable Spent
func (b *SpentBuilder) Build() Spent {
	return Spent{
		sum:          b.sum,
		currencyCode: b.currencyCode,
	}
}

// Build creates an immutable Recurrence
func (b *RecurrenceBuilder) Build() Recurrence {
	// Create copies to ensure immutability
	repetitionsCopy := make([]RecurrenceRepetition, len(b.repetitions))
	copy(repetitionsCopy, b.repetitions)
	
	transactionsCopy := make([]RecurrenceTransaction, len(b.transactions))
	copy(transactionsCopy, b.transactions)
	
	return Recurrence{
		id:              b.id,
		typ:             b.typ,
		title:           b.title,
		description:     b.description,
		firstDate:       b.firstDate,
		latestDate:      b.latestDate,
		repeatUntil:     b.repeatUntil,
		nrOfRepetitions: b.nrOfRepetitions,
		applyRules:      b.applyRules,
		active:          b.active,
		notes:           b.notes,
		repetitions:     repetitionsCopy,
		transactions:    transactionsCopy,
	}
}

// PaidDateBuilder is a builder for PaidDate
type PaidDateBuilder struct {
	date                 *time.Time
	transactionGroupId   *string
	transactionJournalId *string
}

// NewPaidDateBuilder creates a new PaidDateBuilder
func NewPaidDateBuilder() *PaidDateBuilder {
	return &PaidDateBuilder{}
}

// WithDate sets the date
func (b *PaidDateBuilder) WithDate(date *time.Time) *PaidDateBuilder {
	b.date = date
	return b
}

// WithTransactionGroupId sets the transaction group ID
func (b *PaidDateBuilder) WithTransactionGroupId(id *string) *PaidDateBuilder {
	b.transactionGroupId = id
	return b
}

// WithTransactionJournalId sets the transaction journal ID
func (b *PaidDateBuilder) WithTransactionJournalId(id *string) *PaidDateBuilder {
	b.transactionJournalId = id
	return b
}

// Build creates a PaidDate
func (b *PaidDateBuilder) Build() PaidDate {
	return PaidDate{
		date:                 b.date,
		transactionGroupId:   b.transactionGroupId,
		transactionJournalId: b.transactionJournalId,
	}
}

// InsightCategoryEntryBuilder is a builder for InsightCategoryEntry
type InsightCategoryEntryBuilder struct {
	id           string
	name         string
	amount       string
	currencyCode string
}

// NewInsightCategoryEntryBuilder creates a new InsightCategoryEntryBuilder
func NewInsightCategoryEntryBuilder() *InsightCategoryEntryBuilder {
	return &InsightCategoryEntryBuilder{}
}

// WithId sets the ID
func (b *InsightCategoryEntryBuilder) WithId(id string) *InsightCategoryEntryBuilder {
	b.id = id
	return b
}

// WithName sets the name
func (b *InsightCategoryEntryBuilder) WithName(name string) *InsightCategoryEntryBuilder {
	b.name = name
	return b
}

// WithAmount sets the amount
func (b *InsightCategoryEntryBuilder) WithAmount(amount string) *InsightCategoryEntryBuilder {
	b.amount = amount
	return b
}

// WithCurrencyCode sets the currency code
func (b *InsightCategoryEntryBuilder) WithCurrencyCode(currencyCode string) *InsightCategoryEntryBuilder {
	b.currencyCode = currencyCode
	return b
}

// Build creates an InsightCategoryEntry
func (b *InsightCategoryEntryBuilder) Build() InsightCategoryEntry {
	return InsightCategoryEntry{
		id:           b.id,
		name:         b.name,
		amount:       b.amount,
		currencyCode: b.currencyCode,
	}
}

// InsightTotalEntryBuilder is a builder for InsightTotalEntry
type InsightTotalEntryBuilder struct {
	amount       string
	currencyCode string
}

// NewInsightTotalEntryBuilder creates a new InsightTotalEntryBuilder
func NewInsightTotalEntryBuilder() *InsightTotalEntryBuilder {
	return &InsightTotalEntryBuilder{}
}

// WithAmount sets the amount
func (b *InsightTotalEntryBuilder) WithAmount(amount string) *InsightTotalEntryBuilder {
	b.amount = amount
	return b
}

// WithCurrencyCode sets the currency code
func (b *InsightTotalEntryBuilder) WithCurrencyCode(currencyCode string) *InsightTotalEntryBuilder {
	b.currencyCode = currencyCode
	return b
}

// Build creates an InsightTotalEntry
func (b *InsightTotalEntryBuilder) Build() InsightTotalEntry {
	return InsightTotalEntry{
		amount:       b.amount,
		currencyCode: b.currencyCode,
	}
}

// InsightCategoryResponseBuilder is a builder for InsightCategoryResponse
type InsightCategoryResponseBuilder struct {
	entries []InsightCategoryEntry
}

// NewInsightCategoryResponseBuilder creates a new InsightCategoryResponseBuilder
func NewInsightCategoryResponseBuilder() *InsightCategoryResponseBuilder {
	return &InsightCategoryResponseBuilder{
		entries: []InsightCategoryEntry{},
	}
}

// WithEntries sets the entries
func (b *InsightCategoryResponseBuilder) WithEntries(entries []InsightCategoryEntry) *InsightCategoryResponseBuilder {
	b.entries = entries
	return b
}

// AddEntry adds a single entry
func (b *InsightCategoryResponseBuilder) AddEntry(entry InsightCategoryEntry) *InsightCategoryResponseBuilder {
	b.entries = append(b.entries, entry)
	return b
}

// Build creates an InsightCategoryResponse
func (b *InsightCategoryResponseBuilder) Build() InsightCategoryResponse {
	// Create a copy to ensure immutability
	entriesCopy := make([]InsightCategoryEntry, len(b.entries))
	copy(entriesCopy, b.entries)
	
	return InsightCategoryResponse{
		entries: entriesCopy,
	}
}

// InsightTotalResponseBuilder is a builder for InsightTotalResponse
type InsightTotalResponseBuilder struct {
	entries []InsightTotalEntry
}

// NewInsightTotalResponseBuilder creates a new InsightTotalResponseBuilder
func NewInsightTotalResponseBuilder() *InsightTotalResponseBuilder {
	return &InsightTotalResponseBuilder{
		entries: []InsightTotalEntry{},
	}
}

// WithEntries sets the entries
func (b *InsightTotalResponseBuilder) WithEntries(entries []InsightTotalEntry) *InsightTotalResponseBuilder {
	b.entries = entries
	return b
}

// AddEntry adds a single entry
func (b *InsightTotalResponseBuilder) AddEntry(entry InsightTotalEntry) *InsightTotalResponseBuilder {
	b.entries = append(b.entries, entry)
	return b
}

// Build creates an InsightTotalResponse
func (b *InsightTotalResponseBuilder) Build() InsightTotalResponse {
	// Create a copy to ensure immutability
	entriesCopy := make([]InsightTotalEntry, len(b.entries))
	copy(entriesCopy, b.entries)
	
	return InsightTotalResponse{
		entries: entriesCopy,
	}
}