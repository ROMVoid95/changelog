// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

var DefaultConfig []byte

// Group is a grouping of PRs
type Group struct {
	Name    string   `yaml:"name"`
	Labels  []string `yaml:"labels"`
	Default bool     `yaml:"default"`
}

// Config is the changelog settings
type Config struct {
	Repo       string         `yaml:"repo"`
	Service    string         `yaml:"service"`
	BaseURL    string         `yaml:"base-url"`
	Groups     []Group        `yaml:"groups"`
	SkipLabels string         `yaml:"skip-labels"`
	SkipRegex  *regexp.Regexp `yaml:"-"`
}

// Load a config from a path, defaulting to changelog.example.yml
func New(configPath string) (*Config, error) {
	var err error
	var configContent []byte
	if len(configPath) == 0 {
		configContent = DefaultConfig
	} else {
		configContent, err = ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
	}

	var cfg *Config
	if err = yaml.Unmarshal(configContent, &cfg); err != nil {
		return nil, err
	}

	if len(cfg.SkipLabels) > 0 {
		if cfg.SkipRegex, err = regexp.Compile(cfg.SkipLabels); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
