# 💬 Discussion: Models (GORM)

> **Feature**: `11` — Models (GORM)
> **Status**: 🟢 COMPLETE
> **Branch**: `docs/11-models`
> **Depends On**: #09 (Database Connection ✅)
> **Date Started**: 2026-03-06
> **Date Completed**: 2026-03-06

---

## Summary

Define the GORM model layer for the RGo framework. This feature provides a `BaseModel` struct with common fields (ID, timestamps), a `User` model and a `Post` model as concrete examples demonstrating GORM tags, relationships (HasMany, BelongsTo), and JSON serialization. Models live in `database/models/` and embed `BaseModel` for consistent primary key and timestamp behavior.

---

## Functional Requirements

- As a **framework developer**, I want a `BaseModel` struct with `ID`, `CreatedAt`, `UpdatedAt` fields so that all models share consistent primary key and timestamp behavior
- As a **framework developer**, I want a `User` model with name, email, password, role, active fields and proper GORM tags so that user data is structured and validated at the database level
- As a **framework developer**, I want a `Post` model with title, slug, body, user_id fields so that content models demonstrate foreign key relationships
- As a **framework developer**, I want `User.Posts` (HasMany) and `Post.User` (BelongsTo) relationships defined via GORM tags so that relationship patterns are established for the framework
- As a **framework developer**, I want JSON tags on all model fields so that models serialize correctly for API responses
- As a **framework developer**, I want `Password` excluded from JSON output (`json:"-"`) so that sensitive data is never leaked in API responses
- As a **framework user**, I want example models to reference when creating my own models

## Current State / Reference

### What Exists
- **Database Connection (#09 ✅)**: `database.Connect()` returns `*gorm.DB`, connection pool configured
- **GORM (#09 ✅)**: `gorm.io/gorm v1.31.1` already installed
- **`database/models/`**: Empty directory with `.gitkeep`
- **`database/querybuilder/`**: Empty directory with `.gitkeep` — not in scope for this feature

### Blueprint Reference
The blueprint (Models (GORM) section, lines 2942–3002) shows:
1. `BaseModel` struct — `ID` (uint, primaryKey), `CreatedAt`, `UpdatedAt` (time.Time)
2. `User` struct — embeds `BaseModel`, fields: Name, Email, Password, Role, Active, Posts (HasMany)
3. `Post` struct — embeds `BaseModel`, fields: Title, Slug, Body, UserID, User (BelongsTo)
4. Relationship examples — `db.Preload("Posts").Find(&users)`, nested preloading
5. Hooks example — `BeforeCreate` for password hashing

### Blueprint Adaptations
| Blueprint | Our Implementation | Reason |
|---|---|---|
| `BeforeCreate` hook hashes password using `helpers.HashPassword()` | Deferred — no helpers package yet | Feature #19 (Helpers) and #22 (Crypto) are not built. Hook will be a placeholder or omitted. |
| Relationship usage examples (Preload, etc.) | Not in model files — documented only | Usage patterns are not model definitions; they belong in service/controller code. |

## Proposed Approach

### BaseModel Struct

A reusable base that all models embed. Uses GORM's convention for `ID`, `CreatedAt`, `UpdatedAt`. We define our own `BaseModel` instead of using `gorm.Model` because:
1. `gorm.Model` includes `DeletedAt` (soft deletes) — we want that as an opt-in concern (Feature #52, Soft Deletes)
2. Custom JSON tags — `gorm.Model` doesn't include JSON tags
3. Framework convention — users see the pattern and replicate it

### User Model

Standard user model with authentication-related fields. GORM tags define database constraints (size, unique index, not null, defaults). The `Password` field uses `json:"-"` to prevent JSON serialization.

### Post Model

Content model demonstrating foreign key relationships. `UserID` is the foreign key, `User` is the BelongsTo relationship, and `User.Posts` is the HasMany side.

### No Hooks in This Feature

The blueprint's `BeforeCreate` hook depends on `helpers.HashPassword()` which doesn't exist yet (Feature #19/#22). Hooks will be added when those dependencies are built. This keeps Feature #11 focused on model definitions only.

## Edge Cases & Risks

- [x] `gorm.Model` vs custom `BaseModel` — we use custom because `gorm.Model` includes `DeletedAt` and lacks JSON tags
- [x] Hook dependencies — `BeforeCreate` needs helpers package; deferred to later feature
- [x] Circular imports — `models` package must not import anything from `core/` or `app/`; it's a leaf package
- [x] JSON tags — `Password` must be `json:"-"` to prevent leaking in API responses

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #09 — Database Connection | Feature | ✅ Done |
| `gorm.io/gorm` | External | ✅ Installed (v1.31.1) |

## Open Questions

_(All resolved)_

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-06 | Custom `BaseModel` instead of `gorm.Model` | `gorm.Model` includes `DeletedAt` (soft deletes are Feature #52). Custom gives us JSON tags and opt-in soft deletes later. |
| 2026-03-06 | No `BeforeCreate` hook | Depends on helpers/crypto packages (#19, #22) not yet built. Will be added in those features. |
| 2026-03-06 | User + Post as concrete models | Blueprint defines these two. They demonstrate HasMany/BelongsTo patterns and serve as examples. |
| 2026-03-06 | Models are pure data structs — no business logic | Models define schema only. Validation, services, and business logic belong in their respective layers. |

## Discussion Complete ✅

**Summary**: Feature #11 creates `BaseModel`, `User`, and `Post` model structs in `database/models/` with GORM tags, relationships, and JSON serialization. No hooks (deferred), no migrations (Feature #12), pure data definition.
**Completed**: 2026-03-06
**Next**: Create architecture doc → `11-models-architecture.md`
