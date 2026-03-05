# ✅ Tasks: Error Handling

> **Feature**: `04` — Error Handling
> **Architecture**: [`04-error-handling-architecture.md`](04-error-handling-architecture.md)
> **Branch**: `feature/04-error-handling`
> **Status**: � COMPLETE
> **Progress**: 14/14 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [x] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Core Error Types

> AppError struct, error/Unwrap interface compliance.

- [x] **A.1** — Create `core/errors/errors.go` with package declaration and imports
- [x] **A.2** — Define `AppError` struct with `Code int`, `Message string`, `Err error` fields
- [x] **A.3** — Implement `Error() string` method returning `Message`
- [x] **A.4** — Implement `Unwrap() error` method returning `Err` (handles nil)
- [x] 📍 **Checkpoint A** — `AppError` compiles, implements `error` interface, `go vet` clean

---

## Phase B — Constructor Helpers

> Helper functions for common HTTP error scenarios.

- [x] **B.1** — Implement `NotFound(message string) *AppError` — 404
- [x] **B.2** — Implement `BadRequest(message string) *AppError` — 400
- [x] **B.3** — Implement `Internal(err error) *AppError` — 500 with wrapped error
- [x] **B.4** — Implement `Unauthorized(message string) *AppError` — 401
- [x] **B.5** — Implement `Forbidden(message string) *AppError` — 403
- [x] **B.6** — Implement `Conflict(message string) *AppError` — 409
- [x] **B.7** — Implement `Unprocessable(message string) *AppError` — 422
- [x] 📍 **Checkpoint B** — All 7 constructors compile, return correct codes, `go vet` clean

---

## Phase C — Error Response Helper

> Debug-aware response formatting.

- [x] **C.1** — Implement `ErrorResponse() map[string]any` on `*AppError` — returns `{"success": false, "error": message}`, adds `"internal"` key only when `config.IsDebug() && e.Err != nil`
- [x] 📍 **Checkpoint C** — `ErrorResponse()` compiles, `go vet` clean

---

## Phase D — Testing

> Execute the test plan, verify all acceptance criteria.

- [x] **D.1** — Create `core/errors/errors_test.go` with all test cases from test plan
- [x] **D.2** — Run `go test ./core/errors/...` — all tests pass
- [x] **D.3** — Run `go vet ./...` — no issues
- [x] 📍 **Checkpoint D** — All test cases pass, zero vet warnings

---

## Phase E — Documentation & Cleanup

> Changelog, roadmap, self-review.

- [x] **E.1** — Update changelog doc with implementation summary
- [x] **E.2** — Self-review all diffs — code is clean, idiomatic Go
- [x] 📍 **Checkpoint E** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [x] All phases complete
- [x] All checkpoints verified
- [x] Final commit with descriptive message
- [x] Merge to `main`
- [x] Push `main`
- [x] **Keep the feature branch** — do not delete
- [x] Update project roadmap progress
- [x] Create review doc → `04-error-handling-review.md`
