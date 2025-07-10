package repository

import (
	"errors"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/models"
	"gorm.io/gorm"
)

// GormTopicRepository implements TopicRepository using GORM
type GormTopicRepository struct {
	db *gorm.DB
}

// NewGormTopicRepository creates a new GormTopicRepository
func NewGormTopicRepository(db *gorm.DB) *GormTopicRepository {
	return &GormTopicRepository{db: db}
}

// GetTopicByID retrieves a topic by ID
func (r *GormTopicRepository) GetTopicByID(id uint) (*models.Topic, error) {
	var topic models.Topic
	result := r.db.First(&topic, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &topic, nil
}

// CreateTopic creates a new topic
func (r *GormTopicRepository) CreateTopic(topic *models.Topic) error {
	return r.db.Create(topic).Error
}

// UpdateTopic updates a topic
func (r *GormTopicRepository) UpdateTopic(topic *models.Topic) error {
	return r.db.Save(topic).Error
}

// DeleteTopic deletes a topic
func (r *GormTopicRepository) DeleteTopic(id uint) error {
	return r.db.Delete(&models.Topic{}, id).Error
}

// GetTopics retrieves topics with pagination and filters
func (r *GormTopicRepository) GetTopics(page, pageSize int, filters map[string]interface{}) ([]models.Topic, int64, error) {
	var topics []models.Topic
	var total int64
	
	query := r.db
	
	// Apply filters
	for key, value := range filters {
		query = query.Where(key, value)
	}
	
	// Count total
	if err := query.Model(&models.Topic{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply pagination
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&topics)
	
	return topics, total, result.Error
}

// GetTopicsByCategory retrieves topics by category with pagination
func (r *GormTopicRepository) GetTopicsByCategory(categoryID uint, page, pageSize int) ([]models.Topic, int64, error) {
	return r.GetTopics(page, pageSize, map[string]interface{}{"category_id": categoryID})
}

// GetTopicsByUser retrieves topics by user with pagination
func (r *GormTopicRepository) GetTopicsByUser(userID uint, page, pageSize int) ([]models.Topic, int64, error) {
	return r.GetTopics(page, pageSize, map[string]interface{}{"user_id": userID})
}

// GetTopicsByBookSection retrieves topics by book section
func (r *GormTopicRepository) GetTopicsByBookSection(bookID, chapterID, sectionID uint) ([]models.Topic, error) {
	var topics []models.Topic
	
	// Use the link table to find topics related to a section
	query := r.db.Table("topics").
		Joins("JOIN topic_content_links ON topics.id = topic_content_links.topic_id")
	
	if sectionID > 0 {
		query = query.Where("topic_content_links.content_type = ? AND topic_content_links.content_id = ?", 
			models.SectionReference, sectionID)
	} else if chapterID > 0 {
		query = query.Where("topic_content_links.content_type = ? AND topic_content_links.content_id = ?", 
			models.ChapterReference, chapterID)
	} else if bookID > 0 {
		query = query.Where("topic_content_links.content_type = ? AND topic_content_links.content_id = ?", 
			models.BookReference, bookID)
	} else {
		return nil, errors.New("at least one of bookID, chapterID, or sectionID must be specified")
	}
	
	result := query.Order("topics.created_at DESC").Find(&topics)
	return topics, result.Error
}

// GetAllTopics retrieves all topics with pagination
func (r *GormTopicRepository) GetAllTopics(page, pageSize int) ([]models.Topic, error) {
	var topics []models.Topic
	
	// Apply pagination
	offset := (page - 1) * pageSize
	result := r.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&topics)
	
	return topics, result.Error
}

// IncrementTopicViewCount increments a topic's view count
func (r *GormTopicRepository) IncrementTopicViewCount(id uint) error {
	return r.db.Model(&models.Topic{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// UpdateTopicLastPostTime updates a topic's last post time
func (r *GormTopicRepository) UpdateTopicLastPostTime(id uint) error {
	return r.db.Model(&models.Topic{}).Where("id = ?", id).
		UpdateColumns(map[string]interface{}{
			"last_activity_at": gorm.Expr("NOW()"),
			"reply_count": gorm.Expr("reply_count + ?", 1),
		}).Error
}

// PinTopic pins or unpins a topic
func (r *GormTopicRepository) PinTopic(id uint, pinned bool) error {
	return r.db.Model(&models.Topic{}).Where("id = ?", id).
		Update("is_pinned", pinned).Error
}

// LockTopic locks or unlocks a topic
func (r *GormTopicRepository) LockTopic(id uint, locked bool) error {
	return r.db.Model(&models.Topic{}).Where("id = ?", id).
		Update("is_locked", locked).Error
}