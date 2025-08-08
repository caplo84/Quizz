package repository

import (
    "context"
    "github.com/caplo84/quizz-backend/internal/models"
)

type TopicRepository interface {
    GetAllTopics(ctx context.Context) ([]models.Topic, error)
    GetTopicByID(ctx context.Context, id uint) (*models.Topic, error)
    CreateTopic(ctx context.Context, topic *models.Topic) error
}

type QuizRepository interface {
    GetQuizzesByTopic(ctx context.Context, topicID uint) ([]models.Quiz, error)
    GetQuizBySlug(ctx context.Context, slug string) (*models.Quiz, error)
    GetQuizQuestions(ctx context.Context, quizID uint) ([]models.Question, error)
    CreateQuiz(ctx context.Context, quiz *models.Quiz) error
    UpdateQuiz(ctx context.Context, quiz *models.Quiz) error
    DeleteQuiz(ctx context.Context, id uint) error
}

type AttemptRepository interface {
    CreateAttempt(ctx context.Context, attempt *models.Attempt) error
    GetAttemptByID(ctx context.Context, id uint) (*models.Attempt, error)
    UpdateAttemptAnswers(ctx context.Context, attempt *models.Attempt) error
}