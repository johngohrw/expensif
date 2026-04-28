package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"expensif/internal/domain"
	"expensif/internal/rate"
	"expensif/internal/repository"
)

var (
	ErrInvalidAmount      = errors.New("amount must be greater than 0")
	ErrMissingCategory    = errors.New("category is required")
	ErrMissingDescription = errors.New("description is required")
	ErrNoRates            = errors.New("no exchange rates available")
)

type Service struct {
	repo       repository.Repository
	rateClient *rate.Client
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo, rateClient: rate.NewClient()}
}

// --- Expenses ---

func (s *Service) CreateExpense(ctx context.Context, amount float64, category, description, date, currency string, paidByID int64) (int64, error) {
	if amount <= 0 {
		return 0, ErrInvalidAmount
	}
	category = strings.TrimSpace(category)
	description = strings.TrimSpace(description)
	if category == "" {
		return 0, ErrMissingCategory
	}
	if description == "" {
		return 0, ErrMissingDescription
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if currency == "" {
		currency = "USD"
	}
	e := domain.Expense{
		Amount:   amount,
		Category: category,
		Description: description,
		Date:     date,
		Currency: currency,
		PaidByID: paidByID,
	}
	return s.repo.CreateExpense(ctx, e)
}

func (s *Service) ListExpenses(ctx context.Context, limit int) ([]domain.Expense, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.repo.ListExpenses(ctx, limit)
}

func (s *Service) GetExpense(ctx context.Context, id int64) (*domain.Expense, error) {
	return s.repo.GetExpense(ctx, id)
}

func (s *Service) UpdateExpense(ctx context.Context, id int64, amount float64, category, description, date, currency string, paidByID int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	category = strings.TrimSpace(category)
	description = strings.TrimSpace(description)
	if category == "" {
		return ErrMissingCategory
	}
	if description == "" {
		return ErrMissingDescription
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if currency == "" {
		currency = "USD"
	}
	e := domain.Expense{
		ID:          id,
		Amount:      amount,
		Category:    category,
		Description: description,
		Date:        date,
		Currency:    currency,
		PaidByID:    paidByID,
	}
	return s.repo.UpdateExpense(ctx, e)
}

func (s *Service) DeleteExpense(ctx context.Context, id int64) error {
	return s.repo.DeleteExpense(ctx, id)
}

func (s *Service) ListCategories(ctx context.Context) ([]string, error) {
	return s.repo.ListCategories(ctx)
}

func (s *Service) SummaryByCategory(ctx context.Context) (map[string]float64, error) {
	return s.repo.SummaryByCategory(ctx)
}

func (s *Service) TotalExpenses(ctx context.Context) (float64, error) {
	return s.repo.TotalExpenses(ctx)
}

func (s *Service) DailyGroups(ctx context.Context, limit int) ([]domain.DailyGroup, error) {
	expenses, err := s.ListExpenses(ctx, limit)
	if err != nil {
		return nil, err
	}
	groups := make(map[string][]domain.Expense)
	for _, e := range expenses {
		groups[e.Date] = append(groups[e.Date], e)
	}
	var dailyGroups []domain.DailyGroup
	for date, exps := range groups {
		var dayTotal float64
		for _, e := range exps {
			dayTotal += e.Amount
		}
		dailyGroups = append(dailyGroups, domain.DailyGroup{
			Date:     date,
			Expenses: exps,
			Total:    dayTotal,
		})
	}
	sort.Slice(dailyGroups, func(i, j int) bool {
		return dailyGroups[i].Date > dailyGroups[j].Date
	})
	return dailyGroups, nil
}

// --- Preferences ---

func (s *Service) Preferences(ctx context.Context) (*domain.Preferences, error) {
	p, err := s.repo.GetPreferences(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.Preferences{Currency: "USD"}, nil
		}
		return nil, err
	}
	return p, nil
}

func (s *Service) SavePreferences(ctx context.Context, currency string, userID int64) error {
	if currency == "" {
		currency = "USD"
	}
	return s.repo.SavePreferences(ctx, domain.Preferences{
		Currency: currency,
		UserID:   userID,
	})
}

func (s *Service) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *Service) CreateUser(ctx context.Context, name string) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errors.New("name is required")
	}
	return s.repo.CreateUser(ctx, name)
}

func (s *Service) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, id int64, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is required")
	}
	return s.repo.UpdateUser(ctx, id, name)
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}

// --- Exchange Rates ---

func (s *Service) RefreshRates(ctx context.Context) error {
	rates, date, err := s.rateClient.Latest(ctx, "USD")
	if err != nil {
		return fmt.Errorf("fetch rates: %w", err)
	}
	return s.repo.SaveRates(ctx, "USD", date, rates)
}

func (s *Service) GetRatesForConversion(ctx context.Context) (map[string]float64, string, error) {
	today := time.Now().Format("2006-01-02")
	rates, err := s.repo.GetRates(ctx, "USD", today)
	if err == nil && len(rates) > 0 {
		return rates, today, nil
	}
	rates, date, err := s.repo.GetLatestRates(ctx, "USD")
	if err != nil {
		return nil, "", ErrNoRates
	}
	if len(rates) == 0 {
		return nil, "", ErrNoRates
	}
	return rates, date, nil
}

func (s *Service) ConvertWithRates(amount float64, from, to string, rates map[string]float64) (float64, error) {
	if from == to {
		return amount, nil
	}
	fromRate, ok := rates[from]
	if !ok {
		return 0, fmt.Errorf("no rate for %s", from)
	}
	toRate, ok := rates[to]
	if !ok {
		return 0, fmt.Errorf("no rate for %s", to)
	}
	return amount * toRate / fromRate, nil
}

func (s *Service) ConvertExpensesTotal(ctx context.Context, expenses []domain.Expense, target string) (float64, string, error) {
	rates, date, err := s.GetRatesForConversion(ctx)
	if err != nil {
		return 0, "", err
	}
	var total float64
	for _, e := range expenses {
		conv, err := s.ConvertWithRates(e.Amount, e.Currency, target, rates)
		if err != nil {
			continue // skip unconvertible currencies
		}
		total += conv
	}
	return total, date, nil
}
