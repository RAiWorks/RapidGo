# 💬 Discussion: Crypto & Security Utilities

> **Feature**: `22` — Crypto & Security Utilities
> **Status**: COMPLETE
> **Date**: 2026-03-06

---

## 1. What Are We Building?

A `core/crypto/` package providing 8 cryptographic utility functions using only Go stdlib. These are framework-level building blocks for token generation, data hashing, HMAC signing/verification, and AES-256-GCM encryption/decryption.

## 2. Current State

- `core/crypto/` directory exists but is empty (`.gitkeep` only)
- `app/helpers/random.go` has `RandomString(n int)` — hex-encoded crypto/rand (overlaps with `RandomHex` but different package/scope)
- `core/session/cookie_store.go` already has AES-256-GCM encrypt/decrypt — private functions, not reusable
- No HMAC or SHA-256 hashing utilities exist
- Framework doc `docs/framework/security/crypto.md` already defines the API

## 3. Blueprint Scope

The blueprint (lines 2679–2803) defines exactly 8 functions:

| Function | Purpose | Stdlib Packages |
|----------|---------|-----------------|
| `RandomBytes(n int) ([]byte, error)` | Raw random bytes | `crypto/rand`, `io` |
| `RandomHex(n int) string` | Hex-encoded random | `encoding/hex` |
| `RandomBase64(n int) string` | URL-safe base64 random | `encoding/base64` |
| `SHA256Hash(data string) string` | SHA-256 hex digest | `crypto/sha256` |
| `HMACSign(message, key string) string` | HMAC-SHA256 signature | `crypto/hmac` |
| `HMACVerify(message, signature, key string) bool` | Constant-time verify | `crypto/hmac` |
| `Encrypt(plaintext string, key []byte) (string, error)` | AES-256-GCM encrypt | `crypto/aes`, `crypto/cipher` |
| `Decrypt(encrypted string, key []byte) (string, error)` | AES-256-GCM decrypt | `crypto/aes`, `crypto/cipher` |

**Zero external dependencies** — all stdlib.

## 4. Approach

Single file `core/crypto/crypto.go` with all 8 functions. Matches the blueprint exactly.

### Overlap with existing code
- `helpers.RandomString` and `crypto.RandomHex` both produce hex-encoded random strings. They coexist — helpers is app-level convenience, crypto is framework-level primitive.
- `CookieStore` has private encrypt/decrypt. Could be refactored to use `crypto.Encrypt/Decrypt` but that's a future cleanup, not in scope for #22.

## 5. Edge Cases

1. **`RandomBytes` failure**: `crypto/rand` read failure returns error
2. **`Encrypt` key length**: Must be exactly 32 bytes for AES-256; `aes.NewCipher` validates this
3. **`Decrypt` short ciphertext**: Must check `len(data) < nonceSize` before slicing
4. **`HMACVerify` constant-time**: Uses `hmac.Equal` to prevent timing attacks

## 6. Dependencies

None — all Go stdlib. No new `go.mod` entries.

## 7. Open Questions — RESOLVED

| Question | Resolution |
|----------|-----------|
| Refactor CookieStore to use crypto package? | No — future cleanup, not #22 scope |
| Remove helpers.RandomString? | No — different package, different audience |

---

**Summary**: Clean, focused feature — 8 stdlib crypto functions in one file. No external dependencies. No cross-feature impact.
