package service

import (
        "encoding/json"
        "fmt"
        "log"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/repository"
)

// CitationService handles citation-related operations
type CitationService struct {
        citationRepo *repository.CitationRepository
        bookRepo     repository.BookRepository
}

// NewCitationService creates a new citation service
func NewCitationService(
        citationRepo *repository.CitationRepository,
        bookRepo repository.BookRepository,
) *CitationService {
        return &CitationService{
                citationRepo: citationRepo,
                bookRepo:     bookRepo,
        }
}

// CreateCitation adds a new citation
func (s *CitationService) CreateCitation(citation *models.Citation) error {
        return s.citationRepo.CreateCitation(citation)
}

// GetCitationByID retrieves a citation by ID
func (s *CitationService) GetCitationByID(id uint) (*models.Citation, error) {
        return s.citationRepo.GetCitationByID(id)
}

// GetCitationsByBook retrieves all citations for a book
func (s *CitationService) GetCitationsByBook(bookID uint) ([]models.Citation, error) {
        return s.citationRepo.GetCitationsByBook(bookID)
}

// RecordCitationUsage records a citation being used in the book
func (s *CitationService) RecordCitationUsage(citationID, bookID, chapterID, sectionID uint) error {
        usage := &models.CitationUsage{
                CitationID: citationID,
                BookID:     bookID,
                ChapterID:  chapterID,
                SectionID:  sectionID,
        }
        return s.citationRepo.RecordCitationUsage(usage)
}

// GenerateBibliography creates a bibliography for a book
func (s *CitationService) GenerateBibliography(bookID uint) (string, error) {
        // Generate the bibliography content
        bibliography, err := s.citationRepo.GenerateBibliography(bookID)
        if err != nil {
                return "", err
        }

        // Update the book's back matter with the generated bibliography
        backmatter, err := s.bookRepo.GetBackMatterByBookID(bookID, true)
        if err != nil {
                // If back matter doesn't exist, create it
                newBackMatter := &models.BookBackMatter{
                        BookID:       bookID,
                        Bibliography: bibliography,
                }
                if err := s.bookRepo.UpdateBackMatter(newBackMatter); err != nil {
                        return "", err
                }
                return bibliography, nil
        }

        // Update existing back matter
        backmatter.Bibliography = bibliography
        if err := s.bookRepo.UpdateBackMatter(backmatter); err != nil {
                return "", err
        }

        return bibliography, nil
}

// GetCitationStats retrieves citation statistics for a book
func (s *CitationService) GetCitationStats(bookID uint) (*models.CitationStats, error) {
        return s.citationRepo.GetCitationStats(bookID)
}

// ImportCitationsFromJSON imports citations from a JSON file
func (s *CitationService) ImportCitationsFromJSON(bookID uint, jsonData string) error {
        var citations []models.Citation
        if err := json.Unmarshal([]byte(jsonData), &citations); err != nil {
                return fmt.Errorf("failed to parse citation JSON: %w", err)
        }

        // Process each citation
        for _, c := range citations {
                c.BookID = bookID // Ensure the book ID is set correctly
                if err := s.citationRepo.CreateCitation(&c); err != nil {
                        log.Printf("Error importing citation %s: %v", c.CitationKey, err)
                        // Continue with other citations
                }
        }

        return nil
}

// UpdateAppendixWithBibliography updates the appendix with the bibliography content
func (s *CitationService) UpdateAppendixWithBibliography(bookID uint) error {
        // Generate the bibliography
        bibliography, err := s.citationRepo.GenerateBibliography(bookID)
        if err != nil {
                return err
        }

        // Get the book's back matter
        backmatter, err := s.bookRepo.GetBackMatterByBookID(bookID, true)
        if err != nil {
                return err
        }

        // Parse the existing appendix JSON
        var appendixItems []models.AppendixItem
        if backmatter.AppendixJSON != "" {
                if err := json.Unmarshal([]byte(backmatter.AppendixJSON), &appendixItems); err != nil {
                        return fmt.Errorf("failed to parse appendix JSON: %w", err)
                }
        }

        // Check if we already have a Bibliography appendix
        bibliographyExists := false
        for i, item := range appendixItems {
                if item.Title == "Bibliography and Citations" {
                        // Update the existing bibliography appendix
                        appendixItems[i].Content = bibliography
                        bibliographyExists = true
                        break
                }
        }

        // If no bibliography appendix exists, add one
        if !bibliographyExists {
                appendixItems = append(appendixItems, models.AppendixItem{
                        Title:   "Bibliography and Citations",
                        Content: bibliography,
                })
        }

        // Marshal the updated appendix items back to JSON
        appendixJSON, err := json.Marshal(appendixItems)
        if err != nil {
                return fmt.Errorf("failed to marshal appendix JSON: %w", err)
        }

        // Update the back matter
        backmatter.AppendixJSON = string(appendixJSON)
        return s.bookRepo.UpdateBackMatter(backmatter)
}

// AddCitationFromText parses citation text and adds to the database
func (s *CitationService) AddCitationFromText(bookID uint, citationType, citationText string) (*models.Citation, error) {
        // Simple implementation for demo - in real app would use more sophisticated parsing
        citation := &models.Citation{
                BookID:      bookID,
                CitationKey: fmt.Sprintf("%s-%d", citationType, bookID), // Example key, would be better in real app
                RefNumber:   0, // Will be set by repo
                Type:        citationType,
                Author:      "Parsed Author", // Would be parsed from citationText
                Year:        "2023",          // Would be parsed from citationText
                Title:       citationText,    // Would be parsed from citationText
                Source:      "Parsed Source", // Would be parsed from citationText
        }

        if err := s.citationRepo.CreateCitation(citation); err != nil {
                return nil, err
        }

        return citation, nil
}