package model

type Port struct {
	Protocol string `yaml:"protocol"`
	PortNr   string `yaml:"portNr"`
}

type PortsData struct {
	Ports []Port `yaml:"ports"`
}
