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

var host string = "http://localhost:5000"

type NestorAPIError struct {
	Errors []string `json:"errors"`
}

func (e NestorAPIError) Error() string {
	formattedErrors := []string{"Oops, encountered these errors:"}

	for _, e1 := range e.Errors {
		formattedErrors = append(formattedErrors, fmt.Sprintf("· %s", e1))
	}

	return strings.Join(formattedErrors, "\n")
}

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
			var ne NestorAPIError
			err := json.Unmarshal(contents, &ne)

			if err == nil {
				if len(ne.Errors) > 0 {
					return "", ne
				}
			}

			return string(contents), nil
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
