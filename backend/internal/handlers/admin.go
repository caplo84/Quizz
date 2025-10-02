package handlers

import (
	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/caplo84/quizz-backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type AdminHandler struct {
	adminService      services.AdminService
	githubSyncService services.GitHubSyncService
}

func NewAdminHandler(adminService services.AdminService, githubSyncService services.GitHubSyncService) *AdminHandler {
	return &AdminHandler{
		adminService:      adminService,
		githubSyncService: githubSyncService,
	}
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
			"error": "Failed to download topic images",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All topic images download completed successfully",
	})
}
