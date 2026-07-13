---
type: task
blocked_by: [01, 04]
claimed_by: claude-code-session-716e8e8c
claimed_at: 2026-07-13T00:00:00Z
assets:
  [
    ../../daily-timeline/assets/day-entry-ledger.html.approved,
    ../../daily-timeline/assets/todays-row.html.approved,
    ../../daily-timeline/assets/upcoming-continuation.html.approved,
  ]
---

# The ledger

## Question

Deliver the spec's whole visual layer, from the three approved assets: the day-entry
partial (five columns, two unbroken rails, 28px rows, the three alignment
invariants), muted empty-day rows, today's gray-900 dot and non-receding empty row
(`{{if eq .Date $.Today}}` against page data), the future continuation with tinted
dates and the capped-at-3 overflow row linking to the Calendar View, month-break
dividers via the adjacent-groups Go helper, and the `NoExpensesEver` template branch
swapping in the verbatim "No expenses yet." block.

This **drops the data-table island from the Daily View** — expense rows become
server-rendered links to the edit page, carrying the `?return=` the danger zone
(previous ticket) already honours. The `?date=` view renders populated days through
the same partial verbatim and keeps its fuller empty state. Delete any prototype
leftovers this replaces.

Blocked by the danger zone as well as the service: landing the ledger first would
remove the daily view's row delete before its demoted home exists.

Tests: only the cheap same-partial assertion (timeline, `?date=` — the fragment
endpoint joins it next ticket). The chrome itself is deliberately untested — verify
by eyeballing the rendered page against the assets on a seeded scratch DB, including
the empty-account, dormant-account, and future-heavy cases.

Done when: the daily view is card-free, JS-free, and matches the approved markup;
`?date=` and the empty states behave per spec; checks green.
