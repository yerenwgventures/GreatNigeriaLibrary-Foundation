package models

import (
	"time"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

// HealthCheckResponse represents a health check response
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// FileUploadResponse represents a file upload response
type FileUploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query    string `json:"query" form:"query" binding:"required"`
	Category string `json:"category" form:"category"`
	Page     int    `json:"page" form:"page" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// SearchResult represents a search result item
type SearchResult struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Score       float64   `json:"score"`
	CreatedAt   time.Time `json:"created_at"`
}

// SearchResponse represents a search response
type SearchResponse struct {
	Success    bool               `json:"success"`
	Query      string             `json:"query"`
	Results    []SearchResult     `json:"results"`
	Pagination PaginationResponse `json:"pagination"`
	Took       int64              `json:"took_ms"`
}

// NotificationRequest represents a notification request
type NotificationRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
	Data    string `json:"data,omitempty"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Data      string    `json:"data"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// ActivityLog represents an activity log entry
type ActivityLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index"`
	Action      string    `json:"action" gorm:"size:100;not null"`
	Resource    string    `json:"resource" gorm:"size:100"`
	ResourceID  string    `json:"resource_id" gorm:"size:100"`
	IPAddress   string    `json:"ip_address" gorm:"size:45"`
	UserAgent   string    `json:"user_agent" gorm:"size:1000"`
	Metadata    string    `json:"metadata" gorm:"type:json"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"index"`
	Action     string    `json:"action" gorm:"size:100;not null"`
	Table      string    `json:"table" gorm:"size:100;not null"`
	RecordID   string    `json:"record_id" gorm:"size:100;not null"`
	OldValues  string    `json:"old_values" gorm:"type:json"`
	NewValues  string    `json:"new_values" gorm:"type:json"`
	IPAddress  string    `json:"ip_address" gorm:"size:45"`
	UserAgent  string    `json:"user_agent" gorm:"size:1000"`
	CreatedAt  time.Time `json:"created_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// SystemConfig represents system configuration
type SystemConfig struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Key         string    `json:"key" gorm:"size:255;not null;uniqueIndex"`
	Value       string    `json:"value" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:text"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Feature represents a feature flag
type Feature struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255;not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	IsEnabled   bool      `json:"is_enabled" gorm:"default:false"`
	Config      string    `json:"config" gorm:"type:json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RateLimitInfo represents rate limit information
type RateLimitInfo struct {
	Limit     int   `json:"limit"`
	Remaining int   `json:"remaining"`
	Reset     int64 `json:"reset"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse represents validation errors response
type ValidationErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}
