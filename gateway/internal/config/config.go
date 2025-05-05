package config

type Config struct {
	AppEnv     string     `mapstructure:"APP_ENV"`
	AppID      string     `mapstructure:"APP_ID"`
	HTTPServer HTTPServer `mapstructure:",squash"`
	CORS       CORS       `mapstructure:",squash"`
	Validator  Validator  `mapstructure:",squash"`
	JWT        JWT        `mapstructure:",squash"`
}
