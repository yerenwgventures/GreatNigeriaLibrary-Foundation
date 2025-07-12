package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/handlers"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/repository"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/discussion/service"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/auth"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/config"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/database"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/middleware"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/redis"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Starting Great Nigeria Discussion Service")

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

	// Run GORM auto-migration for discussion service
	migrationService := database.NewMigrationService(db, logger)
	if err := migrationService.MigrateDiscussionService(); err != nil {
		logger.Fatal("Failed to run discussion service migrations: " + err.Error())
	}

	// Initialize repositories
	discussionRepo := repository.NewDiscussionRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// Initialize services
	discussionService := service.NewDiscussionService(discussionRepo, logger)
	commentService := service.NewCommentService(commentRepo, logger)
	likeService := service.NewLikeService(likeRepo, logger)

	// Initialize handlers
	discussionHandler := handlers.NewDiscussionHandler(discussionService, logger)
	commentHandler := handlers.NewCommentHandler(commentService, logger)
	likeHandler := handlers.NewLikeHandler(likeService, logger)

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

	// Public routes - can be accessed without authentication
	public := router.Group("/public")
	{
		public.GET("/discussions",
			middleware.PermissionRequired(authManager, auth.PermissionReadDiscussion, logger),
			discussionHandler.ListDiscussions)
		public.GET("/discussions/:id",
			middleware.PermissionRequired(authManager, auth.PermissionReadDiscussion, logger),
			discussionHandler.GetDiscussion)
		public.GET("/discussions/:id/comments",
			middleware.PermissionRequired(authManager, auth.PermissionReadDiscussion, logger),
			commentHandler.ListComments)
	}

	// Discussion reading routes - require authentication and read permission
	discussions := router.Group("/discussions")
	discussions.Use(middleware.AuthRequired(jwtManager, logger))
	discussions.Use(middleware.PermissionRequired(authManager, auth.PermissionReadDiscussion, logger))
	{
		discussions.GET("/", discussionHandler.ListDiscussions)
		discussions.GET("/:id", discussionHandler.GetDiscussion)
		discussions.GET("/:id/comments", commentHandler.ListComments)
	}

	// Discussion participation routes - require authentication and create permission
	participate := router.Group("/participate")
	participate.Use(middleware.AuthRequired(jwtManager, logger))
	{
		participate.POST("/discussions",
			middleware.PermissionRequired(authManager, auth.PermissionCreateDiscussion, logger),
			discussionHandler.CreateDiscussion)
		participate.PATCH("/discussions/:id",
			middleware.ResourceOwnerOrPermission(authManager, auth.PermissionUpdateDiscussion, "discussion_owner_id", logger),
			discussionHandler.UpdateDiscussion)
		participate.DELETE("/discussions/:id",
			middleware.ResourceOwnerOrPermission(authManager, auth.PermissionDeleteDiscussion, "discussion_owner_id", logger),
			discussionHandler.DeleteDiscussion)
		participate.POST("/discussions/:id/comments",
			middleware.PermissionRequired(authManager, auth.PermissionCreateDiscussion, logger),
			commentHandler.CreateComment)
		participate.PATCH("/comments/:id",
			middleware.ResourceOwnerOrPermission(authManager, auth.PermissionUpdateDiscussion, "comment_owner_id", logger),
			commentHandler.UpdateComment)
		participate.DELETE("/comments/:id",
			middleware.ResourceOwnerOrPermission(authManager, auth.PermissionDeleteDiscussion, "comment_owner_id", logger),
			commentHandler.DeleteComment)
		participate.POST("/discussions/:id/like",
			middleware.PermissionRequired(authManager, auth.PermissionReadDiscussion, logger),
			likeHandler.LikeDiscussion)
		participate.DELETE("/discussions/:id/like",
			middleware.PermissionRequired(authManager, auth.PermissionReadDiscussion, logger),
			likeHandler.UnlikeDiscussion)
		participate.POST("/comments/:id/like",
			middleware.PermissionRequired(authManager, auth.PermissionReadDiscussion, logger),
			likeHandler.LikeComment)
	}

	// Moderation routes - require moderator role and moderation permissions
	moderate := router.Group("/moderate")
	moderate.Use(middleware.AuthRequired(jwtManager, logger))
	moderate.Use(middleware.RoleRequired(int(auth.RoleModerator), logger))
	{
		moderate.DELETE("/discussions/:id",
			middleware.PermissionRequired(authManager, auth.PermissionModerateDiscussion, logger),
			discussionHandler.DeleteDiscussion)
		moderate.DELETE("/comments/:id",
			middleware.PermissionRequired(authManager, auth.PermissionModerateDiscussion, logger),
			commentHandler.DeleteComment)
		moderate.POST("/discussions/:id/lock",
			middleware.PermissionRequired(authManager, auth.PermissionModerateDiscussion, logger),
			func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Discussion locking not yet implemented"})
			})
		moderate.POST("/discussions/:id/pin",
			middleware.PermissionRequired(authManager, auth.PermissionModerateDiscussion, logger),
			func(c *gin.Context) {
				c.JSON(501, gin.H{"message": "Discussion pinning not yet implemented"})
			})
	}

	// Admin discussion routes - require admin permissions
	admin := router.Group("/admin")
	admin.Use(middleware.AuthRequired(jwtManager, logger))
	admin.Use(middleware.RoleRequired(int(auth.RoleAdmin), logger))
	{
		admin.GET("/discussions/stats", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Discussion stats not yet implemented"})
		})
		admin.GET("/discussions/reports", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "Discussion reports not yet implemented"})
		})
	}
		protected.DELETE("/comments/:id/like", likeHandler.UnlikeComment)
	}

	// Start server
	port := os.Getenv("DISCUSSION_SERVICE_PORT")
	if port == "" {
		port = "8003" // Default port for discussion service
	}

	logger.Info("Discussion Service starting on port " + port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		logger.Fatal("Failed to start Discussion Service: " + err.Error())
	}
}
