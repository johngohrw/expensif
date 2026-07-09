# Validate expense dates on the write path

Type: grilling
Status: open
Blocked by: none

## Question

`validateExpenseInput` (`internal/service/service.go:47-66`) checks the amount, the
category, and the description. It does not check the date. It only substitutes today
when the date is empty.

Proven by probe, not inference — all returned `201`:

    {"date":"2031-01-01"}   accepted, renders above today
    {"date":"banana"}       accepted, renders as a day header reading "banana"

`templates/form.html:46` is an `<input type="date">` with no `max`, so the browser
won't stop a future date either, and the API has no browser at all.

[What happens to expenses dated after today](./03-future-dated-expenses.md) decided
that garbage dates get fixed at the source rather than defended against on read. This
ticket specifies that fix.

Resolve:

- **What is a valid date?** Must it `time.Parse("2006-01-02")`? Anything else — a
  lower bound, a sanity ceiling?
- **Are future dates valid?** Ticket 03 says *yes* — they render in Upcoming. So
  validation must reject unparseable dates **without** rejecting future ones. Do not
  quietly conflate the two; a `max=` on the form input would break Upcoming.
- **Reject or coerce?** Return a `ErrInvalidDate` sentinel alongside `ErrInvalidAmount`
  and friends, or silently fall back to today? The existing errors reject; follow suit
  unless there's a reason not to.
- **Both write paths.** `CreateExpense` (`service.go:68`) and `UpdateExpense`
  (`service.go:100`) both call `validateExpenseInput`, so one change covers both —
  confirm that, and confirm `isValidationErr` maps the new sentinel to a `400` in
  `handlers_api.go:61`.
- **The form.** Does `form.html` need anything, given `<input type="date">` already
  constrains a browser user to well-formed dates? The gap is the API, not the form.
- **Existing rows.** Ticket 03 accepted the risk that already-written garbage rows
  keep rendering. Is that still acceptable once you look at it directly? A migration
  that nulls or repairs unparseable dates is the alternative. Note the local dev
  database is empty; deployed ones may not be.
- **Tests.** `internal/service` has tests. A table test over the date cases is cheap.

## Scope note

This is a write-path change and the map's destination is a daily-view spec. It is in
scope only because the daily view's Upcoming section cannot be specified while an
unparseable date can outrank every real one. Keep it to date validation; do not let it
grow into a general validation audit.
