# Archie

Archie is a docs-first AI CLI that helps engineers turn vague ideas into clear, executable technical designs.

Archie is NOT a coding tool. It is a design, specification, and execution-orchestration system built on Markdown files.

## Project Architecture

### Core Components

```
cmd/           → CLI entry points (init, review, design, spec, etc.)
internal/      → Core logic
  status/      → Feature status parsing and display
  export/      → Documentation export and bundling
resources/     → AI context files
  agents/      → AGENTS.md - role definition and constraints
  docs/        → schema.md + schema/*.md - output format specs
  commands/    → archie-*.yaml - command workflows
```

### Key Files

| File | Purpose |
|------|---------|
| `AGENTS.md` | AI role, state machine, output discipline |
| `schema.md` | Document structure rules |
| `schema/*.md` | Templates for each document type |
| `archie-*.yaml` | Command workflows with guard checks |

## Design Principles

### Feature State Machine

```
NOT_REVIEWED → [review] → UNDER_REVIEW → READY_FOR_DESIGN
READY_FOR_DESIGN → [design] → UNDER_DESIGN → DESIGNED
DESIGNED → [spec] → SPEC_READY
SPEC_READY → [plan] → IMPLEMENTING → FINISHED
Any → BLOCKED (if unchecked blocker exists)
```

### Minimum Design Chain

**Feature → Workflow → Spec**

Other artifacts (API, storage, metrics) support this chain but are not strictly required.

## Code Requirements

### Go Code Style

- Follow standard Go conventions (gofmt, golint)
- Error handling: wrap errors with context using `fmt.Errorf("context: %w", err)`
- Logging: use structured logging
- Naming: clear, descriptive names; avoid abbreviations

### Documentation Files

When generating or modifying Archie documentation:

1. **Conciseness**: Record only key decisions, not thought processes
2. **Field Limits**:
   - One-liner: max 80 chars
   - Summary items: 2-3 sentences
   - Requirements: 1 sentence each
   - Changelog: 1 line per entry
3. **Forbidden Content**:
   - Thought processes ("after consideration...")
   - Redundant modifiers ("very important", "please note")
   - Duplicated information across files
   - Obvious explanations

### Command YAML Structure

Each command file follows:

```yaml
title: <CommandName>
description: <one-line description>
content: |
  <purpose - one line>

  ## Guard Check
  - Status MUST be: <allowed statuses>
  - Required: <files that must exist>

  ## Reads
  - <files to read>

  ## Writes
  - <files to write>

  ## Workflow
  <numbered steps>

  ## Core Rules
  <key constraints>
```

### Schema File Structure

- Use `<placeholder>` format, not `[ ]`
- Use tables for structured data
- No `---` separators between sections
- Minimal prose, focus on structure

## Testing

Run tests: `go test ./...`
Build: `go build -o archie .`

## Common Tasks

### Adding a New Command

1. Create `resources/commands/archie-<name>.yaml`
2. Add guard check, reads/writes, workflow steps
3. Register in command registry (cmd/)

### Modifying Schema

1. Update `resources/docs/schema/<file>.md`
2. Update `resources/docs/schema.md` lookup table if needed
3. Ensure commands reference the correct schema

### State Machine Changes

1. Update AGENTS.md Feature State Model
2. Update guard checks in affected commands
3. Verify state transitions are valid

## Output Quality Checklist

Before writing any documentation:

- Can half the content be removed without losing key info?
- Is there content duplicated from other files?
- Are there explanations instead of decisions?
- Does it follow the schema template exactly?
