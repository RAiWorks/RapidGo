# 📋 Review: Database Migrations

> **Feature**: `12` — Database Migrations
> **Branch**: `feature/12-database-migrations`
> **Merged**: 2026-03-06
> **Commit**: `a0bc022` (main)

---

## Summary

Feature #12 adds a two-tier database migration system: GORM `AutoMigrate` for model-to-table sync, plus a file-based migration engine with Up/Down functions for changes AutoMigrate can't handle. Four CLI commands drive the system. A central model registry keeps migration targets in one place.

## Files Changed

| File | Type | Description |
|---|---|---|
| `database/models/registry.go` | Created | `All()` — returns all models for AutoMigrate |
| `database/migrations/migrator.go` | Created | Migration engine: `Run()`, `Rollback()`, `Status()`, `Register()`, `ResetRegistry()` |
| `core/cli/migrate.go` | Created | `RapidGo migrate` — AutoMigrate + pending file-based migrations |
| `core/cli/migrate_rollback.go` | Created | `RapidGo migrate:rollback` — undo last batch |
| `core/cli/migrate_status.go` | Created | `RapidGo migrate:status` — applied/pending table |
| `core/cli/make_migration.go` | Created | `RapidGo make:migration <name>` — generate timestamped migration file |
| `core/cli/root.go` | Modified | Added 4 commands to `init()` |
| `core/cli/cli_test.go` | Modified | Added TC-10 (`TestToSnakeCase`) |
| `database/migrations/migrations_test.go` | Created | 9 tests for migration engine |

## Dependencies Added

None — all dependencies already present from Features #09 and #10.

## Test Results

| Package | Tests | Status |
|---|---|---|
| `database/migrations` | 9 | ✅ PASS |
| `core/cli` (TC-10) | 1 | ✅ PASS |
| **Full regression** | **141** | **✅ PASS** |

## Deviations from Plan

| What Changed | Why |
|---|---|
| Added `ResetRegistry()` | Tests need isolated registry state |
| TC-08 version names adjusted | Original names sorted incorrectly for index-based assertions |
