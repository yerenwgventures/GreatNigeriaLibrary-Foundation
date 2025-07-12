package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Foundation services
	authHandlers "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/auth/handlers"
	authRepository "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/auth/repository"
	authService "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/auth/service"

	// Shared packages
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/auth"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/config"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/database"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/middleware"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/redis"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	appLogger := logger.New(logger.ParseLogLevel(cfg.Logging.Level))

	// Initialize database
	dbConfig := &database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Username:        cfg.Database.Username,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.Database,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}

	dbConn, err := database.NewConnection(dbConfig)
	if err != nil {
		appLogger.Fatal("Failed to connect to database: " + err.Error())
	}
	defer dbConn.Close()

	// Run comprehensive GORM auto-migration for all foundation services
	migrationService := database.NewMigrationService(dbConn.DB, appLogger)
	if err := migrationService.MigrateFoundation(); err != nil {
		appLogger.Fatal("Failed to run foundation migrations: " + err.Error())
	}

	// Initialize repositories
	userRepo := authRepository.NewUserRepository(dbConn.DB, appLogger)
	// Add other foundation repositories

	// Initialize services
	userSvc := authService.NewUserService(userRepo, appLogger)
	// Add other foundation services

	// Initialize enhanced JWT manager and authorization manager
	jwtManager := userSvc.GetJWTManager()
	authManager := auth.NewAuthorizationManager()

	// Initialize handlers
	userHandler := authHandlers.NewUserHandler(userSvc, appLogger)
	// Add other foundation handlers

	// Setup Gin router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger(appLogger))
	router.Use(middleware.Recovery(appLogger))
	router.Use(middleware.SecurityHeaders())

	// API routes
	api := router.Group("/api/v1")

	// Foundation routes only
	setupFoundationRoutes(api, userHandler, jwtManager, authManager, appLogger)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "great-nigeria-library-foundation",
			"version": "1.0.0",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appLogger.Info("Starting Great Nigeria Library Foundation on port " + port)
	if err := router.Run(":" + port); err != nil {
		appLogger.Fatal("Failed to start server: " + err.Error())
	}
}

func setupFoundationRoutes(api *gin.RouterGroup, userHandler *authHandlers.UserHandler, jwtManager *auth.JWTManager, authManager *auth.AuthorizationManager, logger *logger.Logger) {
	// Public auth routes - no authentication required
	authPublic := api.Group("/auth")
	{
		authPublic.POST("/register", userHandler.Register)
		authPublic.POST("/login", userHandler.Login)
	}

	// Protected auth routes - require authentication
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.AuthRequired(jwtManager, logger))
	{
		authProtected.POST("/logout", userHandler.Logout)
		authProtected.GET("/profile", userHandler.GetProfile)
		authProtected.PUT("/profile",
			middleware.PermissionRequired(authManager, auth.PermissionUpdateProfile, logger),
			userHandler.UpdateProfile)
	}

	// Public content routes - no authentication required
	contentPublic := api.Group("/content/public")
	{
		contentPublic.GET("/books", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Demo books endpoint - Foundation version (Public)",
				"books": []gin.H{
					{
						"id":    1,
						"title": "Platform User Guide",
						"type":  "demo",
						"access": "public",
					},
				},
			})
		})
	}

	// Protected content routes - require authentication and content permissions
	contentProtected := api.Group("/content")
	contentProtected.Use(middleware.AuthRequired(jwtManager, logger))
	contentProtected.Use(middleware.PermissionRequired(authManager, auth.PermissionReadContent, logger))
	{
		contentProtected.GET("/books", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Demo books endpoint - Foundation version (Protected)",
				"books": []gin.H{
					{
						"id":    1,
						"title": "Platform User Guide",
						"type":  "demo",
						"access": "protected",
					},
					{
						"id":    2,
						"title": "Nigerian History Overview",
						"type":  "educational",
						"access": "protected",
					},
					{
						"id":    3,
						"title": "Advanced Research Methods",
						"type":  "premium",
						"access": "protected",
					},
				},
			})
		})
	}

	// Admin content routes - require admin permissions
	contentAdmin := api.Group("/content/admin")
	contentAdmin.Use(middleware.AuthRequired(jwtManager, logger))
	contentAdmin.Use(middleware.RoleRequired(int(auth.RoleAdmin), logger))
	{
		contentAdmin.GET("/stats", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Content statistics - Admin only",
				"stats": gin.H{
					"total_books": 3,
					"public_books": 1,
					"protected_books": 2,
					"total_views": 1250,
				},
			})
		})
		contentAdmin.POST("/books",
			middleware.PermissionRequired(authManager, auth.PermissionCreateContent, logger),
			func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Content creation not yet implemented"})
			})
	}

	// Discussion routes
	discussion := api.Group("/discussion")
	{
		discussion.GET("/forums", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Discussion forums - Foundation version",
			})
		})
	}
}
