package services

import (
	"context"
	"errors"
	"testing"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuizRepository mocks the repository.QuizRepository interface
type MockQuizRepository struct {
	mock.Mock
}

func (m *MockQuizRepository) GetQuizzesByTopic(ctx context.Context, topicID uint) ([]models.Quiz, error) {
	args := m.Called(ctx, topicID)
	return args.Get(0).([]models.Quiz), args.Error(1)
}

func (m *MockQuizRepository) GetQuizBySlug(ctx context.Context, slug string) (*models.Quiz, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quiz), args.Error(1)
}

func (m *MockQuizRepository) GetQuizByExternalID(ctx context.Context, externalID string) (*models.Quiz, error) {
	args := m.Called(ctx, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quiz), args.Error(1)
}

func (m *MockQuizRepository) GetQuizQuestions(ctx context.Context, quizID uint) ([]models.Question, error) {
	args := m.Called(ctx, quizID)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockQuizRepository) GetAllQuestions(ctx context.Context) ([]models.Question, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockQuizRepository) GetRandomQuestions(ctx context.Context, topicID uint, limit int, excludeIDs []uint) ([]models.Question, error) {
	args := m.Called(ctx, topicID, limit, excludeIDs)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockQuizRepository) GetQuestionsByIDs(ctx context.Context, ids []uint) ([]models.Question, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockQuizRepository) Create(ctx context.Context, quiz *models.Quiz) error {
	args := m.Called(ctx, quiz)
	return args.Error(0)
}

func (m *MockQuizRepository) CreateQuiz(ctx context.Context, quiz *models.Quiz) error {
	args := m.Called(ctx, quiz)
	return args.Error(0)
}

func (m *MockQuizRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}

func (m *MockQuizRepository) CreateChoice(ctx context.Context, choice *models.Choice) error {
	args := m.Called(ctx, choice)
	return args.Error(0)
}

func (m *MockQuizRepository) UpdateQuiz(ctx context.Context, quiz *models.Quiz) error {
	args := m.Called(ctx, quiz)
	return args.Error(0)
}

func (m *MockQuizRepository) DeleteQuiz(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuizRepository) GetQuizByID(ctx context.Context, id uint) (*models.Quiz, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quiz), args.Error(1)
}

// ---- QuizService tests ----

func TestQuizService_GetQuizzesByTopic_Success(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	expected := []models.Quiz{{ID: 1, Title: "Go Basics"}, {ID: 2, Title: "Go Advanced"}}
	repo.On("GetQuizzesByTopic", ctx, uint(3)).Return(expected, nil)

	result, err := svc.GetQuizzesByTopic(ctx, 3)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestQuizService_GetQuizzesByTopic_Error(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	repo.On("GetQuizzesByTopic", ctx, uint(3)).Return([]models.Quiz{}, errors.New("db error"))

	result, err := svc.GetQuizzesByTopic(ctx, 3)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestQuizService_GetQuizBySlug_CacheMiss_DBHit(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	// Cache miss
	cache.On("Get", mock.Anything, "quiz:go-basics").Return(nil, errors.New("cache miss"))
	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}
	repo.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)

	result, err := svc.GetQuizBySlug(ctx, "go-basics")
	assert.NoError(t, err)
	assert.Equal(t, "go-basics", result.Slug)
}

func TestQuizService_GetQuizBySlug_DBError(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	cache.On("Get", mock.Anything, "quiz:not-found").Return(nil, errors.New("cache miss"))
	repo.On("GetQuizBySlug", mock.Anything, "not-found").Return(nil, errors.New("not found"))

	result, err := svc.GetQuizBySlug(ctx, "not-found")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestQuizService_GetQuizByID_Success(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	quiz := &models.Quiz{ID: 5, Title: "Python Quiz"}
	repo.On("GetQuizByID", ctx, uint(5)).Return(quiz, nil)

	result, err := svc.GetQuizByID(ctx, 5)
	assert.NoError(t, err)
	assert.Equal(t, uint(5), result.ID)
}

func TestQuizService_GetQuizByID_NotFound(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	repo.On("GetQuizByID", ctx, uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.GetQuizByID(ctx, 99)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestQuizService_GetQuizQuestions_Success(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	questions := []models.Question{{ID: 1}, {ID: 2}}
	repo.On("GetQuizQuestions", ctx, uint(1)).Return(questions, nil)

	result, err := svc.GetQuizQuestions(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestQuizService_GetRandomQuestions_Success(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	questions := []models.Question{{ID: 3}, {ID: 4}, {ID: 5}}
	repo.On("GetRandomQuestions", ctx, uint(2), 3, []uint{1, 2}).Return(questions, nil)

	result, err := svc.GetRandomQuestions(ctx, 2, 3, []uint{1, 2})
	assert.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestQuizService_GetQuestionsByIDs_Success(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	questions := []models.Question{{ID: 1}, {ID: 3}}
	repo.On("GetQuestionsByIDs", ctx, []uint{1, 3}).Return(questions, nil)

	result, err := svc.GetQuestionsByIDs(ctx, []uint{1, 3})
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestQuizService_CreateQuiz_Success(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	quiz := &models.Quiz{Title: "New Quiz", Slug: "new-quiz"}
	repo.On("CreateQuiz", ctx, quiz).Return(nil)

	err := svc.CreateQuiz(ctx, quiz)
	assert.NoError(t, err)
}

func TestQuizService_CreateQuiz_Error(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	quiz := &models.Quiz{Title: "New Quiz", Slug: "new-quiz"}
	repo.On("CreateQuiz", ctx, quiz).Return(errors.New("duplicate key"))

	err := svc.CreateQuiz(ctx, quiz)
	assert.Error(t, err)
}

func TestQuizService_UpdateQuiz_InvalidateCache(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}
	repo.On("UpdateQuiz", ctx, quiz).Return(nil)
	cache.On("Delete", ctx, "quiz:go-basics").Return(nil)

	err := svc.UpdateQuiz(ctx, quiz)
	assert.NoError(t, err)
	cache.AssertCalled(t, "Delete", ctx, "quiz:go-basics")
}

func TestQuizService_UpdateQuiz_Error(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}
	repo.On("UpdateQuiz", ctx, quiz).Return(errors.New("db error"))

	err := svc.UpdateQuiz(ctx, quiz)
	assert.Error(t, err)
}

func TestQuizService_DeleteQuiz_Success(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}
	repo.On("GetQuizByID", mock.Anything, uint(1)).Return(quiz, nil)
	repo.On("DeleteQuiz", ctx, uint(1)).Return(nil)
	cache.On("Delete", mock.Anything, "quiz:go-basics").Return(nil)

	err := svc.DeleteQuiz(ctx, 1)
	assert.NoError(t, err)
}

func TestQuizService_DeleteQuiz_Error(t *testing.T) {
	repo := new(MockQuizRepository)
	cache := new(MockCache)
	svc := NewQuizService(repo, cache)
	ctx := context.Background()

	repo.On("GetQuizByID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
	repo.On("DeleteQuiz", ctx, uint(1)).Return(errors.New("db error"))

	err := svc.DeleteQuiz(ctx, 1)
	assert.Error(t, err)
}
