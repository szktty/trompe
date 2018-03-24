package main

import (
	"flag"
	"fmt"
	"github.com/szktty/trompe"
	"github.com/szktty/trompe/parser"
	"io/ioutil"
	"os"
	//"path/filepath"
	"strings"
)

var debugModeOpt = flag.Bool("d", false, "debug mode")
var verboseModeOpt = flag.Bool("v", false, "verbose mode")
var versionModeOpt = flag.Bool("version", false, "print version")
var printOpt = flag.Bool("p", false, "output compiled code to standart output")
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
		fmt.Printf("Usage: trompec [options] files\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	trompe.Init()

	if *debugAstOpt {
		file := flag.Arg(0)
		node := parser.Parse(file)
		fmt.Printf("%s\n", trompe.NodeDesc(node))
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
	objFile := trompe.NewObjectFile(file)
	objFile.AddCompiledCode(code)
	objFile.AddAttr(trompe.NewObjectAttr("main",
		trompe.NewObjectValue(trompe.ObjectValueTypeCode,
			fmt.Sprintf("%d", code.Id))))
	data, err := objFile.Marshal()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	if *printOpt {
		var buf strings.Builder
		buf.Write(data)
		fmt.Println(buf.String())
	} else {
		path := file + "o" // .tmo
		ioutil.WriteFile(path, data, 0)
	}
}
