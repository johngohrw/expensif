# Expensif — Project Specification (2026-06-02)

## 1. Overview

**Expensif** is a personal expense tracker designed for individual or small-household use. It tracks expenses with categories, descriptions, currencies, and payers, supports multi-currency conversion via live exchange rates, and provides daily, list, and calendar views. The app is deployed via Docker on Unraid.

The architecture follows a **Go + React Islands** pattern: server-rendered HTML templates with selective React hydration for interactive components (category pills, description pills, data tables, mobile navigation). This keeps the initial page load lightweight while enabling rich client-side interactivity where needed.

---

## 2. Stack

| Layer | Technology | Version | Notes |
|-------|-----------|---------|-------|
| Backend language | Go | 1.25.0 | Standard library HTTP (`net/http`) |
| Database | SQLite | 3.x (via modernc.org/sqlite) | Embedded, zero-config |
| Frontend framework | React | 18.3.1 | Islands architecture — not SPA |
| Frontend language | TypeScript | 5.6.3 | Strict mode |
| Build tool | Vite | 5.4.11 | Builds to `static/`, generates manifest.json |
| Styling | Tailwind CSS | 3.x (CDN) | Loaded via CDN, ~30KB overhead |
| Templating | Go `html/template` | stdlib | Server-side rendering |
| Exchange rates | frankfurter.dev | v1 API | Free, no API key required |
| Containerization | Docker | Multi-stage | Node → Go → Alpine runtime |
| CI/CD | GitHub Actions | — | Builds linux/amd64, pushes to GHCR |
| Hot reload | Air | — | Go auto-rebuild on save |
| Testing (Go) | `testing` + stdlib | — | 40 tests, table-driven |
| Testing (UI) | Vitest + RTL | 4.1.5 / 16.3.2 | 7 tests, jsdom environment |

---

## 3. Directory Structure

```
expensif/
├── cmd/server/              # Go entry point
│   └── main.go              # DB init, service wiring, server start, graceful shutdown
├── internal/
│   ├── assets/
│   │   └── assets.go        # Vite manifest loader, ScriptTag + DevClient helpers
│   ├── db/
│   │   └── db.go            # SQLite connection + migrations
│   ├── domain/
│   │   ├── models.go        # All domain structs (Expense, User, Preferences, etc.)
│   │   └── currency.go      # Currency symbol map
│   ├── rate/
│   │   └── client.go        # frankfurter.dev HTTP client
│   ├── repository/
│   │   ├── repository.go    # Repository interfaces (split: Expense, User, Preference, Rate)
│   │   └── sqlite.go        # SQLite implementation of all interfaces
│   ├── service/
│   │   └── service.go       # Business logic: validation, defaults, conversion
│   └── web/
│       ├── handlers_api.go      # JSON API handlers (REST)
│       ├── handlers_api_test.go # 20 API handler tests
│       ├── handlers_html.go     # HTML page handlers (server-rendered)
│       ├── middleware.go        # RecoverPanic + RequestLog
│       ├── mock_repo_test.go    # In-memory mock repo for tests
│       ├── renderer.go          # Template engine setup + helpers
│       ├── server.go            # HTTP mux + route registration
│       └── server_test.go       # Server lifecycle tests
├── templates/               # Go HTML templates
│   ├── base.html            # Layout: nav, flash messages, island scripts
│   ├── daily.html           # Daily grouped view (default homepage)
│   ├── list.html            # All expenses table
│   ├── calendar.html        # Monthly calendar grid
│   ├── form.html            # Shared expense form (amount, category, description, currency, date, paid by)
│   ├── add.html             # New expense wrapper
│   ├── edit.html            # Edit expense wrapper
│   ├── preferences.html     # Settings: user, currency, timezone
│   ├── users.html           # User list + DataTable
│   ├── user_form.html       # New/edit user form
│   └── partials/            # Reusable components
│       ├── button.html      # Shared button partial (6 variants)
│       ├── data-table.html  # DataTable island container
│       ├── mobile-pad.html  # Mobile horizontal padding wrapper
│       └── page-title.html  # Page heading with optional action
├── ui/                      # Frontend source
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── src/
│       ├── components/
│       │   ├── Button.tsx
│       │   ├── CategoryPills.tsx      # Fetches /api/categories, renders PillSelect
│       │   ├── DescriptionPills.tsx   # Fetches /api/expenses/descriptions?category=
│       │   ├── MobileNav.tsx          # Hamburger + Panel side drawer
│       │   ├── Panel.tsx              # Generic slide-out panel
│       │   ├── PillSelect.tsx         # Generic scrollable pill selector
│       │   └── DataTable/             # Table island with column-driven config
│       │       ├── DataTable.tsx
│       │       ├── renderers.tsx
│       │       └── DataTable.test.tsx
│       ├── entries/           # Hydration entry points (one per island)
│       │   ├── category-pills.tsx
│       │   ├── data-table.tsx
│       │   ├── description-pills.tsx
│       │   └── mobile-nav.tsx
│       ├── lib/
│       │   └── hydrate.tsx    # Shared hydrateIsland helper
│       ├── test-setup.ts
│       └── types/
│           └── index.ts
├── static/                  # Vite build output (gitignored)
├── docker-compose.yml
├── Dockerfile               # Multi-stage: Node → Go → Alpine
├── .github/workflows/docker.yml  # CI for GHCR
├── Makefile
├── .air.toml                # Hot reload config
└── .env.example
```

---

## 4. Schemas & Models

### SQLite Tables (from `internal/db/db.go` migrations)

**`expenses`**
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT |
| amount | REAL | NOT NULL |
| category | TEXT | NOT NULL |
| description | TEXT | |
| date | TEXT | YYYY-MM-DD |
| currency | TEXT | DEFAULT 'USD' |
| paid_by | INTEGER | FK → users(id) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |

Indexes: `idx_expenses_category`, `idx_expenses_created_at`

**`preferences`**
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY CHECK(id = 1) |
| currency | TEXT | DEFAULT 'USD' |
| user_id | INTEGER | FK → users(id) |
| timezone | TEXT | DEFAULT '' |

**`users`**
| Column | Type | Constraints |
|--------|------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT |
| name | TEXT | NOT NULL UNIQUE |

**`exchange_rates`**
| Column | Type | Constraints |
|--------|------|-------------|
| base_currency | TEXT | Part of PK |
| target_currency | TEXT | Part of PK |
| rate | REAL | NOT NULL |
| date | TEXT | Part of PK |
| fetched_at | DATETIME | DEFAULT CURRENT_TIMESTAMP |

### Go Domain Types (`internal/domain/models.go`)

```go
type Expense struct {
    ID              int64     // json:"id"
    Amount          float64   // json:"amount"
    Category        string    // json:"category"
    Description     string    // json:"description"
    Date            string    // json:"date" — YYYY-MM-DD
    Currency        string    // json:"currency"
    PaidByID        int64     // json:"paidById,omitempty"
    PaidByName      string    // json:"paidByName,omitempty" — computed via SQL JOIN
    CreatedAt       time.Time // json:"createdAt"
    ConvertedAmount float64   // json:"convertedAmount" — computed at render time
}

type Preferences struct {
    Currency string // json:"currency"
    UserID   int64  // json:"userId"
    Timezone string // json:"timezone"
}

type User struct {
    ID   int64  // json:"id"
    Name string // json:"name"
}

type DailyGroup struct {
    Date           string    // json:"date"
    Expenses       []Expense // json:"expenses"
    Total          float64   // json:"total"
    ConvertedTotal float64   // json:"convertedTotal"
}

type CalendarMonth struct {
    Label    string         // json:"label"
    HasToday bool           // json:"hasToday"
    Cells    []CalendarCell // json:"cells"
}

type DescriptionCount struct {
    Description string // json:"description"
    Count       int    // json:"count"
}
```

---

## 5. Routes

### HTML Routes (Server-Rendered)

| Method | Route | Handler | Description |
|--------|-------|---------|-------------|
| GET | `/` | `HandleDaily` | Daily grouped view (default homepage) |
| GET | `/calendar` | `HandleCalendar` | Monthly calendar grid |
| GET | `/expenses` | `HandleList` | All expenses table |
| GET | `/expenses/new` | `HandleAdd` | New expense form |
| POST | `/expenses/new` | `HandleCreate` | Create expense |
| GET | `/expenses/edit/{id}` | `HandleEdit` | Edit expense form |
| POST | `/expenses/edit/{id}` | `HandleUpdate` | Update expense |
| POST | `/expenses/delete/{id}` | `HandleDelete` | Delete expense |
| GET | `/preferences` | `HandlePreferences` | Settings page |
| POST | `/preferences` | `HandleSavePreferences` | Save settings |
| GET | `/users` | `HandleUsers` | User management |
| GET | `/users/new` | `HandleUserNew` | New user form |
| POST | `/users/new` | `HandleUserCreate` | Create user |
| GET | `/users/edit/{id}` | `HandleUserEdit` | Edit user form |
| POST | `/users/edit/{id}` | `HandleUserUpdate` | Update user |
| POST | `/users/delete/{id}` | `HandleUserDelete` | Delete user |
| GET | `/static/` | FileServer | Vite build output |

### JSON API Routes

| Method | Route | Handler | Description |
|--------|-------|---------|-------------|
| GET | `/api/expenses` | `HandleList` | List expenses (query: `?n=limit`) |
| POST | `/api/expenses` | `HandleCreate` | Create expense (JSON body) |
| GET | `/api/expenses/{id}` | `HandleGet` | Get single expense |
| PUT | `/api/expenses/{id}` | `HandleUpdate` | Update expense (JSON body) |
| DELETE | `/api/expenses/{id}` | `HandleDelete` | Delete expense |
| GET | `/api/categories` | `HandleCategories` | List recent categories |
| GET | `/api/expenses/descriptions` | `HandleDescriptions` | Top descriptions by category (query: `?category=`) |
| GET | `/api/summary` | `HandleSummary` | Total + per-category totals |

---

## 6. Features

### Expense Tracking
- Add/edit/delete expenses with amount, category, description, currency, date, payer
- 17 supported currencies with symbol rendering
- Default date = "today" in user's timezone
- Category pills: horizontally scrollable quick-select from recent categories
- Description pills: dynamically fetched top descriptions for selected category (debounced 500ms)

### Views
- **Daily** (`/`): Grouped by date with per-day totals, ghost-variant DataTable
- **All Expenses** (`/expenses`): Full table with converted totals, default-variant DataTable
- **Calendar** (`/calendar`): Monthly grid from (earliest expense | 1yr ago) to (today + 1yr), sticky month labels, today badge, auto-scroll to today

### Multi-Currency
- Live exchange rates fetched from frankfurter.dev every 6 hours
- Rates stored in SQLite, cached per-day
- Expense amounts converted to preference currency for display
- Conversion shown in list view with rate date

### User Management
- Create/edit/delete users
- "Paid By" field on expenses links to users
- User deletion: atomic transaction (delete user + clear expenses.paid_by)

### Preferences
- Preferred currency (17 options)
- Default user/payer
- Timezone (22 IANA zones) — affects "today", "yesterday", calendar range

### Mobile UX
- Responsive: full-bleed cards on mobile, rounded on desktop
- Sticky navbar
- Hamburger side-panel navigation (left)
- Preferences link bottom-aligned in mobile nav
- Mobile content padding wrapper (`px-3 sm:px-0`)

---

## 7. Key Design Decisions

### React Islands over SPA
Only interactive elements (category pills, description pills, data tables, mobile nav) are hydrated as React components. Static buttons, forms, and text remain server-rendered Go templates. This minimizes JavaScript payload while allowing rich interactivity where needed.

### Split Repository Interfaces
The original 20-method `Repository` god interface was split into `ExpenseRepository`, `UserRepository`, `PreferenceRepository`, `RateRepository`. A composite `Repository` interface is retained for convenience. This improves testability and makes dependencies explicit.

### Column-Driven DataTable
The DataTable island accepts declarative column config from Go templates (`columns`, `data`, `variant`, `meta`, `actions`). This enables ghost/default variants and makes the table reusable across list, daily, and users pages without duplicating rendering logic.

### Vite Dev Proxy
Vite dev server runs on `:8081` and proxies non-static requests to Go on `:8080`. Internal Vite paths (`/@vite`, `/src/`, `/node_modules/`) bypass the proxy via a `bypass()` function. React Refresh preamble is injected by Go's `DevClient()` template function.

### Asset Manifest Matching
Production script tags are generated from Vite's `manifest.json`. Matching uses `Name` field OR `Src` path substring to handle both camelCase (`categoryPills`) and kebab-case (`category-pills`) naming conventions.

### Timezone Handling
Dates are stored as bare `YYYY-MM-DD` strings in SQLite. The user's timezone preference (or server's `TZ` env var / `time.Local`) is used to compute "today", "yesterday", and calendar ranges. Template functions `humanDate` and `formatDate` accept a timezone parameter.

---

## 8. Architecture Structure

```
┌─────────────────────────────────────────────────────────────────┐
│                         Browser                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │  Category   │  │ Description │  │  DataTable  │  MobileNav │
│  │   Pills     │  │   Pills     │  │   Island    │            │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘            │
│         │                │                │                     │
│         └────────────────┴────────────────┘                     │
│                          │                                       │
│                   Vite dev (:8081)                               │
│                          │                                       │
└──────────────────────────┼───────────────────────────────────────┘
                           │ proxy non-static
┌──────────────────────────┼───────────────────────────────────────┐
│                      Go Server (:8080)                           │
│  ┌───────────────────────┴───────────────────────────────────┐  │
│  │                    http.NewServeMux                        │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌────────────────────┐ │  │
│  │  │ HTML Routes │  │  API Routes │  │   Static Files     │ │  │
│  │  │ (templates) │  │   (JSON)    │  │   (/static/)       │ │  │
│  │  └──────┬──────┘  └──────┬──────┘  └────────────────────┘ │  │
│  │         │                │                                  │  │
│  │  ┌──────▼──────┐  ┌──────▼──────┐                          │  │
│  │  │HTMLHandler  │  │ APIHandler  │                          │  │
│  │  └──────┬──────┘  └──────┬──────┘                          │  │
│  │         │                │                                  │  │
│  │  ┌──────▼────────────────▼──────┐                          │  │
│  │  │         Service              │                          │  │
│  │  │  (validation, defaults,      │                          │  │
│  │  │   conversion logic)          │                          │  │
│  │  └──────┬───────────────────────┘                          │  │
│  │         │                                                   │  │
│  │  ┌──────▼──────────────────────────────────────┐           │  │
│  │  │              Repository.Repos                │           │  │
│  │  │  ┌────────┐ ┌──────┐ ┌──────────┐ ┌───────┐ │           │  │
│  │  │  │Expense │ │ User │ │Preference│ │ Rates │ │           │  │
│  │  │  └────────┘ └──────┘ └──────────┘ └───────┘ │           │  │
│  │  └──────┬───────────────────────────────────────┘           │  │
│  │         │                                                   │  │
│  │  ┌──────▼──────┐  ┌─────────────┐  ┌──────────────────┐   │  │
│  │  │  SQLite     │  │ rate.Client │  │  Template Engine │   │  │
│  │  │  (embedded) │  │(frankfurter)│  │  (html/template) │   │  │
│  │  └─────────────┘  └─────────────┘  └──────────────────┘   │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 9. Issues & Technical Debt

1. **Currency conversion in HTML handlers** (`HandleList`, `HandleDaily`) duplicates conversion logic. Should be moved to `service.DailyGroups()` / `service.ListExpenses()` with a target currency parameter. (Identified Apr 29, still open)

2. **~~`rate.Client` is hardcoded, not injectable~~** — Fixed. `service.RateFetcher` is a consumer-owned interface; `*rate.Client` structurally satisfies it. `service.New()` now accepts the rate client as a parameter. (Identified Apr 29, resolved Jun 2026)

3. **0 HTML handler tests** — At the time this was written all Go tests were API handlers; service-layer tests have since been added, but the primary user-facing surface (HTML rendering, form submission, redirects) still has zero coverage. (Identified Apr 29, still open)

4. **`SummaryByCategory` sums mixed currencies** — The API returns raw amounts without conversion. Categories with multi-currency expenses produce meaningless totals. (Identified Apr 29, still open)

5. **`ConvertedAmount` on domain model** — This is a presentation concern leaking into `domain.Expense`. It should be computed at render time, not stored on the struct. (Identified Apr 29, still open)

6. **No DataTable SSR fallback / CLS prevention** — Only `<noscript>` placeholder exists. Server-rendered `<table>` inside `data-table-root` would prevent layout shift on hydration. (Identified Apr 30, still open)

7. **`humanDate` / `currencySymbol` duplicated** — Exist in both Go (`renderer.go`) and JS (`renderers.tsx`). They need to stay in sync manually. (Identified Apr 30, still open)

8. **~~Console.log in DataTable.tsx~~** — Stale debug logging removed. (Identified Apr 30, resolved Jun 2026)

9. **Action registry is empty** — `onClick` actions in DataTable column config won't execute without registry entries. Not a bug (feature not yet needed). (Identified Apr 30)

10. **Mobile nav requires JS** — No-JS mobile users won't see nav links. Acceptable tradeoff for islands architecture. (Identified Jun 1)

11. **`float64` for monetary amounts** — Known accepted risk. No plans to switch to integer cents or decimal library.

---

## 10. Usage

### Development

```bash
# Install deps
cd ui && npm install

# Start both servers (Go with Air hot reload + Vite dev)
make dev

# Or manually:
# Terminal 1: DEV=true go run ./cmd/server
# Terminal 2: cd ui && npm run dev
# Access http://localhost:8080 (Vite proxies to Go)
# LAN: http://<ip>:8080 (auto-detected by Makefile)
```

### Production Build

```bash
make prod    # npm run build + go build
./bin/server # Runs on :8080 with production assets
```

### Docker

```bash
make docker-up    # docker compose up -d
make docker-down  # docker compose down
```

### Environment Variables

| Var | Default | Description |
|-----|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATA_DIR` | `~/.expensif` | SQLite database directory |
| `DEV` | `false` | Development mode (enables Vite HMR) |
| `VITE_DEV_HOST` | `localhost:8081` | Vite dev server address |
| `TZ` | system local | Server timezone fallback |

### Testing

```bash
go test ./...          # 40 tests, all pass
go vet ./...           # clean
cd ui && npm run build # TypeScript + Vite build
cd ui && npx vitest run # 7/7 frontend tests pass
```

---

## 11. Tests

### Go Tests (40 total)

| File | Tests | Coverage |
|------|-------|----------|
| `handlers_api_test.go` | 20 | API CRUD, validation, categories, summary, `isValidationErr` |
| `server_test.go` | 2 | Server startup, shutdown, graceful handling |
| `service_test.go` | 12 | Rate refresh, conversion math, validation |
| `mock_repo_test.go` | N/A | Mock implementation (used by above) |

**Notable tests:**
- `TestAPIList_LimitQuery` — query param parsing
- `TestAPICreate_MissingCategory` / `TestAPICreate_InvalidAmount` — validation errors
- `TestIsValidationErr` — 8 subtests covering direct + wrapped errors

**Gap:** Zero HTML handler tests. No integration tests for form submission flow, template rendering, or redirects.

### Frontend Tests (7 total)

| File | Tests |
|------|-------|
| `DataTable.test.tsx` | 7 |

**Test subjects:**
- Renders headers and rows for default variant
- Renders badge type correctly
- Renders currency with meta (converted + original amounts)
- Ghost variant has no outer card wrapper
- Renders action buttons (link with param replacement)
- Custom renderer override works
- Date formatting ("Today")

**Gap:** No tests for PillSelect, CategoryPills, DescriptionPills, MobileNav, or Panel.

---

## 12. Notable Files

| File | Purpose |
|------|---------|
| `cmd/server/main.go` | Entry point: DB init, service wiring, background rate fetcher, graceful shutdown |
| `internal/web/server.go` | Route registration, middleware stack (RecoverPanic + RequestLog) |
| `internal/web/handlers_html.go` | All HTML page handlers (~400 lines) |
| `internal/web/handlers_api.go` | All JSON API handlers (~150 lines) |
| `internal/service/service.go` | Business logic, validation, currency conversion |
| `internal/repository/sqlite.go` | All SQL queries (~350 lines) |
| `internal/assets/assets.go` | Vite manifest parsing, production script tag generation |
| `internal/web/renderer.go` | Template engine: func map, page data struct, render method |
| `ui/vite.config.ts` | Build config: 4 island entry points, dev proxy, manifest generation |
| `ui/src/lib/hydrate.tsx` | Shared `hydrateIsland(name, Component)` helper |
| `templates/base.html` | Layout: nav (desktop + mobile), flash messages, island script tags |
| `templates/form.html` | Shared expense form: amount, category+pills, description+pills, currency, date, paid by |
| `docker-compose.yml` | Production deployment: named volume for SQLite |
| `.github/workflows/docker.yml` | CI: build linux/amd64, push to GHCR |

---

*Generated from codebase exploration on 2026-06-02. See `.plan/handoffs/` for session history and decision logs.*
