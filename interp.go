package trompe

type Context struct {
	Parent *Context
	Clos   *Closure
}

type Stack struct {
	Locals []Value
	Index  int // -1 start
}

func (s *Stack) Top() Value {
	return s.Locals[s.Index]
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

type Program struct {
	Path string
}

type ProgCounter struct {
	Count  int
	Ctx    *Context
	Labels map[string]int
}

func CreateProgCounter(ctx *Context) ProgCounter {
	return ProgCounter{Count: -1, Ctx: ctx, Labels: make(map[string]int)}
}

func (pc *ProgCounter) HasNext() bool {
	return pc.Count < len(pc.Ctx.Clos.Ops)
}

func (pc *ProgCounter) Next() int {
	pc.Count += 1
	return pc.Ctx.Clos.Ops[pc.Count-1]
}

func (pc *ProgCounter) AddLabel(name string) {
	pc.Labels[name] = pc.Count
}

func Eval(prog *Program, stack *Stack, ctx *Context) {
	pc := CreateProgCounter(ctx)
	op := OpNop
	i := 0
	for pc.HasNext() {
		op = pc.Next()
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
			stack.Push(&ValInt{i})
		case OpLoadLit:
			i = pc.Next()
			stack.Push(ctx.Clos.Lits[i])
		case OpLoadLocal:
			i = pc.Next()
			stack.Push(stack.Get(i))
		case OpStore:
			i = pc.Next()
			stack.Set(i, stack.Top())
		case OpPop:
			stack.Pop()
		case OpLabel:
			i = pc.Next()
			pc.AddLabel(ctx.Clos.Lits[i].String())
		default:
			break
		}
	}
}
