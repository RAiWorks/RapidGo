package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware generates and validates per-session CSRF tokens.
// Safe methods (GET, HEAD, OPTIONS) are skipped.
// State-changing methods require a valid token in the _csrf_token form
// field or X-CSRF-Token header.
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, _ := c.Get("session")
		data := sess.(map[string]interface{})

		// Generate token if not present
		token, ok := data["_csrf_token"].(string)
		if !ok || token == "" {
			b := make([]byte, 32)
			rand.Read(b)
			token = hex.EncodeToString(b)
			data["_csrf_token"] = token
			c.Set("session", data)
		}

		// Make token available to templates
		c.Set("csrf_token", token)

		// Skip safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Validate token on POST/PUT/PATCH/DELETE
		submitted := c.PostForm("_csrf_token")
		if submitted == "" {
			submitted = c.GetHeader("X-CSRF-Token")
		}
		if submitted != token {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token mismatch",
			})
			return
		}

		c.Next()
	}
}
