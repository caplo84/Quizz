package services

import (
    "context"
    "github.com/caplo84/quizz-backend/internal/models"
    "github.com/caplo84/quizz-backend/internal/repository"
)

type QuizService interface {
    GetQuizzesByTopic(ctx context.Context, topicID uint) ([]models.Quiz, error)
    GetQuizBySlug(ctx context.Context, slug string) (*models.Quiz, error)
    GetQuizQuestions(ctx context.Context, quizID uint) ([]models.Question, error)
    CreateQuiz(ctx context.Context, quiz *models.Quiz) error
    UpdateQuiz(ctx context.Context, quiz *models.Quiz) error
    DeleteQuiz(ctx context.Context, id uint) error
}

type quizService struct {
    repo repository.QuizRepository
}

func NewQuizService(repo repository.QuizRepository) QuizService {
    return &quizService{repo: repo}
}

func (s *quizService) GetQuizzesByTopic(ctx context.Context, topicID uint) ([]models.Quiz, error) {
    return s.repo.GetQuizzesByTopic(ctx, topicID)
}

func (s *quizService) GetQuizBySlug(ctx context.Context, slug string) (*models.Quiz, error) {
    return s.repo.GetQuizBySlug(ctx, slug)
}

func (s *quizService) GetQuizQuestions(ctx context.Context, quizID uint) ([]models.Question, error) {
    return s.repo.GetQuizQuestions(ctx, quizID)
}

func (s *quizService) CreateQuiz(ctx context.Context, quiz *models.Quiz) error {
    // Add validation or business logic here if needed
    return s.repo.CreateQuiz(ctx, quiz)
}

func (s *quizService) UpdateQuiz(ctx context.Context, quiz *models.Quiz) error {
    return s.repo.UpdateQuiz(ctx, quiz)
}

func (s *quizService) DeleteQuiz(ctx context.Context, id uint) error {
    return s.repo.DeleteQuiz(ctx, id)
}