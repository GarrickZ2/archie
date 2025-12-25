# Storage

## Conventions
- Primary DB: (e.g. MySQL)
- Other stores allowed: (Redis/Mongo/Dynamo) (if any)
- Naming rules:
- Migration rules:
- Retention/TTL rules:
- PII / compliance rules:

## MySQL

### <feature-key>

#### Final Schema (Latest)
##### Table: <table_name>
- Purpose:
- Primary key:
- Indexes:
- Access patterns:

Columns:
- <col_name> <type> [NULL/NOT NULL] [DEFAULT ...]
    - Comment:
    - Notes: (range/encoding)

#### Change Plan (What to do)
- Action: CREATE | ALTER | DROP | BACKFILL
- Steps:
    1. ...
    2. ...
- Rollback:
- Risk/Locking notes:

#### Change Log (Append-only)
- YYYY-MM-DD (who): <what changed> (links: MR/PR)

## Redis (optional)

### <feature-key>
- Key patterns:
- Value schema:
- TTL / eviction:
- Access patterns:
- Failure/consistency notes:
- Change Log: