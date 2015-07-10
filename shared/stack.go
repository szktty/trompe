package trompe

type Stack struct {
	Slots []Value
	Capa  int
	Ptr   int
}

var ExtraNumSlots = 4

func NewStack(capa int) *Stack {
	capa += ExtraNumSlots
	return &Stack{Slots: make([]Value, capa), Capa: capa, Ptr: -1}
}

func (s *Stack) Increase(size int) {
	old := s.Slots
	s.Capa += size + ExtraNumSlots
	s.Slots = make([]Value, s.Capa)
	copy(s.Slots, old)
}

func (s *Stack) Top() Value {
	return s.Slots[s.Ptr]
}

func (s *Stack) TopArray() []Value {
	if v, ok := s.Top().([]Value); ok {
		return v
	} else {
		panic("not array")
	}
}

func (s *Stack) Push(obj Value) {
	if v, ok := obj.(int); ok {
		obj = int64(v)
	}
	s.Ptr++
	s.Slots[s.Ptr] = obj
}

func (s *Stack) PushLocal(i int) {
	s.Push(s.Slots[i])
}

func (s *Stack) PushLocalIndirect(b int) {
	ary := s.Slots[b/16].([]Value)
	s.Push(ary[b%16])
}

func (s *Stack) PushUnit() {
	s.Push(UnitValue)
}

func (s *Stack) PushNil() {
	s.Push(NilValue)
}

func (s *Stack) PushBool(v bool) {
	s.Push(v)
}

func (s *Stack) Pop() Value {
	s.Ptr--
	return s.Slots[s.Ptr+1]
}

func (s *Stack) SwapPop() {
	s.Slots[s.Ptr-1] = s.Slots[s.Ptr]
	s.Ptr--
}

func (s *Stack) PopBool() bool {
	return s.Pop().(bool)
}

func (s *Stack) PopInt() int64 {
	if v, ok := s.Pop().(int64); ok {
		return v
	} else {
		panic("not int")
	}
}

func (s *Stack) PopArray() []Value {
	if v, ok := s.Pop().([]Value); ok {
		return v
	} else {
		panic("not array")
	}
}

func (s *Stack) PopValues(n int) []Value {
	vs := make([]Value, n)
	for i := 0; i < n; i++ {
		vs[n-1-i] = s.Pop()
	}
	return vs
}

func (s *Stack) StorePopLocal(i int) {
	s.Slots[i] = s.Pop()
}

func (s *Stack) StorePopLocalIndirect(b int) {
	ary := s.Slots[b/16].([]Value)
	ary[b%16] = s.Pop()
}
