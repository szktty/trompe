package trompe

import "fmt"

type Context struct {
	Parent *Context
	Clos   Closure
	Args   []Value
	Len    int
	Env    *Env
}

// TODO: Env
func CreateContext(parent *Context, clos Closure, args []Value, len int) Context {
	return Context{
		Parent: parent,
		Clos:   clos,
		Args:   args,
		Len:    len,
		Env:    nil,
	}
}

func (ctx *Context) CompiledCode() *CompiledCode {
	return ctx.Clos.(*CompiledCode)
}

func (ctx *Context) Literal(i int) Value {
	return ctx.CompiledCode().Lits[i]
}

type Stack struct {
	Locals []Value
	Index  int // -1 start
}

func CreateStack(len int) Stack {
	return Stack{Locals: make([]Value, len), Index: -1}
}

func (s *Stack) Top() Value {
	return s.Locals[s.Index]
}

func (s *Stack) TopPop() Value {
	top := s.Locals[s.Index]
	s.Pop()
	return top
}

func (s *Stack) Get(i int) Value {
	return s.Locals[i]
}

func (s *Stack) Set(i int, value Value) {
	s.Locals[i] = value
}

func (s *Stack) Push(value Value) {
	s.Index++
	if len(s.Locals) < s.Index {
		s.Locals = append(s.Locals, value)
	} else {
		s.Locals[s.Index] = value
	}
}

func (s *Stack) Pop() {
	s.Index--
}

func (s *Stack) Inspect() {
	fmt.Printf("stack:\n")
	for i, value := range s.Locals {
		if i > s.Index {
			continue
		}
		if i == s.Index {
			fmt.Printf(" -> ")
		} else {
			fmt.Printf("    ")
		}
		fmt.Printf("%d: %s\n", i, value.Desc())
	}
	fmt.Printf("\n")
}

type Program struct {
	Path string
}

type ProgCounter struct {
	Count  int
	Ctx    *Context
	Labels map[string]int
}

func CreateProgCounter(ctx *Context) ProgCounter {
	return ProgCounter{Count: 0, Ctx: ctx, Labels: make(map[string]int)}
}

func (pc *ProgCounter) HasNext() bool {
	return pc.Count < len(pc.Ctx.CompiledCode().Ops)
}

func (pc *ProgCounter) Next() int {
	pc.Count += 1
	return pc.Ctx.CompiledCode().Ops[pc.Count-1]
}

func (pc *ProgCounter) AddLabel(name string) {
	pc.Labels[name] = pc.Count
}

func (pc *ProgCounter) Jump(label string) {
	pc.Count = pc.Labels[label]
}

func (prog *Program) Eval(ctx *Context) Value {
	var op int
	var i int
	var top Value
	var retVal Value
	pc := CreateProgCounter(ctx)
	cont := true
	stack := CreateStack(16)
	args := make([]Value, 16)
	for cont && pc.HasNext() {
		stack.Inspect()
		op = pc.Next()
		switch op {
		case OpNop:
			break
		case OpLoadUnit:
			stack.Push(SharedValUnit)
		case OpLoadTrue:
			stack.Push(SharedValTrue)
		case OpLoadFalse:
			stack.Push(SharedValFalse)
		case OpLoadInt:
			i = pc.Next()
			stack.Push(&ValInt{i})
		case OpLoadLit:
			i = pc.Next()
			stack.Push(ctx.Literal(i))
		case OpLoadLocal:
			i = pc.Next()
			stack.Push(stack.Get(i))
		case OpStore:
			i = pc.Next()
			stack.Set(i, stack.Top())
		case OpPop:
			stack.Pop()
		case OpReturn:
			return stack.Top()
		case OpLabel:
			i = pc.Next()
			pc.AddLabel(ctx.Literal(i).String())
		case OpJump:
			i = pc.Next()
			pc.Jump(ctx.Literal(i).String())
		case OpBranchTrue:
			i = pc.Next()
			top = stack.Top()
			if top.Bool() {
				pc.Jump(ctx.Literal(i).String())
			}
		case OpBranchFalse:
			i = pc.Next()
			top = stack.Top()
			if !top.Bool() {
				pc.Jump(ctx.Literal(i).String())
			}
		case OpCall:
			i = pc.Next()
			for j := 0; j < i; j++ {
				args[j] = stack.TopPop()
			}
			clos := stack.TopPop().Closure()
			newCtx := CreateContext(ctx, clos, args, i)
			retVal = clos.Apply(prog, &newCtx)
			stack.Push(retVal)
		case OpSome:
			top = stack.TopPop()
			stack.Push(CreateValOpt(top))
		case OpList:
			i = pc.Next()
			list := ListNil
			for j := 0; j < i; j++ {
				list = list.Cons(stack.TopPop())
			}
			stack.Push(CreateValList(list))
		default:
			panic("unknown opcode")
		}
	}

	return stack.Top()
}
