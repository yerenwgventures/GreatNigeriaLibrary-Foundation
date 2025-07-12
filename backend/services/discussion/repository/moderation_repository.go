package repository

import (
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/models"
        "gorm.io/gorm"
)

// ModerationRepository defines the interface for moderation operations
type ModerationRepository interface {
        // Content moderation rules
        CreateModerationRule(rule *models.ContentModerationRule) error
        GetModerationRules(active bool) ([]models.ContentModerationRule, error)
        GetModerationRuleByID(id uint) (*models.ContentModerationRule, error)
        UpdateModerationRule(rule *models.ContentModerationRule) error
        DeleteModerationRule(id uint) error
        GetActiveRulesByType(appliesTo string) ([]models.ContentModerationRule, error)
        
        // Content filtering results
        CreateFilterResult(result *models.ContentFilterResult) error
        GetFilterResultByID(id uint) (*models.ContentFilterResult, error)
        GetFilterResultsByContent(contentType string, contentID uint) ([]models.ContentFilterResult, error)
        GetFilterResultsByUser(userID uint, limit, offset int) ([]models.ContentFilterResult, error)
        UpdateFilterResult(result *models.ContentFilterResult) error
        
        // Moderation queue
        AddToModerationQueue(item *models.ModAdvancedQueue) error
        GetModerationQueueItems(status string, limit, offset int) ([]models.ModAdvancedQueue, error)
        GetModerationQueueItemByID(id uint) (*models.ModAdvancedQueue, error)
        GetModerationQueueCountByStatus(status string) (int64, error)
        AssignModerationQueueItem(itemID, moderatorID uint) error
        UpdateModerationQueueItem(item *models.ModAdvancedQueue) error
        
        // User trust scores
        GetUserTrustScore(userID uint) (*models.UserTrustScore, error)
        CreateUserTrustScore(score *models.UserTrustScore) error
        UpdateUserTrustScore(score *models.UserTrustScore) error
        GetUsersByTrustLevel(level models.UserTrustLevel) ([]models.UserTrustScore, error)
        
        // Moderator privileges
        GetModeratorPrivileges(userID uint) (*models.ModeratorPrivilege, error)
        CreateModeratorPrivileges(privileges *models.ModeratorPrivilege) error
        UpdateModeratorPrivileges(privileges *models.ModeratorPrivilege) error
        GetAllModerators() ([]models.ModeratorPrivilege, error)
        IsUserModerator(userID uint) (bool, error)
        
        // User moderation actions
        CreateUserModerationAction(action *models.UserModerationAction) error
        GetUserModerationActions(userID uint) ([]models.UserModerationAction, error)
        GetUserModerationActionByID(id uint) (*models.UserModerationAction, error)
        GetActiveActionsForUser(userID uint) ([]models.UserModerationAction, error)
        DeactivateUserModerationAction(actionID uint) error
        
        // Prohibited words
        CreateProhibitedWord(word *models.ProhibitedWord) error
        GetProhibitedWords(active bool) ([]models.ProhibitedWord, error)
        GetProhibitedWordByID(id uint) (*models.ProhibitedWord, error)
        UpdateProhibitedWord(word *models.ProhibitedWord) error
        DeleteProhibitedWord(id uint) error
}

// GormModerationRepository implements the ModerationRepository interface
type GormModerationRepository struct {
        db *gorm.DB
}

// NewGormModerationRepository creates a new moderation repository
func NewGormModerationRepository(db *gorm.DB) *GormModerationRepository {
        return &GormModerationRepository{db: db}
}

// CreateModerationRule creates a new content moderation rule
func (r *GormModerationRepository) CreateModerationRule(rule *models.ContentModerationRule) error {
        return r.db.Create(rule).Error
}

// GetModerationRules retrieves all content moderation rules, optionally filtered by active status
func (r *GormModerationRepository) GetModerationRules(active bool) ([]models.ContentModerationRule, error) {
        var rules []models.ContentModerationRule
        query := r.db
        
        if active {
                query = query.Where("is_active = ?", true)
        }
        
        result := query.Find(&rules)
        return rules, result.Error
}

// GetModerationRuleByID retrieves a content moderation rule by ID
func (r *GormModerationRepository) GetModerationRuleByID(id uint) (*models.ContentModerationRule, error) {
        var rule models.ContentModerationRule
        result := r.db.First(&rule, id)
        if result.Error != nil {
                return nil, result.Error
        }
        return &rule, nil
}

// UpdateModerationRule updates a content moderation rule
func (r *GormModerationRepository) UpdateModerationRule(rule *models.ContentModerationRule) error {
        return r.db.Save(rule).Error
}

// DeleteModerationRule deletes a content moderation rule
func (r *GormModerationRepository) DeleteModerationRule(id uint) error {
        return r.db.Delete(&models.ContentModerationRule{}, id).Error
}

// GetActiveRulesByType retrieves active rules for a specific content type
func (r *GormModerationRepository) GetActiveRulesByType(appliesTo string) ([]models.ContentModerationRule, error) {
        var rules []models.ContentModerationRule
        result := r.db.Where("is_active = ? AND applies_to = ?", true, appliesTo).Find(&rules)
        return rules, result.Error
}

// CreateFilterResult creates a new content filter result
func (r *GormModerationRepository) CreateFilterResult(result *models.ContentFilterResult) error {
        return r.db.Create(result).Error
}

// GetFilterResultByID retrieves a content filter result by ID
func (r *GormModerationRepository) GetFilterResultByID(id uint) (*models.ContentFilterResult, error) {
        var result models.ContentFilterResult
        dbResult := r.db.First(&result, id)
        if dbResult.Error != nil {
                return nil, dbResult.Error
        }
        return &result, nil
}

// GetFilterResultsByContent retrieves filter results for a specific content
func (r *GormModerationRepository) GetFilterResultsByContent(contentType string, contentID uint) ([]models.ContentFilterResult, error) {
        var results []models.ContentFilterResult
        dbResult := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).Find(&results)
        return results, dbResult.Error
}

// GetFilterResultsByUser retrieves filter results for a specific user with pagination
func (r *GormModerationRepository) GetFilterResultsByUser(userID uint, limit, offset int) ([]models.ContentFilterResult, error) {
        var results []models.ContentFilterResult
        dbResult := r.db.Where("user_id = ?", userID).
                Order("created_at DESC").
                Limit(limit).
                Offset(offset).
                Find(&results)
        return results, dbResult.Error
}

// UpdateFilterResult updates a content filter result
func (r *GormModerationRepository) UpdateFilterResult(result *models.ContentFilterResult) error {
        return r.db.Save(result).Error
}

// AddToModerationQueue adds an item to the moderation queue
func (r *GormModerationRepository) AddToModerationQueue(item *models.ModAdvancedQueue) error {
        return r.db.Create(item).Error
}

// GetModerationQueueItems retrieves items from the moderation queue with pagination
func (r *GormModerationRepository) GetModerationQueueItems(status string, limit, offset int) ([]models.ModAdvancedQueue, error) {
        var items []models.ModAdvancedQueue
        query := r.db
        
        if status != "" {
                query = query.Where("status = ?", status)
        }
        
        result := query.Order("priority DESC, created_at ASC").
                Limit(limit).
                Offset(offset).
                Find(&items)
        
        return items, result.Error
}

// GetModerationQueueItemByID retrieves a moderation queue item by ID
func (r *GormModerationRepository) GetModerationQueueItemByID(id uint) (*models.ModAdvancedQueue, error) {
        var item models.ModAdvancedQueue
        result := r.db.First(&item, id)
        if result.Error != nil {
                return nil, result.Error
        }
        return &item, nil
}

// GetModerationQueueCountByStatus gets the count of moderation queue items by status
func (r *GormModerationRepository) GetModerationQueueCountByStatus(status string) (int64, error) {
        var count int64
        query := r.db.Model(&models.ModAdvancedQueue{})
        
        if status != "" {
                query = query.Where("status = ?", status)
        }
        
        result := query.Count(&count)
        return count, result.Error
}

// AssignModerationQueueItem assigns a moderation queue item to a moderator
func (r *GormModerationRepository) AssignModerationQueueItem(itemID, moderatorID uint) error {
        return r.db.Model(&models.ModAdvancedQueue{}).
                Where("id = ?", itemID).
                Updates(map[string]interface{}{
                        "assigned_to": moderatorID,
                        "updated_at":  time.Now(),
                }).Error
}

// UpdateModerationQueueItem updates a moderation queue item
func (r *GormModerationRepository) UpdateModerationQueueItem(item *models.ModAdvancedQueue) error {
        return r.db.Save(item).Error
}

// GetUserTrustScore retrieves a user's trust score
func (r *GormModerationRepository) GetUserTrustScore(userID uint) (*models.UserTrustScore, error) {
        var score models.UserTrustScore
        result := r.db.Where("user_id = ?", userID).First(&score)
        if result.Error != nil {
                if result.Error == gorm.ErrRecordNotFound {
                        // Create default trust score for new users
                        defaultScore := models.UserTrustScore{
                                UserID:           userID,
                                TrustLevel:       models.TrustLevelNewUser,
                                TrustScore:       0,
                                ContentScore:     0,
                                CommunityScore:   0,
                                ModeratorScore:   0,
                                ReportCount:      0,
                                WarningCount:     0,
                                ContentRejections: 0,
                                LastScoreUpdate:  time.Now(),
                                LastCalculatedAt: time.Now(),
                        }
                        
                        if err := r.CreateUserTrustScore(&defaultScore); err != nil {
                                return nil, err
                        }
                        
                        return &defaultScore, nil
                }
                return nil, result.Error
        }
        return &score, nil
}

// CreateUserTrustScore creates a user trust score
func (r *GormModerationRepository) CreateUserTrustScore(score *models.UserTrustScore) error {
        return r.db.Create(score).Error
}

// UpdateUserTrustScore updates a user trust score
func (r *GormModerationRepository) UpdateUserTrustScore(score *models.UserTrustScore) error {
        return r.db.Save(score).Error
}

// GetUsersByTrustLevel retrieves users by trust level
func (r *GormModerationRepository) GetUsersByTrustLevel(level models.UserTrustLevel) ([]models.UserTrustScore, error) {
        var scores []models.UserTrustScore
        result := r.db.Where("trust_level = ?", level).Find(&scores)
        return scores, result.Error
}

// GetModeratorPrivileges retrieves moderator privileges for a user
func (r *GormModerationRepository) GetModeratorPrivileges(userID uint) (*models.ModeratorPrivilege, error) {
        var privileges models.ModeratorPrivilege
        result := r.db.Where("user_id = ?", userID).First(&privileges)
        if result.Error != nil {
                return nil, result.Error
        }
        return &privileges, nil
}

// CreateModeratorPrivileges creates moderator privileges for a user
func (r *GormModerationRepository) CreateModeratorPrivileges(privileges *models.ModeratorPrivilege) error {
        return r.db.Create(privileges).Error
}

// UpdateModeratorPrivileges updates moderator privileges for a user
func (r *GormModerationRepository) UpdateModeratorPrivileges(privileges *models.ModeratorPrivilege) error {
        return r.db.Save(privileges).Error
}

// GetAllModerators retrieves all moderators
func (r *GormModerationRepository) GetAllModerators() ([]models.ModeratorPrivilege, error) {
        var moderators []models.ModeratorPrivilege
        result := r.db.Where("is_active = ?", true).Find(&moderators)
        return moderators, result.Error
}

// IsUserModerator checks if a user is a moderator
func (r *GormModerationRepository) IsUserModerator(userID uint) (bool, error) {
        var count int64
        result := r.db.Model(&models.ModeratorPrivilege{}).
                Where("user_id = ? AND is_active = ?", userID, true).
                Count(&count)
        return count > 0, result.Error
}

// CreateUserModerationAction creates a moderation action for a user
func (r *GormModerationRepository) CreateUserModerationAction(action *models.UserModerationAction) error {
        return r.db.Create(action).Error
}

// GetUserModerationActions retrieves moderation actions for a user
func (r *GormModerationRepository) GetUserModerationActions(userID uint) ([]models.UserModerationAction, error) {
        var actions []models.UserModerationAction
        result := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&actions)
        return actions, result.Error
}

// GetUserModerationActionByID retrieves a user moderation action by ID
func (r *GormModerationRepository) GetUserModerationActionByID(id uint) (*models.UserModerationAction, error) {
        var action models.UserModerationAction
        result := r.db.First(&action, id)
        if result.Error != nil {
                return nil, result.Error
        }
        return &action, nil
}

// GetActiveActionsForUser retrieves active moderation actions for a user
func (r *GormModerationRepository) GetActiveActionsForUser(userID uint) ([]models.UserModerationAction, error) {
        var actions []models.UserModerationAction
        now := time.Now()
        
        result := r.db.Where("user_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)", 
                userID, true, now).Find(&actions)
                
        return actions, result.Error
}

// DeactivateUserModerationAction deactivates a user moderation action
func (r *GormModerationRepository) DeactivateUserModerationAction(actionID uint) error {
        return r.db.Model(&models.UserModerationAction{}).
                Where("id = ?", actionID).
                Updates(map[string]interface{}{
                        "is_active":  false,
                        "updated_at": time.Now(),
                }).Error
}

// CreateProhibitedWord creates a prohibited word
func (r *GormModerationRepository) CreateProhibitedWord(word *models.ProhibitedWord) error {
        return r.db.Create(word).Error
}

// GetProhibitedWords retrieves prohibited words, optionally filtered by active status
func (r *GormModerationRepository) GetProhibitedWords(active bool) ([]models.ProhibitedWord, error) {
        var words []models.ProhibitedWord
        query := r.db
        
        if active {
                query = query.Where("is_active = ?", true)
        }
        
        result := query.Find(&words)
        return words, result.Error
}

// GetProhibitedWordByID retrieves a prohibited word by ID
func (r *GormModerationRepository) GetProhibitedWordByID(id uint) (*models.ProhibitedWord, error) {
        var word models.ProhibitedWord
        result := r.db.First(&word, id)
        if result.Error != nil {
                return nil, result.Error
        }
        return &word, nil
}

// UpdateProhibitedWord updates a prohibited word
func (r *GormModerationRepository) UpdateProhibitedWord(word *models.ProhibitedWord) error {
        return r.db.Save(word).Error
}

// DeleteProhibitedWord deletes a prohibited word
func (r *GormModerationRepository) DeleteProhibitedWord(id uint) error {
        return r.db.Delete(&models.ProhibitedWord{}, id).Error
}