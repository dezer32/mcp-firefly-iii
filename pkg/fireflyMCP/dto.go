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
	Id     string      `json:"id"`
	Active bool        `json:"active"`
	Name   string      `json:"name"`
	Notes  interface{} `json:"notes"`
	Spent  Spent       `json:"spent"`
}
type BudgetList struct {
	Data       []Budget   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Category struct {
	Id    string      `json:"id"`
	Name  string      `json:"name"`
	Notes interface{} `json:"notes"`
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
