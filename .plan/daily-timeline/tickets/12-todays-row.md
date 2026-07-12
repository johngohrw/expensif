---
type: prototype
blocked_by: [05, 08]
---

# Today's row

## Question

Is today visually distinguished from any other day on the timeline, and if so, how?

[The day-entry ledger](./08-day-entry-ledger-redesign.md) settled the row unit — five
columns, two rails, 28px, month-break dividers — and made **no provision for today**.
Its approved markup ([`assets/day-entry-ledger.html.approved`](../assets/day-entry-ledger.html.approved))
draws every day identically. On a timeline anchored at today, the anchor is invisible.

The precedent is next door and it disagrees. `HandleCalendar` carries an `IsToday` flag
through to `CalendarCell.IsToday` and the calendar renders today distinctly. The ledger
has no equivalent, and [ticket 04](./04-date-indexed-daily-groups.md) deliberately did
**not** add one — it deferred the field to this ticket, noting the handler already knows
today and can pass it in page data for `{{if eq .Date $.Today}}` with no change to the
service signature. That escape hatch is still open; this ticket decides whether to use it.

This is a prototype ticket: produce markup against the approved ledger and react to it.

> **"Whether" is closed — [ticket 14](./14-upcoming-section-chrome.md) settled it.** There
> is **no Upcoming section**: the ledger continues past today into the future with no
> heading, no divider and no panel. Today's row is therefore *the only thing* separating
> "hasn't happened yet" from "already spent", and ticket 14 **requires** it to be
> distinguished. The alternative — letting the future days' blue-tinted dates carry the
> boundary alone — was put to the user and rejected: on a page with no future expenses
> (the common case) the tint never appears, so the boundary would exist only on pages that
> don't need one.
>
> **This ticket now decides only *how*.** The first bullet below is struck.

Resolve, one at a time:

- ~~**Distinguished at all?**~~ **Closed by ticket 14 — yes, necessarily.** A timeline that
  always starts at today does *not* simply start at today: when Upcoming is non-empty,
  future days sit above it, which is exactly the case that breaks "top row = today".
- **Distinguished *differently* when empty?** Today with no expenses yet is the single
  most common state of the most important row — it is the row the user came to fill.
  Ticket 01's muted empty line is designed to recede into the background. Today's empty
  row receding is either correct (nothing to report) or exactly wrong (the primary call
  to action, greyed out).
- **How, within the ledger's language.** Rail weight, a date-gutter treatment, a label
  ("Today"), or type weight — but not a card, not a background fill that breaks the two
  unbroken rails, and nothing hover-only (ticket 01 rejected hover for mobile).
- **It must not collide with the future tint.** Ticket 14 spends `text-blue-500` on future
  dates. Whatever marks today has to read as *distinct from* "upcoming", not as more of it
  — the two sit adjacent, with today's row directly below the nearest future day.
- **The second consumer.** [Ticket 05](./05-converge-handledaily-branches.md) renders
  `?date=` through the *same* day-entry partial. So `?date=2026-07-12` — today, reached
  from the calendar — draws whatever this ticket decides. Is that right, or does the
  marker belong to the timeline only, making it a template parameter rather than a
  property of the day?
- **What carries it.** Page data (`$.Today`, compared in the template) versus a field on
  `DailyGroup`. Ticket 04 argued for the former; the second consumer above is the fact
  that could overturn it.
