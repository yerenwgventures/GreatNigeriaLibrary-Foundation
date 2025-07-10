package service

import (
        "encoding/csv"
        "encoding/json"
        "fmt"
        "io"
        "strings"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/repository"
)

// ContentAdminService defines the interface for content administration operations
type ContentAdminService interface {
        // Import functions
        ImportBooksFromJSON(reader io.Reader) ([]models.Book, error)
        ImportChaptersFromJSON(reader io.Reader, bookID uint) ([]models.Chapter, error)
        ImportSectionsFromJSON(reader io.Reader, chapterID uint) ([]models.Section, error)
        ImportBooksFromCSV(reader io.Reader) ([]models.Book, error)
        ImportChaptersFromCSV(reader io.Reader, bookID uint) ([]models.Chapter, error)
        ImportSectionsFromCSV(reader io.Reader, chapterID uint) ([]models.Section, error)
        
        // Export functions
        ExportBooksToJSON(writer io.Writer, bookIDs []uint) error
        ExportChaptersToJSON(writer io.Writer, bookID uint) error
        ExportSectionsToJSON(writer io.Writer, chapterID uint) error
        ExportBooksToCSV(writer io.Writer, bookIDs []uint) error
        ExportChaptersToCSV(writer io.Writer, bookID uint) error
        ExportSectionsToCSV(writer io.Writer, chapterID uint) error
        
        // Content revision functions
        CreateBookRevision(bookID uint, changes map[string]interface{}, notes string) (*models.BookRevision, error)
        CreateChapterRevision(chapterID uint, changes map[string]interface{}, notes string) (*models.ChapterRevision, error)
        CreateSectionRevision(sectionID uint, changes map[string]interface{}, notes string) (*models.SectionRevision, error)
        GetBookRevisions(bookID uint) ([]models.BookRevision, error)
        GetChapterRevisions(chapterID uint) ([]models.ChapterRevision, error)
        GetSectionRevisions(sectionID uint) ([]models.SectionRevision, error)
        RestoreBookRevision(revisionID uint) error
        RestoreChapterRevision(revisionID uint) error
        RestoreSectionRevision(revisionID uint) error
        
        // Content scheduling and publishing
        ScheduleBookPublishing(bookID uint, publishDate time.Time) error
        ScheduleChapterPublishing(chapterID uint, publishDate time.Time) error
        ScheduleSectionPublishing(sectionID uint, publishDate time.Time) error
        GetScheduledContent() ([]interface{}, error)
        PublishContent(contentType string, contentID uint) error
        UnpublishContent(contentType string, contentID uint) error
}

// ContentAdminServiceImpl implements the ContentAdminService interface
type ContentAdminServiceImpl struct {
        bookRepo    repository.BookRepository
        chapterRepo repository.ChapterRepository
        sectionRepo repository.SectionRepository
}

// NewContentAdminService creates a new content admin service
func NewContentAdminService(
        bookRepo repository.BookRepository,
        chapterRepo repository.ChapterRepository,
        sectionRepo repository.SectionRepository,
) ContentAdminService {
        return &ContentAdminServiceImpl{
                bookRepo:    bookRepo,
                chapterRepo: chapterRepo,
                sectionRepo: sectionRepo,
        }
}

// ImportBooksFromJSON imports books from a JSON file
func (s *ContentAdminServiceImpl) ImportBooksFromJSON(reader io.Reader) ([]models.Book, error) {
        var books []models.Book
        
        // Decode JSON from reader
        decoder := json.NewDecoder(reader)
        if err := decoder.Decode(&books); err != nil {
                return nil, fmt.Errorf("error decoding JSON: %w", err)
        }
        
        // Validate and save each book
        importedBooks := make([]models.Book, 0, len(books))
        for _, book := range books {
                // Clear any ID field to avoid overwriting existing books
                book.ID = 0
                
                // Set default values if needed
                if book.CreatedAt.IsZero() {
                        book.CreatedAt = time.Now()
                }
                if book.UpdatedAt.IsZero() {
                        book.UpdatedAt = time.Now()
                }
                
                // Save the book
                var err error
                err = s.bookRepo.CreateBook(&book)
                if err != nil {
                        return importedBooks, fmt.Errorf("error saving book '%s': %w", book.Title, err)
                }
                
                importedBooks = append(importedBooks, book)
        }
        
        return importedBooks, nil
}

// ImportChaptersFromJSON imports chapters from a JSON file
func (s *ContentAdminServiceImpl) ImportChaptersFromJSON(reader io.Reader, bookID uint) ([]models.Chapter, error) {
        var chapters []models.Chapter
        
        // Decode JSON from reader
        decoder := json.NewDecoder(reader)
        if err := decoder.Decode(&chapters); err != nil {
                return nil, fmt.Errorf("error decoding JSON: %w", err)
        }
        
        // Check if book exists
        _, err := s.bookRepo.GetBookByID(bookID)
        if err != nil {
                return nil, fmt.Errorf("book with ID %d not found: %w", bookID, err)
        }
        
        // Validate and save each chapter
        importedChapters := make([]models.Chapter, 0, len(chapters))
        for _, chapter := range chapters {
                // Clear any ID field to avoid overwriting existing chapters
                chapter.ID = 0
                
                // Set the book ID
                chapter.BookID = bookID
                
                // Set default values if needed
                if chapter.CreatedAt.IsZero() {
                        chapter.CreatedAt = time.Now()
                }
                if chapter.UpdatedAt.IsZero() {
                        chapter.UpdatedAt = time.Now()
                }
                
                // Convert to BookChapter
                bookChapter := models.BookChapter{
                        BookID:      chapter.BookID,
                        Title:       chapter.Title,
                        Number:      chapter.Number,
                        Description: chapter.Description,
                        Published:   false,  // Default to not published
                        CreatedAt:   chapter.CreatedAt,
                        UpdatedAt:   chapter.UpdatedAt,
                }
                
                // Save the chapter
                err = s.chapterRepo.CreateChapter(&bookChapter)
                if err != nil {
                        return importedChapters, fmt.Errorf("error saving chapter '%s': %w", chapter.Title, err)
                }
                
                importedChapters = append(importedChapters, chapter)
        }
        
        return importedChapters, nil
}

// ImportSectionsFromJSON imports sections from a JSON file
func (s *ContentAdminServiceImpl) ImportSectionsFromJSON(reader io.Reader, chapterID uint) ([]models.Section, error) {
        var sections []models.Section
        
        // Decode JSON from reader
        decoder := json.NewDecoder(reader)
        if err := decoder.Decode(&sections); err != nil {
                return nil, fmt.Errorf("error decoding JSON: %w", err)
        }
        
        // Check if chapter exists
        chapter, err := s.chapterRepo.GetChapterByID(chapterID)
        if err != nil {
                return nil, fmt.Errorf("chapter with ID %d not found: %w", chapterID, err)
        }
        
        // Validate and save each section
        importedSections := make([]models.Section, 0, len(sections))
        for _, section := range sections {
                // Clear any ID field to avoid overwriting existing sections
                section.ID = 0
                
                // Set the chapter ID and book ID
                section.ChapterID = chapterID
                section.BookID = chapter.BookID
                
                // Set default values if needed
                if section.CreatedAt.IsZero() {
                        section.CreatedAt = time.Now()
                }
                if section.UpdatedAt.IsZero() {
                        section.UpdatedAt = time.Now()
                }
                
                // Convert to BookSection
                bookSection := models.BookSection{
                        BookID:      section.BookID,
                        ChapterID:   section.ChapterID,
                        Title:       section.Title,
                        Number:      section.Number,
                        Content:     section.Content,
                        Published:   false,  // Default to not published
                        CreatedAt:   section.CreatedAt,
                        UpdatedAt:   section.UpdatedAt,
                }
                
                // Save the section
                err = s.sectionRepo.CreateSection(&bookSection)
                if err != nil {
                        return importedSections, fmt.Errorf("error saving section '%s': %w", section.Title, err)
                }
                
                importedSections = append(importedSections, section)
        }
        
        return importedSections, nil
}

// ImportBooksFromCSV imports books from a CSV file
func (s *ContentAdminServiceImpl) ImportBooksFromCSV(reader io.Reader) ([]models.Book, error) {
        // Create CSV reader
        csvReader := csv.NewReader(reader)
        
        // Read header
        header, err := csvReader.Read()
        if err != nil {
                return nil, fmt.Errorf("error reading CSV header: %w", err)
        }
        
        // Map header indices
        headerMap := make(map[string]int)
        for i, h := range header {
                headerMap[strings.ToLower(strings.TrimSpace(h))] = i
        }
        
        // Check required headers
        requiredHeaders := []string{"title", "description"}
        for _, h := range requiredHeaders {
                if _, ok := headerMap[h]; !ok {
                        return nil, fmt.Errorf("required header '%s' not found in CSV", h)
                }
        }
        
        // Read and process rows
        importedBooks := make([]models.Book, 0)
        for {
                row, err := csvReader.Read()
                if err == io.EOF {
                        break
                }
                if err != nil {
                        return importedBooks, fmt.Errorf("error reading CSV row: %w", err)
                }
                
                // Create book from row
                book := models.Book{
                        Title:       row[headerMap["title"]],
                        Description: row[headerMap["description"]],
                        CreatedAt:   time.Now(),
                        UpdatedAt:   time.Now(),
                }
                
                // Optional fields
                if idx, ok := headerMap["author"]; ok && idx < len(row) {
                        book.Author = row[idx]
                }
                if idx, ok := headerMap["cover_image"]; ok && idx < len(row) {
                        book.CoverImage = row[idx]
                }
                if idx, ok := headerMap["published"]; ok && idx < len(row) {
                        book.Published = row[idx] == "true" || row[idx] == "1" || row[idx] == "yes"
                }
                
                // Save the book
                saveErr := s.bookRepo.CreateBook(&book)
                if saveErr != nil {
                        return importedBooks, fmt.Errorf("error saving book '%s': %w", book.Title, saveErr)
                }
                
                importedBooks = append(importedBooks, book)
        }
        
        return importedBooks, nil
}

// ImportChaptersFromCSV imports chapters from a CSV file
func (s *ContentAdminServiceImpl) ImportChaptersFromCSV(reader io.Reader, bookID uint) ([]models.Chapter, error) {
        // Check if book exists
        _, err := s.bookRepo.GetBookByID(bookID)
        if err != nil {
                return nil, fmt.Errorf("book with ID %d not found: %w", bookID, err)
        }
        
        // Create CSV reader
        csvReader := csv.NewReader(reader)
        
        // Read header
        header, err := csvReader.Read()
        if err != nil {
                return nil, fmt.Errorf("error reading CSV header: %w", err)
        }
        
        // Map header indices
        headerMap := make(map[string]int)
        for i, h := range header {
                headerMap[strings.ToLower(strings.TrimSpace(h))] = i
        }
        
        // Check required headers
        requiredHeaders := []string{"title", "number"}
        for _, h := range requiredHeaders {
                if _, ok := headerMap[h]; !ok {
                        return nil, fmt.Errorf("required header '%s' not found in CSV", h)
                }
        }
        
        // Read and process rows
        importedChapters := make([]models.Chapter, 0)
        for {
                row, err := csvReader.Read()
                if err == io.EOF {
                        break
                }
                if err != nil {
                        return importedChapters, fmt.Errorf("error reading CSV row: %w", err)
                }
                
                // Parse chapter number
                var chapterNumber int
                fmt.Sscanf(row[headerMap["number"]], "%d", &chapterNumber)
                
                // Create chapter from row
                chapter := models.Chapter{
                        BookID:      bookID,
                        Title:       row[headerMap["title"]],
                        Number:      chapterNumber,
                        CreatedAt:   time.Now(),
                        UpdatedAt:   time.Now(),
                }
                
                // Optional fields
                if idx, ok := headerMap["description"]; ok && idx < len(row) {
                        chapter.Description = row[idx]
                }
                
                // Create BookChapter from Chapter
                published := false
                if idx, ok := headerMap["published"]; ok && idx < len(row) {
                        published = row[idx] == "true" || row[idx] == "1" || row[idx] == "yes"
                }
                
                bookChapter := models.BookChapter{
                        BookID:      chapter.BookID,
                        Title:       chapter.Title,
                        Number:      chapter.Number,
                        Description: chapter.Description,
                        Published:   published,
                        CreatedAt:   chapter.CreatedAt,
                        UpdatedAt:   chapter.UpdatedAt,
                }
                
                // Save the chapter
                err = s.chapterRepo.CreateChapter(&bookChapter)
                if err != nil {
                        return importedChapters, fmt.Errorf("error saving chapter '%s': %w", chapter.Title, err)
                }
                
                importedChapters = append(importedChapters, chapter)
        }
        
        return importedChapters, nil
}

// ImportSectionsFromCSV imports sections from a CSV file
func (s *ContentAdminServiceImpl) ImportSectionsFromCSV(reader io.Reader, chapterID uint) ([]models.Section, error) {
        // Check if chapter exists
        chapter, err := s.chapterRepo.GetChapterByID(chapterID)
        if err != nil {
                return nil, fmt.Errorf("chapter with ID %d not found: %w", chapterID, err)
        }
        
        // Create CSV reader
        csvReader := csv.NewReader(reader)
        
        // Read header
        header, err := csvReader.Read()
        if err != nil {
                return nil, fmt.Errorf("error reading CSV header: %w", err)
        }
        
        // Map header indices
        headerMap := make(map[string]int)
        for i, h := range header {
                headerMap[strings.ToLower(strings.TrimSpace(h))] = i
        }
        
        // Check required headers
        requiredHeaders := []string{"title", "content", "number"}
        for _, h := range requiredHeaders {
                if _, ok := headerMap[h]; !ok {
                        return nil, fmt.Errorf("required header '%s' not found in CSV", h)
                }
        }
        
        // Read and process rows
        importedSections := make([]models.Section, 0)
        for {
                row, err := csvReader.Read()
                if err == io.EOF {
                        break
                }
                if err != nil {
                        return importedSections, fmt.Errorf("error reading CSV row: %w", err)
                }
                
                // Parse section number
                var sectionNumber int
                fmt.Sscanf(row[headerMap["number"]], "%d", &sectionNumber)
                
                // Create section from row
                section := models.Section{
                        ChapterID:   chapterID,
                        BookID:      chapter.BookID,
                        Title:       row[headerMap["title"]],
                        Content:     row[headerMap["content"]],
                        Number:      sectionNumber,
                        CreatedAt:   time.Now(),
                        UpdatedAt:   time.Now(),
                }
                
                // Create BookSection from Section
                published := false
                if idx, ok := headerMap["published"]; ok && idx < len(row) {
                        published = row[idx] == "true" || row[idx] == "1" || row[idx] == "yes"
                }
                
                bookSection := models.BookSection{
                        BookID:      section.BookID,
                        ChapterID:   section.ChapterID,
                        Title:       section.Title,
                        Number:      section.Number,
                        Content:     section.Content,
                        Published:   published,
                        CreatedAt:   section.CreatedAt,
                        UpdatedAt:   section.UpdatedAt,
                }
                
                // Save the section
                err = s.sectionRepo.CreateSection(&bookSection)
                if err != nil {
                        return importedSections, fmt.Errorf("error saving section '%s': %w", section.Title, err)
                }
                
                importedSections = append(importedSections, section)
        }
        
        return importedSections, nil
}

// ExportBooksToJSON exports books to a JSON file
func (s *ContentAdminServiceImpl) ExportBooksToJSON(writer io.Writer, bookIDs []uint) error {
        var books []models.Book
        var err error
        
        // If bookIDs is empty, get all books
        if len(bookIDs) == 0 {
                books, err = s.bookRepo.GetAllBooks(true) // Include unpublished books
        } else {
                // Get specified books
                books = make([]models.Book, 0, len(bookIDs))
                for _, id := range bookIDs {
                        book, err := s.bookRepo.GetBookByID(id)
                        if err != nil {
                                return fmt.Errorf("error getting book with ID %d: %w", id, err)
                        }
                        books = append(books, *book)
                }
        }
        
        if err != nil {
                return fmt.Errorf("error getting books: %w", err)
        }
        
        // Encode to JSON
        encoder := json.NewEncoder(writer)
        encoder.SetIndent("", "  ")
        return encoder.Encode(books)
}

// ExportChaptersToJSON exports chapters of a book to a JSON file
func (s *ContentAdminServiceImpl) ExportChaptersToJSON(writer io.Writer, bookID uint) error {
        // Get book with chapters
        book, err := s.bookRepo.GetBookWithChapters(bookID)
        if err != nil {
                return fmt.Errorf("error getting book with ID %d: %w", bookID, err)
        }
        
        // Encode to JSON
        encoder := json.NewEncoder(writer)
        encoder.SetIndent("", "  ")
        return encoder.Encode(book.Chapters)
}

// ExportSectionsToJSON exports sections of a chapter to a JSON file
func (s *ContentAdminServiceImpl) ExportSectionsToJSON(writer io.Writer, chapterID uint) error {
        // Get chapter with sections
        chapter, err := s.chapterRepo.GetChapterWithSections(chapterID)
        if err != nil {
                return fmt.Errorf("error getting chapter with ID %d: %w", chapterID, err)
        }
        
        // Encode to JSON
        encoder := json.NewEncoder(writer)
        encoder.SetIndent("", "  ")
        return encoder.Encode(chapter.Sections)
}

// ExportBooksToCSV exports books to a CSV file
func (s *ContentAdminServiceImpl) ExportBooksToCSV(writer io.Writer, bookIDs []uint) error {
        var books []models.Book
        var err error
        
        // If bookIDs is empty, get all books
        if len(bookIDs) == 0 {
                books, err = s.bookRepo.GetAllBooks(true) // Include unpublished books
        } else {
                // Get specified books
                books = make([]models.Book, 0, len(bookIDs))
                for _, id := range bookIDs {
                        book, err := s.bookRepo.GetBookByID(id)
                        if err != nil {
                                return fmt.Errorf("error getting book with ID %d: %w", id, err)
                        }
                        books = append(books, *book)
                }
        }
        
        if err != nil {
                return fmt.Errorf("error getting books: %w", err)
        }
        
        // Create CSV writer
        csvWriter := csv.NewWriter(writer)
        defer csvWriter.Flush()
        
        // Write header
        header := []string{"ID", "Title", "Description", "Author", "CoverImage", "Published", "CreatedAt", "UpdatedAt"}
        if err := csvWriter.Write(header); err != nil {
                return fmt.Errorf("error writing CSV header: %w", err)
        }
        
        // Write rows
        for _, book := range books {
                row := []string{
                        fmt.Sprintf("%d", book.ID),
                        book.Title,
                        book.Description,
                        book.Author,
                        book.CoverImage,
                        fmt.Sprintf("%t", book.Published),
                        book.CreatedAt.Format(time.RFC3339),
                        book.UpdatedAt.Format(time.RFC3339),
                }
                if err := csvWriter.Write(row); err != nil {
                        return fmt.Errorf("error writing CSV row: %w", err)
                }
        }
        
        return nil
}

// ExportChaptersToCSV exports chapters of a book to a CSV file
func (s *ContentAdminServiceImpl) ExportChaptersToCSV(writer io.Writer, bookID uint) error {
        // Get book with chapters
        book, err := s.bookRepo.GetBookWithChapters(bookID)
        if err != nil {
                return fmt.Errorf("error getting book with ID %d: %w", bookID, err)
        }
        
        // Create CSV writer
        csvWriter := csv.NewWriter(writer)
        defer csvWriter.Flush()
        
        // Write header
        header := []string{"ID", "BookID", "Title", "Description", "Number", "Published", "CreatedAt", "UpdatedAt"}
        if err := csvWriter.Write(header); err != nil {
                return fmt.Errorf("error writing CSV header: %w", err)
        }
        
        // Write rows
        for _, chapter := range book.Chapters {
                row := []string{
                        fmt.Sprintf("%d", chapter.ID),
                        fmt.Sprintf("%d", chapter.BookID),
                        chapter.Title,
                        chapter.Description,
                        fmt.Sprintf("%d", chapter.Number),
                        fmt.Sprintf("%t", chapter.Published),
                        chapter.CreatedAt.Format(time.RFC3339),
                        chapter.UpdatedAt.Format(time.RFC3339),
                }
                if err := csvWriter.Write(row); err != nil {
                        return fmt.Errorf("error writing CSV row: %w", err)
                }
        }
        
        return nil
}

// ExportSectionsToCSV exports sections of a chapter to a CSV file
func (s *ContentAdminServiceImpl) ExportSectionsToCSV(writer io.Writer, chapterID uint) error {
        // Get chapter with sections
        chapter, err := s.chapterRepo.GetChapterWithSections(chapterID)
        if err != nil {
                return fmt.Errorf("error getting chapter with ID %d: %w", chapterID, err)
        }
        
        // Create CSV writer
        csvWriter := csv.NewWriter(writer)
        defer csvWriter.Flush()
        
        // Write header
        header := []string{"ID", "ChapterID", "BookID", "Title", "Number", "Content", "Format", "TimeToRead", "Published", "CreatedAt", "UpdatedAt"}
        if err := csvWriter.Write(header); err != nil {
                return fmt.Errorf("error writing CSV header: %w", err)
        }
        
        // Write rows
        for _, section := range chapter.Sections {
                row := []string{
                        fmt.Sprintf("%d", section.ID),
                        fmt.Sprintf("%d", section.ChapterID),
                        fmt.Sprintf("%d", section.BookID),
                        section.Title,
                        fmt.Sprintf("%d", section.Number),
                        section.Content,
                        section.Format,
                        fmt.Sprintf("%d", section.TimeToRead),
                        fmt.Sprintf("%t", section.Published),
                        section.CreatedAt.Format(time.RFC3339),
                        section.UpdatedAt.Format(time.RFC3339),
                }
                if err := csvWriter.Write(row); err != nil {
                        return fmt.Errorf("error writing CSV row: %w", err)
                }
        }
        
        return nil
}

// CreateBookRevision creates a new revision for a book
func (s *ContentAdminServiceImpl) CreateBookRevision(bookID uint, changes map[string]interface{}, notes string) (*models.BookRevision, error) {
        // Get current book
        book, err := s.bookRepo.GetBookByID(bookID)
        if err != nil {
                return nil, fmt.Errorf("error getting book with ID %d: %w", bookID, err)
        }
        
        // Create revision record with current content
        bookContent, err := json.Marshal(book)
        if err != nil {
                return nil, fmt.Errorf("error marshaling book: %w", err)
        }
        
        revision := &models.BookRevision{
                BookID:    bookID,
                Content:   string(bookContent),
                Notes:     notes,
                CreatedAt: time.Now(),
        }
        
        // Save revision
        if err := s.bookRepo.CreateBookRevision(revision); err != nil {
                return nil, fmt.Errorf("error saving book revision: %w", err)
        }
        
        // Apply changes to book
        for field, value := range changes {
                switch field {
                case "title":
                        if title, ok := value.(string); ok {
                                book.Title = title
                        }
                case "description":
                        if description, ok := value.(string); ok {
                                book.Description = description
                        }
                case "author":
                        if author, ok := value.(string); ok {
                                book.Author = author
                        }
                case "cover_image":
                        if coverImage, ok := value.(string); ok {
                                book.CoverImage = coverImage
                        }
                case "published":
                        if published, ok := value.(bool); ok {
                                book.Published = published
                        }
                }
        }
        
        // Update the book
        book.UpdatedAt = time.Now()
        if err := s.bookRepo.UpdateBook(book); err != nil {
                return nil, fmt.Errorf("error updating book: %w", err)
        }
        
        return revision, nil
}

// CreateChapterRevision creates a new revision for a chapter
func (s *ContentAdminServiceImpl) CreateChapterRevision(chapterID uint, changes map[string]interface{}, notes string) (*models.ChapterRevision, error) {
        // Get current chapter
        chapter, err := s.chapterRepo.GetChapterByID(chapterID)
        if err != nil {
                return nil, fmt.Errorf("error getting chapter with ID %d: %w", chapterID, err)
        }
        
        // Create revision record with current content
        chapterContent, err := json.Marshal(chapter)
        if err != nil {
                return nil, fmt.Errorf("error marshaling chapter: %w", err)
        }
        
        revision := &models.ChapterRevision{
                ChapterID: chapterID,
                Content:   string(chapterContent),
                Notes:     notes,
                CreatedAt: time.Now(),
        }
        
        // Save revision
        if err := s.chapterRepo.CreateChapterRevision(revision); err != nil {
                return nil, fmt.Errorf("error saving chapter revision: %w", err)
        }
        
        // Apply changes to chapter
        for field, value := range changes {
                switch field {
                case "title":
                        if title, ok := value.(string); ok {
                                chapter.Title = title
                        }
                case "description":
                        if description, ok := value.(string); ok {
                                chapter.Description = description
                        }
                case "number":
                        if number, ok := value.(int); ok {
                                chapter.Number = number
                        }
                case "published":
                        if published, ok := value.(bool); ok {
                                chapter.Published = published
                        }
                }
        }
        
        // Update the chapter
        chapter.UpdatedAt = time.Now()
        if err := s.chapterRepo.UpdateChapter(chapter); err != nil {
                return nil, fmt.Errorf("error updating chapter: %w", err)
        }
        
        return revision, nil
}

// CreateSectionRevision creates a new revision for a section
func (s *ContentAdminServiceImpl) CreateSectionRevision(sectionID uint, changes map[string]interface{}, notes string) (*models.SectionRevision, error) {
        // Get current section
        section, err := s.sectionRepo.GetSectionByID(sectionID)
        if err != nil {
                return nil, fmt.Errorf("error getting section with ID %d: %w", sectionID, err)
        }
        
        // Create revision record with current content
        sectionContent, err := json.Marshal(section)
        if err != nil {
                return nil, fmt.Errorf("error marshaling section: %w", err)
        }
        
        revision := &models.SectionRevision{
                SectionID: sectionID,
                Content:   string(sectionContent),
                Notes:     notes,
                CreatedAt: time.Now(),
        }
        
        // Save revision
        if err := s.sectionRepo.CreateSectionRevision(revision); err != nil {
                return nil, fmt.Errorf("error saving section revision: %w", err)
        }
        
        // Apply changes to section
        for field, value := range changes {
                switch field {
                case "title":
                        if title, ok := value.(string); ok {
                                section.Title = title
                        }
                case "content":
                        if content, ok := value.(string); ok {
                                section.Content = content
                        }
                case "number":
                        if number, ok := value.(int); ok {
                                section.Number = number
                        }
                case "format":
                        if format, ok := value.(string); ok {
                                section.Format = format
                        }
                case "time_to_read":
                        if timeToRead, ok := value.(int); ok {
                                section.TimeToRead = timeToRead
                        }
                case "published":
                        if published, ok := value.(bool); ok {
                                section.Published = published
                        }
                }
        }
        
        // Update the section
        section.UpdatedAt = time.Now()
        if err := s.sectionRepo.UpdateSection(section); err != nil {
                return nil, fmt.Errorf("error updating section: %w", err)
        }
        
        return revision, nil
}

// GetBookRevisions retrieves all revisions for a book
func (s *ContentAdminServiceImpl) GetBookRevisions(bookID uint) ([]models.BookRevision, error) {
        return s.bookRepo.GetBookRevisions(bookID)
}

// GetChapterRevisions retrieves all revisions for a chapter
func (s *ContentAdminServiceImpl) GetChapterRevisions(chapterID uint) ([]models.ChapterRevision, error) {
        return s.chapterRepo.GetChapterRevisions(chapterID)
}

// GetSectionRevisions retrieves all revisions for a section
func (s *ContentAdminServiceImpl) GetSectionRevisions(sectionID uint) ([]models.SectionRevision, error) {
        return s.sectionRepo.GetSectionRevisions(sectionID)
}

// RestoreBookRevision restores a book to a previous revision
func (s *ContentAdminServiceImpl) RestoreBookRevision(revisionID uint) error {
        // Get the revision
        revision, err := s.bookRepo.GetBookRevisionByID(revisionID)
        if err != nil {
                return fmt.Errorf("error getting book revision with ID %d: %w", revisionID, err)
        }
        
        // Parse the revision content
        var book models.Book
        if err := json.Unmarshal([]byte(revision.Content), &book); err != nil {
                return fmt.Errorf("error unmarshaling book revision: %w", err)
        }
        
        // Update revision date
        book.UpdatedAt = time.Now()
        
        // Update the book
        if err := s.bookRepo.UpdateBook(&book); err != nil {
                return fmt.Errorf("error restoring book from revision: %w", err)
        }
        
        return nil
}

// RestoreChapterRevision restores a chapter to a previous revision
func (s *ContentAdminServiceImpl) RestoreChapterRevision(revisionID uint) error {
        // Get the revision
        revision, err := s.chapterRepo.GetChapterRevisionByID(revisionID)
        if err != nil {
                return fmt.Errorf("error getting chapter revision with ID %d: %w", revisionID, err)
        }
        
        // Parse the revision content
        var chapter models.BookChapter
        if err := json.Unmarshal([]byte(revision.Content), &chapter); err != nil {
                return fmt.Errorf("error unmarshaling chapter revision: %w", err)
        }
        
        // Update revision date
        chapter.UpdatedAt = time.Now()
        
        // Update the chapter
        if err := s.chapterRepo.UpdateChapter(&chapter); err != nil {
                return fmt.Errorf("error restoring chapter from revision: %w", err)
        }
        
        return nil
}

// RestoreSectionRevision restores a section to a previous revision
func (s *ContentAdminServiceImpl) RestoreSectionRevision(revisionID uint) error {
        // Get the revision
        revision, err := s.sectionRepo.GetSectionRevisionByID(revisionID)
        if err != nil {
                return fmt.Errorf("error getting section revision with ID %d: %w", revisionID, err)
        }
        
        // Parse the revision content
        var section models.BookSection
        if err := json.Unmarshal([]byte(revision.Content), &section); err != nil {
                return fmt.Errorf("error unmarshaling section revision: %w", err)
        }
        
        // Update revision date
        section.UpdatedAt = time.Now()
        
        // Update the section
        if err := s.sectionRepo.UpdateSection(&section); err != nil {
                return fmt.Errorf("error restoring section from revision: %w", err)
        }
        
        return nil
}

// ScheduleBookPublishing schedules a book to be published at a future date
func (s *ContentAdminServiceImpl) ScheduleBookPublishing(bookID uint, publishDate time.Time) error {
        // Get the book
        book, err := s.bookRepo.GetBookByID(bookID)
        if err != nil {
                return fmt.Errorf("error getting book with ID %d: %w", bookID, err)
        }
        
        // Update the book's scheduled publish date
        book.ScheduledPublishAt = &publishDate
        book.UpdatedAt = time.Now()
        
        // Update the book
        if err := s.bookRepo.UpdateBook(book); err != nil {
                return fmt.Errorf("error scheduling book publishing: %w", err)
        }
        
        return nil
}

// ScheduleChapterPublishing schedules a chapter to be published at a future date
func (s *ContentAdminServiceImpl) ScheduleChapterPublishing(chapterID uint, publishDate time.Time) error {
        // Get the chapter
        chapter, err := s.chapterRepo.GetChapterByID(chapterID)
        if err != nil {
                return fmt.Errorf("error getting chapter with ID %d: %w", chapterID, err)
        }
        
        // Update the chapter's scheduled publish date
        chapter.ScheduledPublishAt = &publishDate
        chapter.UpdatedAt = time.Now()
        
        // Update the chapter
        if err := s.chapterRepo.UpdateChapter(chapter); err != nil {
                return fmt.Errorf("error scheduling chapter publishing: %w", err)
        }
        
        return nil
}

// ScheduleSectionPublishing schedules a section to be published at a future date
func (s *ContentAdminServiceImpl) ScheduleSectionPublishing(sectionID uint, publishDate time.Time) error {
        // Get the section
        section, err := s.sectionRepo.GetSectionByID(sectionID)
        if err != nil {
                return fmt.Errorf("error getting section with ID %d: %w", sectionID, err)
        }
        
        // Update the section's scheduled publish date
        section.ScheduledPublishAt = &publishDate
        section.UpdatedAt = time.Now()
        
        // Update the section
        if err := s.sectionRepo.UpdateSection(section); err != nil {
                return fmt.Errorf("error scheduling section publishing: %w", err)
        }
        
        return nil
}

// GetScheduledContent gets all content scheduled to be published
func (s *ContentAdminServiceImpl) GetScheduledContent() ([]interface{}, error) {
        scheduled := make([]interface{}, 0)
        
        // Get scheduled books
        books, err := s.bookRepo.GetScheduledBooks()
        if err != nil {
                return nil, fmt.Errorf("error getting scheduled books: %w", err)
        }
        for _, book := range books {
                scheduled = append(scheduled, map[string]interface{}{
                        "type":        "book",
                        "id":          book.ID,
                        "title":       book.Title,
                        "publishDate": book.ScheduledPublishAt,
                })
        }
        
        // Get scheduled chapters
        chapters, err := s.chapterRepo.GetScheduledChapters()
        if err != nil {
                return nil, fmt.Errorf("error getting scheduled chapters: %w", err)
        }
        for _, chapter := range chapters {
                scheduled = append(scheduled, map[string]interface{}{
                        "type":        "chapter",
                        "id":          chapter.ID,
                        "title":       chapter.Title,
                        "bookID":      chapter.BookID,
                        "publishDate": chapter.ScheduledPublishAt,
                })
        }
        
        // Get scheduled sections
        sections, err := s.sectionRepo.GetScheduledSections()
        if err != nil {
                return nil, fmt.Errorf("error getting scheduled sections: %w", err)
        }
        for _, section := range sections {
                scheduled = append(scheduled, map[string]interface{}{
                        "type":        "section",
                        "id":          section.ID,
                        "title":       section.Title,
                        "bookID":      section.BookID,
                        "chapterID":   section.ChapterID,
                        "publishDate": section.ScheduledPublishAt,
                })
        }
        
        return scheduled, nil
}

// PublishContent publishes content by setting its published flag to true
func (s *ContentAdminServiceImpl) PublishContent(contentType string, contentID uint) error {
        switch contentType {
        case "book":
                // Get the book
                book, err := s.bookRepo.GetBookByID(contentID)
                if err != nil {
                        return fmt.Errorf("error getting book with ID %d: %w", contentID, err)
                }
                
                // Update the book
                book.Published = true
                book.ScheduledPublishAt = nil
                book.UpdatedAt = time.Now()
                
                if err := s.bookRepo.UpdateBook(book); err != nil {
                        return fmt.Errorf("error publishing book: %w", err)
                }
                
        case "chapter":
                // Get the chapter
                chapter, err := s.chapterRepo.GetChapterByID(contentID)
                if err != nil {
                        return fmt.Errorf("error getting chapter with ID %d: %w", contentID, err)
                }
                
                // Update the chapter
                chapter.Published = true
                chapter.ScheduledPublishAt = nil
                chapter.UpdatedAt = time.Now()
                
                if err := s.chapterRepo.UpdateChapter(chapter); err != nil {
                        return fmt.Errorf("error publishing chapter: %w", err)
                }
                
        case "section":
                // Get the section
                section, err := s.sectionRepo.GetSectionByID(contentID)
                if err != nil {
                        return fmt.Errorf("error getting section with ID %d: %w", contentID, err)
                }
                
                // Update the section
                section.Published = true
                section.ScheduledPublishAt = nil
                section.UpdatedAt = time.Now()
                
                if err := s.sectionRepo.UpdateSection(section); err != nil {
                        return fmt.Errorf("error publishing section: %w", err)
                }
                
        default:
                return fmt.Errorf("invalid content type: %s", contentType)
        }
        
        return nil
}

// UnpublishContent unpublishes content by setting its published flag to false
func (s *ContentAdminServiceImpl) UnpublishContent(contentType string, contentID uint) error {
        switch contentType {
        case "book":
                // Get the book
                book, err := s.bookRepo.GetBookByID(contentID)
                if err != nil {
                        return fmt.Errorf("error getting book with ID %d: %w", contentID, err)
                }
                
                // Update the book
                book.Published = false
                book.UpdatedAt = time.Now()
                
                if err := s.bookRepo.UpdateBook(book); err != nil {
                        return fmt.Errorf("error unpublishing book: %w", err)
                }
                
        case "chapter":
                // Get the chapter
                chapter, err := s.chapterRepo.GetChapterByID(contentID)
                if err != nil {
                        return fmt.Errorf("error getting chapter with ID %d: %w", contentID, err)
                }
                
                // Update the chapter
                chapter.Published = false
                chapter.UpdatedAt = time.Now()
                
                if err := s.chapterRepo.UpdateChapter(chapter); err != nil {
                        return fmt.Errorf("error unpublishing chapter: %w", err)
                }
                
        case "section":
                // Get the section
                section, err := s.sectionRepo.GetSectionByID(contentID)
                if err != nil {
                        return fmt.Errorf("error getting section with ID %d: %w", contentID, err)
                }
                
                // Update the section
                section.Published = false
                section.UpdatedAt = time.Now()
                
                if err := s.sectionRepo.UpdateSection(section); err != nil {
                        return fmt.Errorf("error unpublishing section: %w", err)
                }
                
        default:
                return fmt.Errorf("invalid content type: %s", contentType)
        }
        
        return nil
}