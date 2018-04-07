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
	return ValClosType
}

func (prim *Primitive) Desc() string {
	return fmt.Sprintf("<prim %p>", prim)
}

func (prim *Primitive) Bool() bool {
	panic("Prim")
}

func (prim *Primitive) Int() int {
	panic("Prim")
}

func (prim *Primitive) String() string {
	panic("Prim")
}

func (prim *Primitive) Closure() Closure {
	return prim
}

func (prim *Primitive) List() *List {
	panic("Prim")
}

func (prim *Primitive) Tuple() []Value {
	panic("Prim")
}

func (prim *Primitive) Iter() ValIter {
	return nil
}

func (prim *Primitive) Arity() int {
	return prim.arity
}

func (prim *Primitive) Apply(interp *Interp, ctx *Context, env *Env) (Value, error) {
	return prim.Func(ctx, ctx.Args, ctx.NumArgs)
}
