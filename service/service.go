// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"strings"
)

const defaultGitea = "https://gitea.com"

// Load returns a service from a string
func New(serviceType, repo, baseURL, milestone, token string) (Service, error) {
	switch strings.ToLower(serviceType) {
	case "github":
		return &GitHub{
			Milestone: milestone,
			Token:     token,
			Repo:      repo,
		}, nil
	case "gitea":
		ownerRepo := strings.Split(repo, "/")
		if strings.TrimSpace(baseURL) == "" {
			baseURL = defaultGitea
		}
		return &Gitea{
			Milestone: milestone,
			Token:     token,
			BaseURL:   baseURL,
			Owner:     ownerRepo[0],
			Repo:      ownerRepo[1],
		}, nil
	default:
		return nil, fmt.Errorf("unknown service type %s", serviceType)
	}
}

// Service defines how a struct can be a Changelog Service
type Service interface {
	Generate() (string, []PullRequest, error)
	Contributors() (ContributorList, error)
}

// Label is the minimum information needed for a PR label
type Label struct {
	Name string
}

// PullRequest is the minimum information needed to make a changelog entry
type PullRequest struct {
	Title  string
	Index  int64
	Labels []Label
}

// Contributor is a project contributor
type Contributor struct {
	Name    string
	Profile string
}

// ContributorList is a slice of Contributors that can be sorted
type ContributorList []Contributor

// Len is the length of the ContributorList
func (cl ContributorList) Len() int {
	return len(cl)
}

// Less determines whether a Contributor comes before another Contributor
func (cl ContributorList) Less(i, j int) bool {
	return cl[i].Name < cl[j].Name
}

// Swap swaps Contributors in a ContributorList
func (cl ContributorList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}
