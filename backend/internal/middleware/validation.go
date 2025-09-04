package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateJSON validates JSON request body
func ValidateJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// ValidateSlug validates URL parameters
func ValidateSlug() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		if slug != "" {
			// Basic slug validation
			if len(slug) < 1 || len(slug) > 100 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid slug format",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
