package trompe

type Primitive struct {
	Name  string
	Func  func(*Context, []Value, int) (Value, error)
	Arity int
}

var sharedPrims = map[string]Value{}

func (prim *Primitive) Apply(ctx *Context) (Value, error) {
	return prim.Func(ctx, ctx.Args, ctx.NumArgs)
}

func GetPrim(name string) Value {
	return sharedPrims[name]
}

func SetPrim(
	name string,
	f func(*Context, []Value, int) (Value, error),
	arity int) {
	sharedPrims[name] = CreateValClos(&Primitive{Name: name, Func: f, Arity: arity})
}
