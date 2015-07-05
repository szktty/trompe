package main

import (
	"fmt"
	"github.com/jteeuwen/go-pkg-optarg"
	"github.com/szktty/trompe/shared"
	"os"
)

func main() {
	optarg.Header("General options")
	optarg.Add("h", "help", "Print this message and exit", false)
	optarg.Add("s", "syntax", "Check syntax only", false)
	optarg.Add("t", "typing", "Check syntax and types only", false)
	optarg.Add("t", "compile", "Compilation only", false)
	optarg.Add("v", "verbose", "Print verbose message", false)
	optarg.Add("V", "version", "Print version and exit", false)
	optarg.Add("c", "code", "Program <code> passed in string. Omit specified program files", "")
	optarg.Add("i", "interface", "Print inferred interface", false)
	optarg.Add("I", "include", "Add <arg> to the list of search-path directory", "")

	optarg.Header("Warning options")
	optarg.Add("", "warn-error", "Treates warnings as error", false)
	optarg.Add("", "warn-unused",
		"Warn for unused parameters or variables [default]", true)
	optarg.Add("", "warn-nokeyword",
		"Warn if an argument without keyword is passed to a keyword parameter [default]", true)
	optarg.Add("", "warn-camelcase",
		"Warn for camel-case identifier [default]", true)
	optarg.Add("", "nowarn", "Ignore all warnings", false)
	optarg.Add("", "nowarn-unused", "Disable -warn-unused", false)
	optarg.Add("", "nowarn-nokeyword", "Disable -warn-nokeyword", false)
	optarg.Add("", "nowarn-camelcase", "Disable -warn-camelcase", false)

	optarg.Header("Debug options")
	optarg.Add("d", "debug", "Print debug message", false)
	optarg.Add("", "dinstr", "Print bytecode description of compiled functions", false)
	optarg.Add("", "dinstr-all", "Print bytecode description of all compiled codes (include blocks)", false)

	phase := trompe.LoadingPhaseAll
	opts := make([]*optarg.Option, 0)
	for opt := range optarg.Parse() {
		switch opt.Name {
		case "debug":
			trompe.LoggingGroups |= trompe.LogGroupDebug
		case "syntax":
			phase = trompe.LoadingPhaseSyntax
		case "typing":
			phase = trompe.LoadingPhaseTyping
		case "compile":
			phase = trompe.LoadingPhaseCompilation
		case "verbose":
			trompe.LoggingGroups |= trompe.LogGroupVerbose
		case "version":
			fmt.Printf("%s\n", trompe.Version())
			return
		case "help":
			optarg.Usage()
			return
		default:
			opts = append(opts, opt)
		}
	}

	if len(optarg.Remainder) == 0 {
		optarg.Usage()
		return
	}

	f := optarg.Remainder[0]
	if !trompe.FileExists(f) {
		trompe.PrintError("%s: No such file", f)
	}

	if trompe.LogGroupEnabled(trompe.LogGroupVerbose) {
		fmt.Printf("search path:\n")
		for _, path := range trompe.SearchPath() {
			if path == "" {
				fmt.Printf("    \"%s\" (current directory)\n", path)
			} else {
				fmt.Printf("    \"%s\"\n", path)
			}
		}
	}

	trompe.Init()
	s := trompe.NewState()
	if _, err := s.LoadFile(f, phase, opts); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
