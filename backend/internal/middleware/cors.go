package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

// CORS handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		environment := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
		if environment == "" {
			environment = strings.ToLower(strings.TrimSpace(os.Getenv("GIN_MODE")))
		}

		if allowedOrigins == "" {
			if environment == "production" || environment == "release" {
				allowedOrigins = ""
			} else {
				allowedOrigins = "*" // Development convenience fallback.
			}
		}

		origin := c.Request.Header.Get("Origin")
		allowCredentials := false
		if allowedOrigins == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origins := strings.Split(allowedOrigins, ",")
			for _, o := range origins {
				if strings.TrimSpace(o) == origin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					allowCredentials = true
					break
				}
			}
		}

		if allowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
