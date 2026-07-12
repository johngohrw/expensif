---
type: task
blocked_by: [03]
---

# Architecture docs catch up with what shipped

## Question

Deliver the spec's **Further Notes** doc obligations, now that the code exists:

- Amend ADR 0001 with a note (not a new ADR) that the scroll island is the first
  non-React island — vanilla TS fetching server-rendered HTML — and why that doesn't
  reverse the decision.
- Touch up `AGENTS.md`'s "React Islands" architecture sentence to match.
- Update `CONTEXT.md`: the Daily View definition (date-indexed, every day appears),
  and add the terms the implementation firmed up — candidates from the planning map:
  **Ledger**, **Upcoming**, **Window**, **Rail** — per the `domain-modeling` skill.

Also sweep for leftovers the earlier tickets should have deleted (prototype handlers
or templates, dead template branches, the unused data-table registration on the daily
view) and remove any found.

Done when: a fresh agent reading `AGENTS.md`, `CONTEXT.md`, and the ADRs gets an
accurate picture of the shipped daily view.
