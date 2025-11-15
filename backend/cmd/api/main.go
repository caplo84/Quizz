package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appLogger "github.com/caplo84/quizz-backend/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/caplo84/quizz-backend/internal/cache"
	"github.com/caplo84/quizz-backend/internal/handlers"
	"github.com/caplo84/quizz-backend/internal/metrics"
	"github.com/caplo84/quizz-backend/internal/middleware"
	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/caplo84/quizz-backend/internal/repository"
	"github.com/caplo84/quizz-backend/internal/services"
	"github.com/caplo84/quizz-backend/internal/services/datasources"
	"github.com/caplo84/quizz-backend/pkg/utils"
)

// Application holds the application dependencies
type Application struct {
	Config *utils.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Cache  cache.Cache
	Router *gin.Engine
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	// Load configuration using utility
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Printf("Configuration load failed, using fallback config: %v\n", err)
		config = buildFallbackConfig()
	}

	// Initialize structured logging based on config
	if err := appLogger.InitializeLogger(config.Logging.Level, config.Logging.Format); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	appLogger.Log.WithFields(appLogger.Fields{
		"environment": config.Server.Environment,
		"version":     "1.0.0",
	}).Info("Starting Quiz Backend API")

	// Set Gin mode based on config
	gin.SetMode(config.Server.GinMode)

	// Initialize application
	app := &Application{
		Config: config,
	}

	// Initialize database with logging
	if err := app.initDatabase(); err != nil {
		appLogger.Log.WithError(err).Error("Database initialization failed - continuing startup without database")
	}

	// Initialize Redis
	if err := app.initRedis(); err != nil {
		appLogger.Log.WithError(err).Error("Redis initialization failed - continuing startup without Redis")
	}

	// Initialize cache with Redis when available, otherwise use in-memory fallback.
	if app.Redis != nil {
		app.Cache = cache.NewRedisCache(app.Redis)
		appLogger.Log.Info("Using Redis cache")
	} else {
		app.Cache = cache.NewMemoryCache()
		appLogger.Log.Warn("Using in-memory cache fallback")
	}

	// Setup routes
	app.Router = app.setupRouter()

	// Start server with graceful shutdown
	app.startServer()
}

// initDatabase initializes the database connection
func (app *Application) initDatabase() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		app.Config.Database.Host,
		app.Config.Database.Port,
		app.Config.Database.User,
		app.Config.Database.Password,
		app.Config.Database.DBName,
		app.Config.Database.SSLMode,
	)

	// Configure GORM logger
	var gormLoggerInstance gormLogger.Interface
	if app.Config.Server.Environment == "development" {
		gormLoggerInstance = gormLogger.Default.LogMode(gormLogger.Info)
	} else {
		gormLoggerInstance = gormLogger.Default.LogMode(gormLogger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLoggerInstance,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	app.DB = db
	log.Println("✅ Database connected successfully")

	if app.Config.Server.Environment == "development" {
		if err := app.autoMigrate(); err != nil {
			log.Printf("⚠️  Auto-migration warning: %v", err)
		}
	}

	return nil
}

// autoMigrate runs GORM auto-migration for development
func (app *Application) autoMigrate() error {
	// Only run in development environment
	if app.Config.Server.Environment != "development" {
		appLogger.Log.Info("Skipping auto-migration in non-development environment")
		return nil
	}

	appLogger.Log.Info("Running GORM auto-migration for development")

	// Use silent mode to reduce noise while keeping error reporting
	migrationDB := app.DB.Session(&gorm.Session{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})

	// Migrate all models - GORM handles dependencies automatically
	err := migrationDB.AutoMigrate(
		&models.Topic{},
		&models.Quiz{},
		&models.Question{},
		&models.Choice{},
		&models.Attempt{},
	)

	if err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	appLogger.Log.Info("✅ Auto-migration completed successfully")
	return nil
}

// initRedis initializes the Redis connection
func (app *Application) initRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", app.Config.Redis.Host, app.Config.Redis.Port),
		Password: app.Config.Redis.Password,
		DB:       app.Config.Redis.DB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	app.Redis = rdb
	log.Println("✅ Redis connected successfully")
	return nil
}

// Update the setupRouter method to use middleware and services correctly
func (app *Application) setupRouter() *gin.Engine {
	router := gin.New()

	// Add CORS middleware
	router.Use(middleware.CORS())

	// Add Prometheus middleware
	router.Use(metrics.PrometheusMiddleware())

	// Add other middleware
	router.Use(middleware.RequestLoggingMiddleware())
	router.Use(gin.Recovery())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(app.DB, app.Redis)

	// Simple root endpoint for platform health checks.
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Health check endpoints
	health := router.Group("/health")
	{
		health.GET("/live", healthHandler.LivenessProbe)
		health.GET("/ready", healthHandler.ReadinessProbe)
	}

	// Prometheus metrics endpoint
	router.GET("/metrics", metrics.MetricsHandler())

	// GitHub client configuration
	githubConfig := datasources.GitHubConfig{
		Token:      os.Getenv("GITHUB_TOKEN"), // Set this environment variable
		Owner:      "Ebazhanov",
		Repository: "linkedin-skill-assessments-quizzes",
		BaseURL:    "https://api.github.com",
	}

	githubClient := datasources.NewGitHubClient(githubConfig)

	// Initialize repositories and services with cache
	topicRepo := repository.NewTopicRepository(app.DB)
	topicService := services.NewTopicService(topicRepo, app.Cache)
	topicHandler := handlers.NewTopicHandler(topicService)

	quizRepo := repository.NewQuizRepository(app.DB)
	quizService := services.NewQuizService(quizRepo, app.Cache)
	quizHandler := handlers.NewQuizHandler(quizService, topicService)

	attemptRepo := repository.NewAttemptRepository(app.DB)
	attemptService := services.NewAttemptService(attemptRepo, quizService, app.Cache)
	attemptHandler := handlers.NewAttemptHandler(attemptService, quizService)

	adminService := services.NewAdminService(quizRepo, app.Cache, githubClient, topicRepo)
	aiAnswerService := services.NewAIAnswerServiceFromEnv()
	questionCorrector := services.NewQuizCorrectorService(app.DB, githubClient, aiAnswerService)

	// Create GitHub sync service
	githubSyncService := services.NewGitHubSyncService(githubClient, quizRepo, topicRepo)

	adminHandler := handlers.NewAdminHandler(adminService, githubSyncService, questionCorrector)

	// Serve static files for quiz images
	router.Static("/static", "./static")

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.HealthCheck)
		v1.GET("/topics", topicHandler.GetTopics)

		v1.GET("/topics/:topic/quizzes", quizHandler.GetQuizzes)
		v1.GET("/topics/:topic/questions/random", quizHandler.GetRandomQuestions)
		v1.GET("/questions/by-ids", quizHandler.GetQuestionsByIDs) // New endpoint
		v1.GET("/quizzes/:slug", quizHandler.GetQuizBySlug)
		v1.GET("/quizzes/:slug/questions", quizHandler.GetQuizQuestions)

		v1.POST("/quizzes/:slug/attempts", attemptHandler.CreateAttempt)
		v1.PUT("/quizzes/:slug/attempts/:id", attemptHandler.SubmitAttempt)
		v1.GET("/quizzes/:slug/attempts/:id", attemptHandler.GetAttempt)

		// Admin routes
		admin := v1.Group("/admin")
		{
			admin.GET("/quizzes/:id", adminHandler.GetQuizByID)
			admin.POST("/quizzes", adminHandler.CreateQuiz)
			admin.PUT("/quizzes/:id", adminHandler.UpdateQuiz)
			admin.DELETE("/quizzes/:id", adminHandler.DeleteQuiz)

			admin.POST("/topics", adminHandler.CreateTopic)
			admin.PUT("/topics/:id", adminHandler.UpdateTopic)
			admin.DELETE("/topics/:id", adminHandler.DeleteTopic)
		}
	}

	adminRoutes := router.Group("/api/admin")
	{
		adminRoutes.POST("/sync/github", adminHandler.SyncGitHubData)
		adminRoutes.GET("/sync/github/status", adminHandler.GetGitHubSyncStatus)
		adminRoutes.POST("/download-all-topic-images", adminHandler.DownloadAllTopicImages)
		adminRoutes.POST("/questions/correct", adminHandler.CorrectQuestions)
		adminRoutes.GET("/ai/settings", adminHandler.GetAISettings)
		adminRoutes.PUT("/ai/settings", adminHandler.UpdateAISettings)
	}

	router.GET("/health", healthHandler.HealthCheck)

	return router
}

func (app *Application) startServer() {
	port := app.resolveServerPort()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.Router,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Log.WithFields(appLogger.Fields{
			"port":        port,
			"env":         app.Config.Server.Environment,
			"port_env_set": os.Getenv("PORT") != "",
		}).Info("Server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Log.WithError(err).Fatal("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Log.Info("Server shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Log.WithError(err).Fatal("Server forced to shutdown")
	}

	// Close database connection
	if app.DB != nil {
		if sqlDB, err := app.DB.DB(); err == nil {
			sqlDB.Close()
		}
	}

	// Close Redis connection
	if app.Redis != nil {
		app.Redis.Close()
	}

	appLogger.Log.Info("Server exited")
}

func (app *Application) resolveServerPort() string {
	if envPort := os.Getenv("PORT"); envPort != "" {
		return envPort
	}

	if app.Config != nil && app.Config.Server.Port != "" {
		return app.Config.Server.Port
	}

	return "8080"
}

func buildFallbackConfig() *utils.Config {
	serverEnv := os.Getenv("APP_ENV")
	if serverEnv == "" {
		serverEnv = "development"
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		if serverEnv == "production" {
			ginMode = gin.ReleaseMode
		} else {
			ginMode = gin.DebugMode
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &utils.Config{
		Server: utils.ServerConfig{
			Port:        port,
			Environment: serverEnv,
			GinMode:     ginMode,
		},
		Database: utils.DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "quiz_user"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   getEnvOrDefault("DB_NAME", "quiz_db"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		Redis: utils.RedisConfig{
			Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     getEnvOrDefault("REDIS_PORT", "6379"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		},
		Logging: utils.LoggingConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "json"),
			Output: "stdout",
		},
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
