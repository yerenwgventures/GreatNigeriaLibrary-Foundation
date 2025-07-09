package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
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
