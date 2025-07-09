package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/logger"
)

// LikeHandler handles like-related requests
type LikeHandler struct {
	likeService LikeService
	logger      *logger.Logger
}

// LikeService defines the interface for like service operations
type LikeService interface {
	LikeDiscussion(userID, discussionID uint) error
	UnlikeDiscussion(userID, discussionID uint) error
	LikeComment(userID, commentID uint) error
	UnlikeComment(userID, commentID uint) error
	GetUserLikes(userID uint, entityType string) ([]uint, error)
}

// NewLikeHandler creates a new like handler
func NewLikeHandler(likeService LikeService, logger *logger.Logger) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
		logger:      logger,
	}
}

// LikeDiscussion handles liking a discussion
func (h *LikeHandler) LikeDiscussion(c *gin.Context) {
	discussionIDParam := c.Param("id")
	discussionID, err := strconv.ParseUint(discussionIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid discussion ID"))
		return
	}

	userID, _ := c.Get("user_id")

	err = h.likeService.LikeDiscussion(userID.(uint), uint(discussionID))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to like discussion")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to like discussion"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discussion liked successfully"})
}

// UnlikeDiscussion handles unliking a discussion
func (h *LikeHandler) UnlikeDiscussion(c *gin.Context) {
	discussionIDParam := c.Param("id")
	discussionID, err := strconv.ParseUint(discussionIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid discussion ID"))
		return
	}

	userID, _ := c.Get("user_id")

	err = h.likeService.UnlikeDiscussion(userID.(uint), uint(discussionID))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to unlike discussion")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to unlike discussion"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discussion unliked successfully"})
}

// LikeComment handles liking a comment
func (h *LikeHandler) LikeComment(c *gin.Context) {
	commentIDParam := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid comment ID"))
		return
	}

	userID, _ := c.Get("user_id")

	err = h.likeService.LikeComment(userID.(uint), uint(commentID))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to like comment")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to like comment"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment liked successfully"})
}

// UnlikeComment handles unliking a comment
func (h *LikeHandler) UnlikeComment(c *gin.Context) {
	commentIDParam := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid comment ID"))
		return
	}

	userID, _ := c.Get("user_id")

	err = h.likeService.UnlikeComment(userID.(uint), uint(commentID))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to unlike comment")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to unlike comment"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment unliked successfully"})
}

// GetUserLikes handles getting all likes for a user
func (h *LikeHandler) GetUserLikes(c *gin.Context) {
	entityType := c.DefaultQuery("entity_type", "discussion")
	if entityType != "discussion" && entityType != "comment" {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid entity type. Must be 'discussion' or 'comment'"))
		return
	}

	userID, _ := c.Get("user_id")

	likes, err := h.likeService.GetUserLikes(userID.(uint), entityType)
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to get user likes")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get user likes"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entity_type": entityType,
		"liked_ids":   likes,
	})
}
