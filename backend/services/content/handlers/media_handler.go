package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/service"
)

// MediaHandler handles requests for media generation
type MediaHandler struct {
	mediaGenerator service.MediaGenerator
}

// NewMediaHandler creates a new MediaHandler
func NewMediaHandler(mediaGenerator service.MediaGenerator) *MediaHandler {
	return &MediaHandler{
		mediaGenerator: mediaGenerator,
	}
}

// RegisterRoutes registers the media-related routes
func (h *MediaHandler) RegisterRoutes(router *gin.RouterGroup) {
	media := router.Group("/sections")
	{
		media.POST("/:id/audio", h.GenerateAudio)
		media.POST("/:id/photos", h.GeneratePhotos)
		media.POST("/:id/video", h.GenerateVideo)
		media.POST("/:id/pdf", h.GeneratePDF)
		media.GET("/:id/share", h.GetShareableLink)
	}
}

// GenerateAudio handles POST /api/books/sections/:id/audio
func (h *MediaHandler) GenerateAudio(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	audioURL, duration, err := h.mediaGenerator.GenerateAudioFromText(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate audio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audioUrl":    audioURL,
		"duration":    duration,
		"title":       "Audio Book",
		"generatedAt": time.Now().Format(time.RFC3339),
	})
}

// GeneratePhotos handles POST /api/books/sections/:id/photos
func (h *MediaHandler) GeneratePhotos(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	photoURLs, err := h.mediaGenerator.GeneratePhotoCollection(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate photos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photoUrls":   photoURLs,
		"count":       len(photoURLs),
		"title":       "Photo Book",
		"generatedAt": time.Now().Format(time.RFC3339),
	})
}

// GenerateVideo handles POST /api/books/sections/:id/video
func (h *MediaHandler) GenerateVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	videoURL, duration, err := h.mediaGenerator.GenerateVideoSlideshow(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate video"})
		return
	}

	// Generate a thumbnail URL (in a real implementation, this would be a real thumbnail)
	thumbnailURL := videoURL + ".thumbnail.jpg"

	c.JSON(http.StatusOK, gin.H{
		"videoUrl":     videoURL,
		"duration":     duration,
		"title":        "Video Book",
		"thumbnailUrl": thumbnailURL,
		"generatedAt":  time.Now().Format(time.RFC3339),
	})
}

// GeneratePDF handles POST /api/books/sections/:id/pdf
func (h *MediaHandler) GeneratePDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	pdfURL, pageCount, err := h.mediaGenerator.GeneratePDF(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pdfUrl":      pdfURL,
		"pageCount":   pageCount,
		"title":       "PDF Book",
		"generatedAt": time.Now().Format(time.RFC3339),
	})
}

// GetShareableLink handles GET /api/books/sections/:id/share
func (h *MediaHandler) GetShareableLink(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
		return
	}

	mediaType := c.Query("type")
	if mediaType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Media type is required"})
		return
	}

	shareableLink, mediaURL, err := h.mediaGenerator.GetShareableLink(uint(id), mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate shareable link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shareableLink": shareableLink,
		"mediaUrl":      mediaURL,
	})
}
