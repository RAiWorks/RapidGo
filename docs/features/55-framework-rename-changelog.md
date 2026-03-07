# 📝 Changelog: Framework Rename / Rebrand

> **Feature**: `55` — Framework Rename / Rebrand
> **Branch**: `feature/55-framework-rename`
> **Started**: 2026-03-07
> **Completed**: —

---

## Log

<!-- Add entries as you work. Most recent first. -->

- **2026-03-07** — All phases complete (A–F). Full test suite passes (30 packages, 0 failures). CLI outputs `RapidGo Framework v0.2.0`. No stale references in Go source or docs.
- **2026-03-07** — Phase E: Full test suite run. Found `cli_test.go:44` still checking for `"RGo"` — fixed to `"RapidGo"`. All 30 packages pass on re-run.
- **2026-03-07** — Phase D: Mass-replaced all `.md` docs. 49 files updated. Zero stale refs outside rename docs.
- **2026-03-07** — Phase C: Updated `resources/views/home.html` — heading, GitHub link, footer.
- **2026-03-07** — Phase B: Updated `.env` (5 refs), `Makefile` (binary target), `Caddyfile` (comment). Created `LICENSE` (MIT). Ran `go mod tidy`. Dockerfile had no RGo refs — no change needed. `.gitignore` uses `bin/` — no change needed (B.5 NO-OP confirmed).
- **2026-03-07** — Phase A: Updated `go.mod` module path. Mass-replaced imports in 37 `.go` files. Updated CLI strings in `root.go`, `version.go`, `serve.go`. Updated `home_controller.go`, `database/connection.go`, test assertions in 4 test files. `go build ./...` passes with zero errors.
- **2026-03-07** — Created `feature/55-framework-rename` branch from `main`.

---

## Deviations from Plan

<!-- Things that went differently than the architecture/tasks docs planned -->

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| `cli_test.go` assertion | Covered by A.7 | Missed in initial pass — test used `"RGo"` in a string slice, not a standalone literal | Caught by Phase E test run, fixed immediately |
| Dockerfile (B.3) | Update any comment referencing RGo | No RGo references found — no changes needed | Binary is named `server`, no framework-specific comments |
| `.gitignore` (B.5) | Update if binary name referenced | NO-OP confirmed — uses `bin/` directory pattern | Already identified in cross-check |

## Key Decisions Made During Build

<!-- Runtime decisions NOT in the discussion/architecture docs -->

| Decision | Context | Date |
|---|---|---|
| `cli_test.go` TC-03 assertion updated | Test checked for `"RGo"` in help output — changed to `"RapidGo"` | 2026-03-07 |
