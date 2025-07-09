package models

import (
	"time"

	"gorm.io/gorm"
)

// ReportStatus defines the status of a content report
type ReportStatus string

const (
	// StatusPending for new reports awaiting review
	StatusPending ReportStatus = "pending"
	
	// StatusInReview for reports currently being reviewed
	StatusInReview ReportStatus = "in_review"
	
	// StatusResolved for reports that have been addressed
	StatusResolved ReportStatus = "resolved"
	
	// StatusRejected for reports that were determined not to violate guidelines
	StatusRejected ReportStatus = "rejected"
)

// ReportCategory defines the category of a content report
type ReportCategory string

const (
	// CategorySpam for spam reports
	CategorySpam ReportCategory = "spam"
	
	// CategoryHarassment for harassment reports
	CategoryHarassment ReportCategory = "harassment"
	
	// CategoryHateSpeech for hate speech reports
	CategoryHateSpeech ReportCategory = "hate_speech"
	
	// CategoryViolence for violence or threatening content reports
	CategoryViolence ReportCategory = "violence"
	
	// CategoryIllegalContent for illegal content reports
	CategoryIllegalContent ReportCategory = "illegal_content"
	
	// CategoryPrivacyViolation for privacy violation reports
	CategoryPrivacyViolation ReportCategory = "privacy_violation"
	
	// CategoryCopyright for copyright violation reports
	CategoryCopyright ReportCategory = "copyright"
	
	// CategoryMisinformation for misinformation reports
	CategoryMisinformation ReportCategory = "misinformation"
	
	// CategoryOther for other reports
	CategoryOther ReportCategory = "other"
)

// ReportResolutionType defines how a report was resolved
type ReportResolutionType string

const (
	// ResolutionNoAction for reports where no action was taken
	ResolutionNoAction ReportResolutionType = "no_action"
	
	// ResolutionWarning for reports where a warning was issued
	ResolutionWarning ReportResolutionType = "warning"
	
	// ResolutionContentRemoved for reports where content was removed
	ResolutionContentRemoved ReportResolutionType = "content_removed"
	
	// ResolutionContentEdited for reports where content was edited
	ResolutionContentEdited ReportResolutionType = "content_edited"
	
	// ResolutionUserSuspended for reports where user was suspended
	ResolutionUserSuspended ReportResolutionType = "user_suspended"
	
	// ResolutionUserBanned for reports where user was banned
	ResolutionUserBanned ReportResolutionType = "user_banned"
)

// ContentReport represents a report of inappropriate content
type ContentReport struct {
	gorm.Model
	ReporterID     uint           `json:"reporterId" gorm:"index:idx_reporter"`
	ContentType    string         `json:"contentType" gorm:"index:idx_content"` // "topic", "comment", etc.
	ContentID      uint           `json:"contentId" gorm:"index:idx_content"`
	Category       ReportCategory `json:"category" gorm:"index:idx_category"`
	Reason         string         `json:"reason" gorm:"type:text"`
	AdditionalInfo string         `json:"additionalInfo" gorm:"type:text"`
	Status         ReportStatus   `json:"status" gorm:"index:idx_status;default:'pending'"`
	AssignedTo     *uint          `json:"assignedTo"`
	ReviewedBy     *uint          `json:"reviewedBy"`
	ReviewedAt     *time.Time     `json:"reviewedAt"`
	Resolution     *ReportResolutionType `json:"resolution"`
	ResolutionNotes string        `json:"resolutionNotes" gorm:"type:text"`
	ReporterNotified bool         `json:"reporterNotified" gorm:"default:false"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

// ReportEvidence represents supporting evidence for a content report
type ReportEvidence struct {
	gorm.Model
	ReportID   uint      `json:"reportId" gorm:"index:idx_report"`
	Type       string    `json:"type"` // "screenshot", "link", "message", etc.
	Content    string    `json:"content" gorm:"type:text"`
	FilePath   string    `json:"filePath"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"createdAt"`
}

// ReportComment represents a comment on a content report
type ReportComment struct {
	gorm.Model
	ReportID  uint      `json:"reportId" gorm:"index:idx_report_comment"`
	UserID    uint      `json:"userId" gorm:"index:idx_comment_user"`
	Comment   string    `json:"comment" gorm:"type:text"`
	IsInternal bool     `json:"isInternal" gorm:"default:true"` // Is this comment only visible to staff
	CreatedAt time.Time `json:"createdAt"`
}

// ReportActionLog represents a log of actions taken on a report
type ReportActionLog struct {
	gorm.Model
	ReportID    uint      `json:"reportId" gorm:"index:idx_report_action"`
	UserID      uint      `json:"userId"`
	Action      string    `json:"action"` // "created", "assigned", "status_changed", "commented", "resolved", etc.
	OldValue    string    `json:"oldValue"`
	NewValue    string    `json:"newValue"`
	CreatedAt   time.Time `json:"createdAt"`
}