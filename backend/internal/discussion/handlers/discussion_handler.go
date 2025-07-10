package handlers

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/service"
)

// DiscussionHandler handles discussion-related API requests
type DiscussionHandler struct {
        discussionService service.DiscussionService
        pointsIntegration *ForumPointsIntegration
}

// NewDiscussionHandler creates a new discussion handler
func NewDiscussionHandler(discussionService service.DiscussionService) *DiscussionHandler {
        return &DiscussionHandler{
                discussionService: discussionService,
        }
}

// WithPointsIntegration adds points integration to the discussion handler
func (h *DiscussionHandler) WithPointsIntegration(pointsIntegration *ForumPointsIntegration) *DiscussionHandler {
        h.pointsIntegration = pointsIntegration
        return h
}

// RegisterRoutes registers all discussion routes
func (h *DiscussionHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
        // Public routes (no authentication required)
        discussions := router.Group("/api/v1/discussions")
        {
                // Categories
                discussions.GET("/categories", h.GetCategories)
                discussions.GET("/categories/:id", h.GetCategoryByID)
                discussions.GET("/categories/slug/:slug", h.GetCategoryBySlug)

                // Topics
                discussions.GET("", h.GetTopics)
                discussions.GET("/:id", h.GetTopicByID)
                discussions.GET("/category/:id", h.GetTopicsByCategory)
                discussions.GET("/user/:id", h.GetTopicsByUser)
                discussions.GET("/book/:bookID", h.GetTopicsByBook)
                discussions.GET("/book/:bookID/chapter/:chapterID", h.GetTopicsByChapter)
                discussions.GET("/book/:bookID/chapter/:chapterID/section/:sectionID", h.GetTopicsBySection)

                // Comments
                discussions.GET("/:id/comments", h.GetCommentsByTopic)
                discussions.GET("/comments/:id", h.GetCommentByID)
                discussions.GET("/comments/:id/replies", h.GetRepliesByComment)

                // Tags
                discussions.GET("/tags", h.GetAllTags)
                discussions.GET("/:id/tags", h.GetTagsByTopic)

                // Reactions
                discussions.GET("/:id/reactions", h.GetTopicReactions)
                discussions.GET("/comments/:id/reactions", h.GetCommentReactions)
        }

        // Authenticated routes
        authDiscussions := router.Group("/api/v1/discussions")
        authDiscussions.Use(authMiddleware)
        {
                // Topics
                authDiscussions.POST("", h.CreateTopic)
                authDiscussions.PATCH("/:id", h.UpdateTopic)
                authDiscussions.DELETE("/:id", h.DeleteTopic)
                authDiscussions.POST("/:id/view", h.ViewTopic)

                // Comments
                authDiscussions.POST("/:id/comments", h.CreateComment)
                authDiscussions.PATCH("/comments/:id", h.UpdateComment)
                authDiscussions.DELETE("/comments/:id", h.DeleteComment)

                // Reactions
                authDiscussions.POST("/:id/reactions", h.AddTopicReaction)
                authDiscussions.DELETE("/:id/reactions/:type", h.RemoveTopicReaction)
                authDiscussions.POST("/comments/:id/reactions", h.AddCommentReaction)
                authDiscussions.DELETE("/comments/:id/reactions/:type", h.RemoveCommentReaction)

                // User stats
                authDiscussions.GET("/stats", h.GetUserStats)
        }

        // Admin routes
        adminDiscussions := router.Group("/api/v1/admin/discussions")
        adminDiscussions.Use(authMiddleware, adminMiddleware)
        {
                // Categories
                adminDiscussions.POST("/categories", h.CreateCategory)
                adminDiscussions.PATCH("/categories/:id", h.UpdateCategory)
                adminDiscussions.DELETE("/categories/:id", h.DeleteCategory)

                // Topic moderation
                adminDiscussions.POST("/topics/:id/pin", h.PinTopic)
                adminDiscussions.POST("/topics/:id/lock", h.LockTopic)

                // Tags
                adminDiscussions.POST("/tags", h.CreateTag)
                adminDiscussions.POST("/topics/:id/tags/:tagId", h.AddTagToTopic)
                adminDiscussions.DELETE("/topics/:id/tags/:tagId", h.RemoveTagFromTopic)
        }
}

// GetCategories handles GET /api/v1/discussions/categories
func (h *DiscussionHandler) GetCategories(c *gin.Context) {
        categories, err := h.discussionService.GetCategories()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GetCategoryByID handles GET /api/v1/discussions/categories/:id
func (h *DiscussionHandler) GetCategoryByID(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
                return
        }

        category, err := h.discussionService.GetCategoryByID(uint(id))
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": category})
}

// GetCategoryBySlug handles GET /api/v1/discussions/categories/slug/:slug
func (h *DiscussionHandler) GetCategoryBySlug(c *gin.Context) {
        slug := c.Param("slug")
        category, err := h.discussionService.GetCategoryBySlug(slug)
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": category})
}

// CreateCategory handles POST /api/v1/admin/discussions/categories
func (h *DiscussionHandler) CreateCategory(c *gin.Context) {
        var request struct {
                Name        string `json:"name" binding:"required"`
                Description string `json:"description"`
                ParentID    *uint  `json:"parentId"`
                SortOrder   int    `json:"sortOrder"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        category, err := h.discussionService.CreateCategory(
                request.Name,
                request.Description,
                request.ParentID,
                request.SortOrder,
        )
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": category})
}

// UpdateCategory handles PATCH /api/v1/admin/discussions/categories/:id
func (h *DiscussionHandler) UpdateCategory(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
                return
        }

        var request struct {
                Name        string `json:"name" binding:"required"`
                Description string `json:"description"`
                ParentID    *uint  `json:"parentId"`
                SortOrder   int    `json:"sortOrder"`
                IsActive    bool   `json:"isActive"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        category, err := h.discussionService.UpdateCategory(
                uint(id),
                request.Name,
                request.Description,
                request.ParentID,
                request.SortOrder,
                request.IsActive,
        )
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory handles DELETE /api/v1/admin/discussions/categories/:id
func (h *DiscussionHandler) DeleteCategory(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
                return
        }

        err = h.discussionService.DeleteCategory(uint(id))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetTopics handles GET /api/v1/discussions
func (h *DiscussionHandler) GetTopics(c *gin.Context) {
        // Parse query parameters
        pageStr := c.DefaultQuery("page", "1")
        pageSizeStr := c.DefaultQuery("pageSize", "10")

        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
                page = 1
        }

        pageSize, err := strconv.Atoi(pageSizeStr)
        if err != nil || pageSize < 1 || pageSize > 100 {
                pageSize = 10
        }

        // Create filters based on query parameters
        filters := make(map[string]interface{})
        if isPinnedStr := c.Query("isPinned"); isPinnedStr != "" {
                isPinned := isPinnedStr == "true"
                filters["is_pinned"] = isPinned
        }

        if isLockedStr := c.Query("isLocked"); isLockedStr != "" {
                isLocked := isLockedStr == "true"
                filters["is_locked"] = isLocked
        }

        topics, total, err := h.discussionService.GetTopics(page, pageSize, filters)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "data": topics,
                "meta": gin.H{
                        "total":    total,
                        "page":     page,
                        "pageSize": pageSize,
                        "pages":    (total + int64(pageSize) - 1) / int64(pageSize),
                },
        })
}

// GetTopicByID handles GET /api/v1/discussions/:id
func (h *DiscussionHandler) GetTopicByID(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        topic, err := h.discussionService.GetTopicByID(uint(id))
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": topic})
}

// GetTopicsByCategory handles GET /api/v1/discussions/category/:id
func (h *DiscussionHandler) GetTopicsByCategory(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
                return
        }

        // Parse pagination parameters
        pageStr := c.DefaultQuery("page", "1")
        pageSizeStr := c.DefaultQuery("pageSize", "10")

        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
                page = 1
        }

        pageSize, err := strconv.Atoi(pageSizeStr)
        if err != nil || pageSize < 1 || pageSize > 100 {
                pageSize = 10
        }

        topics, total, err := h.discussionService.GetTopicsByCategory(uint(id), page, pageSize)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "data": topics,
                "meta": gin.H{
                        "total":    total,
                        "page":     page,
                        "pageSize": pageSize,
                        "pages":    (total + int64(pageSize) - 1) / int64(pageSize),
                },
        })
}

// GetTopicsByUser handles GET /api/v1/discussions/user/:id
func (h *DiscussionHandler) GetTopicsByUser(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
                return
        }

        // Parse pagination parameters
        pageStr := c.DefaultQuery("page", "1")
        pageSizeStr := c.DefaultQuery("pageSize", "10")

        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
                page = 1
        }

        pageSize, err := strconv.Atoi(pageSizeStr)
        if err != nil || pageSize < 1 || pageSize > 100 {
                pageSize = 10
        }

        topics, total, err := h.discussionService.GetTopicsByUser(uint(id), page, pageSize)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "data": topics,
                "meta": gin.H{
                        "total":    total,
                        "page":     page,
                        "pageSize": pageSize,
                        "pages":    (total + int64(pageSize) - 1) / int64(pageSize),
                },
        })
}

// GetTopicsByBook handles GET /api/v1/discussions/book/:bookID
func (h *DiscussionHandler) GetTopicsByBook(c *gin.Context) {
        bookIDStr := c.Param("bookID")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        topics, err := h.discussionService.GetTopicsByBookSection(uint(bookID), 0, 0)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": topics})
}

// GetTopicsByChapter handles GET /api/v1/discussions/book/:bookID/chapter/:chapterID
func (h *DiscussionHandler) GetTopicsByChapter(c *gin.Context) {
        bookIDStr := c.Param("bookID")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        chapterIDStr := c.Param("chapterID")
        chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
                return
        }

        topics, err := h.discussionService.GetTopicsByBookSection(uint(bookID), uint(chapterID), 0)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": topics})
}

// GetTopicsBySection handles GET /api/v1/discussions/book/:bookID/chapter/:chapterID/section/:sectionID
func (h *DiscussionHandler) GetTopicsBySection(c *gin.Context) {
        bookIDStr := c.Param("bookID")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        chapterIDStr := c.Param("chapterID")
        chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
                return
        }

        sectionIDStr := c.Param("sectionID")
        sectionID, err := strconv.ParseUint(sectionIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
                return
        }

        topics, err := h.discussionService.GetTopicsByBookSection(uint(bookID), uint(chapterID), uint(sectionID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": topics})
}

// CreateTopic handles POST /api/v1/discussions
func (h *DiscussionHandler) CreateTopic(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        var request struct {
                Title      string `json:"title" binding:"required"`
                Content    string `json:"content" binding:"required"`
                CategoryID uint   `json:"categoryId" binding:"required"`
                BookID     *uint  `json:"bookId"`
                ChapterID  *uint  `json:"chapterId"`
                SectionID  *uint  `json:"sectionId"`
                TagIDs     []uint `json:"tagIds"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        topic, err := h.discussionService.CreateTopic(
                userID.(uint),
                request.CategoryID,
                request.Title,
                request.Content,
                request.BookID,
                request.ChapterID,
                request.SectionID,
                request.TagIDs,
        )
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create topic"})
                return
        }

        // Award points for creating a topic if points integration is enabled
        if h.pointsIntegration != nil {
                // Get category slug for more detailed points records
                categorySlug := ""
                category, err := h.discussionService.GetCategoryByID(request.CategoryID)
                if err == nil && category != nil {
                        categorySlug = category.Slug
                }
                
                // Determine quality level of the content
                quality := h.pointsIntegration.DetermineContentQuality(request.Content)
                
                // Award points in background - note the method has different name in ForumPointsIntegration
                h.pointsIntegration.AwardPointsForNewTopic(userID.(uint), topic.ID, categorySlug, quality)
        }

        c.JSON(http.StatusCreated, gin.H{"data": topic})
}

// UpdateTopic handles PATCH /api/v1/discussions/:id
func (h *DiscussionHandler) UpdateTopic(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        var request struct {
                Title      string `json:"title" binding:"required"`
                Content    string `json:"content" binding:"required"`
                CategoryID uint   `json:"categoryId" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        topic, err := h.discussionService.UpdateTopic(
                uint(id),
                userID.(uint),
                request.Title,
                request.Content,
                request.CategoryID,
        )
        if err != nil {
                if _, ok := err.(models.DiscussionError); ok {
                        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update topic"})
                }
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": topic})
}

// DeleteTopic handles DELETE /api/v1/discussions/:id
func (h *DiscussionHandler) DeleteTopic(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get role for admin check
        role, _ := c.Get("userRole")
        isAdmin := role == "admin"

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        err = h.discussionService.DeleteTopic(uint(id), userID.(uint), isAdmin)
        if err != nil {
                if _, ok := err.(models.DiscussionError); ok {
                        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete topic"})
                }
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Topic deleted successfully"})
}

// ViewTopic handles POST /api/v1/discussions/:id/view
func (h *DiscussionHandler) ViewTopic(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        err = h.discussionService.ViewTopic(uint(id), userID.(uint))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record view"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "View recorded successfully"})
}

// PinTopic handles POST /api/v1/admin/discussions/topics/:id/pin
func (h *DiscussionHandler) PinTopic(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        var request struct {
                Pinned bool `json:"pinned"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        err = h.discussionService.PinTopic(uint(id), request.Pinned)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update topic"})
                return
        }

        // If this topic is being pinned (featured) and points integration is enabled, award points to the topic creator
        if h.pointsIntegration != nil && request.Pinned {
                // Get the topic to find out who created it and its category
                topic, err := h.discussionService.GetTopicByID(uint(id))
                if err == nil && topic != nil {
                        // Get category slug for more detailed points records
                        categorySlug := ""
                        category, err := h.discussionService.GetCategoryByID(topic.CategoryID)
                        if err == nil && category != nil {
                                categorySlug = category.Slug
                        }
                        
                        // Award featured topic points to the creator
                        h.pointsIntegration.AwardPointsForFeaturedTopic(topic.UserID, uint(id), categorySlug)
                }
        }

        c.JSON(http.StatusOK, gin.H{"message": "Topic updated successfully"})
}

// LockTopic handles POST /api/v1/admin/discussions/topics/:id/lock
func (h *DiscussionHandler) LockTopic(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        var request struct {
                Locked bool `json:"locked"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        err = h.discussionService.LockTopic(uint(id), request.Locked)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update topic"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Topic updated successfully"})
}

// GetCommentsByTopic handles GET /api/v1/discussions/:id/comments
func (h *DiscussionHandler) GetCommentsByTopic(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        // Parse pagination parameters
        pageStr := c.DefaultQuery("page", "1")
        pageSizeStr := c.DefaultQuery("pageSize", "10")

        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
                page = 1
        }

        pageSize, err := strconv.Atoi(pageSizeStr)
        if err != nil || pageSize < 1 || pageSize > 100 {
                pageSize = 10
        }

        comments, total, err := h.discussionService.GetCommentsByTopic(uint(id), page, pageSize)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "data": comments,
                "meta": gin.H{
                        "total":    total,
                        "page":     page,
                        "pageSize": pageSize,
                        "pages":    (total + int64(pageSize) - 1) / int64(pageSize),
                },
        })
}

// GetCommentByID handles GET /api/v1/discussions/comments/:id
func (h *DiscussionHandler) GetCommentByID(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
                return
        }

        comment, err := h.discussionService.GetCommentByID(uint(id))
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": comment})
}

// GetRepliesByComment handles GET /api/v1/discussions/comments/:id/replies
func (h *DiscussionHandler) GetRepliesByComment(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
                return
        }

        replies, err := h.discussionService.GetRepliesByComment(uint(id))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch replies"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": replies})
}

// CreateComment handles POST /api/v1/discussions/:id/comments
func (h *DiscussionHandler) CreateComment(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        topicIDStr := c.Param("id")
        topicID, err := strconv.ParseUint(topicIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        var request struct {
                Content  string `json:"content" binding:"required"`
                ParentID *uint  `json:"parentId"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        comment, err := h.discussionService.CreateComment(
                userID.(uint),
                uint(topicID),
                request.Content,
                request.ParentID,
        )
        if err != nil {
                if _, ok := err.(models.DiscussionError); ok {
                        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
                }
                return
        }

        // Award points for creating a comment/reply if points integration is enabled
        if h.pointsIntegration != nil {
                // Determine quality level of the content
                quality := h.pointsIntegration.DetermineContentQuality(request.Content)
                
                // Check if this is a reply to another comment or a top-level comment
                isReply := request.ParentID != nil
                
                // Award points in background - note the method has different name in ForumPointsIntegration
                h.pointsIntegration.AwardPointsForReply(userID.(uint), uint(topicID), comment.ID, isReply, quality)
        }

        c.JSON(http.StatusCreated, gin.H{"data": comment})
}

// UpdateComment handles PATCH /api/v1/discussions/comments/:id
func (h *DiscussionHandler) UpdateComment(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get role for admin check
        role, _ := c.Get("userRole")
        isAdmin := role == "admin"

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
                return
        }

        var request struct {
                Content string `json:"content" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        comment, err := h.discussionService.UpdateComment(
                uint(id),
                userID.(uint),
                request.Content,
                isAdmin,
        )
        if err != nil {
                if _, ok := err.(models.DiscussionError); ok {
                        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
                }
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": comment})
}

// DeleteComment handles DELETE /api/v1/discussions/comments/:id
func (h *DiscussionHandler) DeleteComment(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get role for admin check
        role, _ := c.Get("userRole")
        isAdmin := role == "admin"

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
                return
        }

        err = h.discussionService.DeleteComment(uint(id), userID.(uint), isAdmin)
        if err != nil {
                if _, ok := err.(models.DiscussionError); ok {
                        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
                }
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// GetAllTags handles GET /api/v1/discussions/tags
func (h *DiscussionHandler) GetAllTags(c *gin.Context) {
        tags, err := h.discussionService.GetAllTags()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": tags})
}

// GetTagsByTopic handles GET /api/v1/discussions/:id/tags
func (h *DiscussionHandler) GetTagsByTopic(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        tags, err := h.discussionService.GetTagsByTopic(uint(id))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": tags})
}

// CreateTag handles POST /api/v1/admin/discussions/tags
func (h *DiscussionHandler) CreateTag(c *gin.Context) {
        var request struct {
                Name     string `json:"name" binding:"required"`
                Color    string `json:"color"`
                IsSystem bool   `json:"isSystem"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        tag, err := h.discussionService.CreateTag(
                request.Name,
                request.Color,
                request.IsSystem,
        )
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": tag})
}

// AddTagToTopic handles POST /api/v1/admin/discussions/topics/:id/tags/:tagId
func (h *DiscussionHandler) AddTagToTopic(c *gin.Context) {
        topicIDStr := c.Param("id")
        topicID, err := strconv.ParseUint(topicIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        tagIDStr := c.Param("tagId")
        tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
                return
        }

        err = h.discussionService.AddTagToTopic(uint(topicID), uint(tagID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag to topic"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Tag added to topic successfully"})
}

// RemoveTagFromTopic handles DELETE /api/v1/admin/discussions/topics/:id/tags/:tagId
func (h *DiscussionHandler) RemoveTagFromTopic(c *gin.Context) {
        topicIDStr := c.Param("id")
        topicID, err := strconv.ParseUint(topicIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        tagIDStr := c.Param("tagId")
        tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
                return
        }

        err = h.discussionService.RemoveTagFromTopic(uint(topicID), uint(tagID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag from topic"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Tag removed from topic successfully"})
}

// AddTopicReaction handles POST /api/v1/discussions/:id/reactions
func (h *DiscussionHandler) AddTopicReaction(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        var request struct {
                ReactionType string `json:"reactionType" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        err = h.discussionService.AddReaction(
                userID.(uint),
                uint(id),
                "topic",
                request.ReactionType,
        )
        if err != nil {
                if _, ok := err.(models.DiscussionError); ok {
                        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
                }
                return
        }

        // If this is an upvote and points integration is enabled, award points to the topic creator
        if h.pointsIntegration != nil && (request.ReactionType == "upvote" || request.ReactionType == "like") {
                // Get the topic to find out who created it
                topic, err := h.discussionService.GetTopicByID(uint(id))
                if err == nil && topic != nil {
                        // Award points to the topic creator for receiving an upvote
                        h.pointsIntegration.AwardPointsForUpvotes(topic.UserID, uint(id), 0, 1)
                }
        }

        c.JSON(http.StatusOK, gin.H{"message": "Reaction added successfully"})
}

// RemoveTopicReaction handles DELETE /api/v1/discussions/:id/reactions/:type
func (h *DiscussionHandler) RemoveTopicReaction(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        reactionType := c.Param("type")

        err = h.discussionService.RemoveReaction(
                userID.(uint),
                uint(id),
                "topic",
                reactionType,
        )
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Reaction removed successfully"})
}

// AddCommentReaction handles POST /api/v1/discussions/comments/:id/reactions
func (h *DiscussionHandler) AddCommentReaction(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
                return
        }

        var request struct {
                ReactionType string `json:"reactionType" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        err = h.discussionService.AddReaction(
                userID.(uint),
                uint(id),
                "comment",
                request.ReactionType,
        )
        if err != nil {
                if _, ok := err.(models.DiscussionError); ok {
                        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
                }
                return
        }

        // If this is an upvote and points integration is enabled, award points to the comment creator
        if h.pointsIntegration != nil && (request.ReactionType == "upvote" || request.ReactionType == "like") {
                // Get the comment to find out who created it
                comment, err := h.discussionService.GetCommentByID(uint(id))
                if err == nil && comment != nil {
                        // Find the topic ID for this comment
                        topicID := comment.TopicID

                        // Award points to the comment creator for receiving an upvote
                        h.pointsIntegration.AwardPointsForUpvotes(comment.UserID, topicID, uint(id), 1)
                }
        }

        c.JSON(http.StatusOK, gin.H{"message": "Reaction added successfully"})
}

// RemoveCommentReaction handles DELETE /api/v1/discussions/comments/:id/reactions/:type
func (h *DiscussionHandler) RemoveCommentReaction(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
                return
        }

        reactionType := c.Param("type")

        err = h.discussionService.RemoveReaction(
                userID.(uint),
                uint(id),
                "comment",
                reactionType,
        )
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Reaction removed successfully"})
}

// GetTopicReactions handles GET /api/v1/discussions/:id/reactions
func (h *DiscussionHandler) GetTopicReactions(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
                return
        }

        reactions, err := h.discussionService.GetReactionsSummary(uint(id), "topic")
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reactions"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": reactions})
}

// GetCommentReactions handles GET /api/v1/discussions/comments/:id/reactions
func (h *DiscussionHandler) GetCommentReactions(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
                return
        }

        reactions, err := h.discussionService.GetReactionsSummary(uint(id), "comment")
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reactions"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": reactions})
}

// GetUserStats handles GET /api/v1/discussions/stats
func (h *DiscussionHandler) GetUserStats(c *gin.Context) {
        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        stats, err := h.discussionService.GetUserDiscussionStats(userID.(uint))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user stats"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": stats})
}