/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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
		StaffingRequest:  false,
		ControllerOnline: true,
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
