package models

import (
        "time"

        "gorm.io/gorm"
)

// SubscriptionFrequency defines how often a user receives updates
type SubscriptionFrequency string

const (
        // FrequencyInstant for immediate notifications
        FrequencyInstant SubscriptionFrequency = "instant"
        
        // FrequencyDaily for once-a-day digests
        FrequencyDaily SubscriptionFrequency = "daily"
        
        // FrequencyWeekly for once-a-week digests
        FrequencyWeekly SubscriptionFrequency = "weekly"
)

// SubscriptionType defines what the user is subscribing to
type SubscriptionType string

const (
        // TopicSubscription for subscribing to a specific topic
        TopicSubscription SubscriptionType = "topic"
        
        // CategorySubscription for subscribing to an entire category
        CategorySubscription SubscriptionType = "category"
        
        // TagSubscription for subscribing to a specific tag
        TagSubscription SubscriptionType = "tag"
)

// AdvancedSubscription represents a user's subscription to a topic, category, or tag with enhanced features
type AdvancedSubscription struct {
        gorm.Model
        UserID            uint                 `json:"userId" gorm:"index:idx_subscription_user"`
        Type              SubscriptionType     `json:"type" gorm:"index:idx_subscription_type"`
        ReferenceID       uint                 `json:"referenceId" gorm:"index:idx_subscription_reference"` // TopicID, CategoryID, or TagID
        Frequency         SubscriptionFrequency `json:"frequency"`
        EmailNotification bool                 `json:"emailNotification" gorm:"default:true"`
        PushNotification  bool                 `json:"pushNotification" gorm:"default:true"`
        InAppNotification bool                 `json:"inAppNotification" gorm:"default:true"`
        Muted             bool                 `json:"muted" gorm:"default:false"` // If true, no notifications will be sent
        LastNotifiedAt    *time.Time           `json:"lastNotifiedAt"`
        ExpiresAt         *time.Time           `json:"expiresAt"` // Optional expiration date
        CreatedAt         time.Time            `json:"createdAt"`
        UpdatedAt         time.Time            `json:"updatedAt"`
}

// SubscriptionDigest represents a digest of updates for a subscription
type SubscriptionDigest struct {
        gorm.Model
        UserID          uint      `json:"userId" gorm:"index:idx_digest_user"`
        FrequencyType   SubscriptionFrequency `json:"frequencyType"`
        DigestContent   string    `json:"digestContent" gorm:"type:text"` // JSON content of the digest
        DeliveryStatus  string    `json:"deliveryStatus"`                 // pending, sent, failed
        ScheduledFor    time.Time `json:"scheduledFor"`
        SentAt          *time.Time `json:"sentAt"`
        ErrorMessage    string    `json:"errorMessage"`
        CreatedAt       time.Time `json:"createdAt"`
        UpdatedAt       time.Time `json:"updatedAt"`
}

// SubscriptionPreference represents user-wide preferences for subscriptions
type SubscriptionPreference struct {
        gorm.Model
        UserID                 uint                 `json:"userId" gorm:"uniqueIndex"`
        DefaultFrequency       SubscriptionFrequency `json:"defaultFrequency"`
        DefaultEmailEnabled    bool                 `json:"defaultEmailEnabled" gorm:"default:true"`
        DefaultPushEnabled     bool                 `json:"defaultPushEnabled" gorm:"default:true"`
        DefaultInAppEnabled    bool                 `json:"defaultInAppEnabled" gorm:"default:true"`
        DigestDay              int                  `json:"digestDay"`              // 0-6 for weekly digests (Sunday-Saturday)
        DigestHour             int                  `json:"digestHour"`             // 0-23 for daily and weekly digests
        AutoSubscribeToReplies bool                 `json:"autoSubscribeToReplies" gorm:"default:true"`
        AutoSubscribeToCreated bool                 `json:"autoSubscribeToCreated" gorm:"default:true"`
        CreatedAt              time.Time            `json:"createdAt"`
        UpdatedAt              time.Time            `json:"updatedAt"`
}