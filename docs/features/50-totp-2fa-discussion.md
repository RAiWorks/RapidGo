# 💬 Discussion: Two-Factor Authentication (TOTP)

> **Feature**: `50` — Two-Factor Authentication (TOTP)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## What Are We Building?

A framework-level TOTP (Time-based One-Time Password) package that provides RFC 6238-compliant two-factor authentication. The `core/totp` package wraps [`pquerna/otp`](https://github.com/pquerna/otp) to offer secret generation, QR provisioning URI creation, code validation, and one-time backup code management. The User model gains TOTP fields, and a database migration adds the columns to the `users` table.

---

## Why?

- **Security**: passwords alone are insufficient — TOTP adds a second factor that requires physical access to an authenticator device
- **Industry standard**: RFC 6238 TOTP is supported by Google Authenticator, Authy, Microsoft Authenticator, 1Password, and all major platforms
- **Roadmap priority**: listed as **High** in the roadmap under "Enhanced security"
- **Framework completeness**: Laravel, Django, Rails, and Phoenix all provide 2FA utilities — RapidGo should too
- **Building blocks exist**: `core/crypto` has AES-256-GCM encryption for secret storage, `core/auth/jwt.go` handles tokens, sessions handle state

---

## Prior Art

| System | Approach | Notes |
|---|---|---|
| Laravel (Fortify) | `google2fa-laravel` package, trait-based | QR generation + validation, recovery codes, two-step enable flow |
| Django | `django-otp` / `django-two-factor-auth` | Pluggable backends (TOTP, HOTP, SMS), middleware-based enforcement |
| Rails (Devise) | `devise-two-factor` gem | Attr-encrypted secret, `otp_required_for_login` flag |
| Go (`pquerna/otp`) | Standalone RFC 6238/4226 library | Key generation, validation, QR image support, no framework opinions |

---

## Constraints

1. **Use `pquerna/otp`** — mature, RFC-compliant, no external dependencies beyond stdlib, QR provisioning URI support, Apache 2.0 licensed
2. **Framework provides building blocks, not controllers/routes** — `core/totp` package offers generate/validate/backup-code functions; application code builds auth flows on top (same pattern as `core/auth/jwt.go`)
3. **TOTP secret encrypted at rest** — use existing `core/crypto.Encrypt()` (AES-256-GCM) with a config key; never store plaintext secrets in the database
4. **Opt-in per user** — `TOTPEnabled` boolean flag on User model; existing users unaffected until they explicitly enable 2FA
5. **Backup codes** — generate a set of one-time recovery codes; store as bcrypt hashes (same pattern as passwords); each code usable once
6. **No SMS/email OTP** — TOTP only; SMS/email are separate features with external dependencies
7. **No controllers or routes** — those are application-level; the framework provides `core/totp` + model fields + migration
8. **Config via environment** — `TOTP_ISSUER` (app name for QR code), `TOTP_ENCRYPTION_KEY` (32-byte hex key for encrypting secrets)

---

## Decision Log

| # | Decision | Rationale |
|---|---|---|
| 1 | Use `pquerna/otp` library | RFC 6238 compliant, production-grade, QR support, no CGO, Apache 2.0, widely adopted |
| 2 | `core/totp/` as standalone package | Follows framework pattern (like `core/auth/`, `core/crypto/`); no coupling to HTTP layer |
| 3 | Encrypt TOTP secret with AES-256-GCM | Reuse `core/crypto.Encrypt()`; protects against DB dumps; key from env `TOTP_ENCRYPTION_KEY` |
| 4 | Bcrypt backup code hashes | Same pattern as password storage; one-way hash; constant-time comparison |
| 5 | User model gets 4 fields | `TOTPEnabled bool`, `TOTPSecret string` (encrypted), `TOTPVerifiedAt *time.Time`, `BackupCodesHash string` (JSON array of hashes) |
| 6 | No routes/controllers in this feature | Framework provides primitives; app developers compose flows. Matches `core/auth/jwt.go` pattern (no login controller shipped) |
| 7 | Provisioning URI instead of QR image | `key.URL()` returns `otpauth://` URI; clients can render QR codes themselves (a URI is more flexible than a PNG) |
| 8 | Configurable issuer name | `TOTP_ISSUER` env var; defaults to "RapidGo"; appears in authenticator app as the account label |

---

## Open Questions

_None — all resolved._

---

## Next

Architecture doc → `50-totp-2fa-architecture.md`
