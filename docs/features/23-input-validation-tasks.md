# 📋 Tasks: Input Validation

> **Feature**: `23` — Input Validation
> **Branch**: `feature/23-input-validation`
> **Status**: NOT STARTED

---

## Phase 1 — Implementation

- [ ] Create `core/validation/validation.go` with Errors type, Validator struct, 9 methods
- [ ] Remove `core/validation/.gitkeep` (replaced by real file)
- [ ] Verify: `go build ./core/validation/...` compiles

**Checkpoint**: All types and methods compile with stdlib-only imports.

## Phase 2 — Tests

- [ ] Create `core/validation/validation_test.go`
  - [ ] TC-01: Required fails on empty string
  - [ ] TC-02: Required fails on whitespace-only
  - [ ] TC-03: Required passes on non-empty
  - [ ] TC-04: MinLength fails below minimum
  - [ ] TC-05: MinLength passes at minimum
  - [ ] TC-06: MaxLength fails above maximum
  - [ ] TC-07: MaxLength passes at maximum
  - [ ] TC-08: Email fails on invalid
  - [ ] TC-09: Email passes on valid
  - [ ] TC-10: URL fails on non-URL
  - [ ] TC-11: URL passes on valid URL
  - [ ] TC-12: Matches fails on non-match
  - [ ] TC-13: Matches passes on match
  - [ ] TC-14: In fails on disallowed value
  - [ ] TC-15: In passes on allowed value
  - [ ] TC-16: Confirmed fails on mismatch
  - [ ] TC-17: Confirmed passes on match
  - [ ] TC-18: IP fails on invalid IP
  - [ ] TC-19: IP passes on valid IP
  - [ ] TC-20: Errors.Add and Errors.First
  - [ ] TC-21: Chaining multiple rules
  - [ ] TC-22: Valid returns true when no errors
- [ ] Run full `go test ./... -count=1` — all pass

**Checkpoint**: All 22 tests pass. Full regression green.

## Phase 3 — Finalize

- [ ] Update changelog
- [ ] Run `go vet ./...` — clean
- [ ] Commit and push
