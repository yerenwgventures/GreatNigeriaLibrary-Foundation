package handlers

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// ContentAccessService defines the interface for content access service
type ContentAccessService interface {
        // Content access management
        SetContentAccess(contentType string, contentID uint, visibility models.ContentVisibility, minPoints int, isPremium bool) error
        GetContentAccess(contentType string, contentID uint) (*models.ContentAccess, error)
        
        // Content access rule management
        CreateContentRule(rule *models.ContentAccessRule, createdBy uint) error
        UpdateContentRule(rule *models.ContentAccessRule, updatedBy uint) error
        DeleteContentRule(id uint, deletedBy uint) error
        GetContentRules(contentType string, includeInactive bool) ([]models.ContentAccessRule, error)
        
        // User content permissions
        GrantUserPermission(permission *models.UserContentPermission) error
        RevokeUserPermission(id uint, revokedBy uint) error
        GetUserPermissions(userID uint, contentType string) ([]models.UserContentPermission, error)
        
        // User privacy settings
        GetUserPrivacySettings(userID uint) (*models.UserPrivacySettings, error)
        UpdateUserPrivacySettings(userID uint, settings *models.UserPrivacySettings) error
        
        // Content access validation
        CheckContentAccess(userID uint, contentType string, contentID uint) (bool, string)
}

// ContentAccessHandler handles HTTP requests for content access
type ContentAccessHandler struct {
        service ContentAccessService
        logger  *logger.Logger
}

// NewContentAccessHandler creates a new content access handler
func NewContentAccessHandler(service ContentAccessService, logger *logger.Logger) *ContentAccessHandler {
        return &ContentAccessHandler{
                service: service,
                logger:  logger,
        }
}

// SetContentAccess sets access control for content
func (h *ContentAccessHandler) SetContentAccess(c *gin.Context) {
        // Admin only endpoint
        userID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        if !isAdmin.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("Admin privileges required"))
                return
        }

        var req struct {
                ContentType string                   `json:"content_type" binding:"required"`
                ContentID   uint                     `json:"content_id" binding:"required"`
                Visibility  models.ContentVisibility `json:"visibility"`
                MinPoints   int                      `json:"min_points"`
                IsPremium   bool                     `json:"is_premium"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        err := h.service.SetContentAccess(req.ContentType, req.ContentID, req.Visibility, req.MinPoints, req.IsPremium)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to set content access")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to set content access"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Content access set successfully",
        })
}

// GetContentAccess gets access control settings for content
func (h *ContentAccessHandler) GetContentAccess(c *gin.Context) {
        // Moderator or admin endpoint
        isAdmin, _ := c.Get("is_admin")
        isModerator, _ := c.Get("is_moderator")
        if !isAdmin.(bool) && !isModerator.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("Moderator privileges required"))
                return
        }

        contentType := c.Query("content_type")
        if contentType == "" {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Content type is required"))
                return
        }

        contentIDStr := c.Query("content_id")
        if contentIDStr == "" {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Content ID is required"))
                return
        }

        contentID, err := strconv.ParseUint(contentIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid content ID"))
                return
        }

        access, err := h.service.GetContentAccess(contentType, uint(contentID))
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to get content access")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get content access"))
                return
        }

        if access == nil {
                c.JSON(http.StatusOK, gin.H{
                        "message": "No specific access control set for this content (default is public)",
                        "access": map[string]interface{}{
                                "content_type":   contentType,
                                "content_id":     contentID,
                                "visibility":     models.ContentVisibilityPublic,
                                "min_points_req": 0,
                                "is_premium":     false,
                        },
                })
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "access": access,
        })
}

// CreateContentRule creates a new content access rule
func (h *ContentAccessHandler) CreateContentRule(c *gin.Context) {
        // Admin only endpoint
        userID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        if !isAdmin.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("Admin privileges required"))
                return
        }

        var rule models.ContentAccessRule
        if err := c.ShouldBindJSON(&rule); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        err := h.service.CreateContentRule(&rule, userID.(uint))
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to create content rule")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to create content rule"))
                return
        }

        c.JSON(http.StatusCreated, gin.H{
                "message": "Content rule created successfully",
                "rule":    rule,
        })
}

// UpdateContentRule updates a content access rule
func (h *ContentAccessHandler) UpdateContentRule(c *gin.Context) {
        // Admin only endpoint
        userID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        if !isAdmin.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("Admin privileges required"))
                return
        }

        var rule models.ContentAccessRule
        if err := c.ShouldBindJSON(&rule); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        // Check if ID is provided
        if rule.ID == 0 {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Rule ID is required"))
                return
        }

        err := h.service.UpdateContentRule(&rule, userID.(uint))
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to update content rule")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to update content rule"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Content rule updated successfully",
                "rule":    rule,
        })
}

// DeleteContentRule deletes a content access rule
func (h *ContentAccessHandler) DeleteContentRule(c *gin.Context) {
        // Admin only endpoint
        userID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        if !isAdmin.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("Admin privileges required"))
                return
        }

        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid rule ID"))
                return
        }

        err = h.service.DeleteContentRule(uint(id), userID.(uint))
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to delete content rule")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to delete content rule"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Content rule deleted successfully",
        })
}

// GetContentRules gets content access rules
func (h *ContentAccessHandler) GetContentRules(c *gin.Context) {
        // Moderator or admin endpoint
        isAdmin, _ := c.Get("is_admin")
        isModerator, _ := c.Get("is_moderator")
        if !isAdmin.(bool) && !isModerator.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("Moderator privileges required"))
                return
        }

        contentType := c.Query("content_type")
        includeInactiveParam := c.DefaultQuery("include_inactive", "false")
        includeInactive := includeInactiveParam == "true"

        rules, err := h.service.GetContentRules(contentType, includeInactive)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to get content rules")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get content rules"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "rules": rules,
        })
}

// GrantUserPermission grants a permission to a user
func (h *ContentAccessHandler) GrantUserPermission(c *gin.Context) {
        // Moderator or admin endpoint
        userID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        isModerator, _ := c.Get("is_moderator")
        if !isAdmin.(bool) && !isModerator.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("Moderator privileges required"))
                return
        }

        var permission models.UserContentPermission
        if err := c.ShouldBindJSON(&permission); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        // Set granted by
        permission.GrantedByID = userID.(uint)

        err := h.service.GrantUserPermission(&permission)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to grant user permission")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to grant user permission"))
                return
        }

        c.JSON(http.StatusCreated, gin.H{
                "message":    "Permission granted successfully",
                "permission": permission,
        })
}

// RevokeUserPermission revokes a permission from a user
func (h *ContentAccessHandler) RevokeUserPermission(c *gin.Context) {
        // Moderator or admin endpoint
        userID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        isModerator, _ := c.Get("is_moderator")
        if !isAdmin.(bool) && !isModerator.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("Moderator privileges required"))
                return
        }

        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid permission ID"))
                return
        }

        err = h.service.RevokeUserPermission(uint(id), userID.(uint))
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to revoke user permission")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to revoke user permission"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Permission revoked successfully",
        })
}

// GetUserPermissions gets permissions for a user
func (h *ContentAccessHandler) GetUserPermissions(c *gin.Context) {
        // Check if user is requesting their own permissions or is admin/moderator
        requestUserID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        isModerator, _ := c.Get("is_moderator")

        userIDParam := c.Param("user_id")
        userID, err := strconv.ParseUint(userIDParam, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid user ID"))
                return
        }

        // If not self and not admin/moderator, forbid
        if requestUserID.(uint) != uint(userID) && !isAdmin.(bool) && !isModerator.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("You can only view your own permissions"))
                return
        }

        contentType := c.Query("content_type")

        permissions, err := h.service.GetUserPermissions(uint(userID), contentType)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to get user permissions")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get user permissions"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "permissions": permissions,
        })
}

// GetUserPrivacySettings gets privacy settings for a user
func (h *ContentAccessHandler) GetUserPrivacySettings(c *gin.Context) {
        // Check if user is requesting their own settings or is admin
        requestUserID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")

        userIDParam := c.Param("user_id")
        userID, err := strconv.ParseUint(userIDParam, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid user ID"))
                return
        }

        // If not self and not admin, forbid
        if requestUserID.(uint) != uint(userID) && !isAdmin.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("You can only view your own privacy settings"))
                return
        }

        settings, err := h.service.GetUserPrivacySettings(uint(userID))
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to get user privacy settings")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get user privacy settings"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "settings": settings,
        })
}

// UpdateUserPrivacySettings updates privacy settings for a user
func (h *ContentAccessHandler) UpdateUserPrivacySettings(c *gin.Context) {
        // Check if user is updating their own settings or is admin
        requestUserID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")

        userIDParam := c.Param("user_id")
        userID, err := strconv.ParseUint(userIDParam, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid user ID"))
                return
        }

        // If not self and not admin, forbid
        if requestUserID.(uint) != uint(userID) && !isAdmin.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("You can only update your own privacy settings"))
                return
        }

        var settings models.UserPrivacySettings
        if err := c.ShouldBindJSON(&settings); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        // Ensure user ID in settings matches URL
        settings.UserID = uint(userID)

        err = h.service.UpdateUserPrivacySettings(uint(userID), &settings)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to update user privacy settings")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to update user privacy settings"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message":  "Privacy settings updated successfully",
                "settings": settings,
        })
}

// CheckContentAccess checks if a user has access to content
func (h *ContentAccessHandler) CheckContentAccess(c *gin.Context) {
        var req struct {
                ContentType string `json:"content_type" binding:"required"`
                ContentID   uint   `json:"content_id" binding:"required"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        // Get user ID from context if authenticated
        userID, exists := c.Get("user_id")
        var userIDValue uint
        if exists {
                userIDValue = userID.(uint)
        }

        hasAccess, message := h.service.CheckContentAccess(userIDValue, req.ContentType, req.ContentID)
        
        c.JSON(http.StatusOK, gin.H{
                "has_access": hasAccess,
                "message":    message,
        })
}