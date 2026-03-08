package models

import (
	"reflect"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// testItem is a test-only model that embeds BaseModel.
type testItem struct {
	BaseModel
	Name string `gorm:"size:255"`
}

// TC-01: BaseModel has expected fields
func TestBaseModel_Fields(t *testing.T) {
	rt := reflect.TypeOf(BaseModel{})

	checks := []struct {
		name     string
		wantType string
	}{
		{"ID", "uint"},
		{"CreatedAt", "time.Time"},
		{"UpdatedAt", "time.Time"},
		{"DeletedAt", "gorm.DeletedAt"},
	}

	for _, c := range checks {
		f, ok := rt.FieldByName(c.name)
		if !ok {
			t.Fatalf("expected field %q on BaseModel", c.name)
		}
		if f.Type.String() != c.wantType {
			t.Fatalf("expected %q to be %s, got %s", c.name, c.wantType, f.Type.String())
		}
	}
}

// TC-02: Embedding BaseModel inherits ID
func TestBaseModel_EmbedInheritsID(t *testing.T) {
	item := testItem{BaseModel: BaseModel{ID: 42}}
	if item.ID != 42 {
		t.Fatalf("expected item.ID == 42, got %d", item.ID)
	}
}

// TC-03: AutoMigrate succeeds with BaseModel-embedded struct
func TestBaseModel_AutoMigrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&testItem{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
}

// TC-04: GORM creates and queries with BaseModel fields
func TestBaseModel_CreateAndQuery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	db.AutoMigrate(&testItem{})

	item := testItem{Name: "alpha"}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	var found testItem
	if err := db.First(&found, item.ID).Error; err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if found.Name != "alpha" {
		t.Fatalf("expected name 'alpha', got %q", found.Name)
	}
	if found.ID == 0 {
		t.Fatal("expected ID > 0")
	}
	if found.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
}
