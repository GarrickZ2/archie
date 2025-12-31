# Storage

## Conventions
- Primary DB: [ ]
- Other stores: [ ]
- Naming rules: [ ]
- Migration strategy: [ ]
- Retention/TTL rules: [ ]
- PII / compliance rules: [ ]

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

#### Change Plan
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
