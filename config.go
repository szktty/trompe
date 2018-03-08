package trompe

import "fmt"

var DebugMode = false
var VerboseMode = false

func Debug(format string, arg ...interface{}) {
	if DebugMode {
		fmt.Printf(format, arg...)
		fmt.Println()
	}
}

func Verbose(format string, arg ...interface{}) {
	if VerboseMode || DebugMode {
		fmt.Printf(format, arg...)
		fmt.Println()
	}
}
