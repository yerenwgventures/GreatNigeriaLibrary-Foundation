package service

import (
        "errors"
        "fmt"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/repository"
)

// SubscriptionService defines the interface for subscription operations
type SubscriptionService interface {
        Subscribe(userID uint, refType models.SubscriptionType, refID uint, frequency models.SubscriptionFrequency, emailNotif, pushNotif, inAppNotif bool) (*models.AdvancedSubscription, error)
        Unsubscribe(subscriptionID uint) error
        UnsubscribeFromReference(userID uint, refType models.SubscriptionType, refID uint) error
        GetUserSubscriptions(userID uint) ([]models.AdvancedSubscription, error)
        GetSubscriptionsByType(userID uint, subscriptionType models.SubscriptionType) ([]models.AdvancedSubscription, error)
        UpdateSubscriptionSettings(subscriptionID uint, frequency models.SubscriptionFrequency, emailNotif, pushNotif, inAppNotif, muted bool) (*models.AdvancedSubscription, error)
        
        GetUserPreferences(userID uint) (*models.SubscriptionPreference, error)
        UpdateUserPreferences(userID uint, defaults models.SubscriptionPreference) (*models.SubscriptionPreference, error)
        
        GenerateDigests(frequency models.SubscriptionFrequency) (int, error)
        ProcessPendingDigests() (int, error)
        GetUserDigestHistory(userID uint, page, pageSize int) ([]models.SubscriptionDigest, error)
}

// SubscriptionServiceImpl implements the SubscriptionService interface
type SubscriptionServiceImpl struct {
        subscriptionRepo repository.SubscriptionRepository
        discussionRepo   repository.DiscussionRepository
        topicRepo        repository.TopicRepository
        categoryRepo     repository.CategoryRepository
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(
        subscriptionRepo repository.SubscriptionRepository,
        discussionRepo repository.DiscussionRepository,
        topicRepo repository.TopicRepository,
        categoryRepo repository.CategoryRepository,
) SubscriptionService {
        return &SubscriptionServiceImpl{
                subscriptionRepo: subscriptionRepo,
                discussionRepo:   discussionRepo,
                topicRepo:        topicRepo,
                categoryRepo:     categoryRepo,
        }
}

// Subscribe creates a new subscription
func (s *SubscriptionServiceImpl) Subscribe(
        userID uint,
        refType models.SubscriptionType,
        refID uint,
        frequency models.SubscriptionFrequency,
        emailNotif, pushNotif, inAppNotif bool,
) (*models.AdvancedSubscription, error) {
        // Check if subscription already exists
        existingSub, err := s.subscriptionRepo.GetSubscriptionByReference(userID, refType, refID)
        if err == nil && existingSub != nil {
                // Update existing subscription with new values
                existingSub.Frequency = frequency
                existingSub.EmailNotification = emailNotif
                existingSub.PushNotification = pushNotif
                existingSub.InAppNotification = inAppNotif
                existingSub.Muted = false
                existingSub.UpdatedAt = time.Now()
                
                if err := s.subscriptionRepo.UpdateSubscription(existingSub); err != nil {
                        return nil, fmt.Errorf("failed to update existing subscription: %w", err)
                }
                
                return existingSub, nil
        }
        
        // Validate reference exists
        if err := s.validateReference(refType, refID); err != nil {
                return nil, err
        }
        
        // Create new subscription using AdvancedSubscription
        subscription := &models.AdvancedSubscription{
                UserID:            userID,
                Type:              refType,
                ReferenceID:       refID,
                Frequency:         frequency,
                EmailNotification: emailNotif,
                PushNotification:  pushNotif,
                InAppNotification: inAppNotif,
                Muted:             false,
                CreatedAt:         time.Now(),
                UpdatedAt:         time.Now(),
        }
        
        if err := s.subscriptionRepo.CreateSubscription(subscription); err != nil {
                return nil, fmt.Errorf("failed to create subscription: %w", err)
        }
        
        return subscription, nil
}

// Unsubscribe deletes a subscription
func (s *SubscriptionServiceImpl) Unsubscribe(subscriptionID uint) error {
        // Check if subscription exists
        _, err := s.subscriptionRepo.GetSubscription(subscriptionID)
        if err != nil {
                return fmt.Errorf("subscription not found: %w", err)
        }
        
        // Delete subscription
        if err := s.subscriptionRepo.DeleteSubscription(subscriptionID); err != nil {
                return fmt.Errorf("failed to delete subscription: %w", err)
        }
        
        return nil
}

// UnsubscribeFromReference deletes a subscription by reference
func (s *SubscriptionServiceImpl) UnsubscribeFromReference(userID uint, refType models.SubscriptionType, refID uint) error {
        // Get subscription by reference
        subscription, err := s.subscriptionRepo.GetSubscriptionByReference(userID, refType, refID)
        if err != nil {
                return fmt.Errorf("subscription not found: %w", err)
        }
        
        // Delete subscription
        if err := s.subscriptionRepo.DeleteSubscription(subscription.ID); err != nil {
                return fmt.Errorf("failed to delete subscription: %w", err)
        }
        
        return nil
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (s *SubscriptionServiceImpl) GetUserSubscriptions(userID uint) ([]models.AdvancedSubscription, error) {
        // Get subscriptions directly from repository
        return s.subscriptionRepo.GetUserSubscriptions(userID)
}

// GetSubscriptionsByType retrieves subscriptions of a specific type for a user
func (s *SubscriptionServiceImpl) GetSubscriptionsByType(userID uint, subscriptionType models.SubscriptionType) ([]models.AdvancedSubscription, error) {
        // Get subscriptions directly from repository
        return s.subscriptionRepo.GetSubscriptionsByType(userID, subscriptionType)
}

// UpdateSubscriptionSettings updates subscription settings
func (s *SubscriptionServiceImpl) UpdateSubscriptionSettings(
        subscriptionID uint,
        frequency models.SubscriptionFrequency,
        emailNotif, pushNotif, inAppNotif, muted bool,
) (*models.AdvancedSubscription, error) {
        // Get subscription
        subscription, err := s.subscriptionRepo.GetSubscription(subscriptionID)
        if err != nil {
                return nil, fmt.Errorf("subscription not found: %w", err)
        }
        
        // Update settings directly
        subscription.Frequency = frequency
        subscription.EmailNotification = emailNotif
        subscription.PushNotification = pushNotif
        subscription.InAppNotification = inAppNotif
        subscription.Muted = muted
        subscription.UpdatedAt = time.Now()
        
        // Save changes
        if err := s.subscriptionRepo.UpdateSubscription(subscription); err != nil {
                return nil, fmt.Errorf("failed to update subscription: %w", err)
        }
        
        return subscription, nil
}

// GetUserPreferences retrieves subscription preferences for a user
func (s *SubscriptionServiceImpl) GetUserPreferences(userID uint) (*models.SubscriptionPreference, error) {
        return s.subscriptionRepo.GetSubscriptionPreference(userID)
}

// UpdateUserPreferences updates subscription preferences for a user
func (s *SubscriptionServiceImpl) UpdateUserPreferences(
        userID uint,
        defaults models.SubscriptionPreference,
) (*models.SubscriptionPreference, error) {
        // Get current preferences
        preference, err := s.subscriptionRepo.GetSubscriptionPreference(userID)
        if err != nil {
                return nil, fmt.Errorf("failed to get preferences: %w", err)
        }
        
        // Update preferences
        preference.DefaultFrequency = defaults.DefaultFrequency
        preference.DefaultEmailEnabled = defaults.DefaultEmailEnabled
        preference.DefaultPushEnabled = defaults.DefaultPushEnabled
        preference.DefaultInAppEnabled = defaults.DefaultInAppEnabled
        preference.DigestDay = defaults.DigestDay
        preference.DigestHour = defaults.DigestHour
        preference.AutoSubscribeToReplies = defaults.AutoSubscribeToReplies
        preference.AutoSubscribeToCreated = defaults.AutoSubscribeToCreated
        preference.UpdatedAt = time.Now()
        
        // Save changes
        if err := s.subscriptionRepo.UpdateSubscriptionPreference(preference); err != nil {
                return nil, fmt.Errorf("failed to update preferences: %w", err)
        }
        
        return preference, nil
}

// DigestContent represents the structure of a digest's content
type DigestContent struct {
        Topics       []DigestTopic       `json:"topics"`
        Categories   []DigestCategory    `json:"categories"`
        Tags         []DigestTag         `json:"tags"`
        TotalUpdates int                 `json:"totalUpdates"`
}

// DigestTopic represents a topic update in a digest
type DigestTopic struct {
        TopicID      uint      `json:"topicId"`
        Title        string    `json:"title"`
        NewPosts     int       `json:"newPosts"`
        LastActivity time.Time `json:"lastActivity"`
        URL          string    `json:"url"`
}

// DigestCategory represents a category update in a digest
type DigestCategory struct {
        CategoryID   uint           `json:"categoryId"`
        Name         string         `json:"name"`
        NewTopics    int            `json:"newTopics"`
        NewPosts     int            `json:"newPosts"`
        LastActivity time.Time      `json:"lastActivity"`
        TopTopics    []DigestTopic  `json:"topTopics"`
}

// DigestTag represents a tag update in a digest
type DigestTag struct {
        TagID        uint           `json:"tagId"`
        Name         string         `json:"name"`
        NewTopics    int            `json:"newTopics"`
        NewPosts     int            `json:"newPosts"`
        TopTopics    []DigestTopic  `json:"topTopics"`
}

// GenerateDigests generates digests for a specific frequency
func (s *SubscriptionServiceImpl) GenerateDigests(frequency models.SubscriptionFrequency) (int, error) {
        // Get all subscriptions with the given frequency
        // For simplicity, we'll use directly SQLs here as this would typically involve complex joins
        // and aggregations better handled in raw SQL
        
        // For now, just simulate the process with a simple implementation
        // In a real system, we'd have a more efficient batch processing approach
        
        // TODO: Implement actual digest generation logic
        // - Query all users with subscriptions at the given frequency
        // - For each user, gather all content updates since their last notification
        // - Generate a digest
        // - Schedule the digest for delivery
        
        // Simulate some successful digest creations
        createdCount := 0
        
        // Return number of digests created
        return createdCount, nil
}

// ProcessPendingDigests processes pending digests for delivery
func (s *SubscriptionServiceImpl) ProcessPendingDigests() (int, error) {
        // Get pending digests scheduled for now or earlier
        now := time.Now()
        digests, err := s.subscriptionRepo.GetPendingDigests(models.FrequencyDaily, now)
        if err != nil {
                return 0, fmt.Errorf("failed to fetch pending digests: %w", err)
        }
        
        processedCount := 0
        
        // Process each digest
        for _, digest := range digests {
                // In a real implementation, we'd:
                // 1. Parse the digest content
                // 2. Format appropriate email/notification
                // 3. Send it via email/push notification service
                // 4. Update the status
                
                // Simulate successful delivery
                err := s.subscriptionRepo.UpdateDigestStatus(digest.ID, "sent", &now, "")
                if err != nil {
                        // Log the error but continue processing other digests
                        fmt.Printf("Error updating digest status: %v\n", err)
                        continue
                }
                
                processedCount++
        }
        
        // Return number of digests processed
        return processedCount, nil
}

// GetUserDigestHistory retrieves digest history for a user with pagination
func (s *SubscriptionServiceImpl) GetUserDigestHistory(userID uint, page, pageSize int) ([]models.SubscriptionDigest, error) {
        if page < 1 {
                page = 1
        }
        
        if pageSize < 1 {
                pageSize = 10
        }
        
        offset := (page - 1) * pageSize
        return s.subscriptionRepo.GetUserDigests(userID, pageSize, offset)
}

// validateReference validates that the reference exists
func (s *SubscriptionServiceImpl) validateReference(refType models.SubscriptionType, refID uint) error {
        switch refType {
        case models.TopicSubscription:
                // Validate topic exists
                _, err := s.topicRepo.GetTopicByID(refID)
                if err != nil {
                        return fmt.Errorf("topic not found: %w", err)
                }
        case models.CategorySubscription:
                // Validate category exists
                _, err := s.categoryRepo.GetCategoryByID(refID)
                if err != nil {
                        return fmt.Errorf("category not found: %w", err)
                }
        case models.TagSubscription:
                // Validate tag exists
                // (assuming we have a tag repository)
                // _, err := s.tagRepo.GetTagByID(refID)
                // if err != nil {
                //     return fmt.Errorf("tag not found: %w", err)
                // }
                // For now, we'll skip detailed tag validation
        default:
                return errors.New("invalid subscription type")
        }
        
        return nil
}