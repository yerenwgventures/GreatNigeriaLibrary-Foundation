package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
)

// ProgressHandler handles progress-related requests
type ProgressHandler struct {
	progressService ProgressService
	logger          *logger.Logger
}

// ProgressService defines the interface for progress service operations
type ProgressService interface {
	UpdateProgress(userID, bookID, chapterID, sectionID uint, lastPosition int, percentComplete float64) (*models.UserProgressResponse, error)
	GetProgress(userID, bookID uint) (*models.UserProgressResponse, error)
}

// NewProgressHandler creates a new progress handler
func NewProgressHandler(progressService ProgressService, logger *logger.Logger) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		logger:          logger,
	}
}

// UpdateProgress handles updating a user's reading progress
func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
	bookIDParam := c.Param("id")
	bookID, err := strconv.ParseUint(bookIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid book ID"))
		return
	}

	var req struct {
		ChapterID       uint    `json:"chapter_id" binding:"required"`
		SectionID       uint    `json:"section_id" binding:"required"`
		LastPosition    int     `json:"last_position"`
		PercentComplete float64 `json:"percent_complete"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")

	progress, err := h.progressService.UpdateProgress(
		userID.(uint),
		uint(bookID),
		req.ChapterID,
		req.SectionID,
		req.LastPosition,
		req.PercentComplete,
	)
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to update progress")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to update progress"))
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetProgress handles getting a user's reading progress for a book
func (h *ProgressHandler) GetProgress(c *gin.Context) {
	bookIDParam := c.Param("id")
	bookID, err := strconv.ParseUint(bookIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid book ID"))
		return
	}

	userID, _ := c.Get("user_id")

	progress, err := h.progressService.GetProgress(userID.(uint), uint(bookID))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to get progress")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get progress"))
		return
	}

	if progress == nil {
		// Return empty progress object if no progress exists
		c.JSON(http.StatusOK, &models.UserProgressResponse{
			BookID:          uint(bookID),
			PercentComplete: 0,
		})
		return
	}

	c.JSON(http.StatusOK, progress)
}
