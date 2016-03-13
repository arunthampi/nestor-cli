// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Power for Nestor",
	Run:   runNew,
}

func init() {
	RootCmd.AddCommand(newCmd)
}

var defaultAppFile string = `
module.exports = function(robot) {
  robot.respond(/hello/, function(msg, done) {
    msg.reply("hello back", done);
  });
};
`
var readmeUrl string = "https://raw.githubusercontent.com/zerobotlabs/nestorbot/master/README.md"

type AppNestorJson struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Permalink   string `json:"permalink"`
}

func runNew(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		color.Red("Usage: nestor new <power-name>")
		os.Exit(1)
	}

	appName := strings.TrimSpace(args[0])
	permalink := strings.ToLower(strings.Replace(appName, " ", "-", -1))

	currentDir, err := os.Getwd()
	if err != nil {
		color.Red("Unexpected error occurred while creating your power: %+v", err)
		os.Exit(1)
	}

	base := path.Join(currentDir, permalink)
	if _, err = os.Stat(base); err == nil {
		color.Red("A power with this name already exists in this directory, pick another name or another directory to create your power")
		os.Exit(1)
	}

	color.Green("+ Creating directory...")
	err = os.Mkdir(base, 0755)
	if err != nil {
		color.Red("- Unexpected error occurred while creating your power: %+v", err)
		os.Exit(1)
	}

	binary, lookErr := exec.LookPath("npm")
	if lookErr != nil {
		color.Red("- Could not find npm in your $PATH. Please install nodejs from https://nodejs.org")
		os.Exit(1)
	}

	color.Green("+ Running npm init...")
	npmInitCommand := exec.Command(binary, "init", "-y")
	npmInitCommand.Dir = base
	err = npmInitCommand.Run()
	if err != nil {
		color.Red("- Unexpected Error while creating package.json for your power")
		os.Exit(1)
	}

	color.Green("+ Running npm install nestorbot...")
	npmInstallCommand := exec.Command(binary, "install", "nestorbot", "--save")
	npmInstallCommand.Dir = base
	err = npmInstallCommand.Run()
	if err != nil {
		color.Red("- Unexpected Error while downloading nestorbot")
		os.Exit(1)
	}

	color.Green("+ Creating nestor.json...")
	var app AppNestorJson
	app.Name = appName
	app.Description = appName
	app.Permalink = permalink

	marshaledApp, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		color.Red("- Unexpected Error while creating nestor.json")
		os.Exit(1)
	}

	nestorJsonPath := path.Join(base, "nestor.json")
	err = ioutil.WriteFile(nestorJsonPath, marshaledApp, 0644)
	if err != nil {
		color.Red("- Unexpected Error while creating nestor.json: %+v", err)
		os.Exit(1)
	}

	color.Green("+ Creating example at index.js...")
	indexFilePath := path.Join(base, "index.js")
	err = ioutil.WriteFile(indexFilePath, []byte(defaultAppFile), 0644)
	if err != nil {
		color.Red("- Unexpected Error while creating index.js: %+v", err)
		os.Exit(1)
	}

	color.Green("+ Downloading README...")
	resp, err := http.Get(readmeUrl)
	if err != nil {
		color.Red("- Unexpected Error while fetching README: %+v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	readmePath := path.Join(base, "README.md")
	readme, err := os.Create(readmePath)
	if err != nil {
		color.Red("- Unexpected Error while creating README: %+v", err)
		os.Exit(1)
	}
	defer readme.Close()

	_, err = io.Copy(readme, resp.Body)
	if err != nil {
		color.Red("- Unexpected Error while creating README: %+v", err)
		os.Exit(1)
	}

	color.Green("Successfully created power %s here: %s", appName, base)
	fmt.Printf("\ncd %s and run `nestor save` and `nestor shell` to start working on your power\n", base)
}
