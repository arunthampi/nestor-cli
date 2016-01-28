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
	Email string `json:"email"`
	Token string `json:"token"`
	Err   string `json:"error"`
}

var UnexpectedServerError error = fmt.Errorf("Unexpected response from the Nestor API. Try again after a while or contact help@asknestor.me")

func Login(email string, password string) (*LoginInfo, error) {
	l := LoginInfo{Email: email}

	params := url.Values{}
	params.Set("user[email]", email)
	params.Set("user[password]", password)

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
	u, err := url.ParseRequestURI(Host)
	if err != nil {
		return "", err
	}

	u.Path = path
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	params.Set("format", "json")
	r, _ := http.NewRequest(method, urlStr, strings.NewReader(params.Encode()))

	if method != "GET" {
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	}

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
