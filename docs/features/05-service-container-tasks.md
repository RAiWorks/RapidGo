# ✅ Tasks: Service Container

> **Feature**: `05` — Service Container
> **Architecture**: [`05-service-container-architecture.md`](05-service-container-architecture.md)
> **Branch**: `feature/05-service-container`
> **Status**: 🔴 NOT STARTED
> **Progress**: 0/18 tasks complete

---

## Pre-Flight Checklist

- [x] Discussion doc is marked COMPLETE
- [x] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [x] Dependent features are merged to `main`
- [x] Test plan doc created
- [x] Changelog doc created (empty)

---

## Phase A — Container Core

> Container struct, registration methods, resolution methods.

- [ ] **A.1** — Create `core/container/container.go` with package declaration, imports, `Factory` type, `Container` struct
- [ ] **A.2** — Implement `New()` constructor
- [ ] **A.3** — Implement `Bind()` — transient factory registration
- [ ] **A.4** — Implement `Singleton()` — shared instance registration with lazy init
- [ ] **A.5** — Implement `Instance()` — pre-created object registration
- [ ] **A.6** — Implement `Make()` — resolve by name (instances → bindings → panic)
- [ ] **A.7** — Implement `MustMake[T]()` — generic typed resolution
- [ ] **A.8** — Implement `Has()` — existence check
- [ ] 📍 **Checkpoint A** — Container compiles, `go vet` clean

---

## Phase B — Provider Interface

> Provider contract for two-phase service lifecycle.

- [ ] **B.1** — Create `core/container/provider.go` with `Provider` interface (`Register`, `Boot`)
- [ ] 📍 **Checkpoint B** — Provider interface compiles, `go vet` clean

---

## Phase C — App Bootstrap

> App struct for provider management and boot sequence.

- [ ] **C.1** — Create `core/app/app.go` with `App` struct, `New()`, `Register()`, `Boot()`, `Make()`
- [ ] 📍 **Checkpoint C** — App compiles, `go vet` clean

---

## Phase D — Testing

> Execute the test plan, verify all acceptance criteria.

- [ ] **D.1** — Create `core/container/container_test.go` with container test cases
- [ ] **D.2** — Create `core/app/app_test.go` with app bootstrap test cases
- [ ] **D.3** — Run `go test ./core/container/... ./core/app/...` — all tests pass
- [ ] **D.4** — Run `go vet ./...` — no issues
- [ ] 📍 **Checkpoint D** — All test cases pass, zero vet warnings

---

## Phase E — Documentation & Cleanup

> Changelog, roadmap, self-review.

- [ ] **E.1** — Update changelog doc with implementation summary
- [ ] **E.2** — Self-review all diffs — code is clean, idiomatic Go
- [ ] 📍 **Checkpoint E** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Update project roadmap progress
- [ ] Create review doc → `05-service-container-review.md`
