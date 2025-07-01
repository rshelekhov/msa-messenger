package config

type App struct {
	ID             string `mapstructure:"APP_ID"`
	Env            string `mapstructure:"APP_ENV"`
	ServiceName    string `mapstructure:"APP_SERVICE_NAME"`
	ServiceVersion string `mapstructure:"APP_SERVICE_VERSION"`
	EnableMetrics  bool   `mapstructure:"APP_ENABLE_METRICS"`
	OTLPEndpoint   string `mapstructure:"APP_OTLP_ENDPOINT"`
}
