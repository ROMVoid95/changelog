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
	client, err := gitea.NewClient(ge.BaseURL, gitea.SetToken(ge.Token))
	if err != nil {
		return "", nil, err
	}

	entries := make([]Entry, 0)

	milestone, _, err := client.GetMilestoneByName(ge.Owner, ge.Repo, ge.Milestone)
	if err != nil {
		return "", nil, err
	}

	from := "pulls"
	if ge.Issues {
		from = "issues"
	}

	tagURL := getGiteaTagURL(client, ge.BaseURL, ge.Owner, ge.Repo, ge.Milestone, from, milestone.ID)

	perPage := ge.perPage(client)
	for p := 1; ; p++ {
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

		issues, _, err := client.ListRepoIssues(ge.Owner, ge.Repo, options)
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
	}

	return tagURL, entries, nil
}

func getGiteaTagURL(c *gitea.Client, baseURL, owner, repo, mileName, from string, mileID int64) string {
	if err := c.CheckServerVersionConstraint(">=1.12"); err != nil {
		return fmt.Sprintf("## [%s](%s/%s/%s/%s?q=&type=all&state=closed&milestone=%d) - %s", mileName, baseURL, owner, repo, from, mileID, time.Now().Format("2006-01-02"))
	}
	return fmt.Sprintf("## [%s](%s/%s/%s/releases/tag/%s) - %s", mileName, baseURL, owner, repo, mileName, time.Now().Format("2006-01-02"))
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
	client, err := gitea.NewClient(ge.BaseURL, gitea.SetToken(ge.Token))
	if err != nil {
		return nil, err
	}

	contributorsMap := make(map[string]bool)

	milestone, _, err := client.GetMilestoneByName(ge.Owner, ge.Repo, ge.Milestone)
	if err != nil {
		return nil, err
	}

	perPage := ge.perPage(client)
	for p := 1; ; p++ {
		results, _, err := client.ListRepoPullRequests(ge.Owner, ge.Repo, gitea.ListPullRequestsOptions{
			ListOptions: gitea.ListOptions{
				Page:     p,
				PageSize: perPage,
			},
			State:     "closed",
			Milestone: milestone.ID,
		})
		if err != nil {
			return nil, err
		}

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

func (ge *Gitea) perPage(client *gitea.Client) int {
	// set low value so it will work in most cases
	perPage := 10
	if client.CheckServerVersionConstraint(">=1.13.0") == nil {
		conf, _, err := client.GetGlobalAPISettings()
		if err != nil {
			return perPage
		}
		return conf.MaxResponseItems
	}
	return perPage
}
