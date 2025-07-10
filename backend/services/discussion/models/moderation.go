package models

import (
        "time"

        "gorm.io/gorm"
)

// ModerationAction represents the type of moderation action
type ModerationAction string

const (
        // ActionNone when no moderation action is required
        ActionNone ModerationAction = "none"
        
        // ActionWarning when a user receives a warning
        ActionWarning ModerationAction = "warning"
        
        // ActionFlagged when content is flagged for review
        ActionFlagged ModerationAction = "flagged"
        
        // ActionHidden when content is hidden from public view
        ActionHidden ModerationAction = "hidden"
        
        // ActionRemoved when content is completely removed
        ActionRemoved ModerationAction = "removed"
        
        // ActionMuted when a user is muted temporarily
        ActionMuted ModerationAction = "muted"
        
        // ActionSuspended when a user is suspended temporarily
        ActionSuspended ModerationAction = "suspended"
        
        // ActionBanned when a user is banned permanently
        ActionBanned ModerationAction = "banned"
        
        // ActionTemporaryBan for temporary bans
        ActionTemporaryBan ModerationAction = "temporary_ban"
        
        // ActionPermanentBan for permanent bans
        ActionPermanentBan ModerationAction = "permanent_ban"
        
        // ActionApprove when content is approved
        ActionApprove ModerationAction = "approve"
        
        // ActionReject when content is rejected
        ActionReject ModerationAction = "reject"
        
        // ActionSendToQueue when content should be sent to the moderation queue
        ActionSendToQueue ModerationAction = "send_to_queue"
        
        // ActionAutomaticFilter when content is automatically filtered
        ActionAutomaticFilter ModerationAction = "automatic_filter"
)

// ModAdvancedQueue represents the enhanced moderation queue (to avoid name collision with discussion.go)
type ModAdvancedQueue struct {
        gorm.Model
        ContentType    string     `json:"contentType"`
        ContentID      uint       `json:"contentId"`
        UserID         uint       `json:"userId"`
        Reason         string     `json:"reason" gorm:"type:text"`
        FilterResultID *uint      `json:"filterResultId"`
        Status         string     `json:"status"` // pending, in_review, approved, rejected
        Priority       int        `json:"priority" gorm:"default:3"`
        AssignedTo     *uint      `json:"assignedTo"`
        ReviewedBy     *uint      `json:"reviewedBy"`
        ReviewedAt     *time.Time `json:"reviewedAt"`
        Decision       string     `json:"decision"`
        Notes          string     `json:"notes" gorm:"type:text"`
        CreatedAt      time.Time  `json:"createdAt"`
        UpdatedAt      time.Time  `json:"updatedAt"`
}

// AdvancedModerationQueue represents the queue of items awaiting moderation (legacy version)
type AdvancedModerationQueue struct {
        gorm.Model
        ContentType   string     `json:"contentType"`
        ContentID     uint       `json:"contentId"`
        ReportedBy    uint       `json:"reportedBy"`
        Reason        string     `json:"reason" gorm:"type:text"`
        Status        string     `json:"status"` // pending, approved, rejected
        ModeratedBy   *uint      `json:"moderatedBy"`
        ModeratedAt   *time.Time `json:"moderatedAt"`
        Notes         string     `json:"notes" gorm:"type:text"`
        Priority      int        `json:"priority" gorm:"default:0"`
        ReportedCount int        `json:"reportedCount" gorm:"default:1"`
}

// ContentFilterResult represents the result of content filtering
type ContentFilterResult struct {
        gorm.Model
        ContentType           string           `json:"contentType"`
        ContentID             uint             `json:"contentId"`
        UserID                uint             `json:"userId"`
        TriggeredRules        string           `json:"triggeredRules" gorm:"type:text"` // JSON array of rule IDs
        Action                ModerationAction `json:"action"`
        FilteredContent       string           `json:"filteredContent" gorm:"type:text"`
        CleanedContent        string           `json:"cleanedContent" gorm:"type:text"`
        AutomaticallyProcessed bool             `json:"automaticallyProcessed" gorm:"default:true"`
        ModeratorID           *uint            `json:"moderatorId"`
        ReviewedAt            *time.Time       `json:"reviewedAt"`
        CreatedAt             time.Time        `json:"createdAt"`
        UpdatedAt             time.Time        `json:"updatedAt"`
        // Keeping previous fields for backward compatibility
        FilterType            string           `json:"filterType"`
        FilterScore           float64          `json:"filterScore"`
        IsFlagged             bool             `json:"isFlagged"`
        ReviewedBy            *uint            `json:"reviewedBy"`
        IsApproved            bool             `json:"isApproved"`
        ReviewNotes           string           `json:"reviewNotes" gorm:"type:text"`
        DetectedProblems      string           `json:"detectedProblems" gorm:"type:text"`
        AutoModeratedActions  string           `json:"autoModeratedActions" gorm:"type:text"`
}

// ContentModerationRule represents a rule for content moderation
type ContentModerationRule struct {
        gorm.Model
        Name          string    `json:"name"`
        Description   string    `json:"description" gorm:"type:text"`
        RuleType      string    `json:"ruleType"`
        PatternType   string    `json:"patternType"` // exact, regex, wildcard
        Pattern       string    `json:"pattern" gorm:"type:text"`
        Action        string    `json:"action"` // String representation of ModerationAction
        Severity      int       `json:"severity"`
        CategoryID    *uint     `json:"categoryId"`
        IsActive      bool      `json:"isActive" gorm:"default:true"`
        CreatedBy     uint      `json:"createdBy"`
        LastUpdatedBy uint      `json:"lastUpdatedBy"`
        AppliesTo     string    `json:"appliesTo"` // topics, comments, all
}

// UserTrustLevel represents the trust level of a user
type UserTrustLevel string

const (
        // TrustLevelNewUser for new users with minimal trust
        TrustLevelNewUser UserTrustLevel = "new_user"
        
        // TrustLevelBasic for users with basic trust
        TrustLevelBasic UserTrustLevel = "basic"
        
        // TrustLevelMember for established community members
        TrustLevelMember UserTrustLevel = "member"
        
        // TrustLevelRegular for regular contributors
        TrustLevelRegular UserTrustLevel = "regular"
        
        // TrustLevelLeader for community leaders
        TrustLevelLeader UserTrustLevel = "leader"
)

// UserTrustScore represents a user's trust score
type UserTrustScore struct {
        gorm.Model
        UserID               uint          `json:"userId" gorm:"uniqueIndex"`
        TrustLevel           UserTrustLevel `json:"trustLevel"`
        Score                float64       `json:"score"`
        TrustScore           float64       `json:"trustScore"`
        ContentScore         float64       `json:"contentScore"`
        CommunityScore       float64       `json:"communityScore"`
        ModeratorScore       float64       `json:"moderatorScore"`
        ReportCount          int           `json:"reportCount" gorm:"default:0"`
        WarningCount         int           `json:"warningCount" gorm:"default:0"`
        ContentRejections    int           `json:"contentRejections" gorm:"default:0"`
        LastCalculatedAt     time.Time     `json:"lastCalculatedAt"`
        LastScoreUpdate      time.Time     `json:"lastScoreUpdate"`
        PositiveFactors      string        `json:"positiveFactors" gorm:"type:text"`
        NegativeFactors      string        `json:"negativeFactors" gorm:"type:text"`
        ManualAdjustment     float64       `json:"manualAdjustment" gorm:"default:0"`
        ManualAdjustedBy     *uint         `json:"manualAdjustedBy"`
        ManualAdjustmentNote string        `json:"manualAdjustmentNote" gorm:"type:text"`
}

// ModeratorPrivilege represents the privileges of a moderator
type ModeratorPrivilege struct {
        gorm.Model
        UserID              uint      `json:"userId" gorm:"index:idx_mod_privilege_user"`
        IsGlobalModerator   bool      `json:"isGlobalModerator" gorm:"default:false"`
        CanApproveContent   bool      `json:"canApproveContent" gorm:"default:true"`
        CanRejectContent    bool      `json:"canRejectContent" gorm:"default:true"`
        CanDeleteContent    bool      `json:"canDeleteContent" gorm:"default:false"`
        CanBanUsers         bool      `json:"canBanUsers" gorm:"default:false"`
        CanEditAnyContent   bool      `json:"canEditAnyContent" gorm:"default:false"`
        CanEditContent      bool      `json:"canEditContent" gorm:"default:false"`
        CanManageRules      bool      `json:"canManageRules" gorm:"default:false"`
        CanAssignModerators bool      `json:"canAssignModerators" gorm:"default:false"`
        CanAccessDashboard  bool      `json:"canAccessDashboard" gorm:"default:true"`
        IsActive            bool      `json:"isActive" gorm:"default:true"`
        AssignedBy          uint      `json:"assignedBy"`
        AssignedAt          time.Time `json:"assignedAt"`
        ExpiresAt           *time.Time `json:"expiresAt"`
        Notes               string    `json:"notes" gorm:"type:text"`
        CreatedAt           time.Time `json:"createdAt"`
        UpdatedAt           time.Time `json:"updatedAt"`
}

// UserModerationAction represents an action taken against a user
type UserModerationAction struct {
        gorm.Model
        UserID             uint            `json:"userId" gorm:"index:idx_mod_action_user"`
        ActionType         string          `json:"actionType"` // warn, mute, suspend, ban
        Reason             string          `json:"reason" gorm:"type:text"`
        AppliedBy          uint            `json:"appliedBy"`
        AppliedAt          time.Time       `json:"appliedAt"`
        Duration           *int            `json:"duration"` // in hours, null for permanent
        ExpiresAt          *time.Time      `json:"expiresAt"`
        IsActive           bool            `json:"isActive" gorm:"default:true"`
        RevokedBy          *uint           `json:"revokedBy"`
        RevokedAt          *time.Time      `json:"revokedAt"`
        RevocationReason   string          `json:"revocationReason" gorm:"type:text"`
        RelatedContentID   *uint           `json:"relatedContentId"`
        RelatedContentType string          `json:"relatedContentType"`
        ModeratorID        uint            `json:"moderatorId"`
        Notes              string          `json:"notes" gorm:"type:text"`
        CreatedAt          time.Time       `json:"createdAt"`
        UpdatedAt          time.Time       `json:"updatedAt"`
}

// ProhibitedWord represents a prohibited word or phrase
type ProhibitedWord struct {
        gorm.Model
        Word           string    `json:"word" gorm:"uniqueIndex"`
        Severity       int       `json:"severity" gorm:"default:1"`
        Action         string    `json:"action"` // warn, flag, block
        IsRegex        bool      `json:"isRegex" gorm:"default:false"`
        CategoryID     *uint     `json:"categoryId"` // null means applies to all categories
        IsActive       bool      `json:"isActive" gorm:"default:true"`
        CreatedBy      uint      `json:"createdBy"`
        LastUpdatedBy  uint      `json:"lastUpdatedBy"`
        Replacements   string    `json:"replacements" gorm:"type:text"` // comma-separated list of allowed replacements
        Replacement    string    `json:"replacement" gorm:"type:text"`  // single replacement
        IsAutoReplace  bool      `json:"isAutoReplace" gorm:"default:false"`
        CreatedAt      time.Time `json:"createdAt"`
        UpdatedAt      time.Time `json:"updatedAt"`
}