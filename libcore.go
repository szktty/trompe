package trompe

import "fmt"

func LibCorePrimShow(prog *Program, args []Value, nargs int) Value {
	// TODO: check number of args
	fmt.Printf("%s\n", args[0].Desc())
	return LangUnit
}
