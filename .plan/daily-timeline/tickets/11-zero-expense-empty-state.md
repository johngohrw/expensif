---
type: grilling
blocked_by: [01, 08]
claimed_by: claude-code-session-2026-07-12
claimed_at: 2026-07-11T18:40:00Z
---

# The zero-expense-ever empty state

## Question

What does the daily view show for an account that has never recorded an expense?

Today this is unambiguous. `templates/daily.html:55` branches on `{{if .DailyGroups}}`
and, when there are none, swaps the **entire page** for a centred "No expenses yet."
plus a primary add button. It works because an expense-indexed timeline with no
expenses has no groups at all.

The date-indexed timeline destroys that branch's premise. After
[Re-shape DailyGroups around dates](./04-date-indexed-daily-groups.md), a brand-new
account has **thirty gap-filled groups**, not zero — every day in the window exists.
`{{if .DailyGroups}}` is now always true, and a new user's first sight of the app is
thirty muted "no expenses" ledger lines from [ticket 01](./01-muted-empty-day-design.md),
each with its own `+`, and nothing else.

Resolve, one at a time:

- **Which state wins**, and on what condition. Keeping the page-swap means asking a
  question the window cannot answer — "has this account *ever* had an expense" is not
  "is this window empty". `GetEarliestExpenseDate` (`repository/sqlite.go`, already
  stubbed in both test doubles) answers it in one call, and
  [ticket 07](./07-scroll-island-contract.md) already relies on it for the scroll
  terminal. Is that the condition, or is a 30-muted-row page an acceptable — even
  honest — first run?
- **If the page-swap survives**, where does it live? It is a third top-level branch in
  `HandleDaily`, and [ticket 05](./05-converge-handledaily-branches.md) refused a third
  branch once already (for `<noscript>`). Does that refusal bind here?
- **The near-empty account.** An account whose only expense is four months old has a
  *window* full of empty days but is not new. It gets the timeline, not the empty state
  — confirm that, because it is the case that distinguishes "empty window" from "empty
  account", and it is the one a naive `len(expenses) == 0` check gets wrong.
- **The `?date=` view is not this.** An empty single-day filter already has its own
  fuller "No expenses for this day" state (ticket 05). This ticket must not change it.
- **What the empty state says.** The current copy is "No expenses yet." + "Add one".
  Ticket 08's ledger has a different visual language; does the copy or the affordance
  change with it?
