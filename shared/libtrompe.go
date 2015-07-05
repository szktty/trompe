package trompe

func ModuleTrompe() *Module {
	m := NewModule("Trompe")
	m.SetType("unit", TUnit)
	m.SetType("bool", TBool)
	m.SetType("char", TChar)
	m.SetType("string", TString)
	m.SetType("int", TInt)
	m.SetType("float", TFloat)
	m.SetType("list", TList)
	m.SetType("exn", TExn)
	m.SetType("format", TFormat)
	m.SetType("formatter", TFormatter)
	//m.SetType("regexp", TRegexp)
	m.SetExn("Failure", TString)
	m.AddInclude(ModulePervasives())
	m.AddModule(ModuleList())
	m.AddModule(ModuleString())
	return m
}
