package service

import (
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/repository"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// ProgressService handles progress-related business logic
type ProgressService struct {
        progressRepo repository.ProgressRepository
        logger       *logger.Logger
}

// NewProgressService creates a new progress service
func NewProgressService(progressRepo repository.ProgressRepository, logger *logger.Logger) *ProgressService {
        return &ProgressService{
                progressRepo: progressRepo,
                logger:       logger,
        }
}

// UpdateProgress updates a user's reading progress
func (s *ProgressService) UpdateProgress(userID, bookID, chapterID, sectionID uint, lastPosition int, percentComplete float64) (*models.UserProgressResponse, error) {
        progress, err := s.progressRepo.UpdateProgress(userID, bookID, chapterID, sectionID, lastPosition, percentComplete)
        if err != nil {
                s.logger.WithError(err).Error("Failed to update progress")
                return nil, errors.ErrInternalServer("Failed to update progress")
        }

        return &models.UserProgressResponse{
                BookID:          progress.BookID,
                ChapterID:       progress.ChapterID,
                SectionID:       progress.SectionID,
                LastPosition:    progress.LastPosition,
                PercentComplete: progress.PercentComplete,
                LastUpdated:     progress.LastUpdated,
        }, nil
}

// GetProgress gets a user's reading progress for a book
func (s *ProgressService) GetProgress(userID, bookID uint) (*models.UserProgressResponse, error) {
        progress, err := s.progressRepo.GetByUserAndBook(userID, bookID)
        if err != nil {
                s.logger.WithError(err).Error("Failed to get progress")
                return nil, errors.ErrInternalServer("Failed to get progress")
        }

        if progress == nil {
                return nil, nil
        }

        return &models.UserProgressResponse{
                BookID:          progress.BookID,
                ChapterID:       progress.ChapterID,
                SectionID:       progress.SectionID,
                LastPosition:    progress.LastPosition,
                PercentComplete: progress.PercentComplete,
                LastUpdated:     progress.LastUpdated,
        }, nil
}

// CheckCompletion checks if a book is completed
func (s *ProgressService) CheckCompletion(userID, bookID uint) (bool, error) {
        isCompleted, err := s.progressRepo.CheckBookCompletion(userID, bookID)
        if err != nil {
                s.logger.WithError(err).Error("Failed to check book completion")
                return false, errors.ErrInternalServer("Failed to check book completion")
        }
        return isCompleted, nil
}

// MarkTopicCompleted marks a topic as completed and records completion time
func (s *ProgressService) MarkTopicCompleted(userID, bookID, chapterID, sectionID uint, points int) (*models.TopicCompletion, error) {
        // For a real implementation, this would interact with the database
        // For now, we'll just create and return a mock completion object
        completion := &models.TopicCompletion{
                UserID:      userID,
                BookID:      bookID,
                ChapterID:   chapterID,
                SectionID:   sectionID,
                Points:      points,
                CompletedAt: time.Now(),
        }
        
        return completion, nil
}
