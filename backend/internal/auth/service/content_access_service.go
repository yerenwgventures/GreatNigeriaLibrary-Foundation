package service

import (
	"github.com/sirupsen/logrus"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/auth/repository"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// ContentAccessService defines the interface for content access service operations
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

// ContentAccessServiceImpl implements ContentAccessService
type ContentAccessServiceImpl struct {
	contentRepo repository.ContentAccessRepository
	userRepo    repository.UserRepository
	logger      *logrus.Logger
}

// NewContentAccessService creates a new content access service
func NewContentAccessService(contentRepo repository.ContentAccessRepository, userRepo repository.UserRepository, logger *logrus.Logger) *ContentAccessServiceImpl {
	return &ContentAccessServiceImpl{
		contentRepo: contentRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

// SetContentAccess sets the access level for specific content
func (s *ContentAccessServiceImpl) SetContentAccess(contentType string, contentID uint, visibility models.ContentVisibility, minPoints int, isPremium bool) error {
	// First check if content access record exists
	contentAccess, err := s.contentRepo.GetContentAccess(contentType, contentID)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"content_type": contentType,
			"content_id":   contentID,
		}).Error("Failed to get content access")
		return errors.ErrInternalServer("Failed to check content access")
	}

	// If content access doesn't exist, create a new one
	if contentAccess == nil {
		contentAccess = &models.ContentAccess{
			ContentType:  contentType,
			ContentID:    contentID,
			Visibility:   visibility,
			MinPointsReq: minPoints,
			IsPremium:    isPremium,
		}
		err = s.contentRepo.CreateContentAccess(contentAccess)
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"content_type": contentType,
				"content_id":   contentID,
			}).Error("Failed to create content access")
			return errors.ErrInternalServer("Failed to set content access")
		}
	} else {
		// Update existing content access
		contentAccess.Visibility = visibility
		contentAccess.MinPointsReq = minPoints
		contentAccess.IsPremium = isPremium

		err = s.contentRepo.UpdateContentAccess(contentAccess)
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"content_type": contentType,
				"content_id":   contentID,
			}).Error("Failed to update content access")
			return errors.ErrInternalServer("Failed to update content access")
		}
	}

	return nil
}

// GetContentAccess gets the access settings for specific content
func (s *ContentAccessServiceImpl) GetContentAccess(contentType string, contentID uint) (*models.ContentAccess, error) {
	contentAccess, err := s.contentRepo.GetContentAccess(contentType, contentID)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"content_type": contentType,
			"content_id":   contentID,
		}).Error("Failed to get content access")
		return nil, errors.ErrInternalServer("Failed to get content access")
	}

	return contentAccess, nil
}

// CreateContentRule creates a new content access rule
func (s *ContentAccessServiceImpl) CreateContentRule(rule *models.ContentAccessRule, createdBy uint) error {
	rule.CreatedBy = createdBy
	err := s.contentRepo.CreateContentRule(rule)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"content_type": rule.ContentType,
			"created_by":   createdBy,
		}).Error("Failed to create content rule")
		return errors.ErrInternalServer("Failed to create content access rule")
	}

	return nil
}

// UpdateContentRule updates a content access rule
func (s *ContentAccessServiceImpl) UpdateContentRule(rule *models.ContentAccessRule, updatedBy uint) error {
	// First check if rule exists
	existingRule, err := s.contentRepo.GetContentRule(rule.ID)
	if err != nil {
		s.logger.WithError(err).WithField("rule_id", rule.ID).Error("Failed to get content rule")
		return errors.ErrInternalServer("Failed to check content rule")
	}

	if existingRule == nil {
		return errors.ErrNotFound("Content rule not found")
	}

	// Update rule fields
	existingRule.Name = rule.Name
	existingRule.Description = rule.Description
	existingRule.AppliesTo = rule.AppliesTo
	existingRule.Visibility = rule.Visibility
	existingRule.MinPointsReq = rule.MinPointsReq
	existingRule.MinMemberLevel = rule.MinMemberLevel
	existingRule.IsPremiumOnly = rule.IsPremiumOnly
	existingRule.IsModeratorOnly = rule.IsModeratorOnly
	existingRule.IsAdminOnly = rule.IsAdminOnly
	existingRule.IsActive = rule.IsActive
	existingRule.Priority = rule.Priority

	err = s.contentRepo.UpdateContentRule(existingRule)
	if err != nil {
		s.logger.WithError(err).WithField("rule_id", rule.ID).Error("Failed to update content rule")
		return errors.ErrInternalServer("Failed to update content access rule")
	}

	return nil
}

// DeleteContentRule deletes a content access rule
func (s *ContentAccessServiceImpl) DeleteContentRule(id uint, deletedBy uint) error {
	// First check if rule exists
	existingRule, err := s.contentRepo.GetContentRule(id)
	if err != nil {
		s.logger.WithError(err).WithField("rule_id", id).Error("Failed to get content rule")
		return errors.ErrInternalServer("Failed to check content rule")
	}

	if existingRule == nil {
		return errors.ErrNotFound("Content rule not found")
	}

	err = s.contentRepo.DeleteContentRule(id)
	if err != nil {
		s.logger.WithError(err).WithField("rule_id", id).Error("Failed to delete content rule")
		return errors.ErrInternalServer("Failed to delete content access rule")
	}

	return nil
}

// GetContentRules gets content access rules by content type
func (s *ContentAccessServiceImpl) GetContentRules(contentType string, includeInactive bool) ([]models.ContentAccessRule, error) {
	rules, err := s.contentRepo.GetContentRules(contentType, !includeInactive)
	if err != nil {
		s.logger.WithError(err).WithField("content_type", contentType).Error("Failed to get content rules")
		return nil, errors.ErrInternalServer("Failed to get content access rules")
	}

	return rules, nil
}

// GrantUserPermission grants a specific permission to a user
func (s *ContentAccessServiceImpl) GrantUserPermission(permission *models.UserContentPermission) error {
	// First check if user exists
	user, err := s.userRepo.GetByID(permission.UserID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", permission.UserID).Error("Failed to get user")
		return errors.ErrInternalServer("Failed to check user")
	}

	if user == nil {
		return errors.ErrNotFound("User not found")
	}

	// Create permission
	err = s.contentRepo.CreateUserPermission(permission)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":      permission.UserID,
			"content_type": permission.ContentType,
			"content_id":   permission.ContentID,
		}).Error("Failed to create user permission")
		return errors.ErrInternalServer("Failed to grant permission")
	}

	return nil
}

// RevokeUserPermission revokes a specific permission from a user
func (s *ContentAccessServiceImpl) RevokeUserPermission(id uint, revokedBy uint) error {
	err := s.contentRepo.DeleteUserPermission(id)
	if err != nil {
		s.logger.WithError(err).WithField("permission_id", id).Error("Failed to delete user permission")
		return errors.ErrInternalServer("Failed to revoke permission")
	}

	return nil
}

// GetUserPermissions gets all permissions for a user
func (s *ContentAccessServiceImpl) GetUserPermissions(userID uint, contentType string) ([]models.UserContentPermission, error) {
	permissions, err := s.contentRepo.GetUserPermissions(userID, contentType)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":      userID,
			"content_type": contentType,
		}).Error("Failed to get user permissions")
		return nil, errors.ErrInternalServer("Failed to get user permissions")
	}

	return permissions, nil
}

// GetUserPrivacySettings gets privacy settings for a user
func (s *ContentAccessServiceImpl) GetUserPrivacySettings(userID uint) (*models.UserPrivacySettings, error) {
	// First check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user")
		return nil, errors.ErrInternalServer("Failed to check user")
	}

	if user == nil {
		return nil, errors.ErrNotFound("User not found")
	}

	// Get or create privacy settings
	settings, err := s.contentRepo.GetUserPrivacySettings(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user privacy settings")
		return nil, errors.ErrInternalServer("Failed to get privacy settings")
	}

	return settings, nil
}

// UpdateUserPrivacySettings updates privacy settings for a user
func (s *ContentAccessServiceImpl) UpdateUserPrivacySettings(userID uint, settings *models.UserPrivacySettings) error {
	// First check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user")
		return errors.ErrInternalServer("Failed to check user")
	}

	if user == nil {
		return errors.ErrNotFound("User not found")
	}

	// Ensure the settings belong to the user
	settings.UserID = userID

	// Update settings
	err = s.contentRepo.UpdateUserPrivacySettings(settings)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to update user privacy settings")
		return errors.ErrInternalServer("Failed to update privacy settings")
	}

	return nil
}

// CheckContentAccess checks if a user has access to specific content
func (s *ContentAccessServiceImpl) CheckContentAccess(userID uint, contentType string, contentID uint) (bool, string) {
	// If no user ID provided, check for public access
	if userID == 0 {
		// Check if content is public
		contentAccess, err := s.contentRepo.GetContentAccess(contentType, contentID)
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"content_type": contentType,
				"content_id":   contentID,
			}).Error("Failed to get content access")
			return false, "Failed to check content access"
		}

		// If no access record or public visibility, allow access
		if contentAccess == nil || contentAccess.Visibility == models.ContentVisibilityPublic {
			return true, ""
		}

		return false, "You must be logged in to access this content"
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user")
		return false, "Failed to check user"
	}

	if user == nil {
		return false, "User not found"
	}

	// Check if user is active
	if !user.IsActive {
		return false, "Your account is inactive"
	}

	// For admin users, always grant access
	if user.IsUserAdmin() {
		return true, ""
	}

	// For moderators, grant access to most content
	if user.IsUserModerator() {
		// Check if content requires admin access
		contentAccess, err := s.contentRepo.GetContentAccess(contentType, contentID)
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"content_type": contentType,
				"content_id":   contentID,
			}).Error("Failed to get content access")
			return false, "Failed to check content access"
		}

		// If no specific access record, allow access to moderators
		if contentAccess == nil {
			return true, ""
		}

		// If content requires admin access, deny access to moderators
		if contentAccess.Visibility == models.ContentVisibilityAdmins {
			return false, "This content is only accessible to administrators"
		}

		return true, ""
	}

	// Check for user-specific permissions
	userPermissions, err := s.contentRepo.GetUserPermissions(userID, contentType)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":      userID,
			"content_type": contentType,
		}).Error("Failed to get user permissions")
		// Continue checking general access, don't fail because of permission error
	} else {
		// Check if user has a specific permission for this content
		for _, permission := range userPermissions {
			if permission.ContentID == contentID && permission.CanView {
				return true, ""
			}
		}
	}

	// Get content access settings
	contentAccess, err := s.contentRepo.GetContentAccess(contentType, contentID)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"content_type": contentType,
			"content_id":   contentID,
		}).Error("Failed to get content access")
		return false, "Failed to check content access"
	}

	// If no specific access record, allow access (default is public)
	if contentAccess == nil {
		return true, ""
	}

	// Check visibility level
	switch contentAccess.Visibility {
	case models.ContentVisibilityPublic:
		return true, ""

	case models.ContentVisibilityRegistered:
		// Already checked that user exists and is active
		return true, ""

	case models.ContentVisibilityEngaged:
		if user.IsUserEngaged() || user.PointsBalance >= contentAccess.MinPointsReq {
			return true, ""
		}
		return false, "You need to be an Engaged user to access this content"

	case models.ContentVisibilityActive:
		if user.IsUserActive() || user.PointsBalance >= contentAccess.MinPointsReq {
			return true, ""
		}
		return false, "You need to be an Active user to access this content"

	case models.ContentVisibilityPremium:
		if user.IsUserPremium() {
			return true, ""
		}
		return false, "This content requires a Premium membership"

	case models.ContentVisibilityModerators:
		return false, "This content is only accessible to moderators"

	case models.ContentVisibilityAdmins:
		return false, "This content is only accessible to administrators"

	default:
		// Unknown visibility level, default to allowing access
		return true, ""
	}
}
