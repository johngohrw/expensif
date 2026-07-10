# 2026-07-10 — Wayfinder audited against upstream, and split into method + adapter

Second session of the day. The first one (`2026-07-10-wayfinder-schema.md`) gave the map
a frontmatter schema and a linter. This one audited that work against Matt Pocock's
original and found the schema had been solving the wrong problem.

## What this session did

**Audited the skill against upstream.** The original lives at
`/Users/rengwu/Desktop/pocock-skills/reference/skills/engineering/wayfinder/SKILL.md` —
a local checkout, not fetched. Read it before touching the skill again.

**Split the skill into method and adapter**, restoring upstream's layering.
`FORMAT.md` is now [`TRACKER-MARKDOWN.md`](../../.pocock-skills/wayfinder/TRACKER-MARKDOWN.md).
**Fixed a parser bug that silently resolved open tickets.** See `9a0c5d6` in the tool repo.

Commits, expensif: `e97857f`, `842a859`, `82280d2`. Tool repo: `9a0c5d6`, `4a209d7`.
Their messages carry the reasoning; the divergences from upstream are listed in
[`.pocock-skills/README.md`](../../.pocock-skills/README.md). Do not re-derive any of it.

**No expensif application code was touched.** `git diff main..HEAD -- '*.go'` is empty.
The daily-timeline effort remains all planning.

## The finding, in one paragraph

Upstream stores tickets **on an issue tracker** and keeps storage mechanics out of the
skill entirely. Our port moved them to `.plan/` markdown and inlined the mechanics. That
is why this skill grew a schema, seven checks and a linter — upstream needs none of them,
because a tracker is a database: ids cannot collide, blocking cannot dangle, an assignee
is the claim. Every check here is a constraint enforced by hand. The method itself never
degraded; it survives near-verbatim.

Corollary worth holding onto: upstream requires *native* blocking specifically because
**it renders the frontier visually in the tracker's own UI**. Pocock's wayfinder already
assumes a GUI. The markdown port removed it. Building the wayfinder interface restores
his intent rather than departing from it.

## State of the tree — read this first

**Both repos are on unmerged branches, not `main`.**

- expensif: `wayfinder-tracker-adapter`
- `/Users/rengwu/Desktop/Projects/wayfinder`: `derive-status-from-body`

Both are fast-forward merges. **If you resolve a ticket from `main` you will be working
against a format that no longer exists** — `main` still has `status:` in every ticket.
Merge or check out the branch before doing anything.

Trees are clean. Tool tests pass; `wayfinder lint` reports the map clean. Expensif's Go
tests were not run and did not need to be.

## Where to pick up

The frontier is derived from the tickets — scan `.plan/daily-timeline/tickets/`. Claim per
the skill; it and its adapter are authoritative on how.

**The frontier will hand you the wrong ticket.** First-by-number offers
[Should HandleDaily's two branches converge](../daily-timeline/tickets/05-converge-handledaily-branches.md),
but [Test strategy for the date-indexed timeline](../daily-timeline/tickets/10-test-strategy.md)
should come first: 05, 06 and 07 all concern how the new query is *consumed*, which is
easier once its contract has tests pinning it.

**The map should say so structurally rather than leaving it in this handoff**, where it
will rot. Since the lower-numbers rule was rejected, a newer ticket may block an older
one: add `10` to the `blocked_by` of 05, 06 and 07 and the frontier states its own order.
The user was offered this and the session ended before an answer. It is a one-line edit
per ticket plus a lint run. Do it, or take 10 and say why.

## Things you would otherwise learn the hard way

- **Two of the four statuses have never been used on a real map.** No ticket has ever
  carried a `## Ruled out` section, and none has ever carried `claimed_by`. Both paths
  have unit tests, written by the same author as the code they test. The first real
  out-of-scope ruling and the first real mid-flight claim are first runs. Watch them.
- **`claimed_at` must be RFC 3339 or lint errors.** Generate it
  (`date -u +%Y-%m-%dT%H:%M:%SZ`); do not type it by hand.
- **The fence rule exists only in the tool.** The adapter requires structural scans to
  ignore fenced code blocks, and the `grep` it prints is fence-blind on purpose, marked as
  a convenience. An agent working without the tool will reach for that grep and it is
  *wrong on exactly the ticket most likely to exist here* — one that quotes the ticket
  format. No cheap fix is known.
- **Do not commit a half-written answer.** An `## Answer` heading with nothing under it is
  a lint error by design: it means a session died mid-write, and the ticket correctly
  stays claimed rather than reading as resolved.
- **`assets:` has a doc/reality mismatch, and nothing checks it.** The adapter says
  repo-relative; the one real usage (ticket 08) is relative to `tickets/`, which is what
  resolves on disk. Decide which is right and say so.
- **Nothing in expensif points at the tool.** `AGENTS.md` still does not mention
  `/Users/rengwu/Desktop/Projects/wayfinder/`, so a fresh agent cannot find it. Deferred
  across three sessions now; the adapter's "run the linter if the repo has one" stays
  aspirational until someone decides.
- **The evidence base is thin, and thinner than the polish suggests.** The linter has run
  clean across one map, one project, whose author also wrote the checker. The fence bug —
  which marked an open ticket resolved and dropped it from the frontier — surfaced only
  because the user asked a passing question about how completion is detected. Assume there
  is another one. Charting a second map, on a subject that is *not* ticket formats, would
  be the single most valuable next test.

## Decisions taken, so you don't relitigate them

Four, all the user's, all in the commit messages and the skill: the interface is a
**read-only viewer**; status stays **derived from the body**, with the scan made
fence-aware; the **method/adapter seam is restored**; and the markdown adapter permits
**one session at a time** (which is what makes the delta check sound). The three
divergences from upstream were kept and are now documented rather than silently carried.

Deliberately not built: the viewer itself, and a second (GitHub) adapter that would prove
the seam.

## Environment

- The tool: `cd ../wayfinder && go run ./cmd/wayfinder {status,lint} ../expensif/.plan/daily-timeline`.
  Dependency-free, its own git repo, not a submodule. `lint` exits 1 on errors.
- `make dev` runs `air` (Go, :8080) and Vite (:8081). `DEV=true` is set in the user's
  shell. The Go server caches templates at startup; restart after template edits.
- The user's dev database (`~/.expensif/expenses.db`) is empty. Prototyping used a scratch
  DB via `DATA_DIR=<tmp> PORT=8090 go run ./cmd/server`. Do not seed the user's database.

## Suggested skills

Apply only if available in your environment; all live in `.pocock-skills/<name>/SKILL.md`
and are **not** auto-discovered — read the file and follow it.

- `wayfinder` — the method, plus `TRACKER-MARKDOWN.md` for the shapes, the fence rule and
  the checklist. Read both before touching `.plan/daily-timeline/`. One ticket per session.
- `grill-me` — how ticket 10 wants to be resolved; it is a grilling ticket with six
  sub-decisions already enumerated in its body.
- `domain-modeling` — terms have firmed up ("the rail", "Upcoming", "day entry", "empty
  day", "undermined", "the adapter") and `CONTEXT.md` still does not exist.
- `codebase-design` — tickets 05 and 07 are seam questions.
- `writing-great-skills` — if you edit the skill again. The audit was run against it.
- `handoff` — when you finish.
