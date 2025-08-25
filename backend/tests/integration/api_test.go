package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/caplo84/quizz-backend/internal/handlers"
	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/caplo84/quizz-backend/internal/repository"
	"github.com/caplo84/quizz-backend/internal/services"
)

func stringPtr(s string) *string {
	return &s
}

func setupTestApp(t *testing.T) (*gin.Engine, *gorm.DB) {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate
	err = db.AutoMigrate(&models.Topic{}, &models.Quiz{}, &models.Question{}, &models.Choice{}, &models.Attempt{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Setup repositories and services
	topicRepo := repository.NewTopicRepository(db)
	topicService := services.NewTopicService(topicRepo)
	topicHandler := handlers.NewTopicHandler(topicService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/topics", topicHandler.GetTopics)
	}

	return router, db
}

func TestAPI_GetTopics_Integration(t *testing.T) {
	router, db := setupTestApp(t)

	// Seed test data
	topics := []models.Topic{
		{Name: "Programming", Slug: "programming", Description: stringPtr("Programming topics")},
		{Name: "Science", Slug: "science", Description: stringPtr("Science topics")},
	}

	for _, topic := range topics {
		db.Create(&topic)
	}

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/topics", nil)
	router.ServeHTTP(w, req)

	// Debug: Print the actual response
	t.Logf("Response body: %s", w.Body.String())

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// If using standardized response format, expect this structure:
	var response struct {
		Data []models.Topic `json:"data"`
	}

	// Or if just returning the array directly:
	// var response []models.Topic

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Data, 2) // Change to response.Data if using standardized format
	// assert.Len(t, response, 2) // Use this if returning array directly
}

func TestAPI_HealthCheck_Integration(t *testing.T) {
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"services": gin.H{
				"database": "connected",
			},
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}
