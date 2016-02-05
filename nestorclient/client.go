package nestorclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var host string = "http://localhost:5000"

func CallAPI(path string, method string, params url.Values, expectedStatusCode int) (string, error) {
	urlStr, token := parseURLStringAndToken(path, method, params)
	r, _ := http.NewRequest(method, urlStr, strings.NewReader(params.Encode()))

	if method != "GET" {
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	}

	if token != "" {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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

	u, err := url.ParseRequestURI(host)
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
