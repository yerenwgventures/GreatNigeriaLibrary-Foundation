package models

import (
        "time"

        "gorm.io/gorm"
)

// ContentFormat defines the format of content text
type ContentFormat string

const (
        // FormatMarkdown for markdown formatted content
        FormatMarkdown ContentFormat = "markdown"
        
        // FormatHTML for HTML formatted content
        FormatHTML ContentFormat = "html"
        
        // FormatPlain for plain text content
        FormatPlain ContentFormat = "plain"
        
        // FormatRichText for rich text editor content
        FormatRichText ContentFormat = "rich_text"
)

// RichTextContent represents rich text content in a topic or comment
type RichTextContent struct {
        gorm.Model
        ContentType    string        `json:"contentType"` // topic, comment
        ContentID      uint          `json:"contentId"`
        Format         ContentFormat `json:"format"`
        RawContent     string        `json:"rawContent" gorm:"type:text"`
        RenderedHTML   string        `json:"renderedHtml" gorm:"type:text"`
        HasMentions    bool          `json:"hasMentions" gorm:"default:false"`
        HasAttachments bool          `json:"hasAttachments" gorm:"default:false"`
        ContainsMedia  bool          `json:"containsMedia" gorm:"default:false"`
        ContainsCode   bool          `json:"containsCode" gorm:"default:false"`
        LastProcessedAt time.Time    `json:"lastProcessedAt"`
        CreatedAt      time.Time     `json:"createdAt"`
        UpdatedAt      time.Time     `json:"updatedAt"`
}

// Attachment represents a file attachment in rich text content
type Attachment struct {
        gorm.Model
        ContentType      string    `json:"contentType"` // topic, comment
        ContentID        uint      `json:"contentId"`
        FileName         string    `json:"fileName"`
        FileSize         int64     `json:"fileSize"`
        FileType         string    `json:"fileType"`
        StoragePath      string    `json:"storagePath"`
        ThumbnailPath    string    `json:"thumbnailPath"`
        UploadedBy       uint      `json:"uploadedBy"`
        UserID           uint      `json:"userId"`      // Alias for UploadedBy for compatibility
        IsImage          bool      `json:"isImage" gorm:"default:false"`
        ImageWidth       int       `json:"imageWidth"`
        ImageHeight      int       `json:"imageHeight"`
        Width            int       `json:"width"`       // Alias for ImageWidth for compatibility
        Height           int       `json:"height"`      // Alias for ImageHeight for compatibility
        URL              string    `json:"url"`         // Public URL for the attachment
        DownloadCount    int       `json:"downloadCount" gorm:"default:0"`
        IsSafe           bool      `json:"isSafe" gorm:"default:true"`
        ModeratedBy      *uint     `json:"moderatedBy"`
        ModeratedAt      *time.Time `json:"moderatedAt"`
        CreatedAt        time.Time `json:"createdAt"`   // Explicitly included for compatibility
}

// CodeBlock represents a code block in rich text content
type CodeBlock struct {
        gorm.Model
        ContentType    string    `json:"contentType"` // topic, comment
        ContentID      uint      `json:"contentId"`
        Code           string    `json:"code" gorm:"type:text"`
        Language       string    `json:"language"`
        HighlightedHTML string    `json:"highlightedHtml" gorm:"type:text"`
        LineCount      int       `json:"lineCount"`
}

// Quote represents a quote in rich text content
type Quote struct {
        gorm.Model
        ContentType    string    `json:"contentType"` // topic, comment
        ContentID      uint      `json:"contentId"`
        QuotedText     string    `json:"quotedText" gorm:"type:text"`
        QuotedTopicID  *uint     `json:"quotedTopicId"`
        QuotedCommentID *uint     `json:"quotedCommentId"`
        QuotedUserID   *uint     `json:"quotedUserId"`
        Citation       string    `json:"citation"`
}

// TextMention represents a mention of a user in text content
type TextMention struct {
        gorm.Model
        ContentType    string    `json:"contentType"` // topic, comment
        ContentID      uint      `json:"contentId"`
        MentionedUserID uint      `json:"mentionedUserId"`
        MentionedAt    time.Time `json:"mentionedAt"`
        IsRead         bool      `json:"isRead" gorm:"default:false"`
        ReadAt         *time.Time `json:"readAt"`
}

// Note: The Mention struct is defined in discussion.go