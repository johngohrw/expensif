# Expensif — Project Brief

## Overview

Expensif is a Go-based expense tracking web application with server-rendered HTML UI and a JSON API. It is a single binary with no JavaScript framework — vanilla JS only (Tailwind CSS via CDN). Data persists in SQLite.

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Language | Go | 1.25 |
| Database | SQLite (pure Go, no CGO) | `modernc.org/sqlite` v1.50.0 |
| CSS | Tailwind CSS | via CDN |
| Templates | Go `html/template` | stdlib |
| Date formatting | `github.com/dustin/go-humanize` | v1.0.1 |
| Exchange rates | Frankfurter API | `api.frankfurter.dev/v1` |

## Directory Structure

```
.
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── db/db.go                    # SQLite init + migrations
│   ├── domain/
│   │   ├── models.go               # Expense, Preferences, User, CategorySummary, DailyGroup
│   │   └── currency.go             # CurrencySymbol(code) → "$")
│   ├── rate/client.go              # Frankfurter HTTP client
│   ├── repository/
│   │   ├── repository.go           # Repository interface
│   │   └── sqlite.go               # SQLite implementation
│   ├── service/service.go          # Business logic, validation, conversions
│   └── web/
│       ├── server.go               # HTTP routing, Server struct
│       ├── middleware.go           # RecoverPanic + RequestLog middleware
│       ├── handlers_html.go        # HTML page handlers
│       ├── handlers_api.go         # JSON API handlers
│       ├── handlers_api_test.go    # 18 API tests
│       ├── mock_repo_test.go       # In-memory mock Repository
│       ├── server_test.go          # 2 Server lifecycle tests
│       └── renderer.go             # Template parsing + PageData
├── templates/
│   ├── base.html                   # Layout, nav, flash messages
│   ├── form.html                   # Shared add/edit expense form
│   ├── add.html / edit.html        # Thin wrappers around form.html
│   ├── list.html                   # Expenses table with summary
│   ├── daily.html                  # Timeline view grouped by day
│   ├── preferences.html            # Currency + default user
│   ├── users.html                  # User management list
│   └── user_form.html              # Create/edit user form
├── go.mod / go.sum
└── .gitignore
```

## Database Schema

SQLite file: `~/.expensif/expenses.db`

```sql
CREATE TABLE expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    category TEXT NOT NULL,
    description TEXT,          -- service requires non-empty, but DB allows NULL
    date TEXT,
    currency TEXT DEFAULT 'USD',
    paid_by INTEGER,           -- nullable FK to users.id (no FK constraint enforced)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- format: "YYYY-MM-DD HH:MM:SS"
);

CREATE TABLE preferences (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    currency TEXT DEFAULT 'USD',
    user_id INTEGER             -- nullable; default user for new expenses
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE exchange_rates (
    base_currency TEXT NOT NULL,
    target_currency TEXT NOT NULL,
    rate REAL NOT NULL,
    date TEXT NOT NULL,
    fetched_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (base_currency, target_currency, date)
);

CREATE INDEX idx_expenses_category ON expenses(category);
CREATE INDEX idx_expenses_created_at ON expenses(created_at);
```

**Migrations:** All schema changes are handled via `internal/db/db.go` migrations. Columns added incrementally via `ALTER TABLE` with `pragma_table_info` checks. The `paid_by` column had a TEXT→INTEGER migration that rebuilds the table.

## Domain Models

```go
type Expense struct {
    ID              int64
    Amount          float64
    Category        string
    Description     string
    Date            string // YYYY-MM-DD
    Currency        string
    PaidByID        int64  `json:"paid_by_id,omitempty"`
    PaidByName      string `json:"paid_by_name,omitempty"` // populated by repo via JOIN
    CreatedAt       time.Time
    ConvertedAmount float64 `json:"-"` // computed at render time
}

type Preferences struct {
    Currency string
    UserID   int64
}

type User struct {
    ID   int64
    Name string
}

type CategorySummary struct {
    Name   string
    Amount float64
}

type DailyGroup struct {
    Date           string
    Expenses       []Expense
    Total          float64
    ConvertedTotal float64
}
```

## Routes

### HTML Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/` | HandleList | Expenses list with totals |
| GET | `/daily` | HandleDaily | Daily timeline view |
| GET | `/preferences` | HandlePreferences | Settings page |
| POST | `/preferences` | HandleSavePreferences | Save settings |
| GET | `/expenses/new` | HandleAdd | Add form (supports `?date=` prefill) |
| POST | `/expenses/new` | HandleCreate | Create (`action=another` → redirect to add again) |
| GET | `/expenses/edit/{id}` | HandleEdit | Edit form |
| POST | `/expenses/edit/{id}` | HandleUpdate | Update |
| POST | `/expenses/delete/{id}` | HandleDelete | Delete (JS confirm) |
| GET | `/users` | HandleUsers | User management |
| GET | `/users/new` | HandleUserNew | New user form |
| POST | `/users/new` | HandleUserCreate | Create user |
| GET | `/users/edit/{id}` | HandleUserEdit | Edit user form |
| POST | `/users/edit/{id}` | HandleUserUpdate | Update user |
| POST | `/users/delete/{id}` | HandleUserDelete | Delete user |

### JSON API Routes
| Method | Path | Handler |
|--------|------|---------|
| GET | `/api/expenses` | HandleList (default limit: 50) |
| POST | `/api/expenses` | HandleCreate |
| GET | `/api/expenses/{id}` | HandleGet |
| PUT | `/api/expenses/{id}` | HandleUpdate |
| DELETE | `/api/expenses/{id}` | HandleDelete |
| GET | `/api/categories` | HandleCategories |
| GET | `/api/summary` | HandleSummary |

Query param on `/api/expenses`: `?n=10` to override default limit.

## Features

### Expense CRUD
- Create, list (limit 100 in HTML), get by ID, update, delete
- Amount, category, description, date, currency per expense
- Category autosuggest via `/api/categories` (frequency-sorted pills in last 3 months)

### Currency
- **17 supported currencies:** USD, MYR, JPY, CNY, THB, EUR, GBP, SGD, KRW, AUD, CAD, INR, VND, PHP, IDR, HKD, TWD
- **Per-transaction currency:** each expense has its own currency
- **Conversion:** rates fetched from Frankfurter API, stored in SQLite
- **Background refresh:** every 6 hours, starting immediately on boot
- **Cross-rate math:** via USD base (`amount * toRate / fromRate`)

### Humanized Dates
- List view shows "Today", "Yesterday", then `humanize.Time()` fallback
- Raw date in HTML `title` attribute

### Daily View
- Expenses grouped by day, newest first
- Daily total (converted when rates available)
- "+ Add expense" link pre-filled with that day's date

### User System
- Full CRUD for users (`/users`)
- **"Paid By" field** on expenses — links to user via `paid_by INTEGER`
- Preferences stores `user_id` as default for new expenses
- On user delete, `paid_by` on linked expenses nulled (atomic transaction in SQLite repo)
- `PaidByName` populated via `LEFT JOIN users` at repository read time

### Preferences
- Currency dropdown + default user dropdown
- Singleton row (`id = 1`)

### UI/UX
- Tailwind CSS via CDN
- Three-dot dropdown menu for Preferences
- "+ New Expense" as primary blue button
- Shared `form.html` for add/edit with "Create & add another" option
- Flash messages for errors/success
- Form buttons: Cancel | Save / Cancel | Create | Create & add another

### Middleware
- **RecoverPanic:** catches panics, logs stack trace, returns 500
- **RequestLog:** logs method, path, status code, duration

### Graceful Shutdown
- Listens for `SIGINT`/`SIGTERM`
- Cancels background context (stops rate fetcher)
- Calls `server.Shutdown()` with 5-second timeout

## Testing

- **20 tests** total in `internal/web/`
- **18 API tests** covering all JSON endpoints (CRUD, validation, categories, summary)
- **2 Server tests** (`TestServerShutdownBeforeRun`, `TestServerRunAndShutdown`)
- **In-memory mock repo** (`mockRepo`) satisfies full `Repository` interface with thread-safe operations

Run: `go test ./...`

## Key Design Decisions

1. **No SPA / No JS framework** — pure server-rendered HTML, vanilla JS only (fetch for category pills)
2. **SQLite in home dir** — `~/.expensif/expenses.db`
3. **Preferences as singleton** — `id = 1` constraint
4. **No foreign key constraint** on `expenses.paid_by` → `users.id` — manual cleanup via transaction on delete
5. **float64 for amounts** — accepted rounding risk (display rounds to 2 decimals)
6. **Cross-rate via USD** — all conversions go through USD base rates
7. **Rate fetching** — only USD base rates fetched; other bases not supported

## Architecture Layers

```
HTTP Request
    ↓
Middleware (RecoverPanic → RequestLog)
    ↓
Server (routing)
    ↓
HTMLHandler / APIHandler
    ↓
Service (validation, business logic, conversion)
    ↓
Repository (interface)
    ↓
SQLite Repo (SQL queries, transactions)
    ↓
SQLite DB (~/.expensif/expenses.db)
```

## PageData / Template Data Model

`PageData` (in `renderer.go`) is the shared data struct for all templates:
- `Active`, `Flash`, `FlashError`
- `Expenses`, `Expense`, `DailyGroups`
- `Total`, `ConvertedTotal`, `RateDate`, `ShowConverted`
- `Categories []CategorySummary`
- `Currency`, `CurrencySymbol`, `UserID`, `Users`, `User`
- `Today`, `PaidByID`

Templates: `base.html` wraps `title` + `content` blocks. `form.html` is shared via `dict` helper.

## Template FuncMap
- `dict` — key-value map builder
- `humanDate` — "Today"/"Yesterday"/humanize
- `currencySymbol` — maps code to symbol

## Accepted Tradeoffs / Known Issues

| Issue | Decision |
|-------|----------|
| `float64` for money | Skipped — rounding to 2 decimals is acceptable for personal use |
| No CSRF tokens | Acceptable for local/single-user app |
| No auth | Single-user local app |
| Rate fetch only USD base | Acceptable — most currencies convertible via USD |
| `log.Fatalf` on startup errors | Skips `defer database.Close()` — acceptable for fatal startup |

## What Changed Recently (April 28 PM Session)

1. **Fixed `CreatedAt` parsing** — SQLite timestamp format now correctly parsed
2. **Atomic `DeleteUser`** — single transaction, removed `ClearExpensePaidBy` from interface
3. **Middleware** — panic recovery + request logging
4. **`PaidByName` via JOIN** — eliminated N+1 query, removed `HydratePaidByNames`
5. **Graceful shutdown** — signal handling, cancellable background goroutine, `Server.Shutdown()`
6. **Cleanup** — removed dead `Server.mux` field, fixed `.gitignore` duplicate

## Running the App

```bash
cd /Users/rengwu/Desktop/Projects/expensif
go run ./cmd/server
# or build:
go build -o expensif ./cmd/server
./expensif
```

Environment: `PORT` defaults to 8080.

## Next Session Notes

- The `.context/` folder contains session summaries and this brief. The `golang-react-islands-spec-2026-04-28.md` is a spec doc for a potential React Islands migration — **not yet implemented**.
- All tests pass. `go vet` is clean. Build succeeds.
- The codebase is in a stable, well-factored state.
- If adding new features, consider: splitting the fat `Repository` interface, adding `WithTx` support, or adding more HTML handler tests.
