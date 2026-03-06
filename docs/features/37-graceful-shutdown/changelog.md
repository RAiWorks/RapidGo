# Feature #37 — Graceful Shutdown: Changelog

## [Unreleased]

### Added
- `core/server/server.go` — `Config` struct and `ListenAndServe()` function.
- Configurable Read/Write/Idle/Shutdown timeouts.
- Signal-based graceful shutdown (SIGINT, SIGTERM).

### Changed
- `core/cli/serve.go` — replaced `r.Run()` with `server.ListenAndServe()`.
