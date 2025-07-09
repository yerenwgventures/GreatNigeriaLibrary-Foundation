package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
)

// SessionService defines the interface for session-related operations
type SessionService interface {
	// Get all active sessions for a user
	GetSessions(userID uint) ([]models.SessionResponse, error)
	
	// Revoke a specific session for a user
	RevokeSession(userID uint, sessionID string) error
	
	// Revoke all sessions for a user except the current one
	RevokeAllSessions(userID uint, currentSessionID string) error
	
	// Admin functions
	PerformMaintenance() (int, error)
}

// SessionHandler handles HTTP requests for session management
type SessionHandler struct {
	sessionService SessionService
	logger         *logger.Logger
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(sessionService SessionService, logger *logger.Logger) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		logger:         logger,
	}
}

// GetSessions retrieves all active sessions for the authenticated user
func (h *SessionHandler) GetSessions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	sessions, err := h.sessionService.GetSessions(userID.(uint))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to get user sessions")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get user sessions"))
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

// RevokeSession revokes a specific session
func (h *SessionHandler) RevokeSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		SessionID string `json:"session_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
		return
	}
	
	// Get the current session ID from the context
	currentSessionID, exists := c.Get("session_id")
	
	// If trying to revoke the current session, handle differently
	if exists && currentSessionID.(string) == req.SessionID {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Cannot revoke current session. Use logout instead."))
		return
	}
	
	err := h.sessionService.RevokeSession(userID.(uint), req.SessionID)
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to revoke session")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to revoke session"))
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Session revoked successfully",
	})
}

// RevokeAllSessions revokes all sessions for the user except the current one
func (h *SessionHandler) RevokeAllSessions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	currentSessionID, _ := c.Get("session_id")
	
	err := h.sessionService.RevokeAllSessions(userID.(uint), currentSessionID.(string))
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to revoke all sessions")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to revoke all sessions"))
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "All other sessions revoked successfully",
	})
}

// PerformMaintenance performs maintenance on sessions (admin only)
func (h *SessionHandler) PerformMaintenance(c *gin.Context) {
	count, err := h.sessionService.PerformMaintenance()
	if err != nil {
		if e, ok := err.(*errors.APIError); ok {
			c.JSON(e.Status, e)
			return
		}
		h.logger.WithError(err).Error("Failed to perform session maintenance")
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to perform session maintenance"))
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Session maintenance completed successfully",
		"expired_sessions_removed": count,
	})
}