# Feature #36 — Health Checks: Changelog

## [Unreleased]

### Added
- `core/health/health.go` — `Routes(r, db)` registers `/health` (liveness) and `/health/ready` (readiness) endpoints.
- Liveness returns HTTP 200 `{"status":"ok"}`.
- Readiness pings database: HTTP 200 on success, HTTP 503 on failure with error detail.
- Health routes auto-registered in `RouterProvider.Boot()`.
