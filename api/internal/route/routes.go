package route

import (
	"vm2cont/api/internal/route/analyze"
	"vm2cont/api/internal/route/dockerize"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	// Routes for the analyze endpoint
	router.POST("/analyze/fs", analyze.AnalyzeFileSystem)
	router.POST("/analyze/services", analyze.AnalyzeSystemServices)
	router.POST("analyze/ports", analyze.AnalyzeExposedPorts)
	router.POST("analyze/complete", analyze.CreateCompleteAnalysisProfile)

	// Routes for the dockerize endpoint
	router.POST("/dockerize/dockerfile", dockerize.CreateDockerfile)
	router.POST("/dockerize/image", dockerize.CreateDockerImage)
	router.POST("/dockerize/container", dockerize.CreateDockerContainer)
	router.POST("dockerize/complete", dockerize.CreateCompleteDockerization)
}
