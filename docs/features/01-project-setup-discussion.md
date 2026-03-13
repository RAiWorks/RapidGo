# 💬 Discussion: Project Setup & Structure

> **Feature**: `01` — Project Setup & Structure
> **Status**: 🟢 COMPLETE
> **Branch**: `feature/01-project-setup`
> **Depends On**: — (no dependencies — this is the foundation)
> **Date Started**: 2026-03-05
> **Date Completed**: 2026-03-05

---

## Summary

Initialize the RapidGo Framework as a Go module, create the full directory structure defined in the blueprint, add the application entry point (`cmd/main.go`), and provide essential project files (`.env`, `Makefile`, `.gitignore`). This feature produces a compilable, runnable skeleton that prints a startup banner — the foundation every other feature builds on.

---

## Functional Requirements

- As a framework developer, I want a properly initialized Go module so that all subsequent features have a clean import path (`github.com/raiworks/rapidgo`)
- As a framework developer, I want the complete directory tree pre-created so that feature branches can add files without creating folders ad hoc
- As a framework developer, I want a minimal `cmd/main.go` that compiles and runs so that the build pipeline is proven from day one
- As a framework developer, I want a `.env` file with placeholder values so that the config system (Feature #02) has a file to load
- As a framework developer, I want a `Makefile` with standard targets so that common commands are consistent across environments
- As a framework developer, I want `.gitignore` rules so that build artifacts, binaries, and environment secrets are never committed

## Current State / Reference

### What Exists

Nothing — greenfield project. The repository currently contains only documentation (`docs/` folder with framework reference docs, mastery process, context, and roadmap).

### What Works Well

- Complete blueprint defines the exact directory structure (43 directories, 5 root files)
- All 62 framework reference docs already written — we know exactly what every package will contain
- Technology stack decided (Go 1.21+, Gin, GORM, Cobra, etc.)

### What Needs Improvement

N/A — no existing code to improve.

## Proposed Approach

1. **Initialize Go module** — `go mod init github.com/raiworks/rapidgo`
2. **Create full directory tree** — All directories from the blueprint, with `.gitkeep` in empty directories
3. **Create `cmd/main.go`** — Minimal `package main` with a startup banner print
4. **Create `.env`** — Placeholder environment variables covering all config keys the framework will need
5. **Create `Makefile`** — Targets: `build`, `run`, `test`, `clean`, `fmt`, `vet`, `lint`
6. **Update `.gitignore`** — Go binaries, `.env` overrides, IDE files, OS files, storage artifacts
7. **Verify** — `go build ./...` succeeds, `go run cmd/main.go` prints the banner

This is a pure scaffold — no business logic, no dependencies beyond the Go standard library. Third-party dependencies (`go.sum`) will be introduced by subsequent features as they `go get` their required packages.

## Edge Cases & Risks

- [x] **Module path casing** — GitHub repo is `raiworks/rapidgo` (mixed case). Go modules are case-sensitive; the module path must match the repository URL exactly: `github.com/raiworks/rapidgo`
- [x] **Empty directories in Git** — Git doesn't track empty directories. Use `.gitkeep` placeholder files in every leaf directory that has no Go files
- [x] **`.env` in version control** — The `.env` file with placeholder/default values IS committed (it serves as documentation of all config keys). A `.env.local` override (which contains real secrets) is gitignored
- [x] **No dependencies yet** — `go.sum` will be empty/absent until Feature #02+ adds packages. That's expected and correct
- [x] **`main.go` does nothing useful** — Intentionally minimal. The app lifecycle, config loading, and server startup are separate features (#02, #05, etc.)

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Go 1.21+ | Toolchain | ✅ Required on dev machine |
| Git | Toolchain | ✅ Required for version control |
| Make (optional) | Toolchain | ⚡ Nice to have — Makefile targets can be run manually |

No feature dependencies — this is Feature #01.

## Open Questions

All resolved during discussion:

- [x] **What should the module path be?** → `github.com/raiworks/rapidgo` (matches GitHub repo URL)
- [x] **Should we add all directories now or only what's needed?** → All directories. The blueprint defines the full structure and having it pre-created keeps feature branches focused on code, not folder creation
- [x] **Should `.env` be committed?** → Yes, with placeholder values. Gitignore `.env.local` for real secrets
- [x] **Should `main.go` import any framework packages?** → No. Just `fmt` for the startup banner. Framework wiring happens in later features
- [x] **What Go version constraint?** → `go 1.21` in `go.mod` (minimum for `log/slog` which is in our tech stack)

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-05 | Module path: `github.com/raiworks/rapidgo` | Matches GitHub repository URL, standard Go convention |
| 2026-03-05 | Create all blueprint directories upfront | Avoids directory creation scattered across 40+ features; keeps feature branches clean |
| 2026-03-05 | Use `.gitkeep` for empty directories | Git doesn't track empty dirs; `.gitkeep` is the standard convention |
| 2026-03-05 | Commit `.env` with placeholders, gitignore `.env.local` | `.env` serves as config documentation; secrets go in `.env.local` |
| 2026-03-05 | `go 1.21` minimum | Required for `log/slog` (standard library structured logging) |
| 2026-03-05 | Minimal `main.go` — no framework imports | Framework wiring is Feature #05/#06 (Service Container & Providers) |

## Discussion Complete ✅

**Summary**: Feature #01 creates the Go module, full directory structure, minimal entry point, and project files — a compilable skeleton with zero third-party dependencies.
**Completed**: 2026-03-05
**Next**: Create architecture doc → `01-project-setup-architecture.md`
