package handlers

import (
	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/caplo84/quizz-backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AdminHandler struct {
	adminService      services.AdminService
	githubSyncService services.GitHubSyncService
	questionCorrector services.QuestionCorrector
}

func NewAdminHandler(adminService services.AdminService, githubSyncService services.GitHubSyncService, questionCorrector services.QuestionCorrector) *AdminHandler {
	return &AdminHandler{
		adminService:      adminService,
		githubSyncService: githubSyncService,
		questionCorrector: questionCorrector,
	}
}

type QuestionCorrectionRequest struct {
	QuizSlug            string  `json:"quiz_slug"`
	DryRun              *bool   `json:"dry_run"`
	BatchSize           *int    `json:"batch_size"`
	ConfidenceThreshold float64 `json:"confidence_threshold"`
	Verbose             bool    `json:"verbose"`
	ReviewOnly          bool    `json:"review_only"`
}

type AISettingsUpdateRequest struct {
	Provider            string `json:"provider"`
	CloudflareAPIToken  string `json:"cloudflare_api_token"`
	CloudflareAccountID string `json:"cloudflare_account_id"`
	CloudflareModel     string `json:"cloudflare_ai_model"`
	OllamaBaseURL       string `json:"ollama_base_url"`
	OllamaModel         string `json:"ollama_model"`
	CloudflareAIBaseURL string `json:"cloudflare_ai_base_url"`
}

// CreateQuiz handles POST /admin/quizzes
func (h *AdminHandler) CreateQuiz(c *gin.Context) {
	var quiz models.Quiz

	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.adminService.CreateQuiz(c.Request.Context(), &quiz); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create quiz",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": quiz,
	})
}

// UpdateQuiz handles PUT /admin/quizzes/:id
func (h *AdminHandler) UpdateQuiz(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid quiz ID",
		})
		return
	}

	var quiz models.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	quiz.ID = uint(id)

	if err := h.adminService.UpdateQuiz(c.Request.Context(), &quiz); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": quiz,
	})
}

// DeleteQuiz handles DELETE /admin/quizzes/:id
func (h *AdminHandler) DeleteQuiz(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid quiz ID",
		})
		return
	}

	if err := h.adminService.DeleteQuiz(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz deleted successfully",
	})
}

// GetQuizByID handles GET /admin/quizzes/:id
func (h *AdminHandler) GetQuizByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid quiz ID",
		})
		return
	}

	quiz, err := h.adminService.GetQuizByID(c.Request.Context(), uint(id))
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

// CreateTopic handles POST /admin/topics
func (h *AdminHandler) CreateTopic(c *gin.Context) {
	var topic models.Topic

	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := h.adminService.CreateTopic(c.Request.Context(), &topic); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create topic",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": topic,
	})
}

// UpdateTopic handles PUT /admin/topics/:id
func (h *AdminHandler) UpdateTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid topic ID",
		})
		return
	}

	var topic models.Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	topic.ID = uint(id)

	if err := h.adminService.UpdateTopic(c.Request.Context(), &topic); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update topic",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": topic,
	})
}

// DeleteTopic handles DELETE /admin/topics/:id
func (h *AdminHandler) DeleteTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid topic ID",
		})
		return
	}

	if err := h.adminService.DeleteTopic(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete topic",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topic deleted successfully",
	})
}

// SyncGitHubData handles POST /admin/sync-github
func (h *AdminHandler) SyncGitHubData(c *gin.Context) {
	ctx := c.Request.Context()

	// Trigger GitHub sync
	err := h.githubSyncService.SyncFromGitHub(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "GitHub sync failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "GitHub sync completed successfully",
		"timestamp": time.Now(),
	})
}

// GetGitHubSyncStatus handles GET /admin/sync-github/status
func (h *AdminHandler) GetGitHubSyncStatus(c *gin.Context) {
	// Check rate limits
	rateLimit, err := h.githubSyncService.GetRateLimit(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rate limit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rate_limit": rateLimit,
		"last_sync":  "implement this", // You can track this in database
	})
}

// DownloadAllTopicImages handles POST /admin/download-all-topic-images
func (h *AdminHandler) DownloadAllTopicImages(c *gin.Context) {
	err := h.adminService.DownloadAllTopicImages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to download topic images",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All topic images download completed successfully",
	})
}

// CorrectQuestions handles POST /api/admin/questions/correct
func (h *AdminHandler) CorrectQuestions(c *gin.Context) {
	if h.questionCorrector == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Question corrector service is not configured",
		})
		return
	}

	var req QuestionCorrectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// allow empty body by ignoring EOF-style bind errors
		req = QuestionCorrectionRequest{}
	}

	dryRun := true
	if req.DryRun != nil {
		dryRun = *req.DryRun
	}

	batchSize := 100
	if req.BatchSize != nil && *req.BatchSize > 0 {
		batchSize = *req.BatchSize
	}

	confidenceThreshold := req.ConfidenceThreshold
	if confidenceThreshold <= 0 {
		confidenceThreshold = 0.7
	}

	opts := services.CorrectionOptions{
		QuizSlug:            req.QuizSlug,
		DryRun:              dryRun,
		BatchSize:           batchSize,
		ConfidenceThreshold: confidenceThreshold,
		Verbose:             req.Verbose,
		ReviewOnly:          req.ReviewOnly,
	}

	report, err := h.questionCorrector.CorrectAllQuizzes(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Question correction failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Question correction completed",
		"report":  report,
	})
}

// GetAISettings handles GET /api/admin/ai/settings
func (h *AdminHandler) GetAISettings(c *gin.Context) {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("AI_PROVIDER")))
	if provider == "" {
		provider = "ollama"
	}

	c.JSON(http.StatusOK, gin.H{
		"provider": provider,
		"cloudflare": gin.H{
			"account_id":       strings.TrimSpace(os.Getenv("CLOUDFLARE_ACCOUNT_ID")),
			"model":            firstNonEmpty(strings.TrimSpace(os.Getenv("CLOUDFLARE_AI_MODEL")), "@cf/meta/llama-3.1-8b-instruct"),
			"base_url":         strings.TrimSpace(os.Getenv("CLOUDFLARE_AI_BASE_URL")),
			"token_configured": strings.TrimSpace(os.Getenv("CLOUDFLARE_API_TOKEN")) != "" || strings.TrimSpace(os.Getenv("CF_API_TOKEN")) != "",
		},
		"ollama": gin.H{
			"base_url": firstNonEmpty(strings.TrimSpace(os.Getenv("OLLAMA_BASE_URL")), "http://localhost:11434"),
			"model":    firstNonEmpty(strings.TrimSpace(os.Getenv("OLLAMA_MODEL")), "qwen2.5:7b"),
		},
	})
}

// UpdateAISettings handles PUT /api/admin/ai/settings
func (h *AdminHandler) UpdateAISettings(c *gin.Context) {
	var req AISettingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	provider := strings.ToLower(strings.TrimSpace(req.Provider))
	if provider != "" {
		switch provider {
		case "ollama", "cloudflare":
			_ = os.Setenv("AI_PROVIDER", provider)
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unsupported provider. Allowed values: ollama, cloudflare",
			})
			return
		}
	}

	if v := strings.TrimSpace(req.CloudflareAPIToken); v != "" {
		_ = os.Setenv("CLOUDFLARE_API_TOKEN", v)
		_ = os.Setenv("CF_API_TOKEN", v)
	}
	if v := strings.TrimSpace(req.CloudflareAccountID); v != "" {
		_ = os.Setenv("CLOUDFLARE_ACCOUNT_ID", v)
		_ = os.Setenv("CF_ACCOUNT_ID", v)
	}
	if v := strings.TrimSpace(req.CloudflareModel); v != "" {
		_ = os.Setenv("CLOUDFLARE_AI_MODEL", v)
	}
	if v := strings.TrimSpace(req.CloudflareAIBaseURL); v != "" {
		_ = os.Setenv("CLOUDFLARE_AI_BASE_URL", v)
	}
	if v := strings.TrimSpace(req.OllamaBaseURL); v != "" {
		_ = os.Setenv("OLLAMA_BASE_URL", v)
	}
	if v := strings.TrimSpace(req.OllamaModel); v != "" {
		_ = os.Setenv("OLLAMA_MODEL", v)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "AI settings updated",
	})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
