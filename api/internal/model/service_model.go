package model

type Service struct {
	Name     string `yaml:"name"`
	SubState string `yaml:"substate"`
	Command  string `yaml:"command"`
}

type ServicesData struct {
	Services []Service `yaml:"services"`
}
