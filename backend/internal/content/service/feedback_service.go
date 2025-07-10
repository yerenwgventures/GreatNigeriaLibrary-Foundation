package service

import (
        "context"
        "errors"
        "fmt"

        "github.com/sirupsen/logrus"
        "gorm.io/gorm"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/repository"
)

// FeedbackService defines the interface for feedback-related operations
type FeedbackService interface {
        SubmitMoodFeedback(ctx context.Context, userID, bookID, chapterID, sectionID uint, value int, moodCategory string, comment string, learningStyle string) (*repository.ContentFeedback, error)
        SubmitDifficultyFeedback(ctx context.Context, userID, bookID, chapterID, sectionID uint, value int, difficultyCategory string, comment string, recommendNext bool) (*repository.ContentFeedback, error)
        GetUserContentFeedback(ctx context.Context, userID uint) ([]repository.ContentFeedback, error)
        GetUserRecentMoods(ctx context.Context, userID uint, limit int) ([]repository.ContentFeedback, error)
        GetRecommendedContent(ctx context.Context, userID uint, baseContentID uint) ([]uint, error)
        GetContentFeedbackSummary(ctx context.Context, bookID, chapterID, sectionID uint) (map[string]interface{}, error)
        GetDetailedFeedbackAnalysis(ctx context.Context, bookID uint) (map[string]interface{}, error)
        DeleteMoodFeedback(ctx context.Context, feedbackID, userID uint) error
        DeleteDifficultyFeedback(ctx context.Context, feedbackID, userID uint) error
}

// DefaultFeedbackService implements the FeedbackService interface
type DefaultFeedbackService struct {
        feedbackRepo repository.FeedbackRepository
        bookRepo     repository.BookRepository
        logger       *logrus.Logger
}

// NewFeedbackService creates a new feedback service instance
func NewFeedbackService(feedbackRepo repository.FeedbackRepository, bookRepo repository.BookRepository, logger *logrus.Logger) FeedbackService {
        return &DefaultFeedbackService{
                feedbackRepo: feedbackRepo,
                bookRepo:     bookRepo,
                logger:       logger,
        }
}

// SubmitMoodFeedback handles the submission of mood feedback
func (s *DefaultFeedbackService) SubmitMoodFeedback(ctx context.Context, userID, bookID, chapterID, sectionID uint, value int, moodCategory string, comment string, learningStyle string) (*repository.ContentFeedback, error) {
        // Validate input parameters
        if userID == 0 {
                return nil, errors.New("user ID is required")
        }
        if bookID == 0 {
                return nil, errors.New("book ID is required")
        }
        if value < 1 || value > 5 {
                return nil, errors.New("mood value must be between 1 and 5")
        }

        // Validate mood category if provided
        if moodCategory != "" {
                validCategories := map[string]bool{
                        string(repository.MoodTypeHappy):      true,
                        string(repository.MoodTypeInspired):   true,
                        string(repository.MoodTypeConfused):   true,
                        string(repository.MoodTypeFrustrated): true,
                        string(repository.MoodTypeMotivated):  true,
                }
                if !validCategories[moodCategory] {
                        return nil, errors.New("invalid mood category")
                }
        } else {
                // Assign default mood category based on value
                switch {
                case value == 5:
                        moodCategory = string(repository.MoodTypeMotivated)
                case value == 4:
                        moodCategory = string(repository.MoodTypeHappy)
                case value == 3:
                        moodCategory = string(repository.MoodTypeInspired)
                case value == 2:
                        moodCategory = string(repository.MoodTypeConfused)
                case value == 1:
                        moodCategory = string(repository.MoodTypeFrustrated)
                }
        }
        
        // Validate learning style if provided
        if learningStyle != "" {
                validStyles := map[string]bool{
                        string(repository.LearningStyleVisual):      true,
                        string(repository.LearningStyleAuditory):    true,
                        string(repository.LearningStyleKinesthetic): true,
                        string(repository.LearningStyleReading):     true,
                }
                if !validStyles[learningStyle] {
                        return nil, errors.New("invalid learning style")
                }
        }

        // Validate the content exists
        exists, err := s.validateContent(bookID, chapterID, sectionID)
        if err != nil {
                return nil, err
        }
        if !exists {
                return nil, errors.New("specified content does not exist")
        }

        // Submit the feedback
        feedback, err := s.feedbackRepo.SubmitMoodFeedback(ctx, userID, bookID, chapterID, sectionID, value, moodCategory, comment, learningStyle)
        if err != nil {
                s.logger.WithError(err).Error("Failed to submit mood feedback")
                return nil, err
        }

        return feedback, nil
}

// SubmitDifficultyFeedback handles the submission of difficulty feedback
func (s *DefaultFeedbackService) SubmitDifficultyFeedback(ctx context.Context, userID, bookID, chapterID, sectionID uint, value int, difficultyCategory string, comment string, recommendNext bool) (*repository.ContentFeedback, error) {
        // Validate input parameters
        if userID == 0 {
                return nil, errors.New("user ID is required")
        }
        if bookID == 0 {
                return nil, errors.New("book ID is required")
        }
        if value < 1 || value > 5 {
                return nil, errors.New("difficulty value must be between 1 and 5")
        }

        // Validate difficulty category if provided
        if difficultyCategory != "" {
                validCategories := map[string]bool{
                        string(repository.DifficultyCategoryEasy):     true,
                        string(repository.DifficultyCategoryModerate): true,
                        string(repository.DifficultyCategoryHard):     true,
                        string(repository.DifficultyCategoryComplex):  true,
                }
                if !validCategories[difficultyCategory] {
                        return nil, errors.New("invalid difficulty category")
                }
        } else {
                // Assign default difficulty category based on value
                switch {
                case value == 1:
                        difficultyCategory = string(repository.DifficultyCategoryEasy)
                case value == 2:
                        difficultyCategory = string(repository.DifficultyCategoryEasy)
                case value == 3:
                        difficultyCategory = string(repository.DifficultyCategoryModerate)
                case value == 4:
                        difficultyCategory = string(repository.DifficultyCategoryHard)
                case value == 5:
                        difficultyCategory = string(repository.DifficultyCategoryComplex)
                }
        }

        // Validate the content exists
        exists, err := s.validateContent(bookID, chapterID, sectionID)
        if err != nil {
                return nil, err
        }
        if !exists {
                return nil, errors.New("specified content does not exist")
        }

        // Submit the feedback
        feedback, err := s.feedbackRepo.SubmitDifficultyFeedback(ctx, userID, bookID, chapterID, sectionID, value, difficultyCategory, comment, recommendNext)
        if err != nil {
                s.logger.WithError(err).Error("Failed to submit difficulty feedback")
                return nil, err
        }

        // If user requests recommendations based on this feedback
        if recommendNext {
                // Queue a background task to generate recommendations
                // This is not implemented in the current scope but could be added later
                s.logger.Info("Recommendation requested for content ID: %d", sectionID)
        }

        return feedback, nil
}

// GetUserContentFeedback retrieves all feedback submitted by a user
func (s *DefaultFeedbackService) GetUserContentFeedback(ctx context.Context, userID uint) ([]repository.ContentFeedback, error) {
        if userID == 0 {
                return nil, errors.New("user ID is required")
        }

        feedbacks, err := s.feedbackRepo.GetUserContentFeedback(ctx, userID)
        if err != nil {
                s.logger.WithError(err).Error("Failed to retrieve user feedback")
                return nil, err
        }

        return feedbacks, nil
}

// GetUserRecentMoods retrieves the most recent mood feedback from a user
func (s *DefaultFeedbackService) GetUserRecentMoods(ctx context.Context, userID uint, limit int) ([]repository.ContentFeedback, error) {
        if userID == 0 {
                return nil, errors.New("user ID is required")
        }

        moods, err := s.feedbackRepo.GetUserRecentMoods(ctx, userID, limit)
        if err != nil {
                s.logger.WithError(err).Error("Failed to retrieve user recent moods")
                return nil, err
        }

        return moods, nil
}

// GetRecommendedContent provides content recommendations based on user feedback
func (s *DefaultFeedbackService) GetRecommendedContent(ctx context.Context, userID uint, baseContentID uint) ([]uint, error) {
        if userID == 0 {
                return nil, errors.New("user ID is required")
        }

        recommendedContentIDs, err := s.feedbackRepo.GetRecommendedContent(ctx, userID, baseContentID)
        if err != nil {
                s.logger.WithError(err).Error("Failed to retrieve recommended content")
                return nil, err
        }

        return recommendedContentIDs, nil
}

// GetContentFeedbackSummary retrieves aggregated feedback statistics for specific content
func (s *DefaultFeedbackService) GetContentFeedbackSummary(ctx context.Context, bookID, chapterID, sectionID uint) (map[string]interface{}, error) {
        // If content is specified, validate it exists
        if bookID > 0 {
                exists, err := s.validateContent(bookID, chapterID, sectionID)
                if err != nil {
                        return nil, err
                }
                if !exists {
                        return nil, errors.New("specified content does not exist")
                }
        }

        summary, err := s.feedbackRepo.GetContentFeedbackSummary(ctx, bookID, chapterID, sectionID)
        if err != nil {
                s.logger.WithError(err).Error("Failed to retrieve feedback summary")
                return nil, err
        }

        return summary, nil
}

// GetDetailedFeedbackAnalysis provides a comprehensive analysis of feedback for a book
func (s *DefaultFeedbackService) GetDetailedFeedbackAnalysis(ctx context.Context, bookID uint) (map[string]interface{}, error) {
        if bookID == 0 {
                return nil, errors.New("book ID is required")
        }

        // Validate the book exists
        exists, err := s.validateContent(bookID, 0, 0)
        if err != nil {
                return nil, err
        }
        if !exists {
                return nil, errors.New("specified book does not exist")
        }

        // Get detailed analysis
        analysis, err := s.feedbackRepo.GetDetailedFeedbackAnalysis(ctx, bookID)
        if err != nil {
                s.logger.WithError(err).Error("Failed to retrieve detailed feedback analysis")
                return nil, err
        }

        // Enhance the analysis with book metadata
        book, err := s.bookRepo.GetBookByID(bookID)
        if err == nil && book != nil {
                analysis["bookTitle"] = book.Title
                analysis["bookAuthor"] = book.Author
                analysis["bookDescription"] = book.Description
        }

        // Add adaptive learning suggestions based on the analysis
        analysis["adaptiveLearning"] = s.generateAdaptiveLearningRecommendations(analysis)

        return analysis, nil
}

// generateAdaptiveLearningRecommendations creates customized learning recommendations
// based on feedback analysis
func (s *DefaultFeedbackService) generateAdaptiveLearningRecommendations(analysis map[string]interface{}) map[string]interface{} {
        recommendations := map[string]interface{}{
                "contentAdjustments": []string{},
                "learningPathways":    []string{},
                "supplementaryMaterials": []string{},
        }

        // Extract data from analysis
        moodAnalysis, moodOk := analysis["moodAnalysis"].(map[string]interface{})
        difficultyAnalysis, diffOk := analysis["difficultyAnalysis"].(map[string]interface{})
        learningStyles, stylesOk := analysis["learningStyles"].(map[string]interface{})
        
        contentAdjustments := make([]string, 0)
        learningPathways := make([]string, 0)
        supplementaryMaterials := make([]string, 0)

        // Generate content adjustment recommendations based on difficulty analysis
        if diffOk {
                if distribution, ok := difficultyAnalysis["categoryDistribution"].(map[string]int); ok {
                        // Check if there's a lot of "hard" or "complex" feedback
                        hardCount := distribution[string(repository.DifficultyCategoryHard)]
                        complexCount := distribution[string(repository.DifficultyCategoryComplex)]
                        totalCount := 0
                        
                        for _, count := range distribution {
                                totalCount += count
                        }
                        
                        if totalCount > 0 {
                                hardPercentage := float64(hardCount+complexCount) / float64(totalCount)
                                
                                if hardPercentage >= 0.6 {
                                        contentAdjustments = append(contentAdjustments, 
                                                "Consider simplifying complex sections or providing more explanatory examples")
                                        contentAdjustments = append(contentAdjustments, 
                                                "Add more progressive learning paths that build up to difficult concepts")
                                        supplementaryMaterials = append(supplementaryMaterials,
                                                "Provide additional reference materials for challenging topics")
                                }
                        }
                }
        }

        // Generate learning pathways based on mood analysis
        if moodOk {
                if distribution, ok := moodAnalysis["categoryDistribution"].(map[string]int); ok {
                        // Check if there's a lot of "confused" or "frustrated" feedback
                        confusedCount := distribution[string(repository.MoodTypeConfused)]
                        frustratedCount := distribution[string(repository.MoodTypeFrustrated)]
                        totalCount := 0
                        
                        for _, count := range distribution {
                                totalCount += count
                        }
                        
                        if totalCount > 0 {
                                negativePercentage := float64(confusedCount+frustratedCount) / float64(totalCount)
                                
                                if negativePercentage >= 0.4 {
                                        learningPathways = append(learningPathways, 
                                                "Implement a more structured step-by-step approach to difficult concepts")
                                        contentAdjustments = append(contentAdjustments,
                                                "Consider revising sections with high confusion rates")
                                        supplementaryMaterials = append(supplementaryMaterials,
                                                "Create visual guides and diagrams for complex information")
                                }
                        }
                }
        }

        // Generate learning style-based recommendations
        if stylesOk {
                if distribution, ok := learningStyles["distribution"].(map[string]int); ok {
                        // Find the dominant learning style
                        var dominantStyle string
                        maxCount := 0
                        
                        for style, count := range distribution {
                                if count > maxCount {
                                        maxCount = count
                                        dominantStyle = style
                                }
                        }
                        
                        if dominantStyle != "" {
                                switch dominantStyle {
                                case "visual":
                                        learningPathways = append(learningPathways,
                                                "Emphasize visual learning with more diagrams, charts, and images")
                                        supplementaryMaterials = append(supplementaryMaterials,
                                                "Provide video tutorials and visual guides as supplementary content")
                                case "auditory":
                                        learningPathways = append(learningPathways,
                                                "Enhance content with audio explanations and discussions")
                                        supplementaryMaterials = append(supplementaryMaterials,
                                                "Create podcast-style content for key concepts")
                                case "kinesthetic":
                                        learningPathways = append(learningPathways,
                                                "Incorporate more interactive exercises and practical activities")
                                        supplementaryMaterials = append(supplementaryMaterials,
                                                "Develop hands-on projects and simulations")
                                case "reading":
                                        learningPathways = append(learningPathways,
                                                "Structure content with clear textual explanations and references")
                                        supplementaryMaterials = append(supplementaryMaterials,
                                                "Provide additional reading materials and case studies")
                                }
                        }
                }
        }

        // Add default recommendations if none were generated
        if len(contentAdjustments) == 0 {
                contentAdjustments = append(contentAdjustments, 
                        "Continue monitoring feedback to identify potential content improvements")
        }
        
        if len(learningPathways) == 0 {
                learningPathways = append(learningPathways,
                        "Consider offering multiple learning paths to accommodate different learning preferences")
        }
        
        if len(supplementaryMaterials) == 0 {
                supplementaryMaterials = append(supplementaryMaterials,
                        "Provide diverse supplementary materials to support various learning styles")
        }

        recommendations["contentAdjustments"] = contentAdjustments
        recommendations["learningPathways"] = learningPathways
        recommendations["supplementaryMaterials"] = supplementaryMaterials

        return recommendations
}

// DeleteMoodFeedback deletes a mood feedback entry
func (s *DefaultFeedbackService) DeleteMoodFeedback(ctx context.Context, feedbackID, userID uint) error {
        if feedbackID == 0 {
                return errors.New("feedback ID is required")
        }
        if userID == 0 {
                return errors.New("user ID is required")
        }

        err := s.feedbackRepo.DeleteFeedback(ctx, feedbackID, userID, repository.FeedbackTypeMood)
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return errors.New("feedback not found or you are not authorized to delete it")
                }
                s.logger.WithError(err).Error("Failed to delete mood feedback")
                return err
        }

        return nil
}

// DeleteDifficultyFeedback deletes a difficulty feedback entry
func (s *DefaultFeedbackService) DeleteDifficultyFeedback(ctx context.Context, feedbackID, userID uint) error {
        if feedbackID == 0 {
                return errors.New("feedback ID is required")
        }
        if userID == 0 {
                return errors.New("user ID is required")
        }

        err := s.feedbackRepo.DeleteFeedback(ctx, feedbackID, userID, repository.FeedbackTypeDifficulty)
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return errors.New("feedback not found or you are not authorized to delete it")
                }
                s.logger.WithError(err).Error("Failed to delete difficulty feedback")
                return err
        }

        return nil
}

// validateContent checks if the specified content exists
func (s *DefaultFeedbackService) validateContent(bookID, chapterID, sectionID uint) (bool, error) {
        // Check if book exists
        if bookID > 0 {
                book, err := s.bookRepo.GetBookByID(bookID)
                if err != nil {
                        if errors.Is(err, gorm.ErrRecordNotFound) {
                                return false, nil
                        }
                        return false, fmt.Errorf("error checking book existence: %w", err)
                }
                if book == nil {
                        return false, nil
                }

                // If chapter ID is provided, check if it exists in the book
                if chapterID > 0 {
                        chapter, err := s.bookRepo.GetChapterByID(chapterID)
                        if err != nil {
                                if errors.Is(err, gorm.ErrRecordNotFound) {
                                        return false, nil
                                }
                                return false, fmt.Errorf("error checking chapter existence: %w", err)
                        }
                        if chapter == nil {
                                return false, nil
                        }

                        // If section ID is provided, check if it exists in the chapter
                        if sectionID > 0 {
                                section, err := s.bookRepo.GetSectionByID(sectionID)
                                if err != nil {
                                        if errors.Is(err, gorm.ErrRecordNotFound) {
                                                return false, nil
                                        }
                                        return false, fmt.Errorf("error checking section existence: %w", err)
                                }
                                if section == nil {
                                        return false, nil
                                }
                        }
                }
        }

        return true, nil
}