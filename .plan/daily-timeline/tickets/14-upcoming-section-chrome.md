---
type: prototype
blocked_by: [03, 08]
---

# The Upcoming section's chrome

## Question

What does the Upcoming section *look* like?

[What happens to expenses dated after today](./03-future-dated-expenses.md) fixed its
**content**: future-dated expenses render in an "Upcoming" section above the timeline,
grouped by day, **not** gap-filled (no walk into next year), ordered so the furthest-out
day sits at the top and the day nearest today sits last, flowing into today's row below.
It fixed none of its chrome. `UpcomingGroups(ctx, after)` already exists as a decision
([ticket 04](./04-date-indexed-daily-groups.md)); nothing renders it.

The raw material is [ticket 08](./08-day-entry-ledger-redesign.md)'s ledger and its
approved markup ([`assets/day-entry-ledger.html.approved`](../assets/day-entry-ledger.html.approved)).
The constraint that shapes everything: **the common case is empty.** Most users have no
future-dated expenses, so whatever this section is, its usual state is nothing at all.

This is a prototype ticket: draw it against the approved ledger and react.

Resolve, one at a time:

- **The empty case first, because it is the normal case.** Does the section vanish
  entirely (no heading, no divider, no reserved space), or does it persist as a
  zero-height affordance? Vanishing is almost certainly right — but say it, because the
  timeline's top edge then changes shape depending on data, and
  [ticket 12](./12-todays-row.md) is deciding whether "top row" means "today".
- **Heading and divider.** Ticket 08's ledger has exactly one horizontal rule in its
  vocabulary — the month-break divider. Does "Upcoming" reuse it, or does a section
  boundary need a heavier mark than a month boundary? Two competing dividers a few rows
  apart is a real risk when an upcoming expense sits in next month.
- **Do upcoming days carry a `+` row?** Ticket 01's empty ledger line always shows `+`.
  But Upcoming is *ungapped* — an empty future day does not exist in it, so there is no
  empty row to carry a `+`, and there is no affordance to add an expense on a specific
  future day from this section. Is that a gap, or correct (that is what the calendar and
  `?date=` are for)?
- **Does it collapse when long?** A user who logs a year of rent has twelve upcoming
  days pushing today's row below the fold — inverting the page's whole purpose. A cap, a
  "show more", or nothing?
- **The two rails.** Ticket 08's rails are described as *unbroken*. Does Upcoming share
  them (one continuous rail through the section boundary into the timeline), or is it a
  separate rail system? This is the decision that makes Upcoming feel like part of the
  ledger or a box sitting on top of it.
