package repository

import (
        "errors"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
        "github.com/sirupsen/logrus"
        "gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
        db     *gorm.DB
        logger *logrus.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB, logger *logrus.Logger) *UserRepository {
        return &UserRepository{
                db:     db,
                logger: logger,
        }
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
        r.logger.WithField("username", user.Username).Info("Creating new user")
        return r.db.Create(user).Error
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
        var user models.User
        if err := r.db.First(&user, id).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        r.logger.WithField("id", id).Info("User not found")
                        return nil, nil
                }
                r.logger.WithError(err).WithField("id", id).Error("Error getting user by ID")
                return nil, err
        }
        return &user, nil
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
        var user models.User
        if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        r.logger.WithField("email", email).Info("User not found by email")
                        return nil, nil
                }
                r.logger.WithError(err).WithField("email", email).Error("Error getting user by email")
                return nil, err
        }
        return &user, nil
}

// GetByUsername gets a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
        var user models.User
        if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        r.logger.WithField("username", username).Info("User not found by username")
                        return nil, nil
                }
                r.logger.WithError(err).WithField("username", username).Error("Error getting user by username")
                return nil, err
        }
        return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
        r.logger.WithField("id", user.ID).Info("Updating user")
        return r.db.Save(user).Error
}

// UpdateLastLogin updates a user's last login time
func (r *UserRepository) UpdateLastLogin(id uint) error {
        r.logger.WithField("id", id).Info("Updating user last login time")
        return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", time.Now()).Error
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
        r.logger.WithField("id", id).Info("Updating user password")
        return r.db.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(id uint) error {
        r.logger.WithField("id", id).Info("Deleting user")
        return r.db.Delete(&models.User{}, id).Error
}

// GetUserStats gets a user's stats (post and comment count)
func (r *UserRepository) GetUserStats(userID uint) (int, int, error) {
        var postCount int64
        var commentCount int64

        if err := r.db.Model(&models.Discussion{}).Where("user_id = ?", userID).Count(&postCount).Error; err != nil {
                r.logger.WithError(err).WithField("user_id", userID).Error("Error getting post count")
                return 0, 0, err
        }

        if err := r.db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&commentCount).Error; err != nil {
                r.logger.WithError(err).WithField("user_id", userID).Error("Error getting comment count")
                return 0, 0, err
        }

        return int(postCount), int(commentCount), nil
}

// CheckEmailExists checks if an email exists
func (r *UserRepository) CheckEmailExists(email string) (bool, error) {
        var count int64
        if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
                r.logger.WithError(err).WithField("email", email).Error("Error checking if email exists")
                return false, err
        }
        return count > 0, nil
}

// CheckUsernameExists checks if a username exists
func (r *UserRepository) CheckUsernameExists(username string) (bool, error) {
        var count int64
        if err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
                r.logger.WithError(err).WithField("username", username).Error("Error checking if username exists")
                return false, err
        }
        return count > 0, nil
}

// GetUserWithStats gets a user with their stats
func (r *UserRepository) GetUserWithStats(id uint) (*models.UserWithStats, error) {
        var user models.User
        if err := r.db.First(&user, id).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        r.logger.WithField("id", id).Info("User not found")
                        return nil, nil
                }
                r.logger.WithError(err).WithField("id", id).Error("Error getting user with stats")
                return nil, err
        }

        postCount, commentCount, err := r.GetUserStats(id)
        if err != nil {
                return nil, err
        }

        var bookmarksCount int64
        if err := r.db.Model(&models.Bookmark{}).Where("user_id = ?", id).Count(&bookmarksCount).Error; err != nil {
                r.logger.WithError(err).WithField("user_id", id).Error("Error getting bookmarks count")
                return nil, err
        }

        userWithStats := &models.UserWithStats{
                User:           user,
                PostCount:      postCount,
                CommentCount:   commentCount,
                BookmarksCount: int(bookmarksCount),
        }

        return userWithStats, nil
}

// ListUsers gets a list of users with pagination
func (r *UserRepository) ListUsers(page, pageSize int) ([]models.User, int64, error) {
        var users []models.User
        var total int64

        offset := (page - 1) * pageSize

        // Get total count
        if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
                r.logger.WithError(err).Error("Error counting users")
                return nil, 0, err
        }

        // Get users with pagination
        if err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
                r.logger.WithError(err).Error("Error getting users with pagination")
                return nil, 0, err
        }

        return users, total, nil
}

// CreatePasswordResetToken creates a password reset token for a user
func (r *UserRepository) CreatePasswordResetToken(userID uint, token string, expiresAt time.Time) error {
        r.logger.WithField("user_id", userID).Info("Creating password reset token")
        
        passwordResetToken := models.PasswordResetToken{
                UserID:    userID,
                Token:     token,
                ExpiresAt: expiresAt,
                Used:      false,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
        }
        
        return r.db.Create(&passwordResetToken).Error
}

// GetPasswordResetToken gets a password reset token by token string
func (r *UserRepository) GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
        var passwordResetToken models.PasswordResetToken
        
        if err := r.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&passwordResetToken).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        r.logger.WithField("token", token).Info("Password reset token not found or expired")
                        return nil, nil
                }
                r.logger.WithError(err).WithField("token", token).Error("Error getting password reset token")
                return nil, err
        }
        
        return &passwordResetToken, nil
}

// MarkPasswordResetTokenAsUsed marks a password reset token as used
func (r *UserRepository) MarkPasswordResetTokenAsUsed(tokenID uint) error {
        r.logger.WithField("token_id", tokenID).Info("Marking password reset token as used")
        
        return r.db.Model(&models.PasswordResetToken{}).Where("id = ?", tokenID).Update("used", true).Error
}

// DeleteExpiredPasswordResetTokens deletes expired password reset tokens
func (r *UserRepository) DeleteExpiredPasswordResetTokens() error {
        r.logger.Info("Deleting expired password reset tokens")
        
        return r.db.Where("expires_at < ? OR used = ?", time.Now(), true).Delete(&models.PasswordResetToken{}).Error
}

// CreateEmailVerificationToken creates an email verification token for a user
func (r *UserRepository) CreateEmailVerificationToken(userID uint, token string, expiresAt time.Time) error {
        r.logger.WithField("user_id", userID).Info("Creating email verification token")
        
        emailVerificationToken := models.EmailVerificationToken{
                UserID:    userID,
                Token:     token,
                ExpiresAt: expiresAt,
                Used:      false,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
        }
        
        return r.db.Create(&emailVerificationToken).Error
}

// GetEmailVerificationToken gets an email verification token by token string
func (r *UserRepository) GetEmailVerificationToken(token string) (*models.EmailVerificationToken, error) {
        var emailVerificationToken models.EmailVerificationToken
        
        if err := r.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&emailVerificationToken).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        r.logger.WithField("token", token).Info("Email verification token not found or expired")
                        return nil, nil
                }
                r.logger.WithError(err).WithField("token", token).Error("Error getting email verification token")
                return nil, err
        }
        
        return &emailVerificationToken, nil
}

// MarkEmailVerificationTokenAsUsed marks an email verification token as used
func (r *UserRepository) MarkEmailVerificationTokenAsUsed(tokenID uint) error {
        r.logger.WithField("token_id", tokenID).Info("Marking email verification token as used")
        
        return r.db.Model(&models.EmailVerificationToken{}).Where("id = ?", tokenID).Update("used", true).Error
}

// DeleteExpiredEmailVerificationTokens deletes expired email verification tokens
func (r *UserRepository) DeleteExpiredEmailVerificationTokens() error {
        r.logger.Info("Deleting expired email verification tokens")
        
        return r.db.Where("expires_at < ? OR used = ?", time.Now(), true).Delete(&models.EmailVerificationToken{}).Error
}

// UpdateUserVerificationStatus updates a user's verification status
func (r *UserRepository) UpdateUserVerificationStatus(userID uint, isVerified bool) error {
        r.logger.WithFields(map[string]interface{}{
                "user_id": userID,
                "status": isVerified,
        }).Info("Updating user verification status")
        
        return r.db.Model(&models.User{}).Where("id = ?", userID).Update("is_verified", isVerified).Error
}

// UpdateUserRole updates a user's role
func (r *UserRepository) UpdateUserRole(userID uint, role int) error {
        r.logger.WithFields(map[string]interface{}{
                "user_id": userID,
                "role": role,
        }).Info("Updating user role")
        
        // If changing to admin role, also update legacy IsAdmin flag
        if role >= models.RoleAdmin {
                return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
                        "role": role,
                        "is_admin": true,
                }).Error
        }
        
        return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

// GetUsersByRole gets users by their role
func (r *UserRepository) GetUsersByRole(role int, page, pageSize int) ([]models.User, int64, error) {
        var users []models.User
        var total int64
        
        offset := (page - 1) * pageSize
        
        // Get total count by role
        if err := r.db.Model(&models.User{}).Where("role = ?", role).Count(&total).Error; err != nil {
                r.logger.WithError(err).WithField("role", role).Error("Error counting users by role")
                return nil, 0, err
        }
        
        // Get users with pagination
        if err := r.db.Where("role = ?", role).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
                r.logger.WithError(err).WithField("role", role).Error("Error getting users by role with pagination")
                return nil, 0, err
        }
        
        return users, total, nil
}

// DeleteUser soft-deletes a user by ID
func (r *UserRepository) DeleteUser(id uint) error {
        r.logger.WithField("user_id", id).Info("Soft-deleting user account")
        
        // Use transaction to ensure all related user data is soft-deleted
        return r.db.Transaction(func(tx *gorm.DB) error {
                // First check if user exists
                var user models.User
                if err := tx.First(&user, id).Error; err != nil {
                        if errors.Is(err, gorm.ErrRecordNotFound) {
                                return err
                        }
                        r.logger.WithError(err).WithField("user_id", id).Error("Error finding user for deletion")
                        return err
                }
                
                // Soft delete user (GORM will set DeletedAt)
                if err := tx.Delete(&user).Error; err != nil {
                        r.logger.WithError(err).WithField("user_id", id).Error("Error soft-deleting user")
                        return err
                }
                
                // Could add additional cleanup operations here if needed in the future
                // For example, cascade soft-delete user's comments, posts, etc.
                
                return nil
        })
}

// VerifyPassword checks if the provided password matches the stored hash for a user
func (r *UserRepository) VerifyPassword(id uint, providedPassword string) (bool, error) {
        var user models.User
        if err := r.db.Select("id, password").First(&user, id).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return false, errors.New("user not found")
                }
                r.logger.WithError(err).WithField("user_id", id).Error("Error finding user for password verification")
                return false, err
        }
        
        // Password verification will be done in the service layer
        return true, nil
}
