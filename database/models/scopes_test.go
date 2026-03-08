package models

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// testScopesItem embeds BaseModel (which provides DeletedAt for soft-delete).
type testScopesItem struct {
	BaseModel
	Label string `gorm:"size:255"`
}

func setupScopesTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&testScopesItem{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

// T01: WithTrashed includes soft-deleted records
func TestWithTrashed_IncludesDeletedRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&testScopesItem{Label: "Alice"})
	db.Create(&testScopesItem{Label: "Bob"})

	// Soft-delete Alice
	db.Where("label = ?", "Alice").Delete(&testScopesItem{})

	var items []testScopesItem
	db.Scopes(WithTrashed).Find(&items)
	if len(items) != 2 {
		t.Fatalf("expected 2 items (including deleted), got %d", len(items))
	}
}

// T02: WithTrashed includes active records
func TestWithTrashed_IncludesActiveRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&testScopesItem{Label: "Alice"})

	var items []testScopesItem
	db.Scopes(WithTrashed).Find(&items)
	if len(items) != 1 {
		t.Fatalf("expected 1 active item, got %d", len(items))
	}
	if items[0].Label != "Alice" {
		t.Fatalf("expected 'Alice', got '%s'", items[0].Label)
	}
}

// T03: OnlyTrashed returns only soft-deleted records
func TestOnlyTrashed_ReturnsOnlyDeletedRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&testScopesItem{Label: "Alice"})
	db.Create(&testScopesItem{Label: "Bob"})

	// Soft-delete Bob
	db.Where("label = ?", "Bob").Delete(&testScopesItem{})

	var items []testScopesItem
	db.Scopes(OnlyTrashed).Find(&items)
	if len(items) != 1 {
		t.Fatalf("expected 1 trashed item, got %d", len(items))
	}
	if items[0].Label != "Bob" {
		t.Fatalf("expected 'Bob', got '%s'", items[0].Label)
	}
}

// T04: OnlyTrashed excludes active records
func TestOnlyTrashed_ExcludesActiveRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&testScopesItem{Label: "Alice"})
	db.Create(&testScopesItem{Label: "Bob"})

	// No deletions
	var items []testScopesItem
	db.Scopes(OnlyTrashed).Find(&items)
	if len(items) != 0 {
		t.Fatalf("expected 0 trashed items, got %d", len(items))
	}
}

// T05: Default query excludes soft-deleted records
func TestDefaultQuery_ExcludesDeletedRecords(t *testing.T) {
	db := setupScopesTestDB(t)
	db.Create(&testScopesItem{Label: "Alice"})
	db.Create(&testScopesItem{Label: "Bob"})

	// Soft-delete Alice
	db.Where("label = ?", "Alice").Delete(&testScopesItem{})

	var items []testScopesItem
	db.Find(&items)
	if len(items) != 1 {
		t.Fatalf("expected 1 active item, got %d", len(items))
	}
	if items[0].Label != "Bob" {
		t.Fatalf("expected 'Bob', got '%s'", items[0].Label)
	}
}
