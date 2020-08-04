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
	Issues    bool
}

// Generate returns a Gitea changelog
func (ge *Gitea) Generate() (string, []Entry, error) {
	client := gitea.NewClient(ge.BaseURL, ge.Token)

	entries := make([]Entry, 0)

	milestoneID, err := ge.milestoneID(client)
	if err != nil {
		return "", nil, err
	}

	from := "pulls"
	if ge.Issues {
		from = "issues"
	}

	tagURL := fmt.Sprintf("## [%s](%s/%s/%s/%s?q=&type=all&state=closed&milestone=%d) - %s", ge.Milestone, ge.BaseURL, ge.Owner, ge.Repo, from, milestoneID, time.Now().Format("2006-01-02"))

	p := 1
	// https://github.com/go-gitea/gitea/blob/d92781bf941972761177ac9e07441f8893758fd3/models/repo.go#L63
	// https://github.com/go-gitea/gitea/blob/e3c3b33ea7a5a223e22688c3f0eb2d3dab9f991c/models/pull_list.go#L104
	// FIXME Gitea has this hard-coded at 40
	perPage := 40
	for {
		options := gitea.ListIssueOption{
			ListOptions: gitea.ListOptions{
				Page:     p,
				PageSize: perPage,
			},
			Milestones: []string{ge.Milestone},
			State:      gitea.StateClosed,
			Type:       gitea.IssueTypePull,
		}
		if ge.Issues {
			options.Type = gitea.IssueTypeIssue
		}

		issues, err := client.ListRepoIssues(ge.Owner, ge.Repo, options)
		if err != nil {
			return "", nil, err
		}

		for _, issue := range issues {
			if issue != nil {
				if options.Type == gitea.IssueTypePull && issue.PullRequest != nil && !(issue.PullRequest.HasMerged) {
					continue
				}

				entry := convertToEntry(*issue)
				entries = append(entries, entry)
			}
		}

		if len(issues) != perPage {
			break
		}

		p++
	}

	return tagURL, entries, nil
}

func convertToEntry(issue gitea.Issue) Entry {
	entry := Entry{
		Index: issue.Index,
		Title: issue.Title,
	}

	labels := make([]Label, len(issue.Labels))
	for idx, lbl := range issue.Labels {
		labels[idx] = Label{
			Name: lbl.Name,
		}
	}

	entry.Labels = labels

	return entry
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
			ListOptions: gitea.ListOptions{
				Page:     p,
				PageSize: perPage,
			},
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
