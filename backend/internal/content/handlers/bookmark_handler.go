package handlers

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/service"
)

// BookmarkHandler defines handlers for bookmark-related endpoints
type BookmarkHandler struct {
        bookmarkService service.BookmarkService
}

// NewBookmarkHandler creates a new bookmark handler instance
func NewBookmarkHandler(bookmarkService service.BookmarkService) *BookmarkHandler {
        return &BookmarkHandler{
                bookmarkService: bookmarkService,
        }
}

// CreateBookmark handles the POST /books/:id/bookmarks endpoint
func (h *BookmarkHandler) CreateBookmark(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        bookIDStr := c.Param("id")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        // Parse request body
        var bookmarkRequest struct {
                ChapterID   uint   `json:"chapterId"`
                SectionID   uint   `json:"sectionId"`
                Title       string `json:"title"`
                Description string `json:"description"`
                Color       string `json:"color"`
                Position    int    `json:"position"`
        }

        if err := c.ShouldBindJSON(&bookmarkRequest); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        // Create bookmark object
        bookmark := &models.Bookmark{
                UserID:      userID.(uint),
                BookID:      uint(bookID),
                ChapterID:   bookmarkRequest.ChapterID,
                SectionID:   bookmarkRequest.SectionID,
                Title:       bookmarkRequest.Title,
                Description: bookmarkRequest.Description,
                Color:       bookmarkRequest.Color,
                Position:    bookmarkRequest.Position,
        }

        if err := h.bookmarkService.CreateBookmark(bookmark); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark"})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "Bookmark created successfully", "data": bookmark})
}

// GetBookmarks handles the GET /books/:id/bookmarks endpoint
func (h *BookmarkHandler) GetBookmarks(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        bookIDStr := c.Param("id")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        bookmarks, err := h.bookmarkService.GetBookmarks(userID.(uint), uint(bookID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookmarks"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": bookmarks})
}

// UpdateBookmark handles the PUT /books/:id/bookmarks/:bookmarkId endpoint
func (h *BookmarkHandler) UpdateBookmark(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        bookIDStr := c.Param("id")
        bookmarkIDStr := c.Param("bookmarkId")

        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        bookmarkID, err := strconv.ParseUint(bookmarkIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
                return
        }

        // First, get the existing bookmark to verify ownership
        existingBookmark, err := h.bookmarkService.GetBookmarkByID(uint(bookmarkID), userID.(uint))
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found or doesn't belong to user"})
                return
        }

        // Verify the bookmark is for the specified book
        if existingBookmark.BookID != uint(bookID) {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Bookmark does not belong to the specified book"})
                return
        }

        // Parse request body
        var updateRequest struct {
                Title       string `json:"title"`
                Description string `json:"description"`
                Color       string `json:"color"`
        }

        if err := c.ShouldBindJSON(&updateRequest); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        // Update fields
        existingBookmark.Title = updateRequest.Title
        existingBookmark.Description = updateRequest.Description
        existingBookmark.Color = updateRequest.Color

        // Save the update
        if err := h.bookmarkService.UpdateBookmark(existingBookmark); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bookmark"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Bookmark updated successfully", "data": existingBookmark})
}

// DeleteBookmark handles the DELETE /books/:id/bookmarks/:bookmarkId endpoint
func (h *BookmarkHandler) DeleteBookmark(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        bookmarkIDStr := c.Param("bookmarkId")
        bookmarkID, err := strconv.ParseUint(bookmarkIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
                return
        }

        if err := h.bookmarkService.DeleteBookmark(uint(bookmarkID), userID.(uint)); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bookmark"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
}

// ShareBookmark handles the POST /bookmarks/:id/share endpoint
func (h *BookmarkHandler) ShareBookmark(c *gin.Context) {
        // Extract user ID from JWT token
        fromUserID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        bookmarkIDStr := c.Param("id")
        bookmarkID, err := strconv.ParseUint(bookmarkIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
                return
        }

        // Parse request body
        var shareRequest struct {
                ToUserID uint `json:"toUserId"`
        }

        if err := c.ShouldBindJSON(&shareRequest); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        // Share the bookmark
        err = h.bookmarkService.ShareBookmark(
                uint(bookmarkID),
                fromUserID.(uint),
                shareRequest.ToUserID,
        )

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share bookmark"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Bookmark shared successfully"})
}

// RegisterRoutes registers the bookmark-related routes
func (h *BookmarkHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
        // All bookmark routes require authentication
        bookmarks := router.Group("/api/v1")
        bookmarks.Use(authMiddleware)
        {
                // Book-specific bookmarks
                bookmarks.POST("/books/:id/bookmarks", h.CreateBookmark)
                bookmarks.GET("/books/:id/bookmarks", h.GetBookmarks)
                bookmarks.PUT("/books/:id/bookmarks/:bookmarkId", h.UpdateBookmark)
                bookmarks.DELETE("/books/:id/bookmarks/:bookmarkId", h.DeleteBookmark)
                
                // General bookmark management
                bookmarks.POST("/bookmarks/:id/share", h.ShareBookmark)
        }
}