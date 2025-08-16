package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/caplo84/quizz-backend/internal/services"
    "github.com/caplo84/quizz-backend/pkg/utils"
)

type TopicHandler struct {
    Service services.TopicService
}

func NewTopicHandler(service services.TopicService) *TopicHandler {
    return &TopicHandler{Service: service}
}

func (h *TopicHandler) GetTopics(c *gin.Context) {
    topics, err := h.Service.GetAllTopics(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    utils.SuccessResponse(c, http.StatusOK, topics)
}