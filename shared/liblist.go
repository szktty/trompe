package trompe

import (
	"fmt"
)

var TypePolyList1 = TApp(TcList, TArgs(TVar("a")))
var TypePolyList2 = TApp(TcList, TArgs(TVar("a"), TVar("b")))

func TypePolyListArrow(arg ...Type) Type {
	return WrapInTypePoly1(TApp(TcArrow, TArgs(arg...)))
}

func ModuleList() *Module {
	m := NewModule("List")

	// val length : 'a list -> int
	m.SetFieldType("length", TypePolyListArrow(TypePolyList1, TInt))
	m.SetPrim("length", List_length)

	// val hd : 'a list -> 'a
	m.SetFieldType("hd", TypePolyListArrow(TypePolyList1, TVar("a")))
	m.SetPrim("hd", List_hd)

	// val tl : 'a list -> 'a list
	m.SetFieldType("tl", TypePolyListArrow(TypePolyList1, TypePolyList1))
	m.SetPrim("tl", List_tl)

	// val iter : 'a list -> ('a -> unit) -> unit
	m.SetFieldType("iter", TypePolyListArrow(TypePolyList1,
		TApp(TcArrow, TArgs(TVar("a"), TUnit)), TUnit))
	m.SetPrim("iter", List_iter)

	// val filter : 'a list -> ('a -> bool) -> 'a list
	m.SetFieldType("filter", TypePolyListArrow(TypePolyList1,
		TApp(TcArrow, TArgs(TVar("a"), TBool)), TypePolyList1))

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

func List_tl(state *State, parent *Context, args []Value) (Value, error) {
	list := args[0].(*List)
	if list == NilValue {
		return nil, fmt.Errorf("list is nil")
	} else {
		return list.Tail, nil
	}
}

func List_iter(state *State, parent *Context, args []Value) (Value, error) {
	list := args[0].(*List)
	blk := args[1].(*BlockClosure)
	for list != NilValue {
		state.Exec(parent.Module, parent, blk, []Value{list.Head})
		list = list.Tail
	}
	return UnitValue, nil
}
