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
	"os"
	"strconv"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/Bowery/prompt"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/app"
	"github.com/zerobotlabs/nestor-cli/login"
	"github.com/zerobotlabs/nestor-cli/version"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a version of your app to Nestor",
	Run:   runDeploy,
}

var latestDeploy bool = false

func init() {
	f := deployCmd.Flags()
	f.BoolVarP(&latestDeploy, "latest", "l", false, "Latest Version")
}

func runDeploy(cmd *cobra.Command, args []string) {
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
		color.Red("You haven't saved your app yet. Run `nestor save` before you can deploy your app\n")
		os.Exit(1)
	}

	versions, err := version.FetchVersions(a, l)
	if err != nil {
		color.Red("Error fetching versions for your app\n")
		os.Exit(1)
	}

	if latestDeploy {
		pickedVersion := versions[0]
		if pickedVersion.CurrentlyDeployed {
			color.Green("Version %s already deployed. Nothing to do here\n", pickedVersion.Ref)
		} else {
			err = pickedVersion.Deploy(a, l)

			if err != nil {
				color.Red("Error deploying %s. Please try again later or contact hello@asknestor.me\n", pickedVersion.Ref)
				os.Exit(1)
			}

			color.Green("Deployed version %s successfully\n", pickedVersion.Ref)
		}
	} else {
		table := version.TableizeVersions(versions)
		fmt.Printf("\n")
		table.Render()
		fmt.Printf("\n")

		ok := false
		intIndex := 0

		for !ok {
			index, promptErr := prompt.Basic(fmt.Sprintf("Pick a version to deploy (1-%d): ", len(versions)), true)
			if promptErr != nil {
				os.Exit(1)
			}

			intIndex, err = strconv.Atoi(index)
			if err == nil && intIndex > 0 && intIndex <= len(versions) {
				ok = true
			}
		}

		pickedVersion := versions[intIndex-1]
		if pickedVersion.CurrentlyDeployed {
			color.Green("Version %s already deployed. Nothing to do here\n", pickedVersion.Ref)
		} else {
			err = pickedVersion.Deploy(a, l)

			if err != nil {
				color.Red("Error deploying %s. Please try again later or contact hello@asknestor.me\n", pickedVersion.Ref)
				os.Exit(1)
			}

			color.Green("Deployed version %s successfully\n", pickedVersion.Ref)
		}
	}
}

func init() {
	RootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
