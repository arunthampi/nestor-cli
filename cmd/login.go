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
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/login"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Nestor with your username and password",
	Run:   runLogin,
}

var unexpectedErrorWhileLoggingInErr error = fmt.Errorf("Unexpected error while logging in")
var unexpectedErrorWhileLoggingOutErr error = fmt.Errorf("Unexpected error while logging out")
var unexpectedErrorWhileFetchingTeamsErr error = fmt.Errorf("Unexpected error while fetching teams")

func runLogin(cmd *cobra.Command, args []string) {
	if loginInfo := login.SavedLoginInfo(); loginInfo != nil {
		color.Red("You are already logged in as %s. To logout, type \"nestor logout\"\n", loginInfo.Email)
		os.Exit(1)
	}

	email := getEmail()
	password := getPassword()

	loginInfo, err := login.Login(email, password)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	err = loginInfo.Save()
	if err != nil {
		color.Red(unexpectedErrorWhileLoggingInErr.Error())
		os.Exit(1)
	}

	color.Green("Successfully logged in as %s\n", email)

	// Make user pick a default team right at the beginning
	runTeam(cmd, args)
}

// Prompts the user for an email
func getEmail() string {
	email, err := prompt.Basic("Your email: ", true)
	if err != nil {
		os.Exit(1)
	}

	return email
}

// Prompts the user for a password
func getPassword() string {
	password, err := prompt.Password("Your password (if you're not sure what your password is, set it at https://www.asknestor.me/users/me/edit): ")
	if err != nil {
		os.Exit(1)
	}

	return password
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
