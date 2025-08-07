package middleware

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

// Logger is a middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Stop timer
        latency := time.Since(start)
        if latency > time.Minute {
            latency = latency.Truncate(time.Second)
        }

        // Collect log fields
        statusCode := c.Writer.Status()
        clientIP := c.ClientIP()
        method := c.Request.Method
        errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

        // Prepare log event
        event := log.Info()
        if statusCode >= 400 {
            event = log.Error().Str("error", errorMessage)
        }

        // Log the request
        event.
            Str("method", method).
            Str("path", path).
            Int("status", statusCode).
            Str("ip", clientIP).
            Str("latency", fmt.Sprintf("%v", latency)).
            Str("user-agent", c.Request.UserAgent()).
            Str("referer", c.Request.Referer()).
            Str("query", raw).
            Msg("Request processed")
    }
}