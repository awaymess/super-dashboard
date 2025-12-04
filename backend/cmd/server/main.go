package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/superdashboard/backend/internal/config"
	"github.com/superdashboard/backend/internal/handler"
	"github.com/superdashboard/backend/internal/repository"
	"github.com/superdashboard/backend/internal/service"
	"github.com/superdashboard/backend/pkg/database"
	"github.com/superdashboard/backend/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	log.Info().
		Str("env", cfg.Env).
		Str("port", cfg.Port).
		Bool("useMockData", cfg.UseMockData).
		Msg("Configuration loaded")

	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.Default())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Super Dashboard API v1",
				"version": "1.0.0",
			})
		})
	}

	// Initialize database and services only if not using mock data
	if !cfg.UseMockData && cfg.DatabaseURL != "" {
		db, err := database.Connect(cfg.DatabaseURL)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}

		// Run migrations
		if err := database.AutoMigrate(db); err != nil {
			log.Fatal().Err(err).Msg("Failed to run database migrations")
		}

		// Initialize repositories
		userRepo := repository.NewUserRepository(db)

		// Initialize services
		authService := service.NewAuthService(userRepo, cfg.JWTSecret)

		// Initialize handlers
		authHandler := handler.NewAuthHandler(authService)

		// Register routes
		authHandler.RegisterAuthRoutes(v1)

		log.Info().Msg("Database-backed services initialized")
	} else {
		log.Info().Msg("Running with mock data mode - auth endpoints not available")
	}

	// Start server
	addr := ":" + cfg.Port
	log.Info().Str("addr", addr).Msg("Starting server")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
