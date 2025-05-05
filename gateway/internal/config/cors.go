package config

type CORS struct {
	AllowedOrigins   []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods   []string `mapstructure:"CORS_ALLOWED_METHODS"`
	AllowedHeaders   []string `mapstructure:"CORS_ALLOWED_HEADERS"`
	ExposedHeaders   []string `mapstructure:"CORS_EXPOSED_HEADERS"`
	AllowCredentials bool     `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	MaxAge           int      `mapstructure:"CORS_MAX_AGE"`
	Debug            bool     `mapstructure:"CORS_DEBUG"`
}
