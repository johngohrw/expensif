# The infinite-scroll island's contract

Type: grilling
Status: open
Blocked by: 01, 04

## Question

Specify the new island end to end. It is the app's first pagination island and the
first island to fetch anything — every existing island (`data-table`,
`category-pills`, `description-pills`, `mobile-nav`) is handed its props as JSON at
mount and never talks to the server.

Resolve:

- **Endpoint.** Path and shape. `internal/web/handlers_api.go` exists; does this live
  there as `GET /api/daily?before=YYYY-MM-DD&days=30`, returning `[]DailyGroup`? Does
  it reuse the same service call `HandleDaily` uses (it should — see ticket 04)?
- **Trigger.** `IntersectionObserver` on a foot sentinel, or an explicit button? The
  decision said "infinite scroll", but a button is a legitimate degenerate case worth
  re-confirming once the cost is concrete.
- **States.** Loading, error, and the terminal marker ("No expenses before Mar 4,
  2024" — see ticket 02's termination rule). What happens on a failed fetch: retry,
  silent stop, or a visible error with a retry affordance?
- **How does the island know where to start and when to stop?** It needs the oldest
  date currently on screen and the earliest expense date. Both are server-known at
  first render — pass them as island props via `data-props`, like `hydrate.tsx:12`
  reads, rather than making the client discover them.
- **Mounting.** `hydrate.tsx` finds a single container by `[data-island="name"]`;
  `data-table.tsx` uses its own `querySelectorAll` loop and does not use
  `hydrateIsland` at all. Which pattern does this island follow, and does it need to
  mount the `DataTable` component directly as a child rather than via `data-table`'s
  document-wide scan? (It does — a scan cannot see nodes appended later.)
- **No-JS.** With the first window server-rendered, a JS-less user sees 30 days and no
  way back. Acceptable, or does a `<noscript>` link to a plain paginated fallback earn
  its keep? Note the page's tables already require JS, so a JS-less user sees 30 muted
  cards and nothing else.

---

## Update — after ticket 08 (the ledger)

**The premise of this ticket has changed.** It was written assuming the island renders
`DataTable` as a child of each appended day-card. There is no `DataTable` on the daily
view any more — verified on the prototype: 0 `data-table-root` nodes, against 8 on the
control page. The island would render plain ledger rows. Simpler.

The no-JS question is no longer a footnote. With expenses server-rendered, the daily
view works without JavaScript **for the first time** — a JS-less user gets a fully
functional 30-day timeline, not a `<noscript>` message. Infinite scroll would be the
only thing on the page that needs JS.

So re-open the mechanism, honestly:

- Ticket 02 chose an infinite-scroll island over a plain `?days=` link. One argument
  against the plain link was that this page already requires JS, so a no-JS fallback
  bought nothing. **That argument is now false.**
- A plain "Load older" link would leave the entire daily view JS-free and delete the
  island, the endpoint, and ticket 06's drift problem in one move.
- The user chose the island deliberately and may still want it. But the decision was
  made under a fact that no longer holds, and it deserves one honest re-ask before
  anyone writes the endpoint.

If the island survives: it renders ledger rows, needs the oldest on-screen date and
the earliest expense date as props, and the mount question below still stands.
