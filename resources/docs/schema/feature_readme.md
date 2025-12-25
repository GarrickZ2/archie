# Feature Registry

| Feature | Status | Owner | Depends On | Spec |
|--------|--------|-------|------------|------|
| feature-a | READY_FOR_DESIGN | @alice | feature-b | ❌ |
| feature-b | BLOCKED | @bob | infra-x | ❌ |
| feature-c | SPEC_READY | @me | - | ✅ |

Rules:
- Source of truth for feature content is `features/<feature>.md`
- This file is auto-updated by Archie