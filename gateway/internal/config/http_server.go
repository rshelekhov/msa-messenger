package config

import "time"

type HTTPServer struct {
	Address     string        `mapstructure:"HTTP_SERVER_ADDRESS" envDefault:"localhost:8080"`
	Timeout     time.Duration `mapstructure:"HTTP_SERVER_TIMEOUT" envDefault:"10s"`
	IdleTimeout time.Duration `mapstructure:"HTTP_SERVER_IDLE_TIMEOUT" envDefault:"60s"`
}
