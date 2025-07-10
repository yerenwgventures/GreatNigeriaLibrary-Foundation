package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/service"
)

// BookHandlers contains handlers for book-related operations
type BookHandlers struct {
	bookService       service.BookService
	bookImportService service.BookImportService
}

// NewBookHandlers creates a new BookHandlers instance
func NewBookHandlers(bookService service.BookService, bookImportService service.BookImportService) *BookHandlers {
	return &BookHandlers{
		bookService:       bookService,
		bookImportService: bookImportService,
	}
}

// RegisterRoutes registers book-related routes
func (h *BookHandlers) RegisterRoutes(router *gin.RouterGroup) {
	books := router.Group("/books")
	{
		books.GET("", h.GetAllBooks)
		books.GET("/:id", h.GetBookByID)
		books.GET("/:id/content", h.GetBookContent)
		books.GET("/:id/frontmatter", h.GetBookFrontMatter)
		books.GET("/:id/backmatter", h.GetBookBackMatter)
		books.GET("/:id/chapters", h.GetBookChapters)

		books.GET("/chapters/:id", h.GetChapterByID)
		books.GET("/chapters/:id/sections", h.GetChapterSections)

		books.GET("/sections/:id", h.GetSectionByID)
		books.GET("/sections/:id/forumtopics", h.GetSectionForumTopics)
		books.GET("/sections/:id/actionsteps", h.GetSectionActionSteps)
		books.GET("/sections/:id/interactive", h.GetSectionInteractiveElements)
		books.GET("/sections/:id/subsections", h.GetSectionSubsections)

		books.GET("/subsections/:id", h.GetSubsectionByID)

		// User interaction endpoints
		books.POST("/:id/progress", h.TrackUserProgress)
		books.GET("/:id/progress", h.GetUserBookProgress)

		books.POST("/:id/bookmarks", h.CreateBookmark)
		books.GET("/:id/bookmarks", h.GetUserBookmarks)
		books.DELETE("/bookmarks/:id", h.DeleteBookmark)

		books.POST("/:id/notes", h.CreateNote)
		books.GET("/:id/notes", h.GetUserNotes)
		books.PUT("/notes/:id", h.UpdateNote)
		books.DELETE("/notes/:id", h.DeleteNote)

		books.POST("/:id/feedback/mood", h.SubmitMoodFeedback)
		books.POST("/:id/feedback/difficulty", h.SubmitDifficultyFeedback)
		books.GET("/:id/feedback/summary", h.GetFeedbackSummary)

		// Admin endpoints
		admin := books.Group("/admin")
		{
			admin.POST("/import/book1", h.ImportBook1)
		}
	}
}

// GetAllBooks handles GET /api/books
func (h *BookHandlers) GetAllBooks(c *gin.Context) {
	includeUnpublishedStr := c.DefaultQuery("include_unpublished", "false")
	includeUnpublished := includeUnpublishedStr == "true"

	// Only admins can see unpublished books
	if includeUnpublished {
		// In a real implementation, check user role here
		// For now, we'll allow it
	}

	// Get the includeChapters parameter from the query string
	includeChaptersStr := c.DefaultQuery("include_chapters", "false")
	includeChapters := includeChaptersStr == "true"

	var books interface{}
	var err error

	if includeChapters {
		// Use GetAllBooksWithChapters to include chapters
		books, err = h.bookService.GetAllBooksWithChapters(includeUnpublished)
	} else {
		// Use GetAllBooks for better performance when chapters aren't needed
		books, err = h.bookService.GetAllBooks(includeUnpublished)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetBookByID handles GET /api/books/:id
func (h *BookHandlers) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Use GetBookWithChapters instead of GetBookByID to include chapters
	book, err := h.bookService.GetBookWithChapters(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// GetBookContent handles GET /api/books/:id/content
func (h *BookHandlers) GetBookContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	includeUnpublishedStr := c.DefaultQuery("include_unpublished", "false")
	includeUnpublished := includeUnpublishedStr == "true"

	// Only admins can see unpublished content
	if includeUnpublished {
		// In a real implementation, check user role here
		// For now, we'll allow it
	}

	content, err := h.bookService.GetBookContent(uint(id), includeUnpublished)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book content not found"})
		return
	}

	c.JSON(http.StatusOK, content)
}

// GetBookFrontMatter handles GET /api/books/:id/frontmatter
func (h *BookHandlers) GetBookFrontMatter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	includeUnpublishedStr := c.DefaultQuery("include_unpublished", "false")
	includeUnpublished := includeUnpublishedStr == "true"

	frontMatter, err := h.bookService.GetBookFrontMatter(uint(id), includeUnpublished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch front matter"})
		return
	}

	c.JSON(http.StatusOK, frontMatter)
}

// GetBookBackMatter handles GET /api/books/:id/backmatter
func (h *BookHandlers) GetBookBackMatter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	includeUnpublishedStr := c.DefaultQuery("include_unpublished", "false")
	includeUnpublished := includeUnpublishedStr == "true"

	backMatter, err := h.bookService.GetBookBackMatter(uint(id), includeUnpublished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch back matter"})
		return
	}

	c.JSON(http.StatusOK, backMatter)
}

// GetBookChapters handles GET /api/books/:id/chapters
func (h *BookHandlers) GetBookChapters(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	includeUnpublishedStr := c.DefaultQuery("include_unpublished", "false")
	includeUnpublished := includeUnpublishedStr == "true"

	chapters, err := h.bookService.GetBookChapters(uint(id), includeUnpublished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chapters"})
		return
	}

	c.JSON(http.StatusOK, chapters)
}

// GetChapterByID handles GET /api/books/chapters/:id
func (h *BookHandlers) GetChapterByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	chapter, err := h.bookService.GetChapter(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	c.JSON(http.StatusOK, chapter)
}

// GetChapterSections handles GET /api/books/chapters/:id/sections
func (h *BookHandlers) GetChapterSections(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	includeUnpublishedStr := c.DefaultQuery("include_unpublished", "false")
	includeUnpublished := includeUnpublishedStr == "true"

	sections, err := h.bookService.GetChapterSections(uint(id), includeUnpublished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sections"})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// GetSectionByID handles GET /api/books/sections/:id
func (h *BookHandlers) GetSectionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	section, err := h.bookService.GetSection(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	// Get interactive elements for this section if requested
	includeInteractive := c.Query("include_interactive") == "true"
	fmt.Printf("includeInteractive parameter: %v\n", includeInteractive)

	if includeInteractive {
		fmt.Printf("Fetching interactive elements for section ID: %d\n", id)
		interactiveElements, err := h.bookService.GetInteractiveElementsBySectionID(uint(id))
		if err != nil {
			fmt.Printf("Error fetching interactive elements: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interactive elements"})
			return
		}
		fmt.Printf("Found %d interactive elements\n", len(interactiveElements))

		// Return section with interactive elements
		response := gin.H{
			"section":              section,
			"interactive_elements": interactiveElements,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, section)
}

// GetSectionForumTopics handles GET /api/books/sections/:id/forumtopics
func (h *BookHandlers) GetSectionForumTopics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	topics, err := h.bookService.GetForumTopicsForSection(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch forum topics"})
		return
	}

	c.JSON(http.StatusOK, topics)
}

// GetSectionActionSteps handles GET /api/books/sections/:id/actionsteps
func (h *BookHandlers) GetSectionActionSteps(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	steps, err := h.bookService.GetActionStepsForSection(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch action steps"})
		return
	}

	c.JSON(http.StatusOK, steps)
}

// GetSectionSubsections handles GET /api/books/sections/:id/subsections
func (h *BookHandlers) GetSectionSubsections(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	includeUnpublishedStr := c.DefaultQuery("include_unpublished", "false")
	includeUnpublished := includeUnpublishedStr == "true"

	// Get the section with its subsections
	sectionWithSubsections, err := h.bookService.GetSectionWithSubsections(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch section with subsections"})
		return
	}

	// Alternative way to get just the subsections
	subsections, err := h.bookService.GetSubsectionsBySectionID(uint(id), includeUnpublished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subsections"})
		return
	}

	// Determine response based on query parameter
	fullSectionStr := c.DefaultQuery("full_section", "false")
	fullSection := fullSectionStr == "true"

	if fullSection {
		c.JSON(http.StatusOK, sectionWithSubsections)
	} else {
		c.JSON(http.StatusOK, subsections)
	}
}

// GetSubsectionByID handles GET /api/books/subsections/:id
func (h *BookHandlers) GetSubsectionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subsection ID"})
		return
	}

	subsection, err := h.bookService.GetSubsection(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subsection not found"})
		return
	}

	// Check if we should return rendered content
	renderStr := c.DefaultQuery("render", "false")
	render := renderStr == "true"

	if render {
		renderedContent, err := h.bookService.GetRenderedSubsection(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render subsection content"})
			return
		}

		// Return subsection with rendered content
		response := gin.H{
			"subsection":       subsection,
			"rendered_content": renderedContent,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, subsection)
}

// GetSectionInteractiveElements handles GET /api/books/sections/:id/interactive
func (h *BookHandlers) GetSectionInteractiveElements(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	elements, err := h.bookService.GetInteractiveElementsBySectionID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interactive elements"})
		return
	}

	c.JSON(http.StatusOK, elements)
}

// TrackUserProgress handles POST /api/books/:id/progress
func (h *BookHandlers) TrackUserProgress(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var request struct {
		UserID    uint    `json:"user_id" binding:"required"`
		SectionID uint    `json:"section_id" binding:"required"`
		IsRead    bool    `json:"is_read"`
		Progress  float64 `json:"progress"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.bookService.TrackUserProgress(
		request.UserID,
		uint(bookID),
		request.SectionID,
		request.IsRead,
		request.Progress,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress tracked successfully"})
}

// GetUserBookProgress handles GET /api/books/:id/progress
func (h *BookHandlers) GetUserBookProgress(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	progress, err := h.bookService.GetUserBookProgress(uint(userID), uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"progress": progress})
}

// CreateBookmark handles POST /api/books/:id/bookmarks
func (h *BookHandlers) CreateBookmark(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var request struct {
		UserID      uint   `json:"user_id" binding:"required"`
		SectionID   uint   `json:"section_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookmark, err := h.bookService.CreateBookmark(
		request.UserID,
		uint(bookID),
		request.SectionID,
		request.Title,
		request.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark"})
		return
	}

	c.JSON(http.StatusCreated, bookmark)
}

// GetUserBookmarks handles GET /api/books/:id/bookmarks
func (h *BookHandlers) GetUserBookmarks(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	bookmarks, err := h.bookService.GetUserBookmarks(uint(userID), uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookmarks"})
		return
	}

	c.JSON(http.StatusOK, bookmarks)
}

// DeleteBookmark handles DELETE /api/books/bookmarks/:id
func (h *BookHandlers) DeleteBookmark(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	err = h.bookService.DeleteBookmark(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bookmark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
}

// CreateNote handles POST /api/books/:id/notes
func (h *BookHandlers) CreateNote(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var request struct {
		UserID    uint   `json:"user_id" binding:"required"`
		SectionID uint   `json:"section_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
		IsPrivate bool   `json:"is_private"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.bookService.CreateNote(
		request.UserID,
		uint(bookID),
		request.SectionID,
		request.Content,
		request.IsPrivate,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetUserNotes handles GET /api/books/:id/notes
func (h *BookHandlers) GetUserNotes(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	notes, err := h.bookService.GetUserNotes(uint(userID), uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// UpdateNote handles PUT /api/books/notes/:id
func (h *BookHandlers) UpdateNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var request struct {
		Content   string `json:"content" binding:"required"`
		IsPrivate bool   `json:"is_private"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.bookService.UpdateNote(uint(id), request.Content, request.IsPrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNote handles DELETE /api/books/notes/:id
func (h *BookHandlers) DeleteNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	err = h.bookService.DeleteNote(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

// SubmitMoodFeedback handles POST /api/books/:id/feedback/mood
func (h *BookHandlers) SubmitMoodFeedback(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var request struct {
		UserID    uint                   `json:"user_id" binding:"required"`
		SectionID uint                   `json:"section_id"`
		Rating    int                    `json:"rating" binding:"required,min=1,max=5"`
		Comment   string                 `json:"comment"`
		Metadata  map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback, err := h.bookService.SubmitFeedback(
		request.UserID,
		uint(bookID),
		request.SectionID,
		"mood",
		request.Rating,
		request.Comment,
		request.Metadata,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit mood feedback"})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// SubmitDifficultyFeedback handles POST /api/books/:id/feedback/difficulty
func (h *BookHandlers) SubmitDifficultyFeedback(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var request struct {
		UserID    uint                   `json:"user_id" binding:"required"`
		SectionID uint                   `json:"section_id"`
		Rating    int                    `json:"rating" binding:"required,min=1,max=5"`
		Comment   string                 `json:"comment"`
		Metadata  map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback, err := h.bookService.SubmitFeedback(
		request.UserID,
		uint(bookID),
		request.SectionID,
		"difficulty",
		request.Rating,
		request.Comment,
		request.Metadata,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit difficulty feedback"})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// GetFeedbackSummary handles GET /api/books/:id/feedback/summary
func (h *BookHandlers) GetFeedbackSummary(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	stats, err := h.bookService.GetFeedbackStats(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feedback stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ImportBook1 handles POST /api/books/admin/import/book1
func (h *BookHandlers) ImportBook1(c *gin.Context) {
	err := h.bookImportService.ImportBook1()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import Book 1"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book 1 imported successfully"})
}
