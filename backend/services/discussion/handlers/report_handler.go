package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/service"
)

// ReportHandler defines the handler for report endpoints
type ReportHandler struct {
	reportService service.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// RegisterRoutes registers the routes for report management
func (h *ReportHandler) RegisterRoutes(router *gin.RouterGroup) {
	reports := router.Group("/reports")
	{
		// User-facing endpoints
		reports.POST("", h.CreateReport)
		reports.GET("/me", h.GetMyReports)
		reports.GET("/:id", h.GetReport)
		
		// Evidence operations
		reports.POST("/:id/evidence", h.AddEvidence)
		reports.GET("/:id/evidence", h.GetEvidence)
		reports.DELETE("/evidence/:evidenceId", h.DeleteEvidence)
		
		// Comment operations
		reports.POST("/:id/comments", h.AddComment)
		reports.GET("/:id/comments", h.GetComments)
		reports.DELETE("/comments/:commentId", h.DeleteComment)
		
		// Moderation endpoints (would be protected by moderator role)
		moderator := reports.Group("/moderation")
		{
			moderator.GET("/pending", h.GetPendingReports)
			moderator.GET("/in-review", h.GetInReviewReports)
			moderator.GET("/resolved", h.GetResolvedReports)
			moderator.GET("/rejected", h.GetRejectedReports)
			moderator.GET("/category/:category", h.GetReportsByCategory)
			moderator.GET("/stats", h.GetReportStats)
			moderator.POST("/:id/assign", h.AssignReport)
			moderator.POST("/:id/status", h.UpdateReportStatus)
			moderator.POST("/:id/resolve", h.ResolveReport)
			moderator.GET("/:id/logs", h.GetActionLogs)
		}
	}
}

// CreateReportRequest represents a request to create a content report
type CreateReportRequest struct {
	ContentType    string `json:"contentType" binding:"required"`
	ContentID      uint   `json:"contentId" binding:"required"`
	Category       string `json:"category" binding:"required"`
	Reason         string `json:"reason" binding:"required"`
	AdditionalInfo string `json:"additionalInfo"`
}

// CreateReport creates a new content report
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Validate content type
	if req.ContentType != "topic" && req.ContentType != "comment" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
		return
	}
	
	// Parse category
	var category models.ReportCategory
	switch req.Category {
	case "spam":
		category = models.CategorySpam
	case "harassment":
		category = models.CategoryHarassment
	case "hate_speech":
		category = models.CategoryHateSpeech
	case "violence":
		category = models.CategoryViolence
	case "illegal_content":
		category = models.CategoryIllegalContent
	case "privacy_violation":
		category = models.CategoryPrivacyViolation
	case "copyright":
		category = models.CategoryCopyright
	case "misinformation":
		category = models.CategoryMisinformation
	case "other":
		category = models.CategoryOther
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report category"})
		return
	}
	
	// Create report
	report, err := h.reportService.CreateReport(
		userID.(uint),
		req.ContentType,
		req.ContentID,
		category,
		req.Reason,
		req.AdditionalInfo,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, report)
}

// GetMyReports retrieves reports created by the authenticated user
func (h *ReportHandler) GetMyReports(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get pagination parameters
	page, pageSize := getPaginationParams(c)
	
	// Get reports
	reports, total, err := h.reportService.GetReportsByReporter(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"pagination": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetReport retrieves a report by ID
func (h *ReportHandler) GetReport(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get report
	report, err := h.reportService.GetReportByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Check if user has access to this report
	// User must be the reporter or a moderator assigned to this report
	if report.ReporterID != userID.(uint) && (report.AssignedTo == nil || *report.AssignedTo != userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this report"})
		return
	}
	
	c.JSON(http.StatusOK, report)
}

// AddEvidenceRequest represents a request to add evidence to a report
type AddEvidenceRequest struct {
	Type     string `json:"type" binding:"required"`
	Content  string `json:"content"`
	FilePath string `json:"filePath"`
	URL      string `json:"url"`
}

// AddEvidence adds evidence to a report
func (h *ReportHandler) AddEvidence(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	var req AddEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Add evidence
	evidence, err := h.reportService.AddReportEvidence(
		uint(id),
		userID.(uint),
		req.Type,
		req.Content,
		req.FilePath,
		req.URL,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, evidence)
}

// GetEvidence retrieves evidence for a report
func (h *ReportHandler) GetEvidence(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get report to check permissions
	report, err := h.reportService.GetReportByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Check if user has access to this report
	if report.ReporterID != userID.(uint) && (report.AssignedTo == nil || *report.AssignedTo != userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view evidence for this report"})
		return
	}
	
	// Get evidence
	evidence, err := h.reportService.GetReportEvidence(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, evidence)
}

// DeleteEvidence deletes evidence from a report
func (h *ReportHandler) DeleteEvidence(c *gin.Context) {
	// Get evidence ID from path
	evidenceIDStr := c.Param("evidenceId")
	evidenceID, err := strconv.ParseUint(evidenceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Delete evidence
	if err := h.reportService.DeleteEvidence(uint(evidenceID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Evidence deleted successfully"})
}

// AddCommentRequest represents a request to add a comment to a report
type AddCommentRequest struct {
	Comment    string `json:"comment" binding:"required"`
	IsInternal bool   `json:"isInternal"`
}

// AddComment adds a comment to a report
func (h *ReportHandler) AddComment(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Add comment
	comment, err := h.reportService.AddReportComment(
		uint(id),
		userID.(uint),
		req.Comment,
		req.IsInternal,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, comment)
}

// GetComments retrieves comments for a report
func (h *ReportHandler) GetComments(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Check if internal comments should be included
	includeInternalStr := c.DefaultQuery("includeInternal", "false")
	includeInternal := includeInternalStr == "true"
	
	// Get comments
	comments, err := h.reportService.GetReportComments(uint(id), includeInternal, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, comments)
}

// DeleteComment deletes a comment from a report
func (h *ReportHandler) DeleteComment(c *gin.Context) {
	// Get comment ID from path
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Delete comment
	if err := h.reportService.DeleteReportComment(uint(commentID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// GetPendingReports retrieves pending reports
func (h *ReportHandler) GetPendingReports(c *gin.Context) {
	// Get pagination parameters
	page, pageSize := getPaginationParams(c)
	
	// Get reports
	reports, total, err := h.reportService.GetReportsByStatus(models.StatusPending, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"pagination": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetInReviewReports retrieves in-review reports
func (h *ReportHandler) GetInReviewReports(c *gin.Context) {
	// Get pagination parameters
	page, pageSize := getPaginationParams(c)
	
	// Get reports
	reports, total, err := h.reportService.GetReportsByStatus(models.StatusInReview, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"pagination": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetResolvedReports retrieves resolved reports
func (h *ReportHandler) GetResolvedReports(c *gin.Context) {
	// Get pagination parameters
	page, pageSize := getPaginationParams(c)
	
	// Get reports
	reports, total, err := h.reportService.GetReportsByStatus(models.StatusResolved, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"pagination": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetRejectedReports retrieves rejected reports
func (h *ReportHandler) GetRejectedReports(c *gin.Context) {
	// Get pagination parameters
	page, pageSize := getPaginationParams(c)
	
	// Get reports
	reports, total, err := h.reportService.GetReportsByStatus(models.StatusRejected, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"pagination": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetReportsByCategory retrieves reports by category
func (h *ReportHandler) GetReportsByCategory(c *gin.Context) {
	// Get category from path
	categoryStr := c.Param("category")
	
	// Parse category
	var category models.ReportCategory
	switch categoryStr {
	case "spam":
		category = models.CategorySpam
	case "harassment":
		category = models.CategoryHarassment
	case "hate_speech":
		category = models.CategoryHateSpeech
	case "violence":
		category = models.CategoryViolence
	case "illegal_content":
		category = models.CategoryIllegalContent
	case "privacy_violation":
		category = models.CategoryPrivacyViolation
	case "copyright":
		category = models.CategoryCopyright
	case "misinformation":
		category = models.CategoryMisinformation
	case "other":
		category = models.CategoryOther
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report category"})
		return
	}
	
	// Get status filter if provided
	statusStr := c.Query("status")
	var status *models.ReportStatus
	if statusStr != "" {
		var s models.ReportStatus
		switch statusStr {
		case "pending":
			s = models.StatusPending
		case "in_review":
			s = models.StatusInReview
		case "resolved":
			s = models.StatusResolved
		case "rejected":
			s = models.StatusRejected
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status filter"})
			return
		}
		status = &s
	}
	
	// Get pagination parameters
	page, pageSize := getPaginationParams(c)
	
	// Get reports
	reports, total, err := h.reportService.GetReportsByCategoryAndStatus(category, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"pagination": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetReportStats retrieves report statistics
func (h *ReportHandler) GetReportStats(c *gin.Context) {
	// Get report stats
	stats, err := h.reportService.GetReportStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// AssignReportRequest represents a request to assign a report to a moderator
type AssignReportRequest struct {
	ModeratorID uint `json:"moderatorId" binding:"required"`
}

// AssignReport assigns a report to a moderator
func (h *ReportHandler) AssignReport(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	var req AssignReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Assign report
	if err := h.reportService.AssignReportToModerator(uint(id), req.ModeratorID, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Report assigned successfully"})
}

// UpdateReportStatusRequest represents a request to update a report's status
type UpdateReportStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateReportStatus updates a report's status
func (h *ReportHandler) UpdateReportStatus(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	var req UpdateReportStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Parse status
	var status models.ReportStatus
	switch req.Status {
	case "pending":
		status = models.StatusPending
	case "in_review":
		status = models.StatusInReview
	case "resolved":
		status = models.StatusResolved
	case "rejected":
		status = models.StatusRejected
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Update status
	if err := h.reportService.UpdateReportStatus(uint(id), status, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Report status updated successfully"})
}

// ResolveReportRequest represents a request to resolve a report
type ResolveReportRequest struct {
	Resolution string `json:"resolution" binding:"required"`
	Notes      string `json:"notes"`
}

// ResolveReport resolves a report
func (h *ReportHandler) ResolveReport(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	var req ResolveReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Parse resolution
	var resolution models.ReportResolutionType
	switch req.Resolution {
	case "no_action":
		resolution = models.ResolutionNoAction
	case "warning":
		resolution = models.ResolutionWarning
	case "content_removed":
		resolution = models.ResolutionContentRemoved
	case "content_edited":
		resolution = models.ResolutionContentEdited
	case "user_suspended":
		resolution = models.ResolutionUserSuspended
	case "user_banned":
		resolution = models.ResolutionUserBanned
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resolution"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Resolve report
	if err := h.reportService.ResolveReport(uint(id), resolution, req.Notes, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Report resolved successfully"})
}

// GetActionLogs retrieves action logs for a report
func (h *ReportHandler) GetActionLogs(c *gin.Context) {
	// Get report ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}
	
	// Get action logs
	logs, err := h.reportService.GetActionLogs(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, logs)
}

// Helper functions

// getPaginationParams extracts pagination parameters from request
func getPaginationParams(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	return page, pageSize
}