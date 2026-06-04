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
