# 💬 Discussion: Logging

> **Feature**: `03` — Logging
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/03-logging`
> **Depends On**: #01 (Project Setup & Structure ✅), #02 (Configuration System ✅)
> **Date Started**: 2026-03-05
> **Date Completed**: 2026-03-05

---

## Summary

Implement the framework's structured logging system using Go's standard library `log/slog`. The logger reads config from `.env` (`LOG_LEVEL`, `LOG_FORMAT`, `LOG_OUTPUT`), supports JSON and text output formats, environment-based defaults, and file output to `storage/logs/`. This is a zero-dependency feature (slog is stdlib since Go 1.21).

---

## Functional Requirements

- As a framework developer, I want a `Setup()` function that initializes the global slog logger so that all subsequent `slog.Info()` / `slog.Error()` calls use the configured handler
- As a framework developer, I want the log level configurable via `LOG_LEVEL` env var so that I can control verbosity per environment
- As a framework developer, I want JSON format in production and text format in development so that logs are machine-parseable in prod and human-readable in dev
- As a framework developer, I want the log format configurable via `LOG_FORMAT` env var so that I can override the environment default
- As a framework developer, I want log output configurable via `LOG_OUTPUT` env var (`stdout` or `file`) so that production can write logs to `storage/logs/`
- As a framework developer, I want `Setup()` to return `*slog.Logger` so that components can hold a logger reference if needed
- As a framework developer, I want `slog.SetDefault()` called during setup so that the global `slog.Info()` / `slog.Error()` functions work everywhere without passing a logger instance

## Current State / Reference

### What Exists

- `core/logger/` directory with `.gitkeep` (Feature #01)
- `storage/logs/` directory with `.gitkeep` (Feature #01)
- `.env` has `LOG_LEVEL=debug`, `LOG_FORMAT=json`, `LOG_OUTPUT=stdout`
- `core/config/` package fully functional (Feature #02) — `Env()`, `EnvBool()`, `IsProduction()`, `IsDevelopment()`, `IsTesting()`
- Blueprint shows a minimal `Setup()` with hardcoded JSON/Info
- Framework reference doc shows environment-based handler switching

### What Works Well

- `log/slog` is stdlib — zero external dependencies
- Config system already provides `Env()` for reading LOG_* keys
- `IsProduction()` / `IsDevelopment()` available for environment-based defaults

### What Needs Improvement

N/A — greenfield code in `core/logger/`.

## Proposed Approach

1. **Create `core/logger/logger.go`** — `Setup()` function that:
   - Reads `LOG_LEVEL` from config, maps to `slog.Level` (default: `debug` in dev, `info` in prod)
   - Reads `LOG_FORMAT` from config (`json` or `text`), creates appropriate handler
   - Reads `LOG_OUTPUT` from config (`stdout` or `file`), selects output writer
   - If `file`, opens `storage/logs/app.log` with append mode
   - Calls `slog.SetDefault(logger)` to set global default
   - Returns `*slog.Logger`
2. **Create `core/logger/level.go`** — `parseLevel()` helper to map string → `slog.Level`
3. **Update `cmd/main.go`** — Call `logger.Setup()` after `config.Load()`, replace `fmt.Println` banner with `slog.Info` for the startup message
4. **Create `core/logger/logger_test.go`** — Unit tests

### Why NOT zerolog/zap?

The blueprint recommends `log/slog` as the primary choice. zerolog and zap are listed as alternatives "when profiling shows logging is a bottleneck." slog is sufficient for the framework's needs and keeps the dependency count at zero for this feature.

## Edge Cases & Risks

- [x] **Invalid `LOG_LEVEL` value** — Default to `info` if the string doesn't match any known level. Don't panic.
- [x] **Invalid `LOG_FORMAT` value** — Default to `json` if not `json` or `text`. JSON is the safer default.
- [x] **`LOG_OUTPUT=file` but `storage/logs/` doesn't exist** — Create the directory with `os.MkdirAll`. Log a warning if file creation fails, fall back to stdout.
- [x] **File permissions** — Use `0644` for log files (readable by owner and group, as recommended by docs). The security note about 0600/0640 applies to environments with strict multitenancy — 0644 is fine for single-app deployments.
- [x] **No log rotation** — Out of scope for this feature. Log rotation is handled externally (logrotate on Linux, or a future feature). Log file is opened in append mode.
- [x] **Request logging middleware** — Out of scope. Mentioned in reference doc as a usage pattern, not a core logger feature. Will be implemented in Feature #08 (Middleware Pipeline).
- [x] **`cmd/main.go` banner** — Keep `fmt.Println` for the ASCII banner (it's display output, not a log entry). Add one `slog.Info("server initialized", ...)` after the banner.

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Feature #01 — Project Setup & Structure | Feature | ✅ Done |
| Feature #02 — Configuration System | Feature | ✅ Done |
| `log/slog` | Standard Library (Go 1.21+) | ✅ Available |

## Open Questions

All resolved:

- [x] **slog, zerolog, or zap?** → slog. Stdlib, zero deps, sufficient for framework needs.
- [x] **Should `Setup()` accept parameters or read config directly?** → Read config directly via `config.Env()`. Keeps the call site simple: `logger.Setup()` with no args.
- [x] **Should we support multiple output destinations (stdout + file)?** → No. Single output for now. `stdout` or `file`, not both. Multi-output can be added later with `io.MultiWriter`.
- [x] **Should the banner use slog or fmt?** → Keep `fmt` for the ASCII banner (display), add `slog.Info` for the structured "server initialized" log entry.
- [x] **File logging path** — `storage/logs/app.log`. Simple, predictable, matches the `storage/logs/` directory from Feature #01.
- [x] **Log rotation** — Out of scope. External tooling (logrotate) handles this. File opened in append mode.

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-05 | Use `log/slog` (stdlib) | Zero deps, Go 1.21+ standard, blueprint recommends it |
| 2026-03-05 | Config-driven setup via `LOG_LEVEL`, `LOG_FORMAT`, `LOG_OUTPUT` | `.env` already has these keys; config system available |
| 2026-03-05 | JSON default format | Machine-parseable, safer default for production |
| 2026-03-05 | Text format for dev when `LOG_FORMAT=text` | Human-readable during development |
| 2026-03-05 | File output to `storage/logs/app.log` | Matches project structure, append mode, no rotation |
| 2026-03-05 | `parseLevel()` maps string → `slog.Level` | Clean separation, fallback to `slog.LevelInfo` on unknown |
| 2026-03-05 | Keep `fmt` for banner, add `slog.Info` for structured log | Banner is display output, not a log event |

## Discussion Complete ✅

**Summary**: Feature #03 creates the `core/logger` package with config-driven `Setup()` using `log/slog`, supporting JSON/text formats, stdout/file output, and configurable log levels. Zero new dependencies.
**Completed**: 2026-03-05
**Next**: Create architecture doc → `03-logging-architecture.md`
