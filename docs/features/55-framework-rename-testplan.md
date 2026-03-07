# 🧪 Test Plan: Framework Rename / Rebrand

> **Feature**: `55` — Framework Rename / Rebrand
> **Tasks**: [`55-framework-rename-tasks.md`](55-framework-rename-tasks.md)
> **Date**: 2026-03-07

---

## Acceptance Criteria

- [ ] Every `.go` file compiles with the new module path `github.com/RAiWorks/RapidGo`
- [ ] All existing tests pass without modification (beyond the rename itself)
- [ ] CLI binary is named `rapidgo` and outputs `RapidGo Framework v0.2.0`
- [ ] Welcome page shows "RapidGo" branding with correct GitHub link
- [ ] No stale "RGo" references remain in Go source files (except `RGO_TEST_*` env vars in helpers_test.go)
- [ ] No stale "RGo" references remain in documentation (except historical context where appropriate)
- [ ] `LICENSE` file exists with MIT license and RAi Works copyright
- [ ] `.env` defaults use `RapidGo` and `rapidgo_` prefixes

---

## Test Cases

### TC-01: Compilation

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | All Go files have updated imports |
| **Steps** | 1. Run `go build ./...` |
| **Expected Result** | Zero compilation errors |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-02: Full Test Suite

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | TC-01 passes |
| **Steps** | 1. Run `go test ./... -count=1` |
| **Expected Result** | All tests pass across all packages |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-03: CLI Version Output

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Binary built as `bin/rapidgo` |
| **Steps** | 1. Run `./bin/rapidgo version` |
| **Expected Result** | Output: `RapidGo Framework v0.2.0` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-04: CLI Serve Banner

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Binary built, `.env` configured |
| **Steps** | 1. Run `./bin/rapidgo serve` → 2. Read startup banner |
| **Expected Result** | Banner shows `RapidGo Framework` and `github.com/RAiWorks/RapidGo` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-05: Welcome Page Branding

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Server running on port 8080 |
| **Steps** | 1. Visit `http://localhost:8080` → 2. Check page title → 3. Check heading → 4. Check GitHub link → 5. Check footer |
| **Expected Result** | All show "RapidGo", GitHub link points to `github.com/RAiWorks/RapidGo` |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-06: No Stale Go References

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | All replacements applied |
| **Steps** | 1. Run `grep -r "github.com/RAiWorks/RGo[^a-zA-Z]" --include="*.go"` → 2. Run `grep -r '"rgo"' --include="*.go"` |
| **Expected Result** | Zero matches for module path. Only `RGO_TEST_VAR`/`RGO_TEST_MISSING` in `helpers_test.go` for the second grep |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-07: No Stale Doc References

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | All doc replacements applied |
| **Steps** | 1. Run `grep -r "RGo" docs/ \| grep -v RapidGo \| grep -v "RGO_"` |
| **Expected Result** | Zero matches |
| **Status** | ⬜ Not Run |
| **Notes** | Some feature docs may reference "RGo" in historical context — acceptable if alongside "RapidGo" |

### TC-08: LICENSE File Exists

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | LICENSE file created |
| **Steps** | 1. Verify `LICENSE` file exists at project root → 2. Verify contains "MIT License" → 3. Verify contains "RAi Works" → 4. Verify contains "(https://rai.works)" |
| **Expected Result** | All four checks pass |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-09: Makefile Binary Target

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Makefile updated |
| **Steps** | 1. Run `make build` → 2. Check `bin/rapidgo` exists → 3. Run `bin/rapidgo version` |
| **Expected Result** | Binary created at `bin/rapidgo`, outputs version string |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-10: Environment Config Defaults

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `.env` updated |
| **Steps** | 1. Verify `APP_NAME=RapidGo` → 2. Verify `DB_NAME=rapidgo_dev` → 3. Verify `CACHE_PREFIX=rapidgo_` → 4. Verify `MAIL_FROM_NAME=RapidGo` |
| **Expected Result** | All four defaults use new name |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-11: Go Module Path

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `go.mod` updated |
| **Steps** | 1. Read `go.mod` first line → 2. Run `go mod tidy` → 3. Verify no errors |
| **Expected Result** | Module path is `github.com/RAiWorks/RapidGo`, tidy succeeds |
| **Status** | ⬜ Not Run |
| **Notes** | — |

### TC-12: Docker Build

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Dockerfile references correct module |
| **Steps** | 1. Run `docker build .` (if Docker available) |
| **Expected Result** | Build succeeds |
| **Status** | ⬜ Not Run |
| **Notes** | Optional — depends on Docker availability on dev machine |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | `RGO_TEST_VAR` env vars in helpers_test.go | NOT renamed — these are test-specific env var names, not framework branding |
| 2 | Old commits reference "RGo" in messages | Expected — no history rewriting |
| 3 | Feature doc filenames (e.g., `01-project-setup-discussion.md`) | NOT renamed — filenames stay, only content updated |
| 4 | `go.sum` has old module references | Regenerated by `go mod tidy` — old references removed automatically |
| 5 | Docker binary named `server` not `rapidgo` | Stays as `server` — Dockerfile uses generic name, no change needed |

---

## Test Summary

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 10 | — | — | — |
| Edge Cases | 2 | — | — | — |
| **Total** | **12** | — | — | — |

**Result**: ⬜ NOT RUN
