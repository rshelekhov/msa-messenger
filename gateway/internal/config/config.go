package config

type Config struct {
	App        App        `mapstructure:",squash"`
	HTTPServer HTTPServer `mapstructure:",squash"`
	CORS       CORS       `mapstructure:",squash"`
	Validator  Validator  `mapstructure:",squash"`
	JWT        JWT        `mapstructure:",squash"`
}
