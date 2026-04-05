package handlers

import (
	"bytes"
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

// MockAttemptService implements services.AttemptService for testing
type MockAttemptService struct {
	mock.Mock
}

func (m *MockAttemptService) CreateAttempt(ctx context.Context, attempt *models.Attempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockAttemptService) GetAttemptByID(ctx context.Context, id uint) (*models.Attempt, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attempt), args.Error(1)
}

func (m *MockAttemptService) UpdateAttemptAnswers(ctx context.Context, attempt *models.Attempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func newAttemptTestRouter(attemptSvc *MockAttemptService, quizSvc *MockQuizService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewAttemptHandler(attemptSvc, quizSvc)
	r := gin.New()
	r.POST("/quizzes/:slug/attempts", h.CreateAttempt)
	r.PUT("/quizzes/:slug/attempts/:id", h.SubmitAttempt)
	r.GET("/quizzes/:slug/attempts/:id", h.GetAttempt)
	return r
}

// ---- CreateAttempt ----

func TestAttemptHandler_CreateAttempt_Success(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	quiz := &models.Quiz{ID: 1, Slug: "go-basics", Topic: models.Topic{Name: "Go"}}
	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)
	attemptSvc.On("CreateAttempt", mock.Anything, mock.AnythingOfType("*models.Attempt")).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/quizzes/go-basics/attempts", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	quizSvc.AssertExpectations(t)
	attemptSvc.AssertExpectations(t)
}

func TestAttemptHandler_CreateAttempt_QuizNotFound(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	quizSvc.On("GetQuizBySlug", mock.Anything, "missing").Return((*models.Quiz)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/quizzes/missing/attempts", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAttemptHandler_CreateAttempt_ServiceError(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	quiz := &models.Quiz{ID: 1, Slug: "go-basics", Topic: models.Topic{Name: "Go"}}
	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)
	attemptSvc.On("CreateAttempt", mock.Anything, mock.AnythingOfType("*models.Attempt")).Return(errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/quizzes/go-basics/attempts", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ---- GetAttempt ----

func TestAttemptHandler_GetAttempt_Success(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	attempt := &models.Attempt{ID: 5, QuizID: 1}
	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}

	attemptSvc.On("GetAttemptByID", mock.Anything, uint(5)).Return(attempt, nil)
	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/go-basics/attempts/5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAttemptHandler_GetAttempt_InvalidID(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/go-basics/attempts/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAttemptHandler_GetAttempt_NotFound(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	attemptSvc.On("GetAttemptByID", mock.Anything, uint(99)).Return((*models.Attempt)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/go-basics/attempts/99", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAttemptHandler_GetAttempt_QuizMismatch(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	// Attempt belongs to quiz 2, but quiz slug resolves to quiz 1
	attempt := &models.Attempt{ID: 5, QuizID: 2}
	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}

	attemptSvc.On("GetAttemptByID", mock.Anything, uint(5)).Return(attempt, nil)
	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quizzes/go-basics/attempts/5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---- SubmitAttempt ----

func TestAttemptHandler_SubmitAttempt_Success(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	attempt := &models.Attempt{ID: 5, QuizID: 1}
	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}

	attemptSvc.On("GetAttemptByID", mock.Anything, uint(5)).Return(attempt, nil)
	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)
	attemptSvc.On("UpdateAttemptAnswers", mock.Anything, mock.AnythingOfType("*models.Attempt")).Return(nil)

	body := map[string]interface{}{
		"answers": map[string]interface{}{"1": "a", "2": "b"},
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/quizzes/go-basics/attempts/5", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	attemptSvc.AssertExpectations(t)
}

func TestAttemptHandler_SubmitAttempt_InvalidID(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/quizzes/go-basics/attempts/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAttemptHandler_SubmitAttempt_AttemptNotFound(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	attemptSvc.On("GetAttemptByID", mock.Anything, uint(5)).Return((*models.Attempt)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/quizzes/go-basics/attempts/5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAttemptHandler_SubmitAttempt_InvalidBody(t *testing.T) {
	attemptSvc := new(MockAttemptService)
	quizSvc := new(MockQuizService)
	router := newAttemptTestRouter(attemptSvc, quizSvc)

	attempt := &models.Attempt{ID: 5, QuizID: 1}
	quiz := &models.Quiz{ID: 1, Slug: "go-basics"}

	attemptSvc.On("GetAttemptByID", mock.Anything, uint(5)).Return(attempt, nil)
	quizSvc.On("GetQuizBySlug", mock.Anything, "go-basics").Return(quiz, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/quizzes/go-basics/attempts/5", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
