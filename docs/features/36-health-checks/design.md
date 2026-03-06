# Feature #36 — Health Checks: Design

## Package

`core/health/health.go`

## Public API

```go
// Routes registers /health and /health/ready on the given router.
func Routes(r *router.Router, db *gorm.DB)
```

### GET /health (Liveness)

- Always returns HTTP 200:
  ```json
  {"status": "ok"}
  ```

### GET /health/ready (Readiness)

- Calls `db.DB()` → `sqlDB.Ping()`.
- On success → HTTP 200:
  ```json
  {"status": "ready", "db": "connected"}
  ```
- On failure → HTTP 503:
  ```json
  {"status": "error", "db": "<error message>"}
  ```

## Integration

In `RouterProvider.Boot()`, after route registration:

```go
db := container.MustMake[*gorm.DB](c, "db")
health.Routes(r, db)
```

## File Layout

```
core/health/
  health.go       — Routes function with two handlers
  health_test.go   — HTTP tests for both endpoints
```
