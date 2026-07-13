---
type: task
blocked_by: []
claimed_by: claude-code-session-efb86036
claimed_at: 2026-07-13T00:00:00Z
assets: [../../daily-timeline/assets/edit-delete-danger-zone.html.approved]
---

# The edit page's danger zone

## Question

Deliver the spec's **Delete and the edit page** decisions from the approved asset: a
"Danger zone" at the foot of the edit page as a **sibling** form (the page is already
a form; nesting is invalid HTML) posting to the existing delete route behind a
`confirm()` naming the expense ID. The return path is a server-issued `?return=`: the
edit handler passes it through to a hidden field, and the delete handler validates it
as a local path (begins with `/`, not `//`, no scheme) with a `/` fallback.

Independent of the ledger: nothing sends `?return=` yet (the ledger ticket threads it
into the row's Edit link), but the plumbing is fully verifiable by hand-typed URL.
Touch neither the List View nor its data-table — it is the intentional bulk-delete
surface.

Tests: the delete handler's `?return=` validation (accepted local path, rejected
`//host`, rejected scheme, absent → `/`), following the existing handler-test shape.
Verify live: delete an expense from the edit page on a scratch DB and land where
`?return=` says.

Done when: delete works from the edit page with confirm and a validated return, and
checks are green.
