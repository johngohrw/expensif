---
name: project-spec-generator
description: Explore a codebase deeply and generate a comprehensive project spec document describing stack, structure, architecture, schemas, routes, features, and design decisions. Use when the user wants to understand a project, onboard to it, or document its current state.
---

# Project Spec Generator

Deeply explore a codebase and produce a structured project specification document saved to `.context/project-spec-<YYYY-MM-DD>.md`.

## Prerequisites

If a `.context/` folder exists, read any timestamped or named context files first to understand recent decisions, state, and known issues.

## Phase 1: Exploration

Explore the codebase organically. Walk through files as a developer would when trying to understand and modify the system. Follow this loose order:

1. **Top-level structure** — List the root directory. Identify monorepo tools, workspace files, and package manager.
2. **Manifest files** — Read `package.json`, `go.mod`, `Cargo.toml`, `pyproject.toml`, or equivalent to understand dependencies and scripts.
3. **Workspace config** — Read `pnpm-workspace.yaml`, `nx.json`, `turbo.json`, etc.
4. **Entry points** — Find and read main entry files (e.g., `src/index.ts`, `src/main.tsx`, `app/main.py`).
5. **Core configuration** — Read framework configs (`astro.config.mjs`, `next.config.mjs`, `vite.config.ts`, `payload.config.ts`, etc.).
6. **Schemas / Models** — Read database schemas, ORM models, Payload collections, Prisma files, or domain types.
7. **Routes / API** — Map HTTP routes, API endpoints, page files, and server handlers.
8. **Key components / UI** — Read significant UI files, layouts, and pages to understand features.
9. **Interfaces / Types** — Read type definitions and interfaces; they reveal intended architecture.
10. **Tests** — Search for and read test files to understand what is valued and what is ignored.
11. **Git history** — Run `git log --oneline --all -30` (or more) to understand evolution, rushed decisions, and refactor history.
12. **Environment / Infra** — Check `.env.example`, Dockerfiles, CI configs, deployment scripts.
13. **TODOs / FIXMEs** — Search for `TODO`, `FIXME`, `HACK`, `XXX` comments.

Keep a running list of observations. Be specific about file paths and patterns.

## Phase 2: Synthesis

From your observations, synthesize a **comprehensive project spec** in Markdown with the following sections:

### Required Sections

| Section | What to Include |
|---------|-----------------|
| **1. Overview** | What the project does, its purpose, and who it serves. |
| **2. Stack** | Table of technologies with versions (frameworks, languages, databases, infra). |
| **3. Directory Structure** | Tree or bulleted list of key directories and their roles. |
| **4. Schemas & Models** | All collections, tables, entities, types. Field names, types, relationships, constraints. |
| **5. Routes** | All API routes, page routes, admin routes. Note static vs dynamic, SSR vs SSG. |
| **6. Features** | Major functional capabilities. Split by domain if needed. |
| **7. Key Design Decisions** | Why certain architectural choices were made (e.g., monorepo, headless CMS, specific DB). |
| **8. Architecture Structure** | ASCII or text diagram showing how packages/services/data flow interact. |
| **9. Issues & Technical Debt** | Numbered list of problems found: no tests, dead code, security risks, coupling, missing packages, etc. |
| **10. Usage** | How to install, run, build, test, and deploy. Environment variables required. |
| **11. Tests** | What testing exists (or doesn't). Recommended testing strategy if absent. |
| **12. Notable Files** | Table of the most important files and their purposes. |

### Rules for the Document

- Be specific with file paths and line numbers where relevant.
- If something is disabled, missing, or commented out, explicitly state that.
- Do not speculate beyond what the code shows; it is okay to say "appears unused" or "likely boilerplate."
- Use tables and code blocks for readability.
- Include the date in the filename: `.context/project-spec-<YYYY-MM-DD>.md`.

## Phase 3: Save

Write the final Markdown document to:

```
.context/project-spec-<YYYY-MM-DD>.md
```

If the `.context/` directory does not exist, create it.

Then confirm to the user:
- The file path where the spec was saved.
- A 2-3 sentence summary of what the project is.
- The top 3 issues or most interesting architectural findings.
