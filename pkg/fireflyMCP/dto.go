package fireflyMCP

import "time"

type Pagination struct {
	Count       int `json:"count"`
	Total       int `json:"total"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
}
type Spent struct {
	Sum          string `json:"sum"`
	CurrencyCode string `json:"currency_code"`
}

type Budget struct {
	Id     string  `json:"id"`
	Active bool    `json:"active"`
	Name   string  `json:"name"`
	Notes  *string `json:"notes"`
	Spent  Spent   `json:"spent"`
}
type BudgetList struct {
	Data       []Budget   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Category struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Notes *string `json:"notes"`
}

type CategoryList struct {
	Data       []Category `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Account struct {
	Id     string  `json:"id"`
	Active bool    `json:"active"`
	Name   string  `json:"name"`
	Notes  *string `json:"notes"`
	Type   string  `json:"type"`
}

type AccountList struct {
	Data       []Account  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Transaction struct {
	Id              string  `json:"id"`
	Amount          string  `json:"amount"`
	BillId          *string `json:"bill_id"`
	BillName        *string `json:"bill_name"`
	BudgetId        *string `json:"budget_id"`
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

type TransactionGroup struct {
	Id           string        `json:"id"`
	GroupTitle   string        `json:"group_title"`
	Transactions []Transaction `json:"transactions"`
}

type TransactionList struct {
	Data       []TransactionGroup `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

type BasicSummary struct {
	Key           string `json:"key"`
	Title         string `json:"title"`
	CurrencyCode  string `json:"currency_code"`
	MonetaryValue string `json:"monetary_value"`
}

type BasicSummaryList struct {
	Data []BasicSummary `json:"data"`
}

type InsightCategoryEntry struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type InsightTotalEntry struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type InsightCategoryResponse struct {
	Entries []InsightCategoryEntry `json:"entries"`
}

type InsightTotalResponse struct {
	Entries []InsightTotalEntry `json:"entries"`
}

type BudgetSpent struct {
	Sum            string `json:"sum"`
	CurrencyCode   string `json:"currency_code"`
	CurrencySymbol string `json:"currency_symbol"`
}

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

type BudgetLimitList struct {
	Data       []BudgetLimit `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

type Tag struct {
	Id          string  `json:"id"`
	Tag         string  `json:"tag"`
	Description *string `json:"description"`
}

type TagList struct {
	Data       []Tag      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type PaidDate struct {
	Date                 *time.Time `json:"date"`
	TransactionGroupId   *string    `json:"transaction_group_id"`
	TransactionJournalId *string    `json:"transaction_journal_id"`
}

type Bill struct {
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

type BillList struct {
	Data       []Bill     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type RecurrenceRepetition struct {
	Id          string  `json:"id"`
	Type        string  `json:"type"`
	Moment      string  `json:"moment"`
	Skip        int     `json:"skip"`
	Weekend     int     `json:"weekend"`
	Description *string `json:"description"`
}

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

type Recurrence struct {
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

type RecurrenceList struct {
	Data       []Recurrence `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

// TransactionStoreRequest represents the request body for creating a new transaction
type TransactionStoreRequest struct {
	ErrorIfDuplicateHash bool                      `json:"error_if_duplicate_hash,omitempty" jsonschema:"Break if transaction with same hash already exists (default: false)"` // Break if transaction already exists
	ApplyRules           bool                      `json:"apply_rules,omitempty" jsonschema:"Whether to apply processing rules when creating transaction (default: false)"`    // Whether to apply rules when submitting
	FireWebhooks         bool                      `json:"fire_webhooks,omitempty" jsonschema:"Whether to fire webhooks for this transaction (default: true)"`                 // Whether to fire webhooks (default: true)
	GroupTitle           string                    `json:"group_title,omitempty" jsonschema:"Title for the transaction group (for split transactions)"`                        // Title for split transactions
	Transactions         []TransactionSplitRequest `json:"transactions" jsonschema:"Array of transactions to create (required, at least one)"`                                 // Array of transactions (required)
}

// TransactionSplitRequest represents a single transaction in a transaction group
type TransactionSplitRequest struct {
	Type                string   `json:"type" jsonschema:"Transaction type: withdrawal, deposit, transfer (required)"`                                     // Transaction type: withdrawal, deposit, transfer (required)
	Date                string   `json:"date" jsonschema:"Transaction date in RFC3339 format, e.g. 2024-01-15T00:00:00Z (required)"`                          // Transaction date in RFC3339 format (required)
	Amount              string   `json:"amount" jsonschema:"Transaction amount as string (e.g. '100.00') (required)"`                                      // Transaction amount (required)
	Description         string   `json:"description" jsonschema:"Transaction description (required)"`                                                      // Transaction description (required)
	SourceId            *string  `json:"source_id,omitempty" jsonschema:"Source account ID (use either source_id or source_name)"`                         // Source account ID
	SourceName          *string  `json:"source_name,omitempty" jsonschema:"Source account name (use either source_id or source_name)"`                     // Source account name
	DestinationId       *string  `json:"destination_id,omitempty" jsonschema:"Destination account ID (use either destination_id or destination_name)"`     // Destination account ID
	DestinationName     *string  `json:"destination_name,omitempty" jsonschema:"Destination account name (use either destination_id or destination_name)"` // Destination account name
	CategoryId          *string  `json:"category_id,omitempty" jsonschema:"Category ID (use either category_id or category_name)"`                         // Category ID
	CategoryName        *string  `json:"category_name,omitempty" jsonschema:"Category name (use either category_id or category_name)"`                     // Category name
	BudgetId            *string  `json:"budget_id,omitempty" jsonschema:"Budget ID (use either budget_id or budget_name)"`                                 // Budget ID
	BudgetName          *string  `json:"budget_name,omitempty" jsonschema:"Budget name (use either budget_id or budget_name)"`                             // Budget name
	Tags                []string `json:"tags,omitempty" jsonschema:"Array of tag names to attach to transaction"`                                          // Transaction tags
	CurrencyId          *string  `json:"currency_id,omitempty" jsonschema:"Currency ID for the transaction"`                                               // Currency ID
	CurrencyCode        *string  `json:"currency_code,omitempty" jsonschema:"Currency code (e.g. 'USD', 'EUR')"`                                           // Currency code
	ForeignAmount       *string  `json:"foreign_amount,omitempty" jsonschema:"Amount in foreign currency as string"`                                       // Amount in foreign currency
	ForeignCurrencyId   *string  `json:"foreign_currency_id,omitempty" jsonschema:"Foreign currency ID"`                                                   // Foreign currency ID
	ForeignCurrencyCode *string  `json:"foreign_currency_code,omitempty" jsonschema:"Foreign currency code (e.g. 'USD', 'EUR')"`                           // Foreign currency code
	BillId              *string  `json:"bill_id,omitempty" jsonschema:"Bill ID to link this transaction to"`                                               // Bill ID
	BillName            *string  `json:"bill_name,omitempty" jsonschema:"Bill name to link this transaction to"`                                           // Bill name
	PiggyBankId         *string  `json:"piggy_bank_id,omitempty" jsonschema:"Piggy bank ID for savings transfers"`                                         // Piggy bank ID
	PiggyBankName       *string  `json:"piggy_bank_name,omitempty" jsonschema:"Piggy bank name for savings transfers"`                                     // Piggy bank name
	Notes               *string  `json:"notes,omitempty" jsonschema:"Additional notes or comments for the transaction"`                                    // Transaction notes
	Reconciled          *bool    `json:"reconciled,omitempty" jsonschema:"Whether the transaction has been reconciled (default: false)"`                   // Whether transaction is reconciled
	Order               *int     `json:"order,omitempty" jsonschema:"Order of this split in the transaction group"`                                        // Order in the list
}

// ReceiptItemRequest represents a single item on a shopping receipt
type ReceiptItemRequest struct {
	Amount       string   `json:"amount" jsonschema:"Item amount as string (e.g. '12.99') (required)"`
	Description  string   `json:"description" jsonschema:"Item description (required)"`
	CategoryId   *string  `json:"category_id,omitempty" jsonschema:"Category ID (overrides default_category)"`
	CategoryName *string  `json:"category_name,omitempty" jsonschema:"Category name (overrides default_category)"`
	BudgetId     *string  `json:"budget_id,omitempty" jsonschema:"Budget ID"`
	BudgetName   *string  `json:"budget_name,omitempty" jsonschema:"Budget name"`
	Tags         []string `json:"tags,omitempty" jsonschema:"Additional tags for this item (merged with receipt tags)"`
	Notes        *string  `json:"notes,omitempty" jsonschema:"Item-specific notes"`
}

// ReceiptStoreRequest represents a request to store a shopping receipt with multiple items
type ReceiptStoreRequest struct {
	StoreName           string               `json:"store_name" jsonschema:"Name of the store/shop (required, becomes destination expense account)"`
	ReceiptDate         string               `json:"receipt_date" jsonschema:"Date of the receipt in RFC3339 format, e.g. 2024-01-15T00:00:00Z (required)"`
	Items               []ReceiptItemRequest `json:"items" jsonschema:"Array of receipt items/transactions (required, at least one)"`
	TotalAmount         *string              `json:"total_amount,omitempty" jsonschema:"Expected total amount - if provided, validates that sum of items equals this"`
	SourceAccountId     *string              `json:"source_account_id,omitempty" jsonschema:"Source account ID for payment (use either source_account_id or source_account_name)"`
	SourceAccountName   *string              `json:"source_account_name,omitempty" jsonschema:"Source account name for payment (use either source_account_id or source_account_name)"`
	DefaultCategoryId   *string              `json:"default_category_id,omitempty" jsonschema:"Default category ID for items without category"`
	DefaultCategoryName *string              `json:"default_category_name,omitempty" jsonschema:"Default category name for items without category"`
	CurrencyCode        *string              `json:"currency_code,omitempty" jsonschema:"Currency code for all transactions (e.g. 'USD', 'EUR')"`
	Tags                []string             `json:"tags,omitempty" jsonschema:"Tags to apply to all transactions"`
	Notes               *string              `json:"notes,omitempty" jsonschema:"Notes for the receipt (applied to first transaction)"`
	ApplyRules          bool                 `json:"apply_rules,omitempty" jsonschema:"Whether to apply processing rules (default: false)"`
	FireWebhooks        bool                 `json:"fire_webhooks,omitempty" jsonschema:"Whether to fire webhooks (default: true)"`
}
