# 📋 Review: Crypto & Security Utilities

> **Feature**: `22` — Crypto & Security Utilities
> **Branch**: `feature/22-crypto-security`
> **Merged**: 2026-03-06
> **Commit**: `3d3b7b5` (impl)

---

## Summary

Feature #22 adds 8 cryptographic utility functions to `core/crypto/crypto.go` using only Go stdlib packages. Covers random generation (bytes, hex, base64), SHA-256 hashing, HMAC-SHA256 signing/verification, and AES-256-GCM encryption/decryption.

## Files Changed

| File | Type | Description |
|---|---|---|
| `core/crypto/crypto.go` | Created | 8 crypto functions — RandomBytes, RandomHex, RandomBase64, SHA256Hash, HMACSign, HMACVerify, Encrypt, Decrypt |
| `core/crypto/crypto_test.go` | Created | 15 tests (TC-01 to TC-15) |
| `core/crypto/.gitkeep` | Deleted | Replaced by real implementation |
| `docs/features/22-crypto-security-changelog.md` | Modified | Updated with build log and deviation |

## Dependencies Added

None — all Go stdlib.

## Test Results

- **15 new tests** — all pass
- **Full regression**: all packages pass, 0 failures
- **`go vet`**: clean

## Deviations

| What | Why |
|---|---|
| Added explicit 32-byte key length check to `Encrypt`/`Decrypt` | `aes.NewCipher` accepts 16/24/32-byte keys; AES-256 requires exactly 32. TC-15 validates this. |

## Status: ✅ SHIPPED
