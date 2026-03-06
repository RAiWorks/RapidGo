# 📋 Tasks: Session Management

> **Feature**: `20` — Session Management
> **Architecture**: [`20-session-management-architecture.md`](20-session-management-architecture.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Phase A — Implementation

| # | Task | File | Est. |
|---|---|---|---|
| A1 | Create `Store` interface | `core/session/store.go` | 3m |
| A2 | Create `MemoryStore` | `core/session/memory_store.go` | 5m |
| A3 | Create `FileStore` | `core/session/file_store.go` | 5m |
| A4 | Create `DBStore` + `SessionRecord` | `core/session/db_store.go` | 8m |
| A5 | Create `CookieStore` (AES-256-GCM) | `core/session/cookie_store.go` | 10m |
| A6 | Create `NewStore` factory | `core/session/factory.go` | 5m |
| A7 | Create `Manager` + flash messages | `core/session/manager.go` | 10m |
| A8 | Create `SessionMiddleware` | `core/middleware/session.go` | 5m |
| A9 | Create `SessionProvider` | `app/providers/session_provider.go` | 3m |
| A10 | Register `SessionProvider` in `NewApp()` | `core/cli/root.go` | 2m |

## Phase B — Tests

| # | Task | File | Est. |
|---|---|---|---|
| B1 | Test all stores (Memory, File, Cookie), Manager, flash messages | `core/session/session_test.go` | 20m |
| B2 | Test `SessionMiddleware` | `core/middleware/middleware_test.go` (extend) | 5m |

## Phase C — Quality Assurance

| # | Task | Est. |
|---|---|---|
| C1 | Run full test suite — 0 failures | 3m |
| C2 | Cross-check implementation vs architecture | 5m |
| C3 | Document deviations (if any) | 3m |
| C4 | Commit and push to `feature/20-session-management` | 2m |
| C5 | Merge to main | 2m |
