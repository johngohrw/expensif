---
type: grilling
status: open
blocked_by: [02, 03]
---

# Re-shape DailyGroups around dates, not expenses

## Question

What is the new signature and body of the service call behind the daily view?

`DailyGroups(ctx, limit int)` (`internal/service/service.go:143`) is expense-indexed:
list N expenses, bucket into a map keyed by date, sort descending. Days without
expenses cannot exist in a map built by iterating expenses. The gap-filling loop has
to come from somewhere that knows about *dates*.

`HandleCalendar` already writes this loop —
`for d := first; !d.After(last); d = d.AddDate(0, 0, 1)` — but it writes it *in the
handler*, reading from a `dayTotals` map it built itself. So the codebase has a
working answer that lives at the wrong layer.

Resolve, applying `codebase-design` (this is a module-depth question):

- Signature: `DailyGroups(ctx, end string, days int)`? `DailyGroupsInRange(ctx,
  start, end string)`? Something that takes the window from ticket 02 directly?
- Where does the gap-filling loop live — service, or handler like the calendar does?
  If service, should `HandleCalendar` be refactored to share it, or is that a
  different effort? (Be honest: shared-helper pressure is real, but so is coupling
  two views that are diverging.)
- Does `DailyGroup` (`internal/domain/models.go:34`) need an `IsEmpty` or `IsToday`
  field, or is `len(Expenses) == 0` enough for the template?
- Currency conversion currently happens in the handler after the fact
  (`handlers_html.go:271-285`), looping over groups. Empty groups have nothing to
  convert. Does that loop need guarding, or does it fall out clean?
- Timezone: "today" comes from `nowInTZ(prefs.Timezone)` but `Expense.Date` is a bare
  `2006-01-02` string. Where is the window's end computed, and in whose zone?
