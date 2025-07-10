package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/auth"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/response"
)

// ProfileCompletionHandler handles profile completion-related requests
type ProfileCompletionHandler struct {
	userService UserServiceInterface
}

// NewProfileCompletionHandler creates a new profile completion handler
func NewProfileCompletionHandler(userService UserServiceInterface) *ProfileCompletionHandler {
	return &ProfileCompletionHandler{
		userService: userService,
	}
}

// GetProfileCompletionStatus gets the profile completion status for a user
func (h *ProfileCompletionHandler) GetProfileCompletionStatus(c *gin.Context) {
	// Get the current user ID from the token
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// Get the profile completion status response
	status, err := h.userService.GetProfileCompletionResponse(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get profile completion status", err)
		return
	}

	response.Success(c, http.StatusOK, "Profile completion status retrieved", status)
}

// GetUserProfileCompletionStatus gets the profile completion status for a specific user (admin/moderator only)
func (h *ProfileCompletionHandler) GetUserProfileCompletionStatus(c *gin.Context) {
	// Get the current user ID from the token (reviewer)
	reviewerID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// Check if the reviewer has moderator permissions
	reviewer, err := h.userService.GetUserByID(reviewerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get reviewer information", err)
		return
	}

	if !reviewer.IsUserModerator() {
		response.Error(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	// Get the target user ID from the URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Get the profile completion status response
	status, err := h.userService.GetProfileCompletionResponse(uint(userID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get profile completion status", err)
		return
	}

	response.Success(c, http.StatusOK, "Profile completion status retrieved", status)
}

// UpdateProfileCompletionField updates a profile completion field
func (h *ProfileCompletionHandler) UpdateProfileCompletionFromActivity(c *gin.Context) {
	// Get the current user ID from the token
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// Get the activity type from the request
	type ActivityRequest struct {
		ActivityType string `json:"activity_type" binding:"required"`
	}

	var req ActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Update the profile completion based on activity
	if err := h.userService.UpdateProfileCompletionFromActivity(userID, req.ActivityType); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update profile completion", err)
		return
	}

	// Get the updated profile completion status
	status, err := h.userService.GetProfileCompletionResponse(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get updated profile completion status", err)
		return
	}

	response.Success(c, http.StatusOK, "Profile completion updated", status)
}

// CheckProfileCompletionReminder checks if a profile completion reminder should be sent
func (h *ProfileCompletionHandler) CheckProfileCompletionReminder(c *gin.Context) {
	// Get the current user ID from the token
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// Check if a reminder should be sent
	shouldSend, err := h.userService.SendProfileCompletionReminder(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to check profile completion reminder", err)
		return
	}

	// Get the profile completion status
	status, err := h.userService.GetProfileCompletionResponse(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get profile completion status", err)
		return
	}

	response.Success(c, http.StatusOK, "Profile completion reminder check completed", gin.H{
		"should_send_reminder": shouldSend,
		"completion_status":    status,
	})
}