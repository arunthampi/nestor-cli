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
	"path"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/Bowery/prompt"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log Out of Nestor",
	Run:   runLogout,
}

func runLogout(cmd *cobra.Command, args []string) {
	var l *nestorclient.LoginInfo

	if l = savedLoginInfo(); l == nil {
		fmt.Printf("You are not logged in. To login, type \"nestor login\"\n")
		os.Exit(1)
	}

	shouldLogout, err := prompt.Ask(fmt.Sprintf("Are you sure you want to log out as %s", l.Email))
	if err != nil {
		fmt.Printf("Unexpected error while getting your confirmation\n")
		os.Exit(1)
	}

	if shouldLogout {
		err = removeToken()
		if err != nil {
			fmt.Println(unexpectedErrorWhileLoggingOutErr)
			os.Exit(1)
		}
		fmt.Printf("Successfully logged out as %s\n", l.Email)
	} else {
		os.Exit(1)
	}
}

func removeToken() error {
	p := path.Join("/tmp", nestorRoot, tokenFileName)
	return os.Remove(p)
}

func init() {
	RootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
