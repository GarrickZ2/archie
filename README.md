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
- 🧩 **Composable sub-agents for API, workflow, storage, metrics, tasks, test plans**

---

## Core Idea

Archie treats a **feature** as the atomic unit and enforces the minimum design chain:

> **Feature → Workflow → Spec**

Everything else (API, storage, metrics, tasks, test plans, deployment notes) supports this chain.

---

## Installation

```bash
go install github.com/GarrickZ2/archie@latest
```

### Verify Installation

```bash
archie --help
```

---

## How to Use Archie

Archie has **two modes of operation**:

### Mode 1: CLI Commands (Terminal)

Direct workspace management commands you run in your terminal:

| Command | Description |
|---------|-------------|
| `archie init` | Initialize workspace structure and install agent commands |
| `archie setup` | Interactive TUI to edit background and manage features |
| `archie status` | Show project status with interactive feature browser |
| `archie export` | Export documentation to single markdown file |

### Mode 2: Agent Commands (Coding Assistant)

AI-powered design commands that work **inside your coding agent** through conversational slash commands.

**Supported Coding Agents:**
- 🤖 **Claude Code** (`.claude/commands/`)
- 🤖 **Cursor** (`.cursor/commands/`)
- 🤖 **Windsurf** (`.windsurf/workflows/`)
- 🤖 **Gemini Code Assist** (`.gemini/commands/`)
- 🤖 **Qwen Code** (`.qwen/commands/`)
- 🤖 Custom agents (via `archie custom-agent`)

**How it works:**
1. Run `archie init` to install agent command files
2. Open project in your coding assistant (e.g., Claude Code, Cursor)
3. Use slash commands + conversation to invoke Archie agents

**Available Agent Commands:**

| Slash Command | Description | Status Required |
|---------------|-------------|-----------------|
| `/init` | Initialize project workspace | N/A |
| `/review` | Clarify requirements and dependencies | NOT_REVIEWED, UNDER_REVIEW, BLOCKED |
| `/design` | Generate design artifacts (workflow, API, storage, metrics) | READY_FOR_DESIGN, UNDER_DESIGN |
| `/test-plan` | Create comprehensive test plan | DESIGNED, SPEC_READY, IMPLEMENTING, FINISHED |
| `/spec` | Generate coding-ready specification | DESIGNED |
| `/plan` | Manage execution tasks and timeline | DESIGNED, SPEC_READY, IMPLEMENTING, FINISHED |
| `/revise` | Modify/create/merge/split/delete features | Any |
| `/ask` | Ask questions about the project | Any |
| `/fix` | Fix validation errors or schema issues | Any |

---

## Quick Start

### Step 1: Install Archie (CLI)

```bash
go install github.com/GarrickZ2/archie@latest
```

### Step 2: Initialize Your Project (CLI)

```bash
mkdir my-project
cd my-project
archie init
```

This creates:
- Workspace file structure (`features/`, `workflow/`, `spec/`, etc.)
- Schema templates (`.archie/docs/schema/`)
- Agent command files for your coding assistant

### Step 3: Set Up Project Context (CLI)

```bash
archie setup
```

Interactive TUI to:
- Edit `background.md` (project context, tech stack, hard rules)
- Create and manage feature files

### Step 4: Use Agents in Your Coding Assistant

Open the project in **Claude Code**, **Cursor**, **Windsurf**, or your preferred coding agent.

**Example conversation in Claude Code:**

```
You: /review user-authentication

Claude: I'll help you review the user-authentication feature.
        Let me clarify the requirements and dependencies...

        [Conducts AI-assisted review, discovers dependencies,
         identifies blockers, updates feature documentation]

        ✅ Review complete. Status → READY_FOR_DESIGN

---

You: /design user-authentication

Claude: I'll generate the design artifacts.

        [Creates workflow diagrams, API specs, storage schemas,
         metrics definitions, task breakdown]

        ✅ Design complete. Generated:
           - workflow/user-authentication/main.mmd
           - api/api.md (updated)
           - storage.md (updated)
           - metrics.md (updated)

---

You: /test-plan user-authentication

Claude: I'll create a comprehensive test plan.

        [Generates test case matrix (Unit/Integration/E2E),
         failure scenarios, coverage mindmap, automation strategy]

        ✅ Test plan complete.
           - testplan/user-authentication.md

---

You: /spec user-authentication

Claude: I'll generate the coding-ready specification.

        [Aggregates all design artifacts into executable spec]

        ✅ Spec ready. Status → SPEC_READY
           - spec/user-authentication.spec.md
```

### Step 5: Monitor Progress (CLI)

```bash
archie status
```

Interactive TUI showing:
- Overall project health
- Feature status distribution
- Blocked features
- Detailed feature information

### Step 6: Export Documentation (CLI)

```bash
archie export
```

Generates a single markdown file with:
- Selected documentation
- Feature specifications
- TOC, statistics, and dependency graph

---

## What Archie Does

With Archie you can:

- ✅ Initialize a project from scratch or messy notes
- ✅ Review features and clarify requirements
- ✅ Design workflows, APIs, storage, and observability
- ✅ Generate comprehensive test plans
- ✅ Generate coding-ready specs
- ✅ Track blockers, tasks, and progress
- ✅ Export clean documentation bundles
- ✅ Manage feature lifecycle (create/merge/split/delete)

All via a CLI with interactive TUI support.

---

## Project Structure

```
.
├── background.md           # Project context, tech stack, hard rules
├── features/
│   ├── README.md          # Feature registry
│   └── <feature-key>.md   # Feature definition
├── workflow/
│   └── <feature-key>/
│       ├── workflow.md    # Workflow narrative
│       └── main.mmd       # Main flow diagram (Mermaid)
├── spec/
│   └── <feature-key>.spec.md  # Coding-ready specification
├── testplan/
│   └── <feature-key>.md       # Test case checklist
├── tasks.md               # Execution tasks per feature
├── metrics.md            # Observability per feature
├── storage.md            # Database schemas per feature
├── api/
│   └── api.md            # API index and change records
├── dependency.md         # Dependency catalog
└── deployment.md         # Release notes and checklist
```

---

## Complete Workflow Example

```bash
# 1. Install and initialize (CLI)
go install github.com/GarrickZ2/archie@latest
mkdir my-ecommerce
cd my-ecommerce
archie init

# 2. Set up project (CLI)
archie setup
# → Edit background.md with tech stack, hard rules
# → Create feature: "checkout-discount"

# 3. Open in Claude Code and use agents (Agent Mode)
```

**In Claude Code:**
```
You: /review checkout-discount

Claude: Let me review the checkout-discount feature...
        [Reviews requirements, discovers dependencies, identifies blockers]
        ✅ Status → READY_FOR_DESIGN

You: /design checkout-discount

Claude: Generating design artifacts...
        ✅ Created workflow diagrams, API specs, storage schemas

You: /test-plan checkout-discount

Claude: Creating comprehensive test plan...
        ✅ Generated 15 unit tests, 8 integration tests, 3 E2E tests

You: /spec checkout-discount

Claude: Generating coding specification...
        ✅ Spec ready for implementation

You: /plan checkout-discount

Claude: Creating execution plan...
        ✅ 12 tasks created, estimated timeline: 5-7 days
```

```bash
# 4. Monitor progress (CLI)
archie status
# → Interactive TUI showing all features and their status

# 5. Export for review (CLI)
archie export
# → Generated: archie-export-2024-01-15.md
```

---

## Feature State Machine

```
NOT_REVIEWED → UNDER_REVIEW → READY_FOR_DESIGN → UNDER_DESIGN →
DESIGNED → SPEC_READY → IMPLEMENTING → FINISHED

Special state: BLOCKED (can occur at any stage)
```

---

## Key Concepts

### Commands
High-level orchestrators that manage the feature lifecycle.

### Subagents
Specialized capability units:
- Workflow Designer
- API Designer
- Storage Designer
- Metrics Designer
- Task Manager
- Test Plan Generator

### Templates
Archie uses schema templates (`.archie/docs/schema/`) so teams can maintain consistent document structure while keeping it machine-parseable.

### State Machine
Every feature follows a strict state progression ensuring design quality before implementation.

---

## Advanced Features

### Feature Management with `/revise` (Agent Mode)

The `/revise` agent command supports powerful feature lifecycle operations:

**In your coding assistant:**
```
You: /revise --create payment-gateway

Agent: Creating new feature: payment-gateway
       ✅ Created features/payment-gateway.md
       ✅ Status: NOT_REVIEWED

You: /revise --merge user-login,user-signup --into user-auth

Agent: Merging features...
       ✅ Merged features/user-login.md + features/user-signup.md
       ✅ Created features/user-auth.md
       ✅ Archived source features

You: /revise --split checkout-flow

Agent: How would you like to split this feature?
       [Interactive conversation to determine split boundaries]
       ✅ Created 3 new features from checkout-flow

You: /revise --change checkout-discount --status IMPLEMENTING

Agent: Updating status...
       ✅ checkout-discount: DESIGNED → IMPLEMENTING
```

### CLI Commands Reference

#### Initialize Workspace
```bash
archie init
```
Creates workspace structure and installs agent commands for all supported coding assistants.

#### Interactive Setup
```bash
archie setup
```
TUI interface to edit project context and manage features.

#### Status Monitoring
```bash
# Interactive status browser
archie status

# Compact status report
archie status --compact
```

#### Documentation Export
```bash
# Interactive export with selection
archie export

# Export to specific file
archie export -o docs/design.md

# Export without table of contents
archie export --no-toc

# Export without statistics
archie export --no-stats

# Export without dependency graph
archie export --no-dep-graph
```

---

## Who Is This For?

- Backend / Full-stack engineers
- Tech leads
- Early-stage startups
- Infra & platform teams
- Anyone tired of chaotic design docs

---

## Philosophy

- **Local-first**: Everything stored in your filesystem
- **Markdown-native**: Human-readable, version-controllable
- **Deterministic over magical**: Predictable, explicit behavior
- **Explicit over implicit**: No hidden state, clear contracts
- **Humans stay in control**: AI assists, you decide

---

## Status

Archie is under active design and early development.

Contributions, feedback, and ideas are welcome.

---

## License

MIT License - see LICENSE file for details.

---

## Contributing

We welcome contributions! Please:
1. Fork the repository
2. Create a feature branch
3. Submit a pull request with clear description

For bugs or feature requests, please open an issue.
