package trompe

type Primitive struct {
	Name  string
	Func  func(*Context, []Value, int) (Value, error)
	Arity int
}

var sharedPrims = map[string]*Primitive{}

func NewPrim(name string,
	f func(*Context, []Value, int) (Value, error),
	arity int) *Primitive {
	return &Primitive{name, f, arity}
}

func (prim *Primitive) Apply(ctx *Context) (Value, error) {
	return prim.Func(ctx, ctx.Args, ctx.NumArgs)
}

func GetPrim(name string) *Primitive {
	return sharedPrims[name]
}

func SetPrim(
	name string,
	f func(*Context, []Value, int) (Value, error),
	arity int) {
	sharedPrims[name] = NewPrim(name, f, arity)
}
