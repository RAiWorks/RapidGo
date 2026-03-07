# 📋 Tasks: Two-Factor Authentication (TOTP)

> **Feature**: `50` — Two-Factor Authentication (TOTP)
> **Architecture**: [`50-totp-2fa-architecture.md`](50-totp-2fa-architecture.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Phase A — TOTP Core Package

| # | Task | Detail |
|---|---|---|
| A1 | Add `pquerna/otp` dependency | `go get github.com/pquerna/otp` |
| A2 | Create `core/totp/totp.go` | `Key` struct, `GenerateKey()`, `ValidateCode()`, `GenerateBackupCodes()`, `HashBackupCode()`, `VerifyBackupCode()` |
| A3 | Create `core/totp/totp_test.go` | Tests for all 5 public functions |

**Exit**: `core/totp` package compiles, all tests pass.

---

## Phase B — User Model & Migration

| # | Task | Detail |
|---|---|---|
| B1 | Add TOTP fields to `User` model | `TOTPEnabled`, `TOTPSecret`, `TOTPVerifiedAt`, `BackupCodesHash` in `database/models/user.go` |
| B2 | Create migration file | `database/migrations/20260308000002_add_totp_fields.go` — adds 4 columns to `users` table |

**Exit**: User model has TOTP fields, migration compiles, existing model tests pass.

---

## Phase C — Verification

| # | Task | Detail |
|---|---|---|
| C1 | Run all tests | `go test ./...` — all packages pass |
| C2 | Verify build | `go build -o bin/rapidgo.exe ./cmd` — clean compile |

**Exit**: All packages green, binary builds.

---

## Next

Test plan → `50-totp-2fa-testplan.md`
