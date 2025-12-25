# Archie

**Archie** is a docs-first AI CLI that helps engineers turn vague ideas into clear, executable technical designs — and keep projects moving with structure and discipline.

Archie is not a coding tool.  
It is a **design, specification, and execution-orchestration system** built on Markdown files.

---

## Why Archie?

Most projects fail before coding starts:
- requirements are unclear,
- designs are scattered,
- decisions are lost in chat,
- progress is hard to track.

Archie fixes this by enforcing:

- 📄 **Docs as Source of Truth**
- 🔄 **A strict feature state machine**
- 🧠 **AI-assisted review, design, and spec generation**
- 📦 **Everything stored locally, in Markdown**
- 🧩 **Composable sub-agents for API, workflow, storage, metrics, tasks**

---

## Core Idea

Archie treats a **feature** as the atomic unit and enforces the minimum design chain:

> **Feature → Workflow → Spec**

Everything else (API, storage, metrics, tasks, deployment notes) supports this chain.

---

## What Archie Does

With Archie you can:

- Initialize a project from scratch or messy notes
- Review features and clarify requirements
- Design workflows, APIs, storage, and observability
- Generate coding-ready specs
- Track blockers, tasks, and progress
- Export clean documentation bundles

All via a CLI.

---

## Example Workflow

```bash
archie init
archie review checkout-discount
archie design checkout-discount
archie spec checkout-discount
archie update-progress checkout-discount --note "Implementation started"
```

---

## Project Structure (Simplified)

```
.
├── background.md
├── features/
│   └── checkout-discount.md
├── workflow/
│   └── checkout-discount/
├── spec/
│   └── checkout-discount.spec.md
├── tasks.md
├── metrics.md
└── api/
```

---

## Key Concepts

### Commands
High-level orchestrators:
- `init`
- `review`
- `design`
- `spec`
- `update-progress`

### Subagents
Specialized capability units:
- Workflow Designer
- API Designer
- Storage Designer
- Metrics Designer
- Task Manager

### Templates
Archie uses a lightweight template system so teams can customize document structure while keeping it machine-parseable.

---

## Who Is This For?

- Backend / Full-stack engineers
- Tech leads
- Early-stage startups
- Infra & platform teams
- Anyone tired of chaotic design docs

---

## Philosophy

- Local-first
- Markdown-native
- Deterministic over magical
- Explicit over implicit
- Humans stay in control

---

## Status

Archie is under active design and early development.

Contributions, feedback, and ideas are welcome.
