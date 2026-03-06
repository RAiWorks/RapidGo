# Feature #37 — Graceful Shutdown: Discussion

## Overview

Replace the current `r.Run()` call with a proper `http.Server` that supports configurable timeouts, signal handling (SIGINT/SIGTERM), and graceful shutdown with a deadline.

## Blueprint Reference

- `http.Server` with ReadTimeout (15s), WriteTimeout (15s), IdleTimeout (60s).
- Starts server in goroutine, waits for SIGINT/SIGTERM.
- `srv.Shutdown(ctx)` with 30s deadline.
- Closes DB connection after shutdown.

## Current State

`core/cli/serve.go` calls `r.Run(":" + port)` — Gin's default, no graceful shutdown. `core/server/` has only `.gitkeep`.

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Package | `core/server/` | Blueprint convention, replaces `.gitkeep` |
| Timeout source | Env vars with defaults | `SERVER_READ_TIMEOUT`, `SERVER_WRITE_TIMEOUT`, `SERVER_IDLE_TIMEOUT`, `SERVER_SHUTDOWN_TIMEOUT` |
| Signal handling | `signal.NotifyContext` | Cleaner than manual channel pattern since Go 1.16+ |
| DB close | Via container "db" if bound | Gracefully close DB pool on shutdown |
