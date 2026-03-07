package providers

import (
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/session"
	"gorm.io/gorm"
)

// SessionProvider registers the session manager.
type SessionProvider struct{}

// Register creates a session manager singleton from the configured session driver.
func (p *SessionProvider) Register(c *container.Container) {
	c.Singleton("session", func(c *container.Container) interface{} {
		db := container.MustMake[*gorm.DB](c, "db")
		store, _ := session.NewStore(db)
		return session.NewManager(store)
	})
}

// Boot is a no-op for SessionProvider.
func (p *SessionProvider) Boot(c *container.Container) {}
