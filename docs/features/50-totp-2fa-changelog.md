# 📝 Changelog: Two-Factor Authentication (TOTP)

> **Feature**: `50` — Two-Factor Authentication (TOTP)
> **Status**: ✅ SHIPPED
> **Date**: 2026-03-07
> **Commit**: `a28379c`

---

## Added

- `core/totp/totp.go` — TOTP package: `GenerateKey()`, `ValidateCode()`, `GenerateBackupCodes()`, `HashBackupCode()`, `VerifyBackupCode()`
- `core/totp/totp_test.go` — 17 unit tests for TOTP operations
- `database/migrations/20260308000002_add_totp_fields.go` — migration adding `totp_enabled`, `totp_secret`, `totp_verified_at`, `backup_codes_hash` to `users` table
- `go.mod` / `go.sum` — `github.com/pquerna/otp` dependency

## Changed

- `database/models/user.go` — User model gains `TOTPEnabled`, `TOTPSecret`, `TOTPVerifiedAt`, `BackupCodesHash` fields

## Files

| File | Action |
|---|---|
| `core/totp/totp.go` | NEW |
| `core/totp/totp_test.go` | NEW |
| `database/models/user.go` | MODIFIED |
| `database/migrations/20260308000002_add_totp_fields.go` | NEW |
| `go.mod` | MODIFIED |
| `go.sum` | MODIFIED |

## Migration Guide

- New `TOTP_ENCRYPTION_KEY` env var required if enabling TOTP (64-char hex string = 32-byte AES key)
- Optional `TOTP_ISSUER` env var (defaults to `"RapidGo"`)
- Run `migrate` to add TOTP columns to the `users` table
- Existing users are unaffected (`totp_enabled` defaults to `false`)
