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
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/app"
	"github.com/zerobotlabs/nestor-cli/login"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Saves your app to Nestor",
	Run:   runSave,
}

func runSave(cmd *cobra.Command, args []string) {
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
		color.Red("- Error fetching details for app\n")
	}

	color.Green("+ Building deployment artifact...\n")
	err = a.BuildArtifact()
	if err != nil {
		color.Red("- Error while building deployment artifact for your app\n")
	}

	// Check if you need to do coffee compilation
	err = a.CompileCoffeescript()
	if err != nil {
		color.Red("- There was an error compiling coffeescript in your app\n")
		os.Exit(1)
	}

	err = a.CalculateLocalSha256()
	if err != nil {
		color.Red("- There was an error calculating whether your app needs to be uploaded\n")
		os.Exit(1)
	}

	if a.LocalSha256 != a.RemoteSha256 {
		color.Green("+ Generating zip...\n")
		zip, err := a.ZipBytes()
		if err != nil {
			color.Red("- Error creating a zip of your app's deployment artifact\n")
			os.Exit(1)
		}

		color.Green("+ Uploading zip...\n")
		// Upload app contents
		buffer := bytes.NewBuffer(zip)
		err = a.Upload(buffer, l)
		if err != nil {
			color.Red("- Error while uploading deployment artifact: %+v\n", err)
			os.Exit(1)
		}
	}

	// Make API call to Nestor with contents from JSON file along with S3 URL so that the API can create a functioning bot app
	color.Green("+ Saving app to Nestor...\n")
	err = a.SaveToNestor(l)
	if err != nil {
		color.Red("- Error while saving app to nestor: %+v\n", err)
		os.Exit(1)
	}

	color.Green("+ Successfully saved app to Nestor!\n")
	fmt.Printf("\nYou can test your app by running `nestor shell`\n")
	fmt.Printf("To deploy your app to Slack, run `nestor deploy --latest`\n")
}

func init() {
	RootCmd.AddCommand(saveCmd)
}

func pathToNestorJson(args []string) (string, error) {
	nestorJsonPath := func(base string) (string, error) {
		var err error
		p := path.Join(base, "nestor.json")
		if _, err = os.Stat(p); err != nil {
			return "", err
		}

		return p, err
	}

	if len(args) > 0 {
		return nestorJsonPath(args[0])
	} else {
		base, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return nestorJsonPath(base)
	}
}
