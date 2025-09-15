package handlers

import (
	"net/http"
	"strconv"

	"github.com/caplo84/quizz-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService  services.QuizService
	topicService services.TopicService
}

func NewQuizHandler(quizService services.QuizService, topicService services.TopicService) *QuizHandler {
	return &QuizHandler{
		quizService:  quizService,
		topicService: topicService,
	}
}

// GetQuizzes handles GET /quizzes and GET /topics/:topic/quizzes
func (h *QuizHandler) GetQuizzes(c *gin.Context) {
	var topicID uint

	// Check if this is the /topics/:topic/quizzes route
	topicSlug := c.Param("topic")
	if topicSlug != "" {
		// Get topic ID from slug
		topic, err := h.topicService.GetTopicBySlug(c.Request.Context(), topicSlug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Topic not found",
			})
			return
		}
		topicID = topic.ID
	} else {
		// Fallback to optional topic_id query parameter for /quizzes route
		topicIDParam := c.Query("topic_id")
		if topicIDParam != "" {
			if id, err := strconv.ParseUint(topicIDParam, 10, 32); err == nil {
				topicID = uint(id)
			}
		}
	}

	quizzes, err := h.quizService.GetQuizzesByTopic(c.Request.Context(), topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch quizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": quizzes,
	})
}

// GetQuizBySlug handles GET /quizzes/:slug
func (h *QuizHandler) GetQuizBySlug(c *gin.Context) {
	slug := c.Param("slug")

	quiz, err := h.quizService.GetQuizBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Quiz not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": quiz,
	})
}

// GetQuizQuestions handles GET /quizzes/:slug/questions
func (h *QuizHandler) GetQuizQuestions(c *gin.Context) {
	slug := c.Param("slug")

	quiz, err := h.quizService.GetQuizBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Quiz not found",
		})
		return
	}

	questions, err := h.quizService.GetQuizQuestions(c.Request.Context(), quiz.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch questions",
		})
		return
	}

	// Remove correct answers from response
	for i := range questions {
		for j := range questions[i].Choices {
			questions[i].Choices[j].IsCorrect = false
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": questions,
	})
}
