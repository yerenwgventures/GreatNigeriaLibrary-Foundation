package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/service"
)

// ContentLinkHandler defines the handler for content link endpoints
type ContentLinkHandler struct {
	contentLinkService service.ContentLinkService
}

// NewContentLinkHandler creates a new content link handler
func NewContentLinkHandler(contentLinkService service.ContentLinkService) *ContentLinkHandler {
	return &ContentLinkHandler{
		contentLinkService: contentLinkService,
	}
}

// RegisterRoutes registers the routes for content link operations
func (h *ContentLinkHandler) RegisterRoutes(router *gin.RouterGroup) {
	contentLinks := router.Group("/content-links")
	{
		// Topic content links
		contentLinks.POST("/topics", h.LinkTopicToContent)
		contentLinks.GET("/topics/:topicId", h.GetTopicContentLinks)
		contentLinks.GET("/content/:contentType/:contentId/topics", h.GetLinkedTopics)
		contentLinks.PUT("/topics/:linkId/highlight", h.HighlightTopicLink)
		contentLinks.DELETE("/topics/:linkId", h.RemoveTopicLink)
		
		// Comment content links (citations)
		contentLinks.POST("/comments", h.AddContentCitation)
		contentLinks.GET("/comments/:commentId", h.GetCommentContentLinks)
		contentLinks.GET("/content/:contentType/:contentId/comments", h.GetContentCitations)
		contentLinks.PUT("/comments/:linkId", h.UpdateContentCitation)
		contentLinks.DELETE("/comments/:linkId", h.RemoveContentCitation)
		
		// Auto-generated topic management
		contentLinks.POST("/templates", h.CreateTopicTemplate)
		contentLinks.GET("/templates/:contentType", h.GetTopicTemplates)
		contentLinks.PUT("/templates/:templateId", h.UpdateTopicTemplate)
		contentLinks.DELETE("/templates/:templateId", h.DeleteTopicTemplate)
		
		// Auto-generate topics
		contentLinks.POST("/generate-topics/:contentType/:contentId", h.GenerateTopicsForContent)
		
		// Discussion recommendations
		contentLinks.POST("/recommendations", h.AddContentRecommendation)
		contentLinks.GET("/content/:contentType/:contentId/recommendations", h.GetContentRecommendations)
		contentLinks.DELETE("/recommendations/:recommendationId", h.RemoveContentRecommendation)
		contentLinks.POST("/content/:contentType/:contentId/generate-recommendations", h.GenerateRecommendationsForContent)
	}
}

// LinkTopicRequest represents a request to link a topic to content
type LinkTopicRequest struct {
	TopicID       uint   `json:"topicId" binding:"required"`
	ContentType   string `json:"contentType" binding:"required"`
	ContentID     uint   `json:"contentId" binding:"required"`
	IsHighlighted bool   `json:"isHighlighted"`
}

// LinkTopicToContent links a topic to content
func (h *ContentLinkHandler) LinkTopicToContent(c *gin.Context) {
	var req LinkTopicRequest
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
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch req.ContentType {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Link topic to content
	link, err := h.contentLinkService.LinkTopicToContent(
		req.TopicID,
		contentType,
		req.ContentID,
		userID.(uint),
		req.IsHighlighted,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, link)
}

// GetTopicContentLinks retrieves content links for a topic
func (h *ContentLinkHandler) GetTopicContentLinks(c *gin.Context) {
	// Get topic ID from path
	topicIDStr := c.Param("topicId")
	topicID, err := strconv.ParseUint(topicIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}
	
	// Get links
	links, err := h.contentLinkService.GetTopicContentLinks(uint(topicID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, links)
}

// GetLinkedTopics retrieves topics linked to a content
func (h *ContentLinkHandler) GetLinkedTopics(c *gin.Context) {
	// Get content type and ID from path
	contentTypeStr := c.Param("contentType")
	contentIDStr := c.Param("contentId")
	
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch contentTypeStr {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Get topics
	links, err := h.contentLinkService.GetLinkedTopics(contentType, uint(contentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, links)
}

// HighlightTopicLinkRequest represents a request to highlight a topic link
type HighlightTopicLinkRequest struct {
	Highlighted bool `json:"highlighted"`
}

// HighlightTopicLink highlights a topic link
func (h *ContentLinkHandler) HighlightTopicLink(c *gin.Context) {
	// Get link ID from path
	linkIDStr := c.Param("linkId")
	linkID, err := strconv.ParseUint(linkIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}
	
	var req HighlightTopicLinkRequest
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
	
	// Highlight link
	if err := h.contentLinkService.HighlightTopicLink(uint(linkID), req.Highlighted, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Topic link updated successfully"})
}

// RemoveTopicLink removes a topic link
func (h *ContentLinkHandler) RemoveTopicLink(c *gin.Context) {
	// Get link ID from path
	linkIDStr := c.Param("linkId")
	linkID, err := strconv.ParseUint(linkIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Remove link
	if err := h.contentLinkService.RemoveTopicLink(uint(linkID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Topic link removed successfully"})
}

// AddContentCitationRequest represents a request to add a content citation
type AddContentCitationRequest struct {
	CommentID       uint   `json:"commentId" binding:"required"`
	ContentType     string `json:"contentType" binding:"required"`
	ContentID       uint   `json:"contentId" binding:"required"`
	CitationText    string `json:"citationText" binding:"required"`
	CitationContext string `json:"citationContext"`
}

// AddContentCitation adds a content citation
func (h *ContentLinkHandler) AddContentCitation(c *gin.Context) {
	var req AddContentCitationRequest
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
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch req.ContentType {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Add citation
	citation, err := h.contentLinkService.AddContentCitation(
		req.CommentID,
		contentType,
		req.ContentID,
		userID.(uint),
		req.CitationText,
		req.CitationContext,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, citation)
}

// GetCommentContentLinks retrieves content links for a comment
func (h *ContentLinkHandler) GetCommentContentLinks(c *gin.Context) {
	// Get comment ID from path
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	
	// Get links
	links, err := h.contentLinkService.GetCommentContentLinks(uint(commentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, links)
}

// GetContentCitations retrieves citations for a content
func (h *ContentLinkHandler) GetContentCitations(c *gin.Context) {
	// Get content type and ID from path
	contentTypeStr := c.Param("contentType")
	contentIDStr := c.Param("contentId")
	
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch contentTypeStr {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Get citations
	citations, err := h.contentLinkService.GetContentCitations(contentType, uint(contentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, citations)
}

// UpdateContentCitationRequest represents a request to update a content citation
type UpdateContentCitationRequest struct {
	CitationText    string `json:"citationText" binding:"required"`
	CitationContext string `json:"citationContext"`
}

// UpdateContentCitation updates a content citation
func (h *ContentLinkHandler) UpdateContentCitation(c *gin.Context) {
	// Get link ID from path
	linkIDStr := c.Param("linkId")
	linkID, err := strconv.ParseUint(linkIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}
	
	var req UpdateContentCitationRequest
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
	
	// Update citation
	if err := h.contentLinkService.UpdateContentCitation(uint(linkID), req.CitationText, req.CitationContext, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Citation updated successfully"})
}

// RemoveContentCitation removes a content citation
func (h *ContentLinkHandler) RemoveContentCitation(c *gin.Context) {
	// Get link ID from path
	linkIDStr := c.Param("linkId")
	linkID, err := strconv.ParseUint(linkIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Remove citation
	if err := h.contentLinkService.RemoveContentCitation(uint(linkID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Citation removed successfully"})
}

// CreateTopicTemplateRequest represents a request to create a topic template
type CreateTopicTemplateRequest struct {
	Name           string `json:"name" binding:"required"`
	ContentType    string `json:"contentType" binding:"required"`
	TitleTemplate  string `json:"titleTemplate" binding:"required"`
	BodyTemplate   string `json:"bodyTemplate" binding:"required"`
}

// CreateTopicTemplate creates a topic template
func (h *ContentLinkHandler) CreateTopicTemplate(c *gin.Context) {
	var req CreateTopicTemplateRequest
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
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch req.ContentType {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Create template
	template, err := h.contentLinkService.CreateTopicTemplate(
		req.Name,
		contentType,
		req.TitleTemplate,
		req.BodyTemplate,
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, template)
}

// GetTopicTemplates retrieves topic templates
func (h *ContentLinkHandler) GetTopicTemplates(c *gin.Context) {
	// Get content type from path
	contentTypeStr := c.Param("contentType")
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch contentTypeStr {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Get templates
	templates, err := h.contentLinkService.GetTopicTemplates(contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, templates)
}

// UpdateTopicTemplateRequest represents a request to update a topic template
type UpdateTopicTemplateRequest struct {
	Name          string `json:"name" binding:"required"`
	TitleTemplate string `json:"titleTemplate" binding:"required"`
	BodyTemplate  string `json:"bodyTemplate" binding:"required"`
	IsActive      bool   `json:"isActive"`
}

// UpdateTopicTemplate updates a topic template
func (h *ContentLinkHandler) UpdateTopicTemplate(c *gin.Context) {
	// Get template ID from path
	templateIDStr := c.Param("templateId")
	templateID, err := strconv.ParseUint(templateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}
	
	var req UpdateTopicTemplateRequest
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
	
	// Update template
	template, err := h.contentLinkService.UpdateTopicTemplate(
		uint(templateID),
		req.Name,
		req.TitleTemplate,
		req.BodyTemplate,
		req.IsActive,
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, template)
}

// DeleteTopicTemplate deletes a topic template
func (h *ContentLinkHandler) DeleteTopicTemplate(c *gin.Context) {
	// Get template ID from path
	templateIDStr := c.Param("templateId")
	templateID, err := strconv.ParseUint(templateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Delete template
	if err := h.contentLinkService.DeleteTopicTemplate(uint(templateID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// GenerateTopicsForContent generates topics for a content
func (h *ContentLinkHandler) GenerateTopicsForContent(c *gin.Context) {
	// Get content type and ID from path
	contentTypeStr := c.Param("contentType")
	contentIDStr := c.Param("contentId")
	
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch contentTypeStr {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Generate topics
	topicIDs, err := h.contentLinkService.GenerateTopicsForContent(contentType, uint(contentID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Topics generated successfully",
		"topicIds": topicIDs,
	})
}

// AddContentRecommendationRequest represents a request to add a content recommendation
type AddContentRecommendationRequest struct {
	ContentType string  `json:"contentType" binding:"required"`
	ContentID   uint    `json:"contentId" binding:"required"`
	TopicID     uint    `json:"topicId" binding:"required"`
	Score       float64 `json:"score" binding:"required"`
}

// AddContentRecommendation adds a content recommendation
func (h *ContentLinkHandler) AddContentRecommendation(c *gin.Context) {
	var req AddContentRecommendationRequest
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
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch req.ContentType {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Add recommendation
	recommendation, err := h.contentLinkService.AddContentRecommendation(
		contentType,
		req.ContentID,
		req.TopicID,
		req.Score,
		userID.(uint),
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, recommendation)
}

// GetContentRecommendations gets recommendations for a content
func (h *ContentLinkHandler) GetContentRecommendations(c *gin.Context) {
	// Get content type and ID from path
	contentTypeStr := c.Param("contentType")
	contentIDStr := c.Param("contentId")
	
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch contentTypeStr {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Get limit from query
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	
	// Get recommendations
	recommendations, err := h.contentLinkService.GetContentRecommendations(contentType, uint(contentID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, recommendations)
}

// RemoveContentRecommendation removes a content recommendation
func (h *ContentLinkHandler) RemoveContentRecommendation(c *gin.Context) {
	// Get recommendation ID from path
	recommendationIDStr := c.Param("recommendationId")
	recommendationID, err := strconv.ParseUint(recommendationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recommendation ID"})
		return
	}
	
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Remove recommendation
	if err := h.contentLinkService.RemoveContentRecommendation(uint(recommendationID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Recommendation removed successfully"})
}

// GenerateRecommendationsForContent generates recommendations for a content
func (h *ContentLinkHandler) GenerateRecommendationsForContent(c *gin.Context) {
	// Get content type and ID from path
	contentTypeStr := c.Param("contentType")
	contentIDStr := c.Param("contentId")
	
	contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}
	
	// Parse content type
	var contentType models.ContentReferenceType
	switch contentTypeStr {
	case "book":
		contentType = models.BookReference
	case "chapter":
		contentType = models.ChapterReference
	case "section":
		contentType = models.SectionReference
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
	
	// Generate recommendations
	count, err := h.contentLinkService.GenerateRecommendationsForContent(contentType, uint(contentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Recommendations generated successfully",
		"count":   count,
	})
}