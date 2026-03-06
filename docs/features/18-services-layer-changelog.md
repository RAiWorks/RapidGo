# 📝 Changelog: Services Layer

> **Feature**: `18` — Services Layer
> **Branch**: `feature/18-services-layer`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

### Phase A — UserService
- Created `app/services/user_service.go` — `UserService`, `NewUserService`, `GetByID`, `Create`, `Update`, `Delete`
- `go build` clean, `go vet` clean

### Phase B — Testing
- Created `app/services/user_service_test.go` — 8 test cases (TC-01 through TC-08)
- `go test ./app/services/... -v` — 8/8 pass
- `go test ./... -count=1` — 185 total tests, 0 failures

### Phase C — Changelog & Cross-Check
- Code vs architecture: exact match, 0 deviations
- Changelog updated

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| None | — | — | — |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| None | Implementation matched architecture exactly | 2026-03-06 |
