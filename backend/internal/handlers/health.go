package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "context"
)

type HealthHandler struct {
    DB    *gorm.DB
    Redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
    return &HealthHandler{DB: db, Redis: redis}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
    // Check DB
    sqlDB, err := h.DB.DB()
    if err != nil || sqlDB.Ping() != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": "down"})
        return
    }
    // Check Redis
    if err := h.Redis.Ping(context.Background()).Err(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "redis": "down"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}