package middleware

import (
	"crypto/rand"
	"fmt"

	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-ID"

// RequestID returns middleware that assigns a unique ID to each request.
// If the request already has an X-Request-ID header, it is preserved.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = generateUUID()
		}
		c.Set("request_id", id)
		c.Header(requestIDHeader, id)
		c.Next()
	}
}

// generateUUID produces a UUID v4 string using crypto/rand.
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
