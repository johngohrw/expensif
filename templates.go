package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
)

type CategorySummary struct {
	Name   string
	Amount float64
}

type DailyGroup struct {
	Date     string
	Expenses []Expense
	Total    float64
}

type PageData struct {
	Active         string
	Flash          string
	FlashError     bool
	Expenses       []Expense
	Expense        *Expense
	Total          float64
	Categories     []CategorySummary
	Today          string
	DailyGroups    []DailyGroup
	DarkMode       bool
	CurrencySymbol string
	Currency       string
}

var pageTemplates map[string]*template.Template

func parseTemplates() {
	funcMap := template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"humanDate": func(dateStr string) string {
			t, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return dateStr
			}
			today := time.Now().Truncate(24 * time.Hour)
			date := t.Truncate(24 * time.Hour)
			diff := int(today.Sub(date).Hours() / 24)
			switch diff {
			case 0:
				return "Today"
			case 1:
				return "Yesterday"
			default:
				return humanize.Time(t)
			}
		},
		"currencySymbol": currencySymbol,
	}

	parsePage := func(files ...string) *template.Template {
		t := template.New("").Funcs(funcMap)
		allFiles := append([]string{"templates/base.html"}, files...)
		return template.Must(t.ParseFiles(allFiles...))
	}

	pageTemplates = map[string]*template.Template{
		"list":  parsePage("templates/list.html"),
		"daily": parsePage("templates/daily.html"),
		"add":   parsePage("templates/form.html", "templates/add.html"),
		"edit":  parsePage("templates/form.html", "templates/edit.html"),
		"prefs": parsePage("templates/preferences.html"),
	}
}

func basePageData(active string) PageData {
	prefs := getPreferences()
	return PageData{
		Active:         active,
		DarkMode:       prefs.DarkMode,
		CurrencySymbol: currencySymbol(prefs.Currency),
		Currency:       prefs.Currency,
		Today:          time.Now().Format("2006-01-02"),
	}
}

func renderPage(w http.ResponseWriter, name string, data PageData) {
	t, ok := pageTemplates[name]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleListHTML(w http.ResponseWriter, r *http.Request) {
	expenses, err := listExpenses(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, _ := totalExpenses()
	byCat, _ := summaryByCategory()

	var categories []CategorySummary
	for name, amount := range byCat {
		categories = append(categories, CategorySummary{Name: name, Amount: amount})
	}

	data := basePageData("list")
	data.Expenses = expenses
	data.Total = total
	data.Categories = categories
	renderPage(w, "list", data)
}

func handleDailyHTML(w http.ResponseWriter, r *http.Request) {
	expenses, err := listExpenses(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	groups := make(map[string][]Expense)
	for _, e := range expenses {
		groups[e.Date] = append(groups[e.Date], e)
	}

	var dailyGroups []DailyGroup
	for date, exps := range groups {
		var dayTotal float64
		for _, e := range exps {
			dayTotal += e.Amount
		}
		dailyGroups = append(dailyGroups, DailyGroup{
			Date:     date,
			Expenses: exps,
			Total:    dayTotal,
		})
	}

	// Sort descending by date
	for i := 0; i < len(dailyGroups); i++ {
		for j := i + 1; j < len(dailyGroups); j++ {
			if dailyGroups[i].Date < dailyGroups[j].Date {
				dailyGroups[i], dailyGroups[j] = dailyGroups[j], dailyGroups[i]
			}
		}
	}

	data := basePageData("daily")
	data.DailyGroups = dailyGroups
	renderPage(w, "daily", data)
}

func handleAddHTML(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	data := basePageData("add")
	data.Today = date
	renderPage(w, "add", data)
}

func handleCreateHTML(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	date := r.FormValue("date")
	category := r.FormValue("category")
	description := r.FormValue("description")
	currency := r.FormValue("currency")

	if amount <= 0 || category == "" || description == "" {
		data := basePageData("add")
		data.Flash = "Amount, description, and category are required"
		data.FlashError = true
		renderPage(w, "add", data)
		return
	}

	_, err := addExpense(amount, category, description, date, currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("action") == "another" {
		http.Redirect(w, r, "/expenses/new", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handleEditHTML(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	expense, err := getExpense(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	data := basePageData("edit")
	data.Expense = expense
	renderPage(w, "edit", data)
}

func handleUpdateHTML(w http.ResponseWriter, r *http.Request) {
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

	if amount <= 0 || category == "" || description == "" {
		expense, _ := getExpense(id)
		data := basePageData("edit")
		data.Expense = expense
		data.Flash = "Amount, description, and category are required"
		data.FlashError = true
		renderPage(w, "edit", data)
		return
	}

	if err := updateExpense(id, amount, category, description, date, currency); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleDeleteHTML(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := deleteExpense(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlePreferencesHTML(w http.ResponseWriter, r *http.Request) {
	data := basePageData("prefs")
	renderPage(w, "prefs", data)
}

func handleSavePreferences(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	currency := r.FormValue("currency")
	if currency == "" {
		currency = "USD"
	}
	darkMode := r.FormValue("dark_mode") == "1"

	if err := savePreferences(currency, darkMode); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := basePageData("prefs")
	data.Flash = "Preferences saved"
	renderPage(w, "prefs", data)
}
