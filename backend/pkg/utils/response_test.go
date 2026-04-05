package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	data := map[string]string{"key": "value"}
	SuccessResponse(c, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Nil(t, resp.Error)
}

func TestSuccessResponse_StatusCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/", nil)

	SuccessResponse(c, http.StatusCreated, gin.H{"id": 1})

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	ErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, "resource not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, ErrCodeNotFound, resp.Error.Code)
	assert.Equal(t, "resource not found", resp.Error.Message)
}

func TestErrorResponseWithDetails(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	ErrorResponseWithDetails(c, http.StatusBadRequest, ErrCodeValidation, "validation failed", "field 'name' is required")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, ErrCodeValidation, resp.Error.Code)
	assert.Equal(t, "field 'name' is required", resp.Error.Details)
}

func TestPaginatedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	items := []string{"a", "b", "c"}
	PaginatedResponse(c, items, 1, 10, 3)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Meta)
	assert.Equal(t, 1, resp.Meta.Page)
	assert.Equal(t, 10, resp.Meta.Limit)
	assert.Equal(t, 3, resp.Meta.Total)
	assert.Equal(t, 1, resp.Meta.TotalPages)
}

func TestPaginatedResponse_MultiplePages(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	PaginatedResponse(c, []int{}, 2, 5, 23)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 5, resp.Meta.TotalPages) // ceil(23/5) = 5
}

func TestErrorCodes(t *testing.T) {
	assert.Equal(t, "VALIDATION_ERROR", ErrCodeValidation)
	assert.Equal(t, "NOT_FOUND", ErrCodeNotFound)
	assert.Equal(t, "UNAUTHORIZED", ErrCodeUnauthorized)
	assert.Equal(t, "FORBIDDEN", ErrCodeForbidden)
	assert.Equal(t, "INTERNAL_ERROR", ErrCodeInternal)
	assert.Equal(t, "RATE_LIMIT_EXCEEDED", ErrCodeRateLimit)
	assert.Equal(t, "BAD_REQUEST", ErrCodeBadRequest)
}
