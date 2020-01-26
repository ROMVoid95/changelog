// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

// GitHub defines a GitHub service
type GitHub struct {
	Milestone string
	Token     string
	Repo      string
}

// Generate returns a GitHub changelog
func (gh *GitHub) Generate() (string, []PullRequest, error) {
	tagURL := fmt.Sprintf("## [%s](https://github.com/%s/releases/tag/v%s) - %s", gh.Milestone, gh.Repo, gh.Milestone, time.Now().Format("2006-01-02"))

	client := github.NewClient(nil)
	ctx := context.Background()

	prs := make([]PullRequest, 0)

	query := fmt.Sprintf(`repo:%s is:merged milestone:"%s"`, gh.Repo, gh.Milestone)
	p := 1
	perPage := 100
	for {
		result, _, err := client.Search.Issues(ctx, query, &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    p,
				PerPage: perPage,
			},
		})
		if err != nil {
			return "", nil, err
		}
		p++

		for _, pr := range result.Issues {
			if pr.IsPullRequest() {
				p := PullRequest{
					Title: pr.GetTitle(),
					Index: int64(pr.GetNumber()),
				}

				labels := make([]Label, len(pr.Labels))
				for idx, lbl := range pr.Labels {
					labels[idx] = Label{
						Name: lbl.GetName(),
					}
				}
				p.Labels = labels

				prs = append(prs, p)
			}
		}

		if len(result.Issues) != perPage {
			break
		}
	}

	return tagURL, prs, nil
}

// Contributors returns a list of contributors from GitHub
func (gh *GitHub) Contributors() (ContributorList, error) {
	client := github.NewClient(nil)
	ctx := context.Background()

	contributorsMap := make(map[string]bool)
	query := fmt.Sprintf(`repo:%s is:merged milestone:"%s"`, gh.Repo, gh.Milestone)
	p := 1
	perPage := 100
	for {
		result, _, err := client.Search.Issues(ctx, query, &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    p,
				PerPage: perPage,
			},
		})
		if err != nil {
			return nil, err
		}
		p++

		for _, pr := range result.Issues {
			contributorsMap[pr.GetUser().GetLogin()] = true
		}

		if len(result.Issues) != perPage {
			break
		}
	}

	contributors := make(ContributorList, 0, len(contributorsMap))
	for contributor := range contributorsMap {
		contributors = append(contributors, Contributor{
			Name:    contributor,
			Profile: fmt.Sprintf("https://github.com/%s", contributor),
		})
	}

	return contributors, nil
}
