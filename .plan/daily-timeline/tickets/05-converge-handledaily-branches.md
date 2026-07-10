---
type: grilling
blocked_by: [04]
---

# Should HandleDaily's two branches converge

## Question

`HandleDaily` (`internal/web/handlers_html.go:237`) is two handlers wearing one coat.

The `?date=` branch calls `ExpensesInRange(date, date)`, sums a total by hand, and
hand-builds a one-element `[]domain.DailyGroup`. The unfiltered branch calls
`DailyGroups(ctx, 100)`. They share only the currency-conversion loop that runs after
them both.

Once ticket 04 gives the service a date-windowed call, the filtered branch is a
window of size 1. The duplication becomes obvious and removable — *or* the two
branches turn out to be genuinely different pages that happen to share a template.

Resolve:

- Is `?date=` a window of one day, or a different view? It has a back-link
  (`BackHref`/`BackLabel`, set from `?from=calendar`) that the timeline has no
  concept of, and it is the only current caller of the empty-day markup.
- If it converges: does a single-day window with no expenses still show the muted
  card, or the fuller "No expenses for this day" state it shows today? Two different
  designs for the same zero state is a smell — but landing on the calendar's date
  cell and seeing a *muted* card may read as a broken link.
- If it does not converge, say so plainly in the spec and leave the branch alone.
  Convergence is not a goal; a smaller total surface is.

The answer determines whether ticket 01's design has one consumer or two.
