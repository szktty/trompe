package main

import (
	"flag"
	"fmt"
	"github.com/szktty/trompe"
	"os"
)

var debugModeOpt = flag.Bool("d", false, "debug mode")
var verboseModeOpt = flag.Bool("v", false, "verbose mode")
var versionModeOpt = flag.Bool("version", false, "print version")

func main() {
	flag.Parse()

	trompe.DebugMode = *debugModeOpt
	trompe.VerboseMode = *verboseModeOpt

	if *versionModeOpt {
		fmt.Printf("%s\n", trompe.Version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Printf("Usage: trompe [options] files\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	trompe.Init()

	if flag.Arg(0) == "test" {
		trompe.TestCompiledCodeHelloWorld()
		trompe.TestCompiledCodeFizzBuzzCompare()
		trompe.TestCompiledCodeFizzBuzzMatch()
	}
}
