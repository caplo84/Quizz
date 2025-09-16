package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTopicService for testing
type MockTopicService struct {
	mock.Mock
}

func (m *MockTopicService) GetAllTopics(ctx context.Context) ([]models.Topic, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Topic), args.Error(1)
}

func (m *MockTopicService) GetTopicByID(ctx context.Context, id uint) (*models.Topic, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Topic), args.Error(1)
}

func (m *MockTopicService) CreateTopic(ctx context.Context, topic *models.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicService) GetTopicBySlug(ctx context.Context, slug string) (*models.Topic, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(*models.Topic), args.Error(1)
}

func TestTopicHandler_GetTopics_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTopicService)
	handler := NewTopicHandler(mockService)

	expectedTopics := []models.Topic{
		{ID: 1, Name: "Programming", Slug: "programming"},
		{ID: 2, Name: "Science", Slug: "science"},
	}

	mockService.On("GetAllTopics", mock.Anything).Return(expectedTopics, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/topics", nil)
	c.Request = req

	handler.GetTopics(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestTopicHandler_GetTopics_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockTopicService)
	handler := NewTopicHandler(mockService)

	mockService.On("GetAllTopics", mock.Anything).Return([]models.Topic{}, errors.New("service error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/topics", nil)
	c.Request = req

	handler.GetTopics(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
