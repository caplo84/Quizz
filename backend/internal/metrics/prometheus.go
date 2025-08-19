package metrics

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    // HTTP Request metrics
    HttpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "quiz_api_http_requests_total",
            Help: "Total HTTP requests processed",
        },
        []string{"method", "endpoint", "status"},
    )

    HttpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "quiz_api_http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1.0, 2.5, 5.0, 10.0},
        },
        []string{"method", "endpoint"},
    )

    // Quiz-specific business metrics
    QuizAttemptsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "quiz_attempts_total",
            Help: "Total quiz attempts",
        },
        []string{"topic", "completed"},
    )

    QuizScores = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "quiz_scores",
            Help: "Quiz completion scores",
            Buckets: []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
        },
        []string{"topic"},
    )

    DatabaseOperations = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "quiz_db_operation_duration_seconds",
            Help: "Database operation duration",
            Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1.0},
        },
        []string{"operation", "table"},
    )

    CacheOperations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "quiz_cache_operations_total",
            Help: "Cache operations",
        },
        []string{"operation", "result"}, // get/set, hit/miss
    )

    ActiveConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "quiz_api_active_connections",
            Help: "Number of active connections",
        },
    )
)

func init() {
    // Register all metrics with Prometheus
    prometheus.MustRegister(HttpRequestsTotal)
    prometheus.MustRegister(HttpRequestDuration)
    prometheus.MustRegister(QuizAttemptsTotal)
    prometheus.MustRegister(QuizScores)
    prometheus.MustRegister(DatabaseOperations)
    prometheus.MustRegister(CacheOperations)
    prometheus.MustRegister(ActiveConnections)
}

// PrometheusMiddleware collects HTTP metrics
func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Track active connections
        ActiveConnections.Inc()
        defer ActiveConnections.Dec()

        c.Next()

        // Record metrics after request completion
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        path := c.FullPath()
        if path == "" {
            path = "unknown"
        }

        HttpRequestsTotal.WithLabelValues(
            c.Request.Method,
            path,
            status,
        ).Inc()

        HttpRequestDuration.WithLabelValues(
            c.Request.Method,
            path,
        ).Observe(duration)
    }
}

// MetricsHandler exposes Prometheus metrics endpoint
func MetricsHandler() gin.HandlerFunc {
    h := promhttp.Handler()
    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}

// RecordQuizAttempt records business metrics for quiz attempts
func RecordQuizAttempt(topic string, completed bool, score float64) {
    completedStr := "false"
    if completed {
        completedStr = "true"
        QuizScores.WithLabelValues(topic).Observe(score)
    }
    
    QuizAttemptsTotal.WithLabelValues(topic, completedStr).Inc()
}

// RecordDatabaseOperation records database performance metrics
func RecordDatabaseOperation(operation, table string, duration time.Duration) {
    DatabaseOperations.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordCacheOperation records cache hit/miss metrics
func RecordCacheOperation(operation, result string) {
    CacheOperations.WithLabelValues(operation, result).Inc()
}