# 2026-07-10 — Wayfinder gets a schema, a linter, and one resolved ticket

## What this session did

Two efforts, one feeding the other.

**Made the wayfinder map machine-readable**, so a tool (eventually a GUI) can synthesise
its state without an LLM. Ticket headers moved into YAML frontmatter; three new fields
appeared; the skill grew a facts-vs-prose rule and shed 73 lines to a new `FORMAT.md`.
See `3af8f1e`, `c9097b9`, `0059ccb`.

**Resolved [Re-shape DailyGroups around dates](../daily-timeline/tickets/04-date-indexed-daily-groups.md)**,
the keystone of the daily-timeline map. See `3be8db4`, `f13356b`.

**No Go code was written or changed in this repo.** `git diff ae84568..HEAD -- '*.go'` is
empty. The daily-timeline effort is still all planning.

Do not re-read the decisions here — they live in their artifacts. The map's
Decisions-so-far indexes them; `.pocock-skills/wayfinder/SKILL.md` and its
`TRACKER-MARKDOWN.md` adapter hold the format contract.

## The tool lives outside this repo

**`/Users/rengwu/Desktop/Projects/wayfinder/`** — a dependency-free Go CLI, its own git
repo, not a submodule and not referenced from anywhere in expensif. Nothing here will
lead you to it.

```
cd ../wayfinder
go run ./cmd/wayfinder status ../expensif/.plan/daily-timeline
go run ./cmd/wayfinder lint   ../expensif/.plan/daily-timeline   # exit 1 on errors
```

`lint` checks the invariants in
[`TRACKER-MARKDOWN.md`](../../.pocock-skills/wayfinder/TRACKER-MARKDOWN.md); `status`
prints the frontier.
It parses the old loose-header format too, warning on it. Module path is bare
(`module wayfinder`), so `go install` from a remote needs a real path first.

`FORMAT.md` says "where the repo has a tool that performs these, run it." In this repo
that sentence is currently aspirational, because nothing points at the tool. Deciding
whether to hardcode a sibling path into `AGENTS.md` was left to the user.

## Where to pick up

The frontier is four tickets: 05, 06, 07 and the new
[Test strategy for the date-indexed timeline](../daily-timeline/tickets/10-test-strategy.md).
First-by-number says 05.

**Take 10 instead**, and say why you did. The daily view has zero test coverage, ticket
04 deleted the only function that would have needed characterising, and 05/06/07 all
concern how the new query is *consumed* — easier once its contract has tests pinning it.
10 also inherited the sharpest brief, since 04 handed it a real interface.

Claim the ticket, and commit the claim, before doing any work — the skill and its
adapter are authoritative on how.

## Design work proposed, deliberately not implemented

The user asked for analysis, then said not to implement. Do not redo the analysis;
either execute it or leave it.

**Superseded later the same day. Execute none of it.** A subsequent session took the
unrepresentability route instead: `status` is deleted as a stored field and derived from
the body's `## Answer` / `## Ruled out` heading, and deleting a ticket is now banned.
That subsumes the first proposal — there is no status left to flip — and the second was
**rejected**, because it forbids a newly-created ticket from blocking an older one, an
edge this very map wanted when ticket 10 appeared. The third is fixed.
`TRACKER-MARKDOWN.md` is authoritative; the bullets below are kept only as the
reasoning that led there.

- **Flip `status: resolved` last.** Step 4 of *Work through the map* lists "append the
  answer, set status, append the map line" — an agent (me) set the status first,
  intending to write the answer next, and produced three lint errors. If the status flip
  is the *final* act, a resolved ticket with no `## Answer` is unreachable on the happy
  path. Checks 2 and 3 stop being checked and become unrepresentable.
- **`blocked_by` names only lower numbers.** Verified true for all ten tickets today,
  and structurally so: a ticket is created after its blockers exist. State it as a rule
  and a `blocked_by` cycle becomes unrepresentable, retiring half of check 1. It also
  gives the ticket numbering a real job.
- **The checklist's closing line overpromises.** "The skill needs no tool and assumes
  none" is true of intent, false of reliability — see the next section.

## Things you would otherwise learn the hard way

- **The checklist does not reliably self-maintain.** Evidence, all from this session: I
  violated three of its seven checks one turn after writing them, with them in context;
  a rename left a stale `Blocked by:` in Chart step 4; the tool's own README went stale
  within minutes of my writing a rule against that; `AGENTS.md` had said "4 of 9
  resolved" for a week when it was 5. Checks 1, 3, 4 and 6 are whole-graph — they need
  state across N files and degrade as N grows. Check 7's last clause ranges over the
  whole repo and **neither the agent nor the linter can discharge it**, since the linter
  reads only `.plan/`.
- **A fog patch anchored to a resolved ticket is a lint error, by design.** It forces
  graduation. This fired on me for real, not in a test.
- **Dated handoffs may be corrected where they carry live instructions**, but not where
  they record what a session did. The `2026-07-09` handoff told the next session to set
  `Status: claimed` — the pre-frontmatter spelling — and `AGENTS.md` sends fresh agents
  to the handoffs first. Fixed in `0059ccb`. The `session-*.md` archive is never
  rewritten.
- **Ticket identity is still a bare number.** Two parallel sessions both reach for the
  next one and `blocked_by` goes ambiguous. The linter *detects* the collision, which is
  after both have committed. Slug identity would fix it and would ripple through every
  `blocked_by` in the skill. Unfixed; it only bites with real concurrency, and this map
  has never had two sessions at once.
- **`undermined_by` exists because a green checkmark lies.** Ticket 02's island decision
  rests on a premise ticket 08 destroyed. It still stands — nobody reopened it — and
  `wayfinder status` now prints it under "Undermined". Do not silently re-decide it.
- **The GUI is unbuilt on purpose.** The parser exists so a canvas *could* render the
  graph. The linter has now run clean across exactly one map, through one session, whose
  author also wrote the checker. That is a start, not evidence.

## State of the tree

Both repos clean. `wayfinder lint` reports the map clean, 10 tickets, no drift; its own
`go test ./...` passes. Expensif's Go tests were not run and did not need to be — no Go
changed.

## Environment

- `make dev` runs `air` (Go, :8080) and Vite (:8081). `DEV=true` is set in the user's
  shell. The Go server caches templates at startup; restart after template edits.
- The user's dev database (`~/.expensif/expenses.db`) is empty. Prototyping used a
  scratch DB via `DATA_DIR=<tmp> PORT=8090 go run ./cmd/server`. Do not seed the user's
  database.

## Suggested skills

Apply only if available in your environment; all live in `.pocock-skills/<name>/SKILL.md`
and are **not** auto-discovered — read the file and follow it.

- `wayfinder` — the map's protocol, plus `TRACKER-MARKDOWN.md` for the schema and the
  seven checks. Read both before touching `.plan/daily-timeline/`. One ticket per session.
- `writing-great-skills` — if you act on the deferred design work above; the proposals
  came out of auditing wayfinder against it.
- `domain-modeling` — terms have firmed up ("the rail", "Upcoming", "day entry", "empty
  day", "undermined") and `CONTEXT.md` still does not exist.
- `codebase-design` — tickets 05 and 07 are seam questions; 04's answer used its
  deletion test to keep `HandleCalendar` out of the refactor.
- `handoff` — when you finish.
