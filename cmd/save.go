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

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
	"github.com/zerobotlabs/nestor-cli/utils"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Saves your app to Nestor",
	Run:   runSave,
}

func runSave(cmd *cobra.Command, args []string) {
	var l *nestorclient.LoginInfo
	var app nestorclient.App

	// Check if you are logged in first
	if l = utils.SavedLoginInfo(); l == nil {
		fmt.Printf("You are not logged in. To login, type \"nestor login\"\n")
		os.Exit(1)
	}

	// Check if you have a valid nestor.json file
	nestorJsonPath, err := pathToNestorJson(args)
	if err != nil {
		fmt.Printf("Could not find nestor.json in the path specified\n")
		os.Exit(1)
	}

	app.ManifestPath = nestorJsonPath

	err = app.ParseManifest()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	// Check if existing app exists and if so, then we should be making calls to the "UPDATE" function
	// We are ignoring the error for now but at some point we will have to show an error that is not annoying
	err = app.FetchDetails(l)
	if err != nil {
		fmt.Printf("Error fetching details for app\n")
	}

	fmt.Printf("Building deployment artifact...\n")
	err = app.BuildArtifact()
	if err != nil {
		fmt.Printf("Error while building deployment artifact for your app\n")
	}

	// Check if you need to do coffee compilation
	err = app.CompileCoffeescript()
	if err != nil {
		fmt.Printf("There was an error compiling coffeescript in your app\n")
		os.Exit(1)
	}

	fmt.Printf("Calculating SHA256 of artifact...\n")
	err = app.CalculateLocalSha256()
	if err != nil {
		fmt.Printf("Error while calculating SHA256 of artifact\n")
		os.Exit(1)
	}

	if app.LocalSha256 != app.RemoteSha256 {
		fmt.Printf("Generating zip...\n")
		zip, err := app.ZipBytes()
		if err != nil {
			fmt.Printf("Error creating a zip of your app's deployment artifact\n")
			os.Exit(1)
		}

		fmt.Printf("Uploading zip...\n")
		// Upload app contents
		buffer := bytes.NewBuffer(zip)
		err = app.Upload(buffer, l)
		if err != nil {
			fmt.Printf("Error while uploading deployment artifact: %+v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Not uploading your app, since there are no changes to the code\n")
	}

	// Make API call to Nestor with contents from JSON file along with S3 URL so that the API can create a functioning bot app
	fmt.Printf("Saving app to Nestor...\n")
	err = app.SaveToNestor(l)
	if err != nil {
		fmt.Printf("Error while saving app to nestor: %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully saved app to Nestor!\n")
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
