# Centralized Error Handling Guide

This document explains the centralized error handling system implemented in the Great Nigeria Library Foundation backend.

## Overview

The centralized error handling system provides:

- **Consistent Error Responses**: All errors follow the same JSON structure
- **Automatic Logging**: Errors are logged with full context automatically
- **Panic Recovery**: All panics are caught and handled gracefully
- **Validation Support**: Automatic formatting of validation errors
- **Simplified Handlers**: Handlers can focus on business logic

## Architecture

### Core Components

1. **Error Handling Middleware** (`middleware/error_handler.go`)
2. **Panic Recovery Middleware** (`middleware/middleware.go`)
3. **Response Utilities** (`response/response.go`)
4. **Error Types** (`errors/errors.go`)

### Middleware Stack

```go
router := gin.New()
router.Use(middleware.PanicRecovery(logger))  // Catches panics
router.Use(middleware.ErrorHandler(logger))   // Handles errors
router.Use(middleware.RequestLogger())        // Logs requests
router.Use(middleware.SecurityHeaders())      // Adds security headers
```

## Error Response Format

All errors return a consistent JSON structure:

```json
{
  "success": false,
  "message": "Human-readable error message",
  "error": "Detailed error description",
  "code": "ERROR_TYPE_CODE",
  "timestamp": "2023-12-07T10:30:00Z",
  "validation": [
    {
      "field": "email",
      "message": "Email is required",
      "value": ""
    }
  ]
}
```

## Error Types

### Application Errors (`AppError`)

```go
// Predefined errors
errors.ErrUnauthorized
errors.ErrForbidden
errors.ErrNotFound("Resource")
errors.ErrInternalServer("Operation failed")
errors.ErrBadRequest("Invalid input")

// Custom errors
&errors.AppError{
    Code:    http.StatusBadRequest,
    Message: "Custom error message",
    Type:    "CUSTOM_ERROR",
}
```

### Validation Errors

```go
// Single validation error
validationErr := &errors.ValidationErrors{}
validationErr.Add("email", "Email is required", "")
validationErr.Add("password", "Password too short", "123")

// Automatic validation (from gin binding)
if err := c.ShouldBindJSON(&req); err != nil {
    c.Error(err) // Automatically formatted
    c.Abort()
    return
}
```

## Handler Patterns

### Basic Error Handling

```go
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // Binding errors are automatically handled
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(err)
        c.Abort()
        return
    }
    
    // Service call
    user, err := h.service.CreateUser(&req)
    if err != nil {
        c.Error(err) // Middleware handles the response
        c.Abort()
        return
    }
    
    // Success response
    response.Created(c, "User created successfully", user)
}
```

### Using Response Helpers

```go
func (h *Handler) GetUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        response.AppError(c, errors.ErrBadRequest("Invalid user ID"))
        return
    }
    
    user, err := h.service.GetUser(uint(id))
    if err != nil {
        response.HandleError(c, err)
        return
    }
    
    response.OK(c, "User retrieved successfully", user)
}
```

### Conditional Error Handling

```go
func (h *Handler) UpdateUser(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    // Use conditional error helper
    if response.ConditionalError(c, userID == 0, errors.ErrUnauthorized) {
        return
    }
    
    // Continue with normal flow
    // ...
}
```

## Logging

### Automatic Context

The error handling middleware automatically logs errors with:

- HTTP method and path
- Client IP address
- User agent
- Query parameters
- User ID and username (if authenticated)
- Error details and stack traces (for panics)

### Log Levels

- **Info**: Client errors (4xx status codes)
- **Error**: Server errors (5xx status codes)
- **Error**: Panics with full stack traces

### Example Log Output

```json
{
  "level": "error",
  "timestamp": "2023-12-07T10:30:00Z",
  "message": "Application error: User not found",
  "method": "GET",
  "path": "/api/users/123",
  "client_ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "user_id": 456,
  "username": "john.doe",
  "error_type": "NOT_FOUND"
}
```

## Migration Guide

### Before (Old Pattern)

```go
func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    user, err := h.service.Login(&req)
    if err != nil {
        if appErr, ok := err.(*errors.AppError); ok {
            c.JSON(appErr.Code, appErr)
            return
        }
        h.logger.Error("Login failed: " + err.Error())
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(200, gin.H{"user": user})
}
```

### After (New Pattern)

```go
func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(err)
        c.Abort()
        return
    }
    
    user, err := h.service.Login(&req)
    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }
    
    response.OK(c, "Login successful", gin.H{"user": user})
}
```

## Best Practices

### 1. Use Appropriate Error Types

```go
// Good: Use specific error types
c.Error(errors.ErrNotFound("User"))
c.Error(errors.ErrUnauthorized)

// Avoid: Generic errors without context
c.Error(fmt.Errorf("error occurred"))
```

### 2. Let Middleware Handle Responses

```go
// Good: Let middleware handle the response
if err != nil {
    c.Error(err)
    c.Abort()
    return
}

// Avoid: Manual JSON responses
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()})
    return
}
```

### 3. Use Response Helpers

```go
// Good: Use response helpers for consistency
response.OK(c, "Success", data)
response.Created(c, "Resource created", resource)

// Avoid: Direct JSON responses
c.JSON(200, gin.H{"data": data})
```

### 4. Validate Early

```go
// Good: Validate input early
if err := c.ShouldBindJSON(&req); err != nil {
    c.Error(err)
    c.Abort()
    return
}

// Good: Custom validation
if req.Email == "" {
    validationErr := &errors.ValidationErrors{}
    validationErr.Add("email", "Email is required", "")
    c.Error(validationErr)
    c.Abort()
    return
}
```

## Testing Error Handling

### Unit Tests

```go
func TestHandler_CreateUser_ValidationError(t *testing.T) {
    // Setup
    handler := NewHandler(mockService, mockLogger)
    router := gin.New()
    router.Use(middleware.ErrorHandler(mockLogger))
    router.POST("/users", handler.CreateUser)
    
    // Test invalid request
    req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"email": ""}`))
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusBadRequest, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.False(t, response["success"].(bool))
    assert.Equal(t, "VALIDATION_ERROR", response["code"])
}
```

### Integration Tests

```go
func TestErrorHandling_Integration(t *testing.T) {
    // Test that panics are recovered
    // Test that different error types return correct responses
    // Test that logging works correctly
}
```

## Troubleshooting

### Common Issues

1. **Errors not being caught**: Ensure middleware is registered before routes
2. **Inconsistent responses**: Use response helpers instead of direct JSON
3. **Missing context**: Ensure user context is set in authentication middleware
4. **Validation not working**: Check that binding tags are correct

### Debug Mode

Enable debug logging to see detailed error information:

```go
config := ErrorHandlerConfig{
    Logger:       logger,
    IncludeStack: true,
    HideInternal: false,
}
router.Use(NewErrorHandler(config))
```

For more information, see the [API Documentation](api-reference.md) and [Configuration Guide](configuration.md).
