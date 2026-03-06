# ✅ Tasks: Models (GORM)

> **Feature**: `11` — Models (GORM)
> **Architecture**: [`11-models-architecture.md`](11-models-architecture.md)
> **Branch**: `feature/11-models`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/10 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — BaseModel

> Define the reusable base struct.

- [ ] **A.1** — Create `database/models/base.go`: `BaseModel` struct with ID, CreatedAt, UpdatedAt
- [ ] 📍 **Checkpoint A** — `go build ./database/models/...` succeeds

---

## Phase B — User & Post Models

> Define the concrete application models.

- [ ] **B.1** — Create `database/models/user.go`: `User` struct with GORM tags and Posts relationship
- [ ] **B.2** — Create `database/models/post.go`: `Post` struct with GORM tags and User relationship
- [ ] 📍 **Checkpoint B** — `go build ./database/...` succeeds, `go vet ./database/...` clean

---

## Phase C — Testing

> Verify model definitions with GORM/SQLite integration tests.

- [ ] **C.1** — Create `database/models/models_test.go` with struct, GORM, and relationship tests
- [ ] **C.2** — Run `go test ./database/models/...` — all tests pass
- [ ] **C.3** — Run `go test ./...` + `go vet ./...` — full regression, no failures
- [ ] 📍 **Checkpoint C** — All tests pass, zero vet warnings

---

## Phase D — Documentation & Cleanup

> Changelog, self-review.

- [ ] **D.1** — Update changelog doc with implementation summary
- [ ] **D.2** — Self-review all diffs — code is clean, idiomatic Go
- [ ] 📍 **Checkpoint D** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Update project roadmap progress
- [ ] Create review doc → `11-models-review.md`
