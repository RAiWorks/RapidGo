# Feature #37 — Graceful Shutdown: Plan

## Tasks

1. Create `core/server/server.go` with `Config` struct and `ListenAndServe()`.
2. Update `core/cli/serve.go` to use `server.ListenAndServe()`.
3. Remove `core/server/.gitkeep`.
4. Write tests: server starts and responds, shutdown on context cancel.
5. Run full regression + go vet.
6. Commit, merge to main, push.

## Test Plan

| TC | Description | Expected |
|----|-------------|----------|
| 01 | Server starts and responds to HTTP request | 200 OK |
| 02 | Server shuts down cleanly on context cancel | No error, exits |
| 03 | Config defaults populate correctly | Non-zero timeouts |
