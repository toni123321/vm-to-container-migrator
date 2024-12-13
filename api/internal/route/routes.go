package route

import (
	"vm2cont/api/internal/route/analyze"
	"vm2cont/api/internal/route/dockerize"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	// Routes for the analyze endpoint
	router.POST("/analyze/fs", analyze.AnalyzeApplicationFiles)
	router.POST("/analyze/services", analyze.AnalyzeServices)
	router.POST("analyze/ports", analyze.AnalyzeExposedPorts)
	router.POST("analyze/complete/single-approach", analyze.CreateCompleteAnalysisProfile)
	router.POST("analyze/complete/mixed-approach", analyze.CreateCompleteAnalysisProfileMixedApproach)

	// Routes for the dockerize endpoint
	router.POST("/dockerize/dockerfile", dockerize.CreateDockerfile)
	router.POST("/dockerize/image", dockerize.CreateDockerImage)
	router.POST("/dockerize/container", dockerize.CreateDockerContainer)
	router.POST("dockerize/complete", dockerize.CreateCompleteDockerization)
}
