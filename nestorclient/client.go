package nestorclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var Host string = "http://localhost:5000"

type LoginInfo struct {
	Email         string `json:"email"`
	Token         string `json:"token"`
	Err           string `json:"error"`
	DefaultTeamId string
}

type Team struct {
	Id   string `json:"uid"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

var UnexpectedServerError error = fmt.Errorf("Unexpected response from the Nestor API. Try again after a while or contact help@asknestor.me")

func GetTeams(loginInfo *LoginInfo) ([]Team, error) {
	var teams = []Team{}

	params := url.Values{
		"Authorization": []string{loginInfo.Token},
	}

	response, err := callAPI("/teams", "GET", params, 200)

	if err != nil {
		return nil, UnexpectedServerError
	}

	if err = json.Unmarshal([]byte(response), &teams); err != nil {
		// If JSON parsing fails that means it's a server error too
		return nil, UnexpectedServerError
	}

	return teams, nil
}

func HydrateApp(app *App, loginInfo *LoginInfo) error {
	params := url.Values{
		"Authorization":  []string{loginInfo.Token},
		"app[permalink]": []string{app.Permalink},
	}

	response, err := callAPI(fmt.Sprintf("/teams/%s/apps/search", loginInfo.DefaultTeamId), "GET", params, 200)

	if err != nil {
		return UnexpectedServerError
	}

	if err = json.Unmarshal([]byte(response), app); err != nil {
		// If JSON parsing fails that means it's a server error too
		return UnexpectedServerError
	}

	return nil
}

func UploadUrl(loginInfo *LoginInfo) (*url.URL, error) {
	type _urlPayload struct {
		Url string `json:"url"`
	}

	var urlPayload _urlPayload

	params := url.Values{
		"Authorization": []string{loginInfo.Token},
	}
	response, err := callAPI(fmt.Sprintf("/teams/%s/apps/issue_upload_url", loginInfo.DefaultTeamId), "POST", params, 200)
	if err != nil {
		return nil, UnexpectedServerError
	}

	if err = json.Unmarshal([]byte(response), &urlPayload); err != nil {
		// If JSON parsing fails that means it's a server error too
		return nil, UnexpectedServerError
	}

	return url.Parse(urlPayload.Url)
}

func Login(email string, password string) (*LoginInfo, error) {
	l := LoginInfo{Email: email}

	params := url.Values{
		"user[email]":    []string{email},
		"user[password]": []string{password},
	}

	response, err := callAPI("/users/issue_token", "POST", params, 200)

	if err != nil {
		return nil, UnexpectedServerError
	}

	if err = json.Unmarshal([]byte(response), &l); err != nil {
		// If JSON parsing fails that means it's a server error too
		return nil, UnexpectedServerError
	}

	// Wrap the error from the API in an error struct
	if l.Err != "" {
		return nil, fmt.Errorf(l.Err)
	}

	return &l, nil
}

func callAPI(path string, method string, params url.Values, expectedStatusCode int) (string, error) {
	urlStr, token := parseURLStringAndToken(path, method, params)
	r, _ := http.NewRequest(method, urlStr, strings.NewReader(params.Encode()))

	if method != "GET" {
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	}

	if token != "" {
		r.Header.Add("Authorization", token)
	}

	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		} else {
			if resp.StatusCode != expectedStatusCode {
				return "", fmt.Errorf("Expected Status Code: %d, Got: %d\n", expectedStatusCode, resp.StatusCode)
			} else {
				return string(contents), err
			}
		}
	}
}

func parseURLStringAndToken(path string, method string, params url.Values) (string, string) {
	var token string

	u, err := url.ParseRequestURI(Host)
	if err != nil {
		return "", ""
	}

	params.Set("format", "json")
	// If token is present in params, then take it out set as a header
	if token = params.Get("Authorization"); token != "" {
		params.Del("Authorization")
	}

	u.Path = path

	if method == "GET" {
		u.RawQuery = params.Encode()
	}
	urlStr := fmt.Sprintf("%v", u)

	return urlStr, token
}
