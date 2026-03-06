# 🧪 Test Plan: Models (GORM)

> **Feature**: `11` — Models (GORM)
> **Tasks**: [`11-models-tasks.md`](11-models-tasks.md)
> **Date**: 2026-03-06

---

## Acceptance Criteria

- [ ] `BaseModel` has `ID` (uint), `CreatedAt` (time.Time), `UpdatedAt` (time.Time) fields
- [ ] `User` embeds `BaseModel` and has Name, Email, Password, Role, Active, Posts fields
- [ ] `Post` embeds `BaseModel` and has Title, Slug, Body, UserID, User fields
- [ ] `Password` field has `json:"-"` tag (never serialized)
- [ ] GORM AutoMigrate succeeds for User and Post models (SQLite :memory:)
- [ ] GORM can create and query User records
- [ ] GORM can create Post with foreign key to User
- [ ] GORM Preload loads User→Posts relationship
- [ ] All tests pass with `go test ./database/models/...`
- [ ] All tests pass with `go test ./...` (full regression)
- [ ] `go vet ./...` reports no issues

---

## Test Cases

### TC-01: BaseModel has expected fields

**File**: `database/models/models_test.go`
**Function**: `TestBaseModel_Fields`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `BaseModel{}` | No error |
| 2 | Assert `ID` field exists and is `uint` | Pass (via reflect) |
| 3 | Assert `CreatedAt` field exists and is `time.Time` | Pass |
| 4 | Assert `UpdatedAt` field exists and is `time.Time` | Pass |

---

### TC-02: User embeds BaseModel

**File**: `database/models/models_test.go`
**Function**: `TestUser_EmbedsBaseModel`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `User{BaseModel: BaseModel{ID: 42}}` | No error |
| 2 | Assert `user.ID == 42` | Pass |

---

### TC-03: User password excluded from JSON

**File**: `database/models/models_test.go`
**Function**: `TestUser_PasswordExcludedFromJSON`

| Step | Action | Expected |
|---|---|---|
| 1 | Create `User{Password: "secret"}` with Name, Email set | No error |
| 2 | Marshal to JSON | No error |
| 3 | Assert JSON does not contain `"password"` key | Pass |
| 4 | Assert JSON contains `"name"` key | Pass |

---

### TC-04: GORM AutoMigrate succeeds

**File**: `database/models/models_test.go`
**Function**: `TestModels_AutoMigrate`

| Step | Action | Expected |
|---|---|---|
| 1 | Open SQLite `:memory:` via GORM | No error |
| 2 | Run `db.AutoMigrate(&User{}, &Post{})` | No error |

---

### TC-05: GORM creates and queries User

**File**: `database/models/models_test.go`
**Function**: `TestUser_CreateAndQuery`

| Step | Action | Expected |
|---|---|---|
| 1 | AutoMigrate User, Post | No error |
| 2 | Create `User{Name: "Alice", Email: "alice@test.com", Password: "hash"}` | No error |
| 3 | Query user by email | Found, Name == "Alice" |
| 4 | Assert `user.ID > 0` | Pass (auto-incremented) |
| 5 | Assert `user.CreatedAt` is not zero | Pass |

---

### TC-06: GORM creates Post with foreign key

**File**: `database/models/models_test.go`
**Function**: `TestPost_CreateWithForeignKey`

| Step | Action | Expected |
|---|---|---|
| 1 | AutoMigrate User, Post | No error |
| 2 | Create User | No error |
| 3 | Create `Post{Title: "Hello", Slug: "hello", Body: "World", UserID: user.ID}` | No error |
| 4 | Query post | Found, `post.UserID == user.ID` |

---

### TC-07: GORM Preload loads User→Posts

**File**: `database/models/models_test.go`
**Function**: `TestUser_PreloadPosts`

| Step | Action | Expected |
|---|---|---|
| 1 | AutoMigrate, create User, create 2 Posts for that user | No error |
| 2 | Query user with `Preload("Posts")` | No error |
| 3 | Assert `len(user.Posts) == 2` | Pass |

---

### TC-08: Post defaults and constraints

**File**: `database/models/models_test.go`
**Function**: `TestUser_Defaults`

| Step | Action | Expected |
|---|---|---|
| 1 | AutoMigrate, create User with only required fields (Name, Email, Password) | No error |
| 2 | Query user | Found |
| 3 | Assert `user.Role == "user"` | Pass (default) |
| 4 | Assert `user.Active == true` | Pass (default) |

---

## Test Summary

| TC | Function | Type | Status |
|---|---|---|---|
| TC-01 | `TestBaseModel_Fields` | Unit | ⬜ |
| TC-02 | `TestUser_EmbedsBaseModel` | Unit | ⬜ |
| TC-03 | `TestUser_PasswordExcludedFromJSON` | Unit | ⬜ |
| TC-04 | `TestModels_AutoMigrate` | Integration | ⬜ |
| TC-05 | `TestUser_CreateAndQuery` | Integration | ⬜ |
| TC-06 | `TestPost_CreateWithForeignKey` | Integration | ⬜ |
| TC-07 | `TestUser_PreloadPosts` | Integration | ⬜ |
| TC-08 | `TestUser_Defaults` | Integration | ⬜ |

**Total**: 8 test cases
