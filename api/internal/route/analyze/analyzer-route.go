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

func CreateCompleteAnalysisProfileMixedApproach(context *gin.Context) {
	var req struct {
		User                    string `json:"user" binding:"required"`
		Host                    string `json:"host" binding:"required"`
		PrivateKeyPath          string `json:"privateKeyPath" binding:"required"`
		ApplicationFileStrategy string `json:"applicationFileStrategy" binding:"required"`
		ExposedPortsStrategy    string `json:"exposedPortsStrategy" binding:"required"`
		ServicesStrategy        string `json:"servicesStrategy" binding:"required"`
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

	analyzer, err := GetAnalyzerFactory("mixed")

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to get analyzer", Error: err.Error()})
		return
	}

	mixedAnalyzer := analyzer.(*MixedAnalyzerImpl)

	// Check if JSON values for strategies are valid
	for _, strategy := range []string{req.ApplicationFileStrategy, req.ExposedPortsStrategy, req.ServicesStrategy} {
		if strategy != "fs" && strategy != "process" {
			context.JSON(http.StatusInternalServerError, Response{Message: "Invalid strategy", Error: "Invalid strategy"})
			return
		}
	}

	var strategyImplAppFiles IAnalyzerFactory
	var strategyImplExposedPorts IAnalyzerFactory
	var strategyImplServices IAnalyzerFactory

	if req.ApplicationFileStrategy == "fs" {
		strategyImplAppFiles = &FsAnalyzerImpl{}
	} else {
		strategyImplAppFiles = &ProcessAnalyzerImpl{}
	}

	if req.ExposedPortsStrategy == "fs" {
		strategyImplExposedPorts = &FsAnalyzerImpl{}
	} else {
		strategyImplExposedPorts = &ProcessAnalyzerImpl{}
	}

	if req.ServicesStrategy == "fs" {
		strategyImplServices = &FsAnalyzerImpl{}
	} else {
		strategyImplServices = &ProcessAnalyzerImpl{}
	}

	mixedAnalyzer.SetApplicationFileStrategy(strategyImplAppFiles)
	mixedAnalyzer.SetExposedPortsStrategy(strategyImplExposedPorts)
	mixedAnalyzer.SetServicesStrategy(strategyImplServices)

	applicationFiles, err := mixedAnalyzer.collectApplicationFiles(req.User, req.Host, sourceDir, destinationDir, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect file system", Error: err.Error()})
		return
	}

	services, err := mixedAnalyzer.collectServices(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect system services", Error: err.Error()})
		return
	}

	exposedPorts, err := mixedAnalyzer.collectExposedPorts(req.User, req.Host, req.PrivateKeyPath)

	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to collect exposed ports", Error: err.Error()})
		return
	}

	// Print JSON with all the results collected (applicationFiles, services, exposedPorts) following this format: {"applicationFiles": "result", "services": "result", "exposedPorts": "result"}
	context.JSON(http.StatusOK, gin.H{"applicationFiles": applicationFiles, "services": services, "exposedPorts": exposedPorts})
}
