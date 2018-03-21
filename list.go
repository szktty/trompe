package trompe

type List struct {
	Value Value
	Next  *List
}

var ListNil = &List{nil, nil}

func CreateList(value Value) *List {
	return &List{Value: value, Next: nil}
}

func (l *List) Len() int {
	i := 0
	for l.Next != nil {
		l = l.Next
		i++
	}
	return i
}

func (l *List) Cons(value Value) *List {
	return &List{Value: value, Next: l}
}
