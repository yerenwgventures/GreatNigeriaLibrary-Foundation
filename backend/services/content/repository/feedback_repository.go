package repository

import (
        "context"
        "errors"
        "time"

        "github.com/sirupsen/logrus"
        "gorm.io/gorm"
)

// FeedbackType represents the type of feedback
type FeedbackType string

const (
        // FeedbackTypeMood represents mood feedback
        FeedbackTypeMood FeedbackType = "mood"
        // FeedbackTypeDifficulty represents difficulty feedback
        FeedbackTypeDifficulty FeedbackType = "difficulty"
)

// MoodType represents different mood categories
type MoodType string

const (
        // MoodTypeHappy represents a happy mood (😄)
        MoodTypeHappy MoodType = "happy"
        // MoodTypeInspired represents an inspired mood (😮)
        MoodTypeInspired MoodType = "inspired"  
        // MoodTypeConfused represents a confused mood (😕)
        MoodTypeConfused MoodType = "confused"
        // MoodTypeFrustrated represents a frustrated mood (😣)
        MoodTypeFrustrated MoodType = "frustrated"
        // MoodTypeMotivated represents a motivated mood (🚀)
        MoodTypeMotivated MoodType = "motivated"
)

// MoodEmoji maps mood types to emoji representations
var MoodEmoji = map[MoodType]string{
        MoodTypeHappy:      "😄",
        MoodTypeInspired:   "😮",
        MoodTypeConfused:   "😕",
        MoodTypeFrustrated: "😣",
        MoodTypeMotivated:  "🚀",
}

// DifficultyCategory represents content difficulty categories
type DifficultyCategory string

const (
        // DifficultyCategoryEasy represents easy content (🟢)
        DifficultyCategoryEasy DifficultyCategory = "easy"
        // DifficultyCategoryModerate represents moderate content (🟡)
        DifficultyCategoryModerate DifficultyCategory = "moderate"
        // DifficultyCategoryHard represents hard content (🟠)
        DifficultyCategoryHard DifficultyCategory = "hard"
        // DifficultyCategoryComplex represents complex content (🔴)
        DifficultyCategoryComplex DifficultyCategory = "complex"
)

// DifficultyEmoji maps difficulty categories to emoji representations
var DifficultyEmoji = map[DifficultyCategory]string{
        DifficultyCategoryEasy:     "🟢",
        DifficultyCategoryModerate: "🟡",
        DifficultyCategoryHard:     "🟠",
        DifficultyCategoryComplex:  "🔴",
}

// LearningStyle represents different learning approaches
type LearningStyle string

const (
        // LearningStyleVisual for those who learn better with images, diagrams
        LearningStyleVisual LearningStyle = "visual"
        // LearningStyleAuditory for those who learn better with spoken explanations
        LearningStyleAuditory LearningStyle = "auditory"
        // LearningStyleKinesthetic for those who learn by doing
        LearningStyleKinesthetic LearningStyle = "kinesthetic"
        // LearningStyleReading for those who learn through reading text
        LearningStyleReading LearningStyle = "reading"
)

// ContentFeedback represents a user's feedback for content
type ContentFeedback struct {
        ID             uint      `gorm:"primaryKey"`
        UserID         uint      `gorm:"index;not null"`
        BookID         uint      `gorm:"index;not null"`
        ChapterID      uint      `gorm:"index"`
        SectionID      uint      `gorm:"index"`
        Type           string    `gorm:"size:20;not null"` // "mood" or "difficulty"
        Value          int       `gorm:"not null"`         // 1-5 for both mood and difficulty
        MoodCategory   string    `gorm:"size:20"`          // Optional mood category
        DifficultyCat  string    `gorm:"size:20"`          // Optional difficulty category
        Comment        string    `gorm:"type:text"`
        LearningStyle  string    `gorm:"size:30"`          // User's learning style preference
        RecommendNext  bool      `gorm:"default:false"`    // Whether to recommend next content
        CreatedAt      time.Time `gorm:"not null"`
        UpdatedAt      time.Time `gorm:"not null"`
}

// FeedbackRepository interface defines methods for feedback operations
type FeedbackRepository interface {
        SubmitMoodFeedback(ctx context.Context, userID, bookID, chapterID, sectionID uint, value int, moodCategory string, comment string, learningStyle string) (*ContentFeedback, error)
        SubmitDifficultyFeedback(ctx context.Context, userID, bookID, chapterID, sectionID uint, value int, difficultyCategory string, comment string, recommendNext bool) (*ContentFeedback, error)
        GetUserContentFeedback(ctx context.Context, userID uint) ([]ContentFeedback, error)
        GetUserRecentMoods(ctx context.Context, userID uint, limit int) ([]ContentFeedback, error)
        GetRecommendedContent(ctx context.Context, userID uint, baseContentID uint) ([]uint, error)
        GetContentFeedbackSummary(ctx context.Context, bookID, chapterID, sectionID uint) (map[string]interface{}, error)
        GetDetailedFeedbackAnalysis(ctx context.Context, bookID uint) (map[string]interface{}, error)
        DeleteFeedback(ctx context.Context, feedbackID, userID uint, feedbackType FeedbackType) error
}

// GormFeedbackRepository implements FeedbackRepository using GORM
type GormFeedbackRepository struct {
        db     *gorm.DB
        logger *logrus.Logger
}

// NewFeedbackRepository creates a new feedback repository instance
func NewFeedbackRepository(db *gorm.DB, logger *logrus.Logger) FeedbackRepository {
        return &GormFeedbackRepository{
                db:     db,
                logger: logger,
        }
}

// SubmitMoodFeedback saves a mood feedback record
func (r *GormFeedbackRepository) SubmitMoodFeedback(ctx context.Context, userID, bookID, chapterID, sectionID uint, value int, moodCategory string, comment string, learningStyle string) (*ContentFeedback, error) {
        // Check if user already submitted mood feedback for this content
        var existingFeedback ContentFeedback
        result := r.db.Where("user_id = ? AND book_id = ? AND chapter_id = ? AND section_id = ? AND type = ?",
                userID, bookID, chapterID, sectionID, FeedbackTypeMood).First(&existingFeedback)

        // Map the value to a mood category if not provided
        if moodCategory == "" {
                switch {
                case value == 5:
                        moodCategory = string(MoodTypeMotivated)
                case value == 4:
                        moodCategory = string(MoodTypeHappy)
                case value == 3:
                        moodCategory = string(MoodTypeInspired)
                case value == 2:
                        moodCategory = string(MoodTypeConfused)
                case value == 1:
                        moodCategory = string(MoodTypeFrustrated)
                }
        }

        feedback := ContentFeedback{
                UserID:        userID,
                BookID:        bookID,
                ChapterID:     chapterID,
                SectionID:     sectionID,
                Type:          string(FeedbackTypeMood),
                Value:         value,
                MoodCategory:  moodCategory,
                Comment:       comment,
                LearningStyle: learningStyle,
                CreatedAt:     time.Now(),
                UpdatedAt:     time.Now(),
        }

        // If feedback already exists, update it
        if result.Error == nil {
                existingFeedback.Value = value
                existingFeedback.MoodCategory = moodCategory
                existingFeedback.Comment = comment
                existingFeedback.LearningStyle = learningStyle
                existingFeedback.UpdatedAt = time.Now()
                
                if err := r.db.Save(&existingFeedback).Error; err != nil {
                        r.logger.WithError(err).Error("Failed to update mood feedback")
                        return nil, err
                }
                return &existingFeedback, nil
        }

        // Create new feedback
        if err := r.db.Create(&feedback).Error; err != nil {
                r.logger.WithError(err).Error("Failed to create mood feedback")
                return nil, err
        }

        return &feedback, nil
}

// SubmitDifficultyFeedback saves a difficulty feedback record
func (r *GormFeedbackRepository) SubmitDifficultyFeedback(ctx context.Context, userID, bookID, chapterID, sectionID uint, value int, difficultyCategory string, comment string, recommendNext bool) (*ContentFeedback, error) {
        // Check if user already submitted difficulty feedback for this content
        var existingFeedback ContentFeedback
        result := r.db.Where("user_id = ? AND book_id = ? AND chapter_id = ? AND section_id = ? AND type = ?",
                userID, bookID, chapterID, sectionID, FeedbackTypeDifficulty).First(&existingFeedback)

        // Map the value to a difficulty category if not provided
        if difficultyCategory == "" {
                switch {
                case value == 1:
                        difficultyCategory = string(DifficultyCategoryEasy)
                case value == 2:
                        difficultyCategory = string(DifficultyCategoryEasy)
                case value == 3:
                        difficultyCategory = string(DifficultyCategoryModerate)
                case value == 4:
                        difficultyCategory = string(DifficultyCategoryHard)
                case value == 5:
                        difficultyCategory = string(DifficultyCategoryComplex)
                }
        }

        feedback := ContentFeedback{
                UserID:        userID,
                BookID:        bookID,
                ChapterID:     chapterID,
                SectionID:     sectionID,
                Type:          string(FeedbackTypeDifficulty),
                Value:         value,
                DifficultyCat: difficultyCategory,
                Comment:       comment,
                RecommendNext: recommendNext,
                CreatedAt:     time.Now(),
                UpdatedAt:     time.Now(),
        }

        // If feedback already exists, update it
        if result.Error == nil {
                existingFeedback.Value = value
                existingFeedback.DifficultyCat = difficultyCategory
                existingFeedback.Comment = comment
                existingFeedback.RecommendNext = recommendNext
                existingFeedback.UpdatedAt = time.Now()
                
                if err := r.db.Save(&existingFeedback).Error; err != nil {
                        r.logger.WithError(err).Error("Failed to update difficulty feedback")
                        return nil, err
                }
                return &existingFeedback, nil
        }

        // Create new feedback
        if err := r.db.Create(&feedback).Error; err != nil {
                r.logger.WithError(err).Error("Failed to create difficulty feedback")
                return nil, err
        }

        return &feedback, nil
}

// GetUserContentFeedback retrieves all feedback submitted by a user
func (r *GormFeedbackRepository) GetUserContentFeedback(ctx context.Context, userID uint) ([]ContentFeedback, error) {
        var feedbacks []ContentFeedback
        if err := r.db.Where("user_id = ?", userID).Find(&feedbacks).Error; err != nil {
                r.logger.WithError(err).Error("Failed to retrieve user feedback")
                return nil, err
        }
        return feedbacks, nil
}

// GetUserRecentMoods retrieves the most recent mood feedback from a user
func (r *GormFeedbackRepository) GetUserRecentMoods(ctx context.Context, userID uint, limit int) ([]ContentFeedback, error) {
        var moods []ContentFeedback
        
        if limit <= 0 {
                limit = 10 // Default limit if not specified or invalid
        }
        
        query := r.db.Where("user_id = ? AND type = ?", userID, FeedbackTypeMood).
                Order("created_at DESC").
                Limit(limit)
        
        if err := query.Find(&moods).Error; err != nil {
                r.logger.WithError(err).Error("Failed to retrieve user recent moods")
                return nil, err
        }
        
        return moods, nil
}

// GetRecommendedContent provides content recommendations based on user feedback
func (r *GormFeedbackRepository) GetRecommendedContent(ctx context.Context, userID uint, baseContentID uint) ([]uint, error) {
        // This implementation recommends content based on:
        // 1. Current user's learning style and mood
        // 2. Content with similar difficulty level to what the user has positively engaged with
        // 3. Popular content among users with similar feedback patterns
        
        // First, get user's preferred learning style from most recent mood feedbacks
        var learningStyle string
        var recentFeedback ContentFeedback
        
        err := r.db.Where("user_id = ? AND learning_style != ''", userID).
                Order("created_at DESC").
                First(&recentFeedback).Error
        
        if err == nil {
                learningStyle = recentFeedback.LearningStyle
        }
        
        // Start building recommendations
        var recommendedContentIDs []uint
        
        // 1. Find content with similar difficulty level that other users found helpful
        type RecommendedContent struct {
                SectionID uint
                Count     int
        }
        
        var similarDifficultyContent []RecommendedContent
        
        // Get user's average difficulty rating
        var avgDifficulty float64
        err = r.db.Model(&ContentFeedback{}).
                Select("AVG(value) as avg_difficulty").
                Where("user_id = ? AND type = ?", userID, FeedbackTypeDifficulty).
                Scan(&avgDifficulty).Error
        
        if err != nil {
                avgDifficulty = 3.0 // Default if no data
        }
        
        // Find content with similar difficulty that received positive mood ratings
        difficultyRange := 1.0 // Within +/- 1 of user's average
        
        difficultyQuery := r.db.Model(&ContentFeedback{}).
                Select("section_id, COUNT(*) as count").
                Joins("JOIN content_feedback mood ON mood.section_id = content_feedback.section_id AND mood.type = ?", FeedbackTypeMood).
                Where("content_feedback.type = ? AND content_feedback.value BETWEEN ? AND ? AND mood.value >= 4",
                        FeedbackTypeDifficulty,
                        avgDifficulty-difficultyRange,
                        avgDifficulty+difficultyRange).
                Where("content_feedback.section_id != 0").
                Where("content_feedback.section_id != ?", baseContentID).
                Group("section_id").
                Order("count DESC").
                Limit(5)
        
        if err := difficultyQuery.Find(&similarDifficultyContent).Error; err != nil {
                r.logger.WithError(err).Error("Failed to find similar difficulty content")
        } else {
                for _, content := range similarDifficultyContent {
                        recommendedContentIDs = append(recommendedContentIDs, content.SectionID)
                }
        }
        
        // 2. If user has a learning style preference, find content that matches it
        if learningStyle != "" {
                var learningStyleContent []RecommendedContent
                
                learningStyleQuery := r.db.Model(&ContentFeedback{}).
                        Select("section_id, COUNT(*) as count").
                        Where("learning_style = ? AND type = ? AND value >= 4", learningStyle, FeedbackTypeMood).
                        Where("section_id != 0").
                        Where("section_id != ?", baseContentID).
                        Group("section_id").
                        Order("count DESC").
                        Limit(3)
                
                if err := learningStyleQuery.Find(&learningStyleContent).Error; err == nil {
                        for _, content := range learningStyleContent {
                                // Check if this content is already in our recommendations
                                alreadyIncluded := false
                                for _, id := range recommendedContentIDs {
                                        if id == content.SectionID {
                                                alreadyIncluded = true
                                                break
                                        }
                                }
                                
                                if !alreadyIncluded {
                                        recommendedContentIDs = append(recommendedContentIDs, content.SectionID)
                                }
                        }
                }
        }
        
        // 3. Add some popular content in case we don't have enough recommendations
        if len(recommendedContentIDs) < 5 {
                var popularContent []RecommendedContent
                
                popularQuery := r.db.Model(&ContentFeedback{}).
                        Select("section_id, COUNT(*) as count").
                        Where("type = ? AND value >= 4", FeedbackTypeMood).
                        Where("section_id != 0").
                        Where("section_id != ?", baseContentID).
                        Group("section_id").
                        Order("count DESC").
                        Limit(5)
                
                if err := popularQuery.Find(&popularContent).Error; err == nil {
                        for _, content := range popularContent {
                                alreadyIncluded := false
                                for _, id := range recommendedContentIDs {
                                        if id == content.SectionID {
                                                alreadyIncluded = true
                                                break
                                        }
                                }
                                
                                if !alreadyIncluded && len(recommendedContentIDs) < 5 {
                                        recommendedContentIDs = append(recommendedContentIDs, content.SectionID)
                                }
                        }
                }
        }
        
        return recommendedContentIDs, nil
}

// GetContentFeedbackSummary retrieves aggregated feedback statistics for specific content
func (r *GormFeedbackRepository) GetContentFeedbackSummary(ctx context.Context, bookID, chapterID, sectionID uint) (map[string]interface{}, error) {
        // Base query depending on content level
        query := r.db.Model(&ContentFeedback{}).Select("type, avg(value) as average, count(*) as count, min(value) as min, max(value) as max")
        
        if bookID > 0 {
                query = query.Where("book_id = ?", bookID)
                
                if chapterID > 0 {
                        query = query.Where("chapter_id = ?", chapterID)
                        
                        if sectionID > 0 {
                                query = query.Where("section_id = ?", sectionID)
                        }
                }
        }
        
        // Group by type to get mood and difficulty statistics separately
        query = query.Group("type")
        
        // Execute the query
        type Result struct {
                Type    string  `json:"type"`
                Average float64 `json:"average"`
                Count   int     `json:"count"`
                Min     int     `json:"min"`
                Max     int     `json:"max"`
        }
        
        var results []Result
        if err := query.Find(&results).Error; err != nil {
                r.logger.WithError(err).Error("Failed to retrieve feedback summary")
                return nil, err
        }
        
        // Format the results
        summary := map[string]interface{}{
                "totalFeedback": 0,
                "mood":          map[string]interface{}{},
                "difficulty":    map[string]interface{}{},
        }
        
        totalCount := 0
        for _, result := range results {
                totalCount += result.Count
                
                if result.Type == string(FeedbackTypeMood) {
                        summary["mood"] = map[string]interface{}{
                                "average": result.Average,
                                "count":   result.Count,
                                "min":     result.Min,
                                "max":     result.Max,
                        }
                } else if result.Type == string(FeedbackTypeDifficulty) {
                        summary["difficulty"] = map[string]interface{}{
                                "average": result.Average,
                                "count":   result.Count,
                                "min":     result.Min,
                                "max":     result.Max,
                        }
                }
        }
        
        summary["totalFeedback"] = totalCount
        
        return summary, nil
}

// GetDetailedFeedbackAnalysis provides detailed analysis of feedback for a book
func (r *GormFeedbackRepository) GetDetailedFeedbackAnalysis(ctx context.Context, bookID uint) (map[string]interface{}, error) {
        if bookID == 0 {
                return nil, errors.New("book ID is required")
        }
        
        // Create the result structure
        analysis := map[string]interface{}{
                "bookID":            bookID,
                "totalFeedbackCount": 0,
                "moodAnalysis":      map[string]interface{}{},
                "difficultyAnalysis": map[string]interface{}{},
                "learningStyles":    map[string]interface{}{},
                "contentInsights":   []map[string]interface{}{},
                "emotionalJourney":  []map[string]interface{}{},
        }
        
        // 1. Get total feedback count
        var totalCount int64
        err := r.db.Model(&ContentFeedback{}).Where("book_id = ?", bookID).Count(&totalCount).Error
        if err != nil {
                r.logger.WithError(err).Error("Failed to count total feedback")
                return nil, err
        }
        analysis["totalFeedbackCount"] = totalCount
        
        // 2. Mood category distribution
        var moodCategoryResults []struct {
                Category string
                Count    int
        }
        
        moodCategoryQuery := r.db.Model(&ContentFeedback{}).
                Select("mood_category as category, COUNT(*) as count").
                Where("book_id = ? AND type = ? AND mood_category != ''", bookID, FeedbackTypeMood).
                Group("mood_category").
                Order("count DESC")
        
        if err := moodCategoryQuery.Find(&moodCategoryResults).Error; err != nil {
                r.logger.WithError(err).Error("Failed to retrieve mood category distribution")
        } else {
                moodDistribution := make(map[string]int)
                for _, result := range moodCategoryResults {
                        moodDistribution[result.Category] = result.Count
                }
                analysis["moodAnalysis"].(map[string]interface{})["categoryDistribution"] = moodDistribution
        }
        
        // 3. Difficulty category distribution
        var difficultyCategoryResults []struct {
                Category string
                Count    int
        }
        
        difficultyCategoryQuery := r.db.Model(&ContentFeedback{}).
                Select("difficulty_cat as category, COUNT(*) as count").
                Where("book_id = ? AND type = ? AND difficulty_cat != ''", bookID, FeedbackTypeDifficulty).
                Group("difficulty_cat").
                Order("count DESC")
        
        if err := difficultyCategoryQuery.Find(&difficultyCategoryResults).Error; err != nil {
                r.logger.WithError(err).Error("Failed to retrieve difficulty category distribution")
        } else {
                difficultyDistribution := make(map[string]int)
                for _, result := range difficultyCategoryResults {
                        difficultyDistribution[result.Category] = result.Count
                }
                analysis["difficultyAnalysis"].(map[string]interface{})["categoryDistribution"] = difficultyDistribution
        }
        
        // 4. Learning styles preferences
        var learningStyleResults []struct {
                Style string
                Count int
        }
        
        learningStyleQuery := r.db.Model(&ContentFeedback{}).
                Select("learning_style as style, COUNT(*) as count").
                Where("book_id = ? AND learning_style != ''", bookID).
                Group("learning_style").
                Order("count DESC")
        
        if err := learningStyleQuery.Find(&learningStyleResults).Error; err != nil {
                r.logger.WithError(err).Error("Failed to retrieve learning style preferences")
        } else {
                learningStyleDistribution := make(map[string]int)
                for _, result := range learningStyleResults {
                        learningStyleDistribution[result.Style] = result.Count
                }
                analysis["learningStyles"].(map[string]interface{})["distribution"] = learningStyleDistribution
        }
        
        // 5. Content insights - sections with most engagement and their mood/difficulty patterns
        type ContentInsight struct {
                SectionID      uint
                ChapterID      uint
                FeedbackCount  int
                AvgMood        float64
                AvgDifficulty  float64
                CommonMood     string
                CommonDifficulty string
        }
        
        var contentInsights []ContentInsight
        
        // Get sections with most feedback
        sectionQuery := r.db.Model(&ContentFeedback{}).
                Select("section_id, chapter_id, COUNT(*) as feedback_count").
                Where("book_id = ? AND section_id != 0", bookID).
                Group("section_id, chapter_id").
                Order("feedback_count DESC").
                Limit(10)
        
        var topSections []struct {
                SectionID     uint
                ChapterID     uint
                FeedbackCount int
        }
        
        if err := sectionQuery.Find(&topSections).Error; err != nil {
                r.logger.WithError(err).Error("Failed to retrieve top sections")
        } else {
                for _, section := range topSections {
                        insight := ContentInsight{
                                SectionID:     section.SectionID,
                                ChapterID:     section.ChapterID,
                                FeedbackCount: section.FeedbackCount,
                        }
                        
                        // Get average mood
                        var avgMood float64
                        moodQuery := r.db.Model(&ContentFeedback{}).
                                Select("AVG(value) as avg_mood").
                                Where("book_id = ? AND section_id = ? AND type = ?", 
                                        bookID, section.SectionID, FeedbackTypeMood)
                        
                        if err := moodQuery.Scan(&avgMood).Error; err == nil {
                                insight.AvgMood = avgMood
                        }
                        
                        // Get average difficulty
                        var avgDifficulty float64
                        difficultyQuery := r.db.Model(&ContentFeedback{}).
                                Select("AVG(value) as avg_difficulty").
                                Where("book_id = ? AND section_id = ? AND type = ?", 
                                        bookID, section.SectionID, FeedbackTypeDifficulty)
                        
                        if err := difficultyQuery.Scan(&avgDifficulty).Error; err == nil {
                                insight.AvgDifficulty = avgDifficulty
                        }
                        
                        // Get most common mood category
                        var moodCategory struct {
                                Category string
                                Count    int
                        }
                        
                        commonMoodQuery := r.db.Model(&ContentFeedback{}).
                                Select("mood_category as category, COUNT(*) as count").
                                Where("book_id = ? AND section_id = ? AND type = ? AND mood_category != ''", 
                                        bookID, section.SectionID, FeedbackTypeMood).
                                Group("mood_category").
                                Order("count DESC").
                                Limit(1)
                        
                        if err := commonMoodQuery.Scan(&moodCategory).Error; err == nil {
                                insight.CommonMood = moodCategory.Category
                        }
                        
                        // Get most common difficulty category
                        var difficultyCategory struct {
                                Category string
                                Count    int
                        }
                        
                        commonDifficultyQuery := r.db.Model(&ContentFeedback{}).
                                Select("difficulty_cat as category, COUNT(*) as count").
                                Where("book_id = ? AND section_id = ? AND type = ? AND difficulty_cat != ''", 
                                        bookID, section.SectionID, FeedbackTypeDifficulty).
                                Group("difficulty_cat").
                                Order("count DESC").
                                Limit(1)
                        
                        if err := commonDifficultyQuery.Scan(&difficultyCategory).Error; err == nil {
                                insight.CommonDifficulty = difficultyCategory.Category
                        }
                        
                        contentInsights = append(contentInsights, insight)
                }
                
                // Convert to map for JSON serialization
                insightMaps := make([]map[string]interface{}, 0, len(contentInsights))
                for _, insight := range contentInsights {
                        insightMap := map[string]interface{}{
                                "sectionID":        insight.SectionID,
                                "chapterID":        insight.ChapterID,
                                "feedbackCount":    insight.FeedbackCount,
                                "averageMood":      insight.AvgMood,
                                "averageDifficulty": insight.AvgDifficulty,
                                "commonMood":       insight.CommonMood,
                                "commonDifficulty": insight.CommonDifficulty,
                        }
                        
                        // Add emoji for common mood
                        if insight.CommonMood != "" {
                                if moodType := MoodType(insight.CommonMood); moodType != "" {
                                        if emoji, exists := MoodEmoji[moodType]; exists {
                                                insightMap["commonMoodEmoji"] = emoji
                                        }
                                }
                        }
                        
                        // Add emoji for common difficulty
                        if insight.CommonDifficulty != "" {
                                if diffType := DifficultyCategory(insight.CommonDifficulty); diffType != "" {
                                        if emoji, exists := DifficultyEmoji[diffType]; exists {
                                                insightMap["commonDifficultyEmoji"] = emoji
                                        }
                                }
                        }
                        
                        // Add suitable mood description based on average mood
                        if insight.AvgMood >= 4.5 {
                                insightMap["moodDescription"] = "Extremely positive"
                        } else if insight.AvgMood >= 3.5 {
                                insightMap["moodDescription"] = "Positive"
                        } else if insight.AvgMood >= 2.5 {
                                insightMap["moodDescription"] = "Neutral"
                        } else if insight.AvgMood >= 1.5 {
                                insightMap["moodDescription"] = "Negative"
                        } else {
                                insightMap["moodDescription"] = "Very negative"
                        }
                        
                        // Add suitable difficulty description
                        if insight.AvgDifficulty >= 4.5 {
                                insightMap["difficultyDescription"] = "Very challenging"
                        } else if insight.AvgDifficulty >= 3.5 {
                                insightMap["difficultyDescription"] = "Challenging"
                        } else if insight.AvgDifficulty >= 2.5 {
                                insightMap["difficultyDescription"] = "Moderate"
                        } else if insight.AvgDifficulty >= 1.5 {
                                insightMap["difficultyDescription"] = "Easy"
                        } else {
                                insightMap["difficultyDescription"] = "Very easy"
                        }
                        
                        insightMaps = append(insightMaps, insightMap)
                }
                
                analysis["contentInsights"] = insightMaps
        }
        
        // 6. Emotional journey through chapters
        var chapterEmotions []struct {
                ChapterID uint
                AvgMood   float64
        }
        
        emotionQuery := r.db.Model(&ContentFeedback{}).
                Select("chapter_id, AVG(value) as avg_mood").
                Where("book_id = ? AND type = ? AND chapter_id != 0", bookID, FeedbackTypeMood).
                Group("chapter_id").
                Order("chapter_id ASC")
        
        if err := emotionQuery.Find(&chapterEmotions).Error; err == nil {
                emotionalJourney := make([]map[string]interface{}, 0, len(chapterEmotions))
                for _, emotion := range chapterEmotions {
                        emotionMap := map[string]interface{}{
                                "chapterID":  emotion.ChapterID,
                                "averageMood": emotion.AvgMood,
                        }
                        
                        // Get dominant mood for this chapter
                        var dominantMood struct {
                                Category string
                                Count    int
                        }
                        
                        dominantMoodQuery := r.db.Model(&ContentFeedback{}).
                                Select("mood_category as category, COUNT(*) as count").
                                Where("book_id = ? AND chapter_id = ? AND type = ? AND mood_category != ''", 
                                        bookID, emotion.ChapterID, FeedbackTypeMood).
                                Group("mood_category").
                                Order("count DESC").
                                Limit(1)
                        
                        if err := dominantMoodQuery.Scan(&dominantMood).Error; err == nil {
                                emotionMap["dominantMood"] = dominantMood.Category
                                
                                // Add the emoji representation for this mood
                                if moodType := MoodType(dominantMood.Category); moodType != "" {
                                        if emoji, exists := MoodEmoji[moodType]; exists {
                                                emotionMap["dominantMoodEmoji"] = emoji
                                        }
                                }
                        }
                        
                        emotionalJourney = append(emotionalJourney, emotionMap)
                }
                
                analysis["emotionalJourney"] = emotionalJourney
        }
        
        // 7. Add emoji reference maps for frontend use
        moodEmojiMap := make(map[string]string)
        for mood, emoji := range MoodEmoji {
                moodEmojiMap[string(mood)] = emoji
        }
        analysis["moodEmojis"] = moodEmojiMap
        
        difficultyEmojiMap := make(map[string]string)
        for difficulty, emoji := range DifficultyEmoji {
                difficultyEmojiMap[string(difficulty)] = emoji
        }
        analysis["difficultyEmojis"] = difficultyEmojiMap
        
        return analysis, nil
}

// DeleteFeedback deletes a specific feedback entry
func (r *GormFeedbackRepository) DeleteFeedback(ctx context.Context, feedbackID, userID uint, feedbackType FeedbackType) error {
        result := r.db.Where("id = ? AND user_id = ? AND type = ?", feedbackID, userID, feedbackType).Delete(&ContentFeedback{})
        if result.Error != nil {
                r.logger.WithError(result.Error).Error("Failed to delete feedback")
                return result.Error
        }
        
        if result.RowsAffected == 0 {
                r.logger.Warn("No feedback found or user not authorized to delete")
                return gorm.ErrRecordNotFound
        }
        
        return nil
}