# 💬 Discussion: Service Container

> **Feature**: `05` — Service Container
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/05-service-container`
> **Depends On**: #01 (Project Setup ✅)
> **Date Started**: 2026-03-05
> **Date Completed**: 2026-03-05

---

## Summary

Implement the core service container — a dependency injection (DI) mechanism that provides register/resolve patterns for managing services throughout the framework. This includes the `Container` struct with `Bind`, `Singleton`, `Instance`, `Make`, `MustMake`, and `Has` methods, the `Provider` interface with two-phase lifecycle (`Register`/`Boot`), and the `App` struct that orchestrates bootstrap. This is the **foundation for extensibility** — all subsequent features (database, sessions, cache, mail, events) will register their services through this container via providers (Feature #06).

---

## Functional Requirements

- As a **framework developer**, I want a `Container` struct so that I can register and resolve services by name
- As a **framework developer**, I want `Bind()` (transient), `Singleton()` (shared), and `Instance()` (pre-created) registration methods so that services can be created with the appropriate lifecycle
- As a **framework developer**, I want `Make()` to resolve services by name, checking instances first then bindings, panicking if not found
- As a **framework developer**, I want `MustMake[T]()` as a generic typed resolution helper so that resolved services are type-safe without manual casting
- As a **framework developer**, I want `Has()` to check if a service is registered without resolving it
- As a **framework developer**, I want a `Provider` interface with `Register()` and `Boot()` methods so that service registration follows a two-phase lifecycle
- As a **framework developer**, I want an `App` struct that manages provider registration, boot sequence, and service resolution so that `main.go` has a clean bootstrap flow

## Current State / Reference

### What Exists
- **Configuration** (#02 ✅): `config.Load()`, `config.Env()`, `config.IsDebug()` — available for providers to read config
- **Logging** (#03 ✅): `slog` globally configured — available for logging during boot
- **Error handling** (#04 ✅): `core/errors` package — available for error types
- **`core/container/` directory**: Empty placeholder created in Feature #01
- **`core/app/` directory**: Empty placeholder created in Feature #01
- **Blueprint**: Full `Container` implementation with 6 methods, `Provider` interface, `App` struct, and usage examples

### What Works Well
- Blueprint code is complete and well-structured — can be implemented as-is with minor adjustments
- `sync.RWMutex` for thread safety is idiomatic Go
- Two-phase provider lifecycle (`Register` → `Boot`) cleanly separates binding from initialization

### What Needs Improvement
- No DI container exists yet — services are created manually in `main.go`
- `cmd/main.go` currently calls `config.Load()` and `logger.Setup()` directly — will be restructured with `App` bootstrap in a future integration step

## Proposed Approach

Create three packages:

1. **`core/container`** — The DI container with `Bind`, `Singleton`, `Instance`, `Make`, `MustMake[T]`, `Has`
2. **`core/container` (same package)** — The `Provider` interface definition
3. **`core/app`** — The `App` struct with `New()`, `Register()`, `Boot()`, `Make()`

The `Provider` interface lives in the `container` package (not `app`) because providers depend on `Container`, and putting the interface there avoids circular imports.

**NOT in scope for this feature**:
- Built-in providers (DatabaseProvider, SessionProvider, etc.) — that's Feature #06
- Integration with `cmd/main.go` — `main.go` will be updated when providers are registered in Feature #06
- Gin router or HTTP-related services

## Edge Cases & Risks

- [x] `Make()` panics on unregistered service — intentional fail-fast behavior, documented
- [x] `MustMake[T]()` panics on type assertion failure — intentional, documented
- [x] `Singleton` race condition — the blueprint's `Singleton` closes over the factory and guards `instances` with mutex; needs verification under concurrent access
- [x] `Instance()` overwrites existing instance — acceptable behavior, allows test mocking
- [x] `Bind()` overwrites existing binding — acceptable, last-write-wins
- [x] Empty `Boot()` on providers that don't need it — required by interface, just use empty body

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #01 — Project Setup | Feature | ✅ Done |
| `sync` | Stdlib | ✅ Available |
| `fmt` | Stdlib | ✅ Available |

## Open Questions

_All resolved during discussion._

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-05 | `Provider` interface in `container` package, not `app` | Avoids circular imports; providers need `*Container` parameter |
| 2026-03-05 | `Make()` panics on missing service | Fail-fast is correct for DI — all services must be registered before boot. Return-error pattern would pollute every call site |
| 2026-03-05 | `MustMake[T]` as package-level generic function | Go generics don't support generic methods on structs; must be a package function |
| 2026-03-05 | `Singleton` uses closure-based lazy init with mutex | Blueprint pattern; singleton factory runs on first `Make()`, result cached in `instances` |
| 2026-03-05 | No built-in providers in this feature | Feature #06 adds providers; keeps #05 focused on the container mechanics |
| 2026-03-05 | Don't modify `cmd/main.go` yet | App bootstrap integration comes with Feature #06 when providers are registered |

## Discussion Complete ✅

**Summary**: Feature #05 implements the core DI container (`core/container`) with 6 methods, the `Provider` interface, and the `App` struct (`core/app`) for bootstrap orchestration. Scoped to container mechanics only — no built-in providers, no `main.go` changes.
**Completed**: 2026-03-05
**Next**: Create architecture doc → `05-service-container-architecture.md`
