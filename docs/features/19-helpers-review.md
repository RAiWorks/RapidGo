# 📋 Review: Helpers

> **Feature**: `19` — Helpers
> **Branch**: `feature/19-helpers`
> **Merged**: 2026-03-06
> **Commit**: `cbac919` (merge), `12b4852` (impl)

---

## Summary

Feature #19 adds 17 stateless helper functions across 7 files in `app/helpers/`. Covers password hashing, secure random generation, string manipulation, number formatting, time formatting, struct/map utilities, and environment variable access.

## Files Changed

| File | Type | Description |
|---|---|---|
| `app/helpers/password.go` | Created | `HashPassword`, `CheckPassword` (bcrypt) |
| `app/helpers/random.go` | Created | `RandomString` (crypto/rand hex) |
| `app/helpers/string.go` | Created | `Slugify`, `Truncate`, `Contains`, `Title`, `Excerpt`, `StripHTML`, `Mask` |
| `app/helpers/number.go` | Created | `FormatBytes`, `Clamp` |
| `app/helpers/time.go` | Created | `TimeAgo`, `FormatDate` |
| `app/helpers/data.go` | Created | `StructToMap`, `MapKeys` |
| `app/helpers/env.go` | Created | `Env` |
| `app/helpers/helpers_test.go` | Created | 26 tests covering all 17 functions |
| `go.mod` / `go.sum` | Modified | `golang.org/x/crypto` promoted to direct dependency |

## Dependencies Added

| Package | Version | Purpose |
|---|---|---|
| `golang.org/x/crypto` | v0.48.0 | bcrypt for password hashing (promoted from indirect) |

## Test Results

- **New tests**: 26
- **Total tests**: 211 — all pass
- **`go vet`**: clean

## Architecture Compliance

All 17 functions from the architecture document are implemented. No deviations.

## Key Decisions

1. **`crypto/rand`** for `RandomString` — cryptographically secure, not `math/rand`
2. **bcrypt** for password hashing — industry standard, configurable cost via `golang.org/x/crypto`
3. **Pure functions** — no state, no side effects, easily testable

## Status: ✅ SHIPPED
