package team

import (
	"encoding/json"
	"net/url"

	"github.com/zerobotlabs/nestor-cli/errors"
	"github.com/zerobotlabs/nestor-cli/login"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
)

type Team struct {
	Id   string `json:"uid"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

func GetTeams(loginInfo *login.LoginInfo) ([]Team, error) {
	var teams = []Team{}

	params := url.Values{
		"Authorization": []string{loginInfo.Token},
	}

	response, err := nestorclient.CallAPI("/teams", "GET", params, 200)

	if err != nil {
		return nil, errors.UnexpectedServerError
	}

	if err = json.Unmarshal([]byte(response), &teams); err != nil {
		// If JSON parsing fails that means it's a server error too
		return nil, errors.UnexpectedServerError
	}

	return teams, nil
}

func (t *Team) Save(l *login.LoginInfo) error {
	l.DefaultTeamId = t.Id
	err := l.Save()
	if err != nil {
		return err
	}

	return nil
}
