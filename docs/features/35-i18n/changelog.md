# Feature #35 — Localization / i18n: Changelog

## [Unreleased]

### Added
- `core/i18n/i18n.go` — `Translator` struct, `NewTranslator()`, `LoadFile()`, `LoadDir()`, `Get()`.
- 10 test cases (TC-01 to TC-10) for loading, fallback, interpolation, errors, concurrency.

### Removed
- `core/i18n/.gitkeep` — replaced by real implementation.

### Deviation log
| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | No `LoadDir()` method | Added `LoadDir(dir) error` | Convenience for loading all locale files from a directory at once |
