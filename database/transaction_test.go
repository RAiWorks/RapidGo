package database

import (
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// testUser is a minimal model for transaction tests.
// Avoids importing database/models to keep the database package independent.
type testUser struct {
	ID      uint `gorm:"primarykey"`
	Name    string
	Credits int
}

func (testUser) TableName() string { return "users" }

func setupTxTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&testUser{}); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db
}

// TC-01: WithTransaction commits when callback returns nil
func TestWithTransaction_Commit(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 1, Name: "Alice", Credits: 100})

	err := WithTransaction(db, func(tx *gorm.DB) error {
		return tx.Model(&testUser{}).Where("id = ?", 1).Update("credits", 200).Error
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var u testUser
	db.First(&u, 1)
	if u.Credits != 200 {
		t.Fatalf("expected credits 200, got %d", u.Credits)
	}
}

// TC-02: WithTransaction rolls back when callback returns error
func TestWithTransaction_Rollback(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 1, Name: "Alice", Credits: 100})

	err := WithTransaction(db, func(tx *gorm.DB) error {
		tx.Model(&testUser{}).Where("id = ?", 1).Update("credits", 999)
		return errors.New("forced rollback")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var u testUser
	db.First(&u, 1)
	if u.Credits != 100 {
		t.Fatalf("expected credits 100 after rollback, got %d", u.Credits)
	}
}

// TC-03: WithTransaction rolls back on panic inside callback
func TestWithTransaction_PanicRollback(t *testing.T) {
	db := setupTxTestDB(t)
	db.Create(&testUser{ID: 1, Name: "Alice", Credits: 100})

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic to propagate")
			}
		}()
		_ = WithTransaction(db, func(tx *gorm.DB) error {
			tx.Model(&testUser{}).Where("id = ?", 1).Update("credits", 999)
			panic("test panic")
		})
	}()

	// Give SQLite a moment to release any locks
	time.Sleep(10 * time.Millisecond)

	var u testUser
	db.First(&u, 1)
	if u.Credits != 100 {
		t.Fatalf("expected credits 100 after panic rollback, got %d", u.Credits)
	}
}

// TC-04: WithTransaction propagates the callback error to the caller
func TestWithTransaction_ErrorPropagation(t *testing.T) {
	db := setupTxTestDB(t)

	sentinel := errors.New("specific error")
	err := WithTransaction(db, func(tx *gorm.DB) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got: %v", err)
	}
}
