package fireflyMCP

import "time"

type Pagination struct {
	Count       int `json:"count"`
	Total       int `json:"total"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
}
type Spent []struct {
	Sum          string `json:"sum"`
	CurrencyCode string `json:"currency_code"`
}

type Budget struct {
	Id        string      `json:"id"`
	Active    bool        `json:"active"`
	Name      string      `json:"name"`
	Notes     interface{} `json:"notes"`
	Spent     Spent       `json:"spent"`
	UpdatedAt time.Time   `json:"updated_at"`
}
type BudgetList struct {
	Data       []Budget   `json:"data"`
	Pagination Pagination `json:"pagination"`
}
