package trompe

func ModuleList() *Module {
	m := NewModule("List")

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

	m.SetFieldType("length",
		TPoly(Tyvars("a"),
			TApp(
				TcTyFun(Tyvars("a"),
					TApp(TcArrow, TArgs(
						TApp(TcList, TArgs(TVar("a"))),
						TInt))),
				TArgs(TVar("a")))))

	return m
}

func List_length(state *State, parent *Context, args []Value) (Value, error) {
	return 0, nil // TODO
}
