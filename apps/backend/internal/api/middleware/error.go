package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ErrorHandler middleware handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "internal_server_error",
					Message: "An internal server error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// JSONError is a helper function to return JSON error responses
func JSONError(c *gin.Context, statusCode int, errorType string, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}
