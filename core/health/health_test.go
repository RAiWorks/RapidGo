package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RAiWorks/RapidGo/v2/core/router"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupRouter(t *testing.T, db *gorm.DB) *router.Router {
	t.Helper()
	t.Setenv("APP_ENV", "testing")
	r := router.New()
	Routes(r, func() *gorm.DB { return db })
	return r
}

func openDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db
}

// TC-01: GET /health returns 200 with {"status":"ok"}
func TestLiveness_ReturnsOK(t *testing.T) {
	r := setupRouter(t, openDB(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status = %q, want %q", body["status"], "ok")
	}
}

// TC-02: GET /health/ready returns 200 with live DB
func TestReadiness_WithLiveDB(t *testing.T) {
	r := setupRouter(t, openDB(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health/ready status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ready" {
		t.Fatalf("status = %q, want %q", body["status"], "ready")
	}
	if body["db"] != "connected" {
		t.Fatalf("db = %q, want %q", body["db"], "connected")
	}
}

// TC-03: GET /health/ready returns 503 when DB is closed
func TestReadiness_WithClosedDB(t *testing.T) {
	db := openDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB(): %v", err)
	}
	sqlDB.Close()

	r := setupRouter(t, db)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready status = %d, want %d", w.Code, http.StatusServiceUnavailable)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "error" {
		t.Fatalf("status = %q, want %q", body["status"], "error")
	}
	if body["db"] == "" {
		t.Fatal("expected non-empty db error message")
	}
}
