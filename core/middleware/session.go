package middleware

import (
	"github.com/RAiWorks/RapidGo/v2/core/session"
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

		// Persist session after the handler runs
		updated, _ := c.Get("session")
		mgr.Save(c.Writer, id, updated.(map[string]interface{}))
	}
}
