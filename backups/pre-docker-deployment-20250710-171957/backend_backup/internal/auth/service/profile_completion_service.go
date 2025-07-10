package service

import (
	"errors"
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
	"gorm.io/gorm"
)

// GetProfileCompletionStatus gets a user's profile completion status
func (s *UserService) GetProfileCompletionStatus(userID uint) (*models.ProfileCompletionStatus, error) {
	var status models.ProfileCompletionStatus
	result := s.db.Where("user_id = ?", userID).First(&status)
	
	if result.Error != nil {
		// If no status exists, create a new one with default values
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Get user to check existing fields
			user, err := s.GetUserByID(userID)
			if err != nil {
				return nil, err
			}
			
			// Get verification status to check verification fields
			verStatus, _ := s.GetVerificationStatus(userID)
			
			// Create a new profile completion status
			status = models.ProfileCompletionStatus{
				UserID:                  userID,
				HasUsername:             true, // Required at registration
				HasEmail:                true, // Required at registration
				HasFullName:             true, // Required at registration
				HasProfileImage:         user.ProfileImage != "",
				HasBio:                  user.Bio != "",
				HasEmailVerified:        false,
				HasPhoneNumber:          false,
				HasPhoneVerified:        false,
				HasIdentityVerified:     false,
				HasAddressVerified:      false,
				HasSetPreferences:       false,
				HasJoinedGroups:         false,
				HasBookmarkedContent:    false,
				HasParticipatedInForum:  false,
				HasSharedContent:        false,
				HasEarnedPoints:         user.PointsBalance > 0,
				CompletionPercentage:    0,
				LastUpdated:             time.Now(),
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
			}
			
			// If verification status exists, update the completion status
			if verStatus != nil {
				status.HasEmailVerified = verStatus.EmailVerified
				status.HasPhoneVerified = verStatus.PhoneVerified
				status.HasIdentityVerified = verStatus.IdentityVerified
				status.HasAddressVerified = verStatus.AddressVerified
			}
			
			// Calculate initial completion percentage
			status.CalculateCompletionPercentage()
			
			// Save to database
			if err := s.db.Create(&status).Error; err != nil {
				return nil, err
			}
			
			return &status, nil
		}
		
		return nil, result.Error
	}
	
	return &status, nil
}

// UpdateProfileCompletionField updates a specific profile completion field
func (s *UserService) UpdateProfileCompletionField(userID uint, field string, value bool) error {
	// First get or create the profile completion status
	status, err := s.GetProfileCompletionStatus(userID)
	if err != nil {
		return err
	}
	
	// Update the specific field
	updates := map[string]interface{}{
		field:          value,
		"updated_at":   time.Now(),
		"last_updated": time.Now(),
	}
	
	// Save the changes
	if err := s.db.Model(&models.ProfileCompletionStatus{}).
		Where("user_id = ?", userID).
		Updates(updates).Error; err != nil {
		return err
	}
	
	// Recalculate completion percentage
	status, err = s.GetProfileCompletionStatus(userID)
	if err != nil {
		return err
	}
	
	status.CalculateCompletionPercentage()
	
	return s.db.Model(&models.ProfileCompletionStatus{}).
		Where("user_id = ?", userID).
		Update("completion_percentage", status.CompletionPercentage).Error
}

// GetProfileCompletionResponse gets a user's profile completion response for the client
func (s *UserService) GetProfileCompletionResponse(userID uint) (*models.ProfileCompletionResponse, error) {
	// Get the profile completion status
	status, err := s.GetProfileCompletionStatus(userID)
	if err != nil {
		return nil, err
	}
	
	// Create the response
	response := status.CreateProfileCompletionResponse()
	
	return &response, nil
}

// UpdateProfileCompletionFromUserUpdate updates profile completion based on user profile updates
func (s *UserService) UpdateProfileCompletionFromUserUpdate(userID uint, updateData *models.UserUpdateRequest) error {
	// Update profile image field if it's not empty
	if updateData.ProfileImage != "" {
		if err := s.UpdateProfileCompletionField(userID, "has_profile_image", true); err != nil {
			return err
		}
	}
	
	// Update bio field if it's not empty
	if updateData.Bio != "" {
		if err := s.UpdateProfileCompletionField(userID, "has_bio", true); err != nil {
			return err
		}
	}
	
	return nil
}

// UpdateProfileCompletionFromVerification updates profile completion based on verification changes
func (s *UserService) UpdateProfileCompletionFromVerification(userID uint, verType string, status bool) error {
	if !status {
		return nil // No need to update if verification was unsuccessful
	}
	
	// Update the appropriate field based on verification type
	switch verType {
	case "email":
		return s.UpdateProfileCompletionField(userID, "has_email_verified", true)
	case "phone":
		// Update both phone fields
		if err := s.UpdateProfileCompletionField(userID, "has_phone_number", true); err != nil {
			return err
		}
		return s.UpdateProfileCompletionField(userID, "has_phone_verified", true)
	case "identity":
		return s.UpdateProfileCompletionField(userID, "has_identity_verified", true)
	case "address":
		return s.UpdateProfileCompletionField(userID, "has_address_verified", true)
	}
	
	return nil
}

// UpdateProfileCompletionFromActivity updates profile completion based on user activity
func (s *UserService) UpdateProfileCompletionFromActivity(userID uint, activityType string) error {
	// Update based on activity type
	switch activityType {
	case "forum_participation":
		return s.UpdateProfileCompletionField(userID, "has_participated_in_forum", true)
	case "bookmark_content":
		return s.UpdateProfileCompletionField(userID, "has_bookmarked_content", true)
	case "share_content":
		return s.UpdateProfileCompletionField(userID, "has_shared_content", true)
	case "join_group":
		return s.UpdateProfileCompletionField(userID, "has_joined_groups", true)
	case "set_preferences":
		return s.UpdateProfileCompletionField(userID, "has_set_preferences", true)
	case "earn_points":
		return s.UpdateProfileCompletionField(userID, "has_earned_points", true)
	}
	
	return nil
}

// SendProfileCompletionReminder checks if a reminder should be sent and updates the last reminder time
func (s *UserService) SendProfileCompletionReminder(userID uint) (bool, error) {
	status, err := s.GetProfileCompletionStatus(userID)
	if err != nil {
		return false, err
	}
	
	// Only send reminders for profiles under 80% complete
	if status.CompletionPercentage >= 80 {
		return false, nil
	}
	
	// Check if a reminder has been sent in the last 7 days
	if status.LastReminderSent != nil {
		lastReminderTime := *status.LastReminderSent
		if time.Since(lastReminderTime) < time.Hour*24*7 {
			return false, nil
		}
	}
	
	// Update the last reminder time
	now := time.Now()
	if err := s.db.Model(&models.ProfileCompletionStatus{}).
		Where("user_id = ?", userID).
		Update("last_reminder_sent", &now).Error; err != nil {
		return false, err
	}
	
	// Return true to indicate a reminder should be sent
	return true, nil
}