package trompe

import (
	"fmt"
)

type Primitive struct {
	Func  PrimFun
	arity int
}

type PrimFun = func(*Context, []Value, int) (Value, error)

func NewPrim(f func(*Context, []Value, int) (Value, error),
	arity int) *Primitive {
	return &Primitive{f, arity}
}

func (prim *Primitive) Type() int {
	return ValueTypeClos
}

func (prim *Primitive) Desc() string {
	return fmt.Sprintf("<prim %p>", prim)
}

func (prim *Primitive) Arity() int {
	return prim.arity
}

func (prim *Primitive) Apply(interp *Interp, ctx *Context, env *Env) (Value, error) {
	return prim.Func(ctx, ctx.Args, ctx.NumArgs)
}
