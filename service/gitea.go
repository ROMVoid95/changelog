// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"time"

	"code.gitea.io/sdk/gitea"
)

// Gitea defines a Gitea service
type Gitea struct {
	Milestone string
	Token     string
	BaseURL   string
	Owner     string
	Repo      string
}

// Generate returns a Gitea changelog
func (ge *Gitea) Generate() (string, []PullRequest, error) {
	client := gitea.NewClient(ge.BaseURL, ge.Token)

	prs := make([]PullRequest, 0)

	milestoneID, err := ge.milestoneID(client)
	if err != nil {
		return "", nil, err
	}

	tagURL := fmt.Sprintf("## [%s](%s/%s/%s/pulls?q=&type=all&state=closed&milestone=%d) - %s", ge.Milestone, ge.BaseURL, ge.Owner, ge.Repo, milestoneID, time.Now().Format("2006-01-02"))

	p := 1
	// https://github.com/go-gitea/gitea/blob/d92781bf941972761177ac9e07441f8893758fd3/models/repo.go#L63
	// https://github.com/go-gitea/gitea/blob/e3c3b33ea7a5a223e22688c3f0eb2d3dab9f991c/models/pull_list.go#L104
	// FIXME Gitea has this hard-coded at 40
	perPage := 40
	for {
		results, err := client.ListRepoPullRequests(ge.Owner, ge.Repo, gitea.ListPullRequestsOptions{
			Page:      p,
			State:     "closed",
			Milestone: milestoneID,
		})
		if err != nil {
			return "", nil, err
		}
		p++

		for _, pr := range results {
			if pr != nil && pr.HasMerged {
				p := PullRequest{
					Title: pr.Title,
					Index: pr.Index,
				}

				labels := make([]Label, len(pr.Labels))
				for idx, lbl := range pr.Labels {
					labels[idx] = Label{
						Name: lbl.Name,
					}
				}
				p.Labels = labels

				prs = append(prs, p)
			}
		}

		if len(results) != perPage {
			break
		}
	}

	return tagURL, prs, nil
}

// Contributors returns a list of contributors from Gitea
func (ge *Gitea) Contributors() (ContributorList, error) {
	client := gitea.NewClient(ge.BaseURL, ge.Token)

	contributorsMap := make(map[string]bool)

	milestoneID, err := ge.milestoneID(client)
	if err != nil {
		return nil, err
	}

	p := 1
	perPage := 100
	for {
		results, err := client.ListRepoPullRequests(ge.Owner, ge.Repo, gitea.ListPullRequestsOptions{
			Page:      p,
			State:     "closed",
			Milestone: milestoneID,
		})
		if err != nil {
			return nil, err
		}
		p++

		for _, pr := range results {
			if pr != nil && pr.HasMerged {
				contributorsMap[pr.Poster.UserName] = true
			}
		}

		if len(results) != perPage {
			break
		}
	}

	contributors := make(ContributorList, 0, len(contributorsMap))
	for contributor := range contributorsMap {
		contributors = append(contributors, Contributor{
			Name:    contributor,
			Profile: fmt.Sprintf("%s/%s", ge.BaseURL, contributor),
		})
	}

	return contributors, nil
}

func (ge *Gitea) milestoneID(client *gitea.Client) (int64, error) {
	milestones, err := client.ListRepoMilestones(ge.Owner, ge.Repo, gitea.ListMilestoneOption{State: gitea.StateAll})
	if err != nil {
		return 0, err
	}

	for _, ms := range milestones {
		if ms.Title == ge.Milestone {
			return ms.ID, nil
		}
	}

	return 0, fmt.Errorf("no milestone found for %s", ge.Milestone)
}
