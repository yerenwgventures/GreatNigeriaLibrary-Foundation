package service

import (
        "errors"
        "fmt"
        "regexp"
        "strings"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/repository"
)

// RichTextService defines the interface for rich text operations
type RichTextService interface {
        // Rich text content operations
        CreateOrUpdateRichText(contentID uint, contentType string, rawContent string, format models.ContentFormat) (*models.RichTextContent, error)
        GetRichTextContent(contentID uint, contentType string) (*models.RichTextContent, error)
        DeleteRichTextContent(contentID uint, contentType string) error
        
        // Mention operations
        ProcessMentions(contentID uint, contentType string, rawContent string) ([]models.Mention, error)
        GetMentionsForUser(userID uint) ([]models.Mention, error)
        MarkMentionAsNotified(mentionID uint) error
        
        // Attachment operations
        CreateAttachment(contentID uint, contentType string, userID uint, fileName string, fileSize int64, fileType string, storagePath string, url string, isImage bool, width int, height int) (*models.Attachment, error)
        GetAttachmentsByContent(contentID uint, contentType string) ([]models.Attachment, error)
        DeleteAttachment(attachmentID uint) error
        
        // Code block operations
        CreateOrUpdateCodeBlock(contentID uint, contentType string, language string, code string, position int) (*models.CodeBlock, error)
        GetCodeBlocksByContent(contentID uint, contentType string) ([]models.CodeBlock, error)
        DeleteCodeBlock(codeBlockID uint) error
        
        // Quote operations
        CreateQuote(contentID uint, contentType string, quotedContentID uint, quotedType string, quotedUserID uint, quotedUsername string, quotedContent string, position int) (*models.Quote, error)
        GetQuotesByContent(contentID uint, contentType string) ([]models.Quote, error)
        DeleteQuote(quoteID uint) error
}

// RichTextServiceImpl implements the RichTextService interface
type RichTextServiceImpl struct {
        richTextRepo repository.RichTextRepository
        userRepo     repository.UserRepository
}

// NewRichTextService creates a new rich text service
func NewRichTextService(
        richTextRepo repository.RichTextRepository,
        userRepo repository.UserRepository,
) RichTextService {
        return &RichTextServiceImpl{
                richTextRepo: richTextRepo,
                userRepo:     userRepo,
        }
}

// CreateOrUpdateRichText creates or updates rich text content
func (s *RichTextServiceImpl) CreateOrUpdateRichText(
        contentID uint,
        contentType string,
        rawContent string,
        format models.ContentFormat,
) (*models.RichTextContent, error) {
        // Validate content type
        if contentType != "topic" && contentType != "comment" {
                return nil, errors.New("invalid content type")
        }
        
        // Check if content already exists
        existingContent, err := s.richTextRepo.GetRichTextContent(contentID, contentType)
        if err == nil && existingContent != nil {
                // Update existing content
                existingContent.RawContent = rawContent
                existingContent.Format = format
                existingContent.RenderedHTML = s.renderToHTML(rawContent, format)
                existingContent.UpdatedAt = time.Now()
                
                // Check for mentions and attachments
                hasMentions := strings.Contains(rawContent, "@")
                hasAttachments := strings.Contains(rawContent, "![") || strings.Contains(rawContent, "<img")
                
                existingContent.HasMentions = hasMentions
                existingContent.HasAttachments = hasAttachments
                
                if err := s.richTextRepo.UpdateRichTextContent(existingContent); err != nil {
                        return nil, fmt.Errorf("failed to update rich text content: %w", err)
                }
                
                // Process mentions if needed
                if hasMentions {
                        if _, err := s.ProcessMentions(contentID, contentType, rawContent); err != nil {
                                // Log error but continue
                                fmt.Printf("Error processing mentions: %v\n", err)
                        }
                }
                
                return existingContent, nil
        }
        
        // Create new content
        renderedHTML := s.renderToHTML(rawContent, format)
        hasMentions := strings.Contains(rawContent, "@")
        hasAttachments := strings.Contains(rawContent, "![") || strings.Contains(rawContent, "<img")
        
        content := &models.RichTextContent{
                ContentID:      contentID,
                ContentType:    contentType,
                Format:         format,
                RawContent:     rawContent,
                RenderedHTML:   renderedHTML,
                HasMentions:    hasMentions,
                HasAttachments: hasAttachments,
                CreatedAt:      time.Now(),
                UpdatedAt:      time.Now(),
        }
        
        if err := s.richTextRepo.CreateRichTextContent(content); err != nil {
                return nil, fmt.Errorf("failed to create rich text content: %w", err)
        }
        
        // Process mentions if needed
        if hasMentions {
                if _, err := s.ProcessMentions(contentID, contentType, rawContent); err != nil {
                        // Log error but continue
                        fmt.Printf("Error processing mentions: %v\n", err)
                }
        }
        
        return content, nil
}

// GetRichTextContent retrieves rich text content
func (s *RichTextServiceImpl) GetRichTextContent(contentID uint, contentType string) (*models.RichTextContent, error) {
        return s.richTextRepo.GetRichTextContent(contentID, contentType)
}

// DeleteRichTextContent deletes rich text content
func (s *RichTextServiceImpl) DeleteRichTextContent(contentID uint, contentType string) error {
        // Delete associated elements first
        if err := s.richTextRepo.DeleteMentionsByContent(contentID, contentType); err != nil {
                // Log error but continue
                fmt.Printf("Error deleting mentions: %v\n", err)
        }
        
        if err := s.richTextRepo.DeleteAttachmentsByContent(contentID, contentType); err != nil {
                // Log error but continue
                fmt.Printf("Error deleting attachments: %v\n", err)
        }
        
        if err := s.richTextRepo.DeleteCodeBlocksByContent(contentID, contentType); err != nil {
                // Log error but continue
                fmt.Printf("Error deleting code blocks: %v\n", err)
        }
        
        if err := s.richTextRepo.DeleteQuotesByContent(contentID, contentType); err != nil {
                // Log error but continue
                fmt.Printf("Error deleting quotes: %v\n", err)
        }
        
        // Delete the rich text content
        return s.richTextRepo.DeleteRichTextContent(contentID, contentType)
}

// ProcessMentions processes and extracts mentions from content
func (s *RichTextServiceImpl) ProcessMentions(contentID uint, contentType string, rawContent string) ([]models.Mention, error) {
        // Delete existing mentions for this content
        if err := s.richTextRepo.DeleteMentionsByContent(contentID, contentType); err != nil {
                return nil, fmt.Errorf("failed to delete existing mentions: %w", err)
        }
        
        // Extract mentions from content
        // For markdown and HTML, look for @username patterns
        r := regexp.MustCompile(`@([a-zA-Z0-9_-]+)`)
        matches := r.FindAllStringSubmatch(rawContent, -1)
        
        mentions := make([]models.Mention, 0)
        processedUsernames := make(map[string]bool)
        
        for _, match := range matches {
                if len(match) < 2 {
                        continue
                }
                
                username := match[1]
                
                // Skip if already processed this username
                if _, ok := processedUsernames[username]; ok {
                        continue
                }
                
                processedUsernames[username] = true
                
                // Find user by username
                user, err := s.userRepo.GetUserByUsername(username)
                if err != nil {
                        // User not found, skip
                        continue
                }
                
                // Create mention
                mention := models.Mention{
                        UserID:      user.ID,
                        TargetType:  contentType,
                        TargetID:    contentID,
                        MentionedBy: 0, // We don't know who created the content here, can be updated later
                        IsRead:      false,
                }
                
                if err := s.richTextRepo.CreateMention(&mention); err != nil {
                        // Log error but continue
                        fmt.Printf("Error creating mention: %v\n", err)
                        continue
                }
                
                mentions = append(mentions, mention)
        }
        
        return mentions, nil
}

// GetMentionsForUser retrieves mentions for a user
func (s *RichTextServiceImpl) GetMentionsForUser(userID uint) ([]models.Mention, error) {
        return s.richTextRepo.GetMentionsForUser(userID)
}

// MarkMentionAsNotified marks a mention as notified (read)
func (s *RichTextServiceImpl) MarkMentionAsNotified(mentionID uint) error {
        // Since Mention now uses IsRead field instead of Notified, 
        // we're adapting the method to mark it as read
        return s.richTextRepo.MarkMentionAsRead(mentionID)
}

// CreateAttachment creates a new attachment
func (s *RichTextServiceImpl) CreateAttachment(
        contentID uint,
        contentType string,
        userID uint,
        fileName string,
        fileSize int64,
        fileType string,
        storagePath string,
        url string,
        isImage bool,
        width int,
        height int,
) (*models.Attachment, error) {
        attachment := &models.Attachment{
                ContentID:    contentID,
                ContentType:  contentType,
                UploadedBy:   userID,
                UserID:       userID,  // Set both fields for compatibility
                FileName:     fileName,
                FileSize:     fileSize,
                FileType:     fileType,
                StoragePath:  storagePath,
                URL:          url,
                IsImage:      isImage,
                ImageWidth:   width,
                ImageHeight:  height,
                Width:        width,   // Set both fields for compatibility
                Height:       height,  // Set both fields for compatibility
                CreatedAt:    time.Now(),
        }
        
        if err := s.richTextRepo.CreateAttachment(attachment); err != nil {
                return nil, fmt.Errorf("failed to create attachment: %w", err)
        }
        
        // Update has_attachments flag in rich text content
        if richText, err := s.richTextRepo.GetRichTextContent(contentID, contentType); err == nil && richText != nil {
                richText.HasAttachments = true
                if err := s.richTextRepo.UpdateRichTextContent(richText); err != nil {
                        // Log error but continue
                        fmt.Printf("Error updating rich text has_attachments flag: %v\n", err)
                }
        }
        
        return attachment, nil
}

// GetAttachmentsByContent retrieves attachments by content ID and type
func (s *RichTextServiceImpl) GetAttachmentsByContent(contentID uint, contentType string) ([]models.Attachment, error) {
        return s.richTextRepo.GetAttachmentsByContent(contentID, contentType)
}

// DeleteAttachment deletes an attachment
func (s *RichTextServiceImpl) DeleteAttachment(attachmentID uint) error {
        // Get attachment to determine content ID and type
        attachment, err := s.richTextRepo.GetAttachment(attachmentID)
        if err != nil {
                return fmt.Errorf("failed to get attachment: %w", err)
        }
        
        // Delete the attachment
        if err := s.richTextRepo.DeleteAttachment(attachmentID); err != nil {
                return fmt.Errorf("failed to delete attachment: %w", err)
        }
        
        // Check if there are any remaining attachments
        attachments, err := s.richTextRepo.GetAttachmentsByContent(attachment.ContentID, attachment.ContentType)
        if err != nil {
                // Log error but continue
                fmt.Printf("Error getting remaining attachments: %v\n", err)
        } else if len(attachments) == 0 {
                // No more attachments, update has_attachments flag
                if richText, err := s.richTextRepo.GetRichTextContent(attachment.ContentID, attachment.ContentType); err == nil && richText != nil {
                        richText.HasAttachments = false
                        if err := s.richTextRepo.UpdateRichTextContent(richText); err != nil {
                                // Log error but continue
                                fmt.Printf("Error updating rich text has_attachments flag: %v\n", err)
                        }
                }
        }
        
        return nil
}

// CreateOrUpdateCodeBlock creates or updates a code block
func (s *RichTextServiceImpl) CreateOrUpdateCodeBlock(
        contentID uint,
        contentType string,
        language string,
        code string,
        position int,
) (*models.CodeBlock, error) {
        // Get existing code blocks
        codeBlocks, err := s.richTextRepo.GetCodeBlocksByContent(contentID, contentType)
        if err != nil {
                // Log error but continue
                fmt.Printf("Error getting code blocks: %v\n", err)
        }
        
        // Check if there's a code block at this position
        for _, block := range codeBlocks {
                // Note: CodeBlock doesn't have a Position field, so we'll need another way to identify
                // which block to update. For now, we'll assume the first one is the target.
                if true { // Temporary solution until we find a proper way to identify blocks
                        // Update existing code block
                        block.Language = language
                        block.Code = code
                        block.HighlightedHTML = s.highlightCode(code, language)
                        block.LineCount = len(strings.Split(code, "\n"))
                        
                        if err := s.richTextRepo.UpdateCodeBlock(&block); err != nil {
                                return nil, fmt.Errorf("failed to update code block: %w", err)
                        }
                        
                        return &block, nil
                }
        }
        
        // Create new code block
        codeBlock := &models.CodeBlock{
                ContentID:       contentID,
                ContentType:     contentType,
                Language:        language,
                Code:            code,
                HighlightedHTML: s.highlightCode(code, language),
                LineCount:       len(strings.Split(code, "\n")),
                // Position is not in the struct, so we'll add it as a comment
                // Position:    position,
                // CreatedAt is handled by GORM
        }
        
        if err := s.richTextRepo.CreateCodeBlock(codeBlock); err != nil {
                return nil, fmt.Errorf("failed to create code block: %w", err)
        }
        
        return codeBlock, nil
}

// GetCodeBlocksByContent retrieves code blocks by content ID and type
func (s *RichTextServiceImpl) GetCodeBlocksByContent(contentID uint, contentType string) ([]models.CodeBlock, error) {
        return s.richTextRepo.GetCodeBlocksByContent(contentID, contentType)
}

// DeleteCodeBlock deletes a code block
func (s *RichTextServiceImpl) DeleteCodeBlock(codeBlockID uint) error {
        return s.richTextRepo.DeleteCodeBlock(codeBlockID)
}

// CreateQuote creates a new quote
func (s *RichTextServiceImpl) CreateQuote(
        contentID uint,
        contentType string,
        quotedContentID uint,
        quotedType string,
        quotedUserID uint,
        quotedUsername string,
        quotedContent string,
        position int,
) (*models.Quote, error) {
        // Determine the appropriate quoted ID field based on type
        var quotedTopicID, quotedCommentID *uint
        if quotedType == "topic" {
                quotedTopicID = &quotedContentID
        } else if quotedType == "comment" {
                quotedCommentID = &quotedContentID
        }
        
        quote := &models.Quote{
                ContentType:     contentType,
                ContentID:       contentID,
                QuotedText:      quotedContent,
                QuotedTopicID:   quotedTopicID,
                QuotedCommentID: quotedCommentID,
                QuotedUserID:    &quotedUserID,
                Citation:        quotedUsername, // Store the username as citation
        }
        
        if err := s.richTextRepo.CreateQuote(quote); err != nil {
                return nil, fmt.Errorf("failed to create quote: %w", err)
        }
        
        return quote, nil
}

// GetQuotesByContent retrieves quotes by content ID and type
func (s *RichTextServiceImpl) GetQuotesByContent(contentID uint, contentType string) ([]models.Quote, error) {
        return s.richTextRepo.GetQuotesByContent(contentID, contentType)
}

// DeleteQuote deletes a quote
func (s *RichTextServiceImpl) DeleteQuote(quoteID uint) error {
        return s.richTextRepo.DeleteQuote(quoteID)
}



// Helper functions

// renderToHTML renders raw content to HTML based on the format
func (s *RichTextServiceImpl) renderToHTML(rawContent string, format models.ContentFormat) string {
        // In a real implementation, we'd use a library like Blackfriday for Markdown,
        // sanitize HTML input, etc.
        
        // For this example, we'll do a simplified conversion
        switch format {
        case models.FormatMarkdown:
                // Convert basic Markdown to HTML
                html := rawContent
                
                // Bold
                html = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(html, "<strong>$1</strong>")
                
                // Italic
                html = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(html, "<em>$1</em>")
                
                // Links
                html = regexp.MustCompile(`\[(.+?)\]\((.+?)\)`).ReplaceAllString(html, "<a href=\"$2\">$1</a>")
                
                // Images
                html = regexp.MustCompile(`!\[(.+?)\]\((.+?)\)`).ReplaceAllString(html, "<img src=\"$2\" alt=\"$1\">")
                
                // Headers
                html = regexp.MustCompile(`^# (.+)$`).ReplaceAllString(html, "<h1>$1</h1>")
                html = regexp.MustCompile(`^## (.+)$`).ReplaceAllString(html, "<h2>$1</h2>")
                html = regexp.MustCompile(`^### (.+)$`).ReplaceAllString(html, "<h3>$1</h3>")
                
                // Code blocks
                html = regexp.MustCompile("```(.+?)```").ReplaceAllString(html, "<pre><code>$1</code></pre>")
                
                // Inline code
                html = regexp.MustCompile("`(.+?)`").ReplaceAllString(html, "<code>$1</code>")
                
                // Paragraphs (simplified)
                html = "<p>" + html + "</p>"
                
                return html
                
        case models.FormatHTML:
                // For HTML, we'd sanitize the input to prevent XSS
                // For this example, we'll just return as-is
                return rawContent
                
        case models.FormatPlain:
                // Convert plain text to HTML (escaping and adding paragraphs)
                // Escape HTML
                html := strings.ReplaceAll(rawContent, "&", "&amp;")
                html = strings.ReplaceAll(html, "<", "&lt;")
                html = strings.ReplaceAll(html, ">", "&gt;")
                
                // Replace newlines with paragraph breaks
                html = "<p>" + strings.ReplaceAll(html, "\n\n", "</p><p>") + "</p>"
                
                // Replace single newlines with line breaks
                html = strings.ReplaceAll(html, "\n", "<br>")
                
                return html
                
        default:
                return rawContent
        }
}

// highlightCode performs syntax highlighting for code blocks
func (s *RichTextServiceImpl) highlightCode(code string, language string) string {
        // In a real implementation, we'd use a syntax highlighting library
        // For this example, we'll just wrap in a pre and code tags with a language class
        
        // Escape HTML in the code
        escapedCode := strings.ReplaceAll(code, "&", "&amp;")
        escapedCode = strings.ReplaceAll(escapedCode, "<", "&lt;")
        escapedCode = strings.ReplaceAll(escapedCode, ">", "&gt;")
        
        return fmt.Sprintf("<pre><code class=\"language-%s\">%s</code></pre>", language, escapedCode)
}