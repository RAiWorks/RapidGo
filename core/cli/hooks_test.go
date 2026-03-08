package cli

import (
	"testing"

	"github.com/RAiWorks/RapidGo/core/app"
	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/RAiWorks/RapidGo/core/scheduler"
	"github.com/RAiWorks/RapidGo/core/service"
	"gorm.io/gorm"
)

func TestHooksDefaultNil(t *testing.T) {
	// Reset to ensure clean state
	bootstrapFn = nil
	routeRegistrar = nil
	jobRegistrar = nil
	scheduleRegistrar = nil
	modelRegistryFn = nil
	seederFn = nil

	if bootstrapFn != nil {
		t.Error("bootstrapFn should default to nil")
	}
	if routeRegistrar != nil {
		t.Error("routeRegistrar should default to nil")
	}
	if jobRegistrar != nil {
		t.Error("jobRegistrar should default to nil")
	}
	if scheduleRegistrar != nil {
		t.Error("scheduleRegistrar should default to nil")
	}
	if modelRegistryFn != nil {
		t.Error("modelRegistryFn should default to nil")
	}
	if seederFn != nil {
		t.Error("seederFn should default to nil")
	}
}

func TestSetBootstrapStoresFunction(t *testing.T) {
	defer func() { bootstrapFn = nil }()

	called := false
	SetBootstrap(func(a *app.App, mode service.Mode) {
		called = true
	})

	if bootstrapFn == nil {
		t.Fatal("SetBootstrap did not store function")
	}
	bootstrapFn(nil, service.ModeAll)
	if !called {
		t.Error("stored bootstrap function was not called")
	}
}

func TestSetRoutesStoresFunction(t *testing.T) {
	defer func() { routeRegistrar = nil }()

	called := false
	SetRoutes(func(r *router.Router, c *container.Container, mode service.Mode) {
		called = true
	})

	if routeRegistrar == nil {
		t.Fatal("SetRoutes did not store function")
	}
	routeRegistrar(nil, nil, service.ModeAll)
	if !called {
		t.Error("stored route registrar was not called")
	}
}

func TestSetJobRegistrarStoresFunction(t *testing.T) {
	defer func() { jobRegistrar = nil }()

	called := false
	SetJobRegistrar(func() {
		called = true
	})

	if jobRegistrar == nil {
		t.Fatal("SetJobRegistrar did not store function")
	}
	jobRegistrar()
	if !called {
		t.Error("stored job registrar was not called")
	}
}

func TestSetScheduleRegistrarStoresFunction(t *testing.T) {
	defer func() { scheduleRegistrar = nil }()

	called := false
	SetScheduleRegistrar(func(s *scheduler.Scheduler, a *app.App) {
		called = true
	})

	if scheduleRegistrar == nil {
		t.Fatal("SetScheduleRegistrar did not store function")
	}
	scheduleRegistrar(nil, nil)
	if !called {
		t.Error("stored schedule registrar was not called")
	}
}

func TestSetModelRegistryStoresFunction(t *testing.T) {
	defer func() { modelRegistryFn = nil }()

	SetModelRegistry(func() []interface{} {
		return []interface{}{"test"}
	})

	if modelRegistryFn == nil {
		t.Fatal("SetModelRegistry did not store function")
	}
	result := modelRegistryFn()
	if len(result) != 1 {
		t.Errorf("modelRegistryFn returned %d items, want 1", len(result))
	}
}

func TestSetSeederStoresFunction(t *testing.T) {
	defer func() { seederFn = nil }()

	called := false
	SetSeeder(func(db *gorm.DB, name string) error {
		called = true
		return nil
	})

	if seederFn == nil {
		t.Fatal("SetSeeder did not store function")
	}
	if err := seederFn(nil, ""); err != nil {
		t.Errorf("seederFn returned error: %v", err)
	}
	if !called {
		t.Error("stored seeder function was not called")
	}
}
