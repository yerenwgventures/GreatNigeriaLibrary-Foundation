package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/service"
)

// FlagHandler defines the handler for content flag endpoints
type FlagHandler struct {
	flagService service.FlagService
}

// NewFlagHandler creates a new flag handler
func NewFlagHandler(flagService service.FlagService) *FlagHandler {
	return &FlagHandler{
		flagService: flagService,
	}
}

// RegisterRoutes registers the routes for content flags
func (h *FlagHandler) RegisterRoutes(router *gin.RouterGroup) {
	moderation := router.Group("/moderation")
	{
		// Content flagging
		moderation.POST("/flags", h.FlagContent)
		moderation.GET("/flags/content/:type/:id", h.GetFlagsByContent)
		moderation.GET("/flags", h.GetFlagsByStatus)
		moderation.GET("/flags/:id", h.GetFlagByID)
		moderation.PUT("/flags/:id/review", h.ReviewFlag)
		moderation.PUT("/flags/:id/assign", h.AssignFlag)
		
		// Content moderation
		moderation.POST("/status", h.CreateOrUpdateModerationStatus)
		moderation.GET("/status/:type/:id", h.GetModerationStatusByContent)
		moderation.PUT("/status/:type/:id/notify", h.MarkUserNotified)
		moderation.GET("/status/pending/count", h.GetPendingModerationCount)
		
		// User penalties
		moderation.POST("/penalties", h.ApplyUserPenalty)
		moderation.GET("/penalties/user/:id", h.GetUserPenalties)
		moderation.GET("/penalties/user/:id/active", h.GetActivePenalties)
		moderation.PUT("/penalties/:id/remove", h.RemovePenalty)
		moderation.GET("/penalties/user/:id/history", h.GetUserDisciplineHistory)
		moderation.GET("/penalties/user/:id/restricted", h.IsUserRestricted)
	}
}

// FlagContentRequest represents a request to flag content
type FlagContentRequest struct {
	ContentType string `json:"contentType" binding:"required,oneof=topic comment"`
	ContentID   uint   `json:"contentId" binding:"required"`
	FlagType    string `json:"flagType" binding:"required"`
	Description string `json:"description"`
}

// FlagContent flags a content item
func (h *FlagHandler) FlagContent(c *gin.Context) {
	var req FlagContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Parse flag type
	var flagType models.ContentFlagType
	switch req.FlagType {
	case "harassment":
		flagType = models.FlagTypeHarassment
	case "hate_speech":
		flagType = models.FlagTypeHateSpeech
	case "spam":
		flagType = models.FlagTypeSpam
	case "off_topic":
		flagType = models.FlagTypeOffTopic
	case "misleading":
		flagType = models.FlagTypeMisleading
	case "inappropriate":
		flagType = models.FlagTypeInappropriate
	case "violent_content":
		flagType = models.FlagTypeViolentContent
	case "illegal_content":
		flagType = models.FlagTypeIllegalContent
	case "other":
		flagType = models.FlagTypeOther
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flag type"})
		return
	}
	
	// Flag content
	flag, err := h.flagService.FlagContent(
		req.ContentType,
		req.ContentID,
		userID.(uint),
		flagType,
		req.Description,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, flag)
}

// GetFlagsByContent retrieves flags for a content item
func (h *FlagHandler) GetFlagsByContent(c *gin.Context) {
	// Get content type and ID from path
	contentType := c.Param("type")
	contentIDStr := c.Param("id")
	
	// Validate content type
	if contentType != "topic" && contentType != "comment" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Parse content ID
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Get flags
	flags, err := h.flagService.GetFlagsByContent(contentType, uint(contentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, flags)
}

// GetFlagsByStatus retrieves flags by status
func (h *FlagHandler) GetFlagsByStatus(c *gin.Context) {
	// Get status from query
	statusStr := c.DefaultQuery("status", "")
	
	// Parse status
	var status models.FlagStatus
	if statusStr != "" {
		switch statusStr {
		case "pending":
			status = models.FlagStatusPending
		case "reviewed":
			status = models.FlagStatusReviewed
		case "approved":
			status = models.FlagStatusApproved
		case "rejected":
			status = models.FlagStatusRejected
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flag status"})
			return
		}
	}
	
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	// Get flags
	flags, count, err := h.flagService.GetFlagsByStatus(status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"flags":     flags,
		"totalCount": count,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// GetFlagByID retrieves a flag by ID
func (h *FlagHandler) GetFlagByID(c *gin.Context) {
	// Get flag ID from path
	flagIDStr := c.Param("id")
	flagID, err := strconv.ParseUint(flagIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flag ID"})
		return
	}
	
	// Get flag
	flag, err := h.flagService.GetFlagByID(uint(flagID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, flag)
}

// ReviewFlagRequest represents a request to review a flag
type ReviewFlagRequest struct {
	Status      string `json:"status" binding:"required,oneof=reviewed approved rejected"`
	ActionTaken string `json:"actionTaken"`
	Notes       string `json:"notes"`
}

// ReviewFlag reviews a flag
func (h *FlagHandler) ReviewFlag(c *gin.Context) {
	// Get flag ID from path
	flagIDStr := c.Param("id")
	flagID, err := strconv.ParseUint(flagIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flag ID"})
		return
	}
	
	var req ReviewFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Parse status
	var status models.FlagStatus
	switch req.Status {
	case "reviewed":
		status = models.FlagStatusReviewed
	case "approved":
		status = models.FlagStatusApproved
	case "rejected":
		status = models.FlagStatusRejected
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flag status"})
		return
	}
	
	// Review flag
	if err := h.flagService.ReviewFlag(uint(flagID), status, req.ActionTaken, req.Notes, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Flag reviewed successfully"})
}

// AssignFlagRequest represents a request to assign a flag
type AssignFlagRequest struct {
	ModeratorID uint `json:"moderatorId" binding:"required"`
}

// AssignFlag assigns a flag to a moderator
func (h *FlagHandler) AssignFlag(c *gin.Context) {
	// Get flag ID from path
	flagIDStr := c.Param("id")
	flagID, err := strconv.ParseUint(flagIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flag ID"})
		return
	}
	
	var req AssignFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Assign flag
	if err := h.flagService.AssignFlag(uint(flagID), req.ModeratorID, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Flag assigned successfully"})
}

// CreateOrUpdateModerationStatusRequest represents a request to create or update a moderation status
type CreateOrUpdateModerationStatusRequest struct {
	ContentType string `json:"contentType" binding:"required,oneof=topic comment"`
	ContentID   uint   `json:"contentId" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=pending approved rejected hidden"`
	Reason      string `json:"reason"`
	Notes       string `json:"notes"`
}

// CreateOrUpdateModerationStatus creates or updates a moderation status
func (h *FlagHandler) CreateOrUpdateModerationStatus(c *gin.Context) {
	var req CreateOrUpdateModerationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Parse status
	var status models.ContentModerationStatusType
	switch req.Status {
	case "pending":
		status = models.ModerationStatusPending
	case "approved":
		status = models.ModerationStatusApproved
	case "rejected":
		status = models.ModerationStatusRejected
	case "hidden":
		status = models.ModerationStatusHidden
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid moderation status"})
		return
	}
	
	// Create or update moderation status
	moderationStatus, err := h.flagService.CreateOrUpdateModerationStatus(
		req.ContentType,
		req.ContentID,
		status,
		userID.(uint),
		req.Reason,
		req.Notes,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, moderationStatus)
}

// GetModerationStatusByContent retrieves the moderation status for a content item
func (h *FlagHandler) GetModerationStatusByContent(c *gin.Context) {
	// Get content type and ID from path
	contentType := c.Param("type")
	contentIDStr := c.Param("id")
	
	// Validate content type
	if contentType != "topic" && contentType != "comment" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Parse content ID
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Get moderation status
	status, err := h.flagService.GetModerationStatusByContent(contentType, uint(contentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, status)
}

// MarkUserNotified marks a user as notified about a moderation status change
func (h *FlagHandler) MarkUserNotified(c *gin.Context) {
	// Get content type and ID from path
	contentType := c.Param("type")
	contentIDStr := c.Param("id")
	
	// Validate content type
	if contentType != "topic" && contentType != "comment" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Parse content ID
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Mark as notified
	if err := h.flagService.MarkUserNotified(contentType, uint(contentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User marked as notified"})
}

// GetPendingModerationCount retrieves the count of pending moderation items
func (h *FlagHandler) GetPendingModerationCount(c *gin.Context) {
	count, err := h.flagService.GetPendingModerationCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// ApplyUserPenaltyRequest represents a request to apply a penalty to a user
type ApplyUserPenaltyRequest struct {
	UserID             uint   `json:"userId" binding:"required"`
	PenaltyType        string `json:"penaltyType" binding:"required,oneof=warning suspension ban restriction"`
	Reason             string `json:"reason" binding:"required"`
	Description        string `json:"description"`
	Duration           *int   `json:"duration"`
	RelatedContentType *string `json:"relatedContentType"`
	RelatedContentID   *uint  `json:"relatedContentId"`
	Notes              string `json:"notes"`
}

// ApplyUserPenalty applies a penalty to a user
func (h *FlagHandler) ApplyUserPenalty(c *gin.Context) {
	var req ApplyUserPenaltyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Parse penalty type
	var penaltyType models.UserPenaltyType
	switch req.PenaltyType {
	case "warning":
		penaltyType = models.PenaltyTypeWarning
	case "suspension":
		penaltyType = models.PenaltyTypeSuspension
		// Validate duration for suspensions
		if req.Duration == nil || *req.Duration <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duration is required for suspensions"})
			return
		}
	case "ban":
		penaltyType = models.PenaltyTypeBan
	case "restriction":
		penaltyType = models.PenaltyTypeRestriction
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid penalty type"})
		return
	}
	
	// Apply penalty
	penalty, err := h.flagService.ApplyUserPenalty(
		req.UserID,
		penaltyType,
		req.Reason,
		req.Description,
		userID.(uint),
		req.Duration,
		req.RelatedContentType,
		req.RelatedContentID,
		req.Notes,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, penalty)
}

// GetUserPenalties retrieves penalties for a user
func (h *FlagHandler) GetUserPenalties(c *gin.Context) {
	// Get user ID from path
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get penalties
	penalties, err := h.flagService.GetUserPenalties(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, penalties)
}

// GetActivePenalties retrieves active penalties for a user
func (h *FlagHandler) GetActivePenalties(c *gin.Context) {
	// Get user ID from path
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get active penalties
	penalties, err := h.flagService.GetActivePenalties(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, penalties)
}

// RemovePenaltyRequest represents a request to remove a penalty
type RemovePenaltyRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// RemovePenalty removes a penalty
func (h *FlagHandler) RemovePenalty(c *gin.Context) {
	// Get penalty ID from path
	penaltyIDStr := c.Param("id")
	penaltyID, err := strconv.ParseUint(penaltyIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid penalty ID"})
		return
	}
	
	var req RemovePenaltyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Remove penalty
	if err := h.flagService.RemovePenalty(uint(penaltyID), userID.(uint), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Penalty removed successfully"})
}

// GetUserDisciplineHistory retrieves the discipline history for a user
func (h *FlagHandler) GetUserDisciplineHistory(c *gin.Context) {
	// Get user ID from path
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get discipline history
	penalties, err := h.flagService.GetUserDisciplineHistory(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, penalties)
}

// IsUserRestricted checks if a user has any active restrictions
func (h *FlagHandler) IsUserRestricted(c *gin.Context) {
	// Get user ID from path
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Check if user is restricted
	isRestricted, reason, err := h.flagService.IsUserRestricted(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"restricted": isRestricted,
		"reason":     reason,
	})
}