package service

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/repository"
)

// BookmarkService defines interface for bookmark-related business logic
type BookmarkService interface {
        CreateBookmark(bookmark *models.Bookmark) error
        GetBookmarks(userID, bookID uint) ([]models.Bookmark, error)
        GetBookmarksByChapter(userID, chapterID uint) ([]models.Bookmark, error)
        GetBookmarkByID(id, userID uint) (*models.Bookmark, error)
        UpdateBookmark(bookmark *models.Bookmark) error
        DeleteBookmark(id, userID uint) error
        ShareBookmark(bookmarkID, fromUserID, toUserID uint) error
}

// BookmarkServiceImpl implements the BookmarkService interface
type BookmarkServiceImpl struct {
        bookmarkRepo repository.BookmarkRepository
}

// NewBookmarkService creates a new bookmark service instance
func NewBookmarkService(bookmarkRepo repository.BookmarkRepository) BookmarkService {
        return &BookmarkServiceImpl{
                bookmarkRepo: bookmarkRepo,
        }
}

// CreateBookmark creates a new bookmark
func (s *BookmarkServiceImpl) CreateBookmark(bookmark *models.Bookmark) error {
        return s.bookmarkRepo.CreateBookmark(bookmark)
}

// GetBookmarks retrieves all bookmarks for a user and book
func (s *BookmarkServiceImpl) GetBookmarks(userID, bookID uint) ([]models.Bookmark, error) {
        return s.bookmarkRepo.GetBookmarks(userID, bookID)
}

// GetBookmarksByChapter retrieves bookmarks for a specific chapter
func (s *BookmarkServiceImpl) GetBookmarksByChapter(userID, chapterID uint) ([]models.Bookmark, error) {
        return s.bookmarkRepo.GetBookmarksByChapter(userID, chapterID)
}

// GetBookmarkByID retrieves a bookmark by ID
func (s *BookmarkServiceImpl) GetBookmarkByID(id, userID uint) (*models.Bookmark, error) {
        return s.bookmarkRepo.GetBookmarkByID(id, userID)
}

// UpdateBookmark updates an existing bookmark
func (s *BookmarkServiceImpl) UpdateBookmark(bookmark *models.Bookmark) error {
        return s.bookmarkRepo.UpdateBookmark(bookmark)
}

// DeleteBookmark deletes a bookmark
func (s *BookmarkServiceImpl) DeleteBookmark(id, userID uint) error {
        return s.bookmarkRepo.DeleteBookmark(id, userID)
}

// ShareBookmark shares a bookmark with another user
func (s *BookmarkServiceImpl) ShareBookmark(bookmarkID, fromUserID, toUserID uint) error {
        return s.bookmarkRepo.ShareBookmark(bookmarkID, fromUserID, toUserID)
}