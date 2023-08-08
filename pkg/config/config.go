package config

import (
	"os"
	"strings"

	"dario.cat/mergo"
	"sigs.k8s.io/yaml"
)

var Cfg *Config

// Primarily used to define Features and disable them by default
// should facilities not configure them
// We may extend this later with other defaults
var defaultConfig = &Config{
	Features: ConfigFeatures{
		StaffingRequest: false,
	},
}

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

	// Merge in default, which will only "add" undefined values
	// We use this to set feature flags to disabled by default
	err = mergo.Merge(defaultConfig, cfg, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

func sanitizeConfig(cfg *Config) {
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if !strings.HasSuffix(cfg.Storage.BaseURL, "/") {
		cfg.Storage.BaseURL += "/"
	}
}
