package handlers

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/service"
)

// InteractiveElementHandler defines handlers for interactive element endpoints
type InteractiveElementHandler struct {
        elementService service.InteractiveElementService
}

// NewInteractiveElementHandler creates a new interactive element handler instance
func NewInteractiveElementHandler(elementService service.InteractiveElementService) *InteractiveElementHandler {
        return &InteractiveElementHandler{
                elementService: elementService,
        }
}

// GetInteractiveElementsBySection handles the GET /interactive/section/:id endpoint
func (h *InteractiveElementHandler) GetInteractiveElementsBySection(c *gin.Context) {
        // Extract section ID from path
        sectionIDStr := c.Param("id")
        sectionID, err := strconv.ParseUint(sectionIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
                return
        }

        // Get the interactive elements
        elements, err := h.elementService.GetInteractiveElementsBySection(uint(sectionID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interactive elements"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": elements})
}

// GetInteractiveElementByID handles the GET /interactive/:id endpoint
func (h *InteractiveElementHandler) GetInteractiveElementByID(c *gin.Context) {
        // Extract element ID from path
        elementIDStr := c.Param("id")
        elementID, err := strconv.ParseUint(elementIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid element ID"})
                return
        }

        // Get the interactive element
        element, err := h.elementService.GetInteractiveElementByID(uint(elementID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interactive element"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": element})
}

// CreateQuiz handles the POST /interactive/quiz endpoint
func (h *InteractiveElementHandler) CreateQuiz(c *gin.Context) {
        // Check for admin role
        role, exists := c.Get("userRole")
        if !exists || role.(string) != "admin" {
                c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
                return
        }

        // Parse request body
        var request struct {
                SectionID      uint                `json:"sectionId" binding:"required"`
                Title          string              `json:"title" binding:"required"`
                Description    string              `json:"description" binding:"required"`
                Content        models.QuizContent  `json:"content" binding:"required"`
                CompletionType string              `json:"completionType" binding:"required"`
                PointsValue    int                 `json:"pointsValue" binding:"required"`
                RequiredStatus bool                `json:"requiredStatus"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Create the quiz
        element, err := h.elementService.CreateQuiz(
                request.SectionID,
                request.Title,
                request.Description,
                &request.Content,
                request.CompletionType,
                request.PointsValue,
                request.RequiredStatus,
        )

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": element})
}

// CreateReflection handles the POST /interactive/reflection endpoint
func (h *InteractiveElementHandler) CreateReflection(c *gin.Context) {
        // Check for admin role
        role, exists := c.Get("userRole")
        if !exists || role.(string) != "admin" {
                c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
                return
        }

        // Parse request body
        var request struct {
                SectionID      uint                      `json:"sectionId" binding:"required"`
                Title          string                    `json:"title" binding:"required"`
                Description    string                    `json:"description" binding:"required"`
                Content        models.ReflectionContent  `json:"content" binding:"required"`
                CompletionType string                    `json:"completionType" binding:"required"`
                PointsValue    int                       `json:"pointsValue" binding:"required"`
                RequiredStatus bool                      `json:"requiredStatus"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Create the reflection
        element, err := h.elementService.CreateReflection(
                request.SectionID,
                request.Title,
                request.Description,
                &request.Content,
                request.CompletionType,
                request.PointsValue,
                request.RequiredStatus,
        )

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reflection: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": element})
}

// CreateCallToAction handles the POST /interactive/call-to-action endpoint
func (h *InteractiveElementHandler) CreateCallToAction(c *gin.Context) {
        // Check for admin role
        role, exists := c.Get("userRole")
        if !exists || role.(string) != "admin" {
                c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
                return
        }

        // Parse request body
        var request struct {
                SectionID      uint                        `json:"sectionId" binding:"required"`
                Title          string                      `json:"title" binding:"required"`
                Description    string                      `json:"description" binding:"required"`
                Content        models.CallToActionContent  `json:"content" binding:"required"`
                CompletionType string                      `json:"completionType" binding:"required"`
                PointsValue    int                         `json:"pointsValue" binding:"required"`
                RequiredStatus bool                        `json:"requiredStatus"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Create the call-to-action
        element, err := h.elementService.CreateCallToAction(
                request.SectionID,
                request.Title,
                request.Description,
                &request.Content,
                request.CompletionType,
                request.PointsValue,
                request.RequiredStatus,
        )

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create call-to-action: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": element})
}

// CreateDiscussionPrompt handles the POST /interactive/discussion-prompt endpoint
func (h *InteractiveElementHandler) CreateDiscussionPrompt(c *gin.Context) {
        // Check for admin role
        role, exists := c.Get("userRole")
        if !exists || role.(string) != "admin" {
                c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
                return
        }

        // Parse request body
        var request struct {
                SectionID      uint                           `json:"sectionId" binding:"required"`
                Title          string                         `json:"title" binding:"required"`
                Description    string                         `json:"description" binding:"required"`
                Content        models.DiscussionPromptContent `json:"content" binding:"required"`
                CompletionType string                         `json:"completionType" binding:"required"`
                PointsValue    int                            `json:"pointsValue" binding:"required"`
                RequiredStatus bool                           `json:"requiredStatus"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Create the discussion prompt
        element, err := h.elementService.CreateDiscussionPrompt(
                request.SectionID,
                request.Title,
                request.Description,
                &request.Content,
                request.CompletionType,
                request.PointsValue,
                request.RequiredStatus,
        )

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create discussion prompt: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": element})
}

// DeleteInteractiveElement handles the DELETE /interactive/:id endpoint
func (h *InteractiveElementHandler) DeleteInteractiveElement(c *gin.Context) {
        // Check for admin role
        role, exists := c.Get("userRole")
        if !exists || role.(string) != "admin" {
                c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
                return
        }

        // Extract element ID from path
        elementIDStr := c.Param("id")
        elementID, err := strconv.ParseUint(elementIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid element ID"})
                return
        }

        // Delete the element
        err = h.elementService.DeleteInteractiveElement(uint(elementID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete interactive element"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Interactive element deleted successfully"})
}

// SubmitQuizResponse handles the POST /interactive/quiz/:id/submit endpoint
func (h *InteractiveElementHandler) SubmitQuizResponse(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Extract element ID from path
        elementIDStr := c.Param("id")
        elementID, err := strconv.ParseUint(elementIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid element ID"})
                return
        }

        // Parse request body
        var request struct {
                Answers []service.QuizAnswer `json:"answers" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Submit the response
        response, err := h.elementService.SubmitQuizResponse(userID.(uint), uint(elementID), request.Answers)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit quiz response: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": response})
}

// SubmitReflectionResponse handles the POST /interactive/reflection/:id/submit endpoint
func (h *InteractiveElementHandler) SubmitReflectionResponse(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Extract element ID from path
        elementIDStr := c.Param("id")
        elementID, err := strconv.ParseUint(elementIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid element ID"})
                return
        }

        // Parse request body
        var request struct {
                Response string `json:"response" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Submit the response
        response, err := h.elementService.SubmitReflectionResponse(userID.(uint), uint(elementID), request.Response)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit reflection response: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": response})
}

// SubmitCallToActionResponse handles the POST /interactive/call-to-action/:id/submit endpoint
func (h *InteractiveElementHandler) SubmitCallToActionResponse(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Extract element ID from path
        elementIDStr := c.Param("id")
        elementID, err := strconv.ParseUint(elementIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid element ID"})
                return
        }

        // Parse request body
        var request struct {
                ActionType string `json:"actionType" binding:"required"`
                ActionData string `json:"actionData" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Submit the response
        response, err := h.elementService.SubmitCallToActionResponse(userID.(uint), uint(elementID), request.ActionType, request.ActionData)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit call-to-action response: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": response})
}

// SubmitDiscussionResponse handles the POST /interactive/discussion-prompt/:id/submit endpoint
func (h *InteractiveElementHandler) SubmitDiscussionResponse(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Extract element ID from path
        elementIDStr := c.Param("id")
        elementID, err := strconv.ParseUint(elementIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid element ID"})
                return
        }

        // Parse request body
        var request struct {
                Response string `json:"response" binding:"required"`
                TopicID  uint   `json:"topicId" binding:"required"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
                return
        }

        // Submit the response
        response, err := h.elementService.SubmitDiscussionResponse(userID.(uint), uint(elementID), request.Response, request.TopicID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit discussion response: " + err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"data": response})
}

// GetUserResponsesForElement handles the GET /interactive/:id/responses endpoint
func (h *InteractiveElementHandler) GetUserResponsesForElement(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Extract element ID from path
        elementIDStr := c.Param("id")
        elementID, err := strconv.ParseUint(elementIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid element ID"})
                return
        }

        // Get the responses
        responses, err := h.elementService.GetUserResponsesForElement(userID.(uint), uint(elementID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responses"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetUserProgress handles the GET /interactive/progress/:bookId endpoint
func (h *InteractiveElementHandler) GetUserProgress(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Extract book ID from path
        bookIDStr := c.Param("bookId")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        // Get the progress
        progress, err := h.elementService.GetUserProgress(userID.(uint), uint(bookID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": progress})
}

// RegisterRoutes registers the interactive element routes
func (h *InteractiveElementHandler) RegisterRoutes(router *gin.Engine, authMiddleware, adminMiddleware gin.HandlerFunc) {
        // Basic interactive element endpoints
        interactive := router.Group("/api/v1/interactive")
        {
                // Public endpoints for retrieving interactive elements
                interactive.GET("/section/:id", h.GetInteractiveElementsBySection)
                interactive.GET("/:id", h.GetInteractiveElementByID)
                
                // Admin-only endpoints for creating/updating/deleting interactive elements
                admin := interactive.Group("/")
                admin.Use(authMiddleware, adminMiddleware)
                {
                        admin.POST("/quiz", h.CreateQuiz)
                        admin.POST("/reflection", h.CreateReflection)
                        admin.POST("/call-to-action", h.CreateCallToAction)
                        admin.POST("/discussion-prompt", h.CreateDiscussionPrompt)
                        admin.DELETE("/:id", h.DeleteInteractiveElement)
                }
                
                // Authenticated endpoints for submitting responses
                auth := interactive.Group("/")
                auth.Use(authMiddleware)
                {
                        auth.POST("/quiz/:id/submit", h.SubmitQuizResponse)
                        auth.POST("/reflection/:id/submit", h.SubmitReflectionResponse)
                        auth.POST("/call-to-action/:id/submit", h.SubmitCallToActionResponse)
                        auth.POST("/discussion-prompt/:id/submit", h.SubmitDiscussionResponse)
                        auth.GET("/:id/responses", h.GetUserResponsesForElement)
                        auth.GET("/progress/:bookId", h.GetUserProgress)
                }
        }
}