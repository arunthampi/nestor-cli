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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/Bowery/prompt"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Nestor with your username and password",
	Run:   runLogin,
}

var unexpectedErrorWhileLoggingInErr error = fmt.Errorf("Unexpected error while logging in")
var unexpectedErrorWhileLoggingOutErr error = fmt.Errorf("Unexpected error while logging out")
var rootDir string = "/tmp"

const nestorRoot string = ".nestor"
const tokenFileName string = "token"

func runLogin(cmd *cobra.Command, args []string) {
	if loginInfo := savedLoginInfo(); loginInfo != nil {
		fmt.Printf("You are already logged in as %s. To logout, type \"nestor logout\"\n", loginInfo.Email)
		os.Exit(1)
	}

	email := getEmail()
	password := getPassword()

	loginInfo, err := nestorclient.Login(email, password)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = saveToken(loginInfo)
	if err != nil {
		fmt.Println(unexpectedErrorWhileLoggingInErr.Error())
		os.Exit(1)
	}

	fmt.Printf("Successfully logged in as %s\n", email)
}

// Prompts the user for an email
func getEmail() string {
	email, err := prompt.Basic("Your email: ", true)
	if err != nil {
		fmt.Printf("Unexpected error while getting your email")
		os.Exit(1)
	}

	return email
}

// Prompts the user for a password
func getPassword() string {
	password, err := prompt.Password("Your password (if you're not sure what your password is, set it at https://www.asknestor.me/users/me/edit): ")
	if err != nil {
		fmt.Printf("Unexpected error while getting your password")
		os.Exit(1)
	}

	return password
}

func saveToken(loginInfo *nestorclient.LoginInfo) error {
	parentDir := path.Join("/tmp", nestorRoot)
	err := os.MkdirAll(parentDir, 0755)
	if err != nil {
		return err
	}

	loginJson, err := json.Marshal(loginInfo)
	if err != nil {
		return err
	}

	p := path.Join(parentDir, tokenFileName)
	return ioutil.WriteFile(p, loginJson, 0644)
}

func savedLoginInfo() *nestorclient.LoginInfo {
	var l nestorclient.LoginInfo

	p := path.Join("/tmp", nestorRoot, tokenFileName)

	loginJson, err := ioutil.ReadFile(p)
	if err != nil {
		return nil
	}

	if err := json.Unmarshal(loginJson, &l); err != nil {
		return nil
	}

	return &l
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
