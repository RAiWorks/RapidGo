# 📝 Changelog: Soft Deletes

> **Feature**: `52` — Soft Deletes
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-08
> **Commit**: _pending_

---

## Added

- `database/models/base.go` — `DeletedAt gorm.DeletedAt` field added to `BaseModel` with `gorm:"index"` and `json:"deleted_at,omitempty"` tags
- `database/models/scopes.go` — `WithTrashed` and `OnlyTrashed` GORM scope helpers
- `database/models/scopes_test.go` — 5 tests for scope behavior (T01–T05)
- `database/migrations/20260308000001_add_soft_deletes.go` — migration adding `deleted_at` column + index to `users` and `posts` tables
- `app/services/user_service.go` — `HardDelete(id)` and `Restore(id)` methods on `UserService`

## Changed

- `database/models/base.go` — `BaseModel` now includes `DeletedAt` (breaking: `Delete()` calls become soft deletes)
- `app/services/user_service_test.go` — TC-08 updated for soft delete behavior (T12); 6 new test cases added (T06–T11)

## Files

| File | Action |
|---|---|
| `database/models/base.go` | MODIFIED |
| `database/models/scopes.go` | NEW |
| `database/models/scopes_test.go` | NEW |
| `database/migrations/20260308000001_add_soft_deletes.go` | NEW |
| `app/services/user_service.go` | MODIFIED |
| `app/services/user_service_test.go` | MODIFIED |

## Migration Guide

**Breaking change**: `db.Delete(&Model{}, id)` now performs a soft delete (sets `deleted_at`) instead of permanently removing the row.

- To permanently delete: use `db.Unscoped().Delete(&Model{}, id)` or `UserService.HardDelete(id)`
- To include deleted records in queries: use `db.Scopes(models.WithTrashed).Find(&results)`
- To query only deleted records: use `db.Scopes(models.OnlyTrashed).Find(&results)`
- To restore a soft-deleted record: use `UserService.Restore(id)`
