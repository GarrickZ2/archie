# API / RPC Index

## Conventions
- IDL: Thrift only (or mixed allowed)
- Compatibility rules:
- Error model:

## Services

### <service-name>
- Owner:
- Purpose:
- IDL: api/<service-name>.thrift
- Consumers:

#### Change Records (append-only)
##### YYYY-MM-DD: <feature-key> - <change title>
- Type: ADD | MODIFY | DEPRECATE | REMOVE
- Scope: <methods/structs/enums>
- Compatibility: backward-compatible? (yes/no) + why
- Migration plan:
- Links: features/<feature-key>.md, spec/<feature-key>.spec.md, MR/PR

#### Methods (summary)
- <Method>: purpose + link to IDL section anchor (or line ref)