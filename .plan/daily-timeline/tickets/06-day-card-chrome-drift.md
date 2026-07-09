# Contain the day-card chrome drift

Type: grilling
Status: open
Blocked by: 01, 04

## Question

[Window size and how older days load](./02-window-size-and-pagination.md) put the first
30-day window in `daily.html` (Go template) and every window below it in a React
island. Day-card chrome — the date header, `humanDate` subtitle, footer total, add
affordance, and the muted empty-day treatment from ticket 01 — is therefore written
twice, in two languages, and must render identically at the seam. A user scrolling
past day 30 must not see the cards change.

This is an accepted cost, not a mistake. The job is to contain it, not undo it.

Note what is *not* at risk: the tables. `templates/partials/data-table.html` emits no
table markup at all, only a `<div data-table-root>` and JSON props — the React
`DataTable` renders every table on this page already. Only the chrome duplicates.

Resolve:

- Is the duplication load-bearing enough to warrant a test that renders one
  `DailyGroup` through both paths and asserts the same fields appear? Or is a visual
  check at the seam enough?
- Does currency conversion produce the same string on both sides? The Go path formats
  with `printf "%.2f"` and `$.CurrencySymbol` (`daily.html:37`); the React path would
  format in TSX. Rounding and symbol placement are easy to diverge on.
- `formatDate` / `humanDate` are Go template funcs (`renderer.go:77,99`) taking a
  timezone. What renders the date in the island — a JS reimplementation, or does the
  JSON endpoint ship pre-formatted display strings alongside the raw date?
- Cheaper alternative worth pricing before accepting: does the endpoint return
  **display-ready** `DailyGroup`s (formatted date, formatted total, formatted label)
  so the TSX card is a dumb renderer and all formatting stays in Go?
