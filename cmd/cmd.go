// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"os"
	"path/filepath"
)

var (
	MilestoneFlag  string
	ConfigPathFlag string
	TokenFlag      string
	DetailsFlag    bool
	AfterFlag      int64
)

func getDefaultConfigFile() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	config := filepath.Join(pwd, ".changelog.yml")
	info, err := os.Stat(config)
	if err == nil && !info.IsDir() {
		return config
	}
	return ""
}
