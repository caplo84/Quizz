package services

import (
	"context"
	"errors"
	"testing"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTopicRepository for testing
type MockTopicRepository struct {
	mock.Mock
}

func (m *MockTopicRepository) GetAllTopics(ctx context.Context) ([]models.Topic, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Topic), args.Error(1)
}

func (m *MockTopicRepository) GetTopicByID(ctx context.Context, id uint) (*models.Topic, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Topic), args.Error(1)
}

func (m *MockTopicRepository) CreateTopic(ctx context.Context, topic *models.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func TestTopicService_GetAllTopics(t *testing.T) {
	mockRepo := new(MockTopicRepository)
	service := NewTopicService(mockRepo)
	ctx := context.Background()

	// Test successful case
	expectedTopics := []models.Topic{
		{ID: 1, Name: "Programming", Slug: "programming"},
		{ID: 2, Name: "Science", Slug: "science"},
	}

	mockRepo.On("GetAllTopics", ctx).Return(expectedTopics, nil)

	topics, err := service.GetAllTopics(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedTopics, topics)
	mockRepo.AssertExpectations(t)
}

func TestTopicService_GetAllTopics_Error(t *testing.T) {
	mockRepo := new(MockTopicRepository)
	service := NewTopicService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetAllTopics", ctx).Return([]models.Topic{}, errors.New("database error"))

	topics, err := service.GetAllTopics(ctx)

	assert.Error(t, err)
	assert.Empty(t, topics)
	mockRepo.AssertExpectations(t)
}

func TestTopicService_GetTopicByID(t *testing.T) {
	mockRepo := new(MockTopicRepository)
	service := NewTopicService(mockRepo)
	ctx := context.Background()

	expectedTopic := &models.Topic{ID: 1, Name: "Programming", Slug: "programming"}

	mockRepo.On("GetTopicByID", ctx, uint(1)).Return(expectedTopic, nil)

	topic, err := service.GetTopicByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedTopic, topic)
	mockRepo.AssertExpectations(t)
}
