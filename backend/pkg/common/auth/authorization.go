package auth

import (
	"errors"
	"fmt"
	"strings"
)

// Role represents user roles in the system
type Role int

const (
	RoleGuest Role = iota
	RoleUser
	RoleModerator
	RoleAdmin
	RoleSuperAdmin
)

// String returns the string representation of a role
func (r Role) String() string {
	switch r {
	case RoleGuest:
		return "guest"
	case RoleUser:
		return "user"
	case RoleModerator:
		return "moderator"
	case RoleAdmin:
		return "admin"
	case RoleSuperAdmin:
		return "superadmin"
	default:
		return "unknown"
	}
}

// Permission represents a specific permission in the system
type Permission string

const (
	// User permissions
	PermissionReadProfile   Permission = "user:read_profile"
	PermissionUpdateProfile Permission = "user:update_profile"
	PermissionDeleteProfile Permission = "user:delete_profile"

	// Content permissions
	PermissionReadContent   Permission = "content:read"
	PermissionCreateContent Permission = "content:create"
	PermissionUpdateContent Permission = "content:update"
	PermissionDeleteContent Permission = "content:delete"
	PermissionPublishContent Permission = "content:publish"

	// Discussion permissions
	PermissionReadDiscussion   Permission = "discussion:read"
	PermissionCreateDiscussion Permission = "discussion:create"
	PermissionUpdateDiscussion Permission = "discussion:update"
	PermissionDeleteDiscussion Permission = "discussion:delete"
	PermissionModerateDiscussion Permission = "discussion:moderate"

	// Group permissions
	PermissionReadGroup   Permission = "group:read"
	PermissionCreateGroup Permission = "group:create"
	PermissionUpdateGroup Permission = "group:update"
	PermissionDeleteGroup Permission = "group:delete"
	PermissionManageGroup Permission = "group:manage"

	// Admin permissions
	PermissionManageUsers    Permission = "admin:manage_users"
	PermissionManageContent  Permission = "admin:manage_content"
	PermissionManageSystem   Permission = "admin:manage_system"
	PermissionViewAnalytics  Permission = "admin:view_analytics"
	PermissionManageSettings Permission = "admin:manage_settings"
)

// AuthorizationManager handles role-based access control
type AuthorizationManager struct {
	rolePermissions map[Role][]Permission
}

// NewAuthorizationManager creates a new authorization manager
func NewAuthorizationManager() *AuthorizationManager {
	am := &AuthorizationManager{
		rolePermissions: make(map[Role][]Permission),
	}
	am.initializeRolePermissions()
	return am
}

// initializeRolePermissions sets up default role-permission mappings
func (am *AuthorizationManager) initializeRolePermissions() {
	// Guest permissions (minimal)
	am.rolePermissions[RoleGuest] = []Permission{
		PermissionReadContent,
		PermissionReadDiscussion,
		PermissionReadGroup,
	}

	// User permissions (basic user actions)
	am.rolePermissions[RoleUser] = []Permission{
		// Inherit guest permissions
		PermissionReadContent,
		PermissionReadDiscussion,
		PermissionReadGroup,
		
		// User-specific permissions
		PermissionReadProfile,
		PermissionUpdateProfile,
		PermissionCreateContent,
		PermissionUpdateContent, // Own content only
		PermissionCreateDiscussion,
		PermissionUpdateDiscussion, // Own discussions only
		PermissionCreateGroup,
		PermissionUpdateGroup, // Own groups only
	}

	// Moderator permissions (content moderation)
	am.rolePermissions[RoleModerator] = []Permission{
		// Inherit user permissions
		PermissionReadContent,
		PermissionReadDiscussion,
		PermissionReadGroup,
		PermissionReadProfile,
		PermissionUpdateProfile,
		PermissionCreateContent,
		PermissionUpdateContent,
		PermissionCreateDiscussion,
		PermissionUpdateDiscussion,
		PermissionCreateGroup,
		PermissionUpdateGroup,
		
		// Moderator-specific permissions
		PermissionDeleteContent,   // Can delete inappropriate content
		PermissionDeleteDiscussion, // Can delete inappropriate discussions
		PermissionModerateDiscussion, // Can moderate discussions
		PermissionManageGroup,     // Can manage groups
	}

	// Admin permissions (system administration)
	am.rolePermissions[RoleAdmin] = []Permission{
		// Inherit moderator permissions
		PermissionReadContent,
		PermissionReadDiscussion,
		PermissionReadGroup,
		PermissionReadProfile,
		PermissionUpdateProfile,
		PermissionCreateContent,
		PermissionUpdateContent,
		PermissionDeleteContent,
		PermissionCreateDiscussion,
		PermissionUpdateDiscussion,
		PermissionDeleteDiscussion,
		PermissionModerateDiscussion,
		PermissionCreateGroup,
		PermissionUpdateGroup,
		PermissionDeleteGroup,
		PermissionManageGroup,
		
		// Admin-specific permissions
		PermissionPublishContent,
		PermissionManageUsers,
		PermissionManageContent,
		PermissionViewAnalytics,
		PermissionManageSettings,
	}

	// SuperAdmin permissions (full system access)
	am.rolePermissions[RoleSuperAdmin] = []Permission{
		// All permissions
		PermissionReadContent,
		PermissionCreateContent,
		PermissionUpdateContent,
		PermissionDeleteContent,
		PermissionPublishContent,
		PermissionReadDiscussion,
		PermissionCreateDiscussion,
		PermissionUpdateDiscussion,
		PermissionDeleteDiscussion,
		PermissionModerateDiscussion,
		PermissionReadGroup,
		PermissionCreateGroup,
		PermissionUpdateGroup,
		PermissionDeleteGroup,
		PermissionManageGroup,
		PermissionReadProfile,
		PermissionUpdateProfile,
		PermissionDeleteProfile,
		PermissionManageUsers,
		PermissionManageContent,
		PermissionManageSystem,
		PermissionViewAnalytics,
		PermissionManageSettings,
	}
}

// HasPermission checks if a role has a specific permission
func (am *AuthorizationManager) HasPermission(role Role, permission Permission) bool {
	permissions, exists := am.rolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if a role has any of the specified permissions
func (am *AuthorizationManager) HasAnyPermission(role Role, permissions []Permission) bool {
	for _, permission := range permissions {
		if am.HasPermission(role, permission) {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if a role has all of the specified permissions
func (am *AuthorizationManager) HasAllPermissions(role Role, permissions []Permission) bool {
	for _, permission := range permissions {
		if !am.HasPermission(role, permission) {
			return false
		}
	}
	return true
}

// GetRolePermissions returns all permissions for a role
func (am *AuthorizationManager) GetRolePermissions(role Role) []Permission {
	permissions, exists := am.rolePermissions[role]
	if !exists {
		return []Permission{}
	}
	return permissions
}

// CanAccessResource checks if a user can access a specific resource
func (am *AuthorizationManager) CanAccessResource(userRole Role, userID uint, resource string, action string, resourceOwnerID uint) bool {
	permission := Permission(fmt.Sprintf("%s:%s", resource, action))
	
	// Check if user has the general permission
	if am.HasPermission(userRole, permission) {
		return true
	}

	// Check ownership-based permissions for certain actions
	if userID == resourceOwnerID {
		switch action {
		case "update", "delete":
			// Users can update/delete their own content
			readPermission := Permission(fmt.Sprintf("%s:read", resource))
			return am.HasPermission(userRole, readPermission)
		}
	}

	return false
}

// ValidatePermissions validates a list of permission strings
func (am *AuthorizationManager) ValidatePermissions(permissionStrings []string) ([]Permission, error) {
	var permissions []Permission
	var invalidPermissions []string

	for _, permStr := range permissionStrings {
		permission := Permission(permStr)
		
		// Check if permission is valid by checking if any role has it
		isValid := false
		for _, rolePerms := range am.rolePermissions {
			for _, p := range rolePerms {
				if p == permission {
					isValid = true
					break
				}
			}
			if isValid {
				break
			}
		}

		if isValid {
			permissions = append(permissions, permission)
		} else {
			invalidPermissions = append(invalidPermissions, permStr)
		}
	}

	if len(invalidPermissions) > 0 {
		return permissions, fmt.Errorf("invalid permissions: %s", strings.Join(invalidPermissions, ", "))
	}

	return permissions, nil
}

// RoleFromInt converts an integer to a Role
func RoleFromInt(roleInt int) (Role, error) {
	switch roleInt {
	case 0:
		return RoleGuest, nil
	case 1:
		return RoleUser, nil
	case 2:
		return RoleModerator, nil
	case 3:
		return RoleAdmin, nil
	case 4:
		return RoleSuperAdmin, nil
	default:
		return RoleGuest, errors.New("invalid role")
	}
}

// RoleFromString converts a string to a Role
func RoleFromString(roleStr string) (Role, error) {
	switch strings.ToLower(roleStr) {
	case "guest":
		return RoleGuest, nil
	case "user":
		return RoleUser, nil
	case "moderator":
		return RoleModerator, nil
	case "admin":
		return RoleAdmin, nil
	case "superadmin":
		return RoleSuperAdmin, nil
	default:
		return RoleGuest, errors.New("invalid role string")
	}
}
