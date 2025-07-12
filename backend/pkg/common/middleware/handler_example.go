package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/response"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// ExampleHandler demonstrates how to use centralized error handling
type ExampleHandler struct {
	service ExampleService
	logger  Logger
}

type ExampleService interface {
	Login(req *models.UserLoginRequest) (*models.UserResponse, *models.TokenPair, error)
}

// BEFORE: Old error handling pattern (what we're replacing)
func (h *ExampleHandler) LoginOld(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// OLD: Direct JSON response with inconsistent format
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	user, tokens, err := h.service.Login(&req)
	if err != nil {
		// OLD: Manual error type checking and logging
		if e, ok := err.(*errors.AppError); ok {
			c.JSON(e.Code, e)
			return
		}
		h.logger.Error("Failed to login user")
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// OLD: Inconsistent success response format
	c.JSON(200, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// AFTER: New centralized error handling pattern
func (h *ExampleHandler) LoginNew(c *gin.Context) {
	var req models.UserLoginRequest
	
	// NEW: Use response helper for binding errors
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BindingError(c, err)
		return
	}

	// Get client context (example - adjust based on your actual model fields)
	// req.IP = c.ClientIP()
	// if req.DeviceInfo == "" {
	//     req.DeviceInfo = c.GetHeader("User-Agent")
	// }

	// Call service
	user, tokens, err := h.service.Login(&req)
	
	// NEW: Simple error handling - let middleware handle the details
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// NEW: Consistent success response
	response.OK(c, "Login successful", gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// ALTERNATIVE: Even simpler with direct error passing
func (h *ExampleHandler) LoginSimplest(c *gin.Context) {
	var req models.UserLoginRequest
	
	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrBadRequest("Invalid request format"))
		c.Abort()
		return
	}

	// Set context (example - adjust based on your actual model fields)
	// req.IP = c.ClientIP()
	// if req.DeviceInfo == "" {
	//     req.DeviceInfo = c.GetHeader("User-Agent")
	// }

	// Call service - any error will be handled by middleware
	user, tokens, err := h.service.Login(&req)
	if err != nil {
		c.Error(err) // Middleware will handle logging and response
		c.Abort()
		return
	}

	// Success response
	response.OK(c, "Login successful", gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// Example of validation error handling
func (h *ExampleHandler) ValidateExample(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	// Binding errors are automatically handled by middleware
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err) // Validator errors are automatically formatted
		c.Abort()
		return
	}

	// Custom validation
	if req.Email == "admin@example.com" {
		// Custom validation error
		validationErr := &errors.ValidationErrors{}
		validationErr.Add("email", "This email is reserved", req.Email)
		c.Error(validationErr)
		c.Abort()
		return
	}

	response.OK(c, "Validation passed", nil)
}

// Example of different error types
func (h *ExampleHandler) ErrorTypesExample(c *gin.Context) {
	action := c.Param("action")

	switch action {
	case "not-found":
		c.Error(errors.ErrNotFound("Resource"))
		c.Abort()
		
	case "unauthorized":
		c.Error(errors.ErrUnauthorized)
		c.Abort()
		
	case "forbidden":
		c.Error(errors.ErrForbidden)
		c.Abort()
		
	case "validation":
		validationErr := &errors.ValidationErrors{}
		validationErr.Add("field1", "Field is required", "")
		validationErr.Add("field2", "Field must be numeric", "abc")
		c.Error(validationErr)
		c.Abort()
		
	case "internal":
		c.Error(errors.ErrInternalServer("Something went wrong"))
		c.Abort()
		
	default:
		response.OK(c, "No error triggered", nil)
	}
}

// Example of conditional error handling
func (h *ExampleHandler) ConditionalExample(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Use conditional error helper
	if response.ConditionalError(c, userID == 0, errors.ErrUnauthorized) {
		return // Error was sent, handler exits
	}

	// Continue with normal flow
	response.OK(c, "User authenticated", gin.H{"user_id": userID})
}

// Example showing how middleware catches panics
func (h *ExampleHandler) PanicExample(c *gin.Context) {
	trigger := c.Query("trigger")
	
	if trigger == "panic" {
		panic("This is a test panic - middleware will catch it")
	}
	
	response.OK(c, "No panic occurred", nil)
}

/*
SUMMARY OF BENEFITS:

1. CONSISTENT ERROR RESPONSES:
   - All errors follow the same JSON structure
   - Timestamps and request IDs are automatically added
   - Error codes are standardized

2. CENTRALIZED LOGGING:
   - All errors are logged with consistent context
   - User information is automatically included
   - Stack traces for panics

3. SIMPLIFIED HANDLERS:
   - No need for manual error type checking
   - No need for manual logging
   - No need for manual JSON response formatting

4. BETTER VALIDATION:
   - Automatic validation error formatting
   - Human-readable validation messages
   - Support for custom validation errors

5. PANIC RECOVERY:
   - All panics are caught and logged
   - Stack traces are captured
   - Graceful error responses

6. MAINTAINABILITY:
   - Error handling logic is centralized
   - Easy to modify error response format
   - Consistent behavior across all endpoints

USAGE IN MAIN APPLICATION:
```go
// In your main.go or router setup
router := gin.New()

// Add the centralized error handling middleware
router.Use(middleware.PanicRecovery(logger))
router.Use(middleware.ErrorHandler(logger))

// Your handlers can now use simple error handling
router.POST("/login", userHandler.Login)
```
*/
