# 🧪 Test Plan: Two-Factor Authentication (TOTP)

> **Feature**: `50` — Two-Factor Authentication (TOTP)
> **Tasks**: [`50-totp-2fa-tasks.md`](50-totp-2fa-tasks.md)
> **Status**: ✅ SHIPPED
> **Date**: 2026-03-07

---

## Test File

- `core/totp/totp_test.go`

---

## Unit Tests

### 1. Key Generation

| # | Test | Expectation |
|---|---|---|
| T01 | `TestGenerateKey_ReturnsKey` | Returns non-nil `*Key` with non-empty `Secret` and `URL` |
| T02 | `TestGenerateKey_URLContainsIssuer` | `URL` contains the issuer name in `otpauth://` format |
| T03 | `TestGenerateKey_URLContainsAccount` | `URL` contains the account name |
| T04 | `TestGenerateKey_SecretIsBase32` | `Secret` is valid base32-encoded string |
| T05 | `TestGenerateKey_UniqueSecrets` | Two consecutive calls produce different secrets |

### 2. Code Validation

| # | Test | Expectation |
|---|---|---|
| T06 | `TestValidateCode_ValidCode` | Returns `true` for a code generated from the same secret at current time |
| T07 | `TestValidateCode_InvalidCode` | Returns `false` for an incorrect code |
| T08 | `TestValidateCode_WrongSecret` | Returns `false` when code is validated against a different secret |
| T09 | `TestValidateCode_EmptyCode` | Returns `false` for empty string code |

### 3. Backup Code Generation

| # | Test | Expectation |
|---|---|---|
| T10 | `TestGenerateBackupCodes_ReturnsCorrectCount` | Returns exactly `n` codes |
| T11 | `TestGenerateBackupCodes_FormatXXXXDashXXXX` | Each code matches `XXXX-XXXX` pattern (hex chars + dash) |
| T12 | `TestGenerateBackupCodes_AllUnique` | No duplicate codes in the returned slice |
| T13 | `TestGenerateBackupCodes_ZeroCount` | Returns empty slice for `n=0` |

### 4. Backup Code Hashing & Verification

| # | Test | Expectation |
|---|---|---|
| T14 | `TestHashBackupCode_ReturnsHash` | Returns non-empty string (bcrypt hash) |
| T15 | `TestVerifyBackupCode_CorrectCode` | Returns `true` for matching code and hash |
| T16 | `TestVerifyBackupCode_WrongCode` | Returns `false` for non-matching code |
| T17 | `TestVerifyBackupCode_DifferentCase` | Backup codes are case-insensitive (normalized to uppercase before hash/verify) |

---

## Acceptance Criteria

1. All 17 tests pass
2. All existing tests across all packages still pass (`go test ./...`)
3. `core/totp` package provides all 5 public functions from architecture doc
4. User model has 4 new TOTP fields with correct GORM + JSON tags
5. Migration file compiles and follows existing migration patterns
6. `go build` succeeds with no errors

---

## Next

Changelog → `50-totp-2fa-changelog.md`
