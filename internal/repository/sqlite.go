package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"expensif/internal/domain"
)

type sqliteRepo struct {
	db *sql.DB
}

func NewSQLite(db *sql.DB) Repository {
	return &sqliteRepo{db: db}
}

func (r *sqliteRepo) CreateExpense(ctx context.Context, e domain.Expense) (int64, error) {
	if e.Date == "" {
		e.Date = time.Now().Format("2006-01-02")
	}
	if e.Currency == "" {
		e.Currency = "USD"
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO expenses (amount, category, description, date, currency) VALUES (?, ?, ?, ?, ?)`,
		e.Amount, e.Category, e.Description, e.Date, e.Currency,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteRepo) ListExpenses(ctx context.Context, limit int) ([]domain.Expense, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, amount, category, description, date, currency, created_at FROM expenses ORDER BY date DESC, created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		var e domain.Expense
		var createdAt string
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Description, &e.Date, &e.Currency, &createdAt); err != nil {
			return nil, err
		}
		e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}

func (r *sqliteRepo) GetExpense(ctx context.Context, id int64) (*domain.Expense, error) {
	var e domain.Expense
	var createdAt string
	err := r.db.QueryRowContext(ctx,
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

func (r *sqliteRepo) UpdateExpense(ctx context.Context, e domain.Expense) error {
	if e.Date == "" {
		e.Date = time.Now().Format("2006-01-02")
	}
	if e.Currency == "" {
		e.Currency = "USD"
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE expenses SET amount = ?, category = ?, description = ?, date = ?, currency = ? WHERE id = ?`,
		e.Amount, e.Category, e.Description, e.Date, e.Currency, e.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("no expense with id %d", e.ID)
	}
	return nil
}

func (r *sqliteRepo) DeleteExpense(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE id = ?`, id)
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

func (r *sqliteRepo) ListCategories(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
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

func (r *sqliteRepo) SummaryByCategory(ctx context.Context) (map[string]float64, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT category, SUM(amount) FROM expenses GROUP BY category`)
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

func (r *sqliteRepo) TotalExpenses(ctx context.Context) (float64, error) {
	var total sql.NullFloat64
	err := r.db.QueryRowContext(ctx, `SELECT SUM(amount) FROM expenses`).Scan(&total)
	if err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return total.Float64, nil
}

func (r *sqliteRepo) GetPreferences(ctx context.Context) (*domain.Preferences, error) {
	var p domain.Preferences
	err := r.db.QueryRowContext(ctx, `SELECT currency FROM preferences WHERE id = 1`).Scan(&p.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("get preferences: %w", err)
	}
	return &p, nil
}

func (r *sqliteRepo) SavePreferences(ctx context.Context, p domain.Preferences) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO preferences (id, currency) VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET currency = excluded.currency
	`, p.Currency)
	if err != nil {
		return fmt.Errorf("save preferences: %w", err)
	}
	return nil
}

// --- Exchange Rates ---

func (r *sqliteRepo) SaveRates(ctx context.Context, base string, date string, rates map[string]float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO exchange_rates (base_currency, target_currency, rate, date)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for target, rate := range rates {
		if _, err := stmt.ExecContext(ctx, base, target, rate, date); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *sqliteRepo) GetRates(ctx context.Context, base string, date string) (map[string]float64, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT target_currency, rate FROM exchange_rates WHERE base_currency = ? AND date = ?`,
		base, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rates := make(map[string]float64)
	for rows.Next() {
		var target string
		var rate float64
		if err := rows.Scan(&target, &rate); err != nil {
			return nil, err
		}
		rates[target] = rate
	}
	return rates, rows.Err()
}

func (r *sqliteRepo) GetLatestRates(ctx context.Context, base string) (map[string]float64, string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT target_currency, rate, date
		FROM exchange_rates
		WHERE base_currency = ?
		  AND date = (SELECT MAX(date) FROM exchange_rates WHERE base_currency = ?)
	`, base, base)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	rates := make(map[string]float64)
	var date string
	for rows.Next() {
		var target string
		var rate float64
		if err := rows.Scan(&target, &rate, &date); err != nil {
			return nil, "", err
		}
		rates[target] = rate
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}
	if len(rates) == 0 {
		return nil, "", sql.ErrNoRows
	}
	return rates, date, nil
}
