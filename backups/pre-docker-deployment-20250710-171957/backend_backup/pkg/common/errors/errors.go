package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// HTTPStatus returns the HTTP status code for the error
func (e *AppError) HTTPStatus() int {
	return e.Code
}

// ErrorType represents error types
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "validation"
	ErrorTypeAuth         ErrorType = "authentication"
	ErrorTypeAuthorization ErrorType = "authorization"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeInternal     ErrorType = "internal"
	ErrorTypeExternal     ErrorType = "external"
	ErrorTypeRateLimit    ErrorType = "rate_limit"
)

// Predefined errors
var (
	// Authentication errors
	ErrInvalidCredentials = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid email or password",
		Type:    string(ErrorTypeAuth),
	}

	ErrTokenExpired = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Token has expired",
		Type:    string(ErrorTypeAuth),
	}

	ErrTokenInvalid = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid token",
		Type:    string(ErrorTypeAuth),
	}

	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized access",
		Type:    string(ErrorTypeAuth),
	}

	// Authorization errors
	ErrForbidden = &AppError{
		Code:    http.StatusForbidden,
		Message: "Access forbidden",
		Type:    string(ErrorTypeAuthorization),
	}

	ErrInsufficientPermissions = &AppError{
		Code:    http.StatusForbidden,
		Message: "Insufficient permissions",
		Type:    string(ErrorTypeAuthorization),
	}

	// Validation errors
	ErrValidationFailed = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Validation failed",
		Type:    string(ErrorTypeValidation),
	}

	ErrInvalidInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid input provided",
		Type:    string(ErrorTypeValidation),
	}

	ErrMissingRequiredField = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Missing required field",
		Type:    string(ErrorTypeValidation),
	}

	// Not found errors
	ErrUserNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "User not found",
		Type:    string(ErrorTypeNotFound),
	}

	ErrResourceNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Resource not found",
		Type:    string(ErrorTypeNotFound),
	}

	// Conflict errors
	ErrUserAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "User already exists",
		Type:    string(ErrorTypeConflict),
	}

	ErrEmailAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Email already exists",
		Type:    string(ErrorTypeConflict),
	}

	ErrUsernameAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Username already exists",
		Type:    string(ErrorTypeConflict),
	}

	ErrResourceAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "Resource already exists",
		Type:    string(ErrorTypeConflict),
	}

	// Internal errors
	ErrInternalServerError = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		Type:    string(ErrorTypeInternal),
	}

	ErrDatabaseError = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Database error",
		Type:    string(ErrorTypeInternal),
	}

	ErrServiceUnavailable = &AppError{
		Code:    http.StatusServiceUnavailable,
		Message: "Service temporarily unavailable",
		Type:    string(ErrorTypeInternal),
	}

	// Rate limiting errors
	ErrRateLimitExceeded = &AppError{
		Code:    http.StatusTooManyRequests,
		Message: "Rate limit exceeded",
		Type:    string(ErrorTypeRateLimit),
	}
)

// Error creation functions
func ErrValidation(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Type:    string(ErrorTypeValidation),
	}
}

func ErrNotFound(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Type:    string(ErrorTypeNotFound),
	}
}

func ErrConflict(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
		Type:    string(ErrorTypeConflict),
	}
}

func ErrInternalServer(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Type:    string(ErrorTypeInternal),
	}
}

func ErrBadRequest(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Type:    string(ErrorTypeValidation),
	}
}

func ErrUnauthorizedAccess(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Type:    string(ErrorTypeAuth),
	}
}

func ErrForbiddenAccess(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
		Type:    string(ErrorTypeAuthorization),
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError converts an error to AppError if possible
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return ErrInternalServer(err.Error())
}

// WrapError wraps an error with additional context
func WrapError(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Code:    appErr.Code,
			Message: fmt.Sprintf("%s: %s", message, appErr.Message),
			Type:    appErr.Type,
		}
	}
	return ErrInternalServer(fmt.Sprintf("%s: %s", message, err.Error()))
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Error implements the error interface for ValidationErrors
func (v *ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %s", v.Errors[0].Message)
}

// Add adds a validation error
func (v *ValidationErrors) Add(field, message, value string) {
	v.Errors = append(v.Errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors returns true if there are validation errors
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// NewValidationErrors creates a new ValidationErrors instance
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}
