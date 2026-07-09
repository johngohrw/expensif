# Daily Timeline — continuous days, empty days included

## Destination

A spec at `.plan/daily-timeline/spec.md` describing a daily view whose timeline is
**date-indexed rather than expense-indexed**: every day in the window appears, days
with no expenses render muted with an add affordance. The spec is detailed enough
that one fresh session can implement it end to end.

## Notes

Domain: Go server-rendered templates + React Islands (partial hydration), SQLite.

Existing ground truth, read before deciding anything:

- `internal/service/service.go:143` — `DailyGroups(ctx, limit)`. Takes the last 100
  **expenses** and buckets them by date into a map. Only days that own an expense
  ever become a group. This is the thing being replaced.
- `internal/web/handlers_html.go:237` — `HandleDaily`. Two disjoint branches: a
  `?date=` single-day filter that builds its own one-element group, and the
  unfiltered list that calls `DailyGroups`.
- `internal/web/handlers_html.go:126` — `HandleCalendar`. **The precedent.** Already
  walks every day in a bounded range (`139-152`), renders zero-spend days, and
  aggregates via one `ExpensesInRange` call. Steal its shape; do not re-invent it.
- `templates/daily.html:44` — an empty-day branch already exists ("No expenses for
  this day. / Add one"). It is currently unreachable except via `?date=`.

Skills every session should consult: `codebase-design` before changing the service
boundary; `domain-modeling` whenever a term firms up (`CONTEXT.md` does not exist
yet — create it lazily).

Standing preference: server-render unless a decision *forces* an island. Collapsing
was ruled out precisely to keep this page static.

Bounds are settled up front and constrain every ticket: **anchor at today, walk
backwards a fixed window, paginate older windows on demand. No future days.**

## Decisions so far

<!-- one line per resolved ticket: gist + link -->

_None yet._

## Not yet specified

- **The zero-expense-ever empty state.** `daily.html:55` currently swaps the whole
  page for "No expenses yet." Once every day is a card, a brand-new account renders
  N muted cards instead. Which wins? Hangs on the empty-day design.
- **Today's card.** Whether today is visually distinguished from any other day, and
  whether it is distinguished *differently* when empty. `HandleCalendar` has an
  `IsToday` flag and `CalendarCell.IsToday`; daily has no equivalent. Hangs on the
  empty-day design.
- **Test strategy.** `internal/web/handlers_api_test.go` exists; a date-indexed
  timeline wants tests for gap-filling, window edges, and the DST/timezone boundary.
  Can't be specified until the query shape lands.
- **Query cost.** Whether one `ExpensesInRange` per window is enough, or whether the
  view wants a `GROUP BY date` aggregate the way the calendar's `dayTotals` does.
  Hangs on the window size.
- **Timezone day boundaries.** `nowInTZ(prefs.Timezone)` defines "today", but
  `Expense.Date` is a bare `2006-01-02` string. Whether a window edge can straddle a
  DST transition, and whether that matters, isn't sharp yet.

## Out of scope

- **Collapsing long runs of empty days** into an expandable "no expenses, Mar 3–17"
  row. Ruled out by the user while scoping. It is the one decision that would force
  this page to become a React island; leaving it out keeps the timeline
  server-rendered. Returns only as a fresh effort.
