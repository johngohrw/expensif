package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"

	"expensif/internal/assets"
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
	UserID         int64
	Users          []domain.User
	User           *domain.User
	PaidByID       int64
	Islands        []string // Names of React islands to hydrate on this page
}

type Renderer struct {
	templates map[string]*template.Template
}

func NewRenderer(templatesDir string, dev bool, manifest assets.Manifest) (*Renderer, error) {
	helper := &assets.AssetHelper{Dev: dev, Manifest: manifest}

	funcMap := template.FuncMap{
		"default": func(def, v interface{}) interface{} {
			switch val := v.(type) {
			case string:
				if val == "" {
					return def
				}
			case nil:
				return def
			}
			return v
		},
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
			todayStr := time.Now().Format("2006-01-02")
			if dateStr == todayStr {
				return "Today"
			}
			yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
			if dateStr == yesterdayStr {
				return "Yesterday"
			}
			return humanize.Time(t)
		},
		"formatDate": func(dateStr string) string {
			t, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
			if err != nil {
				return dateStr
			}
			return t.Format("Jan 2, Monday")
		},
		"currencySymbol": domain.CurrencySymbol,
		"script":    func(entry string) template.HTML { return helper.ScriptTag(entry) },
		"devClient": func() template.HTML { return helper.DevClient() },
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"json": func(v interface{}) (string, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
		"list": func(items ...interface{}) []interface{} {
			return items
		},
		"jsonSafe": func(v interface{}) template.JS {
			b, err := json.Marshal(v)
			if err != nil {
				return template.JS("null")
			}
			return template.JS(b)
		},
	}

	parsePage := func(files ...string) *template.Template {
		t := template.New("").Funcs(funcMap)
		allFiles := append([]string{filepath.Join(templatesDir, "base.html")}, files...)
		partialFiles, _ := filepath.Glob(filepath.Join(templatesDir, "partials", "*.html"))
		allFiles = append(allFiles, partialFiles...)
		return template.Must(t.ParseFiles(allFiles...))
	}

	templates := map[string]*template.Template{
		"list":      parsePage(filepath.Join(templatesDir, "list.html")),
		"daily":     parsePage(filepath.Join(templatesDir, "daily.html")),
		"add":       parsePage(filepath.Join(templatesDir, "form.html"), filepath.Join(templatesDir, "add.html")),
		"edit":      parsePage(filepath.Join(templatesDir, "form.html"), filepath.Join(templatesDir, "edit.html")),
		"prefs":     parsePage(filepath.Join(templatesDir, "preferences.html")),
		"users":     parsePage(filepath.Join(templatesDir, "users.html")),
		"user_form": parsePage(filepath.Join(templatesDir, "user_form.html")),
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
