package app

import "github.com/RAiWorks/RapidGo/v2/core/container"

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
