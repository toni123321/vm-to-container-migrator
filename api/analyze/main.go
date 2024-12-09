package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

const BASE_OUTPUT_DIR = "./migration-output/"

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

func create_base_output_dir() {
	// create the base output directory if it does not exist
	if _, err := os.Stat(BASE_OUTPUT_DIR); os.IsNotExist(err) {
		os.Mkdir(BASE_OUTPUT_DIR, 0755)
	}
}

func connect_ssh(user string, host string, port string, privateKeyPath string) *ssh.Client {
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

	return client
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
	// Define the exclude service regex list
	excludeSvcRegexList := []string{
		`^apparmor\.service$`,
		`^apport\.service$`,
		`^blk-availability\.service$`,
		`^chrony\.service$`,
		`^cloud-config\.service$`,
		`^cloud-final\.service$`,
		`^cloud-init-local\.service$`,
		`^cloud-init\.service$`,
		`^console-setup\.service$`,
		`^cron\.service$`,
		`^dbus\.service$`,
		`^finalrd\.service$`,
		`^getty@tty1\.service$`,
		`^google-guest-agent\.service$`,
		`^google-osconfig-agent\.service$`,
		`^google-shutdown-scripts\.service$`,
		`^keyboard-setup\.service$`,
		`^kmod-static-nodes\.service$`,
		`^lvm2-monitor\.service$`,
		`^multipathd\.service$`,
		`^networkd-dispatcher\.service$`,
		`^packagekit\.service$`,
		`^plymouth-quit-wait\.service$`,
		`^plymouth-quit\.service$`,
		`^plymouth-read-write\.service$`,
		`^polkit\.service$`,
		`^rsyslog\.service$`,
		`^serial-getty@ttyS0\.service$`,
		`^setvtrgb\.service$`,
		`^snapd\.apparmor\.service$`,
		`^snapd\.seeded\.service$`,
		`^snapd\.service$`,
		`^ssh\.service$`,
		`^systemd-binfmt\.service$`,
		`^systemd-fsck-root\.service$`,
		`^systemd-fsck@dev-disk-by\\x2dlabel-UEFI\.service$`,
		`^systemd-journal-flush\.service$`,
		`^systemd-journald\.service$`,
		`^systemd-logind\.service$`,
		`^systemd-machine-id-commit\.service$`,
		`^systemd-modules-load\.service$`,
		`^systemd-networkd-wait-online\.service$`,
		`^systemd-networkd\.service$`,
		`^systemd-random-seed\.service$`,
		`^systemd-remount-fs\.service$`,
		`^systemd-resolved\.service$`,
		`^systemd-sysctl\.service$`,
		`^systemd-sysusers\.service$`,
		`^systemd-tmpfiles-setup-dev\.service$`,
		`^systemd-tmpfiles-setup\.service$`,
		`^systemd-udev-trigger\.service$`,
		`^systemd-udevd\.service$`,
		`^systemd-update-utmp\.service$`,
		`^systemd-user-sessions\.service$`,
		`^ufw\.service$`,
		`^unattended-upgrades\.service$`,
		`^user-runtime-dir@\d+\.service$`, // Matches user-runtime-dir@<any UID>.service
		`^user@\d+\.service$`,             // Matches user@<any UID>.service
	}

	// Compile the exclude service regexes and store them in a set
	excludeSvcRegexSet := make(map[*regexp.Regexp]struct{}, len(excludeSvcRegexList))
	for _, svcRegex := range excludeSvcRegexList {
		regexObj, err := regexp.Compile(svcRegex)
		if err != nil {
			log.Fatalf("Failed to compile regex: %v", err)
		}
		excludeSvcRegexSet[regexObj] = struct{}{}
	}

	var services []Service
	for _, service := range output {
		// Check for empty line
		if service != "" {
			// Split the service line into fields separated by whitespace delimeter
			fields := strings.Split(service, "|")
			if len(fields) >= 3 {
				name := fields[0]
				subState := fields[1]
				command := fields[2]

				// Check if the service name matches any of the exclude regexes
				excluded := false

				for regexObj := range excludeSvcRegexSet {
					if regexObj.MatchString(name) {
						excluded = true
						break
					}
				}

				// If the service is not excluded, add it to the list of services
				if !excluded {
					services = append(services, Service{Name: name, SubState: subState, Command: command})
				}
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
	// Connect to the VM through SSH
	client := connect_ssh(user, host, port, privateKeyPath)

	if client == nil {
		log.Fatalf("Failed to connect to the VM")
	} else {
		fmt.Println("Connected to the VM successfully!")
	}

	fmt.Println("Gathering system services ...")

	// Run systemctl command to get active running services and their commands
	// cmd := "systemctl --type=service --state=active --no-pager --no-legend"
	cmd := `systemctl list-units --no-pager -ql --type=service --state=running | awk '{printf "%s\0", $1}' | xargs -r0 -I{serviceName} bash -c 'name={serviceName}; sub=$(systemctl show -p SubState --value "$name"); cmd=$(systemctl cat "$name" 2>/dev/null | grep -i "ExecStart=" | awk -F= "{print \$2}" | sed "s/daemon on/daemon off/g" | sed "s/master_process on;//g"); echo "$name|$sub|$cmd"'`
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

func parse_exposed_ports(output []string) []Port {
	var ports []Port
	for _, port := range output {
		if port != "" {
			fields := strings.Fields(port)
			if len(fields) == 2 {
				protocol := fields[0]
				portNr := fields[1]
				ports = append(ports, Port{Protocol: protocol, PortNr: portNr})
			}
		}
	}
	return ports
}

func save_exposed_ports(ports []Port, path string) {
	portsData := PortsData{Ports: ports}
	yamFile, err := yaml.Marshal(portsData)
	if err != nil {
		log.Fatalf("Failed to move exposed ports to YAML format: %v", err)
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

func collect_exposed_ports(user string, host string, port string, privateKeyPath string) {
	// Connect to the VM through SSH
	client := connect_ssh(user, host, port, privateKeyPath)

	if client == nil {
		log.Fatalf("Failed to connect to the VM")
	} else {
		fmt.Println("Connected to the VM successfully!")
	}

	cmd := "sudo ss -tuln | grep LISTEN | grep -vE '(:22 )' | awk '!existing_values[$4]++' | awk -F ' ' '{print $1,$5}' | awk -F'[ :]' '{print $1, $3}'"

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

	// Parse the output by sending it as list of ports strings
	ports := parse_exposed_ports(strings.Split(string(output), "\n"))

	// Save the ports to a YAML file
	save_exposed_ports(ports, "exposed-ports.yaml")

	fmt.Println("Exposed ports collected successfully!")
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

	// Step 3: Collect exposed ports from source VM
	collect_exposed_ports(
		"antoniomihailov2001",
		"34.173.30.91",
		"22",
		"/home/toni/.ssh/id_ed25519_gcloud_source_vm",
	)
}
