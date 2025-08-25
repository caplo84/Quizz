package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		// Strict transport security
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		// Content security policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:")
		// Prevent embedding in iframes
		c.Header("X-Content-Security-Policy", "frame-ancestors 'none'")
		// Prevent browsers from performing MIME sniffing
		c.Header("X-Download-Options", "noopen")
		// Disable caching of sensitive data
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Next()
	}
}
