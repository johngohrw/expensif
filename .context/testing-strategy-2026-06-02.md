# Expensif — Testing Strategy Discussion (2026-06-02)

## Current State

**Go (28 tests):** All API handler tests. Service layer tested only indirectly. HTML handlers, repository SQL, and template rendering have zero coverage.

**Frontend (7 tests):** Only DataTable. PillSelect, CategoryPills, DescriptionPills, MobileNav, Panel — none tested.

**Integration:** None. No end-to-end flows.

---

## Strategy: Layers, Not Uniformity

Don't aim for 80% blanket coverage. Aim for **confidence per layer**.

### Layer 1: Service Layer (highest value, easiest win)

The service has pure business logic currently only tested through API handler indirection. Direct service tests would cover:

- Validation rules (`amount > 0`, `category != ""`)
- Timezone-aware "today" default
- `DailyGroups` grouping + sorting
- Currency conversion math

**Blocker:** `rate.Client` is hardcoded in `service.New()`. Make it an interface — this is existing architecture debt #4. Once injectable, pass a mock rate client that returns fixed rates. Test conversion math deterministically.

**Example test:**
```go
func TestCreateExpense_DefaultsToTodayInTimezone(t *testing.T) {
    svc := service.New(mockRepos, mockRateClient)
    id, _ := svc.CreateExpense(ctx, 10.0, "food", "lunch", "", "USD", 0)
    exp, _ := svc.GetExpense(ctx, id)
    // assert exp.Date == "2026-06-02" (or whatever tz says)
}
```

### Layer 2: HTML Handler Tests (high value, medium effort)

The blocker is the template renderer. `HTMLHandler.render()` calls `renderer.Render()` which reads templates from disk. In tests, file I/O makes tests slow and flaky.

**Solution:** Introduce a `TemplateRenderer` interface:

```go
type TemplateRenderer interface {
    Render(w io.Writer, name string, data PageData) error
}
```

In production, `*Renderer` satisfies this. In tests, use a `mockRenderer` that captures `PageData` into a buffer. Now assert:

- `HandleCreate` → redirects to `/` with 303
- `HandleCreate` with validation error → re-renders "add" template with `FlashError = true`
- `HandleAdd` → passes correct `Today` based on timezone
- `HandleDaily` → passes `DailyGroups` with converted amounts

**Quick wins:** POST handlers are mostly redirects. These are ~5 lines each to test.

### Layer 3: Repository Tests with In-Memory SQLite (medium value, easy)

`modernc.org/sqlite` supports `:memory:` databases. Spin one up per test, run migrations, test SQL directly.

**What to test:**
- `ListDescriptionsByCategory` — case-insensitive matching, empty description exclusion, limit=20
- `ListCategories` — 3-month cutoff
- `SummaryByCategory` — aggregation correctness
- `DeleteUser` — atomicity (user gone, expenses.paid_by cleared)
- Migrations — adding columns on older schemas

This catches SQL bugs that mock tests miss. The mock repo doesn't run real SQL — it simulates in Go. The SQLite implementation could have a typo in a `GROUP BY` that the mock wouldn't catch.

### Layer 4: Frontend Island Tests (medium value, medium effort)

**PillSelect:** Render with options, click a pill, assert `onSelect` called. Test active state styling, disabled pills.

**CategoryPills:** Mock `fetch("/api/categories")`, assert pills render, assert click sets input value + dispatches change event.

**DescriptionPills:** Mock fetch, test debounce (advance timers), test loading state, test empty category clears pills.

**Panel:** Open/close state transitions, backdrop click calls `onClose`, Escape key, scroll lock (harder in jsdom but doable).

Use `msw` (Mock Service Worker) or simple `global.fetch = vi.fn()` for API mocking. `@testing-library/react` + `userEvent` for interactions.

### Layer 5: Integration / E2E (high value, high effort)

For an islands architecture, **Playwright** is the right tool. It tests the actual hydration pipeline: Go serves HTML → Vite chunks load → React hydrates → islands become interactive.

**Flows to cover:**
1. Add expense → redirected to Daily → expense appears in today's card
2. Click category pill → description pills update → click description pill → input filled
3. Calendar → scrolls to today → today's cell has red badge
4. Mobile viewport → hamburger opens → navigate to Preferences → save → toast appears

**When to add this:** After Layers 1-3 are solid. E2E is expensive to maintain; don't start here.

---

## Priority Order

| Priority | What | Effort | Why First |
|----------|------|--------|-----------|
| 1 | Make `rate.Client` interface | 10 min | Unlocks service testing |
| 2 | Service layer tests (validation, tz, DailyGroups) | 1-2 hrs | Pure logic, fast, high confidence |
| 3 | HTML handler POST redirect tests | 1 hr | Easy wins, catches form bugs |
| 4 | Repository in-memory SQLite tests | 2 hrs | Catches real SQL bugs |
| 5 | CategoryPills + DescriptionPills tests | 2 hrs | Interactive core of the form |
| 6 | Playwright E2E (1-2 flows) | 4 hrs | Sanity check for hydration |

---

## What NOT to Test

- **Template HTML output** — Too brittle. If you change a `div` to a `span`, tests break. Assert on `PageData` instead.
- **Vite build** — `npm run build` passing is sufficient.
- **CSS/visuals** — Not worth visual regression for a personal tool.
- **Rate client HTTP calls** — It's a thin wrapper around frankfurter.dev. The contract is simple enough.

---

## Honest Assessment

The biggest testing gap is **HTML handlers**. They contain the most user-facing logic (form validation, redirects, data preparation) but are the hardest to test because of the template coupling. The `TemplateRenderer` interface extraction is a one-time refactor that pays off immediately.

For the frontend, focus on **CategoryPills and DescriptionPills** since they're the most complex interactive behavior. Panel and MobileNav are mostly presentational — manual testing is fine for now.
