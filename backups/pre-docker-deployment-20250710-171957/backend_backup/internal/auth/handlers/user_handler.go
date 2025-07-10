package handlers

import (
        "net/http"
        "strconv"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/errors"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// UserHandler handles user-related requests
type UserHandler struct {
        userService UserService
        logger      *logger.Logger
}

// UserService defines the interface for user service operations
type UserService interface {
        // Authentication operations
        Register(req *models.UserRegisterRequest) (*models.UserResponse, *models.TokenPair, error)
        Login(req *models.UserLoginRequest) (*models.UserResponse, *models.TokenPair, error)
        RefreshToken(refreshToken string) (*models.TokenPair, error)
        RefreshTokenWithSession(refreshToken, clientIP, deviceInfo string) (*models.TokenPair, error)
        Logout(userID uint) error
        LogoutSession(userID uint, sessionID string) error
        LogoutAllSessions(userID uint) error
        
        // User data operations
        GetUserByID(id uint) (*models.User, error)
        UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
        GetUserProfile(id uint) (*models.UserProfile, error)
        ListUsers(page, pageSize int) ([]models.UserResponse, int64, error)
        
        // Password operations
        ResetPassword(email string) error
        ConfirmPasswordReset(req *models.PasswordResetConfirmRequest) error
        VerifyPassword(userID uint, password string) (bool, error)
        
        // OAuth operations
        OAuthLoginURL(provider string) (string, error)
        OAuthCallback(provider, code string) (*models.UserResponse, *models.TokenPair, error)
        
        // Email verification methods
        SendEmailVerification(email string) error
        VerifyEmail(token string) error
        ResendVerificationEmail(email string) error
        
        // User role operations
        UpdateUserRole(userID uint, role int) error
        GetUsersByRole(role int, page, pageSize int) ([]models.UserResponse, int64, error)
        
        // Account management
        DeleteUser(userID uint, password string) error
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService UserService, logger *logger.Logger) *UserHandler {
        return &UserHandler{
                userService: userService,
                logger:      logger,
        }
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
        var req models.UserRegisterRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        if !req.AcceptTerms {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("You must accept the terms and conditions"))
                return
        }

        // Get client IP address
        req.IP = c.ClientIP()
        
        // Get User-Agent header if device info not provided
        if req.DeviceInfo == "" {
                req.DeviceInfo = c.GetHeader("User-Agent")
        }

        user, tokens, err := h.userService.Register(&req)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to register user")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to register user"))
                return
        }

        c.JSON(http.StatusCreated, gin.H{
                "user":   user,
                "tokens": tokens,
        })
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
        var req models.UserLoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        // Get client IP address
        req.IP = c.ClientIP()
        
        // Get User-Agent header if device info not provided
        if req.DeviceInfo == "" {
                req.DeviceInfo = c.GetHeader("User-Agent")
        }

        user, tokens, err := h.userService.Login(&req)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to login user")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to login user"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "user":   user,
                "tokens": tokens,
        })
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
        var req struct {
                RefreshToken string `json:"refresh_token" binding:"required"`
                DeviceInfo   string `json:"device_info,omitempty"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        // Get client IP address and User-Agent if not provided
        clientIP := c.ClientIP()
        if req.DeviceInfo == "" {
                req.DeviceInfo = c.GetHeader("User-Agent")
        }

        tokens, err := h.userService.RefreshTokenWithSession(req.RefreshToken, clientIP, req.DeviceInfo)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to refresh token")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to refresh token"))
                return
        }

        c.JSON(http.StatusOK, tokens)
}

// GetUser gets a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid user ID"))
                return
        }

        // Check if user is requesting their own info or is an admin
        tokenUserID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        
        if uint(id) != tokenUserID.(uint) && !isAdmin.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("You don't have permission to access this user"))
                return
        }

        user, err := h.userService.GetUserByID(uint(id))
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to get user")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get user"))
                return
        }

        c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateUser handles updating a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid user ID"))
                return
        }

        // Check if user is updating their own info or is an admin
        tokenUserID, _ := c.Get("user_id")
        isAdmin, _ := c.Get("is_admin")
        
        if uint(id) != tokenUserID.(uint) && !isAdmin.(bool) {
                c.JSON(http.StatusForbidden, errors.ErrForbidden("You don't have permission to update this user"))
                return
        }

        var req models.UserUpdateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        user, err := h.userService.UpdateUser(uint(id), &req)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to update user")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to update user"))
                return
        }

        c.JSON(http.StatusOK, user)
}

// GetUserProfile gets a user's public profile
func (h *UserHandler) GetUserProfile(c *gin.Context) {
        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid user ID"))
                return
        }

        profile, err := h.userService.GetUserProfile(uint(id))
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to get user profile")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get user profile"))
                return
        }

        c.JSON(http.StatusOK, profile)
}

// ResetPassword handles password reset requests
func (h *UserHandler) ResetPassword(c *gin.Context) {
        var req struct {
                Email string `json:"email" binding:"required,email"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        err := h.userService.ResetPassword(req.Email)
        if err != nil {
                // Don't expose whether the email exists for security
                h.logger.WithError(err).WithFields(map[string]interface{}{
                        "email": req.Email,
                }).Info("Password reset requested")
                
                // Always return success even if the email doesn't exist
                c.JSON(http.StatusOK, gin.H{
                        "message": "If your email is registered, you will receive password reset instructions",
                })
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "If your email is registered, you will receive password reset instructions",
        })
}

// ConfirmPasswordReset handles password reset confirmation
func (h *UserHandler) ConfirmPasswordReset(c *gin.Context) {
        var req models.PasswordResetConfirmRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }

        err := h.userService.ConfirmPasswordReset(&req)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to reset password")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to reset password"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Password has been reset successfully",
        })
}

// Logout handles user logout
func (h *UserHandler) Logout(c *gin.Context) {
        userID, _ := c.Get("user_id")
        sessionID, exists := c.Get("session_id")
        
        var err error
        // If session ID exists, log out from that specific session only
        if exists && sessionID != nil {
                err = h.userService.LogoutSession(userID.(uint), sessionID.(string))
        } else {
                // If no session ID, log out from current device (based on token)
                err = h.userService.Logout(userID.(uint))
        }
        
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to logout user")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to logout user"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Successfully logged out",
        })
}

// OAuthLogin handles OAuth login
func (h *UserHandler) OAuthLogin(c *gin.Context) {
        provider := c.Param("provider")
        
        url, err := h.userService.OAuthLoginURL(provider)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to generate OAuth URL")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to generate OAuth URL"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "url": url,
        })
}

// OAuthCallback handles OAuth callback
func (h *UserHandler) OAuthCallback(c *gin.Context) {
        provider := c.Param("provider")
        code := c.Query("code")
        
        if code == "" {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Authorization code is missing"))
                return
        }

        user, tokens, err := h.userService.OAuthCallback(provider, code)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to process OAuth callback")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to process OAuth callback"))
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "user":   user,
                "tokens": tokens,
        })
}

// SendEmailVerification handles sending a verification email
func (h *UserHandler) SendEmailVerification(c *gin.Context) {
        var req models.EmailVerificationRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }
        
        err := h.userService.SendEmailVerification(req.Email)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to send verification email")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to send verification email"))
                return
        }
        
        // Always return success even if email doesn't exist for security reasons
        c.JSON(http.StatusOK, models.EmailVerificationResponse{
                Message: "If your email is registered, you will receive a verification email",
        })
}

// VerifyEmail handles verifying an email with a token
func (h *UserHandler) VerifyEmail(c *gin.Context) {
        var req models.EmailVerificationConfirmRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }
        
        err := h.userService.VerifyEmail(req.Token)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to verify email")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to verify email"))
                return
        }
        
        c.JSON(http.StatusOK, models.EmailVerificationResponse{
                Message: "Email verified successfully",
        })
}

// ResendVerificationEmail handles resending a verification email
func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
        var req models.EmailVerificationRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }
        
        err := h.userService.ResendVerificationEmail(req.Email)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to resend verification email")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to resend verification email"))
                return
        }
        
        // Always return success even if email doesn't exist for security reasons
        c.JSON(http.StatusOK, models.EmailVerificationResponse{
                Message: "If your email is registered, you will receive a verification email",
        })
}

// ListUsers handles listing users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
        // Parse pagination parameters
        pageStr := c.DefaultQuery("page", "1")
        pageSizeStr := c.DefaultQuery("page_size", "20")
        
        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
                page = 1
        }
        
        pageSize, err := strconv.Atoi(pageSizeStr)
        if err != nil || pageSize < 1 || pageSize > 100 {
                pageSize = 20
        }
        
        // Get users from service
        users, total, err := h.userService.ListUsers(page, pageSize)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to list users")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to list users"))
                return
        }
        
        // Return paginated response
        c.JSON(http.StatusOK, gin.H{
                "users": users,
                "pagination": gin.H{
                        "total": total,
                        "page": page,
                        "page_size": pageSize,
                        "total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
                },
        })
}

// UpdateUserRole handles updating a user's role
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
        // Get user ID from path
        idParam := c.Param("id")
        id, err := strconv.ParseUint(idParam, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid user ID"))
                return
        }
        
        // Parse role data
        var req struct {
                Role int `json:"role" binding:"required"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest(err.Error()))
                return
        }
        
        // Validate role value
        if req.Role < models.RoleBasicUser || (req.Role > models.RoleModerator && req.Role != models.RoleAdmin) {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid role value"))
                return
        }
        
        // Update user role
        err = h.userService.UpdateUserRole(uint(id), req.Role)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to update user role")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to update user role"))
                return
        }
        
        c.JSON(http.StatusOK, gin.H{
                "message": "User role updated successfully",
                "role": req.Role,
                "role_name": models.GetRoleNameByID(req.Role),
        })
}

// GetUsersByRole handles getting users by role with pagination
func (h *UserHandler) GetUsersByRole(c *gin.Context) {
        // Get role from path
        roleParam := c.Param("role")
        role, err := strconv.Atoi(roleParam)
        if err != nil {
                c.JSON(http.StatusBadRequest, errors.ErrInvalidRequest("Invalid role ID"))
                return
        }
        
        // Parse pagination parameters
        pageStr := c.DefaultQuery("page", "1")
        pageSizeStr := c.DefaultQuery("page_size", "20")
        
        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
                page = 1
        }
        
        pageSize, err := strconv.Atoi(pageSizeStr)
        if err != nil || pageSize < 1 || pageSize > 100 {
                pageSize = 20
        }
        
        // Get users by role
        users, total, err := h.userService.GetUsersByRole(role, page, pageSize)
        if err != nil {
                if e, ok := err.(*errors.APIError); ok {
                        c.JSON(e.Status, e)
                        return
                }
                h.logger.WithError(err).Error("Failed to get users by role")
                c.JSON(http.StatusInternalServerError, errors.ErrInternalServer("Failed to get users by role"))
                return
        }
        
        // Return paginated response
        c.JSON(http.StatusOK, gin.H{
                "role": role,
                "role_name": models.GetRoleNameByID(role),
                "users": users,
                "pagination": gin.H{
                        "total": total,
                        "page": page,
                        "page_size": pageSize,
                        "total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
                },
        })
}
