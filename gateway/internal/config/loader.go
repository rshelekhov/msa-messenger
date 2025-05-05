package config

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const CONFIG_PATH = "CONFIG_PATH"

func MustLoad() *Config {
	cfg := &Config{}
	configPath := fetchConfigPath()

	// Initialize Viper
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// If a path to the config is specified, load from the file
	// In Kubernetes, all settings come through environment variables
	// from ConfigMap and Secret, so config files are only needed for local development
	if configPath != "" {
		// Determine the file format based on the extension
		ext := strings.TrimPrefix(strings.ToLower(configPath[strings.LastIndex(configPath, "."):]), ".")
		v.SetConfigType(ext)
		v.SetConfigFile(configPath)

		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("warning: error reading config file: %s", err)
		}
	}

	// Fill the structure
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	return cfg
}

func fetchConfigPath() string {
	var v string
	flag.StringVar(&v, "config", "", "path to config file")
	flag.Parse()

	if v == "" {
		v = os.Getenv(CONFIG_PATH)
	}

	return v
}
