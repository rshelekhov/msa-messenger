package config

import "time"

type HTTPServer struct {
	Address     string        `yaml:"Address" default:"localhost:8080"`
	Timeout     time.Duration `yaml:"Timeout" default:"10s"`
	IdleTimeout time.Duration `yaml:"IdleTimeout" default:"60s"`
}
