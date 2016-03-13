package version

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/dustin/go-humanize"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/zerobotlabs/nestor-cli/app"
	"github.com/zerobotlabs/nestor-cli/errors"
	"github.com/zerobotlabs/nestor-cli/login"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
)

type Version struct {
	Ref               string `json:"ref"`
	CreatedTimestamp  int64  `json:"created_timestamp"`
	CurrentlyDeployed bool   `json:"currently_deployed"`
}

func TableizeVersions(versions []Version) *tablewriter.Table {
	var elems [][]string

	for i, v := range versions {
		currentlyDeployed := ""

		if v.CurrentlyDeployed {
			currentlyDeployed = "*"
		}

		elems = append(elems, []string{fmt.Sprintf("%d", i+1), v.Ref, v.HumanizedCreatedAt(), currentlyDeployed})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"No.", "Ref", "Saved At", "Deployed"})
	table.SetBorder(false)
	table.AppendBulk(elems)

	return table
}

func FetchVersions(a app.App, l *login.LoginInfo) ([]Version, error) {
	params := url.Values{
		"Authorization": []string{l.Token},
	}
	var versions []Version

	response, err := nestorclient.CallAPI(fmt.Sprintf("/teams/%s/powers/%s/versions", l.DefaultTeamId, a.Permalink), "GET", params, 200)
	if err != nil {
		return versions, err
	}

	if err = json.Unmarshal([]byte(response), &versions); err != nil {
		// If JSON parsing fails that means it's a server error too
		return versions, errors.UnexpectedServerError
	}

	return versions, nil
}

func (v *Version) HumanizedCreatedAt() string {
	t := time.Unix(v.CreatedTimestamp, 0)
	return humanize.Time(t)
}

func (v *Version) Deploy(a app.App, l *login.LoginInfo) error {
	params := url.Values{
		"Authorization": []string{l.Token},
		"version[ref]":  []string{v.Ref},
	}

	_, err := nestorclient.CallAPI(fmt.Sprintf("/teams/%s/powers/%s/deploys", l.DefaultTeamId, a.Permalink), "POST", params, 201)
	if err != nil {
		return err
	}

	return nil
}
