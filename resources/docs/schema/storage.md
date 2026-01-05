# Storage

## Conventions
| Item | Value |
|------|-------|
| Primary DB | <MySQL/PostgreSQL/etc> |
| Other stores | <Redis/ES/etc> |
| Naming | <rules> |
| Migration | <strategy> |
| Retention/TTL | <rules> |
| PII/Compliance | <rules> |

## Data: <table_name>

- Type: MySQL/Redis/etc
- Purpose: <one sentence>

### Schema
```sql
CREATE TABLE ...
```

### Change Plan
- Action: CREATE | ALTER | DROP | BACKFILL
- Steps: <numbered list>
- Rollback: <plan>
- Risk: <locking notes>

### Change Log
<!-- ARCHIE:APPEND_ONLY -->
- YYYY-MM-DD (who): <change> (MR/PR link)
<!-- ARCHIE:END -->
