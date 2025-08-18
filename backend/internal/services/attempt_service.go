package services

import (
    "context"
    "github.com/caplo84/quizz-backend/internal/cache"
    "github.com/caplo84/quizz-backend/internal/models"
    "github.com/caplo84/quizz-backend/internal/repository"
)

type AttemptService interface {
    CreateAttempt(ctx context.Context, attempt *models.Attempt) error
    GetAttemptByID(ctx context.Context, id uint) (*models.Attempt, error)
    UpdateAttemptAnswers(ctx context.Context, attempt *models.Attempt) error
}

type attemptService struct {
    repo       repository.AttemptRepository
    quizService QuizService
    cache      cache.Cache
}

func NewAttemptService(repo repository.AttemptRepository, quizService QuizService, cache cache.Cache) AttemptService {
    return &attemptService{
        repo:       repo,
        quizService: quizService,
        cache:      cache,
    }
}

func (s *attemptService) CreateAttempt(ctx context.Context, attempt *models.Attempt) error {
    // Add business logic, validation, etc.
    return s.repo.CreateAttempt(ctx, attempt)
}

func (s *attemptService) GetAttemptByID(ctx context.Context, id uint) (*models.Attempt, error) {
    return s.repo.GetAttemptByID(ctx, id)
}

func (s *attemptService) UpdateAttemptAnswers(ctx context.Context, attempt *models.Attempt) error {
    // Add answer validation, scoring, etc.
    return s.repo.UpdateAttemptAnswers(ctx, attempt)
}