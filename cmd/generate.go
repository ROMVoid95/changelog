// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"code.gitea.io/changelog/config"
	"code.gitea.io/changelog/service"

	"github.com/urfave/cli/v2"
)

var Generate = &cli.Command{
	Name:   "generate",
	Usage:  "Generates a changelog",
	Action: runGenerate,
}

func runGenerate(cmd *cli.Context) error {
	cfg, err := config.New(ConfigPathFlag)
	if err != nil {
		return err
	}

	labels := make(map[string]string)
	entries := make(map[string][]service.PullRequest)
	var defaultGroup string
	for _, g := range cfg.Groups {
		entries[g.Name] = []service.PullRequest{}
		for _, l := range g.Labels {
			labels[l] = g.Name
		}
		if g.Default {
			defaultGroup = g.Name
		}
	}

	if defaultGroup == "" {
		defaultGroup = cfg.Groups[len(cfg.Groups)-1].Name
	}

	s, err := service.New(cfg.Service, cfg.Repo, cfg.BaseURL, MilestoneFlag, TokenFlag)
	if err != nil {
		return err
	}

	title, prs, err := s.Generate()
	if err != nil {
		return err
	}

PRLoop: // labels in Go, let's get old school
	for _, pr := range prs {
		if pr.Index < AfterFlag {
			continue
		}

		var label string
		for _, lb := range pr.Labels {
			if cfg.SkipRegex != nil && cfg.SkipRegex.MatchString(lb.Name) {
				continue PRLoop
			}

			if g, ok := labels[lb.Name]; ok && len(label) == 0 {
				label = g
			}
		}

		if len(label) > 0 {
			entries[label] = append(entries[label], pr)
		} else {
			entries[defaultGroup] = append(entries[defaultGroup], pr)
		}
	}

	fmt.Println(title)
	for _, g := range cfg.Groups {
		if len(entries[g.Name]) == 0 {
			continue
		}

		if DetailsFlag {
			fmt.Println("<details><summary>" + g.Name + "</summary>")
			fmt.Println()
			for _, entry := range entries[g.Name] {
				fmt.Printf("* %s (#%d)\n", entry.Title, entry.Index)
			}
			fmt.Println("</details>")
		} else {
			fmt.Println("* " + g.Name)
			for _, entry := range entries[g.Name] {
				fmt.Printf("  * %s (#%d)\n", entry.Title, entry.Index)
			}
		}
	}

	return nil
}
