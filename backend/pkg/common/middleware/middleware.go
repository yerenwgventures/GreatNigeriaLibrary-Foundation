package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
)

// Logger interface for middleware
type Logger interface {
	Info(msg string)
	Error(msg string)
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// JWTManager interface for middleware
type JWTManager interface {
	ValidateToken(tokenString string) (*Claims, error)
	ExtractUserID(tokenString string) (uint, error)
	IsTokenRevoked(tokenString string) bool
}

// AuthorizationManager interface for middleware
type AuthorizationManager interface {
	HasPermission(role Role, permission Permission) bool
	HasAnyPermission(role Role, permissions []Permission) bool
	CanAccessResource(userRole Role, userID uint, resource string, action string, resourceOwnerID uint) bool
}

// Claims represents enhanced JWT claims
type Claims struct {
	UserID       uint     `json:"user_id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Role         int      `json:"role"`
	Permissions  []string `json:"permissions,omitempty"`
	SessionID    string   `json:"session_id,omitempty"`
	TokenType    string   `json:"token_type"`
	DeviceID     string   `json:"device_id,omitempty"`
	IPAddress    string   `json:"ip_address,omitempty"`
	TokenVersion int      `json:"token_version"`
}

// Role represents user roles
type Role int

// Permission represents a specific permission
type Permission string

// CORS middleware
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Allow specific origins or all origins in development
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"https://greatnigeria.com",
			"https://www.greatnigeria.com",
		}
		
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		
		if allowed || origin == "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// RequestLogger middleware logs HTTP requests
func RequestLogger(logger Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(map[string]interface{}{
			"method":     param.Method,
			"path":       param.Path,
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"user_agent": param.Request.UserAgent(),
		}).Info("HTTP Request")
		
		return ""
	})
}

// Recovery middleware recovers from panics with enhanced error handling
func Recovery(logger Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic with full context
		logger.WithFields(map[string]interface{}{
			"error":      recovered,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"query":      c.Request.URL.RawQuery,
			"headers":    c.Request.Header,
		}).Error("Panic recovered")

		// Return standardized error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":   false,
			"message":   "Internal server error",
			"error":     "An unexpected error occurred. Please try again later.",
			"code":      "INTERNAL_SERVER_ERROR",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
}

// ErrorHandler middleware provides centralized error handling
func ErrorHandler(logger Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Process the request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error (most recent)
			err := c.Errors.Last()

			// Log the error with context
			logFields := map[string]interface{}{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
				"query":      c.Request.URL.RawQuery,
			}

			// Add user context if available
			if userID, exists := c.Get("user_id"); exists {
				logFields["user_id"] = userID
			}
			if username, exists := c.Get("username"); exists {
				logFields["username"] = username
			}

			// Handle different error types
			switch e := err.Err.(type) {
			case *errors.AppError:
				// Application error - log and return structured response
				logger.WithFields(logFields).WithField("error_type", e.Type).Error(e.Message)

				if !c.Writer.Written() {
					c.JSON(e.Code, gin.H{
						"success":   false,
						"message":   e.Message,
						"error":     e.Message,
						"code":      e.Type,
						"timestamp": time.Now().UTC().Format(time.RFC3339),
					})
				}

			case *errors.ValidationErrors:
				// Validation errors - return detailed validation response
				logger.WithFields(logFields).Error("Validation failed")

				if !c.Writer.Written() {
					c.JSON(http.StatusBadRequest, gin.H{
						"success":    false,
						"message":    "Validation failed",
						"error":      "Please check your input and try again",
						"code":       "VALIDATION_ERROR",
						"timestamp":  time.Now().UTC().Format(time.RFC3339),
						"validation": e.Errors,
					})
				}

			default:
				// Unknown error - log as internal server error
				logger.WithFields(logFields).Error(fmt.Sprintf("Unhandled error: %v", e))

				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success":   false,
						"message":   "Internal server error",
						"error":     "An unexpected error occurred. Please try again later.",
						"code":      "INTERNAL_SERVER_ERROR",
						"timestamp": time.Now().UTC().Format(time.RFC3339),
					})
				}
			}
		}
	})
}

// PanicRecovery middleware provides enhanced panic recovery with stack traces
func PanicRecovery(logger Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Capture stack trace
		stack := debug.Stack()

		// Log the panic with full context and stack trace
		logger.WithFields(map[string]interface{}{
			"panic":      recovered,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"query":      c.Request.URL.RawQuery,
			"stack":      string(stack),
		}).Error("Panic recovered with stack trace")

		// Return standardized error response
		if !c.Writer.Written() {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":   false,
				"message":   "Internal server error",
				"error":     "An unexpected error occurred. Please try again later.",
				"code":      "PANIC_RECOVERED",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		}
	})
}

// RateLimit middleware (basic implementation)
func RateLimit() gin.HandlerFunc {
	// This is a basic implementation
	// In production, you'd want to use Redis or similar
	return gin.HandlerFunc(func(c *gin.Context) {
		// Basic rate limiting logic would go here
		c.Next()
	})
}

// AuthRequired middleware validates JWT tokens with enhanced security
func AuthRequired(jwtManager JWTManager, logger Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.WithField("path", c.Request.URL.Path).Error("Missing authorization header")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.WithField("header", authHeader).Error("Invalid authorization header format")
			c.Error(errors.ErrTokenInvalid)
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Enhanced token validation
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			logger.WithFields(map[string]interface{}{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
				"ip":    c.ClientIP(),
			}).Error("Token validation failed")
			c.Error(errors.ErrTokenInvalid)
			c.Abort()
			return
		}

		// Additional security checks
		if err := validateTokenSecurity(claims, c, logger); err != nil {
			logger.WithFields(map[string]interface{}{
				"error":   err.Error(),
				"user_id": claims.UserID,
				"path":    c.Request.URL.Path,
				"ip":      c.ClientIP(),
			}).Error("Token security validation failed")
			c.Error(errors.ErrTokenInvalid)
			c.Abort()
			return
		}

		// Set enhanced user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)
		c.Set("session_id", claims.SessionID)
		c.Set("device_id", claims.DeviceID)
		c.Set("token_version", claims.TokenVersion)
		c.Set("token", token) // Store token for potential revocation

		logger.WithFields(map[string]interface{}{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
			"path":     c.Request.URL.Path,
			"ip":       c.ClientIP(),
		}).Info("User authenticated successfully")

		c.Next()
	})
}

// validateTokenSecurity performs additional security validations
func validateTokenSecurity(claims *Claims, c *gin.Context, logger Logger) error {
	// Validate token type (should be access token for API requests)
	if claims.TokenType != "access" && claims.TokenType != "" {
		return fmt.Errorf("invalid token type: %s", claims.TokenType)
	}

	// Optional: Validate IP address if stored in token
	if claims.IPAddress != "" && claims.IPAddress != c.ClientIP() {
		logger.WithFields(map[string]interface{}{
			"token_ip":   claims.IPAddress,
			"request_ip": c.ClientIP(),
			"user_id":    claims.UserID,
		}).Info("IP address mismatch detected")
		// Note: In production, you might want to be more strict about this
	}

	return nil
}

// OptionalAuth middleware validates JWT tokens but doesn't require them
func OptionalAuth(jwtManager JWTManager, logger Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("session_id", claims.SessionID)

		c.Next()
	})
}

// RoleRequired middleware checks if user has required role (enhanced)
func RoleRequired(requiredRole int, logger Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			logger.WithField("path", c.Request.URL.Path).Error("Role check failed: no role in context")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		role, ok := userRole.(int)
		if !ok || role < requiredRole {
			userID, _ := c.Get("user_id")
			logger.WithFields(map[string]interface{}{
				"user_id":       userID,
				"user_role":     role,
				"required_role": requiredRole,
				"path":          c.Request.URL.Path,
				"ip":            c.ClientIP(),
			}).Error("Insufficient role permissions")

			c.Error(errors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	})
}

// PermissionRequired middleware checks if user has specific permission
func PermissionRequired(authManager AuthorizationManager, requiredPermission Permission, logger Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			logger.WithField("path", c.Request.URL.Path).Error("Permission check failed: no role in context")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		roleInt, ok := userRole.(int)
		if !ok {
			logger.WithField("path", c.Request.URL.Path).Error("Permission check failed: invalid role type")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		role := Role(roleInt)
		if !authManager.HasPermission(role, requiredPermission) {
			userID, _ := c.Get("user_id")
			logger.WithFields(map[string]interface{}{
				"user_id":            userID,
				"user_role":          role,
				"required_permission": requiredPermission,
				"path":               c.Request.URL.Path,
				"ip":                 c.ClientIP(),
			}).Error("Insufficient permissions")

			c.Error(errors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	})
}

// AnyPermissionRequired middleware checks if user has any of the specified permissions
func AnyPermissionRequired(authManager AuthorizationManager, requiredPermissions []Permission, logger Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			logger.WithField("path", c.Request.URL.Path).Error("Permission check failed: no role in context")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		roleInt, ok := userRole.(int)
		if !ok {
			logger.WithField("path", c.Request.URL.Path).Error("Permission check failed: invalid role type")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		role := Role(roleInt)
		if !authManager.HasAnyPermission(role, requiredPermissions) {
			userID, _ := c.Get("user_id")
			logger.WithFields(map[string]interface{}{
				"user_id":             userID,
				"user_role":           role,
				"required_permissions": requiredPermissions,
				"path":                c.Request.URL.Path,
				"ip":                  c.ClientIP(),
			}).Error("Insufficient permissions")

			c.Error(errors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	})
}

// ResourceOwnerOrPermission middleware checks if user owns resource or has permission
func ResourceOwnerOrPermission(authManager AuthorizationManager, permission Permission, resourceOwnerIDKey string, logger Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			logger.WithField("path", c.Request.URL.Path).Error("Resource access check failed: no user ID in context")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		userIDUint, ok := userID.(uint)
		if !ok {
			logger.WithField("path", c.Request.URL.Path).Error("Resource access check failed: invalid user ID type")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Get resource owner ID from context or URL parameter
		var resourceOwnerID uint
		if ownerID, exists := c.Get(resourceOwnerIDKey); exists {
			if ownerIDUint, ok := ownerID.(uint); ok {
				resourceOwnerID = ownerIDUint
			}
		}

		// Check if user owns the resource
		if userIDUint == resourceOwnerID {
			c.Next()
			return
		}

		// Check if user has the required permission
		userRole, exists := c.Get("role")
		if !exists {
			logger.WithField("path", c.Request.URL.Path).Error("Resource access check failed: no role in context")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		roleInt, ok := userRole.(int)
		if !ok {
			logger.WithField("path", c.Request.URL.Path).Error("Resource access check failed: invalid role type")
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		role := Role(roleInt)
		if !authManager.HasPermission(role, permission) {
			logger.WithFields(map[string]interface{}{
				"user_id":          userIDUint,
				"resource_owner":   resourceOwnerID,
				"required_permission": permission,
				"path":             c.Request.URL.Path,
				"ip":               c.ClientIP(),
			}).Error("Insufficient permissions for resource access")

			c.Error(errors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	})
}

// AdminRequired middleware checks if user is admin
func AdminRequired(logger Logger) gin.HandlerFunc {
	return RoleRequired(3, logger) // Assuming 3 is admin role
}

// ModeratorRequired middleware checks if user is moderator or admin
func ModeratorRequired(logger Logger) gin.HandlerFunc {
	return RoleRequired(2, logger) // Assuming 2 is moderator role
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	})
}

// Timeout middleware adds request timeout
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// This would implement request timeout logic
		// For now, just pass through
		c.Next()
	})
}

// GetUserID helper function to extract user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	
	id, ok := userID.(uint)
	return id, ok
}

// GetUserRole helper function to extract user role from context
func GetUserRole(c *gin.Context) (int, bool) {
	userRole, exists := c.Get("role")
	if !exists {
		return 0, false
	}
	
	role, ok := userRole.(int)
	return role, ok
}

// GetUsername helper function to extract username from context
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	
	name, ok := username.(string)
	return name, ok
}
