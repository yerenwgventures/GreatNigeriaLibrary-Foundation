package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/handlers"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/repository"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/discussion/service"
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
	logger.Info("Starting Great Nigeria Discussion Service")

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

	// Set up Gin router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())

	// Public routes - can be accessed without authentication
	router.GET("/discussions", discussionHandler.ListDiscussions)
	router.GET("/discussions/:id", discussionHandler.GetDiscussion)
	router.GET("/discussions/:id/comments", commentHandler.ListComments)

	// Protected routes - require authentication
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/discussions", discussionHandler.CreateDiscussion)
		protected.PATCH("/discussions/:id", discussionHandler.UpdateDiscussion)
		protected.DELETE("/discussions/:id", discussionHandler.DeleteDiscussion)
		protected.POST("/discussions/:id/comments", commentHandler.CreateComment)
		protected.PATCH("/comments/:id", commentHandler.UpdateComment)
		protected.DELETE("/comments/:id", commentHandler.DeleteComment)
		protected.POST("/discussions/:id/like", likeHandler.LikeDiscussion)
		protected.DELETE("/discussions/:id/like", likeHandler.UnlikeDiscussion)
		protected.POST("/comments/:id/like", likeHandler.LikeComment)
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
