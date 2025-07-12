package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Code      string      `json:"code,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ErrorResponse represents a detailed error response
type ErrorResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	Error       string      `json:"error"`
	Code        string      `json:"code"`
	Timestamp   string      `json:"timestamp"`
	Validation  interface{} `json:"validation,omitempty"`
	RequestID   string      `json:"request_id,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, err error) {
	response := APIResponse{
		Success: false,
		Message: message,
	}
	
	if err != nil {
		response.Error = err.Error()
	}
	
	c.JSON(statusCode, response)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, err error) {
	Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string, err error) {
	Error(c, http.StatusUnauthorized, message, err)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string, err error) {
	Error(c, http.StatusForbidden, message, err)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string, err error) {
	Error(c, http.StatusNotFound, message, err)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string, err error) {
	Error(c, http.StatusInternalServerError, message, err)
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// OK sends a 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ErrorWithCode sends an error response with custom error code
func ErrorWithCode(c *gin.Context, statusCode int, message, errorCode string, err error) {
	response := ErrorResponse{
		Success:   false,
		Message:   message,
		Code:      errorCode,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err != nil {
		response.Error = err.Error()
	} else {
		response.Error = message
	}

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	c.JSON(statusCode, response)
}

// AppError sends an AppError response using centralized error handling
func AppError(c *gin.Context, appErr *errors.AppError) {
	c.Error(appErr)
	c.Abort()
}

// ValidationError sends a validation error response using centralized error handling
func ValidationError(c *gin.Context, validationErr *errors.ValidationErrors) {
	c.Error(validationErr)
	c.Abort()
}

// BindingError handles JSON binding errors and sends appropriate response
func BindingError(c *gin.Context, err error) {
	c.Error(errors.ErrBadRequest("Invalid request format: " + err.Error()))
	c.Abort()
}

// DatabaseError handles database errors and sends appropriate response
func DatabaseError(c *gin.Context, err error, operation string) {
	c.Error(errors.ErrInternalServer("Database operation failed: " + operation))
	c.Abort()
}

// ServiceError handles service layer errors and sends appropriate response
func ServiceError(c *gin.Context, err error, service string) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.Error(appErr)
	} else {
		c.Error(errors.ErrInternalServer("Service error in " + service + ": " + err.Error()))
	}
	c.Abort()
}

// PaginatedSuccess sends a successful paginated response
func PaginatedSuccess(c *gin.Context, message string, data interface{}, pagination interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    message,
		"data":       data,
		"pagination": pagination,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

// ConditionalError sends error only if condition is true, otherwise continues
func ConditionalError(c *gin.Context, condition bool, appErr *errors.AppError) bool {
	if condition {
		c.Error(appErr)
		c.Abort()
		return true
	}
	return false
}

// HandleError is a generic error handler that determines the best response
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *errors.AppError:
		c.Error(e)
	case *errors.ValidationErrors:
		c.Error(e)
	default:
		c.Error(errors.ErrInternalServer(err.Error()))
	}
	c.Abort()
}
