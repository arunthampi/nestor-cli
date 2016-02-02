package login

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/zerobotlabs/nestor-cli/errors"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
)

const nestorRoot string = ".nestor"
const tokenFileName string = "login"

type LoginInfo struct {
	Email         string `json:"email"`
	Token         string `json:"token"`
	Err           string `json:"error"`
	DefaultTeamId string
}

func Login(email string, password string) (*LoginInfo, error) {
	l := LoginInfo{Email: email}

	params := url.Values{
		"user[email]":    []string{email},
		"user[password]": []string{password},
	}

	response, err := nestorclient.CallAPI("/users/issue_token", "POST", params, 200)

	if err != nil {
		return nil, errors.UnexpectedServerError
	}

	if err = json.Unmarshal([]byte(response), &l); err != nil {
		// If JSON parsing fails that means it's a server error too
		return nil, errors.UnexpectedServerError
	}

	// Wrap the error from the API in an error struct
	if l.Err != "" {
		return nil, fmt.Errorf(l.Err)
	}

	return &l, nil
}

func (loginInfo *LoginInfo) Save() error {
	h, err := homedir.Dir()
	if err != nil {
		return err
	}

	parentDir := path.Join(h, nestorRoot)
	err = os.MkdirAll(parentDir, 0755)
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

func SavedLoginInfo() *LoginInfo {
	var l LoginInfo

	h, err := homedir.Dir()
	if err != nil {
		return nil
	}

	p := path.Join(h, nestorRoot, tokenFileName)

	loginJson, err := ioutil.ReadFile(p)
	if err != nil {
		return nil
	}

	if err := json.Unmarshal(loginJson, &l); err != nil {
		return nil
	}

	return &l
}

func Delete() error {
	h, err := homedir.Dir()
	if err != nil {
		return err
	}

	p := path.Join(h, nestorRoot, tokenFileName)
	return os.Remove(p)
}
