package main

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

const BASE_INPUT_DIR = "../analyze/migration-output/"
const BASE_OUTPUT_DIR = "./dockerize-output/"

type Port struct {
	Protocol string `yaml:"protocol"`
	PortNr   string `yaml:"portNr"`
}

type PortsData struct {
	Ports []Port `yaml:"ports"`
}

type Service struct {
	Name     string `yaml:"name"`
	SubState string `yaml:"substate"`
	Command  string `yaml:"command"`
}

// ServicesData is the top-level structure to store services
type ServicesData struct {
	Services []Service `yaml:"services"`
}

// func convert_to_dockerfile(fsPath string, servicesPath string, portsPath string, dockerfilePath string) {

// }

func create_base_output_dir() {
	// create the base output directory if it does not exist
	if _, err := os.Stat(BASE_OUTPUT_DIR); os.IsNotExist(err) {
		os.Mkdir(BASE_OUTPUT_DIR, 0755)
	}
}

func allow_sudo_for_tar(usr string) {
	// Check if the current user can run tar without sudo already
	cmdCheck := "sudo -l | grep tar"
	cmdCheckExec := exec.Command("sh", "-c", cmdCheck)
	if err := cmdCheckExec.Run(); err == nil {
		fmt.Println("Current user can run tar without sudo!")
		return
	}

	// Current user can run tar without sudo through visudo
	cmd := fmt.Sprintf("echo \"%s ALL=(ALL) NOPASSWD: /bin/tar\" | sudo EDITOR='tee -a' visudo", usr)
	cmdExec := exec.Command("sh", "-c", cmd)

	fmt.Println("Allowing sudo for the `tar` command ...")

	err := cmdExec.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error allowing sudo for the `tar` command: %s\n", err)
	}

	fmt.Println("Successfully allowed sudo for the `tar` command!")
}
func create_tar_archive() {
	// Create a tar archive of the source VM's file system
	tarPath := fmt.Sprintf("%ssource-vm-fs.tar.gz", BASE_OUTPUT_DIR)
	fsPath := fmt.Sprintf("%ssource-vm-fs", BASE_INPUT_DIR)

	cmd := fmt.Sprintf("sudo tar -czf %s -C %s .", tarPath, fsPath)
	cmdExec := exec.Command("sh", "-c", cmd)

	fmt.Println("Creating the tar archieve ...")

	err := cmdExec.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tar archive: %s\n", err)
	}
	fmt.Println("Successfully created tar archive of the source VM's file system!")

}

func get_expose_ports(fname string) PortsData {
	// Read the file containing the exposed ports
	fullFnamePath := fmt.Sprintf("%s%s", BASE_INPUT_DIR, fname)
	data, err := os.ReadFile(fullFnamePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading the file: %s\n", err)
	}

	// Unmarshal YAML data into PortsData struct
	var portsData PortsData
	err = yaml.Unmarshal(data, &portsData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling the data: %s\n", err)
	}

	return portsData
}

func generate_expose_ports_commands(fname string) []string {
	portsData := get_expose_ports(fname)

	var exposePortsCmds []string
	for _, port := range portsData.Ports {
		exposePortsCmds = append(exposePortsCmds, fmt.Sprintf("EXPOSE %s", port.PortNr))
	}

	// Return a slice of strings containing the expose ports commands
	return exposePortsCmds
}

func generate_run_services_commands(fname string) []string {
	// Read the file containing the services
	fullFnamePath := fmt.Sprintf("%s%s", BASE_INPUT_DIR, fname)
	data, err := os.ReadFile(fullFnamePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading the file: %s\n", err)
	}

	// Unmarhsal YAML data into ServicesData struct
	var servicesData ServicesData
	err = yaml.Unmarshal(data, &servicesData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling the data: %s\n", err)
	}

	// Create a slice of strings containing the run services commands
	var runServicesCmds []string
	for _, service := range servicesData.Services {
		runServicesCmds = append(runServicesCmds, fmt.Sprintf("%s &", service.Command))
	}

	// Return a slice of strings containing the run services commands
	return runServicesCmds
}

func save_run_services_to_sh(commands []string, outputPath string) error {
	fullOutputPath := fmt.Sprintf("%s%s", BASE_OUTPUT_DIR, outputPath)
	// Create or open the .sh file
	file, err := os.Create(fullOutputPath)
	if err != nil {
		return fmt.Errorf("Error creating the .sh file: %s", err)
	}
	defer file.Close()

	// Write the shebang line
	_, err = file.WriteString("#!/bin/sh\n")
	if err != nil {
		return fmt.Errorf("Error writing to the .sh file: %s", err)
	}

	// Write each command to the file
	for _, cmd := range commands {
		_, err = file.WriteString(fmt.Sprintf("%s\n", cmd))
		if err != nil {
			return fmt.Errorf("Error writing to the .sh file: %s", err)
		}
	}

	// Keep the script running
	_, err = file.WriteString("tail -f /dev/null\n")
	if err != nil {
		return fmt.Errorf("Error writing to the .sh file: %s", err)
	}

	fmt.Println("Successfully saved run services commands to the .sh file!")
	return nil
}

func generate_dockerfile(tarPath string, exposePortsCmds []string, runServicesShPath string, outputPath string) {
	// Create or open the Dockerfile
	fullOutputPath := fmt.Sprintf("%s%s", BASE_OUTPUT_DIR, outputPath)
	file, err := os.Create(fullOutputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating the Dockerfile: %s\n", err)
	}
	defer file.Close()

	fmt.Println("Generating the Dockerfile ...")

	// Write the base image using Ubuntu 22.04 as an example
	_, err = file.WriteString("FROM ubuntu:22.04\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to the Dockerfile: %s\n", err)
	}

	// Copy the source VM's file system to the Docker image using ADD and the following format: ADD source-vm-fs.tar.gz /
	fullTarPath := fmt.Sprintf("%s%s", BASE_OUTPUT_DIR, tarPath)
	_, err = file.WriteString(fmt.Sprintf("ADD %s /\n", fullTarPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to the Dockerfile: %s\n", err)
	}

	fullRunServicesShPath := fmt.Sprintf("%s%s", BASE_OUTPUT_DIR, runServicesShPath)
	// Copy the sh script for running the services to the Docker image
	_, err = file.WriteString(fmt.Sprintf("COPY %s /\n", fullRunServicesShPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to the Dockerfile: %s\n", err)
	}

	// Write the expose ports commands to the Dockerfile
	for _, cmd := range exposePortsCmds {
		_, err = file.WriteString(fmt.Sprintf("%s\n", cmd))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to the Dockerfile: %s\n", err)
		}
	}

	execRunServicesPath := fmt.Sprintf("/%s", runServicesShPath)
	// Make the sh script for running the services executable
	_, err = file.WriteString(fmt.Sprintf("RUN chmod +x %s\n", execRunServicesPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to the Dockerfile: %s\n", err)
	}

	// Execute the sh script for running the services using ENTRYPOINT
	_, err = file.WriteString(fmt.Sprintf("ENTRYPOINT [\"%s\"]\n", execRunServicesPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to the Dockerfile: %s\n", err)
	}

	fmt.Println("Successfully generated the Dockerfile!")
}

func build_docker_image(dockerfilePath string) string {
	fullDockerfilePath := fmt.Sprintf("%s%s", BASE_OUTPUT_DIR, dockerfilePath)
	imageName := "dockerized-vm"

	// Build the Docker image using the Dockerfile
	cmd := exec.Command("docker", "build", "-t", imageName, "-f", fullDockerfilePath, ".")

	fmt.Println("Building the Docker image ...")

	// Capture the output of the command
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error building the Docker image: %s\n", err)
	}

	fmt.Println("Successfully built the Docker image!")

	return imageName
}

func run_docker_container(imageName string, fnameExposePorts string) {
	// Run the Docker container using the built image
	containerName := "dockerized-vm-container"

	portsData := get_expose_ports(fnameExposePorts)

	portPrefix := "80" // Default prefix for host ports
	exposePorts := []string{}
	for _, port := range portsData.Ports {
		hostPort := portPrefix + port.PortNr
		containerPort := port.PortNr
		exposePorts = append(exposePorts, "-p", fmt.Sprintf("%s:%s", hostPort, containerPort))
	}

	args := append([]string{"run", "-d", "--name", containerName}, exposePorts...)
	args = append(args, imageName)

	cmdExec := exec.Command("docker", args...)

	// Capture the output of the command
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running the Docker container: %s\n", err)
	}

	fmt.Println("Running the Docker container ...")
}

func main() {
	fmt.Println("Dockerize module!")

	// Create the base output directory if it does not exist
	create_base_output_dir()

	// Step 1: Create a tar archive of the source VM's file system
	allow_sudo_for_tar("toni")
	create_tar_archive()

	// Generate the expose ports commands
	exposePortsCmds := generate_expose_ports_commands("exposed-ports.yaml")
	// Generate the run services commands
	runServicesCmds := generate_run_services_commands("sys-services.yaml")
	// Save the run services commands to a .sh file
	save_run_services_to_sh(runServicesCmds, "run-services.sh")

	// Create a Dockerfile using the tar archieve, the expose ports commands, and execute CMD using run services commands .sh file
	generate_dockerfile("source-vm-fs.tar.gz", exposePortsCmds, "run-services.sh", "Dockerfile")

	// Build the Docker image using the Dockerfile
	dockerImgName := build_docker_image("Dockerfile")

	// Run the Docker container using the built image
	run_docker_container(dockerImgName, "exposed-ports.yaml")
}
