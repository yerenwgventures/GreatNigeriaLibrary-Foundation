package handlers

import (
        "bytes"
        "encoding/json"
        "fmt"
        "net/http"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/discussion/service"
        "github.com/sirupsen/logrus"
)

// ForumPointsIntegration provides an interface to award points for forum activities
type ForumPointsIntegration struct {
        discussionService service.DiscussionService
        logger            *logrus.Logger
        pointsAPIURL      string
        enabled           bool
}

// Quality constants for content
const (
        QualityLow    = "low"
        QualityMedium = "medium"
        QualityHigh   = "high"
)

// NewForumPointsIntegration creates a new forum points integration
func NewForumPointsIntegration(
        discussionService service.DiscussionService,
        logger *logrus.Logger,
        pointsAPIURL string,
        enabled bool,
) *ForumPointsIntegration {
        return &ForumPointsIntegration{
                discussionService: discussionService,
                logger:            logger,
                pointsAPIURL:      pointsAPIURL,
                enabled:           enabled,
        }
}

// AwardPointsForNewTopic awards points when a user creates a new topic
func (p *ForumPointsIntegration) AwardPointsForNewTopic(userID uint, topicID uint, category string, quality string) error {
        if !p.enabled {
                p.logger.Info("Points integration is disabled, not awarding points for new topic")
                return nil
        }

        // Prepare request to points service
        reqBody := map[string]interface{}{
                "user_id":  userID,
                "topic_id": topicID,
                "category": category,
                "quality":  quality,
        }

        return p.sendPointsRequest("/api/points/discussion/topic", reqBody)
}

// AwardPointsForReply awards points when a user creates a reply
func (p *ForumPointsIntegration) AwardPointsForReply(userID uint, topicID uint, commentID uint, isReply bool, quality string) error {
        if !p.enabled {
                p.logger.Info("Points integration is disabled, not awarding points for reply")
                return nil
        }

        // Prepare request to points service
        reqBody := map[string]interface{}{
                "user_id":    userID,
                "topic_id":   topicID,
                "comment_id": commentID,
                "is_reply":   isReply,
                "quality":    quality,
        }

        return p.sendPointsRequest("/api/points/discussion/reply", reqBody)
}

// AwardPointsForUpvotes awards points when a user receives upvotes
func (p *ForumPointsIntegration) AwardPointsForUpvotes(userID uint, topicID uint, commentID uint, count int) error {
        if !p.enabled {
                p.logger.Info("Points integration is disabled, not awarding points for upvotes")
                return nil
        }

        // Prepare request to points service
        reqBody := map[string]interface{}{
                "user_id":    userID,
                "topic_id":   topicID,
                "comment_id": commentID,
                "count":      count,
        }

        return p.sendPointsRequest("/api/points/discussion/upvote", reqBody)
}

// AwardPointsForFeaturedTopic awards points when a user's topic is featured
func (p *ForumPointsIntegration) AwardPointsForFeaturedTopic(userID uint, topicID uint, category string) error {
        if !p.enabled {
                p.logger.Info("Points integration is disabled, not awarding points for featured topic")
                return nil
        }

        // Prepare request to points service
        reqBody := map[string]interface{}{
                "user_id":  userID,
                "topic_id": topicID,
                "category": category,
        }

        return p.sendPointsRequest("/api/points/discussion/featured", reqBody)
}

// Helper method to send requests to the points service
func (p *ForumPointsIntegration) sendPointsRequest(endpoint string, data map[string]interface{}) error {
        jsonData, err := json.Marshal(data)
        if err != nil {
                p.logger.WithError(err).Error("Failed to marshal points request data")
                return err
        }

        url := p.pointsAPIURL + endpoint
        p.logger.WithFields(logrus.Fields{
                "url":  url,
                "data": string(jsonData),
        }).Debug("Sending points request")

        // Create and send the HTTP request
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
                p.logger.WithError(err).Error("Failed to create points request")
                return err
        }

        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                p.logger.WithError(err).Error("Failed to send points request")
                return err
        }
        defer resp.Body.Close()

        // Check response status
        if resp.StatusCode != http.StatusOK {
                p.logger.WithFields(logrus.Fields{
                        "status": resp.StatusCode,
                        "url":    url,
                }).Error("Points service returned non-OK status")
                return fmt.Errorf("points service returned status %d", resp.StatusCode)
        }

        var respData map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
                p.logger.WithError(err).Error("Failed to decode points response")
                return err
        }

        p.logger.WithFields(logrus.Fields{
                "response": respData,
                "endpoint": endpoint,
        }).Debug("Points request successful")

        return nil
}

// IsEnabled returns whether the points integration is enabled
func (p *ForumPointsIntegration) IsEnabled() bool {
        return p.enabled
}

// SetEnabled enables or disables the points integration
func (p *ForumPointsIntegration) SetEnabled(enabled bool) {
        p.enabled = enabled
}

// DetermineContentQuality analyzes content and returns its quality level
func (p *ForumPointsIntegration) DetermineContentQuality(content string) string {
        // In a real implementation, this would analyze the content for:
        // - Length/comprehensiveness
        // - Formatting quality (headings, paragraphs, etc.)
        // - Media included
        // - Links to sources
        // - Etc.
        
        if len(content) == 0 {
                return QualityLow
        }
        
        wordCount := countWords(content)
        
        // Simple quality heuristic based on length
        if wordCount > 250 {
                return QualityHigh
        } else if wordCount > 100 {
                return QualityMedium
        }
        
        return QualityLow
}

// countWords is a helper function to count the number of words in a string
func countWords(s string) int {
        // Split the string on whitespace
        words := 0
        inWord := false
        
        for _, r := range s {
                if isWordChar(r) {
                        if !inWord {
                                words++
                                inWord = true
                        }
                } else {
                        inWord = false
                }
        }
        
        return words
}

// isWordChar returns true if the rune is a word character (letter, number)
func isWordChar(r rune) bool {
        return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}