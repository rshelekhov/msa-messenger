package config

type CORS struct {
	AllowedOrigins   []string `yaml:"AllowedOrigins"`
	AllowedMethods   []string `yaml:"AllowedMethods"`
	AllowedHeaders   []string `yaml:"AllowedHeaders"`
	ExposedHeaders   []string `yaml:"ExposedHeaders"`
	AllowCredentials bool     `yaml:"AllowCredentials"`
	MaxAge           int      `yaml:"MaxAge"`
	Debug            bool     `yaml:"Debug"`
}
