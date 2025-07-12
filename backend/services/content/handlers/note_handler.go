package handlers

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/service"
)

// NoteHandler defines handlers for note-related endpoints
type NoteHandler struct {
        noteService service.NoteService
}

// NewNoteHandler creates a new note handler instance
func NewNoteHandler(noteService service.NoteService) *NoteHandler {
        return &NoteHandler{
                noteService: noteService,
        }
}

// CreateNote handles the POST /books/:id/notes endpoint
func (h *NoteHandler) CreateNote(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        bookIDStr := c.Param("id")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        // Parse request body
        var noteRequest struct {
                ChapterID uint   `json:"chapterId"`
                SectionID uint   `json:"sectionId"`
                Content   string `json:"content"`
                Color     string `json:"color"`
                Position  int    `json:"position"`
                Tags      string `json:"tags"`
        }

        if err := c.ShouldBindJSON(&noteRequest); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        // Create note object
        note := &models.BookNote{
                UserID:    userID.(uint),
                BookID:    uint(bookID),
                ChapterID: noteRequest.ChapterID,
                SectionID: noteRequest.SectionID,
                Content:   noteRequest.Content,
                Color:     noteRequest.Color,
                Position:  noteRequest.Position,
                Tags:      noteRequest.Tags,
        }

        if err := h.noteService.CreateNote(note); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "Note created successfully", "data": note})
}

// GetNotes handles the GET /books/:id/notes endpoint
func (h *NoteHandler) GetNotes(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        bookIDStr := c.Param("id")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        notes, err := h.noteService.GetNotes(userID.(uint), uint(bookID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": notes})
}

// GetSectionNotes handles the GET /books/:id/sections/:sectionId/notes endpoint
func (h *NoteHandler) GetSectionNotes(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        sectionIDStr := c.Param("sectionId")
        sectionID, err := strconv.ParseUint(sectionIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid section ID"})
                return
        }

        notes, err := h.noteService.GetNotesBySection(userID.(uint), uint(sectionID))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"data": notes})
}

// UpdateNote handles the PUT /notes/:id endpoint
func (h *NoteHandler) UpdateNote(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        noteIDStr := c.Param("id")
        noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
                return
        }

        // First, get the existing note to verify ownership
        existingNote, err := h.noteService.GetNoteByID(uint(noteID), userID.(uint))
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or doesn't belong to user"})
                return
        }

        // Parse request body
        var updateRequest struct {
                Content string `json:"content"`
                Color   string `json:"color"`
                Tags    string `json:"tags"`
        }

        if err := c.ShouldBindJSON(&updateRequest); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
                return
        }

        // Update fields
        existingNote.Content = updateRequest.Content
        existingNote.Color = updateRequest.Color
        existingNote.Tags = updateRequest.Tags

        // Save the update
        if err := h.noteService.UpdateNote(existingNote); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully", "data": existingNote})
}

// DeleteNote handles the DELETE /notes/:id endpoint
func (h *NoteHandler) DeleteNote(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        noteIDStr := c.Param("id")
        noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
                return
        }

        if err := h.noteService.DeleteNote(uint(noteID), userID.(uint)); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

// ExportNotes handles the GET /books/:id/notes/export endpoint
func (h *NoteHandler) ExportNotes(c *gin.Context) {
        // Extract user ID from JWT token
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        bookIDStr := c.Param("id")
        bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
                return
        }

        format := c.DefaultQuery("format", "txt")
        // Validate format
        if format != "txt" && format != "json" && format != "markdown" && format != "md" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export format. Use txt, json, or markdown."})
                return
        }

        exportContent, err := h.noteService.ExportNotes(userID.(uint), uint(bookID), format)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export notes"})
                return
        }

        // Set appropriate Content-Type based on format
        var contentType string
        switch format {
        case "json":
                contentType = "application/json"
        case "markdown", "md":
                contentType = "text/markdown"
        default:
                contentType = "text/plain"
        }

        // Set filename
        filename := "notes-export." + format
        c.Header("Content-Disposition", "attachment; filename="+filename)
        c.Data(http.StatusOK, contentType, []byte(exportContent))
}

// RegisterRoutes registers the note-related routes
func (h *NoteHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
        // All note routes require authentication
        notes := router.Group("/api/v1")
        notes.Use(authMiddleware)
        {
                // Book-specific notes
                notes.POST("/books/:id/notes", h.CreateNote)
                notes.GET("/books/:id/notes", h.GetNotes)
                notes.GET("/books/:id/notes/export", h.ExportNotes)
                notes.GET("/books/:id/sections/:sectionId/notes", h.GetSectionNotes)
                
                // General note management
                notes.PUT("/notes/:id", h.UpdateNote)
                notes.DELETE("/notes/:id", h.DeleteNote)
        }
}