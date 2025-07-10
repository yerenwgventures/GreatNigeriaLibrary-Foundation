package handlers

import (
        "net/http"
        "strconv"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/auth"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/response"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// VerificationHandler handles verification-related requests
type VerificationHandler struct {
        userService UserServiceInterface
}

// NewVerificationHandler creates a new verification handler
func NewVerificationHandler(userService UserServiceInterface) *VerificationHandler {
        return &VerificationHandler{
                userService: userService,
        }
}

// GetVerificationStatus returns the verification status for a user
func (h *VerificationHandler) GetVerificationStatus(c *gin.Context) {
        // Get the current user ID from the token
        userID, err := auth.GetUserIDFromContext(c)
        if err != nil {
                response.Error(c, http.StatusUnauthorized, "Unauthorized", err)
                return
        }

        // Get the verification status
        status, err := h.userService.GetVerificationStatus(userID)
        if err != nil {
                response.Error(c, http.StatusInternalServerError, "Failed to get verification status", err)
                return
        }

        response.Success(c, http.StatusOK, "Verification status retrieved", status)
}

// SubmitVerificationRequest submits a verification request
func (h *VerificationHandler) SubmitVerificationRequest(c *gin.Context) {
        // Get the current user ID from the token
        userID, err := auth.GetUserIDFromContext(c)
        if err != nil {
                response.Error(c, http.StatusUnauthorized, "Unauthorized", err)
                return
        }

        var req models.VerificationRequestCreateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                response.Error(c, http.StatusBadRequest, "Invalid request data", err)
                return
        }

        // Create the verification request
        verReq := models.VerificationRequest{
                UserID:       userID,
                Type:         req.Type,
                Status:       "pending",
                DocumentType: req.DocumentType,
                DocumentURL:  req.DocumentURL,
                Notes:        req.Notes,
                CreatedAt:    time.Now(),
                UpdatedAt:    time.Now(),
        }

        if err := h.userService.CreateVerificationRequest(verReq); err != nil {
                response.Error(c, http.StatusInternalServerError, "Failed to submit verification request", err)
                return
        }

        response.Success(c, http.StatusCreated, "Verification request submitted", gin.H{
                "request_type": req.Type,
                "status":       "pending",
        })
}

// GetVerificationRequests returns all verification requests for a user
func (h *VerificationHandler) GetVerificationRequests(c *gin.Context) {
        // Get the current user ID from the token
        userID, err := auth.GetUserIDFromContext(c)
        if err != nil {
                response.Error(c, http.StatusUnauthorized, "Unauthorized", err)
                return
        }

        // Get the verification requests
        requests, err := h.userService.GetVerificationRequests(userID)
        if err != nil {
                response.Error(c, http.StatusInternalServerError, "Failed to get verification requests", err)
                return
        }

        response.Success(c, http.StatusOK, "Verification requests retrieved", requests)
}

// ReviewVerificationRequest reviews a verification request (admin/moderator only)
func (h *VerificationHandler) ReviewVerificationRequest(c *gin.Context) {
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

        // Get the request ID from the URL parameter
        requestIDStr := c.Param("id")
        requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
        if err != nil {
                response.Error(c, http.StatusBadRequest, "Invalid request ID", err)
                return
        }

        var req models.VerificationRequestUpdateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                response.Error(c, http.StatusBadRequest, "Invalid request data", err)
                return
        }

        // Update the verification request
        err = h.userService.ReviewVerificationRequest(uint(requestID), reviewerID, req.Status, req.ReviewNotes)
        if err != nil {
                response.Error(c, http.StatusInternalServerError, "Failed to update verification request", err)
                return
        }

        // If approved, update user verification status
        if req.Status == "approved" {
                verReq, err := h.userService.GetVerificationRequestByID(uint(requestID))
                if err != nil {
                        response.Error(c, http.StatusInternalServerError, "Failed to get verification request", err)
                        return
                }

                // Update the appropriate verification status field
                switch verReq.Type {
                case "email":
                        err = h.userService.UpdateVerificationField(verReq.UserID, "email_verified", true)
                case "phone":
                        err = h.userService.UpdateVerificationField(verReq.UserID, "phone_verified", true)
                case "identity":
                        err = h.userService.UpdateVerificationField(verReq.UserID, "identity_verified", true)
                case "address":
                        err = h.userService.UpdateVerificationField(verReq.UserID, "address_verified", true)
                }

                if err != nil {
                        response.Error(c, http.StatusInternalServerError, "Failed to update verification status", err)
                        return
                }
                
                // Update profile completion status based on verification change
                if err := h.userService.UpdateProfileCompletionFromVerification(verReq.UserID, verReq.Type, true); err != nil {
                        // Log but don't fail the request
                        // TODO: Implement logging
                }

                // Update overall user verification status if email is verified (minimum requirement)
                user, err := h.userService.GetUserByID(verReq.UserID)
                if err != nil {
                        response.Error(c, http.StatusInternalServerError, "Failed to get user", err)
                        return
                }

                verStatus, err := h.userService.GetVerificationStatus(verReq.UserID)
                if err != nil {
                        response.Error(c, http.StatusInternalServerError, "Failed to get verification status", err)
                        return
                }

                if verStatus.EmailVerified && !user.IsVerified {
                        if err := h.userService.SetUserVerified(verReq.UserID, true); err != nil {
                                response.Error(c, http.StatusInternalServerError, "Failed to update user verification status", err)
                                return
                        }
                }

                // Check for trust level updates
                h.checkAndUpdateTrustLevel(verReq.UserID)
        }

        response.Success(c, http.StatusOK, "Verification request updated", gin.H{
                "request_id": requestID,
                "status":     req.Status,
        })
}

// GetUserBadges returns all badges for a user
func (h *VerificationHandler) GetUserBadges(c *gin.Context) {
        // Get user ID from URL parameter
        userIDStr := c.Param("id")
        userID, err := strconv.ParseUint(userIDStr, 10, 64)
        if err != nil {
                // If not provided or invalid, use the current user's ID
                currentUserID, err := auth.GetUserIDFromContext(c)
                if err != nil {
                        response.Error(c, http.StatusUnauthorized, "Unauthorized", err)
                        return
                }
                userID = uint64(currentUserID)
        }

        // Check if the requested user exists
        user, err := h.userService.GetUserByID(uint(userID))
        if err != nil {
                response.Error(c, http.StatusNotFound, "User not found", err)
                return
        }

        // Get the user's badges
        badges, err := h.userService.GetUserBadges(uint(userID))
        if err != nil {
                response.Error(c, http.StatusInternalServerError, "Failed to get user badges", err)
                return
        }

        // Filter out hidden badges if the requester is not the user or an admin
        currentUserID, _ := auth.GetUserIDFromContext(c)
        currentUser, _ := h.userService.GetUserByID(currentUserID)
        
        if currentUserID != uint(userID) && !currentUser.IsUserAdmin() {
                var publicBadges []models.UserBadge
                for _, badge := range badges {
                        if !badge.IsHidden && badge.IsPublic {
                                publicBadges = append(publicBadges, badge)
                        }
                }
                badges = publicBadges
        }

        response.Success(c, http.StatusOK, "User badges retrieved", gin.H{
                "user_id":   userID,
                "username":  user.Username,
                "badges":    badges,
                "trust_level": user.TrustLevel,
                "trust_level_name": user.GetTrustLevelName(),
        })
}

// checkAndUpdateTrustLevel checks if a user is eligible for a trust level update
func (h *VerificationHandler) checkAndUpdateTrustLevel(userID uint) error {
        // Get verification status
        verStatus, err := h.userService.GetVerificationStatus(userID)
        if err != nil {
                return err
        }

        // Get user
        user, err := h.userService.GetUserByID(userID)
        if err != nil {
                return err
        }

        // Define criteria for each trust level
        // These are just examples - real implementation would include more complex logic
        // based on user activity, post quality, etc.
        var newTrustLevel models.TrustLevel

        // Basic trust level: Email verified
        if verStatus.EmailVerified {
                newTrustLevel = models.TrustLevelBasic
        }

        // Member trust level: Email + Phone verified
        if verStatus.EmailVerified && verStatus.PhoneVerified {
                newTrustLevel = models.TrustLevelMember
        }

        // Regular trust level: Email + Phone + either identity or address verified
        if verStatus.EmailVerified && verStatus.PhoneVerified && 
                (verStatus.IdentityVerified || verStatus.AddressVerified) {
                newTrustLevel = models.TrustLevelRegular
        }

        // Trusted member: All verification types completed
        if verStatus.EmailVerified && verStatus.PhoneVerified && 
                verStatus.IdentityVerified && verStatus.AddressVerified {
                newTrustLevel = models.TrustLevelTrusted
        }

        // Leader trust level requires manual promotion by admins and additional criteria
        // (like tenure, contribution quality, etc.)

        // Only update if the new trust level is higher than the current one
        if newTrustLevel > user.TrustLevel {
                return h.userService.UpdateUserTrustLevel(userID, newTrustLevel)
        }

        return nil
}