package repository

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "gorm.io/gorm"
)

// PointsRepository defines the interface for points-related data operations
type PointsRepository interface {
        // Points transaction operations
        CreatePointsTransaction(transaction *models.PointsTransaction) error
        GetPointsTransactionByID(id uint) (*models.PointsTransaction, error)
        GetPointsTransactionsByUserID(userID uint) ([]models.PointsTransaction, error)
        
        // User points balance operations
        GetUserPointsBalance(userID uint) (int, error)
        UpdateUserPoints(userID uint, points int) error

        // Award points to a user and create a transaction record
        AwardPoints(userID uint, sourceType string, sourceID uint, points int, description string) error
}

// GormPointsRepository implements PointsRepository using GORM
type GormPointsRepository struct {
        db *gorm.DB
}

// NewGormPointsRepository creates a new GormPointsRepository
func NewGormPointsRepository(db *gorm.DB) *GormPointsRepository {
        return &GormPointsRepository{db: db}
}

// CreatePointsTransaction creates a new points transaction in the database
func (r *GormPointsRepository) CreatePointsTransaction(transaction *models.PointsTransaction) error {
        return r.db.Create(transaction).Error
}

// GetPointsTransactionByID retrieves a points transaction by its ID
func (r *GormPointsRepository) GetPointsTransactionByID(id uint) (*models.PointsTransaction, error) {
        var transaction models.PointsTransaction
        if err := r.db.First(&transaction, id).Error; err != nil {
                return nil, err
        }
        return &transaction, nil
}

// GetPointsTransactionsByUserID retrieves all points transactions for a specific user
func (r *GormPointsRepository) GetPointsTransactionsByUserID(userID uint) ([]models.PointsTransaction, error) {
        var transactions []models.PointsTransaction
        if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error; err != nil {
                return nil, err
        }
        return transactions, nil
}

// GetUserPointsBalance calculates and returns the current points balance for a user
func (r *GormPointsRepository) GetUserPointsBalance(userID uint) (int, error) {
        var balance int
        var userProfile models.UserProfile
        
        // Try to retrieve the user profile first
        if err := r.db.Where("user_id = ?", userID).First(&userProfile).Error; err != nil {
                if err == gorm.ErrRecordNotFound {
                        // Create a new profile if it doesn't exist
                        userProfile = models.UserProfile{
                                UserID: userID,
                                Points: 0,
                        }
                        if err := r.db.Create(&userProfile).Error; err != nil {
                                return 0, err
                        }
                        return 0, nil
                }
                return 0, err
        }
        
        balance = userProfile.Points
        return balance, nil
}

// UpdateUserPoints updates the points balance for a user
func (r *GormPointsRepository) UpdateUserPoints(userID uint, points int) error {
        var userProfile models.UserProfile
        
        // Try to retrieve the user profile first
        if err := r.db.Where("user_id = ?", userID).First(&userProfile).Error; err != nil {
                if err == gorm.ErrRecordNotFound {
                        // Create a new profile if it doesn't exist
                        userProfile = models.UserProfile{
                                UserID: userID,
                                Points: points,
                        }
                        return r.db.Create(&userProfile).Error
                }
                return err
        }
        
        // Update the existing profile
        userProfile.Points = points
        return r.db.Save(&userProfile).Error
}

// AwardPoints awards points to a user and creates a transaction record
func (r *GormPointsRepository) AwardPoints(userID uint, sourceType string, sourceID uint, points int, description string) error {
        // Start a database transaction
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Get current points balance
                var userProfile models.UserProfile
                if err := tx.Where("user_id = ?", userID).First(&userProfile).Error; err != nil {
                        if err == gorm.ErrRecordNotFound {
                                // Create a new profile if it doesn't exist
                                userProfile = models.UserProfile{
                                        UserID: userID,
                                        Points: points, // Initial points
                                }
                                if err := tx.Create(&userProfile).Error; err != nil {
                                        return err
                                }
                        } else {
                                return err
                        }
                } else {
                        // Update existing profile
                        userProfile.Points += points
                        if err := tx.Save(&userProfile).Error; err != nil {
                                return err
                        }
                }
                
                // Create a transaction record
                transaction := models.PointsTransaction{
                        UserID:          userID,
                        Points:          points,
                        TransactionType: models.PointsEarned, // Assuming this is for earning points
                        ReferenceType:   sourceType,
                        ReferenceID:     &sourceID,
                        Description:     description,
                }
                
                return tx.Create(&transaction).Error
        })
}