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

func AnalyzeApplicationFiles(context *gin.Context) {
	var req struct {
		User             string `json:"user" binding:"required"`
		Host             string `json:"host" binding:"required"`
		PrivateKeyPath   string `json:"privateKeyPath" binding:"required"`
		AnalyzerApproach string `json:"analyzerApproach" binding:"required"`
	}

	// Check if JSON request is valid
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	// Create base output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	analyzer, err := GetAnalyzerFactory(req.AnalyzerApproach)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to get analyzer", Error: err.Error()})
		return
	}

	sourceDir := "/"
	destinationDir := "source-vm-fs"
	result, err := analyzer.collectApplicationFiles(req.User, req.Host, sourceDir, destinationDir, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect file system", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}

func AnalyzeServices(context *gin.Context) {
	var req struct {
		User             string `json:"user" binding:"required"`
		Host             string `json:"host" binding:"required"`
		PrivateKeyPath   string `json:"privateKeyPath" binding:"required"`
		AnalyzerApproach string `json:"analyzerApproach" binding:"required"`
	}

	// Check if JSON request is valid
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	// Create base output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	analyzer, err := GetAnalyzerFactory(req.AnalyzerApproach)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to get analyzer", Error: err.Error()})
		return
	}

	result, err := analyzer.collectServices(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect system services", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}

func AnalyzeExposedPorts(context *gin.Context) {
	var req struct {
		User             string `json:"user" binding:"required"`
		Host             string `json:"host" binding:"required"`
		PrivateKeyPath   string `json:"privateKeyPath" binding:"required"`
		AnalyzerApproach string `json:"analyzerApproach" binding:"required"`
	}

	// Check if JSON request is valid
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, Response{Message: "Invalid request", Error: err.Error()})
		return
	}

	// Create base output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	analyzer, err := GetAnalyzerFactory(req.AnalyzerApproach)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to get analyzer", Error: err.Error()})
		return
	}

	result, err := analyzer.collectExposedPorts(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect exposed ports", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: result})
	}
}

func CreateCompleteAnalysisProfile(context *gin.Context) {
	var req struct {
		User             string `json:"user" binding:"required"`
		Host             string `json:"host" binding:"required"`
		PrivateKeyPath   string `json:"privateKeyPath" binding:"required"`
		AnalyzerApproach string `json:"analyzerApproach" binding:"required"`
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

	analyzer, err := GetAnalyzerFactory(req.AnalyzerApproach)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to get analyzer", Error: err.Error()})
		return
	}

	_, err = analyzer.collectApplicationFiles(req.User, req.Host, sourceDir, destinationDir, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect file system", Error: err.Error()})
		return
	}

	_, err = analyzer.collectServices(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect system services", Error: err.Error()})
		return
	}

	_, err = analyzer.collectExposedPorts(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect exposed ports", Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, Response{Message: "Successfully collected the file system, system services, and exposed ports"})
}
