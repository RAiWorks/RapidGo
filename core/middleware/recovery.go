package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery returns middleware that catches panics and returns a 500 error.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					"error", err,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				)

				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "internal server error",
					})
				} else {
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
