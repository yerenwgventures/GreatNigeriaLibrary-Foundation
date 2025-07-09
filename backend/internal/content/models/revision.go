package models

import (
        "time"

        "gorm.io/gorm"
)

// Using BookRevision, ChapterRevision, and SectionRevision from book.go

// ContentRevisionLog represents a log of content revisions
type ContentRevisionLog struct {
        gorm.Model
        ContentType    string    `json:"contentType"`    // book, chapter, or section
        ContentID      uint      `json:"contentId"`
        RevisionID     uint      `json:"revisionId"`
        UserID         uint      `json:"userId"`
        Action         string    `json:"action"`         // create, update, delete, restore
        ChangeSummary  string    `json:"changeSummary"`  // Summary of changes
        Timestamp      time.Time `json:"timestamp"`
}