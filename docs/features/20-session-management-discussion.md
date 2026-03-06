# 💬 Discussion: Session Management

> **Feature**: `20` — Session Management
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06
> **Depends on**: #02 (Config), #08 (Middleware), #09 (Database)

---

## Scope

Feature #20 implements the full session management system with a driver-based architecture. All five storage backends are included, plus the `Manager`, `Store` interface, `NewStore` factory, session middleware, flash messages, and a `SessionProvider`.

### What's In

| Component | Location |
|---|---|
| `Store` interface | `core/session/store.go` |
| `MemoryStore` | `core/session/memory_store.go` |
| `FileStore` | `core/session/file_store.go` |
| `DBStore` + `SessionRecord` | `core/session/db_store.go` |
| `CookieStore` | `core/session/cookie_store.go` |
| `RedisStore` | `core/session/redis_store.go` |
| `NewStore` factory | `core/session/factory.go` |
| `Manager` (Start/Save/Destroy/Flash) | `core/session/manager.go` |
| `SessionMiddleware` | `core/middleware/session.go` |
| `SessionProvider` | `app/providers/session_provider.go` |
| Provider registration in `NewApp()` | `core/cli/root.go` |

### What's Out

- Redis as a new dependency (`github.com/redis/go-redis/v9`) — the `RedisStore` will be implemented but the dependency will be **optional/deferred**. We compile-gate it behind a build tag `redis` to avoid forcing every user to pull in Redis. **Decision**: For simplicity in this initial implementation, we will skip the Redis store entirely and add it in a future feature when Redis is introduced as a proper dependency. The factory will return an error for `SESSION_DRIVER=redis`.
- Cookie store encryption — requires `APP_KEY` and AES-256-GCM. We will implement a simple version using `crypto/aes` + `crypto/cipher` from stdlib (no new deps).
- Session table migration via `make:migration` — the `DBStore` uses `AutoMigrate` for the `SessionRecord` model.

### Key Decisions

1. **No Redis dependency**: `RedisStore` is **deferred** to avoid adding `go-redis/v9`. The factory returns a descriptive error for `SESSION_DRIVER=redis`.
2. **CookieStore**: Implemented using `crypto/aes` + `crypto/cipher` (AES-256-GCM) from stdlib. Requires a 32-byte `APP_KEY`.
3. **MemoryStore**: Used as default for tests. The `SessionProvider` will default to `memory` driver when `SESSION_DRIVER` is not set.
4. **Provider order**: SessionProvider registers after DatabaseProvider (position 4, before Middleware), since session middleware needs the manager.
5. **Flash messages**: `Flash`, `GetFlash`, `FlashErrors`, `FlashOldInput` are methods on `Manager` operating on session data maps.
