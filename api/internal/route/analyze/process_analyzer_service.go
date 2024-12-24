package analyze

import (
	"fmt"
)

// Implementation of IAnalyzerFactory interface
type ProcessAnalyzerImpl struct{}

func (p *ProcessAnalyzerImpl) collectApplicationFiles(user string, host string, sourceDir string, destinationDir string, privateKeyPath string) (string, error) {
	// return "Application files collected through process analysis", nil

	// Simulate the process of collecting application files through process analysis
	// destination includes the destination directory, concatenated with the base output directory
	var destination string = BASE_ANALYZE_OUTPUT_DIR + destinationDir

	fmt.Println("File system collected successfully!")

	result := fmt.Sprintf("The Application files were collected through file system analysis\nPath to file system: %s", destination)
	return result, nil
}

func (p *ProcessAnalyzerImpl) collectExposedPorts(user string, host string, privateKeyPath string) (string, error) {
	// Simulate the process of collecting exposed ports through process analysis
	path := BASE_ANALYZE_OUTPUT_DIR + "exposed-ports.yaml"

	fmt.Println("Exposed ports collected successfully!")

	result := fmt.Sprintf("The Exposed ports were collected through file system analysis\nPath to exposed ports: %s", path)
	return result, nil
}

func (p *ProcessAnalyzerImpl) collectServices(user string, host string, privateKeyPath string) (string, error) {
	// Simulate the process of collecting services through process analysis
	path := BASE_ANALYZE_OUTPUT_DIR + "services.yaml"

	fmt.Println("Services collected successfully!")

	result := fmt.Sprintf("The System services were collected through file system analysis\nPath to system services: %s", path)
	return result, nil
}
