package middleware

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/pkg/errors"
)

// ErrorResponse represents the error response structure
type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
    Code    int    `json:"code"`
}

// ErrorHandler is a middleware that handles errors from HTTP handlers
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // Process request

        // Check if there are any errors to handle
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            var statusCode int
            var message string
            var details string

            // Handle different types of errors
            switch {
            case errors.Is(err, ErrNotFound):
                statusCode = http.StatusNotFound
                message = "Resource not found"
            case errors.Is(err, ErrUnauthorized):
                statusCode = http.StatusUnauthorized
                message = "Unauthorized"
            case errors.Is(err, ErrForbidden):
                statusCode = http.StatusForbidden
                message = "Forbidden"
            case errors.Is(err, ErrBadRequest):
                statusCode = http.StatusBadRequest
                message = "Bad request"
            case errors.Is(err, ErrValidation):
                statusCode = http.StatusUnprocessableEntity
                message = "Validation error"
            default:
                statusCode = http.StatusInternalServerError
                message = "Internal server error"
                details = err.Error()
                
                // Log unexpected errors
                log.Printf("Unexpected error: %+v", err)
            }

            // Send error response
            c.JSON(statusCode, ErrorResponse{
                Error:   message,
                Details: details,
                Code:    statusCode,
            })
            return
        }
    }
}

// Error types
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrBadRequest   = errors.New("bad request")
    ErrValidation   = errors.New("validation error")
)

// Error helper functions
func NewValidationError(err error) error {
    return errors.Wrap(ErrValidation, err.Error())
}

func NewNotFoundError(resource string) error {
    return errors.Wrapf(ErrNotFound, "%s not found", resource)
}

func NewUnauthorizedError(reason string) error {
    return errors.Wrapf(ErrUnauthorized, "unauthorized: %s", reason)
}