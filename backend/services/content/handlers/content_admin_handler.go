package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/service"
)

// ContentAdminHandler handles content administration endpoints
type ContentAdminHandler struct {
	adminService service.ContentAdminService
}

// NewContentAdminHandler creates a new content admin handler
func NewContentAdminHandler(adminService service.ContentAdminService) *ContentAdminHandler {
	return &ContentAdminHandler{
		adminService: adminService,
	}
}

// RegisterRoutes registers the routes for content administration
func (h *ContentAdminHandler) RegisterRoutes(router *gin.Engine) {
	contentAdmin := router.Group("/api/v1/content-admin")
	{
		// Import endpoints
		contentAdmin.POST("/import/books", h.ImportBooks)
		contentAdmin.POST("/import/chapters/:bookID", h.ImportChapters)
		contentAdmin.POST("/import/sections/:chapterID", h.ImportSections)

		// Export endpoints
		contentAdmin.GET("/export/books", h.ExportBooks)
		contentAdmin.GET("/export/chapters/:bookID", h.ExportChapters)
		contentAdmin.GET("/export/sections/:chapterID", h.ExportSections)

		// Revision endpoints
		contentAdmin.POST("/revisions/books/:bookID", h.CreateBookRevision)
		contentAdmin.POST("/revisions/chapters/:chapterID", h.CreateChapterRevision)
		contentAdmin.POST("/revisions/sections/:sectionID", h.CreateSectionRevision)
		contentAdmin.GET("/revisions/books/:bookID", h.GetBookRevisions)
		contentAdmin.GET("/revisions/chapters/:chapterID", h.GetChapterRevisions)
		contentAdmin.GET("/revisions/sections/:sectionID", h.GetSectionRevisions)
		contentAdmin.POST("/revisions/books/:revisionID/restore", h.RestoreBookRevision)
		contentAdmin.POST("/revisions/chapters/:revisionID/restore", h.RestoreChapterRevision)
		contentAdmin.POST("/revisions/sections/:revisionID/restore", h.RestoreSectionRevision)

		// Publishing endpoints
		contentAdmin.POST("/publishing/books/:bookID/schedule", h.ScheduleBookPublishing)
		contentAdmin.POST("/publishing/chapters/:chapterID/schedule", h.ScheduleChapterPublishing)
		contentAdmin.POST("/publishing/sections/:sectionID/schedule", h.ScheduleSectionPublishing)
		contentAdmin.GET("/publishing/scheduled", h.GetScheduledContent)
		contentAdmin.POST("/publishing/:contentType/:contentID/publish", h.PublishContent)
		contentAdmin.POST("/publishing/:contentType/:contentID/unpublish", h.UnpublishContent)
	}
}

// ImportBooks imports books from a file
func (h *ContentAdminHandler) ImportBooks(c *gin.Context) {
	// Get format from query parameter
	format := c.DefaultQuery("format", "json")

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File required"})
		return
	}

	// Open file
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileReader.Close()

	var books interface{}

	// Import based on format
	switch format {
	case "json":
		books, err = h.adminService.ImportBooksFromJSON(fileReader)
	case "csv":
		books, err = h.adminService.ImportBooksFromCSV(fileReader)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Books imported successfully", "books": books})
}

// ImportChapters imports chapters from a file
func (h *ContentAdminHandler) ImportChapters(c *gin.Context) {
	// Get book ID from path
	bookIDStr := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Get format from query parameter
	format := c.DefaultQuery("format", "json")

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File required"})
		return
	}

	// Open file
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileReader.Close()

	var chapters interface{}

	// Import based on format
	switch format {
	case "json":
		chapters, err = h.adminService.ImportChaptersFromJSON(fileReader, uint(bookID))
	case "csv":
		chapters, err = h.adminService.ImportChaptersFromCSV(fileReader, uint(bookID))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapters imported successfully", "chapters": chapters})
}

// ImportSections imports sections from a file
func (h *ContentAdminHandler) ImportSections(c *gin.Context) {
	// Get chapter ID from path
	chapterIDStr := c.Param("chapterID")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	// Get format from query parameter
	format := c.DefaultQuery("format", "json")

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File required"})
		return
	}

	// Open file
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileReader.Close()

	var sections interface{}

	// Import based on format
	switch format {
	case "json":
		sections, err = h.adminService.ImportSectionsFromJSON(fileReader, uint(chapterID))
	case "csv":
		sections, err = h.adminService.ImportSectionsFromCSV(fileReader, uint(chapterID))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sections imported successfully", "sections": sections})
}

// ExportBooks exports books to a file
func (h *ContentAdminHandler) ExportBooks(c *gin.Context) {
	// Get format from query parameter
	format := c.DefaultQuery("format", "json")

	// Get book IDs from query parameters
	var bookIDs []uint
	bookIDsStr := c.QueryArray("bookID")
	for _, idStr := range bookIDsStr {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}
		bookIDs = append(bookIDs, uint(id))
	}

	// Set content type based on format
	switch format {
	case "json":
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=books.json")
		err := h.adminService.ExportBooksToJSON(c.Writer, bookIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	case "csv":
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=books.csv")
		err := h.adminService.ExportBooksToCSV(c.Writer, bookIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format"})
	}
}

// ExportChapters exports chapters to a file
func (h *ContentAdminHandler) ExportChapters(c *gin.Context) {
	// Get book ID from path
	bookIDStr := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Get format from query parameter
	format := c.DefaultQuery("format", "json")

	// Set content type based on format
	switch format {
	case "json":
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=book_%d_chapters.json", bookID))
		err := h.adminService.ExportChaptersToJSON(c.Writer, uint(bookID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	case "csv":
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=book_%d_chapters.csv", bookID))
		err := h.adminService.ExportChaptersToCSV(c.Writer, uint(bookID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format"})
	}
}

// ExportSections exports sections to a file
func (h *ContentAdminHandler) ExportSections(c *gin.Context) {
	// Get chapter ID from path
	chapterIDStr := c.Param("chapterID")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	// Get format from query parameter
	format := c.DefaultQuery("format", "json")

	// Set content type based on format
	switch format {
	case "json":
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=chapter_%d_sections.json", chapterID))
		err := h.adminService.ExportSectionsToJSON(c.Writer, uint(chapterID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	case "csv":
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=chapter_%d_sections.csv", chapterID))
		err := h.adminService.ExportSectionsToCSV(c.Writer, uint(chapterID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format"})
	}
}

// CreateBookRevisionRequest represents a request to create a book revision
type CreateBookRevisionRequest struct {
	Changes map[string]interface{} `json:"changes" binding:"required"`
	Notes   string                 `json:"notes" binding:"required"`
}

// CreateBookRevision creates a new revision for a book
func (h *ContentAdminHandler) CreateBookRevision(c *gin.Context) {
	// Get book ID from path
	bookIDStr := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var req CreateBookRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create revision
	revision, err := h.adminService.CreateBookRevision(uint(bookID), req.Changes, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, revision)
}

// CreateChapterRevisionRequest represents a request to create a chapter revision
type CreateChapterRevisionRequest struct {
	Changes map[string]interface{} `json:"changes" binding:"required"`
	Notes   string                 `json:"notes" binding:"required"`
}

// CreateChapterRevision creates a new revision for a chapter
func (h *ContentAdminHandler) CreateChapterRevision(c *gin.Context) {
	// Get chapter ID from path
	chapterIDStr := c.Param("chapterID")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var req CreateChapterRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create revision
	revision, err := h.adminService.CreateChapterRevision(uint(chapterID), req.Changes, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, revision)
}

// CreateSectionRevisionRequest represents a request to create a section revision
type CreateSectionRevisionRequest struct {
	Changes map[string]interface{} `json:"changes" binding:"required"`
	Notes   string                 `json:"notes" binding:"required"`
}

// CreateSectionRevision creates a new revision for a section
func (h *ContentAdminHandler) CreateSectionRevision(c *gin.Context) {
	// Get section ID from path
	sectionIDStr := c.Param("sectionID")
	sectionID, err := strconv.ParseUint(sectionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	var req CreateSectionRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create revision
	revision, err := h.adminService.CreateSectionRevision(uint(sectionID), req.Changes, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, revision)
}

// GetBookRevisions gets all revisions for a book
func (h *ContentAdminHandler) GetBookRevisions(c *gin.Context) {
	// Get book ID from path
	bookIDStr := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Get revisions
	revisions, err := h.adminService.GetBookRevisions(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, revisions)
}

// GetChapterRevisions gets all revisions for a chapter
func (h *ContentAdminHandler) GetChapterRevisions(c *gin.Context) {
	// Get chapter ID from path
	chapterIDStr := c.Param("chapterID")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	// Get revisions
	revisions, err := h.adminService.GetChapterRevisions(uint(chapterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, revisions)
}

// GetSectionRevisions gets all revisions for a section
func (h *ContentAdminHandler) GetSectionRevisions(c *gin.Context) {
	// Get section ID from path
	sectionIDStr := c.Param("sectionID")
	sectionID, err := strconv.ParseUint(sectionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	// Get revisions
	revisions, err := h.adminService.GetSectionRevisions(uint(sectionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, revisions)
}

// RestoreBookRevision restores a book to a previous revision
func (h *ContentAdminHandler) RestoreBookRevision(c *gin.Context) {
	// Get revision ID from path
	revisionIDStr := c.Param("revisionID")
	revisionID, err := strconv.ParseUint(revisionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revision ID"})
		return
	}

	// Restore revision
	if err := h.adminService.RestoreBookRevision(uint(revisionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book restored to revision successfully"})
}

// RestoreChapterRevision restores a chapter to a previous revision
func (h *ContentAdminHandler) RestoreChapterRevision(c *gin.Context) {
	// Get revision ID from path
	revisionIDStr := c.Param("revisionID")
	revisionID, err := strconv.ParseUint(revisionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revision ID"})
		return
	}

	// Restore revision
	if err := h.adminService.RestoreChapterRevision(uint(revisionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter restored to revision successfully"})
}

// RestoreSectionRevision restores a section to a previous revision
func (h *ContentAdminHandler) RestoreSectionRevision(c *gin.Context) {
	// Get revision ID from path
	revisionIDStr := c.Param("revisionID")
	revisionID, err := strconv.ParseUint(revisionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revision ID"})
		return
	}

	// Restore revision
	if err := h.adminService.RestoreSectionRevision(uint(revisionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Section restored to revision successfully"})
}

// SchedulePublishingRequest represents a request to schedule publishing
type SchedulePublishingRequest struct {
	PublishDate string `json:"publishDate" binding:"required"` // ISO 8601 date format
}

// ScheduleBookPublishing schedules a book to be published
func (h *ContentAdminHandler) ScheduleBookPublishing(c *gin.Context) {
	// Get book ID from path
	bookIDStr := c.Param("bookID")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var req SchedulePublishingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse publish date
	publishDate, err := time.Parse(time.RFC3339, req.PublishDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publish date format. Use ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)"})
		return
	}

	// Schedule publishing
	if err := h.adminService.ScheduleBookPublishing(uint(bookID), publishDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book publishing scheduled successfully", "publishDate": publishDate})
}

// ScheduleChapterPublishing schedules a chapter to be published
func (h *ContentAdminHandler) ScheduleChapterPublishing(c *gin.Context) {
	// Get chapter ID from path
	chapterIDStr := c.Param("chapterID")
	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var req SchedulePublishingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse publish date
	publishDate, err := time.Parse(time.RFC3339, req.PublishDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publish date format. Use ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)"})
		return
	}

	// Schedule publishing
	if err := h.adminService.ScheduleChapterPublishing(uint(chapterID), publishDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter publishing scheduled successfully", "publishDate": publishDate})
}

// ScheduleSectionPublishing schedules a section to be published
func (h *ContentAdminHandler) ScheduleSectionPublishing(c *gin.Context) {
	// Get section ID from path
	sectionIDStr := c.Param("sectionID")
	sectionID, err := strconv.ParseUint(sectionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	var req SchedulePublishingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse publish date
	publishDate, err := time.Parse(time.RFC3339, req.PublishDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publish date format. Use ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)"})
		return
	}

	// Schedule publishing
	if err := h.adminService.ScheduleSectionPublishing(uint(sectionID), publishDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Section publishing scheduled successfully", "publishDate": publishDate})
}

// GetScheduledContent gets all content scheduled to be published
func (h *ContentAdminHandler) GetScheduledContent(c *gin.Context) {
	// Get scheduled content
	scheduled, err := h.adminService.GetScheduledContent()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scheduled)
}

// PublishContent publishes content immediately
func (h *ContentAdminHandler) PublishContent(c *gin.Context) {
	// Get content type and ID from path
	contentType := c.Param("contentType")
	contentIDStr := c.Param("contentID")
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	// Validate content type
	if contentType != "book" && contentType != "chapter" && contentType != "section" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}

	// Publish content
	if err := h.adminService.PublishContent(contentType, uint(contentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s published successfully", contentType)})
}

// UnpublishContent unpublishes content
func (h *ContentAdminHandler) UnpublishContent(c *gin.Context) {
	// Get content type and ID from path
	contentType := c.Param("contentType")
	contentIDStr := c.Param("contentID")
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	// Validate content type
	if contentType != "book" && contentType != "chapter" && contentType != "section" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}

	// Unpublish content
	if err := h.adminService.UnpublishContent(contentType, uint(contentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s unpublished successfully", contentType)})
}
