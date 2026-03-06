# 📝 Changelog: Session Management

> **Feature**: `20` — Session Management
> **Date**: 2026-03-06

---

## Added

- `core/session/store.go` — `Store` interface (Read, Write, Destroy, GC)
- `core/session/memory_store.go` — `MemoryStore` (dev/testing, sync.RWMutex)
- `core/session/file_store.go` — `FileStore` (JSON files, 0600 perms)
- `core/session/db_store.go` — `DBStore` + `SessionRecord` (GORM, AutoMigrate)
- `core/session/cookie_store.go` — `CookieStore` (AES-256-GCM, stdlib only)
- `core/session/factory.go` — `NewStore` factory (memory, file, db, cookie; redis deferred)
- `core/session/manager.go` — `Manager` with Start/Save/Destroy + Flash/GetFlash/FlashErrors/FlashOldInput
- `core/middleware/session.go` — `SessionMiddleware` (auto load/save per request)
- `app/providers/session_provider.go` — `SessionProvider` (registers session manager as "session")
- `core/session/session_test.go` — 24 tests for stores, manager, flash, factory

## Changed

- `core/cli/root.go` — Added `SessionProvider` registration after `DatabaseProvider`

## Deviations

_None documented yet. Will be updated during BUILD phase._
