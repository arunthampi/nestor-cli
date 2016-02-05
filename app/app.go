package app

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/jpillora/archive"
	"github.com/zerobotlabs/nestor-cli/errors"
	"github.com/zerobotlabs/nestor-cli/login"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
	"github.com/zerobotlabs/nestor-cli/shim"
)

type App struct {
	Id              int64           `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Permalink       string          `json:"permalink"`
	Public          bool            `json:"public"`
	EnvironmentKeys map[string]bool `json:"environment_keys"`
	RemoteSha256    string          `json:"sha_256"`
	GitRevision     string          `json:"git_revision"`
	LocalSha256     string
	UploadKey       string
	ManifestPath    string
	ArtifactPath    string
}

func (a *App) SourcePath() string {
	return path.Dir(a.ManifestPath)
}

func (a *App) ParseManifest() error {
	contents, err := ioutil.ReadFile(a.ManifestPath)
	if err != nil {
		return fmt.Errorf("Error reading nestor.json file: %s", a.ManifestPath)
	} else {
		if err = json.Unmarshal([]byte(contents), a); err != nil {
			fmt.Errorf("Error reading nestor.json file: %s", a.ManifestPath)
			os.Exit(1)
		}
	}

	// TODO: Better error descriptions
	if a.Name == "" || a.Permalink == "" || a.Description == "" {
		return fmt.Errorf("You need to set a valid 'name', 'permalink' and 'description' in your nestor.json")
	}

	return nil
}

func (a *App) Hydrate(loginInfo *login.LoginInfo) error {
	params := url.Values{
		"Authorization":  []string{loginInfo.Token},
		"app[permalink]": []string{a.Permalink},
	}

	response, err := nestorclient.CallAPI(fmt.Sprintf("/teams/%s/apps/search", loginInfo.DefaultTeamId), "GET", params, 200)

	if err != nil {
		return errors.UnexpectedServerError
	}

	if err = json.Unmarshal([]byte(response), a); err != nil {
		// If JSON parsing fails that means it's a server error too
		return errors.UnexpectedServerError
	}

	return nil
}

func (a *App) UploadUrl(loginInfo *login.LoginInfo) (*url.URL, error) {
	type _urlPayload struct {
		Url string `json:"url"`
	}

	var urlPayload _urlPayload

	params := url.Values{
		"Authorization": []string{loginInfo.Token},
	}
	response, err := nestorclient.CallAPI(fmt.Sprintf("/teams/%s/apps/issue_upload_url", loginInfo.DefaultTeamId), "POST", params, 200)
	if err != nil {
		return nil, errors.UnexpectedServerError
	}

	if err = json.Unmarshal([]byte(response), &urlPayload); err != nil {
		// If JSON parsing fails that means it's a server error too
		return nil, errors.UnexpectedServerError
	}

	return url.Parse(urlPayload.Url)
}

func (a *App) CompileCoffeescript() error {
	coffeeFiles := []string{}

	if a.ArtifactPath == "" {
		return fmt.Errorf("ArtifactPath not set")
	}

	err := filepath.Walk(a.ArtifactPath, func(p string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() && strings.HasSuffix(p, ".coffee") {
			coffeeFiles = append(coffeeFiles, p)
		}
		return err
	})

	if err != nil {
		return err
	}

	if len(coffeeFiles) == 0 {
		return nil
	}

	color.Green("+ Compiling Coffeescript...\n")

	binary, lookErr := exec.LookPath("coffee")
	if lookErr != nil {
		color.Red("- Could not find coffee in your $PATH. Please install coffee-script with 'npm install -g coffee-script'\n")
		return lookErr
	}

	for _, coffeeFile := range coffeeFiles {
		_, err := exec.Command(binary, "-c", coffeeFile).Output()
		if err != nil {
			return err
		}

		err = os.Remove(coffeeFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Upload(body io.Reader, l *login.LoginInfo) error {
	uploadUrl, err := a.UploadUrl(l)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("PUT", uploadUrl.String(), body)
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(r)

	if err != nil {
		return nil
	}

	a.UploadKey = uploadUrl.Path

	return nil
}

// Stolen from https://github.com/apex/apex/function/function.go
func (a *App) Zip() (io.Reader, error) {
	buf := new(bytes.Buffer)
	zip := archive.NewZipWriter(buf)

	if err := zip.AddDir(a.ArtifactPath); err != nil {
		return nil, err
	}

	if err := zip.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func (a *App) CalculateLocalSha256() error {
	h := sha256.New()

	walkFunc := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			contents, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			h.Write(contents)
		}

		return nil
	}

	err := filepath.Walk(a.ArtifactPath, walkFunc)
	if err != nil {
		return err
	}

	a.LocalSha256 = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return nil
}

// ZipBytes returns the generated zip as bytes.
func (a *App) ZipBytes() ([]byte, error) {
	r, err := a.Zip()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (a *App) SaveToNestor(l *login.LoginInfo) error {
	var returnedApp App
	var response string
	var err error

	envKeysPayload, err := json.Marshal(a.EnvironmentKeys)
	if err != nil {
		return err
	}

	params := url.Values{
		"Authorization":         []string{l.Token},
		"app[name]":             []string{a.Name},
		"app[permalink]":        []string{a.Permalink},
		"app[upload_key]":       []string{a.UploadKey},
		"app[sha_256]":          []string{a.LocalSha256},
		"app[environment_keys]": []string{string(envKeysPayload)},
	}

	// If the app hasn't been created yet, then create it, otherwise update it
	if a.Id == 0 {
		response, err = nestorclient.CallAPI(fmt.Sprintf("/teams/%s/apps", l.DefaultTeamId), "POST", params, 201)
	} else {
		response, err = nestorclient.CallAPI(fmt.Sprintf("/teams/%s/apps/%d", l.DefaultTeamId, a.Id), "PATCH", params, 200)
	}

	if err = json.Unmarshal([]byte(response), &returnedApp); err != nil {
		return errors.UnexpectedServerError
	}

	a.Id = returnedApp.Id

	return err
}

// Recursively copy from one directory to another without copying over stuff from
// shit we don't need
func (a *App) BuildArtifact() error {
	tmpdir, err := ioutil.TempDir("", "nestor")
	scriptDir := path.Join(tmpdir, "script")

	a.ArtifactPath = tmpdir

	if err != nil {
		return nil
	}

	walkFunc := func(p string, info os.FileInfo, err error) error {
		p = strings.TrimPrefix(p, a.SourcePath())

		if (p == "/.git") && info.IsDir() {
			return filepath.SkipDir
		}

		if info.IsDir() {
			os.MkdirAll(path.Join(scriptDir, p), info.Mode())
		} else if info.Mode().IsRegular() {
			copyFileContents := func(src, dst string) (err error) {
				in, err := os.Open(src)
				if err != nil {
					return
				}
				defer in.Close()
				out, err := os.Create(dst)
				if err != nil {
					return
				}
				defer func() {
					cerr := out.Close()
					if err == nil {
						err = cerr
					}
				}()
				if _, err = io.Copy(out, in); err != nil {
					return
				}
				err = out.Sync()
				return nil
			}

			// Ignore nestor.json, we don't want to upload a zip just because
			// app metadata has changed
			if p == "/nestor.json" {
				return nil
			}

			return copyFileContents(path.Join(a.SourcePath(), p), path.Join(scriptDir, p))
		}

		return nil
	}

	err = filepath.Walk(a.SourcePath(), walkFunc)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(a.ArtifactPath, "index.js"), shim.MustAsset("index.js"), 0644)
	if err != nil {
		return err
	}

	return nil
}
