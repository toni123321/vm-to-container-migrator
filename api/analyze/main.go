package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

const BASE_OUTPUT_DIR = "./migration-output/"

type Service struct {
	Name     string `yaml:"name"`
	SubState string `yaml:"substate"`
}

// ServicesData is the top-level structure to store services
type ServicesData struct {
	Services []Service `yaml:"services"`
}

func create_base_output_dir() {
	// create the base output directory if it does not exist
	if _, err := os.Stat(BASE_OUTPUT_DIR); os.IsNotExist(err) {
		os.Mkdir(BASE_OUTPUT_DIR, 0755)
	}
}

func collect_fs(user string, host string, sourceDir string, destinationDir string, privateKeyPath string) {
	// source includes the user and host of the VM, plus the source directory
	var source string = fmt.Sprintf("%s@%s:%s", user, host, sourceDir)
	// destination includes the destination directory, concatenated with the base output directory
	var destination string = BASE_OUTPUT_DIR + destinationDir

	excludes := []string{
		"--exclude=/proc/*",
		"--exclude=/boot/*",
		"--exclude=/sys/*",
		"--exclude=/dev/*",
		"--exclude=/lib/modules/*",
		"--exclude=/usr/share/man/*",
		"--exclude=/usr/share/doc/*",
		"--exclude=/var/cache/*",
		"--exclude=/var/backups/*",
		"--exclude=/var/log/*",
		"--exclude=/var/tmp/*",
		"--exclude=/var/run/*",
		"--exclude=/var/lib/lxcfs/*",
		"--exclude=/run/*",
	}

	args := []string{
		"-avz",                                                                     // Use archive mode, verbose, compress
		"--progress",                                                               // Show progress during transfer
		"--stats",                                                                  // Show file transfer statistics
		"-e", fmt.Sprintf("ssh -i %s -o StrictHostKeyChecking=no", privateKeyPath), // Use ssh for secure connection
		"--rsync-path", "sudo rsync", // Use sudo to run rsync with root permissions
	}

	args = append(args, excludes...) // Add exclude options
	args = append(args, source, destination)

	// Run rsync command
	fmt.Printf("Running rsync command: rsync %v\n", args)

	fmt.Println("Gathering file system...")

	cmd := exec.Command("rsync", args...)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error runing rsync: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File system synchronized successfully!")
}

func parse_sys_services(output []string) []Service {
	var services []Service
	for _, service := range output {
		if service != "" {
			fields := strings.Fields(service)
			if len(fields) >= 4 {
				name := fields[0]
				subState := fields[3]
				services = append(services, Service{Name: name, SubState: subState})
			}
		}
	}
	return services
}

func save_sys_services(services []Service, path string) {
	serviceData := ServicesData{Services: services}
	yamFile, err := yaml.Marshal(serviceData)
	if err != nil {
		log.Fatalf("Failed to move services to YAML format: %v", err)
	}

	// concatenate the base output directory with the path
	path = BASE_OUTPUT_DIR + path
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(yamFile)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

}

func collect_sys_services(user string, host string, port string, privateKeyPath string) {

	// Create SSH client configuration
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the VM
	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer client.Close()

	fmt.Println("Gatheirng system services...")

	// Run systemctl command to get active services
	cmd := "systemctl --type=service --state=active --no-pager --no-legend"
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Run the command
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}

	// Parse the output by sending it as list of service strings
	services := parse_sys_services(strings.Split(string(output), "\n"))

	// Save the services to a YAML file
	save_sys_services(services, "sys-services.yaml")

	fmt.Println("System services collected successfully!")
}

func main() {
	fmt.Println("Analyze module!")

	// Create the base output directory if it does not exist
	create_base_output_dir()

	// Step 1: Collect file system from source VM using default filters
	collect_fs(
		"antoniomihailov2001",
		"34.173.30.91",
		"/",
		"source-vm-fs",
		"/home/toni/.ssh/id_ed25519_gcloud_source_vm",
	)

	// Step 2: Collect active and running system services from source VM
	collect_sys_services(
		"antoniomihailov2001",
		"34.173.30.91",
		"22",
		"/home/toni/.ssh/id_ed25519_gcloud_source_vm",
	)
}
