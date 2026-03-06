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
