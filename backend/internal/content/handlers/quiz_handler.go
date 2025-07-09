package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/models"
)

// QuizHandler handles requests for quiz-related operations
type QuizHandler struct {
	// In a real implementation, this would have a quiz service
}

// NewQuizHandler creates a new QuizHandler
func NewQuizHandler() *QuizHandler {
	return &QuizHandler{}
}

// RegisterRoutes registers the quiz-related routes
func (h *QuizHandler) RegisterRoutes(router *gin.RouterGroup) {
	quiz := router.Group("/sections")
	{
		quiz.GET("/:id/quiz", h.GetQuizQuestions)
		quiz.POST("/:id/quiz/submit", h.SubmitQuizAnswers)
	}
}

// GetQuizQuestions handles GET /api/books/sections/:id/quiz
func (h *QuizHandler) GetQuizQuestions(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	// In a real implementation, this would fetch questions from a database
	// For now, we'll return mock data
	questions := []models.QuizQuestion{
		{
			ID:       "1",
			Question: "What is the capital of Nigeria?",
			Options: []string{
				"Lagos",
				"Abuja",
				"Kano",
				"Port Harcourt",
			},
			CorrectOptionIndex: 1,
			Explanation:        "Abuja is the capital city of Nigeria, having replaced Lagos in 1991.",
		},
		{
			ID:       "2",
			Question: "Which river is the longest in Nigeria?",
			Options: []string{
				"Niger River",
				"Benue River",
				"Osun River",
				"Cross River",
			},
			CorrectOptionIndex: 0,
			Explanation:        "The Niger River is the longest river in Nigeria and the third-longest river in Africa.",
		},
		{
			ID:       "3",
			Question: "What year did Nigeria gain independence?",
			Options: []string{
				"1957",
				"1960",
				"1963",
				"1966",
			},
			CorrectOptionIndex: 1,
			Explanation:        "Nigeria gained independence from British colonial rule on October 1, 1960.",
		},
	}

	c.JSON(http.StatusOK, questions)
}

// SubmitQuizAnswers handles POST /api/books/sections/:id/quiz/submit
func (h *QuizHandler) SubmitQuizAnswers(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	var request struct {
		Answers map[string]int `json:"answers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, this would validate answers against stored questions
	// For now, we'll return a mock response
	c.JSON(http.StatusOK, gin.H{
		"score":         2,
		"totalQuestions": 3,
		"correctAnswers": []string{"1", "3"},
		"incorrectAnswers": []string{"2"},
	})
}
