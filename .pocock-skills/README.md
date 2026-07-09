# Skills

General-purpose agent skills, adapted from [Matt Pocock's skills](https://github.com/mattpocock/skills) (MIT — see [LICENSE](./LICENSE)). Rewritten to be:

- **Standalone** — no setup skill, no per-repo configuration step. Each skill states its own conventions inline.
- **Language-agnostic** — examples are pseudocode or multi-language; tooling references defer to the project's own task runner and static checks.
- **Harness/model-agnostic** — no dependence on a specific agent product. Where a skill benefits from subagents, it degrades gracefully: run isolated parallel subagents if the environment supports them, otherwise run the same work as sequential, self-contained passes.

## Shared conventions

- **`.plan/`** — the project's planning memory, committed to version control. Specs (`.plan/<slug>/spec.md`), ticket breakdowns (`.plan/<slug>/tickets.md`), wayfinder maps (`.plan/<slug>/map.md` + `.plan/<slug>/tickets/NN-<slug>.md`), and handoffs (`.plan/handoffs/`).
- **`CONTEXT.md`** at the repo root — the domain glossary; **`docs/adr/`** — architecture decision records. Both created lazily by `domain-modeling`; skills that read them proceed silently when they don't exist.

## Skills

**User-invoked** — reached by typing their name; their job is to orchestrate a session.

- **[grill-me](./grill-me/SKILL.md)** — a relentless interview to sharpen a plan or design, one question at a time.
- **[grill-with-docs](./grill-with-docs/SKILL.md)** — grill-me plus live domain-model upkeep: updates `CONTEXT.md` and ADRs as decisions crystallise.
- **[handoff](./handoff/SKILL.md)** — compact the current conversation into a handoff document (`.plan/handoffs/`) so a fresh agent can continue.
- **[improve-codebase-architecture](./improve-codebase-architecture/SKILL.md)** — scan for deepening opportunities, present as a visual HTML report, then grill through the chosen one.
- **[to-spec](./to-spec/SKILL.md)** — synthesize the current conversation into a spec at `.plan/<slug>/spec.md`.
- **[to-tickets](./to-tickets/SKILL.md)** — break a plan or spec into tracer-bullet tickets at `.plan/<slug>/tickets.md`.
- **[wayfinder](./wayfinder/SKILL.md)** — chart a big, foggy effort as a map of investigation tickets under `.plan/<slug>/`, resolved one per session.
- **[writing-great-skills](./writing-great-skills/SKILL.md)** — reference for writing and editing skills well.

**Model-invoked** — reachable by the agent on its own (and by other skills), or by typing their name.

- **[codebase-design](./codebase-design/SKILL.md)** — shared vocabulary and principles for designing deep modules.
- **[domain-modeling](./domain-modeling/SKILL.md)** — actively build and sharpen the project's domain model (`CONTEXT.md`, `docs/adr/`).
- **[prototype](./prototype/SKILL.md)** — throwaway code that answers a design question: an interactive terminal app for logic/state questions, or radically different UI variants behind one switcher.
- **[research](./research/SKILL.md)** — investigate a question against primary sources and capture cited findings as markdown.
- **[review-code](./review-code/SKILL.md)** — two-axis review of a diff: Standards (repo conventions + smell baseline) and Spec (does it implement what was asked?).
