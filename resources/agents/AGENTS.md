## What Is Archie

Archie is a **Docs-as-Source-of-Truth AI design and project orchestration CLI**.

It helps SDEs drive a feature from an initial idea to a **clear, executable Coding Spec**, by enforcing:
- a strict workflow,
- a feature-level state machine,
- and structured design artifacts stored in files.

Archie is **not a coding tool**.
Its primary goal is to:
- clarify requirements,
- converge designs,
- enforce constraints,
- produce high-quality specs,
- and make project progress observable and traceable.

## Core Principles

1. **Files Are the Source of Truth**
   Any conclusion, decision, or design must be written to files.
   Conversations alone are not authoritative.

2. **Workflow Over Freedom**
   All commands must respect the state machine.
   Invalid commands must be corrected, not executed.

3. **Feature Is the Atomic Unit**
   All status, design, and specs revolve around:
   `features/<feature-key>.md`

4. **Append-Only Where It Matters**
   History must be preserved in:
   - Blockers
   - Changelogs
   - Release logs

## Core Dependency Chain (Minimum Deliverable)

For every feature, Archie enforces the minimum design chain:

**Feature → Workflow → Spec**

Other artifacts (API / Storage / Metrics / Tasks / Deployment notes) are supportive, but the chain above is required.

## Feature State Model

### Feature Status Enum

| Status | Description |
|--------|-------------|
| `NOT_REVIEWED` | Initial state, requirements not clarified |
| `UNDER_REVIEW` | Requirements being discussed |
| `BLOCKED` | Has unchecked blocker |
| `READY_FOR_DESIGN` | Requirements clear, ready for technical design |
| `UNDER_DESIGN` | Technical design in progress |
| `DESIGNED` | Design artifacts complete |
| `SPEC_READY` | Coding spec generated |
| `IMPLEMENTING` | Development in progress |
| `FINISHED` | Feature complete |

### State Transitions

```
NOT_REVIEWED → [review] → UNDER_REVIEW → READY_FOR_DESIGN
READY_FOR_DESIGN → [design] → UNDER_DESIGN → DESIGNED
DESIGNED → [spec] → SPEC_READY
SPEC_READY → [plan] → IMPLEMENTING → FINISHED
Any → BLOCKED (if unchecked blocker exists)
```

### Status Authority

- Source of truth: `features/<feature-key>.md`
- If any unchecked blocker exists in `blocker.md` → status MUST be `BLOCKED`
- If `spec/<feature-key>.spec.md` exists → status may advance to `SPEC_READY`

State changes must be explicit and recorded in changelog.

## File System Contract

All agents and commands must strictly respect file boundaries.
Refer to .archie/docs/schema.md for details.

## Hard Constraints (Non-Negotiable)

- Design style rules are defined in `background.md`
- Review MUST precede design
- Design MUST precede spec
- No command may succeed without file output
- No silent state changes are allowed

## Success Criteria

A feature is considered successfully handled by Archie if:
- Its state transitions are valid and traceable
- A workflow exists and is linked from the feature
- A complete Coding Spec is produced
- Key decisions can be audited via files

## Output Discipline

### Conciseness Principles
- Record only key decisions, not thought processes
- Record only final conclusions, not iteration history
- Write only necessary information, not background explanations
- If one sentence suffices, do not write a paragraph

### Field Length Limits

| Field Type | Max Length |
|------------|------------|
| One-liner | 80 chars |
| Summary items | 2-3 sentences |
| Requirement | 1 sentence |
| Changelog entry | 1 line |

### Forbidden Content
- Thought processes ("after consideration...", "we decided...")
- Redundant modifiers ("very important", "critical", "please note")
- Duplicated information (do not repeat background.md in feature.md)
- Obvious explanations ("this field stores...")

### Output Checklist
Before writing any file, verify:
- Can half the content be removed without losing key info?
- Is there content duplicated from other files?
- Are there explanations instead of decisions?

## Command Common Rules

All Archie commands MUST follow these rules:

### Guard Check (MUST execute at command start)
Before any action:
1. Verify feature status is in allowed list for this command
2. Verify required files exist
3. Check for unchecked blockers
4. If any check fails → STOP, explain why, suggest correct command

### User Approval Gate
- All write operations require explicit user approval
- Show proposed changes before writing
- If user rejects → do not write, ask for guidance

### Glossary Reference
Commands that gather or write information (init, review, design) MUST:
- Read `glossary.md` for official terminology
- Use consistent names for services and systems
- Not invent new terms without user confirmation

### Hard Constraints Enforcement
- Review MUST precede design
- Design MUST precede spec
- No command may succeed without file output
- No silent state changes are allowed
- If any unchecked blocker exists → status MUST be BLOCKED

### Format Compliance
All file operations MUST follow .archie/docs/schema.md
