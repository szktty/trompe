package trompe

import "fmt"

type Context struct {
	Parent  *Context
	Clos    Closure
	Args    []Value
	NumArgs int
	Env     *Env
}

// TODO: Env
func CreateContext(parent *Context, clos Closure, args []Value, numArgs int) Context {
	return Context{
		Parent:  parent,
		Clos:    clos,
		Args:    args,
		NumArgs: numArgs,
		Env:     nil,
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
	Labels map[int]int
}

func CreateProgCounter(ctx *Context) ProgCounter {
	return ProgCounter{Count: 0, Ctx: ctx, Labels: make(map[int]int, 16)}
}

func (pc *ProgCounter) HasNext() bool {
	return pc.Count < len(pc.Ctx.CompiledCode().Ops)
}

func (pc *ProgCounter) Next() int {
	pc.Count += 1
	return pc.Ctx.CompiledCode().Ops[pc.Count-1]
}

func (pc *ProgCounter) AddLabel(n int) {
	pc.Labels[n] = pc.Count
}

func (pc *ProgCounter) Jump(n int) {
	pc.Count = pc.Labels[n]
}

func (prog *Program) Eval(ctx *Context) (Value, error) {
	var op int
	var i int
	var top Value
	var retVal Value
	var err error
	pc := CreateProgCounter(ctx)
	cont := true
	stack := CreateStack(16)
	args := make([]Value, 16)
	for cont && pc.HasNext() {
		op = pc.Next()
		fmt.Printf("op: %s\n", GetOpName(op))
		stack.Inspect()
		switch op {
		case OpNop:
			break
		case OpLoadUnit:
			stack.Push(LangUnit)
		case OpLoadTrue:
			stack.Push(LangTrue)
		case OpLoadFalse:
			stack.Push(LangFalse)
		case OpLoadInt:
			i = pc.Next()
			stack.Push(&ValInt{i})
		case OpLoadLit:
			i = pc.Next()
			stack.Push(ctx.Literal(i))
		case OpLoadLocal:
			i = pc.Next()
			stack.Push(stack.Get(i))
		case OpLoadAttr:
			i = pc.Next()
			name := ctx.Literal(i).String()
			attr := ctx.Env.GetAttr(name)
			if attr == nil {
				err = CreateKeyError(ctx, name)
				break
			}
			stack.Push(attr)
		case OpLoadPrim:
			i = pc.Next()
			prim := GetPrim(ctx.Literal(i).String())
			stack.Push(prim)
		case OpStore:
			i = pc.Next()
			stack.Set(i, stack.Top())
		case OpPop:
			stack.Pop()
		case OpDup:
			stack.Push(stack.Top())
		case OpReturn:
			break
		case OpReturnUnit:
			stack.Push(LangUnit)
			break
		case OpLabel:
			i = pc.Next()
			pc.AddLabel(i)
		case OpJump:
			i = pc.Next()
			pc.Jump(i)
		case OpBranchTrue:
			i = pc.Next()
			top = stack.Top()
			if top.Bool() {
				pc.Jump(i)
			}
		case OpBranchFalse:
			i = pc.Next()
			top = stack.Top()
			if !top.Bool() {
				pc.Jump(i)
			}
		case OpCall:
			i = pc.Next()
			for j := 0; j < i; j++ {
				args[j] = stack.TopPop()
			}
			clos := stack.TopPop().Closure()
			newCtx := CreateContext(ctx, clos, args, i)
			retVal, err = clos.Apply(prog, &newCtx)
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
			panic("unsupported opcode")
		}
	}

	if err != nil {
		return nil, err
	} else {
		return stack.Top(), nil
	}
}
