package trompe

import (
	"fmt"
)

func ModuleString() *Module {
	m := NewModule("String")
	m.SetFieldType("length", TApp(TcArrow, TArgs(TString, TInt)))
	m.SetPrim("length", String_length)

	m.SetFieldType("of_int", TApp(TcArrow, TArgs(TInt, TString)))
	m.SetPrim("of_int", String_of_int)

	m.SetFieldType("to_int", TApp(TcArrow, TArgs(TString, TInt)))
	m.SetPrim("to_int", String_of_int)

	//m.SetFieldType("get", TApp(TcKeyArrowv("t", "pos"), TArgs(TString, TInt, TChar)))
	m.SetFieldType("get", TApp(TcArrow, TArgs(TString, TInt, TChar)))
	m.SetPrim("get", String_get)

	return m
}

func String_length(state *State, parent *Context, args []Value) (Value, error) {
	return int64(len(args[0].(string))), nil
}

func String_of_int(state *State, parent *Context, args []Value) (Value, error) {
	return fmt.Sprintf("%d", args[0].(int64)), nil
}

// TODO
func String_to_int(state *State, parent *Context, args []Value) (Value, error) {
	return fmt.Sprintf("%d", args[0].(int64)), nil
}

func String_get(state *State, parent *Context, args []Value) (Value, error) {
	s := []rune(args[0].(string))
	i := args[1].(int64)
	return s[i], nil
}
