package main

import (
	"net/http"
	"path/filepath"

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

	// Initialize services based on configuration
	if cfg.UseMockData {
		// Use mock repositories
		log.Info().Msg("Initializing mock data repositories")

		// Find mock data directory
		mockDir := findMockDir()

		// Initialize match repository
		matchRepo, err := repository.NewMockMatchRepository(filepath.Join(mockDir, "matches.json"))
		if err != nil {
			log.Warn().Err(err).Msg("Failed to load mock match data")
		} else {
			matchHandler := handler.NewMatchHandler(matchRepo)
			matchHandler.RegisterMatchRoutes(v1)
			log.Info().Msg("Match endpoints registered with mock data")
		}

		// Initialize stock repository
		stockRepo, err := repository.NewMockStockRepository(filepath.Join(mockDir, "stocks.json"))
		if err != nil {
			log.Warn().Err(err).Msg("Failed to load mock stock data")
		} else {
			stockHandler := handler.NewStockHandler(stockRepo)
			stockHandler.RegisterStockRoutes(v1)
			log.Info().Msg("Stock endpoints registered with mock data")
		}

		log.Info().Msg("Running with mock data mode")
	} else if cfg.DatabaseURL != "" {
		// Use database repositories
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
		log.Warn().Msg("No database URL configured and not in mock mode")
	}

	// Start server
	addr := ":" + cfg.Port
	log.Info().Str("addr", addr).Msg("Starting server")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

// findMockDir finds the mock data directory.
func findMockDir() string {
	// Try relative paths from different working directories
	paths := []string{
		"mock",
		"../mock",
		"backend/mock",
	}
	for _, p := range paths {
		if _, err := filepath.Glob(filepath.Join(p, "*.json")); err == nil {
			return p
		}
	}
	return "mock"
}
