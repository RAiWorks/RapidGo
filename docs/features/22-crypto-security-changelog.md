# 📝 Changelog: Crypto & Security Utilities

> **Feature**: `22` — Crypto & Security Utilities
> **Branch**: `feature/22-crypto-security`
> **Started**: 2026-03-06
> **Completed**: 2026-03-06

---

## Log

- **BUILD**: All 15 tests pass, full regression green (268 tests), `go vet` clean
- **BUILD**: Created `core/crypto/crypto_test.go` — 15 tests (TC-01 to TC-15)
- **BUILD**: Created `core/crypto/crypto.go` — 8 functions, stdlib only
- **BUILD**: Removed `core/crypto/.gitkeep`, verified compilation
- **BUILD**: Created feature branch `feature/22-crypto-security`

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| Key-length guard in Encrypt/Decrypt | Architecture relied on `aes.NewCipher` rejection | Added explicit 32-byte check before `aes.NewCipher` | `aes.NewCipher` accepts 16/24/32-byte keys; AES-256 requires exactly 32. TC-15 validates this. |

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
