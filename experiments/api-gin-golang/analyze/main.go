package main

import (
	analyze "api/analyze/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router
	router := gin.Default()

	// Set up routes
	analyze.SetUpRoutes(router)

	// Add a health check route
	router.GET("/health", healthcheck)
	// Start the server
	router.Run(":8001")
}

func healthcheck(context *gin.Context) {
	// Return 204 No Content
	context.Status(204)
}
