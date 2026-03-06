package router

import "github.com/gin-gonic/gin"

// RouteGroup wraps a Gin router group for sub-route registration.
type RouteGroup struct {
	group *gin.RouterGroup
}

// Get registers a GET route in the group.
func (g *RouteGroup) Get(path string, handlers ...gin.HandlerFunc) {
	g.group.GET(path, handlers...)
}

// Post registers a POST route in the group.
func (g *RouteGroup) Post(path string, handlers ...gin.HandlerFunc) {
	g.group.POST(path, handlers...)
}

// Put registers a PUT route in the group.
func (g *RouteGroup) Put(path string, handlers ...gin.HandlerFunc) {
	g.group.PUT(path, handlers...)
}

// Delete registers a DELETE route in the group.
func (g *RouteGroup) Delete(path string, handlers ...gin.HandlerFunc) {
	g.group.DELETE(path, handlers...)
}

// Patch registers a PATCH route in the group.
func (g *RouteGroup) Patch(path string, handlers ...gin.HandlerFunc) {
	g.group.PATCH(path, handlers...)
}

// Options registers an OPTIONS route in the group.
func (g *RouteGroup) Options(path string, handlers ...gin.HandlerFunc) {
	g.group.OPTIONS(path, handlers...)
}

// Group creates a nested sub-group with a shared prefix and optional middleware.
func (g *RouteGroup) Group(prefix string, handlers ...gin.HandlerFunc) *RouteGroup {
	return &RouteGroup{group: g.group.Group(prefix, handlers...)}
}

// Use adds middleware to the group.
func (g *RouteGroup) Use(middleware ...gin.HandlerFunc) {
	g.group.Use(middleware...)
}
