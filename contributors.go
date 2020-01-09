// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/google/go-github/github"
	"github.com/urfave/cli/v2"
)

var cmdContributors = &cli.Command{
	Name:        "contributors",
	Usage:       "generate contributors list",
	Description: "generate contributors list",
	Action:      runContributors,
}

func runContributors(cmd *cli.Context) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	client := github.NewClient(nil)
	ctx := context.Background()

	contributorsMap := make(map[string]bool)
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

		for _, pr := range result.Issues {
			contributorsMap[*pr.User.Login] = true
		}

		if len(result.Issues) != perPage {
			break
		}
	}

	contributors := make([]string, 0, len(contributorsMap))
	for contributor, _ := range contributorsMap {
		contributors = append(contributors, contributor)
	}

	sort.Strings(contributors)

	for _, contributor := range contributors {
		fmt.Printf("* [@%s](https://github.com/%s)\n", contributor, contributor)
	}

	return nil
}
