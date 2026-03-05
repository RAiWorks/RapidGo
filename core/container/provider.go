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
