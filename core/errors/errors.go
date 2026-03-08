package errors

import (
	"github.com/RAiWorks/RapidGo/v2/core/config"
)

// AppError represents a structured application error with HTTP status context.
type AppError struct {
	Code    int    // HTTP status code (e.g., 404, 500)
	Message string // User-facing safe message
	Err     error  // Internal error (logged, never exposed in production)
}

// Error returns the user-facing message. Implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error for errors.Is/As support.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NotFound creates a 404 error.
func NotFound(message string) *AppError {
	return &AppError{Code: 404, Message: message}
}

// BadRequest creates a 400 error.
func BadRequest(message string) *AppError {
	return &AppError{Code: 400, Message: message}
}

// Internal creates a 500 error wrapping an internal error.
func Internal(err error) *AppError {
	return &AppError{Code: 500, Message: "internal server error", Err: err}
}

// Unauthorized creates a 401 error.
func Unauthorized(message string) *AppError {
	return &AppError{Code: 401, Message: message}
}

// Forbidden creates a 403 error.
func Forbidden(message string) *AppError {
	return &AppError{Code: 403, Message: message}
}

// Conflict creates a 409 error.
func Conflict(message string) *AppError {
	return &AppError{Code: 409, Message: message}
}

// Unprocessable creates a 422 error.
func Unprocessable(message string) *AppError {
	return &AppError{Code: 422, Message: message}
}

// ErrorResponse returns a map for JSON error responses. In debug mode,
// it includes internal error details. In production, only the safe message.
func (e *AppError) ErrorResponse() map[string]any {
	resp := map[string]any{
		"success": false,
		"error":   e.Message,
	}
	if config.IsDebug() && e.Err != nil {
		resp["internal"] = e.Err.Error()
	}
	return resp
}
