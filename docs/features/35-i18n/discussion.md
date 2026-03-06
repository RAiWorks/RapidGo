# Feature #35 — Localization / i18n: Discussion

## What problem does this solve?

Applications serving multiple languages need a way to load translation strings and resolve them by locale + key. A `Translator` provides this with fallback to a default locale and support for template variables in messages.

## Why now?

Configuration (#02) is shipped. The translator reads JSON files from disk and uses Go's `text/template` for variable interpolation — no external dependencies required.

## What does the blueprint specify?

- `Translator` struct with `messages map[string]map[string]string`, `fallback` locale, `sync.RWMutex`.
- `NewTranslator(fallback)` constructor.
- `LoadFile(locale, path)` — reads a JSON file into the locale's message map.
- `Get(locale, key, args...)` — resolves message, falls back to fallback locale, returns key if not found.
- Template interpolation via `text/template` when args are provided.
- `resolve(locale, key)` — internal lookup helper.
- JSON translation files in `resources/lang/{locale}.json`.

## Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Package | `core/i18n/` | Blueprint-specified; existing stub directory |
| File format | JSON | Blueprint-specified; simple and standard |
| Template engine | `text/template` | Blueprint-specified; stdlib, no external deps |
| Fallback | Configurable default locale | Blueprint-specified; returns raw key as last resort |
| Thread safety | `sync.RWMutex` | Blueprint-specified; safe for concurrent reads |
| LoadDir | Add `LoadDir(dir)` helper | Convenience for loading all `*.json` files from a directory |

## What is out of scope?

- Pluralization rules (app-level or future enhancement).
- Nested JSON keys (blueprint uses flat keys like `"errors.not_found"`).
- Accept-Language header parsing (middleware concern).
- Database-backed translations.
