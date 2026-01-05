# SCHEMA.md

Defines how Archie documents work as a system.

All schema templates: `.archie/docs/schema/`

## 1. Schema Lookup Rules

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
| `glossary.md` | `.archie/docs/schema/glossary.md` |
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
- `api/<ServiceName>.thrift` and `api/<ServiceName>.proto` are not Markdown-schema validated, but must follow:
  - hard rules in `background.md`
  - conventions and change records in `api/api.md`
  - one service per file (all methods in one .thrift/.proto file)

### 1.5 Test Plan

| Workspace File | Schema File |
|----------------|------------|
| `testplan/<feature-key>.md` | `.archie/docs/schema/testplan.md` |

## 2. Folder Structure Contracts

### 2.1 `features/`
Required:
```
features/
  <feature-key>.md
```

### 2.2 `workflow/`
Required:
```
workflow/
  <feature-key>/
    workflow.md
```

Optional diagrams (naming should reflect function):
```
workflow/
  <feature-key>/
    state.mmd        # state machine diagram
    sequence.mmd     # sequence diagram
    flowchart.mmd    # flowchart diagram
```

Constraints:
- Every feature MUST have workflow.md
- Diagram files are optional, naming should reflect their purpose
- Mermaid diagrams must be valid syntax

### 2.3 `api/`
Required:
```
api/
  api.md
  <ServiceName>.thrift
  <ServiceName>.proto
```

Constraints:
- **One service per file** - All methods for a service in one .thrift/.proto file
- **api.md is lightweight** - Index + change records only (no decisions/analysis)
- **Complete interface definitions** in .thrift/.proto files
- Decisions and analysis belong in features/ or spec/ files

### 2.4 `assets/`
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

### 2.5 `testplan/`
Required (minimum):
```
testplan/
  <feature-key>.md
```

Constraints:
- Every feature should have testplan after design phase
- Test case IDs must follow convention: <feature-key>-{U|I|E}-<###>
- Simple checklist format with minimal fields
