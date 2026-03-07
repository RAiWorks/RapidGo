package providers

import (
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/middleware"
	"github.com/RAiWorks/RapidGo/core/service"
)

// MiddlewareProvider registers built-in middleware aliases and groups.
type MiddlewareProvider struct {
	Mode service.Mode
}

// Register is a no-op — middleware has no singleton to register.
func (p *MiddlewareProvider) Register(c *container.Container) {}

// Boot registers middleware aliases relevant to the current mode.
func (p *MiddlewareProvider) Boot(c *container.Container) {
	// Always register — universally useful
	middleware.RegisterAlias("recovery", middleware.Recovery())
	middleware.RegisterAlias("requestid", middleware.RequestID())
	middleware.RegisterAlias("cors", middleware.CORS())
	middleware.RegisterAlias("error_handler", middleware.ErrorHandler())
	middleware.RegisterAlias("ratelimit", middleware.RateLimitMiddleware())

	// Web-only middleware
	if p.Mode.Has(service.ModeWeb) {
		middleware.RegisterAlias("csrf", middleware.CSRFMiddleware())
	}

	// Auth — needed by both web (session) and api (JWT)
	middleware.RegisterAlias("auth", middleware.AuthMiddleware())

	middleware.RegisterGroup("global",
		middleware.Recovery(),
		middleware.RequestID(),
	)
}
