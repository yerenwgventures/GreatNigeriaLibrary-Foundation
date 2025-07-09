package repository

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/models"
        "gorm.io/gorm"
)

// ContentScoringRepository defines the interface for content scoring operations
type ContentScoringRepository interface {
        // Content scores
        CreateContentScore(score *models.ContentScore) error
        GetContentScores(contentType models.ContentType, contentID uint) ([]models.ContentScore, error)
        GetContentScoresByType(contentType models.ContentType, contentID uint, scoreType models.ContentScoreType) ([]models.ContentScore, error)
        GetLatestContentScore(contentType models.ContentType, contentID uint, scoreType models.ContentScoreType) (*models.ContentScore, error)
        UpdateContentScore(score *models.ContentScore) error
        DeleteContentScore(id uint) error
        
        // Score criteria
        CreateScoreCriteria(criteria *models.ContentScoreCriteria) error
        GetScoreCriteria(scoreType models.ContentScoreType) ([]models.ContentScoreCriteria, error)
        GetScoreCriteriaByID(id uint) (*models.ContentScoreCriteria, error)
        UpdateScoreCriteria(criteria *models.ContentScoreCriteria) error
        DeleteScoreCriteria(id uint) error
        
        // Score values
        CreateScoreValue(value *models.ContentScoreValue) error
        GetScoreValues(scoreID uint) ([]models.ContentScoreValue, error)
        UpdateScoreValue(value *models.ContentScoreValue) error
        DeleteScoreValue(id uint) error
        
        // Quality metrics
        CreateQualityMetrics(metrics *models.ContentQualityMetrics) error
        GetQualityMetrics(contentType models.ContentType, contentID uint) (*models.ContentQualityMetrics, error)
        UpdateQualityMetrics(metrics *models.ContentQualityMetrics) error
        
        // Relevance metrics
        CreateRelevanceMetrics(metrics *models.ContentRelevanceMetrics) error
        GetRelevanceMetrics(contentType models.ContentType, contentID uint) (*models.ContentRelevanceMetrics, error)
        UpdateRelevanceMetrics(metrics *models.ContentRelevanceMetrics) error
        
        // Safety metrics
        CreateSafetyMetrics(metrics *models.ContentSafetyMetrics) error
        GetSafetyMetrics(contentType models.ContentType, contentID uint) (*models.ContentSafetyMetrics, error)
        UpdateSafetyMetrics(metrics *models.ContentSafetyMetrics) error
        GetFlaggedContent() ([]models.ContentSafetyMetrics, error)
        
        // Automated analysis
        CreateAutomatedAnalysis(analysis *models.ContentAutomatedAnalysis) error
        GetAutomatedAnalysis(contentType models.ContentType, contentID uint) (*models.ContentAutomatedAnalysis, error)
        UpdateAutomatedAnalysis(analysis *models.ContentAutomatedAnalysis) error
}

// GormContentScoringRepository implements the ContentScoringRepository interface using GORM
type GormContentScoringRepository struct {
        db *gorm.DB
}

// NewGormContentScoringRepository creates a new repository instance
func NewGormContentScoringRepository(db *gorm.DB) *GormContentScoringRepository {
        return &GormContentScoringRepository{db: db}
}

// CreateContentScore creates a new content score
func (r *GormContentScoringRepository) CreateContentScore(score *models.ContentScore) error {
        return r.db.Create(score).Error
}

// GetContentScores retrieves all content scores for a piece of content
func (r *GormContentScoringRepository) GetContentScores(contentType models.ContentType, contentID uint) ([]models.ContentScore, error) {
        var scores []models.ContentScore
        result := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).Find(&scores)
        return scores, result.Error
}

// GetContentScoresByType retrieves all content scores of a specific type for a piece of content
func (r *GormContentScoringRepository) GetContentScoresByType(contentType models.ContentType, contentID uint, scoreType models.ContentScoreType) ([]models.ContentScore, error) {
        var scores []models.ContentScore
        result := r.db.Where("content_type = ? AND content_id = ? AND score_type = ?", contentType, contentID, scoreType).Find(&scores)
        return scores, result.Error
}

// GetLatestContentScore retrieves the latest content score of a specific type for a piece of content
func (r *GormContentScoringRepository) GetLatestContentScore(contentType models.ContentType, contentID uint, scoreType models.ContentScoreType) (*models.ContentScore, error) {
        var score models.ContentScore
        result := r.db.Where("content_type = ? AND content_id = ? AND score_type = ?", contentType, contentID, scoreType).
                Order("scored_at DESC").
                First(&score)
        if result.Error != nil {
                return nil, result.Error
        }
        return &score, nil
}

// UpdateContentScore updates a content score
func (r *GormContentScoringRepository) UpdateContentScore(score *models.ContentScore) error {
        return r.db.Save(score).Error
}

// DeleteContentScore deletes a content score
func (r *GormContentScoringRepository) DeleteContentScore(id uint) error {
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Delete all score values
                if err := tx.Where("content_score_id = ?", id).Delete(&models.ContentScoreValue{}).Error; err != nil {
                        return err
                }
                
                // Delete the score
                return tx.Delete(&models.ContentScore{}, id).Error
        })
}

// CreateScoreCriteria creates a new score criteria
func (r *GormContentScoringRepository) CreateScoreCriteria(criteria *models.ContentScoreCriteria) error {
        return r.db.Create(criteria).Error
}

// GetScoreCriteria retrieves all score criteria for a score type
func (r *GormContentScoringRepository) GetScoreCriteria(scoreType models.ContentScoreType) ([]models.ContentScoreCriteria, error) {
        var criteria []models.ContentScoreCriteria
        result := r.db.Where("score_type = ?", scoreType).Find(&criteria)
        return criteria, result.Error
}

// GetScoreCriteriaByID retrieves a score criteria by ID
func (r *GormContentScoringRepository) GetScoreCriteriaByID(id uint) (*models.ContentScoreCriteria, error) {
        var criteria models.ContentScoreCriteria
        result := r.db.First(&criteria, id)
        if result.Error != nil {
                return nil, result.Error
        }
        return &criteria, nil
}

// UpdateScoreCriteria updates a score criteria
func (r *GormContentScoringRepository) UpdateScoreCriteria(criteria *models.ContentScoreCriteria) error {
        return r.db.Save(criteria).Error
}

// DeleteScoreCriteria deletes a score criteria
func (r *GormContentScoringRepository) DeleteScoreCriteria(id uint) error {
        return r.db.Delete(&models.ContentScoreCriteria{}, id).Error
}

// CreateScoreValue creates a new score value
func (r *GormContentScoringRepository) CreateScoreValue(value *models.ContentScoreValue) error {
        return r.db.Create(value).Error
}

// GetScoreValues retrieves all score values for a content score
func (r *GormContentScoringRepository) GetScoreValues(scoreID uint) ([]models.ContentScoreValue, error) {
        var values []models.ContentScoreValue
        result := r.db.Where("content_score_id = ?", scoreID).Find(&values)
        return values, result.Error
}

// UpdateScoreValue updates a score value
func (r *GormContentScoringRepository) UpdateScoreValue(value *models.ContentScoreValue) error {
        return r.db.Save(value).Error
}

// DeleteScoreValue deletes a score value
func (r *GormContentScoringRepository) DeleteScoreValue(id uint) error {
        return r.db.Delete(&models.ContentScoreValue{}, id).Error
}

// CreateQualityMetrics creates new quality metrics
func (r *GormContentScoringRepository) CreateQualityMetrics(metrics *models.ContentQualityMetrics) error {
        return r.db.Create(metrics).Error
}

// GetQualityMetrics retrieves quality metrics for a piece of content
func (r *GormContentScoringRepository) GetQualityMetrics(contentType models.ContentType, contentID uint) (*models.ContentQualityMetrics, error) {
        var metrics models.ContentQualityMetrics
        result := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).First(&metrics)
        if result.Error != nil {
                return nil, result.Error
        }
        return &metrics, nil
}

// UpdateQualityMetrics updates quality metrics
func (r *GormContentScoringRepository) UpdateQualityMetrics(metrics *models.ContentQualityMetrics) error {
        return r.db.Save(metrics).Error
}

// CreateRelevanceMetrics creates new relevance metrics
func (r *GormContentScoringRepository) CreateRelevanceMetrics(metrics *models.ContentRelevanceMetrics) error {
        return r.db.Create(metrics).Error
}

// GetRelevanceMetrics retrieves relevance metrics for a piece of content
func (r *GormContentScoringRepository) GetRelevanceMetrics(contentType models.ContentType, contentID uint) (*models.ContentRelevanceMetrics, error) {
        var metrics models.ContentRelevanceMetrics
        result := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).First(&metrics)
        if result.Error != nil {
                return nil, result.Error
        }
        return &metrics, nil
}

// UpdateRelevanceMetrics updates relevance metrics
func (r *GormContentScoringRepository) UpdateRelevanceMetrics(metrics *models.ContentRelevanceMetrics) error {
        return r.db.Save(metrics).Error
}

// CreateSafetyMetrics creates new safety metrics
func (r *GormContentScoringRepository) CreateSafetyMetrics(metrics *models.ContentSafetyMetrics) error {
        return r.db.Create(metrics).Error
}

// GetSafetyMetrics retrieves safety metrics for a piece of content
func (r *GormContentScoringRepository) GetSafetyMetrics(contentType models.ContentType, contentID uint) (*models.ContentSafetyMetrics, error) {
        var metrics models.ContentSafetyMetrics
        result := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).First(&metrics)
        if result.Error != nil {
                return nil, result.Error
        }
        return &metrics, nil
}

// UpdateSafetyMetrics updates safety metrics
func (r *GormContentScoringRepository) UpdateSafetyMetrics(metrics *models.ContentSafetyMetrics) error {
        return r.db.Save(metrics).Error
}

// GetFlaggedContent retrieves all flagged content
func (r *GormContentScoringRepository) GetFlaggedContent() ([]models.ContentSafetyMetrics, error) {
        var metrics []models.ContentSafetyMetrics
        result := r.db.Where("flagged_content = ?", true).Find(&metrics)
        return metrics, result.Error
}

// CreateAutomatedAnalysis creates new automated analysis
func (r *GormContentScoringRepository) CreateAutomatedAnalysis(analysis *models.ContentAutomatedAnalysis) error {
        return r.db.Create(analysis).Error
}

// GetAutomatedAnalysis retrieves automated analysis for a piece of content
func (r *GormContentScoringRepository) GetAutomatedAnalysis(contentType models.ContentType, contentID uint) (*models.ContentAutomatedAnalysis, error) {
        var analysis models.ContentAutomatedAnalysis
        result := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).First(&analysis)
        if result.Error != nil {
                return nil, result.Error
        }
        return &analysis, nil
}

// UpdateAutomatedAnalysis updates automated analysis
func (r *GormContentScoringRepository) UpdateAutomatedAnalysis(analysis *models.ContentAutomatedAnalysis) error {
        return r.db.Save(analysis).Error
}