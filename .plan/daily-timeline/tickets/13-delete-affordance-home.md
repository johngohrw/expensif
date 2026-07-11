---
type: prototype
blocked_by: [08]
claimed_by: claude-code-session-2026-07-12
claimed_at: 2026-07-11T18:12:00Z
---

# Where delete lives once the ledger lands

## Question

[The day-entry ledger](./08-day-entry-ledger-redesign.md) demoted delete "to the edit
page". **The edit page has no delete affordance.** Verified, not assumed:

- `templates/edit.html` is **nine lines** — a page title and a `<form>` wrapping the
  shared `form` partial. `grep -i delete` returns nothing.
- Delete exists today **only inside the `data-table` island's actions column**
  (`templates/daily.html:25`): a `type: "form"` action posting to
  `/expenses/delete/{id}`, `variant: "danger"`, with a `confirm` dialog
  (`"Delete expense #{id}?"`).
- Ticket 08 **drops `data-table` from the daily view**. So delete does not *move* —
  it **disappears**, and the page it was demoted to has nothing to receive it.
- The route survives regardless: `POST /expenses/delete/{id}` → `HandleDelete`
  (`internal/web/server.go:30`). `list.html` keeps its `data-table`, so delete still
  exists *there*. This ticket is about the daily view's loss of it, not the route's.

So the demotion was a decision to build something that does not exist. This ticket
builds it — or overturns the demotion.

This is a prototype ticket: put an affordance on `edit.html`, look at it, react.

Resolve, one at a time:

- **Does delete belong on the edit page at all?** Deleting via a page whose title is
  "Edit Expense #12" and whose submit button says "Save Changes" is a mode error waiting
  to happen. The alternative is that ticket 08's demotion was wrong and the ledger row
  needs *some* destructive affordance — but that reopens an approved design, so it needs
  a real argument, not a preference.
- **What it looks like on a nine-line page.** Danger button below the form? A separate
  form (it must be — a nested `<form>` is invalid HTML, and the page is already one
  `<form method="POST">`)? Where, relative to "Save Changes", so it is reachable but not
  mis-clickable?
- **The confirm step.** The current dialog is a JS `confirm()` supplied by the island.
  Ticket 08's whole point is that the daily view becomes **JS-free**. Does the edit page
  keep a JS confirm (it may — it is a different page), or does delete become a two-step
  server-rendered confirmation, or does it need no confirm at all once it is off the
  one-click table row?
- **The return path.** `HandleDelete` redirects — to where? From the daily view, back to
  the daily view was implicit. From the edit page, "back to daily" is a guess: the user
  may have arrived from `/list`, or from `?date=`.
- **Discoverability.** Delete moves from one click on the row to: notice the row, click
  Edit, find delete on the next page. For the "I fat-fingered the amount" case that is
  fine. Confirm there is no case where it is not.
