package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

// TestHandleDaily_TimelineAndSingleDayShareDayEntryPartial is the spec's one
// "same partial" assertion: the timeline and the ?date= view both draw a
// populated day through the day-entry partial. data-day is the partial's
// machine-readable marker — nothing else emits it — so its presence in both
// responses is evidence of one implementation, without asserting on chrome.
func TestHandleDaily_TimelineAndSingleDayShareDayEntryPartial(t *testing.T) {
	repo := newMockRepo()
	repo.prefs = domain.Preferences{Currency: "USD", Timezone: "UTC"}
	repo.CreateExpense(t.Context(), domain.Expense{Amount: 5, Category: "food", Description: "x", Date: "2026-07-05", Currency: "USD"})

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
	h.now = func() time.Time { return time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC) }

	for _, target := range []string{"/", "/?date=2026-07-05"} {
		rec := httptest.NewRecorder()
		h.HandleDaily(rec, httptest.NewRequest("GET", target, nil))
		if rec.Code != 200 {
			t.Fatalf("%s: expected 200, got %d: %s", target, rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), `data-day="2026-07-05"`) {
			t.Fatalf("%s: the day did not render through the day-entry partial", target)
		}
	}
}

// TestHandleDelete_ReturnPath covers the danger zone's return path. The value
// arrives in a hidden form field, so it is attacker-supplied in principle: a
// delete that honoured it blindly would be an open redirect off the back of a
// destructive POST.
func TestHandleDelete_ReturnPath(t *testing.T) {
	tests := []struct {
		name   string
		form   string
		wantTo string
	}{
		{"local path is honoured", "return=/daily?date=2026-07-01", "/daily?date=2026-07-01"},
		{"absent falls back to root", "", "/"},
		{"empty falls back to root", "return=", "/"},
		{"protocol-relative host is rejected", "return=" + url.QueryEscape("//evil.example"), "/"},
		{"absolute URL with a scheme is rejected", "return=" + url.QueryEscape("https://evil.example/x"), "/"},
		{"backslash-folded host is rejected", "return=" + url.QueryEscape(`/\evil.example`), "/"},
		{"bare relative path is rejected", "return=" + url.QueryEscape("evil.example"), "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepo()
			id, err := repo.CreateExpense(t.Context(), domain.Expense{
				Amount: 5, Category: "food", Description: "x", Date: "2026-07-01", Currency: "USD",
			})
			if err != nil {
				t.Fatalf("seed expense: %v", err)
			}

			svc := service.New(repository.Repos{
				Expenses:    repo,
				Users:       repo,
				Preferences: repo,
				Rates:       repo,
			}, &mockRateClient{})
			h := NewHTMLHandler(svc, nil)

			req := httptest.NewRequest("POST", "/expenses/delete/"+strconv.FormatInt(id, 10), strings.NewReader(tt.form))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetPathValue("id", strconv.FormatInt(id, 10))

			rec := httptest.NewRecorder()
			h.HandleDelete(rec, req)

			if rec.Code != http.StatusSeeOther {
				t.Fatalf("expected 303, got %d: %s", rec.Code, rec.Body.String())
			}
			if got := rec.Header().Get("Location"); got != tt.wantTo {
				t.Fatalf("redirected to %q, want %q", got, tt.wantTo)
			}
			// The delete itself must happen either way — a rejected return
			// path redirects home, it does not veto the destructive action.
			if _, err := repo.GetExpense(t.Context(), id); err == nil {
				t.Fatal("expense survived the delete")
			}
		})
	}
}
