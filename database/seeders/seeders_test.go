package seeders

import (
	"errors"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	return db
}

// mockSeeder tracks whether Seed was called.
type mockSeeder struct {
	name   string
	called bool
	err    error
}

func (m *mockSeeder) Name() string         { return m.name }
func (m *mockSeeder) Seed(db *gorm.DB) error { m.called = true; return m.err }

// TC-01: Register adds a seeder to the registry.
func TestRegister_AddsSeeder(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	Register(&mockSeeder{name: "test"})

	names := Names()
	if len(names) != 1 {
		t.Fatalf("expected 1 seeder, got %d", len(names))
	}
	if names[0] != "test" {
		t.Fatalf("expected name 'test', got %q", names[0])
	}
}

// TC-02: RunAll executes all registered seeders.
func TestRunAll_ExecutesSeeders(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	a := &mockSeeder{name: "a"}
	b := &mockSeeder{name: "b"}
	Register(a)
	Register(b)

	db := setupTestDB(t)
	if err := RunAll(db); err != nil {
		t.Fatalf("RunAll failed: %v", err)
	}
	if !a.called {
		t.Fatal("expected seeder 'a' to be called")
	}
	if !b.called {
		t.Fatal("expected seeder 'b' to be called")
	}
}

// TC-03: RunAll stops on first error.
func TestRunAll_StopsOnError(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	a := &mockSeeder{name: "a"}
	b := &mockSeeder{name: "b", err: errors.New("seed failed")}
	Register(a)
	Register(b)

	db := setupTestDB(t)
	err := RunAll(db)
	if err == nil {
		t.Fatal("expected error from RunAll")
	}
	if !strings.Contains(err.Error(), "seed failed") {
		t.Fatalf("expected error to contain 'seed failed', got %q", err.Error())
	}
}

// TC-04: RunByName executes the named seeder only.
func TestRunByName_FindsSeeder(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	a := &mockSeeder{name: "first"}
	b := &mockSeeder{name: "second"}
	Register(a)
	Register(b)

	db := setupTestDB(t)
	if err := RunByName(db, "second"); err != nil {
		t.Fatalf("RunByName failed: %v", err)
	}
	if a.called {
		t.Fatal("expected seeder 'first' NOT to be called")
	}
	if !b.called {
		t.Fatal("expected seeder 'second' to be called")
	}
}

// TC-05: RunByName returns error for unknown name.
func TestRunByName_NotFound(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	err := RunByName(setupTestDB(t), "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown seeder")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected 'not found' in error, got %q", err.Error())
	}
}

// TC-06: RunAll with empty registry returns nil (and logs a warning).
func TestRunAll_EmptyRegistry(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	db := setupTestDB(t)
	if err := RunAll(db); err != nil {
		t.Fatalf("RunAll on empty registry should return nil, got: %v", err)
	}
}
