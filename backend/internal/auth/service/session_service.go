package service

import (
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// SessionRepository defines the interface for session data operations
type SessionRepository interface {
	GetSessionByID(sessionID string) (*models.Session, error)
	GetSessionsByUserID(userID uint) ([]models.Session, error)
	CreateSession(session *models.Session) error
	UpdateSession(session *models.Session) error
	DeleteSession(sessionID string) error
	DeleteSessionsByUserID(userID uint, exceptSessionID string) (int, error)
	DeleteExpiredSessions() (int, error)
}

// SessionService implements session management functionality
type SessionService struct {
	sessionRepo    SessionRepository
	userRepo       UserRepository
	logger         *logger.Logger
	sessionTimeout time.Duration      // Default session timeout
	longLivedTime  time.Duration      // Long-lived session timeout (for "remember me")
}

// NewSessionService creates a new session service
func NewSessionService(
	sessionRepo SessionRepository,
	userRepo UserRepository,
	logger *logger.Logger,
) *SessionService {
	return &SessionService{
		sessionRepo:    sessionRepo,
		userRepo:       userRepo,
		logger:         logger,
		sessionTimeout: 24 * time.Hour,    // Default 24 hours
		longLivedTime:  30 * 24 * time.Hour, // 30 days for "remember me"
	}
}

// GetSessions retrieves all active sessions for a user
func (s *SessionService) GetSessions(userID uint) ([]models.SessionResponse, error) {
	// Check if user exists
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.ErrNotFound("User not found")
	}

	sessions, err := s.sessionRepo.GetSessionsByUserID(user.ID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get sessions")
		return nil, errors.ErrInternalServer("Failed to get sessions")
	}

	// Convert to response format
	var responseSessions []models.SessionResponse
	for _, session := range sessions {
		responseSessions = append(responseSessions, models.SessionResponse{
			ID:           session.ID,
			UserID:       session.UserID,
			DeviceInfo:   session.DeviceInfo,
			DeviceType:   session.DeviceType,
			IP:           session.IP,
			LastActivity: session.LastActivity,
			CreatedAt:    session.CreatedAt,
			ExpiresAt:    session.ExpiresAt,
			IsCurrentSession: false, // This will be set by the client
		})
	}

	return responseSessions, nil
}

// CreateSession creates a new session for a user
func (s *SessionService) CreateSession(userID uint, deviceInfo, deviceType, ip string, rememberMe bool) (*models.Session, error) {
	// Check if user exists
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.ErrNotFound("User not found")
	}

	// Determine session expiration based on remember me flag
	var expiresAt time.Time
	if rememberMe {
		expiresAt = time.Now().Add(s.longLivedTime)
	} else {
		expiresAt = time.Now().Add(s.sessionTimeout)
	}

	// Create new session
	session := &models.Session{
		ID:           models.GenerateUUID(), // Use a UUID for the session ID
		UserID:       user.ID,
		DeviceInfo:   deviceInfo,
		DeviceType:   deviceType,
		IP:           ip,
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
		ExpiresAt:    expiresAt,
		IsActive:     true,
	}

	// Save session
	if err := s.sessionRepo.CreateSession(session); err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to create session")
		return nil, errors.ErrInternalServer("Failed to create session")
	}

	return session, nil
}

// RefreshSession extends the expiration of an existing session
func (s *SessionService) RefreshSession(sessionID string, newIP, deviceInfo string) (*models.Session, error) {
	// Get session by ID
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return nil, errors.ErrNotFound("Session not found")
	}

	// Check if session is active
	if !session.IsActive {
		return nil, errors.ErrUnauthorized("Session is not active")
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.ErrUnauthorized("Session has expired")
	}

	// Determine new expiration time (preserve original duration)
	duration := session.ExpiresAt.Sub(session.CreatedAt)
	session.ExpiresAt = time.Now().Add(duration)
	session.LastActivity = time.Now()
	
	// Update IP and device info if provided
	if newIP != "" {
		session.IP = newIP
	}
	
	if deviceInfo != "" {
		session.DeviceInfo = deviceInfo
	}

	// Update session
	if err := s.sessionRepo.UpdateSession(session); err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to refresh session")
		return nil, errors.ErrInternalServer("Failed to refresh session")
	}

	return session, nil
}

// RevokeSession revokes a specific session for a user
func (s *SessionService) RevokeSession(userID uint, sessionID string) error {
	// Get session by ID
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return errors.ErrNotFound("Session not found")
	}

	// Verify the session belongs to the user
	if session.UserID != userID {
		return errors.ErrForbidden("You don't have permission to revoke this session")
	}

	// Delete the session
	if err := s.sessionRepo.DeleteSession(sessionID); err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to delete session")
		return errors.ErrInternalServer("Failed to revoke session")
	}

	return nil
}

// RevokeAllSessions revokes all sessions for a user except the current one
func (s *SessionService) RevokeAllSessions(userID uint, currentSessionID string) error {
	// Delete all sessions for user except current
	count, err := s.sessionRepo.DeleteSessionsByUserID(userID, currentSessionID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to delete sessions")
		return errors.ErrInternalServer("Failed to revoke sessions")
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":     userID,
		"sessions_revoked": count,
	}).Info("Revoked all sessions except current")

	return nil
}

// PerformMaintenance removes expired sessions (admin function)
func (s *SessionService) PerformMaintenance() (int, error) {
	count, err := s.sessionRepo.DeleteExpiredSessions()
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete expired sessions")
		return 0, errors.ErrInternalServer("Failed to perform session maintenance")
	}

	s.logger.WithField("expired_sessions_removed", count).Info("Session maintenance completed")
	return count, nil
}