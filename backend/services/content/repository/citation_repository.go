package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
	"gorm.io/gorm"
)

// CitationRepository handles database operations for citations
type CitationRepository struct {
	DB *gorm.DB
}

// NewCitationRepository creates a new citation repository
func NewCitationRepository(db *gorm.DB) *CitationRepository {
	return &CitationRepository{
		DB: db,
	}
}

// CreateCitation adds a new citation to the database
func (r *CitationRepository) CreateCitation(citation *models.Citation) error {
	return r.DB.Create(citation).Error
}

// GetCitationByID retrieves a citation by its ID
func (r *CitationRepository) GetCitationByID(id uint) (*models.Citation, error) {
	var citation models.Citation
	if err := r.DB.First(&citation, id).Error; err != nil {
		return nil, err
	}
	return &citation, nil
}

// GetCitationByKey retrieves a citation by its citation key within a book
func (r *CitationRepository) GetCitationByKey(bookID uint, citationKey string) (*models.Citation, error) {
	var citation models.Citation
	if err := r.DB.Where("book_id = ? AND citation_key = ?", bookID, citationKey).First(&citation).Error; err != nil {
		return nil, err
	}
	return &citation, nil
}

// GetCitationsByBook retrieves all citations for a specific book
func (r *CitationRepository) GetCitationsByBook(bookID uint) ([]models.Citation, error) {
	var citations []models.Citation
	if err := r.DB.Where("book_id = ?", bookID).Order("ref_number").Find(&citations).Error; err != nil {
		return nil, err
	}
	return citations, nil
}

// GetCitationsByType retrieves citations of a specific type for a book
func (r *CitationRepository) GetCitationsByType(bookID uint, citationType string) ([]models.Citation, error) {
	var citations []models.Citation
	if err := r.DB.Where("book_id = ? AND type = ?", bookID, citationType).
		Order("author, year").Find(&citations).Error; err != nil {
		return nil, err
	}
	return citations, nil
}

// GetCitationsUsage retrieves usage information for a citation
func (r *CitationRepository) GetCitationsUsage(citationID uint) ([]models.CitationUsage, error) {
	var usages []models.CitationUsage
	if err := r.DB.Where("citation_id = ?", citationID).Find(&usages).Error; err != nil {
		return nil, err
	}
	return usages, nil
}

// RecordCitationUsage records the usage of a citation in a specific location
func (r *CitationRepository) RecordCitationUsage(usage *models.CitationUsage) error {
	// Check if this usage already exists to avoid duplicates
	var count int64
	r.DB.Model(&models.CitationUsage{}).
		Where("citation_id = ? AND book_id = ? AND chapter_id = ? AND section_id = ?",
			usage.CitationID, usage.BookID, usage.ChapterID, usage.SectionID).
		Count(&count)

	if count > 0 {
		return nil // Usage already recorded
	}

	// Create the usage record
	if err := r.DB.Create(usage).Error; err != nil {
		return err
	}

	// Increment the cited_count for the citation
	return r.DB.Model(&models.Citation{}).
		Where("id = ?", usage.CitationID).
		UpdateColumn("cited_count", gorm.Expr("cited_count + ?", 1)).
		Error
}

// GetCitationStats retrieves citation statistics for a book
func (r *CitationRepository) GetCitationStats(bookID uint) (*models.CitationStats, error) {
	var totalCount int64
	r.DB.Model(&models.Citation{}).Where("book_id = ?", bookID).Count(&totalCount)

	// Get counts by type
	type TypeCount struct {
		Type  string
		Count int
	}
	var typeCounts []TypeCount
	r.DB.Model(&models.Citation{}).
		Select("type, COUNT(*) as count").
		Where("book_id = ?", bookID).
		Group("type").
		Scan(&typeCounts)

	byType := make(map[string]int)
	for _, tc := range typeCounts {
		byType[tc.Type] = tc.Count
	}

	// Get most cited
	var mostCited []models.Citation
	r.DB.Where("book_id = ?", bookID).
		Order("cited_count DESC").
		Limit(10).
		Find(&mostCited)

	// Calculate academic rate (books, journals)
	var academicCount int64
	r.DB.Model(&models.Citation{}).
		Where("book_id = ? AND (type = 'book' OR type = 'journal')", bookID).
		Count(&academicCount)

	academicRate := 0.0
	if totalCount > 0 {
		academicRate = float64(academicCount) / float64(totalCount) * 100
	}

	return &models.CitationStats{
		BookID:       bookID,
		TotalCount:   int(totalCount),
		ByType:       byType,
		MostCited:    mostCited,
		AcademicRate: academicRate,
	}, nil
}

// GenerateBibliography creates a bibliography for a book
func (r *CitationRepository) GenerateBibliography(bookID uint) (string, error) {
	var book models.Book
	if err := r.DB.First(&book, bookID).Error; err != nil {
		return "", fmt.Errorf("book not found: %w", err)
	}

	// Get all citations for this book
	citations, err := r.GetCitationsByBook(bookID)
	if err != nil {
		return "", err
	}

	// Generate the bibliography content in proper academic format
	var bibliography strings.Builder
	bibliography.WriteString(fmt.Sprintf("# Bibliography for %s\n\n", book.Title))

	// Group citations by type
	types := []string{"book", "journal", "report", "government", "interview", "survey", "media"}
	sectionHeadings := map[string]string{
		"book":       "### Books",
		"journal":    "### Journal Articles",
		"report":     "### Reports and Institutional Publications",
		"government": "### Government Sources",
		"interview":  "### Interviews and Focus Groups",
		"survey":     "### Survey Data",
		"media":      "### Media Sources",
	}

	// Academic sources section
	bibliography.WriteString("## Academic Sources\n\n")

	// Process book and journal citations under Academic Sources
	for _, t := range []string{"book", "journal"} {
		var typeCitations []models.Citation
		for _, c := range citations {
			if c.Type == t {
				typeCitations = append(typeCitations, c)
			}
		}

		if len(typeCitations) > 0 {
			bibliography.WriteString(sectionHeadings[t] + "\n\n")
			for _, c := range typeCitations {
				if t == "book" {
					bibliography.WriteString(fmt.Sprintf("%s. (%s). *%s*. %s. [%d]\n\n",
						c.Author, c.Year, c.Title, c.Source, c.RefNumber))
				} else {
					bibliography.WriteString(fmt.Sprintf("%s. (%s). '%s'. *%s*. [%d]\n\n",
						c.Author, c.Year, c.Title, c.Source, c.RefNumber))
				}
			}
		}
	}

	// Reports and institutional publications
	bibliography.WriteString("## Reports and Institutional Publications\n\n")
	for _, t := range []string{"report"} {
		var typeCitations []models.Citation
		for _, c := range citations {
			if c.Type == t {
				typeCitations = append(typeCitations, c)
			}
		}

		if len(typeCitations) > 0 {
			for _, c := range typeCitations {
				bibliography.WriteString(fmt.Sprintf("%s. (%s). *%s*. %s. [%d]\n\n",
					c.Author, c.Year, c.Title, c.Source, c.RefNumber))
			}
		}
	}

	// Government sources
	bibliography.WriteString("## Government Sources\n\n")
	for _, t := range []string{"government"} {
		var typeCitations []models.Citation
		for _, c := range citations {
			if c.Type == t {
				typeCitations = append(typeCitations, c)
			}
		}

		if len(typeCitations) > 0 {
			for _, c := range typeCitations {
				bibliography.WriteString(fmt.Sprintf("%s. (%s). *%s*. %s. [%d]\n\n",
					c.Author, c.Year, c.Title, c.Source, c.RefNumber))
			}
		}
	}

	// Field research
	bibliography.WriteString("## Field Research\n\n")

	// Interviews and focus groups
	bibliography.WriteString("### Interviews and Focus Groups\n\n")
	bibliography.WriteString("*Note: Names have been changed to protect privacy*\n\n")
	for _, t := range []string{"interview"} {
		var typeCitations []models.Citation
		for _, c := range citations {
			if c.Type == t {
				typeCitations = append(typeCitations, c)
			}
		}

		if len(typeCitations) > 0 {
			for _, c := range typeCitations {
				bibliography.WriteString(fmt.Sprintf("%s. (%s). %s. %s. [%d]\n\n",
					c.Author, c.Year, c.Title, c.Source, c.RefNumber))
			}
		}
	}

	// Survey data
	bibliography.WriteString("### Survey Data\n\n")
	for _, t := range []string{"survey"} {
		var typeCitations []models.Citation
		for _, c := range citations {
			if c.Type == t {
				typeCitations = append(typeCitations, c)
			}
		}

		if len(typeCitations) > 0 {
			for _, c := range typeCitations {
				bibliography.WriteString(fmt.Sprintf("%s. (%s). %s. %s. [%d]\n\n",
					c.Author, c.Year, c.Title, c.Source, c.RefNumber))
			}
		}
	}

	// Media sources
	bibliography.WriteString("## Media Sources\n\n")
	for _, t := range []string{"media"} {
		var typeCitations []models.Citation
		for _, c := range citations {
			if c.Type == t {
				typeCitations = append(typeCitations, c)
			}
		}

		if len(typeCitations) > 0 {
			for _, c := range typeCitations {
				bibliography.WriteString(fmt.Sprintf("%s. (%s). '%s'. *%s*. [%d]\n\n",
					c.Author, c.Year, c.Title, c.Source, c.RefNumber))
			}
		}
	}

	// Add citation statistics
	stats, _ := r.GetCitationStats(bookID)
	if stats != nil {
		bibliography.WriteString("## Citation Statistics\n\n")
		bibliography.WriteString("| Source Type | Count | Percentage |\n")
		bibliography.WriteString("|-------------|-------|------------|\n")

		for _, t := range types {
			count, ok := stats.ByType[t]
			if ok && count > 0 {
				percentage := float64(count) / float64(stats.TotalCount) * 100
				bibliography.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n",
					formatTypeLabel(t), count, percentage))
			}
		}
	}

	// Update bibliography metadata
	meta := models.BibliographyMetadata{
		BookID:        bookID,
		Title:         fmt.Sprintf("Bibliography for %s", book.Title),
		Description:   "Comprehensive bibliography of all cited works",
		LastGenerated: time.Now(),
	}
	
	// Upsert the metadata
	r.DB.Where("book_id = ?", bookID).Delete(&models.BibliographyMetadata{})
	r.DB.Create(&meta)

	return bibliography.String(), nil
}

// formatTypeLabel converts database citation type to display label
func formatTypeLabel(t string) string {
	switch t {
	case "book":
		return "Books"
	case "journal":
		return "Journal Articles"
	case "report":
		return "Research Reports"
	case "government":
		return "Government Sources"
	case "interview":
		return "Interviews/Focus Groups"
	case "survey":
		return "Survey Data"
	case "media":
		return "Media Sources"
	default:
		return strings.Title(t)
	}
}