// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/github"
	"github.com/urfave/cli/v2"
)

var cmdGenerate = &cli.Command{
	Name:        "generate",
	Usage:       "generate changelog",
	Description: "generate changelog",
	Action:      runGenerate,
}

func runGenerate(cmd *cli.Context) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	client := github.NewClient(nil)
	ctx := context.Background()

	labels := make(map[string]string)
	changelogs := make(map[string][]github.Issue)
	var defaultGroup string
	for _, g := range config.Groups {
		changelogs[g.Name] = []github.Issue{}
		for _, l := range g.Labels {
			labels[l] = g.Name
		}
		if g.Default {
			defaultGroup = g.Name
		}
	}

	if defaultGroup == "" {
		defaultGroup = config.Groups[len(config.Groups)-1].Name
	}

	query := fmt.Sprintf(`repo:%s is:merged milestone:"%s"`, config.Repo, milestone)
	p := 1
	perPage := 100
	for {
		result, _, err := client.Search.Issues(ctx, query, &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    p,
				PerPage: perPage,
			},
		})
		p++
		if err != nil {
			log.Fatal(err.Error())
		}

	PRLoop: // labels in Go, let's get old school
		for _, pr := range result.Issues {
			var label string
			for _, lb := range pr.Labels {
				if config.SkipRegex != nil && config.SkipRegex.MatchString(lb.GetName()) {
					continue PRLoop
				}

				if g, ok := labels[lb.GetName()]; ok && len(label) == 0 {
					label = g
				}
			}

			if len(label) > 0 {
				changelogs[label] = append(changelogs[label], pr)
			} else {
				changelogs[defaultGroup] = append(changelogs[defaultGroup], pr)
			}
		}

		if len(result.Issues) != perPage {
			break
		}
	}

	fmt.Printf("## [%s](https://github.com/%s/releases/tag/v%s) - %s\n", milestone, config.Repo, milestone, time.Now().Format("2006-01-02"))
	for _, g := range config.Groups {
		if len(changelogs[g.Name]) == 0 {
			continue
		}

		fmt.Println("* " + g.Name)
		for _, pr := range changelogs[g.Name] {
			fmt.Printf("  * %s (#%d)\n", *pr.Title, *pr.Number)
		}
	}

	return nil
}
