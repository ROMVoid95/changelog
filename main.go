// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

//go:generate go run changelog.example.go
//go:generate go fmt ./...

import (
	"fmt"
	"os"

	"code.gitea.io/changelog/cmd"

	"github.com/urfave/cli/v2"
)

var (
	// Version of changelog
	Version = "development"
)

func main() {
	app := &cli.App{
		Name:    "changelog",
		Usage:   "Changelog tools for Gitea",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "milestone",
				Aliases:     []string{"m"},
				Usage:       "Targeted milestone",
				Destination: &cmd.MilestoneFlag,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Specify a config file",
				Destination: &cmd.ConfigPathFlag,
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Usage:       "Access token for private repositories/instances",
				Destination: &cmd.TokenFlag,
			},
			&cli.BoolFlag{
				Name:        "details",
				Aliases:     []string{"d"},
				Usage:       "Generate detail lists instead of long lists",
				Destination: &cmd.DetailsFlag,
			},
			&cli.Int64Flag{
				Name:        "after",
				Aliases:     []string{"a"},
				Usage:       "Only select PRs after a given index (continuing a previous changelog)",
				Destination: &cmd.AfterFlag,
			},
			&cli.BoolFlag{
				Name:        "issues",
				Aliases:     []string{"i"},
				Usage:       "Generate changelog from issues (otherwise from pulls)",
				Destination: &cmd.IssuesFlag,
			},
		},
		Commands: []*cli.Command{
			cmd.Generate,
			cmd.Contributors,
			cmd.Init,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Failed to run app with %s: %v\n", os.Args[1:], err)
	}
}
