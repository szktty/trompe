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

// TODO: private
type ProgCounter struct {
	Count  int
	Code   *CompiledCode
	Labels map[int]int
	isSkip bool
	dest   int
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
	if pc.isSkip && pc.dest == n {
		pc.isSkip = false
		fmt.Printf("stop skip \n")
	}
}

func (pc *ProgCounter) Jump(n int) {
	if _, ok := pc.Labels[n]; !ok {
		pc.isSkip = true
		pc.dest = n
		return
	}
	fmt.Printf("jump %p\n", pc.Labels[n])
	pc.Count = pc.Labels[n]
}

type Interp struct {
	Top *Module
}

func NewInterp(top *Module) *Interp {
	return &Interp{Top: top}
}

func Run(file string, code *CompiledCode) (Value, error) {
	m := NewModule(nil, file)
	ctx := NewContext(nil, m, code, nil, 0)
	ip := NewInterp(m)
	return ip.Eval(&ctx, NewEnv(m.Env), code)
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

		if pc.isSkip {
			fmt.Printf("--- skip\n")
		}

		switch op {
		case OpNop:
			break
		case OpLoadUnit:
			if !pc.isSkip {
				stack.Push(SharedUnit)
			}
		case OpLoadTrue:
			if !pc.isSkip {
				stack.Push(SharedTrue)
			}
		case OpLoadFalse:
			if !pc.isSkip {
				stack.Push(SharedFalse)
			}
		case OpLoadInt:
			i = pc.Next()
			if !pc.isSkip {
				stack.Push(NewInt(i))
			}
		case OpLoadLit:
			i = pc.Next()
			if !pc.isSkip {
				stack.Push(code.Lits[i])
			}
		case OpLoadLocal:
			i = pc.Next()
			if !pc.isSkip {
				name := code.Syms[i]
				fmt.Printf("-- load local %s\n", name)
				value := env.Get(name)
				if value == nil {
					err = NewKeyError(ctx, name)
					panic(err.Error())
					break
				}
				stack.Push(value)
			}
		case OpLoadAttr:
			i = pc.Next()
			if !pc.isSkip {
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
			}
		case OpLoadModule:
			if !pc.isSkip {
				stack.Push(NewRef(ctx.Module.Path(), ctx.Module))
			}
		case OpStoreLocal:
			i = pc.Next()
			if !pc.isSkip {
				stack.Set(i, stack.Top())
			}
		case OpStoreAttr:
			i = pc.Next()
			if !pc.isSkip {
				name := code.Syms[i]
				v := stack.TopPop()
				ref, _ := ValueToRef(top)
				m := ref.Module()
				m.Env.Set(name, v)
			}
		case OpPop:
			if !pc.isSkip {
				stack.Pop()
			}
		case OpDup:
			if !pc.isSkip {
				stack.Push(stack.Top())
			}
		case OpReturn:
			if !pc.isSkip {
				cont = false
			}
		case OpReturnUnit:
			if !pc.isSkip {
				stack.Push(SharedUnit)
				cont = false
			}
		case OpLabel:
			i = pc.Next()
			pc.AddLabel(i)
		case OpJump:
			i = pc.Next()
			if !pc.isSkip {
				pc.Jump(i)
			}
		case OpBranchTrue:
			i = pc.Next()
			if !pc.isSkip {
				top = stack.TopPop()
				b, _ := ValueToBool(top)
				if b.Value {
					pc.Jump(i)
				}
			}
		case OpBranchFalse:
			i = pc.Next()
			if !pc.isSkip {
				top = stack.TopPop()
				b, _ := ValueToBool(top)
				if !b.Value {
					pc.Jump(i)
				}
			}
		case OpBranchNext:
			i = pc.Next()
			if !pc.isSkip {
				top = stack.Top()
				iter, ok := ValueToIter(top)
				if !ok {
					panic(fmt.Sprintf("not iter %s", top.Desc()))
				}
				if next := iter.Next(); next != nil {
					stack.Push(next)
				} else {
					stack.Pop() // pop iterator
					fmt.Printf("branch next -> %d\n", i)
					pc.Jump(i)
				}
			}
		case OpMatch:
			if !pc.isSkip {
				ptn := stack.TopPop()
				top = stack.TopPop()
				if ptn, ok := ptn.(*Pattern); ok {
					if ptn.Eval(env, top) {
						stack.Push(SharedTrue)
					} else {
						stack.Push(SharedFalse)
					}
				} else {
					panic("not pattern")
				}
			}
		case OpIter:
			top = stack.TopPop()
			if !pc.isSkip {
				iter := NewIter(top)
				if iter == nil {
					panic("cannot get iterator")
				}
				stack.Push(iter)
			}
		case OpBegin:
			if !pc.isSkip {
				env = NewEnv(env)
			}
		case OpEnd:
			if !pc.isSkip {
				env = env.Parent
			}
		case OpCall:
			i = pc.Next()
			if !pc.isSkip {
				for j := i; j > 0; j-- {
					args[j-1] = stack.TopPop()
				}
				clos, ok := ValueToClos(stack.TopPop())
				if !ok {
					panic("not closure")
				}
				if err := ValidateArity(ctx, i, clos.Arity()); err != nil {
					return nil, err
				}
				newCtx := NewContext(ctx, ctx.Module, clos, args, i)
				retVal, err = clos.Apply(ip, &newCtx, NewEnv(env))
				stack.Push(retVal)
			}
		case OpPanic:
			i = pc.Next()
			if !pc.isSkip {
				switch i {
				case OpPanicMatch:
					panic("pattern match error")
				default:
					panic(fmt.Sprintf("unknown panic %d", i))
				}
			}
		case OpSome:
			if !pc.isSkip {
				top = stack.TopPop()
				stack.Push(NewOption(top))
			}
		case OpList:
			i = pc.Next()
			if !pc.isSkip {
				list := ListNil
				for j := 0; j < i; j++ {
					list = list.Cons(stack.TopPop())
				}
				stack.Push(NewList(list))
			}
		case OpClosedRange:
			if !pc.isSkip {
				r := stack.TopPop()
				l := stack.TopPop()
				li, _ := ValueToInt(l)
				ri, _ := ValueToInt(r)
				stack.Push(NewRange(li.Value, ri.Value, true))
			}
		case OpHalfOpenRange:
			if !pc.isSkip {
				r := stack.TopPop()
				l := stack.TopPop()
				li, _ := ValueToInt(l)
				ri, _ := ValueToInt(r)
				stack.Push(NewRange(li.Value, ri.Value, false))
			}
		default:
			panic(fmt.Sprintf("unsupported opcode %s", GetOpName(op)))
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
