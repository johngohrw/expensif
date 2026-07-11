---
type: grilling
blocked_by: [04]
claimed_by: claude-code-session-2026-07-12
claimed_at: 2026-07-11T17:51:06Z
---

# Test strategy for the date-indexed timeline

## Question

What is tested, at which seam, and what is the first test written?

Graduated from the map's fog once [Re-shape DailyGroups around dates](./04-date-indexed-daily-groups.md)
fixed the query shape. It could not be specified before that: the seam under test did
not exist.

The daily view has **zero test coverage today** — nothing in `internal/web/*_test.go`
or `internal/service/service_test.go` touches `daily` or `DailyGroups`. Ticket 04
deletes `DailyGroups(ctx, limit)` outright, so there is nothing to characterize; the
tests are written against the new interface, not the old one.

Read `.plan/expensif/testing-strategy.md` before deciding anything.

Resolve, one at a time:

- **Which seam.** `DailyGroupsInRange` is a service call over a repository interface
  that `service_test.go:328` already stubs (`expenseRepoStub.ListExpensesInRange`).
  Is the service the test surface, or does `HandleDaily` get exercised through
  `mock_repo_test.go` as well? Ticket 04's answer says the interface is the test
  surface — does that settle it, or does the handler's currency loop want its own test?
- **The gap-filling cases.** Empty window, window with expenses only on the first and
  last day, a single day, `start > end`. What is the contract when `start > end` — an
  error, or an empty slice?
- **The `Expenses` is never nil invariant.** Ticket 04 makes this a promise the island
  depends on. Is it asserted once, or on every gap-filled day in every test?
- **The DST/timezone boundary.** Ticket 04 argues no window edge can straddle a DST
  transition, because the walk runs on UTC-parsed dates and the timezone is consulted
  once to name today. Is that argument tested, and if so where — the walk is
  timezone-free, so the only thing left to test lives in the handler.
- **`UpcomingGroups`.** Ungapped, descending, dated strictly after today. Does the
  boundary (an expense dated exactly today) have a test?
- **Where the first test goes**, and whether it lands before or with the implementation.
