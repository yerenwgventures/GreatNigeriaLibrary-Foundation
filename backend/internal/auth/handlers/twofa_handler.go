package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/auth/service"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
)

// TwoFAHandler handles two-factor authentication related requests
type TwoFAHandler struct {
	twoFAService service.TwoFAService
	logger       *logger.Logger
}

// NewTwoFAHandler creates a new TwoFAHandler instance
func NewTwoFAHandler(twoFAService service.TwoFAService, logger *logger.Logger) *TwoFAHandler {
	return &TwoFAHandler{
		twoFAService: twoFAService,
		logger:       logger,
	}
}

// SetupTwoFA initializes 2FA for a user
func (h *TwoFAHandler) SetupTwoFA(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Setup 2FA
	setupResponse, err := h.twoFAService.SetupTwoFA(userID.(uint))
	if err != nil {
		h.logger.Error("Failed to setup 2FA", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup 2FA"})
		return
	}
	
	c.JSON(http.StatusOK, setupResponse)
}

// VerifyTwoFA verifies a 2FA token without enabling 2FA
func (h *TwoFAHandler) VerifyTwoFA(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Parse request
	var req models.TwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	// Verify token
	valid, err := h.twoFAService.VerifyTwoFA(userID.(uint), req.Token)
	if err != nil {
		h.logger.Error("Failed to verify 2FA token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
		return
	}
	
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token", "valid": false})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Token verified successfully", "valid": true})
}

// EnableTwoFA enables 2FA for a user
func (h *TwoFAHandler) EnableTwoFA(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Parse request
	var req models.TwoFactorEnableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	// Enable 2FA
	err := h.twoFAService.EnableTwoFA(userID.(uint), req.Token)
	if err != nil {
		h.logger.Error("Failed to enable 2FA", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Generate backup codes
	backupCodes, err := h.twoFAService.GenerateBackupCodes(userID.(uint))
	if err != nil {
		h.logger.Error("Failed to generate backup codes", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Two-factor authentication enabled successfully",
		"backup_codes": backupCodes,
	})
}

// DisableTwoFA disables 2FA for a user
func (h *TwoFAHandler) DisableTwoFA(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Parse request
	var req models.TwoFactorDisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	// Disable 2FA
	err := h.twoFAService.DisableTwoFA(userID.(uint), req.Token, req.Password)
	if err != nil {
		h.logger.Error("Failed to disable 2FA", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Two-factor authentication disabled successfully"})
}

// GetTwoFAStatus gets the current 2FA status for a user
func (h *TwoFAHandler) GetTwoFAStatus(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Get 2FA status
	status, err := h.twoFAService.GetTwoFAStatus(userID.(uint))
	if err != nil {
		h.logger.Error("Failed to get 2FA status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get 2FA status"})
		return
	}
	
	c.JSON(http.StatusOK, status)
}

// GenerateBackupCodes generates new backup codes for a user
func (h *TwoFAHandler) GenerateBackupCodes(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Generate backup codes
	backupCodes, err := h.twoFAService.GenerateBackupCodes(userID.(uint))
	if err != nil {
		h.logger.Error("Failed to generate backup codes", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Backup codes generated successfully",
		"backup_codes": backupCodes,
	})
}

// ValidateBackupCode validates a backup code
func (h *TwoFAHandler) ValidateBackupCode(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Parse request
	var req models.TwoFactorBackupCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	// Validate backup code
	valid, err := h.twoFAService.ValidateBackupCode(userID.(uint), req.BackupCode)
	if err != nil {
		h.logger.Error("Failed to validate backup code", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate backup code"})
		return
	}
	
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup code", "valid": false})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Backup code validated successfully", "valid": true})
}