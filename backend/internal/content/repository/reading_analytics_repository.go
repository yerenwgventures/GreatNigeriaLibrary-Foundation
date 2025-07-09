package repository

import (
        "errors"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/models"
        "gorm.io/gorm"
)

// ReadingAnalyticsRepository defines the interface for reading analytics operations
type ReadingAnalyticsRepository interface {
        RecordReadingSession(session *models.ReadingSession) error
        GetReadingSessions(userID, bookID uint, startDate, endDate time.Time) ([]models.ReadingSession, error)
        GetReadingAnalytics(userID, bookID uint) (*models.ReadingAnalytics, error)
        GetRecentlyViewedSections(userID uint, limit int) ([]uint, error)
        UpdateReadingAnalytics(userID, bookID uint) error
}

// GormReadingAnalyticsRepository implements the ReadingAnalyticsRepository interface with GORM
type GormReadingAnalyticsRepository struct {
        db *gorm.DB
}

// NewGormReadingAnalyticsRepository creates a new reading analytics repository instance
func NewGormReadingAnalyticsRepository(db *gorm.DB) *GormReadingAnalyticsRepository {
        return &GormReadingAnalyticsRepository{db: db}
}

// RecordReadingSession records a new reading session
func (r *GormReadingAnalyticsRepository) RecordReadingSession(session *models.ReadingSession) error {
        // Calculate the duration if not set
        if session.Duration == 0 && !session.EndTime.IsZero() && !session.StartTime.IsZero() {
                session.Duration = int(session.EndTime.Sub(session.StartTime).Seconds())
        }
        
        // Create the session
        if err := r.db.Create(session).Error; err != nil {
                return err
        }
        
        // Update reading progress with time spent
        var progress models.ReadingProgress
        result := r.db.Where("user_id = ? AND book_id = ?", session.UserID, session.BookID).First(&progress)
        
        if result.Error != nil {
                if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                        // Create a new progress record
                        progress = models.ReadingProgress{
                                UserID:       session.UserID,
                                BookID:       session.BookID,
                                ChapterID:    session.ChapterID,
                                SectionID:    session.SectionID,
                                LastReadAt:   session.EndTime,
                                LastReadDay:  session.EndTime.Truncate(24 * time.Hour),
                                TimeSpent:    session.Duration,
                                SessionCount: 1,
                        }
                        return r.db.Create(&progress).Error
                }
                return result.Error
        }
        
        // Update existing progress
        progress.TimeSpent += session.Duration
        progress.SessionCount++
        progress.LastReadAt = session.EndTime
        progress.LastReadDay = session.EndTime.Truncate(24 * time.Hour)
        
        // If the session has chapter and section info, update the progress position
        if session.ChapterID > 0 && session.SectionID > 0 {
                progress.ChapterID = session.ChapterID
                progress.SectionID = session.SectionID
        }
        
        return r.db.Save(&progress).Error
}

// GetReadingSessions retrieves reading sessions for a user and book within the given date range
func (r *GormReadingAnalyticsRepository) GetReadingSessions(userID, bookID uint, startDate, endDate time.Time) ([]models.ReadingSession, error) {
        var sessions []models.ReadingSession
        
        query := r.db.Where("user_id = ?", userID)
        
        // Apply book filter if specified
        if bookID > 0 {
                query = query.Where("book_id = ?", bookID)
        }
        
        // Apply date range filters
        if !startDate.IsZero() {
                query = query.Where("start_time >= ?", startDate)
        }
        
        if !endDate.IsZero() {
                query = query.Where("end_time <= ?", endDate)
        }
        
        // Get sessions ordered by start time descending
        result := query.Order("start_time DESC").Find(&sessions)
        
        return sessions, result.Error
}

// GetReadingAnalytics retrieves or generates reading analytics for a user and book
func (r *GormReadingAnalyticsRepository) GetReadingAnalytics(userID, bookID uint) (*models.ReadingAnalytics, error) {
        var analytics models.ReadingAnalytics
        
        // Try to find existing analytics
        result := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&analytics)
        
        if result.Error != nil {
                if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                        // If not found, generate new analytics
                        err := r.UpdateReadingAnalytics(userID, bookID)
                        if err != nil {
                                return nil, err
                        }
                        
                        // Then fetch the newly created analytics
                        result = r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&analytics)
                        if result.Error != nil {
                                return nil, result.Error
                        }
                } else {
                        return nil, result.Error
                }
        }
        
        return &analytics, nil
}

// GetRecentlyViewedSections retrieves the section IDs of recently viewed sections
func (r *GormReadingAnalyticsRepository) GetRecentlyViewedSections(userID uint, limit int) ([]uint, error) {
        if limit <= 0 {
                limit = 5 // Default limit
        }
        
        var sessions []models.ReadingSession
        result := r.db.
                Where("user_id = ? AND section_id > 0", userID).
                Order("end_time DESC").
                Limit(limit).
                Find(&sessions)
        
        if result.Error != nil {
                return nil, result.Error
        }
        
        // Extract unique section IDs
        sectionIDs := make([]uint, 0, limit)
        sectionMap := make(map[uint]bool)
        
        for _, session := range sessions {
                if !sectionMap[session.SectionID] {
                        sectionMap[session.SectionID] = true
                        sectionIDs = append(sectionIDs, session.SectionID)
                }
        }
        
        return sectionIDs, nil
}

// UpdateReadingAnalytics generates or updates reading analytics for a user and book
func (r *GormReadingAnalyticsRepository) UpdateReadingAnalytics(userID, bookID uint) error {
        // Get all sessions for this user and book
        var sessions []models.ReadingSession
        result := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Find(&sessions)
        if result.Error != nil {
                return result.Error
        }
        
        if len(sessions) == 0 {
                return errors.New("no reading sessions found")
        }
        
        // Calculate analytics
        totalTime := 0
        longestSession := 0
        sessionCount := len(sessions)
        var lastActivity time.Time
        
        // Day of week counters
        dayCounters := make(map[string]int)
        timeSlotCounters := make(map[string]int)
        
        for _, session := range sessions {
                // Sum total time
                totalTime += session.Duration
                
                // Track longest session
                if session.Duration > longestSession {
                        longestSession = session.Duration
                }
                
                // Track last activity
                if session.EndTime.After(lastActivity) {
                        lastActivity = session.EndTime
                }
                
                // Count activity by day of week
                day := session.StartTime.Weekday().String()
                dayCounters[day]++
                
                // Count activity by time slot
                hour := session.StartTime.Hour()
                var timeSlot string
                
                switch {
                case hour >= 5 && hour < 12:
                        timeSlot = "Morning"
                case hour >= 12 && hour < 17:
                        timeSlot = "Afternoon"
                case hour >= 17 && hour < 21:
                        timeSlot = "Evening"
                default:
                        timeSlot = "Night"
                }
                
                timeSlotCounters[timeSlot]++
        }
        
        // Find most active day
        mostActiveDay := ""
        maxDayCount := 0
        for day, count := range dayCounters {
                if count > maxDayCount {
                        maxDayCount = count
                        mostActiveDay = day
                }
        }
        
        // Find preferred time slot
        preferredTimeSlot := ""
        maxTimeSlotCount := 0
        for timeSlot, count := range timeSlotCounters {
                if count > maxTimeSlotCount {
                        maxTimeSlotCount = count
                        preferredTimeSlot = timeSlot
                }
        }
        
        // Calculate average session time
        avgSessionTime := 0
        if sessionCount > 0 {
                avgSessionTime = totalTime / sessionCount
        }
        
        // Create or update analytics
        var analytics models.ReadingAnalytics
        
        result = r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&analytics)
        
        if result.Error != nil {
                if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                        // Create new analytics
                        analytics = models.ReadingAnalytics{
                                UserID:            userID,
                                BookID:            bookID,
                                TotalTimeSpent:    totalTime,
                                AvgSessionTime:    avgSessionTime,
                                LongestSession:    longestSession,
                                SessionCount:      sessionCount,
                                LastActivity:      lastActivity,
                                MostActiveDay:     mostActiveDay,
                                PreferredTimeSlot: preferredTimeSlot,
                        }
                        
                        return r.db.Create(&analytics).Error
                }
                return result.Error
        }
        
        // Update existing analytics
        analytics.TotalTimeSpent = totalTime
        analytics.AvgSessionTime = avgSessionTime
        analytics.LongestSession = longestSession
        analytics.SessionCount = sessionCount
        analytics.LastActivity = lastActivity
        analytics.MostActiveDay = mostActiveDay
        analytics.PreferredTimeSlot = preferredTimeSlot
        
        return r.db.Save(&analytics).Error
}