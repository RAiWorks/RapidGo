---
title: "pgx Database Driver Integration"
version: "0.1.0"
status: "Final"
date: "2026-03-15"
last_updated: "2026-03-15"
authors:
  - "raiworks"
supersedes: ""
---

# pgx Database Driver Integration

## Abstract

This guide shows how to use the `jackc/pgx/v5` PostgreSQL driver with
RapidGo. The framework provides hooks and a service container — you
import pgx directly, register a connection pool via a provider, and
resolve it wherever you need it.

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Install pgx](#2-install-pgx)
3. [Create a pgx Provider](#3-create-a-pgx-provider)
4. [Register the Provider](#4-register-the-provider)
5. [Use the Pool](#5-use-the-pool)
6. [Health Checks](#6-health-checks)
7. [Running Alongside GORM](#7-running-alongside-gorm)
8. [References](#8-references)

## 1. Prerequisites

- RapidGo v2.6.0+ (for `Logger` interface and `LoadConfig[T]`)
- PostgreSQL 12+
- Go 1.21+

## 2. Install pgx

```bash
go get github.com/jackc/pgx/v5
```

## 3. Create a pgx Provider

Create `app/providers/pgx_provider.go` in your starter project:

```go
package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raiworks/rapidgo/v2/core/config"
	"github.com/raiworks/rapidgo/v2/core/container"
)

// PgxConfig holds pgx-specific configuration.
type PgxConfig struct {
	Host     string `env:"PGX_HOST" default:"localhost"`
	Port     int    `env:"PGX_PORT" default:"5432"`
	Database string `env:"PGX_DATABASE" validate:"required"`
	User     string `env:"PGX_USER" validate:"required"`
	Password string `env:"PGX_PASSWORD" validate:"required"`
	SSLMode  string `env:"PGX_SSL_MODE" default:"disable"`
	MaxConns int32  `env:"PGX_MAX_CONNS" default:"10"`
}

type PgxProvider struct{}

func (p *PgxProvider) Register(c *container.Container) {
	c.Singleton("pgx", func(_ *container.Container) interface{} {
		cfg, err := config.LoadConfig[PgxConfig]()
		if err != nil {
			panic(fmt.Sprintf("pgx config error: %v", err))
		}

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port,
			cfg.Database, cfg.SSLMode,
		)

		poolCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			panic(fmt.Sprintf("pgx parse config: %v", err))
		}
		poolCfg.MaxConns = cfg.MaxConns

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
		if err != nil {
			panic(fmt.Sprintf("pgx connect: %v", err))
		}
		return pool
	})
}

func (p *PgxProvider) Boot(c *container.Container) {
	// Pool is lazy-initialized on first Make("pgx") call.
}
```

## 4. Register the Provider

In your bootstrap function (typically `cmd/main.go` or `app/plugins.go`):

```go
cli.SetBootstrap(func(a *app.App, mode service.Mode) {
	// ... other providers ...
	a.Register(&providers.PgxProvider{})
})
```

Add the environment variables to `.env`:

```env
PGX_HOST=localhost
PGX_PORT=5432
PGX_DATABASE=myapp
PGX_USER=postgres
PGX_PASSWORD=secret
PGX_SSL_MODE=disable
PGX_MAX_CONNS=10
```

## 5. Use the Pool

Resolve the pool from the container wherever you need it:

```go
pool := container.MustMake[*pgxpool.Pool](c, "pgx")

rows, err := pool.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
if err != nil {
	return err
}
defer rows.Close()

for rows.Next() {
	var id int
	var name string
	if err := rows.Scan(&id, &name); err != nil {
		return err
	}
	// ...
}
```

Or use `TryMake` for safe resolution:

```go
pool, err := container.TryMake[*pgxpool.Pool](c, "pgx")
if err != nil {
	// pgx not registered — fall back to GORM or return error
}
```

## 6. Health Checks

Add a pgx health check to RapidGo's health endpoint:

```go
func (p *PgxProvider) Boot(c *container.Container) {
	pool := container.MustMake[*pgxpool.Pool](c, "pgx")

	// Register a health check
	if c.Has("health") {
		h := container.MustMake[*health.Checker](c, "health")
		h.AddCheck("pgx", func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return pool.Ping(ctx)
		})
	}
}
```

## 7. Running Alongside GORM

pgx and GORM can coexist. Use GORM for your ORM models and pgx for
performance-critical queries:

```go
// GORM for standard CRUD
db := container.MustMake[*gorm.DB](c, "db")
db.Create(&user)

// pgx for raw performance
pool := container.MustMake[*pgxpool.Pool](c, "pgx")
pool.QueryRow(ctx, "SELECT count(*) FROM events WHERE ts > $1", cutoff)
```

Both connect to the same PostgreSQL database — just configure the same
host/port/database in both `DB_*` and `PGX_*` env vars.

## 8. References

- [pgx documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [pgxpool documentation](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)
- [RapidGo Container](../references/container.md)
- [RapidGo Extending the Framework](extending-framework.md)
