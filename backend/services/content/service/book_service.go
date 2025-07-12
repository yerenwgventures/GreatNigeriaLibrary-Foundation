package service

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/repository"
)

// BookService interface defines the operations for book content
type BookService interface {
        GetAllBooks(includeUnpublished bool) ([]models.Book, error)
        GetAllBooksWithChapters(includeUnpublished bool) ([]models.Book, error)
        GetBookByID(id uint) (*models.Book, error)
        GetBookWithChapters(id uint) (*models.Book, error)
        GetBookChapters(bookID uint, includeUnpublished bool) ([]models.BookChapter, error)
        GetChaptersByBookID(bookID uint, includeUnpublished bool) ([]models.BookChapter, error)
        GetChapter(id uint) (*models.BookChapter, error)
        GetChapterByID(id uint) (*models.BookChapter, error)
        GetChapterSections(chapterID uint, includeUnpublished bool) ([]models.BookSection, error)
        GetChapterWithSections(chapterID uint) (*models.BookChapter, error)
        GetSection(id uint) (*models.BookSection, error)
        GetSectionByID(id uint) (*models.BookSection, error)
        GetSectionWithSubsections(sectionID uint) (*models.BookSection, error)
        GetRenderedSection(sectionID uint) (string, error)
        
        // Subsection operations
        GetSubsection(id uint) (*models.BookSubsection, error)
        GetSubsectionByID(id uint) (*models.BookSubsection, error)
        GetSubsectionsBySectionID(sectionID uint, includeUnpublished bool) ([]models.BookSubsection, error)
        GetRenderedSubsection(subsectionID uint) (string, error)
        
        GetRelatedForumTopics(sectionID uint) ([]models.ForumTopic, error)
        GetForumTopicsForSection(sectionID uint) ([]models.ForumTopic, error)
        SaveReadingProgress(progress *models.ReadingProgress) error
        SearchBooks(query string, userID uint) ([]models.Book, error)
        GetRecommendations(userID uint, contentID uint) ([]models.Book, error)
        // Add other required methods from existing implementation
        GetBookFrontMatter(bookID uint, includeUnpublished bool) ([]models.BookFrontMatter, error)
        GetBookBackMatter(bookID uint, includeUnpublished bool) (*models.BookBackMatter, error)
        GetActionStepsForSection(sectionID uint) ([]models.ActionStep, error)
        TrackUserProgress(userID, bookID, sectionID uint, isRead bool, progress float64) error
        GetUserBookProgress(userID, bookID uint) (float64, error)
        CreateBookmark(userID, bookID, sectionID uint, title, description string) (*models.BookBookmark, error)
        GetUserBookmarks(userID, bookID uint) ([]models.BookBookmark, error)
        DeleteBookmark(id uint) error
        CreateNote(userID, bookID, sectionID uint, content string, isPrivate bool) (*models.BookNote, error)
        GetUserNotes(userID, bookID uint) ([]models.BookNote, error)
        UpdateNote(id uint, content string, isPrivate bool) (*models.BookNote, error)
        DeleteNote(id uint) error
        SubmitFeedback(userID, bookID, sectionID uint, feedbackType string, rating int, content string, metadata map[string]interface{}) (*models.BookFeedback, error)
        GetFeedbackStats(bookID uint) (map[string]interface{}, error)
        GetBookContent(bookID uint, includeUnpublished bool) (map[string]interface{}, error)
        // Interactive elements
        GetInteractiveElementsBySectionID(sectionID uint) ([]models.InteractiveElement, error)
}

// BookServiceImpl handles operations related to books and their content
type BookServiceImpl struct {
        repository repository.BookRepository
        contentRenderer ContentRenderer // Changed from *ContentRenderer to ContentRenderer
}

// NewBookService creates a new book service
func NewBookService(repo repository.BookRepository, renderer ContentRenderer) BookService {
        return &BookServiceImpl{
                repository: repo,
                contentRenderer: renderer,
        }
}

// GetAllBooks retrieves all books
func (s *BookServiceImpl) GetAllBooks(includeUnpublished bool) ([]models.Book, error) {
        return s.repository.GetAllBooks(includeUnpublished)
}

// GetAllBooksWithChapters retrieves all books with their chapters
func (s *BookServiceImpl) GetAllBooksWithChapters(includeUnpublished bool) ([]models.Book, error) {
        return s.repository.GetAllBooksWithChapters(includeUnpublished)
}

// GetBookByID retrieves a book by its ID
func (s *BookServiceImpl) GetBookByID(id uint) (*models.Book, error) {
        return s.repository.GetBookByID(id)
}

// GetBookWithChapters retrieves a book with its chapters
func (s *BookServiceImpl) GetBookWithChapters(id uint) (*models.Book, error) {
        return s.repository.GetBookWithChapters(id)
}

// GetBookFrontMatter retrieves all front matter for a book
func (s *BookServiceImpl) GetBookFrontMatter(bookID uint, includeUnpublished bool) ([]models.BookFrontMatter, error) {
        return s.repository.GetFrontMatterByBookID(bookID, includeUnpublished)
}

// GetBookBackMatter retrieves the back matter for a book
func (s *BookServiceImpl) GetBookBackMatter(bookID uint, includeUnpublished bool) (*models.BookBackMatter, error) {
        return s.repository.GetBackMatterByBookID(bookID, includeUnpublished)
}

// GetChaptersByBookID retrieves all chapters for a book
func (s *BookServiceImpl) GetChaptersByBookID(bookID uint, includeUnpublished bool) ([]models.BookChapter, error) {
        return s.repository.GetChaptersByBookID(bookID, includeUnpublished)
}

// GetBookChapters is an alias for GetChaptersByBookID for compatibility
func (s *BookServiceImpl) GetBookChapters(bookID uint, includeUnpublished bool) ([]models.BookChapter, error) {
        return s.GetChaptersByBookID(bookID, includeUnpublished)
}

// GetChapterByID retrieves a chapter by its ID
func (s *BookServiceImpl) GetChapterByID(id uint) (*models.BookChapter, error) {
        return s.repository.GetChapterByID(id)
}

// GetChapter is an alias for GetChapterByID for compatibility
func (s *BookServiceImpl) GetChapter(id uint) (*models.BookChapter, error) {
        return s.GetChapterByID(id)
}

// GetChapterWithSections gets a chapter with all its sections
func (s *BookServiceImpl) GetChapterWithSections(chapterID uint) (*models.BookChapter, error) {
        return s.repository.GetChapterWithSections(chapterID)
}

// GetSectionByID retrieves a section by its ID
func (s *BookServiceImpl) GetSectionByID(id uint) (*models.BookSection, error) {
        return s.repository.GetSectionByID(id)
}

// GetSection is an alias for GetSectionByID for compatibility
func (s *BookServiceImpl) GetSection(id uint) (*models.BookSection, error) {
        return s.GetSectionByID(id)
}

// GetSectionWithSubsections gets a section with all its subsections
func (s *BookServiceImpl) GetSectionWithSubsections(sectionID uint) (*models.BookSection, error) {
        // Get the section
        section, err := s.GetSectionByID(sectionID)
        if err != nil {
                return nil, err
        }

        // Get subsections for the section
        subsections, err := s.GetSubsectionsBySectionID(sectionID, false) // Only published
        if err != nil {
                return nil, err
        }

        // Assign subsections to the section
        section.Subsections = subsections
        
        return section, nil
}

// GetSubsection is an alias for GetSubsectionByID for compatibility
func (s *BookServiceImpl) GetSubsection(id uint) (*models.BookSubsection, error) {
        return s.GetSubsectionByID(id)
}

// GetSubsectionByID retrieves a subsection by its ID
func (s *BookServiceImpl) GetSubsectionByID(id uint) (*models.BookSubsection, error) {
        return s.repository.GetSubsectionByID(id)
}

// GetSubsectionsBySectionID retrieves all subsections for a section
func (s *BookServiceImpl) GetSubsectionsBySectionID(sectionID uint, includeUnpublished bool) ([]models.BookSubsection, error) {
        return s.repository.GetSubsectionsBySectionID(sectionID, includeUnpublished)
}

// GetRenderedSubsection returns the rendered HTML content for a subsection
func (s *BookServiceImpl) GetRenderedSubsection(subsectionID uint) (string, error) {
        subsection, err := s.repository.GetSubsectionByID(subsectionID)
        if err != nil {
                return "", err
        }
        
        // Render the content to HTML using the content renderer
        renderedContent, err := s.contentRenderer.RenderMarkdown(subsection.Content)
        if err != nil {
                return "", err
        }
        
        return renderedContent, nil
}

// GetRenderedSection returns the rendered HTML content for a section
func (s *BookServiceImpl) GetRenderedSection(sectionID uint) (string, error) {
        section, err := s.repository.GetSectionByID(sectionID)
        if err != nil {
                return "", err
        }
        
        // Render the content to HTML using the content renderer
        // Use RenderMarkdown instead of RenderContent since that's in the interface
        renderedContent, err := s.contentRenderer.RenderMarkdown(section.Content)
        if err != nil {
                return "", err
        }
        
        return renderedContent, nil
}

// GetRelatedForumTopics retrieves all forum topics for a section
func (s *BookServiceImpl) GetRelatedForumTopics(sectionID uint) ([]models.ForumTopic, error) {
        return s.repository.GetForumTopicsBySectionID(sectionID)
}

// GetForumTopicsForSection is an alias for GetRelatedForumTopics
func (s *BookServiceImpl) GetForumTopicsForSection(sectionID uint) ([]models.ForumTopic, error) {
        return s.GetRelatedForumTopics(sectionID)
}

// GetChapterSections retrieves sections for a chapter
func (s *BookServiceImpl) GetChapterSections(chapterID uint, includeUnpublished bool) ([]models.BookSection, error) {
        return s.repository.GetSectionsByChapterID(chapterID, includeUnpublished)
}

// GetActionStepsForSection retrieves all action steps for a section
func (s *BookServiceImpl) GetActionStepsForSection(sectionID uint) ([]models.ActionStep, error) {
        return s.repository.GetActionStepsBySectionID(sectionID)
}

// SearchBooks searches for books matching the query
func (s *BookServiceImpl) SearchBooks(query string, userID uint) ([]models.Book, error) {
        // Pass empty tag list and default limit of 10
        return s.repository.SearchBooks(query, []string{}, 10)
}

// GetRecommendations gets book recommendations for a user
func (s *BookServiceImpl) GetRecommendations(userID uint, contentID uint) ([]models.Book, error) {
        // Use the repository method, with a standard limit of 5 recommendations
        return s.repository.GetRecommendations(userID, 5)
}

// TrackUserProgress tracks a user's progress in reading a book section
func (s *BookServiceImpl) TrackUserProgress(userID, bookID, sectionID uint, isRead bool, progress float64) error {
        bookProgress := &models.BookProgress{
                UserID:    userID,
                BookID:    bookID,
                SectionID: sectionID,
                IsRead:    isRead,
                Progress:  progress,
        }
        return s.repository.CreateOrUpdateBookProgress(bookProgress)
}

// GetUserBookProgress gets a user's overall progress for a book
func (s *BookServiceImpl) GetUserBookProgress(userID, bookID uint) (float64, error) {
        return s.repository.GetBookProgressPercentage(userID, bookID)
}

// SaveReadingProgress saves a user's reading progress
func (s *BookServiceImpl) SaveReadingProgress(progress *models.ReadingProgress) error {
        // First check if the reading progress already exists
        existingProgress, err := s.repository.GetReadingProgress(progress.UserID, progress.BookID)
        if err != nil {
                return err
        }
        
        if existingProgress != nil {
                // Update existing progress
                existingProgress.LastReadAt = progress.LastReadAt
                existingProgress.TimeSpent = progress.TimeSpent
                existingProgress.SessionCount = progress.SessionCount
                existingProgress.StreakDays = progress.StreakDays
                
                // Save the updated progress
                return s.repository.UpdateReadingProgress(existingProgress)
        }
        
        // Create new progress
        return s.repository.CreateReadingProgress(progress)
}

// CreateBookmark creates a bookmark for a user
func (s *BookServiceImpl) CreateBookmark(userID, bookID, sectionID uint, title, description string) (*models.BookBookmark, error) {
        bookmark := &models.BookBookmark{
                UserID:      userID,
                BookID:      bookID,
                SectionID:   sectionID,
                Title:       title,
                Description: description,
        }
        err := s.repository.CreateBookmark(bookmark)
        if err != nil {
                return nil, err
        }
        return bookmark, nil
}

// GetUserBookmarks gets all bookmarks for a user in a book
func (s *BookServiceImpl) GetUserBookmarks(userID, bookID uint) ([]models.BookBookmark, error) {
        return s.repository.GetBookmarksByUserAndBook(userID, bookID)
}

// DeleteBookmark deletes a bookmark
func (s *BookServiceImpl) DeleteBookmark(id uint) error {
        return s.repository.DeleteBookmark(id)
}

// CreateNote creates a note for a user
func (s *BookServiceImpl) CreateNote(userID, bookID, sectionID uint, content string, isPrivate bool) (*models.BookNote, error) {
        note := &models.BookNote{
                UserID:    userID,
                BookID:    bookID,
                SectionID: sectionID,
                Content:   content,
                // IsPrivate field removed as it's not in the model
        }
        err := s.repository.CreateNote(note)
        if err != nil {
                return nil, err
        }
        return note, nil
}

// GetUserNotes gets all notes for a user in a book
func (s *BookServiceImpl) GetUserNotes(userID, bookID uint) ([]models.BookNote, error) {
        return s.repository.GetNotesByUserAndBook(userID, bookID)
}

// UpdateNote updates a note
func (s *BookServiceImpl) UpdateNote(id uint, content string, isPrivate bool) (*models.BookNote, error) {
        note, err := s.repository.GetNoteByID(id)
        if err != nil {
                return nil, err
        }
        
        note.Content = content
        // IsPrivate field removed as it's not in the model
        
        err = s.repository.UpdateNote(note)
        if err != nil {
                return nil, err
        }
        
        return note, nil
}

// DeleteNote deletes a note
func (s *BookServiceImpl) DeleteNote(id uint) error {
        return s.repository.DeleteNote(id)
}

// SubmitFeedback submits feedback for a book or section
func (s *BookServiceImpl) SubmitFeedback(userID, bookID, sectionID uint, feedbackType string, rating int, content string, metadata map[string]interface{}) (*models.BookFeedback, error) {
        feedback := &models.BookFeedback{
                UserID:       userID,
                BookID:       bookID,
                SectionID:    sectionID,
                Rating:       rating,
                Comment:      content,
                // FeedbackType, Content, and Metadata fields removed as they're not in the model
        }
        
        err := s.repository.CreateFeedback(feedback)
        if err != nil {
                return nil, err
        }
        
        return feedback, nil
}

// GetFeedbackStats gets aggregated feedback statistics for a book
func (s *BookServiceImpl) GetFeedbackStats(bookID uint) (map[string]interface{}, error) {
        return s.repository.GetFeedbackStats(bookID)
}

// GetInteractiveElementsBySectionID retrieves all interactive elements for a section
func (s *BookServiceImpl) GetInteractiveElementsBySectionID(sectionID uint) ([]models.InteractiveElement, error) {
        return s.repository.GetInteractiveElementsBySectionID(sectionID)
}

// GetBookContent retrieves a complete book with all its content
func (s *BookServiceImpl) GetBookContent(bookID uint, includeUnpublished bool) (map[string]interface{}, error) {
        // Get the book
        book, err := s.repository.GetBookByID(bookID)
        if err != nil {
                return nil, err
        }
        
        // Get front matter
        frontMatter, err := s.repository.GetFrontMatterByBookID(bookID, includeUnpublished)
        if err != nil {
                return nil, err
        }
        
        // Get chapters
        chapters, err := s.repository.GetChaptersByBookID(bookID, includeUnpublished)
        if err != nil {
                return nil, err
        }
        
        // For each chapter, get sections
        chapterContent := make([]map[string]interface{}, 0)
        for _, chapter := range chapters {
                sections, err := s.repository.GetSectionsByChapterID(chapter.ID, includeUnpublished)
                if err != nil {
                        return nil, err
                }
                
                // For each section, get forum topics and action steps
                sectionContent := make([]map[string]interface{}, 0)
                for _, section := range sections {
                        forumTopics, err := s.repository.GetForumTopicsBySectionID(section.ID)
                        if err != nil {
                                return nil, err
                        }
                        
                        actionSteps, err := s.repository.GetActionStepsBySectionID(section.ID)
                        if err != nil {
                                return nil, err
                        }
                        
                        interactiveElements, err := s.repository.GetInteractiveElementsBySectionID(section.ID)
                        if err != nil {
                                return nil, err
                        }
                        
                        // Get subsections for Book 3 (Comprehensive Edition)
                        var subsections []models.BookSubsection
                        if bookID == 3 { // Book 3 is the only one with subsections
                                subsections, err = s.repository.GetSubsectionsBySectionID(section.ID, includeUnpublished)
                                if err != nil {
                                        return nil, err
                                }
                        }
                        
                        sectionContent = append(sectionContent, map[string]interface{}{
                                "section":              section,
                                "forum_topics":         forumTopics,
                                "action_steps":         actionSteps,
                                "interactive_elements": interactiveElements,
                                "subsections":          subsections,
                        })
                }
                
                chapterContent = append(chapterContent, map[string]interface{}{
                        "chapter":  chapter,
                        "sections": sectionContent,
                })
        }
        
        // Compile the full book content
        bookContent := map[string]interface{}{
                "book":         book,
                "front_matter": frontMatter,
                "chapters":     chapterContent,
        }
        
        return bookContent, nil
}