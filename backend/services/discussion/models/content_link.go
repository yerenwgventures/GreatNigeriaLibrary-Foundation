package models

import (
        "time"

        "gorm.io/gorm"
)

// ContentReferenceType defines the type of content being referenced
type ContentReferenceType string

const (
        // BookReference for references to books
        BookReference ContentReferenceType = "book"
        
        // ChapterReference for references to chapters
        ChapterReference ContentReferenceType = "chapter"
        
        // SectionReference for references to sections
        SectionReference ContentReferenceType = "section"
        
        // ResourceReference for resource references
        ResourceReference ContentReferenceType = "resource"
        
        // ProjectReference for project references
        ProjectReference ContentReferenceType = "project"
        
        // ReportReference for report references
        ReportReference ContentReferenceType = "report"
)

// TopicContentLink represents a link between a discussion topic and book content
type TopicContentLink struct {
        gorm.Model
        TopicID         uint                `json:"topicId" gorm:"index:idx_topic_content_link"`
        ContentType     ContentReferenceType `json:"contentType" gorm:"index:idx_topic_content_link"`
        ContentID       uint                `json:"contentId" gorm:"index:idx_topic_content_link"`
        CreatedBy       uint                `json:"createdBy"`
        IsAutoGenerated bool                `json:"isAutoGenerated" gorm:"default:false"`
        IsHighlighted   bool                `json:"isHighlighted" gorm:"default:false"`
        CreatedAt       time.Time           `json:"createdAt"`
        UpdatedAt       time.Time           `json:"updatedAt"`
}

// CommentContentLink represents a link between a comment and book content
type CommentContentLink struct {
        gorm.Model
        CommentID       uint                `json:"commentId" gorm:"index:idx_comment_content_link"`
        ContentType     ContentReferenceType `json:"contentType" gorm:"index:idx_comment_content_link"`
        ContentID       uint                `json:"contentId" gorm:"index:idx_comment_content_link"`
        CreatedBy       uint                `json:"createdBy"`
        CitationText    string              `json:"citationText" gorm:"type:text"`
        CitationContext string              `json:"citationContext" gorm:"type:text"`
        CreatedAt       time.Time           `json:"createdAt"`
        UpdatedAt       time.Time           `json:"updatedAt"`
}

// ContentDiscussionRecommendation represents a recommendation for related discussions
type ContentDiscussionRecommendation struct {
        gorm.Model
        ContentType     ContentReferenceType `json:"contentType" gorm:"index:idx_content_recommendation"`
        ContentID       uint                `json:"contentId" gorm:"index:idx_content_recommendation"`
        TopicID         uint                `json:"topicId"`
        RecommendationScore float64         `json:"recommendationScore"`
        IsManuallyAdded bool                `json:"isManuallyAdded" gorm:"default:false"`
        AddedBy         *uint               `json:"addedBy"`
        CreatedAt       time.Time           `json:"createdAt"`
        UpdatedAt       time.Time           `json:"updatedAt"`
}

// AutoGeneratedTopicTemplate represents a template for auto-generating topics
type AutoGeneratedTopicTemplate struct {
        gorm.Model
        Name           string    `json:"name"`
        ContentType    ContentReferenceType `json:"contentType"`
        TitleTemplate  string    `json:"titleTemplate" gorm:"type:text"`
        BodyTemplate   string    `json:"bodyTemplate" gorm:"type:text"`
        IsActive       bool      `json:"isActive" gorm:"default:true"`
        CreatedBy      uint      `json:"createdBy"`
        CreatedAt      time.Time `json:"createdAt"`
        UpdatedAt      time.Time `json:"updatedAt"`
}