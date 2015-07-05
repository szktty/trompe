package trompe

import (
	"fmt"
	"os"
)

type RuntimeError struct {
	Loc *Loc
	Msg string
}

const (
	_LogGroup         = iota
	LogGroupVerbose   = 1 << 1
	LogGroupLexing    = 1 << 2
	LogGroupParsing   = 1 << 3
	LogGroupTyping    = 1 << 4
	LogGroupCompiling = 1 << 5
	LogGroupExec      = 1 << 6
	LogGroupDebug     = LogGroupLexing | LogGroupParsing | LogGroupTyping | LogGroupCompiling | LogGroupExec
)

var LogGroupNames = map[uint64]string{
	LogGroupDebug:     "Debug",
	LogGroupVerbose:   "Verbose",
	LogGroupLexing:    "Lexing",
	LogGroupParsing:   "Parsing",
	LogGroupTyping:    "Typing",
	LogGroupCompiling: "Compiling",
	LogGroupExec:      "Exec",
}

var LoggingGroups uint64

func LogGroupEnabled(flag uint64) bool {
	return LoggingGroups&flag > 0
}

func Logf(flag uint64, f string, v ...interface{}) {
	if LoggingGroups&flag > 0 {
		fmt.Print("# ")
		fmt.Printf(f, v...)
		fmt.Println()
	}
}

func Debugf(f string, v ...interface{}) {
	Logf(LogGroupDebug, f, v...)
}

func Verbosef(f string, v ...interface{}) {
	Logf(LogGroupVerbose, f, v...)
}

func Panicf(f string, v ...interface{}) {
	panic(fmt.Errorf(f, v...))
}

func LogTypingf(f string, v ...interface{}) {
	Logf(LogGroupTyping, f, v...)
}

func LogCompilingf(f string, v ...interface{}) {
	Logf(LogGroupCompiling, f, v...)
}

func PrintAndExit(code int, f string, v ...interface{}) {
	fmt.Printf(f, v...)
	fmt.Println()
	os.Exit(code)
}

func PrintError(f string, v ...interface{}) {
	fmt.Printf("Error: ")
	PrintAndExit(1, f, v...)
}

func PrintWarn(file string, line int, f string, v ...interface{}) {
	fmt.Printf("%s: line %s: Warn: ", file, line)
	fmt.Printf(f, v...)
}

/*
func PrintStackTrace(err error) {
	switch e := err.(type) {
	case *runtimeError:
		fmt.Printf("Traceback:\n")
		ctx := e.ctx
		for ctx != nil {
			fmt.Printf("    File %s, Line %d, %s\n",
				ctx.code.file, ctx.currentLine()+1, ctx.code.name)
			ctx = ctx.parent
		}
		fmt.Printf("Error: %s: line %d: %s\n", e.ctx.code.file,
			e.ctx.currentLine()+1, e.msg)
	default:
		panic("unknown error")
	}
}
*/

func NewRuntimeError(loc *Loc, err error) error {
	return &RuntimeError{Loc: loc, Msg: err.Error()}
}

func RuntimeErrorf(loc *Loc, f string, v ...interface{}) error {
	msg := fmt.Sprintf(f, v...)
	return &RuntimeError{Loc: loc, Msg: msg}
}

func (err *RuntimeError) Error() string {
	return fmt.Sprintf("File \"%s\", line %d, column %d:\nError: %s\n",
		err.Loc.File, err.Loc.Start.Line+1, err.Loc.Start.Col+1, err.Msg)
}
