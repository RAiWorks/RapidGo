# 🧪 Test Plan: Logging

> **Feature**: `03` — Logging
> **Tasks**: [`03-logging-tasks.md`](03-logging-tasks.md)
> **Status**: 🟢 Ready for Execution
> **Test Cases**: 10 (6 unit + 2 integration + 2 edge)

---

## Unit Tests — `core/logger/logger_test.go`

### TC-01: parseLevel() with valid levels

| Field | Detail |
|---|---|
| **Action** | Call `parseLevel()` for "debug", "info", "warn", "error" |
| **Expected** | Returns `slog.LevelDebug`, `slog.LevelInfo`, `slog.LevelWarn`, `slog.LevelError` respectively |
| **Pass Criteria** | All 4 mappings correct |

### TC-02: parseLevel() with unknown string

| Field | Detail |
|---|---|
| **Action** | Call `parseLevel("invalid")` |
| **Expected** | Returns `slog.LevelInfo` (fallback) |
| **Pass Criteria** | Return value equals `slog.LevelInfo` |

### TC-03: Setup() with JSON format

| Field | Detail |
|---|---|
| **Precondition** | `LOG_FORMAT=json`, `LOG_LEVEL=info`, `LOG_OUTPUT=stdout` |
| **Action** | Call `Setup()`, write a log entry, capture output |
| **Expected** | Output is valid JSON with `"level"` and `"msg"` fields |
| **Pass Criteria** | Output parses as JSON successfully |

### TC-04: Setup() with text format

| Field | Detail |
|---|---|
| **Precondition** | `LOG_FORMAT=text`, `LOG_LEVEL=info`, `LOG_OUTPUT=stdout` |
| **Action** | Call `Setup()`, write a log entry, capture output |
| **Expected** | Output is slog text format (key=value pairs) |
| **Pass Criteria** | Output contains `level=INFO` and `msg=` |

### TC-05: Setup() with file output

| Field | Detail |
|---|---|
| **Precondition** | `LOG_OUTPUT=file`, temp directory for storage/logs |
| **Action** | Call `Setup()`, write a log entry |
| **Expected** | `storage/logs/app.log` is created with log content |
| **Pass Criteria** | File exists, contains the log entry |

### TC-06: Setup() log level filtering

| Field | Detail |
|---|---|
| **Precondition** | `LOG_LEVEL=warn`, `LOG_FORMAT=json`, `LOG_OUTPUT=stdout` |
| **Action** | Call `Setup()`, write `slog.Info()` and `slog.Warn()` |
| **Expected** | Info message is suppressed, Warn message appears |
| **Pass Criteria** | Only warn-level output captured |

---

## Integration Tests — Manual Verification

### TC-07: go run with default .env

| Field | Detail |
|---|---|
| **Precondition** | `.env` has `LOG_LEVEL=debug`, `LOG_FORMAT=json`, `LOG_OUTPUT=stdout` |
| **Action** | Run `go run ./cmd/` |
| **Expected** | Banner displays via fmt, then JSON log line with "server initialized" |
| **Pass Criteria** | Both banner and structured log visible in output |

### TC-08: go vet passes

| Field | Detail |
|---|---|
| **Action** | Run `go vet ./...` |
| **Expected** | Zero warnings across all packages |
| **Pass Criteria** | Exit code 0, no output |

---

## Edge Cases

### TC-09: Setup() with invalid LOG_LEVEL

| Field | Detail |
|---|---|
| **Precondition** | `LOG_LEVEL=garbage` |
| **Action** | Call `Setup()` |
| **Expected** | Defaults to `slog.LevelInfo`, no panic |
| **Pass Criteria** | Logger works at info level |

### TC-10: Setup() with invalid LOG_FORMAT

| Field | Detail |
|---|---|
| **Precondition** | `LOG_FORMAT=yaml` |
| **Action** | Call `Setup()` |
| **Expected** | Defaults to JSON handler, no panic |
| **Pass Criteria** | Output is valid JSON |

---

## Execution Notes

- Unit tests (TC-01 through TC-06, TC-09, TC-10) run via `go test ./core/logger/... -v`
- Tests that capture slog output should redirect to a buffer using a custom handler or capture stdout
- File output tests (TC-05) should use `t.TempDir()` to avoid polluting the project's `storage/logs/`
- Integration tests (TC-07, TC-08) run manually in terminal
