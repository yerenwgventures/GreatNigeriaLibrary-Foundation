package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
)

// CommentHandler handles comment-related requests
type CommentHandler struct {
	commentService CommentService
	logger         *logger.Logger
}

// CommentService defines the interface for comment service operations
type CommentService interface {
	ListComments(discussionID uint, page, limit int) ([]models.CommentResponse, int, error)
	CreateComment(userID, discussionID uint, req *models.CommentCreateRequest) (*models.CommentResponse, error)
	UpdateComment(id, userID uint, req *models.CommentUpdateRequest) (*models.CommentResponse, error)
	DeleteComment(id, userID uint) error
	GetCommentReplies(commentID uint) ([]models.CommentResponse, error)
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(commentService CommentService, logger *logger.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		logger:         logger,
	}
}

// ListComments handles listing comments for a discussion
func (h *CommentHandler) ListComments(c *gin.Context) {
	discussionIDParam := c.Param("id")
	discussionID, err := strconv.ParseUint(discussionIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid discussion ID"))
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	// Ensure valid pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	comments, total, err := h.commentService.ListComments(uint(discussionID), page, limit)
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to list comments")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to list comments"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// CreateComment handles creating a new comment
func (h *CommentHandler) CreateComment(c *gin.Context) {
	discussionIDParam := c.Param("id")
	discussionID, err := strconv.ParseUint(discussionIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid discussion ID"))
		return
	}

	var req models.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")

	comment, err := h.commentService.CreateComment(userID.(uint), uint(discussionID), &req)
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to create comment")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to create comment"))
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// UpdateComment handles updating a comment
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid comment ID"))
		return
	}

	var req models.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")

	comment, err := h.commentService.UpdateComment(uint(id), userID.(uint), &req)
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to update comment")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to update comment"))
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment handles deleting a comment
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid comment ID"))
		return
	}

	userID, _ := c.Get("user_id")

	err = h.commentService.DeleteComment(uint(id), userID.(uint))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to delete comment")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to delete comment"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// GetCommentReplies handles getting replies to a comment
func (h *CommentHandler) GetCommentReplies(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid comment ID"))
		return
	}

	replies, err := h.commentService.GetCommentReplies(uint(id))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to get comment replies")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get comment replies"))
		return
	}

	c.JSON(http.StatusOK, replies)
}
