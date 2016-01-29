package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/zerobotlabs/nestor-cli/nestorclient"
)

const nestorRoot string = ".nestor"
const tokenFileName string = "login"

func SaveLoginInfo(loginInfo *nestorclient.LoginInfo) error {
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

func SavedLoginInfo() *nestorclient.LoginInfo {
	var l nestorclient.LoginInfo

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

func RemoveLoginInfo() error {
	h, err := homedir.Dir()
	if err != nil {
		return err
	}

	p := path.Join(h, nestorRoot, tokenFileName)
	return os.Remove(p)
}

// Sha256 returns a base64 encoded SHA256 hash of `b`.
func Sha256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
