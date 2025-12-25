# SCHEMA.md

This document defines **how Archie documents work as a system**:
- what each workspace file is for,
- when you should edit it,
- which schema template to follow,
- and the constraints Archie enforces.

All schema templates live under:

```
.archie/docs/schema/
```

---

## 1. Schema Lookup Rules (Quick Index)

When writing or editing a workspace file, consult the corresponding schema file.

### 1.1 Workspace Root Documents

| Workspace File | Schema File |
|----------------|------------|
| `background.md` | `.archie/docs/schema/background.md` |
| `dependency.md` | `.archie/docs/schema/dependency.md` |
| `deployment.md` | `.archie/docs/schema/deployment.md` |
| `metrics.md` | `.archie/docs/schema/metrics.md` |
| `storage.md` | `.archie/docs/schema/storage.md` |
| `tasks.md` | `.archie/docs/schema/tasks.md` |

### 1.2 Feature / Spec

| Workspace File | Schema File |
|----------------|------------|
| `features/README.md` | `.archie/docs/schema/feature_readme.md` |
| `features/<feature-key>.md` | `.archie/docs/schema/feature.md` |
| `spec/<feature-key>.spec.md` | `.archie/docs/schema/spec.md` |

### 1.3 Workflow

| Workspace File | Schema File |
|----------------|------------|
| `workflow/<feature-key>/workflow.md` | `.archie/docs/schema/workflow.md` |

### 1.4 API

| Workspace File | Schema File |
|----------------|------------|
| `api/api.md` | `.archie/docs/schema/api.md` |

Notes:
- `api/<service>.thrift` is not Markdown-schema validated, but must follow:
  - hard rules in `background.md`
  - conventions and change records in `api/api.md`

---

## 2. File Catalog: Purpose, When to Edit, Constraints

### 2.1 `background.md`
**Purpose**
- Defines project-wide context, hard rules, tech stack, quality bar, and domain knowledge.

**When to edit**
- At project start (`init`)
- When global conventions or tech choices change
- When a decision impacts multiple features

**Follow schema**
- `.archie/docs/schema/background.md`

**Constraints**
- Must include “Hard Rules / Conventions” (non-negotiable)
- Open Questions / Assumptions should be append-only (recommended)
- Agents must raise a blocker if a feature design violates hard rules

---

### 2.2 `features/README.md` (Feature Registry)
**Purpose**
- A machine-readable registry of all features in the workspace.

**When to edit**
- Usually **never manually**; maintained by Archie (`init`, new feature creation, renames).

**Follow schema**
- `.archie/docs/schema/feature_readme.md`

**Constraints**
- Archie may regenerate this file entirely
- Must list each feature key and the corresponding file path

---

### 2.3 `features/<feature-key>.md`
**Purpose**
- The authoritative feature definition: status, scope, requirements, acceptance criteria, and links to artifacts.

**When to edit**
- `review`: clarify requirements, AC, dependencies
- `revise`: update requirements/scope; keep changelog
- `design/spec`: update artifact links + status/changelog only
- `update-progress`: update status/changelog only

**Follow schema**
- `.archie/docs/schema/feature.md`

**Constraints**
- Must contain required headings and required fields (Status + Spec fields)
- Status must match Archie state machine constraints
- Changelog must be append-only
- If any unchecked blocker exists for the feature → status must be `BLOCKED`

---

### 2.4 `workflow/<feature-key>/workflow.md` and `workflow/<feature-key>/*.mmd`
**Purpose**
- Defines the feature workflow in narrative + diagrams.
- This is mandatory for every feature.

**When to edit**
- `design`: always create/update workflow artifacts
- `revise`: update workflow when requirements change
- If diagrams are wrong or missing, fix here before generating spec

**Follow schema**
- `.archie/docs/schema/workflow.md`

**Constraints**
- Every feature must have:
  - `workflow/<feature-key>/workflow.md`
  - `workflow/<feature-key>/main.mmd`
- Mermaid diagrams must render without syntax errors
- Diagram naming is standardized: `main.mmd` required; `state.mmd` / `sequence.mmd` optional

---

### 2.5 `spec/<feature-key>.spec.md`
**Purpose**
- Coding-ready specification derived from feature + workflow + supporting artifacts.

**When to edit**
- `spec`: generate/refresh spec from design artifacts
- `revise`: update spec only when requirements/design changed
- Avoid manual edits unless your workflow expects it (team choice)

**Follow schema**
- `.archie/docs/schema/spec.md`

**Constraints**
- Must reference workflow files
- Must include required sections (interfaces/workflow/observability/rollout/test plan)
- Changelog must be append-only
- Feature `Spec.Readiness` can be `READY` only if spec passes schema validation

---

### 2.6 `tasks.md`
**Purpose**
- Execution plan (TODOs) + lightweight timeline per feature.

**When to edit**
- `design`: create initial tasks for implementation readiness
- `update-progress`: generate tasks if missing/incomplete; update task statuses; update timeline estimate
- Humans may update owners/ETAs (recommended)

**Follow schema**
- `.archie/docs/schema/tasks.md`

**Constraints**
- Must contain `By Feature` section and per-feature subsections
- Each task must have checkbox + log block (append-only)
- `update-progress` must ensure a `Timeline` subsection exists for the feature

---

### 2.7 `metrics.md`
**Purpose**
- Observability contract per feature (SLIs/SLOs/alerts/runbooks).

**When to edit**
- `design`: always produce feature metrics
- When reliability requirements change

**Follow schema**
- `.archie/docs/schema/metrics.md`

**Constraints**
- Every metric must include Description + Definition + Target
- Workflows with failure paths must have metrics covering success/error/latency (recommended minimum)

---

### 2.8 `storage.md`
**Purpose**
- Storage schema and change plan (create/alter/backfill) per feature.

**When to edit**
- `design`: whenever the feature requires persistence or data model changes
- When schema evolves

**Follow schema**
- `.archie/docs/schema/storage.md`

**Constraints**
- Must include latest schema and a change plan for the feature
- Change log should be append-only (recommended)
- Must obey data compliance rules defined in `background.md`

---

### 2.9 `api/api.md` (and `api/<service>.thrift`)
**Purpose**
- `api/api.md`: API index, service ownership, and change records per feature.
- `api/<service>.thrift`: IDL definition (when used).

**When to edit**
- `design`: whenever adding/modifying RPC/API
- When maintaining compatibility or deprecations

**Follow schema**
- `.archie/docs/schema/api.md` (for the index file)

**Constraints**
- API changes must be recorded in **append-only change records**
- For modifications: document old vs new behavior + migration plan
- Must obey API hard rules in `background.md` (e.g., Thrift-only)

---

### 2.10 `dependency.md`
**Purpose**
- A catalog of upstream/downstream dependencies shared across features (not per-feature deep detail).

**When to edit**
- When a new dependency is introduced
- When dependency properties change (owner/SLA/QPS/constraints)
- When you want a centralized view of system integration points

**Follow schema**
- `.archie/docs/schema/dependency.md`

**Constraints**
- Should list which features depend on each dependency (recommended)
- Avoid per-feature duplication: feature file should reference dependency name only

---

### 2.11 `deployment.md`
**Purpose**
- Release environment notes, deployment checklist, release order, and append-only release log.

**When to edit**
- Before and during rollout
- When release process changes
- After deployment to record outcomes

**Follow schema**
- `.archie/docs/schema/deployment.md`

**Constraints**
- Release log must be append-only
- Checklist should be used as the operational source of truth

---

## 3. Folder Structure Contracts

### 3.1 `features/`
Required:
```
features/
  README.md
  <feature-key>.md
```

### 3.2 `workflow/`
Required (minimum):
```
workflow/
  <feature-key>/
    workflow.md
```

Optional diagrams:
```
workflow/
  <feature-key>/
    state.mmd
    sequence.mmd
```

### 3.3 `assets/`
Required:
```
assets/
  images/
  exports/
```

Recommended:
```
assets/
  images/<feature-key>/
  exports/<YYYY-MM-DD-or-release>/
```

Constraints:
- Export output must go under `assets/exports/`
- Images should be grouped by feature for easy referencing

