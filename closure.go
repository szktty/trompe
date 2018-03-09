package trompe

type Closure interface {
	Apply(*Program, *Context) Value
}

type CompiledCode struct {
	Lits []Value
	Ops  []Opcode
}

func (code *CompiledCode) Apply(prog *Program, ctx *Context) Value {
	return Eval(prog, ctx)
}

type Primitive struct {
	Func func(*Program, []Value, int) Value
}

func (prim *Primitive) Apply(prog *Program, ctx *Context) Value {
	return prim.Func(prog, ctx.Args, ctx.Len)
}
