# JWT Authentication & Authorization Guide

This document explains the enhanced JWT authentication and authorization system implemented in the Great Nigeria Library Foundation backend.

## Overview

The JWT authentication system has been completely rewritten with enhanced security features, role-based access control, and Redis-based token management.

### Key Features

- **Enhanced JWT Tokens**: Extended claims with permissions, device tracking, and security metadata
- **Redis-Based Token Revocation**: Real-time token invalidation and session management
- **Role-Based Access Control (RBAC)**: Granular permissions system with predefined roles
- **Session Management**: Device tracking and session lifecycle management
- **Security Enhancements**: IP validation, token versioning, and comprehensive logging

## Architecture

### Core Components

1. **JWTManager** (`pkg/common/auth/jwt.go`): Enhanced JWT token management with Redis support
2. **AuthorizationManager** (`pkg/common/auth/authorization.go`): Role-based access control system
3. **Enhanced Middleware** (`pkg/common/middleware/middleware.go`): Authentication and authorization middleware
4. **Redis Integration** (`pkg/common/redis/redis.go`): Token storage and session management

### Authentication Flow

```
User Login Request
       ↓
Credential Validation
       ↓
Session Creation
       ↓
Enhanced JWT Generation
       ↓
Redis Token Storage
       ↓
Response with Tokens
```

## JWT Token Structure

### Enhanced Claims

```go
type Claims struct {
    UserID       uint     `json:"user_id"`
    Username     string   `json:"username"`
    Email        string   `json:"email"`
    Role         int      `json:"role"`
    Permissions  []string `json:"permissions,omitempty"`
    SessionID    string   `json:"session_id,omitempty"`
    TokenType    string   `json:"token_type"` // "access" or "refresh"
    DeviceID     string   `json:"device_id,omitempty"`
    IPAddress    string   `json:"ip_address,omitempty"`
    TokenVersion int      `json:"token_version"`
    jwt.RegisteredClaims
}
```

### Token Types

- **Access Token**: Short-lived (15 minutes) for API access
- **Refresh Token**: Long-lived (7 days) for token renewal

### Security Features

- **Token Versioning**: Invalidate all user tokens by incrementing version
- **Device Tracking**: Associate tokens with specific devices
- **IP Validation**: Optional IP address validation
- **Revocation Support**: Real-time token blacklisting via Redis

## Role-Based Access Control

### Predefined Roles

```go
const (
    RoleGuest      Role = 0  // Limited read access
    RoleUser       Role = 1  // Basic user permissions
    RoleModerator  Role = 2  // Content moderation
    RoleAdmin      Role = 3  // System administration
    RoleSuperAdmin Role = 4  // Full system access
)
```

### Permission System

#### User Permissions
- `user:read_profile` - View user profiles
- `user:update_profile` - Update own profile
- `user:delete_profile` - Delete own profile

#### Content Permissions
- `content:read` - View content
- `content:create` - Create new content
- `content:update` - Update content
- `content:delete` - Delete content
- `content:publish` - Publish content

#### Discussion Permissions
- `discussion:read` - View discussions
- `discussion:create` - Create discussions
- `discussion:update` - Update discussions
- `discussion:delete` - Delete discussions
- `discussion:moderate` - Moderate discussions

#### Group Permissions
- `group:read` - View groups
- `group:create` - Create groups
- `group:update` - Update groups
- `group:delete` - Delete groups
- `group:manage` - Manage group members

#### Admin Permissions
- `admin:manage_users` - User management
- `admin:manage_content` - Content management
- `admin:manage_system` - System configuration
- `admin:view_analytics` - View analytics
- `admin:manage_settings` - Manage settings

### Role-Permission Mapping

```go
// User Role Permissions
RoleUser: []Permission{
    PermissionReadContent,
    PermissionCreateContent,
    PermissionUpdateContent, // Own content only
    PermissionReadProfile,
    PermissionUpdateProfile,
    // ... more permissions
}

// Admin Role Permissions
RoleAdmin: []Permission{
    // All user permissions plus:
    PermissionDeleteContent,
    PermissionPublishContent,
    PermissionManageUsers,
    PermissionManageContent,
    // ... more permissions
}
```

## Authentication Middleware

### AuthRequired Middleware

Validates JWT tokens and sets user context:

```go
router.Use(middleware.AuthRequired(jwtManager, logger))
```

**Features:**
- Token validation with enhanced security checks
- User context injection
- IP address validation (optional)
- Device tracking
- Comprehensive logging

### Authorization Middleware

#### Role-Based Authorization

```go
// Require specific role
router.Use(middleware.RoleRequired(auth.RoleAdmin, logger))

// Require specific permission
router.Use(middleware.PermissionRequired(authManager, auth.PermissionManageUsers, logger))

// Require any of multiple permissions
permissions := []auth.Permission{
    auth.PermissionUpdateContent,
    auth.PermissionDeleteContent,
}
router.Use(middleware.AnyPermissionRequired(authManager, permissions, logger))
```

#### Resource-Based Authorization

```go
// Allow resource owner or users with specific permission
router.Use(middleware.ResourceOwnerOrPermission(
    authManager,
    auth.PermissionDeleteContent,
    "content_owner_id",
    logger,
))
```

## Redis Integration

### Token Storage

```go
// Store token metadata
tokenStore.StoreToken(ctx, tokenHash, metadata, expiration)

// Check if token is revoked
isRevoked := tokenStore.IsTokenRevoked(ctx, tokenHash)

// Revoke specific token
tokenStore.RevokeToken(ctx, tokenHash, expiration)
```

### Session Management

```go
// Store session
tokenStore.StoreSession(ctx, sessionID, sessionData, expiration)

// Get user sessions
sessions := tokenStore.GetUserSessions(ctx, userID)

// Revoke all user sessions
tokenStore.RevokeAllUserSessions(ctx, userID)
```

### Token Versioning

```go
// Increment user token version (invalidates all tokens)
newVersion := tokenStore.IncrementUserTokenVersion(ctx, userID)

// Get current token version
version := tokenStore.GetUserTokenVersion(ctx, userID)
```

## Configuration

### Environment Variables

```bash
# JWT Configuration
JWT_SECRET=your-super-secret-key
JWT_ISSUER=great-nigeria-library
ACCESS_TOKEN_EXPIRATION=15m
REFRESH_TOKEN_EXPIRATION=168h
ENABLE_TOKEN_REVOCATION=true
ENABLE_SESSION_TRACKING=true
TOKEN_SECURITY_CHECKS=true

# Redis Configuration
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5
REDIS_MAX_RETRIES=3
```

### YAML Configuration

```yaml
auth:
  jwt_secret: "your-super-secret-key"
  jwt_issuer: "great-nigeria-library"
  access_token_expiration: "15m"
  refresh_token_expiration: "168h"
  enable_token_revocation: true
  enable_session_tracking: true
  token_security_checks: true

redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  database: 0
  pool_size: 10
  min_idle_conns: 5
  max_retries: 3
```

## Usage Examples

### Basic Authentication

```go
// Protect route with authentication
protected := router.Group("/api/v1")
protected.Use(middleware.AuthRequired(jwtManager, logger))
{
    protected.GET("/profile", userHandler.GetProfile)
    protected.PUT("/profile", userHandler.UpdateProfile)
}
```

### Role-Based Protection

```go
// Admin-only routes
admin := router.Group("/api/v1/admin")
admin.Use(middleware.AuthRequired(jwtManager, logger))
admin.Use(middleware.RoleRequired(auth.RoleAdmin, logger))
{
    admin.GET("/users", adminHandler.ListUsers)
    admin.DELETE("/users/:id", adminHandler.DeleteUser)
}
```

### Permission-Based Protection

```go
// Content management routes
content := router.Group("/api/v1/content")
content.Use(middleware.AuthRequired(jwtManager, logger))
{
    // Anyone can read
    content.GET("/:id", contentHandler.GetContent)
    
    // Require create permission
    content.POST("/", 
        middleware.PermissionRequired(authManager, auth.PermissionCreateContent, logger),
        contentHandler.CreateContent,
    )
    
    // Resource owner or delete permission
    content.DELETE("/:id",
        middleware.ResourceOwnerOrPermission(
            authManager,
            auth.PermissionDeleteContent,
            "content_owner_id",
            logger,
        ),
        contentHandler.DeleteContent,
    )
}
```

### Token Management

```go
// Generate enhanced tokens
tokens, err := jwtManager.GenerateTokensWithMetadata(
    userID,
    username,
    email,
    role,
    permissions,
    sessionID,
    deviceID,
    ipAddress,
)

// Validate token with security checks
claims, err := jwtManager.ValidateToken(tokenString)

// Revoke specific token
err := jwtManager.RevokeToken(tokenString)

// Revoke all user tokens
err := jwtManager.RevokeAllUserTokens(userID)
```

## Security Best Practices

### Token Security

1. **Use Strong Secrets**: Generate cryptographically secure JWT secrets
2. **Short Expiration**: Keep access tokens short-lived (15 minutes)
3. **Secure Storage**: Store refresh tokens securely on client side
4. **HTTPS Only**: Always use HTTPS in production
5. **Token Rotation**: Implement refresh token rotation

### Authorization Security

1. **Principle of Least Privilege**: Grant minimum required permissions
2. **Resource Ownership**: Check resource ownership before permissions
3. **Input Validation**: Validate all authorization inputs
4. **Audit Logging**: Log all authorization decisions
5. **Regular Review**: Periodically review role assignments

### Redis Security

1. **Authentication**: Use Redis AUTH if available
2. **Network Security**: Secure Redis network access
3. **Key Expiration**: Set appropriate TTL for all keys
4. **Monitoring**: Monitor Redis for unusual activity

## Troubleshooting

### Common Issues

1. **Token Validation Failures**
   ```
   Error: Token validation failed: token has been revoked
   ```
   **Solution**: Check Redis connectivity and token revocation status

2. **Permission Denied**
   ```
   Error: Insufficient permissions for resource access
   ```
   **Solution**: Verify user role and required permissions

3. **Redis Connection Issues**
   ```
   Warning: Failed to connect to Redis, JWT features will be limited
   ```
   **Solution**: Check Redis configuration and connectivity

### Debug Mode

Enable debug logging to troubleshoot authentication issues:

```bash
LOG_LEVEL=debug
```

This will provide detailed logs for:
- Token validation steps
- Permission checks
- Redis operations
- Session management

For more information, see the [Configuration Guide](configuration.md) and [Security Guide](security.md).
