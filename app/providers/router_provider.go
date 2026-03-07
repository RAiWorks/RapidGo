package providers

import (
	"os"
	"path/filepath"

	"github.com/RAiWorks/RapidGo/core/config"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/health"
	"github.com/RAiWorks/RapidGo/core/metrics"
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/RAiWorks/RapidGo/core/service"
	"github.com/RAiWorks/RapidGo/routes"
	"gorm.io/gorm"
)

// RouterProvider creates the router and registers route definitions.
type RouterProvider struct {
	Mode service.Mode
}

// Register creates a new Router and registers it as "router" in the container.
func (p *RouterProvider) Register(c *container.Container) {
	c.Instance("router", router.New())
}

// Boot sets up templates, static serving, and loads route definitions based on mode.
func (p *RouterProvider) Boot(c *container.Container) {
	r := container.MustMake[*router.Router](c, "router")

	// Template engine and static serving — only for web mode
	if p.Mode.Has(service.ModeWeb) {
		r.SetFuncMap(router.DefaultFuncMap())
		viewsDir := filepath.Join("resources", "views")
		if info, err := os.Stat(viewsDir); err == nil && info.IsDir() {
			r.LoadTemplates(viewsDir)
		}
		if info, err := os.Stat("resources/static"); err == nil && info.IsDir() {
			r.Static("/static", "./resources/static")
		}
		if info, err := os.Stat("storage/uploads"); err == nil && info.IsDir() {
			r.Static("/uploads", "./storage/uploads")
		}
	}

	// Route definitions — conditional on mode
	if p.Mode.Has(service.ModeWeb) {
		routes.RegisterWeb(r)
	}
	if p.Mode.Has(service.ModeAPI) {
		routes.RegisterAPI(r)
	}
	if p.Mode.Has(service.ModeWS) {
		routes.RegisterWS(r)
	}

	// Health check — available in any HTTP mode when DB is present
	if c.Has("db") {
		health.Routes(r, func() *gorm.DB {
			return container.MustMake[*gorm.DB](c, "db")
		})
	}

	// Prometheus metrics — opt-in via METRICS_ENABLED
	if config.EnvBool("METRICS_ENABLED", false) {
		m := metrics.New()
		r.Use(m.Middleware())
		r.Get(config.Env("METRICS_PATH", "/metrics"), metrics.Handler())
	}
}
