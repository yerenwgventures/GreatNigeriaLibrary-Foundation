package handlers

import (
        "net/http"
        "strconv"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/service"
)

// ReadingAnalyticsHandler defines handlers for reading analytics endpoints
type ReadingAnalyticsHandler struct {
        analyticsService service.ReadingAnalyticsService
}

// NewReadingAnalyticsHandler creates a new reading analytics handler instance
func NewReadingAnalyticsHandler(analyticsService service.ReadingAnalyticsService) *ReadingAnalyticsHandler {
        return &ReadingAnalyticsHandler{
                analyticsService: analyticsService,
        }
}

// StartReadingSession handles the POST /analytics/sessions/start endpoint
func (h *ReadingAnalyticsHandler) StartReadingSession(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Parse request body
        var sessionRequest struct {
                BookID    uint   `json:"bookId" binding:"required"`
                ChapterID uint   `json:"chapterId" binding:"required"`
                SectionID uint   `json:"sectionId" binding:"required"`
                Source    string `json:"source" binding:"required"`
        }

        if err := c.ShouldBindJSON(&sessionRequest); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Start the session
        session, err := h.analyticsService.StartReadingSession(
                userID.(uint),
                sessionRequest.BookID,
                sessionRequest.ChapterID,
                sessionRequest.SectionID,
                sessionRequest.Source,
        )

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start reading session"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Reading session started",
                "data": gin.H{
                        "sessionId":  session.ID,
                        "startTime":  session.StartTime,
                        "bookId":     session.BookID,
                        "chapterId":  session.ChapterID,
                        "sectionId":  session.SectionID,
                },
        })
}

// EndReadingSession handles the POST /analytics/sessions/:id/end endpoint
func (h *ReadingAnalyticsHandler) EndReadingSession(c *gin.Context) {
        // Extract user ID from JWT token
        _, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get session ID from path
        sessionIDStr := c.Param("id")
        sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
                return
        }

        // End the session
        err = h.analyticsService.EndReadingSession(uint(sessionID), time.Now())
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end reading session"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Reading session ended",
                "data": gin.H{
                        "sessionId": sessionID,
                        "endTime":   time.Now(),
                },
        })
}

// GetReadingSessions handles the GET /analytics/sessions endpoint
func (h *ReadingAnalyticsHandler) GetReadingSessions(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Parse query parameters
        bookIDStr := c.DefaultQuery("bookId", "0")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        // Parse date range parameters
        startDateStr := c.DefaultQuery("startDate", "")
        endDateStr := c.DefaultQuery("endDate", "")

        var startDate, endDate time.Time
        if startDateStr != "" {
                startDate, err = time.Parse(time.RFC3339, startDateStr)
                if err != nil {
                        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
                        return
                }
        }

        if endDateStr != "" {
                endDate, err = time.Parse(time.RFC3339, endDateStr)
                if err != nil {
                        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
                        return
                }
        }

        // Get the sessions
        sessions, err := h.analyticsService.GetReadingSessions(userID.(uint), uint(bookID), startDate, endDate)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reading sessions"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": sessions})
}

// GetReadingAnalytics handles the GET /analytics/books/:id endpoint
func (h *ReadingAnalyticsHandler) GetReadingAnalytics(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get book ID from path
        bookIDStr := c.Param("id")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        // Get the analytics
        analytics, err := h.analyticsService.GetReadingAnalytics(userID.(uint), uint(bookID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reading analytics"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": analytics})
}

// GetRecentlyViewedSections handles the GET /analytics/recently-viewed endpoint
func (h *ReadingAnalyticsHandler) GetRecentlyViewedSections(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Parse limit parameter
        limitStr := c.DefaultQuery("limit", "5")
        limit, err := strconv.Atoi(limitStr)
        if err != nil || limit <= 0 {
                limit = 5
        }

        // Get recently viewed sections
        sections, err := h.analyticsService.GetRecentlyViewedSections(userID.(uint), limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recently viewed sections"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": sections})
}

// GetReadingTimeStats handles the GET /analytics/books/:id/stats endpoint
func (h *ReadingAnalyticsHandler) GetReadingTimeStats(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get book ID from path
        bookIDStr := c.Param("id")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        // Get the stats
        stats, err := h.analyticsService.GetReadingTimeStats(userID.(uint), uint(bookID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reading time stats"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": stats})
}

// RegisterRoutes registers the reading analytics routes
func (h *ReadingAnalyticsHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
        // All reading analytics routes require authentication
        analytics := router.Group("/api/v1/analytics")
        analytics.Use(authMiddleware)
        {
                // Reading session management
                analytics.POST("/sessions/start", h.StartReadingSession)
                analytics.POST("/sessions/:id/end", h.EndReadingSession)
                analytics.GET("/sessions", h.GetReadingSessions)
                
                // Book analytics
                analytics.GET("/books/:id", h.GetReadingAnalytics)
                analytics.GET("/books/:id/stats", h.GetReadingTimeStats)
                
                // Recent activity
                analytics.GET("/recently-viewed", h.GetRecentlyViewedSections)
        }
}