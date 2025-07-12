package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/handlers"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/repository"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/service"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/auth"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/config"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/redis"
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

	// Load configuration (YAML with environment overrides)
	cfg, err := config.LoadFromYAML("config.yaml")
	if err != nil {
		// Fallback to environment-only configuration
		logger.Info("YAML config not found, using environment variables only")
		cfg, err = config.LoadConfig()
		if err != nil {
			logger.Fatal("Failed to load configuration: " + err.Error())
		}
	}

	// Connect to database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}

	// Run GORM auto-migration for content service
	migrationService := database.NewMigrationService(db, logger)
	if err := migrationService.MigrateContentService(); err != nil {
		logger.Fatal("Failed to run content service migrations: " + err.Error())
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

	// Initialize enhanced JWT manager and authorization manager
	var jwtManager *auth.JWTManager
	var authManager *auth.AuthorizationManager

	// Initialize Redis client if enabled
	if cfg.Redis.Enabled {
		redisConfig := &redis.Config{
			Host:         cfg.Redis.Host,
			Port:         cfg.Redis.Port,
			Password:     cfg.Redis.Password,
			Database:     cfg.Redis.Database,
			PoolSize:     cfg.Redis.PoolSize,
			MinIdleConns: cfg.Redis.MinIdleConns,
			MaxRetries:   cfg.Redis.MaxRetries,
			DialTimeout:  cfg.Redis.DialTimeout,
			ReadTimeout:  cfg.Redis.ReadTimeout,
			WriteTimeout: cfg.Redis.WriteTimeout,
		}

		redisClient, err := redis.NewClient(redisConfig)
		if err != nil {
			logger.WithError(err).Warn("Failed to connect to Redis, JWT features will be limited")
			jwtManager = auth.NewJWTManagerWithoutRedis(
				cfg.Auth.JWTSecret,
				cfg.Auth.AccessTokenExpiration,
				cfg.Auth.RefreshTokenExpiration,
				cfg.Auth.JWTIssuer,
			)
		} else {
			jwtManager = auth.NewJWTManager(
				cfg.Auth.JWTSecret,
				cfg.Auth.AccessTokenExpiration,
				cfg.Auth.RefreshTokenExpiration,
				redisClient.Client,
				cfg.Auth.JWTIssuer,
			)
			logger.Info("Redis connected successfully for JWT token management")
		}
	} else {
		jwtManager = auth.NewJWTManagerWithoutRedis(
			cfg.Auth.JWTSecret,
			cfg.Auth.AccessTokenExpiration,
			cfg.Auth.RefreshTokenExpiration,
			cfg.Auth.JWTIssuer,
		)
	}

	authManager = auth.NewAuthorizationManager()

	// Set up Gin router with centralized error handling
	router := gin.New()

	// Add centralized error handling middleware
	router.Use(middleware.PanicRecovery(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.RequestLogger())
	router.Use(middleware.SecurityHeaders())

	// Public API routes - no authentication required
	public := router.Group("/public")
	{
		public.GET("/books", bookHandlers.GetAllBooks)
		public.GET("/books/:id", bookHandlers.GetBookByID)
		public.GET("/books/:id/chapters", bookHandlers.GetBookChapters)
		public.GET("/books/chapters/:id", bookHandlers.GetChapterByID)
		public.GET("/books/sections/:id", bookHandlers.GetSectionByID)
		public.GET("/feedback/summary", feedbackHandler.GetContentFeedbackSummary)
	}

	// Content reading routes - require basic authentication
	content := router.Group("/content")
	content.Use(middleware.AuthRequired(jwtManager, logger))
	{
		content.GET("/books",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			bookHandlers.GetAllBooks)
		content.GET("/books/:id",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			bookHandlers.GetBookByID)
		content.GET("/books/:id/chapters",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			bookHandlers.GetBookChapters)
		content.GET("/books/chapters/:id",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			bookHandlers.GetChapterByID)
		content.GET("/books/sections/:id",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			bookHandlers.GetSectionByID)
	}

	// Register interactive elements routes
	apiGroup := router.Group("/api")
	quizHandler.RegisterRoutes(apiGroup)
	mediaHandler.RegisterRoutes(apiGroup)

	// User content interaction routes - require authentication and content permissions
	userContent := router.Group("/user")
	userContent.Use(middleware.AuthRequired(jwtManager, logger))
	{
		// Book progress routes
		userContent.POST("/books/:id/progress",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			progressHandler.UpdateProgress)
		userContent.GET("/books/:id/progress",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			progressHandler.GetProgress)

		// Bookmark routes
		userContent.POST("/books/:id/bookmarks",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			bookmarkHandler.CreateBookmark)
		userContent.GET("/books/:id/bookmarks",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			bookmarkHandler.GetBookmarks)
		userContent.DELETE("/books/:id/bookmarks/:bookmarkId",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			bookmarkHandler.DeleteBookmark)

		// Note routes
		userContent.POST("/books/:id/notes",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			noteHandler.CreateNote)
		userContent.GET("/books/:id/notes",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			noteHandler.GetNotes)
		userContent.GET("/notes/:noteId",
			middleware.ResourceOwnerOrPermission(authManager, auth.PermissionReadContent, "note_owner_id", logger),
			noteHandler.GetNoteByID)
		userContent.PUT("/notes/:noteId",
			middleware.ResourceOwnerOrPermission(authManager, auth.PermissionUpdateContent, "note_owner_id", logger),
			noteHandler.UpdateNote)
		userContent.DELETE("/notes/:noteId",
			middleware.ResourceOwnerOrPermission(authManager, auth.PermissionDeleteContent, "note_owner_id", logger),
			noteHandler.DeleteNote)
		userContent.GET("/notes/categories",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			noteHandler.GetNoteCategories)
		userContent.POST("/notes/export",
			middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger),
			noteHandler.ExportNotes)
	}

	// Content creation routes - require content creation permissions
	create := router.Group("/create")
	create.Use(middleware.AuthRequired(jwtManager, logger))
	create.Use(middleware.PermissionRequired(authManager, auth.PermissionCreateContent, logger))
	{
		// Future content creation endpoints will go here
		create.POST("/books", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Content creation not yet implemented"})
		})
	}

	// Content management routes - require content management permissions
	manage := router.Group("/manage")
	manage.Use(middleware.AuthRequired(jwtManager, logger))
	manage.Use(middleware.AnyPermissionRequired(authManager, []auth.Permission{
		auth.PermissionUpdateContent,
		auth.PermissionDeleteContent,
		auth.PermissionPublishContent,
	}, logger))
	{
		// Future content management endpoints will go here
		manage.PUT("/books/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Content management not yet implemented"})
		})
		manage.DELETE("/books/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Content management not yet implemented"})
		})
	}

	// Admin content routes - require admin permissions
	admin := router.Group("/admin")
	admin.Use(middleware.AuthRequired(jwtManager, logger))
	admin.Use(middleware.RoleRequired(int(auth.RoleAdmin), logger))
	{
		admin.GET("/content/stats", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Admin content stats not yet implemented"})
		})
		admin.POST("/content/publish/:id",
			middleware.PermissionRequired(authManager, auth.PermissionPublishContent, logger),
			func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Content publishing not yet implemented"})
			})
	}

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
