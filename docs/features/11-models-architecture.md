# 🏗️ Architecture: Models (GORM)

> **Feature**: `11` — Models (GORM)
> **Discussion**: [`11-models-discussion.md`](11-models-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

The Models feature defines GORM model structs in `database/models/`. A reusable `BaseModel` provides consistent primary key and timestamp fields. `User` and `Post` models demonstrate GORM tags, constraints, JSON serialization, and relationship patterns (HasMany, BelongsTo). Models are pure data definitions — no business logic, no hooks, no migrations.

## File Structure

```
database/models/
├── base.go          # BaseModel struct (ID, CreatedAt, UpdatedAt)
├── user.go          # User model (name, email, password, role, active, posts)
└── post.go          # Post model (title, slug, body, user_id, user relationship)

database/models/
└── models_test.go   # Tests for model struct definitions
```

### Files Created (3)
| File | Package | Lines (est.) |
|---|---|---|
| `database/models/base.go` | `models` | ~15 |
| `database/models/user.go` | `models` | ~15 |
| `database/models/post.go` | `models` | ~15 |

### Files Modified (0)
No existing files are modified.

---

## Component Design

### BaseModel (`database/models/base.go`)

**Responsibility**: Provide common fields for all models — primary key and timestamps.
**Package**: `models`

```go
package models

import "time"

// BaseModel provides common fields for all models.
// Embed this in your model structs instead of gorm.Model
// to get ID, CreatedAt, and UpdatedAt without soft deletes.
type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

**Design notes**:
- Custom struct instead of `gorm.Model` — excludes `DeletedAt` (soft deletes are Feature #52, opt-in)
- JSON tags included — `gorm.Model` omits these
- `ID` is `uint` with `gorm:"primaryKey"` — GORM convention for auto-increment
- `CreatedAt` / `UpdatedAt` — automatically managed by GORM

### User Model (`database/models/user.go`)

**Responsibility**: Define the user data schema with authentication fields and relationship to posts.
**Package**: `models`

```go
package models

// User represents an application user.
type User struct {
	BaseModel
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     string `gorm:"size:50;default:user" json:"role"`
	Active   bool   `gorm:"default:true" json:"active"`
	Posts    []Post  `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}
```

**Design notes**:
- `Password` uses `json:"-"` — never serialized in API responses
- `Email` has `uniqueIndex` — database-level uniqueness constraint
- `Role` defaults to `"user"` — standard role for new registrations
- `Active` defaults to `true` — accounts are active by default
- `Posts` is the HasMany side — `foreignKey:UserID` tells GORM which FK to use
- `Posts` uses `json:"posts,omitempty"` — only included when explicitly preloaded

### Post Model (`database/models/post.go`)

**Responsibility**: Define content model with foreign key relationship to user.
**Package**: `models`

```go
package models

// Post represents a content entry authored by a user.
type Post struct {
	BaseModel
	Title  string `gorm:"size:255;not null" json:"title"`
	Slug   string `gorm:"size:255;uniqueIndex" json:"slug"`
	Body   string `gorm:"type:text" json:"body"`
	UserID uint   `gorm:"index;not null" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

**Design notes**:
- `Slug` has `uniqueIndex` — URL-friendly unique identifier
- `Body` uses `gorm:"type:text"` — long content, not VARCHAR
- `UserID` is the foreign key with `index` — optimizes JOIN queries
- `User` is the BelongsTo side — loaded via `db.Preload("User")`
- `User` uses `json:"user,omitempty"` — only included when preloaded

---

## Blueprint Adaptations

| Blueprint | Our Implementation | Reason |
|---|---|---|
| `BeforeCreate` hook with `helpers.HashPassword()` | Omitted — deferred to Feature #19/#22 | Dependencies don't exist yet. Hook will be added when Helpers/Crypto are built. |
| Relationship usage examples (Preload) | Not in model files | Usage patterns belong in services/controllers, not model definitions. Documented in discussion. |
| All models in one file | Split into `base.go`, `user.go`, `post.go` | Go convention — one type per file for clarity and maintainability. |

---

## Impact on Existing Code

| Area | Impact |
|---|---|
| `database/models/` | **New files** — replaces `.gitkeep` with actual model definitions |
| `database/connection.go` | **No change** — models don't import the connection package |
| Provider chain | **No change** — no new providers needed |
| `cmd/main.go` | **No change** |
| `.env` | **No change** |
