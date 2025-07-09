package repository

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/models"
        "gorm.io/gorm"
        "time"
)

// SectionRepository defines the interface for section data operations
type SectionRepository interface {
        CreateSection(section *models.BookSection) error
        GetSectionByID(id uint) (*models.BookSection, error)
        GetSectionsByChapterID(chapterID uint, includeUnpublished bool) ([]models.BookSection, error)
        UpdateSection(section *models.BookSection) error
        DeleteSection(id uint) error
        
        // Revision operations
        CreateSectionRevision(revision *models.SectionRevision) error
        GetSectionRevisions(sectionID uint) ([]models.SectionRevision, error)
        GetSectionRevisionByID(id uint) (*models.SectionRevision, error)
        
        // Publishing operations
        GetScheduledSections() ([]models.BookSection, error)
}

// GormSectionRepository implements SectionRepository using GORM
type GormSectionRepository struct {
        db *gorm.DB
}

// NewGormSectionRepository creates a new GormSectionRepository
func NewGormSectionRepository(db *gorm.DB) *GormSectionRepository {
        return &GormSectionRepository{db: db}
}

// CreateSection creates a new section in the database
func (r *GormSectionRepository) CreateSection(section *models.BookSection) error {
        return r.db.Create(section).Error
}

// GetSectionByID retrieves a section by its ID
func (r *GormSectionRepository) GetSectionByID(id uint) (*models.BookSection, error) {
        var section models.BookSection
        if err := r.db.First(&section, id).Error; err != nil {
                return nil, err
        }
        return &section, nil
}

// GetSectionsByChapterID retrieves all sections for a specific chapter
func (r *GormSectionRepository) GetSectionsByChapterID(chapterID uint, includeUnpublished bool) ([]models.BookSection, error) {
        var sections []models.BookSection
        query := r.db.Where("chapter_id = ?", chapterID)
        
        if !includeUnpublished {
                query = query.Where("is_published = ?", true)
        }
        
        if err := query.Order("section_number").Find(&sections).Error; err != nil {
                return nil, err
        }
        
        return sections, nil
}

// UpdateSection updates an existing section in the database
func (r *GormSectionRepository) UpdateSection(section *models.BookSection) error {
        return r.db.Save(section).Error
}

// DeleteSection deletes a section from the database
func (r *GormSectionRepository) DeleteSection(id uint) error {
        return r.db.Delete(&models.BookSection{}, id).Error
}

// CreateSectionRevision creates a new section revision in the database
func (r *GormSectionRepository) CreateSectionRevision(revision *models.SectionRevision) error {
        return r.db.Create(revision).Error
}

// GetSectionRevisions retrieves all revisions for a specific section
func (r *GormSectionRepository) GetSectionRevisions(sectionID uint) ([]models.SectionRevision, error) {
        var revisions []models.SectionRevision
        if err := r.db.Where("section_id = ?", sectionID).Order("created_at DESC").Find(&revisions).Error; err != nil {
                return nil, err
        }
        return revisions, nil
}

// GetSectionRevisionByID retrieves a specific section revision by ID
func (r *GormSectionRepository) GetSectionRevisionByID(id uint) (*models.SectionRevision, error) {
        var revision models.SectionRevision
        if err := r.db.First(&revision, id).Error; err != nil {
                return nil, err
        }
        return &revision, nil
}

// GetScheduledSections retrieves all sections scheduled to be published
func (r *GormSectionRepository) GetScheduledSections() ([]models.BookSection, error) {
        var sections []models.BookSection
        if err := r.db.Where("scheduled_publish_at IS NOT NULL AND scheduled_publish_at > ?", time.Now()).Find(&sections).Error; err != nil {
                return nil, err
        }
        return sections, nil
}