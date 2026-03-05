# 📋 Tasks: Logging

> **Feature**: `03` — Logging
> **Architecture**: [`03-logging-architecture.md`](03-logging-architecture.md)
> **Status**: 🟢 Ready for Build
> **Tasks**: 14 tasks across 4 phases

---

## Phase A — Core Logger Package (5 tasks)

- [ ] **A-01**: Create `core/logger/level.go` with `parseLevel()` function mapping string → `slog.Level`
- [ ] **A-02**: Create `core/logger/logger.go` with `Setup()` function (reads config, creates handler, sets default)
- [ ] **A-03**: Verify `Setup()` reads `LOG_LEVEL`, `LOG_FORMAT`, `LOG_OUTPUT` from config
- [ ] **A-04**: Verify file output creates `storage/logs/app.log` in append mode with `0644` permissions
- [ ] **A-05**: Run `go vet ./core/logger/...` — must pass with zero warnings

## Phase B — Main Entry Point Update (3 tasks)

- [ ] **B-01**: Update `cmd/main.go` to import and call `logger.Setup()` after `config.Load()`
- [ ] **B-02**: Add `slog.Info("server initialized", ...)` with app name, port, and env after banner
- [ ] **B-03**: Run `go build -o bin/rgo.exe ./cmd/` — must compile successfully

## Phase C — Tests (4 tasks)

- [ ] **C-01**: Create `core/logger/logger_test.go` with test for `parseLevel()` — all 4 levels + unknown fallback
- [ ] **C-02**: Add test for `Setup()` with `LOG_FORMAT=json` — verify JSON handler is used
- [ ] **C-03**: Add test for `Setup()` with `LOG_FORMAT=text` — verify text handler is used
- [ ] **C-04**: Add test for `Setup()` with `LOG_OUTPUT=file` — verify file is created at `storage/logs/app.log`

## Phase D — Integration Validation (2 tasks)

- [ ] **D-01**: Run `go test ./core/logger/... -v` — all tests must pass
- [ ] **D-02**: Run `go run ./cmd/` — verify banner displays AND structured log line appears
