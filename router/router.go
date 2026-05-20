package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Expense Management API is running",
		})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	return r
}
