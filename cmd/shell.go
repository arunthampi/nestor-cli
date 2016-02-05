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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/peterh/liner"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/app"
	"github.com/zerobotlabs/nestor-cli/exec"
	"github.com/zerobotlabs/nestor-cli/login"
)

const nestorRoot string = ".nestor"
const historyFileName string = "history"

var historyFile string

func init() {
	h, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	historyFile = path.Join(h, nestorRoot, historyFileName)
	_, err = os.Stat(historyFile)
	if err != nil && os.IsNotExist(err) {
		ioutil.WriteFile(historyFile, []byte(""), 0644)
	}
}

func saveHistory(line *liner.State, hf *os.File) {
	line.WriteHistory(hf)
	line.Close()
	hf.Close()
}

func runShell(cmd *cobra.Command, args []string) {
	var l *login.LoginInfo
	var a app.App

	// Check if you are logged in first
	if l = login.SavedLoginInfo(); l == nil {
		color.Red("You are not logged in. To login, type \"nestor login\"\n")
		os.Exit(1)
	}

	// Check if you have a valid nestor.json file
	nestorJsonPath, err := pathToNestorJson(args)
	if err != nil {
		color.Red("Could not find nestor.json in the path specified\n")
		os.Exit(1)
	}

	a.ManifestPath = nestorJsonPath

	err = a.ParseManifest()
	if err != nil {
		color.Red("%s\n", err.Error())
		os.Exit(1)
	}

	// Check if existing app exists and if so, then we should be making calls to the "UPDATE" function
	// We are ignoring the error for now but at some point we will have to show an error that is not annoying
	err = a.Hydrate(l)
	if err != nil {
		color.Red("Error fetching details for app\n")
	}

	if a.Id == 0 {
		color.Red("You haven't saved your app yet. Run `nestor save` before you can debug your app\n")
		os.Exit(1)
	}

	ok := false

	line := liner.NewLiner()

	line.SetCtrlCAborts(true)
	hf, err := os.OpenFile(historyFile, os.O_RDWR, 0644)
	if err != nil {
		color.Red("Unexpected error opening shell\n")
		line.Close()
		os.Exit(1)
	}

	line.ReadHistory(hf)

	for !ok {
		if command, err := line.Prompt("nestor> "); err == nil {
			command = strings.TrimSpace(command)
			if command == "quit" || command == "exit" {
				fmt.Println("Goodbye!")
				saveHistory(line, hf)
				os.Exit(1)
			}

			if command == "" {
				continue
			}

			output := exec.Output{}
			err := output.Exec(&a, l, command)
			if err != nil {
				color.Red("unexpected error while running your app. Please try again later or contact hello@asknestor.me\n", err)
			}

			line.AppendHistory(command)

			if output.Logs != "" {
				color.Yellow(output.Logs)
			}

			for _, send := range output.ToSend {
				fmt.Println(send)
			}
			for _, reply := range output.ToReply {
				color.Blue("<@user>: ")
				fmt.Println(reply)
			}
		}
	}

	saveHistory(line, hf)
}

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Open an interactive shell to debug your Nestor app",
	Run:   runShell,
}

func init() {
	RootCmd.AddCommand(shellCmd)
}
