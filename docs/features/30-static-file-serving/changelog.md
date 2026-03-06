# Feature #30 — Static File Serving: Changelog

## [Unreleased]

### Status
Feature already fully implemented in Features #07 (Router) and #17 (Views & Templates).

### Existing (no changes)
- `core/router/router.go` — `Static()` and `StaticFile()` methods.
- `app/providers/router_provider.go` — `/static` and `/uploads` routes.
- `core/router/view_test.go` — TC-04, TC-05 covering both methods.

### Added
- Feature #30 documentation confirming blueprint coverage.

### Deviation log
| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | `StaticFile("/favicon.ico", ...)` included | Omitted | No favicon.ico file exists yet; intentional per #17 architecture decision |
