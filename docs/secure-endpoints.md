# Secure Endpoints Guide

This document provides a comprehensive overview of the secured API endpoints implemented in the Great Nigeria Library Foundation backend.

## Overview

All API endpoints have been secured using the enhanced JWT authentication and authorization system with role-based access control (RBAC) and granular permissions.

### Security Layers

1. **Authentication**: JWT token validation with enhanced security checks
2. **Authorization**: Role-based and permission-based access control
3. **Resource Ownership**: Users can access/modify their own resources
4. **Audit Logging**: All access attempts are logged with full context

## Authentication Service Endpoints

### Public Endpoints (No Authentication Required)

```
POST /auth/register              - User registration
POST /auth/login                 - User login
POST /auth/refresh-token         - Token refresh
POST /auth/password/reset        - Password reset request
POST /auth/password/reset/confirm - Password reset confirmation
GET  /auth/oauth/:provider       - OAuth login
GET  /auth/oauth/:provider/callback - OAuth callback
POST /auth/email/verify/send     - Send email verification
POST /auth/email/verify/confirm  - Confirm email verification
POST /auth/email/verify/resend   - Resend verification email
POST /content/check-access       - Check content access (public endpoint)
```

### User Management Endpoints (Authentication Required)

```
GET    /users/:id                - Get user details
PATCH  /users/:id                - Update user (owner or admin)
GET    /users/:id/profile        - Get user profile
```

**Security**: 
- Requires valid JWT token
- PATCH requires resource ownership or `user:update_profile` permission

### Account Management Endpoints (Authentication Required)

```
DELETE /account/delete           - Delete account (requires delete permission)
GET    /account/2fa/status       - Get 2FA status
POST   /account/2fa/setup        - Setup 2FA
POST   /account/2fa/verify       - Verify 2FA
POST   /account/2fa/enable       - Enable 2FA
POST   /account/2fa/disable      - Disable 2FA
POST   /account/2fa/backup-codes - Generate backup codes
POST   /account/2fa/validate-backup - Validate backup code
GET    /account/sessions         - Get user sessions
POST   /account/sessions/revoke  - Revoke session
POST   /account/sessions/revoke-all - Revoke all sessions
GET    /account/privacy          - Get privacy settings
PUT    /account/privacy          - Update privacy settings
GET    /account/permissions      - Get user permissions
```

**Security**: 
- Requires valid JWT token
- DELETE operations require `user:delete_profile` permission

### Moderator Endpoints (Moderator Role Required)

```
GET    /moderator/tools          - Get moderator tools
GET    /moderator/content/access - Get content access (manage content permission)
GET    /moderator/content/rules  - Get content rules (manage content permission)
POST   /moderator/content/permissions - Grant user permission (manage content permission)
DELETE /moderator/content/permissions/:id - Revoke user permission (manage content permission)
```

**Security**: 
- Requires valid JWT token
- Requires Moderator role (level 2) or higher
- Content operations require `admin:manage_content` permission

### Admin Endpoints (Admin Role Required)

```
GET    /admin/users              - List users (manage users permission)
PATCH  /admin/users/:id/role     - Update user role (manage users permission)
GET    /admin/users/role/:role   - Get users by role (manage users permission)
POST   /admin/content/access     - Set content access (manage content permission)
POST   /admin/content/rules      - Create content rule (manage content permission)
PUT    /admin/content/rules      - Update content rule (manage content permission)
DELETE /admin/content/rules/:id  - Delete content rule (manage content permission)
POST   /admin/sessions/maintenance - Perform maintenance
```

**Security**: 
- Requires valid JWT token
- Requires Admin role (level 3) or higher
- User management requires `admin:manage_users` permission
- Content management requires `admin:manage_content` permission

## Content Service Endpoints

### Public Content Endpoints (No Authentication Required)

```
GET /public/books               - Get all books (public access)
GET /public/books/:id           - Get book by ID (public access)
GET /public/books/:id/chapters  - Get book chapters (public access)
GET /public/books/chapters/:id  - Get chapter by ID (public access)
GET /public/books/sections/:id  - Get section by ID (public access)
GET /public/feedback/summary    - Get content feedback summary
```

### Protected Content Endpoints (Authentication + Read Permission Required)

```
GET /content/books              - Get all books (read content permission)
GET /content/books/:id          - Get book by ID (read content permission)
GET /content/books/:id/chapters - Get book chapters (read content permission)
GET /content/books/chapters/:id - Get chapter by ID (read content permission)
GET /content/books/sections/:id - Get section by ID (read content permission)
```

**Security**: 
- Requires valid JWT token
- Requires `content:read` permission

### User Content Interaction Endpoints (Authentication Required)

```
POST   /user/books/:id/progress     - Update reading progress (read content permission)
GET    /user/books/:id/progress     - Get reading progress (read content permission)
POST   /user/books/:id/bookmarks    - Create bookmark (read content permission)
GET    /user/books/:id/bookmarks    - Get bookmarks (read content permission)
DELETE /user/books/:id/bookmarks/:bookmarkId - Delete bookmark (read content permission)
POST   /user/books/:id/notes        - Create note (read content permission)
GET    /user/books/:id/notes        - Get notes (read content permission)
GET    /user/notes/:noteId          - Get note by ID (owner or read content permission)
PUT    /user/notes/:noteId          - Update note (owner or update content permission)
DELETE /user/notes/:noteId          - Delete note (owner or delete content permission)
GET    /user/notes/categories       - Get note categories (read content permission)
POST   /user/notes/export           - Export notes (read content permission)
```

**Security**: 
- Requires valid JWT token
- Most operations require `content:read` permission
- Note modification requires resource ownership or appropriate content permissions

### Content Creation Endpoints (Create Permission Required)

```
POST /create/books              - Create book (create content permission)
```

**Security**: 
- Requires valid JWT token
- Requires `content:create` permission

### Content Management Endpoints (Management Permissions Required)

```
PUT    /manage/books/:id         - Update book (update/delete/publish content permission)
DELETE /manage/books/:id         - Delete book (update/delete/publish content permission)
```

**Security**: 
- Requires valid JWT token
- Requires any of: `content:update`, `content:delete`, or `content:publish` permissions

### Admin Content Endpoints (Admin Role Required)

```
GET  /admin/content/stats       - Get content statistics
POST /admin/content/publish/:id - Publish content (publish content permission)
```

**Security**: 
- Requires valid JWT token
- Requires Admin role (level 3) or higher
- Publishing requires `content:publish` permission

## Discussion Service Endpoints

### Public Discussion Endpoints (Read Permission Required)

```
GET /public/discussions         - List discussions (read discussion permission)
GET /public/discussions/:id     - Get discussion (read discussion permission)
GET /public/discussions/:id/comments - List comments (read discussion permission)
```

**Security**: 
- Requires `discussion:read` permission (even for public access)

### Protected Discussion Endpoints (Authentication + Read Permission Required)

```
GET /discussions/               - List discussions (read discussion permission)
GET /discussions/:id            - Get discussion (read discussion permission)
GET /discussions/:id/comments   - List comments (read discussion permission)
```

**Security**: 
- Requires valid JWT token
- Requires `discussion:read` permission

### Discussion Participation Endpoints (Authentication Required)

```
POST   /participate/discussions     - Create discussion (create discussion permission)
PATCH  /participate/discussions/:id - Update discussion (owner or update discussion permission)
DELETE /participate/discussions/:id - Delete discussion (owner or delete discussion permission)
POST   /participate/discussions/:id/comments - Create comment (create discussion permission)
PATCH  /participate/comments/:id    - Update comment (owner or update discussion permission)
DELETE /participate/comments/:id    - Delete comment (owner or delete discussion permission)
POST   /participate/discussions/:id/like - Like discussion (read discussion permission)
DELETE /participate/discussions/:id/like - Unlike discussion (read discussion permission)
POST   /participate/comments/:id/like - Like comment (read discussion permission)
```

**Security**: 
- Requires valid JWT token
- Create operations require `discussion:create` permission
- Update/Delete operations require resource ownership or appropriate discussion permissions
- Like operations require `discussion:read` permission

### Discussion Moderation Endpoints (Moderator Role Required)

```
DELETE /moderate/discussions/:id    - Delete discussion (moderate discussion permission)
DELETE /moderate/comments/:id       - Delete comment (moderate discussion permission)
POST   /moderate/discussions/:id/lock - Lock discussion (moderate discussion permission)
POST   /moderate/discussions/:id/pin  - Pin discussion (moderate discussion permission)
```

**Security**: 
- Requires valid JWT token
- Requires Moderator role (level 2) or higher
- Requires `discussion:moderate` permission

### Admin Discussion Endpoints (Admin Role Required)

```
GET /admin/discussions/stats    - Get discussion statistics
GET /admin/discussions/reports  - Get discussion reports
```

**Security**: 
- Requires valid JWT token
- Requires Admin role (level 3) or higher

## Foundation App Endpoints

### Public Auth Endpoints (No Authentication Required)

```
POST /api/v1/auth/register      - User registration
POST /api/v1/auth/login         - User login
```

### Protected Auth Endpoints (Authentication Required)

```
POST /api/v1/auth/logout        - User logout
GET  /api/v1/auth/profile       - Get user profile
PUT  /api/v1/auth/profile       - Update user profile (update profile permission)
```

**Security**: 
- Requires valid JWT token
- Profile update requires `user:update_profile` permission

### Public Content Endpoints (No Authentication Required)

```
GET /api/v1/content/public/books - Get demo books (public access)
```

### Protected Content Endpoints (Authentication + Read Permission Required)

```
GET /api/v1/content/books       - Get demo books (read content permission)
```

**Security**: 
- Requires valid JWT token
- Requires `content:read` permission

### Admin Content Endpoints (Admin Role Required)

```
GET  /api/v1/content/admin/stats - Get content statistics
POST /api/v1/content/admin/books - Create content (create content permission)
```

**Security**: 
- Requires valid JWT token
- Requires Admin role (level 3) or higher
- Content creation requires `content:create` permission

## Security Features

### JWT Token Security

- **Enhanced Claims**: Extended JWT tokens with permissions, device tracking, and security metadata
- **Token Revocation**: Real-time token blacklisting via Redis
- **Session Management**: Device tracking and session lifecycle management
- **IP Validation**: Optional IP address validation for enhanced security

### Role-Based Access Control

- **Guest (0)**: Read-only access to public content
- **User (1)**: Basic user operations and content interaction
- **Moderator (2)**: Content moderation and community management
- **Admin (3)**: System administration and user management
- **SuperAdmin (4)**: Full system access

### Permission-Based Authorization

- **Granular Permissions**: 20+ specific permissions across different domains
- **Resource Ownership**: Users can access/modify their own resources
- **Permission Inheritance**: Higher roles inherit lower role permissions
- **Dynamic Validation**: Runtime permission checking and validation

### Audit and Logging

- **Request Logging**: All API requests logged with user context
- **Authentication Events**: Login, logout, token refresh events logged
- **Authorization Decisions**: Permission checks and access denials logged
- **Security Events**: Failed authentication attempts and suspicious activity logged

## Error Responses

All secured endpoints return standardized error responses:

```json
{
  "success": false,
  "message": "Access denied",
  "error": "Insufficient permissions",
  "code": "FORBIDDEN",
  "timestamp": "2023-12-07T10:30:00Z"
}
```

### Common Error Codes

- `UNAUTHORIZED` (401): Missing or invalid authentication token
- `FORBIDDEN` (403): Insufficient permissions for the requested operation
- `TOKEN_EXPIRED` (401): JWT token has expired
- `TOKEN_REVOKED` (401): JWT token has been revoked
- `INVALID_ROLE` (403): User role insufficient for the requested operation

## Best Practices

### For API Consumers

1. **Always Include Authorization Header**: `Authorization: Bearer <token>`
2. **Handle Token Expiration**: Implement automatic token refresh
3. **Respect Rate Limits**: Follow API rate limiting guidelines
4. **Check Permissions**: Verify user permissions before making requests
5. **Handle Errors Gracefully**: Implement proper error handling for all security responses

### For Developers

1. **Principle of Least Privilege**: Grant minimum required permissions
2. **Resource Ownership Checks**: Always verify resource ownership before permissions
3. **Input Validation**: Validate all inputs before processing
4. **Audit Logging**: Log all security-relevant events
5. **Regular Security Reviews**: Periodically review and update security configurations

## Testing Secured Endpoints

### Authentication Testing

```bash
# 1. Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword123",
    "full_name": "Test User"
  }'

# 2. Login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "securepassword123"
  }'

# Response will include access_token and refresh_token
```

### Authorization Testing

```bash
# 3. Access protected endpoint with token
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <your_access_token>"

# 4. Test permission-based access
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <your_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Updated Test User"
  }'

# 5. Test role-based access (should fail for regular user)
curl -X GET http://localhost:8080/admin/users \
  -H "Authorization: Bearer <your_access_token>"
```

### Error Testing

```bash
# 6. Test without authentication (should return 401)
curl -X GET http://localhost:8080/api/v1/auth/profile

# 7. Test with invalid token (should return 401)
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer invalid_token"

# 8. Test insufficient permissions (should return 403)
curl -X DELETE http://localhost:8080/admin/users/123 \
  -H "Authorization: Bearer <regular_user_token>"
```

For more information, see the [JWT Authentication & Authorization Guide](jwt-authentication-authorization.md) and [Configuration Guide](configuration.md).
