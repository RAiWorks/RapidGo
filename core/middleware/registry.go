package middleware

import "github.com/gin-gonic/gin"

var (
	routeMiddleware  = map[string]gin.HandlerFunc{}
	middlewareGroups = map[string][]gin.HandlerFunc{}
)

// RegisterAlias registers a named middleware that can be referenced by string.
func RegisterAlias(name string, handler gin.HandlerFunc) {
	routeMiddleware[name] = handler
}

// RegisterGroup registers a named group of middleware.
func RegisterGroup(name string, handlers ...gin.HandlerFunc) {
	middlewareGroups[name] = handlers
}

// Resolve returns a middleware handler by alias name.
// Panics if the alias is not registered.
func Resolve(name string) gin.HandlerFunc {
	if h, ok := routeMiddleware[name]; ok {
		return h
	}
	panic("middleware not found: " + name)
}

// ResolveGroup returns all middleware in a named group.
// Returns nil if the group is not registered.
func ResolveGroup(name string) []gin.HandlerFunc {
	if g, ok := middlewareGroups[name]; ok {
		return g
	}
	return nil
}

// ResetRegistry clears all registered aliases and groups. For testing only.
func ResetRegistry() {
	routeMiddleware = map[string]gin.HandlerFunc{}
	middlewareGroups = map[string][]gin.HandlerFunc{}
}
