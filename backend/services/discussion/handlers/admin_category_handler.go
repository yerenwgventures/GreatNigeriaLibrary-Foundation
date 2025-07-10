package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/service"
)

// AdminCategoryHandler handles admin-specific API requests for forum categories
type AdminCategoryHandler struct {
	discussionService service.DiscussionService
}

// NewAdminCategoryHandler creates a new admin category handler
func NewAdminCategoryHandler(discussionService service.DiscussionService) *AdminCategoryHandler {
	return &AdminCategoryHandler{
		discussionService: discussionService,
	}
}

// RegisterRoutes registers all admin category routes
func (h *AdminCategoryHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	// Admin routes for category configuration
	adminCategories := router.Group("/api/v1/admin/discussions/categories")
	adminCategories.Use(authMiddleware, adminMiddleware)
	{
		// Category configuration
		adminCategories.GET("/:id/config", h.GetCategoryConfig)
		adminCategories.PUT("/:id/config", h.UpdateCategoryConfig)

		// Posting rules
		adminCategories.GET("/:id/posting-rules", h.GetPostingRules)
		adminCategories.PUT("/:id/posting-rules", h.UpdatePostingRules)

		// Auto moderation settings
		adminCategories.GET("/:id/moderation", h.GetAutoModerationSettings)
		adminCategories.PUT("/:id/moderation", h.UpdateAutoModerationSettings)

		// Category moderators
		adminCategories.GET("/:id/moderators", h.GetCategoryModerators)
		adminCategories.POST("/:id/moderators", h.AddCategoryModerator)
		adminCategories.PUT("/:id/moderators/:userId", h.UpdateCategoryModerator)
		adminCategories.DELETE("/:id/moderators/:userId", h.RemoveCategoryModerator)
	}
}

// GetCategoryConfig retrieves configuration for a category
func (h *AdminCategoryHandler) GetCategoryConfig(c *gin.Context) {
	// Parse category ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get category config
	config, err := h.discussionService.GetCategoryConfig(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateCategoryConfig handles requests to update category configuration
func (h *AdminCategoryHandler) UpdateCategoryConfig(c *gin.Context) {
	// Parse category ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Parse request body
	var config models.CategoryConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update config
	updatedConfig, err := h.discussionService.CreateOrUpdateCategoryConfig(uint(categoryID), config)
	if err != nil {
		if err == models.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedConfig)
}

// GetPostingRules retrieves posting rules for a category
func (h *AdminCategoryHandler) GetPostingRules(c *gin.Context) {
	// Parse category ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get posting rules
	rules, err := h.discussionService.GetPostingRules(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

// UpdatePostingRules handles requests to update posting rules
func (h *AdminCategoryHandler) UpdatePostingRules(c *gin.Context) {
	// Parse category ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Parse request body
	var rules models.PostingRules
	if err := c.ShouldBindJSON(&rules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update rules
	updatedRules, err := h.discussionService.CreateOrUpdatePostingRules(uint(categoryID), rules)
	if err != nil {
		if err == models.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedRules)
}

// GetAutoModerationSettings retrieves auto-moderation settings for a category
func (h *AdminCategoryHandler) GetAutoModerationSettings(c *gin.Context) {
	// Parse category ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get auto-moderation settings
	settings, err := h.discussionService.GetAutoModerationSettings(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateAutoModerationSettings handles requests to update auto-moderation settings
func (h *AdminCategoryHandler) UpdateAutoModerationSettings(c *gin.Context) {
	// Parse category ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Parse request body
	var settings models.AutoModerationSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update settings
	updatedSettings, err := h.discussionService.CreateOrUpdateAutoModerationSettings(uint(categoryID), settings)
	if err != nil {
		if err == models.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedSettings)
}

// GetCategoryModerators retrieves moderators for a category
func (h *AdminCategoryHandler) GetCategoryModerators(c *gin.Context) {
	// Parse category ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get category moderators
	moderators, err := h.discussionService.GetCategoryModerators(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, moderators)
}

// AddCategoryModerator handles requests to add a moderator to a category
func (h *AdminCategoryHandler) AddCategoryModerator(c *gin.Context) {
	// Parse category ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Parse request body
	var request struct {
		UserID            uint `json:"userId" binding:"required"`
		CanPinPosts       bool `json:"canPinPosts"`
		CanLockPosts      bool `json:"canLockPosts"`
		CanMovePosts      bool `json:"canMovePosts"`
		CanDeletePosts    bool `json:"canDeletePosts"`
		CanModerateUsers  bool `json:"canModerateUsers"`
		CanEditPosts      bool `json:"canEditPosts"`
		CanApproveContent bool `json:"canApproveContent"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create moderator with permissions
	moderator := models.CategoryModerator{
		UserID:            request.UserID,
		CanPinPosts:       request.CanPinPosts,
		CanLockPosts:      request.CanLockPosts,
		CanMovePosts:      request.CanMovePosts,
		CanDeletePosts:    request.CanDeletePosts,
		CanModerateUsers:  request.CanModerateUsers,
		CanEditPosts:      request.CanEditPosts,
		CanApproveContent: request.CanApproveContent,
	}

	// Add moderator
	addedModerator, err := h.discussionService.AddCategoryModerator(uint(categoryID), request.UserID, moderator)
	if err != nil {
		if err == models.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, addedModerator)
}

// UpdateCategoryModerator handles requests to update a moderator's permissions
func (h *AdminCategoryHandler) UpdateCategoryModerator(c *gin.Context) {
	// Parse category ID and user ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse request body
	var request struct {
		CanPinPosts       bool `json:"canPinPosts"`
		CanLockPosts      bool `json:"canLockPosts"`
		CanMovePosts      bool `json:"canMovePosts"`
		CanDeletePosts    bool `json:"canDeletePosts"`
		CanModerateUsers  bool `json:"canModerateUsers"`
		CanEditPosts      bool `json:"canEditPosts"`
		CanApproveContent bool `json:"canApproveContent"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create moderator with updated permissions
	moderator := models.CategoryModerator{
		CanPinPosts:       request.CanPinPosts,
		CanLockPosts:      request.CanLockPosts,
		CanMovePosts:      request.CanMovePosts,
		CanDeletePosts:    request.CanDeletePosts,
		CanModerateUsers:  request.CanModerateUsers,
		CanEditPosts:      request.CanEditPosts,
		CanApproveContent: request.CanApproveContent,
	}

	// Update moderator
	updatedModerator, err := h.discussionService.UpdateCategoryModerator(uint(categoryID), uint(userID), moderator)
	if err != nil {
		if err == models.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedModerator)
}

// RemoveCategoryModerator handles requests to remove a moderator from a category
func (h *AdminCategoryHandler) RemoveCategoryModerator(c *gin.Context) {
	// Parse category ID and user ID
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Remove moderator
	err = h.discussionService.RemoveCategoryModerator(uint(categoryID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Moderator removed successfully"})
}
