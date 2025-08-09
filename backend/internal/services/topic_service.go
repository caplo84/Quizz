package services

import (
    "context"
    "github.com/caplo84/quizz-backend/internal/models"
    "github.com/caplo84/quizz-backend/internal/repository"
)

type TopicService interface {
    GetAllTopics(ctx context.Context) ([]models.Topic, error)
    GetTopicByID(ctx context.Context, id uint) (*models.Topic, error)
    CreateTopic(ctx context.Context, topic *models.Topic) error
}

type topicService struct {
    repo repository.TopicRepository
}

func NewTopicService(repo repository.TopicRepository) TopicService {
    return &topicService{repo: repo}
}

func (s *topicService) GetAllTopics(ctx context.Context) ([]models.Topic, error) {
    return s.repo.GetAllTopics(ctx)
}

func (s *topicService) GetTopicByID(ctx context.Context, id uint) (*models.Topic, error) {
    return s.repo.GetTopicByID(ctx, id)
}

func (s *topicService) CreateTopic(ctx context.Context, topic *models.Topic) error {
    // Add validation or business logic here if needed
    return s.repo.CreateTopic(ctx, topic)
}