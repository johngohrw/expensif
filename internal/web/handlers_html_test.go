package web

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"expensif/internal/domain"
	"expensif/internal/repository"
	"expensif/internal/service"
)

// TestHandleDaily_TodayNamedInPreferencesTimezone freezes the handler's
// clock late on a UTC day and asserts the window is anchored at today in
// the user's timezone. At 2026-07-11T23:00Z it is already July 12 in
// Pacific/Auckland; a handler that consults UTC ends the timeline a day
// early and today's row never renders.
func TestHandleDaily_TodayNamedInPreferencesTimezone(t *testing.T) {
	repo := newMockRepo()
	repo.prefs = domain.Preferences{Currency: "USD", Timezone: "Pacific/Auckland"}
	// One past expense so the timeline renders rather than the
	// zero-expense-ever empty state.
	repo.CreateExpense(t.Context(), domain.Expense{Amount: 5, Category: "food", Description: "x", Date: "2026-07-01", Currency: "USD"})

	svc := service.New(repository.Repos{
		Expenses:    repo,
		Users:       repo,
		Preferences: repo,
		Rates:       repo,
	}, &mockRateClient{})

	renderer, err := NewRenderer("../../templates", false, nil)
	if err != nil {
		t.Fatalf("renderer init failed: %v", err)
	}
	h := NewHTMLHandler(svc, renderer)
	h.now = func() time.Time { return time.Date(2026, 7, 11, 23, 0, 0, 0, time.UTC) }

	rec := httptest.NewRecorder()
	h.HandleDaily(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	// The gap-filled window must include July 12 (today in Auckland). The
	// date string is the observable; asserting the day exists is behaviour,
	// not chrome.
	if !strings.Contains(rec.Body.String(), "Jul 12") {
		t.Fatal("window does not include 2026-07-12: today was named in UTC, not the Preferences timezone")
	}
}
