package repository

import (
        "time"
        
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
        "gorm.io/gorm"
)

// ProgressRepository defines the interface for user progress data operations
type ProgressRepository interface {
        CreateOrUpdateProgress(progress *models.UserProgress) error
        GetUserProgress(userID, bookID uint, contentType string) (*models.UserProgress, error)
        GetAllUserProgress(userID uint) ([]models.UserProgress, error)
        DeleteUserProgress(progressID uint) error
        
        // Additional methods required by ProgressService
        UpdateProgress(userID, bookID, chapterID, sectionID uint, lastPosition int, percentComplete float64) (*models.UserProgress, error)
        GetByUserAndBook(userID, bookID uint) (*models.UserProgress, error)
        CheckBookCompletion(userID, bookID uint) (bool, error)
}

// GormProgressRepository implements ProgressRepository using GORM
type GormProgressRepository struct {
        db *gorm.DB
}

// NewGormProgressRepository creates a new GormProgressRepository
func NewGormProgressRepository(db *gorm.DB) *GormProgressRepository {
        return &GormProgressRepository{db: db}
}

// CreateOrUpdateProgress creates or updates a user progress record
func (r *GormProgressRepository) CreateOrUpdateProgress(progress *models.UserProgress) error {
        // Check if a record already exists
        var existingProgress models.UserProgress
        result := r.db.Where("user_id = ? AND book_id = ?",
                progress.UserID, progress.BookID).First(&existingProgress)

        if result.Error == nil {
                // Record exists, update it
                existingProgress.ChapterID = progress.ChapterID
                existingProgress.SectionID = progress.SectionID
                existingProgress.LastPosition = progress.LastPosition
                existingProgress.PercentComplete = progress.PercentComplete
                existingProgress.LastUpdated = progress.LastUpdated
                return r.db.Save(&existingProgress).Error
        } else if result.Error == gorm.ErrRecordNotFound {
                // Record doesn't exist, create new one
                return r.db.Create(progress).Error
        } else {
                // Some other error occurred
                return result.Error
        }
}

// GetUserProgress retrieves a user's progress for a specific content item
func (r *GormProgressRepository) GetUserProgress(userID, bookID uint, contentType string) (*models.UserProgress, error) {
        var progress models.UserProgress
        if err := r.db.Where("user_id = ? AND book_id = ?",
                userID, bookID).First(&progress).Error; err != nil {
                return nil, err
        }
        return &progress, nil
}

// GetAllUserProgress retrieves all progress records for a specific user
func (r *GormProgressRepository) GetAllUserProgress(userID uint) ([]models.UserProgress, error) {
        var progressList []models.UserProgress
        if err := r.db.Where("user_id = ?", userID).Find(&progressList).Error; err != nil {
                return nil, err
        }
        return progressList, nil
}

// DeleteUserProgress deletes a progress record
func (r *GormProgressRepository) DeleteUserProgress(progressID uint) error {
        return r.db.Delete(&models.UserProgress{}, progressID).Error
}

// UpdateProgress updates a user's reading progress
func (r *GormProgressRepository) UpdateProgress(userID, bookID, chapterID, sectionID uint, lastPosition int, percentComplete float64) (*models.UserProgress, error) {
        // Check if progress exists
        var progress models.UserProgress
        err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&progress).Error
        
        if err == gorm.ErrRecordNotFound {
                // Create new progress
                progress = models.UserProgress{
                        UserID:          userID,
                        BookID:          bookID,
                        ChapterID:       chapterID,
                        SectionID:       sectionID,
                        LastPosition:    lastPosition, 
                        PercentComplete: percentComplete,
                        LastUpdated:     time.Now(),
                }
                if err := r.db.Create(&progress).Error; err != nil {
                        return nil, err
                }
        } else if err != nil {
                return nil, err
        } else {
                // Update existing progress
                progress.ChapterID = chapterID
                progress.SectionID = sectionID
                progress.LastPosition = lastPosition
                progress.PercentComplete = percentComplete
                progress.LastUpdated = time.Now()
                
                if err := r.db.Save(&progress).Error; err != nil {
                        return nil, err
                }
        }
        
        return &progress, nil
}

// GetByUserAndBook retrieves progress for a specific user and book
func (r *GormProgressRepository) GetByUserAndBook(userID, bookID uint) (*models.UserProgress, error) {
        var progress models.UserProgress
        err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&progress).Error
        
        if err == gorm.ErrRecordNotFound {
                return nil, nil
        } else if err != nil {
                return nil, err
        }
        
        return &progress, nil
}

// CheckBookCompletion checks if a book is completed by the user (over 95% complete)
func (r *GormProgressRepository) CheckBookCompletion(userID, bookID uint) (bool, error) {
        var progress models.UserProgress
        err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&progress).Error
        
        if err == gorm.ErrRecordNotFound {
                return false, nil
        } else if err != nil {
                return false, err
        }
        
        // Consider a book completed if progress is over 95%
        return progress.PercentComplete >= 95.0, nil
}