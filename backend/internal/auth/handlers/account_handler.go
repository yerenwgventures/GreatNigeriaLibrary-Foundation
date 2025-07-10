package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// AccountHandler handles account management requests
type AccountHandler struct {
	userService UserService
	logger      *logger.Logger
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(userService UserService, logger *logger.Logger) *AccountHandler {
	return &AccountHandler{
		userService: userService,
		logger:      logger,
	}
}

// DeleteAccount handles account deletion requests
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, errors.ErrUnauthorized("Authentication required"))
		return
	}

	userID, ok := userIDStr.(uint)
	if !ok {
		h.logger.Error("Failed to parse user ID from context")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Invalid user information"))
		return
	}

	// Parse request body to get password confirmation
	var req models.AccountDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid account deletion request")
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid request format"))
		return
	}

	// Validate password and delete account
	if err := h.userService.DeleteUser(userID, req.Password); err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).WithField("user_id", userID).Error("Failed to delete account")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to delete account"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Account successfully deleted",
	})
}