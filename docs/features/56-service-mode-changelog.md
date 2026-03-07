# 📝 Changelog: Service Mode

> **Feature**: `56` — Service Mode Architecture
> **Branch**: `feature/56-service-mode`

---

## Build Log

- **Phase A** — Created `core/service/mode.go` (Mode bitmask, ParseMode, Has, String, Services, PortEnvKey), `core/service/mode_test.go` (6 test functions, 34+ sub-tests), `routes/ws.go` (empty placeholder). Checkpoint: all tests pass, build compiles.
- **Phase B** — Updated `root.go` (NewApp(mode)), `serve.go` (complete rewrite with --mode flag, RunE, serveSingle/serveMulti), `router_provider.go` (Mode field, conditional routes/templates/static), `middleware_provider.go` (Mode field, conditional CSRF), `.env` (RAPIDGO_MODE + port vars), 4 CLI callers (ModeAll), 2 test files. Checkpoint: all 31 packages pass.
- **Phase C** — Added `ServiceConfig` struct and `ListenAndServeMulti()` to `core/server/server.go`, 3 new server tests (TC-04/05/06). Checkpoint: all 31 packages pass.
- **Impl Cross-Check** — Subagent found 2 deviations; both fixed. Full regression passed.
- **Binary Verification** — `rapidgo version` works, `--mode` flag visible in serve help, `--mode=invalid` fails fast with clear error.

---

## Deviations from Plan

| Planned | Actual | Reason |
|---|---|---|
| `applyRoutesForMode()` handles only route registration | Also sets up templates and static files for web mode | Multi-port web mode needs its own Gin router with templates/static loaded — caught during impl cross-check |
| Server tests call `ListenAndServeMulti()` directly | Tests use `http.Server` directly to avoid signal handling in CI | Consistent with existing server test pattern (TC-10 note in testplan). `ListenAndServeMulti` calls `signal.NotifyContext` which conflicts with test harness |
| `config.Load()` called inside `NewApp()` | `config.Load()` called in `serve.RunE` before `config.Env("RAPIDGO_MODE")` | Must load env before reading env var — architecture doc already specified this order |

---

## Decisions During Build

| Decision | Context | Date |
|---|---|---|
| Use `RunE` instead of `Run` for serve command | Needed to return `ParseMode` error to Cobra for clean CLI error display | Build phase |
| `allSamePort()` check before multi-port path | Avoid unnecessary multi-server overhead when all modes resolve to same port | Build phase |
| Health routes registered unconditionally if DB exists | Health check useful regardless of mode — no reason to restrict | Build phase |
