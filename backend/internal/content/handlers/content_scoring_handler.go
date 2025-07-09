package handlers

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/service"
)

// ContentScoringHandler handles content scoring endpoints
type ContentScoringHandler struct {
        scoringService service.ContentScoringService
}

// NewContentScoringHandler creates a new content scoring handler
func NewContentScoringHandler(scoringService service.ContentScoringService) *ContentScoringHandler {
        return &ContentScoringHandler{
                scoringService: scoringService,
        }
}

// RegisterRoutes registers the routes for content scoring
func (h *ContentScoringHandler) RegisterRoutes(router *gin.Engine) {
        contentScoring := router.Group("/api/v1/content-scoring")
        {
                // General scoring endpoints
                contentScoring.POST("/score", h.ScoreContent)
                contentScoring.GET("/scores/:contentType/:contentID", h.GetContentScores)
                contentScoring.GET("/scores/:contentType/:contentID/:scoreType", h.GetLatestScore)
                
                // Score criteria endpoints
                contentScoring.POST("/criteria", h.CreateScoreCriteria)
                contentScoring.GET("/criteria/:scoreType", h.GetScoreCriteria)
                contentScoring.PUT("/criteria/:id", h.UpdateScoreCriteria)
                contentScoring.DELETE("/criteria/:id", h.DeleteScoreCriteria)
                
                // Quality metrics endpoints
                contentScoring.POST("/quality-metrics", h.UpdateQualityMetrics)
                contentScoring.GET("/quality-metrics/:contentType/:contentID", h.GetQualityMetrics)
                
                // Relevance metrics endpoints
                contentScoring.POST("/relevance-metrics", h.UpdateRelevanceMetrics)
                contentScoring.GET("/relevance-metrics/:contentType/:contentID", h.GetRelevanceMetrics)
                
                // Safety metrics endpoints
                contentScoring.POST("/safety-metrics", h.UpdateSafetyMetrics)
                contentScoring.GET("/safety-metrics/:contentType/:contentID", h.GetSafetyMetrics)
                contentScoring.GET("/flagged-content", h.GetFlaggedContent)
                
                // Automated analysis endpoints
                contentScoring.POST("/automated-analysis", h.StoreAutomatedAnalysis)
                contentScoring.GET("/automated-analysis/:contentType/:contentID", h.GetAutomatedAnalysis)
        }
}

// ScoreContentRequest represents a request to score content
type ScoreContentRequest struct {
        ContentType    string                `json:"contentType" binding:"required"`
        ContentID      uint                  `json:"contentId" binding:"required"`
        ScoreType      string                `json:"scoreType" binding:"required"`
        Score          float64               `json:"score" binding:"required"`
        ScoredBy       uint                  `json:"scoredBy" binding:"required"`
        ReviewNotes    string                `json:"reviewNotes"`
        ReviewCategory string                `json:"reviewCategory"`
        CriteriaValues map[uint]float64      `json:"criteriaValues"`
}

// ScoreContent scores a piece of content
func (h *ContentScoringHandler) ScoreContent(c *gin.Context) {
        var req ScoreContentRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch req.ContentType {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Parse score type
        var scoreType models.ContentScoreType
        switch req.ScoreType {
        case "quality":
                scoreType = models.QualityScore
        case "relevance":
                scoreType = models.RelevanceScore
        case "safety":
                scoreType = models.SafetyScore
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score type"})
                return
        }
        
        // Score content
        contentScore, err := h.scoringService.ScoreContent(
                contentType,
                req.ContentID,
                scoreType,
                req.Score,
                req.ScoredBy,
                req.ReviewNotes,
                req.ReviewCategory,
                req.CriteriaValues,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, contentScore)
}

// GetContentScores gets all content scores for a piece of content
func (h *ContentScoringHandler) GetContentScores(c *gin.Context) {
        contentTypeStr := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch contentTypeStr {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Get content scores
        contentScores, err := h.scoringService.GetContentScores(contentType, uint(contentID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, contentScores)
}

// GetLatestScore gets the latest content score of a specific type for a piece of content
func (h *ContentScoringHandler) GetLatestScore(c *gin.Context) {
        contentTypeStr := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        scoreTypeStr := c.Param("scoreType")
        
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch contentTypeStr {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Parse score type
        var scoreType models.ContentScoreType
        switch scoreTypeStr {
        case "quality":
                scoreType = models.QualityScore
        case "relevance":
                scoreType = models.RelevanceScore
        case "safety":
                scoreType = models.SafetyScore
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score type"})
                return
        }
        
        // Get latest content score
        contentScore, err := h.scoringService.GetLatestContentScore(contentType, uint(contentID), scoreType)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, contentScore)
}

// CreateScoreCriteriaRequest represents a request to create a score criteria
type CreateScoreCriteriaRequest struct {
        Name        string  `json:"name" binding:"required"`
        Description string  `json:"description" binding:"required"`
        ScoreType   string  `json:"scoreType" binding:"required"`
        Weight      float64 `json:"weight" binding:"required"`
}

// CreateScoreCriteria creates a new score criteria
func (h *ContentScoringHandler) CreateScoreCriteria(c *gin.Context) {
        var req CreateScoreCriteriaRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Parse score type
        var scoreType models.ContentScoreType
        switch req.ScoreType {
        case "quality":
                scoreType = models.QualityScore
        case "relevance":
                scoreType = models.RelevanceScore
        case "safety":
                scoreType = models.SafetyScore
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score type"})
                return
        }
        
        // Create score criteria
        criteria, err := h.scoringService.CreateScoreCriteria(req.Name, req.Description, scoreType, req.Weight)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusCreated, criteria)
}

// GetScoreCriteria gets all score criteria for a score type
func (h *ContentScoringHandler) GetScoreCriteria(c *gin.Context) {
        scoreTypeStr := c.Param("scoreType")
        
        // Parse score type
        var scoreType models.ContentScoreType
        switch scoreTypeStr {
        case "quality":
                scoreType = models.QualityScore
        case "relevance":
                scoreType = models.RelevanceScore
        case "safety":
                scoreType = models.SafetyScore
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score type"})
                return
        }
        
        // Get score criteria
        criteria, err := h.scoringService.GetScoreCriteria(scoreType)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, criteria)
}

// UpdateScoreCriteriaRequest represents a request to update a score criteria
type UpdateScoreCriteriaRequest struct {
        Name        string  `json:"name" binding:"required"`
        Description string  `json:"description" binding:"required"`
        Weight      float64 `json:"weight" binding:"required"`
}

// UpdateScoreCriteria updates a score criteria
func (h *ContentScoringHandler) UpdateScoreCriteria(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid criteria ID"})
                return
        }
        
        var req UpdateScoreCriteriaRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Update score criteria
        criteria, err := h.scoringService.UpdateScoreCriteria(uint(id), req.Name, req.Description, req.Weight)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, criteria)
}

// DeleteScoreCriteria deletes a score criteria
func (h *ContentScoringHandler) DeleteScoreCriteria(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid criteria ID"})
                return
        }
        
        // Delete score criteria
        if err := h.scoringService.DeleteScoreCriteria(uint(id)); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, gin.H{"message": "Score criteria deleted successfully"})
}

// UpdateQualityMetricsRequest represents a request to update quality metrics
type UpdateQualityMetricsRequest struct {
        ContentType          string  `json:"contentType" binding:"required"`
        ContentID            uint    `json:"contentId" binding:"required"`
        GrammarScore         float64 `json:"grammarScore" binding:"required"`
        ReadabilityScore     float64 `json:"readabilityScore" binding:"required"`
        StructureScore       float64 `json:"structureScore" binding:"required"`
        ClarityScore         float64 `json:"clarityScore" binding:"required"`
        FactualAccuracy      float64 `json:"factualAccuracy" binding:"required"`
        ComprehensivenessScore float64 `json:"comprehensivenessScore" binding:"required"`
}

// UpdateQualityMetrics updates quality metrics for a piece of content
func (h *ContentScoringHandler) UpdateQualityMetrics(c *gin.Context) {
        var req UpdateQualityMetricsRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch req.ContentType {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Update quality metrics
        metrics, err := h.scoringService.UpdateQualityMetrics(
                contentType,
                req.ContentID,
                req.GrammarScore,
                req.ReadabilityScore,
                req.StructureScore,
                req.ClarityScore,
                req.FactualAccuracy,
                req.ComprehensivenessScore,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, metrics)
}

// GetQualityMetrics gets quality metrics for a piece of content
func (h *ContentScoringHandler) GetQualityMetrics(c *gin.Context) {
        contentTypeStr := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch contentTypeStr {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Get quality metrics
        metrics, err := h.scoringService.GetQualityMetrics(contentType, uint(contentID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, metrics)
}

// UpdateRelevanceMetricsRequest represents a request to update relevance metrics
type UpdateRelevanceMetricsRequest struct {
        ContentType        string  `json:"contentType" binding:"required"`
        ContentID          uint    `json:"contentId" binding:"required"`
        TopicRelevance     float64 `json:"topicRelevance" binding:"required"`
        AudienceRelevance  float64 `json:"audienceRelevance" binding:"required"`
        TimelinessScore    float64 `json:"timelinessScore" binding:"required"`
        PracticalityScore  float64 `json:"practicalityScore" binding:"required"`
        ContextRelevance   float64 `json:"contextRelevance" binding:"required"`
}

// UpdateRelevanceMetrics updates relevance metrics for a piece of content
func (h *ContentScoringHandler) UpdateRelevanceMetrics(c *gin.Context) {
        var req UpdateRelevanceMetricsRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch req.ContentType {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Update relevance metrics
        metrics, err := h.scoringService.UpdateRelevanceMetrics(
                contentType,
                req.ContentID,
                req.TopicRelevance,
                req.AudienceRelevance,
                req.TimelinessScore,
                req.PracticalityScore,
                req.ContextRelevance,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, metrics)
}

// GetRelevanceMetrics gets relevance metrics for a piece of content
func (h *ContentScoringHandler) GetRelevanceMetrics(c *gin.Context) {
        contentTypeStr := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch contentTypeStr {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Get relevance metrics
        metrics, err := h.scoringService.GetRelevanceMetrics(contentType, uint(contentID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, metrics)
}

// UpdateSafetyMetricsRequest represents a request to update safety metrics
type UpdateSafetyMetricsRequest struct {
        ContentType          string  `json:"contentType" binding:"required"`
        ContentID            uint    `json:"contentId" binding:"required"`
        LanguageScore        float64 `json:"languageScore" binding:"required"`
        BiasScore            float64 `json:"biasScore" binding:"required"`
        SensitivityScore     float64 `json:"sensitivityScore" binding:"required"`
        LegalComplianceScore float64 `json:"legalComplianceScore" binding:"required"`
        EthicalScore         float64 `json:"ethicalScore" binding:"required"`
        Flagged              bool    `json:"flagged"`
        FlagReason           string  `json:"flagReason"`
}

// UpdateSafetyMetrics updates safety metrics for a piece of content
func (h *ContentScoringHandler) UpdateSafetyMetrics(c *gin.Context) {
        var req UpdateSafetyMetricsRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch req.ContentType {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Update safety metrics
        metrics, err := h.scoringService.UpdateSafetyMetrics(
                contentType,
                req.ContentID,
                req.LanguageScore,
                req.BiasScore,
                req.SensitivityScore,
                req.LegalComplianceScore,
                req.EthicalScore,
                req.Flagged,
                req.FlagReason,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, metrics)
}

// GetSafetyMetrics gets safety metrics for a piece of content
func (h *ContentScoringHandler) GetSafetyMetrics(c *gin.Context) {
        contentTypeStr := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch contentTypeStr {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Get safety metrics
        metrics, err := h.scoringService.GetSafetyMetrics(contentType, uint(contentID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, metrics)
}

// GetFlaggedContent gets all flagged content
func (h *ContentScoringHandler) GetFlaggedContent(c *gin.Context) {
        // Get flagged content
        metrics, err := h.scoringService.GetFlaggedContent()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, metrics)
}

// StoreAutomatedAnalysisRequest represents a request to store automated analysis
type StoreAutomatedAnalysisRequest struct {
        ContentType        string  `json:"contentType" binding:"required"`
        ContentID          uint    `json:"contentId" binding:"required"`
        ReadabilityMetrics string  `json:"readabilityMetrics" binding:"required"`
        LanguageAnalysis   string  `json:"languageAnalysis" binding:"required"`
        TopicAnalysis      string  `json:"topicAnalysis" binding:"required"`
        KeywordsDetected   string  `json:"keywordsDetected" binding:"required"`
        SentimentScore     float64 `json:"sentimentScore" binding:"required"`
        CohesionScore      float64 `json:"cohesionScore" binding:"required"`
}

// StoreAutomatedAnalysis stores automated analysis for a piece of content
func (h *ContentScoringHandler) StoreAutomatedAnalysis(c *gin.Context) {
        var req StoreAutomatedAnalysisRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch req.ContentType {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Store automated analysis
        analysis, err := h.scoringService.StoreAutomatedAnalysis(
                contentType,
                req.ContentID,
                req.ReadabilityMetrics,
                req.LanguageAnalysis,
                req.TopicAnalysis,
                req.KeywordsDetected,
                req.SentimentScore,
                req.CohesionScore,
        )
        
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, analysis)
}

// GetAutomatedAnalysis gets automated analysis for a piece of content
func (h *ContentScoringHandler) GetAutomatedAnalysis(c *gin.Context) {
        contentTypeStr := c.Param("contentType")
        contentIDStr := c.Param("contentID")
        
        contentID, err := strconv.ParseUint(contentIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
                return
        }
        
        // Parse content type
        var contentType models.ContentType
        switch contentTypeStr {
        case "book":
                contentType = models.BookContent
        case "chapter":
                contentType = models.ChapterContent
        case "section":
                contentType = models.SectionContent
        default:
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
                return
        }
        
        // Get automated analysis
        analysis, err := h.scoringService.GetAutomatedAnalysis(contentType, uint(contentID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        
        c.JSON(http.StatusOK, analysis)
}