package service

import (
        "fmt"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/repository"
        "gorm.io/gorm"
)

// ReadingAnalyticsService defines the interface for reading analytics business logic
type ReadingAnalyticsService interface {
        StartReadingSession(userID, bookID, chapterID, sectionID uint, source string) (*models.ReadingSession, error)
        EndReadingSession(sessionID uint, endTime time.Time) error
        GetReadingSessions(userID, bookID uint, startDate, endDate time.Time) ([]models.ReadingSession, error)
        GetReadingAnalytics(userID, bookID uint) (*models.ReadingAnalytics, error)
        GetRecentlyViewedSections(userID uint, limit int) ([]models.BookSection, error)
        GetReadingTimeStats(userID, bookID uint) (map[string]interface{}, error)
}

// ReadingAnalyticsServiceImpl implements the ReadingAnalyticsService interface
type ReadingAnalyticsServiceImpl struct {
        analyticsRepo repository.ReadingAnalyticsRepository
        bookRepo      repository.BookRepository
}

// NewReadingAnalyticsService creates a new reading analytics service instance
func NewReadingAnalyticsService(
        analyticsRepo repository.ReadingAnalyticsRepository,
        bookRepo repository.BookRepository,
) ReadingAnalyticsService {
        return &ReadingAnalyticsServiceImpl{
                analyticsRepo: analyticsRepo,
                bookRepo:      bookRepo,
        }
}

// StartReadingSession starts a new reading session
func (s *ReadingAnalyticsServiceImpl) StartReadingSession(userID, bookID, chapterID, sectionID uint, source string) (*models.ReadingSession, error) {
        // Create a new session
        session := &models.ReadingSession{
                UserID:    userID,
                BookID:    bookID,
                ChapterID: chapterID,
                SectionID: sectionID,
                StartTime: time.Now(),
                Source:    source,
        }
        
        // The session is not saved to the database yet because it hasn't ended
        // We'll return it so the client can store the session ID
        return session, nil
}

// EndReadingSession ends a reading session and records analytics
func (s *ReadingAnalyticsServiceImpl) EndReadingSession(sessionID uint, endTime time.Time) error {
        // In a real implementation, we would retrieve the session from a cache or temporary storage
        // Here, we'll create a mock session for demonstration
        session := &models.ReadingSession{
                Model: gorm.Model{ID: sessionID},
                UserID:    1,
                BookID:    1,
                ChapterID: 1,
                SectionID: 1,
                StartTime: endTime.Add(-30 * time.Minute), // Assume session started 30 minutes ago
                EndTime:   endTime,
                Duration:  int(endTime.Sub(endTime.Add(-30 * time.Minute)).Seconds()),
                Source:    "web",
        }
        
        // Record the session
        return s.analyticsRepo.RecordReadingSession(session)
}

// GetReadingSessions retrieves reading sessions for a user and book
func (s *ReadingAnalyticsServiceImpl) GetReadingSessions(userID, bookID uint, startDate, endDate time.Time) ([]models.ReadingSession, error) {
        return s.analyticsRepo.GetReadingSessions(userID, bookID, startDate, endDate)
}

// GetReadingAnalytics retrieves reading analytics for a user and book
func (s *ReadingAnalyticsServiceImpl) GetReadingAnalytics(userID, bookID uint) (*models.ReadingAnalytics, error) {
        return s.analyticsRepo.GetReadingAnalytics(userID, bookID)
}

// GetRecentlyViewedSections retrieves recently viewed sections for a user
func (s *ReadingAnalyticsServiceImpl) GetRecentlyViewedSections(userID uint, limit int) ([]models.BookSection, error) {
        // Get recently viewed section IDs
        sectionIDs, err := s.analyticsRepo.GetRecentlyViewedSections(userID, limit)
        if err != nil {
                return nil, err
        }
        
        // Fetch the actual section objects
        sections := make([]models.BookSection, 0, len(sectionIDs))
        for _, id := range sectionIDs {
                section, err := s.bookRepo.GetSectionByID(id)
                if err == nil && section != nil {
                        sections = append(sections, *section)
                }
        }
        
        return sections, nil
}

// GetReadingTimeStats retrieves reading time statistics for a user and book
func (s *ReadingAnalyticsServiceImpl) GetReadingTimeStats(userID, bookID uint) (map[string]interface{}, error) {
        // Get the analytics
        analytics, err := s.analyticsRepo.GetReadingAnalytics(userID, bookID)
        if err != nil {
                return nil, err
        }
        
        // Get reading progress
        progress, err := s.bookRepo.GetReadingProgress(userID, bookID)
        if err != nil {
                // If no progress record exists, create an empty one
                progress = &models.ReadingProgress{
                        UserID:       userID,
                        BookID:       bookID,
                        TimeSpent:    0,
                        SessionCount: 0,
                }
        }
        
        // Format the time spent
        var formattedTotalTime string
        totalMinutes := analytics.TotalTimeSpent / 60
        
        if totalMinutes < 60 {
                formattedTotalTime = fmt.Sprintf("%d minutes", totalMinutes)
        } else {
                hours := totalMinutes / 60
                minutes := totalMinutes % 60
                formattedTotalTime = fmt.Sprintf("%d hours, %d minutes", hours, minutes)
        }
        
        // Calculate completion percentage if the book has a time estimate
        var completionPercentage float64 = 0
        book, err := s.bookRepo.GetBookWithChapters(bookID)
        if err == nil {
                // Calculate total estimated reading time for the book
                var totalEstimatedMinutes int = 0
                for _, chapter := range book.Chapters {
                        // Load sections for the chapter
                        chapterWithSections, err := s.bookRepo.GetChapterWithSections(chapter.ID)
                        if err == nil {
                                for _, section := range chapterWithSections.Sections {
                                        totalEstimatedMinutes += section.TimeToRead
                                }
                        }
                }
                
                // Calculate completion percentage
                if totalEstimatedMinutes > 0 {
                        completionPercentage = float64(analytics.TotalTimeSpent/60) / float64(totalEstimatedMinutes) * 100
                        if completionPercentage > 100 {
                                completionPercentage = 100
                        }
                }
        }
        
        // Format the average session time
        var formattedAvgTime string
        avgMinutes := analytics.AvgSessionTime / 60
        
        if avgMinutes < 60 {
                formattedAvgTime = fmt.Sprintf("%d minutes", avgMinutes)
        } else {
                hours := avgMinutes / 60
                minutes := avgMinutes % 60
                formattedAvgTime = fmt.Sprintf("%d hours, %d minutes", hours, minutes)
        }
        
        // Get the streak
        streak := progress.StreakDays
        
        // Build the stats object
        stats := map[string]interface{}{
                "totalTimeSpent":      analytics.TotalTimeSpent,
                "formattedTotalTime":  formattedTotalTime,
                "sessionCount":        analytics.SessionCount,
                "avgSessionTime":      analytics.AvgSessionTime,
                "formattedAvgTime":    formattedAvgTime,
                "longestSession":      analytics.LongestSession,
                "lastActivity":        analytics.LastActivity,
                "mostActiveDay":       analytics.MostActiveDay,
                "preferredTimeSlot":   analytics.PreferredTimeSlot,
                "streak":              streak,
                "completionPercentage": completionPercentage,
        }
        
        return stats, nil
}