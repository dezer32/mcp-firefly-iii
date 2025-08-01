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
