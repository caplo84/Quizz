package services

import (
    "context"
    "github.com/caplo84/quizz-backend/internal/cache"
    "github.com/caplo84/quizz-backend/internal/models"
    "github.com/caplo84/quizz-backend/internal/repository"
)

type AdminService interface {
    CreateQuiz(ctx context.Context, quiz *models.Quiz) error
    UpdateQuiz(ctx context.Context, quiz *models.Quiz) error
    DeleteQuiz(ctx context.Context, id uint) error
    // Add more admin operations as needed
}

type adminService struct {
    quizRepo repository.QuizRepository
    cache    cache.Cache
}

func NewAdminService(quizRepo repository.QuizRepository, cache cache.Cache) AdminService {
    return &adminService{
        quizRepo: quizRepo,
        cache:    cache,
    }
}

func (s *adminService) CreateQuiz(ctx context.Context, quiz *models.Quiz) error {
    // Add admin-specific logic, cache invalidation, etc.
    return s.quizRepo.CreateQuiz(ctx, quiz)
}

func (s *adminService) UpdateQuiz(ctx context.Context, quiz *models.Quiz) error {
    return s.quizRepo.UpdateQuiz(ctx, quiz)
}

func (s *adminService) DeleteQuiz(ctx context.Context, id uint) error {
    return s.quizRepo.DeleteQuiz(ctx, id)
}