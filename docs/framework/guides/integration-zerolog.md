---
title: "Zerolog Logger Integration"
version: "0.1.0"
status: "Final"
date: "2026-03-15"
last_updated: "2026-03-15"
authors:
  - "raiworks"
supersedes: ""
---

# Zerolog Logger Integration

## Abstract

This guide shows how to use `rs/zerolog` with RapidGo by implementing
the `Logger` interface. Zerolog provides zero-allocation JSON logging,
making it a good choice for high-throughput services.

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Install Zerolog](#2-install-zerolog)
3. [Create a Zerolog Adapter](#3-create-a-zerolog-adapter)
4. [Register via Provider](#4-register-via-provider)
5. [Configuration](#5-configuration)
6. [References](#6-references)

## 1. Prerequisites

- RapidGo v2.6.0+ (for `Logger` interface)
- Go 1.21+

## 2. Install Zerolog

```bash
go get github.com/rs/zerolog
```

## 3. Create a Zerolog Adapter

Create `app/helpers/zerolog_logger.go`:

```go
package helpers

import (
	"os"

	"github.com/raiworks/rapidgo/v2/core/logger"
	"github.com/rs/zerolog"
)

// ZerologLogger adapts zerolog.Logger to the RapidGo Logger interface.
type ZerologLogger struct {
	zl zerolog.Logger
}

// NewZerologLogger wraps a zerolog.Logger.
func NewZerologLogger(zl zerolog.Logger) *ZerologLogger {
	return &ZerologLogger{zl: zl}
}

func (l *ZerologLogger) Debug(msg string, args ...any) {
	l.zl.Debug().Fields(toFields(args)).Msg(msg)
}

func (l *ZerologLogger) Info(msg string, args ...any) {
	l.zl.Info().Fields(toFields(args)).Msg(msg)
}

func (l *ZerologLogger) Warn(msg string, args ...any) {
	l.zl.Warn().Fields(toFields(args)).Msg(msg)
}

func (l *ZerologLogger) Error(msg string, args ...any) {
	l.zl.Error().Fields(toFields(args)).Msg(msg)
}

func (l *ZerologLogger) With(args ...any) logger.Logger {
	ctx := l.zl.With()
	fields := toFields(args)
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &ZerologLogger{zl: ctx.Logger()}
}

// toFields converts key-value pairs to a map for zerolog.
func toFields(args []any) map[string]interface{} {
	fields := make(map[string]interface{}, len(args)/2)
	for i := 0; i+1 < len(args); i += 2 {
		if key, ok := args[i].(string); ok {
			fields[key] = args[i+1]
		}
	}
	return fields
}

// NewZerologFromEnv creates a zerolog.Logger from environment settings.
func NewZerologFromEnv(level, format string) zerolog.Logger {
	var zl zerolog.Logger

	if format == "console" {
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().Timestamp().Logger()
	} else {
		zl = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	switch level {
	case "debug":
		zl = zl.Level(zerolog.DebugLevel)
	case "warn":
		zl = zl.Level(zerolog.WarnLevel)
	case "error":
		zl = zl.Level(zerolog.ErrorLevel)
	default:
		zl = zl.Level(zerolog.InfoLevel)
	}

	return zl
}

// Compile-time check: ZerologLogger implements Logger.
var _ logger.Logger = (*ZerologLogger)(nil)
```

## 4. Register via Provider

Create `app/providers/logger_provider.go`:

```go
package providers

import (
	"github.com/raiworks/rapidgo/v2/core/config"
	"github.com/raiworks/rapidgo/v2/core/container"
	"myapp/app/helpers"
)

type ZerologProvider struct{}

func (p *ZerologProvider) Register(c *container.Container) {
	c.Singleton("logger", func(_ *container.Container) interface{} {
		level := config.Env("LOG_LEVEL", "info")
		format := config.Env("LOG_FORMAT", "json")
		zl := helpers.NewZerologFromEnv(level, format)
		return helpers.NewZerologLogger(zl)
	})
}

func (p *ZerologProvider) Boot(c *container.Container) {}
```

Register it in your bootstrap to override the default slog logger:

```go
cli.SetBootstrap(func(a *app.App, mode service.Mode) {
	a.Register(&providers.ZerologProvider{})
	// ... other providers ...
})
```

## 5. Configuration

Uses the same env vars as the default logger:

```env
LOG_LEVEL=info      # debug, info, warn, error
LOG_FORMAT=json     # json or console
```

Usage is identical to any other Logger implementation:

```go
log := container.MustMake[logger.Logger](c, "logger")
log.Info("request handled", "method", "GET", "path", "/api/users", "status", 200)
```

## 6. References

- [Zerolog documentation](https://pkg.go.dev/github.com/rs/zerolog)
- [RapidGo Logger](../references/logger.md)
- [RapidGo Extending the Framework](extending-framework.md)
