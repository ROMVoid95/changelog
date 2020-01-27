// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

func init() {
	DefaultConfig = []byte(`# The full repository name
repo: go-gitea/gitea

# Service type (gitea or github)
service: github

# Base URL for Gitea instance if using gitea service type (optional)
# Default: https://gitea.com
base-url:

# Changelog groups and which labeled PRs to add to each group
groups:
  - 
    name: BREAKING
    labels:
      - kind/breaking
  - 
    name: FEATURES
    labels:
      - kind/feature
  -
    name: BUGFIXES
    labels:
      - kind/bug
  - 
    name: ENHANCEMENTS
    labels:
      - kind/enhancement
      - kind/refactor
      - kind/ui
  -
    name: SECURITY
    labels:
      - kind/security
  - 
    name: TESTING
    labels:
      - kind/testing
  - 
    name: TRANSLATION
    labels:
      - kind/translation
  - 
    name: BUILD
    labels:
      - kind/build
      - kind/lint
  - 
    name: DOCS
    labels:
    - kind/docs
  - 
    name: MISC
    default: true

# regex indicating which labels to skip for the changelog
skip-labels: skip-changelog|backport\/.+`)
}
