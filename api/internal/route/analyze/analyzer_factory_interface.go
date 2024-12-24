package analyze

import (
	"fmt"
)

// Common interface for analyzing
// Abstract factory design pattern used
type IAnalyzerFactory interface {
	collectApplicationFiles(user string, host string, sourceDir string, destinationDir string, privateKeyPath string) (string, error)
	collectExposedPorts(user string, host string, privateKeyPath string) (string, error)
	collectServices(user string, host string, privateKeyPath string) (string, error)
}

// Abstract Factory Implementation
func GetAnalyzerFactory(analyzerType string) (IAnalyzerFactory, error) {
	switch analyzerType {
	case "fs":
		return &FsAnalyzerImpl{}, nil
	case "process":
		return &ProcessAnalyzerImpl{}, nil
	case "mixed":
		return &MixedAnalyzerImpl{}, nil
	default:
		return nil, fmt.Errorf("invalid analyzer type")
	}
}
