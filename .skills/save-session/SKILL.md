---
name: save-session
description: Save the current development session to a timestamped markdown file in .context/. Captures what was done, decisions made, files changed, test state, and open items for continuity. Use at the end of any session before switching tasks or ending the conversation.
---

# Save Session

Persist the current session's work, decisions, and state to a timestamped markdown file in `.context/` so future sessions can pick up where this one left off.

## Steps

1. **Determine the filename**
   ```bash
   date +%Y-%m-%d
   ```
   Use `.context/session-YYYY-MM-DD.md`. If a file with that name already exists, append a suffix like `-pm` or `-2`.

2. **Gather session content** by reviewing:
   - What the user asked for at the start of the session
   - What skills or approaches were used
   - Files created, modified, or deleted
   - Architectural decisions made and their rationale
   - Test results (pass/fail counts, coverage notes)
   - Bugs found or fixed
   - Candidates identified but not yet implemented
   - Open questions or next steps

3. **Write the file** with this structure:

   ```markdown
   # Project Name — Session Summary (YYYY-MM-DD)

   ## Agenda
   Brief description of what this session set out to do.

   ## Changes Made

   ### Change 1 — Short title
   **Problem:** What was the issue or goal?
   **Files:** Which files were touched?
   **Change:** What specifically changed? (plain English, not full diffs)
   **Tests:** Impact on test suite (count before/after, pass/fail)

   ### Change 2 — ...

   ## Current Test State
   - `go test ./...` / `npm test` — results
   - Lint / vet results
   - Build results

   ## Decisions & Rationale
   Any architectural or design decisions made during the session and why.

   ## Open Items / Next Session Notes
   - Unfinished candidates or TODOs
   - Recommended next steps
   - Blockers or questions
   ```

4. **Be concise but specific.** Mention exact file paths, function names, and test names. Do not paste full code diffs — summarize the change in 2-4 sentences.

5. **Cross-reference .context/ files** if relevant (e.g., "See also `golang-react-islands-spec-2026-04-28.md`").

## Rules

- Always write to `.context/`, never overwrite existing files without checking.
- Include test counts before and after if tests were added or removed.
- If no changes were made (e.g., pure exploration), still document findings and candidates.
- Note any skills created during the session and their locations.
- If the session ended with the user picking a candidate for next time, record which one and why.
