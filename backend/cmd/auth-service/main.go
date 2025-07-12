package main

import (
        "log"
        "os"

        "github.com/gin-gonic/gin"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/auth/handlers"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/auth/repository"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/auth/service"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/config"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/database"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/middleware"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
        "github.com/joho/godotenv"
)

func main() {
        // Load environment variables
        if err := godotenv.Load(); err != nil {
                log.Println("Warning: .env file not found")
        }

        // Initialize logger
        logger := logger.NewLogger()
        logger.Info("Starting Great Nigeria Auth Service")

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

        // Initialize repositories
        userRepo := repository.NewUserRepository(db, logger.Logger)
        twoFARepo := repository.NewTwoFARepository(db, logger.Logger)
        sessionRepo := repository.NewSessionRepository(db, logger.Logger)
        contentAccessRepo := repository.NewGormContentAccessRepository(db)

        // Initialize services
        userService := service.NewUserService(userRepo, logger)
        twoFAService := service.NewTwoFAService(twoFARepo, userService, logger)
        sessionService := service.NewSessionService(sessionRepo, userRepo, logger)
        contentAccessService := service.NewContentAccessService(contentAccessRepo, userRepo, logger.Logger)

        // Initialize handlers
        userHandler := handlers.NewUserHandler(userService, logger)
        accountHandler := handlers.NewAccountHandler(userService, logger)
        twoFAHandler := handlers.NewTwoFAHandler(twoFAService, logger)
        sessionHandler := handlers.NewSessionHandler(sessionService, logger)
        contentAccessHandler := handlers.NewContentAccessHandler(contentAccessService, logger)
        verificationHandler := handlers.NewVerificationHandler(userService)
        profileCompletionHandler := handlers.NewProfileCompletionHandler(userService)

        // Set up Gin router
        router := gin.Default()
        router.Use(gin.Recovery())
        router.Use(middleware.RequestLogger())

        // Add health check endpoint
        router.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{
                        "status": "ok",
                        "service": "auth-service",
                })
        })

        // Define API routes
        authRoutes := router.Group("/auth")
        {
                authRoutes.POST("/register", userHandler.Register)
                authRoutes.POST("/login", userHandler.Login)
                authRoutes.POST("/refresh-token", userHandler.RefreshToken)
                authRoutes.POST("/password/reset", userHandler.ResetPassword)
                authRoutes.POST("/password/reset/confirm", userHandler.ConfirmPasswordReset)
                authRoutes.POST("/logout", userHandler.Logout)
                authRoutes.GET("/oauth/:provider", userHandler.OAuthLogin)
                authRoutes.GET("/oauth/:provider/callback", userHandler.OAuthCallback)
        }
        
        // Content access routes (public endpoint)
        contentRoutes := router.Group("/content")
        {
                contentRoutes.POST("/check-access", contentAccessHandler.CheckContentAccess)
        }
                
        // Email verification routes
        emailRoutes := router.Group("/auth/email")
        {
                emailRoutes.POST("/verify/send", userHandler.SendEmailVerification)
                emailRoutes.POST("/verify/confirm", userHandler.VerifyEmail)
                emailRoutes.POST("/verify/resend", userHandler.ResendVerificationEmail)
        }

        // Basic user routes - require authentication
        userRoutes := router.Group("/users")
        userRoutes.Use(middleware.JWTAuth())
        {
                userRoutes.GET("/:id", userHandler.GetUser)
                userRoutes.PATCH("/:id", userHandler.UpdateUser)
                userRoutes.GET("/:id/profile", userHandler.GetUserProfile)
        }
        
        // Account management routes - require authentication
        accountRoutes := router.Group("/account")
        accountRoutes.Use(middleware.JWTAuth())
        {
                accountRoutes.DELETE("/delete", accountHandler.DeleteAccount)
                
                // Two-factor authentication routes
                accountRoutes.GET("/2fa/status", twoFAHandler.GetTwoFAStatus)
                accountRoutes.POST("/2fa/setup", twoFAHandler.SetupTwoFA)
                accountRoutes.POST("/2fa/verify", twoFAHandler.VerifyTwoFA)
                accountRoutes.POST("/2fa/enable", twoFAHandler.EnableTwoFA)
                accountRoutes.POST("/2fa/disable", twoFAHandler.DisableTwoFA)
                accountRoutes.POST("/2fa/backup-codes", twoFAHandler.GenerateBackupCodes)
                accountRoutes.POST("/2fa/validate-backup", twoFAHandler.ValidateBackupCode)
                
                // Session management routes
                accountRoutes.GET("/sessions", sessionHandler.GetSessions)
                accountRoutes.POST("/sessions/revoke", sessionHandler.RevokeSession)
                accountRoutes.POST("/sessions/revoke-all", sessionHandler.RevokeAllSessions)
                
                // Privacy settings routes
                accountRoutes.GET("/privacy", contentAccessHandler.GetUserPrivacySettings)
                accountRoutes.PUT("/privacy", contentAccessHandler.UpdateUserPrivacySettings)
                
                // User permissions routes
                accountRoutes.GET("/permissions", contentAccessHandler.GetUserPermissions)
                
                // Verification routes
                accountRoutes.GET("/verification/status", verificationHandler.GetVerificationStatus)
                accountRoutes.POST("/verification/request", verificationHandler.SubmitVerificationRequest)
                accountRoutes.GET("/verification/requests", verificationHandler.GetVerificationRequests)
                accountRoutes.GET("/badges", verificationHandler.GetUserBadges)
                
                // Profile completion routes
                accountRoutes.GET("/profile-completion", profileCompletionHandler.GetProfileCompletionStatus)
                accountRoutes.POST("/profile-completion/activity", profileCompletionHandler.UpdateProfileCompletionFromActivity)
                accountRoutes.GET("/profile-completion/reminder", profileCompletionHandler.CheckProfileCompletionReminder)
        }
        
        // Admin session routes - require admin access
        adminSessionRoutes := router.Group("/admin/sessions")
        adminSessionRoutes.Use(middleware.AdminAuth())
        {
                adminSessionRoutes.POST("/maintenance", sessionHandler.PerformMaintenance)
        }
        
        // Create role handlers
        roleHandlers := handlers.NewRoleHandlers(userService)
        
        // Engaged user routes - require engaged user role or higher
        engagedUserRoutes := router.Group("/engaged")
        engagedUserRoutes.Use(middleware.RoleAuth(models.RoleEngagedUser))
        {
                engagedUserRoutes.GET("/features", roleHandlers.GetEngagedUserFeatures)
        }
        
        // Active user routes - require active user role or higher
        activeUserRoutes := router.Group("/active")
        activeUserRoutes.Use(middleware.RoleAuth(models.RoleActiveUser))
        {
                activeUserRoutes.GET("/features", roleHandlers.GetActiveUserFeatures)
        }
        
        // Premium user routes - require premium user role or higher
        premiumUserRoutes := router.Group("/premium")
        premiumUserRoutes.Use(middleware.RoleAuth(models.RolePremiumUser))
        {
                premiumUserRoutes.GET("/features", roleHandlers.GetPremiumUserFeatures)
        }
        
        // Moderator routes - require moderator role or higher
        moderatorRoutes := router.Group("/moderator")
        moderatorRoutes.Use(middleware.RoleAuth(models.RoleModerator))
        {
                moderatorRoutes.GET("/tools", roleHandlers.GetModeratorTools)
                
                // Content access management (read-only and permissions)
                moderatorRoutes.GET("/content/access", contentAccessHandler.GetContentAccess)
                moderatorRoutes.GET("/content/rules", contentAccessHandler.GetContentRules)
                moderatorRoutes.POST("/content/permissions", contentAccessHandler.GrantUserPermission)
                moderatorRoutes.DELETE("/content/permissions/:id", contentAccessHandler.RevokeUserPermission)
                
                // Verification request review
                moderatorRoutes.GET("/verification/requests", verificationHandler.GetVerificationRequests)
                moderatorRoutes.POST("/verification/requests/:id/review", verificationHandler.ReviewVerificationRequest)
                moderatorRoutes.GET("/users/:id/badges", verificationHandler.GetUserBadges)
                moderatorRoutes.GET("/users/:id/profile-completion", profileCompletionHandler.GetUserProfileCompletionStatus)
        }

        // Add admin routes with admin-only access
        adminRoutes := router.Group("/admin")
        adminRoutes.Use(middleware.AdminAuth())
        {
                // User management
                adminRoutes.GET("/users", userHandler.ListUsers)
                adminRoutes.PATCH("/users/:id/role", userHandler.UpdateUserRole)
                adminRoutes.GET("/users/role/:role", userHandler.GetUsersByRole)
                
                // Content access management
                adminRoutes.POST("/content/access", contentAccessHandler.SetContentAccess)
                adminRoutes.POST("/content/rules", contentAccessHandler.CreateContentRule)
                adminRoutes.PUT("/content/rules", contentAccessHandler.UpdateContentRule)
                adminRoutes.DELETE("/content/rules/:id", contentAccessHandler.DeleteContentRule)
        }

        // Start server
        port := os.Getenv("AUTH_SERVICE_PORT")
        if port == "" {
                port = "8001" // Default port for auth service
        }

        logger.Info("Auth Service starting on port " + port)
        if err := router.Run("0.0.0.0:" + port); err != nil {
                logger.Fatal("Failed to start Auth Service: " + err.Error())
        }
}
