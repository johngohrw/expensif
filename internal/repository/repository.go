package repository

import (
	"context"

	"expensif/internal/domain"
)

type Repository interface {
	CreateExpense(ctx context.Context, e domain.Expense) (int64, error)
	ListExpenses(ctx context.Context, limit int) ([]domain.Expense, error)
	GetExpense(ctx context.Context, id int64) (*domain.Expense, error)
	UpdateExpense(ctx context.Context, e domain.Expense) error
	DeleteExpense(ctx context.Context, id int64) error
	ListCategories(ctx context.Context) ([]string, error)
	SummaryByCategory(ctx context.Context) (map[string]float64, error)
	TotalExpenses(ctx context.Context) (float64, error)
	GetPreferences(ctx context.Context) (*domain.Preferences, error)
	SavePreferences(ctx context.Context, p domain.Preferences) error

	SaveRates(ctx context.Context, base string, date string, rates map[string]float64) error
	GetRates(ctx context.Context, base string, date string) (map[string]float64, error)
	GetLatestRates(ctx context.Context, base string) (map[string]float64, string, error)

	ListUsers(ctx context.Context) ([]domain.User, error)
	SaveUser(ctx context.Context, name string) error
	CreateUser(ctx context.Context, name string) (int64, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	UpdateUser(ctx context.Context, id int64, name string) error
	DeleteUser(ctx context.Context, id int64) error
	ClearExpensePaidBy(ctx context.Context, userName string) error
}
