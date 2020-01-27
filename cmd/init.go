// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"code.gitea.io/changelog/config"

	"github.com/urfave/cli/v2"
)

var (
	Init = &cli.Command{
		Name:  "init",
		Usage: "Initialize a default .changelog.yml",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Usage:       "Name of the changelog config",
				Value:       ".changelog.yml",
				Destination: &nameFlag,
			},
		},
		Action: runInit,
	}
	nameFlag string
)

func runInit(cmd *cli.Context) error {
	if _, err := os.Stat(nameFlag); err == nil {
		return fmt.Errorf("file '%s' already exists", nameFlag)
	}

	if err := ioutil.WriteFile(nameFlag, config.DefaultConfig, os.ModePerm); err != nil {
		return err
	}

	fmt.Printf("Config initialized at '%s'\n", nameFlag)
	return nil
}
