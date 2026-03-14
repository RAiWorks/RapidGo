---
title: "Zap Logger Integration"
version: "0.1.0"
status: "Final"
date: "2026-03-15"
last_updated: "2026-03-15"
authors:
  - "raiworks"
supersedes: ""
---

# Zap Logger Integration

## Abstract

RapidGo v2.6.0 introduced a `Logger` interface with a default `slog`
implementation. This guide shows how to swap it with Uber's Zap for
structured, high-performance logging.

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [The Logger Interface](#2-the-logger-interface)
3. [Install Zap](#3-install-zap)
4. [Create a Zap Adapter](#4-create-a-zap-adapter)
5. [Register via Provider](#5-register-via-provider)
6. [Configuration](#6-configuration)
7. [References](#7-references)

## 1. Prerequisites

- RapidGo v2.6.0+ (for `Logger` interface)
- Go 1.21+

## 2. The Logger Interface

RapidGo's `Logger` interface defines 5 methods:

```go
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
}
```

Any struct implementing these methods can replace the default `slog`
logger.

## 3. Install Zap

```bash
go get go.uber.org/zap
```

## 4. Create a Zap Adapter

Create `app/helpers/zap_logger.go`:

```go
package helpers

import (
	"github.com/raiworks/rapidgo/v2/core/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger adapts *zap.SugaredLogger to the RapidGo Logger interface.
type ZapLogger struct {
	sugar *zap.SugaredLogger
}

// NewZapLogger creates a ZapLogger from the given *zap.Logger.
func NewZapLogger(z *zap.Logger) *ZapLogger {
	return &ZapLogger{sugar: z.Sugar()}
}

func (l *ZapLogger) Debug(msg string, args ...any) { l.sugar.Debugw(msg, args...) }
func (l *ZapLogger) Info(msg string, args ...any)  { l.sugar.Infow(msg, args...) }
func (l *ZapLogger) Warn(msg string, args ...any)  { l.sugar.Warnw(msg, args...) }
func (l *ZapLogger) Error(msg string, args ...any) { l.sugar.Errorw(msg, args...) }

func (l *ZapLogger) With(args ...any) logger.Logger {
	return &ZapLogger{sugar: l.sugar.With(args...)}
}

// NewZapFromEnv creates a Zap logger configured from environment variables.
// Reads LOG_LEVEL (debug|info|warn|error) and LOG_FORMAT (json|console).
func NewZapFromEnv(level, format string) (*zap.Logger, error) {
	var cfg zap.Config

	if format == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// Parse level
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}
	cfg.Level.SetLevel(zapLevel)

	return cfg.Build()
}

// Compile-time check: ZapLogger implements Logger.
var _ logger.Logger = (*ZapLogger)(nil)
```

## 5. Register via Provider

Create `app/providers/logger_provider.go`:

```go
package providers

import (
	"github.com/raiworks/rapidgo/v2/core/config"
	"github.com/raiworks/rapidgo/v2/core/container"
	"myapp/app/helpers"
)

type ZapLoggerProvider struct{}

func (p *ZapLoggerProvider) Register(c *container.Container) {
	c.Singleton("logger", func(_ *container.Container) interface{} {
		level := config.Env("LOG_LEVEL", "info")
		format := config.Env("LOG_FORMAT", "json")

		z, err := helpers.NewZapFromEnv(level, format)
		if err != nil {
			panic("zap init: " + err.Error())
		}

		return helpers.NewZapLogger(z)
	})
}

func (p *ZapLoggerProvider) Boot(c *container.Container) {}
```

Register it in your bootstrap — it overrides the default slog binding:

```go
cli.SetBootstrap(func(a *app.App, mode service.Mode) {
	// Register Zap AFTER the default logger provider
	// so it overrides the "logger" binding.
	a.Register(&providers.ZapLoggerProvider{})
	// ... other providers ...
})
```

## 6. Configuration

The adapter reads the same env vars as the default logger:

```env
LOG_LEVEL=info      # debug, info, warn, error
LOG_FORMAT=json     # json or console
```

The `With()` method works with key-value pairs, matching Zap's
sugared logger convention:

```go
log := container.MustMake[logger.Logger](c, "logger")
log.Info("user created", "user_id", 42, "email", "a@b.com")

// Create a child logger with context
reqLog := log.With("request_id", "abc-123")
reqLog.Info("handling request")
```

## 7. References

- [Zap documentation](https://pkg.go.dev/go.uber.org/zap)
- [RapidGo Logger](../references/logger.md)
- [RapidGo Extending the Framework](extending-framework.md)
