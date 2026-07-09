# Expensif — Project Brief

## Overview

Expensif is a personal / small-household expense tracker. It is a Go web application that renders server-side HTML templates and hydrates selective **React Islands** for interactivity. Data persists in SQLite, exchange rates are fetched from `frankfurter.dev`, and the app is deployed as a single Docker image.

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Backend language | Go | 1.25.0 |
| Database | SQLite (pure Go, no CGO) | `modernc.org/sqlite` v1.50.0 |
| HTTP server | `net/http` | stdlib |
| Templates | Go `html/template` | stdlib |
| Frontend islands | React | 18.3.1 |
| Frontend language | TypeScript | 5.6.3 |
| Build tool | Vite | 5.4.11 |
| Styling | Tailwind CSS | via CDN |
| Date formatting | `github.com/dustin/go-humanize` | v1.0.1 |
| Exchange rates | Frankfurter API | `api.frankfurter.dev/v1` |

## Directory Structure

```
.
├── cmd/server/main.go              # Entry point: DB, service, server, background rate fetch
├── internal/
│   ├── assets/assets.go            # Vite manifest loader + dev/prod script helpers
│   ├── db/db.go                    # SQLite connection + migrations
│   ├── domain/                     # Domain structs + currency helpers
│   │   ├── models.go
│   │   └── currency.go
│   ├── rate/client.go              # Frankfurter HTTP client
│   ├── repository/                 # Repository interfaces + SQLite impl
│   │   ├── repository.go
│   │   └── sqlite.go
│   ├── service/service.go          # Business logic, validation, conversion
│   └── web/                        # HTTP handlers, server, renderer, middleware
│       ├── server.go
│       ├── middleware.go
│       ├── handlers_html.go
│       ├── handlers_api.go
│       ├── handlers_api_test.go
│       ├── server_test.go
│       ├── mock_repo_test.go
│       └── renderer.go
├── templates/                      # Go HTML templates + partials
│   ├── base.html
│   ├── daily.html
│   ├── list.html
│   ├── calendar.html
│   ├── form.html
│   ├── add.html / edit.html
│   ├── preferences.html
│   ├── users.html / user_form.html
│   └── partials/
├── ui/                             # React/TypeScript source
│   ├── src/
│   │   ├── components/
│   │   ├── entries/
│   │   ├── lib/
│   │   └── types/
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
├── static/                         # Vite build output (gitignored)
├── go.mod / go.sum
├── Dockerfile / docker-compose.yml
├── Makefile / .air.toml
└── .env.example
```

## Database Schema

SQLite file: `~/.expensif/expenses.db` (or `$DATA_DIR/expenses.db` in Docker).

```sql
CREATE TABLE expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    category TEXT NOT NULL,
    description TEXT,         -- service requires non-empty, DB allows NULL
    date TEXT,                -- YYYY-MM-DD
    currency TEXT DEFAULT 'USD',
    paid_by INTEGER,          -- nullable FK to users.id
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE preferences (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    currency TEXT DEFAULT 'USD',
    user_id INTEGER,          -- default payer for new expenses
    timezone TEXT DEFAULT ''  -- IANA timezone, e.g. 'Asia/Singapore'
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

**Migrations:** All schema changes are handled in `internal/db/db.go`. Columns are added incrementally with `pragma_table_info` checks.

## Domain Models

```go
type Expense struct {
    ID              int64
    Amount          float64
    Category        string
    Description     string
    Date            string // YYYY-MM-DD
    Currency        string
    PaidByID        int64
    PaidByName      string // populated by repo via LEFT JOIN users
    CreatedAt       time.Time
    ConvertedAmount float64 // computed at render time for display
}

type Preferences struct {
    Currency string
    UserID   int64
    Timezone string
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

type CalendarMonth struct {
    Label    string
    HasToday bool
    Cells    []CalendarCell
}

type DescriptionCount struct {
    Description string
    Count       int
}
```

## Routes

### HTML Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/` | `HandleDaily` | Daily grouped view (homepage) |
| GET | `/calendar` | `HandleCalendar` | Monthly calendar heatmap |
| GET | `/expenses` | `HandleList` | All expenses table with totals |
| GET | `/expenses/new` | `HandleAdd` | Add form (supports `?date=` prefill) |
| POST | `/expenses/new` | `HandleCreate` | Create (`action=another` → redirect to add again) |
| GET | `/expenses/edit/{id}` | `HandleEdit` | Edit form |
| POST | `/expenses/edit/{id}` | `HandleUpdate` | Update |
| POST | `/expenses/delete/{id}` | `HandleDelete` | Delete (JS confirm) |
| GET | `/preferences` | `HandlePreferences` | Settings page |
| POST | `/preferences` | `HandleSavePreferences` | Save settings |
| GET | `/users` | `HandleUsers` | User management |
| GET | `/users/new` | `HandleUserNew` | New user form |
| POST | `/users/new` | `HandleUserCreate` | Create user |
| GET | `/users/edit/{id}` | `HandleUserEdit` | Edit user form |
| POST | `/users/edit/{id}` | `HandleUserUpdate` | Update user |
| POST | `/users/delete/{id}` | `HandleUserDelete` | Delete user |

### JSON API Routes

| Method | Path | Handler |
|--------|------|---------|
| GET | `/api/expenses?n=...` | `HandleList` (default limit: 50) |
| POST | `/api/expenses` | `HandleCreate` |
| GET | `/api/expenses/{id}` | `HandleGet` |
| PUT | `/api/expenses/{id}` | `HandleUpdate` |
| DELETE | `/api/expenses/{id}` | `HandleDelete` |
| GET | `/api/categories` | `HandleCategories` |
| GET | `/api/expenses/descriptions?category=...` | `HandleDescriptions` |
| GET | `/api/summary` | `HandleSummary` |

## Features

### Expense CRUD
- Create, list (limit 100 in HTML), get by ID, update, delete
- Amount, category, description, date, currency, payer per expense
- Validation: amount > 0, category and description required
- Default date = "today" in user's timezone; default currency = USD

### Users & Payers
- Full CRUD for users (`/users`)
- "Paid By" field on expenses links to a user via `paid_by INTEGER`
- Preferences stores `user_id` as default payer for new expenses
- User deletion clears `paid_by` on linked expenses inside a transaction

### Multi-Currency
- 17 supported currencies: USD, MYR, JPY, CNY, THB, EUR, GBP, SGD, KRW, AUD, CAD, INR, VND, PHP, IDR, HKD, TWD
- Each expense has its own currency
- Rates fetched from Frankfurter, stored in SQLite, refreshed every 6 hours
- Cross-rate conversion via USD base: `amount * toRate / fromRate`
- Converted totals shown in List and Daily views when rates are available

### Views
- **Daily** (`/`): expenses grouped by date, newest first, per-day totals
- **All Expenses** (`/expenses`): full table with converted totals and category summary
- **Calendar** (`/calendar`): monthly grid from (earliest expense | 1yr ago) to (today + 1yr), sticky month labels, today badge, auto-scroll to today, heat-mapped spend blobs

### Interactive Islands
- **Category pills** on add/edit form: quick-select from recent categories
- **Description pills**: dynamically fetched top descriptions for selected category
- **DataTable**: column-driven, server-configured table rendering
- **Mobile nav**: hamburger side-panel on narrow viewports

### Preferences
- Singleton row (`id = 1`)
- Preferred currency, default payer, and IANA timezone

### UI/UX
- Tailwind CSS via CDN
- Sticky top navigation
- Flash messages for errors/success
- Shared `form.html` for add/edit with "Create & add another" option
- Clear buttons injected on text/number inputs
- Mobile content padding wrapper (`px-3 sm:px-0`)

### Middleware
- **RecoverPanic:** catches panics, logs stack, returns 500
- **RequestLog:** logs method, path, status, duration

### Graceful Shutdown
- Listens for `SIGINT`/`SIGTERM`
- Cancels background rate-fetcher goroutine
- Calls `server.Shutdown()` with a 5-second timeout

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
Repository.Repos (Expense, User, Preference, Rate)
    ↓
SQLite Repo (SQL queries, transactions)
    ↓
SQLite DB (~/.expensif/expenses.db)
```

## PageData / Template Data Model

`PageData` (in `renderer.go`) is the shared data struct for all templates:
- `Active`, `Flash`, `FlashError`
- `Expenses`, `Expense`, `DailyGroups`, `CalendarMonths`
- `Total`, `ConvertedTotal`, `RateDate`, `ShowConverted`, `Categories []CategorySummary`
- `Currency`, `CurrencySymbol`, `UserID`, `Users`, `User`, `PaidByID`
- `Today`, `Timezone`, `FilterDate`, `BackHref`, `BackLabel`
- `Islands []string` — names of React island entry points to hydrate on this page

## Template FuncMap

- `dict` — key-value map builder
- `list` — variadic list builder
- `default` — fallback for empty string/nil
- `humanDate` — "Today"/"Yesterday"/humanize
- `formatDate` — "Jan 2, Monday"
- `currencySymbol` — maps code to symbol
- `formatAmount` — format using currency decimal places
- `currencyDecimalsJSON` — JSON map of currency decimal places
- `script` — inject island `<script>` tag (dev or prod)
- `devClient` — inject Vite HMR client in development
- `safeHTML` — render trusted HTML string
- `json` / `jsonSafe` — JSON encoding for template data

## Testing

- **Go**: 40 tests total
  - 20 API handler tests
  - 2 server lifecycle tests
  - 12 service-layer tests (validation, conversion, rate refresh)
  - In-memory mock repo (`mockRepo`) satisfies the full `Repository` interface
- **Frontend**: 7 tests for `DataTable` in `ui/src/components/DataTable/DataTable.test.tsx`

Run:
```bash
go test ./...
go vet ./...
cd ui && npx vitest run
```

## Current State / Open Items

- **No HTML handler tests** — the primary form/redirect/template surface is untested.
- **No repository SQLite tests** — the mock repo does not exercise real SQL.
- **Currency conversion lives in HTML handlers** (`HandleList`, `HandleDaily`) — should move to the service layer.
- **`ConvertedAmount` on `domain.Expense`** is a presentation leak.
- **`SummaryByCategory` sums raw amounts** across mixed currencies without conversion.
- **`rate.Client` is now injectable** via the `service.RateFetcher` interface (fixed).
- **Form values are preserved** on validation errors (fixed).
- **User deletion** clears `paid_by` atomically (fixed).

## Running the App

```bash
# Development
make dev

# Or manually:
# Terminal 1: DEV=true go run ./cmd/server
# Terminal 2: cd ui && npm run dev
# Access http://localhost:8080

# Production
make prod
./bin/server
```

Environment variables:
- `PORT` — defaults to 8080
- `DATA_DIR` — defaults to `~/.expensif`
- `DEV` — development mode (Vite HMR)
- `VITE_DEV_HOST` — Vite dev server address (default `localhost:8081`)
- `TZ` — server timezone fallback
