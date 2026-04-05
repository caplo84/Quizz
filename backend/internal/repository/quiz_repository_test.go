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

func setupQuizTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Topic{}, &models.Quiz{}, &models.Question{}, &models.Choice{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func seedTopic(t *testing.T, db *gorm.DB) models.Topic {
	topic := models.Topic{Name: "Go", Slug: "go"}
	if err := db.Create(&topic).Error; err != nil {
		t.Fatalf("Failed to seed topic: %v", err)
	}
	return topic
}

func TestQuizRepository_CreateQuiz(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic := seedTopic(t, db)
	desc := "Basics of Go"
	quiz := &models.Quiz{
		Title:           "Go Basics",
		Slug:            "go-basics",
		Description:     &desc,
		TopicID:         topic.ID,
		DifficultyLevel: "easy",
	}

	err := repo.CreateQuiz(ctx, quiz)
	assert.NoError(t, err)
	assert.NotZero(t, quiz.ID)
}

func TestQuizRepository_GetQuizBySlug_Success(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic := seedTopic(t, db)
	quiz := &models.Quiz{Title: "Go Basics", Slug: "go-basics", TopicID: topic.ID}
	db.Create(quiz)

	result, err := repo.GetQuizBySlug(ctx, "go-basics")
	assert.NoError(t, err)
	assert.Equal(t, "Go Basics", result.Title)
}

func TestQuizRepository_GetQuizBySlug_NotFound(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	result, err := repo.GetQuizBySlug(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestQuizRepository_GetQuizByID_Success(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic := seedTopic(t, db)
	quiz := &models.Quiz{Title: "Python", Slug: "python", TopicID: topic.ID}
	db.Create(quiz)

	result, err := repo.GetQuizByID(ctx, quiz.ID)
	assert.NoError(t, err)
	assert.Equal(t, quiz.ID, result.ID)
	assert.Equal(t, "Python", result.Title)
}

func TestQuizRepository_GetQuizByID_NotFound(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	result, err := repo.GetQuizByID(ctx, 9999)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestQuizRepository_GetQuizzesByTopic(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic1 := seedTopic(t, db)
	topic2 := models.Topic{Name: "Python", Slug: "python"}
	db.Create(&topic2)

	db.Create(&models.Quiz{Title: "Go 1", Slug: "go-1", TopicID: topic1.ID})
	db.Create(&models.Quiz{Title: "Go 2", Slug: "go-2", TopicID: topic1.ID})
	db.Create(&models.Quiz{Title: "Py 1", Slug: "py-1", TopicID: topic2.ID})

	result, err := repo.GetQuizzesByTopic(ctx, topic1.ID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestQuizRepository_GetQuizzesByTopic_Empty(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	result, err := repo.GetQuizzesByTopic(ctx, 999)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestQuizRepository_UpdateQuiz(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic := seedTopic(t, db)
	quiz := &models.Quiz{Title: "Old Title", Slug: "old-slug", TopicID: topic.ID}
	db.Create(quiz)

	quiz.Title = "New Title"
	err := repo.UpdateQuiz(ctx, quiz)
	assert.NoError(t, err)

	var updated models.Quiz
	db.First(&updated, quiz.ID)
	assert.Equal(t, "New Title", updated.Title)
}

func TestQuizRepository_DeleteQuiz(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic := seedTopic(t, db)
	quiz := &models.Quiz{Title: "To Delete", Slug: "to-delete", TopicID: topic.ID}
	db.Create(quiz)

	err := repo.DeleteQuiz(ctx, quiz.ID)
	assert.NoError(t, err)

	var count int64
	db.Model(&models.Quiz{}).Where("id = ?", quiz.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestQuizRepository_GetQuizQuestions(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic := seedTopic(t, db)
	quiz := &models.Quiz{Title: "Quiz", Slug: "quiz", TopicID: topic.ID}
	db.Create(quiz)

	q1 := &models.Question{QuizID: quiz.ID, QuestionText: "Q1?", QuestionType: "multiple_choice", Points: 1, OrderIndex: 1}
	q2 := &models.Question{QuizID: quiz.ID, QuestionText: "Q2?", QuestionType: "multiple_choice", Points: 1, OrderIndex: 2}
	db.Create(q1)
	db.Create(q2)

	questions, err := repo.GetQuizQuestions(ctx, quiz.ID)
	assert.NoError(t, err)
	assert.Len(t, questions, 2)
	// Verify ordering
	assert.Equal(t, 1, questions[0].OrderIndex)
	assert.Equal(t, 2, questions[1].OrderIndex)
}

func TestQuizRepository_GetQuestionsByIDs(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic := seedTopic(t, db)
	quiz := &models.Quiz{Title: "Quiz", Slug: "quiz2", TopicID: topic.ID}
	db.Create(quiz)

	q1 := &models.Question{QuizID: quiz.ID, QuestionText: "Q1?", QuestionType: "multiple_choice", Points: 1, OrderIndex: 1, IsActive: true}
	q2 := &models.Question{QuizID: quiz.ID, QuestionText: "Q2?", QuestionType: "multiple_choice", Points: 1, OrderIndex: 2, IsActive: true}
	q3 := &models.Question{QuizID: quiz.ID, QuestionText: "Q3?", QuestionType: "multiple_choice", Points: 1, OrderIndex: 3, IsActive: true}
	db.Create(q1)
	db.Create(q2)
	db.Create(q3)

	questions, err := repo.GetQuestionsByIDs(ctx, []uint{q1.ID, q3.ID})
	assert.NoError(t, err)
	assert.Len(t, questions, 2)
}

func TestQuizRepository_GetQuizByExternalID(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	topic := seedTopic(t, db)
	extID := "ext-123"
	quiz := &models.Quiz{Title: "External Quiz", Slug: "ext-quiz", TopicID: topic.ID, ExternalID: &extID}
	db.Create(quiz)

	result, err := repo.GetQuizByExternalID(ctx, "ext-123")
	assert.NoError(t, err)
	assert.Equal(t, "External Quiz", result.Title)
}

func TestQuizRepository_GetQuizByExternalID_NotFound(t *testing.T) {
	db := setupQuizTestDB(t)
	repo := NewQuizRepository(db)
	ctx := context.Background()

	result, err := repo.GetQuizByExternalID(ctx, "nonexistent-id")
	assert.Error(t, err)
	assert.Nil(t, result)
}
