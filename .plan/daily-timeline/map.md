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

Standing preference: server-render unless a decision *forces* an island. ~~Collapsing
was ruled out precisely to keep this page static.~~ **Superseded** — the user chose an
infinite-scroll island in [Window size and how older days load](./tickets/02-window-size-and-pagination.md).
This page will have a fetching island regardless. Keep the preference as a tiebreaker,
not as a rule.

Note for every session: `templates/partials/data-table.html` emits **no table markup**
— only a root div and JSON props. Every table on this page is already client-rendered
React, and the page already requires JS to show any expense. Do not reason as if the
daily view degrades gracefully today; it does not.

Bounds are settled up front and constrain every ticket: **anchor at today, walk
backwards a fixed window, paginate older windows on demand. No future days.**

## Decisions so far

<!-- one line per resolved ticket: gist + link -->

- [Window size and how older days load](./tickets/02-window-size-and-pagination.md) —
  rolling 30 days from today; older windows appended by a new infinite-scroll React
  island fetching JSON; first window stays server-rendered; scroll stops at the
  earliest expense. No new SQL needed — `ListExpensesInRange` already fits.

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
- **Timezone day boundaries.** `nowInTZ(prefs.Timezone)` defines "today", but
  `Expense.Date` is a bare `2006-01-02` string. Whether a window edge can straddle a
  DST transition, and whether that matters, isn't sharp yet.

## Out of scope

- **Collapsing long runs of empty days** into an expandable "no expenses, Mar 3–17"
  row. Ruled out by the user while scoping. **Its original rationale no longer holds**
  — it was ruled out for forcing an island, and ticket 02 has since added one anyway.
  It stays out of scope because the user ruled it out, not because of the cost. Out of
  scope work does not graduate: if you want it, redraw the destination and start a
  fresh effort.
