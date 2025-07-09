package models

import (
	"gorm.io/gorm"
)

// CategoryConfig represents configuration options for forum categories
type CategoryConfig struct {
	gorm.Model
	CategoryID                  uint   `json:"categoryId" gorm:"index"`
	RequireApproval             bool   `json:"requireApproval" gorm:"default:false"`
	AllowAnonymousPosts         bool   `json:"allowAnonymousPosts" gorm:"default:false"`
	AllowAttachments            bool   `json:"allowAttachments" gorm:"default:true"`
	MaxAttachmentSize           int    `json:"maxAttachmentSize" gorm:"default:5"` // In MB
	AttachmentTypes             string `json:"attachmentTypes" gorm:"default:'jpg,jpeg,png,gif,pdf,docx,xlsx,pptx'"`
	MinimumUserLevel            int    `json:"minimumUserLevel" gorm:"default:0"` // 0 = any registered user
	AutoCloseAfterDays          int    `json:"autoCloseAfterDays" gorm:"default:0"` // 0 = never auto-close
	DefaultTopicTags            string `json:"defaultTopicTags"` // Comma-separated list of default tag IDs
	RequireTopicTag             bool   `json:"requireTopicTag" gorm:"default:false"`
	EnableReactions             bool   `json:"enableReactions" gorm:"default:true"`
	EnableCommentReporting      bool   `json:"enableCommentReporting" gorm:"default:true"`
	RequireTopicApproval        bool   `json:"requireTopicApproval" gorm:"default:false"`
	RequireCommentApproval      bool   `json:"requireCommentApproval" gorm:"default:false"`
	ShowCategoryInNavigation    bool   `json:"showCategoryInNavigation" gorm:"default:true"`
	AutoPurgeTrashedItems       bool   `json:"autoPurgeTrashedItems" gorm:"default:false"`
	AutoPurgeTrashedAfterDays   int    `json:"autoPurgeTrashedAfterDays" gorm:"default:30"`
	IconClass                   string `json:"iconClass" gorm:"default:'fas fa-comments'"`
	CustomCSS                   string `json:"customCss" gorm:"type:text"`
	CategoryLayout              string `json:"categoryLayout" gorm:"default:'standard'"` // standard, compact, cards, etc.
	SidebarContent              string `json:"sidebarContent" gorm:"type:text"`
	HeaderContent               string `json:"headerContent" gorm:"type:text"`
	FooterContent               string `json:"footerContent" gorm:"type:text"`
	WelcomeMessage              string `json:"welcomeMessage" gorm:"type:text"`
	AutoModKeywords             string `json:"autoModKeywords" gorm:"type:text"` // Comma-separated words to flag
	Category                    Category `json:"-" gorm:"foreignKey:CategoryID"`
}

// PostingRules represents posting rules for a category
type PostingRules struct {
	gorm.Model
	CategoryID           uint   `json:"categoryId" gorm:"uniqueIndex"`
	MinimumTitleLength   int    `json:"minimumTitleLength" gorm:"default:5"`
	MaximumTitleLength   int    `json:"maximumTitleLength" gorm:"default:100"`
	MinimumContentLength int    `json:"minimumContentLength" gorm:"default:20"`
	MaximumContentLength int    `json:"maximumContentLength" gorm:"default:10000"`
	DisallowedWords      string `json:"disallowedWords" gorm:"type:text"` // Comma-separated
	RequiredWords        string `json:"requiredWords" gorm:"type:text"` // For specific categories
	PostingGuidelines    string `json:"postingGuidelines" gorm:"type:text"`
	BlockDuplicateURLs   bool   `json:"blockDuplicateUrls" gorm:"default:false"`
	RequireURLApproval   bool   `json:"requireUrlApproval" gorm:"default:false"`
	MaxLinksPerPost      int    `json:"maxLinksPerPost" gorm:"default:5"`
	PostingCooldown      int    `json:"postingCooldown" gorm:"default:0"` // In seconds
	Category             Category `json:"-" gorm:"foreignKey:CategoryID"`
}

// AutoModerationSettings represents automatic moderation settings for a category
type AutoModerationSettings struct {
	gorm.Model
	CategoryID             uint    `json:"categoryId" gorm:"uniqueIndex"`
	EnableAutoModeration   bool    `json:"enableAutoModeration" gorm:"default:true"`
	KeywordFlagging        bool    `json:"keywordFlagging" gorm:"default:true"`
	FlaggedKeywords        string  `json:"flaggedKeywords" gorm:"type:text"`
	SpamDetection          bool    `json:"spamDetection" gorm:"default:true"`
	SpamScoreThreshold     float64 `json:"spamScoreThreshold" gorm:"default:0.7"`
	SentimentAnalysis      bool    `json:"sentimentAnalysis" gorm:"default:false"`
	NegativeSentimentLimit float64 `json:"negativeSentimentLimit" gorm:"default:0.8"`
	EnableTextFilters      bool    `json:"enableTextFilters" gorm:"default:true"`
	TextFilterRules        string  `json:"textFilterRules" gorm:"type:text"` // JSON string of rules
	NotifyModerators       bool    `json:"notifyModerators" gorm:"default:true"`
	AutoApproveUsers       string  `json:"autoApproveUsers" gorm:"type:text"` // Comma-separated user roles
	Category               Category `json:"-" gorm:"foreignKey:CategoryID"`
}

// CategoryModerator represents moderators assigned to a category
type CategoryModerator struct {
	gorm.Model
	CategoryID  uint `json:"categoryId" gorm:"primaryKey"`
	UserID      uint `json:"userId" gorm:"primaryKey"`
	CanPinPosts bool `json:"canPinPosts" gorm:"default:true"`
	CanLockPosts bool `json:"canLockPosts" gorm:"default:true"`
	CanMovePosts bool `json:"canMovePosts" gorm:"default:true"`
	CanDeletePosts bool `json:"canDeletePosts" gorm:"default:true"`
	CanModerateUsers bool `json:"canModerateUsers" gorm:"default:false"`
	CanEditPosts bool `json:"canEditPosts" gorm:"default:true"`
	CanApproveContent bool `json:"canApproveContent" gorm:"default:true"`
	Category   Category `json:"-" gorm:"foreignKey:CategoryID"`
}