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
daily view degrades gracefully today; it does not. **But the ledger changes this** —
[ticket 08](./tickets/08-day-entry-ledger-redesign.md) drops `data-table` from the
daily view, so once implemented the page renders expenses server-side and works
without JS. Ticket 02's island decision was taken before that was true — and
[ticket 07](./tickets/07-scroll-island-contract.md) has since re-put it to the user on the
new facts, who chose the island again. The page is JS-free except for infinite scroll.

Bounds are settled up front and constrain every ticket: **anchor at today, walk
backwards a fixed window, paginate older windows on demand. No future days.**

## Decisions so far

<!-- one line per resolved ticket: gist + link -->

- [Window size and how older days load](./tickets/02-window-size-and-pagination.md) —
  rolling 30 days from today; older windows appended by a new infinite-scroll island;
  first window stays server-rendered; scroll stops at the earliest expense. No new SQL
  needed — `ListExpensesInRange` already fits. *Undermined by 06: the island fetches HTML,
  not JSON, and is not React; the window and termination rule stand.*
- [Muted empty-day card design](./tickets/01-muted-empty-day-design.md) — an empty day
  is a **ledger line**, not a card: date gutter, "no expenses", and an always-visible
  `+`. Hover-only affordances rejected (mobile); full cards rejected (22 of 30 days
  are empty — they are the page's background, not its content).
- [What happens to expenses dated after today](./tickets/03-future-dated-expenses.md) —
  they render in an **"Upcoming" section above the timeline**, grouped by day, *not*
  gap-filled. The window still ends at today. Probed: future dates and even
  `"date":"banana"` are accepted with `201` — `validateExpenseInput` never checks the
  date. Today's daily view already shows future expenses above today, so hiding them
  would be a regression.
- [Redesign the day entry as a ledger row](./tickets/08-day-entry-ledger-redesign.md) —
  the card is gone from the whole timeline. Five columns, two unbroken rails, a 28px
  row unit, month-break dividers. **Drops `data-table` from this page**, so expenses
  become server-rendered and the daily view works without JS. Delete demoted to the
  edit page. Approved markup: [`assets/day-entry-ledger.html.approved`](./assets/day-entry-ledger.html.approved).
- [Validate expense dates on the write path](./tickets/09-validate-expense-dates.md) —
  non-empty dates must parse as `YYYY-MM-DD`; future dates remain valid; empty dates
  still default to today. Both API and HTML form write paths reject unparseable dates
  with `ErrInvalidDate` → `400 Bad Request`. Existing garbage rows are left untouched.
- [Contain the day-card chrome drift](./tickets/06-day-card-chrome-drift.md) — **there is no
  drift.** The endpoint serves HTML from the same Go partial as the first window, so the ledger
  has one implementation. `GET /daily/older?start=&end=`, cursor carried as `data-next-*` on the
  fragment. The foot's four states are server-rendered and toggled by `data-state`, keeping zero
  Tailwind in TypeScript. The island is ~30 lines of vanilla TS, not React. Undermines 02 and 07.
- [The infinite-scroll island's contract](./tickets/07-scroll-island-contract.md) — the island
  survives an honest re-ask against the now JS-free page. Sentinel-triggered, one fetch in
  flight, manual retry on error rather than an armed observer, terminal marker at the earliest
  expense, and a **server-issued cursor** so the island does no date math. No `<noscript>`
  fallback — it would grow `HandleDaily` a third branch. *Undermined by 06: its JSON DTO and
  `mountIsland` are void; everything above stands.*
- [Should HandleDaily's two branches converge](./tickets/05-converge-handledaily-branches.md) —
  **they don't.** `?date=` is a single-day detail view for any date in ±1 year (the calendar
  links every cell), not a one-day window of a timeline that ends at today — so it keeps its
  own branch and `ExpensesInRange`. Convergence is of *rendering*: a populated day draws
  through ticket 08's day-entry partial verbatim; an empty `?date=` shows the fuller "No
  expenses for this day" state, so ticket 01's muted line has one consumer (the timeline) and
  ticket 08's day has two. Shared currency loop stays shared.
- [Re-shape DailyGroups around dates, not expenses](./tickets/04-date-indexed-daily-groups.md) —
  `DailyGroups(ctx, limit)` splits into `DailyGroupsInRange(ctx, start, end)` (gap-filled)
  and `UpcomingGroups(ctx, after)` (ungapped). The day-walk lives in the **service**;
  `HandleCalendar` keeps its own. No new fields on `DailyGroup`, and `Expenses` is never
  nil. The service is timezone-free — bare date strings in, UTC walk inside — so no
  window edge straddles a DST transition. *Undermined by 06: the never-nil invariant's
  rationale was JSON marshalling, and there is no JSON. The invariant stands on a new
  footing — see 10.*
- [Test strategy for the date-indexed timeline](./tickets/10-test-strategy.md) — two seams:
  the **service** for all logic (gap-fill cases, the today/upcoming partition asserted from
  both sides in one test, never-nil asserted once), and the **fragment** through the real
  renderer for the island's contract only — `data-next-*` and `data-state`, never chrome.
  `start > end` is `ErrInvalidRange` → **400**, not an empty slice, so a bad cursor cannot
  masquerade as a spent timeline. `HandleDaily` gains an injectable clock so "anchor at
  today" is testable at a UTC-day boundary; ticket 04's DST claim is not tested (it is
  `time.Parse`'s behaviour). First test — `TestDailyGroupsInRange_GapFillsEmptyDays` —
  lands **red, before** the implementation. The 2026-06-02 testing strategy's renderer
  "blocker" does not exist (`server_test.go:45` already renders in-package), and its
  never-assert-on-HTML rule is narrowed, not broken.

## Not yet specified

<!-- Empty. Every patch has graduated into a ticket; the fog is clear to the destination. -->

The four patches that stood here — the zero-expense-ever empty state, today's row,
delete's new home, and the Upcoming section's chrome — graduated into tickets 11–14 once
[ticket 10](./tickets/10-test-strategy.md) closed the last question they could have hung
on. They are questions of *rendering*, and every decision they rest on (the ledger, the
muted line, the `?date=` consumer, Upcoming's content) is now made. Nothing on the route
to the destination is unspecifiable any more — what remains is takeable.

## Out of scope

- **Collapsing long runs of empty days** into an expandable "no expenses, Mar 3–17"
  row. Ruled out by the user while scoping. **Its original rationale no longer holds**
  — it was ruled out for forcing an island, and ticket 02 has since added one anyway.
  It stays out of scope because the user ruled it out, not because of the cost. Out of
  scope work does not graduate: if you want it, redraw the destination and start a
  fresh effort.
