package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/handlers"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/repository"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/service"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/config"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/database"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting Great Nigeria Content Service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration: " + err.Error())
	}

	// Connect to database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}

	// Initialize repositories
	bookRepo := repository.NewBookRepository(db)
	progressRepo := repository.NewProgressRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)
	feedbackRepo := repository.NewFeedbackRepository(db, logger)
	noteRepo := repository.NewNoteRepository(db, logger)

	// Initialize services
	bookService := service.NewBookService(bookRepo, progressRepo, logger)
	progressService := service.NewProgressService(progressRepo, logger)
	bookmarkService := service.NewBookmarkService(bookmarkRepo, logger)
	feedbackService := service.NewFeedbackService(feedbackRepo, bookRepo, logger)
	noteService := service.NewNoteService(noteRepo, bookRepo, logger)
	bookImportService := service.NewBookImportService(bookRepo, logger)
	contentRenderer := service.NewContentRenderer(bookRepo)
	mediaGenerator := service.NewMediaGenerator(bookRepo, contentRenderer, "http://localhost:5000", "./static/media")

	// Initialize handlers
	bookHandlers := handlers.NewBookHandlers(bookService, bookImportService)
	progressHandler := handlers.NewProgressHandler(progressService, logger)
	bookmarkHandler := handlers.NewBookmarkHandler(bookmarkService, logger)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService, logger)
	noteHandler := handlers.NewNoteHandler(noteService, logger)
	quizHandler := handlers.NewQuizHandler()
	mediaHandler := handlers.NewMediaHandler(mediaGenerator)

	// Set up Gin router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())

	// Define API routes
	router.GET("/books", bookHandlers.GetAllBooks)
	router.GET("/books/:id", bookHandlers.GetBookByID)
	router.GET("/books/:id/chapters", bookHandlers.GetBookChapters)
	router.GET("/books/chapters/:id", bookHandlers.GetChapterByID)
	router.GET("/books/sections/:id", bookHandlers.GetSectionByID)

	// Register interactive elements routes
	apiGroup := router.Group("/api")
	quizHandler.RegisterRoutes(apiGroup)
	mediaHandler.RegisterRoutes(apiGroup)

	// Get content feedback summary for anyone
	router.GET("/feedback/summary", feedbackHandler.GetContentFeedbackSummary)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// Book progress routes
		protected.POST("/books/:id/progress", progressHandler.UpdateProgress)
		protected.GET("/books/:id/progress", progressHandler.GetProgress)

		// Bookmark routes
		protected.POST("/books/:id/bookmarks", bookmarkHandler.CreateBookmark)
		protected.GET("/books/:id/bookmarks", bookmarkHandler.GetBookmarks)
		protected.DELETE("/books/:id/bookmarks/:bookmarkId", bookmarkHandler.DeleteBookmark)

		// Note routes
		protected.POST("/books/:id/notes", noteHandler.CreateNote)
		protected.GET("/books/:id/notes", noteHandler.GetNotes)
		protected.GET("/notes/:noteId", noteHandler.GetNoteByID)
		protected.PUT("/notes/:noteId", noteHandler.UpdateNote)
		protected.DELETE("/notes/:noteId", noteHandler.DeleteNote)
		protected.GET("/notes/categories", noteHandler.GetNoteCategories)
		protected.POST("/notes/export", noteHandler.ExportNotes)

		// Feedback routes
		protected.POST("/feedback/mood", feedbackHandler.SubmitMoodFeedback)
		protected.POST("/feedback/difficulty", feedbackHandler.SubmitDifficultyFeedback)
		protected.GET("/feedback/user", feedbackHandler.GetUserContentFeedback)
		protected.DELETE("/feedback/mood/:id", feedbackHandler.DeleteMoodFeedback)
		protected.DELETE("/feedback/difficulty/:id", feedbackHandler.DeleteDifficultyFeedback)
	}

	// Start server
	port := os.Getenv("CONTENT_SERVICE_PORT")
	if port == "" {
		port = "8002" // Default port for content service
	}

	logger.Info("Content Service starting on port " + port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		logger.Fatal("Failed to start Content Service: " + err.Error())
	}
}
