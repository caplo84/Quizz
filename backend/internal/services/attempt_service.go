package services

import (
    "context"
    "time"
    "github.com/caplo84/quizz-backend/internal/cache"
    "github.com/caplo84/quizz-backend/internal/models"
    "github.com/caplo84/quizz-backend/internal/repository"
    "github.com/caplo84/quizz-backend/internal/metrics"
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
    start := time.Now()
    err := s.repo.CreateAttempt(ctx, attempt)
    metrics.RecordDatabaseOperation("insert", "attempts", time.Since(start))
    return err
}

func (s *attemptService) GetAttemptByID(ctx context.Context, id uint) (*models.Attempt, error) {
    start := time.Now()
    attempt, err := s.repo.GetAttemptByID(ctx, id)
    metrics.RecordDatabaseOperation("select", "attempts", time.Since(start))
    return attempt, err
}

func (s *attemptService) UpdateAttemptAnswers(ctx context.Context, attempt *models.Attempt) error {
    start := time.Now()
    err := s.repo.UpdateAttemptAnswers(ctx, attempt)
    
    if err == nil && attempt.IsCompleted {
        quiz, _ := s.quizService.GetQuizByID(ctx, attempt.QuizID)
        if quiz != nil {
            metrics.RecordQuizAttempt(quiz.Topic.Name, true, float64(attempt.Score))
        }
    }
    
    metrics.RecordDatabaseOperation("update", "attempts", time.Since(start))
    return err
}