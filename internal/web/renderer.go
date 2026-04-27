package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"

	"expensif/internal/domain"
)

type PageData struct {
	Active         string
	Flash          string
	FlashError     bool
	Expenses       []domain.Expense
	Expense        *domain.Expense
	Total          float64
	Categories     []domain.CategorySummary
	Today          string
	DailyGroups    []domain.DailyGroup
	CurrencySymbol string
	Currency       string
	ConvertedTotal float64
	RateDate       string
	ShowConverted  bool
}

type Renderer struct {
	templates map[string]*template.Template
}

func NewRenderer(templatesDir string) (*Renderer, error) {
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
		"currencySymbol": domain.CurrencySymbol,
	}

	parsePage := func(files ...string) *template.Template {
		t := template.New("").Funcs(funcMap)
		allFiles := append([]string{filepath.Join(templatesDir, "base.html")}, files...)
		return template.Must(t.ParseFiles(allFiles...))
	}

	templates := map[string]*template.Template{
		"list":  parsePage(filepath.Join(templatesDir, "list.html")),
		"daily": parsePage(filepath.Join(templatesDir, "daily.html")),
		"add":   parsePage(filepath.Join(templatesDir, "form.html"), filepath.Join(templatesDir, "add.html")),
		"edit":  parsePage(filepath.Join(templatesDir, "form.html"), filepath.Join(templatesDir, "edit.html")),
		"prefs": parsePage(filepath.Join(templatesDir, "preferences.html")),
	}

	return &Renderer{templates: templates}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data PageData) error {
	t, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteTemplate(w, "base", data)
}
