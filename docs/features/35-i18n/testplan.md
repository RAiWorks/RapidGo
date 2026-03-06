# Feature #35 — Localization / i18n: Test Plan

## Test cases

| TC | Description | Method | Expected |
|----|-------------|--------|----------|
| TC-01 | LoadFile loads JSON translations | LoadFile + Get | Returns loaded message |
| TC-02 | Get missing key returns raw key | Get("en", "nope") | `"nope"` |
| TC-03 | Get falls back to fallback locale | Get("fr", key) with only "en" loaded | Returns "en" message |
| TC-04 | Get with template args interpolates | Get(_, "welcome", map) | Variables replaced |
| TC-05 | Get with no args returns plain message | Get(_, "errors.not_found") | No interpolation |
| TC-06 | LoadFile returns error for missing file | LoadFile(_, "nope.json") | error != nil |
| TC-07 | LoadFile returns error for invalid JSON | LoadFile(_, badfile) | error != nil |
| TC-08 | LoadDir loads all JSON files | LoadDir(dir) + Get | All locales available |
| TC-09 | LoadDir skips non-JSON files | Dir with .txt file | No error, .txt ignored |
| TC-10 | Concurrent Get is safe | Parallel goroutines | No race/panic |

## Notes

- Tests use `t.TempDir()` to create temporary translation files.
- TC-04 uses `map[string]string{"Name": "Carlos"}` as template data.
- TC-08 creates en.json and es.json in temp dir.
- TC-10 runs parallel Get calls to verify thread safety.
