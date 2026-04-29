---
name: generate-commit-message
description: Reads staged git changes, recent commit history, and project context to generate a well-structured conventional commit message. Present it to the user for approval before committing.
---

# Generate Commit Message

Analyze the current staged changes, recent commit style, and project context to craft a commit message that accurately describes what was done and why.

## Steps

1. **Gather context**
   ```bash
   git status --short
   git diff --staged --stat
   git diff --staged --name-only
   git log --oneline -10
   ```

2. **Read key changed files** (optional but helpful)
   - For new files: read them to understand their purpose
   - For significant modifications: read the diff or the file

3. **Load project context**
   ```bash
   ls -1 .context/ | sort
   ```
   Read the 2–4 most recent `.context/` files to understand recent work and decisions.

4. **Analyze and synthesize**
   - Group changes into logical themes (e.g., "dev workflow", "new feature", "refactor", "bugfix")
   - Identify the *primary* change that defines the commit
   - Note secondary changes that support it
   - Check if there are any "hidden" fixes (e.g., proxy bug, HMR preamble) that aren't obvious from filenames

5. **Draft the message**
   - Use **Conventional Commits** format: `type(scope?): subject`
   - `type`: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `build`
   - Keep subject line under 72 characters
   - Add a blank line, then bullet points for significant details
   - Each bullet should answer *what* and *why*, not just *what*
   - Reference context files or decisions when relevant (e.g., "per golang-react-islands-spec")

6. **Present to user**
   Show the full message in a code block. Ask:
   > "Does this look right? Say 'commit' to stage it, or tell me what to change."

## Rules

- Do NOT commit automatically — always wait for explicit user confirmation.
- If there are unrelated changes mixed together, warn the user and suggest splitting.
- Match the style of recent commits (check `git log --oneline`).
- If `.context/` exists, use it — it contains the "why" behind the changes.
- For delete operations, mention what was removed and why.
- For new files, mention their purpose, not just their names.
- Keep bullet points parallel in structure (start with verb, same tense).

## Example Output

```
feat: add CategoryPills island, Button component, and dev workflow fixes

- Replace inline vanilla JS category fetch with React CategoryPills island
  hydrated on expense add/edit forms
- Add reusable Button.tsx component with 6 variants (primary, secondary,
  neutral, ghost, danger, pill)
- Add templates/partials/button.html Go partial kept in sync with Button.tsx
- Convert all template buttons to use the shared partial
- Makefile: `make dev` now starts both Go and Vite servers in one command
- Fix Vite dev proxy to bypass Vite internals via `bypass()` function
- Add React Refresh preamble to DevClient for HMR support
- Fix AssetHelper.ScriptTag to match manifest entries by Name or Src path
```
