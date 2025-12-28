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

---

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

---

## Core Dependency Chain (Minimum Deliverable)

For every feature, Archie enforces the minimum design chain:

**Feature → Workflow → Spec**

Other artifacts (API / Storage / Metrics / Tasks / Deployment notes) are supportive, but the chain above is required.

---

## Feature State Model

### Feature Status Enum

- `NOT_REVIEWED`
- `UNDER_REVIEW`
- `BLOCKED`
- `READY_FOR_DESIGN`
- `UNDER_DESIGN`
- `DESIGNED`
- `SPEC_READY`
- `IMPLEMENTING`
- `FINISHED`

### Status Authority

- Source of truth: `features/<feature-key>.md`
- If there exists any unchecked blocker in `blocker.md`
  → status **must** be `BLOCKED`
- If `spec/<feature-key>.spec.md` exists and
  `Spec.Readiness = READY`
  → status may advance to `SPEC_READY`

State changes must be explicit and recorded in `features/<feature-key>.md` changelog.

---

## File System Contract

All agents and commands must strictly respect file boundaries.
Please Refer to .archie/docs/schema.md for details.


## Hard Constraints (Non-Negotiable)

- Design style rules are defined in `background.md`
- Review must precede design
- Design must precede spec
- No command may succeed without file output
- No silent state changes are allowed

---

## Success Criteria

A feature is considered successfully handled by Archie if:
- Its state transitions are valid and traceable
- A workflow exists and is linked from the feature
- A complete Coding Spec is produced
- Key decisions can be audited via files

---

## Command Common Rules

All Archie commands must follow these common rules:

### Core Principles Reference
All commands must follow the Core Principles defined above.

### State Machine Reference
All commands must respect the Feature State Model defined above.

### File Contract Reference
All commands must strictly respect file boundaries (see .archie/docs/schema.md).

### Hard Constraints Checklist
Commands must enforce:
- [ ] Review must precede design
- [ ] Design must precede spec
- [ ] No command may succeed without file output
- [ ] No silent state changes are allowed
- [ ] If any unchecked blocker exists → status must be BLOCKED

### Format Compliance
All file operations must strictly follow the format rules defined in .archie/docs/schema.md
