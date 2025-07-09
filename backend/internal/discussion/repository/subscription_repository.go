package repository

import (
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/models"
        "gorm.io/gorm"
)

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
        CreateSubscription(subscription *models.AdvancedSubscription) error
        GetSubscription(id uint) (*models.AdvancedSubscription, error)
        GetUserSubscriptions(userID uint) ([]models.AdvancedSubscription, error)
        GetSubscriptionsByType(userID uint, subscriptionType models.SubscriptionType) ([]models.AdvancedSubscription, error)
        GetSubscriptionByReference(userID uint, refType models.SubscriptionType, refID uint) (*models.AdvancedSubscription, error)
        UpdateSubscription(subscription *models.AdvancedSubscription) error
        DeleteSubscription(id uint) error
        
        CreateSubscriptionPreference(preference *models.SubscriptionPreference) error
        GetSubscriptionPreference(userID uint) (*models.SubscriptionPreference, error)
        UpdateSubscriptionPreference(preference *models.SubscriptionPreference) error
        
        CreateSubscriptionDigest(digest *models.SubscriptionDigest) error
        GetPendingDigests(frequency models.SubscriptionFrequency, before time.Time) ([]models.SubscriptionDigest, error)
        UpdateDigestStatus(id uint, status string, sentAt *time.Time, errorMessage string) error
        GetUserDigests(userID uint, limit int, offset int) ([]models.SubscriptionDigest, error)
}

// GormSubscriptionRepository implements the SubscriptionRepository interface
type GormSubscriptionRepository struct {
        db *gorm.DB
}

// NewGormSubscriptionRepository creates a new subscription repository
func NewGormSubscriptionRepository(db *gorm.DB) *GormSubscriptionRepository {
        return &GormSubscriptionRepository{db: db}
}

// CreateSubscription creates a new subscription
func (r *GormSubscriptionRepository) CreateSubscription(subscription *models.AdvancedSubscription) error {
        return r.db.Create(subscription).Error
}

// GetSubscription retrieves a subscription by ID
func (r *GormSubscriptionRepository) GetSubscription(id uint) (*models.AdvancedSubscription, error) {
        var subscription models.AdvancedSubscription
        result := r.db.First(&subscription, id)
        if result.Error != nil {
                return nil, result.Error
        }
        return &subscription, nil
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (r *GormSubscriptionRepository) GetUserSubscriptions(userID uint) ([]models.AdvancedSubscription, error) {
        var subscriptions []models.AdvancedSubscription
        result := r.db.Where("user_id = ?", userID).Find(&subscriptions)
        return subscriptions, result.Error
}

// GetSubscriptionsByType retrieves subscriptions of a specific type for a user
func (r *GormSubscriptionRepository) GetSubscriptionsByType(userID uint, subscriptionType models.SubscriptionType) ([]models.AdvancedSubscription, error) {
        var subscriptions []models.AdvancedSubscription
        result := r.db.Where("user_id = ? AND type = ?", userID, subscriptionType).Find(&subscriptions)
        return subscriptions, result.Error
}

// GetSubscriptionByReference retrieves a subscription by reference
func (r *GormSubscriptionRepository) GetSubscriptionByReference(userID uint, refType models.SubscriptionType, refID uint) (*models.AdvancedSubscription, error) {
        var subscription models.AdvancedSubscription
        result := r.db.Where("user_id = ? AND type = ? AND reference_id = ?", userID, refType, refID).First(&subscription)
        if result.Error != nil {
                return nil, result.Error
        }
        return &subscription, nil
}

// UpdateSubscription updates a subscription
func (r *GormSubscriptionRepository) UpdateSubscription(subscription *models.AdvancedSubscription) error {
        return r.db.Save(subscription).Error
}

// DeleteSubscription deletes a subscription
func (r *GormSubscriptionRepository) DeleteSubscription(id uint) error {
        return r.db.Delete(&models.AdvancedSubscription{}, id).Error
}

// CreateSubscriptionPreference creates a subscription preference
func (r *GormSubscriptionRepository) CreateSubscriptionPreference(preference *models.SubscriptionPreference) error {
        return r.db.Create(preference).Error
}

// GetSubscriptionPreference retrieves subscription preferences for a user
func (r *GormSubscriptionRepository) GetSubscriptionPreference(userID uint) (*models.SubscriptionPreference, error) {
        var preference models.SubscriptionPreference
        result := r.db.Where("user_id = ?", userID).First(&preference)
        if result.Error != nil {
                // If not found, create default preferences
                if result.Error == gorm.ErrRecordNotFound {
                        preference = models.SubscriptionPreference{
                                UserID:                 userID,
                                DefaultFrequency:       models.FrequencyInstant,
                                DefaultEmailEnabled:    true,
                                DefaultPushEnabled:     true,
                                DefaultInAppEnabled:    true,
                                DigestDay:              0, // Sunday
                                DigestHour:             8, // 8 AM
                                AutoSubscribeToReplies: true,
                                AutoSubscribeToCreated: true,
                                CreatedAt:              time.Now(),
                                UpdatedAt:              time.Now(),
                        }
                        if err := r.CreateSubscriptionPreference(&preference); err != nil {
                                return nil, err
                        }
                        return &preference, nil
                }
                return nil, result.Error
        }
        return &preference, nil
}

// UpdateSubscriptionPreference updates subscription preferences
func (r *GormSubscriptionRepository) UpdateSubscriptionPreference(preference *models.SubscriptionPreference) error {
        return r.db.Save(preference).Error
}

// CreateSubscriptionDigest creates a subscription digest
func (r *GormSubscriptionRepository) CreateSubscriptionDigest(digest *models.SubscriptionDigest) error {
        return r.db.Create(digest).Error
}

// GetPendingDigests retrieves pending digests
func (r *GormSubscriptionRepository) GetPendingDigests(frequency models.SubscriptionFrequency, before time.Time) ([]models.SubscriptionDigest, error) {
        var digests []models.SubscriptionDigest
        result := r.db.Where("frequency_type = ? AND delivery_status = ? AND scheduled_for <= ?", 
                frequency, "pending", before).Find(&digests)
        return digests, result.Error
}

// UpdateDigestStatus updates the status of a digest
func (r *GormSubscriptionRepository) UpdateDigestStatus(id uint, status string, sentAt *time.Time, errorMessage string) error {
        updates := map[string]interface{}{
                "delivery_status": status,
                "updated_at":      time.Now(),
        }
        
        if sentAt != nil {
                updates["sent_at"] = sentAt
        }
        
        if errorMessage != "" {
                updates["error_message"] = errorMessage
        }
        
        return r.db.Model(&models.SubscriptionDigest{}).Where("id = ?", id).Updates(updates).Error
}

// GetUserDigests retrieves digests for a user
func (r *GormSubscriptionRepository) GetUserDigests(userID uint, limit int, offset int) ([]models.SubscriptionDigest, error) {
        var digests []models.SubscriptionDigest
        result := r.db.Where("user_id = ?", userID).
                Order("scheduled_for DESC").
                Limit(limit).
                Offset(offset).
                Find(&digests)
        return digests, result.Error
}