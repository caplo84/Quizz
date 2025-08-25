package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/caplo84/quizz-backend/internal/logger"
)

// RequestLoggingMiddleware logs HTTP requests with structured format
func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()

		// Add request ID to context
		ctx := logger.SetRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// Add request ID to response headers for tracing
		c.Header("X-Request-ID", requestID)

		// Capture request body for logging (if needed)
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Record start time
		start := time.Now()

		// Log incoming request
		logger.Log.WithContext(ctx).WithFields(logger.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"user_agent": c.Request.UserAgent(),
			"ip":         c.ClientIP(),
			"request_id": requestID,
		}).Info("Incoming HTTP request")

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request completion
		logger.Log.LogRequest(ctx, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)

		// Log any errors that occurred
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Log.WithContext(ctx).WithError(err.Err).Error("Request processing error")
			}
		}
	}
}

// StructuredErrorMiddleware handles errors with structured logging
func StructuredErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			logger.Log.WithContext(c.Request.Context()).WithFields(logger.Fields{
				"error_type": err.Type,
				"error_meta": err.Meta,
			}).WithError(err.Err).Error("Request failed with error")
		}
	}
}

// SecurityLoggingMiddleware logs security-related events
func SecurityLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Log authentication attempts
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			logger.Log.LogSecurityEvent(ctx, "auth_attempt", logger.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
			})
		}

		// Check for suspicious patterns
		if isSuspiciousRequest(c) {
			logger.Log.LogSecurityEvent(ctx, "suspicious_request", logger.Fields{
				"path":       c.Request.URL.Path,
				"method":     c.Request.Method,
				"ip":         c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
				"query":      c.Request.URL.RawQuery,
			})
		}

		c.Next()
	}
}

// isSuspiciousRequest checks for suspicious request patterns
func isSuspiciousRequest(c *gin.Context) bool {
	// Add your security checks here
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	// Example: Check for SQL injection patterns
	suspiciousPatterns := []string{
		"' OR '1'='1",
		"UNION SELECT",
		"<script>",
		"javascript:",
		"../",
	}

	for _, pattern := range suspiciousPatterns {
		if contains(path, pattern) || contains(query, pattern) {
			return true
		}
	}

	return false
}

// Helper function to check if string contains substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
