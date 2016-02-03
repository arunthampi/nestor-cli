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
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/login"
	"github.com/zerobotlabs/nestor-cli/team"
)

// teamCmd represents the team command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Sets your Default Slack Team",
	Run:   runTeam,
}

func runTeam(cmd *cobra.Command, args []string) {
	var l *login.LoginInfo

	if l = login.SavedLoginInfo(); l == nil {
		fmt.Printf("You are not logged in. To login, type \"nestor login\"\n")
		os.Exit(1)
	}

	teams, err := team.GetTeams(l)
	if err != nil {
		fmt.Println(unexpectedErrorWhileLoggingInErr.Error())
		os.Exit(1)
	}

	if len(teams) == 1 {
		team := teams[0]
		err := team.Save(l)
		if err != nil {
			fmt.Printf("Error saving default team\n")
			os.Exit(1)
		} else {
			fmt.Printf("Saved %s as your default team\n", team.Name)
		}
	} else {
		table := team.TableizeTeams(teams, l.DefaultTeamId)
		table.Render()
		fmt.Printf("\n")

		ok := false
		intIndex := 0

		for !ok {
			index, promptErr := prompt.Basic(fmt.Sprintf("Pick which team you want to set as your default (1-%d): ", len(teams)), true)
			if promptErr != nil {
				os.Exit(1)
			}

			intIndex, err = strconv.Atoi(index)
			if err == nil && intIndex > 0 && intIndex <= len(teams) {
				ok = true
			}
		}

		team := teams[intIndex-1]
		err := team.Save(l)
		if err != nil {
			fmt.Printf("Error saving default team\n")
			os.Exit(1)
		} else {
			fmt.Printf("Saved %s as your default team\n", team.Name)
		}
	}
}

func init() {
	RootCmd.AddCommand(teamCmd)
}
