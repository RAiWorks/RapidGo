# 📋 Tasks: Configuration System

> **Feature**: `02` — Configuration System
> **Architecture**: [`02-configuration-architecture.md`](02-configuration-architecture.md)
> **Status**: ✅ Complete
> **Tasks**: 18/18 tasks complete across 5 phases

---

## Phase A — Dependency Setup (2 tasks)

- [x] **A-01**: Run `go get github.com/joho/godotenv` to add first third-party dependency
- [x] **A-02**: Verify `go.mod` has `require` block with godotenv entry and `go.sum` is created

## Phase B — Core Config Package (5 tasks)

- [x] **B-01**: Create `core/config/config.go` with `Load()` function
- [x] **B-02**: Create `core/config/env.go` with `Env()`, `EnvInt()`, `EnvBool()` helpers
- [x] **B-03**: Create `core/config/environment.go` with `AppEnv()`, `IsProduction()`, `IsDevelopment()`, `IsTesting()`, `IsDebug()`
- [x] **B-04**: Verify all exported functions have correct signatures matching architecture doc
- [x] **B-05**: Run `go vet ./core/config/...` — must pass with zero warnings

## Phase C — Main Entry Point Update (3 tasks)

- [x] **C-01**: Update `cmd/main.go` to import and call `config.Load()` at start
- [x] **C-02**: Update `cmd/main.go` banner to display `APP_NAME`, `APP_PORT`, `APP_ENV`, `IsDebug()`
- [x] **C-03**: Run `go build -o bin/RapidGo.exe ./cmd/` — must compile successfully

## Phase D — Tests (5 tasks)

- [x] **D-01**: Create `core/config/config_test.go` with tests for `Load()` (with and without `.env`)
- [x] **D-02**: Add tests for `Env()` — key present, key absent (fallback), empty string
- [x] **D-03**: Add tests for `EnvInt()` — valid int, invalid string, empty (fallback)
- [x] **D-04**: Add tests for `EnvBool()` — "true", "1", "false", "0", empty (fallback)
- [x] **D-05**: Add tests for `AppEnv()`, `IsProduction()`, `IsDevelopment()`, `IsTesting()`, `IsDebug()`

## Phase E — Integration Validation (3 tasks)

- [x] **E-01**: Run `go test ./core/config/... -v` — all tests must pass
- [x] **E-02**: Run `go run ./cmd/` — verify banner shows `.env` values (APP_NAME=RapidGo, APP_PORT=8080, etc.)
- [x] **E-03**: Run `go vet ./...` — full project must pass with zero warnings
