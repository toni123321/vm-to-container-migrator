package analyze

import "fmt"

type MixedAnalyzerImpl struct {
	applicationFileStrategy IAnalyzerFactory
	exposedPortsStrategy    IAnalyzerFactory
	servicesStrategy        IAnalyzerFactory
}

// Methods to allow chaning of strategies dynamically
// Strategy design pattern used

func (m *MixedAnalyzerImpl) SetApplicationFileStrategy(strategy IAnalyzerFactory) {
	m.applicationFileStrategy = strategy
}

func (m *MixedAnalyzerImpl) SetExposedPortsStrategy(strategy IAnalyzerFactory) {
	m.exposedPortsStrategy = strategy
}

func (m *MixedAnalyzerImpl) SetServicesStrategy(strategy IAnalyzerFactory) {
	m.servicesStrategy = strategy
}

func (m *MixedAnalyzerImpl) collectApplicationFiles(user string, host string, sourceDir string, destinationDir string, privateKeyPath string) (string, error) {
	if m.applicationFileStrategy == nil {
		return "", fmt.Errorf("ApplicationFileStrategy not set")
	}
	return m.applicationFileStrategy.collectApplicationFiles(user, host, sourceDir, destinationDir, privateKeyPath)
}

func (m *MixedAnalyzerImpl) collectExposedPorts(user string, host string, privateKeyPath string) (string, error) {
	if m.exposedPortsStrategy == nil {
		return "", fmt.Errorf("ExposedPortsStrategy not set")
	}
	return m.exposedPortsStrategy.collectExposedPorts(user, host, privateKeyPath)
}

func (m *MixedAnalyzerImpl) collectServices(user string, host string, privateKeyPath string) (string, error) {
	if m.servicesStrategy == nil {
		return "", fmt.Errorf("ServicesStrategy not set")
	}
	return m.servicesStrategy.collectServices(user, host, privateKeyPath)
}
