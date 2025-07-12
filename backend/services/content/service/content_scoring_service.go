package service

import (
        "errors"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/repository"
)

// ContentScoringService defines the interface for content scoring operations
type ContentScoringService interface {
        // General scoring
        ScoreContent(contentType models.ContentType, contentID uint, scoreType models.ContentScoreType, 
                score float64, scoredBy uint, reviewNotes, reviewCategory string, criteriaValues map[uint]float64) (*models.ContentScore, error)
        GetContentScores(contentType models.ContentType, contentID uint) ([]models.ContentScore, error)
        GetLatestContentScore(contentType models.ContentType, contentID uint, scoreType models.ContentScoreType) (*models.ContentScore, error)
        
        // Score criteria management
        CreateScoreCriteria(name, description string, scoreType models.ContentScoreType, weight float64) (*models.ContentScoreCriteria, error)
        GetScoreCriteria(scoreType models.ContentScoreType) ([]models.ContentScoreCriteria, error)
        UpdateScoreCriteria(id uint, name, description string, weight float64) (*models.ContentScoreCriteria, error)
        DeleteScoreCriteria(id uint) error
        
        // Quality metrics
        UpdateQualityMetrics(contentType models.ContentType, contentID uint, grammarScore, readabilityScore,
                structureScore, clarityScore, factualAccuracy, comprehensivenessScore float64) (*models.ContentQualityMetrics, error)
        GetQualityMetrics(contentType models.ContentType, contentID uint) (*models.ContentQualityMetrics, error)
        
        // Relevance metrics
        UpdateRelevanceMetrics(contentType models.ContentType, contentID uint, topicRelevance, audienceRelevance,
                timelinessScore, practicalityScore, contextRelevance float64) (*models.ContentRelevanceMetrics, error)
        GetRelevanceMetrics(contentType models.ContentType, contentID uint) (*models.ContentRelevanceMetrics, error)
        
        // Safety metrics
        UpdateSafetyMetrics(contentType models.ContentType, contentID uint, languageScore, biasScore,
                sensitivityScore, legalComplianceScore, ethicalScore float64, flagged bool, flagReason string) (*models.ContentSafetyMetrics, error)
        GetSafetyMetrics(contentType models.ContentType, contentID uint) (*models.ContentSafetyMetrics, error)
        GetFlaggedContent() ([]models.ContentSafetyMetrics, error)
        
        // Automated analysis
        StoreAutomatedAnalysis(contentType models.ContentType, contentID uint, readabilityMetrics, languageAnalysis,
                topicAnalysis, keywordsDetected string, sentimentScore, cohesionScore float64) (*models.ContentAutomatedAnalysis, error)
        GetAutomatedAnalysis(contentType models.ContentType, contentID uint) (*models.ContentAutomatedAnalysis, error)
}

// ContentScoringServiceImpl implements the ContentScoringService interface
type ContentScoringServiceImpl struct {
        scoringRepo repository.ContentScoringRepository
        bookRepo    repository.BookRepository
        chapterRepo repository.ChapterRepository
        sectionRepo repository.SectionRepository
}

// NewContentScoringService creates a new content scoring service
func NewContentScoringService(
        scoringRepo repository.ContentScoringRepository,
        bookRepo repository.BookRepository,
        chapterRepo repository.ChapterRepository,
        sectionRepo repository.SectionRepository,
) ContentScoringService {
        return &ContentScoringServiceImpl{
                scoringRepo: scoringRepo,
                bookRepo:    bookRepo,
                chapterRepo: chapterRepo,
                sectionRepo: sectionRepo,
        }
}

// ScoreContent scores a piece of content
func (s *ContentScoringServiceImpl) ScoreContent(
        contentType models.ContentType,
        contentID uint,
        scoreType models.ContentScoreType,
        score float64,
        scoredBy uint,
        reviewNotes, 
        reviewCategory string,
        criteriaValues map[uint]float64,
) (*models.ContentScore, error) {
        // Validate content exists
        if err := s.validateContent(contentType, contentID); err != nil {
                return nil, err
        }
        
        // Create content score
        contentScore := &models.ContentScore{
                ContentType:    contentType,
                ContentID:      contentID,
                ScoreType:      scoreType,
                Score:          score,
                ScoredBy:       scoredBy,
                ScoredAt:       time.Now(),
                ReviewNotes:    reviewNotes,
                ReviewCategory: reviewCategory,
        }
        
        // Save content score
        if err := s.scoringRepo.CreateContentScore(contentScore); err != nil {
                return nil, err
        }
        
        // Save criteria values if provided
        if criteriaValues != nil {
                for criteriaID, value := range criteriaValues {
                        // Get criteria to validate it exists
                        criteria, err := s.scoringRepo.GetScoreCriteriaByID(criteriaID)
                        if err != nil {
                                // Log error but continue
                                // TODO: Add proper logging
                                continue
                        }
                        
                        // Check if criteria is of the correct type
                        if criteria.ScoreType != scoreType {
                                // Log error but continue
                                // TODO: Add proper logging
                                continue
                        }
                        
                        // Create score value
                        scoreValue := &models.ContentScoreValue{
                                ContentScoreID: contentScore.ID,
                                CriteriaID:     criteriaID,
                                Value:          value,
                        }
                        
                        if err := s.scoringRepo.CreateScoreValue(scoreValue); err != nil {
                                // Log error but continue
                                // TODO: Add proper logging
                        }
                }
        }
        
        // Update appropriate metrics based on score type
        switch scoreType {
        case models.QualityScore:
                s.updateQualityMetricsFromScores(contentType, contentID)
        case models.RelevanceScore:
                s.updateRelevanceMetricsFromScores(contentType, contentID)
        case models.SafetyScore:
                s.updateSafetyMetricsFromScores(contentType, contentID)
        }
        
        return contentScore, nil
}

// GetContentScores retrieves all content scores for a piece of content
func (s *ContentScoringServiceImpl) GetContentScores(contentType models.ContentType, contentID uint) ([]models.ContentScore, error) {
        return s.scoringRepo.GetContentScores(contentType, contentID)
}

// GetLatestContentScore retrieves the latest content score of a specific type for a piece of content
func (s *ContentScoringServiceImpl) GetLatestContentScore(contentType models.ContentType, contentID uint, scoreType models.ContentScoreType) (*models.ContentScore, error) {
        return s.scoringRepo.GetLatestContentScore(contentType, contentID, scoreType)
}

// CreateScoreCriteria creates a new score criteria
func (s *ContentScoringServiceImpl) CreateScoreCriteria(name, description string, scoreType models.ContentScoreType, weight float64) (*models.ContentScoreCriteria, error) {
        criteria := &models.ContentScoreCriteria{
                Name:        name,
                Description: description,
                ScoreType:   scoreType,
                Weight:      weight,
        }
        
        if err := s.scoringRepo.CreateScoreCriteria(criteria); err != nil {
                return nil, err
        }
        
        return criteria, nil
}

// GetScoreCriteria retrieves all score criteria for a score type
func (s *ContentScoringServiceImpl) GetScoreCriteria(scoreType models.ContentScoreType) ([]models.ContentScoreCriteria, error) {
        return s.scoringRepo.GetScoreCriteria(scoreType)
}

// UpdateScoreCriteria updates a score criteria
func (s *ContentScoringServiceImpl) UpdateScoreCriteria(id uint, name, description string, weight float64) (*models.ContentScoreCriteria, error) {
        criteria, err := s.scoringRepo.GetScoreCriteriaByID(id)
        if err != nil {
                return nil, err
        }
        
        criteria.Name = name
        criteria.Description = description
        criteria.Weight = weight
        
        if err := s.scoringRepo.UpdateScoreCriteria(criteria); err != nil {
                return nil, err
        }
        
        return criteria, nil
}

// DeleteScoreCriteria deletes a score criteria
func (s *ContentScoringServiceImpl) DeleteScoreCriteria(id uint) error {
        return s.scoringRepo.DeleteScoreCriteria(id)
}

// UpdateQualityMetrics updates quality metrics for a piece of content
func (s *ContentScoringServiceImpl) UpdateQualityMetrics(
        contentType models.ContentType,
        contentID uint,
        grammarScore, 
        readabilityScore,
        structureScore, 
        clarityScore, 
        factualAccuracy, 
        comprehensivenessScore float64,
) (*models.ContentQualityMetrics, error) {
        // Validate content exists
        if err := s.validateContent(contentType, contentID); err != nil {
                return nil, err
        }
        
        // Calculate overall quality score
        overallScore := (grammarScore + readabilityScore + structureScore + clarityScore + factualAccuracy + comprehensivenessScore) / 6
        
        // Try to get existing metrics
        metrics, err := s.scoringRepo.GetQualityMetrics(contentType, contentID)
        if err != nil {
                // Create new metrics
                metrics = &models.ContentQualityMetrics{
                        ContentType:           contentType,
                        ContentID:             contentID,
                        GrammarScore:          grammarScore,
                        ReadabilityScore:      readabilityScore,
                        StructureScore:        structureScore,
                        ClarityScore:          clarityScore,
                        FactualAccuracy:       factualAccuracy,
                        ComprehensivenessScore: comprehensivenessScore,
                        OverallQualityScore:   overallScore,
                }
                
                if err := s.scoringRepo.CreateQualityMetrics(metrics); err != nil {
                        return nil, err
                }
        } else {
                // Update existing metrics
                metrics.GrammarScore = grammarScore
                metrics.ReadabilityScore = readabilityScore
                metrics.StructureScore = structureScore
                metrics.ClarityScore = clarityScore
                metrics.FactualAccuracy = factualAccuracy
                metrics.ComprehensivenessScore = comprehensivenessScore
                metrics.OverallQualityScore = overallScore
                
                if err := s.scoringRepo.UpdateQualityMetrics(metrics); err != nil {
                        return nil, err
                }
        }
        
        return metrics, nil
}

// GetQualityMetrics retrieves quality metrics for a piece of content
func (s *ContentScoringServiceImpl) GetQualityMetrics(contentType models.ContentType, contentID uint) (*models.ContentQualityMetrics, error) {
        return s.scoringRepo.GetQualityMetrics(contentType, contentID)
}

// UpdateRelevanceMetrics updates relevance metrics for a piece of content
func (s *ContentScoringServiceImpl) UpdateRelevanceMetrics(
        contentType models.ContentType,
        contentID uint,
        topicRelevance, 
        audienceRelevance,
        timelinessScore, 
        practicalityScore, 
        contextRelevance float64,
) (*models.ContentRelevanceMetrics, error) {
        // Validate content exists
        if err := s.validateContent(contentType, contentID); err != nil {
                return nil, err
        }
        
        // Calculate overall relevance score
        overallScore := (topicRelevance + audienceRelevance + timelinessScore + practicalityScore + contextRelevance) / 5
        
        // Try to get existing metrics
        metrics, err := s.scoringRepo.GetRelevanceMetrics(contentType, contentID)
        if err != nil {
                // Create new metrics
                metrics = &models.ContentRelevanceMetrics{
                        ContentType:          contentType,
                        ContentID:            contentID,
                        TopicRelevance:       topicRelevance,
                        AudienceRelevance:    audienceRelevance,
                        TimelinessScore:      timelinessScore,
                        PracticalityScore:    practicalityScore,
                        ContextRelevance:     contextRelevance,
                        OverallRelevanceScore: overallScore,
                }
                
                if err := s.scoringRepo.CreateRelevanceMetrics(metrics); err != nil {
                        return nil, err
                }
        } else {
                // Update existing metrics
                metrics.TopicRelevance = topicRelevance
                metrics.AudienceRelevance = audienceRelevance
                metrics.TimelinessScore = timelinessScore
                metrics.PracticalityScore = practicalityScore
                metrics.ContextRelevance = contextRelevance
                metrics.OverallRelevanceScore = overallScore
                
                if err := s.scoringRepo.UpdateRelevanceMetrics(metrics); err != nil {
                        return nil, err
                }
        }
        
        return metrics, nil
}

// GetRelevanceMetrics retrieves relevance metrics for a piece of content
func (s *ContentScoringServiceImpl) GetRelevanceMetrics(contentType models.ContentType, contentID uint) (*models.ContentRelevanceMetrics, error) {
        return s.scoringRepo.GetRelevanceMetrics(contentType, contentID)
}

// UpdateSafetyMetrics updates safety metrics for a piece of content
func (s *ContentScoringServiceImpl) UpdateSafetyMetrics(
        contentType models.ContentType,
        contentID uint,
        languageScore, 
        biasScore,
        sensitivityScore, 
        legalComplianceScore, 
        ethicalScore float64,
        flagged bool,
        flagReason string,
) (*models.ContentSafetyMetrics, error) {
        // Validate content exists
        if err := s.validateContent(contentType, contentID); err != nil {
                return nil, err
        }
        
        // Calculate overall safety score
        overallScore := (languageScore + biasScore + sensitivityScore + legalComplianceScore + ethicalScore) / 5
        
        // Try to get existing metrics
        metrics, err := s.scoringRepo.GetSafetyMetrics(contentType, contentID)
        if err != nil {
                // Create new metrics
                metrics = &models.ContentSafetyMetrics{
                        ContentType:          contentType,
                        ContentID:            contentID,
                        LanguageScore:        languageScore,
                        BiasScore:            biasScore,
                        SensitivityScore:     sensitivityScore,
                        LegalComplianceScore: legalComplianceScore,
                        EthicalScore:         ethicalScore,
                        OverallSafetyScore:   overallScore,
                        FlaggedContent:       flagged,
                        FlagReason:           flagReason,
                }
                
                if err := s.scoringRepo.CreateSafetyMetrics(metrics); err != nil {
                        return nil, err
                }
        } else {
                // Update existing metrics
                metrics.LanguageScore = languageScore
                metrics.BiasScore = biasScore
                metrics.SensitivityScore = sensitivityScore
                metrics.LegalComplianceScore = legalComplianceScore
                metrics.EthicalScore = ethicalScore
                metrics.OverallSafetyScore = overallScore
                metrics.FlaggedContent = flagged
                metrics.FlagReason = flagReason
                
                if err := s.scoringRepo.UpdateSafetyMetrics(metrics); err != nil {
                        return nil, err
                }
        }
        
        return metrics, nil
}

// GetSafetyMetrics retrieves safety metrics for a piece of content
func (s *ContentScoringServiceImpl) GetSafetyMetrics(contentType models.ContentType, contentID uint) (*models.ContentSafetyMetrics, error) {
        return s.scoringRepo.GetSafetyMetrics(contentType, contentID)
}

// GetFlaggedContent retrieves all flagged content
func (s *ContentScoringServiceImpl) GetFlaggedContent() ([]models.ContentSafetyMetrics, error) {
        return s.scoringRepo.GetFlaggedContent()
}

// StoreAutomatedAnalysis stores automated analysis for a piece of content
func (s *ContentScoringServiceImpl) StoreAutomatedAnalysis(
        contentType models.ContentType,
        contentID uint,
        readabilityMetrics, 
        languageAnalysis,
        topicAnalysis, 
        keywordsDetected string,
        sentimentScore, 
        cohesionScore float64,
) (*models.ContentAutomatedAnalysis, error) {
        // Validate content exists
        if err := s.validateContent(contentType, contentID); err != nil {
                return nil, err
        }
        
        // Try to get existing analysis
        analysis, err := s.scoringRepo.GetAutomatedAnalysis(contentType, contentID)
        if err != nil {
                // Create new analysis
                analysis = &models.ContentAutomatedAnalysis{
                        ContentType:        contentType,
                        ContentID:          contentID,
                        ReadabilityMetrics: readabilityMetrics,
                        LanguageAnalysis:   languageAnalysis,
                        TopicAnalysis:      topicAnalysis,
                        KeywordsDetected:   keywordsDetected,
                        SentimentScore:     sentimentScore,
                        CohesionScore:      cohesionScore,
                        AnalysisTimestamp:  time.Now(),
                }
                
                if err := s.scoringRepo.CreateAutomatedAnalysis(analysis); err != nil {
                        return nil, err
                }
        } else {
                // Update existing analysis
                analysis.ReadabilityMetrics = readabilityMetrics
                analysis.LanguageAnalysis = languageAnalysis
                analysis.TopicAnalysis = topicAnalysis
                analysis.KeywordsDetected = keywordsDetected
                analysis.SentimentScore = sentimentScore
                analysis.CohesionScore = cohesionScore
                analysis.AnalysisTimestamp = time.Now()
                
                if err := s.scoringRepo.UpdateAutomatedAnalysis(analysis); err != nil {
                        return nil, err
                }
        }
        
        return analysis, nil
}

// GetAutomatedAnalysis retrieves automated analysis for a piece of content
func (s *ContentScoringServiceImpl) GetAutomatedAnalysis(contentType models.ContentType, contentID uint) (*models.ContentAutomatedAnalysis, error) {
        return s.scoringRepo.GetAutomatedAnalysis(contentType, contentID)
}

// Helper functions

// validateContent checks if the content exists
func (s *ContentScoringServiceImpl) validateContent(contentType models.ContentType, contentID uint) error {
        switch contentType {
        case models.BookContent:
                _, err := s.bookRepo.GetBookByID(contentID)
                return err
        case models.ChapterContent:
                _, err := s.chapterRepo.GetChapterByID(contentID)
                return err
        case models.SectionContent:
                _, err := s.sectionRepo.GetSectionByID(contentID)
                return err
        default:
                return errors.New("invalid content type")
        }
}

// updateQualityMetricsFromScores updates quality metrics based on recent scores
func (s *ContentScoringServiceImpl) updateQualityMetricsFromScores(contentType models.ContentType, contentID uint) {
        // Get all quality scores
        scores, err := s.scoringRepo.GetContentScoresByType(contentType, contentID, models.QualityScore)
        if err != nil || len(scores) == 0 {
                return
        }
        
        // Get criteria for quality scores
        criteria, err := s.scoringRepo.GetScoreCriteria(models.QualityScore)
        if err != nil || len(criteria) == 0 {
                return
        }
        
        // Get most recent score
        latestScore := scores[0]
        for _, score := range scores {
                if score.ScoredAt.After(latestScore.ScoredAt) {
                        latestScore = score
                }
        }
        
        // Get values for the latest score
        values, err := s.scoringRepo.GetScoreValues(latestScore.ID)
        if err != nil || len(values) == 0 {
                return
        }
        
        // Map criteria values
        criteriaMap := make(map[string]float64)
        for _, criteria := range criteria {
                for _, value := range values {
                        if value.CriteriaID == criteria.ID {
                                criteriaMap[criteria.Name] = value.Value
                                break
                        }
                }
        }
        
        // Extract specific criteria values
        grammarScore := getValueOrDefault(criteriaMap, "Grammar", 0)
        readabilityScore := getValueOrDefault(criteriaMap, "Readability", 0)
        structureScore := getValueOrDefault(criteriaMap, "Structure", 0)
        clarityScore := getValueOrDefault(criteriaMap, "Clarity", 0)
        factualAccuracy := getValueOrDefault(criteriaMap, "Factual Accuracy", 0)
        comprehensivenessScore := getValueOrDefault(criteriaMap, "Comprehensiveness", 0)
        
        // Update metrics
        s.UpdateQualityMetrics(
                contentType,
                contentID,
                grammarScore,
                readabilityScore,
                structureScore,
                clarityScore,
                factualAccuracy,
                comprehensivenessScore,
        )
}

// updateRelevanceMetricsFromScores updates relevance metrics based on recent scores
func (s *ContentScoringServiceImpl) updateRelevanceMetricsFromScores(contentType models.ContentType, contentID uint) {
        // Get all relevance scores
        scores, err := s.scoringRepo.GetContentScoresByType(contentType, contentID, models.RelevanceScore)
        if err != nil || len(scores) == 0 {
                return
        }
        
        // Get criteria for relevance scores
        criteria, err := s.scoringRepo.GetScoreCriteria(models.RelevanceScore)
        if err != nil || len(criteria) == 0 {
                return
        }
        
        // Get most recent score
        latestScore := scores[0]
        for _, score := range scores {
                if score.ScoredAt.After(latestScore.ScoredAt) {
                        latestScore = score
                }
        }
        
        // Get values for the latest score
        values, err := s.scoringRepo.GetScoreValues(latestScore.ID)
        if err != nil || len(values) == 0 {
                return
        }
        
        // Map criteria values
        criteriaMap := make(map[string]float64)
        for _, criteria := range criteria {
                for _, value := range values {
                        if value.CriteriaID == criteria.ID {
                                criteriaMap[criteria.Name] = value.Value
                                break
                        }
                }
        }
        
        // Extract specific criteria values
        topicRelevance := getValueOrDefault(criteriaMap, "Topic Relevance", 0)
        audienceRelevance := getValueOrDefault(criteriaMap, "Audience Relevance", 0)
        timelinessScore := getValueOrDefault(criteriaMap, "Timeliness", 0)
        practicalityScore := getValueOrDefault(criteriaMap, "Practicality", 0)
        contextRelevance := getValueOrDefault(criteriaMap, "Context Relevance", 0)
        
        // Update metrics
        s.UpdateRelevanceMetrics(
                contentType,
                contentID,
                topicRelevance,
                audienceRelevance,
                timelinessScore,
                practicalityScore,
                contextRelevance,
        )
}

// updateSafetyMetricsFromScores updates safety metrics based on recent scores
func (s *ContentScoringServiceImpl) updateSafetyMetricsFromScores(contentType models.ContentType, contentID uint) {
        // Get all safety scores
        scores, err := s.scoringRepo.GetContentScoresByType(contentType, contentID, models.SafetyScore)
        if err != nil || len(scores) == 0 {
                return
        }
        
        // Get criteria for safety scores
        criteria, err := s.scoringRepo.GetScoreCriteria(models.SafetyScore)
        if err != nil || len(criteria) == 0 {
                return
        }
        
        // Get most recent score
        latestScore := scores[0]
        for _, score := range scores {
                if score.ScoredAt.After(latestScore.ScoredAt) {
                        latestScore = score
                }
        }
        
        // Get values for the latest score
        values, err := s.scoringRepo.GetScoreValues(latestScore.ID)
        if err != nil || len(values) == 0 {
                return
        }
        
        // Map criteria values
        criteriaMap := make(map[string]float64)
        for _, criteria := range criteria {
                for _, value := range values {
                        if value.CriteriaID == criteria.ID {
                                criteriaMap[criteria.Name] = value.Value
                                break
                        }
                }
        }
        
        // Extract specific criteria values
        languageScore := getValueOrDefault(criteriaMap, "Language", 0)
        biasScore := getValueOrDefault(criteriaMap, "Bias", 0)
        sensitivityScore := getValueOrDefault(criteriaMap, "Sensitivity", 0)
        legalComplianceScore := getValueOrDefault(criteriaMap, "Legal Compliance", 0)
        ethicalScore := getValueOrDefault(criteriaMap, "Ethical Considerations", 0)
        
        // Determine if content should be flagged
        flagged := false
        flagReason := ""
        
        // Automatic flagging if any score is below threshold
        const safetyThreshold = 70.0
        if languageScore < safetyThreshold {
                flagged = true
                flagReason += "Inappropriate language; "
        }
        if biasScore < safetyThreshold {
                flagged = true
                flagReason += "Biased content; "
        }
        if sensitivityScore < safetyThreshold {
                flagged = true
                flagReason += "Insensitive content; "
        }
        if legalComplianceScore < safetyThreshold {
                flagged = true
                flagReason += "Legal compliance issues; "
        }
        if ethicalScore < safetyThreshold {
                flagged = true
                flagReason += "Ethical concerns; "
        }
        
        // Also flag if reviewer flagged it
        if latestScore.ReviewCategory == "Flagged" {
                flagged = true
                if flagReason == "" {
                        flagReason = latestScore.ReviewNotes
                } else {
                        flagReason += "; " + latestScore.ReviewNotes
                }
        }
        
        // Update metrics
        s.UpdateSafetyMetrics(
                contentType,
                contentID,
                languageScore,
                biasScore,
                sensitivityScore,
                legalComplianceScore,
                ethicalScore,
                flagged,
                flagReason,
        )
}

// getValueOrDefault gets a value from a map with a default fallback
func getValueOrDefault(m map[string]float64, key string, defaultValue float64) float64 {
        if value, ok := m[key]; ok {
                return value
        }
        return defaultValue
}