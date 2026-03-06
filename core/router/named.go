package router

import (
	"strings"
	"sync"
)

var (
	namedRoutes = make(map[string]string)
	namedMu     sync.RWMutex
)

// Name registers a route name mapped to a path pattern.
func Name(name, pattern string) {
	namedMu.Lock()
	defer namedMu.Unlock()
	namedRoutes[name] = pattern
}

// Route generates a URL from a named route with parameter substitution.
// Parameters replace :param placeholders in registration order.
// Returns "/" if the route name is not found.
//
// Example: Route("users.show", "42") → "/users/42"
func Route(name string, params ...string) string {
	namedMu.RLock()
	pattern, ok := namedRoutes[name]
	namedMu.RUnlock()
	if !ok {
		return "/"
	}

	result := pattern
	for i := 0; i < len(params); i++ {
		idx := strings.Index(result, ":")
		if idx == -1 {
			break
		}
		end := strings.IndexAny(result[idx:], "/")
		if end == -1 {
			result = result[:idx] + params[i]
		} else {
			result = result[:idx] + params[i] + result[idx+end:]
		}
	}
	return result
}

// ResetNamedRoutes clears all named routes. Used in tests only.
func ResetNamedRoutes() {
	namedMu.Lock()
	defer namedMu.Unlock()
	namedRoutes = make(map[string]string)
}
