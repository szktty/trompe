package trompe

import "fmt"

type Context struct {
	Parent  *Context
	Module  *Module
	Clos    Closure
	Args    []Value
	NumArgs int
}

func NewContext(parent *Context,
	module *Module,
	clos Closure,
	args []Value,
	numArgs int) Context {
	return Context{
		Parent:  parent,
		Module:  module,
		Clos:    clos,
		Args:    args,
		NumArgs: numArgs,
	}
}

type Stack struct {
	Locals []Value
	Index  int // -1 start
}

func NewStack(len int) Stack {
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
		if value == nil {
			fmt.Printf("%d: nil\n", i)
		} else {
			fmt.Printf("%d: %s\n", i, value.Desc())
		}
	}
	fmt.Printf("\n")
}

type Program struct {
	Path string
	Code *CompiledCode
}

type ProgCounter struct {
	Count  int
	Code   *CompiledCode
	Labels map[int]int
}

func NewProgCounter(code *CompiledCode) ProgCounter {
	return ProgCounter{Count: 0, Code: code, Labels: make(map[int]int, 16)}
}

func (pc *ProgCounter) HasNext() bool {
	return pc.Count < len(pc.Code.Ops)
}

func (pc *ProgCounter) Next() int {
	pc.Count += 1
	return pc.Code.Ops[pc.Count-1]
}

func (pc *ProgCounter) AddLabel(n int) {
	pc.Labels[n] = pc.Count
}

func (pc *ProgCounter) Jump(n int) {
	pc.Count = pc.Labels[n]
}

type Interp struct {
}

func NewInterp() *Interp {
	return &Interp{}
}

func (ip *Interp) Eval(ctx *Context, env *Env, code *CompiledCode) (Value, error) {
	var op int
	var i int
	var top Value
	var retVal Value
	var err error
	pc := NewProgCounter(code)
	cont := true
	stack := NewStack(16)
	args := make([]Value, 16)
	for cont && pc.HasNext() {
		op = pc.Next()
		fmt.Printf("next op: %s\n", GetOpName(op))
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
			stack.Push(code.Lits[i])
		case OpLoadLocal:
			i = pc.Next()
			name := code.Syms[i]
			value := env.Get(name)
			if value == nil {
				err = NewKeyError(ctx, name)
				break
			}
			stack.Push(value)
		case OpLoadAttr:
			i = pc.Next()
			name := code.Syms[i]
			ref := stack.TopPop().(*ValModRef)
			m := ref.Module()
			attr := m.Env.Get(name)
			if attr == nil {
				err = NewKeyError(ctx, name)
				break
			}
			stack.Push(attr)
		case OpLoadModule:
			stack.Push(NewValModRefWithModule(ctx.Module))
		case OpStoreLocal:
			i = pc.Next()
			stack.Set(i, stack.Top())
		case OpStoreAttr:
			i = pc.Next()
			name := code.Syms[i]
			value := stack.TopPop()
			ref := stack.TopPop().(*ValModRef)
			m := ref.Module()
			m.Env.Set(name, value)
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
		case OpBegin:
			env = NewEnv(env)
		case OpEnd:
			env = env.Parent
		case OpCall:
			i = pc.Next()
			for j := 0; j < i; j++ {
				args[j] = stack.TopPop()
			}
			clos := stack.TopPop().Closure()
			if err := ValidateArity(ctx, i, clos.Arity()); err != nil {
				return nil, err
			}
			newCtx := NewContext(ctx, ctx.Module, clos, args, i)
			retVal, err = clos.Apply(ip, &newCtx, NewEnv(env))
			stack.Push(retVal)
		case OpSome:
			top = stack.TopPop()
			stack.Push(NewValOpt(top))
		case OpList:
			i = pc.Next()
			list := ListNil
			for j := 0; j < i; j++ {
				list = list.Cons(stack.TopPop())
			}
			stack.Push(NewValList(list))
		case OpClosedRange:
			right := stack.TopPop()
			left := stack.TopPop()
			stack.Push(NewRange(right.Int(), left.Int(), true))
		case OpHalfOpenRange:
			right := stack.TopPop()
			left := stack.TopPop()
			stack.Push(NewRange(right.Int(), left.Int(), false))
		default:
			panic("unsupported opcode")
		}
	}

	if err != nil {
		return nil, err
	} else if stack.Index < 0 {
		return LangUnit, nil
	} else {
		return stack.Top(), nil
	}
}
