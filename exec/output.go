package exec

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zerobotlabs/nestor-cli/app"
	"github.com/zerobotlabs/nestor-cli/login"
)

type Output struct {
	Heartbeat bool     `json:"heartbeat"`
	Error     string   `json:"error"`
	RequestId string   `json:"request_id"`
	Logs      string   `json:"logs"`
	ToSend    []string `json:"to_send"`
	ToReply   []string `json:"to_reply"`
}

var host string = "http://localhost:5400"

func (o *Output) Exec(app *app.App, l *login.LoginInfo, message string) error {
	params := url.Values{
		"message": []string{message},
	}

	urlStr := fmt.Sprintf("%s/teams/%s/apps/%d/exec", host, l.DefaultTeamId, app.Id)
	r, err := http.NewRequest("POST", urlStr, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", l.Token))

	client := &http.Client{}
	resp, err := client.Do(r)

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		if err = json.Unmarshal([]byte(scanner.Text()), o); err != nil {
			return err
		}

		if o.Heartbeat == true {
			continue
		}

		if o.Error != "" {
			resp.Body.Close()
			return nil
		}

		if o.RequestId != "" {
			resp.Body.Close()
			return nil
		}
	}

	resp.Body.Close()
	return nil
}
