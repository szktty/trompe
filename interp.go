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

func (s *Stack) Sharedet(i int) Value {
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
	fmt.Printf("# label[%d] = %d\n", n, pc.Count)
}

func (pc *ProgCounter) Jump(n int) {
	if _, ok := pc.Labels[n]; !ok {
		panic("not exists")
	}
	fmt.Printf("jump %p\n", pc.Labels[n])
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
		fmt.Printf("%d: next op: %s\n", pc.Count, GetOpName(op))
		stack.Inspect()
		switch op {
		case OpNop:
			break
		case OpLoadUnit:
			stack.Push(SharedUnit)
		case OpLoadTrue:
			stack.Push(SharedTrue)
		case OpLoadFalse:
			stack.Push(SharedFalse)
		case OpLoadInt:
			i = pc.Next()
			stack.Push(NewInt(i))
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
			top := stack.TopPop()
			ref, _ := ValueToRef(top)
			m := ref.Module()
			attr := m.Env.Get(name)
			if attr == nil {
				err = NewKeyError(ctx, name)
				break
			}
			stack.Push(attr)
		case OpLoadModule:
			stack.Push(NewRef(ctx.Module.Path(), ctx.Module))
		case OpStoreLocal:
			i = pc.Next()
			stack.Set(i, stack.Top())
		case OpStoreAttr:
			i = pc.Next()
			name := code.Syms[i]
			v := stack.TopPop()
			ref, _ := ValueToRef(top)
			m := ref.Module()
			m.Env.Set(name, v)
		case OpPop:
			stack.Pop()
		case OpDup:
			stack.Push(stack.Top())
		case OpReturn:
			break
		case OpReturnUnit:
			stack.Push(SharedUnit)
			break
		case OpLabel:
			i = pc.Next()
			pc.AddLabel(i)
		case OpJump:
			i = pc.Next()
			pc.Jump(i)
		case OpBranchTrue:
			i = pc.Next()
			top = stack.TopPop()
			b, _ := ValueToBool(top)
			if b.Value {
				pc.Jump(i)
			}
		case OpBranchFalse:
			i = pc.Next()
			top = stack.TopPop()
			b, _ := ValueToBool(top)
			if !b.Value {
				pc.Jump(i)
			}
		case OpBranchNext:
			i = pc.Next()
			top = stack.Top()
			iter, ok := ValueToIter(top)
			if !ok {
				panic("not iter")
			}
			if next := iter.Next(); next != nil {
				stack.Push(next)
			} else {
				stack.Pop() // pop iterator
				fmt.Printf("branch next -> %d\n", i)
				pc.Jump(i)
			}
		case OpIter:
			top = stack.TopPop()
			iter := NewIter(top)
			if iter == nil {
				panic("cannot get iterator")
			}
			stack.Push(iter)
		case OpBegin:
			env = NewEnv(env)
		case OpEnd:
			env = env.Parent
		case OpCall:
			i = pc.Next()
			for j := 0; j < i; j++ {
				args[j] = stack.TopPop()
			}
			clos, _ := ValueToClos(stack.TopPop())
			if err := ValidateArity(ctx, i, clos.Arity()); err != nil {
				return nil, err
			}
			newCtx := NewContext(ctx, ctx.Module, clos, args, i)
			retVal, err = clos.Apply(ip, &newCtx, NewEnv(env))
			stack.Push(retVal)
		case OpSome:
			top = stack.TopPop()
			stack.Push(NewOption(top))
		case OpList:
			i = pc.Next()
			list := ListNil
			for j := 0; j < i; j++ {
				list = list.Cons(stack.TopPop())
			}
			stack.Push(NewList(list))
		case OpClosedRange:
			r := stack.TopPop()
			l := stack.TopPop()
			li, _ := ValueToInt(l)
			ri, _ := ValueToInt(r)
			stack.Push(NewRange(li.Value, ri.Value, true))
		case OpHalfOpenRange:
			r := stack.TopPop()
			l := stack.TopPop()
			li, _ := ValueToInt(l)
			ri, _ := ValueToInt(r)
			stack.Push(NewRange(li.Value, ri.Value, false))
		default:
			panic("unsupported opcode")
		}
	}

	if err != nil {
		return nil, err
	} else if stack.Index < 0 {
		return SharedUnit, nil
	} else {
		return stack.Top(), nil
	}
}
