package service

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/repository"
)

// NoteService defines interface for note-related business logic
type NoteService interface {
        CreateNote(note *models.BookNote) error
        GetNotes(userID, bookID uint) ([]models.BookNote, error)
        GetNotesBySection(userID, sectionID uint) ([]models.BookNote, error)
        GetNoteByID(id, userID uint) (*models.BookNote, error)
        UpdateNote(note *models.BookNote) error
        DeleteNote(id, userID uint) error
        ExportNotes(userID, bookID uint, format string) (string, error)
}

// NoteServiceImpl implements the NoteService interface
type NoteServiceImpl struct {
        noteRepo repository.NoteRepository
}

// NewNoteService creates a new note service instance
func NewNoteService(noteRepo repository.NoteRepository) NoteService {
        return &NoteServiceImpl{
                noteRepo: noteRepo,
        }
}

// CreateNote creates a new note
func (s *NoteServiceImpl) CreateNote(note *models.BookNote) error {
        return s.noteRepo.CreateNote(note)
}

// GetNotes retrieves all notes for a user and book
func (s *NoteServiceImpl) GetNotes(userID, bookID uint) ([]models.BookNote, error) {
        return s.noteRepo.GetNotes(userID, bookID)
}

// GetNotesBySection retrieves notes for a specific section
func (s *NoteServiceImpl) GetNotesBySection(userID, sectionID uint) ([]models.BookNote, error) {
        return s.noteRepo.GetNotesBySection(userID, sectionID)
}

// GetNoteByID retrieves a note by ID
func (s *NoteServiceImpl) GetNoteByID(id, userID uint) (*models.BookNote, error) {
        return s.noteRepo.GetNoteByID(id, userID)
}

// UpdateNote updates an existing note
func (s *NoteServiceImpl) UpdateNote(note *models.BookNote) error {
        return s.noteRepo.UpdateNote(note)
}

// DeleteNote deletes a note
func (s *NoteServiceImpl) DeleteNote(id, userID uint) error {
        return s.noteRepo.DeleteNote(id, userID)
}

// ExportNotes exports notes in the specified format
func (s *NoteServiceImpl) ExportNotes(userID, bookID uint, format string) (string, error) {
        return s.noteRepo.ExportNotes(userID, bookID, format)
}