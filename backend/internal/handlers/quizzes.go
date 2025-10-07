package handlers

import (
	"net/http"
	"strconv"
	"strings"

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

	// Check if client wants to include correct answers (for review page)
	includeAnswers := c.Query("include_answers") == "true"
	
	if !includeAnswers {
		// Remove correct answers from response for quiz play
		for i := range questions {
			for j := range questions[i].Choices {
				questions[i].Choices[j].IsCorrect = false
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": questions,
	})
}

// GetRandomQuestions handles GET /topics/:topic/questions/random?limit=10&exclude=1,2,3
func (h *QuizHandler) GetRandomQuestions(c *gin.Context) {
	topicSlug := c.Param("topic")

	// Get topic by slug
	topic, err := h.topicService.GetTopicBySlug(c.Request.Context(), topicSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Topic not found",
		})
		return
	}

	// Parse limit parameter (default to 10)
	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	// Parse exclude parameter (comma-separated question IDs)
	var excludeQuestionIDs []uint
	if excludeParam := c.Query("exclude"); excludeParam != "" {
		excludeStrings := strings.Split(excludeParam, ",")
		for _, idStr := range excludeStrings {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
				excludeQuestionIDs = append(excludeQuestionIDs, uint(id))
			}
		}
	}

	questions, err := h.quizService.GetRandomQuestions(c.Request.Context(), topic.ID, limit, excludeQuestionIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch random questions",
		})
		return
	}

	// Check if client wants to include correct answers (for review page)
	includeAnswers := c.Query("include_answers") == "true"
	
	if !includeAnswers {
		// Remove correct answers from response for quiz play
		for i := range questions {
			for j := range questions[i].Choices {
				questions[i].Choices[j].IsCorrect = false
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": questions,
		"topic": gin.H{
			"id":   topic.ID,
			"name": topic.Name,
			"slug": topic.Slug,
		},
		"meta": gin.H{
			"total_questions": len(questions),
			"limit":           limit,
			"excluded_count":  len(excludeQuestionIDs),
		},
	})
}

// GetQuestionsByIDs handles GET /questions/by-ids?ids=1,2,3&include_answers=true
func (h *QuizHandler) GetQuestionsByIDs(c *gin.Context) {
	// Parse IDs parameter (comma-separated question IDs)
	idsParam := c.Query("ids")
	if idsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ids parameter is required",
		})
		return
	}

	var questionIDs []uint
	idStrings := strings.Split(idsParam, ",")
	for _, idStr := range idStrings {
		if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
			questionIDs = append(questionIDs, uint(id))
		}
	}

	if len(questionIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No valid question IDs provided",
		})
		return
	}

	questions, err := h.quizService.GetQuestionsByIDs(c.Request.Context(), questionIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch questions",
		})
		return
	}

	// Check if client wants to include correct answers (for review page)
	includeAnswers := c.Query("include_answers") == "true"
	
	if !includeAnswers {
		// Remove correct answers from response for quiz play
		for i := range questions {
			for j := range questions[i].Choices {
				questions[i].Choices[j].IsCorrect = false
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": questions,
	})
}
