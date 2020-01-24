// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import "testing"

var gh = &GitHub{
	Milestone: "1.1.0", // https://github.com/go-gitea/test_repo/milestone/2?closed=1
	Repo:      "go-gitea/test_repo",
}

func TestGitHubGenerate(t *testing.T) {
	_, entries, err := gh.Generate()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(entries) != 1 {
		t.Logf("Expected 1 changelog entry, but got %d", len(entries))
		t.Fail()
	}
}

func TestGitHubContributors(t *testing.T) {
	contributors, err := gh.Contributors()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(contributors) != 1 {
		t.Logf("Expected 1 contributor, but got %d", len(contributors))
		t.Fail()
	}
}
