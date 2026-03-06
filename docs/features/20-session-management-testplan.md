# 🧪 Test Plan: Session Management

> **Feature**: `20` — Session Management
> **Architecture**: [`20-session-management-architecture.md`](20-session-management-architecture.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Test File

`core/session/session_test.go` — package `session`

---

## Test Cases

### MemoryStore

| # | Test | Description | Assert |
|---|---|---|---|
| T01 | `TestMemoryStore_WriteRead` | Write then read data | Data matches |
| T02 | `TestMemoryStore_ReadMissing` | Read non-existent session | Returns empty map |
| T03 | `TestMemoryStore_Destroy` | Destroy then read | Returns empty map |
| T04 | `TestMemoryStore_GC` | Write expired entries, run GC | Expired entries removed |

### FileStore

| # | Test | Description | Assert |
|---|---|---|---|
| T05 | `TestFileStore_WriteRead` | Write then read from temp dir | Data matches |
| T06 | `TestFileStore_ReadMissing` | Read non-existent session | Returns empty map |
| T07 | `TestFileStore_Destroy` | Destroy then read | Returns empty map |

### CookieStore

| # | Test | Description | Assert |
|---|---|---|---|
| T08 | `TestCookieStore_WriteRead` | Encrypt, write, read, decrypt | Data matches |
| T09 | `TestCookieStore_ReadMissing` | Read non-existent | Returns empty map |
| T10 | `TestCookieStore_InvalidKey` | Create with wrong key length | Returns error |
| T11 | `TestCookieStore_Destroy` | Destroy then read | Returns empty map |

### Manager

| # | Test | Description | Assert |
|---|---|---|---|
| T12 | `TestManager_StartNewSession` | No cookie sent | Returns new ID + empty data |
| T13 | `TestManager_StartExistingSession` | Cookie with stored session | Returns existing data |
| T14 | `TestManager_Save` | Save writes cookie and data | Cookie set, data persisted |
| T15 | `TestManager_Destroy` | Destroy clears cookie and data | Cookie cleared, data removed |

### Flash Messages

| # | Test | Description | Assert |
|---|---|---|---|
| T16 | `TestManager_Flash` | Set a flash message | Flash stored in `_flashes` |
| T17 | `TestManager_GetFlash` | Get and remove flash | Value returned, key removed |
| T18 | `TestManager_GetFlash_Missing` | Get non-existent flash | Returns nil, false |
| T19 | `TestManager_FlashErrors` | Store validation errors | Stored under `_errors` key |
| T20 | `TestManager_FlashOldInput` | Store old input | Stored under `_old_input` key |

### Factory

| # | Test | Description | Assert |
|---|---|---|---|
| T21 | `TestNewStore_Memory` | `SESSION_DRIVER=memory` | Returns `*MemoryStore` |
| T22 | `TestNewStore_File` | `SESSION_DRIVER=file` | Returns `*FileStore` |
| T23 | `TestNewStore_Unsupported` | `SESSION_DRIVER=invalid` | Returns error |
| T24 | `TestNewStore_Redis` | `SESSION_DRIVER=redis` | Returns error (not supported) |

### SessionMiddleware

| # | Test (in `middleware_test.go`) | Description | Assert |
|---|---|---|---|
| T25 | `TestSessionMiddleware_SetsSessionData` | Middleware sets session context | `session_id` and `session` keys present |

---

## Coverage Summary

| Package | Functions | Tests | Target |
|---|---|---|---|
| `core/session` | 25+ | 24 | 95%+ |
| `core/middleware` | 1 | 1 | 100% |

**Total tests after feature**: 211 (existing) + 25 = **236**
