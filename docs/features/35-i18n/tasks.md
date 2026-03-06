# Feature #35 — Localization / i18n: Tasks

## Prerequisites

- [x] Configuration system shipped (#02)
- [x] `core/i18n/` directory exists (stub with `.gitkeep`)

## Implementation tasks

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1 | Define `Translator` struct | `core/i18n/i18n.go` | ⬜ |
| 2 | Implement `NewTranslator()`, `LoadFile()`, `resolve()` | `core/i18n/i18n.go` | ⬜ |
| 3 | Implement `Get()` with fallback + template interpolation | `core/i18n/i18n.go` | ⬜ |
| 4 | Implement `LoadDir()` convenience method | `core/i18n/i18n.go` | ⬜ |
| 5 | Remove `.gitkeep` | `core/i18n/.gitkeep` | ⬜ |
| 6 | Write tests | `core/i18n/i18n_test.go` | ⬜ |
| 7 | Full regression + `go vet` | — | ⬜ |
| 8 | Commit, merge, review doc, roadmap update, push | — | ⬜ |

## Acceptance criteria

- `NewTranslator(fallback)` sets fallback locale.
- `LoadFile(locale, path)` reads JSON file into messages map.
- `LoadDir(dir)` loads all `*.json` files from directory.
- `Get(locale, key)` returns message for locale, falls back, returns raw key if not found.
- `Get(locale, key, args)` interpolates template variables.
- Thread-safe with `sync.RWMutex`.
- All existing tests pass (regression).
