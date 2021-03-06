package main

import (
	"flag"
	"fmt"
	"github.com/szktty/trompe"
	"github.com/szktty/trompe/parser"
	"os"
)

var debugModeOpt = flag.Bool("d", false, "debug mode")
var verboseModeOpt = flag.Bool("v", false, "verbose mode")
var versionModeOpt = flag.Bool("version", false, "print version")
var syntaxOpt = flag.Bool("syntax", false, "check syntax only")
var debugAstOpt = flag.Bool("debug-ast", false, "parse a file and print ast")

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

	if *debugAstOpt {
		file := flag.Arg(0)
		node := parser.Parse(file)
		fmt.Printf("%s\n", trompe.NodeDesc(node))
		os.Exit(0)
	}

	if *syntaxOpt {
		file := flag.Arg(0)
		parser.Parse(file)
		os.Exit(0)
	}

	if flag.Arg(0) == "test" {
		trompe.TestCompiledCodeHelloWorld()
		trompe.TestCompiledCodeFizzBuzzCompare()
		trompe.TestCompiledCodeFizzBuzzMatch()
		os.Exit(0)
	}

	file := flag.Arg(0)
	node := parser.Parse(file)
	code := trompe.Compile(file, node)
	fmt.Println(code.Inspect())
	trompe.Run(file, code)
}
