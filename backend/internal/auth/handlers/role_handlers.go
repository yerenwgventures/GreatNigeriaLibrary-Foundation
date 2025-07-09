package handlers

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/common/errors"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
)

// Define minimal UserService interface needed for role handlers
type UserService interface {
        GetUserByID(id uint) (*models.User, error)
}

// RoleHandlers contains handlers for different role-specific routes
type RoleHandlers struct {
        userService UserService
}

// NewRoleHandlers creates a new RoleHandlers instance
func NewRoleHandlers(userService UserService) *RoleHandlers {
        return &RoleHandlers{
                userService: userService,
        }
}

// GetEngagedUserFeatures returns features available to engaged users
func (h *RoleHandlers) GetEngagedUserFeatures(c *gin.Context) {
        // Get user ID from context (set by middleware)
        userID, _ := c.Get("user_id")
        user, err := h.userService.GetUserByID(uint(userID.(float64)))
        if err != nil {
                c.JSON(http.StatusInternalServerError, errors.NewAPIError("Error", "Failed to get user", http.StatusInternalServerError))
                return
        }

        // Check if user meets minimum role requirements (double check in addition to middleware)
        if user.Role < models.RoleEngagedUser {
                c.JSON(http.StatusForbidden, errors.NewAPIError("Forbidden", "Insufficient permissions", http.StatusForbidden))
                return
        }

        // Features available to engaged users
        features := []string{
                "Advanced book discussions",
                "Community blog comments",
                "Chapter summaries",
                "Reading guides",
                "Basic community badges",
        }

        c.JSON(http.StatusOK, gin.H{
                "role":     user.GetRoleName(),
                "features": features,
        })
}

// GetActiveUserFeatures returns features available to active users
func (h *RoleHandlers) GetActiveUserFeatures(c *gin.Context) {
        // Get user ID from context (set by middleware)
        userID, _ := c.Get("user_id")
        user, err := h.userService.GetUserByID(uint(userID.(float64)))
        if err != nil {
                c.JSON(http.StatusInternalServerError, errors.NewAPIError("Error", "Failed to get user", http.StatusInternalServerError))
                return
        }

        // Check if user meets minimum role requirements (double check in addition to middleware)
        if user.Role < models.RoleActiveUser {
                c.JSON(http.StatusForbidden, errors.NewAPIError("Forbidden", "Insufficient permissions", http.StatusForbidden))
                return
        }

        // Features available to active users
        features := []string{
                "All engaged user features",
                "Interactive chapter quizzes",
                "Reading lists with progress tracking",
                "Community discussion creation",
                "Advanced community badges",
                "Personalized reading recommendations",
                "Access to exclusive events",
        }

        c.JSON(http.StatusOK, gin.H{
                "role":     user.GetRoleName(),
                "features": features,
        })
}

// GetPremiumUserFeatures returns features available to premium users
func (h *RoleHandlers) GetPremiumUserFeatures(c *gin.Context) {
        // Get user ID from context (set by middleware)
        userID, _ := c.Get("user_id")
        user, err := h.userService.GetUserByID(uint(userID.(float64)))
        if err != nil {
                c.JSON(http.StatusInternalServerError, errors.NewAPIError("Error", "Failed to get user", http.StatusInternalServerError))
                return
        }

        // Check if user meets minimum role requirements (double check in addition to middleware)
        if user.Role < models.RolePremiumUser {
                c.JSON(http.StatusForbidden, errors.NewAPIError("Forbidden", "Insufficient permissions", http.StatusForbidden))
                return
        }

        // Features available to premium users
        features := []string{
                "All active user features",
                "Exclusive premium content",
                "Early access to new chapters",
                "Downloadable materials",
                "Direct author Q&A sessions",
                "Premium community badges",
                "Personalized reading analytics",
                "Ad-free experience",
                "Priority customer support",
        }

        c.JSON(http.StatusOK, gin.H{
                "role":     user.GetRoleName(),
                "features": features,
        })
}

// GetModeratorTools returns tools available to moderators
func (h *RoleHandlers) GetModeratorTools(c *gin.Context) {
        // Get user ID from context (set by middleware)
        userID, _ := c.Get("user_id")
        user, err := h.userService.GetUserByID(uint(userID.(float64)))
        if err != nil {
                c.JSON(http.StatusInternalServerError, errors.NewAPIError("Error", "Failed to get user", http.StatusInternalServerError))
                return
        }

        // Check if user meets minimum role requirements (double check in addition to middleware)
        if user.Role < models.RoleModerator {
                c.JSON(http.StatusForbidden, errors.NewAPIError("Forbidden", "Insufficient permissions", http.StatusForbidden))
                return
        }

        // Moderation tools
        tools := []string{
                "All premium user features",
                "Discussion moderation dashboard",
                "Comment moderation tools",
                "User report management",
                "Content flagging tools",
                "Community guidelines enforcement",
                "User activity monitoring",
                "Moderation action logs",
                "Content quality verification",
        }

        c.JSON(http.StatusOK, gin.H{
                "role":  user.GetRoleName(),
                "tools": tools,
        })
}