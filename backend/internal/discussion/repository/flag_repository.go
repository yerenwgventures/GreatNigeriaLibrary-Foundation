package repository

import (
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/models"
	"gorm.io/gorm"
)

// FlagRepository defines the interface for content flag operations
type FlagRepository interface {
	// Content flags
	CreateContentFlag(flag *models.ContentFlag) error
	GetFlagsByContent(contentType string, contentID uint) ([]models.ContentFlag, error)
	GetFlagsByStatus(status models.FlagStatus, page, pageSize int) ([]models.ContentFlag, error)
	GetFlagByID(id uint) (*models.ContentFlag, error)
	GetFlagsByUser(userID uint) ([]models.ContentFlag, error)
	UpdateFlag(flag *models.ContentFlag) error
	GetFlagCount(status models.FlagStatus) (int64, error)
	AssignFlag(flagID, moderatorID uint) error
	
	// Content moderation status
	CreateModerationStatus(status *models.ContentModerationStatus) error
	GetModerationStatusByContent(contentType string, contentID uint) (*models.ContentModerationStatus, error)
	UpdateModerationStatus(status *models.ContentModerationStatus) error
	GetPendingModerationCount() (int64, error)
	
	// User penalties
	CreateUserPenalty(penalty *models.UserPenalty) error
	GetUserPenalties(userID uint) ([]models.UserPenalty, error)
	GetActivePenalties(userID uint) ([]models.UserPenalty, error)
	UpdateUserPenalty(penalty *models.UserPenalty) error
	DeactivateUserPenalty(penaltyID uint) error
	GetPenaltyByID(id uint) (*models.UserPenalty, error)
}

// GormFlagRepository implements the FlagRepository interface
type GormFlagRepository struct {
	db *gorm.DB
}

// NewGormFlagRepository creates a new flag repository
func NewGormFlagRepository(db *gorm.DB) *GormFlagRepository {
	return &GormFlagRepository{db: db}
}

// CreateContentFlag creates a new content flag
func (r *GormFlagRepository) CreateContentFlag(flag *models.ContentFlag) error {
	return r.db.Create(flag).Error
}

// GetFlagsByContent retrieves flags for a content item
func (r *GormFlagRepository) GetFlagsByContent(contentType string, contentID uint) ([]models.ContentFlag, error) {
	var flags []models.ContentFlag
	result := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).Find(&flags)
	return flags, result.Error
}

// GetFlagsByStatus retrieves flags by status with pagination
func (r *GormFlagRepository) GetFlagsByStatus(status models.FlagStatus, page, pageSize int) ([]models.ContentFlag, error) {
	var flags []models.ContentFlag
	offset := (page - 1) * pageSize
	
	query := r.db
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	result := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&flags)
		
	return flags, result.Error
}

// GetFlagByID retrieves a flag by ID
func (r *GormFlagRepository) GetFlagByID(id uint) (*models.ContentFlag, error) {
	var flag models.ContentFlag
	result := r.db.First(&flag, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &flag, nil
}

// GetFlagsByUser retrieves flags created by a user
func (r *GormFlagRepository) GetFlagsByUser(userID uint) ([]models.ContentFlag, error) {
	var flags []models.ContentFlag
	result := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&flags)
	return flags, result.Error
}

// UpdateFlag updates a content flag
func (r *GormFlagRepository) UpdateFlag(flag *models.ContentFlag) error {
	return r.db.Save(flag).Error
}

// GetFlagCount gets the count of flags with a specific status
func (r *GormFlagRepository) GetFlagCount(status models.FlagStatus) (int64, error) {
	var count int64
	query := r.db.Model(&models.ContentFlag{})
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	result := query.Count(&count)
	return count, result.Error
}

// AssignFlag assigns a flag to a moderator
func (r *GormFlagRepository) AssignFlag(flagID, moderatorID uint) error {
	return r.db.Model(&models.ContentFlag{}).
		Where("id = ?", flagID).
		Updates(map[string]interface{}{
			"assigned_to": moderatorID,
			"updated_at":  time.Now(),
		}).Error
}

// CreateModerationStatus creates a content moderation status
func (r *GormFlagRepository) CreateModerationStatus(status *models.ContentModerationStatus) error {
	return r.db.Create(status).Error
}

// GetModerationStatusByContent retrieves the moderation status for a content item
func (r *GormFlagRepository) GetModerationStatusByContent(contentType string, contentID uint) (*models.ContentModerationStatus, error) {
	var status models.ContentModerationStatus
	result := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).First(&status)
	if result.Error != nil {
		return nil, result.Error
	}
	return &status, nil
}

// UpdateModerationStatus updates a content moderation status
func (r *GormFlagRepository) UpdateModerationStatus(status *models.ContentModerationStatus) error {
	return r.db.Save(status).Error
}

// GetPendingModerationCount gets the count of pending moderation items
func (r *GormFlagRepository) GetPendingModerationCount() (int64, error) {
	var count int64
	result := r.db.Model(&models.ContentModerationStatus{}).
		Where("status = ?", models.ModerationStatusPending).
		Count(&count)
	return count, result.Error
}

// CreateUserPenalty creates a user penalty
func (r *GormFlagRepository) CreateUserPenalty(penalty *models.UserPenalty) error {
	return r.db.Create(penalty).Error
}

// GetUserPenalties retrieves all penalties for a user
func (r *GormFlagRepository) GetUserPenalties(userID uint) ([]models.UserPenalty, error) {
	var penalties []models.UserPenalty
	result := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&penalties)
	return penalties, result.Error
}

// GetActivePenalties retrieves active penalties for a user
func (r *GormFlagRepository) GetActivePenalties(userID uint) ([]models.UserPenalty, error) {
	var penalties []models.UserPenalty
	now := time.Now()
	
	result := r.db.Where("user_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)", 
		userID, true, now).Find(&penalties)
		
	return penalties, result.Error
}

// UpdateUserPenalty updates a user penalty
func (r *GormFlagRepository) UpdateUserPenalty(penalty *models.UserPenalty) error {
	return r.db.Save(penalty).Error
}

// DeactivateUserPenalty deactivates a user penalty
func (r *GormFlagRepository) DeactivateUserPenalty(penaltyID uint) error {
	return r.db.Model(&models.UserPenalty{}).
		Where("id = ?", penaltyID).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

// GetPenaltyByID retrieves a penalty by ID
func (r *GormFlagRepository) GetPenaltyByID(id uint) (*models.UserPenalty, error) {
	var penalty models.UserPenalty
	result := r.db.First(&penalty, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &penalty, nil
}