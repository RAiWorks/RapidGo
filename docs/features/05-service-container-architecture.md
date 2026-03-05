# 🏗️ Architecture: Service Container

> **Feature**: `05` — Service Container
> **Discussion**: [`05-service-container-discussion.md`](05-service-container-discussion.md)
> **Status**: 🟢 FINALIZED
> **Date**: 2026-03-05

---

## Overview

The service container is a dependency injection (DI) mechanism with three registration patterns (`Bind`, `Singleton`, `Instance`), two resolution methods (`Make`, `MustMake[T]`), an existence check (`Has`), and a `Provider` interface for two-phase service lifecycle. The `App` struct in `core/app` orchestrates provider registration and boot. All operations are thread-safe via `sync.RWMutex`.

## File Structure

```
core/container/
├── container.go        # Container struct, Bind, Singleton, Instance, Make, MustMake, Has
├── provider.go         # Provider interface (Register, Boot)
└── container_test.go   # Unit tests for all container methods

core/app/
├── app.go              # App struct, New(), Register(), Boot(), Make()
└── app_test.go         # Unit tests for App bootstrap lifecycle
```

No existing files are modified.

## Component Design

### Container (`core/container/container.go`)

**Responsibility**: Dependency injection container — register factories, resolve services
**Package**: `core/container`

```go
package container

import (
	"fmt"
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

// Has checks if a service is registered (binding or instance).
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, hasBinding := c.bindings[name]
	_, hasInstance := c.instances[name]
	return hasBinding || hasInstance
}
```

### Provider Interface (`core/container/provider.go`)

**Responsibility**: Define the two-phase lifecycle contract for service providers
**Package**: `core/container`

```go
package container

// Provider defines the lifecycle hooks for service registration.
type Provider interface {
	// Register binds services into the container.
	// Called before the application boots. Only register bindings here.
	Register(c *Container)

	// Boot runs after ALL providers have been registered.
	// May resolve other services from the container.
	Boot(c *Container)
}
```

### App Struct (`core/app/app.go`)

**Responsibility**: Application bootstrap — manage provider registration order and boot sequence
**Package**: `core/app`

```go
package app

import "github.com/RAiWorks/RGo/core/container"

// App is the application container that manages providers and services.
type App struct {
	Container *container.Container
	providers []container.Provider
}

// New creates a new application with an empty container.
func New() *App {
	return &App{
		Container: container.New(),
	}
}

// Register adds a provider and immediately calls its Register method.
func (a *App) Register(provider container.Provider) {
	a.providers = append(a.providers, provider)
	provider.Register(a.Container)
}

// Boot calls Boot on all registered providers in registration order.
func (a *App) Boot() {
	for _, p := range a.providers {
		p.Boot(a.Container)
	}
}

// Make resolves a service from the container.
func (a *App) Make(name string) interface{} {
	return a.Container.Make(name)
}
```

## Data Flow

```
main.go
  → app.New()                          creates Container
  → app.Register(&SomeProvider{})      calls provider.Register(container) → binds factories
  → app.Register(&AnotherProvider{})   same
  → app.Boot()                         calls provider.Boot(container) on all, in order
  → app.Make("service")               resolves service from container
```

### Resolution Order in Make()

```
Make("name")
  → Check instances map → found? return it
  → Check bindings map → found? call factory(container), return result
  → Neither? → panic("service not found: name")
```

### Singleton Lifecycle

```
First Make("db")
  → Not in instances → call factory → store in instances → return
Second Make("db")
  → Found in instances → return cached instance
```

## Configuration

No environment variables needed. The container is pure Go with no external configuration.

## Security Considerations

- Service instances live in-memory only — not exposed beyond the Go process
- Provider factories that access credentials should read from env vars inside closures, not store as struct fields
- `Make()` panic on missing service is intentional fail-fast — all services must be registered before boot

## Trade-offs & Alternatives

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| String-keyed container with `interface{}` | Simple, matches blueprint, works with generics via `MustMake[T]` | No compile-time type safety for `Make()` | ✅ Selected |
| Struct-field-based DI (wire-style) | Compile-time safe, no runtime panics | Requires code generation, over-engineered for this framework | ❌ Too complex |
| Global singletons (package-level vars) | Simplest possible approach | No testability, no swapping, tight coupling | ❌ Not extensible |
| Interface-keyed container (`reflect.Type`) | Type-safe keys | Reflection overhead, non-idiomatic, complex API | ❌ Over-engineered |

## Next

Create tasks doc → `05-service-container-tasks.md`
