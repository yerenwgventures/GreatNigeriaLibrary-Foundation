package repository

import (
        "encoding/json"
        "errors"
        "fmt"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/models"
        "gorm.io/gorm"
)

// BookRepository defines the interface for book data operations
type BookRepository interface {
        // Book operations
        CreateBook(book *models.Book) error
        GetBookByID(id uint) (*models.Book, error)
        GetBookWithChapters(id uint) (*models.Book, error)
        GetAllBooks(includeUnpublished bool) ([]models.Book, error)
        GetAllBooksWithChapters(includeUnpublished bool) ([]models.Book, error)
        UpdateBook(book *models.Book) error
        DeleteBook(id uint) error

        // Chapter operations
        CreateChapter(chapter *models.BookChapter) error
        GetChapterByID(id uint) (*models.BookChapter, error)
        GetChaptersByBookID(bookID uint, includeUnpublished bool) ([]models.BookChapter, error)
        UpdateChapter(chapter *models.BookChapter) error
        DeleteChapter(id uint) error

        // Subsection operations
        CreateSubsection(subsection *models.BookSubsection) error
        GetSubsectionByID(id uint) (*models.BookSubsection, error)
        GetSubsectionsBySectionID(sectionID uint, includeUnpublished bool) ([]models.BookSubsection, error)
        UpdateSubsection(subsection *models.BookSubsection) error
        DeleteSubsection(id uint) error
        
        // Section operations
        CreateSection(section *models.BookSection) error
        GetSectionByID(id uint) (*models.BookSection, error)
        GetSectionsByChapterID(chapterID uint, includeUnpublished bool) ([]models.BookSection, error)
        UpdateSection(section *models.BookSection) error
        DeleteSection(id uint) error

        // Front matter operations
        CreateFrontMatter(frontMatter *models.BookFrontMatter) error
        GetFrontMatterByID(id uint) (*models.BookFrontMatter, error)
        GetFrontMatterByBookID(bookID uint, includeUnpublished bool) ([]models.BookFrontMatter, error)
        UpdateFrontMatter(frontMatter *models.BookFrontMatter) error
        DeleteFrontMatter(id uint) error
        
        // Back matter operations
        GetBackMatterByBookID(bookID uint, includeUnpublished bool) (*models.BookBackMatter, error)
        UpdateBackMatter(backMatter *models.BookBackMatter) error
        DeleteBackMatter(id uint) error

        // Forum topic operations
        CreateForumTopic(topic *models.ForumTopic) error
        GetForumTopicByID(id uint) (*models.ForumTopic, error)
        GetForumTopicsBySectionID(sectionID uint) ([]models.ForumTopic, error)
        UpdateForumTopic(topic *models.ForumTopic) error
        DeleteForumTopic(id uint) error
        
        // Interactive element operations
        CreateInteractiveElement(element *models.InteractiveElement) error
        GetInteractiveElementByID(id uint) (*models.InteractiveElement, error)
        GetInteractiveElementsBySectionID(sectionID uint) ([]models.InteractiveElement, error)
        UpdateInteractiveElement(element *models.InteractiveElement) error
        DeleteInteractiveElement(id uint) error

        // Action step operations
        CreateActionStep(step *models.ActionStep) error
        GetActionStepByID(id uint) (*models.ActionStep, error)
        GetActionStepsBySectionID(sectionID uint) ([]models.ActionStep, error)
        UpdateActionStep(step *models.ActionStep) error
        DeleteActionStep(id uint) error

        // User progress operations
        CreateOrUpdateBookProgress(progress *models.BookProgress) error
        GetBookProgressByUserAndBook(userID, bookID uint) ([]models.BookProgress, error)
        GetBookProgressPercentage(userID, bookID uint) (float64, error)

        // Bookmark operations
        CreateBookmark(bookmark *models.BookBookmark) error
        GetBookmarkByID(id uint) (*models.BookBookmark, error)
        GetBookmarksByUserAndBook(userID, bookID uint) ([]models.BookBookmark, error)
        DeleteBookmark(id uint) error

        // Note operations
        CreateNote(note *models.BookNote) error
        GetNoteByID(id uint) (*models.BookNote, error)
        GetNotesByUserAndBook(userID, bookID uint) ([]models.BookNote, error)
        UpdateNote(note *models.BookNote) error
        DeleteNote(id uint) error

        // Feedback operations
        CreateFeedback(feedback *models.BookFeedback) error
        GetFeedbackByID(id uint) (*models.BookFeedback, error)
        GetFeedbackByUserAndBook(userID, bookID uint) ([]models.BookFeedback, error)
        GetFeedbackStats(bookID uint) (map[string]interface{}, error)
        
        // Revision operations
        CreateBookRevision(revision *models.BookRevision) error
        GetBookRevisions(bookID uint) ([]models.BookRevision, error)
        GetBookRevisionByID(id uint) (*models.BookRevision, error)
        
        // Publishing operations
        GetScheduledBooks() ([]models.Book, error)
        
        // Reading progress operations
        GetReadingProgress(userID, bookID uint) (*models.ReadingProgress, error)
        CreateReadingProgress(progress *models.ReadingProgress) error
        UpdateReadingProgress(progress *models.ReadingProgress) error
        
        // Advanced book operations
        GetChapterWithSections(chapterID uint) (*models.BookChapter, error)
        SearchBooks(query string, tags []string, limit int) ([]models.Book, error)
        GetRecommendations(userID uint, limit int) ([]models.Book, error)
}

// BookRepositoryImpl is the implementation of BookRepository
type BookRepositoryImpl struct {
        db *gorm.DB
}

// NewBookRepository creates a new BookRepository
func NewBookRepository(db *gorm.DB) BookRepository {
        return &BookRepositoryImpl{
                db: db,
        }
}

// CreateBook creates a new book record
func (r *BookRepositoryImpl) CreateBook(book *models.Book) error {
        book.CreatedAt = time.Now()
        book.UpdatedAt = time.Now()
        return r.db.Create(book).Error
}

// GetBookByID retrieves a book by its ID
func (r *BookRepositoryImpl) GetBookByID(id uint) (*models.Book, error) {
        var book models.Book
        err := r.db.Where("id = ?", id).First(&book).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("book not found")
                }
                return nil, err
        }
        return &book, nil
}

// GetBookWithChapters retrieves a book with its chapters
func (r *BookRepositoryImpl) GetBookWithChapters(id uint) (*models.Book, error) {
        // Get the book
        book, err := r.GetBookByID(id)
        if err != nil {
                return nil, err
        }
        
        // Get chapters for the book - always include unpublished chapters
        chapters, err := r.GetChaptersByBookID(id, true) // Always include unpublished chapters
        if err != nil {
                return nil, fmt.Errorf("error fetching chapters: %w", err)
        }
        
        // Assign chapters to the book
        book.Chapters = chapters
        
        return book, nil
}

// GetAllBooks retrieves all books
func (r *BookRepositoryImpl) GetAllBooks(includeUnpublished bool) ([]models.Book, error) {
        var books []models.Book
        query := r.db
        if !includeUnpublished {
                query = query.Where("published = ?", true)
        }
        err := query.Find(&books).Error
        return books, err
}

// GetAllBooksWithChapters retrieves all books with their chapters
func (r *BookRepositoryImpl) GetAllBooksWithChapters(includeUnpublished bool) ([]models.Book, error) {
        // Get all books first
        books, err := r.GetAllBooks(includeUnpublished)
        if err != nil {
                return nil, err
        }
        
        // For each book, get its chapters
        for i := range books {
                chapters, err := r.GetChaptersByBookID(books[i].ID, true) // Always include unpublished chapters
                if err != nil {
                        return nil, fmt.Errorf("error fetching chapters for book ID %d: %w", books[i].ID, err)
                }
                books[i].Chapters = chapters
        }
        
        return books, nil
}

// UpdateBook updates a book record
func (r *BookRepositoryImpl) UpdateBook(book *models.Book) error {
        book.UpdatedAt = time.Now()
        return r.db.Save(book).Error
}

// DeleteBook deletes a book by ID
func (r *BookRepositoryImpl) DeleteBook(id uint) error {
        return r.db.Delete(&models.Book{}, id).Error
}

// CreateChapter creates a new chapter record
func (r *BookRepositoryImpl) CreateChapter(chapter *models.BookChapter) error {
        chapter.CreatedAt = time.Now()
        chapter.UpdatedAt = time.Now()
        return r.db.Create(chapter).Error
}

// GetChapterByID retrieves a chapter by its ID
func (r *BookRepositoryImpl) GetChapterByID(id uint) (*models.BookChapter, error) {
        var chapter models.BookChapter
        err := r.db.Where("id = ?", id).First(&chapter).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("chapter not found")
                }
                return nil, err
        }
        return &chapter, nil
}

// GetChaptersByBookID retrieves all chapters for a book
func (r *BookRepositoryImpl) GetChaptersByBookID(bookID uint, includeUnpublished bool) ([]models.BookChapter, error) {
        var chapters []models.BookChapter
        query := r.db.Where("book_id = ?", bookID)
        if !includeUnpublished {
                query = query.Where("published = ?", true)
        }
        err := query.Order("number").Find(&chapters).Error
        return chapters, err
}

// UpdateChapter updates a chapter record
func (r *BookRepositoryImpl) UpdateChapter(chapter *models.BookChapter) error {
        chapter.UpdatedAt = time.Now()
        return r.db.Save(chapter).Error
}

// DeleteChapter deletes a chapter by ID
func (r *BookRepositoryImpl) DeleteChapter(id uint) error {
        return r.db.Delete(&models.BookChapter{}, id).Error
}

// CreateSection creates a new section record
func (r *BookRepositoryImpl) CreateSection(section *models.BookSection) error {
        section.CreatedAt = time.Now()
        section.UpdatedAt = time.Now()
        return r.db.Create(section).Error
}

// GetSectionByID retrieves a section by its ID
func (r *BookRepositoryImpl) GetSectionByID(id uint) (*models.BookSection, error) {
        var section models.BookSection
        err := r.db.Where("id = ?", id).First(&section).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("section not found")
                }
                return nil, err
        }
        return &section, nil
}

// GetSectionsByChapterID retrieves all sections for a chapter
func (r *BookRepositoryImpl) GetSectionsByChapterID(chapterID uint, includeUnpublished bool) ([]models.BookSection, error) {
        var sections []models.BookSection
        query := r.db.Where("chapter_id = ?", chapterID)
        if !includeUnpublished {
                query = query.Where("published = ?", true)
        }
        err := query.Order("number").Find(&sections).Error
        return sections, err
}

// UpdateSection updates a section record
func (r *BookRepositoryImpl) UpdateSection(section *models.BookSection) error {
        section.UpdatedAt = time.Now()
        return r.db.Save(section).Error
}

// DeleteSection deletes a section by ID
func (r *BookRepositoryImpl) DeleteSection(id uint) error {
        return r.db.Delete(&models.BookSection{}, id).Error
}

// CreateSubsection creates a new subsection record
func (r *BookRepositoryImpl) CreateSubsection(subsection *models.BookSubsection) error {
        subsection.CreatedAt = time.Now()
        subsection.UpdatedAt = time.Now()
        return r.db.Create(subsection).Error
}

// GetSubsectionByID retrieves a subsection by its ID
func (r *BookRepositoryImpl) GetSubsectionByID(id uint) (*models.BookSubsection, error) {
        var subsection models.BookSubsection
        err := r.db.Where("id = ?", id).First(&subsection).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("subsection not found")
                }
                return nil, err
        }
        return &subsection, nil
}

// GetSubsectionsBySectionID retrieves all subsections for a section
func (r *BookRepositoryImpl) GetSubsectionsBySectionID(sectionID uint, includeUnpublished bool) ([]models.BookSubsection, error) {
        var subsections []models.BookSubsection
        query := r.db.Where("section_id = ?", sectionID)
        if !includeUnpublished {
                query = query.Where("published = ?", true)
        }
        err := query.Order("number").Find(&subsections).Error
        return subsections, err
}

// UpdateSubsection updates a subsection record
func (r *BookRepositoryImpl) UpdateSubsection(subsection *models.BookSubsection) error {
        subsection.UpdatedAt = time.Now()
        return r.db.Save(subsection).Error
}

// DeleteSubsection deletes a subsection by ID
func (r *BookRepositoryImpl) DeleteSubsection(id uint) error {
        return r.db.Delete(&models.BookSubsection{}, id).Error
}

// CreateFrontMatter creates a new front matter record
func (r *BookRepositoryImpl) CreateFrontMatter(frontMatter *models.BookFrontMatter) error {
        frontMatter.CreatedAt = time.Now()
        frontMatter.UpdatedAt = time.Now()
        return r.db.Create(frontMatter).Error
}

// GetFrontMatterByID retrieves a front matter by its ID
func (r *BookRepositoryImpl) GetFrontMatterByID(id uint) (*models.BookFrontMatter, error) {
        var frontMatter models.BookFrontMatter
        err := r.db.Where("id = ?", id).First(&frontMatter).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("front matter not found")
                }
                return nil, err
        }
        return &frontMatter, nil
}

// GetFrontMatterByBookID retrieves all front matter for a book
func (r *BookRepositoryImpl) GetFrontMatterByBookID(bookID uint, includeUnpublished bool) ([]models.BookFrontMatter, error) {
        var frontMatters []models.BookFrontMatter
        query := r.db.Where("book_id = ?", bookID)
        // FrontMatter doesn't have a published field, so we include all front matter
        err := query.Order("id").Find(&frontMatters).Error
        return frontMatters, err
}

// UpdateFrontMatter updates a front matter record
func (r *BookRepositoryImpl) UpdateFrontMatter(frontMatter *models.BookFrontMatter) error {
        frontMatter.UpdatedAt = time.Now()
        return r.db.Save(frontMatter).Error
}

// DeleteFrontMatter deletes a front matter by ID
func (r *BookRepositoryImpl) DeleteFrontMatter(id uint) error {
        return r.db.Delete(&models.BookFrontMatter{}, id).Error
}

// GetBackMatterByBookID retrieves the back matter for a book
func (r *BookRepositoryImpl) GetBackMatterByBookID(bookID uint, includeUnpublished bool) (*models.BookBackMatter, error) {
        var backMatter models.BookBackMatter
        query := r.db.Where("book_id = ?", bookID)
        // BackMatter doesn't have a published field, so we include all back matter
        err := query.First(&backMatter).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("back matter not found")
                }
                return nil, err
        }

        // For Book 3, process special epilogue and appendix content
        if bookID == 3 && (backMatter.EpilogueJSON != "" || backMatter.AppendixJSON != "") {
                // Process epilogue content from JSON
                if backMatter.EpilogueJSON != "" {
                        var epilogueItems []models.EpilogueItem
                        if err := r.processJSONIntoStruct(backMatter.EpilogueJSON, &epilogueItems); err == nil {
                                backMatter.Epilogue = epilogueItems
                        } else {
                                fmt.Printf("Error parsing epilogue JSON: %v\n", err)
                        }
                }

                // Process appendix content from JSON
                if backMatter.AppendixJSON != "" {
                        var appendixItems []models.AppendixItem
                        if err := r.processJSONIntoStruct(backMatter.AppendixJSON, &appendixItems); err == nil {
                                backMatter.Appendix = appendixItems
                        } else {
                                fmt.Printf("Error parsing appendix JSON: %v\n", err)
                        }
                }
        }

        return &backMatter, nil
}

// UpdateBackMatter updates a back matter record
func (r *BookRepositoryImpl) UpdateBackMatter(backMatter *models.BookBackMatter) error {
        backMatter.UpdatedAt = time.Now()
        
        // Handle epilogue and appendix items for Book 3
        if backMatter.BookID == 3 {
                // Convert epilogue items to JSON if they exist
                if len(backMatter.Epilogue) > 0 {
                        epilogueJSON, err := json.Marshal(backMatter.Epilogue)
                        if err != nil {
                                return fmt.Errorf("error marshaling epilogue items: %w", err)
                        }
                        backMatter.EpilogueJSON = string(epilogueJSON)
                }
                
                // Convert appendix items to JSON if they exist
                if len(backMatter.Appendix) > 0 {
                        appendixJSON, err := json.Marshal(backMatter.Appendix)
                        if err != nil {
                                return fmt.Errorf("error marshaling appendix items: %w", err)
                        }
                        backMatter.AppendixJSON = string(appendixJSON)
                }
        }
        
        return r.db.Save(backMatter).Error
}

// DeleteBackMatter deletes a back matter by ID
func (r *BookRepositoryImpl) DeleteBackMatter(id uint) error {
        return r.db.Delete(&models.BookBackMatter{}, id).Error
}

// CreateForumTopic creates a new forum topic record
func (r *BookRepositoryImpl) CreateForumTopic(topic *models.ForumTopic) error {
        topic.CreatedAt = time.Now()
        topic.UpdatedAt = time.Now()
        return r.db.Create(topic).Error
}

// CreateInteractiveElement creates a new interactive element
func (r *BookRepositoryImpl) CreateInteractiveElement(element *models.InteractiveElement) error {
        return r.db.Create(element).Error
}

// GetInteractiveElementByID retrieves an interactive element by its ID
func (r *BookRepositoryImpl) GetInteractiveElementByID(id uint) (*models.InteractiveElement, error) {
        var element models.InteractiveElement
        err := r.db.Where("id = ?", id).First(&element).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("interactive element not found")
                }
                return nil, err
        }
        return &element, nil
}

// GetInteractiveElementsBySectionID retrieves all interactive elements for a section
func (r *BookRepositoryImpl) GetInteractiveElementsBySectionID(sectionID uint) ([]models.InteractiveElement, error) {
        var elements []models.InteractiveElement
        err := r.db.Where("section_id = ?", sectionID).Order("position").Find(&elements).Error
        return elements, err
}

// UpdateInteractiveElement updates an interactive element
func (r *BookRepositoryImpl) UpdateInteractiveElement(element *models.InteractiveElement) error {
        return r.db.Save(element).Error
}

// DeleteInteractiveElement deletes an interactive element by ID
func (r *BookRepositoryImpl) DeleteInteractiveElement(id uint) error {
        return r.db.Delete(&models.InteractiveElement{}, id).Error
}

// GetForumTopicByID retrieves a forum topic by its ID
func (r *BookRepositoryImpl) GetForumTopicByID(id uint) (*models.ForumTopic, error) {
        var topic models.ForumTopic
        err := r.db.Where("id = ?", id).First(&topic).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("forum topic not found")
                }
                return nil, err
        }
        return &topic, nil
}

// GetForumTopicsBySectionID retrieves all forum topics for a section
func (r *BookRepositoryImpl) GetForumTopicsBySectionID(sectionID uint) ([]models.ForumTopic, error) {
        var topics []models.ForumTopic
        err := r.db.Where("section_id = ?", sectionID).Order("order_index").Find(&topics).Error
        return topics, err
}

// UpdateForumTopic updates a forum topic record
func (r *BookRepositoryImpl) UpdateForumTopic(topic *models.ForumTopic) error {
        topic.UpdatedAt = time.Now()
        return r.db.Save(topic).Error
}

// DeleteForumTopic deletes a forum topic by ID
func (r *BookRepositoryImpl) DeleteForumTopic(id uint) error {
        return r.db.Delete(&models.ForumTopic{}, id).Error
}

// CreateActionStep creates a new action step record
func (r *BookRepositoryImpl) CreateActionStep(step *models.ActionStep) error {
        step.CreatedAt = time.Now()
        step.UpdatedAt = time.Now()
        return r.db.Create(step).Error
}

// GetActionStepByID retrieves an action step by its ID
func (r *BookRepositoryImpl) GetActionStepByID(id uint) (*models.ActionStep, error) {
        var step models.ActionStep
        err := r.db.Where("id = ?", id).First(&step).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("action step not found")
                }
                return nil, err
        }
        return &step, nil
}

// GetActionStepsBySectionID retrieves all action steps for a section
func (r *BookRepositoryImpl) GetActionStepsBySectionID(sectionID uint) ([]models.ActionStep, error) {
        var steps []models.ActionStep
        err := r.db.Where("section_id = ?", sectionID).Order("order_index").Find(&steps).Error
        return steps, err
}

// UpdateActionStep updates an action step record
func (r *BookRepositoryImpl) UpdateActionStep(step *models.ActionStep) error {
        step.UpdatedAt = time.Now()
        return r.db.Save(step).Error
}

// DeleteActionStep deletes an action step by ID
func (r *BookRepositoryImpl) DeleteActionStep(id uint) error {
        return r.db.Delete(&models.ActionStep{}, id).Error
}

// CreateOrUpdateBookProgress creates or updates a book progress record
func (r *BookRepositoryImpl) CreateOrUpdateBookProgress(progress *models.BookProgress) error {
        var existingProgress models.BookProgress
        err := r.db.Where("user_id = ? AND book_id = ? AND section_id = ?", 
                progress.UserID, progress.BookID, progress.SectionID).First(&existingProgress).Error
        
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        // Create new progress
                        progress.CreatedAt = time.Now()
                        progress.UpdatedAt = time.Now()
                        return r.db.Create(progress).Error
                }
                return err
        }
        
        // Update existing progress
        existingProgress.IsRead = progress.IsRead
        existingProgress.Progress = progress.Progress
        existingProgress.UpdatedAt = time.Now()
        return r.db.Save(&existingProgress).Error
}

// GetBookProgressByUserAndBook retrieves book progress records for a user and book
func (r *BookRepositoryImpl) GetBookProgressByUserAndBook(userID, bookID uint) ([]models.BookProgress, error) {
        var progress []models.BookProgress
        err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Find(&progress).Error
        return progress, err
}

// GetBookProgressPercentage calculates the total progress percentage for a book
func (r *BookRepositoryImpl) GetBookProgressPercentage(userID, bookID uint) (float64, error) {
        // Count total sections
        var totalSections int64
        err := r.db.Model(&models.BookSection{}).
                Joins("JOIN book_chapters ON book_sections.chapter_id = book_chapters.id").
                Where("book_chapters.book_id = ? AND book_sections.published = ?", bookID, true).
                Count(&totalSections).Error
        if err != nil {
                return 0, err
        }
        
        if totalSections == 0 {
                return 0, nil
        }
        
        // Count read sections
        var readSections int64
        err = r.db.Model(&models.BookProgress{}).
                Where("user_id = ? AND book_id = ? AND is_read = ?", userID, bookID, true).
                Count(&readSections).Error
        if err != nil {
                return 0, err
        }
        
        return float64(readSections) / float64(totalSections) * 100, nil
}

// CreateBookmark creates a new bookmark record
func (r *BookRepositoryImpl) CreateBookmark(bookmark *models.BookBookmark) error {
        bookmark.CreatedAt = time.Now()
        bookmark.UpdatedAt = time.Now()
        return r.db.Create(bookmark).Error
}

// GetBookmarkByID retrieves a bookmark by its ID
func (r *BookRepositoryImpl) GetBookmarkByID(id uint) (*models.BookBookmark, error) {
        var bookmark models.BookBookmark
        err := r.db.Where("id = ?", id).First(&bookmark).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("bookmark not found")
                }
                return nil, err
        }
        return &bookmark, nil
}

// GetBookmarksByUserAndBook retrieves bookmarks for a user and book
func (r *BookRepositoryImpl) GetBookmarksByUserAndBook(userID, bookID uint) ([]models.BookBookmark, error) {
        var bookmarks []models.BookBookmark
        err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Find(&bookmarks).Error
        return bookmarks, err
}

// DeleteBookmark deletes a bookmark by ID
func (r *BookRepositoryImpl) DeleteBookmark(id uint) error {
        return r.db.Delete(&models.BookBookmark{}, id).Error
}

// CreateBookRevision creates a new book revision record
func (r *BookRepositoryImpl) CreateBookRevision(revision *models.BookRevision) error {
        revision.CreatedAt = time.Now()
        return r.db.Create(revision).Error
}

// GetBookRevisions retrieves all revisions for a book
func (r *BookRepositoryImpl) GetBookRevisions(bookID uint) ([]models.BookRevision, error) {
        var revisions []models.BookRevision
        err := r.db.Where("book_id = ?", bookID).Order("created_at DESC").Find(&revisions).Error
        return revisions, err
}

// GetBookRevisionByID retrieves a book revision by its ID
func (r *BookRepositoryImpl) GetBookRevisionByID(id uint) (*models.BookRevision, error) {
        var revision models.BookRevision
        err := r.db.Where("id = ?", id).First(&revision).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("book revision not found")
                }
                return nil, err
        }
        return &revision, nil
}

// GetScheduledBooks retrieves all books scheduled to be published
func (r *BookRepositoryImpl) GetScheduledBooks() ([]models.Book, error) {
        var books []models.Book
        err := r.db.Where("scheduled_publish_at IS NOT NULL AND scheduled_publish_at > ?", time.Now()).Find(&books).Error
        return books, err
}

// CreateNote creates a new note record
func (r *BookRepositoryImpl) CreateNote(note *models.BookNote) error {
        note.CreatedAt = time.Now()
        note.UpdatedAt = time.Now()
        return r.db.Create(note).Error
}

// GetNoteByID retrieves a note by its ID
func (r *BookRepositoryImpl) GetNoteByID(id uint) (*models.BookNote, error) {
        var note models.BookNote
        err := r.db.Where("id = ?", id).First(&note).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("note not found")
                }
                return nil, err
        }
        return &note, nil
}

// GetNotesByUserAndBook retrieves notes for a user and book
func (r *BookRepositoryImpl) GetNotesByUserAndBook(userID, bookID uint) ([]models.BookNote, error) {
        var notes []models.BookNote
        err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Find(&notes).Error
        return notes, err
}

// UpdateNote updates a note record
func (r *BookRepositoryImpl) UpdateNote(note *models.BookNote) error {
        note.UpdatedAt = time.Now()
        return r.db.Save(note).Error
}

// DeleteNote deletes a note by ID
func (r *BookRepositoryImpl) DeleteNote(id uint) error {
        return r.db.Delete(&models.BookNote{}, id).Error
}

// CreateFeedback creates a new feedback record
func (r *BookRepositoryImpl) CreateFeedback(feedback *models.BookFeedback) error {
        feedback.CreatedAt = time.Now()
        feedback.UpdatedAt = time.Now()
        return r.db.Create(feedback).Error
}

// GetFeedbackByID retrieves a feedback record by its ID
func (r *BookRepositoryImpl) GetFeedbackByID(id uint) (*models.BookFeedback, error) {
        var feedback models.BookFeedback
        err := r.db.Where("id = ?", id).First(&feedback).Error
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return nil, errors.New("feedback not found")
                }
                return nil, err
        }
        return &feedback, nil
}

// GetFeedbackByUserAndBook retrieves feedback for a user and book
func (r *BookRepositoryImpl) GetFeedbackByUserAndBook(userID, bookID uint) ([]models.BookFeedback, error) {
        var feedback []models.BookFeedback
        err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Find(&feedback).Error
        return feedback, err
}

// GetFeedbackStats retrieves aggregated feedback statistics for a book
func (r *BookRepositoryImpl) GetFeedbackStats(bookID uint) (map[string]interface{}, error) {
        stats := make(map[string]interface{})
        
        // Get mood ratings average
        var moodAvg float64
        err := r.db.Model(&models.BookFeedback{}).
                Where("book_id = ? AND feedback_type = ?", bookID, "mood").
                Select("AVG(rating)").
                Row().Scan(&moodAvg)
        if err != nil {
                return nil, err
        }
        stats["mood_average"] = moodAvg
        
        // Get difficulty ratings average
        var difficultyAvg float64
        err = r.db.Model(&models.BookFeedback{}).
                Where("book_id = ? AND feedback_type = ?", bookID, "difficulty").
                Select("AVG(rating)").
                Row().Scan(&difficultyAvg)
        if err != nil {
                return nil, err
        }
        stats["difficulty_average"] = difficultyAvg
        
        // Get total feedback count
        var totalCount int64
        err = r.db.Model(&models.BookFeedback{}).
                Where("book_id = ?", bookID).
                Count(&totalCount).Error
        if err != nil {
                return nil, err
        }
        stats["total_feedback_count"] = totalCount
        
        return stats, nil
}

// GetReadingProgress retrieves the reading progress for a user and book
func (r *BookRepositoryImpl) GetReadingProgress(userID, bookID uint) (*models.ReadingProgress, error) {
        var progress models.ReadingProgress
        err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&progress).Error
        
        if err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        // No reading progress found, return nil without error
                        return nil, nil
                }
                return nil, err
        }
        
        return &progress, nil
}

// CreateReadingProgress creates a new reading progress record
func (r *BookRepositoryImpl) CreateReadingProgress(progress *models.ReadingProgress) error {
        progress.CreatedAt = time.Now()
        progress.UpdatedAt = time.Now()
        return r.db.Create(progress).Error
}

// UpdateReadingProgress updates an existing reading progress record
func (r *BookRepositoryImpl) UpdateReadingProgress(progress *models.ReadingProgress) error {
        progress.UpdatedAt = time.Now()
        return r.db.Save(progress).Error
}

// GetChapterWithSections gets a chapter with all its sections
func (r *BookRepositoryImpl) GetChapterWithSections(chapterID uint) (*models.BookChapter, error) {
        chapter, err := r.GetChapterByID(chapterID)
        if err != nil {
                return nil, err
        }
        
        sections, err := r.GetSectionsByChapterID(chapterID, true) // include unpublished sections
        if err != nil {
                return nil, err
        }
        
        chapter.Sections = sections
        return chapter, nil
}

// SearchBooks searches for books based on query text and tags
func (r *BookRepositoryImpl) SearchBooks(query string, tags []string, limit int) ([]models.Book, error) {
        db := r.db
        
        // Start with a base query
        baseQuery := db.Model(&models.Book{}).Where("published = ?", true)
        
        // Apply text search if query is provided
        if query != "" {
                // Case-insensitive search on title and description
                textSearch := "%" + query + "%"
                baseQuery = baseQuery.Where("LOWER(title) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", 
                        textSearch, textSearch)
        }
        
        // Apply tag filtering if tags are provided
        if len(tags) > 0 {
                baseQuery = baseQuery.Joins("JOIN book_tags ON books.id = book_tags.book_id").
                        Where("book_tags.tag IN ?", tags).
                        Group("books.id")
        }
        
        // Apply limit if specified
        if limit > 0 {
                baseQuery = baseQuery.Limit(limit)
        }
        
        // Execute the query
        var books []models.Book
        err := baseQuery.Order("created_at DESC").Find(&books).Error
        return books, err
}

// GetRecommendations gets book recommendations for a user
func (r *BookRepositoryImpl) GetRecommendations(userID uint, limit int) ([]models.Book, error) {
        // This would normally use a more sophisticated recommendation algorithm
        // For now, we'll implement a basic version that returns:
        // 1. Books in categories the user has read before
        // 2. Popular books the user hasn't read
        // 3. Recent books
        
        var readBookIDs []uint
        
        // Get IDs of books the user has already read (using reading progress)
        err := r.db.Model(&models.ReadingProgress{}).
                Where("user_id = ?", userID).
                Pluck("book_id", &readBookIDs).Error
        
        if err != nil {
                return nil, err
        }
        
        // If user hasn't read any books, just return popular ones
        if len(readBookIDs) == 0 {
                var popularBooks []models.Book
                err := r.db.Where("published = ?", true).
                        Order("view_count DESC").
                        Limit(limit).
                        Find(&popularBooks).Error
                return popularBooks, err
        }
        
        // Get categories from books the user has read
        var categories []string
        err = r.db.Model(&models.Book{}).
                Joins("JOIN book_categories ON books.id = book_categories.book_id").
                Where("books.id IN ?", readBookIDs).
                Pluck("DISTINCT book_categories.category", &categories).Error
        
        if err != nil {
                return nil, err
        }
        
        // Get recommended books based on those categories, excluding already read books
        var recommendedBooks []models.Book
        query := r.db.Model(&models.Book{}).
                Joins("JOIN book_categories ON books.id = book_categories.book_id").
                Where("book_categories.category IN ? AND books.published = ? AND books.id NOT IN ?", 
                        categories, true, readBookIDs).
                Group("books.id").
                Order("books.view_count DESC").
                Limit(limit)
        
        err = query.Find(&recommendedBooks).Error
        
        // If we couldn't find enough recommendations, add some recent books
        if len(recommendedBooks) < limit {
                remainingLimit := limit - len(recommendedBooks)
                var recentBooks []models.Book
                
                err = r.db.Where("published = ? AND id NOT IN ?", true, readBookIDs).
                        Order("created_at DESC").
                        Limit(remainingLimit).
                        Find(&recentBooks).Error
                
                if err != nil {
                        return nil, err
                }
                
                recommendedBooks = append(recommendedBooks, recentBooks...)
        }
        
        return recommendedBooks, nil
}

// processJSONIntoStruct converts a JSON string into a struct
func (r *BookRepositoryImpl) processJSONIntoStruct(jsonStr string, target interface{}) error {
        if jsonStr == "" {
                return fmt.Errorf("empty JSON string")
        }
        return json.Unmarshal([]byte(jsonStr), target)
}