package trompe

import (
	"fmt"
)

func ModuleList() *Module {
	m := NewModule("List")

	// val length : 'a list -> int
	m.SetFieldType("length",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(
						TApp(TcList, TArgs(TVar("a"))),
						TInt))),
				TArgs(TVar("a")))))
	m.SetPrim("length", List_length)

	// val hd : 'a list -> 'a
	m.SetFieldType("hd",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(
						TApp(TcList, TArgs(TVar("a"))),
						TVar("a")))),
				TArgs(TVar("a")))))
	m.SetPrim("hd", List_hd)

	m.SetFieldType("iter",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(
						TApp(TcList, TArgs(TVar("a"))),
						TApp(TcArrow, TArgs(TVar("a"), TUnit)),
						TUnit))),
				TArgs(TVar("a")))))

	m.SetFieldType("filter",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(
						TApp(TcList, TArgs(TVar("a"))),
						TApp(TcArrow, TArgs(TVar("a"), TBool)),
						TApp(TcList, TArgs(TVar("a")))))),
				TArgs(TVar("a")))))

	return m
}

func List_length(state *State, parent *Context, args []Value) (Value, error) {
	list := args[0].(*List)
	return list.Length(), nil
}

func List_hd(state *State, parent *Context, args []Value) (Value, error) {
	list := args[0].(*List)
	if list == NilValue {
		return nil, fmt.Errorf("list is nil")
	} else {
		return list.Head, nil
	}
}
