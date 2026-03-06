# Feature #35 — Localization / i18n: Architecture

## Component overview

```
core/i18n/i18n.go
    │
    ├── Translator struct
    │       sync.RWMutex + map[string]map[string]string + fallback
    │
    ├── NewTranslator(fallback) *Translator
    │
    ├── LoadFile(locale, path) error
    │       Reads JSON file into locale's message map
    │
    ├── LoadDir(dir) error
    │       Loads all *.json files from directory (filename = locale)
    │
    ├── Get(locale, key, args...) string
    │       Resolves message → fallback → raw key
    │       Interpolates args via text/template
    │
    └── resolve(locale, key) string
            Internal lookup helper
```

## New file

| File | Purpose |
|------|---------|
| `core/i18n/i18n.go` | `Translator`, `NewTranslator()`, `LoadFile()`, `LoadDir()`, `Get()` |

## Removed

| File | Reason |
|------|--------|
| `core/i18n/.gitkeep` | Replaced by real implementation |

## Types

```go
type Translator struct {
    mu       sync.RWMutex
    messages map[string]map[string]string // locale -> key -> message
    fallback string
}
```

## Functions

| Function | Signature | Behaviour |
|----------|-----------|-----------|
| `NewTranslator()` | `func NewTranslator(fallback string) *Translator` | Returns empty translator with fallback locale |
| `LoadFile()` | `func (t *Translator) LoadFile(locale, path string) error` | Reads JSON, stores in messages map |
| `LoadDir()` | `func (t *Translator) LoadDir(dir string) error` | Reads all `*.json` files, locale = filename without ext |
| `Get()` | `func (t *Translator) Get(locale, key string, args ...interface{}) string` | Resolves message with fallback + template interpolation |
| `resolve()` | `func (t *Translator) resolve(locale, key string) string` | Internal: looks up locale→key in map |

## Translation file format

```json
{
    "welcome": "Welcome, {{.Name}}!",
    "errors.not_found": "Resource not found"
}
```

## Usage

```go
trans := i18n.NewTranslator("en")
trans.LoadFile("en", "resources/lang/en.json")
trans.LoadFile("es", "resources/lang/es.json")

trans.Get("es", "welcome", map[string]string{"Name": "Carlos"})
// "¡Bienvenido, Carlos!"

trans.Get("fr", "welcome", map[string]string{"Name": "Pierre"})
// Falls back to "en" → "Welcome, Pierre!"

trans.Get("en", "missing.key")
// Returns "missing.key" (raw key)
```
