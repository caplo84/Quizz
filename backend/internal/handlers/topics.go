package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/caplo84/quizz-backend/internal/services"
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
    c.JSON(http.StatusOK, topics)
}