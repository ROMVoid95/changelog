// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"sort"

	"code.gitea.io/changelog/config"
	"code.gitea.io/changelog/service"

	"github.com/urfave/cli/v2"
)

var Contributors = &cli.Command{
	Name:   "contributors",
	Usage:  "Generates a contributors list",
	Action: runContributors,
}

func runContributors(cmd *cli.Context) error {

	if ConfigPathFlag == "" {
		ConfigPathFlag = getDefaultConfigFile()
	}

	cfg, err := config.New(ConfigPathFlag)
	if err != nil {
		return err
	}

	s, err := service.New(cfg.Service, cfg.Repo, cfg.BaseURL, MilestoneFlag, TokenFlag, IssuesFlag)
	if err != nil {
		return err
	}

	contributors, err := s.Contributors()
	if err != nil {
		return err
	}

	sort.Sort(contributors)

	for _, contributor := range contributors {
		fmt.Printf("* [@%s](%s)\n", contributor.Name, contributor.Profile)
	}

	return nil
}
