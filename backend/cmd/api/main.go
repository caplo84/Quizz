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

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/caplo84/quizz-backend/internal/handlers"
	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/caplo84/quizz-backend/internal/middleware"
	"github.com/caplo84/quizz-backend/internal/repository"
	"github.com/caplo84/quizz-backend/internal/services"
	"github.com/caplo84/quizz-backend/pkg/utils"
)

// Application holds the application dependencies
type Application struct {
	Config *utils.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Router *gin.Engine
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration using utility
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize application
	app := &Application{
		Config: config,
	}

	// Set Gin mode based on environment
	if config.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	if err := app.initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	if err := app.initRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
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
	var gormLogger logger.Interface
	if app.Config.Server.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
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
	return app.DB.AutoMigrate(
		&models.Topic{},
		&models.Quiz{},
		&models.Question{},
		&models.Choice{},
		&models.Attempt{},
	)
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

// Add this method to the Application struct
func (app *Application) setupMiddleware(router *gin.Engine) {
	// Recovery middleware recovers from any panics and writes a 500 if there was one
	router.Use(gin.Recovery())

	// Custom logger middleware
	router.Use(middleware.Logger())

	// Error handling (must be after logger)
	router.Use(middleware.ErrorHandler())

	// Security headers
	router.Use(middleware.SecurityHeaders())

	// CORS
	router.Use(middleware.CORS())

	// JSON validation
	router.Use(middleware.ValidateJSON())

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(app.Redis, 10.0, 20) // 10 req/sec, burst 20
	router.Use(rateLimiter.RateLimit())
}

// Update the setupRouter method to use middleware
func (app *Application) setupRouter() *gin.Engine {
	router := gin.New()
	app.setupMiddleware(router)

	// Initialize repositories and services
	topicRepo := repository.NewTopicRepository(app.DB)
	topicService := services.NewTopicService(topicRepo)
	topicHandler := handlers.NewTopicHandler(topicService)

	quizRepo := repository.NewQuizRepository(app.DB)
	quizService := services.NewQuizService(quizRepo)
	quizHandler := handlers.NewQuizHandler(quizService)

	attemptRepo := repository.NewAttemptRepository(app.DB)
	attemptService := services.NewAttemptService(attemptRepo)
	attemptHandler := handlers.NewAttemptHandler(attemptService, quizService)
	adminService := services.NewAdminService(quizRepo)
	// Health handler
	healthHandler := handlers.NewHealthHandler(app.DB, app.Redis)
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.GET("/topics", topicHandler.GetTopics)
		v1.GET("/quizzes", quizHandler.GetQuizzes)
		v1.GET("/quizzes/:slug", quizHandler.GetQuizBySlug)
		v1.GET("/quizzes/:slug/questions", quizHandler.GetQuizQuestions)

		// Quiz attempt routes
		v1.POST("/quizzes/:slug/attempts", attemptHandler.CreateAttempt)
		v1.PUT("/quizzes/:slug/attempts/:id", attemptHandler.SubmitAttempt)
		v1.GET("/quizzes/:slug/attempts/:id", attemptHandler.GetAttempt)

		// In setupRouter() function, update the admin section:
		admin := v1.Group("/admin")
		admin.Use(middleware.AdminAuth())
		{
			adminHandler := handlers.NewAdminHandler(adminService)
			admin.POST("/quizzes", adminHandler.CreateQuiz)
			admin.PUT("/quizzes/:id", adminHandler.UpdateQuiz)
			admin.DELETE("/quizzes/:id", adminHandler.DeleteQuiz)
		}
	}

	return router
}

// startServer starts the HTTP server with graceful shutdown
func (app *Application) startServer() {
	server := &http.Server{
		Addr:    ":" + app.Config.Server.Port,
		Handler: app.Router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", app.Config.Server.Port)
		log.Printf("🌍 Environment: %s", app.Config.Server.Environment)
		log.Printf("📋 Health check: http://localhost:%s/health", app.Config.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if sqlDB, err := app.DB.DB(); err == nil {
		sqlDB.Close()
	}

	// Close Redis connection
	if app.Redis != nil {
		app.Redis.Close()
	}

	log.Println("✅ Server exited gracefully")
}

// Handler functions (placeholders - will be implemented in handlers package)

// healthCheckHandler handles health check requests
func (app *Application) healthCheckHandler(c *gin.Context) {
	// Test database connection
	sqlDB, err := app.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database connection failed",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database ping failed",
		})
		return
	}

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := app.Redis.Ping(ctx).Err(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Redis connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"timestamp":   time.Now().UTC(),
		"environment": app.Config.Server.Environment,
		"version":     "1.0.0",
		"services": gin.H{
			"database": "connected",
			"redis":    "connected",
		},
	})
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}