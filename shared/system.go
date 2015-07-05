package trompe

import (
	"os"
	"path"
	"strings"
)

var repo = "src/github.com/szktty/trompe"
var TROMPE_PATH = "TROMPE_PATH"

func SearchPath() []string {
	spath := make([]string, 0, 16)
	spath = append(spath, "") // current directory

	trpath := os.Getenv(TROMPE_PATH)
	if trpath == "" {
		gopath := os.Getenv("GOPATH")
		if gopath != "" {
			for _, e := range strings.Split(gopath, ":") {
				spath = append(spath,
					path.Join(e, "src", "github.com", "szktty", "trompe", "lib"))
			}

		}
	} else {
		for _, s := range strings.Split(trpath, ":") {
			spath = append(spath, s)
		}
	}
	return spath
}
