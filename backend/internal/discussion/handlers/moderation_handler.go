package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/service"
)

// ModerationHandler defines the handler for moderation endpoints
type ModerationHandler struct {
	moderationService service.ModerationService
}

// NewModerationHandler creates a new moderation handler
func NewModerationHandler(moderationService service.ModerationService) *ModerationHandler {
	return &ModerationHandler{
		moderationService: moderationService,
	}
}

// RegisterRoutes registers the routes for moderation
func (h *ModerationHandler) RegisterRoutes(router *gin.RouterGroup) {
	moderation := router.Group("/moderation")
	{
		// Moderation rules
		moderation.POST("/rules", h.CreateModerationRule)
		moderation.GET("/rules", h.GetModerationRules)
		moderation.PUT("/rules/:id", h.UpdateModerationRule)
		moderation.DELETE("/rules/:id", h.DeleteModerationRule)
		
		// Content filtering
		moderation.POST("/filter", h.FilterContent)
		moderation.GET("/filter/content/:type/:id", h.GetFilterResults)
		moderation.PUT("/filter/:id/review", h.ReviewFilterResult)
		moderation.GET("/filter/user/:id", h.GetUserFilterResults)
		
		// Moderation queue
		moderation.POST("/queue", h.AddToModerationQueue)
		moderation.GET("/queue", h.GetModerationQueue)
		moderation.PUT("/queue/:id/assign", h.AssignModerationItem)
		moderation.PUT("/queue/:id/resolve", h.ResolveModerationItem)
		moderation.GET("/queue/stats", h.GetModerationQueueStats)
		
		// User trust management
		moderation.GET("/trust/:userId", h.GetUserTrustScore)
		moderation.PUT("/trust/:userId", h.UpdateUserTrustScore)
		moderation.GET("/trust/level/:level", h.GetUsersByTrustLevel)
		
		// Moderator management
		moderation.POST("/moderators", h.GrantModeratorPrivileges)
		moderation.PUT("/moderators/:userId", h.UpdateModeratorPrivileges)
		moderation.DELETE("/moderators/:userId", h.RevokeModeratorPrivileges)
		moderation.GET("/moderators/:userId", h.GetModeratorPrivileges)
		moderation.GET("/moderators", h.GetAllModerators)
		
		// User moderation actions
		moderation.POST("/actions", h.CreateUserModerationAction)
		moderation.GET("/actions/user/:userId", h.GetUserModerationActions)
		moderation.PUT("/actions/:id/revoke", h.RevokeUserModerationAction)
		moderation.GET("/actions/user/:userId/active", h.GetActiveUserActions)
		moderation.GET("/actions/user/:userId/banned", h.IsUserBanned)
		
		// Prohibited words
		moderation.POST("/prohibited-words", h.AddProhibitedWord)
		moderation.GET("/prohibited-words", h.GetProhibitedWords)
		moderation.PUT("/prohibited-words/:id", h.UpdateProhibitedWord)
		moderation.DELETE("/prohibited-words/:id", h.DeleteProhibitedWord)
		moderation.POST("/filter-text", h.FilterTextWithProhibitedWords)
	}
}

// CreateModerationRuleRequest represents a request to create a moderation rule
type CreateModerationRuleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Pattern     string `json:"pattern" binding:"required"`
	PatternType string `json:"patternType" binding:"required,oneof=regex keywords"`
	Action      string `json:"action" binding:"required"`
	Severity    int    `json:"severity" binding:"required,min=1,max=10"`
	AppliesTo   string `json:"appliesTo" binding:"required,oneof=topic comment username"`
}

// CreateModerationRule creates a moderation rule
func (h *ModerationHandler) CreateModerationRule(c *gin.Context) {
	var req CreateModerationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Convert action to ModerationAction type
	var action models.ModerationAction
	switch req.Action {
	case "none":
		action = models.ActionNone
	case "approve":
		action = models.ActionApprove
	case "reject":
		action = models.ActionReject
	case "send_to_queue":
		action = models.ActionSendToQueue
	case "automatic_filter":
		action = models.ActionAutomaticFilter
	case "warning":
		action = models.ActionWarning
	case "temporary_ban":
		action = models.ActionTemporaryBan
	case "permanent_ban":
		action = models.ActionPermanentBan
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}
	
	// Create rule
	rule, err := h.moderationService.CreateModerationRule(
		req.Name,
		req.Description,
		req.Pattern,
		req.PatternType,
		action,
		req.Severity,
		req.AppliesTo,
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, rule)
}

// GetModerationRules retrieves moderation rules
func (h *ModerationHandler) GetModerationRules(c *gin.Context) {
	// Get the active parameter from query string
	activeStr := c.DefaultQuery("active", "true")
	active, err := strconv.ParseBool(activeStr)
	if err != nil {
		active = true // Default to true if invalid
	}
	
	// Get rules
	rules, err := h.moderationService.GetModerationRules(active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, rules)
}

// UpdateModerationRuleRequest represents a request to update a moderation rule
type UpdateModerationRuleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Pattern     string `json:"pattern" binding:"required"`
	PatternType string `json:"patternType" binding:"required,oneof=regex keywords"`
	Action      string `json:"action" binding:"required"`
	Severity    int    `json:"severity" binding:"required,min=1,max=10"`
	IsActive    bool   `json:"isActive"`
	AppliesTo   string `json:"appliesTo" binding:"required,oneof=topic comment username"`
}

// UpdateModerationRule updates a moderation rule
func (h *ModerationHandler) UpdateModerationRule(c *gin.Context) {
	// Get rule ID from URL
	ruleIDStr := c.Param("id")
	ruleID, err := strconv.ParseUint(ruleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	
	var req UpdateModerationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Convert action to ModerationAction type
	var action models.ModerationAction
	switch req.Action {
	case "none":
		action = models.ActionNone
	case "approve":
		action = models.ActionApprove
	case "reject":
		action = models.ActionReject
	case "send_to_queue":
		action = models.ActionSendToQueue
	case "automatic_filter":
		action = models.ActionAutomaticFilter
	case "warning":
		action = models.ActionWarning
	case "temporary_ban":
		action = models.ActionTemporaryBan
	case "permanent_ban":
		action = models.ActionPermanentBan
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}
	
	// Update rule
	rule, err := h.moderationService.UpdateModerationRule(
		uint(ruleID),
		req.Name,
		req.Description,
		req.Pattern,
		req.PatternType,
		action,
		req.Severity,
		req.IsActive,
		req.AppliesTo,
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, rule)
}

// DeleteModerationRule deletes a moderation rule
func (h *ModerationHandler) DeleteModerationRule(c *gin.Context) {
	// Get rule ID from URL
	ruleIDStr := c.Param("id")
	ruleID, err := strconv.ParseUint(ruleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Delete rule
	if err := h.moderationService.DeleteModerationRule(uint(ruleID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Rule deleted successfully"})
}

// FilterContentRequest represents a request to filter content
type FilterContentRequest struct {
	Content     string `json:"content" binding:"required"`
	ContentType string `json:"contentType" binding:"required,oneof=topic comment username"`
}

// FilterContent filters content for prohibited words and pattern matching
func (h *ModerationHandler) FilterContent(c *gin.Context) {
	var req FilterContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Filter content
	result, err := h.moderationService.FilterContent(req.Content, req.ContentType, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// If no filtering was needed
	if result == nil {
		c.JSON(http.StatusOK, gin.H{
			"filtered": false,
			"content": req.Content,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"filtered": true,
		"result": result,
	})
}

// GetFilterResults retrieves filter results for content
func (h *ModerationHandler) GetFilterResults(c *gin.Context) {
	// Get content type and ID from URL
	contentType := c.Param("type")
	contentIDStr := c.Param("id")
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Get filter results
	results, err := h.moderationService.GetFilterResults(contentType, uint(contentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, results)
}

// ReviewFilterResultRequest represents a request to review a filter result
type ReviewFilterResultRequest struct {
	Action string `json:"action" binding:"required"`
}

// ReviewFilterResult reviews a filter result
func (h *ModerationHandler) ReviewFilterResult(c *gin.Context) {
	// Get filter result ID from URL
	resultIDStr := c.Param("id")
	resultID, err := strconv.ParseUint(resultIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter result ID"})
		return
	}
	
	var req ReviewFilterResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Convert action to ModerationAction type
	var action models.ModerationAction
	switch req.Action {
	case "none":
		action = models.ActionNone
	case "approve":
		action = models.ActionApprove
	case "reject":
		action = models.ActionReject
	case "send_to_queue":
		action = models.ActionSendToQueue
	case "automatic_filter":
		action = models.ActionAutomaticFilter
	case "warning":
		action = models.ActionWarning
	case "temporary_ban":
		action = models.ActionTemporaryBan
	case "permanent_ban":
		action = models.ActionPermanentBan
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}
	
	// Review filter result
	if err := h.moderationService.ReviewFilterResult(uint(resultID), action, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Filter result reviewed successfully"})
}

// GetUserFilterResults retrieves filter results for a user
func (h *ModerationHandler) GetUserFilterResults(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	// Get filter results
	results, err := h.moderationService.GetUserFilterResults(uint(userID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, results)
}

// AddToModerationQueueRequest represents a request to add an item to the moderation queue
type AddToModerationQueueRequest struct {
	ContentType    string `json:"contentType" binding:"required"`
	ContentID      uint   `json:"contentId" binding:"required"`
	Reason         string `json:"reason" binding:"required"`
	FilterResultID *uint  `json:"filterResultId"`
	Priority       int    `json:"priority" binding:"min=1,max=5"`
}

// AddToModerationQueue adds an item to the moderation queue
func (h *ModerationHandler) AddToModerationQueue(c *gin.Context) {
	var req AddToModerationQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Add to moderation queue
	item, err := h.moderationService.AddToModerationQueue(
		req.ContentType,
		req.ContentID,
		userID.(uint),
		req.Reason,
		req.FilterResultID,
		req.Priority,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, item)
}

// GetModerationQueue retrieves items from the moderation queue
func (h *ModerationHandler) GetModerationQueue(c *gin.Context) {
	// Get status from query string
	status := c.DefaultQuery("status", "")
	
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	// Get moderation queue items
	items, count, err := h.moderationService.GetModerationQueue(status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": count,
		"page":  page,
		"size":  pageSize,
	})
}

// AssignModerationItemRequest represents a request to assign a moderation queue item
type AssignModerationItemRequest struct {
	ModeratorID uint `json:"moderatorId" binding:"required"`
}

// AssignModerationItem assigns a moderation queue item to a moderator
func (h *ModerationHandler) AssignModerationItem(c *gin.Context) {
	// Get item ID from URL
	itemIDStr := c.Param("id")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}
	
	var req AssignModerationItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Assign moderation item
	if err := h.moderationService.AssignModerationItem(uint(itemID), req.ModeratorID, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Moderation item assigned successfully"})
}

// ResolveModerationItemRequest represents a request to resolve a moderation queue item
type ResolveModerationItemRequest struct {
	Decision string `json:"decision" binding:"required,oneof=approved rejected"`
	Notes    string `json:"notes"`
}

// ResolveModerationItem resolves a moderation queue item
func (h *ModerationHandler) ResolveModerationItem(c *gin.Context) {
	// Get item ID from URL
	itemIDStr := c.Param("id")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}
	
	var req ResolveModerationItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Resolve moderation item
	if err := h.moderationService.ResolveModerationItem(uint(itemID), req.Decision, req.Notes, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Moderation item resolved successfully"})
}

// GetModerationQueueStats retrieves statistics for the moderation queue
func (h *ModerationHandler) GetModerationQueueStats(c *gin.Context) {
	// Get moderation queue stats
	stats, err := h.moderationService.GetModerationQueueStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// GetUserTrustScore retrieves a user's trust score
func (h *ModerationHandler) GetUserTrustScore(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get user trust score
	score, err := h.moderationService.GetUserTrustScore(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, score)
}

// UpdateUserTrustScoreRequest represents a request to update a user's trust score
type UpdateUserTrustScoreRequest struct {
	ContentScore   float64 `json:"contentScore" binding:"required"`
	CommunityScore float64 `json:"communityScore" binding:"required"`
	ModeratorScore float64 `json:"moderatorScore" binding:"required"`
}

// UpdateUserTrustScore updates a user's trust score
func (h *ModerationHandler) UpdateUserTrustScore(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	var req UpdateUserTrustScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	moderatorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Update user trust score
	score, err := h.moderationService.UpdateUserTrustScore(
		uint(userID),
		req.ContentScore,
		req.CommunityScore,
		req.ModeratorScore,
		moderatorID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, score)
}

// GetUsersByTrustLevel retrieves users by trust level
func (h *ModerationHandler) GetUsersByTrustLevel(c *gin.Context) {
	// Get trust level from URL
	levelStr := c.Param("level")
	
	// Convert to UserTrustLevel
	var level models.UserTrustLevel
	switch levelStr {
	case "new_user":
		level = models.TrustLevelNewUser
	case "basic":
		level = models.TrustLevelBasic
	case "member":
		level = models.TrustLevelMember
	case "regular":
		level = models.TrustLevelRegular
	case "leader":
		level = models.TrustLevelLeader
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trust level"})
		return
	}
	
	// Get users by trust level
	users, err := h.moderationService.GetUsersByTrustLevel(level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

// GrantModeratorPrivilegesRequest represents a request to grant moderator privileges
type GrantModeratorPrivilegesRequest struct {
	UserID             uint `json:"userId" binding:"required"`
	CanApproveContent  bool `json:"canApproveContent"`
	CanRejectContent   bool `json:"canRejectContent"`
	CanEditContent     bool `json:"canEditContent"`
	CanDeleteContent   bool `json:"canDeleteContent"`
	CanBanUsers        bool `json:"canBanUsers"`
	CanManageRules     bool `json:"canManageRules"`
	CanAssignModerators bool `json:"canAssignModerators"`
	CanAccessDashboard bool `json:"canAccessDashboard"`
}

// GrantModeratorPrivileges grants moderator privileges to a user
func (h *ModerationHandler) GrantModeratorPrivileges(c *gin.Context) {
	var req GrantModeratorPrivilegesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	grantedByID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Extract privileges map
	privileges := map[string]bool{
		"canApproveContent":   req.CanApproveContent,
		"canRejectContent":    req.CanRejectContent,
		"canEditContent":      req.CanEditContent,
		"canDeleteContent":    req.CanDeleteContent,
		"canBanUsers":         req.CanBanUsers,
		"canManageRules":      req.CanManageRules,
		"canAssignModerators": req.CanAssignModerators,
		"canAccessDashboard":  req.CanAccessDashboard,
	}
	
	// Grant moderator privileges
	moderator, err := h.moderationService.GrantModeratorPrivileges(req.UserID, grantedByID.(uint), privileges)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, moderator)
}

// UpdateModeratorPrivilegesRequest represents a request to update moderator privileges
type UpdateModeratorPrivilegesRequest struct {
	CanApproveContent  bool `json:"canApproveContent"`
	CanRejectContent   bool `json:"canRejectContent"`
	CanEditContent     bool `json:"canEditContent"`
	CanDeleteContent   bool `json:"canDeleteContent"`
	CanBanUsers        bool `json:"canBanUsers"`
	CanManageRules     bool `json:"canManageRules"`
	CanAssignModerators bool `json:"canAssignModerators"`
	CanAccessDashboard bool `json:"canAccessDashboard"`
}

// UpdateModeratorPrivileges updates a moderator's privileges
func (h *ModerationHandler) UpdateModeratorPrivileges(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	var req UpdateModeratorPrivilegesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	updatedByID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Extract privileges map
	privileges := map[string]bool{
		"canApproveContent":   req.CanApproveContent,
		"canRejectContent":    req.CanRejectContent,
		"canEditContent":      req.CanEditContent,
		"canDeleteContent":    req.CanDeleteContent,
		"canBanUsers":         req.CanBanUsers,
		"canManageRules":      req.CanManageRules,
		"canAssignModerators": req.CanAssignModerators,
		"canAccessDashboard":  req.CanAccessDashboard,
	}
	
	// Update moderator privileges
	moderator, err := h.moderationService.UpdateModeratorPrivileges(uint(userID), updatedByID.(uint), privileges)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, moderator)
}

// RevokeModeratorPrivileges revokes moderator privileges from a user
func (h *ModerationHandler) RevokeModeratorPrivileges(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Check authentication
	revokedByID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Revoke moderator privileges
	if err := h.moderationService.RevokeModeratorPrivileges(uint(userID), revokedByID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Moderator privileges revoked successfully"})
}

// GetModeratorPrivileges retrieves moderator privileges for a user
func (h *ModerationHandler) GetModeratorPrivileges(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get moderator privileges
	privileges, err := h.moderationService.GetModeratorPrivileges(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, privileges)
}

// GetAllModerators retrieves all moderators
func (h *ModerationHandler) GetAllModerators(c *gin.Context) {
	// Get all moderators
	moderators, err := h.moderationService.GetAllModerators()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, moderators)
}

// CreateUserModerationActionRequest represents a request to create a moderation action for a user
type CreateUserModerationActionRequest struct {
	UserID             uint   `json:"userId" binding:"required"`
	ActionType         string `json:"actionType" binding:"required,oneof=warning temporary_ban permanent_ban"`
	Reason             string `json:"reason" binding:"required"`
	Duration           int    `json:"duration"`
	RelatedContentID   *uint  `json:"relatedContentId"`
	RelatedContentType *string `json:"relatedContentType"`
	Notes              string `json:"notes"`
}

// CreateUserModerationAction creates a moderation action for a user
func (h *ModerationHandler) CreateUserModerationAction(c *gin.Context) {
	var req CreateUserModerationActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	moderatorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Convert action type to ModerationAction
	var actionType models.ModerationAction
	switch req.ActionType {
	case "warning":
		actionType = models.ActionWarning
	case "temporary_ban":
		actionType = models.ActionTemporaryBan
		// Validate duration for temporary bans
		if req.Duration <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be positive for temporary bans"})
			return
		}
	case "permanent_ban":
		actionType = models.ActionPermanentBan
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action type"})
		return
	}
	
	// Create user moderation action
	action, err := h.moderationService.CreateUserModerationAction(
		req.UserID,
		actionType,
		req.Reason,
		req.Duration,
		req.RelatedContentID,
		req.RelatedContentType,
		moderatorID.(uint),
		req.Notes,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, action)
}

// GetUserModerationActions retrieves moderation actions for a user
func (h *ModerationHandler) GetUserModerationActions(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get user moderation actions
	actions, err := h.moderationService.GetUserModerationActions(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, actions)
}

// RevokeUserModerationActionRequest represents a request to revoke a user moderation action
type RevokeUserModerationActionRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// RevokeUserModerationAction revokes a user moderation action
func (h *ModerationHandler) RevokeUserModerationAction(c *gin.Context) {
	// Get action ID from URL
	actionIDStr := c.Param("id")
	actionID, err := strconv.ParseUint(actionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action ID"})
		return
	}
	
	var req RevokeUserModerationActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	moderatorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Revoke user moderation action
	if err := h.moderationService.RevokeUserModerationAction(uint(actionID), moderatorID.(uint), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User moderation action revoked successfully"})
}

// GetActiveUserActions retrieves active moderation actions for a user
func (h *ModerationHandler) GetActiveUserActions(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get active user actions
	actions, err := h.moderationService.GetActiveUserActions(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, actions)
}

// IsUserBanned checks if a user is banned
func (h *ModerationHandler) IsUserBanned(c *gin.Context) {
	// Get user ID from URL
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Check if user is banned
	banned, err := h.moderationService.IsUserBanned(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"banned": banned})
}

// AddProhibitedWordRequest represents a request to add a prohibited word
type AddProhibitedWordRequest struct {
	Word          string `json:"word" binding:"required"`
	Replacement   string `json:"replacement"`
	IsRegex       bool   `json:"isRegex"`
	IsAutoReplace bool   `json:"isAutoReplace"`
	Severity      int    `json:"severity" binding:"required,min=1,max=10"`
}

// AddProhibitedWord adds a prohibited word
func (h *ModerationHandler) AddProhibitedWord(c *gin.Context) {
	var req AddProhibitedWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Add prohibited word
	word, err := h.moderationService.AddProhibitedWord(
		req.Word,
		req.Replacement,
		req.IsRegex,
		req.IsAutoReplace,
		req.Severity,
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, word)
}

// GetProhibitedWords retrieves prohibited words
func (h *ModerationHandler) GetProhibitedWords(c *gin.Context) {
	// Get the active parameter from query string
	activeStr := c.DefaultQuery("active", "true")
	active, err := strconv.ParseBool(activeStr)
	if err != nil {
		active = true // Default to true if invalid
	}
	
	// Get prohibited words
	words, err := h.moderationService.GetProhibitedWords(active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, words)
}

// UpdateProhibitedWordRequest represents a request to update a prohibited word
type UpdateProhibitedWordRequest struct {
	Word          string `json:"word" binding:"required"`
	Replacement   string `json:"replacement"`
	IsRegex       bool   `json:"isRegex"`
	IsAutoReplace bool   `json:"isAutoReplace"`
	IsActive      bool   `json:"isActive"`
	Severity      int    `json:"severity" binding:"required,min=1,max=10"`
}

// UpdateProhibitedWord updates a prohibited word
func (h *ModerationHandler) UpdateProhibitedWord(c *gin.Context) {
	// Get word ID from URL
	wordIDStr := c.Param("id")
	wordID, err := strconv.ParseUint(wordIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}
	
	var req UpdateProhibitedWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Update prohibited word
	word, err := h.moderationService.UpdateProhibitedWord(
		uint(wordID),
		req.Word,
		req.Replacement,
		req.IsRegex,
		req.IsAutoReplace,
		req.IsActive,
		req.Severity,
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, word)
}

// DeleteProhibitedWord deletes a prohibited word
func (h *ModerationHandler) DeleteProhibitedWord(c *gin.Context) {
	// Get word ID from URL
	wordIDStr := c.Param("id")
	wordID, err := strconv.ParseUint(wordIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}
	
	// Check authentication
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Delete prohibited word
	if err := h.moderationService.DeleteProhibitedWord(uint(wordID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Prohibited word deleted successfully"})
}

// FilterTextRequest represents a request to filter text for prohibited words
type FilterTextRequest struct {
	Text string `json:"text" binding:"required"`
}

// FilterTextWithProhibitedWords filters text for prohibited words
func (h *ModerationHandler) FilterTextWithProhibitedWords(c *gin.Context) {
	var req FilterTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Filter text
	filtered, wasFiltered, filteredWords, err := h.moderationService.FilterTextWithProhibitedWords(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"original":      req.Text,
		"filtered":      filtered,
		"wasFiltered":   wasFiltered,
		"filteredWords": filteredWords,
	})
}