package domain

import "time"

type Expense struct {
	ID              int64
	Amount          float64
	Category        string
	Description     string
	Date            string // YYYY-MM-DD
	Currency        string
	PaidByID        int64  `json:"paid_by_id,omitempty"`
	PaidByName      string `json:"paid_by_name,omitempty"` // computed at render time
	CreatedAt       time.Time
	ConvertedAmount float64 `json:"-"` // computed at render time
}

type Preferences struct {
	Currency string
	UserID   int64
}

type User struct {
	ID   int64
	Name string
}

type CategorySummary struct {
	Name   string
	Amount float64
}

type DailyGroup struct {
	Date           string
	Expenses       []Expense
	Total          float64
	ConvertedTotal float64
}
