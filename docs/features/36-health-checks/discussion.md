# Feature #36 — Health Checks: Discussion

## Overview

Health check endpoints for Docker/Kubernetes liveness and readiness probes.

## Blueprint Reference

Two standard endpoints:

- **`GET /health`** — Liveness probe. Returns `{"status":"ok"}` (HTTP 200). Proves the process is alive.
- **`GET /health/ready`** — Readiness probe. Pings the database. Returns `{"status":"ready","db":"connected"}` (HTTP 200) or `{"status":"error","db":"<message>"}` (HTTP 503).

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Package location | `core/health/` | Aligns with other core modules |
| DB dependency | `*gorm.DB` injected | Matches container pattern (`"db"` binding) |
| Registration | `health.Routes(router, db)` called in `RouterProvider.Boot()` | Simple, centralised |
| Extensibility | Checker interface for future probes | Not needed now — YAGNI |

## Dependencies

- `#07` — Router (Gin wrapper)
- `#09` — Database (GORM `*gorm.DB`)
