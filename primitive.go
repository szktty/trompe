package trompe

type Primitive struct {
	Name  string
	Func  func(*Program, []Value, int) Value
	Arity int
}

var sharedPrims = map[string]Value{}

func (prim *Primitive) Apply(prog *Program, ctx *Context) Value {
	return prim.Func(prog, ctx.Args, ctx.NumArgs)
}

func GetPrim(name string) Value {
	return sharedPrims[name]
}

func SetPrim(name string, f func(*Program, []Value, int) Value, arity int) {
	sharedPrims[name] = CreateValClos(&Primitive{Name: name, Func: f, Arity: arity})
}
