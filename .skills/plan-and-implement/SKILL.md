---
name: plan-and-implement
description: Collaborative feature planning and implementation skill. The user brings an idea; you analyze the codebase, surface options, grill the user (and yourself) on tradeoffs, and only write code once the plan is fully crystallized and explicitly approved. Use for any new feature, significant refactor, or architectural change.
---

# Plan & Implement

Turn a feature idea into a mature, agreed-upon implementation plan before any code is written. This skill is a structured conversation, not a code generator.

## Overview

1. The user describes an idea (feature, refactor, or architectural change).
2. You explore the codebase to understand constraints and existing patterns.
3. You present implementation options with honest tradeoffs.
4. You grill the user and self-grill until the plan is solid.
5. You crystallize the plan and **halt** for explicit user approval.
6. Only after approval do you implement, incrementally and transparently.

---

## Phase 1: Understand the Feature Request

Restate the user's idea in your own words. Identify:

- **The core goal** — what problem does this solve?
- **The user-facing behavior** — what will be different after this?
- **The implied scope** — what files, layers, and systems does this likely touch?
- **Implicit assumptions** — are they assuming a specific tech, pattern, or architecture?

If the request is vague, ask clarifying questions before moving on.

> **Push-back moment #1**: If the idea feels like an XY problem ("I want to use X to do Y"), suggest the simpler underlying approach. If it conflicts with an existing architecture decision in `.context/`, call it out immediately.

---

## Phase 2: Codebase Reconnaissance

Explore the codebase to ground your recommendations in reality. Do not skip this step even if you think you know the project.

1. Read `.context/` files (latest 4–6) for recent decisions and state.
2. List top-level directories and identify the tech stack.
3. Read entry points, core domain types, and interfaces.
4. Read any files you suspect this feature will touch.
5. Check for existing patterns that are relevant:
   - How are similar features already implemented?
   - What is the testing strategy?
   - What is the component / module reuse strategy?
   - How does data flow across layers?
6. Look for landmines:
   - Hardcoded assumptions this feature would violate
   - Missing abstractions you'd have to work around
   - Places where a change would cascade unexpectedly

Take notes. Cite specific file paths and line numbers.

---

## Phase 3: Present Options

Synthesize **2–4 implementation approaches**. At least one should be the "naive" path the user might already have in mind, and at least one should be an alternative that challenges their assumptions.

For each option, present:

| Field | What to include |
|-------|-----------------|
| **Name** | Short label, e.g., "Shared utility", "New service layer", "Extract component" |
| **Approach** | What would be built or changed, in plain English. Not code yet. |
| **Pros** | Why this is good: simplicity, performance, consistency, maintainability. |
| **Cons** | Honest tradeoffs: complexity, bundle size, coupling, dev overhead. |
| **Files touched** | Exact paths that would change or be created. |
| **Risk level** | Low / Medium / High — how likely is this to break existing behavior? |

After presenting options, ask the user which direction resonates, or if they want to combine ideas.

> **Push-back moment #2**: If an option would add dependencies, significant boilerplate, or violate DRY/KISS/YAGNI without clear payoff, say so directly. Provide the simpler alternative and explain why it wins.

---

## Phase 4: Grilling & Collaborative Design

Once the user leans toward an option, drop into deep design conversation.

### Grill the user

- "Does this need to work offline or without dynamic scripting?"
- "How much data / scale are we realistically targeting?"
- "Is this a one-off or will this pattern be reused elsewhere?"
- "What happens if the operation fails? What's the desired UX?"
- "Do you need this to be tested, and if so, at what level (unit, integration, e2e)?"
- "Is there a deadline or deploy constraint I should know about?"

### Self-grill

- "Am I forcing a client-side solution where a server-rendered one is sufficient?"
- "Am I creating an abstraction before there are two use cases?"
- "Would this change make the codebase harder for a new teammate to navigate?"
- "Am I preserving existing behavior, or am I sneaking in a breaking change?"
- "Is there existing code I can reuse or extract instead of writing new code?"

Iterate. It's normal to go back and forth 2–4 times. Revise the plan as new constraints surface.

> **Push-back moment #3**: If the user asks you to skip the plan and "just start coding," politely refuse. Explain that this skill exists to prevent rework, and that 5 minutes of design saves 30 minutes of undoing bad assumptions.

---

## Phase 5: Crystallize the Plan

When the conversation reaches a natural convergence, write a **concise, numbered implementation plan**:

```
## Approved Plan: <Feature Name>

1. **<Step 1>** — e.g., "Create shared `table` component with `variant` param."
2. **<Step 2>** — e.g., "Refactor `overview` and `users` screens to use the new component."
3. **<Step 3>** — e.g., "Extract shared row-rendering logic into reusable helper."
4. **<Step 4>** — e.g., "Run the project's build and test commands to verify."
```

Include:
- A summary of the agreed approach (2–3 sentences)
- The exact files to create, modify, or delete
- Any behavioral changes the user should expect
- Any risks or follow-up work you can already see

Then ask explicitly:

> **"This is the plan. Do not proceed until you confirm. Say 'go' or 'approved' when you're ready, or tell me what to adjust."**

**HALT. Do not write, edit, or create any files until the user confirms.**

---

## Phase 6: Implementation (Post-Approval Only)

After explicit user confirmation:

1. **Work incrementally** — one logical step at a time. Prefer small, reviewable edits.
2. **Build and test after each meaningful change** — run the project's standard build and test commands.
3. **Explain each step** — brief rationale before each edit block.
4. **Be transparent about doubts** — if you hit an unexpected snag, pause and discuss.
5. **Be ready to revert** — if the plan proves flawed mid-implementation, stop, explain why, and ask whether to adjust the plan or roll back.

### Step 6b: Add Tests (Optional but Recommended)

Before declaring the work done, assess whether new or updated tests are appropriate. Do not assume tests are always required — but do not skip them silently either.

1. **Check test coverage** — run the existing suite and see if your changes broke anything or left a gap.
2. **Ask the user** (or decide based on Phase 4 conversation):
   - "Should I add a test for the new behavior?"
   - "What level — unit, integration, or UI/render test?"
3. **Implement tests incrementally** just like production code, using the project's existing testing conventions (e.g., table-driven tests, widget tests, snapshot tests, golden files, etc.).
4. **Run the new tests** and ensure they pass along with the existing suite.

> If the user explicitly declined tests in Phase 4, skip this step but note in your summary: "Skipped tests per user preference."

### After implementation

1. Run the full test suite (including any new tests).
2. Run the project's linter / static analyzer / type checker.
3. Self-grill again:
   - Did I introduce coupling I didn't anticipate?
   - Are names clear and consistent?
   - Is there dead code or leftover scaffolding?
   - Did I break any existing behavior?
4. Suggest a git commit message (Conventional Commits).
5. Prompt for follow-up if the self-grill reveals additional needed changes.

---

## Rules

- **Never write code in Phases 1–5.** Only read, discuss, and plan.
- **Never implement without explicit confirmation.** The user must say "go", "approved", or equivalent.
- **Push back is a feature, not a bug.** If an idea is non-practical, over-engineered, or misaligned with the codebase, say so with clear reasoning.
- **Prefer removing or reusing code over adding new code.**
- **Preserve existing behavior** unless the user explicitly agrees to a breaking change.
- **Document reasoning** in the conversation so the user understands *why*, not just *what*.
- **If the user changes scope mid-plan**, pause and re-crystallize before continuing.
