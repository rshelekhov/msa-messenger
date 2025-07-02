package config

type Config struct {
	App        App        `yaml:"App"`
	HTTPServer HTTPServer `yaml:"HTTPServer"`
	CORS       CORS       `yaml:"CORS"`
	Validator  Validator  `yaml:"Validator"`
	JWT        JWT        `yaml:"JWT"`
}
