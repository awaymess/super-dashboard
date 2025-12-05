package main

// @title Super Dashboard API
// @version 1.0
// @description Integrated Sports Betting & Stock Monitoring API
// @schemes http
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/awaymess/super-dashboard/backend/internal/config"
	"github.com/awaymess/super-dashboard/backend/internal/handler"
	"github.com/awaymess/super-dashboard/backend/internal/middleware"
	"github.com/awaymess/super-dashboard/backend/internal/repository"
	"github.com/awaymess/super-dashboard/backend/internal/service"
	"github.com/awaymess/super-dashboard/backend/pkg/database"
	"github.com/awaymess/super-dashboard/backend/pkg/logger"
	"github.com/awaymess/super-dashboard/backend/pkg/nlp"
	"github.com/awaymess/super-dashboard/backend/pkg/redis"
	"github.com/awaymess/super-dashboard/backend/workers"
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

	// Add security headers middleware
	var securityConfig middleware.SecurityHeadersConfig
	if cfg.Env == "production" {
		securityConfig = middleware.ProductionSecurityHeadersConfig()
	} else {
		securityConfig = middleware.DefaultSecurityHeadersConfig()
	}
	r.Use(middleware.SecurityHeadersMiddleware(securityConfig))

	// Initialize health handler
	healthHandler := handler.NewHealthHandler()
	healthHandler.RegisterHealthRoutes(r)

	// Initialize metrics handler
	metricsHandler := handler.NewMetricsHandler()
	metricsHandler.RegisterMetricsRoutes(r)
	r.Use(metricsHandler.MetricsMiddleware())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Super Dashboard API v1",
				"version": "1.0.0",
			})
		})

		// Ping endpoint for connectivity verification
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message":   "pong",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
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

		// Initialize bet handler (mock mode)
		betHandler := handler.NewBetHandler()
		betHandler.RegisterBetRoutes(v1)
		log.Info().Msg("Betting endpoints registered")

		// Initialize paper trading handler (mock mode - legacy endpoints)
		paperTradingHandler := handler.NewPaperTradingHandler()
		paperTradingHandler.RegisterPaperTradingRoutes(v1)
		log.Info().Msg("Paper trading endpoints registered")

		// Initialize NLP handler (mock mode)
		nlpProvider := nlp.NewMockProvider()
		articleRepo := repository.NewInMemoryArticleRepository()
		nlpService := service.NewNLPService(nlpProvider, articleRepo)
		nlpHandler := handler.NewNLPHandler(nlpService)
		nlpHandler.RegisterNLPRoutes(v1)
		log.Info().Msg("NLP endpoints registered")
		// Initialize paper trading with in-memory repositories (new /paper endpoints)
		portfolioRepo := repository.NewInMemoryPortfolioRepository()
		positionRepo := repository.NewInMemoryPositionRepository()
		orderRepo := repository.NewInMemoryOrderRepository()
		tradeRepo := repository.NewInMemoryTradeRepository()

		// Seed default portfolio with some positions
		if _, err := repository.SeedDefaultPortfolio(portfolioRepo, positionRepo); err != nil {
			log.Warn().Err(err).Msg("Failed to seed default portfolio")
		}

		// Initialize paper trading service with mock price provider
		paperService := service.NewPaperTradingService(portfolioRepo, positionRepo, orderRepo, tradeRepo, nil)
		paperHandler := handler.NewPaperHandler(paperService)
		paperHandler.RegisterPaperRoutes(v1)
		log.Info().Msg("Paper trading API endpoints registered (/api/v1/paper)")

		log.Info().Msg("Running with mock data mode")
	} else if cfg.DatabaseURL != "" {
		// Use database repositories
		db, err := database.Connect(cfg.DatabaseURL)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}

		// Add database health checker with timeout
		healthHandler.AddHealthChecker(func() (string, bool, string) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			
			sqlDB, err := db.DB()
			if err != nil {
				return "database", false, err.Error()
			}
			if err := sqlDB.PingContext(ctx); err != nil {
				return "database", false, err.Error()
			}
			return "database", true, "connected"
		})

		// Run migrations
		if err := database.AutoMigrate(db); err != nil {
			log.Fatal().Err(err).Msg("Failed to run database migrations")
		}

		// Initialize repositories
		userRepo := repository.NewUserRepository(db)
		sessionRepo := repository.NewSessionRepository(db)
		oauthRepo := repository.NewOAuthAccountRepository(db)
		twoFARepo := repository.NewTwoFactorAuthRepository(db)
		auditLogRepo := repository.NewAuditLogRepository(db)
		portfolioRepo := repository.NewPortfolioRepository(db)
		positionRepo := repository.NewPositionRepository(db)
		orderRepo := repository.NewOrderRepository(db)
		tradeRepo := repository.NewTradeRepository(db)

		// Initialize Redis for token storage and rate limiting
		var tokenStore service.TokenStore
		var redisClient *goredis.Client
		if cfg.RedisURL != "" {
			redisWrapper, err := redis.Connect(cfg.RedisURL)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to connect to Redis, continuing without token persistence and distributed rate limiting")
			} else {
				tokenStore = redisWrapper
				// Parse Redis URL to get underlying client for rate limiting
				opts, _ := goredis.ParseURL(cfg.RedisURL)
				if opts != nil {
					redisClient = goredis.NewClient(opts)
				}
				// Add Redis health checker with timeout
				healthHandler.AddHealthChecker(func() (string, bool, string) {
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()
					if err := redisWrapper.Ping(ctx); err != nil {
						return "redis", false, err.Error()
					}
					return "redis", true, "connected"
				})
				log.Info().Msg("Connected to Redis for token storage and rate limiting")
			}
		}

		// Initialize extended auth service with full functionality
		authService := service.NewExtendedAuthService(service.AuthServiceConfig{
			UserRepo:     userRepo,
			SessionRepo:  sessionRepo,
			OAuthRepo:    oauthRepo,
			TwoFARepo:    twoFARepo,
			AuditLogRepo: auditLogRepo,
			TokenStore:   tokenStore,
			JWTSecret:    cfg.JWTSecret,
			IssuerName:   "SuperDashboard",
		})
		paperService := service.NewPaperTradingService(portfolioRepo, positionRepo, orderRepo, tradeRepo, nil)

		// Create auth middleware
		authMiddleware := middleware.AuthMiddleware(authService)

		// Initialize handlers
		authHandler := handler.NewExtendedAuthHandler(authService)
		paperHandler := handler.NewPaperHandler(paperService)

		// Apply rate limiting to auth routes
		authRateLimiter := middleware.AuthRateLimitMiddleware(redisClient)
		apiRateLimiter := middleware.APIRateLimitMiddleware(redisClient)

		// Register auth routes with rate limiting
		authGroup := v1.Group("/auth")
		authGroup.Use(authRateLimiter)
		authHandler.RegisterExtendedAuthRoutes(v1, authMiddleware)
		
		// Register paper routes with API rate limiting
		paperGroup := v1.Group("/paper")
		paperGroup.Use(apiRateLimiter)
		paperHandler.RegisterPaperRoutes(v1)

		log.Info().Msg("Database-backed services initialized with extended auth")
	} else {
		log.Warn().Msg("No database URL configured and not in mock mode")
	}

	// Start background workers
	// Create a cancellable context for workers
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	// Start workers as goroutines
	go workers.StartOddsSync(workerCtx, log.Logger)
	go workers.StartStockSync(workerCtx, log.Logger)
	go workers.StartAlertChecker(workerCtx, log.Logger)
	log.Info().Msg("Background workers started")

	// Start server with graceful shutdown
	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("addr", addr).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// Cancel worker context to stop background workers
	workerCancel()

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
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
