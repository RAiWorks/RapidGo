# 🪞 Review: Project Setup & Structure

> **Feature**: `01` — Project Setup & Structure
> **Branch**: `feature/01-project-setup`
> **Shipped**: 2026-03-05

---

## What Went Well

- The Mastery workflow (Discuss → Design → Plan → Build → Ship → Reflect) worked exactly as intended — zero wasted effort, no backtracking
- Having the full architecture doc with exact file contents meant the build phase was pure execution — no design decisions during implementation
- Pre-creating all 44 directories gives a strong sense of the framework's shape from day one
- Test plan was comprehensive and caught the `make` unavailability on Windows early

## What Went Wrong

- Architecture doc's directory count table had `core/` as 16 when it's actually 17 (1 parent + 16 packages), making the total 44 not 43 — caught and fixed during build
- `make` is not available natively on Windows PowerShell — not a blocker (commands verified manually) but worth noting for future CI/CD setup

## What Was Learned

- Go `go mod init` auto-sets the `go` version to the installed toolchain version (1.25.6), not a minimum — must manually edit `go.mod` for compatibility targets
- PowerShell string replacement with emojis/unicode requires careful handling — regex replacements sometimes don't work cleanly
- `.gitkeep` count (35) is less than directory count (44) because some directories get real Go files instead

## What to Change Next Time

- Double-check table math in architecture docs before finalizing — the count discrepancy could have been caught in the design audit
- Add a "Windows compatibility notes" section to architecture docs for features that produce cross-platform tooling (Makefiles, shell scripts)

---

## Metrics

| Metric | Value |
|---|---|
| Files created | 49 (35 .gitkeep + 4 Go + 4 config + 6 doc updates) |
| Lines of code | 12 (Go), 66 (.env), 29 (Makefile), 32 (.gitignore) |
| Tests passed | 15/15 |
| Deviations from plan | 3 (dir count, make on Windows, .gitignore reference/ rule kept) |

---

## Roadmap Updated

- [x] Feature #01 marked ✅ in `project-roadmap.md`
