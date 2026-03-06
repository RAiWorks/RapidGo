# 🧪 Test Plan: Crypto & Security Utilities

> **Feature**: `22` — Crypto & Security Utilities
> **Total Test Cases**: 15

---

## Random Generation — `core/crypto/crypto_test.go`

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-01 | RandomBytes returns n bytes | `n=16` | `len(result) == 16`, no error |
| TC-02 | RandomBytes returns unique values | Two calls with `n=16` | Different byte slices |
| TC-03 | RandomHex correct length | `n=16` | 32-char hex string |
| TC-04 | RandomBase64 non-empty | `n=24` | Non-empty string, valid base64 |

## Hashing

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-05 | SHA256Hash deterministic | `"hello"` twice | Same 64-char hex both times |
| TC-06 | SHA256Hash length | `"test"` | 64-char hex string |

## HMAC

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-07 | HMACSign deterministic | Same message + key | Same signature |
| TC-08 | HMACVerify valid | Sign then verify | `true` |
| TC-09 | HMACVerify tampered message | Sign "a", verify "b" | `false` |
| TC-10 | HMACVerify wrong key | Sign with key1, verify with key2 | `false` |

## AES-256-GCM

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-11 | Encrypt/Decrypt round-trip | `"secret data"`, 32-byte key | Decrypted == original |
| TC-12 | Encrypt produces unique output | Same plaintext twice | Different ciphertext (random nonce) |
| TC-13 | Decrypt fails wrong key | Encrypt with key1, decrypt with key2 | Error |
| TC-14 | Decrypt fails short ciphertext | `"abc"` as encrypted input | Error: "ciphertext too short" |
| TC-15 | Encrypt fails invalid key | 16-byte key (not 32) | Error |

---

## Acceptance Criteria

1. All 15 tests pass
2. Full regression (`go test ./... -count=1`) — 0 failures
3. `go vet ./...` — clean
4. No new dependencies in `go.mod`
