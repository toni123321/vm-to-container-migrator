package analyze

// Implementation of ProcessAnalyzer interface, and CommonAnalyzer interface

type ProcessAnalyzerImpl struct{}

func (p *ProcessAnalyzerImpl) collectApplicationFiles(user string, host string, sourceDir string, destinationDir string, privateKeyPath string) (string, error) {
	return "", nil
}

func (p *ProcessAnalyzerImpl) collectExposedPorts(user string, host string, privateKeyPath string) (string, error) {
	return "", nil
}

func (p *ProcessAnalyzerImpl) collectServices(user string, host string, privateKeyPath string) (string, error) {
	return "", nil
}
