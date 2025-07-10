package models

import (
	"time"

	"gorm.io/gorm"
)

// MembershipLevel represents user membership levels
type MembershipLevel string

const (
	MembershipBasic   MembershipLevel = "basic"
	MembershipPremium MembershipLevel = "premium"
	MembershipVIP     MembershipLevel = "vip"
)

// TrustLevel represents user trust levels
type TrustLevel string

const (
	TrustLevelNewUser    TrustLevel = "new_user"
	TrustLevelTrusted    TrustLevel = "trusted"
	TrustLevelVerified   TrustLevel = "verified"
	TrustLevelExpert     TrustLevel = "expert"
	TrustLevelModerator  TrustLevel = "moderator"
)

// User represents the unified user model across all services
type User struct {
	gorm.Model
	
	// Basic Information
	Username        string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email           string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password        string    `json:"-" gorm:"size:255;not null"` // Hidden from JSON
	FullName        string    `json:"full_name" gorm:"size:255;not null"`
	DisplayName     string    `json:"display_name" gorm:"size:100"`
	
	// Profile Information
	Bio             string    `json:"bio" gorm:"type:text"`
	ProfileImage    string    `json:"profile_image" gorm:"size:500"`
	ProfileImageURL string    `json:"profile_image_url" gorm:"size:500"` // For discussion module compatibility
	
	// Account Status
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	IsVerified      bool      `json:"is_verified" gorm:"default:false"`
	IsOAuth         bool      `json:"is_oauth" gorm:"default:false"`
	OAuthProvider   string    `json:"oauth_provider" gorm:"size:50"`
	OAuthID         string    `json:"oauth_id" gorm:"size:255"`
	
	// Membership & Points
	MembershipLevel MembershipLevel `json:"membership_level" gorm:"default:'basic'"`
	PointsBalance   int             `json:"points_balance" gorm:"default:0"`
	Reputation      int             `json:"reputation" gorm:"default:0"` // For discussion module
	
	// Activity Tracking
	LastLogin       time.Time `json:"last_login"`
	LastSeenAt      time.Time `json:"last_seen_at"` // For discussion module compatibility
	
	// Role & Permissions
	Role            int       `json:"role" gorm:"default:1"` // 1=user, 2=moderator, 3=admin
	
	// Timestamps (gorm.Model provides ID, CreatedAt, UpdatedAt, DeletedAt)
}

// UserRegisterRequest represents the registration request
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required,min=2,max=255"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserLoginRequest represents the login request
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest represents the user update request
type UserUpdateRequest struct {
	FullName    string `json:"full_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Username    string `json:"username,omitempty"`
}

// UserResponse represents the user response (without sensitive data)
type UserResponse struct {
	ID              uint            `json:"id"`
	Username        string          `json:"username"`
	Email           string          `json:"email"`
	FullName        string          `json:"full_name"`
	DisplayName     string          `json:"display_name"`
	Bio             string          `json:"bio"`
	ProfileImage    string          `json:"profile_image"`
	ProfileImageURL string          `json:"profile_image_url"`
	IsActive        bool            `json:"is_active"`
	IsVerified      bool            `json:"is_verified"`
	MembershipLevel MembershipLevel `json:"membership_level"`
	PointsBalance   int             `json:"points_balance"`
	Reputation      int             `json:"reputation"`
	Role            int             `json:"role"`
	LastLogin       time.Time       `json:"last_login"`
	LastSeenAt      time.Time       `json:"last_seen_at"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// UserProfile represents a public user profile
type UserProfile struct {
	ID              uint            `json:"id"`
	Username        string          `json:"username"`
	DisplayName     string          `json:"display_name"`
	Bio             string          `json:"bio"`
	ProfileImage    string          `json:"profile_image"`
	ProfileImageURL string          `json:"profile_image_url"`
	MembershipLevel MembershipLevel `json:"membership_level"`
	PointsBalance   int             `json:"points_balance"`
	Reputation      int             `json:"reputation"`
	CreatedAt       time.Time       `json:"created_at"`
}

// UserWithStats represents a user with additional statistics
type UserWithStats struct {
	User
	PostCount    int `json:"post_count"`
	CommentCount int `json:"comment_count"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		FullName:        u.FullName,
		DisplayName:     u.DisplayName,
		Bio:             u.Bio,
		ProfileImage:    u.ProfileImage,
		ProfileImageURL: u.ProfileImageURL,
		IsActive:        u.IsActive,
		IsVerified:      u.IsVerified,
		MembershipLevel: u.MembershipLevel,
		PointsBalance:   u.PointsBalance,
		Reputation:      u.Reputation,
		Role:            u.Role,
		LastLogin:       u.LastLogin,
		LastSeenAt:      u.LastSeenAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

// ToProfile converts User to UserProfile
func (u *User) ToProfile() UserProfile {
	return UserProfile{
		ID:              u.ID,
		Username:        u.Username,
		DisplayName:     u.DisplayName,
		Bio:             u.Bio,
		ProfileImage:    u.ProfileImage,
		ProfileImageURL: u.ProfileImageURL,
		MembershipLevel: u.MembershipLevel,
		PointsBalance:   u.PointsBalance,
		Reputation:      u.Reputation,
		CreatedAt:       u.CreatedAt,
	}
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"size:255;not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// EmailVerificationToken represents an email verification token
type EmailVerificationToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"size:255;not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// ContentAccess represents content access control
type ContentAccess struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ContentType string    `json:"content_type" gorm:"size:50;not null"`
	ContentID   uint      `json:"content_id" gorm:"not null"`
	AccessLevel string    `json:"access_level" gorm:"size:50;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ContentAccessRule represents content access rules
type ContentAccessRule struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ContentType string    `json:"content_type" gorm:"size:50;not null"`
	RuleName    string    `json:"rule_name" gorm:"size:100;not null"`
	Conditions  string    `json:"conditions" gorm:"type:json"`
	AccessLevel string    `json:"access_level" gorm:"size:50;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserContentPermission represents user-specific content permissions
type UserContentPermission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	ContentType string    `json:"content_type" gorm:"size:50;not null"`
	ContentID   uint      `json:"content_id" gorm:"not null"`
	Permission  string    `json:"permission" gorm:"size:50;not null"`
	GrantedBy   uint      `json:"granted_by" gorm:"not null"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// UserPrivacySettings represents user privacy settings
type UserPrivacySettings struct {
	ID                    uint      `json:"id" gorm:"primaryKey"`
	UserID                uint      `json:"user_id" gorm:"not null;uniqueIndex"`
	ProfileVisibility     string    `json:"profile_visibility" gorm:"size:50;default:'public'"`
	ActivityVisibility    string    `json:"activity_visibility" gorm:"size:50;default:'friends'"`
	ContactInfoVisibility string    `json:"contact_info_visibility" gorm:"size:50;default:'private'"`
	SearchVisibility      bool      `json:"search_visibility" gorm:"default:true"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
