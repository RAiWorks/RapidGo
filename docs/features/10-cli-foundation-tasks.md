# ✅ Tasks: CLI Foundation

> **Feature**: `10` — CLI Foundation
> **Architecture**: [`10-cli-foundation-architecture.md`](10-cli-foundation-architecture.md)
> **Branch**: `feature/10-cli-foundation`
> **Status**: � COMPLETE
> **Progress**: 12/12 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [x] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Dependencies

> Install Cobra CLI library.

- [x] **A.1** — Run `go get github.com/spf13/cobra`
- [x] 📍 **Checkpoint A** — `go build ./...` succeeds, `go.mod` lists `cobra`

---

## Phase B — Root Command

> Create the CLI package with root command and app bootstrap helper.

- [x] **B.1** — Create `core/cli/root.go`: `rootCmd` definition, `Execute()`, `NewApp()`, `Version` constant
- [x] 📍 **Checkpoint B** — `go build ./...` succeeds (root.go compiles)

---

## Phase C — Serve Command

> Move server startup logic from main.go into a `serve` subcommand.

- [x] **C.1** — Create `core/cli/serve.go`: `serveCmd` with `--port` flag, banner, server start
- [x] **C.2** — Create `core/cli/version.go`: `versionCmd` prints framework version
- [x] **C.3** — Refactor `cmd/main.go` — replace all logic with `cli.Execute()` call
- [x] 📍 **Checkpoint C** — `go build ./...` succeeds, `rgo serve` starts server, `rgo version` prints version

---

## Phase D — Testing

> Comprehensive test suite for CLI commands.

- [x] **D.1** — Create `core/cli/cli_test.go` with root, serve, and version tests
- [x] **D.2** — Run `go test ./core/cli/...` — all tests pass
- [x] **D.3** — Run `go test ./...` + `go vet ./...` — full regression, no failures
- [x] 📍 **Checkpoint D** — All tests pass, zero vet warnings

---

## Phase E — Documentation & Cleanup

> Changelog, self-review.

- [x] **E.1** — Update changelog doc with implementation summary
- [x] **E.2** — Self-review all diffs — code is clean, idiomatic Go
- [x] 📍 **Checkpoint E** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Update project roadmap progress
- [ ] Create review doc → `10-cli-foundation-review.md`
