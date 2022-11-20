package config

import (
	"os"
	"strings"

	"sigs.k8s.io/yaml"
)

var Cfg *Config

func ParseConfig(file string) (*Config, error) {
	// Read config file disk
	config, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// Parse config file
	cfg := &Config{}
	err = yaml.Unmarshal(config, cfg)
	if err != nil {
		return nil, err
	}

	sanitizeConfig(cfg)

	return cfg, nil
}

func sanitizeConfig(cfg *Config) {
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if !strings.HasSuffix(cfg.Storage.BaseURL, "/") {
		cfg.Storage.BaseURL += "/"
	}
}
