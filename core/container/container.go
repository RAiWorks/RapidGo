package container

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// Factory is a function that creates a service instance.
type Factory func(c *Container) interface{}

// Container is the service container for dependency injection.
type Container struct {
	mu        sync.RWMutex
	bindings  map[string]Factory
	instances map[string]interface{}
}

// New creates a new empty container.
func New() *Container {
	return &Container{
		bindings:  make(map[string]Factory),
		instances: make(map[string]interface{}),
	}
}

// Bind registers a factory function for a service. Each Make() call
// invokes the factory and returns a new instance (transient).
func (c *Container) Bind(name string, factory Factory) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bindings[name] = factory
}

// Singleton registers a factory that is only called once.
// The first Make() call creates the instance; subsequent calls
// return the cached instance.
func (c *Container) Singleton(name string, factory Factory) {
	c.Bind(name, func(cont *Container) interface{} {
		cont.mu.RLock()
		if inst, ok := cont.instances[name]; ok {
			cont.mu.RUnlock()
			return inst
		}
		cont.mu.RUnlock()

		inst := factory(cont)
		cont.mu.Lock()
		// Double-check: another goroutine may have initialized while we waited for the lock.
		if existing, ok := cont.instances[name]; ok {
			cont.mu.Unlock()
			return existing
		}
		cont.instances[name] = inst
		cont.mu.Unlock()
		return inst
	})
}

// Instance registers an already-created instance directly.
func (c *Container) Instance(name string, instance interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances[name] = instance
}

// Make resolves a service by name. Checks instances first, then bindings.
// Panics if the service is not registered.
func (c *Container) Make(name string) interface{} {
	c.mu.RLock()
	if inst, ok := c.instances[name]; ok {
		c.mu.RUnlock()
		return inst
	}
	factory, ok := c.bindings[name]
	c.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("service not found: %s", name))
	}
	return factory(c)
}

// MustMake resolves a service and casts to the expected type.
// Panics if the service is not found or the type assertion fails.
func MustMake[T any](c *Container, name string) T {
	return c.Make(name).(T)
}

// TryMake resolves a service by name, returning an error instead of panicking
// if the service is not registered.
func (c *Container) TryMake(name string) (interface{}, error) {
	c.mu.RLock()
	if inst, ok := c.instances[name]; ok {
		c.mu.RUnlock()
		return inst, nil
	}
	factory, ok := c.bindings[name]
	c.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("service not found: %s", name)
	}
	return factory(c), nil
}

// TryMakeTyped resolves a service with type safety, returning an error
// instead of panicking if the service is not found or the type assertion fails.
func TryMake[T any](c *Container, name string) (T, error) {
	var zero T
	raw, err := c.TryMake(name)
	if err != nil {
		return zero, err
	}
	val, ok := raw.(T)
	if !ok {
		return zero, fmt.Errorf("service %s: type assertion failed", name)
	}
	return val, nil
}

// Has checks if a service is registered (binding or instance).
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, hasBinding := c.bindings[name]
	_, hasInstance := c.instances[name]
	return hasBinding || hasInstance
}

// SafeSingleton registers a singleton that catches panics during creation.
// If the factory panics and the app is in local/development mode, the panic
// is logged as a warning and the service is skipped (not registered).
// In production, the panic propagates normally.
func (c *Container) SafeSingleton(name string, factory Factory) {
	c.Singleton(name, func(cont *Container) interface{} {
		defer func() {
			if r := recover(); r != nil {
				if isLocalMode() {
					slog.Warn("service unavailable in local mode, skipping",
						"service", name,
						"error", fmt.Sprint(r),
					)
				} else {
					panic(r)
				}
			}
		}()
		return factory(cont)
	})
}

// isLocalMode checks if APP_ENV is local or development.
func isLocalMode() bool {
	env := os.Getenv("APP_ENV")
	return env == "local" || env == "development" || env == ""
}
