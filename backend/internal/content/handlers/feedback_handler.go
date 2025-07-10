package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/auth"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/repository"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/service"
)

// FeedbackHandler defines the interface for feedback route handling
type FeedbackHandler struct {
	feedbackService service.FeedbackService
	logger          *logrus.Logger
}

// NewFeedbackHandler creates a new feedback handler
func NewFeedbackHandler(feedbackService service.FeedbackService, logger *logrus.Logger) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
		logger:          logger,
	}
}

// RegisterRoutes registers the feedback routes
func (h *FeedbackHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/books/:bookID/feedback/mood", auth.AuthMiddleware(), h.SubmitMoodFeedback)
	router.POST("/books/:bookID/feedback/difficulty", auth.AuthMiddleware(), h.SubmitDifficultyFeedback)
	router.GET("/feedback", auth.AuthMiddleware(), h.GetUserContentFeedback)
	router.GET("/feedback/moods/recent", auth.AuthMiddleware(), h.GetUserRecentMoods)
	router.GET("/content/:contentID/recommendations", auth.AuthMiddleware(), h.GetRecommendedContent)
	router.GET("/books/:bookID/feedback/summary", h.GetContentFeedbackSummary)
	router.GET("/books/:bookID/feedback/analysis", auth.AuthMiddleware(), h.GetDetailedFeedbackAnalysis)
	router.DELETE("/feedback/mood/:id", auth.AuthMiddleware(), h.DeleteMoodFeedback)
	router.DELETE("/feedback/difficulty/:id", auth.AuthMiddleware(), h.DeleteDifficultyFeedback)
	router.GET("/feedback/emoji-maps", h.GetEmojiMaps)
}

// MoodFeedbackRequest represents the request body for mood feedback submission
type MoodFeedbackRequest struct {
	BookID        uint   `json:"bookId" binding:"required"`
	ChapterID     uint   `json:"chapterId"`
	SectionID     uint   `json:"sectionId"`
	Value         int    `json:"value" binding:"required,min=1,max=5"`
	MoodCategory  string `json:"moodCategory"`
	Comment       string `json:"comment"`
	LearningStyle string `json:"learningStyle"`
}

// DifficultyFeedbackRequest represents the request body for difficulty feedback submission
type DifficultyFeedbackRequest struct {
	BookID             uint   `json:"bookId" binding:"required"`
	ChapterID          uint   `json:"chapterId"`
	SectionID          uint   `json:"sectionId"`
	Value              int    `json:"value" binding:"required,min=1,max=5"`
	DifficultyCategory string `json:"difficultyCategory"`
	Comment            string `json:"comment"`
	RecommendNext      bool   `json:"recommendNext"`
}

// FeedbackSummaryQuery represents the query parameters for feedback summary
type FeedbackSummaryQuery struct {
	BookID    uint `form:"bookId"`
	ChapterID uint `form:"chapterId"`
	SectionID uint `form:"sectionId"`
}

// SubmitMoodFeedback handles POST requests to submit mood feedback
func (h *FeedbackHandler) SubmitMoodFeedback(c *gin.Context) {
	var request MoodFeedbackRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Get user ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userID := userClaims.UserID

	// Submit the feedback
	feedback, err := h.feedbackService.SubmitMoodFeedback(
		c.Request.Context(),
		userID,
		request.BookID,
		request.ChapterID,
		request.SectionID,
		request.Value,
		request.MoodCategory,
		request.Comment,
		request.LearningStyle,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to submit mood feedback")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mood feedback submitted successfully",
		"data":    feedback,
	})
}

// SubmitDifficultyFeedback handles POST requests to submit difficulty feedback
func (h *FeedbackHandler) SubmitDifficultyFeedback(c *gin.Context) {
	var request DifficultyFeedbackRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Get user ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userID := userClaims.UserID

	// Submit the feedback
	feedback, err := h.feedbackService.SubmitDifficultyFeedback(
		c.Request.Context(),
		userID,
		request.BookID,
		request.ChapterID,
		request.SectionID,
		request.Value,
		request.DifficultyCategory,
		request.Comment,
		request.RecommendNext,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to submit difficulty feedback")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Difficulty feedback submitted successfully",
		"data":    feedback,
	})
}

// GetUserContentFeedback handles GET requests to retrieve user's feedback
func (h *FeedbackHandler) GetUserContentFeedback(c *gin.Context) {
	// Get user ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userID := userClaims.UserID

	// Get the feedback for the user
	feedback, err := h.feedbackService.GetUserContentFeedback(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user content feedback")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": feedback,
	})
}

// GetUserRecentMoods handles GET requests to retrieve recent mood feedback
func (h *FeedbackHandler) GetUserRecentMoods(c *gin.Context) {
	// Get user ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userID := userClaims.UserID

	// Parse limit from query parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Get recent moods
	moods, err := h.feedbackService.GetUserRecentMoods(c.Request.Context(), userID, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user recent moods")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recent moods: " + err.Error()})
		return
	}

	// Enhance response with emoji information
	type EnhancedMoodFeedback struct {
		ID            uint      `json:"id"`
		UserID        uint      `json:"userId"`
		BookID        uint      `json:"bookId"`
		ChapterID     uint      `json:"chapterId"`
		SectionID     uint      `json:"sectionId"`
		Type          string    `json:"type"`
		Value         int       `json:"value"`
		MoodCategory  string    `json:"moodCategory"`
		MoodEmoji     string    `json:"moodEmoji"`
		Comment       string    `json:"comment"`
		LearningStyle string    `json:"learningStyle"`
		CreatedAt     time.Time `json:"createdAt"`
		UpdatedAt     time.Time `json:"updatedAt"`
	}

	enhancedMoods := make([]EnhancedMoodFeedback, 0, len(moods))
	for _, mood := range moods {
		emoji := ""
		if mood.MoodCategory != "" {
			moodType := repository.MoodType(mood.MoodCategory)
			if emojiVal, exists := repository.MoodEmoji[moodType]; exists {
				emoji = emojiVal
			}
		}

		enhancedMoods = append(enhancedMoods, EnhancedMoodFeedback{
			ID:            mood.ID,
			UserID:        mood.UserID,
			BookID:        mood.BookID,
			ChapterID:     mood.ChapterID,
			SectionID:     mood.SectionID,
			Type:          mood.Type,
			Value:         mood.Value,
			MoodCategory:  mood.MoodCategory,
			MoodEmoji:     emoji,
			Comment:       mood.Comment,
			LearningStyle: mood.LearningStyle,
			CreatedAt:     mood.CreatedAt,
			UpdatedAt:     mood.UpdatedAt,
		})
	}

	// Include emoji reference map for frontend use
	emojiMap := make(map[string]string)
	for moodType, emoji := range repository.MoodEmoji {
		emojiMap[string(moodType)] = emoji
	}

	c.JSON(http.StatusOK, gin.H{
		"data":           enhancedMoods,
		"emojiReference": emojiMap,
	})
}

// GetRecommendedContent handles GET requests to retrieve recommended content
func (h *FeedbackHandler) GetRecommendedContent(c *gin.Context) {
	// Get user ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userID := userClaims.UserID

	// Parse content ID from path parameter
	contentIDStr := c.Param("contentID")
	contentID, err := strconv.ParseUint(contentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	// Get recommended content
	recommendations, err := h.feedbackService.GetRecommendedContent(c.Request.Context(), userID, uint(contentID))
	if err != nil {
		h.logger.WithError(err).Error("Failed to get recommended content")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recommendations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": recommendations,
	})
}

// GetContentFeedbackSummary handles GET requests to retrieve content feedback summary
func (h *FeedbackHandler) GetContentFeedbackSummary(c *gin.Context) {
	var query FeedbackSummaryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	// Get the feedback summary
	summary, err := h.feedbackService.GetContentFeedbackSummary(
		c.Request.Context(),
		query.BookID,
		query.ChapterID,
		query.SectionID,
	)

	if err != nil {
		h.logger.WithError(err).Error("Failed to get content feedback summary")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback summary: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": summary,
	})
}

// DeleteMoodFeedback handles DELETE requests to remove mood feedback
func (h *FeedbackHandler) DeleteMoodFeedback(c *gin.Context) {
	feedbackIDStr := c.Param("id")
	feedbackID, err := strconv.ParseUint(feedbackIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	// Get user ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userID := userClaims.UserID

	// Delete the feedback
	err = h.feedbackService.DeleteMoodFeedback(c.Request.Context(), uint(feedbackID), userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete mood feedback")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mood feedback deleted successfully",
	})
}

// DeleteDifficultyFeedback handles DELETE requests to remove difficulty feedback
func (h *FeedbackHandler) DeleteDifficultyFeedback(c *gin.Context) {
	feedbackIDStr := c.Param("id")
	feedbackID, err := strconv.ParseUint(feedbackIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	// Get user ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userID := userClaims.UserID

	// Delete the feedback
	err = h.feedbackService.DeleteDifficultyFeedback(c.Request.Context(), uint(feedbackID), userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete difficulty feedback")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Difficulty feedback deleted successfully",
	})
}

// GetDetailedFeedbackAnalysis handles GET requests to retrieve comprehensive feedback analysis
func (h *FeedbackHandler) GetDetailedFeedbackAnalysis(c *gin.Context) {
	// Parse book ID from path parameter
	bookIDStr := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Get the detailed analysis
	analysis, err := h.feedbackService.GetDetailedFeedbackAnalysis(c.Request.Context(), uint(bookID))
	if err != nil {
		h.logger.WithError(err).Error("Failed to get detailed feedback analysis")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve detailed feedback analysis: " + err.Error()})
		return
	}

	// Add visualization helper information
	visualizationHelp := map[string]interface{}{
		"moodValues": map[string]string{
			"1": "Very Negative",
			"2": "Negative",
			"3": "Neutral",
			"4": "Positive",
			"5": "Very Positive",
		},
		"difficultyValues": map[string]string{
			"1": "Very Easy",
			"2": "Easy",
			"3": "Moderate",
			"4": "Challenging",
			"5": "Very Challenging",
		},
		"displayGuide": "Use the provided emoji maps to display appropriate graphics for each mood and difficulty category.",
	}

	c.JSON(http.StatusOK, gin.H{
		"data":          analysis,
		"visualization": visualizationHelp,
	})
}

// GetEmojiMaps provides emoji mapping references for frontend usage
func (h *FeedbackHandler) GetEmojiMaps(c *gin.Context) {
	// Create mood emoji mapping
	moodEmojiMap := make(map[string]string)
	for moodType, emoji := range repository.MoodEmoji {
		moodEmojiMap[string(moodType)] = emoji
	}

	// Create difficulty emoji mapping
	difficultyEmojiMap := make(map[string]string)
	for diffType, emoji := range repository.DifficultyEmoji {
		difficultyEmojiMap[string(diffType)] = emoji
	}

	// Create mood value descriptions
	moodValueDescriptions := map[string]string{
		"1": "Very Negative",
		"2": "Negative",
		"3": "Neutral",
		"4": "Positive",
		"5": "Very Positive",
	}

	// Create difficulty value descriptions
	difficultyValueDescriptions := map[string]string{
		"1": "Very Easy",
		"2": "Easy",
		"3": "Moderate",
		"4": "Challenging",
		"5": "Very Challenging",
	}

	// Create default mood mapping based on numeric values
	defaultMoodCategoryMapping := map[string]string{
		"1": string(repository.MoodTypeFrustrated),
		"2": string(repository.MoodTypeConfused),
		"3": string(repository.MoodTypeInspired),
		"4": string(repository.MoodTypeHappy),
		"5": string(repository.MoodTypeMotivated),
	}

	// Create default difficulty mapping based on numeric values
	defaultDifficultyCategoryMapping := map[string]string{
		"1": string(repository.DifficultyCategoryEasy),
		"2": string(repository.DifficultyCategoryEasy),
		"3": string(repository.DifficultyCategoryModerate),
		"4": string(repository.DifficultyCategoryHard),
		"5": string(repository.DifficultyCategoryComplex),
	}

	// Learning style descriptions
	learningStyleDescriptions := map[string]string{
		string(repository.LearningStyleVisual):      "Learns better with visuals and diagrams",
		string(repository.LearningStyleAuditory):    "Learns better with spoken explanations",
		string(repository.LearningStyleKinesthetic): "Learns better by doing and practicing",
		string(repository.LearningStyleReading):     "Learns better through reading text",
	}

	c.JSON(http.StatusOK, gin.H{
		"moodEmojis":                  moodEmojiMap,
		"difficultyEmojis":            difficultyEmojiMap,
		"moodValueDescriptions":       moodValueDescriptions,
		"difficultyValueDescriptions": difficultyValueDescriptions,
		"defaultMoodCategories":       defaultMoodCategoryMapping,
		"defaultDifficultyCategories": defaultDifficultyCategoryMapping,
		"learningStyleDescriptions":   learningStyleDescriptions,
	})
}
