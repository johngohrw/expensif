# Bag of Skills

A collection of agent-agnostic skills following the [Agent Skills standard](https://agentskills.io). These are stateless capabilities that read from and write to a project-local `.context/` directory to maintain continuity across sessions.

This repository is not a standalone project. It is meant to be cloned into a `.skills/` directory within other repositories.

---

## Skills

| Skill                     | Description                                                                                                                                                      |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `save-session`            | Capture the current session — files changed, decisions made, test state, open items — to a timestamped markdown file in `.context/`.                             |
| `resume-session`          | Read the most recent `.context/` files at session start to reconstruct project state, recent changes, and next steps without relying on model context.           |
| `plan-then-build`         | Explore the codebase, present implementation options with tradeoffs, grill the user on constraints, and crystallize an explicit plan before any code is written. |
| `generate-commit-message` | Inspect staged changes, recent commit history, and `.context/` files to draft a Conventional Commit message for user approval.                                   |
| `arch-review`             | Walk the codebase to identify coupling, untested modules, shallow abstractions, and mixed concerns. Presents ranked candidates for refactoring.                  |
| `spec-gen`                | Deeply explore a codebase and emit a comprehensive project spec covering stack, schemas, routes, features, design decisions, and technical debt.                 |

---

## Install

Drop `install.sh` into the target project root and execute it:

```bash
curl -O https://raw.githubusercontent.com/johngohrw/bag-of-skills/main/install.sh
bash install.sh
```

The script performs the following:

1. Aborts if `AGENTS.md` already exists in the project root (to prevent overwriting existing agent instructions).
2. Clones this repository into `.skills/` via a shallow clone.
3. Copies `AGENTS.md` from `.skills/` into the project root.
4. Deletes `install.sh`.

> **Note:** `.skills/` is not added to `.gitignore`. Commit it alongside `AGENTS.md` so your team shares the same skill set.

---

## Update

Each project maintains its own `.skills/` clone. To update skills in a project:

```bash
.skills/update.sh
```

This is a thin wrapper around `git pull` executed from within `.skills/`.

---

## The `.context/` Convention

Most skills in this collection interact with a `.context/` directory at the project root. This directory serves as out-of-band persistent memory: session summaries, architecture decisions, project specs, and other continuity artifacts are written here by agents and read back on subsequent invocations.

Because `.context/` lives inside the project, it can be committed to git or left untracked, depending on whether you want session history to travel with the repository.

---

## `AGENTS.md`

This repository includes an `AGENTS.md` containing introductory context about the skills collection. When `install.sh` runs, it copies this file into the project root. Agents that support context files (e.g. pi, Claude Code) automatically load `AGENTS.md` at startup when walking up from the working directory. This means a freshly started agent knows about the available skills and conventions without explicit prompting.
