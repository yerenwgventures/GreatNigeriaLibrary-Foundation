package service

import (
        "strings"
        "time"

        "github.com/gosimple/slug"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/repository"
        "gorm.io/gorm"
)

// DiscussionService defines the interface for discussion-related business logic
type DiscussionService interface {
        // Categories
        GetCategories() ([]models.Category, error)
        GetCategoryByID(id uint) (*models.Category, error)
        GetCategoryBySlug(slug string) (*models.Category, error)
        CreateCategory(name, description string, parentID *uint, sortOrder int) (*models.Category, error)
        UpdateCategory(id uint, name, description string, parentID *uint, sortOrder int, isActive bool) (*models.Category, error)
        DeleteCategory(id uint) error
        
        // Category Configuration
        GetCategoryConfig(categoryID uint) (*models.CategoryConfig, error)
        CreateOrUpdateCategoryConfig(categoryID uint, config models.CategoryConfig) (*models.CategoryConfig, error)
        
        // Posting Rules
        GetPostingRules(categoryID uint) (*models.PostingRules, error)
        CreateOrUpdatePostingRules(categoryID uint, rules models.PostingRules) (*models.PostingRules, error)
        
        // Auto Moderation Settings
        GetAutoModerationSettings(categoryID uint) (*models.AutoModerationSettings, error)
        CreateOrUpdateAutoModerationSettings(categoryID uint, settings models.AutoModerationSettings) (*models.AutoModerationSettings, error)
        
        // Category Moderators
        GetCategoryModerators(categoryID uint) ([]models.CategoryModerator, error)
        AddCategoryModerator(categoryID, userID uint, permissions models.CategoryModerator) (*models.CategoryModerator, error)
        UpdateCategoryModerator(categoryID, userID uint, permissions models.CategoryModerator) (*models.CategoryModerator, error)
        RemoveCategoryModerator(categoryID, userID uint) error

        // Topics
        GetTopics(page, pageSize int, filters map[string]interface{}) ([]models.Topic, int64, error)
        GetTopicByID(id uint) (*models.Topic, error)
        GetTopicsByCategory(categoryID uint, page, pageSize int) ([]models.Topic, int64, error)
        GetTopicsByUser(userID uint, page, pageSize int) ([]models.Topic, int64, error)
        GetTopicsByBookSection(bookID, chapterID, sectionID uint) ([]models.Topic, error)
        CreateTopic(userID, categoryID uint, title, content string, bookID, chapterID, sectionID *uint, tagIDs []uint) (*models.Topic, error)
        UpdateTopic(id, userID uint, title, content string, categoryID uint) (*models.Topic, error)
        DeleteTopic(id, userID uint, isAdmin bool) error
        PinTopic(id uint, pinned bool) error
        LockTopic(id uint, locked bool) error
        ViewTopic(id, userID uint) error

        // Comments
        GetCommentsByTopic(topicID uint, page, pageSize int) ([]models.Comment, int64, error)
        GetCommentByID(id uint) (*models.Comment, error)
        GetRepliesByComment(commentID uint) ([]models.Comment, error)
        CreateComment(userID, topicID uint, content string, parentID *uint) (*models.Comment, error)
        UpdateComment(id, userID uint, content string, isAdmin bool) (*models.Comment, error)
        DeleteComment(id, userID uint, isAdmin bool) error

        // Reactions
        AddReaction(userID, targetID uint, targetType, reactionType string) error
        RemoveReaction(userID, targetID uint, targetType, reactionType string) error
        GetReactionsSummary(targetID uint, targetType string) ([]models.ReactionSummary, error)

        // Tags
        GetAllTags() ([]models.Tag, error)
        CreateTag(name string, color string, isSystem bool) (*models.Tag, error)
        AddTagToTopic(topicID, tagID uint) error
        RemoveTagFromTopic(topicID, tagID uint) error
        GetTagsByTopic(topicID uint) ([]models.Tag, error)

        // User Stats
        GetUserDiscussionStats(userID uint) (*models.UserDiscussionStats, error)
}

// DiscussionServiceImpl implements the DiscussionService interface
type DiscussionServiceImpl struct {
        discussionRepo repository.DiscussionRepository
}

// NewDiscussionService creates a new discussion service instance
func NewDiscussionService(discussionRepo repository.DiscussionRepository) DiscussionService {
        return &DiscussionServiceImpl{
                discussionRepo: discussionRepo,
        }
}

// GetCategories retrieves all categories
func (s *DiscussionServiceImpl) GetCategories() ([]models.Category, error) {
        return s.discussionRepo.GetCategories()
}

// GetCategoryByID retrieves a category by its ID
func (s *DiscussionServiceImpl) GetCategoryByID(id uint) (*models.Category, error) {
        return s.discussionRepo.GetCategoryByID(id)
}

// GetCategoryBySlug retrieves a category by its slug
func (s *DiscussionServiceImpl) GetCategoryBySlug(slug string) (*models.Category, error) {
        return s.discussionRepo.GetCategoryBySlug(slug)
}

// CreateCategory creates a new category
func (s *DiscussionServiceImpl) CreateCategory(name, description string, parentID *uint, sortOrder int) (*models.Category, error) {
        // Generate slug from name
        slugValue := slug.Make(name)
        
        // Create category
        category := &models.Category{
                Name:        name,
                Description: description,
                Slug:        slugValue,
                ParentID:    parentID,
                SortOrder:   sortOrder,
                IsActive:    true,
        }
        
        err := s.discussionRepo.CreateCategory(category)
        if err != nil {
                return nil, err
        }
        
        return category, nil
}

// UpdateCategory updates an existing category
func (s *DiscussionServiceImpl) UpdateCategory(id uint, name, description string, parentID *uint, sortOrder int, isActive bool) (*models.Category, error) {
        // Get the category
        category, err := s.discussionRepo.GetCategoryByID(id)
        if err != nil {
                return nil, err
        }
        
        // Update category fields
        category.Name = name
        category.Description = description
        category.ParentID = parentID
        category.SortOrder = sortOrder
        category.IsActive = isActive
        
        // Only update slug if name changed
        if name != category.Name {
                category.Slug = slug.Make(name)
        }
        
        err = s.discussionRepo.UpdateCategory(category)
        if err != nil {
                return nil, err
        }
        
        return category, nil
}

// DeleteCategory deletes a category
func (s *DiscussionServiceImpl) DeleteCategory(id uint) error {
        return s.discussionRepo.DeleteCategory(id)
}

// GetCategoryConfig retrieves configuration for a category
func (s *DiscussionServiceImpl) GetCategoryConfig(categoryID uint) (*models.CategoryConfig, error) {
        return s.discussionRepo.GetCategoryConfig(categoryID)
}

// CreateOrUpdateCategoryConfig creates or updates category configuration
func (s *DiscussionServiceImpl) CreateOrUpdateCategoryConfig(categoryID uint, config models.CategoryConfig) (*models.CategoryConfig, error) {
        // Ensure the category exists
        _, err := s.discussionRepo.GetCategoryByID(categoryID)
        if err != nil {
                return nil, models.ErrCategoryNotFound
        }
        
        // Set the category ID in the config
        config.CategoryID = categoryID
        
        // Check if config already exists
        existingConfig, err := s.discussionRepo.GetCategoryConfig(categoryID)
        if err != nil && err != gorm.ErrRecordNotFound {
                return nil, err
        }
        
        if err == gorm.ErrRecordNotFound || existingConfig.ID == 0 {
                // Create new config
                err = s.discussionRepo.CreateCategoryConfig(&config)
                if err != nil {
                        return nil, err
                }
        } else {
                // Update existing config
                config.ID = existingConfig.ID
                err = s.discussionRepo.UpdateCategoryConfig(&config)
                if err != nil {
                        return nil, err
                }
        }
        
        return &config, nil
}

// GetPostingRules retrieves posting rules for a category
func (s *DiscussionServiceImpl) GetPostingRules(categoryID uint) (*models.PostingRules, error) {
        return s.discussionRepo.GetPostingRules(categoryID)
}

// CreateOrUpdatePostingRules creates or updates posting rules for a category
func (s *DiscussionServiceImpl) CreateOrUpdatePostingRules(categoryID uint, rules models.PostingRules) (*models.PostingRules, error) {
        // Ensure the category exists
        _, err := s.discussionRepo.GetCategoryByID(categoryID)
        if err != nil {
                return nil, models.ErrCategoryNotFound
        }
        
        // Set the category ID in the rules
        rules.CategoryID = categoryID
        
        // Check if rules already exist
        existingRules, err := s.discussionRepo.GetPostingRules(categoryID)
        if err != nil && err != gorm.ErrRecordNotFound {
                return nil, err
        }
        
        if err == gorm.ErrRecordNotFound || existingRules.ID == 0 {
                // Create new rules
                err = s.discussionRepo.CreatePostingRules(&rules)
                if err != nil {
                        return nil, err
                }
        } else {
                // Update existing rules
                rules.ID = existingRules.ID
                err = s.discussionRepo.UpdatePostingRules(&rules)
                if err != nil {
                        return nil, err
                }
        }
        
        return &rules, nil
}

// GetAutoModerationSettings retrieves auto-moderation settings for a category
func (s *DiscussionServiceImpl) GetAutoModerationSettings(categoryID uint) (*models.AutoModerationSettings, error) {
        return s.discussionRepo.GetAutoModerationSettings(categoryID)
}

// CreateOrUpdateAutoModerationSettings creates or updates auto-moderation settings for a category
func (s *DiscussionServiceImpl) CreateOrUpdateAutoModerationSettings(categoryID uint, settings models.AutoModerationSettings) (*models.AutoModerationSettings, error) {
        // Ensure the category exists
        _, err := s.discussionRepo.GetCategoryByID(categoryID)
        if err != nil {
                return nil, models.ErrCategoryNotFound
        }
        
        // Set the category ID in the settings
        settings.CategoryID = categoryID
        
        // Check if settings already exist
        existingSettings, err := s.discussionRepo.GetAutoModerationSettings(categoryID)
        if err != nil && err != gorm.ErrRecordNotFound {
                return nil, err
        }
        
        if err == gorm.ErrRecordNotFound || existingSettings.ID == 0 {
                // Create new settings
                err = s.discussionRepo.CreateAutoModerationSettings(&settings)
                if err != nil {
                        return nil, err
                }
        } else {
                // Update existing settings
                settings.ID = existingSettings.ID
                err = s.discussionRepo.UpdateAutoModerationSettings(&settings)
                if err != nil {
                        return nil, err
                }
        }
        
        return &settings, nil
}

// GetCategoryModerators retrieves moderators for a category
func (s *DiscussionServiceImpl) GetCategoryModerators(categoryID uint) ([]models.CategoryModerator, error) {
        return s.discussionRepo.GetCategoryModerators(categoryID)
}

// AddCategoryModerator adds a moderator to a category
func (s *DiscussionServiceImpl) AddCategoryModerator(categoryID, userID uint, permissions models.CategoryModerator) (*models.CategoryModerator, error) {
        // Ensure the category exists
        _, err := s.discussionRepo.GetCategoryByID(categoryID)
        if err != nil {
                return nil, models.ErrCategoryNotFound
        }
        
        // Set category and user IDs
        permissions.CategoryID = categoryID
        permissions.UserID = userID
        
        // Add moderator
        err = s.discussionRepo.AddCategoryModerator(&permissions)
        if err != nil {
                return nil, err
        }
        
        return &permissions, nil
}

// UpdateCategoryModerator updates a moderator's permissions
func (s *DiscussionServiceImpl) UpdateCategoryModerator(categoryID, userID uint, permissions models.CategoryModerator) (*models.CategoryModerator, error) {
        // Ensure the category exists
        _, err := s.discussionRepo.GetCategoryByID(categoryID)
        if err != nil {
                return nil, models.ErrCategoryNotFound
        }
        
        // Set category and user IDs and update
        permissions.CategoryID = categoryID
        permissions.UserID = userID
        
        err = s.discussionRepo.UpdateCategoryModerator(&permissions)
        if err != nil {
                return nil, err
        }
        
        return &permissions, nil
}

// RemoveCategoryModerator removes a moderator from a category
func (s *DiscussionServiceImpl) RemoveCategoryModerator(categoryID, userID uint) error {
        return s.discussionRepo.RemoveCategoryModerator(categoryID, userID)
}

// GetTopics retrieves topics with pagination and filtering
func (s *DiscussionServiceImpl) GetTopics(page, pageSize int, filters map[string]interface{}) ([]models.Topic, int64, error) {
        return s.discussionRepo.GetTopics(page, pageSize, filters)
}

// GetTopicByID retrieves a topic by its ID
func (s *DiscussionServiceImpl) GetTopicByID(id uint) (*models.Topic, error) {
        return s.discussionRepo.GetTopicByID(id)
}

// GetTopicsByCategory retrieves topics in a specific category
func (s *DiscussionServiceImpl) GetTopicsByCategory(categoryID uint, page, pageSize int) ([]models.Topic, int64, error) {
        return s.discussionRepo.GetTopicsByCategory(categoryID, page, pageSize)
}

// GetTopicsByUser retrieves topics created by a specific user
func (s *DiscussionServiceImpl) GetTopicsByUser(userID uint, page, pageSize int) ([]models.Topic, int64, error) {
        return s.discussionRepo.GetTopicsByUser(userID, page, pageSize)
}

// GetTopicsByBookSection retrieves topics related to a specific book section
func (s *DiscussionServiceImpl) GetTopicsByBookSection(bookID, chapterID, sectionID uint) ([]models.Topic, error) {
        return s.discussionRepo.GetTopicsByBookSection(bookID, chapterID, sectionID)
}

// CreateTopic creates a new topic
func (s *DiscussionServiceImpl) CreateTopic(userID, categoryID uint, title, content string, bookID, chapterID, sectionID *uint, tagIDs []uint) (*models.Topic, error) {
        // Create topic
        topic := &models.Topic{
                Title:      title,
                Content:    content,
                UserID:     userID,
                CategoryID: categoryID,
                BookID:     bookID,
                ChapterID:  chapterID,
                SectionID:  sectionID,
                LastPostAt: time.Now(),
        }
        
        err := s.discussionRepo.CreateTopic(topic)
        if err != nil {
                return nil, err
        }
        
        // Add tags if provided
        if len(tagIDs) > 0 {
                for _, tagID := range tagIDs {
                        err := s.discussionRepo.AddTagToTopic(topic.ID, tagID)
                        if err != nil {
                                // Log error but continue (non-critical)
                                // We could consider rolling back the transaction, but for now we'll continue
                        }
                }
        }
        
        return topic, nil
}

// UpdateTopic updates an existing topic
func (s *DiscussionServiceImpl) UpdateTopic(id, userID uint, title, content string, categoryID uint) (*models.Topic, error) {
        // Get the topic
        topic, err := s.discussionRepo.GetTopicByID(id)
        if err != nil {
                return nil, err
        }
        
        // Check if the topic is locked
        if topic.IsLocked {
                return nil, models.ErrTopicLocked
        }
        
        // Check if the user is the topic creator
        if topic.UserID != userID {
                return nil, models.ErrPermissionDenied
        }
        
        // Update topic fields
        topic.Title = title
        topic.Content = content
        topic.CategoryID = categoryID
        
        err = s.discussionRepo.UpdateTopic(topic)
        if err != nil {
                return nil, err
        }
        
        return topic, nil
}

// DeleteTopic deletes a topic
func (s *DiscussionServiceImpl) DeleteTopic(id, userID uint, isAdmin bool) error {
        // Get the topic
        topic, err := s.discussionRepo.GetTopicByID(id)
        if err != nil {
                return err
        }
        
        // Check if the user is the topic creator or an admin
        if topic.UserID != userID && !isAdmin {
                return models.ErrPermissionDenied
        }
        
        return s.discussionRepo.DeleteTopic(id)
}

// PinTopic pins or unpins a topic
func (s *DiscussionServiceImpl) PinTopic(id uint, pinned bool) error {
        return s.discussionRepo.PinTopic(id, pinned)
}

// LockTopic locks or unlocks a topic
func (s *DiscussionServiceImpl) LockTopic(id uint, locked bool) error {
        return s.discussionRepo.LockTopic(id, locked)
}

// ViewTopic increments the view count of a topic
func (s *DiscussionServiceImpl) ViewTopic(id, userID uint) error {
        return s.discussionRepo.RecordTopicView(userID, id)
}

// GetCommentsByTopic retrieves comments for a specific topic
func (s *DiscussionServiceImpl) GetCommentsByTopic(topicID uint, page, pageSize int) ([]models.Comment, int64, error) {
        return s.discussionRepo.GetCommentsByTopic(topicID, page, pageSize)
}

// GetCommentByID retrieves a comment by its ID
func (s *DiscussionServiceImpl) GetCommentByID(id uint) (*models.Comment, error) {
        return s.discussionRepo.GetCommentByID(id)
}

// GetRepliesByComment retrieves all replies to a specific comment
func (s *DiscussionServiceImpl) GetRepliesByComment(commentID uint) ([]models.Comment, error) {
        return s.discussionRepo.GetRepliesByComment(commentID)
}

// CreateComment creates a new comment
func (s *DiscussionServiceImpl) CreateComment(userID, topicID uint, content string, parentID *uint) (*models.Comment, error) {
        // Get the topic to check if it's locked
        topic, err := s.discussionRepo.GetTopicByID(topicID)
        if err != nil {
                return nil, err
        }
        
        if topic.IsLocked {
                return nil, models.ErrTopicLocked
        }
        
        // Check if parent comment exists if parentID is provided
        if parentID != nil {
                _, err := s.discussionRepo.GetCommentByID(*parentID)
                if err != nil {
                        return nil, models.ErrCommentNotFound
                }
        }
        
        // Create comment
        comment := &models.Comment{
                Content:  content,
                UserID:   userID,
                TopicID:  topicID,
                ParentID: parentID,
        }
        
        err = s.discussionRepo.CreateComment(comment)
        if err != nil {
                return nil, err
        }
        
        // Parse mentions (@username) and create mention records
        // This is a simplified implementation
        words := strings.Fields(content)
        for _, word := range words {
                if strings.HasPrefix(word, "@") {
                        // In a real implementation, we would look up the user ID from the username
                        // and create a mention record
                        _ = strings.TrimPrefix(word, "@") // Using _ to acknowledge we're extracting but not using yet
                        // For now, we'll skip this functionality since we don't have user lookup
                }
        }
        
        return comment, nil
}

// UpdateComment updates an existing comment
func (s *DiscussionServiceImpl) UpdateComment(id, userID uint, content string, isAdmin bool) (*models.Comment, error) {
        // Get the comment
        comment, err := s.discussionRepo.GetCommentByID(id)
        if err != nil {
                return nil, err
        }
        
        // Check if the user is the comment creator or an admin
        if comment.UserID != userID && !isAdmin {
                return nil, models.ErrPermissionDenied
        }
        
        // Update comment fields
        comment.Content = content
        
        // Mark as edited
        err = s.discussionRepo.MarkCommentAsEdited(id)
        if err != nil {
                return nil, err
        }
        
        err = s.discussionRepo.UpdateComment(comment)
        if err != nil {
                return nil, err
        }
        
        return comment, nil
}

// DeleteComment deletes a comment
func (s *DiscussionServiceImpl) DeleteComment(id, userID uint, isAdmin bool) error {
        // Get the comment
        comment, err := s.discussionRepo.GetCommentByID(id)
        if err != nil {
                return err
        }
        
        // Check if the user is the comment creator or an admin
        if comment.UserID != userID && !isAdmin {
                return models.ErrPermissionDenied
        }
        
        return s.discussionRepo.DeleteComment(id)
}

// AddReaction adds a reaction to a topic or comment
func (s *DiscussionServiceImpl) AddReaction(userID, targetID uint, targetType, reactionType string) error {
        // Validate target type
        if targetType != "topic" && targetType != "comment" {
                return models.ErrInvalidContent
        }
        
        // Validate target exists
        if targetType == "topic" {
                _, err := s.discussionRepo.GetTopicByID(targetID)
                if err != nil {
                        return models.ErrTopicNotFound
                }
        } else if targetType == "comment" {
                _, err := s.discussionRepo.GetCommentByID(targetID)
                if err != nil {
                        return models.ErrCommentNotFound
                }
        }
        
        // Create reaction
        reaction := &models.Reaction{
                UserID:       userID,
                TargetType:   targetType,
                TargetID:     targetID,
                ReactionType: reactionType,
        }
        
        return s.discussionRepo.AddReaction(reaction)
}

// RemoveReaction removes a reaction from a topic or comment
func (s *DiscussionServiceImpl) RemoveReaction(userID, targetID uint, targetType, reactionType string) error {
        return s.discussionRepo.RemoveReaction(userID, targetID, targetType, reactionType)
}

// GetReactionsSummary retrieves a summary of reactions for a target
func (s *DiscussionServiceImpl) GetReactionsSummary(targetID uint, targetType string) ([]models.ReactionSummary, error) {
        return s.discussionRepo.GetReactionsSummary(targetID, targetType)
}

// GetAllTags retrieves all tags
func (s *DiscussionServiceImpl) GetAllTags() ([]models.Tag, error) {
        return s.discussionRepo.GetAllTags()
}

// CreateTag creates a new tag
func (s *DiscussionServiceImpl) CreateTag(name string, color string, isSystem bool) (*models.Tag, error) {
        // Generate slug from name
        slugValue := slug.Make(name)
        
        // If color not provided, use default
        if color == "" {
                color = "#3498db"
        }
        
        // Create tag
        tag := &models.Tag{
                Name:     name,
                Slug:     slugValue,
                Color:    color,
                IsSystem: isSystem,
        }
        
        err := s.discussionRepo.CreateTag(tag)
        if err != nil {
                return nil, err
        }
        
        return tag, nil
}

// AddTagToTopic adds a tag to a topic
func (s *DiscussionServiceImpl) AddTagToTopic(topicID, tagID uint) error {
        return s.discussionRepo.AddTagToTopic(topicID, tagID)
}

// RemoveTagFromTopic removes a tag from a topic
func (s *DiscussionServiceImpl) RemoveTagFromTopic(topicID, tagID uint) error {
        return s.discussionRepo.RemoveTagFromTopic(topicID, tagID)
}

// GetTagsByTopic retrieves all tags for a topic
func (s *DiscussionServiceImpl) GetTagsByTopic(topicID uint) ([]models.Tag, error) {
        return s.discussionRepo.GetTagsByTopic(topicID)
}

// GetUserDiscussionStats retrieves discussion stats for a user
func (s *DiscussionServiceImpl) GetUserDiscussionStats(userID uint) (*models.UserDiscussionStats, error) {
        return s.discussionRepo.GetUserDiscussionStats(userID)
}