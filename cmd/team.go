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

	"github.com/Bowery/prompt"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
	"github.com/zerobotlabs/nestor-cli/utils"
)

// teamCmd represents the team command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Sets your Default Slack Team",
	Run:   runTeam,
}

func runTeam(cmd *cobra.Command, args []string) {
	var l *nestorclient.LoginInfo

	if l = utils.SavedLoginInfo(); l == nil {
		fmt.Printf("You are not logged in. To login, type \"nestor login\"\n")
		os.Exit(1)
	}

	teams, err := nestorclient.GetTeams(l)
	if err != nil {
		fmt.Println(unexpectedErrorWhileLoggingInErr.Error())
		os.Exit(1)
	}

	if len(teams) == 1 {
		saveDefaultTeam(l, teams[0])
	} else {
		ok := false
		intIndex := 0

		for !ok {
			index, promptErr := prompt.Basic(chooseTeamQuestion(teams), true)
			if promptErr != nil {
				os.Exit(1)
			}

			intIndex, err = strconv.Atoi(index)
			if err == nil && intIndex > 0 && intIndex <= len(teams) {
				ok = true
			}
		}

		saveDefaultTeam(l, teams[intIndex-1])
	}
}

func chooseTeamQuestion(teams []nestorclient.Team) string {
	question := "Pick which team you want to set as your default: "
	for i, team := range teams {
		question += fmt.Sprintf("%d. %s | ", i+1, team.Name)
	}
	question += fmt.Sprintf("Choose %d-%d: ", 1, len(teams))

	return question
}

func saveDefaultTeam(l *nestorclient.LoginInfo, team nestorclient.Team) {
	l.DefaultTeamId = team.Id
	err := utils.SaveLoginInfo(l)
	if err != nil {
		fmt.Printf("Unexpected error while setting your default team\n")
		os.Exit(1)
	} else {
		fmt.Printf("Set your default team as %s (%s)\n", team.Name, team.Url)
	}
}

func init() {
	RootCmd.AddCommand(teamCmd)
}
