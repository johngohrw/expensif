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
		Today:          time.Now().Format("2006-01-02"),
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

	h.render(w, "list", data)
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
	h.render(w, "daily", data)
}

func (h *HTMLHandler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	data := h.basePageData(ctx, "add")
	data.Today = date
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

	_, err := h.svc.CreateExpense(ctx, amount, category, description, date, currency)
	if err != nil {
		data := h.basePageData(ctx, "add")
		data.Flash = err.Error()
		data.FlashError = true
		data.Today = date
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

	err := h.svc.UpdateExpense(ctx, id, amount, category, description, date, currency)
	if err != nil {
		expense, _ := h.svc.GetExpense(ctx, id)
		data := h.basePageData(ctx, "edit")
		data.Expense = expense
		data.Flash = err.Error()
		data.FlashError = true
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

	if err := h.svc.SavePreferences(ctx, currency); err != nil {
		data := h.basePageData(ctx, "prefs")
		data.Flash = "Failed to save preferences"
		data.FlashError = true
		h.render(w, "prefs", data)
		return
	}

	data := h.basePageData(ctx, "prefs")
	data.Flash = "Preferences saved"
	h.render(w, "prefs", data)
}
