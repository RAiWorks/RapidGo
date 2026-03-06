# Feature #37 — Graceful Shutdown: Review

## Summary

Replaced Gin's `r.Run()` with a proper `http.Server` featuring configurable timeouts and OS signal-based graceful shutdown.

## Delivered

| Item | Detail |
|------|--------|
| Package | `core/server/server.go` |
| API | `ListenAndServe(Config) error` |
| Timeouts | ReadTimeout, WriteTimeout, IdleTimeout (15s/15s/60s defaults in serve cmd) |
| Shutdown | `signal.NotifyContext` for SIGINT/SIGTERM, configurable ShutdownTimeout (30s) |
| Integration | `core/cli/serve.go` updated to use `server.ListenAndServe()` |
| Tests | 3 (server responds, graceful shutdown, config fields) |

## Design Notes

- `ListenAndServe` blocks until signal received, then calls `srv.Shutdown(ctx)`.
- Second signal resets to default behaviour (force kill).
- `.gitkeep` removed from `core/server/`.

## Test Results

All 31 packages pass (29 with tests). `go vet` clean.
