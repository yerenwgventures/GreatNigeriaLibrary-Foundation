package repository

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/models"
)

// TopicRepository defines the interface for topic-related database operations
type TopicRepository interface {
        GetTopicByID(id uint) (*models.Topic, error)
        CreateTopic(topic *models.Topic) error
        UpdateTopic(topic *models.Topic) error
        DeleteTopic(id uint) error
        GetTopics(page, pageSize int, filters map[string]interface{}) ([]models.Topic, int64, error)
        GetTopicsByCategory(categoryID uint, page, pageSize int) ([]models.Topic, int64, error)
        GetTopicsByUser(userID uint, page, pageSize int) ([]models.Topic, int64, error)
        GetTopicsByBookSection(bookID, chapterID, sectionID uint) ([]models.Topic, error)
        GetAllTopics(page, pageSize int) ([]models.Topic, error)
        IncrementTopicViewCount(id uint) error
        UpdateTopicLastPostTime(id uint) error
        PinTopic(id uint, pinned bool) error
        LockTopic(id uint, locked bool) error
}

// CommentRepository defines the interface for comment-related database operations
type CommentRepository interface {
        GetCommentByID(id uint) (*models.Comment, error)
        CreateComment(comment *models.Comment) error
        UpdateComment(comment *models.Comment) error
        DeleteComment(id uint) error
        GetCommentsByTopic(topicID uint, page, pageSize int) ([]models.Comment, int64, error)
        GetRepliesByComment(commentID uint) ([]models.Comment, error)
        MarkCommentAsEdited(id uint) error
}

// CategoryRepository defines the interface for category-related database operations
type CategoryRepository interface {
        GetCategories() ([]models.Category, error)
        GetCategoryByID(id uint) (*models.Category, error)
        GetCategoryBySlug(slug string) (*models.Category, error)
        CreateCategory(category *models.Category) error
        UpdateCategory(category *models.Category) error
        DeleteCategory(id uint) error
}

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
        GetUserByID(id uint) (*models.User, error)
        GetUserByUsername(username string) (*models.User, error)
        CreateUser(user *models.User) error
        UpdateUser(user *models.User) error
        GetUserTrustScore(userID uint) (*models.UserTrustScore, error)
        UpdateUserTrustScore(score *models.UserTrustScore) error
}

// Using ContentLinkRepository from content_link_repository.go
// RichTextRepository is defined in rich_text_repository.go