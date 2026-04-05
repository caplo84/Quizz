package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthRouter() *gin.Engine {
	r := gin.New()
	r.GET("/protected", AdminAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func TestAdminAuth_MissingAPIKey(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "")

	router := setupAuthRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminAuth_MissingAuthorizationHeader(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "test-secret-key")

	router := setupAuthRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminAuth_InvalidFormat(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "test-secret-key")

	router := setupAuthRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat test-secret-key")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminAuth_WrongAPIKey_Bearer(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "correct-key")

	router := setupAuthRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer wrong-key")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminAuth_CorrectAPIKey_Bearer(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "my-secret")

	router := setupAuthRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer my-secret")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminAuth_CorrectAPIKey_ApiKey(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "my-secret")

	router := setupAuthRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "ApiKey my-secret")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminAuth_WrongAPIKey_ApiKey(t *testing.T) {
	t.Setenv("ADMIN_API_KEY", "correct-key")

	router := setupAuthRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "ApiKey wrong-key")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
