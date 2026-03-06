package providers

import (
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/router"
	"github.com/RAiWorks/RGo/routes"
)

// RouterProvider creates the router and registers route definitions.
type RouterProvider struct{}

// Register creates a new Router and registers it as "router" in the container.
func (p *RouterProvider) Register(c *container.Container) {
	c.Instance("router", router.New())
}

// Boot loads route definitions from the routes package.
func (p *RouterProvider) Boot(c *container.Container) {
	r := container.MustMake[*router.Router](c, "router")
	routes.RegisterWeb(r)
	routes.RegisterAPI(r)
}
