# 2026-07-09 — Daily timeline: continuous days, empty days included

## What this session did

Charted a wayfinder map at [`.plan/daily-timeline/`](../daily-timeline/map.md) and
resolved four of its nine tickets. The effort: make the daily view show **every day**
in its window, not just days that own an expense, with muted empty days carrying an
add affordance.

**Do not re-read the decisions here — they live in their tickets.** The map's
Decisions-so-far indexes them with a one-line gist and a link. Start there.

## Where to pick up

The frontier is two unblocked tickets. They are independent, so two sessions can run
them in parallel — claim by setting `Status: claimed` and **commit the claim before
doing any work**, or the other session will duplicate it.

- [Re-shape DailyGroups around dates](../daily-timeline/tickets/04-date-indexed-daily-groups.md)
  — **the keystone.** Tickets 05, 06, and 07 all queue behind it. It inherits a
  sharper brief than when it was written: the window ends at today, and Upcoming is a
  second, separate query rather than an extension of the window.
- [Validate expense dates on the write path](../daily-timeline/tickets/09-validate-expense-dates.md)
  — independent of everything else.

Everything else (05, 06, 07) is blocked on 04.

## State of the tree

Clean. `go build ./...`, `go vet ./...`, `go test ./...` all pass. Nothing is
half-finished; no code from this effort has been written.

**No implementation has landed, by design.** Wayfinder plans; it does not do. The
destination is a spec at `.plan/daily-timeline/spec.md`, which does not exist yet.

A throwaway UI prototype was built, iterated with the user over eight rounds, approved,
and then deleted. Its approved markup is preserved verbatim at
[`.plan/daily-timeline/assets/day-entry-ledger.html.approved`](../daily-timeline/assets/day-entry-ledger.html.approved).
Do not go looking for `handlers_html_prototype.go` or `empty-day-prototype.html` —
they are gone on purpose.

## Things you would otherwise learn the hard way

- **The daily view has zero test coverage.** Nothing in `internal/web/*_test.go` or
  `internal/service/service_test.go` touches `daily` or `DailyGroups`. Ticket 04
  changes that function's signature. Consider characterization tests first — it is the
  only work that is both unblocked and unambiguously safe.
- **The design ran ahead of the data layer.** All the visual decisions are made, but
  nothing user-visible can land until 04 does: without the date-indexed query there
  are no empty days to render, so shipping the ledger alone delivers the redesign
  without the feature that motivated it.
- **The repository needs no new SQL.** `ListExpensesInRange(ctx, start, end)`
  (`internal/repository/sqlite.go:69`) already takes a date range with no `LIMIT`, and
  `HandleCalendar` already uses it this way. `HandleCalendar` (`handlers_html.go:126`)
  also already writes the gap-filling day loop — at the handler layer. That is the
  precedent ticket 04 has to accept or reject.
- **`templates/partials/data-table.html` emits no table markup**, only a root div and
  JSON props. Every table on the daily view is client-rendered React today, and the
  page shows a `<noscript>` message rather than any expense without JS. The ledger
  changes this — see the map's Notes.
- **Two decisions rest on premises that later changed.** Both are flagged in-place
  rather than silently re-decided: ticket 07 records that ticket 02's island choice
  was made when the page still required JS (it will not), and the map's Out-of-scope
  section records that collapsing empty-day runs was ruled out for a reason that no
  longer holds. Neither has been reopened. Both should be, honestly, if touched.
- **`{"date":"banana"}` returns `201`.** Verified by probe, not inference. See ticket
  03's answer.

## Environment

- `make dev` runs `air` (Go, :8080) and Vite (:8081) together. `DEV=true` is set in the
  user's shell — templates and the asset manifest behave differently without it.
- The user's dev database (`~/.expensif/expenses.db`) is **empty**. Prototyping used a
  scratch DB via `DATA_DIR=<tmp> PORT=8090 go run ./cmd/server`, seeded through
  `POST /api/expenses`. Do not seed the user's database.
- The Go server caches parsed templates at startup; restart after template edits.

## Suggested skills

Apply only if available in your environment; all live in `.pocock-skills/<name>/SKILL.md`
and are **not** auto-discovered — read the file and follow it.

- `wayfinder` — the map's protocol. Read it before touching `.plan/daily-timeline/`.
  Claim one ticket, resolve one ticket, per session.
- `codebase-design` — ticket 04 is a module-depth question (does the gap-filling loop
  live in the service or the handler, and does `HandleCalendar` share it?).
- `domain-modeling` — terms have firmed up ("the rail", "Upcoming", "day entry",
  "empty day") and `CONTEXT.md` still does not exist.
- `prototype` — if a further UI question arises. Sub-shape A on the existing `/daily`
  route behind `?variant=`, gated on `DEV=true`, worked well.
- `to-spec` — the map's destination. Invoke once the frontier is clear to synthesise
  `.plan/daily-timeline/spec.md`.
