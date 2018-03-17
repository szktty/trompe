package main

import (
	"flag"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/szktty/trompe"
	"github.com/szktty/trompe/parser"
	"os"
)

var debugModeOpt = flag.Bool("d", false, "debug mode")
var verboseModeOpt = flag.Bool("v", false, "verbose mode")
var versionModeOpt = flag.Bool("version", false, "print version")

type TreeShapeListener struct {
	*parser.BaseTrompeListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (this *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	//fmt.Println(ctx.GetText())
}

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
		os.Exit(0)
	}

	file := flag.Arg(0)
	input, _ := antlr.NewFileStream(file)
	lexer := parser.NewTrompeLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewTrompeParser(stream)
	p.BuildParseTrees = true
	tree := p.Chunk()
	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}
