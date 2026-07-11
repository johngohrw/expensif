# 2026-07-11 — The wayfinder tool grew a face

This session split cleanly in two. It opened by resolving one daily-timeline
ticket in **expensif**, then spent the rest building a **GUI for the wayfinder
tool** in the *other* repo. The bulk of the work — and the risk — is in
`/Users/rengwu/Desktop/Projects/wayfinder`, not here.

## Part one — expensif, briefly

Resolved [Should HandleDaily's two branches converge](../daily-timeline/tickets/05-converge-handledaily-branches.md):
**they don't.** `?date=` is a single-day detail view for any date in ±1 year (the
calendar links every cell), not a one-day window of a timeline that ends at today.
It keeps its own branch; what it shares is *rendering* (ticket 08's day-entry
partial) not control flow. Commit `7ffcc9e`, claim `4151de8`. The ticket and
commit carry the reasoning — don't re-derive.

**The daily-timeline map is now 9 of 10 resolved.** Only
[Test strategy](../daily-timeline/tickets/10-test-strategy.md) remains, unblocked
on the frontier. The map's destination — `.plan/daily-timeline/spec.md` — still
does not exist, and there is still **no application `.go`/`.tsx` code** for the
effort. If you return to expensif, ticket 10 then `to-spec` is the path.

## Part two — the wayfinder tool became an app

Everything below is in the **wayfinder repo** on branch `derive-status-from-body`
(a name that now badly undersells it — see *State of the tree*). Eleven commits,
`8b432aa`..`351d067`. **Read the commit messages; they carry the per-step
reasoning and I won't repeat it.**

What exists now: `wayfinder` gained `serve` (browser) and `app` (native WKWebView
window) commands that render a map as a **2.5D pannable star-map** — nodes are
glowing stars coloured by status, dependency edges are directed curved hyperlanes
with flow particles, the frontier pulses, undermined tickets wear a cracked red
halo, fog patches drift as nebulae at the rim. It is a live document: the client
polls a version token and folds edits in with animation, no reload. And it opens
to a **launcher** — splash → Open Folder (native macOS dialog) → pick a project →
a list of its `.plan` maps → click into one.

### Where the decisions live — do not re-litigate

- **`wayfinder/docs/starmap-design.md`** is the design record for the star-map:
  the seven grilled decisions (2.5D canvas not 3D, force-directed-but-rank-biased
  layout, deterministic + idle-drift motion, status-drives-the-star, flow
  hyperlanes, in-canvas select panel, HUD + fog-as-nebula) each chosen against
  alternatives, plus the v1/v2 split. **Read it before touching the canvas.**
- The **launcher** decisions are only in commit `351d067`'s message: native
  osascript folder dialog (chosen over an in-app browser and a dep), recents on
  the splash (persisted via `os.UserConfigDir`), SPA screens over one live canvas.

### Architecture, in one breath

`cmd/wayfinder/server.go` serves a small JSON API (`/api/initial`, `/api/pick`,
`/api/recents`, `/api/maps`, per-effort `/api/graph` + `/api/version`) over the
existing `internal/wayfinder` parser — untouched. `cmd/wayfinder/shell.go` is the
**entire client**: one big Go raw-string const holding HTML + CSS + hand-rolled
vanilla-JS canvas (layout, camera, animation, markdown renderer, live-reload,
three screens). `cmd/wayfinder/project.go` is maps-listing + recents + the
osascript pick. `internal/wayfinder/layout.go` holds `Layers()` (rank by
dependency depth); everything else visual is client-side. **Zero JS dependencies,
no build step.**

## The thing that is NOT verified — read this twice

**No human, and no tool, has seen any of it render.** This was a background
session with no display. I verified *logic and plumbing* exhaustively and *look
and feel* not at all:

- Verified headless: Go build/vet/test, JS syntax (`node --check` on the extracted
  script), layout determinism + finiteness + zero-overlap, live-reload version
  detection (on a **copy** of the effort — never mutate the real `.plan`), every
  API endpoint via curl (with a sandboxed `HOME` so recents didn't pollute
  `~/Library`), and the survivor-move-is-zero property of `relayoutWarm`.
- **Never verified:** the actual pixels, the native window opening, and
  `/api/pick` — calling it pops the macOS folder dialog and blocks, so it is the
  one endpoint I could not drive. Everything up to and from it is tested.

So the **first job of the next session (or the user) is to run
`cd ../wayfinder && go run ./cmd/wayfinder app` and look.** Expect to tune
constants I picked blind: `EASE`/`POSEASE`, bob amplitude, flare strength, the
layout `REP`/`SPRING`/`REST`/`RADIAL`, star colours/sizes, arrowhead size, the
label zoom cutoffs (`0.42`/`0.22`), fog nebula size/placement.

## Things you would otherwise learn the hard way

- **No backticks anywhere in `shell.go`'s JS — including comments.** The client
  lives in a Go raw string delimited by backticks; one stray backtick closes it
  and the JS becomes Go. I hit this exactly once (backticks in a comment →
  `undefined: code`). Guard: `grep -o '` + "`" + `' cmd/wayfinder/shell.go | wc -l`
  must print **2**. Markdown code fences are matched via `String.fromCharCode(96)`
  for this reason.
- **`webview_go` is the repo's only dependency and it is cgo.** `app` needs a cgo
  build (system WebKit on macOS, nothing to install); `serve` is pure stdlib and
  works in any browser. `app` is gated behind a `//go:build cgo` file with a
  no-cgo stub. Keep new native bits on the `app`/macOS side; keep `serve`
  portable.
- **macOS-only surfaces:** `/api/pick` (osascript folder dialog) and `app`
  (WKWebView). A cross-platform folder picker is unbuilt and deferred.
- **Tethers are dormant on daily-timeline.** Fog-to-ticket tethers are implemented
  but every current fog patch is unanchored (the one anchored patch was struck
  when ticket 05 resolved), so none render. To *see* a tether, add a
  `<clears-with: NN>` to a fog patch on a scratch map.
- **A live-reload test must run against a COPY.** Mutating a ticket to watch the
  map react will rewrite real `.plan` files. Copy the effort to a temp dir first.
- **The `app`/`serve` no-arg path is new.** `app <effort>` still jumps to that
  map; `app <project>` opens its list; `app` alone shows the splash. `status`/
  `lint` still default to `.`.

## State of the tree

**Both repos are on unmerged branches, both fast-forward to `main`, both clean.**

- expensif: `wayfinder-tracker-adapter` — ticket 05 on top.
- wayfinder: `derive-status-from-body` — **the branch name is now a lie.** It was
  about deriving ticket status from the body; it now also carries an entire GUI
  app. Consider renaming or just merging it before it grows further. The two
  pre-session commits (`9a0c5d6`, `4a209d7`) are the original status work.

Neither has been merged. `main` in wayfinder still has the CLI-only tool; `main`
in expensif still has the old daily-timeline ticket format. Go tests pass in
wayfinder; nothing in expensif's Go/TS changed this session.

**Still true, and now louder:** nothing in expensif points at the wayfinder tool —
no Makefile target, no doc. Deferred across five sessions. Now that it's a real
app you launch, a `make wayfinder` or a README line would actually earn its place.

## Environment

- Wayfinder tool: `cd ../wayfinder`. `go run ./cmd/wayfinder {app,serve,status,lint} [path]`.
  `serve` honours `PORT` (default 7777). `app` needs cgo (default on macOS).
- Headless verification recipe that worked: extract the client with
  `sed`/regex from the served `/` HTML, `node --check` it, and `eval` its pure
  functions (`layout`, `mdToHtml`, …) in node against a fetched `/api/graph`.
- Recents persist to `os.UserConfigDir()/wayfinder/recents.json`
  (`~/Library/Application Support/wayfinder/` on macOS). Sandbox with `HOME=<tmp>`
  when testing so you don't write real recents.

## Suggested skills

Apply only if available; all live in `.pocock-skills/<name>/SKILL.md` and are not
auto-discovered — read the file and follow it.

- `grill-me` — the user drove the whole star-map and launcher design this way, one
  decision at a time. The natural mode for the next design fork (visual tuning,
  the type-as-celestial-body v2 idea in the design doc, cross-platform picker).
- `wayfinder` — if you return to expensif's daily-timeline: ticket 10 is the last
  frontier ticket. **One ticket per session.**
- `to-spec` — once ticket 10 lands, the daily-timeline map's destination is a spec
  that still does not exist.
- `handoff` — when you finish.
