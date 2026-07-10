---
name: wayfinder
description: Plan a huge chunk of work — more than one agent session can hold — as a shared map of investigation tickets under `.plan/`, and resolve them one at a time until the way to the destination is clear.
disable-model-invocation: true
---

A loose idea has arrived — too big for one agent session, and wrapped in fog: the way from here to the **destination** isn't visible yet. Wayfinding is about finding that way, not charging at the destination. This skill charts the way as a **shared map** of markdown files under `.plan/`, then works its tickets one at a time until the route is clear.

The destination varies per effort, and naming it is the first act of charting — it shapes every ticket. It might be a spec to hand off and iterate on, a decision to lock before planning starts, or a change made in place like a data-structure migration. The map is domain-agnostic — engineering work, course content, whatever fits the shape.

## Plan, don't do

Wayfinder is **planning** by default: each ticket resolves a decision, and the map is done when the way is clear — nothing left to decide before someone goes and does the thing. The pull to just do the work is usually the signal you've reached the edge of the map and it's time to hand off. An effort can override this in its **Notes** — carrying execution into the map itself — but absent that, produce decisions, not deliverables.

## Refer by name

Every map and ticket is a file with a **name** — its title. In everything the human reads — narration, the map's Decisions-so-far — refer to it by that name, never by a bare number or filename. A wall of `03, 04, 05` is illegible; names read at a glance. The path doesn't vanish — a name wraps its link — but it rides *inside* the name, never stands in for it.

## The Map

The map lives at `.plan/<effort-slug>/map.md` — the canonical artifact. Its tickets are files under `.plan/<effort-slug>/tickets/`, one per ticket. The whole directory is committed to version control — it is the shared memory that concurrent and future sessions orient to, so commit map and ticket changes promptly.

The map is an **index**, not a store. It lists the decisions made and points at the tickets that hold their detail; a decision lives in exactly one place — its ticket — so the map never restates it, only gists it and links.

Two rules keep the index honest, because a map that has drifted misleads every session that trusts it, and it does so silently:

- **The map may restate a decision's gist. It may never restate a ticket's status, type, or edges.** A gist is prose, lossy on purpose, and cannot be derived — that is why the map reads in one pass. The facts live in the ticket's frontmatter, and a fact with two homes goes stale in one of them.
- **Never state current progress outside a ticket's frontmatter** — no counts, no statuses, in a README, a contributor guide, or anywhere else that claims to say where the effort stands today. "Four of nine resolved" is true for a week and wrong forever after, and the reader who trusts it cannot tell. Point at the map, or count when asked. A *dated* handoff recording what one session did is history rather than state, and is exempt: it does not drift, because it never claimed to be current.

### The map body

The whole map at low resolution, loaded once per session. Open tickets are **not** listed — they are open ticket files, found by scanning `tickets/`.

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

### Tickets

Each ticket is a file at `.plan/<effort-slug>/tickets/NN-<slug>.md`, numbered from `01`; the number is its identity. Its **frontmatter holds the ticket's facts** — each has exactly one correct value, and this is the only place any of them is written. Its body is the question, sized to one fresh agent session:

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

A session **claims** a ticket by setting `status: claimed` with `claimed_by` and `claimed_at`, and saving, **first**, before any work, so concurrent sessions skip it. Commit the claim promptly — on a shared repo, other clones only see it once it lands. A claim carrying no timestamp cannot be told apart from a session that died mid-ticket, and the frontier will step around that ticket forever.

Blocking is the `blocked_by` list. A ticket is **unblocked** when every ticket it lists is `resolved`; the **frontier** is the open, unblocked, unclaimed tickets — the edge of the known. To find it, scan `tickets/`; first by number wins.

`out_of_scope` is closed but it is **not** resolved, and it never satisfies a `blocked_by` edge. A ticket blocked by an out-of-scope ticket can therefore never unblock — which means one of the two is mis-scoped, and you should say which.

The answer isn't part of the body — it's recorded on resolution (see [Work through the map](#work-through-the-map)). Assets created while resolving a ticket (research notes, prototype code) are saved in the repo, linked from the ticket's `assets:` list, and not pasted in.

### Undermined decisions

A resolved ticket sometimes rests on a premise that a later ticket destroys. Don't silently re-decide it, and don't leave it looking sound: add `undermined_by: [NN]` to the ticket whose premise broke, naming the ticket that broke it, and open its `## Answer` with a line saying so. The decision still stands — nobody has reopened it — but every session that reads the map can now see what it is standing on. A decision recorded as simply `resolved` reads as settled, and that is precisely how a live problem gets laundered into a checkmark.

## Ticket Types

Every ticket is either **HITL** — human in the loop, worked *with* a human who speaks for themselves — or **AFK**, driven by the agent alone. A HITL ticket only resolves through that live exchange; the agent never stands in for the human's side of it (a grilling agent that answers its own questions has broken this).

- **Research** (AFK): Reading documentation, third-party APIs, or local resources like knowledge bases. Use the `research` skill: investigate against primary sources and create a cited markdown summary as a linked asset. Use when knowledge outside the current working directory is required.
- **Prototype** (HITL): Raise the fidelity of the discussion by making a cheap, rough, concrete artifact to react to — an outline, a rough take, a stub, or UI/logic code via the `prototype` skill. Links the prototype as an asset. Use when "how should it look" or "how should it behave" is the key question.
- **Grilling** (HITL): A relentless interview — one question at a time, a recommended answer offered with each, decision branches resolved one by one, facts looked up rather than asked — with the `domain-modeling` skill applied so terms and decisions are captured as they crystallise. The default case.
- **Task** (HITL or AFK): Manual work that must happen before a *decision* can be made — nothing to decide, prototype, or research, but the discussion is blocked until it's done. Signing up for a service so its API can be judged, provisioning access, moving data so its shape can be seen. This is the one type that *does* rather than decides — and it earns its place by unblocking a decision, not by delivering the destination. The agent drives it alone where it can (AFK); otherwise it hands the human a precise checklist (HITL). Resolved when the work is done; the answer records what was done and any resulting facts (credentials location, new URLs, row counts) later tickets depend on.

## Fog of war

The map is _deliberately_ incomplete: don't chart what you can't yet see. Beyond the live tickets lies the **fog of war** — the dim view of decisions and investigations you can tell are coming but can't yet pin down, because they hang on questions still open. Resolving a ticket clears the fog ahead of it, graduating whatever's now specifiable into fresh tickets — one at a time, until the way to the destination is clear and no tickets remain.

The map's **Not yet specified** section is where that dim view is written down: the suspected question, the area to revisit later. It's the undiscovered frontier _toward_ the destination — everything here is in scope, just not sharp enough to ticket. Write as loosely or as fully as the view allows; it doubles as a signpost for collaborators reading where the effort is headed.

Each patch is one bullet whose **bolded lead sentence is its title** — the patch's identity, so it can be referred to and struck once it graduates. Where you already know which open ticket will clear it, anchor it; leave the anchor off when you don't. The title and the anchor are a patch's only structure — everything after them is prose, as loose as the view allows.

```markdown
- **<Patch title>.** <prose> <clears-with: NN>
```

A patch's title must not restate a live ticket. That tracks one question in two places, once sharply and once vaguely, and the vague copy is the one that rots.

**Fog or ticket?** The test is whether you can state the question precisely now — _not_ whether you can answer it now.

- **Ticket when** the question is already sharp — even if it's blocked and you can't act on it yet.
- **Not yet specified when** you can't yet phrase it that sharply. Don't pre-slice the fog into ticket-sized pieces: it's coarser than a ticket, and one patch may graduate into several tickets, or none, once the frontier reaches it.

**Not yet specified** excludes what's already decided (Decisions so far), what's already a live ticket, and what's out of scope (the next section).

## Out of scope

Fog only ever gathers _toward_ the destination. The destination fixes the scope, so work beyond it is **out of scope** — it isn't fog, and it doesn't belong in **Not yet specified**. It gets its own **Out of scope** section on the map: work you've consciously ruled out of _this_ effort. Scope, not sharpness, lands it here.

Out-of-scope work never graduates — the frontier stops at the destination — so it returns only if the destination is redrawn, and then as a fresh effort, not a resumption.

Ruling something out of scope is a scoping act, not a step on the route. When a ticket that already exists turns out to sit past the destination — mis-scoped in while charting, or exposed by a resolution — **close it**: set `status: out_of_scope` with a one-line `## Answer` saying it was ruled out (a closed ticket is unambiguously off the frontier), and leave one line in the map's **Out of scope** section: the gist plus why it's out of scope, linking the closed ticket. It stays out of **Decisions so far**, which records the route actually walked — a scope boundary isn't a step on it.

`out_of_scope` is its own status rather than a flavour of `resolved` for exactly this reason: anything counting the decisions made would otherwise count a boundary as a step.

## What stays prose

Structure only what has exactly one correct value: status, type, edges, claims, anchors. A machine can check those, and a second copy of one is a bug waiting to happen.

Everything else stays unstructured, permanently. The **Destination**, the **Notes**, a ticket's **Question**, its **Answer**, and a decision's one-line gist in Decisions-so-far are lossy human summaries — nothing can derive them, and no field should try to hold them. Give the prose fields a schema and you have built a ticket tracker and thrown away the thing this skill is for.

## Verify before you commit

The map is shared memory: the next session trusts it without re-deriving it, so drift misleads silently. After resolving a ticket, and before committing anything under `.plan/`, check that:

1. Every `blocked_by` names a ticket that exists, no ticket blocks itself, and there is no cycle.
2. A ticket is `resolved` or `out_of_scope` **iff** it has an `## Answer`.
3. Every `resolved` ticket appears exactly once in **Decisions so far**, and everything linked there is `resolved`.
4. Every `out_of_scope` ticket appears in **Out of scope**, and none appears in Decisions-so-far.
5. Every `claimed` ticket carries `claimed_by` and `claimed_at`, and none of those is stale.
6. No fog patch's title restates a live ticket, and no `<clears-with: NN>` names a ticket already `resolved` — that patch should have graduated into a ticket, or been struck.
7. No ticket number is used twice, and nothing outside a ticket's frontmatter states a status or a count.

It is a checklist, not a ceremony — seven greps. Where the repo has a tool that performs it, run that instead; this skill needs no tool and assumes none.

## Invocation

Two modes. Either way, **never resolve more than one ticket per session.**

### Chart the map

User invokes with a loose idea.

1. **Name the destination.** Run a grilling session (see the Grilling ticket type — one question at a time, recommended answers, `domain-modeling` applied) to pin down what this map is finding its way to — the spec, decision, or change. The destination fixes the scope, so it's settled first.
2. **Map the frontier.** Grill again, **breadth-first** this time: fan out across the whole space rather than deep on any one thread, surfacing the open decisions and the first steps takeable now. **If this surfaces no fog** — the way to the destination is already clear, the whole journey small enough for one session — you don't need a map. Stop and ask the user how they'd like to proceed.
3. **Create the map** at `.plan/<effort-slug>/map.md`: Destination and Notes filled in, Decisions-so-far empty, the fog sketched into **Not yet specified**.
4. **Create the tickets you can specify now** as files under `tickets/` — then wire the `Blocked by:` lines in a **second pass** (tickets need numbers before they can reference each other). Wiring sorts them into the frontier and the blocked; everything you can't yet specify stays in the fog — the **Not yet specified** section.
5. Stop — charting the map is one session's work; do not also resolve tickets.

### Work through the map

User invokes with a map (a path under `.plan/`, or the effort's name). A ticket is **optional** — without one, you pick the next decision, not the user.

1. Load the **map** — the low-res view, not every ticket body.
2. Choose the ticket. If the user named one, use it. Otherwise take the first frontier ticket in order. **Claim it**: set `status: claimed` with `claimed_by` and `claimed_at` before any work.
3. Resolve it — **zoom as needed**: read the full body of any related or resolved ticket on demand; invoke the skills the `## Notes` block names. If in doubt, grill (with `domain-modeling` applied).
4. Record the resolution: append the answer under an `## Answer` heading in the ticket file, set `status: resolved`, clear `claimed_by` and `claimed_at`, and **append a context pointer** to the map's Decisions-so-far (gist + link).
5. Add newly-surfaced tickets (create-then-wire); graduate any fog the answer has made specifiable, clearing each graduated patch from **Not yet specified** so it lives only as its new ticket. If the answer reveals a ticket — this one or another — sits beyond the destination, **rule it out of scope** rather than resolving it on the route. If the decision invalidates other parts of the map, update or delete those tickets; if it breaks the premise of a decision already made, mark that ticket `undermined_by` rather than quietly re-deciding it.
6. **Verify before committing** — see the checklist below.

The user may run unblocked tickets in parallel, so expect other sessions to be editing `.plan/` concurrently.
