package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds configuration for the CORS middleware.
type CORSConfig struct {
	AllowOrigins []string // Default: ["*"]
	AllowMethods []string // Default: ["GET","POST","PUT","DELETE","PATCH","OPTIONS"]
	AllowHeaders []string // Default: ["Origin","Content-Type","Accept","Authorization","X-Request-ID"]
	MaxAge       int      // Default: 43200 (12 hours), in seconds
}

// defaultCORSConfig returns the default CORS configuration.
func defaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		MaxAge:       43200,
	}
}

// CORS returns middleware that handles cross-origin requests.
// If no config is provided, sensible defaults are used.
func CORS(configs ...CORSConfig) gin.HandlerFunc {
	cfg := defaultCORSConfig()
	if len(configs) > 0 {
		cfg = configs[0]
	}

	origins := strings.Join(cfg.AllowOrigins, ", ")
	methods := strings.Join(cfg.AllowMethods, ", ")
	headers := strings.Join(cfg.AllowHeaders, ", ")
	maxAge := fmt.Sprintf("%d", cfg.MaxAge)

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Methods", methods)
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Access-Control-Max-Age", maxAge)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
