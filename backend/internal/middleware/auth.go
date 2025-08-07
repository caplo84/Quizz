package middleware

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

// AdminAuth is a middleware that checks for valid admin credentials
func AdminAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get credentials from environment
        username := os.Getenv("ADMIN_USERNAME")
        password := os.Getenv("ADMIN_PASSWORD")

        // Skip auth if credentials are not set
        if username == "" || password == "" {
            c.Next()
            return
        }

        // Get the Basic Auth credentials
        user, pass, hasAuth := c.Request.BasicAuth()

        if !hasAuth || user != username || pass != password {
            c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Unauthorized",
                "code":  http.StatusUnauthorized,
            })
            return
        }

        c.Next()
    }
}

// JWT will be implemented later
func JWT() gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: Implement JWT validation
        c.Next()
    }
}