// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	exampleFile = "changelog.example.yml"
	writeFile   = "config/config_default.go"
	tmpl        = `// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

func init() {
	DefaultConfig = []byte(` + "`" + `%s` + "`" + `)
}
`
)

func main() {
	bytes, err := ioutil.ReadFile(exampleFile)
	if err != nil {
		fmt.Printf("Could not read from %s. Are you in the root directory of the project?", exampleFile)
		os.Exit(1)
	}

	data := fmt.Sprintf(tmpl, string(bytes))

	if err := ioutil.WriteFile(writeFile, []byte(data), os.ModePerm); err != nil {
		fmt.Printf("Could not write to %s.", writeFile)
		os.Exit(1)
	}
}
