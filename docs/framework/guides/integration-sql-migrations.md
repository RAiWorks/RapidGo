---
title: "SQL File Migrations Integration"
version: "0.1.0"
status: "Final"
date: "2026-03-15"
last_updated: "2026-03-15"
authors:
  - "raiworks"
supersedes: ""
---

# SQL File Migrations Integration

## Abstract

RapidGo ships with Go-based migrations via `migrations.Register()`.
This guide shows how to add SQL file migrations alongside Go migrations
using the `golang-migrate/migrate` library.

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Install golang-migrate](#2-install-golang-migrate)
3. [Create SQL Migration Files](#3-create-sql-migration-files)
4. [Create a Migration Runner](#4-create-a-migration-runner)
5. [Add a CLI Command](#5-add-a-cli-command)
6. [When to Use SQL vs Go Migrations](#6-when-to-use-sql-vs-go-migrations)
7. [References](#7-references)

## 1. Prerequisites

- RapidGo v2.4.0+
- PostgreSQL, MySQL, or SQLite

## 2. Install golang-migrate

```bash
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
```

Replace `database/postgres` with `database/mysql` or `database/sqlite3`
for other drivers.

## 3. Create SQL Migration Files

Create a `database/sql_migrations/` directory in your project:

```
database/sql_migrations/
├── 000001_create_events_table.up.sql
├── 000001_create_events_table.down.sql
├── 000002_add_events_index.up.sql
└── 000002_add_events_index.down.sql
```

Example up migration:

```sql
-- 000001_create_events_table.up.sql
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    payload JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);
```

Example down migration:

```sql
-- 000001_create_events_table.down.sql
DROP TABLE IF EXISTS events;
```

## 4. Create a Migration Runner

Create `app/helpers/sql_migrate.go`:

```go
package helpers

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/raiworks/rapidgo/v2/core/config"
)

// NewSQLMigrator creates a migrate instance pointing at the SQL files.
func NewSQLMigrator() (*migrate.Migrate, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Env("DB_USER", "postgres"),
		config.Env("DB_PASSWORD", ""),
		config.Env("DB_HOST", "localhost"),
		config.Env("DB_PORT", "5432"),
		config.Env("DB_NAME", "app"),
		config.Env("DB_SSL_MODE", "disable"),
	)

	return migrate.New("file://database/sql_migrations", dsn)
}

// RunSQLMigrations applies all pending SQL migrations.
func RunSQLMigrations() error {
	m, err := NewSQLMigrator()
	if err != nil {
		return fmt.Errorf("sql migrate init: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("sql migrate up: %w", err)
	}
	return nil
}

// RollbackSQLMigrations rolls back the last SQL migration.
func RollbackSQLMigrations() error {
	m, err := NewSQLMigrator()
	if err != nil {
		return fmt.Errorf("sql migrate init: %w", err)
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("sql migrate rollback: %w", err)
	}
	return nil
}
```

## 5. Add a CLI Command

Register custom CLI commands via the root command in your bootstrap:

```go
package main

import (
	"fmt"

	"github.com/raiworks/rapidgo/v2/core/cli"
	"github.com/spf13/cobra"
	"myapp/app/helpers"
)

func init() {
	cli.RootCmd().AddCommand(&cobra.Command{
		Use:   "migrate:sql",
		Short: "Run SQL file migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := helpers.RunSQLMigrations(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "SQL migrations applied.")
			return nil
		},
	})

	cli.RootCmd().AddCommand(&cobra.Command{
		Use:   "migrate:sql-rollback",
		Short: "Rollback last SQL file migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := helpers.RollbackSQLMigrations(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "SQL migration rolled back.")
			return nil
		},
	})
}
```

Then run:

```bash
go run cmd/main.go migrate:sql
go run cmd/main.go migrate:sql-rollback
```

## 6. When to Use SQL vs Go Migrations

| Use Case | Recommended |
|----------|-------------|
| Schema changes (CREATE, ALTER, DROP) | SQL files |
| Data migrations (backfill, transform) | Go migrations |
| Complex logic with conditionals | Go migrations |
| DBA-reviewed DDL | SQL files |
| Seeding reference data | Go seeders |

Both systems use independent tracking tables (`schema_migrations` for
golang-migrate, `schema_migrations` with batch tracking for RapidGo's
Go migrations). They do not conflict.

## 7. References

- [golang-migrate documentation](https://github.com/golang-migrate/migrate)
- [RapidGo Migrations](../references/migrations.md)
- [RapidGo CLI Commands](../references/cli.md)
