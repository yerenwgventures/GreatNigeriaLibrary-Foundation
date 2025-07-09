package repository

import (
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/models"
        "gorm.io/gorm"
)

// DiscussionRepository defines the interface for discussion-related database operations
type DiscussionRepository interface {
        // Categories
        GetCategories() ([]models.Category, error)
        GetCategoryByID(id uint) (*models.Category, error)
        GetCategoryBySlug(slug string) (*models.Category, error)
        CreateCategory(category *models.Category) error
        UpdateCategory(category *models.Category) error
        DeleteCategory(id uint) error
        
        // Category Configuration
        GetCategoryConfig(categoryID uint) (*models.CategoryConfig, error)
        CreateCategoryConfig(config *models.CategoryConfig) error
        UpdateCategoryConfig(config *models.CategoryConfig) error
        
        // Posting Rules
        GetPostingRules(categoryID uint) (*models.PostingRules, error)
        CreatePostingRules(rules *models.PostingRules) error
        UpdatePostingRules(rules *models.PostingRules) error
        
        // Auto Moderation Settings
        GetAutoModerationSettings(categoryID uint) (*models.AutoModerationSettings, error)
        CreateAutoModerationSettings(settings *models.AutoModerationSettings) error
        UpdateAutoModerationSettings(settings *models.AutoModerationSettings) error
        
        // Category Moderators
        GetCategoryModerators(categoryID uint) ([]models.CategoryModerator, error)
        AddCategoryModerator(moderator *models.CategoryModerator) error
        UpdateCategoryModerator(moderator *models.CategoryModerator) error
        RemoveCategoryModerator(categoryID, userID uint) error

        // Topics
        GetTopics(page, pageSize int, filters map[string]interface{}) ([]models.Topic, int64, error)
        GetTopicByID(id uint) (*models.Topic, error)
        GetTopicsByCategory(categoryID uint, page, pageSize int) ([]models.Topic, int64, error)
        GetTopicsByUser(userID uint, page, pageSize int) ([]models.Topic, int64, error)
        GetTopicsByBookSection(bookID, chapterID, sectionID uint) ([]models.Topic, error)
        CreateTopic(topic *models.Topic) error
        UpdateTopic(topic *models.Topic) error
        DeleteTopic(id uint) error
        IncrementTopicViewCount(id uint) error
        UpdateTopicLastPostTime(id uint) error
        PinTopic(id uint, pinned bool) error
        LockTopic(id uint, locked bool) error

        // Comments
        GetCommentsByTopic(topicID uint, page, pageSize int) ([]models.Comment, int64, error)
        GetCommentByID(id uint) (*models.Comment, error)
        GetRepliesByComment(commentID uint) ([]models.Comment, error)
        CreateComment(comment *models.Comment) error
        UpdateComment(comment *models.Comment) error
        DeleteComment(id uint) error
        MarkCommentAsEdited(id uint) error

        // Reactions
        AddReaction(reaction *models.Reaction) error
        RemoveReaction(userID, targetID uint, targetType, reactionType string) error
        GetReactionsByTarget(targetID uint, targetType string) ([]models.Reaction, error)
        GetReactionsSummary(targetID uint, targetType string) ([]models.ReactionSummary, error)

        // Tags
        GetAllTags() ([]models.Tag, error)
        GetTagByID(id uint) (*models.Tag, error)
        GetTagBySlug(slug string) (*models.Tag, error)
        CreateTag(tag *models.Tag) error
        UpdateTag(tag *models.Tag) error
        DeleteTag(id uint) error
        AddTagToTopic(topicID, tagID uint) error
        RemoveTagFromTopic(topicID, tagID uint) error
        GetTagsByTopic(topicID uint) ([]models.Tag, error)

        // Moderation
        AddToModerationQueue(queue *models.ModerationQueue) error
        GetModerationQueueItems(status string, page, pageSize int) ([]models.ModerationQueue, int64, error)
        UpdateModerationQueueItem(item *models.ModerationQueue) error
        FlagComment(commentID, reportedBy uint, reason string) error
        ApproveFlaggedComment(commentID, moderatorID uint, notes string) error
        RejectFlaggedComment(commentID, moderatorID uint, notes string) error

        // Subscriptions
        SubscribeToTopic(userID, topicID uint, notifyOnReply bool, digestType string) error
        UnsubscribeFromTopic(userID, topicID uint) error
        GetUserSubscriptions(userID uint) ([]models.Subscription, error)
        UpdateSubscriptionSettings(subscription *models.Subscription) error

        // User Stats & Preferences
        GetUserDiscussionStats(userID uint) (*models.UserDiscussionStats, error)
        UpdateUserDiscussionStats(stats *models.UserDiscussionStats) error
        GetNotificationPreferences(userID uint) (*models.NotificationPreference, error)
        UpdateNotificationPreferences(preferences *models.NotificationPreference) error

        // Topic Views & Activity
        RecordTopicView(userID, topicID uint) error
        GetTopicViewCount(topicID uint) (int, error)
        GetRecentlyViewedTopics(userID uint, limit int) ([]models.Topic, error)

        // Mentions
        CreateMention(mention *models.Mention) error
        GetMentionsForUser(userID uint, isRead bool, page, pageSize int) ([]models.Mention, int64, error)
        MarkMentionAsRead(id uint) error
}

// GormDiscussionRepository implements the DiscussionRepository interface using GORM
type GormDiscussionRepository struct {
        db *gorm.DB
}

// NewGormDiscussionRepository creates a new GORM-based repository for discussions
func NewGormDiscussionRepository(db *gorm.DB) *GormDiscussionRepository {
        return &GormDiscussionRepository{
                db: db,
        }
}

// Ensure GormDiscussionRepository implements all required interfaces
var _ TopicRepository = (*GormDiscussionRepository)(nil)
var _ CommentRepository = (*GormDiscussionRepository)(nil)
var _ CategoryRepository = (*GormDiscussionRepository)(nil)
var _ UserRepository = (*GormDiscussionRepository)(nil)

// GetCategories retrieves all active categories
func (r *GormDiscussionRepository) GetCategories() ([]models.Category, error) {
        var categories []models.Category
        err := r.db.Where("is_active = ?", true).Order("sort_order").Find(&categories).Error
        return categories, err
}

// GetCategoryByID retrieves a category by its ID
func (r *GormDiscussionRepository) GetCategoryByID(id uint) (*models.Category, error) {
        var category models.Category
        err := r.db.First(&category, id).Error
        if err != nil {
                return nil, err
        }
        return &category, nil
}

// GetCategoryBySlug retrieves a category by its slug
func (r *GormDiscussionRepository) GetCategoryBySlug(slug string) (*models.Category, error) {
        var category models.Category
        err := r.db.Where("slug = ?", slug).First(&category).Error
        if err != nil {
                return nil, err
        }
        return &category, nil
}

// CreateCategory creates a new category
func (r *GormDiscussionRepository) CreateCategory(category *models.Category) error {
        return r.db.Create(category).Error
}

// UpdateCategory updates an existing category
func (r *GormDiscussionRepository) UpdateCategory(category *models.Category) error {
        return r.db.Save(category).Error
}

// DeleteCategory deletes a category
func (r *GormDiscussionRepository) DeleteCategory(id uint) error {
        return r.db.Delete(&models.Category{}, id).Error
}

// GetCategoryConfig retrieves configuration for a category
func (r *GormDiscussionRepository) GetCategoryConfig(categoryID uint) (*models.CategoryConfig, error) {
        var config models.CategoryConfig
        err := r.db.Where("category_id = ?", categoryID).First(&config).Error
        if err != nil {
                if err == gorm.ErrRecordNotFound {
                        // Return default config if not found
                        return &models.CategoryConfig{
                                CategoryID: categoryID,
                        }, nil
                }
                return nil, err
        }
        return &config, nil
}

// CreateCategoryConfig creates a new category configuration
func (r *GormDiscussionRepository) CreateCategoryConfig(config *models.CategoryConfig) error {
        return r.db.Create(config).Error
}

// UpdateCategoryConfig updates an existing category configuration
func (r *GormDiscussionRepository) UpdateCategoryConfig(config *models.CategoryConfig) error {
        return r.db.Save(config).Error
}

// GetPostingRules retrieves posting rules for a category
func (r *GormDiscussionRepository) GetPostingRules(categoryID uint) (*models.PostingRules, error) {
        var rules models.PostingRules
        err := r.db.Where("category_id = ?", categoryID).First(&rules).Error
        if err != nil {
                if err == gorm.ErrRecordNotFound {
                        // Return default rules if not found
                        return &models.PostingRules{
                                CategoryID: categoryID,
                        }, nil
                }
                return nil, err
        }
        return &rules, nil
}

// CreatePostingRules creates new posting rules for a category
func (r *GormDiscussionRepository) CreatePostingRules(rules *models.PostingRules) error {
        return r.db.Create(rules).Error
}

// UpdatePostingRules updates existing posting rules for a category
func (r *GormDiscussionRepository) UpdatePostingRules(rules *models.PostingRules) error {
        return r.db.Save(rules).Error
}

// GetAutoModerationSettings retrieves auto-moderation settings for a category
func (r *GormDiscussionRepository) GetAutoModerationSettings(categoryID uint) (*models.AutoModerationSettings, error) {
        var settings models.AutoModerationSettings
        err := r.db.Where("category_id = ?", categoryID).First(&settings).Error
        if err != nil {
                if err == gorm.ErrRecordNotFound {
                        // Return default settings if not found
                        return &models.AutoModerationSettings{
                                CategoryID: categoryID,
                        }, nil
                }
                return nil, err
        }
        return &settings, nil
}

// CreateAutoModerationSettings creates new auto-moderation settings for a category
func (r *GormDiscussionRepository) CreateAutoModerationSettings(settings *models.AutoModerationSettings) error {
        return r.db.Create(settings).Error
}

// UpdateAutoModerationSettings updates existing auto-moderation settings for a category
func (r *GormDiscussionRepository) UpdateAutoModerationSettings(settings *models.AutoModerationSettings) error {
        return r.db.Save(settings).Error
}

// GetCategoryModerators retrieves all moderators for a category
func (r *GormDiscussionRepository) GetCategoryModerators(categoryID uint) ([]models.CategoryModerator, error) {
        var moderators []models.CategoryModerator
        err := r.db.Where("category_id = ?", categoryID).Find(&moderators).Error
        return moderators, err
}

// AddCategoryModerator adds a new moderator to a category
func (r *GormDiscussionRepository) AddCategoryModerator(moderator *models.CategoryModerator) error {
        return r.db.Create(moderator).Error
}

// UpdateCategoryModerator updates a category moderator's permissions
func (r *GormDiscussionRepository) UpdateCategoryModerator(moderator *models.CategoryModerator) error {
        return r.db.Save(moderator).Error
}

// RemoveCategoryModerator removes a moderator from a category
func (r *GormDiscussionRepository) RemoveCategoryModerator(categoryID, userID uint) error {
        return r.db.Where("category_id = ? AND user_id = ?", categoryID, userID).Delete(&models.CategoryModerator{}).Error
}

// GetTopics retrieves topics with pagination and filtering
func (r *GormDiscussionRepository) GetTopics(page, pageSize int, filters map[string]interface{}) ([]models.Topic, int64, error) {
        var topics []models.Topic
        var total int64

        // Start with the base query
        query := r.db.Model(&models.Topic{})

        // Apply filters
        for key, value := range filters {
                query = query.Where(key, value)
        }

        // Count total matching records
        err := query.Count(&total).Error
        if err != nil {
                return nil, 0, err
        }

        // Calculate offset and fetch paginated results
        offset := (page - 1) * pageSize
        err = query.
                Preload("Category").
                Order("is_pinned DESC, last_post_at DESC").
                Offset(offset).
                Limit(pageSize).
                Find(&topics).Error

        if err != nil {
                return nil, 0, err
        }

        // Calculate replies count for each topic
        for i := range topics {
                var count int64
                err := r.db.Model(&models.Comment{}).Where("topic_id = ?", topics[i].ID).Count(&count).Error
                if err != nil {
                        return nil, 0, err
                }
                topics[i].RepliesCount = int(count)
        }

        return topics, total, nil
}

// GetTopicByID retrieves a topic by its ID
func (r *GormDiscussionRepository) GetTopicByID(id uint) (*models.Topic, error) {
        var topic models.Topic
        err := r.db.
                Preload("Category").
                Preload("Tags").
                First(&topic, id).Error
        if err != nil {
                return nil, err
        }

        // Get replies count
        var count int64
        err = r.db.Model(&models.Comment{}).Where("topic_id = ?", topic.ID).Count(&count).Error
        if err != nil {
                return nil, err
        }
        topic.RepliesCount = int(count)

        return &topic, nil
}

// GetTopicsByCategory retrieves topics in a specific category
func (r *GormDiscussionRepository) GetTopicsByCategory(categoryID uint, page, pageSize int) ([]models.Topic, int64, error) {
        var topics []models.Topic
        var total int64

        // Count total
        err := r.db.Model(&models.Topic{}).Where("category_id = ?", categoryID).Count(&total).Error
        if err != nil {
                return nil, 0, err
        }

        // Calculate offset and fetch paginated results
        offset := (page - 1) * pageSize
        err = r.db.
                Where("category_id = ?", categoryID).
                Preload("Category").
                Order("is_pinned DESC, last_post_at DESC").
                Offset(offset).
                Limit(pageSize).
                Find(&topics).Error

        if err != nil {
                return nil, 0, err
        }

        // Calculate replies count for each topic
        for i := range topics {
                var count int64
                err := r.db.Model(&models.Comment{}).Where("topic_id = ?", topics[i].ID).Count(&count).Error
                if err != nil {
                        return nil, 0, err
                }
                topics[i].RepliesCount = int(count)
        }

        return topics, total, nil
}

// GetTopicsByUser retrieves topics created by a specific user
func (r *GormDiscussionRepository) GetTopicsByUser(userID uint, page, pageSize int) ([]models.Topic, int64, error) {
        var topics []models.Topic
        var total int64

        // Count total
        err := r.db.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&total).Error
        if err != nil {
                return nil, 0, err
        }

        // Calculate offset and fetch paginated results
        offset := (page - 1) * pageSize
        err = r.db.
                Where("user_id = ?", userID).
                Preload("Category").
                Order("created_at DESC").
                Offset(offset).
                Limit(pageSize).
                Find(&topics).Error

        if err != nil {
                return nil, 0, err
        }

        // Calculate replies count for each topic
        for i := range topics {
                var count int64
                err := r.db.Model(&models.Comment{}).Where("topic_id = ?", topics[i].ID).Count(&count).Error
                if err != nil {
                        return nil, 0, err
                }
                topics[i].RepliesCount = int(count)
        }

        return topics, total, nil
}

// GetTopicsByBookSection retrieves topics related to a specific book section
func (r *GormDiscussionRepository) GetTopicsByBookSection(bookID, chapterID, sectionID uint) ([]models.Topic, error) {
        var topics []models.Topic
        query := r.db.Preload("Category")

        // Build the query based on provided parameters
        if bookID > 0 {
                query = query.Where("book_id = ?", bookID)
        }
        if chapterID > 0 {
                query = query.Where("chapter_id = ?", chapterID)
        }
        if sectionID > 0 {
                query = query.Where("section_id = ?", sectionID)
        }

        err := query.Order("created_at DESC").Find(&topics).Error
        if err != nil {
                return nil, err
        }

        // Calculate replies count for each topic
        for i := range topics {
                var count int64
                err := r.db.Model(&models.Comment{}).Where("topic_id = ?", topics[i].ID).Count(&count).Error
                if err != nil {
                        return nil, err
                }
                topics[i].RepliesCount = int(count)
        }

        return topics, nil
}

// CreateTopic creates a new topic
func (r *GormDiscussionRepository) CreateTopic(topic *models.Topic) error {
        // Set the last post time to now
        topic.LastPostAt = time.Now()
        
        // Begin a transaction
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Create the topic
                if err := tx.Create(topic).Error; err != nil {
                        return err
                }
                
                // Update user stats
                var stats models.UserDiscussionStats
                result := tx.Where("user_id = ?", topic.UserID).First(&stats)
                
                if result.Error != nil {
                        if result.Error == gorm.ErrRecordNotFound {
                                // Create new stats record
                                stats = models.UserDiscussionStats{
                                        UserID:        topic.UserID,
                                        TopicsCreated: 1,
                                        LastActivityAt: time.Now(),
                                }
                                return tx.Create(&stats).Error
                        }
                        return result.Error
                }
                
                // Update existing stats
                stats.TopicsCreated++
                stats.LastActivityAt = time.Now()
                return tx.Save(&stats).Error
        })
}

// UpdateTopic updates an existing topic
func (r *GormDiscussionRepository) UpdateTopic(topic *models.Topic) error {
        return r.db.Save(topic).Error
}

// DeleteTopic deletes a topic and its associated comments
func (r *GormDiscussionRepository) DeleteTopic(id uint) error {
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Delete comments first
                if err := tx.Where("topic_id = ?", id).Delete(&models.Comment{}).Error; err != nil {
                        return err
                }
                
                // Delete reactions
                if err := tx.Where("target_id = ? AND target_type = ?", id, "topic").Delete(&models.Reaction{}).Error; err != nil {
                        return err
                }
                
                // Delete subscriptions
                if err := tx.Where("topic_id = ?", id).Delete(&models.Subscription{}).Error; err != nil {
                        return err
                }
                
                // Delete mentions
                if err := tx.Where("target_id = ? AND target_type = ?", id, "topic").Delete(&models.Mention{}).Error; err != nil {
                        return err
                }
                
                // Delete topic tags
                if err := tx.Where("topic_id = ?", id).Delete(&models.TopicTag{}).Error; err != nil {
                        return err
                }
                
                // Delete the topic itself
                return tx.Delete(&models.Topic{}, id).Error
        })
}

// IncrementTopicViewCount increments the view count of a topic
func (r *GormDiscussionRepository) IncrementTopicViewCount(id uint) error {
        return r.db.Model(&models.Topic{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateTopicLastPostTime updates the last post time of a topic
func (r *GormDiscussionRepository) UpdateTopicLastPostTime(id uint) error {
        return r.db.Model(&models.Topic{}).Where("id = ?", id).UpdateColumn("last_post_at", time.Now()).Error
}

// PinTopic pins or unpins a topic
func (r *GormDiscussionRepository) PinTopic(id uint, pinned bool) error {
        return r.db.Model(&models.Topic{}).Where("id = ?", id).Update("is_pinned", pinned).Error
}

// LockTopic locks or unlocks a topic
func (r *GormDiscussionRepository) LockTopic(id uint, locked bool) error {
        return r.db.Model(&models.Topic{}).Where("id = ?", id).Update("is_locked", locked).Error
}

// GetAllTopics retrieves all topics with pagination
func (r *GormDiscussionRepository) GetAllTopics(page, pageSize int) ([]models.Topic, error) {
        var topics []models.Topic
        
        // Apply pagination
        offset := (page - 1) * pageSize
        result := r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&topics)
        
        return topics, result.Error
}

// GetCommentsByTopic retrieves comments for a specific topic
func (r *GormDiscussionRepository) GetCommentsByTopic(topicID uint, page, pageSize int) ([]models.Comment, int64, error) {
        var comments []models.Comment
        var total int64

        // Count total comments for this topic
        err := r.db.Model(&models.Comment{}).Where("topic_id = ? AND parent_id IS NULL", topicID).Count(&total).Error
        if err != nil {
                return nil, 0, err
        }

        // Calculate offset and fetch paginated results
        offset := (page - 1) * pageSize
        err = r.db.
                Where("topic_id = ? AND parent_id IS NULL", topicID).
                Order("created_at ASC").
                Offset(offset).
                Limit(pageSize).
                Find(&comments).Error

        if err != nil {
                return nil, 0, err
        }

        // For each top-level comment, count replies and load the first few replies
        for i := range comments {
                var replyCount int64
                err := r.db.Model(&models.Comment{}).Where("parent_id = ?", comments[i].ID).Count(&replyCount).Error
                if err != nil {
                        return nil, 0, err
                }
                comments[i].ReplyCount = int(replyCount)

                // Load replies (limited to first 5 for preview)
                err = r.db.
                        Where("parent_id = ?", comments[i].ID).
                        Order("created_at ASC").
                        Limit(5).
                        Find(&comments[i].Replies).Error
                if err != nil {
                        return nil, 0, err
                }
        }

        return comments, total, nil
}

// GetCommentByID retrieves a comment by its ID
func (r *GormDiscussionRepository) GetCommentByID(id uint) (*models.Comment, error) {
        var comment models.Comment
        err := r.db.First(&comment, id).Error
        if err != nil {
                return nil, err
        }
        return &comment, nil
}

// GetRepliesByComment retrieves all replies to a specific comment
func (r *GormDiscussionRepository) GetRepliesByComment(commentID uint) ([]models.Comment, error) {
        var replies []models.Comment
        err := r.db.
                Where("parent_id = ?", commentID).
                Order("created_at ASC").
                Find(&replies).Error
        return replies, err
}

// CreateComment creates a new comment
func (r *GormDiscussionRepository) CreateComment(comment *models.Comment) error {
        // Begin a transaction
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Create the comment
                if err := tx.Create(comment).Error; err != nil {
                        return err
                }
                
                // Update the topic's last post time
                if err := tx.Model(&models.Topic{}).Where("id = ?", comment.TopicID).Update("last_post_at", time.Now()).Error; err != nil {
                        return err
                }
                
                // Update user stats
                var stats models.UserDiscussionStats
                result := tx.Where("user_id = ?", comment.UserID).First(&stats)
                
                if result.Error != nil {
                        if result.Error == gorm.ErrRecordNotFound {
                                // Create new stats record
                                stats = models.UserDiscussionStats{
                                        UserID:         comment.UserID,
                                        CommentsPosted: 1,
                                        LastActivityAt: time.Now(),
                                }
                                return tx.Create(&stats).Error
                        }
                        return result.Error
                }
                
                // Update existing stats
                stats.CommentsPosted++
                stats.LastActivityAt = time.Now()
                return tx.Save(&stats).Error
        })
}

// UpdateComment updates an existing comment
func (r *GormDiscussionRepository) UpdateComment(comment *models.Comment) error {
        return r.db.Save(comment).Error
}

// DeleteComment deletes a comment
func (r *GormDiscussionRepository) DeleteComment(id uint) error {
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Get the comment to check if it has replies
                var comment models.Comment
                if err := tx.First(&comment, id).Error; err != nil {
                        return err
                }
                
                // Delete reactions
                if err := tx.Where("target_id = ? AND target_type = ?", id, "comment").Delete(&models.Reaction{}).Error; err != nil {
                        return err
                }
                
                // Delete mentions
                if err := tx.Where("target_id = ? AND target_type = ?", id, "comment").Delete(&models.Mention{}).Error; err != nil {
                        return err
                }
                
                // If it has replies, just soft delete by setting content to "[deleted]"
                var replyCount int64
                tx.Model(&models.Comment{}).Where("parent_id = ?", id).Count(&replyCount)
                
                if replyCount > 0 {
                        return tx.Model(&comment).Updates(map[string]interface{}{
                                "content":    "[deleted]",
                                "is_deleted": true,
                        }).Error
                }
                
                // If no replies, hard delete
                return tx.Delete(&models.Comment{}, id).Error
        })
}

// MarkCommentAsEdited marks a comment as edited
func (r *GormDiscussionRepository) MarkCommentAsEdited(id uint) error {
        now := time.Now()
        return r.db.Model(&models.Comment{}).Where("id = ?", id).Updates(map[string]interface{}{
                "is_edited": true,
                "edited_at": now,
        }).Error
}

// AddReaction adds a reaction to a topic or comment
func (r *GormDiscussionRepository) AddReaction(reaction *models.Reaction) error {
        // Check if the reaction already exists
        var count int64
        err := r.db.Model(&models.Reaction{}).
                Where("user_id = ? AND target_type = ? AND target_id = ? AND reaction_type = ?",
                        reaction.UserID, reaction.TargetType, reaction.TargetID, reaction.ReactionType).
                Count(&count).Error
        
        if err != nil {
                return err
        }
        
        if count > 0 {
                return models.ErrDuplicateReaction
        }
        
        // Begin a transaction
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Create the reaction
                if err := tx.Create(reaction).Error; err != nil {
                        return err
                }
                
                // Get the target creator's user ID
                var targetUserID uint
                if reaction.TargetType == "topic" {
                        var topic models.Topic
                        if err := tx.Select("user_id").First(&topic, reaction.TargetID).Error; err != nil {
                                return err
                        }
                        targetUserID = topic.UserID
                } else if reaction.TargetType == "comment" {
                        var comment models.Comment
                        if err := tx.Select("user_id").First(&comment, reaction.TargetID).Error; err != nil {
                                return err
                        }
                        targetUserID = comment.UserID
                }
                
                // Update reaction giver's stats
                var giverStats models.UserDiscussionStats
                result := tx.Where("user_id = ?", reaction.UserID).First(&giverStats)
                
                if result.Error != nil {
                        if result.Error == gorm.ErrRecordNotFound {
                                // Create new stats record
                                giverStats = models.UserDiscussionStats{
                                        UserID:         reaction.UserID,
                                        ReactionsGiven: 1,
                                        LastActivityAt: time.Now(),
                                }
                                if err := tx.Create(&giverStats).Error; err != nil {
                                        return err
                                }
                        } else {
                                return result.Error
                        }
                } else {
                        // Update existing stats
                        giverStats.ReactionsGiven++
                        giverStats.LastActivityAt = time.Now()
                        if err := tx.Save(&giverStats).Error; err != nil {
                                return err
                        }
                }
                
                // Update reaction receiver's stats
                if targetUserID != reaction.UserID { // Don't double-count if user reacts to own content
                        var receiverStats models.UserDiscussionStats
                        result := tx.Where("user_id = ?", targetUserID).First(&receiverStats)
                        
                        if result.Error != nil {
                                if result.Error == gorm.ErrRecordNotFound {
                                        // Create new stats record
                                        receiverStats = models.UserDiscussionStats{
                                                UserID:            targetUserID,
                                                ReactionsReceived: 1,
                                        }
                                        return tx.Create(&receiverStats).Error
                                }
                                return result.Error
                        }
                        
                        // Update existing stats
                        receiverStats.ReactionsReceived++
                        return tx.Save(&receiverStats).Error
                }
                
                return nil
        })
}

// RemoveReaction removes a reaction from a topic or comment
func (r *GormDiscussionRepository) RemoveReaction(userID, targetID uint, targetType, reactionType string) error {
        return r.db.Where(
                "user_id = ? AND target_type = ? AND target_id = ? AND reaction_type = ?",
                userID, targetType, targetID, reactionType,
        ).Delete(&models.Reaction{}).Error
}

// GetReactionsByTarget retrieves all reactions for a target
func (r *GormDiscussionRepository) GetReactionsByTarget(targetID uint, targetType string) ([]models.Reaction, error) {
        var reactions []models.Reaction
        err := r.db.
                Where("target_type = ? AND target_id = ?", targetType, targetID).
                Find(&reactions).Error
        return reactions, err
}

// GetReactionsSummary retrieves a summary of reactions for a target
func (r *GormDiscussionRepository) GetReactionsSummary(targetID uint, targetType string) ([]models.ReactionSummary, error) {
        var summaries []models.ReactionSummary
        
        // SQL query to aggregate reactions by type
        rows, err := r.db.Raw(`
                SELECT target_type, target_id, reaction_type, COUNT(*) as count
                FROM reactions
                WHERE target_type = ? AND target_id = ?
                GROUP BY reaction_type
                ORDER BY count DESC
        `, targetType, targetID).Rows()
        
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        
        // Parse the results
        for rows.Next() {
                var summary models.ReactionSummary
                if err := r.db.ScanRows(rows, &summary); err != nil {
                        return nil, err
                }
                summaries = append(summaries, summary)
        }
        
        return summaries, nil
}

// GetAllTags retrieves all tags
func (r *GormDiscussionRepository) GetAllTags() ([]models.Tag, error) {
        var tags []models.Tag
        err := r.db.Find(&tags).Error
        return tags, err
}

// GetTagByID retrieves a tag by its ID
func (r *GormDiscussionRepository) GetTagByID(id uint) (*models.Tag, error) {
        var tag models.Tag
        err := r.db.First(&tag, id).Error
        if err != nil {
                return nil, err
        }
        return &tag, nil
}

// GetTagBySlug retrieves a tag by its slug
func (r *GormDiscussionRepository) GetTagBySlug(slug string) (*models.Tag, error) {
        var tag models.Tag
        err := r.db.Where("slug = ?", slug).First(&tag).Error
        if err != nil {
                return nil, err
        }
        return &tag, nil
}

// CreateTag creates a new tag
func (r *GormDiscussionRepository) CreateTag(tag *models.Tag) error {
        return r.db.Create(tag).Error
}

// UpdateTag updates an existing tag
func (r *GormDiscussionRepository) UpdateTag(tag *models.Tag) error {
        return r.db.Save(tag).Error
}

// DeleteTag deletes a tag
func (r *GormDiscussionRepository) DeleteTag(id uint) error {
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Remove tag from all topics
                if err := tx.Delete(&models.TopicTag{}, "tag_id = ?", id).Error; err != nil {
                        return err
                }
                
                // Delete the tag
                return tx.Delete(&models.Tag{}, id).Error
        })
}

// AddTagToTopic adds a tag to a topic
func (r *GormDiscussionRepository) AddTagToTopic(topicID, tagID uint) error {
        // Check if the association already exists
        var count int64
        err := r.db.Model(&models.TopicTag{}).
                Where("topic_id = ? AND tag_id = ?", topicID, tagID).
                Count(&count).Error
        
        if err != nil {
                return err
        }
        
        if count > 0 {
                return nil // Tag already added
        }
        
        return r.db.Create(&models.TopicTag{
                TopicID: topicID,
                TagID:   tagID,
        }).Error
}

// RemoveTagFromTopic removes a tag from a topic
func (r *GormDiscussionRepository) RemoveTagFromTopic(topicID, tagID uint) error {
        return r.db.Where("topic_id = ? AND tag_id = ?", topicID, tagID).Delete(&models.TopicTag{}).Error
}

// GetTagsByTopic retrieves all tags for a topic
func (r *GormDiscussionRepository) GetTagsByTopic(topicID uint) ([]models.Tag, error) {
        var tags []models.Tag
        err := r.db.
                Joins("JOIN topic_tags ON topic_tags.tag_id = tags.id").
                Where("topic_tags.topic_id = ?", topicID).
                Find(&tags).Error
        return tags, err
}

// AddToModerationQueue adds an item to the moderation queue
func (r *GormDiscussionRepository) AddToModerationQueue(queue *models.ModerationQueue) error {
        // Check if the item is already in the queue
        var count int64
        err := r.db.Model(&models.ModerationQueue{}).
                Where("target_type = ? AND target_id = ? AND status = 'pending'", queue.TargetType, queue.TargetID).
                Count(&count).Error
        
        if err != nil {
                return err
        }
        
        if count > 0 {
                return nil // Already in queue
        }
        
        return r.db.Create(queue).Error
}

// GetModerationQueueItems retrieves items from the moderation queue
func (r *GormDiscussionRepository) GetModerationQueueItems(status string, page, pageSize int) ([]models.ModerationQueue, int64, error) {
        var items []models.ModerationQueue
        var total int64
        
        query := r.db.Model(&models.ModerationQueue{})
        
        if status != "" {
                query = query.Where("status = ?", status)
        }
        
        // Count total
        err := query.Count(&total).Error
        if err != nil {
                return nil, 0, err
        }
        
        // Calculate offset and fetch paginated results
        offset := (page - 1) * pageSize
        err = query.
                Order("created_at DESC").
                Offset(offset).
                Limit(pageSize).
                Find(&items).Error
        
        return items, total, err
}

// UpdateModerationQueueItem updates a moderation queue item
func (r *GormDiscussionRepository) UpdateModerationQueueItem(item *models.ModerationQueue) error {
        return r.db.Save(item).Error
}

// FlagComment flags a comment for moderation
func (r *GormDiscussionRepository) FlagComment(commentID, reportedBy uint, reason string) error {
        // Begin a transaction
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Mark the comment as flagged
                err := tx.Model(&models.Comment{}).
                        Where("id = ?", commentID).
                        Updates(map[string]interface{}{
                                "is_flagged":  true,
                                "flag_reason": reason,
                        }).Error
                
                if err != nil {
                        return err
                }
                
                // Add to moderation queue
                queue := &models.ModerationQueue{
                        TargetType: "comment",
                        TargetID:   commentID,
                        ReportedBy: reportedBy,
                        Reason:     reason,
                        Status:     "pending",
                }
                
                return r.AddToModerationQueue(queue)
        })
}

// ApproveFlaggedComment approves a flagged comment
func (r *GormDiscussionRepository) ApproveFlaggedComment(commentID, moderatorID uint, notes string) error {
        now := time.Now()
        
        // Begin a transaction
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Clear the flagged status
                err := tx.Model(&models.Comment{}).
                        Where("id = ?", commentID).
                        Updates(map[string]interface{}{
                                "is_flagged":  false,
                                "flag_reason": "",
                                "is_approved": true,
                        }).Error
                
                if err != nil {
                        return err
                }
                
                // Update moderation queue item
                return tx.Model(&models.ModerationQueue{}).
                        Where("target_type = ? AND target_id = ? AND status = 'pending'", "comment", commentID).
                        Updates(map[string]interface{}{
                                "status":          "approved",
                                "moderated_by":    moderatorID,
                                "moderated_at":    now,
                                "moderator_notes": notes,
                        }).Error
        })
}

// RejectFlaggedComment rejects a flagged comment
func (r *GormDiscussionRepository) RejectFlaggedComment(commentID, moderatorID uint, notes string) error {
        now := time.Now()
        
        // Begin a transaction
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Update the comment
                err := tx.Model(&models.Comment{}).
                        Where("id = ?", commentID).
                        Updates(map[string]interface{}{
                                "is_approved": false,
                                "content":     "[This comment has been removed by a moderator]",
                        }).Error
                
                if err != nil {
                        return err
                }
                
                // Update moderation queue item
                return tx.Model(&models.ModerationQueue{}).
                        Where("target_type = ? AND target_id = ? AND status = 'pending'", "comment", commentID).
                        Updates(map[string]interface{}{
                                "status":          "rejected",
                                "moderated_by":    moderatorID,
                                "moderated_at":    now,
                                "moderator_notes": notes,
                        }).Error
        })
}

// SubscribeToTopic subscribes a user to a topic
func (r *GormDiscussionRepository) SubscribeToTopic(userID, topicID uint, notifyOnReply bool, digestType string) error {
        // Check if subscription already exists
        var subscription models.Subscription
        result := r.db.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&subscription)
        
        if result.Error != nil {
                if result.Error == gorm.ErrRecordNotFound {
                        // Create new subscription
                        subscription = models.Subscription{
                                UserID:        userID,
                                TopicID:       topicID,
                                NotifyOnReply: notifyOnReply,
                                DigestType:    digestType,
                                IsActive:      true,
                        }
                        return r.db.Create(&subscription).Error
                }
                return result.Error
        }
        
        // Update existing subscription
        subscription.NotifyOnReply = notifyOnReply
        subscription.DigestType = digestType
        subscription.IsActive = true
        return r.db.Save(&subscription).Error
}

// UnsubscribeFromTopic unsubscribes a user from a topic
func (r *GormDiscussionRepository) UnsubscribeFromTopic(userID, topicID uint) error {
        return r.db.Model(&models.Subscription{}).
                Where("user_id = ? AND topic_id = ?", userID, topicID).
                Update("is_active", false).Error
}

// GetUserSubscriptions retrieves all active subscriptions for a user
func (r *GormDiscussionRepository) GetUserSubscriptions(userID uint) ([]models.Subscription, error) {
        var subscriptions []models.Subscription
        err := r.db.
                Where("user_id = ? AND is_active = true", userID).
                Find(&subscriptions).Error
        return subscriptions, err
}

// UpdateSubscriptionSettings updates a subscription's settings
func (r *GormDiscussionRepository) UpdateSubscriptionSettings(subscription *models.Subscription) error {
        return r.db.Save(subscription).Error
}

// GetUserDiscussionStats retrieves discussion stats for a user
func (r *GormDiscussionRepository) GetUserDiscussionStats(userID uint) (*models.UserDiscussionStats, error) {
        var stats models.UserDiscussionStats
        result := r.db.Where("user_id = ?", userID).First(&stats)
        
        if result.Error != nil {
                if result.Error == gorm.ErrRecordNotFound {
                        // Create a new stats record with zeros
                        stats = models.UserDiscussionStats{
                                UserID:          userID,
                                TopicsCreated:   0,
                                CommentsPosted:  0,
                                ReactionsGiven:  0,
                                ReactionsReceived: 0,
                                LastActivityAt: time.Now(),
                        }
                        err := r.db.Create(&stats).Error
                        if err != nil {
                                return nil, err
                        }
                        return &stats, nil
                }
                return nil, result.Error
        }
        
        return &stats, nil
}

// UpdateUserDiscussionStats updates discussion stats for a user
func (r *GormDiscussionRepository) UpdateUserDiscussionStats(stats *models.UserDiscussionStats) error {
        return r.db.Save(stats).Error
}

// GetNotificationPreferences retrieves notification preferences for a user
func (r *GormDiscussionRepository) GetNotificationPreferences(userID uint) (*models.NotificationPreference, error) {
        var preferences models.NotificationPreference
        result := r.db.Where("user_id = ?", userID).First(&preferences)
        
        if result.Error != nil {
                if result.Error == gorm.ErrRecordNotFound {
                        // Create default preferences
                        preferences = models.NotificationPreference{
                                UserID:         userID,
                                EmailOnReply:   true,
                                EmailOnMention: true,
                                EmailDigest:    true,
                                DigestFrequency: "weekly",
                        }
                        err := r.db.Create(&preferences).Error
                        if err != nil {
                                return nil, err
                        }
                        return &preferences, nil
                }
                return nil, result.Error
        }
        
        return &preferences, nil
}

// UpdateNotificationPreferences updates notification preferences for a user
func (r *GormDiscussionRepository) UpdateNotificationPreferences(preferences *models.NotificationPreference) error {
        return r.db.Save(preferences).Error
}

// RecordTopicView records a view of a topic by a user
func (r *GormDiscussionRepository) RecordTopicView(userID, topicID uint) error {
        // Check if viewed recently (last 24 hours)
        var count int64
        err := r.db.Model(&models.TopicView{}).
                Where("user_id = ? AND topic_id = ? AND viewed_at > ?", userID, topicID, time.Now().Add(-24*time.Hour)).
                Count(&count).Error
        
        if err != nil {
                return err
        }
        
        // If not viewed recently, create a new view record and increment topic view count
        if count == 0 {
                // Begin a transaction
                return r.db.Transaction(func(tx *gorm.DB) error {
                        // Create view record
                        view := &models.TopicView{
                                UserID:    userID,
                                TopicID:   topicID,
                                ViewedAt:  time.Now(),
                        }
                        if err := tx.Create(view).Error; err != nil {
                                return err
                        }
                        
                        // Increment topic view count
                        return tx.Model(&models.Topic{}).Where("id = ?", topicID).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
                })
        }
        
        return nil
}

// GetTopicViewCount retrieves the view count of a topic
func (r *GormDiscussionRepository) GetTopicViewCount(topicID uint) (int, error) {
        var topic models.Topic
        err := r.db.Select("view_count").First(&topic, topicID).Error
        if err != nil {
                return 0, err
        }
        return topic.ViewCount, nil
}

// GetRecentlyViewedTopics retrieves recently viewed topics for a user
func (r *GormDiscussionRepository) GetRecentlyViewedTopics(userID uint, limit int) ([]models.Topic, error) {
        var topics []models.Topic
        
        err := r.db.Raw(`
                SELECT t.*
                FROM topics t
                INNER JOIN (
                        SELECT topic_id, MAX(viewed_at) as last_viewed
                        FROM topic_views
                        WHERE user_id = ?
                        GROUP BY topic_id
                        ORDER BY last_viewed DESC
                        LIMIT ?
                ) v ON t.id = v.topic_id
                ORDER BY v.last_viewed DESC
        `, userID, limit).Scan(&topics).Error
        
        return topics, err
}

// CreateMention creates a new mention
func (r *GormDiscussionRepository) CreateMention(mention *models.Mention) error {
        return r.db.Create(mention).Error
}

// GetMentionsForUser retrieves mentions for a user
func (r *GormDiscussionRepository) GetMentionsForUser(userID uint, isRead bool, page, pageSize int) ([]models.Mention, int64, error) {
        var mentions []models.Mention
        var total int64
        
        query := r.db.Model(&models.Mention{}).Where("user_id = ?", userID)
        
        // Filter by read/unread if specified
        if isRead {
                query = query.Where("is_read = true")
        } else {
                query = query.Where("is_read = false")
        }
        
        // Count total
        err := query.Count(&total).Error
        if err != nil {
                return nil, 0, err
        }
        
        // Calculate offset and fetch paginated results
        offset := (page - 1) * pageSize
        err = query.
                Order("created_at DESC").
                Offset(offset).
                Limit(pageSize).
                Find(&mentions).Error
        
        return mentions, total, err
}

// MarkMentionAsRead marks a mention as read
func (r *GormDiscussionRepository) MarkMentionAsRead(id uint) error {
        return r.db.Model(&models.Mention{}).Where("id = ?", id).Update("is_read", true).Error
}

// User Repository methods

// GetUserByID retrieves a user by their ID
func (r *GormDiscussionRepository) GetUserByID(id uint) (*models.User, error) {
        var user models.User
        err := r.db.First(&user, id).Error
        if err != nil {
                return nil, err
        }
        return &user, nil
}

// GetUserByUsername retrieves a user by their username
func (r *GormDiscussionRepository) GetUserByUsername(username string) (*models.User, error) {
        var user models.User
        err := r.db.Where("username = ?", username).First(&user).Error
        if err != nil {
                return nil, err
        }
        return &user, nil
}

// CreateUser creates a new user
func (r *GormDiscussionRepository) CreateUser(user *models.User) error {
        return r.db.Create(user).Error
}

// UpdateUser updates an existing user
func (r *GormDiscussionRepository) UpdateUser(user *models.User) error {
        return r.db.Save(user).Error
}

// GetUserTrustScore retrieves a user's trust score
func (r *GormDiscussionRepository) GetUserTrustScore(userID uint) (*models.UserTrustScore, error) {
        var score models.UserTrustScore
        err := r.db.Where("user_id = ?", userID).First(&score).Error
        if err != nil {
                if err == gorm.ErrRecordNotFound {
                        // Return default score if not found
                        return &models.UserTrustScore{
                                UserID: userID,
                                Score: 50.0, // Default trust score
                        }, nil
                }
                return nil, err
        }
        return &score, nil
}

// UpdateUserTrustScore updates a user's trust score
func (r *GormDiscussionRepository) UpdateUserTrustScore(score *models.UserTrustScore) error {
        return r.db.Save(score).Error
}

// Content Link Repository methods

// CreateTopicContentLink creates a new link between a topic and content
func (r *GormDiscussionRepository) CreateTopicContentLink(link *models.TopicContentLink) error {
        return r.db.Create(link).Error
}

// UpdateTopicContentLink updates an existing topic-content link
func (r *GormDiscussionRepository) UpdateTopicContentLink(link *models.TopicContentLink) error {
        return r.db.Save(link).Error
}

// DeleteTopicContentLink deletes a topic-content link
func (r *GormDiscussionRepository) DeleteTopicContentLink(id uint) error {
        return r.db.Delete(&models.TopicContentLink{}, id).Error
}

// GetTopicContentLinks retrieves content links for a topic
func (r *GormDiscussionRepository) GetTopicContentLinks(topicID uint) ([]models.TopicContentLink, error) {
        var links []models.TopicContentLink
        err := r.db.Where("topic_id = ?", topicID).Find(&links).Error
        return links, err
}

// GetTopicsByContent retrieves topics linked to a content
func (r *GormDiscussionRepository) GetTopicsByContent(contentType models.ContentReferenceType, contentID uint) ([]models.TopicContentLink, error) {
        var links []models.TopicContentLink
        err := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).Find(&links).Error
        return links, err
}

// CreateCommentContentLink creates a new link between a comment and content
func (r *GormDiscussionRepository) CreateCommentContentLink(link *models.CommentContentLink) error {
        return r.db.Create(link).Error
}

// UpdateCommentContentLink updates an existing comment-content link
func (r *GormDiscussionRepository) UpdateCommentContentLink(link *models.CommentContentLink) error {
        return r.db.Save(link).Error
}

// DeleteCommentContentLink deletes a comment-content link
func (r *GormDiscussionRepository) DeleteCommentContentLink(id uint) error {
        return r.db.Delete(&models.CommentContentLink{}, id).Error
}

// GetCommentContentLinks retrieves content links for a comment
func (r *GormDiscussionRepository) GetCommentContentLinks(commentID uint) ([]models.CommentContentLink, error) {
        var links []models.CommentContentLink
        err := r.db.Where("comment_id = ?", commentID).Find(&links).Error
        return links, err
}

// GetCommentsByContent retrieves comments linked to a content
func (r *GormDiscussionRepository) GetCommentsByContent(contentType models.ContentReferenceType, contentID uint) ([]models.CommentContentLink, error) {
        var links []models.CommentContentLink
        err := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).Find(&links).Error
        return links, err
}

// CreateAutoGeneratedTopicTemplate creates a new auto-generated topic template
func (r *GormDiscussionRepository) CreateAutoGeneratedTopicTemplate(template *models.AutoGeneratedTopicTemplate) error {
        return r.db.Create(template).Error
}

// UpdateAutoGeneratedTopicTemplate updates an existing auto-generated topic template
func (r *GormDiscussionRepository) UpdateAutoGeneratedTopicTemplate(template *models.AutoGeneratedTopicTemplate) error {
        return r.db.Save(template).Error
}

// DeleteAutoGeneratedTopicTemplate deletes an auto-generated topic template
func (r *GormDiscussionRepository) DeleteAutoGeneratedTopicTemplate(id uint) error {
        return r.db.Delete(&models.AutoGeneratedTopicTemplate{}, id).Error
}

// GetAutoGeneratedTopicTemplateByID retrieves an auto-generated topic template by ID
func (r *GormDiscussionRepository) GetAutoGeneratedTopicTemplateByID(id uint) (*models.AutoGeneratedTopicTemplate, error) {
        var template models.AutoGeneratedTopicTemplate
        err := r.db.First(&template, id).Error
        if err != nil {
                return nil, err
        }
        return &template, nil
}

// GetAutoGeneratedTopicTemplates retrieves auto-generated topic templates for a content type
func (r *GormDiscussionRepository) GetAutoGeneratedTopicTemplates(contentType models.ContentReferenceType) ([]models.AutoGeneratedTopicTemplate, error) {
        var templates []models.AutoGeneratedTopicTemplate
        err := r.db.Where("content_type = ?", contentType).Find(&templates).Error
        return templates, err
}

// GetActiveAutoGeneratedTopicTemplates retrieves active auto-generated topic templates for a content type
func (r *GormDiscussionRepository) GetActiveAutoGeneratedTopicTemplates(contentType models.ContentReferenceType) ([]models.AutoGeneratedTopicTemplate, error) {
        var templates []models.AutoGeneratedTopicTemplate
        err := r.db.Where("content_type = ? AND is_active = ?", contentType, true).Find(&templates).Error
        return templates, err
}

// CreateContentDiscussionRecommendation creates a new content-discussion recommendation
func (r *GormDiscussionRepository) CreateContentDiscussionRecommendation(recommendation *models.ContentDiscussionRecommendation) error {
        return r.db.Create(recommendation).Error
}

// UpdateContentDiscussionRecommendation updates an existing content-discussion recommendation
func (r *GormDiscussionRepository) UpdateContentDiscussionRecommendation(recommendation *models.ContentDiscussionRecommendation) error {
        return r.db.Save(recommendation).Error
}

// DeleteContentDiscussionRecommendation deletes a content-discussion recommendation
func (r *GormDiscussionRepository) DeleteContentDiscussionRecommendation(id uint) error {
        return r.db.Delete(&models.ContentDiscussionRecommendation{}, id).Error
}

// GetContentDiscussionRecommendations retrieves content-discussion recommendations
func (r *GormDiscussionRepository) GetContentDiscussionRecommendations(contentType models.ContentReferenceType, contentID uint, limit int) ([]models.ContentDiscussionRecommendation, error) {
        var recommendations []models.ContentDiscussionRecommendation
        query := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).Order("relevance_score DESC")
        
        if limit > 0 {
                query = query.Limit(limit)
        }
        
        err := query.Find(&recommendations).Error
        return recommendations, err
}