# 2026-07-10 — The island's contract, and the drift that turned out not to exist

Third session of the day, and the first to touch `.plan/daily-timeline/` rather than the
skill itself. The two earlier ones (`2026-07-10-wayfinder-schema.md`,
`2026-07-10-wayfinder-tracker-adapter.md`) reshaped the wayfinder tooling; this one used it.

**Still no application code.** `git diff main..HEAD -- '*.go' '*.tsx'` is empty. The
daily-timeline effort remains entirely planning, and `.plan/daily-timeline/spec.md` — the
map's destination — does not exist yet.

## What this session did

Resolved two tickets, and undermined two in the process.

- [The infinite-scroll island's contract](../daily-timeline/tickets/07-scroll-island-contract.md)
- [Contain the day-card chrome drift](../daily-timeline/tickets/06-day-card-chrome-drift.md)

Commits: `f79e513` (07), `f5d89bf` (06), with claims at `c782f77` and `27e8083`. The
commit messages carry the reasoning and the tickets carry the detail. **Do not re-derive
either.**

## The finding, in one paragraph

Ticket 02 rejected HTML fragments for exactly one stated reason: `data-table.tsx` scans the
document once at module load, so appended `[data-table-root]` nodes would never mount.
Ticket 08 then removed `data-table` from the daily view, and **nobody re-read 02's rejection
clause against it.** An appended ledger day contains no islands at all — every row in the
approved markup is a plain `<a>` and delete moved to the edit page — so the endpoint can
serve HTML from the same Go partial that renders the first window. The ledger gets one
implementation in one language, and ticket 06 does not get answered so much as dissolved.

## Read this before you touch the map

**Two tickets are marked `undermined_by: [06]`**, and `wayfinder status` will show them:

- **02** — its "island fetches JSON and renders in React" clause, and its HTML-fragment
  rejection. The 30-day window, the island itself, and termination at the earliest expense
  all stand.
- **07** — the display-ready DTO, dropping `convertedTotal`, and `mountIsland`/`createRoot`.
  The cursor, the sentinel with one fetch in flight, manual retry on error, and the refusal
  of a `<noscript>` fallback all stand, and are cheaper now.

Both answers open with a paragraph saying what broke. **07 was resolved earlier in this same
session and undermined an hour later by 06.** Read its answer with the substitution in mind
rather than skipping to its sections.

Ticket 02 also carries an *earlier* undermining, by 08, which 07 discharged by re-asking the
user rather than removing the marker silently. That history is in 02's answer.

## The lesson, which cost a ticket

The user overrode the skill's **one ticket per session** rule to take 06 immediately after
07. It was flagged as a risk before starting, with the specific failure named: the agent had
authored 07's payload decision *and* written 06's update section narrowing its scope, and
would then grill 06 from inside its own conclusion. That is precisely what happened. Reading
02's rejection clause against 08's consequence is the first move a session opening 06 cold
would make; this one got there only on reading the approved markup, after 07 was committed.

The rule is not ceremony. **Do not resolve 05 and 10 in one sitting**, however small they
look — and 10 in particular now depends on decisions this session reversed.

## Where to pick up

The frontier is derived from the tickets — run `wayfinder status`, or scan
`.plan/daily-timeline/tickets/`. Claim per the skill and its adapter; they are authoritative.

**The first-by-number rule is fine here** — unlike the last handoff, there is no hidden
ordering to warn about. The two tickets this session left open are worth knowing about:

- **05** unblocks the most fog: the `?date=` patch in **Not yet specified** is anchored to it,
  and 06's HTML-fragment endpoint gives it a new consideration — `/daily/older` renders the
  same partial, so a converged `HandleDaily` has one more caller to think about.
- **10** inherited a simplification. The both-paths test it was going to owe ticket 06 is
  unbuildable now (there is one path); what replaces it is a cheap assertion that
  `/daily/older` and `HandleDaily` invoke the same partial. Its DTO-shaped sub-questions are
  void — re-read 07's answer preamble before trusting its bullets.

## Things you would otherwise learn the hard way

- **`claimed_by` and `undermined_by` have now had their first real outings** — the last
  handoff flagged both as never-exercised paths. Claims worked. `undermined_by` worked, and
  `wayfinder status` renders an Undermined section. The `## Ruled out` path is **still
  unexercised**; no ticket on this map has ever been ruled out of scope.
- **`claimed_at` must be RFC 3339.** Generate it (`date -u +%Y-%m-%dT%H:%M:%SZ`).
- **Tailwind is the Play CDN** (`templates/base.html:9`), not a build step. It observes DOM
  mutations, so server-inserted HTML is styled. There are no content globs, and no
  `tailwind.config.js` exists. Do not go looking for one.
- **Month breaks need no cross-window state.** Ticket 08 specifies a Go helper comparing
  adjacent groups; because the timeline is gap-filled and contiguous, a break is exactly
  `D.Day == 1`. Cheaper than 08 implies, and correct at a window seam.
- **`data-table.tsx` calls `hydrateRoot` on a container the server leaves empty** — a real
  hydration mismatch React silently recovers from. Found while resolving 07, and
  **deliberately not ticketed**: ticket 08 removes `data-table` from the daily view, so no
  daily-timeline work depends on it. It is separate work, and it is still a bug.
- **`AGENTS.md` calls the architecture "React Islands."** After 06 the scroll island is ~30
  lines of vanilla TypeScript. That sentence wants revisiting when the spec lands.
- **The `assets:` doc/reality mismatch is still open**, unchanged from the last handoff: the
  adapter says repo-relative, ticket 08's real usage is relative to `tickets/`, and nothing
  checks it.
- **Nothing in expensif points at the wayfinder tool.** Deferred across four sessions now.

## State of the tree

**Both repos are still on unmerged branches, not `main`.** Unchanged from the last handoff,
and it still matters: `main` has a ticket format that no longer exists.

- expensif: `wayfinder-tracker-adapter`
- `/Users/rengwu/Desktop/Projects/wayfinder`: `derive-status-from-body`

Both are fast-forward merges. Trees are clean. `wayfinder lint` reports the map clean. Go
tests were not run and did not need to be — no Go changed.

## Environment

- The tool: `cd ../wayfinder && go run ./cmd/wayfinder {status,lint} ../expensif/.plan/daily-timeline`.
  `lint` exits 1 on errors.
- `make dev` runs `air` (Go, :8080) and Vite (:8081). The Go server caches templates at
  startup; restart after template edits.
- The user's dev database (`~/.expensif/expenses.db`) is empty. Prototype against a scratch
  DB (`DATA_DIR=<tmp> PORT=8090 go run ./cmd/server`). **Do not seed the user's database.**

## Suggested skills

Apply only if available in your environment. All live in `.pocock-skills/<name>/SKILL.md` and
are **not** auto-discovered — read the file and follow it.

- `wayfinder` — the method, plus `TRACKER-MARKDOWN.md` for the shapes, the fence rule, claims
  and the verify checklist. Read both before touching `.plan/daily-timeline/`. **One ticket
  per session**; this session is the argument for it.
- `grill-me` — how both remaining tickets want to be resolved. Each enumerates its
  sub-decisions in its body.
- `codebase-design` — ticket 05 is a seam question about `HandleDaily`'s branches.
- `domain-modeling` — terms keep firming up ("the rail", "the foot", "the cursor", "day
  entry", "empty day", "undermined") and `CONTEXT.md` still does not exist.
- `handoff` — when you finish.
