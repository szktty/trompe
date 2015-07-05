package trompe

import (
	"fmt"
	"runtime"
)

var version = "alpha"

func Version() string {
	fmt.Printf("root = %s\n", runtime.GOROOT())
	return fmt.Sprintf("%s (%s on %s)", version, runtime.GOOS, runtime.GOARCH)
}
