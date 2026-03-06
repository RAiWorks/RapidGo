# 📋 Tasks: Crypto & Security Utilities

> **Feature**: `22` — Crypto & Security Utilities
> **Branch**: `feature/22-crypto-security`
> **Status**: NOT STARTED

---

## Phase 1 — Implementation

- [ ] Create `core/crypto/crypto.go` with all 8 functions
- [ ] Remove `core/crypto/.gitkeep` (replaced by real file)
- [ ] Verify: `go build ./core/crypto/...` compiles

**Checkpoint**: All 8 functions compile with stdlib-only imports.

## Phase 2 — Tests

- [ ] Create `core/crypto/crypto_test.go`
  - [ ] TC-01: RandomBytes returns n bytes
  - [ ] TC-02: RandomBytes returns different values each call
  - [ ] TC-03: RandomHex returns correct length (2n hex chars)
  - [ ] TC-04: RandomBase64 returns non-empty base64 string
  - [ ] TC-05: SHA256Hash returns consistent hash for same input
  - [ ] TC-06: SHA256Hash returns 64-char hex string
  - [ ] TC-07: HMACSign returns consistent signature
  - [ ] TC-08: HMACVerify returns true for valid signature
  - [ ] TC-09: HMACVerify returns false for tampered message
  - [ ] TC-10: HMACVerify returns false for wrong key
  - [ ] TC-11: Encrypt/Decrypt round-trip
  - [ ] TC-12: Encrypt produces different ciphertext each call (random nonce)
  - [ ] TC-13: Decrypt fails with wrong key
  - [ ] TC-14: Decrypt fails with short ciphertext
  - [ ] TC-15: Encrypt fails with invalid key length
- [ ] Run full `go test ./... -count=1` — all pass

**Checkpoint**: All 15 tests pass. Full regression green.

## Phase 3 — Finalize

- [ ] Update changelog
- [ ] Run `go vet ./...` — clean
- [ ] Commit and push
