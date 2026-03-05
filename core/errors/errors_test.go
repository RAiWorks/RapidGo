package errors

import (
	"errors"
	"io"
	"os"
	"testing"
)

// --- TC-01: AppError implements error interface ---
func TestAppError_ImplementsErrorInterface(t *testing.T) {
	appErr := &AppError{Code: 404, Message: "not found"}
	var _ error = appErr // compile-time check

	if got := appErr.Error(); got != "not found" {
		t.Errorf("Error() = %q, want %q", got, "not found")
	}
}

// --- TC-02: Unwrap returns wrapped error, errors.Is works ---
func TestUnwrap_ReturnsWrappedError(t *testing.T) {
	appErr := Internal(io.ErrUnexpectedEOF)

	if appErr.Unwrap() != io.ErrUnexpectedEOF {
		t.Errorf("Unwrap() = %v, want %v", appErr.Unwrap(), io.ErrUnexpectedEOF)
	}
	if !errors.Is(appErr, io.ErrUnexpectedEOF) {
		t.Error("errors.Is(appErr, io.ErrUnexpectedEOF) = false, want true")
	}
}

// --- TC-03: Unwrap returns nil when no wrapped error ---
func TestUnwrap_NilErr(t *testing.T) {
	appErr := NotFound("user not found")

	if appErr.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", appErr.Unwrap())
	}
}

// --- TC-04: errors.As extracts AppError ---
func TestErrorsAs_ExtractsAppError(t *testing.T) {
	var err error = NotFound("page not found")

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatal("errors.As failed to extract *AppError")
	}
	if appErr.Code != 404 {
		t.Errorf("Code = %d, want 404", appErr.Code)
	}
}

// --- TC-05 through TC-11: Constructor tests ---
func TestConstructors(t *testing.T) {
	tests := []struct {
		name    string
		create  func() *AppError
		code    int
		message string
		hasErr  bool
	}{
		// TC-05
		{"NotFound", func() *AppError { return NotFound("user not found") }, 404, "user not found", false},
		// TC-06
		{"BadRequest", func() *AppError { return BadRequest("invalid input") }, 400, "invalid input", false},
		// TC-07
		{"Internal", func() *AppError { return Internal(io.EOF) }, 500, "internal server error", true},
		// TC-08
		{"Unauthorized", func() *AppError { return Unauthorized("login required") }, 401, "login required", false},
		// TC-09
		{"Forbidden", func() *AppError { return Forbidden("access denied") }, 403, "access denied", false},
		// TC-10
		{"Conflict", func() *AppError { return Conflict("already exists") }, 409, "already exists", false},
		// TC-11
		{"Unprocessable", func() *AppError { return Unprocessable("validation failed") }, 422, "validation failed", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := tt.create()
			if appErr.Code != tt.code {
				t.Errorf("Code = %d, want %d", appErr.Code, tt.code)
			}
			if appErr.Message != tt.message {
				t.Errorf("Message = %q, want %q", appErr.Message, tt.message)
			}
			if tt.hasErr && appErr.Err == nil {
				t.Error("Err = nil, want non-nil")
			}
			if !tt.hasErr && appErr.Err != nil {
				t.Errorf("Err = %v, want nil", appErr.Err)
			}
		})
	}
}

// --- TC-12: ErrorResponse in debug mode includes internal error ---
func TestErrorResponse_DebugMode(t *testing.T) {
	os.Setenv("APP_DEBUG", "true")
	defer os.Unsetenv("APP_DEBUG")

	appErr := Internal(errors.New("db timeout"))
	resp := appErr.ErrorResponse()

	if resp["success"] != false {
		t.Errorf("success = %v, want false", resp["success"])
	}
	if resp["error"] != "internal server error" {
		t.Errorf("error = %v, want %q", resp["error"], "internal server error")
	}
	if resp["internal"] != "db timeout" {
		t.Errorf("internal = %v, want %q", resp["internal"], "db timeout")
	}
}

// --- TC-13: ErrorResponse in production — no internal details ---
func TestErrorResponse_ProductionMode(t *testing.T) {
	os.Setenv("APP_DEBUG", "false")
	defer os.Unsetenv("APP_DEBUG")

	appErr := Internal(errors.New("db timeout"))
	resp := appErr.ErrorResponse()

	if resp["success"] != false {
		t.Errorf("success = %v, want false", resp["success"])
	}
	if resp["error"] != "internal server error" {
		t.Errorf("error = %v, want %q", resp["error"], "internal server error")
	}
	if _, exists := resp["internal"]; exists {
		t.Errorf("internal key should not exist in production, got %v", resp["internal"])
	}
}

// --- TC-14: ErrorResponse with nil Err in debug mode — no internal key ---
func TestErrorResponse_DebugMode_NilErr(t *testing.T) {
	os.Setenv("APP_DEBUG", "true")
	defer os.Unsetenv("APP_DEBUG")

	appErr := NotFound("not found")
	resp := appErr.ErrorResponse()

	if resp["success"] != false {
		t.Errorf("success = %v, want false", resp["success"])
	}
	if resp["error"] != "not found" {
		t.Errorf("error = %v, want %q", resp["error"], "not found")
	}
	if _, exists := resp["internal"]; exists {
		t.Error("internal key should not exist when Err is nil")
	}
}
