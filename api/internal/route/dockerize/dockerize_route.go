package dockerize

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

func CreateDockerfile(context *gin.Context) {
	exposePortsCmds := generateExposedPortsCommands("exposed-ports.yaml")
	// Generate the run services commands
	runServicesCmds := generateRunServiceCommands("sys-services.yaml")
	// Save the run services commands to a .sh file
	saveRunServicesToSh(runServicesCmds, "run-services.sh")
	sourceTarPath := "source-vm-fs.tar.gz"
	runServicesPath := "run-services.sh"
	dockerfilePath := "Dockerfile"

	// Create the output directory
	utils.CreateBaseOutputDir(BASE_DOCKERIZE_OUTPUT_DIR)

	// Create a Dockerfile using the tar archieve, the expose ports commands, and execute CMD using run services commands .sh file
	generateDockerfile(sourceTarPath, exposePortsCmds, runServicesPath, dockerfilePath)

	context.JSON(http.StatusOK, Response{Message: "Successfully generated the Dockerfile!"})
}

func CreateDockerImage(context *gin.Context) {
	// Create the output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	// Build the Docker image using the Dockerfile
	result, err := buildDockerImage("Dockerfile")
	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to build the Docker image", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: "Successfully built the Docker image with the following name: " + result})
	}
}

func CreateDockerContainer(context *gin.Context) {
	// Create the output directory
	utils.CreateBaseOutputDir(BASE_ANALYZE_OUTPUT_DIR)

	// Create a Docker container using the Docker image
	result, err := runDockerContainer("dockerized-vm", "dockerized-vm")
	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to create the Docker container", Error: err.Error()})
	} else {
		context.JSON(http.StatusOK, Response{Message: "Successfully created the Docker container with the following name: " + result})
	}
}

func CreateCompleteDockerization(context *gin.Context) {
	// Create a tar archive of the source VM's file system
	createTarArchieve(BASE_ANALYZE_OUTPUT_DIR, BASE_DOCKERIZE_OUTPUT_DIR)

	exposePortsCmds := generateExposedPortsCommands("exposed-ports.yaml")
	// Generate the run services commands
	runServicesCmds := generateRunServiceCommands("sys-services.yaml")
	// Save the run services commands to a .sh file
	saveRunServicesToSh(runServicesCmds, "run-services.sh")
	sourceTarPath := "source-vm-fs.tar.gz"
	runServicesPath := "run-services.sh"
	dockerfilePath := "Dockerfile"

	// Create the output directory
	utils.CreateBaseOutputDir(BASE_DOCKERIZE_OUTPUT_DIR)

	// Create a Dockerfile using the tar archieve, the expose ports commands, and execute CMD using run services commands .sh file
	generateDockerfile(sourceTarPath, exposePortsCmds, runServicesPath, dockerfilePath)

	_, err := buildDockerImage("Dockerfile")
	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to build the Docker image", Error: err.Error()})
		return
	}

	result, err := runDockerContainer("dockerized-vm", "exposed-ports.yaml")
	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{Message: "Failed to create the Docker container", Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, Response{Message: "Successfully created the Docker container with the following name: " + result})
}
