# 🏗️ Architecture: Session Management

> **Feature**: `20` — Session Management
> **Discussion**: [`20-session-management-discussion.md`](20-session-management-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-06

---

## Overview

A driver-based session system: `Store` interface → four backends (Memory, File, DB, Cookie) → `Manager` (cookie handling, ID gen, flash) → `SessionMiddleware` → `SessionProvider`.

## File Structure

```
core/session/
├── store.go            # Store interface
├── memory_store.go     # MemoryStore (dev/testing)
├── file_store.go       # FileStore (single-server)
├── db_store.go         # DBStore + SessionRecord (production)
├── cookie_store.go     # CookieStore (AES-256-GCM encrypted)
├── factory.go          # NewStore factory
├── manager.go          # Manager + flash messages
└── session_test.go     # All tests

core/middleware/
└── session.go          # SessionMiddleware

app/providers/
└── session_provider.go # SessionProvider
```

### Files Created (10)
| File | Package | Lines (est.) |
|---|---|---|
| `core/session/store.go` | `session` | ~15 |
| `core/session/memory_store.go` | `session` | ~50 |
| `core/session/file_store.go` | `session` | ~55 |
| `core/session/db_store.go` | `session` | ~65 |
| `core/session/cookie_store.go` | `session` | ~80 |
| `core/session/factory.go` | `session` | ~35 |
| `core/session/manager.go` | `session` | ~120 |
| `core/session/session_test.go` | `session` | ~200 |
| `core/middleware/session.go` | `middleware` | ~25 |
| `app/providers/session_provider.go` | `providers` | ~25 |

### Files Modified (1)
| File | Change |
|---|---|
| `core/cli/root.go` | Add `SessionProvider` registration after `DatabaseProvider` |

---

## Component Design

### Store Interface (`core/session/store.go`)

```go
package session

import "time"

// Store defines the contract every session backend must satisfy.
type Store interface {
	Read(id string) (map[string]interface{}, error)
	Write(id string, data map[string]interface{}, lifetime time.Duration) error
	Destroy(id string) error
	GC(maxLifetime time.Duration) error
}
```

### Memory Store (`core/session/memory_store.go`)

```go
package session

import (
	"sync"
	"time"
)

type memEntry struct {
	Data      map[string]interface{}
	ExpiresAt time.Time
}

type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]memEntry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{sessions: make(map[string]memEntry)}
}

func (s *MemoryStore) Read(id string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.sessions[id]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return make(map[string]interface{}), nil
	}
	return entry.Data, nil
}

func (s *MemoryStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = memEntry{Data: data, ExpiresAt: time.Now().Add(lifetime)}
	return nil
}

func (s *MemoryStore) Destroy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}

func (s *MemoryStore) GC(maxLifetime time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for id, entry := range s.sessions {
		if now.After(entry.ExpiresAt) {
			delete(s.sessions, id)
		}
	}
	return nil
}
```

### File Store (`core/session/file_store.go`)

```go
package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type FileStore struct {
	Path string
}

func (s *FileStore) filepath(id string) string {
	return filepath.Join(s.Path, id+".json")
}

func (s *FileStore) Read(id string) (map[string]interface{}, error) {
	raw, err := os.ReadFile(s.filepath(id))
	if err != nil {
		return make(map[string]interface{}), nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *FileStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	if err := os.MkdirAll(s.Path, 0700); err != nil {
		return err
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filepath(id), raw, 0600)
}

func (s *FileStore) Destroy(id string) error {
	return os.Remove(s.filepath(id))
}

func (s *FileStore) GC(maxLifetime time.Duration) error {
	entries, err := os.ReadDir(s.Path)
	if err != nil {
		return err
	}
	cutoff := time.Now().Add(-maxLifetime)
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(s.Path, e.Name()))
		}
	}
	return nil
}
```

### Database Store (`core/session/db_store.go`)

```go
package session

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

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

func (s *DBStore) Read(id string) (map[string]interface{}, error) {
	var rec SessionRecord
	if err := s.DB.Where("id = ?", id).First(&rec).Error; err != nil {
		return make(map[string]interface{}), nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(rec.Data), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DBStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	rec := SessionRecord{
		ID:         id,
		Data:       string(raw),
		LastActive: time.Now(),
	}
	return s.DB.Save(&rec).Error
}

func (s *DBStore) Destroy(id string) error {
	return s.DB.Where("id = ?", id).Delete(&SessionRecord{}).Error
}

func (s *DBStore) GC(maxLifetime time.Duration) error {
	cutoff := time.Now().Add(-maxLifetime)
	return s.DB.Where("last_active < ?", cutoff).Delete(&SessionRecord{}).Error
}
```

### Cookie Store (`core/session/cookie_store.go`)

```go
package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"
)

// CookieStore stores session data encrypted in the cookie value itself.
// Requires a 32-byte key (AES-256-GCM).
type CookieStore struct {
	Key  []byte
	data map[string]string // id → encrypted payload (in-memory mirror for current request)
}

func NewCookieStore(key []byte) (*CookieStore, error) {
	if len(key) != 32 {
		return nil, errors.New("cookie store requires a 32-byte key")
	}
	return &CookieStore{Key: key, data: make(map[string]string)}, nil
}

func (s *CookieStore) encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *CookieStore) decrypt(encoded string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	return gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
}

func (s *CookieStore) Read(id string) (map[string]interface{}, error) {
	encoded, ok := s.data[id]
	if !ok {
		return make(map[string]interface{}), nil
	}
	plaintext, err := s.decrypt(encoded)
	if err != nil {
		return make(map[string]interface{}), nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CookieStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	encrypted, err := s.encrypt(raw)
	if err != nil {
		return err
	}
	s.data[id] = encrypted
	return nil
}

func (s *CookieStore) Destroy(id string) error {
	delete(s.data, id)
	return nil
}

func (s *CookieStore) GC(maxLifetime time.Duration) error {
	return nil // Expiry handled by cookie MaxAge
}
```

### Store Factory (`core/session/factory.go`)

```go
package session

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

// NewStore resolves the correct session backend from SESSION_DRIVER.
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
	case "memory", "":
		return NewMemoryStore(), nil
	case "cookie":
		key := []byte(os.Getenv("APP_KEY"))
		store, err := NewCookieStore(key)
		if err != nil {
			return nil, fmt.Errorf("cookie session store: %w", err)
		}
		return store, nil
	case "redis":
		return nil, fmt.Errorf("redis session driver requires github.com/redis/go-redis/v9 — not yet included")
	default:
		return nil, fmt.Errorf("unsupported SESSION_DRIVER: %s", driver)
	}
}
```

### Session Manager (`core/session/manager.go`)

```go
package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"time"
)

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

func NewManager(store Store) *Manager {
	lifetime, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))
	if lifetime == 0 {
		lifetime = 120
	}
	secure := os.Getenv("SESSION_SECURE") == "true"
	httpOnly := os.Getenv("SESSION_HTTPONLY") != "false"

	sameSite := http.SameSiteLaxMode
	switch os.Getenv("SESSION_SAMESITE") {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}

	cookieName := os.Getenv("SESSION_COOKIE")
	if cookieName == "" {
		cookieName = "framework_session"
	}

	return &Manager{
		Store:      store,
		CookieName: cookieName,
		Lifetime:   time.Duration(lifetime) * time.Minute,
		Path:       os.Getenv("SESSION_PATH"),
		Domain:     os.Getenv("SESSION_DOMAIN"),
		Secure:     secure,
		HTTPOnly:   httpOnly,
		SameSite:   sameSite,
	}
}

func (m *Manager) generateID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (m *Manager) Start(r *http.Request) (string, map[string]interface{}, error) {
	cookie, err := r.Cookie(m.CookieName)
	if err == nil && cookie.Value != "" {
		data, err := m.Store.Read(cookie.Value)
		if err != nil {
			return "", nil, err
		}
		if len(data) > 0 {
			return cookie.Value, data, nil
		}
	}
	id := m.generateID()
	return id, make(map[string]interface{}), nil
}

func (m *Manager) Save(w http.ResponseWriter, id string, data map[string]interface{}) error {
	if err := m.Store.Write(id, data, m.Lifetime); err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     m.CookieName,
		Value:    id,
		Path:     m.Path,
		Domain:   m.Domain,
		MaxAge:   int(m.Lifetime.Seconds()),
		Secure:   m.Secure,
		HttpOnly: m.HTTPOnly,
		SameSite: m.SameSite,
	})
	return nil
}

func (m *Manager) Destroy(w http.ResponseWriter, id string) error {
	if err := m.Store.Destroy(id); err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:   m.CookieName,
		Value:  "",
		Path:   m.Path,
		Domain: m.Domain,
		MaxAge: -1,
	})
	return nil
}

// Flash writes a flash message that will be available on the next request only.
func (m *Manager) Flash(data map[string]interface{}, key string, value interface{}) {
	flashes, _ := data["_flashes"].(map[string]interface{})
	if flashes == nil {
		flashes = make(map[string]interface{})
	}
	flashes[key] = value
	data["_flashes"] = flashes
}

// GetFlash reads and removes a flash message.
func (m *Manager) GetFlash(data map[string]interface{}, key string) (interface{}, bool) {
	flashes, _ := data["_flashes"].(map[string]interface{})
	if flashes == nil {
		return nil, false
	}
	val, ok := flashes[key]
	if ok {
		delete(flashes, key)
		if len(flashes) == 0 {
			delete(data, "_flashes")
		} else {
			data["_flashes"] = flashes
		}
	}
	return val, ok
}

// FlashErrors stores validation errors for the next request.
func (m *Manager) FlashErrors(data map[string]interface{}, errors map[string][]string) {
	m.Flash(data, "_errors", errors)
}

// FlashOldInput stores form input for re-populating after validation failure.
func (m *Manager) FlashOldInput(data map[string]interface{}, input map[string]string) {
	m.Flash(data, "_old_input", input)
}
```

### Session Middleware (`core/middleware/session.go`)

```go
package middleware

import (
	"github.com/RAiWorks/RGo/core/session"
	"github.com/gin-gonic/gin"
)

// SessionMiddleware automatically loads/saves sessions per request.
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

### Session Provider (`app/providers/session_provider.go`)

```go
package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/session"
	"gorm.io/gorm"
)

// SessionProvider registers the session manager.
type SessionProvider struct{}

func (p *SessionProvider) Register(c *container.Container) {
	c.Singleton("session", func(c *container.Container) interface{} {
		db := container.MustMake[*gorm.DB](c, "db")
		store, _ := session.NewStore(db)
		return session.NewManager(store)
	})
}

func (p *SessionProvider) Boot(c *container.Container) {}
```

---

## Dependencies

| Dependency | Type | Usage |
|---|---|---|
| `gorm.io/gorm` | existing | `DBStore` backend |
| `crypto/aes`, `crypto/cipher`, `crypto/rand` | stdlib | `CookieStore` encryption |
| `encoding/json` | stdlib | Session data serialization |
| `encoding/hex`, `encoding/base64` | stdlib | ID generation, cookie encoding |
| `net/http` | stdlib | Cookie management |
| `github.com/gin-gonic/gin` | existing | Session middleware |

No new external dependencies.

---

## Impact on Existing Code

| Component | Impact |
|---|---|
| `core/session/` | `.gitkeep` stays; 8 new files |
| `core/middleware/` | 1 new file (`session.go`) |
| `app/providers/` | 1 new file (`session_provider.go`) |
| `core/cli/root.go` | Add `SessionProvider` registration after `DatabaseProvider` |
