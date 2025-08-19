package handlers

import (
    "context"
    "net/http"
    "runtime"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
)

type HealthHandler struct {
    db        *gorm.DB
    redis     *redis.Client
    startTime time.Time
}

type HealthResponse struct {
    Status    string                 `json:"status"`
    Timestamp time.Time             `json:"timestamp"`
    Version   string                `json:"version"`
    Uptime    string                `json:"uptime"`
    Services  map[string]string     `json:"services"`
    Metrics   HealthMetrics         `json:"metrics"`
}

type HealthMetrics struct {
    DatabaseConnections int     `json:"database_connections"`
    MemoryUsageMB      float64 `json:"memory_usage_mb"`
    Goroutines         int     `json:"goroutines"`
    CPUCores           int     `json:"cpu_cores"`
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
    return &HealthHandler{
        db:        db,
        redis:     redis,
        startTime: time.Now(),
    }
}

// HealthCheck provides comprehensive health status
func (h *HealthHandler) HealthCheck(c *gin.Context) {
    response := HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Version:   "1.0.0",
        Uptime:    time.Since(h.startTime).String(),
        Services:  make(map[string]string),
    }

    // Check database
    if err := h.checkDatabase(); err != nil {
        response.Status = "unhealthy"
        response.Services["database"] = "error: " + err.Error()
    } else {
        response.Services["database"] = "healthy"
    }

    // Check Redis
    if err := h.checkRedis(); err != nil {
        if response.Status == "healthy" {
            response.Status = "degraded"
        }
        response.Services["redis"] = "error: " + err.Error()
    } else {
        response.Services["redis"] = "healthy"
    }

    // Add system metrics
    response.Metrics = h.getSystemMetrics()

    statusCode := http.StatusOK
    if response.Status == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }

    c.JSON(statusCode, response)
}

// LivenessProbe for Kubernetes
func (h *HealthHandler) LivenessProbe(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":    "alive",
        "timestamp": time.Now(),
    })
}

// ReadinessProbe for Kubernetes
func (h *HealthHandler) ReadinessProbe(c *gin.Context) {
    if err := h.checkDatabase(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "not ready",
            "reason": "database unavailable",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    "ready",
        "timestamp": time.Now(),
    })
}

func (h *HealthHandler) checkDatabase() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    sqlDB, err := h.db.DB()
    if err != nil {
        return err
    }

    return sqlDB.PingContext(ctx)
}

func (h *HealthHandler) checkRedis() error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return h.redis.Ping(ctx).Err()
}

func (h *HealthHandler) getSystemMetrics() HealthMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    sqlDB, _ := h.db.DB()
    dbStats := sqlDB.Stats()

    return HealthMetrics{
        DatabaseConnections: dbStats.OpenConnections,
        MemoryUsageMB:      float64(m.Alloc) / 1024 / 1024,
        Goroutines:         runtime.NumGoroutine(),
        CPUCores:           runtime.NumCPU(),
    }
}