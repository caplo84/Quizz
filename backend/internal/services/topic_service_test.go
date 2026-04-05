package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTopicRepository mocks the repository.TopicRepository interface
type MockTopicRepository struct {
	mock.Mock
}

func (m *MockTopicRepository) GetAllTopics(ctx context.Context) ([]models.Topic, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Topic), args.Error(1)
}

func (m *MockTopicRepository) GetTopicByID(ctx context.Context, id uint) (*models.Topic, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Topic), args.Error(1)
}

func (m *MockTopicRepository) GetBySlug(ctx context.Context, slug string) (*models.Topic, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Topic), args.Error(1)
}

func (m *MockTopicRepository) Create(ctx context.Context, topic *models.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicRepository) CreateTopic(ctx context.Context, topic *models.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicRepository) UpdateTopic(ctx context.Context, topic *models.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicRepository) DeleteTopic(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCache for service tests
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCache) FlushAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ---- TopicService tests ----

func TestTopicService_GetAllTopics_Success(t *testing.T) {
	repo := new(MockTopicRepository)
	svc := NewTopicService(repo, &MockCache{})
	ctx := context.Background()

	expected := []models.Topic{{ID: 1, Name: "Go"}, {ID: 2, Name: "Python"}}
	repo.On("GetAllTopics", ctx).Return(expected, nil)

	result, err := svc.GetAllTopics(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestTopicService_GetAllTopics_Error(t *testing.T) {
	repo := new(MockTopicRepository)
	svc := NewTopicService(repo, &MockCache{})
	ctx := context.Background()

	repo.On("GetAllTopics", ctx).Return([]models.Topic{}, errors.New("db error"))

	result, err := svc.GetAllTopics(ctx)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestTopicService_GetTopicByID_Success(t *testing.T) {
	repo := new(MockTopicRepository)
	svc := NewTopicService(repo, &MockCache{})
	ctx := context.Background()

	topic := &models.Topic{ID: 1, Name: "Go"}
	repo.On("GetTopicByID", ctx, uint(1)).Return(topic, nil)

	result, err := svc.GetTopicByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Go", result.Name)
}

func TestTopicService_GetTopicByID_NotFound(t *testing.T) {
	repo := new(MockTopicRepository)
	svc := NewTopicService(repo, &MockCache{})
	ctx := context.Background()

	repo.On("GetTopicByID", ctx, uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.GetTopicByID(ctx, 99)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestTopicService_GetTopicBySlug_Success(t *testing.T) {
	repo := new(MockTopicRepository)
	svc := NewTopicService(repo, &MockCache{})
	ctx := context.Background()

	topic := &models.Topic{ID: 2, Slug: "python"}
	repo.On("GetBySlug", ctx, "python").Return(topic, nil)

	result, err := svc.GetTopicBySlug(ctx, "python")
	assert.NoError(t, err)
	assert.Equal(t, "python", result.Slug)
}

func TestTopicService_CreateTopic_Success(t *testing.T) {
	repo := new(MockTopicRepository)
	svc := NewTopicService(repo, &MockCache{})
	ctx := context.Background()

	topic := &models.Topic{Name: "Rust", Slug: "rust"}
	repo.On("CreateTopic", ctx, topic).Return(nil)

	err := svc.CreateTopic(ctx, topic)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTopicService_CreateTopic_Error(t *testing.T) {
	repo := new(MockTopicRepository)
	svc := NewTopicService(repo, &MockCache{})
	ctx := context.Background()

	topic := &models.Topic{Name: "Rust", Slug: "rust"}
	repo.On("CreateTopic", ctx, topic).Return(errors.New("duplicate key"))

	err := svc.CreateTopic(ctx, topic)
	assert.Error(t, err)
}
