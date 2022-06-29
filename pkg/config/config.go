package config

import (
	"os"

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

	return cfg, nil
}
