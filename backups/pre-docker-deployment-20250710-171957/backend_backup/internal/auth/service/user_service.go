package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/auth"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/config"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// UserService implements user-related business logic
type UserService struct {
	userRepo       UserRepository
	sessionRepo    SessionRepository // Added session repository
	sessionService *SessionService   // Added session service
	jwtManager     *auth.JWTManager
	oauthManager   *auth.OAuthManager
	logger         *logger.Logger
	config         *config.Config
}

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	UpdateLastLogin(id uint) error
	UpdatePassword(id uint, hashedPassword string) error
	Delete(id uint) error
	DeleteUser(id uint) error
	VerifyPassword(id uint, providedPassword string) (bool, error)
	GetUserStats(userID uint) (int, int, error)
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(username string) (bool, error)
	GetUserWithStats(id uint) (*models.UserWithStats, error)
	ListUsers(page, pageSize int) ([]models.User, int64, error)

	// Password reset operations
	CreatePasswordResetToken(userID uint, token string, expiresAt time.Time) error
	GetPasswordResetToken(token string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenAsUsed(tokenID uint) error
	DeleteExpiredPasswordResetTokens() error

	// Email verification operations
	CreateEmailVerificationToken(userID uint, token string, expiresAt time.Time) error
	GetEmailVerificationToken(token string) (*models.EmailVerificationToken, error)
	MarkEmailVerificationTokenAsUsed(tokenID uint) error
	DeleteExpiredEmailVerificationTokens() error
	UpdateUserVerificationStatus(userID uint, isVerified bool) error

	// User role operations
	UpdateUserRole(userID uint, role int) error
	GetUsersByRole(role int, page, pageSize int) ([]models.User, int64, error)
}

// NewUserService creates a new user service
func NewUserService(userRepo UserRepository, logger *logger.Logger) *UserService {
	// Load configuration for JWT manager
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config for User Service: " + err.Error())
	}

	jwtManager := auth.NewJWTManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessTokenExpiration,
		cfg.Auth.RefreshTokenExpiration,
	)

	// Create OAuth manager
	oauthManager := auth.NewOAuthManager(logger, cfg)

	return &UserService{
		userRepo:     userRepo,
		jwtManager:   jwtManager,
		oauthManager: oauthManager,
		logger:       logger,
		config:       cfg,
	}
}

// SetSessionService allows setting the session service after initialization
// to avoid circular dependency issues
func (s *UserService) SetSessionService(sessionService *SessionService) {
	s.sessionService = sessionService
}

// SetSessionRepository allows setting the session repository after initialization
func (s *UserService) SetSessionRepository(sessionRepo SessionRepository) {
	s.sessionRepo = sessionRepo
}

// Register registers a new user
func (s *UserService) Register(req *models.UserRegisterRequest) (*models.UserResponse, *models.TokenPair, error) {
	s.logger.WithField("email", req.Email).Info("Registering new user")

	// Check if email exists
	emailExists, err := s.userRepo.CheckEmailExists(req.Email)
	if err != nil {
		s.logger.WithError(err).WithField("email", req.Email).Error("Error checking email availability")
		return nil, nil, errors.ErrInternalServer("Error checking email availability")
	}
	if emailExists {
		s.logger.WithField("email", req.Email).Info("Registration failed: Email already exists")
		return nil, nil, errors.ErrResourceExists("Email already exists")
	}

	// Check if username exists
	usernameExists, err := s.userRepo.CheckUsernameExists(req.Username)
	if err != nil {
		s.logger.WithError(err).WithField("username", req.Username).Error("Error checking username availability")
		return nil, nil, errors.ErrInternalServer("Error checking username availability")
	}
	if usernameExists {
		s.logger.WithField("username", req.Username).Info("Registration failed: Username already exists")
		return nil, nil, errors.ErrResourceExists("Username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return nil, nil, errors.ErrInternalServer("Failed to process registration")
	}

	// Create user
	user := &models.User{
		Username:        req.Username,
		Email:           req.Email,
		Password:        string(hashedPassword),
		FullName:        req.FullName,
		MembershipLevel: models.MembershipBasic,
		PointsBalance:   0,
		IsActive:        true,
		IsVerified:      false, // User starts as unverified
		LastLogin:       time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.WithError(err).WithField("username", req.Username).Error("Failed to create user")
		return nil, nil, errors.ErrInternalServer("Failed to process registration")
	}

	// Send verification email
	if err := s.SendEmailVerification(user.Email); err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to send verification email")
		// Non-critical error, continue with registration
	}

	// Check if we should use session-based authentication
	if s.sessionService != nil && req.DeviceInfo != "" {
		return s.registerWithSession(user, req)
	}

	// Generate regular tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate tokens")
		return nil, nil, errors.ErrInternalServer("Failed to generate authentication tokens")
	}

	s.logger.WithField("user_id", user.ID).Info("User registered successfully")
	return &user.ToResponse(), tokens, nil
}

// registerWithSession creates a session during registration and generates session-based tokens
func (s *UserService) registerWithSession(user *models.User, req *models.UserRegisterRequest) (*models.UserResponse, *models.TokenPair, error) {
	// Determine device type based on user agent or explicit device info
	deviceType := "unknown"
	if req.DeviceType != "" {
		deviceType = req.DeviceType
	} else if req.DeviceInfo != "" {
		// Try to determine device type from user agent
		uaLower := strings.ToLower(req.DeviceInfo)
		if strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "android") || strings.Contains(uaLower, "iphone") {
			deviceType = "mobile"
		} else if strings.Contains(uaLower, "tablet") || strings.Contains(uaLower, "ipad") {
			deviceType = "tablet"
		} else {
			deviceType = "desktop"
		}
	}

	// Create a new session
	session, err := s.sessionService.CreateSession(
		user.ID,
		req.DeviceInfo,
		deviceType,
		req.IP,
		req.RememberMe,
	)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to create session during registration")
		// Fall back to regular token generation
		tokens, tokenErr := s.jwtManager.GenerateTokenPair(user)
		if tokenErr != nil {
			s.logger.WithError(tokenErr).WithField("user_id", user.ID).Error("Failed to generate tokens")
			return nil, nil, errors.ErrInternalServer("Failed to generate authentication tokens")
		}
		return &user.ToResponse(), tokens, nil
	}

	// Generate tokens with session information
	tokens, err := s.jwtManager.GenerateTokenPairWithSession(user, session.ID, deviceType)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate tokens with session")
		return nil, nil, errors.ErrInternalServer("Failed to generate authentication tokens")
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":    user.ID,
		"session_id": session.ID,
		"device":     deviceType,
	}).Info("User registered successfully with session")

	return &user.ToResponse(), tokens, nil
}

// Login authenticates a user
func (s *UserService) Login(req *models.UserLoginRequest) (*models.UserResponse, *models.TokenPair, error) {
	s.logger.WithField("email", req.Email).Info("User login attempt")

	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		s.logger.WithError(err).WithField("email", req.Email).Error("Failed to get user by email")
		return nil, nil, errors.ErrInternalServer("Authentication error")
	}
	if user == nil {
		s.logger.WithField("email", req.Email).Info("Login failed: User not found")
		return nil, nil, errors.ErrUnauthorized("Invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		s.logger.WithField("user_id", user.ID).Info("Login failed: User account is deactivated")
		return nil, nil, errors.ErrUnauthorized("Your account has been deactivated")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.WithField("user_id", user.ID).Info("Login failed: Invalid password")
		return nil, nil, errors.ErrUnauthorized("Invalid email or password")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to update last login")
		// Non-critical error, continue
	}

	// Check if we should use session-based authentication
	if s.sessionService != nil {
		return s.loginWithSession(user, req)
	}

	// If session service is not available, fall back to regular authentication
	tokens, err := s.jwtManager.GenerateTokenPair(user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate tokens")
		return nil, nil, errors.ErrInternalServer("Failed to generate authentication tokens")
	}

	s.logger.WithField("user_id", user.ID).Info("User logged in successfully")
	return &user.ToResponse(), tokens, nil
}

// loginWithSession creates a session and generates tokens with session ID
func (s *UserService) loginWithSession(user *models.User, req *models.UserLoginRequest) (*models.UserResponse, *models.TokenPair, error) {
	// Determine device type based on user agent or explicit device info
	deviceType := "unknown"
	if req.DeviceType != "" {
		deviceType = req.DeviceType
	} else if req.DeviceInfo != "" {
		// Try to determine device type from user agent
		uaLower := strings.ToLower(req.DeviceInfo)
		if strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "android") || strings.Contains(uaLower, "iphone") {
			deviceType = "mobile"
		} else if strings.Contains(uaLower, "tablet") || strings.Contains(uaLower, "ipad") {
			deviceType = "tablet"
		} else {
			deviceType = "desktop"
		}
	}

	// Create a new session
	session, err := s.sessionService.CreateSession(
		user.ID,
		req.DeviceInfo,
		deviceType,
		req.IP,
		req.RememberMe,
	)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to create session")
		return nil, nil, errors.ErrInternalServer("Failed to create session")
	}

	// Generate tokens with session information
	tokens, err := s.jwtManager.GenerateTokenPairWithSession(user, session.ID, deviceType)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate tokens with session")
		return nil, nil, errors.ErrInternalServer("Failed to generate authentication tokens")
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":    user.ID,
		"session_id": session.ID,
		"device":     deviceType,
	}).Info("User logged in successfully with session")

	return &user.ToResponse(), tokens, nil
}

// RefreshToken refreshes a token pair
func (s *UserService) RefreshToken(refreshToken string) (*models.TokenPair, error) {
	s.logger.Info("Token refresh request")

	// Check if we should use session-based authentication
	if s.sessionService != nil && s.sessionRepo != nil {
		return s.refreshTokenWithSession(refreshToken)
	}

	// Fallback to regular token refresh
	return s.refreshTokenLegacy(refreshToken)
}

// RefreshTokenWithSession refreshes a token pair with updated session information
func (s *UserService) RefreshTokenWithSession(refreshToken, clientIP, deviceInfo string) (*models.TokenPair, error) {
	s.logger.Info("Token refresh request with session update")

	// First validate the refresh token
	userID, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.logger.WithError(err).Error("Invalid refresh token")
		return nil, errors.ErrUnauthorized("Invalid refresh token")
	}

	// Extract session ID from token
	sessionID, err := s.jwtManager.GetSessionIDFromRefreshToken(refreshToken)
	if err != nil {
		s.logger.WithError(err).Error("Failed to extract session ID from refresh token")
		return nil, errors.ErrUnauthorized("Invalid session")
	}

	// Validate the session
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to get session")
		return nil, errors.ErrInternalServer("Failed to validate session")
	}

	if session == nil || session.UserID != userID || !session.IsActive {
		s.logger.WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Info("Token refresh failed: Invalid or inactive session")
		return nil, errors.ErrUnauthorized("Invalid or expired session")
	}

	// Get user by ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user by ID")
		return nil, errors.ErrInternalServer("Failed to refresh token")
	}
	if user == nil {
		s.logger.WithField("user_id", userID).Info("Token refresh failed: User not found")
		return nil, errors.ErrUnauthorized("User not found")
	}

	// Check if user is active
	if !user.IsActive {
		s.logger.WithField("user_id", user.ID).Info("Token refresh failed: User account is deactivated")
		return nil, errors.ErrUnauthorized("Your account has been deactivated")
	}

	// Update session information if provided
	if clientIP != "" {
		session.LastIP = clientIP
	}

	// Update device info if provided
	if deviceInfo != "" {
		// Determine device type from user agent
		deviceType := session.DeviceType
		uaLower := strings.ToLower(deviceInfo)
		if strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "android") || strings.Contains(uaLower, "iphone") {
			deviceType = "mobile"
		} else if strings.Contains(uaLower, "tablet") || strings.Contains(uaLower, "ipad") {
			deviceType = "tablet"
		} else {
			deviceType = "desktop"
		}

		// Update if different
		if deviceType != session.DeviceType {
			session.DeviceType = deviceType
		}

		// Update device info if provided value is different
		if deviceInfo != session.DeviceInfo {
			session.DeviceInfo = deviceInfo
		}
	}

	// Update session last activity time
	session.LastActivity = time.Now()
	if err := s.sessionRepo.Update(session); err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to update session last activity")
		// Non-critical error, continue with token refresh
	}

	// Generate new tokens with session information
	tokens, err := s.jwtManager.GenerateTokenPairWithSession(user, session.ID, session.DeviceType)
	if err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to generate tokens with session")
		return nil, errors.ErrInternalServer("Failed to generate authentication tokens")
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":    user.ID,
		"session_id": sessionID,
	}).Info("Token refreshed successfully with session update")

	return tokens, nil
}

// refreshTokenLegacy handles token refresh without session tracking
func (s *UserService) refreshTokenLegacy(refreshToken string) (*models.TokenPair, error) {
	// Validate refresh token
	userID, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.logger.WithError(err).Error("Invalid refresh token")
		return nil, errors.ErrUnauthorized("Invalid refresh token")
	}

	// Get user by ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user by ID")
		return nil, errors.ErrInternalServer("Failed to refresh token")
	}
	if user == nil {
		s.logger.WithField("user_id", userID).Info("Token refresh failed: User not found")
		return nil, errors.ErrUnauthorized("User not found")
	}

	// Check if user is active
	if !user.IsActive {
		s.logger.WithField("user_id", user.ID).Info("Token refresh failed: User account is deactivated")
		return nil, errors.ErrUnauthorized("Your account has been deactivated")
	}

	// Generate new tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate tokens")
		return nil, errors.ErrInternalServer("Failed to generate authentication tokens")
	}

	s.logger.WithField("user_id", user.ID).Info("Token refreshed successfully")
	return tokens, nil
}

// refreshTokenWithSession handles token refresh with session validation
func (s *UserService) refreshTokenWithSession(refreshToken string) (*models.TokenPair, error) {
	// Validate refresh token
	userID, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.logger.WithError(err).Error("Invalid refresh token")
		return nil, errors.ErrUnauthorized("Invalid refresh token")
	}

	// Extract session ID from token
	sessionID, err := s.jwtManager.GetSessionIDFromRefreshToken(refreshToken)
	if err != nil {
		s.logger.WithError(err).Error("Failed to extract session ID from refresh token")
		return nil, errors.ErrUnauthorized("Invalid session")
	}

	// Validate the session
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to get session")
		return nil, errors.ErrInternalServer("Failed to validate session")
	}

	if session == nil || session.UserID != userID || !session.IsActive {
		s.logger.WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Info("Token refresh failed: Invalid or inactive session")
		return nil, errors.ErrUnauthorized("Invalid or expired session")
	}

	// Get user by ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user by ID")
		return nil, errors.ErrInternalServer("Failed to refresh token")
	}
	if user == nil {
		s.logger.WithField("user_id", userID).Info("Token refresh failed: User not found")
		return nil, errors.ErrUnauthorized("User not found")
	}

	// Check if user is active
	if !user.IsActive {
		s.logger.WithField("user_id", user.ID).Info("Token refresh failed: User account is deactivated")
		return nil, errors.ErrUnauthorized("Your account has been deactivated")
	}

	// Update session last activity time
	session.LastActivity = time.Now()
	if err := s.sessionRepo.Update(session); err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to update session last activity")
		// Non-critical error, continue with token refresh
	}

	// Generate new tokens with session information
	tokens, err := s.jwtManager.GenerateTokenPairWithSession(user, session.ID, session.DeviceType)
	if err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to generate tokens with session")
		return nil, errors.ErrInternalServer("Failed to generate authentication tokens")
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":    user.ID,
		"session_id": sessionID,
	}).Info("Token refreshed successfully with session")

	return tokens, nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	s.logger.WithField("id", id).Info("Getting user by ID")

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to get user by ID")
		return nil, errors.ErrInternalServer("Failed to get user")
	}
	if user == nil {
		s.logger.WithField("id", id).Info("User not found")
		return nil, errors.ErrNotFound("User not found")
	}
	return user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	s.logger.WithField("id", id).Info("Updating user")

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to get user by ID")
		return nil, errors.ErrInternalServer("Failed to update user")
	}
	if user == nil {
		s.logger.WithField("id", id).Info("Update failed: User not found")
		return nil, errors.ErrNotFound("User not found")
	}

	// Update fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}
	if req.Password != "" {
		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.WithError(err).WithField("id", id).Error("Failed to hash password")
			return nil, errors.ErrInternalServer("Failed to update password")
		}
		user.Password = string(hashedPassword)
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to update user")
		return nil, errors.ErrInternalServer("Failed to update user")
	}

	// Update profile completion status based on user updates
	// Don't fail the user update if this fails, just log the error
	if err := s.UpdateProfileCompletionFromUserUpdate(id, req); err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to update profile completion status")
		// Continue with the user update
	}

	s.logger.WithField("id", id).Info("User updated successfully")
	return &user.ToResponse(), nil
}

// GetUserProfile gets a user's public profile
func (s *UserService) GetUserProfile(id uint) (*models.UserProfile, error) {
	s.logger.WithField("id", id).Info("Getting user profile")

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to get user by ID")
		return nil, errors.ErrInternalServer("Failed to get user profile")
	}
	if user == nil {
		s.logger.WithField("id", id).Info("Profile retrieval failed: User not found")
		return nil, errors.ErrNotFound("User not found")
	}

	postCount, commentCount, err := s.userRepo.GetUserStats(id)
	if err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Failed to get user stats")
		return nil, errors.ErrInternalServer("Failed to get user profile")
	}

	profile := user.ToProfile()
	profile.PostCount = postCount
	profile.CommentCount = commentCount

	return &profile, nil
}

// ResetPassword initiates the password reset process
func (s *UserService) ResetPassword(email string) error {
	s.logger.WithField("email", email).Info("Password reset requested")

	// Check if email exists
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.logger.WithError(err).WithField("email", email).Error("Failed to get user by email")
		return errors.ErrInternalServer("Failed to process password reset")
	}
	if user == nil {
		// Return success even if user doesn't exist for security reasons
		s.logger.WithField("email", email).Info("Password reset requested for non-existent email")
		return nil
	}

	// Delete any existing reset tokens for this user
	s.userRepo.DeleteExpiredPasswordResetTokens()

	// Generate reset token
	resetToken := uuid.New().String()

	// Set expiration time (e.g., 24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Store the token in the database
	err = s.userRepo.CreatePasswordResetToken(user.ID, resetToken, expiresAt)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to create password reset token")
		return errors.ErrInternalServer("Failed to process password reset")
	}

	// Generate reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.config.App.FrontendURL, resetToken)

	// In a production environment, send an email with the reset link
	// For now, just log it
	s.logger.WithFields(map[string]interface{}{
		"email":       email,
		"user_id":     user.ID,
		"reset_token": resetToken,
		"reset_link":  resetLink,
	}).Info("Password reset token generated and ready to be sent")

	// TODO: Implement email sending functionality

	return nil
}

// ConfirmPasswordReset handles the confirmation of a password reset
func (s *UserService) ConfirmPasswordReset(req *models.PasswordResetConfirmRequest) error {
	s.logger.Info("Password reset confirmation requested")

	// Get the token from the database
	token, err := s.userRepo.GetPasswordResetToken(req.Token)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get password reset token")
		return errors.ErrInternalServer("Failed to process password reset confirmation")
	}

	// Check if token exists and is valid
	if token == nil {
		s.logger.Info("Invalid or expired password reset token")
		return errors.ErrUnauthorized("Invalid or expired password reset token")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash new password")
		return errors.ErrInternalServer("Failed to process password reset confirmation")
	}

	// Update the user's password
	err = s.userRepo.UpdatePassword(token.UserID, string(hashedPassword))
	if err != nil {
		s.logger.WithError(err).WithField("user_id", token.UserID).Error("Failed to update user password")
		return errors.ErrInternalServer("Failed to update password")
	}

	// Mark the token as used
	err = s.userRepo.MarkPasswordResetTokenAsUsed(token.ID)
	if err != nil {
		s.logger.WithError(err).WithField("token_id", token.ID).Error("Failed to mark token as used")
		// Non-critical error, continue
	}

	s.logger.WithField("user_id", token.UserID).Info("Password successfully reset")

	return nil
}

// SendEmailVerification sends a verification email to a user
func (s *UserService) SendEmailVerification(email string) error {
	s.logger.WithField("email", email).Info("Email verification requested")

	// Check if email exists
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.logger.WithError(err).WithField("email", email).Error("Failed to get user by email")
		return errors.ErrInternalServer("Failed to process email verification")
	}
	if user == nil {
		// Return success even if user doesn't exist for security reasons
		s.logger.WithField("email", email).Info("Email verification requested for non-existent email")
		return nil
	}

	// If user is already verified, do nothing
	if user.IsVerified {
		s.logger.WithField("user_id", user.ID).Info("User is already verified")
		return nil
	}

	// Delete any existing verification tokens for this user
	s.userRepo.DeleteExpiredEmailVerificationTokens()

	// Generate verification token
	verificationToken := uuid.New().String()

	// Set expiration time (e.g., 48 hours from now)
	expiresAt := time.Now().Add(48 * time.Hour)

	// Store the token in the database
	err = s.userRepo.CreateEmailVerificationToken(user.ID, verificationToken, expiresAt)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to create email verification token")
		return errors.ErrInternalServer("Failed to process email verification")
	}

	// Generate verification link
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", s.config.App.FrontendURL, verificationToken)

	// In a production environment, send an email with the verification link
	// For now, just log it
	s.logger.WithFields(map[string]interface{}{
		"email":              email,
		"user_id":            user.ID,
		"verification_token": verificationToken,
		"verification_link":  verificationLink,
	}).Info("Email verification token generated and ready to be sent")

	// TODO: Implement email sending functionality

	return nil
}

// VerifyEmail verifies a user's email using a verification token
func (s *UserService) VerifyEmail(token string) error {
	s.logger.Info("Email verification confirmation requested")

	// Get the token from the database
	verificationToken, err := s.userRepo.GetEmailVerificationToken(token)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get email verification token")
		return errors.ErrInternalServer("Failed to process email verification")
	}

	// Check if token exists and is valid
	if verificationToken == nil {
		s.logger.Info("Invalid or expired email verification token")
		return errors.ErrUnauthorized("Invalid or expired verification token")
	}

	// Update the user's verification status
	err = s.userRepo.UpdateUserVerificationStatus(verificationToken.UserID, true)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", verificationToken.UserID).Error("Failed to update user verification status")
		return errors.ErrInternalServer("Failed to verify email")
	}

	// Mark the token as used
	err = s.userRepo.MarkEmailVerificationTokenAsUsed(verificationToken.ID)
	if err != nil {
		s.logger.WithError(err).WithField("token_id", verificationToken.ID).Error("Failed to mark token as used")
		// Non-critical error, continue
	}

	s.logger.WithField("user_id", verificationToken.UserID).Info("Email successfully verified")

	return nil
}

// ResendVerificationEmail resends a verification email
func (s *UserService) ResendVerificationEmail(email string) error {
	s.logger.WithField("email", email).Info("Verification email resend requested")

	// Simply reuse the SendEmailVerification method
	return s.SendEmailVerification(email)
}

// UpdateUserRole updates a user's role
func (s *UserService) UpdateUserRole(userID uint, role int) error {
	s.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"role":    role,
	}).Info("Updating user role")

	// Validate role value
	if role < models.RoleBasicUser || (role > models.RoleModerator && role != models.RoleAdmin) {
		s.logger.WithField("role", role).Error("Invalid role value")
		return errors.ErrInvalidRequest("Invalid role value")
	}

	// Get user to verify they exist
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user by ID")
		return errors.ErrInternalServer("Failed to update user role")
	}
	if user == nil {
		s.logger.WithField("user_id", userID).Info("User not found")
		return errors.ErrNotFound("User not found")
	}

	// Update the role in the database
	err = s.userRepo.UpdateUserRole(userID, role)
	if err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id": userID,
			"role":    role,
		}).Error("Failed to update user role")
		return errors.ErrInternalServer("Failed to update user role")
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":   userID,
		"role":      role,
		"role_name": models.GetRoleNameByID(role),
	}).Info("User role updated successfully")

	return nil
}

// GetUsersByRole gets users by their role with pagination
func (s *UserService) GetUsersByRole(role int, page, pageSize int) ([]models.UserResponse, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"role":      role,
		"page":      page,
		"page_size": pageSize,
	}).Info("Getting users by role")

	// Validate role value
	if role < models.RoleBasicUser || (role > models.RoleModerator && role != models.RoleAdmin) {
		s.logger.WithField("role", role).Error("Invalid role value")
		return nil, 0, errors.ErrInvalidRequest("Invalid role value")
	}

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	// Get users from repository
	users, total, err := s.userRepo.GetUsersByRole(role, page, pageSize)
	if err != nil {
		s.logger.WithError(err).WithField("role", role).Error("Failed to get users by role")
		return nil, 0, errors.ErrInternalServer("Failed to get users")
	}

	// Convert to response objects
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return userResponses, total, nil
}

// Logout handles user logout
func (s *UserService) Logout(userID uint) error {
	s.logger.WithField("user_id", userID).Info("User logout")

	// In a real implementation, would blacklist the token in a Redis cache or similar
	// Could also update the user's last activity timestamp

	return nil
}

// OAuthLoginURL generates a URL for OAuth login
func (s *UserService) OAuthLoginURL(provider string) (string, error) {
	s.logger.WithField("provider", provider).Info("OAuth login URL requested")

	// Generate the OAuth URL using the OAuth manager
	authURL, err := s.oauthManager.GetAuthURL(provider)
	if err != nil {
		s.logger.WithError(err).WithField("provider", provider).Error("Failed to generate OAuth URL")
		return "", errors.NewAPIError("OAuth Error", "Failed to generate OAuth URL: "+err.Error(), http.StatusInternalServerError)
	}

	s.logger.WithFields(map[string]interface{}{
		"provider": provider,
		"auth_url": authURL,
	}).Info("OAuth login URL generated")

	return authURL, nil
}

// OAuthCallback handles OAuth callback
func (s *UserService) OAuthCallback(provider, code string) (*models.UserResponse, *models.TokenPair, error) {
	s.logger.WithFields(map[string]interface{}{
		"provider":    provider,
		"code_length": len(code),
	}).Info("OAuth callback received")

	// Exchange the authorization code for user information
	userInfo, err := s.oauthManager.ExchangeCodeAndGetUserInfo(provider, code)
	if err != nil {
		s.logger.WithError(err).WithField("provider", provider).Error("Failed to exchange code for user info")
		return nil, nil, errors.ErrInternalServer("Authentication failed")
	}

	// Check if a user with this email already exists
	existingUser, err := s.userRepo.GetByEmail(userInfo.Email)
	if err != nil {
		s.logger.WithError(err).WithField("email", userInfo.Email).Error("Failed to check for existing user")
		return nil, nil, errors.ErrInternalServer("Authentication failed")
	}

	var user *models.User

	if existingUser != nil {
		// User exists, check if they're already linked to this OAuth provider
		if existingUser.IsOAuth && existingUser.OAuthProvider == provider && existingUser.OAuthID == userInfo.ID {
			// User already linked, just update last login
			user = existingUser
			if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
				s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to update last login")
				// Non-critical error, continue
			}
		} else if existingUser.IsOAuth {
			// User exists but is linked to a different provider
			s.logger.WithField("email", userInfo.Email).Info("User already exists with a different OAuth provider")
			return nil, nil, errors.ErrResourceExists("Account already exists with a different provider")
		} else {
			// User exists with email/password auth, we could link the accounts here if desired
			// For now, just return an error
			s.logger.WithField("email", userInfo.Email).Info("User already exists with password auth")
			return nil, nil, errors.ErrResourceExists("Account already exists with this email")
		}
	} else {
		// Create a new user with OAuth info
		username := "user" + models.GenerateRandomString(8) // Generate a unique username
		user = &models.User{
			Username:        username,
			Email:           userInfo.Email,
			FullName:        userInfo.Name,
			ProfileImage:    userInfo.Picture,
			IsOAuth:         true,
			OAuthProvider:   provider,
			OAuthID:         userInfo.ID,
			IsActive:        true,
			LastLogin:       time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			MembershipLevel: models.MembershipBasic,
			PointsBalance:   0,
		}

		if err := s.userRepo.Create(user); err != nil {
			s.logger.WithError(err).WithField("email", userInfo.Email).Error("Failed to create user via OAuth")
			return nil, nil, errors.ErrInternalServer("Failed to create account")
		}

		s.logger.WithField("user_id", user.ID).Info("Created new user via OAuth")
	}

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate tokens")
		return nil, nil, errors.ErrInternalServer("Failed to generate authentication tokens")
	}

	s.logger.WithField("user_id", user.ID).Info("OAuth login successful")
	return &user.ToResponse(), tokens, nil
}

// ListUsers gets a paginated list of users
func (s *UserService) ListUsers(page, pageSize int) ([]models.UserResponse, int64, error) {
	s.logger.WithFields(map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}).Info("Listing users")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	users, total, err := s.userRepo.ListUsers(page, pageSize)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list users")
		return nil, 0, errors.ErrInternalServer("Failed to list users")
	}

	// Convert users to responses
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return userResponses, total, nil
}

// DeleteUser soft-deletes a user account
func (s *UserService) DeleteUser(userID uint, password string) error {
	s.logger.WithField("user_id", userID).Info("Deleting user account")

	// Verify password first as a security measure
	isValid, err := s.VerifyPassword(userID, password)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to verify password during account deletion")
		return err
	}

	if !isValid {
		s.logger.WithField("user_id", userID).Warn("Account deletion failed due to invalid password")
		return errors.NewAPIError("Unauthorized", "Invalid password", http.StatusUnauthorized)
	}

	// Proceed with account deletion
	if err := s.userRepo.DeleteUser(userID); err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to delete user account")
		return errors.ErrInternalServer("Failed to delete account")
	}

	s.logger.WithField("user_id", userID).Info("User account deleted successfully")
	return nil
}

// VerifyPassword checks if a password is correct for a user
func (s *UserService) VerifyPassword(userID uint, password string) (bool, error) {
	s.logger.WithField("user_id", userID).Info("Verifying user password")

	// Get user by ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user for password verification")
		return false, errors.ErrInternalServer("Failed to verify password")
	}
	if user == nil {
		s.logger.WithField("user_id", userID).Error("User not found for password verification")
		return false, errors.ErrNotFound("User not found")
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Password doesn't match
		return false, nil
	}

	return true, nil
}

// Logout handles user logout
func (s *UserService) Logout(userID uint, refreshToken string) error {
	s.logger.WithField("user_id", userID).Info("Processing logout request")

	// Check if we should use session-based authentication
	if s.sessionService != nil && s.sessionRepo != nil && refreshToken != "" {
		return s.logoutWithSession(userID, refreshToken)
	}

	// Legacy logout has no specific actions
	s.logger.WithField("user_id", userID).Info("User logged out successfully (legacy)")
	return nil
}

// LogoutSession logs out a specific session by its ID
func (s *UserService) LogoutSession(userID uint, sessionID string) error {
	s.logger.WithFields(map[string]interface{}{
		"user_id":    userID,
		"session_id": sessionID,
	}).Info("Processing logout for specific session")

	// Check if we have session service available
	if s.sessionService == nil {
		return errors.ErrInternalServer("Session management is not available")
	}

	// Verify session belongs to the user
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to get session")
		// If we can't get session, consider it an error
		return errors.ErrInternalServer("Failed to validate session")
	}

	// Check if session exists and belongs to the user
	if session == nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Warn("Session not found")
		return errors.ErrNotFound("Session not found")
	}

	if session.UserID != userID {
		s.logger.WithFields(map[string]interface{}{
			"user_id":         userID,
			"session_id":      sessionID,
			"session_user_id": session.UserID,
		}).Warn("Session does not belong to the user attempting to log out")
		return errors.ErrForbidden("You don't have permission to end this session")
	}

	// End the session
	err = s.sessionService.EndSession(sessionID)
	if err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to end session")
		return errors.ErrInternalServer("Failed to end session")
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":    userID,
		"session_id": sessionID,
	}).Info("Session terminated successfully")

	return nil
}

// logoutWithSession handles session termination during logout
func (s *UserService) logoutWithSession(userID uint, refreshToken string) error {
	// Extract session ID from token
	sessionID, err := s.jwtManager.GetSessionIDFromRefreshToken(refreshToken)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to extract session ID from refresh token")
		// If we can't get session ID, just return success
		return nil
	}

	// Verify session belongs to the user
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		s.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Error("Failed to get session")
		// If we can't get session, just return success
		return nil
	}

	// If session is found and belongs to user, terminate it
	if session != nil && session.UserID == userID {
		err = s.sessionService.EndSession(sessionID)
		if err != nil {
			s.logger.WithError(err).WithFields(map[string]interface{}{
				"user_id":    userID,
				"session_id": sessionID,
			}).Error("Failed to end session")
			// Non-critical error, return success anyway
			return nil
		}

		s.logger.WithFields(map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		}).Info("User logged out successfully with session terminated")
	} else if session != nil {
		s.logger.WithFields(map[string]interface{}{
			"user_id":         userID,
			"session_id":      sessionID,
			"session_user_id": session.UserID,
		}).Warn("Session does not belong to the user attempting to log out")
	}

	return nil
}

// LogoutAllSessions terminates all active sessions for a user
func (s *UserService) LogoutAllSessions(userID uint) error {
	s.logger.WithField("user_id", userID).Info("Processing logout all sessions request")

	// Check if we have session service available
	if s.sessionService == nil {
		return errors.ErrInternalServer("Session management is not available")
	}

	// End all sessions for the user
	err := s.sessionService.EndAllUserSessions(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to end all user sessions")
		return errors.ErrInternalServer("Failed to log out from all sessions")
	}

	s.logger.WithField("user_id", userID).Info("All user sessions terminated successfully")
	return nil
}
