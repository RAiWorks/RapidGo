package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raiworks/rapidgo/v2/core/session"
)

// SessionMiddleware automatically loads/saves sessions per request.
// The session cookie is set before c.Next() so it is included in response
// headers even when handlers write the body (e.g. c.HTML()).
func SessionMiddleware(mgr *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, data, err := mgr.Start(c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Set("session_id", id)
		c.Set("session", data)

		// Set the session cookie BEFORE the handler writes the response body.
		mgr.SetCookie(c.Writer, id)

		c.Next()

		// Persist session data to the store after the handler runs.
		updated, _ := c.Get("session")
		mgr.Store.Write(id, updated.(map[string]interface{}), mgr.Lifetime)
	}
}
