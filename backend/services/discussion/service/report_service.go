package service

import (
        "errors"
        "fmt"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/repository"
)

// ReportService defines the interface for content report operations
type ReportService interface {
        // Report operations
        CreateReport(reporterID uint, contentType string, contentID uint, category models.ReportCategory, reason, additionalInfo string) (*models.ContentReport, error)
        GetReportByID(id uint) (*models.ContentReport, error)
        GetReportsByStatus(status models.ReportStatus, page, pageSize int) ([]models.ContentReport, int64, error)
        GetReportsByReporter(userID uint, page, pageSize int) ([]models.ContentReport, int64, error)
        GetReportsByCategoryAndStatus(category models.ReportCategory, status *models.ReportStatus, page, pageSize int) ([]models.ContentReport, int64, error)
        
        // Report moderation
        AssignReportToModerator(reportID uint, moderatorID uint, actionUserID uint) error
        UpdateReportStatus(reportID uint, status models.ReportStatus, actionUserID uint) error
        ResolveReport(reportID uint, resolution models.ReportResolutionType, notes string, actionUserID uint) error
        
        // Notification tracking
        MarkReporterNotified(reportID uint) error
        
        // Dashboard stats
        GetReportStats() (map[string]int64, error)
        
        // Evidence operations
        AddReportEvidence(reportID uint, userID uint, evidenceType, content, filePath, url string) (*models.ReportEvidence, error)
        GetReportEvidence(reportID uint) ([]models.ReportEvidence, error)
        DeleteEvidence(evidenceID uint, userID uint) error
        
        // Comment operations
        AddReportComment(reportID uint, userID uint, comment string, isInternal bool) (*models.ReportComment, error)
        GetReportComments(reportID uint, includeInternal bool, userID uint) ([]models.ReportComment, error)
        DeleteReportComment(commentID uint, userID uint) error
        
        // Action log operations
        GetActionLogs(reportID uint) ([]models.ReportActionLog, error)
}

// ReportServiceImpl implements the ReportService interface
type ReportServiceImpl struct {
        reportRepo     repository.ReportRepository
        userRepo       repository.UserRepository
        discussionRepo repository.DiscussionRepository
        commentRepo    repository.CommentRepository
}

// NewReportService creates a new report service
func NewReportService(
        reportRepo repository.ReportRepository,
        userRepo repository.UserRepository,
        discussionRepo repository.DiscussionRepository,
        commentRepo repository.CommentRepository,
) ReportService {
        return &ReportServiceImpl{
                reportRepo:     reportRepo,
                userRepo:       userRepo,
                discussionRepo: discussionRepo,
                commentRepo:    commentRepo,
        }
}

// CreateReport creates a new content report
func (s *ReportServiceImpl) CreateReport(
        reporterID uint,
        contentType string,
        contentID uint,
        category models.ReportCategory,
        reason,
        additionalInfo string,
) (*models.ContentReport, error) {
        // Validate content type
        if contentType != "topic" && contentType != "comment" {
                return nil, errors.New("invalid content type")
        }
        
        // Validate that the content exists
        if err := s.validateContent(contentType, contentID); err != nil {
                return nil, err
        }
        
        // Check if the user has already reported this content
        existingReports, err := s.reportRepo.GetReportsByContent(contentType, contentID)
        if err != nil {
                return nil, fmt.Errorf("error checking existing reports: %w", err)
        }
        
        for _, report := range existingReports {
                if report.ReporterID == reporterID && report.Status != models.StatusResolved && report.Status != models.StatusRejected {
                        return nil, errors.New("you have already reported this content")
                }
        }
        
        // Create the report
        report := &models.ContentReport{
                ReporterID:     reporterID,
                ContentType:    contentType,
                ContentID:      contentID,
                Category:       category,
                Reason:         reason,
                AdditionalInfo: additionalInfo,
                Status:         models.StatusPending,
                CreatedAt:      time.Now(),
                UpdatedAt:      time.Now(),
        }
        
        if err := s.reportRepo.CreateReport(report); err != nil {
                return nil, fmt.Errorf("error creating report: %w", err)
        }
        
        // Log the action
        actionLog := &models.ReportActionLog{
                ReportID:  report.ID,
                UserID:    reporterID,
                Action:    "created",
                NewValue:  string(models.StatusPending),
                CreatedAt: time.Now(),
        }
        
        if err := s.reportRepo.AddActionLog(actionLog); err != nil {
                // Just log the error but don't fail the operation
                fmt.Printf("Error adding action log: %v\n", err)
        }
        
        return report, nil
}

// GetReportByID retrieves a report by ID
func (s *ReportServiceImpl) GetReportByID(id uint) (*models.ContentReport, error) {
        return s.reportRepo.GetReportByID(id)
}

// GetReportsByStatus retrieves reports by status with pagination
func (s *ReportServiceImpl) GetReportsByStatus(status models.ReportStatus, page, pageSize int) ([]models.ContentReport, int64, error) {
        // Ensure valid pagination
        if page < 1 {
                page = 1
        }
        if pageSize < 1 || pageSize > 100 {
                pageSize = 20
        }
        
        return s.reportRepo.GetReportsByStatus(status, page, pageSize)
}

// GetReportsByReporter retrieves reports by reporter ID with pagination
func (s *ReportServiceImpl) GetReportsByReporter(userID uint, page, pageSize int) ([]models.ContentReport, int64, error) {
        // Ensure valid pagination
        if page < 1 {
                page = 1
        }
        if pageSize < 1 || pageSize > 100 {
                pageSize = 20
        }
        
        return s.reportRepo.GetReportsByReporter(userID, page, pageSize)
}

// GetReportsByCategoryAndStatus retrieves reports by category and optionally status with pagination
func (s *ReportServiceImpl) GetReportsByCategoryAndStatus(category models.ReportCategory, status *models.ReportStatus, page, pageSize int) ([]models.ContentReport, int64, error) {
        // Ensure valid pagination
        if page < 1 {
                page = 1
        }
        if pageSize < 1 || pageSize > 100 {
                pageSize = 20
        }
        
        reports, count, err := s.reportRepo.GetReportsByCategory(category, page, pageSize)
        if err != nil {
                return nil, 0, err
        }
        
        // If status is provided, filter in memory
        // In a real application, we'd modify the repository to filter by both category and status in the DB query
        if status != nil {
                filteredReports := make([]models.ContentReport, 0)
                for _, report := range reports {
                        if report.Status == *status {
                                filteredReports = append(filteredReports, report)
                        }
                }
                
                // Get new count
                statusCat := *status
                newCount, err := s.reportRepo.GetReportCount(&statusCat, &category)
                if err != nil {
                        // Fall back to length if there's an error
                        newCount = int64(len(filteredReports))
                }
                
                return filteredReports, newCount, nil
        }
        
        return reports, count, nil
}

// AssignReportToModerator assigns a report to a moderator
func (s *ReportServiceImpl) AssignReportToModerator(reportID uint, moderatorID uint, actionUserID uint) error {
        // Get the report to check current state
        report, err := s.reportRepo.GetReportByID(reportID)
        if err != nil {
                return fmt.Errorf("error getting report: %w", err)
        }
        
        // Capture old value for action log
        var oldAssignedTo string
        if report.AssignedTo != nil {
                oldAssignedTo = fmt.Sprintf("%d", *report.AssignedTo)
        } else {
                oldAssignedTo = "none"
        }
        
        // Assign the report
        if err := s.reportRepo.AssignReport(reportID, moderatorID); err != nil {
                return fmt.Errorf("error assigning report: %w", err)
        }
        
        // Update status to in review if it was pending
        if report.Status == models.StatusPending {
                if err := s.reportRepo.UpdateReportStatus(reportID, models.StatusInReview, &moderatorID); err != nil {
                        // Log error but continue
                        fmt.Printf("Error updating report status: %v\n", err)
                }
        }
        
        // Log the action
        actionLog := &models.ReportActionLog{
                ReportID:  reportID,
                UserID:    actionUserID,
                Action:    "assigned",
                OldValue:  oldAssignedTo,
                NewValue:  fmt.Sprintf("%d", moderatorID),
                CreatedAt: time.Now(),
        }
        
        if err := s.reportRepo.AddActionLog(actionLog); err != nil {
                // Just log the error but don't fail the operation
                fmt.Printf("Error adding action log: %v\n", err)
        }
        
        return nil
}

// UpdateReportStatus updates the status of a report
func (s *ReportServiceImpl) UpdateReportStatus(reportID uint, status models.ReportStatus, actionUserID uint) error {
        // Get the report to check current state
        report, err := s.reportRepo.GetReportByID(reportID)
        if err != nil {
                return fmt.Errorf("error getting report: %w", err)
        }
        
        // Validate status transition
        if !s.isValidStatusTransition(report.Status, status) {
                return fmt.Errorf("invalid status transition from %s to %s", report.Status, status)
        }
        
        // Update the status
        if err := s.reportRepo.UpdateReportStatus(reportID, status, &actionUserID); err != nil {
                return fmt.Errorf("error updating report status: %w", err)
        }
        
        // Log the action
        actionLog := &models.ReportActionLog{
                ReportID:  reportID,
                UserID:    actionUserID,
                Action:    "status_changed",
                OldValue:  string(report.Status),
                NewValue:  string(status),
                CreatedAt: time.Now(),
        }
        
        if err := s.reportRepo.AddActionLog(actionLog); err != nil {
                // Just log the error but don't fail the operation
                fmt.Printf("Error adding action log: %v\n", err)
        }
        
        return nil
}

// ResolveReport resolves a report
func (s *ReportServiceImpl) ResolveReport(reportID uint, resolution models.ReportResolutionType, notes string, actionUserID uint) error {
        // Get the report to check current state
        report, err := s.reportRepo.GetReportByID(reportID)
        if err != nil {
                return fmt.Errorf("error getting report: %w", err)
        }
        
        // Validate status transition
        if report.Status != models.StatusInReview && report.Status != models.StatusPending {
                return fmt.Errorf("cannot resolve report that is not in review or pending")
        }
        
        // Resolve the report
        if err := s.reportRepo.ResolveReport(reportID, resolution, notes, actionUserID); err != nil {
                return fmt.Errorf("error resolving report: %w", err)
        }
        
        // Take appropriate action based on resolution
        if err := s.executeResolution(report, resolution, actionUserID); err != nil {
                // Log error but continue
                fmt.Printf("Error executing resolution: %v\n", err)
        }
        
        // Log the action
        actionLog := &models.ReportActionLog{
                ReportID:  reportID,
                UserID:    actionUserID,
                Action:    "resolved",
                OldValue:  string(report.Status),
                NewValue:  string(resolution),
                CreatedAt: time.Now(),
        }
        
        if err := s.reportRepo.AddActionLog(actionLog); err != nil {
                // Just log the error but don't fail the operation
                fmt.Printf("Error adding action log: %v\n", err)
        }
        
        return nil
}

// MarkReporterNotified marks a report as having the reporter notified
func (s *ReportServiceImpl) MarkReporterNotified(reportID uint) error {
        return s.reportRepo.MarkReporterNotified(reportID)
}

// GetReportStats gets statistics about reports
func (s *ReportServiceImpl) GetReportStats() (map[string]int64, error) {
        stats := make(map[string]int64)
        
        // Get total count
        total, err := s.reportRepo.GetReportCount(nil, nil)
        if err != nil {
                return nil, fmt.Errorf("error getting total report count: %w", err)
        }
        stats["total"] = total
        
        // Get count by status
        statuses := []models.ReportStatus{
                models.StatusPending,
                models.StatusInReview,
                models.StatusResolved,
                models.StatusRejected,
        }
        
        for _, status := range statuses {
                statusCopy := status // Create a local copy for the pointer
                count, err := s.reportRepo.GetReportCount(&statusCopy, nil)
                if err != nil {
                        // Log error but continue
                        fmt.Printf("Error getting report count for status %s: %v\n", status, err)
                        continue
                }
                stats[string(status)] = count
        }
        
        // Get count by category
        categories := []models.ReportCategory{
                models.CategorySpam,
                models.CategoryHarassment,
                models.CategoryHateSpeech,
                models.CategoryViolence,
                models.CategoryIllegalContent,
                models.CategoryPrivacyViolation,
                models.CategoryCopyright,
                models.CategoryMisinformation,
                models.CategoryOther,
        }
        
        for _, category := range categories {
                categoryCopy := category // Create a local copy for the pointer
                count, err := s.reportRepo.GetReportCount(nil, &categoryCopy)
                if err != nil {
                        // Log error but continue
                        fmt.Printf("Error getting report count for category %s: %v\n", category, err)
                        continue
                }
                stats[string(category)] = count
        }
        
        return stats, nil
}

// AddReportEvidence adds evidence to a report
func (s *ReportServiceImpl) AddReportEvidence(reportID uint, userID uint, evidenceType, content, filePath, url string) (*models.ReportEvidence, error) {
        // Validate that the report exists
        report, err := s.reportRepo.GetReportByID(reportID)
        if err != nil {
                return nil, fmt.Errorf("error getting report: %w", err)
        }
        
        // Check if the user is the reporter or a moderator
        if report.ReporterID != userID && (report.AssignedTo == nil || *report.AssignedTo != userID) {
                return nil, errors.New("unauthorized to add evidence to this report")
        }
        
        // Create the evidence
        evidence := &models.ReportEvidence{
                ReportID:  reportID,
                Type:      evidenceType,
                Content:   content,
                FilePath:  filePath,
                URL:       url,
                CreatedAt: time.Now(),
        }
        
        if err := s.reportRepo.AddReportEvidence(evidence); err != nil {
                return nil, fmt.Errorf("error adding report evidence: %w", err)
        }
        
        // Log the action
        actionLog := &models.ReportActionLog{
                ReportID:  reportID,
                UserID:    userID,
                Action:    "evidence_added",
                NewValue:  fmt.Sprintf("%s evidence", evidenceType),
                CreatedAt: time.Now(),
        }
        
        if err := s.reportRepo.AddActionLog(actionLog); err != nil {
                // Just log the error but don't fail the operation
                fmt.Printf("Error adding action log: %v\n", err)
        }
        
        return evidence, nil
}

// GetReportEvidence retrieves evidence for a report
func (s *ReportServiceImpl) GetReportEvidence(reportID uint) ([]models.ReportEvidence, error) {
        return s.reportRepo.GetReportEvidence(reportID)
}

// DeleteEvidence deletes evidence
func (s *ReportServiceImpl) DeleteEvidence(evidenceID uint, userID uint) error {
        // Implementation would verify permissions first
        return s.reportRepo.DeleteEvidence(evidenceID)
}

// AddReportComment adds a comment to a report
func (s *ReportServiceImpl) AddReportComment(reportID uint, userID uint, comment string, isInternal bool) (*models.ReportComment, error) {
        // Validate that the report exists
        report, err := s.reportRepo.GetReportByID(reportID)
        if err != nil {
                return nil, fmt.Errorf("error getting report: %w", err)
        }
        
        // Check if the user is the reporter or a moderator
        if isInternal && (report.AssignedTo == nil || *report.AssignedTo != userID) {
                return nil, errors.New("unauthorized to add internal comments to this report")
        }
        
        // External comments should only be from the reporter or a moderator
        if !isInternal && report.ReporterID != userID && (report.AssignedTo == nil || *report.AssignedTo != userID) {
                return nil, errors.New("unauthorized to add comments to this report")
        }
        
        // Create the comment
        reportComment := &models.ReportComment{
                ReportID:   reportID,
                UserID:     userID,
                Comment:    comment,
                IsInternal: isInternal,
                CreatedAt:  time.Now(),
        }
        
        if err := s.reportRepo.AddReportComment(reportComment); err != nil {
                return nil, fmt.Errorf("error adding report comment: %w", err)
        }
        
        // Log the action
        actionLog := &models.ReportActionLog{
                ReportID:  reportID,
                UserID:    userID,
                Action:    "comment_added",
                NewValue:  comment,
                CreatedAt: time.Now(),
        }
        
        if err := s.reportRepo.AddActionLog(actionLog); err != nil {
                // Just log the error but don't fail the operation
                fmt.Printf("Error adding action log: %v\n", err)
        }
        
        return reportComment, nil
}

// GetReportComments retrieves comments for a report
func (s *ReportServiceImpl) GetReportComments(reportID uint, includeInternal bool, userID uint) ([]models.ReportComment, error) {
        // Check if the user can view internal comments
        if includeInternal {
                // Validate that the report exists
                report, err := s.reportRepo.GetReportByID(reportID)
                if err != nil {
                        return nil, fmt.Errorf("error getting report: %w", err)
                }
                
                // Check if the user is a moderator assigned to this report
                if report.AssignedTo == nil || *report.AssignedTo != userID {
                        includeInternal = false
                }
        }
        
        // Get all comments
        allComments, err := s.reportRepo.GetReportComments(reportID, false)
        if err != nil {
                return nil, fmt.Errorf("error getting report comments: %w", err)
        }
        
        // Filter out internal comments if needed
        if !includeInternal {
                nonInternalComments := make([]models.ReportComment, 0)
                for _, comment := range allComments {
                        if !comment.IsInternal {
                                nonInternalComments = append(nonInternalComments, comment)
                        }
                }
                return nonInternalComments, nil
        }
        
        return allComments, nil
}

// DeleteReportComment deletes a comment
func (s *ReportServiceImpl) DeleteReportComment(commentID uint, userID uint) error {
        // Implementation would verify permissions first
        return s.reportRepo.DeleteReportComment(commentID)
}

// GetActionLogs retrieves action logs for a report
func (s *ReportServiceImpl) GetActionLogs(reportID uint) ([]models.ReportActionLog, error) {
        return s.reportRepo.GetActionLogs(reportID)
}

// Helper functions

// validateContent checks if the content exists
func (s *ReportServiceImpl) validateContent(contentType string, contentID uint) error {
        switch contentType {
        case "topic":
                _, err := s.discussionRepo.GetTopicByID(contentID)
                return err
        case "comment":
                _, err := s.commentRepo.GetCommentByID(contentID)
                return err
        default:
                return errors.New("invalid content type")
        }
}

// isValidStatusTransition checks if a status transition is valid
func (s *ReportServiceImpl) isValidStatusTransition(from, to models.ReportStatus) bool {
        switch from {
        case models.StatusPending:
                return to == models.StatusInReview || to == models.StatusResolved || to == models.StatusRejected
        case models.StatusInReview:
                return to == models.StatusResolved || to == models.StatusRejected
        case models.StatusResolved, models.StatusRejected:
                return to == models.StatusInReview // Allow reopening
        default:
                return false
        }
}

// executeResolution performs the appropriate action based on the resolution
func (s *ReportServiceImpl) executeResolution(report *models.ContentReport, resolution models.ReportResolutionType, actionUserID uint) error {
        // Execute actions based on resolution
        switch resolution {
        case models.ResolutionNoAction, models.ResolutionWarning:
                // No action needed on the content
                return nil
                
        case models.ResolutionContentRemoved:
                // Remove the reported content
                switch report.ContentType {
                case "topic":
                        // Mark topic as deleted
                        return s.discussionRepo.DeleteTopic(report.ContentID)
                case "comment":
                        // Mark comment as deleted
                        return s.commentRepo.DeleteComment(report.ContentID)
                default:
                        return errors.New("invalid content type")
                }
                
        case models.ResolutionContentEdited:
                // In a real implementation, this would store a reference to edited content
                // or apply automatic content filtering
                return nil
                
        case models.ResolutionUserSuspended, models.ResolutionUserBanned:
                // In a real implementation, this would trigger user suspension or ban
                // through the user service
                return nil
                
        default:
                return errors.New("invalid resolution type")
        }
}