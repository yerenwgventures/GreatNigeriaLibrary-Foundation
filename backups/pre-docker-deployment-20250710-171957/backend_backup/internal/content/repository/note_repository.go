package repository

import (
        "errors"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
        "gorm.io/gorm"
)

// NoteRepository defines the interface for note operations
type NoteRepository interface {
        CreateNote(note *models.BookNote) error
        GetNotes(userID, bookID uint) ([]models.BookNote, error)
        GetNotesBySection(userID, sectionID uint) ([]models.BookNote, error)
        GetNoteByID(id, userID uint) (*models.BookNote, error)
        UpdateNote(note *models.BookNote) error
        DeleteNote(id, userID uint) error
        ExportNotes(userID, bookID uint, format string) (string, error)
}

// GormNoteRepository implements the NoteRepository interface with GORM
type GormNoteRepository struct {
        db *gorm.DB
}

// NewGormNoteRepository creates a new note repository instance
func NewGormNoteRepository(db *gorm.DB) *GormNoteRepository {
        return &GormNoteRepository{db: db}
}

// CreateNote creates a new note
func (r *GormNoteRepository) CreateNote(note *models.BookNote) error {
        return r.db.Create(note).Error
}

// GetNotes retrieves all notes for a user and book
func (r *GormNoteRepository) GetNotes(userID, bookID uint) ([]models.BookNote, error) {
        var notes []models.BookNote
        result := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Order("created_at DESC").Find(&notes)
        return notes, result.Error
}

// GetNotesBySection retrieves notes for a specific section
func (r *GormNoteRepository) GetNotesBySection(userID, sectionID uint) ([]models.BookNote, error) {
        var notes []models.BookNote
        result := r.db.Where("user_id = ? AND section_id = ?", userID, sectionID).Order("position ASC").Find(&notes)
        return notes, result.Error
}

// GetNoteByID retrieves a note by ID, ensuring it belongs to the user
func (r *GormNoteRepository) GetNoteByID(id, userID uint) (*models.BookNote, error) {
        var note models.BookNote
        result := r.db.Where("id = ? AND user_id = ?", id, userID).First(&note)
        if result.Error != nil {
                if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                        return nil, errors.New("note not found")
                }
                return nil, result.Error
        }
        return &note, nil
}

// UpdateNote updates an existing note
func (r *GormNoteRepository) UpdateNote(note *models.BookNote) error {
        // Verify the note exists and belongs to the user
        _, err := r.GetNoteByID(note.ID, note.UserID)
        if err != nil {
                return err
        }
        
        return r.db.Save(note).Error
}

// DeleteNote deletes a note
func (r *GormNoteRepository) DeleteNote(id, userID uint) error {
        result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.BookNote{})
        if result.RowsAffected == 0 {
                return errors.New("note not found or does not belong to user")
        }
        return result.Error
}

// ExportNotes exports notes in the specified format
func (r *GormNoteRepository) ExportNotes(userID, bookID uint, format string) (string, error) {
        // Get all notes
        notes, err := r.GetNotes(userID, bookID)
        if err != nil {
                return "", err
        }
        
        // Get book information
        var book models.Book
        if err := r.db.First(&book, bookID).Error; err != nil {
                return "", err
        }
        
        switch format {
        case "json":
                return r.exportNotesAsJSON(notes, book)
        case "txt":
                return r.exportNotesAsText(notes, book)
        case "markdown", "md":
                return r.exportNotesAsMarkdown(notes, book)
        default:
                return r.exportNotesAsText(notes, book) // Default to text
        }
}

// Helper functions for exporting notes in different formats
func (r *GormNoteRepository) exportNotesAsJSON(notes []models.BookNote, book models.Book) (string, error) {
        // Implementation would use json.Marshal to create a JSON string
        // This would typically be implemented with a proper JSON serialization
        return "{ \"notes\": [] }", nil // Simplified implementation
}

func (r *GormNoteRepository) exportNotesAsText(notes []models.BookNote, book models.Book) (string, error) {
        // Implementation would build a text representation of notes
        // This would typically build a proper formatted string with all notes
        return "Notes Export", nil // Simplified implementation
}

func (r *GormNoteRepository) exportNotesAsMarkdown(notes []models.BookNote, book models.Book) (string, error) {
        // Implementation would build a markdown representation of notes
        // This would typically build proper markdown with formatting
        return "# Notes\n\n", nil // Simplified implementation
}