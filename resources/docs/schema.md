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

### 1.5 Test Plan

| Workspace File | Schema File |
|----------------|------------|
| `testplan/<feature-key>.md` | `.archie/docs/schema/testplan.md` |

---

## 2. File Catalog

| File | Purpose | When to Edit | Constraints |
|------|---------|--------------|-------------|
| `background.md` | Defines project-wide context, hard rules, tech stack, quality bar, and domain knowledge. | At project start (`init`), when global conventions or tech choices change, when a decision impacts multiple features. | Must include "Hard Rules / Conventions" (non-negotiable), Open Questions / Assumptions should be append-only, Agents must raise a blocker if a feature design violates hard rules. |
| `features/README.md` (Feature Registry) | A machine-readable registry of all features in the workspace. | Usually **never manually**; maintained by Archie (`init`, new feature creation, renames). | Archie may regenerate this file entirely, Must list each feature key and the corresponding file path. |
| `features/<feature-key>.md` | The authoritative feature definition: status, scope, requirements, acceptance criteria, and links to artifacts. | `review`: clarify requirements, AC, dependencies + update changelog; `revise`: update requirements/scope + update changelog; `design`: update artifact links + status (NO changelog); `spec`: update spec link/readiness + status + changelog; `plan`: update status + changelog; `test-plan`: update testplan link + changelog. | Must contain required headings and required fields (Status + Spec fields), Status must match Archie state machine constraints, Changelog must be append-only, If any unchecked blocker exists for the feature → status must be `BLOCKED`. |
| `workflow/<feature-key>/workflow.md` and `workflow/<feature-key>/*.mmd` | Defines the feature workflow in narrative + diagrams. This is mandatory for every feature. | `design`: always create/update workflow artifacts; `revise`: update workflow when requirements change; If diagrams are wrong or missing, fix here before generating spec. | Every feature must have: `workflow/<feature-key>/workflow.md` and `workflow/<feature-key>/main.mmd`, Mermaid diagrams must render without syntax errors, Diagram naming is standardized: `main.mmd` required; `state.mmd` / `sequence.mmd` optional. |
| `spec/<feature-key>.spec.md` | Coding-ready specification derived from feature + workflow + supporting artifacts. | `spec`: generate/refresh spec from design artifacts; `revise`: update spec only when requirements/design changed; Avoid manual edits unless your workflow expects it (team choice). | Must reference workflow files, Must include required sections (interfaces/workflow/observability/rollout/test plan), Changelog must be append-only, Feature `Spec.Readiness` can be `READY` only if spec passes schema validation. |
| `tasks.md` | Execution plan (TODOs) + lightweight timeline per feature. | `plan`: generate tasks from design artifacts (with user approval); update task statuses; update timeline estimate; Humans may update owners/ETAs (recommended). | Must contain `By Feature` section and per-feature subsections, Each task must have checkbox + log block (append-only), `plan` must ensure a `Timeline` subsection exists for the feature. |
| `metrics.md` | Observability contract per feature (SLIs/SLOs/alerts/runbooks). | `design`: always produce feature metrics; When reliability requirements change. | Every metric must include Description + Definition + Target, Workflows with failure paths must have metrics covering success/error/latency (recommended minimum). |
| `storage.md` | Storage schema and change plan (create/alter/backfill) per feature. | `design`: whenever the feature requires persistence or data model changes; When schema evolves. | Must include latest schema and a change plan for the feature, Change log should be append-only (recommended), Must obey data compliance rules defined in `background.md`. |
| `api/api.md` (and `api/<service>.thrift`) | `api/api.md`: API index, service ownership, and change records per feature. `api/<service>.thrift`: IDL definition (when used). | `design`: whenever adding/modifying RPC/API; When maintaining compatibility or deprecations. | API changes must be recorded in **append-only change records**, For modifications: document old vs new behavior + migration plan, Must obey API hard rules in `background.md` (e.g., Thrift-only). |
| `dependency.md` | A catalog of upstream/downstream dependencies shared across features (not per-feature deep detail). | When a new dependency is introduced; When dependency properties change (owner/SLA/QPS/constraints); When you want a centralized view of system integration points. | Should list which features depend on each dependency (recommended), Avoid per-feature duplication: feature file should reference dependency name only. |
| `deployment.md` | Release environment notes, deployment checklist, release order, and append-only release log. | Before and during rollout; When release process changes; After deployment to record outcomes. | Release log must be append-only, Checklist should be used as the operational source of truth. |
| `testplan/<feature-key>.md` | Simple test case checklist including test IDs, descriptions, priorities, and failure scenario mapping per feature. | `design`: generate initial test plan from artifacts; `test-plan`: update test cases; `revise`: when workflow/spec changes affect testing. | Must include Test Cases section (Unit/Integration/E2E) with simple format, Every test case must link to workflow/spec sections, Optional Failure Scenarios section for error paths. |

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

### 3.4 `testplan/`
Required (minimum):
```
testplan/
  <feature-key>.md
```

Constraints:
- Every feature should have testplan after design phase
- Test case IDs must follow convention: <feature-key>-{U|I|E}-<###>
- Simple checklist format with minimal fields

---

## 4. Command Workflow Patterns

All Archie commands follow interactive workflows with explicit user confirmation before file operations.

### 4.1 Standard Workflow Pattern

1. **Gather Information** - Collect context through dialogue or read existing files
2. **Process/Generate** - Analyze, plan, or generate content based on inputs
3. **Present to User** - Show proposed changes/generation results
4. **Get Approval** - Ask explicit confirmation ("Should I proceed?")
5. **Execute** - Write files ONLY after approval
6. **Report** - Summarize what was done

### 4.2 Command-Specific Behaviors

#### Interactive Commands (Continuous Dialogue)
- `init`, `review`, `revise`, `design`
- Multi-step dialogue to gather requirements
- Each major step may have approval gates
- Final approval required before any writes

#### Generation Commands (Show Results)
- `plan`, `spec`, `test-plan`
- Generate content from existing artifacts
- Present generated content to user
- Allow modifications before approval
- Write only after explicit approval

#### Read-Only Commands
- `ask`
- No file modifications
- Answer questions based on existing docs
- Cite sources

#### Maintenance Commands
- `fix`
- Detect formatting issues
- Show before/after comparison
- Apply fixes only after approval
- Never change semantic content

#### Intelligent Dispatcher Commands
- `context`
- Accept free-form information
- Analyze and suggest target files/sections
- Format content to match schema
- Show routing plan and get approval
- Write to multiple files after approval

### 4.3 Changelog Policy

Not all commands update changelog in features/<feature-key>.md:

- **Update changelog**: `init`, `review`, `revise`, `spec`, `plan`, `test-plan`, `context` (when adding to features)
- **NO changelog**: `design` (design focuses on creating artifacts, not tracking changes)
- **Rationale**: Changelog tracks requirement changes and milestone events, not design artifact creation
- **Note**: `context` updates changelog only when adding to features/<key>.md, not for background.md

### 4.4 Features File Update Matrix

| Command | Requirements | Scope | Status | Design Artifacts | Spec | Changelog |
|---------|-------------|-------|--------|-----------------|------|-----------|
| init | ✓ | ✓ | ✓ (NOT_REVIEWED) | - | - | ✓ |
| review | ✓ | ✓ | ✓ (→READY_FOR_DESIGN) | - | - | ✓ |
| revise | ✓ | ✓ | May rollback | - | - | ✓ |
| design | - | - | ✓ (→DESIGNED) | ✓ (links) | - | ❌ |
| plan | - | - | ✓ (→IMPLEMENTING/FINISHED) | - | - | ✓ |
| spec | - | - | ✓ (→SPEC_READY) | - | ✓ (link+readiness) | ✓ |
| test-plan | - | - | - | ✓ (testplan link) | - | ✓ |
| context | May add | May add | - | May add | - | ✓ (if feature) |
| fix | - | - | - | - | - | - |

---

## 5. Approval and Safety Rules

1. **User Approval Required**: All commands that write files must get explicit user approval before writes
2. **Show Before Write**: Always present what will be written (summary, diff, or complete content)
3. **No Auto-Execution**: Even "auto-triggered" behaviors (like generating tasks if missing) require user approval
4. **Preserve Content**: When normalizing/restructuring, preserve original content in Legacy section or assets/
5. **Schema Compliance**: All writes must follow schema templates defined in `.archie/docs/schema/`
6. **Blocker Creation**: When issues cannot be auto-resolved, create blockers instead of guessing

---

---
