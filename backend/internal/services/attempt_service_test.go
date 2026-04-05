package services

import (
	"context"
	"errors"
	"testing"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAttemptRepository mocks the repository.AttemptRepository interface
type MockAttemptRepository struct {
	mock.Mock
}

func (m *MockAttemptRepository) CreateAttempt(ctx context.Context, attempt *models.Attempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockAttemptRepository) GetAttemptByID(ctx context.Context, id uint) (*models.Attempt, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attempt), args.Error(1)
}

func (m *MockAttemptRepository) UpdateAttemptAnswers(ctx context.Context, attempt *models.Attempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

// ---- AttemptService tests ----

func TestAttemptService_CreateAttempt_Success(t *testing.T) {
	repo := new(MockAttemptRepository)
	quizRepo := new(MockQuizRepository)
	quizSvc := NewQuizService(quizRepo, &MockCache{})
	svc := NewAttemptService(repo, quizSvc, &MockCache{})
	ctx := context.Background()

	attempt := &models.Attempt{QuizID: 1}
	repo.On("CreateAttempt", ctx, attempt).Return(nil)

	err := svc.CreateAttempt(ctx, attempt)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAttemptService_CreateAttempt_Error(t *testing.T) {
	repo := new(MockAttemptRepository)
	quizRepo := new(MockQuizRepository)
	quizSvc := NewQuizService(quizRepo, &MockCache{})
	svc := NewAttemptService(repo, quizSvc, &MockCache{})
	ctx := context.Background()

	attempt := &models.Attempt{QuizID: 1}
	repo.On("CreateAttempt", ctx, attempt).Return(errors.New("insert failed"))

	err := svc.CreateAttempt(ctx, attempt)
	assert.Error(t, err)
}

func TestAttemptService_GetAttemptByID_Success(t *testing.T) {
	repo := new(MockAttemptRepository)
	quizRepo := new(MockQuizRepository)
	quizSvc := NewQuizService(quizRepo, &MockCache{})
	svc := NewAttemptService(repo, quizSvc, &MockCache{})
	ctx := context.Background()

	expected := &models.Attempt{ID: 5, QuizID: 2}
	repo.On("GetAttemptByID", ctx, uint(5)).Return(expected, nil)

	result, err := svc.GetAttemptByID(ctx, 5)
	assert.NoError(t, err)
	assert.Equal(t, uint(5), result.ID)
}

func TestAttemptService_GetAttemptByID_NotFound(t *testing.T) {
	repo := new(MockAttemptRepository)
	quizRepo := new(MockQuizRepository)
	quizSvc := NewQuizService(quizRepo, &MockCache{})
	svc := NewAttemptService(repo, quizSvc, &MockCache{})
	ctx := context.Background()

	repo.On("GetAttemptByID", ctx, uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.GetAttemptByID(ctx, 99)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAttemptService_UpdateAttemptAnswers_Success(t *testing.T) {
	repo := new(MockAttemptRepository)
	quizRepo := new(MockQuizRepository)
	quizSvc := NewQuizService(quizRepo, &MockCache{})
	svc := NewAttemptService(repo, quizSvc, &MockCache{})
	ctx := context.Background()

	attempt := &models.Attempt{ID: 1, QuizID: 2, IsCompleted: false}
	repo.On("UpdateAttemptAnswers", ctx, attempt).Return(nil)

	err := svc.UpdateAttemptAnswers(ctx, attempt)
	assert.NoError(t, err)
}

func TestAttemptService_UpdateAttemptAnswers_Error(t *testing.T) {
	repo := new(MockAttemptRepository)
	quizRepo := new(MockQuizRepository)
	quizSvc := NewQuizService(quizRepo, &MockCache{})
	svc := NewAttemptService(repo, quizSvc, &MockCache{})
	ctx := context.Background()

	attempt := &models.Attempt{ID: 1, QuizID: 2}
	repo.On("UpdateAttemptAnswers", ctx, attempt).Return(errors.New("update failed"))

	err := svc.UpdateAttemptAnswers(ctx, attempt)
	assert.Error(t, err)
}

func TestAttemptService_UpdateAttemptAnswers_Completed(t *testing.T) {
	repo := new(MockAttemptRepository)
	quizRepo := new(MockQuizRepository)
	cache := &MockCache{}
	quizSvc := NewQuizService(quizRepo, cache)
	svc := NewAttemptService(repo, quizSvc, cache)
	ctx := context.Background()

	quiz := &models.Quiz{ID: 2, Topic: models.Topic{Name: "Go"}}
	attempt := &models.Attempt{ID: 1, QuizID: 2, IsCompleted: true, Score: 8}

	repo.On("UpdateAttemptAnswers", ctx, attempt).Return(nil)
	quizRepo.On("GetQuizByID", mock.Anything, uint(2)).Return(quiz, nil)

	err := svc.UpdateAttemptAnswers(ctx, attempt)
	assert.NoError(t, err)
}
