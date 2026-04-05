package repository

import (
	"context"
	"testing"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupAttemptTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Topic{}, &models.Quiz{}, &models.Attempt{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func seedQuizForAttempt(t *testing.T, db *gorm.DB) models.Quiz {
	topic := models.Topic{Name: "Test", Slug: "test"}
	db.Create(&topic)
	quiz := models.Quiz{Title: "Test Quiz", Slug: "test-quiz", TopicID: topic.ID}
	if err := db.Create(&quiz).Error; err != nil {
		t.Fatalf("Failed to seed quiz: %v", err)
	}
	return quiz
}

func TestAttemptRepository_CreateAttempt(t *testing.T) {
	db := setupAttemptTestDB(t)
	repo := NewAttemptRepository(db)
	ctx := context.Background()

	quiz := seedQuizForAttempt(t, db)
	attempt := &models.Attempt{
		QuizID: quiz.ID,
		Status: "in_progress",
	}

	err := repo.CreateAttempt(ctx, attempt)
	assert.NoError(t, err)
	assert.NotZero(t, attempt.ID)
}

func TestAttemptRepository_GetAttemptByID_Success(t *testing.T) {
	db := setupAttemptTestDB(t)
	repo := NewAttemptRepository(db)
	ctx := context.Background()

	quiz := seedQuizForAttempt(t, db)
	attempt := &models.Attempt{QuizID: quiz.ID, Status: "in_progress"}
	db.Create(attempt)

	result, err := repo.GetAttemptByID(ctx, attempt.ID)
	assert.NoError(t, err)
	assert.Equal(t, attempt.ID, result.ID)
	assert.Equal(t, quiz.ID, result.QuizID)
}

func TestAttemptRepository_GetAttemptByID_NotFound(t *testing.T) {
	db := setupAttemptTestDB(t)
	repo := NewAttemptRepository(db)
	ctx := context.Background()

	result, err := repo.GetAttemptByID(ctx, 9999)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAttemptRepository_UpdateAttemptAnswers(t *testing.T) {
	db := setupAttemptTestDB(t)
	repo := NewAttemptRepository(db)
	ctx := context.Background()

	quiz := seedQuizForAttempt(t, db)
	attempt := &models.Attempt{QuizID: quiz.ID, Status: "in_progress"}
	db.Create(attempt)

	// Update the attempt
	attempt.Status = "completed"
	attempt.IsCompleted = true
	attempt.Score = 7
	attempt.Answers = models.UserAnswers{"1": "a", "2": "b"}

	err := repo.UpdateAttemptAnswers(ctx, attempt)
	assert.NoError(t, err)

	var updated models.Attempt
	db.First(&updated, attempt.ID)
	assert.Equal(t, "completed", updated.Status)
	assert.True(t, updated.IsCompleted)
	assert.Equal(t, 7, updated.Score)
}

func TestAttemptRepository_CreateAttempt_WithUserInfo(t *testing.T) {
	db := setupAttemptTestDB(t)
	repo := NewAttemptRepository(db)
	ctx := context.Background()

	quiz := seedQuizForAttempt(t, db)
	name := "Alice"
	identifier := "user-001"
	attempt := &models.Attempt{
		QuizID:         quiz.ID,
		UserName:       &name,
		UserIdentifier: &identifier,
		Status:         "in_progress",
	}

	err := repo.CreateAttempt(ctx, attempt)
	assert.NoError(t, err)
	assert.NotZero(t, attempt.ID)

	// Verify user info was stored
	result, err := repo.GetAttemptByID(ctx, attempt.ID)
	assert.NoError(t, err)
	assert.Equal(t, &name, result.UserName)
	assert.Equal(t, &identifier, result.UserIdentifier)
}
