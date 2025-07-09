package models

import (
        "time"

        "gorm.io/gorm"
)

// ContentScoreType represents the type of content score
type ContentScoreType string

const (
        // QualityScore represents the quality of content
        QualityScore ContentScoreType = "quality"
        
        // RelevanceScore represents how relevant content is
        RelevanceScore ContentScoreType = "relevance"
        
        // SafetyScore represents the safety/appropriateness of content
        SafetyScore ContentScoreType = "safety"
)

// Using ContentType from book.go

const (
        // BookContent represents a book
        BookContent ContentType = "book"
        
        // ChapterContent represents a chapter
        ChapterContent ContentType = "chapter"
        
        // SectionContent represents a section
        SectionContent ContentType = "section"
)

// ContentScore represents a score for a piece of content
type ContentScore struct {
        gorm.Model
        ContentType    ContentType     `json:"contentType" gorm:"index:idx_content_score_type"`
        ContentID      uint            `json:"contentId" gorm:"index:idx_content_score_type"`
        ScoreType      ContentScoreType `json:"scoreType" gorm:"index:idx_content_score_type"`
        Score          float64         `json:"score"` // 0-100 score
        ScoredBy       uint            `json:"scoredBy"`
        ScoredAt       time.Time       `json:"scoredAt"`
        ReviewNotes    string          `json:"reviewNotes" gorm:"type:text"`
        ReviewCategory string          `json:"reviewCategory"` // Category of review (e.g., grammar, factual accuracy)
}

// ContentScoreCriteria represents criteria for scoring content
type ContentScoreCriteria struct {
        gorm.Model
        Name        string           `json:"name"`
        Description string           `json:"description"`
        ScoreType   ContentScoreType `json:"scoreType"`
        Weight      float64          `json:"weight"` // Weight in the overall score (0-1)
}

// ContentScoreValue represents a score for a specific criterion
type ContentScoreValue struct {
        gorm.Model
        ContentScoreID uint    `json:"contentScoreId" gorm:"index"`
        CriteriaID     uint    `json:"criteriaId" gorm:"index"`
        Value          float64 `json:"value"` // 0-100 score
        Notes          string  `json:"notes" gorm:"type:text"`
}

// ContentQualityMetrics represents metrics for quality scoring
type ContentQualityMetrics struct {
        gorm.Model
        ContentType        ContentType `json:"contentType" gorm:"index:idx_content_quality_metrics"`
        ContentID          uint        `json:"contentId" gorm:"index:idx_content_quality_metrics"`
        GrammarScore       float64     `json:"grammarScore"`       // 0-100 score
        ReadabilityScore   float64     `json:"readabilityScore"`   // 0-100 score
        StructureScore     float64     `json:"structureScore"`     // 0-100 score
        ClarityScore       float64     `json:"clarityScore"`       // 0-100 score
        FactualAccuracy    float64     `json:"factualAccuracy"`    // 0-100 score
        ComprehensivenessScore float64 `json:"comprehensivenessScore"` // 0-100 score
        OverallQualityScore    float64 `json:"overallQualityScore"`    // 0-100 score
}

// ContentRelevanceMetrics represents metrics for relevance scoring
type ContentRelevanceMetrics struct {
        gorm.Model
        ContentType        ContentType `json:"contentType" gorm:"index:idx_content_relevance_metrics"`
        ContentID          uint        `json:"contentId" gorm:"index:idx_content_relevance_metrics"`
        TopicRelevance     float64     `json:"topicRelevance"`     // 0-100 score
        AudienceRelevance  float64     `json:"audienceRelevance"`  // 0-100 score
        TimelinessScore    float64     `json:"timelinessScore"`    // 0-100 score
        PracticalityScore  float64     `json:"practicalityScore"`  // 0-100 score
        ContextRelevance   float64     `json:"contextRelevance"`   // 0-100 score
        OverallRelevanceScore float64  `json:"overallRelevanceScore"` // 0-100 score
}

// ContentSafetyMetrics represents metrics for safety/appropriateness scoring
type ContentSafetyMetrics struct {
        gorm.Model
        ContentType        ContentType `json:"contentType" gorm:"index:idx_content_safety_metrics"`
        ContentID          uint        `json:"contentId" gorm:"index:idx_content_safety_metrics"`
        LanguageScore      float64     `json:"languageScore"`      // 0-100 score
        BiasScore          float64     `json:"biasScore"`          // 0-100 score
        SensitivityScore   float64     `json:"sensitivityScore"`   // 0-100 score
        LegalComplianceScore float64   `json:"legalComplianceScore"` // 0-100 score
        EthicalScore       float64     `json:"ethicalScore"`       // 0-100 score
        OverallSafetyScore float64     `json:"overallSafetyScore"` // 0-100 score
        FlaggedContent     bool        `json:"flaggedContent"`     // Whether content is flagged for review
        FlagReason         string      `json:"flagReason" gorm:"type:text"`
}

// ContentAutomatedAnalysis represents automated analysis of content
type ContentAutomatedAnalysis struct {
        gorm.Model
        ContentType        ContentType `json:"contentType" gorm:"index:idx_content_automated_analysis"`
        ContentID          uint        `json:"contentId" gorm:"index:idx_content_automated_analysis"`
        ReadabilityMetrics string      `json:"readabilityMetrics" gorm:"type:text"` // JSON metrics for readability (Flesch-Kincaid, etc.)
        LanguageAnalysis   string      `json:"languageAnalysis" gorm:"type:text"`   // JSON analysis of language
        TopicAnalysis      string      `json:"topicAnalysis" gorm:"type:text"`      // JSON analysis of topics covered
        KeywordsDetected   string      `json:"keywordsDetected" gorm:"type:text"`   // JSON list of keywords detected
        SentimentScore     float64     `json:"sentimentScore"`                     // -1 to 1 sentiment score
        CohesionScore      float64     `json:"cohesionScore"`                      // 0-100 cohesion score
        AnalysisTimestamp  time.Time   `json:"analysisTimestamp"`
}