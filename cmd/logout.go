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

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/Bowery/prompt"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
	"github.com/zerobotlabs/nestor-cli/utils"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log Out of Nestor",
	Run:   runLogout,
}

func runLogout(cmd *cobra.Command, args []string) {
	var l *nestorclient.LoginInfo

	if l = utils.SavedLoginInfo(); l == nil {
		fmt.Printf("You are not logged in. To login, type \"nestor login\"\n")
		os.Exit(1)
	}

	shouldLogout, err := prompt.Ask(fmt.Sprintf("Are you sure you want to log out as %s", l.Email))
	if err != nil {
		os.Exit(1)
	}

	if shouldLogout {
		err = utils.RemoveLoginInfo()
		if err != nil {
			fmt.Println(unexpectedErrorWhileLoggingOutErr)
			os.Exit(1)
		}
		fmt.Printf("Successfully logged out as %s\n", l.Email)
	} else {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
