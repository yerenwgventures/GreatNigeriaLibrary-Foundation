package models

import (
	"time"
)

// BookRecommendation represents a recommended book for a user
type BookRecommendation struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index"`
	BookID      uint      `json:"book_id" gorm:"index"`
	Book        Book      `json:"book" gorm:"foreignKey:BookID"`
	Score       float64   `json:"score"`
	ReasonCode  string    `json:"reason_code"`
	Reason      string    `json:"reason"`
	IsRead      bool      `json:"is_read" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RecommendationReason defines possible recommendation reasons
type RecommendationReason string

const (
	ReasonSimilarTopic  RecommendationReason = "similar_topic"
	ReasonPopular       RecommendationReason = "popular"
	ReasonRecentlyAdded RecommendationReason = "recently_added"
	ReasonContinuation  RecommendationReason = "continuation"
	ReasonInterest      RecommendationReason = "user_interest"
)