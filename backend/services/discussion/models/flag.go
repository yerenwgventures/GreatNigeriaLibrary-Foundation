package models

import (
	"time"

	"gorm.io/gorm"
)

// ContentFlagType defines the types of content that can be flagged
type ContentFlagType string

const (
	// FlagTypeHarassment for harassment or bullying
	FlagTypeHarassment ContentFlagType = "harassment"
	
	// FlagTypeHateSpeech for hate speech
	FlagTypeHateSpeech ContentFlagType = "hate_speech"
	
	// FlagTypeSpam for spam content
	FlagTypeSpam ContentFlagType = "spam"
	
	// FlagTypeOffTopic for off-topic content
	FlagTypeOffTopic ContentFlagType = "off_topic"
	
	// FlagTypeMisleading for misleading content
	FlagTypeMisleading ContentFlagType = "misleading"
	
	// FlagTypeInappropriate for inappropriate content
	FlagTypeInappropriate ContentFlagType = "inappropriate"
	
	// FlagTypeViolentContent for violent content
	FlagTypeViolentContent ContentFlagType = "violent_content"
	
	// FlagTypeIllegalContent for illegal content
	FlagTypeIllegalContent ContentFlagType = "illegal_content"
	
	// FlagTypeOther for other reasons
	FlagTypeOther ContentFlagType = "other"
)

// FlagStatus defines the possible statuses for a content flag
type FlagStatus string

const (
	// FlagStatusPending for pending flags
	FlagStatusPending FlagStatus = "pending"
	
	// FlagStatusReviewed for reviewed flags
	FlagStatusReviewed FlagStatus = "reviewed"
	
	// FlagStatusApproved for approved flags
	FlagStatusApproved FlagStatus = "approved"
	
	// FlagStatusRejected for rejected flags
	FlagStatusRejected FlagStatus = "rejected"
)

// ContentFlag represents a flag on a content item
type ContentFlag struct {
	gorm.Model
	ContentType      string         `json:"contentType" gorm:"index:idx_content_flag"`
	ContentID        uint           `json:"contentId" gorm:"index:idx_content_flag"`
	FlagType         ContentFlagType `json:"flagType"`
	Description      string         `json:"description" gorm:"type:text"`
	UserID           uint           `json:"userId"`
	Status           FlagStatus     `json:"status" gorm:"default:pending"`
	AssignedTo       *uint          `json:"assignedTo"`
	ReviewedBy       *uint          `json:"reviewedBy"`
	ReviewedAt       *time.Time     `json:"reviewedAt"`
	ActionTaken      string         `json:"actionTaken"`
	Notes            string         `json:"notes" gorm:"type:text"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

// ContentModerationStatusType defines the possible moderation statuses for content
type ContentModerationStatusType string

const (
	// ModerationStatusPending for pending moderation
	ModerationStatusPending ContentModerationStatusType = "pending"
	
	// ModerationStatusApproved for approved content
	ModerationStatusApproved ContentModerationStatusType = "approved"
	
	// ModerationStatusRejected for rejected content
	ModerationStatusRejected ContentModerationStatusType = "rejected"
	
	// ModerationStatusHidden for hidden content
	ModerationStatusHidden ContentModerationStatusType = "hidden"
)

// ContentModerationStatus represents the moderation status of content
type ContentModerationStatus struct {
	gorm.Model
	ContentType      string                    `json:"contentType" gorm:"index:idx_content_mod_status"`
	ContentID        uint                      `json:"contentId" gorm:"index:idx_content_mod_status"`
	Status           ContentModerationStatusType `json:"status" gorm:"default:pending"`
	ModeratorID      *uint                     `json:"moderatorId"`
	Reason           string                    `json:"reason" gorm:"type:text"`
	Notes            string                    `json:"notes" gorm:"type:text"`
	UserNotified     bool                      `json:"userNotified" gorm:"default:false"`
	CreatedAt        time.Time                 `json:"createdAt"`
	UpdatedAt        time.Time                 `json:"updatedAt"`
}

// UserPenaltyType defines the type of penalty applied to a user
type UserPenaltyType string

const (
	// PenaltyTypeWarning for warnings
	PenaltyTypeWarning UserPenaltyType = "warning"
	
	// PenaltyTypeSuspension for temporary suspensions
	PenaltyTypeSuspension UserPenaltyType = "suspension"
	
	// PenaltyTypeBan for permanent bans
	PenaltyTypeBan UserPenaltyType = "ban"
	
	// PenaltyTypeRestriction for restricted privileges
	PenaltyTypeRestriction UserPenaltyType = "restriction"
)

// UserPenalty represents a penalty applied to a user
type UserPenalty struct {
	gorm.Model
	UserID           uint           `json:"userId" gorm:"index"`
	PenaltyType      UserPenaltyType `json:"penaltyType"`
	Reason           string         `json:"reason" gorm:"type:text"`
	Description      string         `json:"description" gorm:"type:text"`
	ModeratorID      uint           `json:"moderatorId"`
	Duration         *int           `json:"duration"`
	ExpiresAt        *time.Time     `json:"expiresAt"`
	IsActive         bool           `json:"isActive" gorm:"default:true"`
	RelatedContentType *string       `json:"relatedContentType"`
	RelatedContentID *uint          `json:"relatedContentId"`
	Notes            string         `json:"notes" gorm:"type:text"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}