package repository

import (
        "errors"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "gorm.io/gorm"
)

// BookmarkRepository defines the interface for bookmark operations
type BookmarkRepository interface {
        CreateBookmark(bookmark *models.Bookmark) error
        GetBookmarks(userID, bookID uint) ([]models.Bookmark, error)
        GetBookmarksByChapter(userID, chapterID uint) ([]models.Bookmark, error)
        GetBookmarkByID(id, userID uint) (*models.Bookmark, error)
        UpdateBookmark(bookmark *models.Bookmark) error
        DeleteBookmark(id, userID uint) error
        ShareBookmark(bookmarkID, fromUserID, toUserID uint) error
}

// GormBookmarkRepository implements the BookmarkRepository interface with GORM
type GormBookmarkRepository struct {
        db *gorm.DB
}

// NewGormBookmarkRepository creates a new bookmark repository instance
func NewGormBookmarkRepository(db *gorm.DB) *GormBookmarkRepository {
        return &GormBookmarkRepository{db: db}
}

// CreateBookmark creates a new bookmark
func (r *GormBookmarkRepository) CreateBookmark(bookmark *models.Bookmark) error {
        return r.db.Create(bookmark).Error
}

// GetBookmarks retrieves all bookmarks for a user and book
func (r *GormBookmarkRepository) GetBookmarks(userID, bookID uint) ([]models.Bookmark, error) {
        var bookmarks []models.Bookmark
        result := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Order("created_at DESC").Find(&bookmarks)
        return bookmarks, result.Error
}

// GetBookmarksByChapter retrieves bookmarks for a specific chapter
func (r *GormBookmarkRepository) GetBookmarksByChapter(userID, chapterID uint) ([]models.Bookmark, error) {
        var bookmarks []models.Bookmark
        result := r.db.Where("user_id = ? AND chapter_id = ?", userID, chapterID).Order("position ASC").Find(&bookmarks)
        return bookmarks, result.Error
}

// GetBookmarkByID retrieves a bookmark by ID, ensuring it belongs to the user
func (r *GormBookmarkRepository) GetBookmarkByID(id, userID uint) (*models.Bookmark, error) {
        var bookmark models.Bookmark
        result := r.db.Where("id = ? AND user_id = ?", id, userID).First(&bookmark)
        if result.Error != nil {
                if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                        return nil, errors.New("bookmark not found")
                }
                return nil, result.Error
        }
        return &bookmark, nil
}

// UpdateBookmark updates an existing bookmark
func (r *GormBookmarkRepository) UpdateBookmark(bookmark *models.Bookmark) error {
        // Verify the bookmark exists and belongs to the user
        _, err := r.GetBookmarkByID(bookmark.ID, bookmark.UserID)
        if err != nil {
                return err
        }
        
        return r.db.Save(bookmark).Error
}

// DeleteBookmark deletes a bookmark
func (r *GormBookmarkRepository) DeleteBookmark(id, userID uint) error {
        result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Bookmark{})
        if result.RowsAffected == 0 {
                return errors.New("bookmark not found or does not belong to user")
        }
        return result.Error
}

// ShareBookmark shares a bookmark with another user
func (r *GormBookmarkRepository) ShareBookmark(bookmarkID, fromUserID, toUserID uint) error {
        // Retrieve the original bookmark
        originalBookmark, err := r.GetBookmarkByID(bookmarkID, fromUserID)
        if err != nil {
                return err
        }
        
        // Create a new bookmark for the target user
        sharedBookmark := models.Bookmark{
                UserID:      toUserID,
                BookID:      originalBookmark.BookID,
                ChapterID:   originalBookmark.ChapterID,
                SectionID:   originalBookmark.SectionID,
                Title:       originalBookmark.Title,
                Description: originalBookmark.Description + " (Shared by user " + string(fromUserID) + ")",
                Color:       originalBookmark.Color,
                Position:    originalBookmark.Position,
        }
        
        return r.db.Create(&sharedBookmark).Error
}