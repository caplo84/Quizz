package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/caplo84/quizz-backend/internal/services"
    "github.com/caplo84/quizz-backend/internal/models"
)

type AttemptHandler struct {
    attemptService services.AttemptService
    quizService    services.QuizService
}

func NewAttemptHandler(attemptService services.AttemptService, quizService services.QuizService) *AttemptHandler {
    return &AttemptHandler{
        attemptService: attemptService,
        quizService:    quizService,
    }
}

// CreateAttempt handles POST /quizzes/:slug/attempts
func (h *AttemptHandler) CreateAttempt(c *gin.Context) {
    slug := c.Param("slug")
    
    quiz, err := h.quizService.GetQuizBySlug(c.Request.Context(), slug)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Quiz not found",
        })
        return
    }
    
    attempt := &models.Attempt{
        QuizID: quiz.ID,
    }
    
    if err := h.attemptService.CreateAttempt(c.Request.Context(), attempt); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create attempt",
        })
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "data": attempt,
    })
}

// SubmitAttempt handles PUT /quizzes/:slug/attempts/:id
func (h *AttemptHandler) SubmitAttempt(c *gin.Context) {
    slug := c.Param("slug")
    attemptIDStr := c.Param("id")
    
    attemptID, err := strconv.ParseUint(attemptIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid attempt ID",
        })
        return
    }
    
    attempt, err := h.attemptService.GetAttemptByID(c.Request.Context(), uint(attemptID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Attempt not found",
        })
        return
    }
    
    quiz, err := h.quizService.GetQuizBySlug(c.Request.Context(), slug)
    if err != nil || quiz.ID != attempt.QuizID {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Quiz not found or attempt mismatch",
        })
        return
    }
    
    // Use a simple struct instead of models.Answer
    var requestBody struct {
        Answers models.UserAnswers `json:"answers"`
    }
    
    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }
    
    // Update attempt with answers as string or JSON
    attempt.Answers = requestBody.Answers
    
    if err := h.attemptService.UpdateAttemptAnswers(c.Request.Context(), attempt); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to submit attempt",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": attempt,
    })
}

// GetAttempt handles GET /quizzes/:slug/attempts/:id
func (h *AttemptHandler) GetAttempt(c *gin.Context) {
    slug := c.Param("slug")
    attemptIDStr := c.Param("id")
    
    attemptID, err := strconv.ParseUint(attemptIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid attempt ID",
        })
        return
    }
    
    attempt, err := h.attemptService.GetAttemptByID(c.Request.Context(), uint(attemptID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Attempt not found",
        })
        return
    }
    
    quiz, err := h.quizService.GetQuizBySlug(c.Request.Context(), slug)
    if err != nil || quiz.ID != attempt.QuizID {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Quiz not found or attempt mismatch",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": attempt,
    })
}