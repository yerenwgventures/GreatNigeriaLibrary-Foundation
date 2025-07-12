package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
)

// ErrorHandlerConfig holds configuration for error handling
type ErrorHandlerConfig struct {
	Logger        Logger
	IncludeStack  bool
	HideInternal  bool
	CustomHandler func(*gin.Context, error)
}

// NewErrorHandler creates a new error handling middleware with configuration
func NewErrorHandler(config ErrorHandlerConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			handleError(c, c.Errors.Last().Err, config)
		}
	})
}

// handleError processes different types of errors and returns appropriate responses
func handleError(c *gin.Context, err error, config ErrorHandlerConfig) {
	// Skip if response already written
	if c.Writer.Written() {
		return
	}

	// Log error with context
	logError(c, err, config.Logger)

	// Handle custom error handler
	if config.CustomHandler != nil {
		config.CustomHandler(c, err)
		return
	}

	// Handle different error types
	switch e := err.(type) {
	case *errors.AppError:
		handleAppError(c, e)
	case *errors.ValidationErrors:
		handleValidationErrors(c, e)
	case validator.ValidationErrors:
		handleValidatorErrors(c, e)
	default:
		handleGenericError(c, err, config.HideInternal)
	}
}

// handleAppError handles application-specific errors
func handleAppError(c *gin.Context, err *errors.AppError) {
	c.JSON(err.Code, gin.H{
		"success":   false,
		"message":   err.Message,
		"error":     err.Message,
		"code":      err.Type,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// handleValidationErrors handles custom validation errors
func handleValidationErrors(c *gin.Context, err *errors.ValidationErrors) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success":    false,
		"message":    "Validation failed",
		"error":      "Please check your input and try again",
		"code":       "VALIDATION_ERROR",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"validation": err.Errors,
	})
}

// handleValidatorErrors handles Go validator errors
func handleValidatorErrors(c *gin.Context, err validator.ValidationErrors) {
	validationErrors := make([]map[string]interface{}, 0, len(err))
	
	for _, fieldError := range err {
		validationErrors = append(validationErrors, map[string]interface{}{
			"field":   fieldError.Field(),
			"message": getValidationMessage(fieldError),
			"value":   fieldError.Value(),
			"tag":     fieldError.Tag(),
		})
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"success":    false,
		"message":    "Validation failed",
		"error":      "Please check your input and try again",
		"code":       "VALIDATION_ERROR",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"validation": validationErrors,
	})
}

// handleGenericError handles unknown errors
func handleGenericError(c *gin.Context, err error, hideInternal bool) {
	message := "Internal server error"
	errorDetail := "An unexpected error occurred. Please try again later."
	
	if !hideInternal {
		errorDetail = err.Error()
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"success":   false,
		"message":   message,
		"error":     errorDetail,
		"code":      "INTERNAL_SERVER_ERROR",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// logError logs the error with appropriate context
func logError(c *gin.Context, err error, logger Logger) {
	logFields := map[string]interface{}{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"client_ip":  c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"query":      c.Request.URL.RawQuery,
		"error":      err.Error(),
	}

	// Add user context if available
	if userID, exists := c.Get("user_id"); exists {
		logFields["user_id"] = userID
	}
	if username, exists := c.Get("username"); exists {
		logFields["username"] = username
	}

	// Log based on error type
	switch e := err.(type) {
	case *errors.AppError:
		if e.Code >= 500 {
			logger.WithFields(logFields).Error(fmt.Sprintf("Application error: %s", e.Message))
		} else {
			logger.WithFields(logFields).Info(fmt.Sprintf("Client error: %s", e.Message))
		}
	case *errors.ValidationErrors:
		logger.WithFields(logFields).Info("Validation error occurred")
	case validator.ValidationErrors:
		logger.WithFields(logFields).Info("Request validation failed")
	default:
		logger.WithFields(logFields).Error(fmt.Sprintf("Unhandled error: %v", err))
	}
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", fe.Field(), fe.Param())
	case "numeric":
		return fmt.Sprintf("%s must be a number", fe.Field())
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", fe.Field())
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", fe.Field())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

// AbortWithError is a helper function to abort request with error
func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

// AbortWithAppError is a helper function to abort request with AppError
func AbortWithAppError(c *gin.Context, appErr *errors.AppError) {
	c.Error(appErr)
	c.Abort()
}

// AbortWithValidationError is a helper function to abort request with validation error
func AbortWithValidationError(c *gin.Context, field, message, value string) {
	validationErr := &errors.ValidationErrors{}
	validationErr.Add(field, message, value)
	c.Error(validationErr)
	c.Abort()
}

// HandleBindingError handles JSON binding errors
func HandleBindingError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		c.Error(validationErrors)
	} else {
		c.Error(errors.ErrBadRequest("Invalid request format"))
	}
	c.Abort()
}

// HandleDatabaseError handles database-related errors
func HandleDatabaseError(c *gin.Context, err error, operation string) {
	// Log the database error
	if logger, exists := c.Get("logger"); exists {
		if l, ok := logger.(Logger); ok {
			l.WithFields(map[string]interface{}{
				"operation": operation,
				"error":     err.Error(),
			}).Error("Database operation failed")
		}
	}

	// Return appropriate error based on the error type
	c.Error(errors.ErrInternalServer(fmt.Sprintf("Database operation failed: %s", operation)))
	c.Abort()
}
