package analyze

type Port struct {
	Protocol string `yaml:"protocol"`
	PortNr   string `yaml:"portNr"`
}

type Service struct {
	Name     string `yaml:"name"`
	SubState string `yaml:"substate"`
	Command  string `yaml:"command"`
}

type PortsData struct {
	Ports []Port `yaml:"ports"`
}

type ServicesData struct {
	Services []Service `yaml:"services"`
}
