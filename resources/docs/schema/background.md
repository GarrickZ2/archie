# Project Context

## Purpose

- What problem this project solves
- Who it is for
- Why it exists (business / product / technical motivation)

> This section defines *why* the system exists.
> All feature designs must align with this purpose.

---

## Success Criteria

- What does “success” look like?
- Measurable outcomes (performance, reliability, adoption, etc.)

---

## Tech Stack

This section defines **allowed and preferred technologies**.

- Language / Runtime:
- Frameworks:
- API / RPC:
- Storage:
- Messaging / Async:
- Observability:
- CI / CD:
- Deployment Platform:

---

## Project Conventions (Hard Rules)

> These are **non-negotiable constraints**.
> If a feature design violates any rule here, the agent must raise a Blocker.

### Code Style

- Naming conventions:
- Error handling style:
- Logging conventions:
- Formatting rules:

### Architecture Patterns

- Overall architecture style (e.g. layered, hexagonal, DDD-lite)
- Service boundaries
- Synchronous vs asynchronous preferences
- Idempotency rules

### API / Contract Rules

- IDL format (e.g. Thrift only)
- Backward compatibility requirements
- Versioning strategy
- Error model

### Testing Strategy

- Required test types (unit / integration / e2e)
- Coverage expectations
- Mocking vs real dependencies

### Git / Delivery Workflow

- Branching strategy:
- Commit message conventions:
- PR / MR requirements:
- Release cadence:

---

## Domain Context

Domain-specific knowledge that **cannot be inferred from code alone**.

Include:

- Core business concepts
- Key entities and relationships
- Invariants and rules
- Common edge cases
- Known pitfalls

> Agents should rely on this section to understand *what is normal* and *what is dangerous* in the domain.

---

## Important Constraints

### Technical Constraints

- Performance limits
- Scale assumptions
- Legacy system limitations

### Business Constraints

- Deadlines
- Cost limits
- Risk tolerance

### Compliance / Security Constraints

- PII handling
- Regulatory requirements
- Audit requirements

---

## External Dependencies

### Upstream Systems

- Name:
- Owner / Team:
- What it provides:
- Stability / SLA:
- Known issues:

### Downstream Systems

- Name:
- How they consume this system:
- Compatibility requirements:

---

## Assumptions

Explicit assumptions made during design.

<!-- ARCHIE:APPEND_ONLY -->

- YYYY-MM-DD: <assumption>

<!-- ARCHIE:END -->

---

## Open Questions

Unresolved global questions that affect multiple features.

<!-- ARCHIE:APPEND_ONLY -->

- YYYY-MM-DD: <question>

<!-- ARCHIE:END -->
