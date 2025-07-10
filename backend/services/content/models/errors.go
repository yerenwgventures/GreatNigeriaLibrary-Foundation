package models

import "errors"

// Error types for the content models
var (
	ErrInvalidElementType = errors.New("invalid interactive element type")
	ErrElementNotFound    = errors.New("interactive element not found")
	ErrInvalidContent     = errors.New("invalid content format")
	ErrPermissionDenied   = errors.New("permission denied to access content")
	ErrContentNotFound    = errors.New("content not found")
)