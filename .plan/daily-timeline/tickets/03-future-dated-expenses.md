# What happens to expenses dated after today

Type: grilling
Status: open
Blocked by: none

## Question

Anchoring the timeline at today means a future-dated expense has nowhere to render.
Is that acceptable, and if not, what gives?

This is not hypothetical. `HandleCalendar` deliberately extends its range to
`today.AddDate(1, 0, 0)` (`internal/web/handlers_html.go:151`) — a full year forward.
Nobody pads a calendar with a year of empty future months by accident; the app
expects future-dated expenses to exist. Confirm that reading before deciding.

Today's daily view has the bug in a quieter form: `DailyGroups` sorts by date
descending and takes the last 100 expenses, so a future-dated expense *does* appear,
at the top, above today. Once today becomes the anchor, it silently vanishes.

Resolve:

- Can an expense actually be created with a future date? Check `HandleAdd`
  (`handlers_html.go:292`) and the form for validation. **Look this up, do not ask.**
- If yes: does the daily view show future days when future expenses exist (a ragged
  top edge), pin them into a "Scheduled" section above today, or keep them out and
  accept that the calendar is where you see them?
- If no future expenses are creatable, is the calendar's forward year dead code, and
  does that change the answer to anything else?

The answer constrains the query shape (whether the window's end is `today` or
`max(today, latest expense)`), so it blocks the `DailyGroups` rework.
