---
title: "Seeders"
version: "0.1.0"
status: "Draft"
date: "2026-03-05"
last_updated: "2026-03-05"
authors:
  - "RAiWorks"
supersedes: ""
---

# Seeders

## Abstract

This document covers database seeding — populating the database with
initial or test data using seeder functions and the CLI command.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Seeder Location](#2-seeder-location)
3. [Writing Seeders](#3-writing-seeders)
4. [Seed Runner](#4-seed-runner)
5. [CLI Command](#5-cli-command)
6. [Security Considerations](#6-security-considerations)
7. [References](#7-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Seeder** — A function that inserts predefined records into the
  database.
- **`FirstOrCreate`** — GORM method that inserts only if a matching
  record doesn't exist, making seeders idempotent.

## 2. Seeder Location

Seeders live in `database/seeders/`.

## 3. Writing Seeders

```go
package seeders

import (
    "log"

    "yourframework/database/models"
    "yourframework/app/helpers"

    "gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
    users := []models.User{
        {Name: "Admin", Email: "admin@example.com", Password: "password123", Role: "admin"},
        {Name: "User", Email: "user@example.com", Password: "password123", Role: "user"},
    }
    for _, u := range users {
        hashed, _ := helpers.HashPassword(u.Password)
        u.Password = hashed
        if err := db.FirstOrCreate(&u, models.User{Email: u.Email}).Error; err != nil {
            log.Printf("seed error: %v", err)
        }
    }
}
```

Key practices:
- **Hash passwords** before insertion — never store plain text.
- **Use `FirstOrCreate`** to make seeders safe to run multiple times.
- **Use realistic data** for development; avoid production secrets
  in seed files.

## 4. Seed Runner

The `RunAll` function calls all seeders in order:

```go
func RunAll(db *gorm.DB) {
    SeedUsers(db)
    // SeedPosts(db)
    // SeedCategories(db)
}
```

## 5. CLI Command

Run seeders via the CLI:

```text
framework db:seed
```

This resolves the database connection from the container and calls
`seeders.RunAll(db)`.

## 6. Security Considerations

- Seed data for admin accounts **MUST** use strong passwords, even
  in development.
- Seeder files **MUST NOT** be deployed to production unless
  specifically intended for initial data setup.
- Passwords in seed data **MUST** be hashed — never stored as
  plain text.

## 7. References

- [Database](database.md)
- [Models](models.md)
- [Migrations](migrations.md)
- [CLI Overview](../cli/cli-overview.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
