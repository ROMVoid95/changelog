// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

//go:generate go run changelog.example.go
//go:generate go fmt ./...

import (
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

var defaultConfig []byte

type Config struct {
	Repo   string `yaml:"repo"`
	Groups []struct {
		Name    string   `yaml:"name"`
		Labels  []string `yaml:"labels"`
		Default bool     `yaml:"default"`
	} `yaml:"groups"`
	SkipLabels string         `yaml:"skip-labels"`
	SkipRegex  *regexp.Regexp `yaml:"-"`
}

func LoadConfig() (*Config, error) {
	var err error
	var configContent []byte
	if len(configPath) == 0 {
		configContent = defaultConfig
	} else {
		configContent, err = ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
	}

	var config *Config
	if err = yaml.Unmarshal(configContent, &config); err != nil {
		return nil, err
	}

	if len(config.SkipLabels) > 0 {
		if config.SkipRegex, err = regexp.Compile(config.SkipLabels); err != nil {
			return nil, err
		}
	}

	return config, nil
}
