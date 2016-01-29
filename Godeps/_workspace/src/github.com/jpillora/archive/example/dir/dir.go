package main

import (
	"os"

	"github.com/zerobotlabs/nestor-cli/Godeps/_workspace/src/github.com/jpillora/archive"
)

func main() {
	a, _ := archive.NewWriter("file.tar", os.Stdout) //detects .tar
	a.AddDir("foo")
	a.Close()
}
