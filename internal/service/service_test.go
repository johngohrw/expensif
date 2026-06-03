package service

import (
	"context"
	"errors"
	"testing"

	"expensif/internal/domain"
	"expensif/internal/repository"
)

// mockRateClient is a test double for RateFetcher.
type mockRateClient struct {
	rates map[string]float64
	date  string
	err   error
}

func (m *mockRateClient) Latest(ctx context.Context, base string) (map[string]float64, string, error) {
	return m.rates, m.date, m.err
}

// mockRateRepo implements only the RateRepository interface for testing.
type mockRateRepo struct {
	saved map[string]map[string]float64 // base -> target -> rate
	date  string
}

func (r *mockRateRepo) SaveRates(ctx context.Context, base string, date string, rates map[string]float64) error {
	if r.saved == nil {
		r.saved = make(map[string]map[string]float64)
	}
	r.saved[base] = make(map[string]float64, len(rates))
	for k, v := range rates {
		r.saved[base][k] = v
	}
	r.date = date
	return nil
}

func (r *mockRateRepo) GetRates(ctx context.Context, base string, date string) (map[string]float64, error) {
	return nil, errors.New("no rates for date")
}

func (r *mockRateRepo) GetLatestRates(ctx context.Context, base string) (map[string]float64, string, error) {
	if r.saved == nil || r.saved[base] == nil {
		return nil, "", errors.New("no rates")
	}
	cp := make(map[string]float64, len(r.saved[base]))
	for k, v := range r.saved[base] {
		cp[k] = v
	}
	return cp, r.date, nil
}

// Compile-time interface checks.
var _ RateFetcher = (*mockRateClient)(nil)
var _ repository.RateRepository = (*mockRateRepo)(nil)

// --- RefreshRates ---

func TestRefreshRates_SavesToRepository(t *testing.T) {
	mockRates := &mockRateClient{
		rates: map[string]float64{"EUR": 0.85, "JPY": 150.0},
		date:  "2024-06-01",
	}
	repo := &mockRateRepo{}

	svc := New(repository.Repos{Rates: repo}, mockRates)

	if err := svc.RefreshRates(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.date != "2024-06-01" {
		t.Fatalf("expected date 2024-06-01, got %s", repo.date)
	}
	if repo.saved["USD"]["EUR"] != 0.85 {
		t.Fatalf("expected EUR rate 0.85, got %v", repo.saved["USD"]["EUR"])
	}
	if repo.saved["USD"]["JPY"] != 150.0 {
		t.Fatalf("expected JPY rate 150.0, got %v", repo.saved["USD"]["JPY"])
	}
}

func TestRefreshRates_FetchError_NoSave(t *testing.T) {
	mockRates := &mockRateClient{
		rates: nil,
		date:  "",
		err:   errors.New("network timeout"),
	}
	repo := &mockRateRepo{}

	svc := New(repository.Repos{Rates: repo}, mockRates)

	err := svc.RefreshRates(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if repo.saved != nil {
		t.Fatal("expected no save on fetch error")
	}
}

// --- ConvertWithRates ---

func TestConvertWithRates_SameCurrency(t *testing.T) {
	svc := New(repository.Repos{}, &mockRateClient{})
	result, err := svc.ConvertWithRates(100.0, "USD", "USD", map[string]float64{"USD": 1.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 100.0 {
		t.Fatalf("expected 100.0, got %f", result)
	}
}

func TestConvertWithRates_CrossCurrency(t *testing.T) {
	svc := New(repository.Repos{}, &mockRateClient{})
	// Rates relative to USD: EUR=0.85, GBP=0.79
	rates := map[string]float64{"USD": 1.0, "EUR": 0.85, "GBP": 0.79}

	// 100 USD -> EUR = 100 * 0.85 / 1.0 = 85
	result, err := svc.ConvertWithRates(100.0, "USD", "EUR", rates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 85.0 {
		t.Fatalf("expected 85.0 USD->EUR, got %f", result)
	}

	// 100 GBP -> EUR = 100 * 0.85 / 0.79 = ~107.59
	result, err = svc.ConvertWithRates(100.0, "GBP", "EUR", rates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 100.0 * 0.85 / 0.79
	if result != expected {
		t.Fatalf("expected %f GBP->EUR, got %f", expected, result)
	}
}

func TestConvertWithRates_MissingRate(t *testing.T) {
	svc := New(repository.Repos{}, &mockRateClient{})
	rates := map[string]float64{"USD": 1.0}

	_, err := svc.ConvertWithRates(100.0, "USD", "EUR", rates)
	if err == nil {
		t.Fatal("expected error for missing rate, got nil")
	}
}

// --- ConvertExpensesTotal ---

func TestConvertExpensesTotal_ConvertsAndSums(t *testing.T) {
	repo := &mockRateRepo{
		saved: map[string]map[string]float64{
			"USD": {"EUR": 0.85, "JPY": 150.0, "USD": 1.0},
		},
		date: "2024-06-01",
	}
	svc := New(repository.Repos{Rates: repo}, &mockRateClient{})

	expenses := []domain.Expense{
		{Amount: 100.0, Currency: "USD"},
		{Amount: 100.0, Currency: "EUR"}, // 100 EUR -> USD = 100 * 1.0 / 0.85 = ~117.65
		{Amount: 150.0, Currency: "JPY"}, // 150 JPY -> USD = 150 * 1.0 / 150.0 = 1.0
	}

	total, date, err := svc.ConvertExpensesTotal(context.Background(), expenses, "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if date != "2024-06-01" {
		t.Fatalf("expected date 2024-06-01, got %s", date)
	}
	expected := 100.0 + (100.0/0.85) + 1.0
	if total != expected {
		t.Fatalf("expected total %f, got %f", expected, total)
	}
}

func TestConvertExpensesTotal_NoRates(t *testing.T) {
	repo := &mockRateRepo{}
	svc := New(repository.Repos{Rates: repo}, &mockRateClient{})

	_, _, err := svc.ConvertExpensesTotal(context.Background(), []domain.Expense{
		{Amount: 100.0, Currency: "USD"},
	}, "USD")
	if !errors.Is(err, ErrNoRates) {
		t.Fatalf("expected ErrNoRates, got %v", err)
	}
}

func TestConvertExpensesTotal_SkipsUnconvertible(t *testing.T) {
	repo := &mockRateRepo{
		saved: map[string]map[string]float64{
			"USD": {"USD": 1.0},
		},
		date: "2024-06-01",
	}
	svc := New(repository.Repos{Rates: repo}, &mockRateClient{})

	expenses := []domain.Expense{
		{Amount: 100.0, Currency: "USD"},
		{Amount: 50.0, Currency: "XXX"}, // no rate for XXX — should be skipped
	}

	total, _, err := svc.ConvertExpensesTotal(context.Background(), expenses, "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 100.0 {
		t.Fatalf("expected 100.0 (skipping unconvertible), got %f", total)
	}
}

// --- CreateExpense validation ---

func TestCreateExpense_Defaults(t *testing.T) {
	// We need a minimal expense repository to capture the created expense.
	var captured domain.Expense
	mockExpenses := &expenseRepoStub{create: func(e domain.Expense) (int64, error) {
		captured = e
		return 1, nil
	}}

	svc := New(repository.Repos{Expenses: mockExpenses}, &mockRateClient{})
	ctx := context.Background()

	id, err := svc.CreateExpense(ctx, 10.0, "food", "lunch", "", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
	if captured.Date == "" {
		t.Fatal("expected default date")
	}
	if captured.Currency != "USD" {
		t.Fatalf("expected default currency USD, got %s", captured.Currency)
	}
}

func TestCreateExpense_InvalidAmount(t *testing.T) {
	svc := New(repository.Repos{}, &mockRateClient{})
	_, err := svc.CreateExpense(context.Background(), -5, "food", "lunch", "2024-06-01", "USD", 0)
	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestCreateExpense_MissingCategory(t *testing.T) {
	svc := New(repository.Repos{}, &mockRateClient{})
	_, err := svc.CreateExpense(context.Background(), 10.0, "", "lunch", "2024-06-01", "USD", 0)
	if !errors.Is(err, ErrMissingCategory) {
		t.Fatalf("expected ErrMissingCategory, got %v", err)
	}
}

func TestCreateExpense_MissingDescription(t *testing.T) {
	svc := New(repository.Repos{}, &mockRateClient{})
	_, err := svc.CreateExpense(context.Background(), 10.0, "food", "", "2024-06-01", "USD", 0)
	if !errors.Is(err, ErrMissingDescription) {
		t.Fatalf("expected ErrMissingDescription, got %v", err)
	}
}

// expenseRepoStub is a minimal stub for expense repository tests.
type expenseRepoStub struct {
	create func(e domain.Expense) (int64, error)
}

func (s *expenseRepoStub) CreateExpense(ctx context.Context, e domain.Expense) (int64, error) {
	if s.create != nil {
		return s.create(e)
	}
	return 0, nil
}

func (s *expenseRepoStub) ListExpenses(ctx context.Context, limit int) ([]domain.Expense, error) {
	return nil, nil
}

func (s *expenseRepoStub) GetExpense(ctx context.Context, id int64) (*domain.Expense, error) {
	return nil, nil
}

func (s *expenseRepoStub) UpdateExpense(ctx context.Context, e domain.Expense) error {
	return nil
}

func (s *expenseRepoStub) DeleteExpense(ctx context.Context, id int64) error {
	return nil
}

func (s *expenseRepoStub) ListCategories(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (s *expenseRepoStub) SummaryByCategory(ctx context.Context) (map[string]float64, error) {
	return nil, nil
}

func (s *expenseRepoStub) TotalExpenses(ctx context.Context) (float64, error) {
	return 0, nil
}

func (s *expenseRepoStub) GetEarliestExpenseDate(ctx context.Context) (string, error) {
	return "", nil
}

func (s *expenseRepoStub) ListDescriptionsByCategory(ctx context.Context, category string, limit int) ([]domain.DescriptionCount, error) {
	return nil, nil
}
