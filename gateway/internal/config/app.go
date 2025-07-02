package config

type App struct {
	ID             string `yaml:"ID"`
	Env            string `yaml:"Env"`
	ServiceName    string `yaml:"ServiceName"`
	ServiceVersion string `yaml:"ServiceVersion"`
	EnableMetrics  bool   `yaml:"EnableMetrics"`
	OTLPEndpoint   string `yaml:"OTLPEndpoint"`
}
