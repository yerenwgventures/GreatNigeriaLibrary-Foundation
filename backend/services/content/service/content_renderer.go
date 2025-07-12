package service

import (
        "encoding/json"
        "fmt"
        "html"
        "regexp"
        "strings"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/repository"
)

// ContentRenderer defines the interface for content rendering services
type ContentRenderer interface {
        RenderSection(section *models.BookSection, userID uint) (string, error)
        RenderMarkdown(markdown string) (string, error)
        ProcessInteractiveElements(content string, sectionID, userID uint) (string, error)
        EnhanceContentWithMedia(content string) (string, error)
        InsertTopicLinks(content string) (string, error)
}

// ContentRendererImpl implements the ContentRenderer interface
type ContentRendererImpl struct {
        elementRepo  repository.InteractiveElementRepository
        discussionRepo repository.DiscussionRepository
}

// NewContentRenderer creates a new content renderer instance with all required dependencies
func NewContentRenderer(elementRepo repository.InteractiveElementRepository, discussionRepo repository.DiscussionRepository) ContentRenderer {
        return &ContentRendererImpl{
                elementRepo:    elementRepo,
                discussionRepo: discussionRepo,
        }
}

// RenderSection renders a section's content with all enhancements
func (r *ContentRendererImpl) RenderSection(section *models.BookSection, userID uint) (string, error) {
        content := section.Content
        
        // Apply Markdown rendering if the content format is markdown
        if section.Format == "markdown" {
                renderedContent, err := r.RenderMarkdown(content)
                if err != nil {
                        return "", err
                }
                content = renderedContent
        }
        
        // Process interactive elements
        content, err := r.ProcessInteractiveElements(content, section.ID, userID)
        if err != nil {
                return "", err
        }
        
        // Enhance content with media
        content, err = r.EnhanceContentWithMedia(content)
        if err != nil {
                return "", err
        }
        
        // Insert discussion topic links
        content, err = r.InsertTopicLinks(content)
        if err != nil {
                return "", err
        }
        
        return content, nil
}

// RenderMarkdown converts markdown to HTML
func (r *ContentRendererImpl) RenderMarkdown(markdown string) (string, error) {
        // This is a simplified implementation. For production, use a proper markdown library.
        // In Go, you might use a library like github.com/gomarkdown/markdown
        
        // Convert headers
        headerRegex := regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
        html := headerRegex.ReplaceAllStringFunc(markdown, func(match string) string {
                parts := headerRegex.FindStringSubmatch(match)
                level := len(parts[1])
                text := parts[2]
                return "<h" + string(rune('0'+level)) + ">" + text + "</h" + string(rune('0'+level)) + ">"
        })
        
        // Convert bold
        boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
        html = boldRegex.ReplaceAllString(html, "<strong>$1</strong>")
        
        // Convert italic
        italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
        html = italicRegex.ReplaceAllString(html, "<em>$1</em>")
        
        // Convert paragraphs (simplified)
        paragraphRegex := regexp.MustCompile(`(?m)^([^<].+)$`)
        html = paragraphRegex.ReplaceAllString(html, "<p>$1</p>")
        
        // Convert lists (simplified)
        listItemRegex := regexp.MustCompile(`(?m)^-\s+(.+)$`)
        html = listItemRegex.ReplaceAllString(html, "<li>$1</li>")
        
        // Wrap list items in <ul>
        listRegex := regexp.MustCompile(`(?s)(<li>.+</li>)`)
        html = listRegex.ReplaceAllString(html, "<ul>$1</ul>")
        
        // Convert links
        linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
        html = linkRegex.ReplaceAllString(html, "<a href=\"$2\">$1</a>")
        
        // Convert images
        imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
        html = imageRegex.ReplaceAllString(html, "<img src=\"$2\" alt=\"$1\">")
        
        // Fix double paragraph tags
        html = strings.ReplaceAll(html, "<p><p>", "<p>")
        html = strings.ReplaceAll(html, "</p></p>", "</p>")
        
        return html, nil
}

// ProcessInteractiveElements processes interactive element placeholders in content
func (r *ContentRendererImpl) ProcessInteractiveElements(content string, sectionID, userID uint) (string, error) {
        // Skip if element repository is not available
        if r.elementRepo == nil {
                return content, nil
        }

        // Find interactive element placeholders with regex
        placeholderRegex := regexp.MustCompile(`\{\{interactive:(\d+)\}\}`)
        
        // Get all interactive elements for this section
        elements, err := r.elementRepo.GetInteractiveElementsBySection(sectionID)
        if err != nil {
                return content, nil // Just return original content on error
        }
        
        // Create map of elements by ID for quick lookup
        elementMap := make(map[uint]models.InteractiveElement)
        for _, element := range elements {
                elementMap[element.ID] = element
        }
        
        // Replace each placeholder with the actual interactive element HTML
        content = placeholderRegex.ReplaceAllStringFunc(content, func(match string) string {
                parts := placeholderRegex.FindStringSubmatch(match)
                if len(parts) < 2 {
                        return match // Keep the placeholder if it's malformed
                }
                
                // Parse the element ID
                var elementID uint
                _, err := fmt.Sscanf(parts[1], "%d", &elementID)
                if err != nil {
                        return match // Keep the placeholder if ID is invalid
                }
                
                // Find the element
                element, exists := elementMap[elementID]
                if !exists {
                        return match // Keep the placeholder if element doesn't exist
                }
                
                // Render the element based on its type
                return renderInteractiveElement(&element, userID)
        })
        
        return content, nil
}

// EnhanceContentWithMedia processes media embed placeholders in content
func (r *ContentRendererImpl) EnhanceContentWithMedia(content string) (string, error) {
        // Process image embeds
        imageRegex := regexp.MustCompile(`\{\{image:([^}]+)\}\}`)
        content = imageRegex.ReplaceAllStringFunc(content, func(match string) string {
                parts := imageRegex.FindStringSubmatch(match)
                if len(parts) < 2 {
                        return match
                }
                
                imagePath := parts[1]
                return "<img src=\"" + html.EscapeString(imagePath) + "\" class=\"embedded-image\" alt=\"Embedded image\">"
        })
        
        // Process video embeds
        videoRegex := regexp.MustCompile(`\{\{video:([^}]+)\}\}`)
        content = videoRegex.ReplaceAllStringFunc(content, func(match string) string {
                parts := videoRegex.FindStringSubmatch(match)
                if len(parts) < 2 {
                        return match
                }
                
                videoPath := parts[1]
                
                // Check if it's a YouTube URL
                if strings.Contains(videoPath, "youtube.com") || strings.Contains(videoPath, "youtu.be") {
                        // Extract video ID (simplified)
                        var videoID string
                        if strings.Contains(videoPath, "v=") {
                                idParts := strings.Split(videoPath, "v=")
                                if len(idParts) > 1 {
                                        videoID = strings.Split(idParts[1], "&")[0]
                                }
                        } else if strings.Contains(videoPath, "youtu.be/") {
                                idParts := strings.Split(videoPath, "youtu.be/")
                                if len(idParts) > 1 {
                                        videoID = idParts[1]
                                }
                        }
                        
                        if videoID != "" {
                                return "<div class=\"video-embed youtube-embed\">" +
                                        "<iframe width=\"560\" height=\"315\" src=\"https://www.youtube.com/embed/" + videoID + "\" " +
                                        "frameborder=\"0\" allow=\"accelerometer; autoplay; clipboard-write; encrypted-media; " +
                                        "gyroscope; picture-in-picture\" allowfullscreen></iframe></div>"
                        }
                }
                
                // Default to regular video player
                return "<div class=\"video-embed\"><video controls><source src=\"" + 
                        html.EscapeString(videoPath) + "\" type=\"video/mp4\">Your browser does not support the video tag.</video></div>"
        })
        
        // Process audio embeds
        audioRegex := regexp.MustCompile(`\{\{audio:([^}]+)\}\}`)
        content = audioRegex.ReplaceAllStringFunc(content, func(match string) string {
                parts := audioRegex.FindStringSubmatch(match)
                if len(parts) < 2 {
                        return match
                }
                
                audioPath := parts[1]
                return "<div class=\"audio-embed\"><audio controls><source src=\"" + 
                        html.EscapeString(audioPath) + "\" type=\"audio/mpeg\">Your browser does not support the audio tag.</audio></div>"
        })
        
        return content, nil
}

// InsertTopicLinks processes discussion topic link placeholders in content
func (r *ContentRendererImpl) InsertTopicLinks(content string) (string, error) {
        // Skip if discussion repository is not available
        if r.discussionRepo == nil {
                return content, nil
        }
        
        // Process topic links
        topicRegex := regexp.MustCompile(`\{\{topic:(\d+)\}\}`)
        content = topicRegex.ReplaceAllStringFunc(content, func(match string) string {
                parts := topicRegex.FindStringSubmatch(match)
                if len(parts) < 2 {
                        return match
                }
                
                var topicID uint
                _, err := fmt.Sscanf(parts[1], "%d", &topicID)
                if err != nil {
                        return match
                }
                
                // Get topic details
                topic, err := r.discussionRepo.GetDiscussionTopicByID(topicID)
                if err != nil {
                        return match
                }
                
                return "<div class=\"discussion-topic-link\">" +
                        "<h4>Join the Discussion</h4>" +
                        "<p>" + html.EscapeString(topic.Title) + "</p>" +
                        "<a href=\"/discussion/topic/" + fmt.Sprintf("%d", topic.ID) + "\" class=\"btn btn-primary\">View Discussion</a>" +
                        "</div>"
        })
        
        return content, nil
}

// Helper function to render an interactive element
func renderInteractiveElement(element *models.InteractiveElement, userID uint) string {
        // Common element wrapper start
        htmlContent := "<div class=\"interactive-element " + string(element.Type) + "-element\" data-element-id=\"" + 
                fmt.Sprintf("%d", element.ID) + "\">\n"
        
        // Add title and description
        htmlContent += "<h3 class=\"interactive-title\">" + html.EscapeString(element.Title) + "</h3>\n"
        htmlContent += "<div class=\"interactive-description\">" + html.EscapeString(element.Description) + "</div>\n"
        
        // Render type-specific content
        switch element.Type {
        case models.QuizType:
                htmlContent += renderQuizContent(element)
        case models.ReflectionType:
                htmlContent += renderReflectionContent(element)
        case models.CallToActionType:
                htmlContent += renderCallToActionContent(element)
        case models.DiscussionPromptType:
                htmlContent += renderDiscussionPromptContent(element)
        case models.PollType:
                htmlContent += renderPollContent(element)
        default:
                htmlContent += "<div class=\"interactive-content\">Unknown element type</div>\n"
        }
        
        // Add element wrapper end
        htmlContent += "</div>\n"
        
        return htmlContent
}

// Helper function to render quiz content
func renderQuizContent(element *models.InteractiveElement) string {
        // Parse quiz content
        var quizContent models.QuizContent
        err := json.Unmarshal([]byte(element.Content), &quizContent)
        if err != nil {
                return "<div class=\"error\">Error loading quiz content</div>"
        }
        
        htmlContent := "<div class=\"quiz-content\">\n"
        
        // Add quiz instructions
        htmlContent += "<div class=\"quiz-instructions\">"
        if quizContent.TimeLimit > 0 {
                htmlContent += "<p>Time limit: " + fmt.Sprintf("%d", quizContent.TimeLimit) + " seconds</p>"
        }
        htmlContent += "<p>Pass score: " + fmt.Sprintf("%d", quizContent.PassScore) + "%</p>"
        htmlContent += "</div>\n"
        
        // Render each question
        htmlContent += "<form class=\"quiz-form\" data-quiz-id=\"" + fmt.Sprintf("%d", element.ID) + "\">\n"
        
        for i, question := range quizContent.Questions {
                htmlContent += "<div class=\"quiz-question\" data-question-id=\"" + fmt.Sprintf("%d", question.ID) + "\">\n"
                htmlContent += "<div class=\"question-text\">" + html.EscapeString(question.QuestionText) + "</div>\n"
                
                // Render media if attached
                if question.Media != nil {
                        htmlContent += renderMediaAttachment(question.Media)
                }
                
                // Render question based on type
                switch question.QuestionType {
                case "multiple-choice":
                        htmlContent += renderMultipleChoiceQuestion(question, i)
                case "true-false":
                        htmlContent += renderTrueFalseQuestion(question, i)
                case "fill-blank":
                        htmlContent += renderFillBlankQuestion(question, i)
                case "short-answer":
                        htmlContent += renderShortAnswerQuestion(question, i)
                }
                
                htmlContent += "</div>\n" // End question div
        }
        
        // Add submit button
        htmlContent += "<div class=\"quiz-actions\">\n"
        htmlContent += "<button type=\"submit\" class=\"btn btn-primary quiz-submit\">Submit Answers</button>\n"
        htmlContent += "</div>\n"
        
        htmlContent += "</form>\n" // End form
        htmlContent += "</div>\n" // End quiz content
        
        return htmlContent
}

// Helper function to render reflection content
func renderReflectionContent(element *models.InteractiveElement) string {
        // Parse reflection content
        var reflectionContent models.ReflectionContent
        err := json.Unmarshal([]byte(element.Content), &reflectionContent)
        if err != nil {
                return "<div class=\"error\">Error loading reflection content</div>"
        }
        
        htmlContent := "<div class=\"reflection-content\">\n"
        
        // Add reflection prompt
        htmlContent += "<div class=\"reflection-prompt\">" + html.EscapeString(reflectionContent.Prompt) + "</div>\n"
        
        // Add guiding questions if any
        if len(reflectionContent.GuidingQuestions) > 0 {
                htmlContent += "<div class=\"guiding-questions\">\n"
                htmlContent += "<h4>Guiding Questions:</h4>\n"
                htmlContent += "<ul>\n"
                
                for _, question := range reflectionContent.GuidingQuestions {
                        htmlContent += "<li>" + html.EscapeString(question) + "</li>\n"
                }
                
                htmlContent += "</ul>\n"
                htmlContent += "</div>\n"
        }
        
        // Add response form
        htmlContent += "<form class=\"reflection-form\" data-reflection-id=\"" + fmt.Sprintf("%d", element.ID) + "\">\n"
        
        // Add textarea for response
        htmlContent += "<div class=\"form-group\">\n"
        htmlContent += "<label for=\"reflection-response-" + fmt.Sprintf("%d", element.ID) + "\">Your Reflection:</label>\n"
        htmlContent += "<textarea id=\"reflection-response-" + fmt.Sprintf("%d", element.ID) + "\" " +
                "class=\"form-control reflection-response\" rows=\"5\" "
        
        // Add min/max length attributes if specified
        if reflectionContent.MinResponseLength > 0 {
                htmlContent += "minlength=\"" + fmt.Sprintf("%d", reflectionContent.MinResponseLength) + "\" "
        }
        if reflectionContent.MaxResponseLength > 0 {
                htmlContent += "maxlength=\"" + fmt.Sprintf("%d", reflectionContent.MaxResponseLength) + "\" "
        }
        
        htmlContent += "required></textarea>\n"
        htmlContent += "</div>\n"
        
        // Add sharing options if specified
        if len(reflectionContent.SharingOptions) > 0 {
                htmlContent += "<div class=\"sharing-options\">\n"
                htmlContent += "<h4>Sharing Options:</h4>\n"
                
                for i, option := range reflectionContent.SharingOptions {
                        htmlContent += "<div class=\"form-check\">\n"
                        htmlContent += "<input class=\"form-check-input\" type=\"radio\" name=\"sharing-option\" " +
                                "id=\"sharing-option-" + fmt.Sprintf("%d-%d", element.ID, i) + "\" " +
                                "value=\"" + html.EscapeString(option) + "\" " +
                                func() string {
                    if i == 0 {
                        return "checked"
                    }
                    return ""
                }() + ">\n"
                        htmlContent += "<label class=\"form-check-label\" for=\"sharing-option-" + fmt.Sprintf("%d-%d", element.ID, i) + "\">" +
                                html.EscapeString(option) + "</label>\n"
                        htmlContent += "</div>\n"
                }
                
                htmlContent += "</div>\n"
        }
        
        // Add submit button
        htmlContent += "<div class=\"reflection-actions\">\n"
        htmlContent += "<button type=\"submit\" class=\"btn btn-primary reflection-submit\">Submit Reflection</button>\n"
        htmlContent += "</div>\n"
        
        htmlContent += "</form>\n" // End form
        htmlContent += "</div>\n" // End reflection content
        
        return htmlContent
}

// Helper function to render call-to-action content
func renderCallToActionContent(element *models.InteractiveElement) string {
        // Parse call-to-action content
        var ctaContent models.CallToActionContent
        err := json.Unmarshal([]byte(element.Content), &ctaContent)
        if err != nil {
                return "<div class=\"error\">Error loading call-to-action content</div>"
        }
        
        htmlContent := "<div class=\"call-to-action-content\">\n"
        
        // Add secondary text if specified
        if ctaContent.SecondaryText != "" {
                htmlContent += "<div class=\"cta-secondary-text\">" + html.EscapeString(ctaContent.SecondaryText) + "</div>\n"
        }
        
        // Add button
        htmlContent += "<div class=\"cta-button-container\">\n"
        htmlContent += "<button type=\"button\" class=\"btn btn-primary cta-button\" " + 
                "data-action-type=\"" + html.EscapeString(ctaContent.ActionType) + "\" " +
                "data-element-id=\"" + fmt.Sprintf("%d", element.ID) + "\" " +
                func() string {
                    if ctaContent.URL != "" {
                        return "data-url=\"" + html.EscapeString(ctaContent.URL) + "\" "
                    }
                    return ""
                }() + 
                ">" + html.EscapeString(ctaContent.ButtonText) + "</button>\n"
        htmlContent += "</div>\n"
        
        htmlContent += "</div>\n" // End call-to-action content
        
        return htmlContent
}

// Helper function to render discussion prompt content
func renderDiscussionPromptContent(element *models.InteractiveElement) string {
        // Parse discussion prompt content
        var discussionContent models.DiscussionPromptContent
        err := json.Unmarshal([]byte(element.Content), &discussionContent)
        if err != nil {
                return "<div class=\"error\">Error loading discussion prompt content</div>"
        }
        
        htmlContent := "<div class=\"discussion-prompt-content\">\n"
        
        // Add topic and initial prompt
        htmlContent += "<div class=\"discussion-topic\">" + html.EscapeString(discussionContent.Topic) + "</div>\n"
        htmlContent += "<div class=\"discussion-initial-prompt\">" + html.EscapeString(discussionContent.InitialPrompt) + "</div>\n"
        
        // Add supporting points if any
        if len(discussionContent.SupportingPoints) > 0 {
                htmlContent += "<div class=\"supporting-points\">\n"
                htmlContent += "<h4>Points to Consider:</h4>\n"
                htmlContent += "<ul>\n"
                
                for _, point := range discussionContent.SupportingPoints {
                        htmlContent += "<li>" + html.EscapeString(point) + "</li>\n"
                }
                
                htmlContent += "</ul>\n"
                htmlContent += "</div>\n"
        }
        
        // Add guidelines if specified
        if discussionContent.Guidelines != "" {
                htmlContent += "<div class=\"discussion-guidelines\">\n"
                htmlContent += "<h4>Discussion Guidelines:</h4>\n"
                htmlContent += "<p>" + html.EscapeString(discussionContent.Guidelines) + "</p>\n"
                htmlContent += "</div>\n"
        }
        
        // Add form for response
        htmlContent += "<form class=\"discussion-form\" data-discussion-id=\"" + fmt.Sprintf("%d", element.ID) + "\">\n"
        
        // Add topic ID as hidden field
        if discussionContent.DiscussionForumID > 0 {
                htmlContent += "<input type=\"hidden\" name=\"topic-id\" value=\"" + fmt.Sprintf("%d", discussionContent.DiscussionForumID) + "\">\n"
        }
        
        // Add textarea for response
        htmlContent += "<div class=\"form-group\">\n"
        htmlContent += "<label for=\"discussion-response-" + fmt.Sprintf("%d", element.ID) + "\">Your Response:</label>\n"
        htmlContent += "<textarea id=\"discussion-response-" + fmt.Sprintf("%d", element.ID) + "\" " +
                "class=\"form-control discussion-response\" rows=\"5\" required></textarea>\n"
        htmlContent += "</div>\n"
        
        // Add submit button
        htmlContent += "<div class=\"discussion-actions\">\n"
        htmlContent += "<button type=\"submit\" class=\"btn btn-primary discussion-submit\">Submit Response</button>\n"
        htmlContent += "</div>\n"
        
        htmlContent += "</form>\n" // End form
        htmlContent += "</div>\n" // End discussion prompt content
        
        return htmlContent
}

// Helper function to render poll content
func renderPollContent(element *models.InteractiveElement) string {
        // Parse poll content
        var pollContent models.PollContent
        err := json.Unmarshal([]byte(element.Content), &pollContent)
        if err != nil {
                return "<div class=\"error\">Error loading poll content</div>"
        }
        
        htmlContent := "<div class=\"poll-content\">\n"
        
        // Add poll question
        htmlContent += "<div class=\"poll-question\">" + html.EscapeString(pollContent.Question) + "</div>\n"
        
        // Add form
        htmlContent += "<form class=\"poll-form\" data-poll-id=\"" + fmt.Sprintf("%d", element.ID) + "\">\n"
        
        // Add options
        htmlContent += "<div class=\"poll-options\">\n"
        
        for _, option := range pollContent.Options {
                htmlContent += "<div class=\"form-check\">\n"
                htmlContent += "<input class=\"form-check-input\" type=\"" + func() string {
                    if pollContent.AllowMultiple {
                        return "checkbox"
                    }
                    return "radio"
                }() + "\" " +
                        "name=\"poll-option\" id=\"poll-option-" + fmt.Sprintf("%d-%s", element.ID, option.ID) + "\" " +
                        "value=\"" + html.EscapeString(option.ID) + "\">\n"
                htmlContent += "<label class=\"form-check-label\" for=\"poll-option-" + fmt.Sprintf("%d-%s", element.ID, option.ID) + "\">" +
                        html.EscapeString(option.Text) + "</label>\n"
                
                // Add image if specified
                if option.Image != "" {
                        htmlContent += "<img src=\"" + html.EscapeString(option.Image) + "\" class=\"poll-option-image\" alt=\"\">\n"
                }
                
                htmlContent += "</div>\n"
        }
        
        htmlContent += "</div>\n" // End poll options
        
        // Add submit button
        htmlContent += "<div class=\"poll-actions\">\n"
        htmlContent += "<button type=\"submit\" class=\"btn btn-primary poll-submit\">Submit Vote</button>\n"
        htmlContent += "</div>\n"
        
        htmlContent += "</form>\n" // End form
        
        // Add container for results
        htmlContent += "<div class=\"poll-results\" style=\"display: none;\">\n"
        htmlContent += "<h4>Poll Results:</h4>\n"
        htmlContent += "<div class=\"results-container\"></div>\n"
        htmlContent += "</div>\n"
        
        htmlContent += "</div>\n" // End poll content
        
        return htmlContent
}

// Helper functions for rendering question types
func renderMultipleChoiceQuestion(question models.QuizQuestion, index int) string {
        htmlContent := "<div class=\"question-options\">\n"
        
        for _, option := range question.Options {
                htmlContent += "<div class=\"form-check\">\n"
                htmlContent += "<input class=\"form-check-input\" type=\"radio\" name=\"question-" + fmt.Sprintf("%d", question.ID) + "\" " +
                        "id=\"option-" + fmt.Sprintf("%d-%s", question.ID, option.ID) + "\" " +
                        "value=\"" + html.EscapeString(option.ID) + "\">\n"
                htmlContent += "<label class=\"form-check-label\" for=\"option-" + fmt.Sprintf("%d-%s", question.ID, option.ID) + "\">" +
                        html.EscapeString(option.Text) + "</label>\n"
                
                // Add media if attached to option
                if option.Media != nil {
                        htmlContent += renderMediaAttachment(option.Media)
                }
                
                htmlContent += "</div>\n"
        }
        
        htmlContent += "</div>\n"
        return htmlContent
}

func renderTrueFalseQuestion(question models.QuizQuestion, index int) string {
        htmlContent := "<div class=\"question-options\">\n"
        
        htmlContent += "<div class=\"form-check\">\n"
        htmlContent += "<input class=\"form-check-input\" type=\"radio\" name=\"question-" + fmt.Sprintf("%d", question.ID) + "\" " +
                "id=\"option-" + fmt.Sprintf("%d-true", question.ID) + "\" value=\"true\">\n"
        htmlContent += "<label class=\"form-check-label\" for=\"option-" + fmt.Sprintf("%d-true", question.ID) + "\">True</label>\n"
        htmlContent += "</div>\n"
        
        htmlContent += "<div class=\"form-check\">\n"
        htmlContent += "<input class=\"form-check-input\" type=\"radio\" name=\"question-" + fmt.Sprintf("%d", question.ID) + "\" " +
                "id=\"option-" + fmt.Sprintf("%d-false", question.ID) + "\" value=\"false\">\n"
        htmlContent += "<label class=\"form-check-label\" for=\"option-" + fmt.Sprintf("%d-false", question.ID) + "\">False</label>\n"
        htmlContent += "</div>\n"
        
        htmlContent += "</div>\n"
        return htmlContent
}

func renderFillBlankQuestion(question models.QuizQuestion, index int) string {
        htmlContent := "<div class=\"question-input\">\n"
        htmlContent += "<input type=\"text\" class=\"form-control\" name=\"question-" + fmt.Sprintf("%d", question.ID) + "\" " +
                "id=\"answer-" + fmt.Sprintf("%d", question.ID) + "\" placeholder=\"Your answer...\">\n"
        htmlContent += "</div>\n"
        return htmlContent
}

func renderShortAnswerQuestion(question models.QuizQuestion, index int) string {
        htmlContent := "<div class=\"question-input\">\n"
        htmlContent += "<textarea class=\"form-control\" name=\"question-" + fmt.Sprintf("%d", question.ID) + "\" " +
                "id=\"answer-" + fmt.Sprintf("%d", question.ID) + "\" rows=\"3\" placeholder=\"Your answer...\"></textarea>\n"
        htmlContent += "</div>\n"
        return htmlContent
}

// Helper function to render media attachment
func renderMediaAttachment(media *models.MediaAttachment) string {
        htmlContent := "<div class=\"media-attachment\">\n"
        
        switch media.Type {
        case "image":
                htmlContent += "<img src=\"" + html.EscapeString(media.URL) + "\" alt=\"" + html.EscapeString(media.Alt) + "\" class=\"media-image\">\n"
        case "audio":
                htmlContent += "<audio controls class=\"media-audio\">\n"
                htmlContent += "<source src=\"" + html.EscapeString(media.URL) + "\" type=\"audio/mpeg\">\n"
                htmlContent += "Your browser does not support the audio element.\n"
                htmlContent += "</audio>\n"
        case "video":
                htmlContent += "<video controls class=\"media-video\">\n"
                htmlContent += "<source src=\"" + html.EscapeString(media.URL) + "\" type=\"video/mp4\">\n"
                htmlContent += "Your browser does not support the video element.\n"
                htmlContent += "</video>\n"
        }
        
        htmlContent += "</div>\n"
        return htmlContent
}