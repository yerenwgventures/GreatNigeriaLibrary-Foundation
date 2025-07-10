package models

import (
	"time"
)

// Citation represents a single bibliographic entry
type Citation struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BookID      uint      `json:"book_id" gorm:"not null"`
	CitationKey string    `json:"citation_key" gorm:"not null;uniqueIndex:idx_book_citation_key,priority:1"` // Unique identifier like "achebe1983"
	RefNumber   int       `json:"ref_number" gorm:"not null"`                                                // Number used in the in-text citation
	Author      string    `json:"author" gorm:"not null"`                                                    // Author(s) name(s)
	Year        string    `json:"year" gorm:"not null"`                                                      // Publication year
	Title       string    `json:"title" gorm:"not null"`                                                     // Title of work
	Source      string    `json:"source" gorm:"not null"`                                                    // Source information (publisher, journal, etc.)
	URL         string    `json:"url"`                                                                       // Online link if available
	Type        string    `json:"type" gorm:"not null"`                                                      // Type of reference (book, article, interview, etc.)
	CitedCount  int       `json:"cited_count" gorm:"default:0"`                                              // How many times this source is cited across the book
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CitationUsage tracks where each citation is used within the book
type CitationUsage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CitationID uint      `json:"citation_id" gorm:"not null"`
	BookID     uint      `json:"book_id" gorm:"not null"`
	ChapterID  uint      `json:"chapter_id" gorm:"not null"`
	SectionID  uint      `json:"section_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}

// BibliographyMetadata stores information about bibliographies
type BibliographyMetadata struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	BookID        uint      `json:"book_id" gorm:"not null;uniqueIndex"`
	Title         string    `json:"title" gorm:"not null"`
	Description   string    `json:"description" gorm:"type:text"`
	LastGenerated time.Time `json:"last_generated"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CitationCategory defines the types of citations
type CitationCategory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	DisplayOrder int       `json:"display_order" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CitationStats represents statistics about citations in a book
type CitationStats struct {
	BookID       uint   `json:"book_id"`
	TotalCount   int    `json:"total_count"`
	ByType       map[string]int `json:"by_type"`
	MostCited    []Citation `json:"most_cited"`
	AcademicRate float64 `json:"academic_rate"` // Percentage of academic sources
}