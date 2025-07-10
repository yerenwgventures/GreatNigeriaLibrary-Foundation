package repository

import (
	"errors"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
	"gorm.io/gorm"
)

// TwoFARepository defines methods for two-factor authentication repository operations
type TwoFARepository interface {
	CreateTwoFA(userID uint, secret string) error
	UpdateTwoFA(userID uint, enabled, verified bool) error
	GetTwoFA(userID uint) (*models.TwoFactorAuth, error)
	StoreBackupCodes(userID uint, backupCodes []string) error
	ValidateBackupCode(userID uint, backupCode string) (bool, error)
	DisableTwoFA(userID uint) error
}

// TwoFARepositoryImpl implements TwoFARepository interface
type TwoFARepositoryImpl struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewTwoFARepository creates a new TwoFARepository instance
func NewTwoFARepository(db *gorm.DB, logger *logger.Logger) TwoFARepository {
	return &TwoFARepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// CreateTwoFA creates a new 2FA record for a user
func (r *TwoFARepositoryImpl) CreateTwoFA(userID uint, secret string) error {
	// Check if 2FA record already exists
	var twoFA models.TwoFactorAuth
	result := r.db.Where("user_id = ?", userID).First(&twoFA)
	
	if result.Error == nil {
		// Record exists, update it
		return r.db.Model(&twoFA).Updates(map[string]interface{}{
			"secret":   secret,
			"enabled":  false,
			"verified": false,
		}).Error
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new 2FA record
		twoFA = models.TwoFactorAuth{
			UserID:  userID,
			Secret:  secret,
			Enabled: false,
			Verified: false,
		}
		return r.db.Create(&twoFA).Error
	}
	
	return result.Error
}

// UpdateTwoFA updates the 2FA settings for a user
func (r *TwoFARepositoryImpl) UpdateTwoFA(userID uint, enabled, verified bool) error {
	return r.db.Model(&models.TwoFactorAuth{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"enabled":  enabled,
			"verified": verified,
		}).Error
}

// GetTwoFA gets the 2FA record for a user
func (r *TwoFARepositoryImpl) GetTwoFA(userID uint) (*models.TwoFactorAuth, error) {
	var twoFA models.TwoFactorAuth
	err := r.db.Where("user_id = ?", userID).First(&twoFA).Error
	if err != nil {
		return nil, err
	}
	return &twoFA, nil
}

// StoreBackupCodes stores the backup codes for a user
func (r *TwoFARepositoryImpl) StoreBackupCodes(userID uint, backupCodes []string) error {
	return r.db.Model(&models.TwoFactorAuth{}).
		Where("user_id = ?", userID).
		Update("backup_codes", backupCodes).Error
}

// ValidateBackupCode validates a backup code and removes it if valid
func (r *TwoFARepositoryImpl) ValidateBackupCode(userID uint, backupCode string) (bool, error) {
	var twoFA models.TwoFactorAuth
	if err := r.db.Where("user_id = ?", userID).First(&twoFA).Error; err != nil {
		return false, err
	}
	
	// Check if backupCode exists in the user's backup codes
	for i, code := range twoFA.BackupCodes {
		if code == backupCode {
			// Remove the used backup code
			twoFA.BackupCodes = append(twoFA.BackupCodes[:i], twoFA.BackupCodes[i+1:]...)
			if err := r.db.Save(&twoFA).Error; err != nil {
				return false, err
			}
			return true, nil
		}
	}
	
	return false, nil
}

// DisableTwoFA disables 2FA for a user
func (r *TwoFARepositoryImpl) DisableTwoFA(userID uint) error {
	return r.db.Model(&models.TwoFactorAuth{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"enabled":  false,
			"verified": false,
		}).Error
}