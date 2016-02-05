package update

import (
	"fmt"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/equinox-io/equinox"
	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/fatih/color"
)

const appId = "app_34CVQima2NP"

var publicKey = []byte(`
-----BEGIN ECDSA PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEaXhVnJySrP7acf4zfJafQ0Zd33gHa7sd
FJJrh6YVtPdRDJOa9kbCyEI+jnndtgk++oka5hS93hlzLvjz+KHf4it3tKaVaz3N
j+aGmaSPtrmz0XgA6Tw4ci9RVpvm3XMW
-----END ECDSA PUBLIC KEY-----
`)
var NotAvailableErr error = fmt.Errorf("No update available")

func Update() (string, error) {
	var opts equinox.Options
	if err := opts.SetPublicKeyPEM(publicKey); err != nil {
		return "", err
	}

	// check for the update
	resp, err := equinox.Check(appId, opts)
	switch {
	case err == equinox.NotAvailableErr:
		return "", NotAvailableErr
	case err != nil:
		return "", err
	}

	// fetch the update and apply it
	color.Green("Updating...")
	err = resp.Apply()
	if err != nil {
		return "", err
	}

	return resp.ReleaseVersion, nil
}
