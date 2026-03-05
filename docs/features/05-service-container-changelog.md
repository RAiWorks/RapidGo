# 📝 Changelog: Service Container

> **Feature**: `05` — Service Container
> **Branch**: `feature/05-service-container`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

- **2026-03-06** — Created `core/container/container.go`: `Container` struct, `Factory` type, `New()`, `Bind()`, `Singleton()`, `Instance()`, `Make()`, `MustMake[T]()`, `Has()`
- **2026-03-06** — Created `core/container/provider.go`: `Provider` interface with `Register()` and `Boot()`
- **2026-03-06** — Created `core/app/app.go`: `App` struct with `New()`, `Register()`, `Boot()`, `Make()`
- **2026-03-06** — Created `core/container/container_test.go`: 14 tests (TC-01–TC-12, TC-16, TC-17)
- **2026-03-06** — Created `core/app/app_test.go`: 3 tests (TC-13–TC-15)
- **2026-03-06** — All 17 tests pass, `go vet` clean
- **2026-03-06** — `-race` unavailable (no GCC on Windows) — concurrency test (TC-12) passes without race detector

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| Singleton double-check locking | Blueprint pattern: no re-check after write lock | Added re-check of `instances` after acquiring write lock in `Singleton` | Prevents factory from executing more than once under concurrent access — proper double-check locking pattern |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| Added double-check in Singleton after Lock | Cross-check identified that without re-check, two goroutines could both call the factory | 2026-03-06 |
| Skip `-race` flag | Windows machine lacks GCC/CGO — race detector requires CGO_ENABLED=1 with C compiler | 2026-03-06 |
