# 📝 Changelog: CLI Foundation

> **Feature**: `10` — CLI Foundation
> **Branch**: `feature/10-cli-foundation`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

1. **Phase A — Dependencies**: Installed `github.com/spf13/cobra v1.10.2` (+ pflag, mousetrap). Build clean.
2. **Phase B — Root Command**: Created `core/cli/root.go` — `rootCmd`, `Execute()`, `NewApp()`, `Version` constant. Build clean.
3. **Phase C — Serve + Version + main.go**: Created `core/cli/serve.go` (serve command with `--port` flag), `core/cli/version.go` (version command). Refactored `cmd/main.go` to thin shell calling `cli.Execute()`. Build + vet clean.
4. **Phase D — Testing**: Created `core/cli/cli_test.go` (7 tests). All 123 tests pass, zero vet warnings.
5. **Phase E — Changelog + self-review**: This entry. All code reviewed — clean, idiomatic Go.

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| `app.Application` → `app.App` | Architecture doc used `app.Application` | Actual type is `app.App` | Architecture doc assumed wrong type name; code adapted to real codebase |
| `fmt.Printf` → `fmt.Fprintf(cmd.OutOrStdout(), ...)` in version cmd | Architecture used `fmt.Printf` | Used `cmd.OutOrStdout()` | Enables testability — `SetOut()` redirects output to buffer in tests |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
| Use `cmd.OutOrStdout()` for testable output | `fmt.Printf` writes to os.Stdout which can't be captured by Cobra's `SetOut()` in tests | 2026-03-06 |
