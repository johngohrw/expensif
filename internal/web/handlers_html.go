package web

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"expensif/internal/domain"
	"expensif/internal/service"
)

type HTMLHandler struct {
	svc      *service.Service
	renderer *Renderer
}

func NewHTMLHandler(svc *service.Service, renderer *Renderer) *HTMLHandler {
	return &HTMLHandler{svc: svc, renderer: renderer}
}

func nowInTZ(tz string) time.Time {
	loc := time.Local
	if tz != "" {
		if l, err := time.LoadLocation(tz); err == nil {
			loc = l
		}
	}
	return time.Now().In(loc)
}

func (h *HTMLHandler) basePageData(ctx context.Context, active string) PageData {
	prefs, err := h.svc.Preferences(ctx)
	if err != nil {
		slog.Error("failed to load preferences", "error", err)
		prefs = &domain.Preferences{Currency: "USD"}
	}
	return PageData{
		Active:         active,
		CurrencySymbol: domain.CurrencySymbol(prefs.Currency),
		Currency:       prefs.Currency,
		UserID:         prefs.UserID,
		Timezone:       prefs.Timezone,
		Today:          nowInTZ(prefs.Timezone).Format("2006-01-02"),
		Islands:        []string{"mobile-nav"},
	}
}

func (h *HTMLHandler) render(w http.ResponseWriter, name string, data PageData) {
	if err := h.renderer.Render(w, name, data); err != nil {
		slog.Error("render failed", "template", name, "error", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

func (h *HTMLHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	expenses, err := h.svc.ListExpenses(ctx, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := h.basePageData(ctx, "list")

	// Convert per-row and totals to preference currency
	convertedTotal, rateDate, _ := h.svc.ConvertExpensesTotal(ctx, expenses, data.Currency)
	data.ConvertedTotal = convertedTotal
	data.RateDate = rateDate
	data.ShowConverted = convertedTotal > 0

	catConverted := make(map[string]float64)
	rates, _, rateErr := h.svc.GetRatesForConversion(ctx)
	if rateErr == nil {
		for i := range expenses {
			conv, err := h.svc.ConvertWithRates(expenses[i].Amount, expenses[i].Currency, data.Currency, rates)
			if err == nil {
				expenses[i].ConvertedAmount = conv
				catConverted[expenses[i].Category] += conv
			}
		}
	}

	data.Expenses = expenses

	var categories []domain.CategorySummary
	for name, amount := range catConverted {
		categories = append(categories, domain.CategorySummary{Name: name, Amount: amount})
	}
	data.Categories = categories

	total, _ := h.svc.TotalExpenses(ctx)
	data.Total = total
	data.Islands = append(data.Islands, "data-table")

	h.render(w, "list", data)
}

func (h *HTMLHandler) HandleCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := h.basePageData(ctx, "calendar")

	earliest, _ := h.svc.GetEarliestExpenseDate(ctx)
	loc := time.Local
	if data.Timezone != "" {
		if l, err := time.LoadLocation(data.Timezone); err == nil {
			loc = l
		}
	}
	today := nowInTZ(data.Timezone)

	// Start = earliest expense or 1 year before today, whichever is earlier
	var start time.Time
	if earliest != "" {
		start, _ = time.Parse("2006-01-02", earliest)
	}
	oneYearAgo := today.AddDate(-1, 0, 0)
	if start.IsZero() || start.After(oneYearAgo) {
		start = oneYearAgo
	}
	start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, loc)

	// End = today + 1 year
	end := today.AddDate(1, 0, 0)
	end = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, loc)

	var months []domain.CalendarMonth
	for m := start; !m.After(end); m = m.AddDate(0, 1, 0) {
		var month domain.CalendarMonth
		month.Label = m.Format("Jan 2006")

		first := m
		last := m.AddDate(0, 1, -1)

		// Padding before (Sunday → first day)
		for d := 0; d < int(first.Weekday()); d++ {
			month.Cells = append(month.Cells, domain.CalendarCell{})
		}

		// Days
		for d := first; !d.After(last); d = d.AddDate(0, 0, 1) {
			isToday := d.Format("2006-01-02") == data.Today
			if isToday {
				month.HasToday = true
			}
			month.Cells = append(month.Cells, domain.CalendarCell{
				Day:     d.Day(),
				Date:    d.Format("2006-01-02"),
				IsToday: isToday,
			})
		}

		// Padding after (last day → Saturday)
		for d := int(last.Weekday()) + 1; d <= 6; d++ {
			month.Cells = append(month.Cells, domain.CalendarCell{})
		}

		months = append(months, month)
	}

	data.CalendarMonths = months
	h.render(w, "calendar", data)
}

func (h *HTMLHandler) HandleDaily(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	groups, err := h.svc.DailyGroups(ctx, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := h.basePageData(ctx, "daily")

	// Convert daily totals and per-row amounts
	rates, _, rateErr := h.svc.GetRatesForConversion(ctx)
	if rateErr == nil {
		for i := range groups {
			var convTotal float64
			for j := range groups[i].Expenses {
				conv, err := h.svc.ConvertWithRates(groups[i].Expenses[j].Amount, groups[i].Expenses[j].Currency, data.Currency, rates)
				if err == nil {
					groups[i].Expenses[j].ConvertedAmount = conv
					convTotal += conv
				}
			}
			groups[i].ConvertedTotal = convTotal
		}
	}

	data.DailyGroups = groups
	data.Islands = append(data.Islands, "data-table")
	h.render(w, "daily", data)
}

func (h *HTMLHandler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	date := r.URL.Query().Get("date")
	data := h.basePageData(ctx, "add")
	if date == "" {
		date = data.Today
	}
	data.Today = date
	users, err := h.svc.ListUsers(ctx)
	if err != nil {
		slog.Error("failed to list users", "error", err)
	}
	data.Users = users
	// Pre-select current user from preferences
	data.PaidByID = data.UserID
	data.Islands = append(data.Islands, "category-pills", "description-pills")
	h.render(w, "add", data)
}

func (h *HTMLHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	date := r.FormValue("date")
	category := r.FormValue("category")
	description := r.FormValue("description")
	currency := r.FormValue("currency")
	paidByID, _ := strconv.ParseInt(r.FormValue("paid_by_id"), 10, 64)

	_, err := h.svc.CreateExpense(ctx, amount, category, description, date, currency, paidByID)
	if err != nil {
		data := h.basePageData(ctx, "add")
		data.Flash = err.Error()
		data.FlashError = true
		data.Today = date
		data.Expense = &domain.Expense{
			Amount:      amount,
			Category:    category,
			Description: description,
			Date:        date,
			Currency:    currency,
			PaidByID:    paidByID,
		}
		users, _ := h.svc.ListUsers(ctx)
		data.Users = users
		data.Islands = append(data.Islands, "category-pills", "description-pills")
		h.render(w, "add", data)
		return
	}

	if r.FormValue("action") == "another" {
		http.Redirect(w, r, "/expenses/new", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *HTMLHandler) HandleEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	expense, err := h.svc.GetExpense(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	data := h.basePageData(ctx, "edit")
	data.Expense = expense
	users, err := h.svc.ListUsers(ctx)
	if err != nil {
		slog.Error("failed to list users", "error", err)
	}
	data.Users = users
	data.Islands = append(data.Islands, "category-pills", "description-pills")
	h.render(w, "edit", data)
}

func (h *HTMLHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	date := r.FormValue("date")
	category := r.FormValue("category")
	description := r.FormValue("description")
	currency := r.FormValue("currency")
	paidByID, _ := strconv.ParseInt(r.FormValue("paid_by_id"), 10, 64)

	err := h.svc.UpdateExpense(ctx, id, amount, category, description, date, currency, paidByID)
	if err != nil {
		data := h.basePageData(ctx, "edit")
		data.Expense = &domain.Expense{
			ID:          id,
			Amount:      amount,
			Category:    category,
			Description: description,
			Date:        date,
			Currency:    currency,
			PaidByID:    paidByID,
		}
		data.Flash = err.Error()
		data.FlashError = true
		users, _ := h.svc.ListUsers(ctx)
		data.Users = users
		data.Islands = append(data.Islands, "category-pills", "description-pills")
		h.render(w, "edit", data)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HTMLHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.svc.DeleteExpense(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HTMLHandler) HandlePreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := h.basePageData(ctx, "prefs")
	users, _ := h.svc.ListUsers(ctx)
	data.Users = users
	h.render(w, "prefs", data)
}

func (h *HTMLHandler) HandleSavePreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	currency := r.FormValue("currency")
	if currency == "" {
		currency = "USD"
	}
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	timezone := r.FormValue("timezone")

	if err := h.svc.SavePreferences(ctx, currency, userID, timezone); err != nil {
		data := h.basePageData(ctx, "prefs")
		data.Flash = "Failed to save preferences"
		data.FlashError = true
		users, _ := h.svc.ListUsers(ctx)
		data.Users = users
		h.render(w, "prefs", data)
		return
	}

	data := h.basePageData(ctx, "prefs")
	data.Flash = "Preferences saved"
	users, _ := h.svc.ListUsers(ctx)
	data.Users = users
	h.render(w, "prefs", data)
}

// --- User Management ---

func (h *HTMLHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.svc.ListUsers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := h.basePageData(ctx, "users")
	data.Users = users
	data.Islands = append(data.Islands, "data-table")
	h.render(w, "users", data)
}

func (h *HTMLHandler) HandleUserNew(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := h.basePageData(ctx, "users")
	h.render(w, "user_form", data)
}

func (h *HTMLHandler) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	_, err := h.svc.CreateUser(ctx, name)
	if err != nil {
		data := h.basePageData(ctx, "users")
		data.Flash = err.Error()
		data.FlashError = true
		h.render(w, "user_form", data)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *HTMLHandler) HandleUserEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	data := h.basePageData(ctx, "users")
	data.User = user
	h.render(w, "user_form", data)
}

func (h *HTMLHandler) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	name := r.FormValue("name")
	err := h.svc.UpdateUser(ctx, id, name)
	if err != nil {
		user, _ := h.svc.GetUser(ctx, id)
		data := h.basePageData(ctx, "users")
		data.User = user
		data.Flash = err.Error()
		data.FlashError = true
		h.render(w, "user_form", data)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *HTMLHandler) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.svc.DeleteUser(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
