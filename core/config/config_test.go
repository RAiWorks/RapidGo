package config

import (
	"os"
	"testing"
)

// TC-01: Load() with .env file present (handled by integration test — go run)

// TC-02: Load() without .env file — should not panic
func TestLoad_NoEnvFile(t *testing.T) {
	// Change to a temp dir where no .env exists
	original, _ := os.Getwd()
	tmp := t.TempDir()
	os.Chdir(tmp)
	defer os.Chdir(original)

	// Should not panic — just logs a message
	Load()
}

// TC-03: Env() with key present
func TestEnv_KeyPresent(t *testing.T) {
	t.Setenv("TEST_KEY", "hello")
	got := Env("TEST_KEY", "default")
	if got != "hello" {
		t.Errorf("Env() = %q, want %q", got, "hello")
	}
}

// TC-04: Env() with key absent (fallback)
func TestEnv_KeyAbsent(t *testing.T) {
	got := Env("TEST_MISSING_KEY_XYZ", "fallback_value")
	if got != "fallback_value" {
		t.Errorf("Env() = %q, want %q", got, "fallback_value")
	}
}

// TC-13: Env() with empty string value
func TestEnv_EmptyValue(t *testing.T) {
	t.Setenv("TEST_EMPTY", "")
	got := Env("TEST_EMPTY", "fallback")
	if got != "fallback" {
		t.Errorf("Env() = %q, want %q", got, "fallback")
	}
}

// TC-05: EnvInt() with valid integer
func TestEnvInt_ValidInt(t *testing.T) {
	t.Setenv("TEST_INT", "42")
	got := EnvInt("TEST_INT", 0)
	if got != 42 {
		t.Errorf("EnvInt() = %d, want %d", got, 42)
	}
}

// TC-06: EnvInt() with invalid string
func TestEnvInt_InvalidString(t *testing.T) {
	t.Setenv("TEST_INT_BAD", "not_a_number")
	got := EnvInt("TEST_INT_BAD", 99)
	if got != 99 {
		t.Errorf("EnvInt() = %d, want %d", got, 99)
	}
}

// EnvInt() with empty value (fallback)
func TestEnvInt_Empty(t *testing.T) {
	got := EnvInt("TEST_INT_MISSING_XYZ", 77)
	if got != 77 {
		t.Errorf("EnvInt() = %d, want %d", got, 77)
	}
}

// TC-07: EnvBool() truthy values
func TestEnvBool_Truthy(t *testing.T) {
	t.Setenv("TEST_BOOL_T", "true")
	if !EnvBool("TEST_BOOL_T", false) {
		t.Error("EnvBool(\"true\") should return true")
	}

	t.Setenv("TEST_BOOL_1", "1")
	if !EnvBool("TEST_BOOL_1", false) {
		t.Error("EnvBool(\"1\") should return true")
	}
}

// TC-08: EnvBool() falsy values
func TestEnvBool_Falsy(t *testing.T) {
	t.Setenv("TEST_BOOL_F", "false")
	if EnvBool("TEST_BOOL_F", true) {
		t.Error("EnvBool(\"false\") should return false")
	}

	t.Setenv("TEST_BOOL_0", "0")
	if EnvBool("TEST_BOOL_0", true) {
		t.Error("EnvBool(\"0\") should return false")
	}
}

// TC-14: EnvBool() with empty (fallback)
func TestEnvBool_Empty(t *testing.T) {
	got := EnvBool("TEST_BOOL_MISSING_XYZ", true)
	if !got {
		t.Error("EnvBool() with missing key should return fallback (true)")
	}
}

// TC-09: Environment detection functions
func TestEnvironmentDetection_Production(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	if !IsProduction() {
		t.Error("IsProduction() should return true")
	}
	if IsDevelopment() {
		t.Error("IsDevelopment() should return false")
	}
	if IsTesting() {
		t.Error("IsTesting() should return false")
	}
}

func TestEnvironmentDetection_Development(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	if !IsDevelopment() {
		t.Error("IsDevelopment() should return true")
	}
	if IsProduction() {
		t.Error("IsProduction() should return false")
	}
}

func TestEnvironmentDetection_Testing(t *testing.T) {
	t.Setenv("APP_ENV", "testing")
	if !IsTesting() {
		t.Error("IsTesting() should return true")
	}
	if IsProduction() {
		t.Error("IsProduction() should return false")
	}
}

func TestEnvironmentDetection_Default(t *testing.T) {
	t.Setenv("APP_ENV", "")
	if !IsDevelopment() {
		t.Error("AppEnv() should default to development")
	}
}

func TestIsDebug(t *testing.T) {
	t.Setenv("APP_DEBUG", "true")
	if !IsDebug() {
		t.Error("IsDebug() should return true when APP_DEBUG=true")
	}

	t.Setenv("APP_DEBUG", "false")
	if IsDebug() {
		t.Error("IsDebug() should return false when APP_DEBUG=false")
	}
}
