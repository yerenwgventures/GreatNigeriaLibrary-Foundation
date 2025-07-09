package models

import (
	"time"
)

// Session represents a user session
type Session struct {
	ID           string    `json:"id" gorm:"primaryKey;size:255"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	DeviceType   string    `json:"device_type" gorm:"size:50"`
	DeviceInfo   string    `json:"device_info" gorm:"size:500"`
	IPAddress    string    `json:"ip_address" gorm:"size:45"`
	UserAgent    string    `json:"user_agent" gorm:"size:1000"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	LastActivity time.Time `json:"last_activity"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// TwoFactorAuth represents 2FA settings for a user
type TwoFactorAuth struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;uniqueIndex"`
	Secret       string    `json:"-" gorm:"size:255;not null"` // Hidden from JSON
	BackupCodes  string    `json:"-" gorm:"type:text"`         // Hidden from JSON
	IsEnabled    bool      `json:"is_enabled" gorm:"default:false"`
	LastUsedAt   time.Time `json:"last_used_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// TwoFactorSetupRequest represents a 2FA setup request
type TwoFactorSetupRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// TwoFactorVerifyRequest represents a 2FA verification request
type TwoFactorVerifyRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// TwoFactorResponse represents a 2FA response
type TwoFactorResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// VerificationStatus represents user verification status
type VerificationStatus struct {
	EmailVerified    bool `json:"email_verified"`
	PhoneVerified    bool `json:"phone_verified"`
	IdentityVerified bool `json:"identity_verified"`
	AddressVerified  bool `json:"address_verified"`
}

// VerificationRequest represents a verification request
type VerificationRequest struct {
	UserID    uint      `json:"user_id" gorm:"not null"`
	Type      string    `json:"type" gorm:"size:50;not null"`
	Status    string    `json:"status" gorm:"size:50;default:'pending'"`
	Documents []string  `json:"documents" gorm:"type:json"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserBadge represents a user badge
type UserBadge struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	BadgeType   string    `json:"badge_type" gorm:"size:100;not null"`
	BadgeName   string    `json:"badge_name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	EarnedAt    time.Time `json:"earned_at"`
	IsVisible   bool      `json:"is_visible" gorm:"default:true"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// UserTrustLevel represents a user's trust level
type UserTrustLevel struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;uniqueIndex"`
	Level     TrustLevel `json:"level" gorm:"default:'new_user'"`
	Points    int        `json:"points" gorm:"default:0"`
	UpdatedAt time.Time  `json:"updated_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// OAuthUserInfo represents OAuth user information
type OAuthUserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirmRequest represents a password reset confirmation
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// EmailVerificationRequest represents an email verification request
type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// EmailVerificationConfirmRequest represents email verification confirmation
type EmailVerificationConfirmRequest struct {
	Token string `json:"token" binding:"required"`
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthResponse represents a generic auth response
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
