package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuizService implements services.QuizService for testing
type MockQuizService struct {
	mock.Mock
}

func (m *MockQuizService) GetQuizzesByTopic(ctx context.Context, topicID uint) ([]models.Quiz, error) {
	args := m.Called(ctx, topicID)
	return args.Get(0).([]models.Quiz), args.Error(1)
}

func (m *MockQuizService) GetQuizBySlug(ctx context.Context, slug string) (*models.Quiz, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quiz), args.Error(1)
}

func (m *MockQuizService) GetQuizByID(ctx context.Context, id uint) (*models.Quiz, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Quiz), args.Error(1)
}

func (m *MockQuizService) GetQuizQuestions(ctx context.Context, quizID uint) ([]models.Question, error) {
	args := m.Called(ctx, quizID)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockQuizService) GetRandomQuestions(ctx context.Context, topicID uint, limit int, excludeIDs []uint) ([]models.Question, error) {
	args := m.Called(ctx, topicID, limit, excludeIDs)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockQuizService) GetQuestionsByIDs(ctx context.Context, ids []uint) ([]models.Question, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.Question), args.Error(1)
}

func (m *MockQuizService) CreateQuiz(ctx context.Context, quiz *models.Quiz) error {
	args := m.Called(ctx, quiz)
	return args.Error(0)
}

func (m *MockQuizService) UpdateQuiz(ctx context.Context, quiz *models.Quiz) error {
	args := m.Called(ctx, quiz)
	return args.Error(0)
}

func (m *MockQuizService) DeleteQuiz(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func newQuizTestRouter(quizSvc *MockQuizService, topicSvc *MockTopicService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewQuizHandler(quizSvc, topicSvc)
	r := gin.New()
	r.GET("/quizzes", h.GetQuizzes)
	r.GET("/quizzes/:slug", h.GetQuizBySlug)
	r.GET("/quizzes/:slug/questions", h.GetQuizQuestions)
	r.GET("/topics/:topic/quizzes", h.GetQuizzes)
	r.GET("/topics/:topic/questions/random", h.GetRandomQuestions)
	r.GET("/questions/by-ids", h.GetQuestionsByIDs)
	return r
}

// ---- GetQuizzes ----

func TestQuizHandler_GetQuizzes_NoTopic(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	quizzes := []models.Quiz{{ID: 1, Title: "Go Basics", Slug: "go-basics"}}
	quizSvc.On("GetQuizzesByTopic", mock.Anything, uint(0)).Return(quizzes, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	quizSvc.AssertExpectations(t)
}

func TestQuizHandler_GetQuizzes_WithTopicSlug(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	topic := &models.Topic{ID: 3, Name: "Go", Slug: "go"}
	topicSvc.On("GetTopicBySlug", mock.Anything, "go").Return(topic, nil)
	quizzes := []models.Quiz{{ID: 1, Title: "Go Basics", Slug: "go-basics"}}
	quizSvc.On("GetQuizzesByTopic", mock.Anything, uint(3)).Return(quizzes, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics/go/quizzes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	quizSvc.AssertExpectations(t)
	topicSvc.AssertExpectations(t)
}

func TestQuizHandler_GetQuizzes_TopicNotFound(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	topicSvc.On("GetTopicBySlug", mock.Anything, "unknown").Return((*models.Topic)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics/unknown/quizzes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	topicSvc.AssertExpectations(t)
}

func TestQuizHandler_GetQuizzes_ServiceError(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	quizSvc.On("GetQuizzesByTopic", mock.Anything, uint(0)).Return([]models.Quiz{}, errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ---- GetQuizBySlug ----

func TestQuizHandler_GetQuizBySlug_Success(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	quiz := &models.Quiz{ID: 1, Title: "Go Basics", Slug: "go-basics"}
	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/go-basics", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Go Basics", data["title"])
}

func TestQuizHandler_GetQuizBySlug_NotFound(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	quizSvc.On("GetQuizBySlug", mock.Anything, "missing").Return((*models.Quiz)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/missing", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---- GetQuizQuestions ----

func TestQuizHandler_GetQuizQuestions_Success(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}
	questions := []models.Question{
		{ID: 1, QuestionText: "Q1", Choices: []models.Choice{{ID: 1, IsCorrect: true}, {ID: 2, IsCorrect: false}}},
	}

	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)
	quizSvc.On("GetQuizQuestions", mock.Anything, uint(1)).Return(questions, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/go-basics/questions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestQuizHandler_GetQuizQuestions_HidesCorrectAnswers(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}
	questions := []models.Question{
		{ID: 1, QuestionText: "Q1", Choices: []models.Choice{
			{ID: 1, ChoiceText: "Correct", IsCorrect: true},
			{ID: 2, ChoiceText: "Wrong", IsCorrect: false},
		}},
	}

	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)
	quizSvc.On("GetQuizQuestions", mock.Anything, uint(1)).Return(questions, nil)

	w := httptest.NewRecorder()
	// Without include_answers=true, IsCorrect must be hidden
	req, _ := http.NewRequest("GET", "/quizzes/go-basics/questions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	q := data[0].(map[string]interface{})
	choices := q["choices"].([]interface{})
	for _, ch := range choices {
		choice := ch.(map[string]interface{})
		// is_correct must be false (hidden)
		assert.False(t, choice["is_correct"].(bool))
	}
}

func TestQuizHandler_GetQuizQuestions_IncludeAnswers(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}
	questions := []models.Question{
		{ID: 1, QuestionText: "Q1", Choices: []models.Choice{
			{ID: 1, ChoiceText: "Correct", IsCorrect: true},
		}},
	}

	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)
	quizSvc.On("GetQuizQuestions", mock.Anything, uint(1)).Return(questions, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/go-basics/questions?include_answers=true", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	q := data[0].(map[string]interface{})
	choices := q["choices"].([]interface{})
	choice := choices[0].(map[string]interface{})
	assert.True(t, choice["is_correct"].(bool))
}

func TestQuizHandler_GetQuizQuestions_QuizNotFound(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	quizSvc.On("GetQuizBySlug", mock.Anything, "missing").Return((*models.Quiz)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/missing/questions", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---- GetRandomQuestions ----

func TestQuizHandler_GetRandomQuestions_Success(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	topic := &models.Topic{ID: 2, Name: "Go", Slug: "go"}
	questions := []models.Question{{ID: 1, QuestionText: "Q1"}, {ID: 2, QuestionText: "Q2"}}

	topicSvc.On("GetTopicBySlug", mock.Anything, "go").Return(topic, nil)
	quizSvc.On("GetRandomQuestions", mock.Anything, uint(2), 10, mock.Anything).Return(questions, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics/go/questions/random", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	meta := resp["meta"].(map[string]interface{})
	assert.Equal(t, float64(2), meta["total_questions"])
	assert.Equal(t, float64(10), meta["limit"])
}

func TestQuizHandler_GetRandomQuestions_CustomLimit(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	topic := &models.Topic{ID: 2, Slug: "go"}
	topicSvc.On("GetTopicBySlug", mock.Anything, "go").Return(topic, nil)
	quizSvc.On("GetRandomQuestions", mock.Anything, uint(2), 5, mock.Anything).Return([]models.Question{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics/go/questions/random?limit=5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestQuizHandler_GetRandomQuestions_TopicNotFound(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	topicSvc.On("GetTopicBySlug", mock.Anything, "nope").Return((*models.Topic)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/topics/nope/questions/random", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---- GetQuestionsByIDs ----

func TestQuizHandler_GetQuestionsByIDs_Success(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	questions := []models.Question{{ID: 1}, {ID: 2}}
	quizSvc.On("GetQuestionsByIDs", mock.Anything, mock.MatchedBy(func(ids []uint) bool {
		return len(ids) == 2
	})).Return(questions, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions/by-ids?ids=1,2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestQuizHandler_GetQuestionsByIDs_MissingParam(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions/by-ids", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQuizHandler_GetQuestionsByIDs_InvalidIDs(t *testing.T) {
	quizSvc := new(MockQuizService)
	topicSvc := new(MockTopicService)
	router := newQuizTestRouter(quizSvc, topicSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions/by-ids?ids=abc,xyz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
