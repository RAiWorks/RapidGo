package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raiworks/rapidgo/v2/core/router"
	"gorm.io/gorm"
)

// Routes registers liveness and readiness health-check endpoints.
// The dbFn callback defers database resolution until the first request.
func Routes(r *router.Router, dbFn func() *gorm.DB) {
	r.Get("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Get("/health/ready", func(c *gin.Context) {
		sqlDB, err := dbFn().DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": err.Error()})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready", "db": "connected"})
	})
}
