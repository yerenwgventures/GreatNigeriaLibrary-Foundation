package models

import (
        "time"
)

// Book represents a book in the system
type Book struct {
        ID                uint          `json:"id" gorm:"primaryKey"`
        Title             string        `json:"title" gorm:"not null"`
        Subtitle          string        `json:"subtitle"`
        Author            string        `json:"author"`
        Description       string        `json:"description" gorm:"type:text"`
        CoverImage        string        `json:"cover_image"`
        Published         bool          `json:"published" gorm:"default:false"`
        ScheduledPublishAt *time.Time    `json:"scheduled_publish_at"`
        CreatedAt         time.Time     `json:"created_at"`
        UpdatedAt         time.Time     `json:"updated_at"`
        Chapters          []BookChapter `json:"chapters" gorm:"-"` // Not persisted to the database
}

// BookFrontMatter represents the front matter of a book
type BookFrontMatter struct {
        ID               uint      `json:"id" gorm:"primaryKey"`
        BookID           uint      `json:"book_id" gorm:"not null"`
        Introduction     string    `json:"introduction" gorm:"type:text"`
        Preface          string    `json:"preface" gorm:"type:text"`
        Acknowledgements string    `json:"acknowledgements" gorm:"type:text"`
        SupportAuthor    string    `json:"support_author" gorm:"type:text;column:support_author"`
        CreatedAt        time.Time `json:"created_at"`
        UpdatedAt        time.Time `json:"updated_at"`
}

// EpilogueItem represents a single epilogue item for Book 3
type EpilogueItem struct {
        Title   string `json:"title"`
        Content string `json:"content" gorm:"type:text"`
        Quote   string `json:"quote" gorm:"type:text"`
}

// AppendixItem represents a single appendix item for Book 3
type AppendixItem struct {
        Title   string `json:"title"`
        Content string `json:"content" gorm:"type:text"`
}

// BookBackMatter represents the back matter content of a book (conclusion, appendices, etc.)
type BookBackMatter struct {
        ID           uint           `json:"id" gorm:"primaryKey"`
        BookID       uint           `json:"book_id" gorm:"not null"`
        Conclusion   string         `json:"conclusion" gorm:"type:text"`
        Appendices   string         `json:"appendices" gorm:"type:text"`
        Bibliography string         `json:"bibliography" gorm:"type:text"`
        Glossary     string         `json:"glossary" gorm:"type:text"`
        AboutAuthor  string         `json:"about_author" gorm:"type:text;column:about_author"`
        Epilogue     []EpilogueItem `json:"epilogue" gorm:"-"` // Not persisted directly to the database
        Appendix     []AppendixItem `json:"appendix" gorm:"-"` // Not persisted directly to the database
        EpilogueJSON string         `json:"-" gorm:"type:text;column:epilogue_json"`
        AppendixJSON string         `json:"-" gorm:"type:text;column:appendix_json"`
        CreatedAt    time.Time      `json:"created_at"`
        UpdatedAt    time.Time      `json:"updated_at"`
}

// BookChapter represents a chapter in a book
type BookChapter struct {
        ID                uint          `json:"id" gorm:"primaryKey"`
        BookID            uint          `json:"book_id" gorm:"not null"`
        Title             string        `json:"title" gorm:"not null"`
        Number            int           `json:"number" gorm:"not null"`
        Description       string        `json:"description" gorm:"type:text"`
        Content           string        `json:"content" gorm:"type:text"`
        Published         bool          `json:"published" gorm:"default:false"`
        ScheduledPublishAt *time.Time    `json:"scheduled_publish_at"`
        CreatedAt         time.Time     `json:"created_at"`
        UpdatedAt         time.Time     `json:"updated_at"`
        Sections          []BookSection `json:"sections" gorm:"-"` // Not persisted to the database
}

// BookSection represents a section within a chapter
type BookSection struct {
        ID                uint             `json:"id" gorm:"primaryKey"`
        BookID            uint             `json:"book_id" gorm:"not null"`
        ChapterID         uint             `json:"chapter_id" gorm:"not null"`
        Title             string           `json:"title" gorm:"not null"`
        Number            int              `json:"number" gorm:"not null"`
        Content           string           `json:"content" gorm:"type:text"`
        Format            string           `json:"format" gorm:"default:'markdown'"`
        TimeToRead        int              `json:"time_to_read" gorm:"default:5"`
        Published         bool             `json:"published" gorm:"default:false"`
        ScheduledPublishAt *time.Time       `json:"scheduled_publish_at"`
        CreatedAt         time.Time        `json:"created_at"`
        UpdatedAt         time.Time        `json:"updated_at"`
        Subsections       []BookSubsection `json:"subsections" gorm:"-"` // Not persisted to the database
}

// BookProgress represents a user's progress through a book
type BookProgress struct {
        ID          uint      `json:"id" gorm:"primaryKey"`
        UserID      uint      `json:"user_id" gorm:"not null"`
        BookID      uint      `json:"book_id" gorm:"not null"`
        ChapterID   uint      `json:"chapter_id"`
        SectionID   uint      `json:"section_id"`
        Progress    float64   `json:"progress" gorm:"type:decimal(5,2)"`
        IsRead      bool      `json:"is_read" gorm:"default:false"`
        LastAccess  time.Time `json:"last_access"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

// BookBookmark represents a bookmark in a book
type BookBookmark struct {
        ID          uint      `json:"id" gorm:"primaryKey"`
        UserID      uint      `json:"user_id" gorm:"not null"`
        BookID      uint      `json:"book_id" gorm:"not null"`
        ChapterID   uint      `json:"chapter_id"`
        SectionID   uint      `json:"section_id"`
        Title       string    `json:"title"`
        Description string    `json:"description" gorm:"type:text"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

// BookNote represents a note attached to a specific part of a book
type BookNote struct {
        ID          uint      `json:"id" gorm:"primaryKey"`
        UserID      uint      `json:"user_id" gorm:"not null"`
        BookID      uint      `json:"book_id" gorm:"not null"`
        ChapterID   uint      `json:"chapter_id"`
        SectionID   uint      `json:"section_id"`
        Content     string    `json:"content" gorm:"type:text"`
        Color       string    `json:"color" gorm:"default:'yellow'"`
        Position    int       `json:"position"`
        Tags        string    `json:"tags"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

// BookFeedback represents user feedback on book content
type BookFeedback struct {
        ID          uint      `json:"id" gorm:"primaryKey"`
        UserID      uint      `json:"user_id" gorm:"not null"`
        BookID      uint      `json:"book_id" gorm:"not null"`
        ChapterID   uint      `json:"chapter_id"`
        SectionID   uint      `json:"section_id"`
        Rating      int       `json:"rating"`
        Comment     string    `json:"comment" gorm:"type:text"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

// ForumTopic represents a discussion topic linked to book content
type ForumTopic struct {
        ID          uint      `json:"id" gorm:"primaryKey"`
        Title       string    `json:"title" gorm:"not null"`
        Description string    `json:"description" gorm:"type:text"`
        BookID      uint      `json:"book_id"`
        ChapterID   uint      `json:"chapter_id"`
        SectionID   uint      `json:"section_id"`
        UserID      uint      `json:"user_id" gorm:"not null"`
        Status      string    `json:"status" gorm:"default:'open'"`
        Views       int       `json:"views" gorm:"default:0"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

// ActionStep represents a concrete action step derived from book content
type ActionStep struct {
        ID          uint      `json:"id" gorm:"primaryKey"`
        Title       string    `json:"title" gorm:"not null"`
        Description string    `json:"description" gorm:"type:text"`
        BookID      uint      `json:"book_id"`
        ChapterID   uint      `json:"chapter_id"`
        SectionID   uint      `json:"section_id"`
        UserID      uint      `json:"user_id" gorm:"not null"`
        Status      string    `json:"status" gorm:"default:'pending'"`
        DueDate     time.Time `json:"due_date"`
        CompletedAt time.Time `json:"completed_at"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

// Bookmark represents a user's bookmark of content
type Bookmark struct {
        ID          uint      `json:"id" gorm:"primaryKey"`
        UserID      uint      `json:"user_id" gorm:"not null"`
        BookID      uint      `json:"book_id" gorm:"not null"`
        ChapterID   uint      `json:"chapter_id"`
        SectionID   uint      `json:"section_id"`
        Title       string    `json:"title"`
        Description string    `json:"description" gorm:"type:text"`
        Color       string    `json:"color" gorm:"default:'blue'"`
        Position    int       `json:"position"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

// ContentRevision types
type BookRevision struct {
        ID        uint      `json:"id" gorm:"primaryKey"`
        BookID    uint      `json:"book_id" gorm:"not null"`
        Content   string    `json:"content" gorm:"type:text"`   // The actual content
        Changes   string    `json:"changes" gorm:"type:text"`
        Notes     string    `json:"notes" gorm:"type:text"`
        CreatedBy uint      `json:"created_by" gorm:"not null"`
        CreatedAt time.Time `json:"created_at"`
}

type ChapterRevision struct {
        ID        uint      `json:"id" gorm:"primaryKey"`
        ChapterID uint      `json:"chapter_id" gorm:"not null"`
        Content   string    `json:"content" gorm:"type:text"`   // The actual content
        Changes   string    `json:"changes" gorm:"type:text"`
        Notes     string    `json:"notes" gorm:"type:text"`
        CreatedBy uint      `json:"created_by" gorm:"not null"`
        CreatedAt time.Time `json:"created_at"`
}

type SectionRevision struct {
        ID        uint      `json:"id" gorm:"primaryKey"`
        SectionID uint      `json:"section_id" gorm:"not null"`
        Content   string    `json:"content" gorm:"type:text"`   // The actual content
        Changes   string    `json:"changes" gorm:"type:text"`
        Notes     string    `json:"notes" gorm:"type:text"`
        CreatedBy uint      `json:"created_by" gorm:"not null"`
        CreatedAt time.Time `json:"created_at"`
}

// BookSubsection represents a subsection within a section
type BookSubsection struct {
        ID        uint      `json:"id" gorm:"primaryKey"`
        BookID    uint      `json:"book_id" gorm:"not null"`
        ChapterID uint      `json:"chapter_id" gorm:"not null"`
        SectionID uint      `json:"section_id" gorm:"not null"`
        Title     string    `json:"title" gorm:"not null"`
        Number    string    `json:"number" gorm:"not null"` // Using string for numbers like "1.1.1"
        Content   string    `json:"content" gorm:"type:text"`
        MediaURL  string    `json:"media_url"`
        Quote     string    `json:"quote" gorm:"type:text"`
        Poem      string    `json:"poem" gorm:"type:text"`
        Format    string    `json:"format" gorm:"default:'markdown'"`
        Published bool      `json:"published" gorm:"default:false"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
}

// Define content types for organization
type ContentType string

const (
        ContentTypeBook       ContentType = "book"
        ContentTypeChapter    ContentType = "chapter"
        ContentTypeSection    ContentType = "section"
        ContentTypeSubsection ContentType = "subsection"
)

// Define book content status
type ContentStatus string

const (
        ContentStatusDraft     ContentStatus = "draft"
        ContentStatusPublished ContentStatus = "published"
        ContentStatusArchived  ContentStatus = "archived"
)

// Chapter represents the content structure in the system
type Chapter struct {
        ID          uint          `json:"id" gorm:"primaryKey"`
        BookID      uint          `json:"book_id" gorm:"not null"`
        Title       string        `json:"title" gorm:"not null"`
        Number      int           `json:"number" gorm:"not null"`
        Description string        `json:"description" gorm:"type:text"`
        Status      ContentStatus `json:"status" gorm:"type:varchar(20);default:'draft'"`
        CreatedAt   time.Time     `json:"created_at"`
        UpdatedAt   time.Time     `json:"updated_at"`
}

// Section represents the content structure in the system
type Section struct {
        ID          uint          `json:"id" gorm:"primaryKey"`
        ChapterID   uint          `json:"chapter_id" gorm:"not null"`
        BookID      uint          `json:"book_id" gorm:"not null"`
        Title       string        `json:"title" gorm:"not null"`
        Content     string        `json:"content" gorm:"type:text"`
        Number      int           `json:"number" gorm:"not null"`
        Status      ContentStatus `json:"status" gorm:"type:varchar(20);default:'draft'"`
        CreatedAt   time.Time     `json:"created_at"`
        UpdatedAt   time.Time     `json:"updated_at"`
}

// TableName overrides the table name for Section
func (Section) TableName() string {
        return "book_sections"
}