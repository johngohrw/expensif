# Bag of Skills

This repo is a collection of reusable coding agent skills. It is not a runnable
project — it is a shared library of capabilities.

## Convention
Each skill lives in its own directory with a `SKILL.md` file, following the
[Agent Skills standard](https://agentskills.io).

## Skills included
- `save-session` — persist session state to `.context/`
- `resume-session` — read `.context/` to resume work
- `plan-then-build` — collaborative planning before coding
- `generate-commit-message` — conventional commit drafting
- `arch-review` — deep codebase analysis
- `spec-gen` — generate comprehensive project specs

## Usage in projects
Projects clone this repo into `.skills/`. Skills read and write to the
project's own `.context/` directory.

## Updating
Projects can update their skills with:
```bash
cd .skills && git pull
```
Or run `.skills/update.sh` if available.
