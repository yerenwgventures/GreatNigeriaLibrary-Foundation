package models

import (
	"time"
)

// PointsTransactionType defines the type of points transaction
type PointsTransactionType string

const (
	// Transaction types for points
	PointsEarned  PointsTransactionType = "earned"
	PointsSpent   PointsTransactionType = "spent"
	PointsExpired PointsTransactionType = "expired"
	PointsAdjusted PointsTransactionType = "adjusted" // For manual adjustments by admins
)

// PointsTransaction represents a transaction in a user's points account
type PointsTransaction struct {
	ID             uint                 `json:"id" gorm:"primaryKey"`
	UserID         uint                 `json:"user_id" gorm:"index;not null"`
	Points         int                  `json:"points" gorm:"not null"` // Positive for earned, negative for spent
	TransactionType PointsTransactionType `json:"transaction_type" gorm:"type:varchar(50);not null;index"`
	Description    string               `json:"description" gorm:"size:255"`
	ReferenceID    *uint                `json:"reference_id,omitempty" gorm:"index"` // Reference to related entity (e.g., book, discussion)
	ReferenceType  string               `json:"reference_type,omitempty" gorm:"size:100"`
	CreatedAt      time.Time            `json:"created_at"`
}

// UserProfile extends the core user model with additional profile data
type UserProfile struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	DisplayName      string    `json:"display_name" gorm:"size:100"`
	Bio              string    `json:"bio" gorm:"type:text"`
	ProfileImageURL  string    `json:"profile_image_url" gorm:"size:255"`
	Points           int       `json:"points" gorm:"default:0"`
	MembershipLevel  string    `json:"membership_level" gorm:"size:50;default:'standard'"`
	LastActive       time.Time `json:"last_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName overrides the table name for PointsTransaction
func (PointsTransaction) TableName() string {
	return "points_transactions"
}

// TableName overrides the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profiles"
}