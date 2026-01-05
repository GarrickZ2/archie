# Spec: <feature-key>

## Goal
<What this spec enables - 1 sentence>

## Non-Goals
<Explicitly not covered>

## Background
<Short context - 2-3 sentences>

## Interfaces

### API/RPC
- IDL: api/<Service>.thrift
- Methods:
  - <MethodName>: <brief description>
    - Request: <key fields>
    - Response: <key fields>
    - Errors: <error codes>

### Events/Async
- Topic: <topic name>
- Payload: <key fields>
- Ordering/retry: <strategy>

## Data Model

### Entities
- <Entity>: <fields>, <constraints>, <lifecycle>

### Storage
- Tables: <list>
- Indexes: <list>
- Migration: <plan>

## Workflow

### Main Flow
<Step-by-step flow>

### Edge Cases
| Case | Handling |
|------|----------|
| <case> | <handling> |

### Failure & Retry
- Idempotency: <strategy>
- Retry: <policy>
- Timeout: <value>

## Observability

### Metrics
| Metric | Description |
|--------|-------------|
| <name> | <description> |

### Logs
<Key log points>

### Traces
<Trace boundaries>

## Security
- AuthN/AuthZ: <approach>
- Data sensitivity: <classification>
- Audit: <requirements>

## Rollout
- Feature flag: <name>
- Canary: <strategy>
- Rollback: <plan>

## Test Plan
- Unit: <key tests>
- Integration: <key flows>
- E2E: <key journeys>
- Detail: testplan/<feature-key>.md

## Changelog
- YYYY-MM-DD (who): <one-line description>
