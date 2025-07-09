package service

import (
        "encoding/json"
        "errors"
        "fmt"
        "regexp"
        "strings"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/repository"
)

// ModerationService defines the interface for moderation operations
type ModerationService interface {
        // Content moderation rules
        CreateModerationRule(name, description, pattern, patternType string, action models.ModerationAction, severity int, appliesTo string, userID uint) (*models.ContentModerationRule, error)
        GetModerationRules(activeOnly bool) ([]models.ContentModerationRule, error)
        UpdateModerationRule(ruleID uint, name, description, pattern, patternType string, action models.ModerationAction, severity int, isActive bool, appliesTo string, userID uint) (*models.ContentModerationRule, error)
        DeleteModerationRule(ruleID uint, userID uint) error
        
        // Content filtering
        FilterContent(content, contentType string, userID uint) (*models.ContentFilterResult, error)
        GetFilterResults(contentType string, contentID uint) ([]models.ContentFilterResult, error)
        ReviewFilterResult(resultID uint, action models.ModerationAction, userID uint) error
        GetUserFilterResults(userID uint, page, pageSize int) ([]models.ContentFilterResult, error)
        
        // Moderation queue
        AddToModerationQueue(contentType string, contentID, userID uint, reason string, filterResultID *uint, priority int) (*models.ModAdvancedQueue, error)
        GetModerationQueue(status string, page, pageSize int) ([]models.ModAdvancedQueue, int64, error)
        AssignModerationItem(itemID, moderatorID, assignerID uint) error
        ResolveModerationItem(itemID uint, decision string, notes string, userID uint) error
        GetModerationQueueStats() (map[string]int64, error)
        
        // User trust management
        GetUserTrustScore(userID uint) (*models.UserTrustScore, error)
        UpdateUserTrustScore(userID uint, contentScore, communityScore, moderatorScore float64, moderatorID uint) (*models.UserTrustScore, error)
        RecalculateUserTrustLevel(userID uint) (models.UserTrustLevel, error)
        GetUsersByTrustLevel(level models.UserTrustLevel) ([]models.UserTrustScore, error)
        
        // Moderator management
        GrantModeratorPrivileges(userID, grantedByID uint, privileges map[string]bool) (*models.ModeratorPrivilege, error)
        UpdateModeratorPrivileges(userID, updatedByID uint, privileges map[string]bool) (*models.ModeratorPrivilege, error)
        RevokeModeratorPrivileges(userID, revokedByID uint) error
        GetModeratorPrivileges(userID uint) (*models.ModeratorPrivilege, error)
        GetAllModerators() ([]models.ModeratorPrivilege, error)
        IsUserModerator(userID uint) (bool, error)
        
        // User moderation actions
        CreateUserModerationAction(userID uint, actionType models.ModerationAction, reason string, duration int, relatedContentID *uint, relatedContentType *string, moderatorID uint, notes string) (*models.UserModerationAction, error)
        GetUserModerationActions(userID uint) ([]models.UserModerationAction, error)
        RevokeUserModerationAction(actionID, moderatorID uint, reason string) error
        GetActiveUserActions(userID uint) ([]models.UserModerationAction, error)
        IsUserBanned(userID uint) (bool, error)
        
        // Prohibited words
        AddProhibitedWord(word, replacement string, isRegex, isAutoReplace bool, severity int, userID uint) (*models.ProhibitedWord, error)
        GetProhibitedWords(activeOnly bool) ([]models.ProhibitedWord, error)
        UpdateProhibitedWord(wordID uint, word, replacement string, isRegex, isAutoReplace, isActive bool, severity int, userID uint) (*models.ProhibitedWord, error)
        DeleteProhibitedWord(wordID, userID uint) error
        FilterTextWithProhibitedWords(text string) (string, bool, []string, error)
}

// ModerationServiceImpl implements the ModerationService interface
type ModerationServiceImpl struct {
        moderationRepo repository.ModerationRepository
        userRepo       repository.UserRepository
}

// NewModerationService creates a new moderation service
func NewModerationService(
        moderationRepo repository.ModerationRepository,
        userRepo repository.UserRepository,
) ModerationService {
        return &ModerationServiceImpl{
                moderationRepo: moderationRepo,
                userRepo:       userRepo,
        }
}

// CreateModerationRule creates a new content moderation rule
func (s *ModerationServiceImpl) CreateModerationRule(
        name, description, pattern, patternType string, 
        action models.ModerationAction, 
        severity int, 
        appliesTo string, 
        userID uint,
) (*models.ContentModerationRule, error) {
        // Validate inputs
        if name == "" || pattern == "" {
                return nil, errors.New("name and pattern are required")
        }
        
        if patternType != "regex" && patternType != "keywords" {
                return nil, errors.New("pattern type must be 'regex' or 'keywords'")
        }
        
        // Validate pattern
        if patternType == "regex" {
                _, err := regexp.Compile(pattern)
                if err != nil {
                        return nil, fmt.Errorf("invalid regex pattern: %w", err)
                }
        }
        
        // Validate severity
        if severity < 1 || severity > 10 {
                return nil, errors.New("severity must be between 1 and 10")
        }
        
        // Validate content type
        if appliesTo != "topic" && appliesTo != "comment" && appliesTo != "username" {
                return nil, errors.New("applies to must be 'topic', 'comment', or 'username'")
        }
        
        // Create rule
        rule := &models.ContentModerationRule{
                Name:          name,
                Description:   description,
                Pattern:       pattern,
                PatternType:   patternType,
                Action:        string(action), // Convert ModerationAction to string
                RuleType:      "content", // Default rule type
                Severity:      severity,
                IsActive:      true,
                AppliesTo:     appliesTo,
                CreatedBy:     userID,
                LastUpdatedBy: userID,
        }
        
        if err := s.moderationRepo.CreateModerationRule(rule); err != nil {
                return nil, fmt.Errorf("error creating moderation rule: %w", err)
        }
        
        return rule, nil
}

// GetModerationRules retrieves all content moderation rules
func (s *ModerationServiceImpl) GetModerationRules(activeOnly bool) ([]models.ContentModerationRule, error) {
        return s.moderationRepo.GetModerationRules(activeOnly)
}

// UpdateModerationRule updates a content moderation rule
func (s *ModerationServiceImpl) UpdateModerationRule(
        ruleID uint, 
        name, description, pattern, patternType string, 
        action models.ModerationAction, 
        severity int, 
        isActive bool, 
        appliesTo string, 
        userID uint,
) (*models.ContentModerationRule, error) {
        // Get existing rule
        rule, err := s.moderationRepo.GetModerationRuleByID(ruleID)
        if err != nil {
                return nil, fmt.Errorf("error getting moderation rule: %w", err)
        }
        
        // Validate inputs
        if name == "" || pattern == "" {
                return nil, errors.New("name and pattern are required")
        }
        
        if patternType != "regex" && patternType != "keywords" {
                return nil, errors.New("pattern type must be 'regex' or 'keywords'")
        }
        
        // Validate pattern
        if patternType == "regex" {
                _, err := regexp.Compile(pattern)
                if err != nil {
                        return nil, fmt.Errorf("invalid regex pattern: %w", err)
                }
        }
        
        // Validate severity
        if severity < 1 || severity > 10 {
                return nil, errors.New("severity must be between 1 and 10")
        }
        
        // Validate content type
        if appliesTo != "topic" && appliesTo != "comment" && appliesTo != "username" {
                return nil, errors.New("applies to must be 'topic', 'comment', or 'username'")
        }
        
        // Update rule
        rule.Name = name
        rule.Description = description
        rule.Pattern = pattern
        rule.PatternType = patternType
        rule.Action = string(action) // Convert ModerationAction to string
        rule.Severity = severity
        rule.IsActive = isActive
        rule.AppliesTo = appliesTo
        rule.LastUpdatedBy = userID
        
        if err := s.moderationRepo.UpdateModerationRule(rule); err != nil {
                return nil, fmt.Errorf("error updating moderation rule: %w", err)
        }
        
        return rule, nil
}

// DeleteModerationRule deletes a content moderation rule
func (s *ModerationServiceImpl) DeleteModerationRule(ruleID, userID uint) error {
        // Check if rule exists
        if _, err := s.moderationRepo.GetModerationRuleByID(ruleID); err != nil {
                return fmt.Errorf("error getting moderation rule: %w", err)
        }
        
        return s.moderationRepo.DeleteModerationRule(ruleID)
}

// FilterContent filters content based on moderation rules
func (s *ModerationServiceImpl) FilterContent(content, contentType string, userID uint) (*models.ContentFilterResult, error) {
        // Get active rules for this content type
        rules, err := s.moderationRepo.GetActiveRulesByType(contentType)
        if err != nil {
                return nil, fmt.Errorf("error getting moderation rules: %w", err)
        }
        
        // Check content against rules
        triggeredRules := []uint{}
        highestSeverity := 0
        var highestSeverityAction models.ModerationAction = models.ActionNone
        
        for _, rule := range rules {
                matched := false
                
                if rule.PatternType == "regex" {
                        re, err := regexp.Compile(rule.Pattern)
                        if err != nil {
                                // Log error but continue
                                fmt.Printf("Invalid regex pattern in rule %d: %v\n", rule.ID, err)
                                continue
                        }
                        matched = re.MatchString(content)
                } else {
                        // Keywords match
                        keywords := strings.Split(rule.Pattern, ",")
                        for _, keyword := range keywords {
                                keyword = strings.TrimSpace(keyword)
                                if strings.Contains(strings.ToLower(content), strings.ToLower(keyword)) {
                                        matched = true
                                        break
                                }
                        }
                }
                
                if matched {
                        triggeredRules = append(triggeredRules, rule.ID)
                        
                        // Check for highest severity rule
                        if rule.Severity > highestSeverity {
                                highestSeverity = rule.Severity
                                highestSeverityAction = models.ModerationAction(rule.Action)
                        } else if rule.Severity == highestSeverity {
                                // If severity is the same, prioritize more restrictive actions
                                if s.isMoreRestrictiveAction(models.ModerationAction(rule.Action), highestSeverityAction) {
                                        highestSeverityAction = models.ModerationAction(rule.Action)
                                }
                        }
                }
        }
        
        // If no rules triggered, return nil
        if len(triggeredRules) == 0 {
                return nil, nil
        }
        
        // Replace prohibited words if needed
        cleanedContent, _, _, err := s.FilterTextWithProhibitedWords(content)
        if err != nil {
                // Log error but continue with original content
                fmt.Printf("Error filtering prohibited words: %v\n", err)
                cleanedContent = content
        }
        
        // Create filter result
        rulesJSON, err := json.Marshal(triggeredRules)
        if err != nil {
                return nil, fmt.Errorf("error marshaling triggered rules: %w", err)
        }
        
        result := &models.ContentFilterResult{
                ContentType:           contentType,
                ContentID:             0, // Will be set after content is created
                UserID:                userID,
                TriggeredRules:        string(rulesJSON),
                Action:                highestSeverityAction,
                FilteredContent:       content,
                CleanedContent:        cleanedContent,
                AutomaticallyProcessed: true,
                CreatedAt:             time.Now(),
        }
        
        return result, nil
}

// GetFilterResults retrieves filter results for content
func (s *ModerationServiceImpl) GetFilterResults(contentType string, contentID uint) ([]models.ContentFilterResult, error) {
        return s.moderationRepo.GetFilterResultsByContent(contentType, contentID)
}

// ReviewFilterResult reviews a filter result
func (s *ModerationServiceImpl) ReviewFilterResult(resultID uint, action models.ModerationAction, userID uint) error {
        // Get filter result
        result, err := s.moderationRepo.GetFilterResultByID(resultID)
        if err != nil {
                return fmt.Errorf("error getting filter result: %w", err)
        }
        
        // Update result
        now := time.Now()
        result.Action = action
        result.AutomaticallyProcessed = false
        result.ModeratorID = &userID
        result.ReviewedAt = &now
        
        return s.moderationRepo.UpdateFilterResult(result)
}

// GetUserFilterResults retrieves filter results for a user
func (s *ModerationServiceImpl) GetUserFilterResults(userID uint, page, pageSize int) ([]models.ContentFilterResult, error) {
        // Ensure valid pagination
        if page < 1 {
                page = 1
        }
        
        if pageSize < 1 || pageSize > 100 {
                pageSize = 20
        }
        
        offset := (page - 1) * pageSize
        
        return s.moderationRepo.GetFilterResultsByUser(userID, pageSize, offset)
}

// AddToModerationQueue adds an item to the moderation queue
func (s *ModerationServiceImpl) AddToModerationQueue(
        contentType string, 
        contentID, userID uint, 
        reason string, 
        filterResultID *uint, 
        priority int,
) (*models.ModAdvancedQueue, error) {
        // Validate priority
        if priority < 1 || priority > 5 {
                priority = 3 // Default to medium priority
        }
        
        // Create queue item
        item := &models.ModAdvancedQueue{
                ContentType:    contentType,
                ContentID:      contentID,
                UserID:         userID,
                Reason:         reason,
                FilterResultID: filterResultID,
                Status:         "pending",
                Priority:       priority,
                CreatedAt:      time.Now(),
                UpdatedAt:      time.Now(),
        }
        
        if err := s.moderationRepo.AddToModerationQueue(item); err != nil {
                return nil, fmt.Errorf("error adding to moderation queue: %w", err)
        }
        
        return item, nil
}

// GetModerationQueue retrieves items from the moderation queue
func (s *ModerationServiceImpl) GetModerationQueue(status string, page, pageSize int) ([]models.ModAdvancedQueue, int64, error) {
        // Ensure valid pagination
        if page < 1 {
                page = 1
        }
        
        if pageSize < 1 || pageSize > 100 {
                pageSize = 20
        }
        
        offset := (page - 1) * pageSize
        
        // Get items
        items, err := s.moderationRepo.GetModerationQueueItems(status, pageSize, offset)
        if err != nil {
                return nil, 0, fmt.Errorf("error getting moderation queue items: %w", err)
        }
        
        // Get total count
        count, err := s.moderationRepo.GetModerationQueueCountByStatus(status)
        if err != nil {
                return nil, 0, fmt.Errorf("error getting moderation queue count: %w", err)
        }
        
        return items, count, nil
}

// AssignModerationItem assigns a moderation queue item to a moderator
func (s *ModerationServiceImpl) AssignModerationItem(itemID, moderatorID, assignerID uint) error {
        // Check if the item exists
        item, err := s.moderationRepo.GetModerationQueueItemByID(itemID)
        if err != nil {
                return fmt.Errorf("error getting moderation queue item: %w", err)
        }
        
        // Check if the moderator has privileges
        isModerator, err := s.moderationRepo.IsUserModerator(moderatorID)
        if err != nil {
                return fmt.Errorf("error checking moderator privileges: %w", err)
        }
        
        if !isModerator {
                return errors.New("user is not a moderator")
        }
        
        // Assign item
        if err := s.moderationRepo.AssignModerationQueueItem(itemID, moderatorID); err != nil {
                return fmt.Errorf("error assigning moderation queue item: %w", err)
        }
        
        // Update status to in-review if it's pending
        if item.Status == "pending" {
                item.Status = "in_review"
                item.AssignedTo = &moderatorID
                item.UpdatedAt = time.Now()
                
                if err := s.moderationRepo.UpdateModerationQueueItem(item); err != nil {
                        return fmt.Errorf("error updating moderation queue item status: %w", err)
                }
        }
        
        return nil
}

// ResolveModerationItem resolves a moderation queue item
func (s *ModerationServiceImpl) ResolveModerationItem(itemID uint, decision, notes string, userID uint) error {
        // Check if the item exists
        item, err := s.moderationRepo.GetModerationQueueItemByID(itemID)
        if err != nil {
                return fmt.Errorf("error getting moderation queue item: %w", err)
        }
        
        // Validate decision
        if decision != "approved" && decision != "rejected" {
                return errors.New("decision must be 'approved' or 'rejected'")
        }
        
        // Check if the user has moderator privileges
        isModerator, err := s.moderationRepo.IsUserModerator(userID)
        if err != nil {
                return fmt.Errorf("error checking moderator privileges: %w", err)
        }
        
        if !isModerator {
                return errors.New("user is not a moderator")
        }
        
        // Update item
        now := time.Now()
        item.Status = decision
        item.Decision = decision
        item.Notes = notes
        item.ReviewedBy = &userID
        item.ReviewedAt = &now
        item.UpdatedAt = now
        
        if err := s.moderationRepo.UpdateModerationQueueItem(item); err != nil {
                return fmt.Errorf("error updating moderation queue item: %w", err)
        }
        
        // If there's a filter result ID, update it too
        if item.FilterResultID != nil {
                result, err := s.moderationRepo.GetFilterResultByID(*item.FilterResultID)
                if err != nil {
                        // Log error but continue
                        fmt.Printf("Error getting filter result: %v\n", err)
                } else {
                        result.ModeratorID = &userID
                        result.ReviewedAt = &now
                        result.AutomaticallyProcessed = false
                        
                        // Set action based on decision
                        if decision == "approved" {
                                result.Action = models.ActionApprove
                        } else {
                                result.Action = models.ActionReject
                        }
                        
                        if err := s.moderationRepo.UpdateFilterResult(result); err != nil {
                                // Log error but continue
                                fmt.Printf("Error updating filter result: %v\n", err)
                        }
                }
        }
        
        // Update user trust score
        if decision == "rejected" {
                trustScore, err := s.moderationRepo.GetUserTrustScore(item.UserID)
                if err != nil {
                        // Log error but continue
                        fmt.Printf("Error getting user trust score: %v\n", err)
                } else {
                        trustScore.ContentRejections++
                        trustScore.UpdatedAt = now
                        
                        if err := s.moderationRepo.UpdateUserTrustScore(trustScore); err != nil {
                                // Log error but continue
                                fmt.Printf("Error updating user trust score: %v\n", err)
                        }
                        
                        // Recalculate trust level
                        if _, err := s.RecalculateUserTrustLevel(item.UserID); err != nil {
                                // Log error but continue
                                fmt.Printf("Error recalculating user trust level: %v\n", err)
                        }
                }
        }
        
        return nil
}

// GetModerationQueueStats gets statistics for the moderation queue
func (s *ModerationServiceImpl) GetModerationQueueStats() (map[string]int64, error) {
        stats := make(map[string]int64)
        
        // Get total count
        total, err := s.moderationRepo.GetModerationQueueCountByStatus("")
        if err != nil {
                return nil, fmt.Errorf("error getting total moderation queue count: %w", err)
        }
        stats["total"] = total
        
        // Get count by status
        statuses := []string{"pending", "in_review", "approved", "rejected"}
        for _, status := range statuses {
                count, err := s.moderationRepo.GetModerationQueueCountByStatus(status)
                if err != nil {
                        // Log error but continue
                        fmt.Printf("Error getting moderation queue count for status %s: %v\n", status, err)
                        continue
                }
                stats[status] = count
        }
        
        return stats, nil
}

// GetUserTrustScore retrieves a user's trust score
func (s *ModerationServiceImpl) GetUserTrustScore(userID uint) (*models.UserTrustScore, error) {
        return s.moderationRepo.GetUserTrustScore(userID)
}

// UpdateUserTrustScore updates a user's trust score
func (s *ModerationServiceImpl) UpdateUserTrustScore(
        userID uint, 
        contentScore, communityScore, moderatorScore float64, 
        moderatorID uint,
) (*models.UserTrustScore, error) {
        // Get current score
        score, err := s.moderationRepo.GetUserTrustScore(userID)
        if err != nil {
                return nil, fmt.Errorf("error getting user trust score: %w", err)
        }
        
        // Update scores
        score.ContentScore = contentScore
        score.CommunityScore = communityScore
        score.ModeratorScore = moderatorScore
        score.LastScoreUpdate = time.Now()
        score.UpdatedAt = time.Now()
        
        // Calculate overall trust score
        // Simple weighted average: 40% content, 40% community, 20% moderator
        score.TrustScore = (0.4 * contentScore) + (0.4 * communityScore) + (0.2 * moderatorScore)
        
        // Update the score
        if err := s.moderationRepo.UpdateUserTrustScore(score); err != nil {
                return nil, fmt.Errorf("error updating user trust score: %w", err)
        }
        
        // Recalculate trust level
        newLevel, err := s.RecalculateUserTrustLevel(userID)
        if err != nil {
                // Log error but continue
                fmt.Printf("Error recalculating user trust level: %v\n", err)
        } else {
                score.TrustLevel = newLevel
                if err := s.moderationRepo.UpdateUserTrustScore(score); err != nil {
                        // Log error but continue
                        fmt.Printf("Error updating user trust level: %v\n", err)
                }
        }
        
        return score, nil
}

// RecalculateUserTrustLevel recalculates a user's trust level based on their trust score
func (s *ModerationServiceImpl) RecalculateUserTrustLevel(userID uint) (models.UserTrustLevel, error) {
        // Get current score
        score, err := s.moderationRepo.GetUserTrustScore(userID)
        if err != nil {
                return "", fmt.Errorf("error getting user trust score: %w", err)
        }
        
        // Calculate new level based on score and other factors
        var newLevel models.UserTrustLevel
        
        // Thresholds and penalties
        const (
                BasicThreshold   = 20.0
                MemberThreshold  = 50.0
                RegularThreshold = 75.0
                LeaderThreshold  = 90.0
                
                ReportPenalty      = -5.0
                WarningPenalty     = -10.0
                RejectionPenalty   = -2.0
        )
        
        // Apply penalties
        adjustedScore := score.TrustScore
        adjustedScore += float64(score.ReportCount) * ReportPenalty
        adjustedScore += float64(score.WarningCount) * WarningPenalty
        adjustedScore += float64(score.ContentRejections) * RejectionPenalty
        
        // Determine level
        if adjustedScore >= LeaderThreshold && score.ContentRejections == 0 && score.WarningCount == 0 {
                newLevel = models.TrustLevelLeader
        } else if adjustedScore >= RegularThreshold {
                newLevel = models.TrustLevelRegular
        } else if adjustedScore >= MemberThreshold {
                newLevel = models.TrustLevelMember
        } else if adjustedScore >= BasicThreshold {
                newLevel = models.TrustLevelBasic
        } else {
                newLevel = models.TrustLevelNewUser
        }
        
        // Update level if it changed
        if newLevel != score.TrustLevel {
                score.TrustLevel = newLevel
                score.UpdatedAt = time.Now()
                
                if err := s.moderationRepo.UpdateUserTrustScore(score); err != nil {
                        return newLevel, fmt.Errorf("error updating user trust level: %w", err)
                }
        }
        
        return newLevel, nil
}

// GetUsersByTrustLevel retrieves users by trust level
func (s *ModerationServiceImpl) GetUsersByTrustLevel(level models.UserTrustLevel) ([]models.UserTrustScore, error) {
        return s.moderationRepo.GetUsersByTrustLevel(level)
}

// GrantModeratorPrivileges grants moderator privileges to a user
func (s *ModerationServiceImpl) GrantModeratorPrivileges(
        userID, grantedByID uint, 
        privileges map[string]bool,
) (*models.ModeratorPrivilege, error) {
        // Check if user already has moderator privileges
        existingPrivileges, err := s.moderationRepo.GetModeratorPrivileges(userID)
        if err == nil {
                // User already has privileges
                if existingPrivileges.IsActive {
                        return nil, errors.New("user already has moderator privileges")
                }
                
                // Reactivate and update privileges
                existingPrivileges.IsActive = true
                existingPrivileges.UpdatedAt = time.Now()
                
                // Update specific privileges
                if val, ok := privileges["canApproveContent"]; ok {
                        existingPrivileges.CanApproveContent = val
                }
                if val, ok := privileges["canRejectContent"]; ok {
                        existingPrivileges.CanRejectContent = val
                }
                if val, ok := privileges["canEditContent"]; ok {
                        existingPrivileges.CanEditContent = val
                }
                if val, ok := privileges["canDeleteContent"]; ok {
                        existingPrivileges.CanDeleteContent = val
                }
                if val, ok := privileges["canBanUsers"]; ok {
                        existingPrivileges.CanBanUsers = val
                }
                if val, ok := privileges["canManageRules"]; ok {
                        existingPrivileges.CanManageRules = val
                }
                if val, ok := privileges["canAssignModerators"]; ok {
                        existingPrivileges.CanAssignModerators = val
                }
                if val, ok := privileges["canAccessDashboard"]; ok {
                        existingPrivileges.CanAccessDashboard = val
                }
                
                if err := s.moderationRepo.UpdateModeratorPrivileges(existingPrivileges); err != nil {
                        return nil, fmt.Errorf("error updating moderator privileges: %w", err)
                }
                
                return existingPrivileges, nil
        }
        
        // Create new moderator privileges
        newPrivileges := &models.ModeratorPrivilege{
                UserID:             userID,
                CanApproveContent:  privileges["canApproveContent"],
                CanRejectContent:   privileges["canRejectContent"],
                CanEditContent:     privileges["canEditContent"],
                CanDeleteContent:   privileges["canDeleteContent"],
                CanBanUsers:        privileges["canBanUsers"],
                CanManageRules:     privileges["canManageRules"],
                CanAssignModerators: privileges["canAssignModerators"],
                CanAccessDashboard: privileges["canAccessDashboard"],
                IsActive:           true,
                AssignedBy:         grantedByID,
                CreatedAt:          time.Now(),
                UpdatedAt:          time.Now(),
        }
        
        if err := s.moderationRepo.CreateModeratorPrivileges(newPrivileges); err != nil {
                return nil, fmt.Errorf("error creating moderator privileges: %w", err)
        }
        
        return newPrivileges, nil
}

// UpdateModeratorPrivileges updates a moderator's privileges
func (s *ModerationServiceImpl) UpdateModeratorPrivileges(
        userID, updatedByID uint, 
        privileges map[string]bool,
) (*models.ModeratorPrivilege, error) {
        // Get current privileges
        currentPrivileges, err := s.moderationRepo.GetModeratorPrivileges(userID)
        if err != nil {
                return nil, fmt.Errorf("error getting moderator privileges: %w", err)
        }
        
        // Update privileges
        if val, ok := privileges["canApproveContent"]; ok {
                currentPrivileges.CanApproveContent = val
        }
        if val, ok := privileges["canRejectContent"]; ok {
                currentPrivileges.CanRejectContent = val
        }
        if val, ok := privileges["canEditContent"]; ok {
                currentPrivileges.CanEditContent = val
        }
        if val, ok := privileges["canDeleteContent"]; ok {
                currentPrivileges.CanDeleteContent = val
        }
        if val, ok := privileges["canBanUsers"]; ok {
                currentPrivileges.CanBanUsers = val
        }
        if val, ok := privileges["canManageRules"]; ok {
                currentPrivileges.CanManageRules = val
        }
        if val, ok := privileges["canAssignModerators"]; ok {
                currentPrivileges.CanAssignModerators = val
        }
        if val, ok := privileges["canAccessDashboard"]; ok {
                currentPrivileges.CanAccessDashboard = val
        }
        
        currentPrivileges.UpdatedAt = time.Now()
        
        if err := s.moderationRepo.UpdateModeratorPrivileges(currentPrivileges); err != nil {
                return nil, fmt.Errorf("error updating moderator privileges: %w", err)
        }
        
        return currentPrivileges, nil
}

// RevokeModeratorPrivileges revokes moderator privileges from a user
func (s *ModerationServiceImpl) RevokeModeratorPrivileges(userID, revokedByID uint) error {
        // Get current privileges
        currentPrivileges, err := s.moderationRepo.GetModeratorPrivileges(userID)
        if err != nil {
                return fmt.Errorf("error getting moderator privileges: %w", err)
        }
        
        // Deactivate privileges
        currentPrivileges.IsActive = false
        currentPrivileges.UpdatedAt = time.Now()
        
        if err := s.moderationRepo.UpdateModeratorPrivileges(currentPrivileges); err != nil {
                return fmt.Errorf("error updating moderator privileges: %w", err)
        }
        
        return nil
}

// GetModeratorPrivileges retrieves moderator privileges for a user
func (s *ModerationServiceImpl) GetModeratorPrivileges(userID uint) (*models.ModeratorPrivilege, error) {
        return s.moderationRepo.GetModeratorPrivileges(userID)
}

// GetAllModerators retrieves all moderators
func (s *ModerationServiceImpl) GetAllModerators() ([]models.ModeratorPrivilege, error) {
        return s.moderationRepo.GetAllModerators()
}

// IsUserModerator checks if a user is a moderator
func (s *ModerationServiceImpl) IsUserModerator(userID uint) (bool, error) {
        return s.moderationRepo.IsUserModerator(userID)
}

// CreateUserModerationAction creates a moderation action for a user
func (s *ModerationServiceImpl) CreateUserModerationAction(
        userID uint, 
        actionType models.ModerationAction, 
        reason string, 
        duration int, 
        relatedContentID *uint, 
        relatedContentType *string, 
        moderatorID uint, 
        notes string,
) (*models.UserModerationAction, error) {
        // Validate action type
        if actionType != models.ActionWarning && 
           actionType != models.ActionTemporaryBan && 
           actionType != models.ActionPermanentBan {
                return nil, errors.New("invalid action type")
        }
        
        // Calculate expiration date for temporary actions
        var expiresAt *time.Time
        if actionType == models.ActionTemporaryBan {
                if duration <= 0 {
                        return nil, errors.New("duration must be positive for temporary bans")
                }
                
                expires := time.Now().AddDate(0, 0, duration)
                expiresAt = &expires
        }
        
        // Create action - convert types as needed
        // Convert ModerationAction to string
        actionTypeStr := string(actionType)
        
        // Convert int to *int for Duration
        var durationPtr *int
        if duration > 0 {
                durationPtr = &duration
        }
        
        // Default content type
        contentTypeStr := ""
        if relatedContentType != nil {
                contentTypeStr = *relatedContentType
        }
        
        now := time.Now()
        
        action := &models.UserModerationAction{
                UserID:             userID,
                ActionType:         actionTypeStr,
                Reason:             reason,
                Duration:           durationPtr,
                ExpiresAt:          expiresAt,
                RelatedContentID:   relatedContentID,
                RelatedContentType: contentTypeStr,
                ModeratorID:        moderatorID,
                AppliedBy:          moderatorID,
                AppliedAt:          now,
                IsActive:           true,
                Notes:              notes,
        }
        
        if err := s.moderationRepo.CreateUserModerationAction(action); err != nil {
                return nil, fmt.Errorf("error creating user moderation action: %w", err)
        }
        
        // Update user trust score
        trustScore, err := s.moderationRepo.GetUserTrustScore(userID)
        if err != nil {
                // Log error but continue
                fmt.Printf("Error getting user trust score: %v\n", err)
        } else {
                if actionType == models.ActionWarning {
                        trustScore.WarningCount++
                }
                
                trustScore.UpdatedAt = time.Now()
                
                if err := s.moderationRepo.UpdateUserTrustScore(trustScore); err != nil {
                        // Log error but continue
                        fmt.Printf("Error updating user trust score: %v\n", err)
                }
                
                // Recalculate trust level
                if _, err := s.RecalculateUserTrustLevel(userID); err != nil {
                        // Log error but continue
                        fmt.Printf("Error recalculating user trust level: %v\n", err)
                }
        }
        
        return action, nil
}

// GetUserModerationActions retrieves moderation actions for a user
func (s *ModerationServiceImpl) GetUserModerationActions(userID uint) ([]models.UserModerationAction, error) {
        return s.moderationRepo.GetUserModerationActions(userID)
}

// RevokeUserModerationAction revokes a user moderation action
func (s *ModerationServiceImpl) RevokeUserModerationAction(actionID, moderatorID uint, reason string) error {
        // Get the action
        action, err := s.moderationRepo.GetUserModerationActionByID(actionID)
        if err != nil {
                return fmt.Errorf("error getting user moderation action: %w", err)
        }
        
        // Append revocation reason to notes
        if reason != "" {
                action.Notes += "\n\nRevoked by moderator ID " + fmt.Sprintf("%d", moderatorID) + " on " + 
                        time.Now().Format(time.RFC3339) + " with reason: " + reason
        }
        
        // Deactivate the action
        return s.moderationRepo.DeactivateUserModerationAction(actionID)
}

// GetActiveUserActions retrieves active moderation actions for a user
func (s *ModerationServiceImpl) GetActiveUserActions(userID uint) ([]models.UserModerationAction, error) {
        return s.moderationRepo.GetActiveActionsForUser(userID)
}

// IsUserBanned checks if a user is banned
func (s *ModerationServiceImpl) IsUserBanned(userID uint) (bool, error) {
        actions, err := s.moderationRepo.GetActiveActionsForUser(userID)
        if err != nil {
                return false, fmt.Errorf("error getting active actions: %w", err)
        }
        
        for _, action := range actions {
                // Convert string to ModerationAction for comparison
                actionTypeEnum := models.ModerationAction(action.ActionType)
                
                if actionTypeEnum == models.ActionTemporaryBan || actionTypeEnum == models.ActionPermanentBan {
                        return true, nil
                }
        }
        
        return false, nil
}

// AddProhibitedWord adds a prohibited word
func (s *ModerationServiceImpl) AddProhibitedWord(
        word, replacement string, 
        isRegex, isAutoReplace bool, 
        severity int, 
        userID uint,
) (*models.ProhibitedWord, error) {
        // Validate word
        if word == "" {
                return nil, errors.New("word cannot be empty")
        }
        
        // Validate regex
        if isRegex {
                _, err := regexp.Compile(word)
                if err != nil {
                        return nil, fmt.Errorf("invalid regex pattern: %w", err)
                }
        }
        
        // Validate severity
        if severity < 1 || severity > 10 {
                return nil, errors.New("severity must be between 1 and 10")
        }
        
        // Create prohibited word
        prohibitedWord := &models.ProhibitedWord{
                Word:          word,
                Replacement:   replacement,
                IsRegex:       isRegex,
                IsAutoReplace: isAutoReplace,
                Severity:      severity,
                CreatedBy:     userID,
                IsActive:      true,
                CreatedAt:     time.Now(),
                UpdatedAt:     time.Now(),
        }
        
        if err := s.moderationRepo.CreateProhibitedWord(prohibitedWord); err != nil {
                return nil, fmt.Errorf("error creating prohibited word: %w", err)
        }
        
        return prohibitedWord, nil
}

// GetProhibitedWords retrieves prohibited words
func (s *ModerationServiceImpl) GetProhibitedWords(activeOnly bool) ([]models.ProhibitedWord, error) {
        return s.moderationRepo.GetProhibitedWords(activeOnly)
}

// UpdateProhibitedWord updates a prohibited word
func (s *ModerationServiceImpl) UpdateProhibitedWord(
        wordID uint, 
        word, replacement string, 
        isRegex, isAutoReplace, isActive bool, 
        severity int, 
        userID uint,
) (*models.ProhibitedWord, error) {
        // Get existing word
        prohibitedWord, err := s.moderationRepo.GetProhibitedWordByID(wordID)
        if err != nil {
                return nil, fmt.Errorf("error getting prohibited word: %w", err)
        }
        
        // Validate word
        if word == "" {
                return nil, errors.New("word cannot be empty")
        }
        
        // Validate regex
        if isRegex {
                _, err := regexp.Compile(word)
                if err != nil {
                        return nil, fmt.Errorf("invalid regex pattern: %w", err)
                }
        }
        
        // Validate severity
        if severity < 1 || severity > 10 {
                return nil, errors.New("severity must be between 1 and 10")
        }
        
        // Update word
        prohibitedWord.Word = word
        prohibitedWord.Replacement = replacement
        prohibitedWord.IsRegex = isRegex
        prohibitedWord.IsAutoReplace = isAutoReplace
        prohibitedWord.IsActive = isActive
        prohibitedWord.Severity = severity
        prohibitedWord.UpdatedAt = time.Now()
        
        if err := s.moderationRepo.UpdateProhibitedWord(prohibitedWord); err != nil {
                return nil, fmt.Errorf("error updating prohibited word: %w", err)
        }
        
        return prohibitedWord, nil
}

// DeleteProhibitedWord deletes a prohibited word
func (s *ModerationServiceImpl) DeleteProhibitedWord(wordID, userID uint) error {
        // Check if word exists
        if _, err := s.moderationRepo.GetProhibitedWordByID(wordID); err != nil {
                return fmt.Errorf("error getting prohibited word: %w", err)
        }
        
        return s.moderationRepo.DeleteProhibitedWord(wordID)
}

// FilterTextWithProhibitedWords filters text for prohibited words
func (s *ModerationServiceImpl) FilterTextWithProhibitedWords(text string) (string, bool, []string, error) {
        // Get active prohibited words
        words, err := s.moderationRepo.GetProhibitedWords(true)
        if err != nil {
                return text, false, nil, fmt.Errorf("error getting prohibited words: %w", err)
        }
        
        filtered := text
        wasFiltered := false
        filteredWords := []string{}
        
        // Apply each word filter
        for _, word := range words {
                if word.IsRegex {
                        re, err := regexp.Compile(word.Word)
                        if err != nil {
                                // Log error but continue
                                fmt.Printf("Invalid regex pattern in prohibited word %d: %v\n", word.ID, err)
                                continue
                        }
                        
                        matches := re.FindAllString(filtered, -1)
                        if len(matches) > 0 {
                                wasFiltered = true
                                for _, match := range matches {
                                        if !contains(filteredWords, match) {
                                                filteredWords = append(filteredWords, match)
                                        }
                                }
                                
                                if word.IsAutoReplace {
                                        filtered = re.ReplaceAllString(filtered, word.Replacement)
                                }
                        }
                } else {
                        // Simple keyword replacement
                        if strings.Contains(strings.ToLower(filtered), strings.ToLower(word.Word)) {
                                wasFiltered = true
                                if !contains(filteredWords, word.Word) {
                                        filteredWords = append(filteredWords, word.Word)
                                }
                                
                                if word.IsAutoReplace {
                                        // Case-insensitive replacement
                                        re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(word.Word))
                                        if err != nil {
                                                // Log error but continue
                                                fmt.Printf("Error creating regex for word '%s': %v\n", word.Word, err)
                                                continue
                                        }
                                        
                                        filtered = re.ReplaceAllString(filtered, word.Replacement)
                                }
                        }
                }
        }
        
        return filtered, wasFiltered, filteredWords, nil
}

// Helper functions

// isMoreRestrictiveAction checks if one action is more restrictive than another
func (s *ModerationServiceImpl) isMoreRestrictiveAction(a, b models.ModerationAction) bool {
        actionOrder := map[models.ModerationAction]int{
                models.ActionNone:           0,
                models.ActionApprove:        1,
                models.ActionSendToQueue:    2,
                models.ActionAutomaticFilter: 3,
                models.ActionReject:         4,
                models.ActionWarning:        5,
                models.ActionTemporaryBan:   6,
                models.ActionPermanentBan:   7,
        }
        
        return actionOrder[a] > actionOrder[b]
}

// contains checks if a string slice contains a string
func contains(slice []string, str string) bool {
        for _, s := range slice {
                if s == str {
                        return true
                }
        }
        return false
}