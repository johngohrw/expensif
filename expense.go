package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Expense struct {
	ID          int64
	Amount      float64
	Category    string
	Description string
	Date        string // YYYY-MM-DD
	Currency    string
	CreatedAt   time.Time
}

func addExpense(amount float64, category, description, date, currency string) (int64, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if currency == "" {
		currency = "USD"
	}
	res, err := db.Exec(
		`INSERT INTO expenses (amount, category, description, date, currency) VALUES (?, ?, ?, ?, ?)`,
		amount, category, description, date, currency,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func listExpenses(limit int) ([]Expense, error) {
	rows, err := db.Query(
		`SELECT id, amount, category, description, date, currency, created_at FROM expenses ORDER BY date DESC, created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		var createdAt string
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Description, &e.Date, &e.Currency, &createdAt); err != nil {
			return nil, err
		}
		e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}

func deleteExpense(id int64) error {
	res, err := db.Exec(`DELETE FROM expenses WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("no expense with id %d", id)
	}
	return nil
}

func summaryByCategory() (map[string]float64, error) {
	rows, err := db.Query(`SELECT category, SUM(amount) FROM expenses GROUP BY CATEGORY`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[string]float64)
	for rows.Next() {
		var cat string
		var sum float64
		if err := rows.Scan(&cat, &sum); err != nil {
			return nil, err
		}
		m[cat] = sum
	}
	return m, rows.Err()
}

func getExpense(id int64) (*Expense, error) {
	var e Expense
	var createdAt string
	err := db.QueryRow(
		`SELECT id, amount, category, description, date, currency, created_at FROM expenses WHERE id = ?`, id,
	).Scan(&e.ID, &e.Amount, &e.Category, &e.Description, &e.Date, &e.Currency, &createdAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no expense with id %d", id)
	}
	if err != nil {
		return nil, err
	}
	e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &e, nil
}

func updateExpense(id int64, amount float64, category, description, date, currency string) error {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if currency == "" {
		currency = "USD"
	}
	res, err := db.Exec(
		`UPDATE expenses SET amount = ?, category = ?, description = ?, date = ?, currency = ? WHERE id = ?`,
		amount, category, description, date, currency, id,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("no expense with id %d", id)
	}
	return nil
}

func listCategories() ([]string, error) {
	rows, err := db.Query(`
		SELECT category FROM expenses
		WHERE date >= date('now', '-3 months')
		GROUP BY category
		ORDER BY COUNT(*) DESC, category ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, rows.Err()
}

func totalExpenses() (float64, error) {
	var total sql.NullFloat64
	err := db.QueryRow(`SELECT SUM(amount) FROM expenses`).Scan(&total)
	if err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return total.Float64, nil
}

type Preferences struct {
	Currency string
	DarkMode bool
}

var currentPrefs *Preferences

func loadPreferences() error {
	var p Preferences
	var darkMode int
	err := db.QueryRow(`SELECT currency, dark_mode FROM preferences WHERE id = 1`).Scan(&p.Currency, &darkMode)
	if err != nil {
		return err
	}
	p.DarkMode = darkMode != 0
	currentPrefs = &p
	return nil
}

func getPreferences() *Preferences {
	if currentPrefs == nil {
		loadPreferences()
	}
	if currentPrefs == nil {
		return &Preferences{Currency: "USD", DarkMode: false}
	}
	return currentPrefs
}

func savePreferences(currency string, darkMode bool) error {
	dm := 0
	if darkMode {
		dm = 1
	}
	_, err := db.Exec(`UPDATE preferences SET currency = ?, dark_mode = ? WHERE id = 1`, currency, dm)
	if err != nil {
		return err
	}
	currentPrefs = &Preferences{Currency: currency, DarkMode: darkMode}
	return nil
}

var currencySymbols = map[string]string{
	"USD":  "$",
	"MYR":  "RM",
	"JPY":  "¥",
	"CNY":  "¥",
	"THB":  "฿",
	"EUR":  "€",
	"GBP":  "£",
	"SGD":  "S$",
	"KRW":  "₩",
	"AUD":  "A$",
	"CAD":  "C$",
	"INR":  "₹",
	"VND":  "₫",
	"PHP":  "₱",
	"IDR":  "Rp",
	"HKD":  "HK$",
	"TWD":  "NT$",
}

func currencySymbol(code string) string {
	if s, ok := currencySymbols[code]; ok {
		return s
	}
	return "$"
}
