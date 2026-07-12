---
type: task
blocked_by: []
claimed_by: claude-code-session-2026-07-13
claimed_at: 2026-07-12T19:07:15Z
---

# The date-indexed service, wired end to end

## Question

Deliver the spec's **Service layer** and **clock seam** decisions as a working page:
`DailyGroupsInRange(ctx, start, end)` (gap-filled, newest-first, `Expenses` never
nil) and `UpcomingGroups(ctx, after)` (ungapped) replace `DailyGroups(ctx, limit)`,
which is deleted. `ErrInvalidRange` on an inverted range, `ErrInvalidDate` reused for
unparseable endpoints, both mapped to 400. The daily handler's unfiltered branch
computes today once in the Preferences timezone via an injectable clock, asks for the
rolling 30-day window plus upcoming groups, and passes a `NoExpensesEver` flag
(earliest expense date `== ""`).

This is a tracer bullet, not a service-only slice: the **existing card template
renders the new groups unchanged**, so the page becomes date-indexed — empty days
appear as the already-existing "No expenses for this day" cards, future expenses
still render above today (uncapped; the ledger chrome is the next ticket's job) —
and the change is demoable in a browser.

Tests come **first**, red, per the spec's Testing Decisions:
`TestDailyGroupsInRange_GapFillsEmptyDays`, the gap-fill cases, the today/upcoming
partition asserted from both sides in one test, the never-nil invariant in one named
test, `start > end` → error → 400, and the frozen-clock timezone test
(`2026-07-11T23:00Z`, `Pacific/Auckland`, window ends `2026-07-12`).

Done when: all named tests green, the old signature gone, and the browser shows a
gap-filled 30-day window on a scratch DB.
