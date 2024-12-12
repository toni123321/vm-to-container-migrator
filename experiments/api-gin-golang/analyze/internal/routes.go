package analyze

import "github.com/gin-gonic/gin"

func SetUpRoutes(router *gin.Engine) {
	router.GET("/analyze/fs", getFileSystem)
	router.GET("/analyze/services", getSystemServices)
	router.GET("analyze/ports", getExposedPorts)
}
