---
type: grilling
blocked_by: [01, 04]
---

# Contain the day-card chrome drift

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

---

## Update — after ticket 08 (the ledger)

The card is gone. What duplicates across Go template and TSX is now the **ledger**:
five columns, two rails, the 28px row unit, the month-break divider, and the expense
rows themselves — which the card never duplicated, because `DataTable` rendered them.
The drift surface **grew**, not shrank.

Two facts that reframe this ticket:

- **The repo already accepts this exact trade.** `templates/partials/button.html:1-3`
  opens with *"keep in sync with `ui/src/components/Button.tsx` — both must use
  identical Tailwind class mappings for variant + size."* Go-template/TSX duplication
  is a documented convention here, not a novel hazard.
- The ledger's alignment rests on three invariants (see ticket 08's answer) that are
  invisible in a screenshot and easy to break in a reimplementation — `self-start` on
  the expenses column, padding on rows not columns, explicit `leading-4`. A TSX
  reimplementation that "looks right" in isolation will drift at the seam.

That strengthens the case for the cheaper option this ticket already floated: have the
JSON endpoint return **display-ready** groups (formatted date, formatted total,
formatted amounts) so the TSX side is a dumb renderer and all formatting logic stays
in Go. Only the class strings duplicate then, exactly as `button.html` does it.
