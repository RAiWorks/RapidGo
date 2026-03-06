# Feature #35 — Localization / i18n: Review

## Summary

Implemented a JSON-based translation system with locale fallback and template interpolation, plus a convenience `LoadDir()` method beyond the blueprint.

## Files changed

| File | Action | Purpose |
|------|--------|---------|
| `core/i18n/i18n.go` | Created | `Translator`, `NewTranslator()`, `LoadFile()`, `LoadDir()`, `Get()`, `resolve()` |
| `core/i18n/i18n_test.go` | Created | 10 test cases |
| `core/i18n/.gitkeep` | Deleted | Replaced by real implementation |

## Blueprint compliance

| Blueprint item | Status | Notes |
|----------------|--------|-------|
| `Translator` struct with RWMutex + messages map + fallback | ✅ | Exact match |
| `NewTranslator(fallback)` | ✅ | Exact match |
| `LoadFile(locale, path)` | ✅ | Exact match |
| `Get(locale, key, args...)` with fallback + interpolation | ✅ | Exact match |
| `resolve(locale, key)` internal helper | ✅ | Exact match |

## Deviations

| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | No `LoadDir()` | Added `LoadDir(dir) error` | Convenience for loading all locale files at once |
| 2 | `tmpl.Execute` error ignored | Error checked, returns raw msg on failure | Defensive improvement |

## Test results

| TC | Description | Result |
|----|-------------|--------|
| TC-01 | LoadFile loads JSON translations | ✅ PASS |
| TC-02 | Get missing key returns raw key | ✅ PASS |
| TC-03 | Get falls back to fallback locale | ✅ PASS |
| TC-04 | Get with template args interpolates | ✅ PASS |
| TC-05 | Get with no args returns plain message | ✅ PASS |
| TC-06 | LoadFile returns error for missing file | ✅ PASS |
| TC-07 | LoadFile returns error for invalid JSON | ✅ PASS |
| TC-08 | LoadDir loads all JSON files | ✅ PASS |
| TC-09 | LoadDir skips non-JSON files | ✅ PASS |
| TC-10 | Concurrent Get is safe | ✅ PASS |

## Regression

- All 29 packages pass (`go test ./...`)
- `go vet ./...` clean
- No new dependencies (stdlib only)
