package analyze

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"vm2cont/api/internal/model"
	"vm2cont/api/internal/utils"

	"gopkg.in/yaml.v3"
)

const BASE_ANALYZE_OUTPUT_DIR = "./output/analyze-output/"

// Implementation of IAnalyzerFactory interface
type FsAnalyzerImpl struct{}

func (f *FsAnalyzerImpl) collectApplicationFiles(user string, host string, sourceDir string, destinationDir string, privateKeyPath string) (string, error) {
	// source includes the user and host of the VM, plus the source directory
	var source string = fmt.Sprintf("%s@%s:%s", user, host, sourceDir)
	// destination includes the destination directory, concatenated with the base output directory
	var destination string = BASE_ANALYZE_OUTPUT_DIR + destinationDir

	excludes := []string{
		"--exclude=/proc/*",
		"--exclude=/boot/*",
		"--exclude=/sys/*",
		"--exclude=/dev/*",
		"--exclude=/snap/*",
		"--exclude=/etc/apparmor/*",
		"--exclude=/etc/apparmor.d/*",
		"--exclude=/etc/apport/*",
		"--exclude=/etc/apt/*",
		"--exclude=/etc/chrony/*",
		"--exclude=/etc/cloud/*",
		"--exclude=/etc/cron.d/*",
		"--exclude=/etc/cron.daily/*",
		"--exclude=/etc/cron.hourly/*",
		"--exclude=/etc/cron.monthly/*",
		"--exclude=/etc/cron.weekly/*",
		"--exclude=/etc/dbus-1/*",
		"--exclude=/etc/dpkg/*",
		"--exclude=/lib/modules/*",
		"--exclude=/usr/share/man/*",
		"--exclude=/usr/share/doc/*",
		"--exclude=/usr/lib/snapd/*",
		"--exclude=/usr/lib/systemd/*",
		"--exclude=/usr/lib/apt/*",
		"--exclude=/usr/lib/dpkg/*",
		"--exclude=/usr/lib/apparmor/*",
		"--exclude=/usr/lib/cloud-init/*",
		"--exclude=/usr/lib/google-cloud-sdk/*",
		"--exclude=/var/cache/*",
		"--exclude=/var/backups/*",
		"--exclude=/var/tmp/*",
		"--exclude=/var/run/*",
		"--exclude=/var/log/*/*",
		"--exclude=/var/lib/lxcfs/*",
		"--exclude=/var/lib/snapd/*",
		"--exclude=/var/lib/systemd/*",
		"--exclude=/var/snap/*",
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

	fmt.Println("Gathering file system ...")

	cmd := exec.Command("rsync", args...)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running rsync: %v\n", err)
		return "", err
	}

	fmt.Println("File system collected successfully!")

	result := fmt.Sprintf("The Application files were collected through file system analysis\nPath to file system: %s", destination)
	return result, nil
}

func (f *FsAnalyzerImpl) collectExposedPorts(user string, host string, privateKeyPath string) (string, error) {
	// Connect to the VM through SSH
	client := utils.ConnectToSsh(user, host, "22", privateKeyPath)

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
	ports := parseExposedPorts(strings.Split(string(output), "\n"))

	// Save the ports to a YAML file
	path, err := saveExposedPorts(ports, "exposed-ports.yaml")

	if err != nil {
		log.Fatalf("Failed to save exposed ports: %v", err)
	}

	fmt.Println("Exposed ports collected successfully!")

	result := fmt.Sprintf("The Exposed ports were collected through file system analysis\nPath to exposed ports: %s", path)
	return result, nil
}

func (f *FsAnalyzerImpl) collectServices(user string, host string, privateKeyPath string) (string, error) {
	// Connect to the VM through SSH
	client := utils.ConnectToSsh(user, host, "22", privateKeyPath)

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
	services := parseSysServices(strings.Split(string(output), "\n"))

	// Save the services to a YAML file
	serviceFilePath, err := saveSysServices(services, "sys-services.yaml")

	if err != nil {
		log.Fatalf("Failed to save system services: %v", err)
	}

	fmt.Println("System services collected successfully!")

	result := fmt.Sprintf("The System services were collected through file system analysis\nPath to system services: %s", serviceFilePath)
	return result, nil
}

func parseSysServices(output []string) []model.Service {
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

	var services []model.Service
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
					services = append(services, model.Service{Name: name, SubState: subState, Command: command})
				}
			}
		}
	}
	return services
}

func saveSysServices(services []model.Service, path string) (string, error) {
	serviceData := model.ServicesData{Services: services}
	yamFile, err := yaml.Marshal(serviceData)
	if err != nil {
		log.Fatalf("Failed to move services to YAML format: %v", err)
	}

	// concatenate the base output directory with the path
	path = BASE_ANALYZE_OUTPUT_DIR + path
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(yamFile)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	return path, nil
}

func parseExposedPorts(output []string) []model.Port {
	var ports []model.Port
	for _, port := range output {
		if port != "" {
			fields := strings.Fields(port)
			if len(fields) == 2 {
				protocol := fields[0]
				portNr := fields[1]
				ports = append(ports, model.Port{Protocol: protocol, PortNr: portNr})
			}
		}
	}
	return ports
}

func saveExposedPorts(ports []model.Port, path string) (string, error) {
	portsData := model.PortsData{Ports: ports}
	yamFile, err := yaml.Marshal(portsData)
	if err != nil {
		log.Fatalf("Failed to move exposed ports to YAML format: %v", err)
	}

	// concatenate the base output directory with the path
	path = BASE_ANALYZE_OUTPUT_DIR + path
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(yamFile)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	return path, nil
}
