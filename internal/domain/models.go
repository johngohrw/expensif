package domain

import "time"

type Expense struct {
	ID              int64     `json:"id"`
	Amount          float64   `json:"amount"`
	Category        string    `json:"category"`
	Description     string    `json:"description"`
	Date            string    `json:"date"` // YYYY-MM-DD
	Currency        string    `json:"currency"`
	PaidByID        int64     `json:"paidById,omitempty"`
	PaidByName      string    `json:"paidByName,omitempty"` // computed at render time
	CreatedAt       time.Time `json:"createdAt"`
	ConvertedAmount float64   `json:"convertedAmount"` // computed at render time
}

type Preferences struct {
	Currency string `json:"currency"`
	UserID   int64  `json:"userId"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CategorySummary struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type DailyGroup struct {
	Date           string    `json:"date"`
	Expenses       []Expense `json:"expenses"`
	Total          float64   `json:"total"`
	ConvertedTotal float64   `json:"convertedTotal"`
}
