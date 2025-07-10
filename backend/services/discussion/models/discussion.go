package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a discussion category
type Category struct {
	gorm.Model
	Name        string `json:"name" gorm:"uniqueIndex"`
	Description string `json:"description"`
	Slug        string `json:"slug" gorm:"uniqueIndex"`
	ParentID    *uint  `json:"parentId"` // For subcategories
	Parent      *Category `json:"-" gorm:"foreignKey:ParentID"`
	Topics      []Topic   `json:"-" gorm:"foreignKey:CategoryID"`
	SortOrder   int       `json:"sortOrder" gorm:"default:0"`
	IsActive    bool      `json:"isActive" gorm:"default:true"`
}

// Topic represents a discussion topic
type Topic struct {
	gorm.Model
	Title       string    `json:"title"`
	Content     string    `json:"content" gorm:"type:text"`
	UserID      uint      `json:"userId" gorm:"index"`
	CategoryID  uint      `json:"categoryId" gorm:"index"`
	Category    Category  `json:"-" gorm:"foreignKey:CategoryID"`
	IsPinned    bool      `json:"isPinned" gorm:"default:false"`
	IsLocked    bool      `json:"isLocked" gorm:"default:false"`
	ViewCount   int       `json:"viewCount" gorm:"default:0"`
	LastPostAt  time.Time `json:"lastPostAt"`
	Comments    []Comment `json:"-" gorm:"foreignKey:TopicID"`
	Reactions   []Reaction `json:"-" gorm:"polymorphic:Target;polymorphicValue:topic"`
	BookID      *uint     `json:"bookId"` // Optional reference to a book
	ChapterID   *uint     `json:"chapterId"` // Optional reference to a chapter
	SectionID   *uint     `json:"sectionId"` // Optional reference to a section
	Tags        []Tag     `json:"tags" gorm:"many2many:topic_tags;"`
	RepliesCount int      `json:"repliesCount" gorm:"-"` // Calculated field, not stored
}

// Comment represents a comment on a topic or another comment
type Comment struct {
	gorm.Model
	Content     string    `json:"content" gorm:"type:text"`
	UserID      uint      `json:"userId" gorm:"index"`
	TopicID     uint      `json:"topicId" gorm:"index"`
	Topic       Topic     `json:"-" gorm:"foreignKey:TopicID"`
	ParentID    *uint     `json:"parentId"` // For replies to comments
	Parent      *Comment  `json:"-" gorm:"foreignKey:ParentID"`
	Replies     []Comment `json:"-" gorm:"foreignKey:ParentID"`
	IsApproved  bool      `json:"isApproved" gorm:"default:true"`
	IsFlagged   bool      `json:"isFlagged" gorm:"default:false"`
	FlagReason  string    `json:"flagReason"`
	Reactions   []Reaction `json:"-" gorm:"polymorphic:Target;polymorphicValue:comment"`
	EditedAt    *time.Time `json:"editedAt"`
	IsEdited    bool       `json:"isEdited" gorm:"default:false"`
	ReplyCount  int        `json:"replyCount" gorm:"-"` // Calculated field, not stored
}

// Reaction represents a user's reaction to a topic or comment
type Reaction struct {
	gorm.Model
	UserID     uint   `json:"userId" gorm:"index"`
	TargetType string `json:"targetType" gorm:"index"`
	TargetID   uint   `json:"targetId" gorm:"index"`
	ReactionType string `json:"reactionType"` // like, celebrate, insightful, etc.
}

// Tag represents a topic tag
type Tag struct {
	gorm.Model
	Name     string  `json:"name" gorm:"uniqueIndex"`
	Slug     string  `json:"slug" gorm:"uniqueIndex"`
	Color    string  `json:"color" gorm:"default:#3498db"`
	Topics   []Topic `json:"-" gorm:"many2many:topic_tags;"`
	IsSystem bool    `json:"isSystem" gorm:"default:false"`
}

// TopicTag represents the many-to-many relationship between topics and tags
type TopicTag struct {
	TopicID uint `gorm:"primaryKey"`
	TagID   uint `gorm:"primaryKey"`
}

// ReactionSummary represents aggregated reactions for a target
type ReactionSummary struct {
	TargetType   string `json:"targetType"`
	TargetID     uint   `json:"targetId"`
	ReactionType string `json:"reactionType"`
	Count        int    `json:"count"`
}

// UserDiscussionStats tracks a user's participation in discussions
type UserDiscussionStats struct {
	gorm.Model
	UserID          uint `json:"userId" gorm:"uniqueIndex"`
	TopicsCreated   int  `json:"topicsCreated" gorm:"default:0"`
	CommentsPosted  int  `json:"commentsPosted" gorm:"default:0"`
	ReactionsGiven  int  `json:"reactionsGiven" gorm:"default:0"`
	ReactionsReceived int `json:"reactionsReceived" gorm:"default:0"`
	LastActivityAt time.Time `json:"lastActivityAt"`
}

// ModerationQueue represents topics/comments waiting for moderation
type ModerationQueue struct {
	gorm.Model
	TargetType     string    `json:"targetType"` // topic or comment
	TargetID       uint      `json:"targetId"`
	ReportedBy     uint      `json:"reportedBy"`
	Reason         string    `json:"reason" gorm:"type:text"`
	Status         string    `json:"status"` // pending, approved, rejected
	ModeratedBy    *uint     `json:"moderatedBy"`
	ModeratedAt    *time.Time `json:"moderatedAt"`
	ModeratorNotes string     `json:"moderatorNotes" gorm:"type:text"`
}

// Subscription represents a user's subscription to a topic
type Subscription struct {
	gorm.Model
	UserID        uint   `json:"userId" gorm:"index"`
	TopicID       uint   `json:"topicId" gorm:"index"`
	NotifyOnReply bool   `json:"notifyOnReply" gorm:"default:true"`
	DigestType    string `json:"digestType"` // none, instant, daily, weekly
	IsActive      bool   `json:"isActive" gorm:"default:true"`
}

// NotificationPreference represents a user's notification preferences
type NotificationPreference struct {
	gorm.Model
	UserID         uint  `json:"userId" gorm:"uniqueIndex"`
	EmailOnReply   bool  `json:"emailOnReply" gorm:"default:true"`
	EmailOnMention bool  `json:"emailOnMention" gorm:"default:true"`
	EmailDigest    bool  `json:"emailDigest" gorm:"default:true"`
	DigestFrequency string `json:"digestFrequency" gorm:"default:weekly"`
}

// TopicView tracks when users view topics
type TopicView struct {
	gorm.Model
	UserID  uint      `json:"userId" gorm:"index"`
	TopicID uint      `json:"topicId" gorm:"index"`
	ViewedAt time.Time `json:"viewedAt"`
}

// Mention represents a mention of a user in a topic or comment
type Mention struct {
	gorm.Model
	UserID      uint   `json:"userId" gorm:"index"`
	TargetType  string `json:"targetType"` // topic or comment
	TargetID    uint   `json:"targetId"`
	MentionedBy uint   `json:"mentionedBy"`
	IsRead      bool   `json:"isRead" gorm:"default:false"`
}

// DiscussionError represents custom errors for the discussion module
type DiscussionError struct {
	Code    string
	Message string
}

func (e DiscussionError) Error() string {
	return e.Message
}

// Common error types
var (
	ErrTopicNotFound     = DiscussionError{Code: "topic_not_found", Message: "Topic not found"}
	ErrCommentNotFound   = DiscussionError{Code: "comment_not_found", Message: "Comment not found"}
	ErrCategoryNotFound  = DiscussionError{Code: "category_not_found", Message: "Category not found"}
	ErrPermissionDenied  = DiscussionError{Code: "permission_denied", Message: "Permission denied"}
	ErrTopicLocked       = DiscussionError{Code: "topic_locked", Message: "Topic is locked"}
	ErrInvalidContent    = DiscussionError{Code: "invalid_content", Message: "Invalid content"}
	ErrDuplicateReaction = DiscussionError{Code: "duplicate_reaction", Message: "Duplicate reaction"}
)