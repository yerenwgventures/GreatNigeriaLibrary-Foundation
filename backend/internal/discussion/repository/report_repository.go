package repository

import (
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/models"
	"gorm.io/gorm"
)

// ReportRepository defines the interface for content report operations
type ReportRepository interface {
	// Report operations
	CreateReport(report *models.ContentReport) error
	GetReportByID(id uint) (*models.ContentReport, error)
	GetReportsByStatus(status models.ReportStatus, page, pageSize int) ([]models.ContentReport, int64, error)
	GetReportsByReporter(userID uint, page, pageSize int) ([]models.ContentReport, int64, error)
	GetReportsByContent(contentType string, contentID uint) ([]models.ContentReport, error)
	GetReportsByCategory(category models.ReportCategory, page, pageSize int) ([]models.ContentReport, int64, error)
	AssignReport(reportID uint, moderatorID uint) error
	UpdateReportStatus(reportID uint, status models.ReportStatus, reviewedBy *uint) error
	ResolveReport(reportID uint, resolution models.ReportResolutionType, notes string, reviewedBy uint) error
	MarkReporterNotified(reportID uint) error
	GetReportCount(status *models.ReportStatus, category *models.ReportCategory) (int64, error)
	
	// Evidence operations
	AddReportEvidence(evidence *models.ReportEvidence) error
	GetReportEvidence(reportID uint) ([]models.ReportEvidence, error)
	DeleteEvidence(evidenceID uint) error
	
	// Comment operations
	AddReportComment(comment *models.ReportComment) error
	GetReportComments(reportID uint, internalOnly bool) ([]models.ReportComment, error)
	DeleteReportComment(commentID uint) error
	
	// Action log operations
	AddActionLog(log *models.ReportActionLog) error
	GetActionLogs(reportID uint) ([]models.ReportActionLog, error)
}

// GormReportRepository implements the ReportRepository interface
type GormReportRepository struct {
	db *gorm.DB
}

// NewGormReportRepository creates a new report repository
func NewGormReportRepository(db *gorm.DB) *GormReportRepository {
	return &GormReportRepository{db: db}
}

// CreateReport creates a new content report
func (r *GormReportRepository) CreateReport(report *models.ContentReport) error {
	return r.db.Create(report).Error
}

// GetReportByID retrieves a report by ID
func (r *GormReportRepository) GetReportByID(id uint) (*models.ContentReport, error) {
	var report models.ContentReport
	result := r.db.First(&report, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &report, nil
}

// GetReportsByStatus retrieves reports by status with pagination
func (r *GormReportRepository) GetReportsByStatus(status models.ReportStatus, page, pageSize int) ([]models.ContentReport, int64, error) {
	var reports []models.ContentReport
	var total int64
	
	offset := (page - 1) * pageSize
	
	// Get count
	countResult := r.db.Model(&models.ContentReport{}).Where("status = ?", status).Count(&total)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}
	
	// Get reports
	result := r.db.Where("status = ?", status).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&reports)
	
	return reports, total, result.Error
}

// GetReportsByReporter retrieves reports by reporter ID with pagination
func (r *GormReportRepository) GetReportsByReporter(userID uint, page, pageSize int) ([]models.ContentReport, int64, error) {
	var reports []models.ContentReport
	var total int64
	
	offset := (page - 1) * pageSize
	
	// Get count
	countResult := r.db.Model(&models.ContentReport{}).Where("reporter_id = ?", userID).Count(&total)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}
	
	// Get reports
	result := r.db.Where("reporter_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&reports)
	
	return reports, total, result.Error
}

// GetReportsByContent retrieves reports by content type and ID
func (r *GormReportRepository) GetReportsByContent(contentType string, contentID uint) ([]models.ContentReport, error) {
	var reports []models.ContentReport
	result := r.db.Where("content_type = ? AND content_id = ?", contentType, contentID).
		Order("created_at DESC").
		Find(&reports)
	return reports, result.Error
}

// GetReportsByCategory retrieves reports by category with pagination
func (r *GormReportRepository) GetReportsByCategory(category models.ReportCategory, page, pageSize int) ([]models.ContentReport, int64, error) {
	var reports []models.ContentReport
	var total int64
	
	offset := (page - 1) * pageSize
	
	// Get count
	countResult := r.db.Model(&models.ContentReport{}).Where("category = ?", category).Count(&total)
	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}
	
	// Get reports
	result := r.db.Where("category = ?", category).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&reports)
	
	return reports, total, result.Error
}

// AssignReport assigns a report to a moderator
func (r *GormReportRepository) AssignReport(reportID uint, moderatorID uint) error {
	return r.db.Model(&models.ContentReport{}).
		Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"assigned_to": moderatorID,
			"updated_at":  time.Now(),
		}).Error
}

// UpdateReportStatus updates the status of a report
func (r *GormReportRepository) UpdateReportStatus(reportID uint, status models.ReportStatus, reviewedBy *uint) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	
	if status == models.StatusInReview && reviewedBy != nil {
		updates["assigned_to"] = *reviewedBy
	}
	
	return r.db.Model(&models.ContentReport{}).
		Where("id = ?", reportID).
		Updates(updates).Error
}

// ResolveReport resolves a report
func (r *GormReportRepository) ResolveReport(reportID uint, resolution models.ReportResolutionType, notes string, reviewedBy uint) error {
	now := time.Now()
	
	return r.db.Model(&models.ContentReport{}).
		Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"status":           models.StatusResolved,
			"resolution":       resolution,
			"resolution_notes": notes,
			"reviewed_by":      reviewedBy,
			"reviewed_at":      now,
			"updated_at":       now,
		}).Error
}

// MarkReporterNotified marks a report as having the reporter notified
func (r *GormReportRepository) MarkReporterNotified(reportID uint) error {
	return r.db.Model(&models.ContentReport{}).
		Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"reporter_notified": true,
			"updated_at":        time.Now(),
		}).Error
}

// GetReportCount gets the total count of reports, optionally filtered by status and/or category
func (r *GormReportRepository) GetReportCount(status *models.ReportStatus, category *models.ReportCategory) (int64, error) {
	var count int64
	query := r.db.Model(&models.ContentReport{})
	
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	
	if category != nil {
		query = query.Where("category = ?", *category)
	}
	
	result := query.Count(&count)
	return count, result.Error
}

// AddReportEvidence adds evidence to a report
func (r *GormReportRepository) AddReportEvidence(evidence *models.ReportEvidence) error {
	return r.db.Create(evidence).Error
}

// GetReportEvidence retrieves evidence for a report
func (r *GormReportRepository) GetReportEvidence(reportID uint) ([]models.ReportEvidence, error) {
	var evidence []models.ReportEvidence
	result := r.db.Where("report_id = ?", reportID).Find(&evidence)
	return evidence, result.Error
}

// DeleteEvidence deletes evidence
func (r *GormReportRepository) DeleteEvidence(evidenceID uint) error {
	return r.db.Delete(&models.ReportEvidence{}, evidenceID).Error
}

// AddReportComment adds a comment to a report
func (r *GormReportRepository) AddReportComment(comment *models.ReportComment) error {
	return r.db.Create(comment).Error
}

// GetReportComments retrieves comments for a report
func (r *GormReportRepository) GetReportComments(reportID uint, internalOnly bool) ([]models.ReportComment, error) {
	var comments []models.ReportComment
	query := r.db.Where("report_id = ?", reportID)
	
	if internalOnly {
		query = query.Where("is_internal = ?", true)
	}
	
	result := query.Order("created_at ASC").Find(&comments)
	return comments, result.Error
}

// DeleteReportComment deletes a comment
func (r *GormReportRepository) DeleteReportComment(commentID uint) error {
	return r.db.Delete(&models.ReportComment{}, commentID).Error
}

// AddActionLog adds an action log entry
func (r *GormReportRepository) AddActionLog(log *models.ReportActionLog) error {
	return r.db.Create(log).Error
}

// GetActionLogs retrieves action logs for a report
func (r *GormReportRepository) GetActionLogs(reportID uint) ([]models.ReportActionLog, error) {
	var logs []models.ReportActionLog
	result := r.db.Where("report_id = ?", reportID).Order("created_at ASC").Find(&logs)
	return logs, result.Error
}