package repository

import (
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
	"gorm.io/gorm"
)

// ContentAccessRepository defines the interface for content access repository operations
type ContentAccessRepository interface {
	// ContentAccess operations
	CreateContentAccess(contentAccess *models.ContentAccess) error
	GetContentAccess(contentType string, contentID uint) (*models.ContentAccess, error)
	UpdateContentAccess(contentAccess *models.ContentAccess) error
	DeleteContentAccess(id uint) error
	
	// ContentAccessRule operations
	CreateContentRule(rule *models.ContentAccessRule) error
	GetContentRule(id uint) (*models.ContentAccessRule, error)
	UpdateContentRule(rule *models.ContentAccessRule) error
	DeleteContentRule(id uint) error
	GetContentRules(contentType string, active bool) ([]models.ContentAccessRule, error)
	
	// UserContentPermission operations
	CreateUserPermission(permission *models.UserContentPermission) error
	GetUserPermissions(userID uint, contentType string) ([]models.UserContentPermission, error)
	UpdateUserPermission(permission *models.UserContentPermission) error
	DeleteUserPermission(id uint) error
	
	// UserPrivacySettings operations
	GetUserPrivacySettings(userID uint) (*models.UserPrivacySettings, error)
	UpdateUserPrivacySettings(settings *models.UserPrivacySettings) error
}

// GormContentAccessRepository is a GORM implementation of ContentAccessRepository
type GormContentAccessRepository struct {
	db *gorm.DB
}

// NewGormContentAccessRepository creates a new GormContentAccessRepository
func NewGormContentAccessRepository(db *gorm.DB) *GormContentAccessRepository {
	return &GormContentAccessRepository{db: db}
}

// CreateContentAccess creates a new content access record
func (r *GormContentAccessRepository) CreateContentAccess(contentAccess *models.ContentAccess) error {
	return r.db.Create(contentAccess).Error
}

// GetContentAccess gets a content access record by content type and ID
func (r *GormContentAccessRepository) GetContentAccess(contentType string, contentID uint) (*models.ContentAccess, error) {
	var contentAccess models.ContentAccess
	err := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).First(&contentAccess).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &contentAccess, nil
}

// UpdateContentAccess updates a content access record
func (r *GormContentAccessRepository) UpdateContentAccess(contentAccess *models.ContentAccess) error {
	return r.db.Save(contentAccess).Error
}

// DeleteContentAccess deletes a content access record
func (r *GormContentAccessRepository) DeleteContentAccess(id uint) error {
	return r.db.Delete(&models.ContentAccess{}, id).Error
}

// CreateContentRule creates a new content access rule
func (r *GormContentAccessRepository) CreateContentRule(rule *models.ContentAccessRule) error {
	return r.db.Create(rule).Error
}

// GetContentRule gets a content access rule by ID
func (r *GormContentAccessRepository) GetContentRule(id uint) (*models.ContentAccessRule, error) {
	var rule models.ContentAccessRule
	err := r.db.First(&rule, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

// UpdateContentRule updates a content access rule
func (r *GormContentAccessRepository) UpdateContentRule(rule *models.ContentAccessRule) error {
	return r.db.Save(rule).Error
}

// DeleteContentRule deletes a content access rule
func (r *GormContentAccessRepository) DeleteContentRule(id uint) error {
	return r.db.Delete(&models.ContentAccessRule{}, id).Error
}

// GetContentRules gets content access rules by content type and active status
func (r *GormContentAccessRepository) GetContentRules(contentType string, active bool) ([]models.ContentAccessRule, error) {
	var rules []models.ContentAccessRule
	query := r.db.Where("content_type = ?", contentType)
	if active {
		query = query.Where("is_active = ?", true)
	}
	
	// Order by priority (higher priority first)
	err := query.Order("priority DESC").Find(&rules).Error
	return rules, err
}

// CreateUserPermission creates a new user content permission
func (r *GormContentAccessRepository) CreateUserPermission(permission *models.UserContentPermission) error {
	return r.db.Create(permission).Error
}

// GetUserPermissions gets user content permissions by user ID and content type
func (r *GormContentAccessRepository) GetUserPermissions(userID uint, contentType string) ([]models.UserContentPermission, error) {
	var permissions []models.UserContentPermission
	query := r.db.Where("user_id = ?", userID)
	if contentType != "" {
		query = query.Where("content_type = ?", contentType)
	}
	err := query.Find(&permissions).Error
	return permissions, err
}

// UpdateUserPermission updates a user content permission
func (r *GormContentAccessRepository) UpdateUserPermission(permission *models.UserContentPermission) error {
	return r.db.Save(permission).Error
}

// DeleteUserPermission deletes a user content permission
func (r *GormContentAccessRepository) DeleteUserPermission(id uint) error {
	return r.db.Delete(&models.UserContentPermission{}, id).Error
}

// GetUserPrivacySettings gets user privacy settings by user ID
func (r *GormContentAccessRepository) GetUserPrivacySettings(userID uint) (*models.UserPrivacySettings, error) {
	var settings models.UserPrivacySettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If not found, create default settings
			settings = models.UserPrivacySettings{
				UserID:                userID,
				ProfileVisibility:     "public",
				ActivityVisibility:    "friends",
				ContactInfoVisibility: "private",
				SearchVisibility:      true,
			}
			err = r.db.Create(&settings).Error
			if err != nil {
				return nil, err
			}
			return &settings, nil
		}
		return nil, err
	}
	return &settings, nil
}

// UpdateUserPrivacySettings updates user privacy settings
func (r *GormContentAccessRepository) UpdateUserPrivacySettings(settings *models.UserPrivacySettings) error {
	return r.db.Save(settings).Error
}