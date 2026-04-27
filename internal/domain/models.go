package domain

import "time"

type Expense struct {
	ID              int64
	Amount          float64
	Category        string
	Description     string
	Date            string // YYYY-MM-DD
	Currency        string
	CreatedAt       time.Time
	ConvertedAmount float64 `json:"-"` // computed at render time
}

type Preferences struct {
	Currency string
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
