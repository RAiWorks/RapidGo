---
title: "Sessions"
version: "0.1.0"
status: "Draft"
date: "2026-03-05"
last_updated: "2026-03-05"
authors:
  - "RAiWorks"
supersedes: ""
---

# Sessions

## Abstract

This document covers the session management system — the driver-based
architecture, all four storage backends (database, Redis, file,
memory), the session manager, middleware, store factory, and flash
messages.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Architecture](#2-architecture)
3. [Configuration](#3-configuration)
4. [Store Interface](#4-store-interface)
5. [Database Store](#5-database-store)
6. [Redis Store](#6-redis-store)
7. [File Store](#7-file-store)
8. [Memory Store](#8-memory-store)
9. [Store Factory](#9-store-factory)
10. [Session Manager](#10-session-manager)
11. [Session Middleware](#11-session-middleware)
12. [Flash Messages](#12-flash-messages)
13. [Security Considerations](#13-security-considerations)
14. [References](#14-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Session** — Server-side storage of per-user state across HTTP
  requests.
- **Flash message** — A one-time session value available only on the
  next request.

## 2. Architecture

The session system uses a **driver-based architecture**. All backends
implement the `Store` interface, making them interchangeable via the
`SESSION_DRIVER` configuration.

```text
Session Middleware
       ↓
Session Manager (cookie handling, ID generation)
       ↓
Store Interface
       ↓
┌──────────┬──────────┬──────────┬──────────┐
│ DBStore  │ Redis    │ File     │ Memory   │
│          │ Store    │ Store    │ Store    │
└──────────┴──────────┴──────────┴──────────┘
```

## 3. Configuration

`.env` session variables:

```env
SESSION_DRIVER=db
SESSION_LIFETIME=120
SESSION_COOKIE=framework_session
SESSION_PATH=/
SESSION_DOMAIN=
SESSION_SECURE=false
SESSION_HTTPONLY=true
SESSION_SAMESITE=lax

# Redis driver only
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# File driver only
SESSION_FILE_PATH=storage/sessions
```

### Driver Comparison

| Driver | Use Case | Config Value |
|--------|----------|-------------|
| **Database** | Production default, shared across instances | `db` |
| **Redis** | High-performance, distributed | `redis` |
| **File** | Simple deployments, single server | `file` |
| **Memory** | Development & testing only | `memory` |

## 4. Store Interface

All backends implement this contract:

```go
package session

import "time"

type Store interface {
    Read(id string) (map[string]interface{}, error)
    Write(id string, data map[string]interface{}, lifetime time.Duration) error
    Destroy(id string) error
    GC(maxLifetime time.Duration) error
}
```

## 5. Database Store

Requires a `sessions` table (auto-created via migration):

```sql
CREATE TABLE sessions (
    id          VARCHAR(255) PRIMARY KEY,
    data        TEXT         NOT NULL,
    user_id     BIGINT       NULL,
    ip_address  VARCHAR(45)  NULL,
    user_agent  TEXT         NULL,
    last_active TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sessions_last_active ON sessions (last_active);
```

```go
type DBStore struct {
    DB *gorm.DB
}

type SessionRecord struct {
    ID         string    `gorm:"primaryKey;size:255"`
    Data       string    `gorm:"type:text;not null"`
    UserID     *uint     `gorm:"index"`
    IPAddress  *string   `gorm:"size:45"`
    UserAgent  *string   `gorm:"type:text"`
    LastActive time.Time `gorm:"autoUpdateTime"`
    CreatedAt  time.Time
}

func (SessionRecord) TableName() string { return "sessions" }
```

The database store serializes session data as JSON. GC deletes records
where `last_active` is older than the max lifetime.

## 6. Redis Store

Uses Redis key expiration for automatic cleanup:

```go
type RedisStore struct {
    Client *redis.Client
    Prefix string // e.g. "session:"
}
```

Redis TTL replaces manual GC — expired sessions are automatically
removed by Redis.

## 7. File Store

Stores one JSON file per session in the configured directory:

```go
type FileStore struct {
    Path string // e.g. "storage/sessions"
}
```

Session files use `0600` permissions for security. GC removes files
older than the max lifetime based on modification time.

## 8. Memory Store

For development and testing only — data is lost on restart:

```go
type MemoryStore struct {
    mu       sync.RWMutex
    sessions map[string]memEntry
}
```

## 9. Store Factory

Resolves the correct backend from `SESSION_DRIVER`:

```go
func NewStore(db *gorm.DB) (Store, error) {
    driver := os.Getenv("SESSION_DRIVER")

    switch driver {
    case "db":
        db.AutoMigrate(&SessionRecord{})
        return &DBStore{DB: db}, nil
    case "file":
        path := os.Getenv("SESSION_FILE_PATH")
        if path == "" {
            path = "storage/sessions"
        }
        return &FileStore{Path: path}, nil
    case "redis":
        redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
        client := redis.NewClient(&redis.Options{
            Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
            Password: os.Getenv("REDIS_PASSWORD"),
            DB:       redisDB,
        })
        return &RedisStore{Client: client, Prefix: "session:"}, nil
    case "memory":
        return NewMemoryStore(), nil
    default:
        return nil, fmt.Errorf("unsupported SESSION_DRIVER: %s", driver)
    }
}
```

## 10. Session Manager

Ties together the store, cookie handling, and ID generation:

```go
type Manager struct {
    Store      Store
    CookieName string
    Lifetime   time.Duration
    Path       string
    Domain     string
    Secure     bool
    HTTPOnly   bool
    SameSite   http.SameSite
}
```

Key methods:

| Method | Description |
|--------|-------------|
| `Start(r)` | Retrieve or create a session from the request cookie |
| `Save(w, id, data)` | Persist session data and write the cookie |
| `Destroy(w, id)` | Delete session and clear the cookie |

Session IDs are generated using 32 bytes from `crypto/rand`.

## 11. Session Middleware

Automatically loads and saves sessions per request:

```go
func SessionMiddleware(mgr *session.Manager) gin.HandlerFunc {
    return func(c *gin.Context) {
        id, data, err := mgr.Start(c.Request)
        if err != nil {
            c.AbortWithStatus(500)
            return
        }
        c.Set("session_id", id)
        c.Set("session", data)

        c.Next()

        updated, _ := c.Get("session")
        mgr.Save(c.Writer, id, updated.(map[string]interface{}))
    }
}
```

Access session data in controllers:

```go
sess, _ := c.Get("session")
data := sess.(map[string]interface{})

// Read
username, _ := data["username"].(string)

// Write
data["last_visited"] = time.Now().Format(time.RFC3339)
c.Set("session", data)
```

## 12. Flash Messages

One-time session data used for passing status messages across
redirects:

### Setting Flash Messages

```go
sessionMgr.Flash(data, "success", "User created successfully!")
```

### Reading Flash Messages

```go
success, _ := sessionMgr.GetFlash(data, "success")
// The flash is automatically removed after reading
```

### Validation Errors + Old Input

```go
// On validation failure — redirect with errors and old input
sessionMgr.FlashErrors(data, v.Errors())
sessionMgr.FlashOldInput(data, map[string]string{
    "name":  c.PostForm("name"),
    "email": c.PostForm("email"),
})
c.Set("session", data)
c.Redirect(http.StatusFound, "/users/create")
```

### Reading in the Next Request

```go
errors, _ := sessionMgr.GetFlash(data, "_errors")
oldInput, _ := sessionMgr.GetFlash(data, "_old_input")
success, _ := sessionMgr.GetFlash(data, "success")

c.Set("session", data) // save cleared flashes
c.HTML(200, "users/create.html", gin.H{
    "errors":  errors,
    "old":     oldInput,
    "success": success,
})
```

## 13. Security Considerations

- Session cookies **MUST** use `HttpOnly` and `Secure` flags in
  production.
- `SameSite` **SHOULD** be set to `Lax` or `Strict` to prevent
  CSRF via session cookies.
- Session IDs **MUST** be generated using `crypto/rand`, not
  `math/rand`.
- File store **MUST** use restrictive permissions (`0600` for files,
  `0700` for directories).
- Sessions **SHOULD** be regenerated after authentication to
  prevent session fixation attacks.
- Redis and database connections for session stores **MUST** use
  authentication in production.

## 14. References

- [Authentication](authentication.md)
- [CSRF Protection](csrf.md)
- [Middleware](../http/middleware.md)
- [Configuration](../core/configuration.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
