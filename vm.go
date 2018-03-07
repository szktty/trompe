package trompe

type Value interface{}

type Closure struct {
	Lits []Value
	Ops  []Opcode
}

type Context struct {
	Parent *Context
	Clos   *Closure
}

type Stack struct {
	Locals []*Value
	Index  int // -1 start
}

func (s *Stack) Push(value *Value) {
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

type Program struct {
	Path string
}

func Eval(prog *Program, stack *Stack, ctx *Context) {
	for op := range ctx.Clos.Ops {
		switch op {
		case OpNop:
			break
		case OpPop:
			stack.Pop()
		default:
			break
		}
	}
}
