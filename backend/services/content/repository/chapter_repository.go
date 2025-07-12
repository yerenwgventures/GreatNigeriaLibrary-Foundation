package repository

import (
        "fmt"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "gorm.io/gorm"
        "time"
)

// ChapterRepository defines the interface for chapter data operations
type ChapterRepository interface {
        CreateChapter(chapter *models.BookChapter) error
        GetChapterByID(id uint) (*models.BookChapter, error)
        GetChapterWithSections(id uint) (*models.BookChapter, error)
        GetChaptersByBookID(bookID uint, includeUnpublished bool) ([]models.BookChapter, error)
        UpdateChapter(chapter *models.BookChapter) error
        DeleteChapter(id uint) error
        
        // Revision operations
        CreateChapterRevision(revision *models.ChapterRevision) error
        GetChapterRevisions(chapterID uint) ([]models.ChapterRevision, error)
        GetChapterRevisionByID(id uint) (*models.ChapterRevision, error)
        
        // Publishing operations
        GetScheduledChapters() ([]models.BookChapter, error)
}

// GormChapterRepository implements ChapterRepository using GORM
type GormChapterRepository struct {
        db *gorm.DB
}

// NewGormChapterRepository creates a new GormChapterRepository
func NewGormChapterRepository(db *gorm.DB) *GormChapterRepository {
        return &GormChapterRepository{db: db}
}

// CreateChapter creates a new chapter in the database
func (r *GormChapterRepository) CreateChapter(chapter *models.BookChapter) error {
        return r.db.Create(chapter).Error
}

// GetChapterByID retrieves a chapter by its ID
func (r *GormChapterRepository) GetChapterByID(id uint) (*models.BookChapter, error) {
        var chapter models.BookChapter
        if err := r.db.First(&chapter, id).Error; err != nil {
                return nil, err
        }
        return &chapter, nil
}

// GetChaptersByBookID retrieves all chapters for a specific book
func (r *GormChapterRepository) GetChaptersByBookID(bookID uint, includeUnpublished bool) ([]models.BookChapter, error) {
        var chapters []models.BookChapter
        query := r.db.Where("book_id = ?", bookID)
        
        if !includeUnpublished {
                query = query.Where("is_published = ?", true)
        }
        
        if err := query.Order("chapter_number").Find(&chapters).Error; err != nil {
                return nil, err
        }
        
        return chapters, nil
}

// UpdateChapter updates an existing chapter in the database
func (r *GormChapterRepository) UpdateChapter(chapter *models.BookChapter) error {
        return r.db.Save(chapter).Error
}

// DeleteChapter deletes a chapter from the database
func (r *GormChapterRepository) DeleteChapter(id uint) error {
        return r.db.Delete(&models.BookChapter{}, id).Error
}

// GetChapterWithSections retrieves a chapter with all its sections
func (r *GormChapterRepository) GetChapterWithSections(id uint) (*models.BookChapter, error) {
        // Get the chapter
        chapter, err := r.GetChapterByID(id)
        if err != nil {
                return nil, err
        }
        
        // Get sections for the chapter
        var sections []models.BookSection
        if err := r.db.Where("chapter_id = ?", id).Order("order_index").Find(&sections).Error; err != nil {
                return nil, fmt.Errorf("error fetching sections for chapter %d: %w", id, err)
        }
        
        // Assign sections to the chapter
        chapter.Sections = sections
        
        return chapter, nil
}

// CreateChapterRevision creates a new chapter revision in the database
func (r *GormChapterRepository) CreateChapterRevision(revision *models.ChapterRevision) error {
        return r.db.Create(revision).Error
}

// GetChapterRevisions retrieves all revisions for a specific chapter
func (r *GormChapterRepository) GetChapterRevisions(chapterID uint) ([]models.ChapterRevision, error) {
        var revisions []models.ChapterRevision
        if err := r.db.Where("chapter_id = ?", chapterID).Order("created_at DESC").Find(&revisions).Error; err != nil {
                return nil, err
        }
        return revisions, nil
}

// GetChapterRevisionByID retrieves a specific chapter revision by ID
func (r *GormChapterRepository) GetChapterRevisionByID(id uint) (*models.ChapterRevision, error) {
        var revision models.ChapterRevision
        if err := r.db.First(&revision, id).Error; err != nil {
                return nil, err
        }
        return &revision, nil
}

// GetScheduledChapters retrieves all chapters scheduled to be published
func (r *GormChapterRepository) GetScheduledChapters() ([]models.BookChapter, error) {
        var chapters []models.BookChapter
        if err := r.db.Where("scheduled_publish_at IS NOT NULL AND scheduled_publish_at > ?", time.Now()).Find(&chapters).Error; err != nil {
                return nil, err
        }
        return chapters, nil
}