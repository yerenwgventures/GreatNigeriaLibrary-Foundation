package handlers

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// UserServiceInterface defines the interface for user service operations needed by verification handler
type UserServiceInterface interface {
        // Basic user operations
        GetUserByID(id uint) (*models.User, error)
        
        // Verification operations
        GetVerificationStatus(userID uint) (*models.VerificationStatus, error)
        CreateVerificationRequest(request models.VerificationRequest) error
        GetVerificationRequests(userID uint) ([]models.VerificationRequest, error)
        GetVerificationRequestByID(id uint) (*models.VerificationRequest, error)
        ReviewVerificationRequest(id, reviewerID uint, status, notes string) error
        UpdateVerificationField(userID uint, field string, value bool) error
        SetUserVerified(userID uint, verified bool) error
        
        // Badge operations
        GetUserBadges(userID uint) ([]models.UserBadge, error)
        AwardBadgeToUser(userID, badgeID, awardedBy uint, reason string) error
        
        // Trust level operations
        UpdateUserTrustLevel(userID uint, level models.TrustLevel) error
        
        // Profile completion operations
        GetProfileCompletionStatus(userID uint) (*models.ProfileCompletionStatus, error)
        GetProfileCompletionResponse(userID uint) (*models.ProfileCompletionResponse, error)
        UpdateProfileCompletionField(userID uint, field string, value bool) error
        UpdateProfileCompletionFromUserUpdate(userID uint, updateData *models.UserUpdateRequest) error
        UpdateProfileCompletionFromVerification(userID uint, verType string, status bool) error
        UpdateProfileCompletionFromActivity(userID uint, activityType string) error
        SendProfileCompletionReminder(userID uint) (bool, error)
}