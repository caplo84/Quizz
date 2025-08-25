package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	redis *redis.Client
	limit rate.Limit
	burst int
}

func NewRateLimiter(redis *redis.Client, requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		redis: redis,
		limit: rate.Limit(requestsPerSecond),
		burst: burst,
	}
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Use Redis for distributed rate limiting
		key := "rate_limit:" + clientIP
		current, err := rl.redis.Incr(c.Request.Context(), key).Result()

		if err != nil {
			// If Redis fails, allow the request
			c.Next()
			return
		}

		if current == 1 {
			// First request, set expiration
			rl.redis.Expire(c.Request.Context(), key, time.Minute)
		}

		// Check if limit exceeded (60 requests per minute)
		if current > 60 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
