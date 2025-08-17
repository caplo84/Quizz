package repository

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "github.com/caplo84/quizz-backend/internal/models"
)

func stringPtr(s string) *string {
    return &s
}

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }

    // Auto-migrate the schema
    err = db.AutoMigrate(&models.Topic{})
    if err != nil {
        t.Fatalf("Failed to migrate test database: %v", err)
    }

    return db
}

func TestTopicRepository_GetAllTopics(t *testing.T) {
    db := setupTestDB(t)
    repo := NewTopicRepository(db)
    ctx := context.Background()

    // Seed test data
    topics := []models.Topic{
        {Name: "Programming", Slug: "programming", Description: stringPtr("Programming topics")},
        {Name: "Science", Slug: "science", Description: stringPtr("Science topics")},
    }

    for _, topic := range topics {
        db.Create(&topic)
    }

    // Test GetAllTopics
    result, err := repo.GetAllTopics(ctx)

    assert.NoError(t, err)
    assert.Len(t, result, 2)
    assert.Equal(t, "Programming", result[0].Name)
    assert.Equal(t, "Science", result[1].Name)
}

func TestTopicRepository_GetTopicByID(t *testing.T) {
    db := setupTestDB(t)
    repo := NewTopicRepository(db)
    ctx := context.Background()

    // Create test topic
    topic := models.Topic{Name: "Programming", Slug: "programming", Description: stringPtr("Programming topics")}
    db.Create(&topic)

    // Test GetTopicByID
    result, err := repo.GetTopicByID(ctx, topic.ID)

    assert.NoError(t, err)
    assert.Equal(t, topic.Name, result.Name)
    assert.Equal(t, topic.Slug, result.Slug)
}

func TestTopicRepository_GetTopicByID_NotFound(t *testing.T) {
    db := setupTestDB(t)
    repo := NewTopicRepository(db)
    ctx := context.Background()

    // Test GetTopicByID with non-existent ID
    result, err := repo.GetTopicByID(ctx, 999)

    assert.Error(t, err)
    assert.Nil(t, result)
}

func TestTopicRepository_CreateTopic(t *testing.T) {
    db := setupTestDB(t)
    repo := NewTopicRepository(db)
    ctx := context.Background()

    topic := &models.Topic{
        Name:        "Math",
        Slug:        "math",
        Description: stringPtr("Mathematics topics"),
    }

    err := repo.CreateTopic(ctx, topic)

    assert.NoError(t, err)
    assert.NotZero(t, topic.ID)

    // Verify topic was created
    var count int64
    db.Model(&models.Topic{}).Where("slug = ?", "math").Count(&count)
    assert.Equal(t, int64(1), count)
}