package utils

import (
	"github.com/gin-gonic/gin"
)

// APIError represents a structured error response format
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface for APIError
func (e APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new APIError instance
func NewAPIError(code, message string) APIError {
	return APIError{
		Code:    code,
		Message: message,
	}
}

// RespondWithError sends a formatted error response with custom code and message
func RespondWithError(c *gin.Context, statusCode int, errorCode string, message string) {
	c.JSON(statusCode, gin.H{
		"error": APIError{
			Code:    errorCode,
			Message: message,
		},
	})
}

// RespondWithCustomError examines the error and sends a formatted error response.
// If the error is an APIError, it uses its code and message. Otherwise, it falls back to GENERIC_ERROR.
func RespondWithCustomError(c *gin.Context, statusCode int, err error) {
	if apiErr, ok := err.(APIError); ok {
		c.JSON(statusCode, gin.H{
			"error": apiErr,
		})
		return
	}
	
	c.JSON(statusCode, gin.H{
		"error": APIError{
			Code:    "GENERIC_ERROR",
			Message: err.Error(),
		},
	})
}
