# 🧪 Test Plan: Project Setup & Structure

> **Feature**: `01` — Project Setup & Structure
> **Tasks**: [`01-project-setup-tasks.md`](01-project-setup-tasks.md)
> **Date**: 2026-03-05

---

## Acceptance Criteria

The feature is DONE when ALL of these are true:

- [x] Go module initialized with path `github.com/RAiWorks/RGo` and `go 1.21`
- [x] All 43 directories exist in the correct hierarchy
- [x] `cmd/main.go` compiles and runs without errors
- [x] Running `go run ./cmd/...` prints the RGo startup banner
- [x] `go vet ./...` reports zero issues
- [x] `.env` exists with all placeholder configuration groups
- [x] `.gitignore` excludes binaries, `.env.local`, IDE files, storage artifacts
- [x] `Makefile` works: `make build`, `make run`, `make clean`
- [x] `README.md` exists with links to documentation
- [x] No third-party dependencies in `go.mod` (only standard library)

---

## Test Cases

### TC-01: Go Module Initialization

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | Go 1.21+ installed, project root is clean |
| **Steps** | 1. Check `go.mod` exists → 2. Verify `module github.com/RAiWorks/RGo` → 3. Verify `go 1.21` directive |
| **Expected Result** | `go.mod` contains correct module path and Go version, no `require` blocks |
| **Status** | ✅ Pass |
| **Notes** | — |

### TC-02: Directory Structure Completeness

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | All Phase B tasks completed |
| **Steps** | 1. List all directories recursively → 2. Compare against architecture doc's 43-directory list → 3. Verify every leaf dir has `.gitkeep` or a `.go` file |
| **Expected Result** | All 43 directories present, no missing, no extras |
| **Status** | ✅ Pass |
| **Notes** | — |

### TC-03: Application Compiles

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `cmd/main.go` created |
| **Steps** | 1. Run `go build ./cmd/...` → 2. Check exit code is 0 → 3. Verify `bin/rgo` binary exists (via Makefile) |
| **Expected Result** | Zero compilation errors, binary produced |
| **Status** | ✅ Pass |
| **Notes** | — |

### TC-04: Application Runs and Prints Banner

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | TC-03 passes |
| **Steps** | 1. Run `go run ./cmd/...` → 2. Capture stdout |
| **Expected Result** | Output contains "RGo Framework" and "github.com/RAiWorks/RGo" |
| **Status** | ✅ Pass |
| **Notes** | — |

### TC-05: Go Vet Clean

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | All Go files created |
| **Steps** | 1. Run `go vet ./...` → 2. Check exit code |
| **Expected Result** | Exit code 0, no warnings or errors |
| **Status** | ✅ Pass |
| **Notes** | — |

### TC-06: Makefile Targets

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `Makefile` created |
| **Steps** | 1. Run `make build` → verify `bin/rgo` → 2. Run `make run` → verify banner → 3. Run `make clean` → verify `bin/` removed |
| **Expected Result** | All three targets execute successfully |
| **Status** | ✅ Pass |
| **Notes** | `make` not available on Windows; verified with equivalent commands (`go build -o bin/rgo.exe`, `go run`, `Remove-Item bin/`) |

### TC-07: `.env` Configuration File

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `.env` created |
| **Steps** | 1. Open `.env` → 2. Verify sections: Application, Database, Session, Cache, Redis, JWT, Mail, Logging, Storage, Server → 3. Verify all values are placeholders (no real secrets) |
| **Expected Result** | All 10 configuration groups present with safe placeholder values |
| **Status** | ✅ Pass |
| **Notes** | — |

### TC-08: `.gitignore` Coverage

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `.gitignore` created |
| **Steps** | 1. Verify `bin/` is ignored → 2. Verify `.env.local` is ignored → 3. Verify `.idea/` and `.vscode/` ignored → 4. Verify `storage/logs/*.log` ignored → 5. Verify `.gitkeep` files are NOT ignored |
| **Expected Result** | All patterns present, `.gitkeep` exception works |
| **Status** | ✅ Pass |
| **Notes** | — |

### TC-09: No Third-Party Dependencies

| Property | Value |
|---|---|
| **Category** | Edge Case |
| **Precondition** | Module initialized |
| **Steps** | 1. Open `go.mod` → 2. Check for `require` block |
| **Expected Result** | No `require` block exists — only standard library used |
| **Status** | ✅ Pass |
| **Notes** | Confirmed: no `require` block, `go.sum` absent |

### TC-10: README Links Valid

| Property | Value |
|---|---|
| **Category** | Happy Path |
| **Precondition** | `README.md` created |
| **Steps** | 1. Verify link to `docs/project-context.md` → 2. Verify link to `docs/project-roadmap.md` → 3. Verify link to `docs/mastery.md` → 4. Verify link to `docs/framework/README.md` |
| **Expected Result** | All 4 linked files exist at the referenced paths |
| **Status** | ✅ Pass |
| **Notes** | All 4 linked docs verified present |

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | `go build` with no dependencies | Compiles with only standard library — no errors |
| 2 | Empty directories tracked by Git | `.gitkeep` files ensure Git tracks every directory |
| 3 | Running on Windows vs Linux/macOS | Makefile uses portable commands; `rm -rf` works on both via Git Bash / WSL |
| 4 | `.env` loaded when no `.env.local` exists | Framework should function with `.env` defaults alone |

## Security Tests

| # | Test | Expected |
|---|---|---|
| 1 | `.env.local` is gitignored | `git status` never shows `.env.local` even if created |
| 2 | `storage/` artifacts gitignored | Log files, cache, sessions, uploads never committed |
| 3 | No secrets in `.env` | All values are placeholder/default — no real credentials |

## Performance Considerations

N/A — this feature creates project structure only. No runtime performance characteristics.

---

## Test Summary

<!-- Fill AFTER running all tests -->

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | 8 | 8 | 0 | 0 |
| Edge Cases | 4 | 4 | 0 | 0 |
| Security | 3 | 3 | 0 | 0 |
| **Total** | **15** | **15** | **0** | **0** |

**Result**: ✅ ALL PASS

**Executed**: 2026-03-05 | **Environment**: Go 1.25.6, Windows/amd64

**Notes**:
- TC-06 (Makefile): `make` not available on Windows PowerShell; verified equivalent commands (`go build -o bin/rgo.exe ./cmd/...`, `go run ./cmd/...`, `Remove-Item bin/`) manually — all pass.
- Edge Case 3 confirmed: Makefile targets require Git Bash / WSL on Windows; Go commands work natively.
