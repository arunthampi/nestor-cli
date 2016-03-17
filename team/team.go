package team

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/zerobotlabs/nestor-cli/errors"
	"github.com/zerobotlabs/nestor-cli/login"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
)

type Team struct {
	Id   string `json:"uid"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

func TableizeTeams(teams []Team, defaultTeamId string) *tablewriter.Table {
	var elems [][]string

	for i, t := range teams {
		isDefault := ""

		if defaultTeamId == t.Id {
			isDefault = "*"
		}

		elems = append(elems, []string{fmt.Sprintf("%d", i+1), t.Name, isDefault})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No.", "Team", "Default"})
	table.SetBorder(false)
	table.AppendBulk(elems)

	return table
}

func GetTeams(loginInfo *login.LoginInfo) ([]Team, error) {
	var teams = []Team{}

	params := url.Values{
		"Authorization": []string{loginInfo.Token},
	}

	response, err := nestorclient.CallAPI("/teams", "GET", params, 200)

	if err != nil {
		if ne, ok := err.(nestorclient.NestorAPIError); ok {
			return nil, ne
		}
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
