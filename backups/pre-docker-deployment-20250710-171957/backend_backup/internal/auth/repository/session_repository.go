package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// SessionRepository implements data access for user sessions
type SessionRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB, logger *logger.Logger) *SessionRepository {
	return &SessionRepository{
		db:     db,
		logger: logger,
	}
}

// GetSessionByID retrieves a session by its ID
func (r *SessionRepository) GetSessionByID(sessionID string) (*models.Session, error) {
	var session models.Session
	if err := r.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		r.logger.WithError(err).Error("Failed to get session by ID")
		return nil, err
	}
	return &session, nil
}

// GetSessionsByUserID retrieves all active sessions for a user
func (r *SessionRepository) GetSessionsByUserID(userID uint) ([]models.Session, error) {
	var sessions []models.Session
	if err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&sessions).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get sessions by user ID")
		return nil, err
	}
	return sessions, nil
}

// CreateSession creates a new session record
func (r *SessionRepository) CreateSession(session *models.Session) error {
	if err := r.db.Create(session).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create session")
		return err
	}
	return nil
}

// UpdateSession updates an existing session record
func (r *SessionRepository) UpdateSession(session *models.Session) error {
	if err := r.db.Save(session).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update session")
		return err
	}
	return nil
}

// DeleteSession marks a session as inactive
func (r *SessionRepository) DeleteSession(sessionID string) error {
	// Soft delete by marking as inactive
	if err := r.db.Model(&models.Session{}).Where("id = ?", sessionID).Update("is_active", false).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete session")
		return err
	}
	return nil
}

// DeleteSessionsByUserID marks all sessions for a user as inactive, except the specified one
func (r *SessionRepository) DeleteSessionsByUserID(userID uint, exceptSessionID string) (int, error) {
	var count int64
	result := r.db.Model(&models.Session{}).
		Where("user_id = ? AND id != ? AND is_active = ?", userID, exceptSessionID, true).
		Update("is_active", false)
	
	if result.Error != nil {
		r.logger.WithError(result.Error).Error("Failed to delete sessions by user ID")
		return 0, result.Error
	}
	
	return int(result.RowsAffected), nil
}

// DeleteExpiredSessions marks all expired sessions as inactive
func (r *SessionRepository) DeleteExpiredSessions() (int, error) {
	var count int64
	result := r.db.Model(&models.Session{}).
		Where("expires_at < ? AND is_active = ?", time.Now(), true).
		Update("is_active", false)
	
	if result.Error != nil {
		r.logger.WithError(result.Error).Error("Failed to delete expired sessions")
		return 0, result.Error
	}
	
	return int(result.RowsAffected), nil
}

// HardDeleteOldInactiveSessions permanently deletes inactive sessions older than the specified duration
func (r *SessionRepository) HardDeleteOldInactiveSessions(olderThan time.Duration) (int, error) {
	cutoffTime := time.Now().Add(-olderThan)
	
	var count int64
	result := r.db.Unscoped().
		Where("is_active = ? AND last_activity < ?", false, cutoffTime).
		Delete(&models.Session{})
	
	if result.Error != nil {
		r.logger.WithError(result.Error).Error("Failed to hard delete old inactive sessions")
		return 0, result.Error
	}
	
	return int(result.RowsAffected), nil
}