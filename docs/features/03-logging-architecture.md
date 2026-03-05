# 🏗️ Architecture: Logging

> **Feature**: `03` — Logging
> **Discussion**: [`03-logging-discussion.md`](03-logging-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-05

---

## Overview

Create the `core/logger` package that configures Go's standard `log/slog` based on `.env` settings (`LOG_LEVEL`, `LOG_FORMAT`, `LOG_OUTPUT`). The `Setup()` function initializes the global default logger and returns it. Zero external dependencies.

## File Structure

```
core/logger/
├── logger.go           # Setup(), Close() — configure slog handler, set default
├── level.go            # parseLevel() — map string to slog.Level
└── logger_test.go      # Unit tests for Setup(), Close(), and parseLevel()

cmd/
└── main.go             # MODIFY — add logger.Setup() call after config.Load()
```

## Data Model

N/A — no database entities. Logging is a runtime concern.

## Component Design

### `core/logger/level.go`

**Responsibility**: Map config string to `slog.Level` constant.
**Package**: `logger`

```go
package logger

import "log/slog"

// parseLevel converts a string level name to a slog.Level.
// Returns slog.LevelInfo if the string is unrecognized.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
```

### `core/logger/logger.go`

**Responsibility**: Configure and initialize the global slog logger.
**Package**: `logger`

```go
package logger

import (
	"log/slog"
	"os"

	"github.com/RAiWorks/RGo/core/config"
)

// logFile holds the open log file handle (if LOG_OUTPUT=file).
var logFile *os.File

// Setup initializes the global slog logger based on config values.
// Reads LOG_LEVEL, LOG_FORMAT, LOG_OUTPUT from environment.
// Sets slog.SetDefault() so that slog.Info(), slog.Error() etc. work globally.
// Returns the configured logger instance.
func Setup() *slog.Logger {
	level := parseLevel(config.Env("LOG_LEVEL", "info"))
	format := config.Env("LOG_FORMAT", "json")
	output := config.Env("LOG_OUTPUT", "stdout")

	var writer *os.File
	if output == "file" {
		if err := os.MkdirAll("storage/logs", 0755); err != nil {
			slog.Warn("failed to create log directory, falling back to stdout", "err", err)
			writer = os.Stdout
		} else {
			f, err := os.OpenFile("storage/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				slog.Warn("failed to open log file, falling back to stdout", "err", err)
				writer = os.Stdout
			} else {
				logFile = f
				writer = f
			}
		}
	} else {
		writer = os.Stdout
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if format == "text" {
		handler = slog.NewTextHandler(writer, opts)
	} else {
		handler = slog.NewJSONHandler(writer, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// Close closes the log file if one was opened.
// Should be called on application shutdown.
func Close() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}
```

### `cmd/main.go` (MODIFY)

**Changes**: Add `logger.Setup()` after `config.Load()`, add structured log entry.

```go
package main

import (
	"fmt"
	"log/slog"

	"github.com/RAiWorks/RGo/core/config"
	"github.com/RAiWorks/RGo/core/logger"
)

func main() {
	config.Load()
	logger.Setup()

	appName := config.Env("APP_NAME", "RGo")
	appPort := config.Env("APP_PORT", "8080")
	appEnv := config.AppEnv()

	fmt.Println("=================================")
	fmt.Printf("  %s Framework\n", appName)
	fmt.Println("  github.com/RAiWorks/RGo")
	fmt.Println("=================================")
	fmt.Printf("  Environment: %s\n", appEnv)
	fmt.Printf("  Port: %s\n", appPort)
	fmt.Printf("  Debug: %v\n", config.IsDebug())
	fmt.Println("=================================")

	slog.Info("server initialized",
		"app", appName,
		"port", appPort,
		"env", appEnv,
	)
}
```

## Data Flow

```
Application Start
    │
    ▼
config.Load()          ← .env loaded into os environment
    │
    ▼
logger.Setup()
    │
    ├── config.Env("LOG_LEVEL", "info")   → parseLevel() → slog.Level
    ├── config.Env("LOG_FORMAT", "json")  → select JSONHandler or TextHandler
    ├── config.Env("LOG_OUTPUT", "stdout")→ os.Stdout or storage/logs/app.log
    │
    ├── Create slog.Handler with level + output
    ├── slog.New(handler) → *slog.Logger
    ├── slog.SetDefault(logger) → global default set
    │
    ▼
slog.Info() / slog.Error() / slog.Warn() / slog.Debug()
    → Available globally throughout the framework
```

## Configuration

| Key | Type | Default | Values | Used By |
|---|---|---|---|---|
| `LOG_LEVEL` | string | `info` | `debug`, `info`, `warn`, `error` | `parseLevel()` |
| `LOG_FORMAT` | string | `json` | `json`, `text` | `Setup()` handler selection |
| `LOG_OUTPUT` | string | `stdout` | `stdout`, `file` | `Setup()` writer selection |

These keys already exist in `.env` from Feature #01.

## Security Considerations

- **NEVER log sensitive data**: passwords, API keys, JWT tokens, credit card numbers, PII
- **Log files** at `storage/logs/app.log` use permissions `0644`
- **`storage/logs/` is gitignored** — log files never reach version control
- **No log rotation** in this feature — handled externally (logrotate, Docker log drivers)

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| `log/slog` (stdlib) | Zero deps, Go standard, structured | Less features than zerolog/zap | ✅ Selected |
| zerolog | Zero-allocation, very fast | External dep, overkill for now | ❌ Add later if needed |
| zap | Feature-rich, structured | External dep, heavier API | ❌ Add later if needed |
| Multi-output (stdout + file) | See logs everywhere | Complexity, io.MultiWriter | ❌ Deferred |

## Next

Create tasks doc → `03-logging-tasks.md`
