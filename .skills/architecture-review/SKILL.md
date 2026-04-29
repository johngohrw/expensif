---
name: architecture-review
description: Deeply analyze a codebase for architectural improvements, refactoring opportunities, and testability issues. Walks through exploration, candidate presentation, and a collaborative implementation loop with the user. Use when the codebase feels messy, untested, hard to navigate, or when preparing for feature work that touches core structures.
---

# Architecture Review & Improvement

Analyze the current codebase to find deepening opportunities, refactoring candidates, and ways to make it more testable and AI-navigable. Then collaborate with the user to implement the best ones.

## Prerequisites

If a `.context/` folder exists, read the latest timestamped files first to understand recent decisions, state, and known issues.

## Phase 1: Exploration

Explore the codebase organically. Do not follow a rigid checklist — walk through files as a developer would when trying to understand and modify the system. Note friction points such as:

- **Conceptual clarity** — is it hard to understand what a module does?
- **Overly complex implementation** — could this be simpler?
- **Untested modules** — critical code with no tests?
- **Hard-to-test code** — global state, tight coupling, side effects?
- **Bugs or mistakes** — clear errors or dangerous patterns?
- **Shallow modules with complex interfaces** — thin wrappers that leak complexity?
- **Leaky abstractions** — does a lower layer force upper layers to know too much?
- **Duplicated logic** — same pattern copied in multiple places?
- **Mixed concerns** — business logic bleeding into HTTP handlers, SQL in services, etc.?
- **Naming mismatches** — do names describe what things actually do?

### How to explore

1. List the top-level directory structure.
2. Read `go.mod`, `package.json`, or equivalent to understand dependencies.
3. Walk through the main entry point(s).
4. Read core domain models / types.
5. Read interfaces — they reveal the intended architecture.
6. Read implementations — they reveal the actual architecture.
7. Read tests — they reveal what is valued and what is ignored.
8. Read templates / UI code if relevant.
9. **Inspect git commit history** — `git log --oneline --all`, `git log -p -- <file>`, or `git blame <file>`. Commit messages, authors, and chronology reveal why certain decisions were made, which changes were rushed, and how the architecture evolved over time. A recent large refactor may explain an odd seam; a series of small patches may indicate organic growth that outpaced design.
10. Note any TODO comments, FIXME, or `panic` usage.

Keep a running list of observations. Be specific about file paths and line numbers where possible.

## Phase 2: Presenting Candidates

From your observations, synthesize a **numbered list of opportunities**. Sort by balancing:

- **Ease of implementation** (low-hanging fruit)
- **Importance** (severity of issue, value of improvement)

For each candidate, present:

| Field                 | What to include                                                    |
| --------------------- | ------------------------------------------------------------------ |
| **#**                 | Number for easy reference                                          |
| **Problem**           | Why this is an issue. Be specific.                                 |
| **Files / Modules**   | Exact paths involved.                                              |
| **Proposed Solution** | In plain English — what would change and why. Not code yet.        |
| **Benefits**          | What improves: testability, clarity, maintainability, performance? |
| **Risk / Effort**     | Quick note on how invasive this is.                                |

**Ask the user**: "Which of these would you like to explore further?" You may also ask follow-up questions to clarify priorities (e.g., "Are you more concerned with test coverage or with simplifying the handler layer?").

## Phase 3: Grilling & Implementation Loop

Once the user picks a candidate, drop into a collaborative design conversation.

### Before writing code

1. **Discuss constraints** — What depends on this? What would break? What tests need updating?
2. **Present options** — If there are multiple ways to solve it, lay them out with tradeoffs.
3. **Listen for feedback** — The user may have constraints you don't know (deployment process, team preferences, backward compatibility needs).
4. **Push back if needed** — If the user suggests something non-practical (e.g., "let's rewrite everything in Rust" or "let's add 5 new dependencies"), provide gentle pushback with clear reasoning. Explain why a simpler or more incremental approach is better.
5. **Crystallize the plan** — Agree on the exact approach. Summarize it in 2-3 sentences before touching any files.
6. **Halt until user confirmation** - Stop here. Only proceed further when user gives a positive confirmation.

### During implementation

1. **Implement incrementally** — One logical change at a time. Prefer small, reviewable edits over massive rewrites.
2. **Run tests after each meaningful change** — `go test ./...`, `npm test`, etc.
3. **Build the project** — Ensure it compiles: `go build ./...`, `npm run build`, etc.
4. **Be transparent about doubts** — If you're unsure about a change, say so. Ask the user.
5. **Be ready to revert** — If things go sideways, revert and discuss a different approach. Do not push through a broken state.

### After implementation

1. **Run the full test suite** — All tests must pass.
2. **Run lint / vet** — `go vet ./...`, `eslint`, etc.
3. **Self-grill** — Ask yourself:
   - Did I introduce any new coupling?
   - Are the names clear?
   - Could a new team member understand this?
   - Did I leave any dead code?
   - Are there edge cases I missed?
   - Would this break any existing behavior?
4. **Prompt for follow-up** — If the self-grill reveals additional needed changes, ask the user for permission to proceed. Do not sneak in extra changes without asking.
5. **Suggest a git commit message** — Keep it clear and concise. Examples:
   - `refactor: extract validation from handlers into service layer`
   - `test: add unit tests for expense repository`
   - `fix: make DeleteUser atomic via SQL transaction`

## Rules

- **Never start coding in Phase 2.** Only present candidates.
- **Never implement without user confirmation** in Phase 3. Get explicit agreement on the approach.
- **Preserve all existing behavior** unless the user explicitly agrees to a breaking change.
- **Prefer removing code over adding code.** Simplicity is the goal.
- **Document your reasoning** in the conversation so the user understands _why_ a change was made.
