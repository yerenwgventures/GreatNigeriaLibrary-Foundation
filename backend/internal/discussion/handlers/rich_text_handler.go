package handlers

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/service"
)

// RichTextHandler defines the handler for rich text endpoints
type RichTextHandler struct {
        richTextService service.RichTextService
}

// NewRichTextHandler creates a new rich text handler
func NewRichTextHandler(richTextService service.RichTextService) *RichTextHandler {
        return &RichTextHandler{
                richTextService: richTextService,
        }
}

// RegisterRoutes registers the routes for rich text operations
func (h *RichTextHandler) RegisterRoutes(router *gin.RouterGroup) {
        richText := router.Group("/rich-text")
        {
                // Rich text content operations
                richText.POST("", h.CreateOrUpdateRichText)
                richText.GET("/:contentType/:contentID", h.GetRichTextContent)
                
                // Attachment operations
                richText.POST("/attachments", h.CreateAttachment)
                richText.GET("/attachments/:contentType/:contentID", h.GetAttachments)
                richText.DELETE("/attachments/:id", h.DeleteAttachment)
                
                // Code block operations
                richText.POST("/code-blocks", h.CreateOrUpdateCodeBlock)
                richText.GET("/code-blocks/:contentType/:contentID", h.GetCodeBlocks)
                richText.DELETE("/code-blocks/:id", h.DeleteCodeBlock)
                
                // Quote operations
                richText.POST("/quotes", h.CreateQuote)
                richText.GET("/quotes/:contentType/:contentID", h.GetQuotes)
                richText.DELETE("/quotes/:id", h.DeleteQuote)
                
                // Mention operations
                richText.GET("/mentions", h.GetMentions)
                richText.PUT("/mentions/:id/notify", h.MarkMentionAsNotified)
        }
}

// RichTextRequest represents a request to create or update rich text content
type RichTextRequest struct {
        ContentID   uint   `json:"contentId" binding:"required"`
        ContentType string `json:"contentType" binding:"required"`
        RawContent  string `json:"rawContent" binding:"required"`
        Format      string `json:"format" binding:"required"`
}

// CreateOrUpdateRichText creates or updates rich text content
func (h *RichTextHandler) CreateOrUpdateRichText(c *gin.Context) {
        var req RichTextRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Validate content type
        if req.ContentType != "topic" && req.ContentType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
                return
        }
        
        // Convert format
        var format models.ContentFormat
        switch req.Format {
        case "html":
                format = models.FormatHTML
        case "markdown":
                format = models.FormatMarkdown
        case "plaintext":
                format = models.FormatPlain
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Format must be 'html', 'markdown', or 'plaintext'"})
                return
        }
        
        // Create or update rich text
        richText, err := h.richTextService.CreateOrUpdateRichText(
                req.ContentID,
                req.ContentType,
                req.RawContent,
                format,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, richText)
}

// GetRichTextContent retrieves rich text content
func (h *RichTextHandler) GetRichTextContent(c *gin.Context) {
        contentType := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        // Validate content type
        if contentType != "topic" && contentType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
                return
        }
        
        // Parse content ID
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Get rich text content
        richText, err := h.richTextService.GetRichTextContent(uint(contentID), contentType)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, richText)
}

// AttachmentRequest represents a request to create an attachment
type AttachmentRequest struct {
        ContentID   uint   `json:"contentId" binding:"required"`
        ContentType string `json:"contentType" binding:"required"`
        FileName    string `json:"fileName" binding:"required"`
        FileSize    int64  `json:"fileSize" binding:"required"`
        FileType    string `json:"fileType" binding:"required"`
        StoragePath string `json:"storagePath" binding:"required"`
        URL         string `json:"url" binding:"required"`
        IsImage     bool   `json:"isImage"`
        Width       int    `json:"width"`
        Height      int    `json:"height"`
}

// CreateAttachment creates a new attachment
func (h *RichTextHandler) CreateAttachment(c *gin.Context) {
        var req AttachmentRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Validate content type
        if req.ContentType != "topic" && req.ContentType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
                return
        }
        
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }
        
        // Create attachment
        attachment, err := h.richTextService.CreateAttachment(
                req.ContentID,
                req.ContentType,
                userID.(uint),
                req.FileName,
                req.FileSize,
                req.FileType,
                req.StoragePath,
                req.URL,
                req.IsImage,
                req.Width,
                req.Height,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusCreated, attachment)
}

// GetAttachments retrieves attachments for a content
func (h *RichTextHandler) GetAttachments(c *gin.Context) {
        contentType := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        // Validate content type
        if contentType != "topic" && contentType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
                return
        }
        
        // Parse content ID
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Get attachments
        attachments, err := h.richTextService.GetAttachmentsByContent(uint(contentID), contentType)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, attachments)
}

// DeleteAttachment deletes an attachment
func (h *RichTextHandler) DeleteAttachment(c *gin.Context) {
        idStr := c.Param("id")
        
        // Parse attachment ID
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
                return
        }
        
        // Delete attachment
        if err := h.richTextService.DeleteAttachment(uint(id)); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, gin.H{"message": "Attachment deleted successfully"})
}

// CodeBlockRequest represents a request to create or update a code block
type CodeBlockRequest struct {
        ContentID   uint   `json:"contentId" binding:"required"`
        ContentType string `json:"contentType" binding:"required"`
        Language    string `json:"language" binding:"required"`
        Code        string `json:"code" binding:"required"`
        Position    int    `json:"position" binding:"required"`
}

// CreateOrUpdateCodeBlock creates or updates a code block
func (h *RichTextHandler) CreateOrUpdateCodeBlock(c *gin.Context) {
        var req CodeBlockRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Validate content type
        if req.ContentType != "topic" && req.ContentType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
                return
        }
        
        // Create or update code block
        codeBlock, err := h.richTextService.CreateOrUpdateCodeBlock(
                req.ContentID,
                req.ContentType,
                req.Language,
                req.Code,
                req.Position,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, codeBlock)
}

// GetCodeBlocks retrieves code blocks for a content
func (h *RichTextHandler) GetCodeBlocks(c *gin.Context) {
        contentType := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        // Validate content type
        if contentType != "topic" && contentType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
                return
        }
        
        // Parse content ID
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Get code blocks
        codeBlocks, err := h.richTextService.GetCodeBlocksByContent(uint(contentID), contentType)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, codeBlocks)
}

// DeleteCodeBlock deletes a code block
func (h *RichTextHandler) DeleteCodeBlock(c *gin.Context) {
        idStr := c.Param("id")
        
        // Parse code block ID
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code block ID"})
                return
        }
        
        // Delete code block
        if err := h.richTextService.DeleteCodeBlock(uint(id)); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, gin.H{"message": "Code block deleted successfully"})
}

// QuoteRequest represents a request to create a quote
type QuoteRequest struct {
        ContentID       uint   `json:"contentId" binding:"required"`
        ContentType     string `json:"contentType" binding:"required"`
        QuotedContentID uint   `json:"quotedContentId" binding:"required"`
        QuotedType      string `json:"quotedType" binding:"required"`
        QuotedUserID    uint   `json:"quotedUserId" binding:"required"`
        QuotedUsername  string `json:"quotedUsername" binding:"required"`
        QuotedContent   string `json:"quotedContent" binding:"required"`
        Position        int    `json:"position" binding:"required"`
}

// CreateQuote creates a new quote
func (h *RichTextHandler) CreateQuote(c *gin.Context) {
        var req QuoteRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Validate content types
        if req.ContentType != "topic" && req.ContentType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
                return
        }
        
        if req.QuotedType != "topic" && req.QuotedType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Quoted type must be 'topic' or 'comment'"})
                return
        }
        
        // Create quote
        quote, err := h.richTextService.CreateQuote(
                req.ContentID,
                req.ContentType,
                req.QuotedContentID,
                req.QuotedType,
                req.QuotedUserID,
                req.QuotedUsername,
                req.QuotedContent,
                req.Position,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusCreated, quote)
}

// GetQuotes retrieves quotes for a content
func (h *RichTextHandler) GetQuotes(c *gin.Context) {
        contentType := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        // Validate content type
        if contentType != "topic" && contentType != "comment" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Content type must be 'topic' or 'comment'"})
                return
        }
        
        // Parse content ID
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Get quotes
        quotes, err := h.richTextService.GetQuotesByContent(uint(contentID), contentType)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, quotes)
}

// DeleteQuote deletes a quote
func (h *RichTextHandler) DeleteQuote(c *gin.Context) {
        idStr := c.Param("id")
        
        // Parse quote ID
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
                return
        }
        
        // Delete quote
        if err := h.richTextService.DeleteQuote(uint(id)); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, gin.H{"message": "Quote deleted successfully"})
}

// GetMentions retrieves mentions for the authenticated user
func (h *RichTextHandler) GetMentions(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }
        
        // Get mentions
        mentions, err := h.richTextService.GetMentionsForUser(userID.(uint))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, mentions)
}

// MarkMentionAsNotified marks a mention as notified
func (h *RichTextHandler) MarkMentionAsNotified(c *gin.Context) {
        idStr := c.Param("id")
        
        // Parse mention ID
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mention ID"})
                return
        }
        
        // Mark mention as notified
        if err := h.richTextService.MarkMentionAsNotified(uint(id)); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, gin.H{"message": "Mention marked as notified"})
}