package repository

import (
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
	"gorm.io/gorm"
)

// DiscussionRepository defines the interface for discussion data operations
type DiscussionRepository interface {
	// Discussion topic operations
	CreateDiscussionTopic(topic *models.DiscussionTopic) error
	GetDiscussionTopicByID(id uint) (*models.DiscussionTopic, error)
	GetDiscussionTopicsByContentID(contentID uint, contentType string) ([]models.DiscussionTopic, error)
	UpdateDiscussionTopic(topic *models.DiscussionTopic) error
	DeleteDiscussionTopic(id uint) error

	// Discussion post operations
	CreateDiscussionPost(post *models.DiscussionPost) error
	GetDiscussionPostByID(id uint) (*models.DiscussionPost, error)
	GetDiscussionPostsByTopicID(topicID uint) ([]models.DiscussionPost, error)
	UpdateDiscussionPost(post *models.DiscussionPost) error
	DeleteDiscussionPost(id uint) error
}

// GormDiscussionRepository implements DiscussionRepository using GORM
type GormDiscussionRepository struct {
	db *gorm.DB
}

// NewGormDiscussionRepository creates a new GormDiscussionRepository
func NewGormDiscussionRepository(db *gorm.DB) *GormDiscussionRepository {
	return &GormDiscussionRepository{db: db}
}

// CreateDiscussionTopic creates a new discussion topic in the database
func (r *GormDiscussionRepository) CreateDiscussionTopic(topic *models.DiscussionTopic) error {
	return r.db.Create(topic).Error
}

// GetDiscussionTopicByID retrieves a discussion topic by its ID
func (r *GormDiscussionRepository) GetDiscussionTopicByID(id uint) (*models.DiscussionTopic, error) {
	var topic models.DiscussionTopic
	if err := r.db.First(&topic, id).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

// GetDiscussionTopicsByContentID retrieves all discussion topics for a specific content item
func (r *GormDiscussionRepository) GetDiscussionTopicsByContentID(contentID uint, contentType string) ([]models.DiscussionTopic, error) {
	var topics []models.DiscussionTopic
	if err := r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

// UpdateDiscussionTopic updates an existing discussion topic in the database
func (r *GormDiscussionRepository) UpdateDiscussionTopic(topic *models.DiscussionTopic) error {
	return r.db.Save(topic).Error
}

// DeleteDiscussionTopic deletes a discussion topic from the database
func (r *GormDiscussionRepository) DeleteDiscussionTopic(id uint) error {
	return r.db.Delete(&models.DiscussionTopic{}, id).Error
}

// CreateDiscussionPost creates a new discussion post in the database
func (r *GormDiscussionRepository) CreateDiscussionPost(post *models.DiscussionPost) error {
	return r.db.Create(post).Error
}

// GetDiscussionPostByID retrieves a discussion post by its ID
func (r *GormDiscussionRepository) GetDiscussionPostByID(id uint) (*models.DiscussionPost, error) {
	var post models.DiscussionPost
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetDiscussionPostsByTopicID retrieves all discussion posts for a specific topic
func (r *GormDiscussionRepository) GetDiscussionPostsByTopicID(topicID uint) ([]models.DiscussionPost, error) {
	var posts []models.DiscussionPost
	if err := r.db.Where("topic_id = ?", topicID).Order("created_at").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// UpdateDiscussionPost updates an existing discussion post in the database
func (r *GormDiscussionRepository) UpdateDiscussionPost(post *models.DiscussionPost) error {
	return r.db.Save(post).Error
}

// DeleteDiscussionPost deletes a discussion post from the database
func (r *GormDiscussionRepository) DeleteDiscussionPost(id uint) error {
	return r.db.Delete(&models.DiscussionPost{}, id).Error
}