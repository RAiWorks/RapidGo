# Feature #36 — Health Checks: Plan

## Tasks

1. Create `core/health/health.go` with `Routes(r, db)`.
2. Register health routes in `RouterProvider.Boot()`.
3. Write tests (liveness 200, readiness 200 with working DB, readiness 503 with broken DB).
4. Run full regression + go vet.
5. Commit, merge to main, push.

## Test Plan

| TC | Description | Expected |
|----|-------------|----------|
| 01 | GET /health | 200, `{"status":"ok"}` |
| 02 | GET /health/ready with live DB | 200, `{"status":"ready","db":"connected"}` |
| 03 | GET /health/ready with closed DB | 503, status "error" |

## Estimated Complexity

Low — two handlers, no new dependencies.
