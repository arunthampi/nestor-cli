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
	"os"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/update"
)

func runUpdate(cmd *cobra.Command, args []string) {
	newVersion, err := update.Update()
	switch {
	case err == update.NotAvailableErr:
		color.Green("No update available\n")
		os.Exit(0)
	case err != nil:
		color.Red("There was an error updating nestor. Please try again later or contact help@asknestor.me\n")
		os.Exit(1)
	}

	color.Green("Updated nestor to new version: %s!\n", newVersion)
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the nestor application",
	Run:   runUpdate,
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
