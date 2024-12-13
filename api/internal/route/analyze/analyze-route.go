package analyze

import (
	"net/http"
	"vm2cont/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// Response model for API responses
type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func AnalyzeFileSystem(context *gin.Context) {
	var req struct {
		User           string `json:"user" binding:"required"`
		Host           string `json:"host" binding:"required"`
		PrivateKeyPath string `json:"privateKeyPath" binding:"required"`
	}

	// Check if JSON request is valid
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	// Create base output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	sourceDir := "/"
	destinationDir := "source-vm-fs"
	result, err := collectFs(req.User, req.Host, sourceDir, destinationDir, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect file system", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}

func AnalyzeSystemServices(context *gin.Context) {
	var req struct {
		User           string `json:"user" binding:"required"`
		Host           string `json:"host" binding:"required"`
		PrivateKeyPath string `json:"privateKeyPath" binding:"required"`
	}

	// Check if JSON request is valid
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	// Create base output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	result, err := collectSysServices(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect system services", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}

func AnalyzeExposedPorts(context *gin.Context) {
	var req struct {
		User           string `json:"user" binding:"required"`
		Host           string `json:"host" binding:"required"`
		PrivateKeyPath string `json:"privateKeyPath" binding:"required"`
	}

	// Check if JSON request is valid
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	// Create base output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	result, err := collectExposedPorts(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect exposed ports", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}

func CreateCompleteAnalysisProfile(context *gin.Context) {
	var req struct {
		User           string `json:"user" binding:"required"`
		Host           string `json:"host" binding:"required"`
		PrivateKeyPath string `json:"privateKeyPath" binding:"required"`
	}

	// Check if JSON request is valid
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	// Create base output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	sourceDir := "/"
	destinationDir := "source-vm-fs"

	_, err := collectFs(req.User, req.Host, sourceDir, destinationDir, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect file system", Error: err.Error()})
		return
	}

	_, err = collectSysServices(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect system services", Error: err.Error()})
		return
	}

	_, err = collectExposedPorts(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect exposed ports", Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, Response{Message: "Successfully collected the file system, system services, and exposed ports"})
}
