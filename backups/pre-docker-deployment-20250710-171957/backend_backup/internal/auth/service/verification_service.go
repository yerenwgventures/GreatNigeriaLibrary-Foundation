package service

import (
        "errors"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
        "gorm.io/gorm"
)

// GetVerificationStatus gets a user's verification status
func (s *UserService) GetVerificationStatus(userID uint) (*models.VerificationStatus, error) {
        var status models.VerificationStatus
        result := s.db.Where("user_id = ?", userID).First(&status)
        
        if result.Error != nil {
                // If no status exists, create a new one with default values
                if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                        status = models.VerificationStatus{
                                UserID:           userID,
                                EmailVerified:    false,
                                PhoneVerified:    false,
                                IdentityVerified: false,
                                AddressVerified:  false,
                                TrustLevel:       models.TrustLevelNew,
                                TrustPoints:      0,
                                CreatedAt:        time.Now(),
                                UpdatedAt:        time.Now(),
                        }
                        
                        if err := s.db.Create(&status).Error; err != nil {
                                return nil, err
                        }
                        
                        return &status, nil
                }
                
                return nil, result.Error
        }
        
        return &status, nil
}

// CreateVerificationRequest creates a new verification request
func (s *UserService) CreateVerificationRequest(request models.VerificationRequest) error {
        return s.db.Create(&request).Error
}

// GetVerificationRequests gets all verification requests for a user
func (s *UserService) GetVerificationRequests(userID uint) ([]models.VerificationRequest, error) {
        var requests []models.VerificationRequest
        if err := s.db.Where("user_id = ?", userID).Find(&requests).Error; err != nil {
                return nil, err
        }
        return requests, nil
}

// GetVerificationRequestByID gets a verification request by ID
func (s *UserService) GetVerificationRequestByID(id uint) (*models.VerificationRequest, error) {
        var request models.VerificationRequest
        if err := s.db.First(&request, id).Error; err != nil {
                return nil, err
        }
        return &request, nil
}

// ReviewVerificationRequest updates a verification request with review information
func (s *UserService) ReviewVerificationRequest(id, reviewerID uint, status, notes string) error {
        now := time.Now()
        return s.db.Model(&models.VerificationRequest{}).
                Where("id = ?", id).
                Updates(map[string]interface{}{
                        "status":        status,
                        "reviewed_by_id": reviewerID,
                        "reviewed_at":   now,
                        "review_notes":  notes,
                        "updated_at":    now,
                }).Error
}

// UpdateVerificationField updates a specific verification field
func (s *UserService) UpdateVerificationField(userID uint, field string, value bool) error {
        // First get or create the verification status
        status, err := s.GetVerificationStatus(userID)
        if err != nil {
                return err
        }
        
        // Update the specific field
        updates := map[string]interface{}{
                field:       value,
                "updated_at": time.Now(),
        }
        
        return s.db.Model(&models.VerificationStatus{}).
                Where("user_id = ?", userID).
                Updates(updates).Error
}

// SetUserVerified updates the user's verified status
func (s *UserService) SetUserVerified(userID uint, verified bool) error {
        return s.db.Model(&models.User{}).
                Where("id = ?", userID).
                Update("is_verified", verified).Error
}

// GetUserBadges gets all badges for a user
func (s *UserService) GetUserBadges(userID uint) ([]models.UserBadge, error) {
        var badges []models.UserBadge
        err := s.db.Preload("Badge").Where("user_id = ?", userID).Find(&badges).Error
        return badges, err
}

// AwardBadgeToUser awards a badge to a user
func (s *UserService) AwardBadgeToUser(userID, badgeID, awardedBy uint, reason string) error {
        userBadge := models.UserBadge{
                UserID:    userID,
                BadgeID:   badgeID,
                AwardedAt: time.Now(),
                AwardedBy: awardedBy,
                Reason:    reason,
                IsHidden:  false,
                IsPublic:  true,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
        }
        return s.db.Create(&userBadge).Error
}

// UpdateUserTrustLevel updates a user's trust level
func (s *UserService) UpdateUserTrustLevel(userID uint, level models.TrustLevel) error {
        tx := s.db.Begin()
        
        // Update the user's trust level
        if err := tx.Model(&models.User{}).
                Where("id = ?", userID).
                Update("trust_level", level).Error; err != nil {
                tx.Rollback()
                return err
        }
        
        // Also update the verification status trust level field for consistency
        if err := tx.Model(&models.VerificationStatus{}).
                Where("user_id = ?", userID).
                Updates(map[string]interface{}{
                        "trust_level":     level,
                        "last_promoted_at": time.Now(),
                }).Error; err != nil {
                tx.Rollback()
                return err
        }
        
        // If this is a promotion to a new level, award a badge
        var existingBadge models.UserBadge
        err := tx.Where("user_id = ? AND badge_id = ?", userID, level+10).First(&existingBadge).Error
        
        // If badge doesn't exist and level is valid for a badge, award it
        if errors.Is(err, gorm.ErrRecordNotFound) && level > 0 {
                // Find or create the appropriate trust level badge
                var badge models.Badge
                err = tx.FirstOrCreate(&badge, models.Badge{
                        Name:        getTrustLevelBadgeName(level),
                        Description: getTrustLevelBadgeDescription(level),
                        ImageURL:    getTrustLevelBadgeImage(level),
                        Type:        models.BadgeTypeTrust,
                        IsPublic:    true,
                        IsRare:      level >= models.TrustLevelTrusted, // Higher trust levels are rare
                }).Error
                
                if err != nil {
                        tx.Rollback()
                        return err
                }
                
                // Award the badge
                userBadge := models.UserBadge{
                        UserID:    userID,
                        BadgeID:   badge.ID,
                        AwardedAt: time.Now(),
                        AwardedBy: 0, // System awarded
                        Reason:    "Promoted to " + models.GetTrustLevelNameByID(level),
                        IsHidden:  false,
                        IsPublic:  true,
                        CreatedAt: time.Now(),
                        UpdatedAt: time.Now(),
                }
                
                if err := tx.Create(&userBadge).Error; err != nil {
                        tx.Rollback()
                        return err
                }
        }
        
        return tx.Commit().Error
}

// Helper functions for badge creation

func getTrustLevelBadgeName(level models.TrustLevel) string {
        switch level {
        case models.TrustLevelBasic:
                return "Verified Member"
        case models.TrustLevelMember:
                return "Established Member"
        case models.TrustLevelRegular:
                return "Valued Contributor"
        case models.TrustLevelTrusted:
                return "Trusted Member"
        case models.TrustLevelLeader:
                return "Community Leader"
        default:
                return "Unknown Trust Level"
        }
}

func getTrustLevelBadgeDescription(level models.TrustLevel) string {
        switch level {
        case models.TrustLevelBasic:
                return "This user has verified their email and is a basic member of the community."
        case models.TrustLevelMember:
                return "This user is an established member who regularly participates in the community."
        case models.TrustLevelRegular:
                return "This user is a valued contributor with a history of quality participation."
        case models.TrustLevelTrusted:
                return "This user is highly trusted and helps maintain community standards."
        case models.TrustLevelLeader:
                return "This user is a community leader who drives initiatives and mentors others."
        default:
                return "Badge for an unknown trust level"
        }
}

func getTrustLevelBadgeImage(level models.TrustLevel) string {
        // Return SVG badge images for each trust level
        switch level {
        case models.TrustLevelBasic:
                return "/static/img/badges/trust-level-1.svg"
        case models.TrustLevelMember:
                return "/static/img/badges/trust-level-2.svg"
        case models.TrustLevelRegular:
                return "/static/img/badges/trust-level-3.svg"
        case models.TrustLevelTrusted:
                return "/static/img/badges/trust-level-4.svg"
        case models.TrustLevelLeader:
                return "/static/img/badges/trust-level-5.svg"
        default:
                return "/static/img/badges/trust-level-unknown.svg"
        }
}