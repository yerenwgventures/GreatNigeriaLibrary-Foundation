package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/service"
)

// SubscriptionHandler defines the handler for subscription endpoints
type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(subscriptionService service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// RegisterRoutes registers the routes for subscription management
func (h *SubscriptionHandler) RegisterRoutes(router *gin.RouterGroup) {
	subscriptions := router.Group("/subscriptions")
	{
		// Subscription management
		subscriptions.POST("", h.Subscribe)
		subscriptions.GET("", h.GetUserSubscriptions)
		subscriptions.GET("/:type", h.GetSubscriptionsByType)
		subscriptions.PUT("/:id", h.UpdateSubscription)
		subscriptions.DELETE("/:id", h.Unsubscribe)
		subscriptions.DELETE("/reference", h.UnsubscribeFromReference)
		
		// Subscription preferences
		subscriptions.GET("/preferences", h.GetUserPreferences)
		subscriptions.PUT("/preferences", h.UpdateUserPreferences)
		
		// Digest history
		subscriptions.GET("/digests", h.GetDigestHistory)
	}
}

// SubscribeRequest represents a request to subscribe to a topic, category, or tag
type SubscribeRequest struct {
	ReferenceType string `json:"referenceType" binding:"required"`
	ReferenceID   uint   `json:"referenceId" binding:"required"`
	Frequency     string `json:"frequency" binding:"required"`
	EmailNotification bool `json:"emailNotification"`
	PushNotification  bool `json:"pushNotification"`
	InAppNotification bool `json:"inAppNotification"`
}

// Subscribe handles subscribing to a topic, category, or tag
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var req SubscribeRequest
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
	
	// Convert reference type to subscription type
	var refType models.SubscriptionType
	switch req.ReferenceType {
	case "topic":
		refType = models.TopicSubscription
	case "category":
		refType = models.CategorySubscription
	case "tag":
		refType = models.TagSubscription
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reference type"})
		return
	}
	
	// Convert frequency
	var frequency models.SubscriptionFrequency
	switch req.Frequency {
	case "instant":
		frequency = models.FrequencyInstant
	case "daily":
		frequency = models.FrequencyDaily
	case "weekly":
		frequency = models.FrequencyWeekly
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid frequency"})
		return
	}
	
	// Create subscription
	subscription, err := h.subscriptionService.Subscribe(
		userID.(uint),
		refType,
		req.ReferenceID,
		frequency,
		req.EmailNotification,
		req.PushNotification,
		req.InAppNotification,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, subscription)
}

// GetUserSubscriptions retrieves all subscriptions for the authenticated user
func (h *SubscriptionHandler) GetUserSubscriptions(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get subscriptions
	subscriptions, err := h.subscriptionService.GetUserSubscriptions(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, subscriptions)
}

// GetSubscriptionsByType retrieves subscriptions of a specific type for the authenticated user
func (h *SubscriptionHandler) GetSubscriptionsByType(c *gin.Context) {
	// Get subscription type from path
	typeStr := c.Param("type")
	
	// Convert to subscription type
	var subscriptionType models.SubscriptionType
	switch typeStr {
	case "topic":
		subscriptionType = models.TopicSubscription
	case "category":
		subscriptionType = models.CategorySubscription
	case "tag":
		subscriptionType = models.TagSubscription
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription type"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get subscriptions
	subscriptions, err := h.subscriptionService.GetSubscriptionsByType(userID.(uint), subscriptionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, subscriptions)
}

// UpdateSubscriptionRequest represents a request to update subscription settings
type UpdateSubscriptionRequest struct {
	Frequency     string `json:"frequency" binding:"required"`
	EmailNotification bool `json:"emailNotification"`
	PushNotification  bool `json:"pushNotification"`
	InAppNotification bool `json:"inAppNotification"`
	Muted            bool `json:"muted"`
}

// UpdateSubscription updates subscription settings
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	// Get subscription ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}
	
	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Convert frequency
	var frequency models.SubscriptionFrequency
	switch req.Frequency {
	case "instant":
		frequency = models.FrequencyInstant
	case "daily":
		frequency = models.FrequencyDaily
	case "weekly":
		frequency = models.FrequencyWeekly
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid frequency"})
		return
	}
	
	// Update subscription
	subscription, err := h.subscriptionService.UpdateSubscriptionSettings(
		uint(id),
		frequency,
		req.EmailNotification,
		req.PushNotification,
		req.InAppNotification,
		req.Muted,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, subscription)
}

// Unsubscribe deletes a subscription
func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	// Get subscription ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}
	
	// Delete subscription
	if err := h.subscriptionService.Unsubscribe(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

// UnsubscribeRequest represents a request to unsubscribe from a reference
type UnsubscribeRequest struct {
	ReferenceType string `json:"referenceType" binding:"required"`
	ReferenceID   uint   `json:"referenceId" binding:"required"`
}

// UnsubscribeFromReference deletes a subscription by reference
func (h *SubscriptionHandler) UnsubscribeFromReference(c *gin.Context) {
	var req UnsubscribeRequest
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
	
	// Convert reference type to subscription type
	var refType models.SubscriptionType
	switch req.ReferenceType {
	case "topic":
		refType = models.TopicSubscription
	case "category":
		refType = models.CategorySubscription
	case "tag":
		refType = models.TagSubscription
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reference type"})
		return
	}
	
	// Delete subscription
	if err := h.subscriptionService.UnsubscribeFromReference(userID.(uint), refType, req.ReferenceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

// GetUserPreferences retrieves subscription preferences for the authenticated user
func (h *SubscriptionHandler) GetUserPreferences(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get preferences
	preferences, err := h.subscriptionService.GetUserPreferences(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, preferences)
}

// UpdatePreferencesRequest represents a request to update subscription preferences
type UpdatePreferencesRequest struct {
	DefaultFrequency       string `json:"defaultFrequency" binding:"required"`
	DefaultEmailEnabled    bool   `json:"defaultEmailEnabled"`
	DefaultPushEnabled     bool   `json:"defaultPushEnabled"`
	DefaultInAppEnabled    bool   `json:"defaultInAppEnabled"`
	DigestDay              int    `json:"digestDay"`
	DigestHour             int    `json:"digestHour"`
	AutoSubscribeToReplies bool   `json:"autoSubscribeToReplies"`
	AutoSubscribeToCreated bool   `json:"autoSubscribeToCreated"`
}

// UpdateUserPreferences updates subscription preferences for the authenticated user
func (h *SubscriptionHandler) UpdateUserPreferences(c *gin.Context) {
	var req UpdatePreferencesRequest
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
	
	// Convert frequency
	var frequency models.SubscriptionFrequency
	switch req.DefaultFrequency {
	case "instant":
		frequency = models.FrequencyInstant
	case "daily":
		frequency = models.FrequencyDaily
	case "weekly":
		frequency = models.FrequencyWeekly
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid frequency"})
		return
	}
	
	// Validate digest day (0-6)
	if req.DigestDay < 0 || req.DigestDay > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Digest day must be between 0 (Sunday) and 6 (Saturday)"})
		return
	}
	
	// Validate digest hour (0-23)
	if req.DigestHour < 0 || req.DigestHour > 23 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Digest hour must be between 0 and 23"})
		return
	}
	
	// Create preference object
	preference := models.SubscriptionPreference{
		UserID:                 userID.(uint),
		DefaultFrequency:       frequency,
		DefaultEmailEnabled:    req.DefaultEmailEnabled,
		DefaultPushEnabled:     req.DefaultPushEnabled,
		DefaultInAppEnabled:    req.DefaultInAppEnabled,
		DigestDay:              req.DigestDay,
		DigestHour:             req.DigestHour,
		AutoSubscribeToReplies: req.AutoSubscribeToReplies,
		AutoSubscribeToCreated: req.AutoSubscribeToCreated,
	}
	
	// Update preferences
	updatedPreference, err := h.subscriptionService.UpdateUserPreferences(userID.(uint), preference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, updatedPreference)
}

// GetDigestHistory retrieves digest history for the authenticated user
func (h *SubscriptionHandler) GetDigestHistory(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get page parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	// Get digest history
	digests, err := h.subscriptionService.GetUserDigestHistory(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, digests)
}