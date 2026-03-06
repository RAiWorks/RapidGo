package services

import (
	"testing"

	"github.com/RAiWorks/RGo/database/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database with the User table.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

// TC-01: NewUserService returns a valid service
func TestNewUserService_ReturnsService(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)
	if svc == nil {
		t.Fatal("expected non-nil UserService")
	}
	if svc.DB == nil {
		t.Fatal("expected DB to be set")
	}
}

// TC-02: Create inserts a new user
func TestCreate_ReturnsUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	user, err := svc.Create("Alice", "alice@example.com", "pass123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected user ID > 0")
	}
	if user.Name != "Alice" {
		t.Fatalf("expected name 'Alice', got '%s'", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("expected email 'alice@example.com', got '%s'", user.Email)
	}
}

// TC-03: Create rejects duplicate email
func TestCreate_DuplicateEmail_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	_, err := svc.Create("Alice", "alice@example.com", "pass")
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	_, err = svc.Create("Bob", "alice@example.com", "pass")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	if err.Error() != "email already exists" {
		t.Fatalf("expected 'email already exists', got '%s'", err.Error())
	}
}

// TC-04: GetByID returns existing user
func TestGetByID_ReturnsUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	user, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Alice" {
		t.Fatalf("expected name 'Alice', got '%s'", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("expected email 'alice@example.com', got '%s'", user.Email)
	}
}

// TC-05: GetByID returns error for non-existent ID
func TestGetByID_NotFound_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	_, err := svc.GetByID(9999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// TC-06: Update modifies specified fields
func TestUpdate_UpdatesFields(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	updated, err := svc.Update(created.ID, map[string]interface{}{"name": "Bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Bob" {
		t.Fatalf("expected name 'Bob', got '%s'", updated.Name)
	}
}

// TC-07: Update returns error for non-existent ID
func TestUpdate_NotFound_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	_, err := svc.Update(9999, map[string]interface{}{"name": "X"})
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// TC-08: Delete removes user from database
func TestDelete_RemovesUser(t *testing.T) {
	db := setupTestDB(t)
	svc := NewUserService(db)

	created, _ := svc.Create("Alice", "alice@example.com", "pass")

	err := svc.Delete(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svc.GetByID(created.ID)
	if err == nil {
		t.Fatal("expected error after delete, user should not exist")
	}
}
