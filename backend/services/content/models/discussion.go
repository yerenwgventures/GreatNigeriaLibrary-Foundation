package models

import (
        "time"
)

// DiscussionTopic represents a discussion topic related to a content item
type DiscussionTopic struct {
        ID          uint         `json:"id" gorm:"primaryKey"`
        Title       string       `json:"title" gorm:"size:255;not null"`
        Description string       `json:"description" gorm:"type:text"`
        ContentID   uint         `json:"content_id" gorm:"index;not null"`
        ContentType ContentType  `json:"content_type" gorm:"type:varchar(50);not null;index"`
        UserID      uint         `json:"user_id" gorm:"index;not null"`
        Posts       []DiscussionPost `json:"posts,omitempty" gorm:"foreignKey:TopicID"`
        ViewCount   int          `json:"view_count" gorm:"default:0"`
        IsLocked    bool         `json:"is_locked" gorm:"default:false"`
        IsPinned    bool         `json:"is_pinned" gorm:"default:false"`
        CreatedAt   time.Time    `json:"created_at"`
        UpdatedAt   time.Time    `json:"updated_at"`
}

// DiscussionPost represents a post within a discussion topic
type DiscussionPost struct {
        ID        uint      `json:"id" gorm:"primaryKey"`
        Content   string    `json:"content" gorm:"type:text;not null"`
        TopicID   uint      `json:"topic_id" gorm:"index;not null"`
        UserID    uint      `json:"user_id" gorm:"index;not null"`
        ParentID  *uint     `json:"parent_id,omitempty" gorm:"index"`
        IsEdited  bool      `json:"is_edited" gorm:"default:false"`
        Likes     int       `json:"likes" gorm:"default:0"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the table name for DiscussionTopic
func (DiscussionTopic) TableName() string {
        return "discussion_topics"
}

// TableName overrides the table name for DiscussionPost
func (DiscussionPost) TableName() string {
        return "discussion_posts"
}