package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/service"
)

// CitationHandler handles HTTP requests related to citations
type CitationHandler struct {
	citationService *service.CitationService
}

// NewCitationHandler creates a new citation handler
func NewCitationHandler(citationService *service.CitationService) *CitationHandler {
	return &CitationHandler{
		citationService: citationService,
	}
}

// GetCitationsByBook retrieves all citations for a book
func (h *CitationHandler) GetCitationsByBook(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	citations, err := h.citationService.GetCitationsByBook(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get citations: %v", err)})
		return
	}

	c.JSON(http.StatusOK, citations)
}

// GetCitationByID retrieves a specific citation
func (h *CitationHandler) GetCitationByID(c *gin.Context) {
	citationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid citation ID"})
		return
	}

	citation, err := h.citationService.GetCitationByID(uint(citationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get citation: %v", err)})
		return
	}

	c.JSON(http.StatusOK, citation)
}

// CreateCitation adds a new citation
func (h *CitationHandler) CreateCitation(c *gin.Context) {
	var citation models.Citation
	if err := c.ShouldBindJSON(&citation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid citation data: %v", err)})
		return
	}

	if err := h.citationService.CreateCitation(&citation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create citation: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, citation)
}

// RecordCitationUsage records a citation being used in the book
type CitationUsageRequest struct {
	CitationID uint `json:"citation_id" binding:"required"`
	BookID     uint `json:"book_id" binding:"required"`
	ChapterID  uint `json:"chapter_id" binding:"required"`
	SectionID  uint `json:"section_id" binding:"required"`
}

func (h *CitationHandler) RecordCitationUsage(c *gin.Context) {
	var request CitationUsageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid usage data: %v", err)})
		return
	}

	if err := h.citationService.RecordCitationUsage(
		request.CitationID, request.BookID, request.ChapterID, request.SectionID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to record usage: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Citation usage recorded successfully"})
}

// GenerateBibliography creates a bibliography for a book
func (h *CitationHandler) GenerateBibliography(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	bibliography, err := h.citationService.GenerateBibliography(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate bibliography: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bibliography": bibliography})
}

// GetCitationStats retrieves citation statistics for a book
func (h *CitationHandler) GetCitationStats(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	stats, err := h.citationService.GetCitationStats(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get citation stats: %v", err)})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ImportCitationsFromJSON imports citations from a JSON file
type ImportCitationsRequest struct {
	BookID   uint   `json:"book_id" binding:"required"`
	JSONData string `json:"json_data" binding:"required"`
}

func (h *CitationHandler) ImportCitationsFromJSON(c *gin.Context) {
	var request ImportCitationsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid import data: %v", err)})
		return
	}

	if err := h.citationService.ImportCitationsFromJSON(request.BookID, request.JSONData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to import citations: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Citations imported successfully"})
}

// UpdateAppendixWithBibliography updates the appendix with the bibliography content
func (h *CitationHandler) UpdateAppendixWithBibliography(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := h.citationService.UpdateAppendixWithBibliography(uint(bookID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update appendix: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bibliography appendix updated successfully"})
}

// RegisterRoutes registers routes for this handler
func (h *CitationHandler) RegisterRoutes(router *gin.RouterGroup) {
	citations := router.Group("/citations")
	{
		citations.GET("/book/:id", h.GetCitationsByBook)
		citations.GET("/:id", h.GetCitationByID)
		citations.POST("", h.CreateCitation)
		citations.POST("/usage", h.RecordCitationUsage)
		citations.GET("/bibliography/:id", h.GenerateBibliography)
		citations.GET("/stats/:id", h.GetCitationStats)
		citations.POST("/import", h.ImportCitationsFromJSON)
		citations.POST("/bibliography/appendix/:id", h.UpdateAppendixWithBibliography)
	}
}