# API / RPC Index

## Conventions
- **RPC**: Thrift IDL - one service per .thrift file under `api/`
- **HTTP API**: Protobuf - one service per .proto file under `api/`
- **Compatibility rules**: [define your rules]
- **Error model**: [define your error handling approach]

---

## Services

### <ServiceName>
- **Owner**: <team or person>
- **Purpose**: <one-line description>
- **IDL**: `api/<ServiceName>.thrift` or `api/<ServiceName>.proto`
- **Consumers**: <list of consumers>

#### Change Records (append-only)

##### YYYY-MM-DD: <feature-key> - <change title>
- **Type**: ADD | MODIFY | DEPRECATE | REMOVE
- **Scope**: <methods/structs/enums changed>
- **Compatibility**: yes/no + brief reason
- **Migration**: <if needed, brief plan or link>
- **Links**:
  - Feature: `features/<feature-key>.md`
  - Spec: `spec/<feature-key>.spec.md`

#### Methods (summary)
- `MethodName`: <brief purpose>

---

## Example

### UserService
- **Owner**: Platform Team
- **Purpose**: User account and profile management
- **IDL**: `api/UserService.thrift`
- **Consumers**: WebApp, MobileApp, AdminPanel

#### Change Records (append-only)

##### 2025-01-15: user-profile-v2 - Add profile picture support
- **Type**: ADD
- **Scope**: ProfilePicture struct, UpdateProfilePicture method
- **Compatibility**: yes - additive change
- **Migration**: not required
- **Links**:
  - Feature: `features/user-profile-v2.md`
  - Spec: `spec/user-profile-v2.spec.md`

##### 2025-01-10: user-auth - Initial service creation
- **Type**: ADD
- **Scope**: GetUser, CreateUser, UpdateUser, DeleteUser
- **Compatibility**: N/A - initial version
- **Migration**: N/A
- **Links**:
  - Feature: `features/user-auth.md`
  - Spec: `spec/user-auth.spec.md`

#### Methods (summary)
- `GetUser`: Retrieve user by ID
- `CreateUser`: Create new user account
- `UpdateUser`: Update user profile
- `UpdateProfilePicture`: Update profile picture
- `DeleteUser`: Soft-delete user
