# Window size and how older days load

Type: grilling
Status: claimed
Blocked by: none

## Question

How many days is one window, and how does the next window back arrive?

The bounds are settled: anchor at today, walk backwards, no future days. What is not
settled is `N` and the loading mechanism. `DailyGroups(ctx, 100)`'s expense cap
cannot survive — 100 expenses is an unbounded number of days, and a window of days is
an unbounded number of expenses.

Resolve, one at a time:

- **N.** 30 days? A calendar month? Enough to fill a viewport and no more? Note that
  N is now the page-weight knob: every day in the window costs a card whether or not
  it has an expense.
- **Mechanism.** A `?before=YYYY-MM-DD` query param and a plain "Load older" link at
  the foot (server-rendered, no JS, matches the standing preference)? Or an infinite
  scroll island? The map rules out islands for *collapsing*; that ruling does not
  automatically extend to pagination — but it sets the burden of proof.
- **Append or replace.** Does "Load older" navigate to a new page showing the older
  window, or append to the current one? The former is trivial and stateless; the
  latter needs an island or `hx-`-style swap.
- **Does the current 100-expense limit have a defender?** Check whether anything else
  depends on `DailyGroups`' signature before changing it.

This ticket blocks the query shape, so it also owes an answer to: what does the
handler need to ask the repository for — a window of expenses, or a window of days?
