package analyze

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const BASE_OUTPUT_DIR = "./migration-output/"

// Response model for API responses
type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func getFileSystem(context *gin.Context) {
	var req struct {
		User           string `json:"user" binding:"required"`
		Host           string `json:"host" binding:"required"`
		SourceDir      string `json:"sourceDir" binding:"required"`
		DestinationDir string `json:"destinationDir" binding:"required"`
		PrivateKeyPath string `json:"privateKeyPath" binding:"required"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	createBaseOutputDir(BASE_OUTPUT_DIR)
	result, err := collectFs(req.User, req.Host, req.SourceDir, req.DestinationDir, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect file system", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}

func getSystemServices(context *gin.Context) {
	var req struct {
		User           string `json:"user" binding:"required"`
		Host           string `json:"host" binding:"required"`
		PrivateKeyPath string `json:"privateKeyPath" binding:"required"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	createBaseOutputDir(BASE_OUTPUT_DIR)
	result, err := collectSysServices(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect system services", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}

func getExposedPorts(context *gin.Context) {
	var req struct {
		User           string `json:"user" binding:"required"`
		Host           string `json:"host" binding:"required"`
		PrivateKeyPath string `json:"privateKeyPath" binding:"required"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	createBaseOutputDir(BASE_OUTPUT_DIR)

	result, err := collectExposedPorts(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect exposed ports", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}
