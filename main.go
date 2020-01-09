// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	// Version of changelog
	Version = "0.2"
)

var (
	milestone  string
	configPath string
)

func main() {
	app := &cli.App{
		Name:    "changelog",
		Usage:   "Changelog generator for Gitea",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "milestone",
				Aliases:     []string{"m"},
				Usage:       "Targeted milestone",
				Required:    true,
				Destination: &milestone,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Specify a config file",
				Destination: &configPath,
			},
		},
		Commands: []*cli.Command{
			cmdGenerate,
			cmdContributors,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Failed to run app with %s: %v\n", os.Args[1:], err)
	}
}
