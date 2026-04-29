---
name: get-up-to-speed
description: Reads the latest timestamped context files from .context/ to understand the current project state, tech stack, recent changes, and blockers before continuing development. Use at the start of any session when you need to catch up on where the project stands.
---

# Get Up To Speed

Load the most recent context files and summarize the project state so we can continue development without losing continuity.

## Steps

1. **List context files**
   ```bash
   ls -1 .context/ | sort
   ```

2. **Read the latest 8 files** (or fewer if less exist), prioritizing the most recent timestamps. Use the `read` tool on each.

3. **Synthesize a summary** covering:
   - **Project overview** — what is this, who uses it
   - **Tech stack** — languages, frameworks, database, key dependencies
   - **Directory structure** — where the main code lives
   - **Current feature set** — what's implemented and working
   - **Recent changes** — what was done in the last 1-3 sessions
   - **Known issues / tradeoffs** — accepted risks, TODOs, blockers
   - **Test state** — are tests passing, coverage notes
   - **Next logical steps** — what the context suggests should happen next

4. **Ask clarifying questions** if anything is ambiguous, contradictory, or missing critical detail (e.g., "The spec says X but the code shows Y — should I implement X?").

5. **Confirm readiness** — state clearly that you are caught up and ready to continue, referencing specific files or decisions from the context.

## Rules

- Always read from `.context/`; do not rely on memory from previous turns.
- If `.context/` does not exist, ask the user where session notes are stored.
- If a file is too large, read it in chunks using offset/limit.
- Cross-reference dates in filenames to establish chronological order.
- Note any uncommitted specs or architecture docs that are **not yet implemented** so they are not mistaken for active code.
