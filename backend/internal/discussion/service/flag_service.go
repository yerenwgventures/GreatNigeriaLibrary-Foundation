package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/repository"
)

// FlagService defines the interface for content flag operations
type FlagService interface {
	// Content flagging
	FlagContent(contentType string, contentID uint, userID uint, flagType models.ContentFlagType, description string) (*models.ContentFlag, error)
	GetFlagsByContent(contentType string, contentID uint) ([]models.ContentFlag, error)
	GetFlagsByStatus(status models.FlagStatus, page, pageSize int) ([]models.ContentFlag, int64, error)
	GetFlagByID(id uint) (*models.ContentFlag, error)
	ReviewFlag(flagID uint, status models.FlagStatus, actionTaken string, notes string, moderatorID uint) error
	AssignFlag(flagID, moderatorID, assignerID uint) error
	
	// Moderation status
	CreateOrUpdateModerationStatus(contentType string, contentID uint, status models.ContentModerationStatusType, moderatorID uint, reason, notes string) (*models.ContentModerationStatus, error)
	GetModerationStatusByContent(contentType string, contentID uint) (*models.ContentModerationStatus, error)
	MarkUserNotified(contentType string, contentID uint) error
	GetPendingModerationCount() (int64, error)
	
	// User penalties
	ApplyUserPenalty(userID uint, penaltyType models.UserPenaltyType, reason, description string, moderatorID uint, duration *int, relatedContentType *string, relatedContentID *uint, notes string) (*models.UserPenalty, error)
	GetUserPenalties(userID uint) ([]models.UserPenalty, error)
	GetActivePenalties(userID uint) ([]models.UserPenalty, error)
	RemovePenalty(penaltyID uint, moderatorID uint, reason string) error
	GetUserDisciplineHistory(userID uint) ([]models.UserPenalty, error)
	IsUserRestricted(userID uint) (bool, string, error)
}

// FlagServiceImpl implements the FlagService interface
type FlagServiceImpl struct {
	flagRepo repository.FlagRepository
	topicRepo repository.TopicRepository
	commentRepo repository.CommentRepository
}

// NewFlagService creates a new flag service
func NewFlagService(
	flagRepo repository.FlagRepository,
	topicRepo repository.TopicRepository,
	commentRepo repository.CommentRepository,
) FlagService {
	return &FlagServiceImpl{
		flagRepo: flagRepo,
		topicRepo: topicRepo,
		commentRepo: commentRepo,
	}
}

// FlagContent flags a content item
func (s *FlagServiceImpl) FlagContent(
	contentType string,
	contentID uint, 
	userID uint, 
	flagType models.ContentFlagType, 
	description string,
) (*models.ContentFlag, error) {
	// Validate content type
	if contentType != "topic" && contentType != "comment" {
		return nil, errors.New("invalid content type")
	}
	
	// Validate content exists
	if err := s.validateContent(contentType, contentID); err != nil {
		return nil, err
	}
	
	// Create flag
	flag := &models.ContentFlag{
		ContentType:  contentType,
		ContentID:    contentID,
		UserID:       userID,
		FlagType:     flagType,
		Description:  description,
		Status:       models.FlagStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	if err := s.flagRepo.CreateContentFlag(flag); err != nil {
		return nil, fmt.Errorf("error creating content flag: %w", err)
	}
	
	// Check if we need to create/update content moderation status
	if _, err := s.flagRepo.GetModerationStatusByContent(contentType, contentID); err != nil {
		// Create new moderation status if it doesn't exist
		status := &models.ContentModerationStatus{
			ContentType:  contentType,
			ContentID:    contentID,
			Status:       models.ModerationStatusPending,
			Reason:       "Content flagged by user",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		
		if err := s.flagRepo.CreateModerationStatus(status); err != nil {
			// Log error but continue
			fmt.Printf("Error creating moderation status: %v\n", err)
		}
	}
	
	return flag, nil
}

// GetFlagsByContent retrieves flags for a content item
func (s *FlagServiceImpl) GetFlagsByContent(contentType string, contentID uint) ([]models.ContentFlag, error) {
	// Validate content type
	if contentType != "topic" && contentType != "comment" {
		return nil, errors.New("invalid content type")
	}
	
	// Validate content exists
	if err := s.validateContent(contentType, contentID); err != nil {
		return nil, err
	}
	
	return s.flagRepo.GetFlagsByContent(contentType, contentID)
}

// GetFlagsByStatus retrieves flags by status
func (s *FlagServiceImpl) GetFlagsByStatus(status models.FlagStatus, page, pageSize int) ([]models.ContentFlag, int64, error) {
	// Validate status
	if status != "" && 
	   status != models.FlagStatusPending && 
	   status != models.FlagStatusReviewed && 
	   status != models.FlagStatusApproved && 
	   status != models.FlagStatusRejected {
		return nil, 0, errors.New("invalid flag status")
	}
	
	// Ensure valid pagination
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	// Get flags
	flags, err := s.flagRepo.GetFlagsByStatus(status, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting flags by status: %w", err)
	}
	
	// Get total count
	count, err := s.flagRepo.GetFlagCount(status)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting flag count: %w", err)
	}
	
	return flags, count, nil
}

// GetFlagByID retrieves a flag by ID
func (s *FlagServiceImpl) GetFlagByID(id uint) (*models.ContentFlag, error) {
	return s.flagRepo.GetFlagByID(id)
}

// ReviewFlag reviews a content flag
func (s *FlagServiceImpl) ReviewFlag(
	flagID uint, 
	status models.FlagStatus, 
	actionTaken string, 
	notes string, 
	moderatorID uint,
) error {
	// Get flag
	flag, err := s.flagRepo.GetFlagByID(flagID)
	if err != nil {
		return fmt.Errorf("error getting flag: %w", err)
	}
	
	// Validate status
	if status != models.FlagStatusReviewed && 
	   status != models.FlagStatusApproved && 
	   status != models.FlagStatusRejected {
		return errors.New("invalid flag status for review")
	}
	
	// Update flag
	now := time.Now()
	flag.Status = status
	flag.ReviewedBy = &moderatorID
	flag.ReviewedAt = &now
	flag.ActionTaken = actionTaken
	flag.Notes = notes
	flag.UpdatedAt = now
	
	if err := s.flagRepo.UpdateFlag(flag); err != nil {
		return fmt.Errorf("error updating flag: %w", err)
	}
	
	// If flag is approved, update content moderation status
	if status == models.FlagStatusApproved {
		moderationStatus, err := s.flagRepo.GetModerationStatusByContent(flag.ContentType, flag.ContentID)
		if err != nil {
			// Create new moderation status
			moderationStatus = &models.ContentModerationStatus{
				ContentType:  flag.ContentType,
				ContentID:    flag.ContentID,
				Status:       models.ModerationStatusHidden,
				ModeratorID:  &moderatorID,
				Reason:       fmt.Sprintf("Flag approved: %s", flag.FlagType),
				Notes:        notes,
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			
			if err := s.flagRepo.CreateModerationStatus(moderationStatus); err != nil {
				return fmt.Errorf("error creating moderation status: %w", err)
			}
		} else {
			// Update existing moderation status
			moderationStatus.Status = models.ModerationStatusHidden
			moderationStatus.ModeratorID = &moderatorID
			moderationStatus.Reason = fmt.Sprintf("Flag approved: %s", flag.FlagType)
			moderationStatus.Notes = notes
			moderationStatus.UpdatedAt = now
			
			if err := s.flagRepo.UpdateModerationStatus(moderationStatus); err != nil {
				return fmt.Errorf("error updating moderation status: %w", err)
			}
		}
	}
	
	return nil
}

// AssignFlag assigns a flag to a moderator
func (s *FlagServiceImpl) AssignFlag(flagID, moderatorID, assignerID uint) error {
	// Check if flag exists
	if _, err := s.flagRepo.GetFlagByID(flagID); err != nil {
		return fmt.Errorf("error getting flag: %w", err)
	}
	
	return s.flagRepo.AssignFlag(flagID, moderatorID)
}

// CreateOrUpdateModerationStatus creates or updates the moderation status for a content item
func (s *FlagServiceImpl) CreateOrUpdateModerationStatus(
	contentType string, 
	contentID uint, 
	status models.ContentModerationStatusType, 
	moderatorID uint, 
	reason, notes string,
) (*models.ContentModerationStatus, error) {
	// Validate content type
	if contentType != "topic" && contentType != "comment" {
		return nil, errors.New("invalid content type")
	}
	
	// Validate content exists
	if err := s.validateContent(contentType, contentID); err != nil {
		return nil, err
	}
	
	// Check if moderation status exists
	moderationStatus, err := s.flagRepo.GetModerationStatusByContent(contentType, contentID)
	if err != nil {
		// Create new moderation status
		moderationStatus = &models.ContentModerationStatus{
			ContentType:  contentType,
			ContentID:    contentID,
			Status:       status,
			ModeratorID:  &moderatorID,
			Reason:       reason,
			Notes:        notes,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		
		if err := s.flagRepo.CreateModerationStatus(moderationStatus); err != nil {
			return nil, fmt.Errorf("error creating moderation status: %w", err)
		}
	} else {
		// Update existing moderation status
		moderationStatus.Status = status
		moderationStatus.ModeratorID = &moderatorID
		moderationStatus.Reason = reason
		moderationStatus.Notes = notes
		moderationStatus.UpdatedAt = time.Now()
		
		if err := s.flagRepo.UpdateModerationStatus(moderationStatus); err != nil {
			return nil, fmt.Errorf("error updating moderation status: %w", err)
		}
	}
	
	return moderationStatus, nil
}

// GetModerationStatusByContent retrieves the moderation status for a content item
func (s *FlagServiceImpl) GetModerationStatusByContent(contentType string, contentID uint) (*models.ContentModerationStatus, error) {
	// Validate content type
	if contentType != "topic" && contentType != "comment" {
		return nil, errors.New("invalid content type")
	}
	
	return s.flagRepo.GetModerationStatusByContent(contentType, contentID)
}

// MarkUserNotified marks a user as notified about a moderation status change
func (s *FlagServiceImpl) MarkUserNotified(contentType string, contentID uint) error {
	// Validate content type
	if contentType != "topic" && contentType != "comment" {
		return errors.New("invalid content type")
	}
	
	moderationStatus, err := s.flagRepo.GetModerationStatusByContent(contentType, contentID)
	if err != nil {
		return fmt.Errorf("error getting moderation status: %w", err)
	}
	
	moderationStatus.UserNotified = true
	moderationStatus.UpdatedAt = time.Now()
	
	if err := s.flagRepo.UpdateModerationStatus(moderationStatus); err != nil {
		return fmt.Errorf("error updating moderation status: %w", err)
	}
	
	return nil
}

// GetPendingModerationCount gets the count of pending moderation items
func (s *FlagServiceImpl) GetPendingModerationCount() (int64, error) {
	return s.flagRepo.GetPendingModerationCount()
}

// ApplyUserPenalty applies a penalty to a user
func (s *FlagServiceImpl) ApplyUserPenalty(
	userID uint, 
	penaltyType models.UserPenaltyType, 
	reason, description string, 
	moderatorID uint, 
	duration *int, 
	relatedContentType *string, 
	relatedContentID *uint, 
	notes string,
) (*models.UserPenalty, error) {
	// Validate penalty type
	if penaltyType != models.PenaltyTypeWarning && 
	   penaltyType != models.PenaltyTypeSuspension && 
	   penaltyType != models.PenaltyTypeBan && 
	   penaltyType != models.PenaltyTypeRestriction {
		return nil, errors.New("invalid penalty type")
	}
	
	// Validate duration for suspensions
	if penaltyType == models.PenaltyTypeSuspension {
		if duration == nil || *duration <= 0 {
			return nil, errors.New("duration is required for suspensions")
		}
	}
	
	// Create penalty
	now := time.Now()
	penalty := &models.UserPenalty{
		UserID:            userID,
		PenaltyType:       penaltyType,
		Reason:            reason,
		Description:       description,
		ModeratorID:       moderatorID,
		Duration:          duration,
		IsActive:          true,
		RelatedContentType: relatedContentType,
		RelatedContentID:  relatedContentID,
		Notes:             notes,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	
	// Set expiration date if applicable
	if duration != nil && *duration > 0 {
		expiresAt := now.AddDate(0, 0, *duration)
		penalty.ExpiresAt = &expiresAt
	}
	
	if err := s.flagRepo.CreateUserPenalty(penalty); err != nil {
		return nil, fmt.Errorf("error creating user penalty: %w", err)
	}
	
	return penalty, nil
}

// GetUserPenalties retrieves all penalties for a user
func (s *FlagServiceImpl) GetUserPenalties(userID uint) ([]models.UserPenalty, error) {
	return s.flagRepo.GetUserPenalties(userID)
}

// GetActivePenalties retrieves active penalties for a user
func (s *FlagServiceImpl) GetActivePenalties(userID uint) ([]models.UserPenalty, error) {
	return s.flagRepo.GetActivePenalties(userID)
}

// RemovePenalty removes a penalty from a user
func (s *FlagServiceImpl) RemovePenalty(penaltyID uint, moderatorID uint, reason string) error {
	// Get penalty
	penalty, err := s.flagRepo.GetPenaltyByID(penaltyID)
	if err != nil {
		return fmt.Errorf("error getting penalty: %w", err)
	}
	
	// Append removal reason to notes
	if reason != "" {
		penalty.Notes += "\n\nRemoved by moderator ID " + fmt.Sprintf("%d", moderatorID) + " on " + 
			time.Now().Format(time.RFC3339) + " with reason: " + reason
	}
	
	// Update penalty
	penalty.IsActive = false
	penalty.UpdatedAt = time.Now()
	
	if err := s.flagRepo.UpdateUserPenalty(penalty); err != nil {
		return fmt.Errorf("error updating penalty: %w", err)
	}
	
	return nil
}

// GetUserDisciplineHistory retrieves the discipline history for a user
func (s *FlagServiceImpl) GetUserDisciplineHistory(userID uint) ([]models.UserPenalty, error) {
	return s.flagRepo.GetUserPenalties(userID)
}

// IsUserRestricted checks if a user has any active restrictions
func (s *FlagServiceImpl) IsUserRestricted(userID uint) (bool, string, error) {
	penalties, err := s.flagRepo.GetActivePenalties(userID)
	if err != nil {
		return false, "", fmt.Errorf("error getting active penalties: %w", err)
	}
	
	for _, penalty := range penalties {
		if penalty.PenaltyType == models.PenaltyTypeBan {
			return true, "User is permanently banned", nil
		}
		
		if penalty.PenaltyType == models.PenaltyTypeSuspension {
			return true, fmt.Sprintf("User is suspended until %s", penalty.ExpiresAt.Format("2006-01-02")), nil
		}
		
		if penalty.PenaltyType == models.PenaltyTypeRestriction {
			return true, "User has posting restrictions", nil
		}
	}
	
	return false, "", nil
}

// Helper functions

// validateContent validates that content exists
func (s *FlagServiceImpl) validateContent(contentType string, contentID uint) error {
	if contentType == "topic" {
		if _, err := s.topicRepo.GetTopicByID(contentID); err != nil {
			return fmt.Errorf("topic not found: %w", err)
		}
	} else if contentType == "comment" {
		if _, err := s.commentRepo.GetCommentByID(contentID); err != nil {
			return fmt.Errorf("comment not found: %w", err)
		}
	} else {
		return errors.New("invalid content type")
	}
	
	return nil
}