package repository

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/models"
        "gorm.io/gorm"
)

// RichTextRepository defines the interface for rich text operations
type RichTextRepository interface {
        // Rich text content operations
        CreateRichTextContent(content *models.RichTextContent) error
        GetRichTextContent(contentID uint, contentType string) (*models.RichTextContent, error)
        UpdateRichTextContent(content *models.RichTextContent) error
        DeleteRichTextContent(contentID uint, contentType string) error
        
        // Mention operations
        CreateMention(mention *models.Mention) error
        GetMentionsByContent(contentID uint, contentType string) ([]models.Mention, error)
        GetMentionsForUser(userID uint) ([]models.Mention, error)
        MarkMentionAsNotified(mentionID uint) error
        MarkMentionAsRead(mentionID uint) error // New method for updated field name
        DeleteMentionsByContent(contentID uint, contentType string) error
        
        // Attachment operations
        CreateAttachment(attachment *models.Attachment) error
        GetAttachmentsByContent(contentID uint, contentType string) ([]models.Attachment, error)
        GetAttachment(attachmentID uint) (*models.Attachment, error)
        DeleteAttachment(attachmentID uint) error
        DeleteAttachmentsByContent(contentID uint, contentType string) error
        
        // Code block operations
        CreateCodeBlock(codeBlock *models.CodeBlock) error
        GetCodeBlocksByContent(contentID uint, contentType string) ([]models.CodeBlock, error)
        UpdateCodeBlock(codeBlock *models.CodeBlock) error
        DeleteCodeBlock(codeBlockID uint) error
        DeleteCodeBlocksByContent(contentID uint, contentType string) error
        
        // Quote operations
        CreateQuote(quote *models.Quote) error
        GetQuotesByContent(contentID uint, contentType string) ([]models.Quote, error)
        GetQuotesByQuotedContent(quotedContentID uint, quotedType string) ([]models.Quote, error)
        DeleteQuote(quoteID uint) error
        DeleteQuotesByContent(contentID uint, contentType string) error
}

// GormRichTextRepository implements the RichTextRepository interface
type GormRichTextRepository struct {
        db *gorm.DB
}

// NewGormRichTextRepository creates a new rich text repository
func NewGormRichTextRepository(db *gorm.DB) *GormRichTextRepository {
        return &GormRichTextRepository{db: db}
}

// CreateRichTextContent creates a new rich text content
func (r *GormRichTextRepository) CreateRichTextContent(content *models.RichTextContent) error {
        return r.db.Create(content).Error
}

// GetRichTextContent retrieves rich text content by content ID and type
func (r *GormRichTextRepository) GetRichTextContent(contentID uint, contentType string) (*models.RichTextContent, error) {
        var content models.RichTextContent
        result := r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).First(&content)
        if result.Error != nil {
                return nil, result.Error
        }
        return &content, nil
}

// UpdateRichTextContent updates rich text content
func (r *GormRichTextRepository) UpdateRichTextContent(content *models.RichTextContent) error {
        return r.db.Save(content).Error
}

// DeleteRichTextContent deletes rich text content
func (r *GormRichTextRepository) DeleteRichTextContent(contentID uint, contentType string) error {
        return r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).Delete(&models.RichTextContent{}).Error
}

// CreateMention creates a new mention
func (r *GormRichTextRepository) CreateMention(mention *models.Mention) error {
        return r.db.Create(mention).Error
}

// GetMentionsByContent retrieves mentions by content ID and type
func (r *GormRichTextRepository) GetMentionsByContent(contentID uint, contentType string) ([]models.Mention, error) {
        var mentions []models.Mention
        result := r.db.Where("target_id = ? AND target_type = ?", contentID, contentType).Find(&mentions)
        return mentions, result.Error
}

// GetMentionsForUser retrieves mentions for a user
func (r *GormRichTextRepository) GetMentionsForUser(userID uint) ([]models.Mention, error) {
        var mentions []models.Mention
        result := r.db.Where("user_id = ? AND is_read = ?", userID, false).Find(&mentions)
        return mentions, result.Error
}

// MarkMentionAsNotified marks a mention as notified
func (r *GormRichTextRepository) MarkMentionAsNotified(mentionID uint) error {
        return r.db.Model(&models.Mention{}).Where("id = ?", mentionID).Update("notified", true).Error
}

// MarkMentionAsRead marks a mention as read
func (r *GormRichTextRepository) MarkMentionAsRead(mentionID uint) error {
        return r.db.Model(&models.Mention{}).Where("id = ?", mentionID).Update("is_read", true).Error
}

// DeleteMentionsByContent deletes mentions by content ID and type
func (r *GormRichTextRepository) DeleteMentionsByContent(contentID uint, contentType string) error {
        return r.db.Where("target_id = ? AND target_type = ?", contentID, contentType).Delete(&models.Mention{}).Error
}

// CreateAttachment creates a new attachment
func (r *GormRichTextRepository) CreateAttachment(attachment *models.Attachment) error {
        return r.db.Create(attachment).Error
}

// GetAttachmentsByContent retrieves attachments by content ID and type
func (r *GormRichTextRepository) GetAttachmentsByContent(contentID uint, contentType string) ([]models.Attachment, error) {
        var attachments []models.Attachment
        result := r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).Find(&attachments)
        return attachments, result.Error
}

// GetAttachment retrieves an attachment by ID
func (r *GormRichTextRepository) GetAttachment(attachmentID uint) (*models.Attachment, error) {
        var attachment models.Attachment
        result := r.db.First(&attachment, attachmentID)
        if result.Error != nil {
                return nil, result.Error
        }
        return &attachment, nil
}

// DeleteAttachment deletes an attachment
func (r *GormRichTextRepository) DeleteAttachment(attachmentID uint) error {
        return r.db.Delete(&models.Attachment{}, attachmentID).Error
}

// DeleteAttachmentsByContent deletes attachments by content ID and type
func (r *GormRichTextRepository) DeleteAttachmentsByContent(contentID uint, contentType string) error {
        return r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).Delete(&models.Attachment{}).Error
}

// CreateCodeBlock creates a new code block
func (r *GormRichTextRepository) CreateCodeBlock(codeBlock *models.CodeBlock) error {
        return r.db.Create(codeBlock).Error
}

// GetCodeBlocksByContent retrieves code blocks by content ID and type
func (r *GormRichTextRepository) GetCodeBlocksByContent(contentID uint, contentType string) ([]models.CodeBlock, error) {
        var codeBlocks []models.CodeBlock
        result := r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).Order("position").Find(&codeBlocks)
        return codeBlocks, result.Error
}

// UpdateCodeBlock updates a code block
func (r *GormRichTextRepository) UpdateCodeBlock(codeBlock *models.CodeBlock) error {
        return r.db.Save(codeBlock).Error
}

// DeleteCodeBlock deletes a code block
func (r *GormRichTextRepository) DeleteCodeBlock(codeBlockID uint) error {
        return r.db.Delete(&models.CodeBlock{}, codeBlockID).Error
}

// DeleteCodeBlocksByContent deletes code blocks by content ID and type
func (r *GormRichTextRepository) DeleteCodeBlocksByContent(contentID uint, contentType string) error {
        return r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).Delete(&models.CodeBlock{}).Error
}

// CreateQuote creates a new quote
func (r *GormRichTextRepository) CreateQuote(quote *models.Quote) error {
        return r.db.Create(quote).Error
}

// GetQuotesByContent retrieves quotes by content ID and type
func (r *GormRichTextRepository) GetQuotesByContent(contentID uint, contentType string) ([]models.Quote, error) {
        var quotes []models.Quote
        result := r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).Order("position").Find(&quotes)
        return quotes, result.Error
}

// GetQuotesByQuotedContent retrieves quotes by quoted content ID and type
func (r *GormRichTextRepository) GetQuotesByQuotedContent(quotedContentID uint, quotedType string) ([]models.Quote, error) {
        var quotes []models.Quote
        result := r.db.Where("quoted_content_id = ? AND quoted_type = ?", quotedContentID, quotedType).Find(&quotes)
        return quotes, result.Error
}

// DeleteQuote deletes a quote
func (r *GormRichTextRepository) DeleteQuote(quoteID uint) error {
        return r.db.Delete(&models.Quote{}, quoteID).Error
}

// DeleteQuotesByContent deletes quotes by content ID and type
func (r *GormRichTextRepository) DeleteQuotesByContent(contentID uint, contentType string) error {
        return r.db.Where("content_id = ? AND content_type = ?", contentID, contentType).Delete(&models.Quote{}).Error
}