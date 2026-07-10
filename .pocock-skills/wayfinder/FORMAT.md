# Format — the shapes of a map, its tickets, and its fog

The disclosed reference for [`wayfinder`](SKILL.md). Write every file to the shapes below, and run the checklist at the bottom before committing anything under `.plan/`.

Structure only what has exactly one correct value: status, type, edges, claims, anchors. A machine can check those, and a second copy of one is a bug waiting to happen. Everything else is prose — see **What stays prose** in the skill.

## The map body

`.plan/<effort-slug>/map.md`. The whole map at low resolution, loaded once per session. Open tickets live as ticket files, found by scanning `tickets/`.

```markdown
# <Map name>

## Destination

<what reaching the end of this map looks like — the spec, decision, or change this effort is finding its way to. One or two lines; every session orients to it before choosing a ticket.>

## Notes

<domain; skills every session should consult; standing preferences for this effort>

## Decisions so far

<!-- the index — one line per resolved ticket: enough to judge relevance, then zoom the link for the detail the ticket holds -->

- [<resolved ticket title>](./tickets/NN-slug.md) — <one-line gist of the answer>

## Not yet specified

<!-- see "Fog of war": in-scope fog you can't ticket yet; graduates as the frontier advances -->

## Out of scope

<!-- see "Out of scope": work ruled beyond the destination; closed, never graduates -->
```

## Tickets

`.plan/<effort-slug>/tickets/NN-<slug>.md`, numbered from `01`; the number is its identity. The **frontmatter holds the ticket's facts** — each has exactly one correct value, and this is the only place any of them is written. The body is the question, sized to one fresh agent session.

```markdown
---
type: research | prototype | grilling | task
status: open | claimed | resolved | out_of_scope
blocked_by: [NN, NN]                 # [] when none
claimed_by: <session or agent id>    # only while status: claimed
claimed_at: <RFC 3339 timestamp>     # only while status: claimed
undermined_by: [NN]                  # optional — see Undermined decisions
assets: [<repo-relative path>]       # optional
---

# <Ticket title>

## Question

<the decision or investigation this ticket resolves>
```

The `## Answer` is appended on resolution, not written up front. Assets created while resolving (research notes, prototype code) are saved in the repo, linked from `assets:`, and not pasted in.

## Fog patches

One bullet per patch in the map's **Not yet specified**. The bolded lead sentence is the patch's **title** — its identity, so it can be referred to and struck once it graduates. Anchor it to the open ticket that will clear it where you know which; leave the anchor off where you don't. Title and anchor are a patch's only structure; the rest is prose, as loose as the view allows.

```markdown
- **<Patch title>.** <prose> <clears-with: NN>
```

## Verify before you commit

The map is shared memory: the next session trusts it without re-deriving it, so drift misleads silently. After resolving a ticket, and before committing anything under `.plan/`, check that:

1. Every `blocked_by` names a ticket that exists, no ticket blocks itself, and there is no cycle.
2. A ticket is `resolved` or `out_of_scope` **iff** it has an `## Answer`.
3. Every `resolved` ticket appears exactly once in **Decisions so far**, and everything linked there is `resolved`.
4. Every `out_of_scope` ticket appears in **Out of scope**, and none appears in Decisions-so-far.
5. Every `claimed` ticket carries `claimed_by` and `claimed_at`, and none of those is stale.
6. Every fog patch's title names a question no live ticket holds, and every `<clears-with: NN>` names a ticket not yet resolved — a patch anchored to a resolved ticket should have graduated into a ticket, or been struck.
7. Each ticket number is used once, and progress is stated nowhere but the frontmatter.

Where the repo has a tool that performs these, run it; the skill needs no tool and assumes none. A tool reading only `.plan/` cannot see the last clause of check 7 — a count written into some file elsewhere in the repo — so that one stays yours.
