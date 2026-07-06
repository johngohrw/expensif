# Agent Skills

This project includes a `.skills/` directory containing reusable agent capabilities, following the [Agent Skills standard](https://agentskills.io). These are available for use within this project.

## Available Skills

| Skill | Description |
|-------|-------------|
| `save-session` | Persist session state to `.context/` |
| `resume-session` | Read `.context/` to resume work |
| `plan-then-build` | Collaborative planning before coding |
| `generate-commit-message` | Conventional commit drafting |
| `arch-review` | Deep codebase analysis |
| `spec-gen` | Generate comprehensive project specs |

## `.context/` Convention

Skills read from and write to a `.context/` directory at the project root. This directory serves as persistent memory for session summaries, architecture decisions, and project specs across agent sessions.

If `.context/` does not exist, create it when a skill requires it.

## Updating Skills

To update the skills in this project:

```bash
cd .skills && git pull
```

Or run `.skills/update.sh` if available.

<!-- iudex:begin -->
## iudex — shaping work before the queue

This project uses [iudex](https://github.com/) to drive tickets through a
`queue → activate → impl → QA → human-qa → merge` pipeline. iudex owns everything
from `iudex queue` onward; the skills below own the step *before* that — turning a
raw idea into robust, sliced, dependency-ordered tickets.

```
grill-me / grill-with-docs → prototype → to-prd → to-issues → iudex queue
        (stress-test)        (validate)   (spec)    (slice)     (registered tickets)
```

When the user wants to shape, refine, or break down work, **read and follow the
relevant skill file** under `.iudex/skills/<name>/SKILL.md`:

| Skill | Use it to |
|-------|-----------|
| `grill-me` | Stress-test a raw idea with relentless one-at-a-time questions (no docs needed) |
| `grill-with-docs` | Same, but challenge the plan against `.context/` glossary + ADRs and update them inline |
| `prototype` | Build throwaway code (terminal logic app or UI variants) to validate a design |
| `to-prd` | Synthesize the discussion into a PRD at `.context/prd/<slug>.md` |
| `to-issues` | Slice a plan/PRD into iudex tickets and register them with `iudex queue --deps` |
| `improve-codebase-architecture` | Find deepening opportunities (HTML report) that feed back into the funnel |

Notes:

- These skills are **run from the workspace root** (where `.iudex/` lives), before
  tickets exist. They are not present inside ticket worktrees.
- `.context/` (domain glossary, ADRs, PRDs) is **tracked** project documentation —
  commit it. `.iudex/` is gitignored operational state.
- iudex's event log is the single source of truth: ticket dependencies live only in
  the `iudex queue --deps` command, never in the `t<id>.md` body.
<!-- iudex:end -->
